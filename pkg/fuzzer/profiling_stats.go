//go:build profiling

package fuzzer

import (
	"sync/atomic"
	"time"
)

type StatCount uint64
type StatDuration time.Duration // FIXME

type ProfilingStats struct {
	countModeGenerate    StatCount
	countModeMutate      StatCount
	countModeMutateHints StatCount
	countModeSmash       StatCount

	durationModeGenerate    StatDuration
	durationModeMutate      StatDuration
	durationModeMutateHints StatDuration
	durationModeSmash       StatDuration

	countMutator StatCount
}

func NewProfilingStats() *ProfilingStats {
	return &ProfilingStats{}
}

// func l(level int, msg string, args ...interface{}
// all maps are passed by reference
func (ps *ProfilingStats) StartLogger(log func(int, string, ...interface{}), stats map[string]uint64) {
	// TODO go routine start log on the file
	// TODO go routine logs on the dashboard
	go func() {
		for {
			time.Sleep(10 * time.Second)
			ps.countMutator.inc()
			log(0, "logging from the coroutine (1): %v", ps)
			stats["CUSTOM_CLASS_STAT"] = ps.countMutator.get()
			log(0, "logging from the coroutine (2): %v", ps.countMutator.get())
		}
	}()
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
