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

const (
	ProfilingStatMutatorSquashAny  ProfilingMutatorName = prefix + "mutator squashAny"
	ProfilingStatMutatorSplice     ProfilingMutatorName = prefix + "mutator splice"
	ProfilingStatMutatorInsertCall ProfilingMutatorName = prefix + "mutator insertCall"
	ProfilingStatMutatorMutateArg  ProfilingMutatorName = prefix + "mutator mutateArg"
	ProfilingStatMutatorRemoveCall ProfilingMutatorName = prefix + "mutator removeCall"
)

// type StatDuration time.Duration // FIXME

// https://siongui.github.io/2016/01/30/go-pretty-print-variable/
func Prettify(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ") // to json
	if err == nil {
		return string(b), nil
	}
	return "", err
}
