package bytecode

import (
	"fmt"

	"celang/pkg/evaluator"
)

// VMRunnerFunc executes a compiled function with given arguments via VM.
type VMRunnerFunc func(fn *CompiledFunction, args []evaluator.Object) evaluator.Object

var globalVMRunner VMRunnerFunc

// SetVMRunner registers the VM invocation engine.
func SetVMRunner(runner VMRunnerFunc) {
	globalVMRunner = runner
}

// CompiledFunction represents an executable function compiled to bytecode.
type CompiledFunction struct {
	Name       string
	NumParams  int
	Chunk      *Chunk
	LocalNames []string
}

func (cf *CompiledFunction) Type() evaluator.ObjectType { return "COMPILED_FUNCTION" }
func (cf *CompiledFunction) Inspect() string {
	if cf.Name == "" {
		return fmt.Sprintf("<fungsi anonim (%d params)>", cf.NumParams)
	}
	return fmt.Sprintf("<fungsi %s (%d params)>", cf.Name, cf.NumParams)
}

// Call implements the evaluator.Callable interface.
func (cf *CompiledFunction) Call(args ...evaluator.Object) evaluator.Object {
	if globalVMRunner != nil {
		return globalVMRunner(cf, args)
	}
	return evaluator.NULL
}
