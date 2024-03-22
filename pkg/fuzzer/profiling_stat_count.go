package fuzzer

import "sync/atomic"

type StatCount uint64

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
