//go:build profiling

package fuzzer

import (
	"time"
)

type ProfilingStats struct {
	// counters for the modes of operations. The mutate mode is split
	// into two separate counters: simple mutate operations and mutate
	// operations performed as part of a smash request
	countModes map[ProfilingModeName]*StatCount

	// counters for individual mutators (see mutation.profiling.go)
	countMutators map[ProfilingMutatorName]*StatCount

	// time elapsed executing modes of operations
	durationModes map[ProfilingModeName]*AtomicDuration
}

func NewProfilingStats() *ProfilingStats {
	ps := ProfilingStats{
		countModes:    make(map[ProfilingModeName]*StatCount),
		countMutators: make(map[ProfilingMutatorName]*StatCount),
		durationModes: make(map[ProfilingModeName]*AtomicDuration),
	}

	for _, mode := range allModes() {
		var stat StatCount = 0
		ps.countModes[mode] = &stat
		ps.durationModes[mode] = NewAtomicDuration()
	}

	for _, mutator := range allMutators() {
		var stat StatCount = 0
		ps.countMutators[mutator] = &stat
	}

	return &ps
}

func (ps *ProfilingStats) allCounts() map[string]uint64 {
	// TODO make cleaner
	return map[string]uint64{
		// modes of operation
		string(ProfilingStatModeGenerate):        ps.countModes[ProfilingStatModeGenerate].get(),
		string(ProfilingStatModeMutate):          ps.countModes[ProfilingStatModeMutate].get(),
		string(ProfilingStatModeMutateHints):     ps.countModes[ProfilingStatModeMutateHints].get(),
		string(ProfilingStatModeSmash):           ps.countModes[ProfilingStatModeSmash].get(),
		string(ProfilingStatModeMutateFromSmash): ps.countModes[ProfilingStatModeMutateFromSmash].get(),
		// mutators
		string(ProfilingStatMutatorSquashAny):  ps.countMutators[ProfilingStatMutatorSquashAny].get(),
		string(ProfilingStatMutatorSplice):     ps.countMutators[ProfilingStatMutatorSplice].get(),
		string(ProfilingStatMutatorInsertCall): ps.countMutators[ProfilingStatMutatorInsertCall].get(),
		string(ProfilingStatMutatorMutateArg):  ps.countMutators[ProfilingStatMutatorMutateArg].get(),
		string(ProfilingStatMutatorRemoveCall): ps.countMutators[ProfilingStatMutatorRemoveCall].get(),
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
	ps.countModes[mode].inc()
}

func (ps *ProfilingStats) AddMutatorCounter(mutator ProfilingMutatorName, value int) {
	ps.countMutators[mutator].add(value)
}

func (ps *ProfilingStats) IncMutatorCounter(mutator ProfilingMutatorName) {
	ps.AddMutatorCounter(mutator, 1)
}

func (ps *ProfilingStats) AddModeDuration(mode ProfilingModeName, duration time.Duration) {
	ps.durationModes[mode].Add(duration)
}
