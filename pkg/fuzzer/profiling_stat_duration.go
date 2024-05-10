package fuzzer

import (
	"sync"
	"time"
)

type AtomicDuration struct {
	lock     sync.Mutex
	duration time.Duration
}

func NewAtomicDuration() *AtomicDuration {
	return &AtomicDuration{
		lock:     sync.Mutex{},
		duration: time.Duration(0),
	}
}

func (ad *AtomicDuration) Get() time.Duration {
	ad.lock.Lock()
	defer ad.lock.Unlock()

	return ad.duration
}

func (ad *AtomicDuration) Add(duration time.Duration) {
	ad.lock.Lock()
	defer ad.lock.Unlock()

	ad.duration += duration
}

func (ad *AtomicDuration) Set(duration time.Duration) {
	ad.lock.Lock()
	defer ad.lock.Unlock()

	ad.duration = duration
}

func (ad *AtomicDuration) Reset() {
	ad.Set(time.Duration(0))
}

func (ad *AtomicDuration) String() string {
	ad.lock.Lock()
	defer ad.lock.Unlock()

	return ad.duration.String()
}
