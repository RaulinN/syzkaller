//go:build profiling

package fuzzer

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ProfilingModeName string
type ProfilingMutatorName string

const prefix = "[profiling] "
const (
	ProfilingStatModeGenerate        ProfilingModeName = prefix + "mode generate"
	ProfilingStatModeMutate          ProfilingModeName = prefix + "mode mutate"
	ProfilingStatModeMutateHints     ProfilingModeName = prefix + "mode mutate with hints"
	ProfilingStatModeSmash           ProfilingModeName = prefix + "mode smash"
	ProfilingStatModeMutateFromSmash ProfilingModeName = prefix + "mode mutate (from smash)"
)

const (
	ProfilingStatMutatorSquashAny  ProfilingMutatorName = prefix + "mutator squashAny"
	ProfilingStatMutatorSplice     ProfilingMutatorName = prefix + "mutator splice"
	ProfilingStatMutatorInsertCall ProfilingMutatorName = prefix + "mutator insertCall"
	ProfilingStatMutatorMutateArg  ProfilingMutatorName = prefix + "mutator mutateArg"
	ProfilingStatMutatorRemoveCall ProfilingMutatorName = prefix + "mutator removeCall"
)

type StatCount uint64

// type StatDuration time.Duration // FIXME

type ProfilingStats struct {
	countModeGenerate        StatCount
	countModeMutate          StatCount
	countModeMutateHints     StatCount
	countModeSmash           StatCount
	countModeMutateFromSmash StatCount

	countMutatorSquashAny  StatCount
	countMutatorSplice     StatCount
	countMutatorInsertCall StatCount
	countMutatorMutateArg  StatCount
	countMutatorRemoveCall StatCount

	// durationModeGenerate    StatDuration
	// durationModeMutate      StatDuration
	// durationModeMutateHints StatDuration
	// durationModeSmash       StatDuration
}

func NewProfilingStats() *ProfilingStats {
	return &ProfilingStats{}
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

func (fuzzer *Fuzzer) StartProfilingLogger() {
	// TODO log actual nice string instead of object
	go func() {
		modes := []ProfilingModeName{
			ProfilingStatModeGenerate,
			ProfilingStatModeMutate,
			ProfilingStatModeMutateHints,
			ProfilingStatModeSmash,
			ProfilingStatModeMutateFromSmash,
		}

		mutators := []ProfilingMutatorName{
			ProfilingStatMutatorSquashAny,
			ProfilingStatMutatorSplice,
			ProfilingStatMutatorInsertCall,
			ProfilingStatMutatorMutateArg,
			ProfilingStatMutatorRemoveCall,
		}

		prevCounts := map[string]uint64{}

		for {
			time.Sleep(10 * time.Second)

			counts := fuzzer.profilingStats.allCounts()

			fuzzer.Logf(0, "logging total counts: %v", counts)

			// TODO lock the stats map?

			for _, mode := range modes {
				modeName := string(mode)
				current := counts[modeName]
				delta := current - prevCounts[modeName]

				fuzzer.stats[modeName] = delta
				prevCounts[modeName] = current
			}

			for _, mutator := range mutators {
				mutatorName := string(mutator)
				current := counts[mutatorName]
				delta := current - prevCounts[mutatorName]

				fuzzer.stats[mutatorName] = delta
				prevCounts[mutatorName] = current
			}
		}
	}()
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

func (ps *ProfilingStats) IncMutatorCounter(mutator ProfilingMutatorName) {
	switch mutator {
	case ProfilingStatMutatorSquashAny:
		ps.countMutatorSquashAny.inc()
	case ProfilingStatMutatorSplice:
		ps.countMutatorSplice.inc()
	case ProfilingStatMutatorInsertCall:
		ps.countMutatorInsertCall.inc()
	case ProfilingStatMutatorMutateArg:
		ps.countMutatorMutateArg.inc()
	case ProfilingStatMutatorRemoveCall:
		ps.countMutatorRemoveCall.inc()
	default:
		panic(fmt.Sprintf("missing switch case for mutator '%v' in IncMutatorCounter", string(mutator)))
	}
}

func (s *StatCount) get() uint64 {
	return atomic.LoadUint64((*uint64)(s))
}

func (s *StatCount) inc() {
	s.add(1)
}

func (s *StatCount) add(v int) {
	atomic.AddUint64((*uint64)(s), uint64(v))
}

func (s *StatCount) set(v int) {
	atomic.StoreUint64((*uint64)(s), uint64(v))
}

func (s *StatCount) reset(swapContainer *uint64) {
	atomic.SwapUint64(swapContainer, uint64(0))
}
