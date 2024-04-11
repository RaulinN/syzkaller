//go:build profiling

// Copyright 2024 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package fuzzer

import (
	"fmt"
	"github.com/google/syzkaller/pkg/fuzzer/ablation_flags"
	"math/rand"
	"time"

	"github.com/google/syzkaller/pkg/corpus"
	"github.com/google/syzkaller/pkg/cover"
	"github.com/google/syzkaller/pkg/ipc"
	"github.com/google/syzkaller/pkg/signal"
	"github.com/google/syzkaller/prog"
)

const (
	smashPrio int64 = iota + 1
	genPrio
	triagePrio
	candidatePrio
	candidateTriagePrio
)

type job interface {
	run(fuzzer *Fuzzer)
	saveID(id int64)
	priority() priority
}

type ProgTypes int

const (
	progCandidate ProgTypes = 1 << iota
	progMinimized
	progSmashed
	progInTriage
)

type jobPriority struct {
	prio priority
}

func newJobPriority(base int64) jobPriority {
	prio := append(make(priority, 0, 2), base)
	return jobPriority{prio}
}

func (jp jobPriority) priority() priority {
	return jp.prio
}

// If we prioritize execution requests only by the base priorities of their origin
// jobs, we risk letting 1000s of simultaneous jobs slowly progress in parallel.
// It's better to let same-prio jobs that were started earlier finish first.
// saveID() allows Fuzzer to attach this sub-prio at the moment of job creation.
func (jp *jobPriority) saveID(id int64) {
	jp.prio = append(jp.prio, id)
}

func genEmptyProgRequest(fuzzer *Fuzzer, rnd *rand.Rand) *Request {
	p := fuzzer.target.Generate(rnd,
		0,
		fuzzer.ChoiceTable())
	return &Request{
		Prog:          p,
		NeedSignal:    true,
		stat:          statGenerate,
		requesterStat: statGenerate,
	}
}

func genProgRequest(fuzzer *Fuzzer, rnd *rand.Rand) *Request {
	fuzzer.profilingStats.IncModeCounter(ProfilingStatModeGenerate)
	start := time.Now()

	p := fuzzer.target.Generate(rnd,
		prog.RecommendedCalls,
		fuzzer.ChoiceTable())

	delta := time.Since(start)
	fuzzer.profilingStats.AddModeDuration(ProfilingStatModeGenerate, delta)

	return &Request{
		Prog:          p,
		NeedSignal:    true,
		stat:          statGenerate,
		requesterStat: statGenerate,
	}
}

func profileMutateObserver(fuzzer *Fuzzer, observer map[prog.MutatorIndex]int) {
	for k, v := range observer {
		switch k {
		case prog.MutatorIndexSquashAny:
			fuzzer.profilingStats.AddMutatorCounter(ProfilingStatMutatorSquashAny, v)
		case prog.MutatorIndexSplice:
			fuzzer.profilingStats.AddMutatorCounter(ProfilingStatMutatorSplice, v)
		case prog.MutatorIndexInsertCall:
			fuzzer.profilingStats.AddMutatorCounter(ProfilingStatMutatorInsertCall, v)
		case prog.MutatorIndexMutateArg:
			fuzzer.profilingStats.AddMutatorCounter(ProfilingStatMutatorMutateArg, v)
		case prog.MutatorIndexRemoveCall:
			fuzzer.profilingStats.AddMutatorCounter(ProfilingStatMutatorRemoveCall, v)
		default:
			panic(fmt.Sprintf("mutator index '%v' case switch unknown in profileMutateObserver", k))
		}
	}
}

func mutateProgRequest(fuzzer *Fuzzer, rnd *rand.Rand) *Request {
	p := fuzzer.Config.Corpus.ChooseProgram(rnd)
	if p == nil {
		return nil
	}
	newP := p.Clone()

	// if the mutate mode is disabled (via ablation), skip the mutation and return
	// a copy of the original program
	if ablation_flags.ABLATION_MODE_MUTATE_ENABLED {
		fuzzer.profilingStats.IncModeCounter(ProfilingStatModeMutate)
		start := time.Now()

		obs := newP.MutateWithObserver(rnd,
			prog.RecommendedCalls,
			fuzzer.ChoiceTable(),
			fuzzer.Config.NoMutateCalls,
			fuzzer.Config.Corpus.Programs(),
		)

		delta := time.Since(start)
		fuzzer.profilingStats.AddModeDuration(ProfilingStatModeMutate, delta)
		profileMutateObserver(fuzzer, obs)
	}

	return &Request{
		Prog:          newP,
		NeedSignal:    true,
		stat:          statFuzz,
		requesterStat: statFuzz,
	}
}

