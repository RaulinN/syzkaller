package fuzzer

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProfilingModeName string
type ProfilingMutatorName string

const prefix = "[prof]"
const (
	ProfilingStatModeGenerate        ProfilingModeName = prefix + " mode generate"
	ProfilingStatModeMutate          ProfilingModeName = prefix + " mode mutate (= fuzz)"
	ProfilingStatModeMutateHints     ProfilingModeName = prefix + " mode mutate with hints"
	ProfilingStatModeSmash           ProfilingModeName = prefix + " mode smash"
	ProfilingStatModeMutateFromSmash ProfilingModeName = prefix + " mode mutate (from smash)"
)

func ProfilingStatContribution(requesterStat string, coverageIncrease bool) string {
	op := "increase"
	if !coverageIncrease {
		op = "did not change"
	}
	return fmt.Sprintf("%s %s > #times cov. %s", prefix, requesterStat, op)
}

func ProfilingAllStatsContribution(coverageIncrease bool) string {
	return ProfilingStatContribution("ALL requesterStats", coverageIncrease)
}

func ProfilingStatBasicBlocksCoverage(requesterStat string) string {
	return fmt.Sprintf("%s %s > #basic blocks covered", prefix, requesterStat)
}

const (
	ProfilingStatMutatorSquashAny  ProfilingMutatorName = prefix + " mutator squashAny"
	ProfilingStatMutatorSplice     ProfilingMutatorName = prefix + " mutator splice"
	ProfilingStatMutatorInsertCall ProfilingMutatorName = prefix + " mutator insertCall"
	ProfilingStatMutatorMutateArg  ProfilingMutatorName = prefix + " mutator mutateArg"
	ProfilingStatMutatorRemoveCall ProfilingMutatorName = prefix + " mutator removeCall"
	ProfilingStatMutatorShuffle    ProfilingMutatorName = prefix + " mutator shuffle"
)

// careful, a new slice is generated each time. Don't abuse
func allModes() []ProfilingModeName {
	return []ProfilingModeName{
		ProfilingStatModeGenerate,
		ProfilingStatModeMutate,
		ProfilingStatModeMutateHints,
		ProfilingStatModeSmash,
		ProfilingStatModeMutateFromSmash,
	}
}

// careful, a new slice is generated each time. Don't abuse
func allMutators() []ProfilingMutatorName {
	return []ProfilingMutatorName{
		ProfilingStatMutatorSquashAny,
		ProfilingStatMutatorSplice,
		ProfilingStatMutatorInsertCall,
		ProfilingStatMutatorMutateArg,
		ProfilingStatMutatorRemoveCall,
		ProfilingStatMutatorShuffle,
	}
}

func (fuzzer *Fuzzer) AddTimeSpent(identifier string, duration time.Duration) {
	fuzzer.mu.Lock()
	defer fuzzer.mu.Unlock()

	fuzzer.stats[identifier] += uint64(duration.Nanoseconds())
}

func (fuzzer *Fuzzer) AddModeTimeSpent(mode ProfilingModeName, duration time.Duration) {
	id := fmt.Sprintf("%v > time spent (!= executing) (ns)", mode)
	fuzzer.AddTimeSpent(id, duration)
}

// https://siongui.github.io/2016/01/30/go-pretty-print-variable/
func Prettify(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ") // to json
	if err == nil {
		return string(b), nil
	}
	return "", err
}

func ToJson(v interface{}) (string, error) {
	b, err := json.Marshal(v) // to json
	if err == nil {
		return string(b), nil
	}
	return "", err
}
