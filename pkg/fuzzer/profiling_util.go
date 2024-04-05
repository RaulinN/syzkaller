package fuzzer

import (
	"encoding/json"
)

type ProfilingModeName string
type ProfilingMutatorName string

const prefix = "[profiling] "
const (
	ProfilingStatModeGenerate        ProfilingModeName = prefix + "mode generate"
	ProfilingStatModeMutate          ProfilingModeName = prefix + "mode mutate"
	ProfilingStatModeMutateHints     ProfilingModeName = prefix + "mode mutate with hints"
	ProfilingStatModeSmash           ProfilingModeName = prefix + "mode smash"
	ProfilingStatModeMutateFromSmash ProfilingModeName = prefix + "mode mutate (from smash)"
)

type ProfilingModeCoverageContribution = string

func ProfilingStatContribution(mode ProfilingModeName) ProfilingModeCoverageContribution {
	return string(mode) + " > coverage contribution"
}

const (
	ProfilingStatMutatorSquashAny  ProfilingMutatorName = prefix + "mutator squashAny"
	ProfilingStatMutatorSplice     ProfilingMutatorName = prefix + "mutator splice"
	ProfilingStatMutatorInsertCall ProfilingMutatorName = prefix + "mutator insertCall"
	ProfilingStatMutatorMutateArg  ProfilingMutatorName = prefix + "mutator mutateArg"
	ProfilingStatMutatorRemoveCall ProfilingMutatorName = prefix + "mutator removeCall"
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
	}
}

// https://siongui.github.io/2016/01/30/go-pretty-print-variable/
func Prettify(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ") // to json
	if err == nil {
		return string(b), nil
	}
	return "", err
}
