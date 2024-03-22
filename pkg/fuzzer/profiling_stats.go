//go:build profiling

package fuzzer

import (
	"fmt"
	"time"
)

type ProfilingStats struct {
	// counters for the modes of operations. The mutate mode is split
	// into two separate counters: simple mutate operations and mutate
	// operations performed as part of a smash request
	countModeGenerate        StatCount
	countModeMutate          StatCount
	countModeMutateHints     StatCount
	countModeSmash           StatCount
	countModeMutateFromSmash StatCount

	// counters for individual mutators (see mutation.profiling.go)
	countMutatorSquashAny  StatCount
	countMutatorSplice     StatCount
	countMutatorInsertCall StatCount
	countMutatorMutateArg  StatCount
	countMutatorRemoveCall StatCount

	// time elapsed executing modes of operations
	/*durationModeGenerate        *AtomicDuration
	durationModeMutate          *AtomicDuration
	durationModeMutateHints     *AtomicDuration
	durationModeSmash           *AtomicDuration
	durationModeMutateFromSmash *AtomicDuration*/

	durationModes map[ProfilingModeName]*AtomicDuration
}

func NewProfilingStats() *ProfilingStats {
	ps := ProfilingStats{
		/*durationModeGenerate:        NewAtomicDuration(),
		durationModeMutate:          NewAtomicDuration(),
		durationModeMutateHints:     NewAtomicDuration(),
		durationModeSmash:           NewAtomicDuration(),
		durationModeMutateFromSmash: NewAtomicDuration(),*/
		durationModes: make(map[ProfilingModeName]*AtomicDuration),
	}

	for _, mode := range allModes() {
		ps.durationModes[mode] = NewAtomicDuration()
	}

	return &ps
}

func (ps *ProfilingStats) allCounts() map[string]uint64 {
	return map[string]uint64{
		// modes of operation
		string(ProfilingStatModeGenerate):        ps.countModeGenerate.get(),
		string(ProfilingStatModeMutate):          ps.countModeMutate.get(),
		string(ProfilingStatModeMutateHints):     ps.countModeMutateHints.get(),
		string(ProfilingStatModeSmash):           ps.countModeSmash.get(),
		string(ProfilingStatModeMutateFromSmash): ps.countModeMutateFromSmash.get(),
		// mutators
		string(ProfilingStatMutatorSquashAny):  ps.countMutatorSquashAny.get(),
		string(ProfilingStatMutatorSplice):     ps.countMutatorSplice.get(),
		string(ProfilingStatMutatorInsertCall): ps.countMutatorInsertCall.get(),
		string(ProfilingStatMutatorMutateArg):  ps.countMutatorMutateArg.get(),
		string(ProfilingStatMutatorRemoveCall): ps.countMutatorRemoveCall.get(),
	}
}

func (ps *ProfilingStats) allDurations() map[string]time.Duration {
	return map[string]time.Duration{
		// modes of operation
		string(ProfilingStatModeGenerate):        ps.durationModes[ProfilingStatModeGenerate].Get(),
		string(ProfilingStatModeMutate):          ps.durationModes[ProfilingStatModeMutate].Get(),
		string(ProfilingStatModeMutateHints):     ps.durationModes[ProfilingStatModeMutateHints].Get(),
		string(ProfilingStatModeSmash):           ps.durationModes[ProfilingStatModeSmash].Get(),
		string(ProfilingStatModeMutateFromSmash): ps.durationModes[ProfilingStatModeMutateFromSmash].Get(),
	}
}

func (ps *ProfilingStats) IncModeCounter(mode ProfilingModeName) {
	switch mode {
	case ProfilingStatModeGenerate:
		ps.countModeGenerate.inc()
	case ProfilingStatModeMutate:
		ps.countModeMutate.inc()
	case ProfilingStatModeMutateHints:
		ps.countModeMutateHints.inc()
	case ProfilingStatModeSmash:
		ps.countModeSmash.inc()
	case ProfilingStatModeMutateFromSmash:
		ps.countModeMutateFromSmash.inc()
	default:
		panic(fmt.Sprintf("missing switch case for mode '%v' in IncModeCounter", string(mode)))
	}
}

func (ps *ProfilingStats) AddMutatorCounter(mutator ProfilingMutatorName, value int) {
	switch mutator {
	case ProfilingStatMutatorSquashAny:
		ps.countMutatorSquashAny.add(value)
	case ProfilingStatMutatorSplice:
		ps.countMutatorSplice.add(value)
	case ProfilingStatMutatorInsertCall:
		ps.countMutatorInsertCall.add(value)
	case ProfilingStatMutatorMutateArg:
		ps.countMutatorMutateArg.add(value)
	case ProfilingStatMutatorRemoveCall:
		ps.countMutatorRemoveCall.add(value)
	default:
		panic(fmt.Sprintf("missing switch case for mutator '%v' in IncMutatorCounter", string(mutator)))
	}
}

func (ps *ProfilingStats) IncMutatorCounter(mutator ProfilingMutatorName) {
	ps.AddMutatorCounter(mutator, 1)
}

/*
func (ps *ProfilingStats) AddModeDuration2(mode ProfilingModeName, duration time.Duration) {
	switch mode {
	case ProfilingStatModeGenerate:
		ps.durationModeGenerate.Add(duration)
	case ProfilingStatModeMutate:
		ps.durationModeMutate.Add(duration)
	case ProfilingStatModeMutateHints:
		ps.durationModeMutateHints.Add(duration)
	case ProfilingStatModeSmash:
		ps.durationModeSmash.Add(duration)
	case ProfilingStatModeMutateFromSmash:
		ps.durationModeMutateFromSmash.Add(duration)
	default:
		panic(fmt.Sprintf("missing switch case for mode '%v' in AddModeDuration", string(mode)))
	}
}*/

func (ps *ProfilingStats) AddModeDuration(mode ProfilingModeName, duration time.Duration) {
	ps.durationModes[mode].Add(duration)
}
