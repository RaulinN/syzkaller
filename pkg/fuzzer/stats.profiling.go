//go:build profiling

// Copyright 2024 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package fuzzer

const (
	statGenerate       = "exec gen"
	statFuzz           = "exec fuzz"
	statCandidate      = "exec candidate"
	statTriage         = "exec triage"
	statMinimize       = "exec minimize"
	statSmash          = "exec smash"
	statHint           = "exec hints"
	statSeed           = "exec seeds"
	statCollide        = "exec collide"
	statExecTotal      = "exec total"
	statBufferTooSmall = "buffer too small"
	statFuzzFromSmash  = "exec fuzz (from smash)"
	statSeedFromHint   = "exec seeds (from hint)"
)

func (fuzzer *Fuzzer) GrabStats() map[string]uint64 {
	fuzzer.mu.Lock()
	defer fuzzer.mu.Unlock()
	ret := fuzzer.stats
	fuzzer.stats = map[string]uint64{}
	return ret
}

func (fuzzer *Fuzzer) GrabAllStats() map[string]uint64 {
	r := fuzzer.GrabStats()

	fuzzer.mu.Lock()
	defer fuzzer.mu.Unlock()

	r["running jobs"] = uint64(fuzzer.runningJobs.Load())
	r["queued candidates"] = uint64(fuzzer.queuedCandidates.Load())
	r["exec queue size (prio queue of req)"] = uint64(fuzzer.nextExec.Len())
	r["queued requests (running req)"] = uint64(len(fuzzer.runningExecs))

	return r
}
