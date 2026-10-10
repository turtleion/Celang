package bytecode

import (
	"fmt"

	"celang/pkg/evaluator"
)

// CompiledFunction represents an executable function compiled to bytecode.
type CompiledFunction struct {
	Name      string
	NumParams int
	Chunk     *Chunk
}

func (cf *CompiledFunction) Type() evaluator.ObjectType { return "COMPILED_FUNCTION" }
func (cf *CompiledFunction) Inspect() string {
	if cf.Name == "" {
		return fmt.Sprintf("<fungsi anonim (%d params)>", cf.NumParams)
	}
	return fmt.Sprintf("<fungsi %s (%d params)>", cf.Name, cf.NumParams)
}
