package fuzzer

import "time"

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
			time.Sleep(30 * time.Second)

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