func candidateRequest(input Candidate) *Request {
	flags := progCandidate
	if input.Minimized {
		flags |= progMinimized
	}
	if input.Smashed {
		flags |= progSmashed
	}
	return &Request{
		Prog:          input.Prog,
		NeedSignal:    true,
		stat:          statCandidate,
		flags:         flags,
		requesterStat: statCandidate,
	}
}

// triageJob are programs for which we noticed potential new coverage during
// first execution. But we are not sure yet if the coverage is real or not.
// During triage we understand if these programs in fact give new coverage,
// and if yes, minimize them and add to corpus.
type triageJob struct {
	p         *prog.Prog
	call      int
	info      ipc.CallInfo
	newSignal signal.Signal
	flags     ProgTypes
	jobPriority
	// in case the coverage increase is indeed real, we need to be able to
	// attribute this contribution to the correct execution mode (coming
	// from the request that started the triageJob), hence we store it
	stat          string
	requesterStat string
}

func triageJobPrio(flags ProgTypes) jobPriority {
	if flags&progCandidate > 0 {
		return newJobPriority(candidateTriagePrio)
	}
	return newJobPriority(triagePrio)
}

func (job *triageJob) run(fuzzer *Fuzzer) {
	if job.requesterStat == "" {
		fuzzer.Logf(0, "ERROR! started a triage job that does not have any requestStat!") // FIXME NICOLAS REMOVE
	}
	logCallName := "extra"
	if job.call != -1 {
		callName := job.p.Calls[job.call].Meta.Name
		logCallName = fmt.Sprintf("call #%v %v", job.call, callName)
	}
	fuzzer.Logf(3, "triaging input for %v (new signal=%v)", logCallName, job.newSignal.Len())
	// Compute input coverage and non-flaky signal for minimization.
	info, stop := job.deflake(fuzzer)
	if stop || info.newStableSignal.Empty() {
		return
	}

	if job.flags&progMinimized == 0 {
		stop = job.minimize(fuzzer, info.newStableSignal)
		if stop {
			return
		}
	}
	fuzzer.Logf(2, "added new input for %q to the corpus:\n%s", logCallName, job.p.String())
	if job.flags&progSmashed == 0 {
		fuzzer.startJob(&smashJob{
			p:           job.p.Clone(),
			call:        job.call,
			jobPriority: newJobPriority(smashPrio),
		})
	}
	input := corpus.NewInput{
		Prog:     job.p,
		Call:     job.call,
		Signal:   info.stableSignal,
		Cover:    info.cover.Serialize(),
		RawCover: info.rawCover,
	}

	covIncrease := fuzzer.Config.Corpus.Save(input)
	covChanged := covIncrease > 0
	// At this point, we are certain that the request that started this triage job did indeed
	// increase the coverage. Some triage jobs come from other sources (e.g. seed or candidate,
	// they don't have a requestExecutionMode assigned => we ignore them
	fuzzer.mu.Lock()

	// increase corresponding dashboard stats
	if job.stat == statMinimize {
		// if the job was created by a minimize request, we want to attribute the coverage
		// uncovered via minimize to the corresponding requester

		nameContrib := fmt.Sprintf("%s (via %s)", ProfilingStatContribution(job.requesterStat, covChanged), statMinimize)
		nameBlocks := fmt.Sprintf("%s (via %s)", ProfilingStatBasicBlocksCoverage(job.requesterStat), statMinimize)

		fuzzer.stats[nameContrib]++
		fuzzer.stats[nameBlocks] += covIncrease
	} else {
		fuzzer.stats[ProfilingStatContribution(job.requesterStat, covChanged)]++
		fuzzer.stats[ProfilingStatBasicBlocksCoverage(job.requesterStat)] += covIncrease
	}

	// increase aggregated stats
	fuzzer.stats[ProfilingAllStatsContribution(covChanged)]++

	fuzzer.stats[ProfilingStatBasicBlocksCoverage("TEST ALL STATS")] += covIncrease // FIXME NICOLAS REMOVE
	fuzzer.mu.Unlock()

	if fuzzer.Config.NewInputs != nil {
		select {
		case <-fuzzer.ctx.Done():
		case fuzzer.Config.NewInputs <- input:
		}
	}
}

