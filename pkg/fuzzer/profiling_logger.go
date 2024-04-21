package fuzzer

import "time"

func (fuzzer *Fuzzer) StartProfilingLogger() {
	go func() {
		modes := allModes()
		mutators := allMutators()

		prevCounts := map[string]uint64{}

		// FIXME NICOLAS I believe this will actually be reset every hour (on executor creation) => make it persistent
		for {
			time.Sleep(60 * time.Second)
			now := time.Now().Unix()

			counts := fuzzer.profilingStats.allCounts()
			countsJson, err := ToJson(counts)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding counts map to JSON")
			}

			durations := fuzzer.profilingStats.allDurations()
			displayDurations := map[string]string{}
			for k, v := range durations {
				displayDurations[k] = v.String()
			}
			durationsJson, err := ToJson(displayDurations)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding duration map to JSON")
			}

			fuzzer.Logf(0, "%v;logging total counts:%v", now, countsJson)
			fuzzer.Logf(0, "%v;logging total durations (hh:mm:ss):%v", now, durationsJson)

			// log all fuzzer stats stats
			stats := fuzzer.GrabAllStats()
			statsJson, err := ToJson(stats)
			if err != nil {
				fuzzer.Logf(0, "ERROR encoding stats map to JSON")
			}
			fuzzer.Logf(0, "%v;logging all stats:%v", now, statsJson)

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
