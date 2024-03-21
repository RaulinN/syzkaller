//go:build profiling

package fuzzer

import (
	"sync"
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
	lock        sync.RWMutex
	totalCounts map[ProfilingCounterName]StatCount
	deltaCounts map[ProfilingCounterName]StatCount

	// countModeGenerate    StatCount
	// countModeMutate      StatCount
	// countModeMutateHints StatCount
	// countModeSmash       StatCount

	// durationModeGenerate    StatDuration
	// durationModeMutate      StatDuration
	// durationModeMutateHints StatDuration
	// durationModeSmash       StatDuration

	// countMutator StatCount
}

func NewProfilingStats() *ProfilingStats {
	return &ProfilingStats{
		lock:        sync.RWMutex{},
		totalCounts: map[ProfilingCounterName]StatCount{},
		deltaCounts: map[ProfilingCounterName]StatCount{},
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

		for {
			time.Sleep(10 * time.Second)

			fuzzer.Logf(0, "logging total counts: %v", fuzzer.profilingStats.totalCounts)
			fuzzer.Logf(0, "logging delta counts: %v", fuzzer.profilingStats.deltaCounts)

			// TODO lock the stats map?
			//fuzzer.profilingStats.lock.Lock()

			var oldValue uint64 = 0
			for _, modeName := range modes {
				stat := fuzzer.profilingStats.deltaCounts[modeName]
				stat.reset(&oldValue)

				fuzzer.stats[string(modeName)] = oldValue
			}

			//fuzzer.profilingStats.lock.Unlock()
		}
	}()
}

func (ps *ProfilingStats) IncCounter(counterName ProfilingCounterName) {
	//ps.lock.Lock()
	//defer ps.lock.Unlock()

	stat := ps.totalCounts[counterName]
	stat.inc()
	stat = ps.deltaCounts[counterName]
	stat.inc()
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