type deflakedCover struct {
	stableSignal    signal.Signal
	newStableSignal signal.Signal
	cover           cover.Cover
	rawCover        []uint32
}

func (job *triageJob) deflake(fuzzer *Fuzzer) (info deflakedCover, stop bool) {
	const signalRuns = 3
	var notExecuted int
	for i := 0; i < signalRuns; i++ {
		result := fuzzer.exec(job, &Request{
			Prog:          job.p,
			NeedSignal:    true,
			NeedCover:     true,
			NeedRawCover:  fuzzer.Config.FetchRawCover,
			stat:          statTriage,
			flags:         progInTriage,
			requesterStat: job.requesterStat,
		})
		if result.Stop {
			stop = true
			return
		}
		if !reexecutionSuccess(result.Info, &job.info, job.call) {
			// The call was not executed or failed.
			notExecuted++
			if notExecuted >= signalRuns/2+1 {
				stop = true
				return // if happens too often, give up
			}
			continue
		}
		thisSignal, thisCover := getSignalAndCover(job.p, result.Info, job.call)
		if len(info.rawCover) == 0 && fuzzer.Config.FetchRawCover {
			info.rawCover = thisCover
		}
		if i == 0 {
			info.stableSignal = thisSignal
			info.newStableSignal = job.newSignal.Intersection(thisSignal)
		} else {
			info.stableSignal = info.stableSignal.Intersection(thisSignal)
			info.newStableSignal = info.newStableSignal.Intersection(thisSignal)
		}
		if info.newStableSignal.Empty() {
			return
		}
		info.cover.Merge(thisCover)
	}
	return
}

func (job *triageJob) minimize(fuzzer *Fuzzer, newSignal signal.Signal) (stop bool) {
	const minimizeAttempts = 3
	job.p, job.call = prog.Minimize(job.p, job.call, false,
		func(p1 *prog.Prog, call1 int) bool {
			if stop {
				return false
			}
			for i := 0; i < minimizeAttempts; i++ {
				result := fuzzer.exec(job, &Request{
					Prog:          p1,
					NeedSignal:    true,
					stat:          statMinimize,
					requesterStat: job.requesterStat,
				})
				if result.Stop {
					stop = true
					return false
				}
				info := result.Info
				if !reexecutionSuccess(info, &job.info, call1) {
					// The call was not executed or failed.
					continue
				}
				thisSignal, _ := getSignalAndCover(p1, info, call1)
				if newSignal.Intersection(thisSignal).Len() == newSignal.Len() {
					return true
				}
			}
			return false
		})
	return stop
}

func reexecutionSuccess(info *ipc.ProgInfo, oldInfo *ipc.CallInfo, call int) bool {
	if info == nil || len(info.Calls) == 0 {
		return false
	}
	if call != -1 {
		// Don't minimize calls from successful to unsuccessful.
		// Successful calls are much more valuable.
		if oldInfo.Errno == 0 && info.Calls[call].Errno != 0 {
			return false
		}
		return len(info.Calls[call].Signal) != 0
	}
	return len(info.Extra.Signal) != 0
}

func getSignalAndCover(p *prog.Prog, info *ipc.ProgInfo, call int) (signal.Signal, []uint32) {
	inf := &info.Extra
	if call != -1 {
		inf = &info.Calls[call]
	}
	return signal.FromRaw(inf.Signal, signalPrio(p, inf, call)), inf.Cover
}

type smashJob struct {
	p    *prog.Prog
	call int
	jobPriority
}

