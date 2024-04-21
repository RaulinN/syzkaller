package profiler

import (
	"encoding/json"
	"os"
	"sync"
)

type AblationConfiguration struct {
	// create the configuration only once
	Once sync.Once
	// flags disable modes
	DisableModeGenerate bool `json:"disable_mode_generate"`
	DisableModeHints    bool `json:"disable_mode_hints"`
	DisableModeMutate   bool `json:"disable_mode_mutate"`
	DisableModeSmash    bool `json:"disable_mode_smash"`
	// flags disable mutators
	AnyMutatorEnabled        bool
	DisableMutatorInsertCall bool `json:"disable_mutator_insert_call"`
	DisableMutatorMutateArg  bool `json:"disable_mutator_mutate_arg"`
	DisableMutatorRemoveCall bool `json:"disable_mutator_remove_call"`
	DisableMutatorSplice     bool `json:"disable_mutator_splice"`
	DisableMutatorSquashAny  bool `json:"disable_mutator_squash_any"`
	// flags disable stages
	DisableStageCollide  bool `json:"disable_stage_collide"`
	DisableStageMinimize bool `json:"disable_stage_minimize"`
	// flags reduce mutators
	ReduceMutatorInsertCall bool `json:"reduce_mutator_insert_call"`
	ReduceMutatorMutateArg  bool `json:"reduce_mutator_mutate_arg"`
	ReduceMutatorSplice     bool `json:"reduce_mutator_splice"`
}

var (
	AblationConfig = AblationConfiguration{
		AnyMutatorEnabled: true,
	}
)

func SetupAblationConfig(filename string, cfg *AblationConfiguration) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return err
	}

	cfg.AnyMutatorEnabled = !(cfg.DisableMutatorInsertCall &&
		cfg.DisableMutatorMutateArg &&
		cfg.DisableMutatorRemoveCall &&
		cfg.DisableMutatorSplice &&
		cfg.DisableMutatorSquashAny)

	return nil
}
