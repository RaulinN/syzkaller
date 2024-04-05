package fuzzer

import "time"

func (fuzzer *Fuzzer) StartProfilingLogger() {
	go func() {
		modes := allModes()
		mutators := allMutators()

		prevCounts := map[string]uint64{}

		for {
			time.Sleep(15 * time.Second)

			counts := fuzzer.profilingStats.allCounts()
			prettyCounts, err := Prettify(counts)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding counts map to JSON")
			}
			fuzzer.Logf(0, "logging total counts: %v", prettyCounts)

			durations := fuzzer.profilingStats.allDurations()
			displayDurations := map[string]string{}
			for k, v := range durations {
				displayDurations[k] = v.String()
			}
			prettyDurations, err := Prettify(displayDurations)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding duration map to JSON")
			}
			ptest, _ := Prettify(durations) // FIXME remove
			fuzzer.Logf(0, "logging total durations (1 - ints): %v", ptest)
			fuzzer.Logf(0, "logging total durations (2 - hh:mm:ss): %v", prettyDurations)
			fuzzer.Logf(0, "------------------------------------------------")

			// TODO display durations on dashboard?
			fuzzer.mu.Lock()
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

			fuzzer.mu.Unlock()
		}
	}()
}
