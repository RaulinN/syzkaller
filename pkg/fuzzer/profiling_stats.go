//go:build profiling

package fuzzer

import (
	"sync/atomic"
	"time"
)

type ProfilingCounterName string

const prefix = "[profiling] "
const (
	ProfilingStatModeGenerate    ProfilingCounterName = prefix + "mode generate"
	ProfilingStatModeMutate      ProfilingCounterName = prefix + "mode mutate"
	ProfilingStatModeMutateHints ProfilingCounterName = prefix + "mode mutate with hints"
	ProfilingStatModeSmash       ProfilingCounterName = prefix + "mode smash"
)

type StatCount uint64

// type StatDuration time.Duration // FIXME

type ProfilingStats struct {
	// lock        sync.RWMutex

	//totalCounts map[ProfilingCounterName]StatCount
	//deltaCounts map[ProfilingCounterName]StatCount

	totalModeGenerate    StatCount
	totalModeMutate      StatCount
	totalModeMutateHints StatCount
	totalModeSmash       StatCount

	// durationModeGenerate    StatDuration
	// durationModeMutate      StatDuration
	// durationModeMutateHints StatDuration
	// durationModeSmash       StatDuration

	// countMutator StatCount
}

func NewProfilingStats() *ProfilingStats {
	return &ProfilingStats{
		//lock:        sync.RWMutex{},
		//totalCounts: map[ProfilingCounterName]StatCount{},
		//deltaCounts: map[ProfilingCounterName]StatCount{},
	}
}

func (ps *ProfilingStats) allCounts() map[ProfilingCounterName]uint64 {
	return map[ProfilingCounterName]uint64{
		ProfilingStatModeGenerate:    ps.totalModeGenerate.get(),
		ProfilingStatModeMutate:      ps.totalModeMutate.get(),
		ProfilingStatModeMutateHints: ps.totalModeMutateHints.get(),
		ProfilingStatModeSmash:       ps.totalModeSmash.get(),
	}
}

func (fuzzer *Fuzzer) StartProfilingLogger() {
	// TODO go routine start log on the file
	// TODO go routine logs on the dashboard
	// TODO log actual nice string instead of object
	go func() {
		modes := [4]ProfilingCounterName{
			ProfilingStatModeGenerate,
			ProfilingStatModeMutate,
			ProfilingStatModeMutateHints,
			ProfilingStatModeSmash,
		}

		prevCounts := map[ProfilingCounterName]uint64{}

		for {
			time.Sleep(10 * time.Second)

			counts := fuzzer.profilingStats.allCounts()

			fuzzer.Logf(0, "logging total counts: %v", counts)

			// TODO lock the stats map?
			//fuzzer.profilingStats.lock.Lock()

			for _, modeName := range modes {
				current := counts[modeName]
				delta := current - prevCounts[modeName]
				fuzzer.stats[string(modeName)] = delta
				prevCounts[modeName] = current
			}

			//fuzzer.profilingStats.lock.Unlock()
		}
	}()
}

func (ps *ProfilingStats) IncCounter(counterName ProfilingCounterName) {
	switch counterName {
	case ProfilingStatModeGenerate:
		ps.totalModeGenerate.inc()
	case ProfilingStatModeMutate:
		ps.totalModeMutate.inc()
	case ProfilingStatModeMutateHints:
		ps.totalModeMutateHints.inc()
	case ProfilingStatModeSmash:
		ps.totalModeSmash.inc()
	default:
		// FIXME
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
