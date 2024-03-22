//go:build profiling

package fuzzer

import (
	"fmt"
)

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
