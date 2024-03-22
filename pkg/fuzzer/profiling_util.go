package fuzzer

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