func (job *smashJob) run(fuzzer *Fuzzer) {
	// smashJob simply starts a hintsJob and performs 100 mutations. We can simply omit
	// these operations and return instantly
	if !ablation_flags.ABLATION_MODE_SMASH_ENABLED {
		return
	}

	fuzzer.Logf(2, "smashing the program %s (call=%d):", job.p, job.call)
	if fuzzer.Config.Comparisons && job.call >= 0 {
		fuzzer.startJob(&hintsJob{
			p:           job.p.Clone(),
			call:        job.call,
			jobPriority: newJobPriority(smashPrio),
		})
	}

	fuzzer.profilingStats.IncModeCounter(ProfilingStatModeSmash)
	start := time.Now()

	const iters = 100
	rnd := fuzzer.rand()
	for i := 0; i < iters; i++ {
		p := job.p.Clone()

		fuzzer.profilingStats.IncModeCounter(ProfilingStatModeMutateFromSmash)
		startInside := time.Now()

		obs := p.MutateWithObserver(rnd, prog.RecommendedCalls,
			fuzzer.ChoiceTable(),
			fuzzer.Config.NoMutateCalls,
			fuzzer.Config.Corpus.Programs(),
		)

		deltaInside := time.Since(startInside)
		fuzzer.profilingStats.AddModeDuration(ProfilingStatModeMutateFromSmash, deltaInside)
		profileMutateObserver(fuzzer, obs)

		result := fuzzer.exec(job, &Request{
			Prog:          p,
			NeedSignal:    true,
			stat:          statSmash, // FIXME NICOLAS RECURSION
			requesterStat: statFuzzFromSmash,
		})
		if result.Stop {
			return
		}
		if fuzzer.Config.Collide {
			result := fuzzer.exec(job, &Request{
				Prog:          randomCollide(p, rnd),
				stat:          statCollide,
				requesterStat: statCollide,
			})
			if result.Stop {
				return
			}
		}
	}
	if fuzzer.Config.FaultInjection && job.call >= 0 {
		job.faultInjection(fuzzer)
	}

	delta := time.Since(start)
	fuzzer.profilingStats.AddModeDuration(ProfilingStatModeSmash, delta)
}

func randomCollide(origP *prog.Prog, rnd *rand.Rand) *prog.Prog {
	if rnd.Intn(5) == 0 {
		// Old-style collide with a 20% probability.
		p, err := prog.DoubleExecCollide(origP, rnd)
		if err == nil {
			return p
		}
	}
	if rnd.Intn(4) == 0 {
		// Duplicate random calls with a 20% probability (25% * 80%).
		p, err := prog.DupCallCollide(origP, rnd)
		if err == nil {
			return p
		}
	}
	p := prog.AssignRandomAsync(origP, rnd)
	if rnd.Intn(2) != 0 {
		prog.AssignRandomRerun(p, rnd)
	}
	return p
}

func (job *smashJob) faultInjection(fuzzer *Fuzzer) {
	for nth := 1; nth <= 100; nth++ {
		fuzzer.Logf(2, "injecting fault into call %v, step %v",
			job.call, nth)
		newProg := job.p.Clone()
		newProg.Calls[job.call].Props.FailNth = nth
		result := fuzzer.exec(job, &Request{
			Prog:          job.p,
			stat:          statSmash,
			requesterStat: statSmash, // FIXME NICOLAS maybe distinguish this cacse?
		})
		if result.Stop {
			return
		}
		info := result.Info
		if info != nil && len(info.Calls) > job.call &&
			info.Calls[job.call].Flags&ipc.CallFaultInjected == 0 {
			break
		}
	}
}

type hintsJob struct {
	p    *prog.Prog
	call int
	jobPriority
}

func (job *hintsJob) run(fuzzer *Fuzzer) {
	// similarly to smashJob, we can simply omit the mutations if
	// syzkaller's mutate with hints mode is disabled
	if !ablation_flags.ABLATION_MODE_HINTS_ENABLED {
		return
	}

	// First execute the original program to dump comparisons from KCOV.
	p := job.p
	result := fuzzer.exec(job, &Request{
		Prog:          p,
		NeedHints:     true,
		stat:          statSeed,
		requesterStat: statSeedFromHint,
	})
	if result.Stop || result.Info == nil {
		return
	}
	// Then mutate the initial program for every match between
	// a syscall argument and a comparison operand.
	// Execute each of such mutants to check if it gives new coverage.
	fuzzer.profilingStats.IncModeCounter(ProfilingStatModeMutateHints)
	start := time.Now()

	p.MutateWithHints(
		job.call,
		result.Info.Calls[job.call].Comps,
		func(p *prog.Prog) bool {
			result := fuzzer.exec(job, &Request{
				Prog:          p,
				NeedSignal:    true,
				stat:          statHint, // FIXME NICOLAS RECURSION
				requesterStat: statHint,
			})
			return !result.Stop
		},
	)

	delta := time.Since(start)
	fuzzer.profilingStats.AddModeDuration(ProfilingStatModeMutateHints, delta)
}
