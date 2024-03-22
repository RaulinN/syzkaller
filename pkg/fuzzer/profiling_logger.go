package fuzzer

import "time"

func (fuzzer *Fuzzer) StartProfilingLogger() {
	go func() {
		modes := allModes()
		mutators := allMutators()

		prevCounts := map[string]uint64{}

		for {
			time.Sleep(30 * time.Second)

			counts := fuzzer.profilingStats.allCounts()
			prettyCounts, err := Prettify(counts)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding counts map to JSON")
			}
			fuzzer.Logf(0, "logging total counts: %v", prettyCounts)

			durations := fuzzer.profilingStats.allDurations()
			prettyDurations, err := Prettify(durations)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding duration map to JSON")
			}
			fuzzer.Logf(0, "logging total durations: %v", prettyDurations)

			// TODO lock the stats map?
			// TODO display durations on dashboard?
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
