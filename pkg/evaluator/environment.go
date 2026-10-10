package evaluator

import (
	"fmt"
	"sync"
)

// VariableEntry stores the state and timeline history for a variable.
// This fulfills the "Time-Travel Variables" architecture.
type VariableEntry struct {
	Name      string
	History   []Object
	IsMutable bool
	TypeHint  string
}

// Current returns the most recent value of the variable.
func (ve *VariableEntry) Current() Object {
	if len(ve.History) == 0 {
		return &Null{}
	}
	return ve.History[len(ve.History)-1]
}

// Previous returns the value right before the current one (if available).
func (ve *VariableEntry) Previous() (Object, bool) {
	if len(ve.History) < 2 {
		return nil, false
	}
	return ve.History[len(ve.History)-2], true
}

// Initial returns the first value ever assigned to this variable.
func (ve *VariableEntry) Initial() Object {
	if len(ve.History) == 0 {
		return &Null{}
	}
	return ve.History[0]
}

// Environment manages variable bindings across lexical scopes with thread-safety for goroutines.
type Environment struct {
	mu    sync.RWMutex
	store map[string]*VariableEntry
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]*VariableEntry),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Define declares a new variable in the current scope with an initial value and time-travel history.
func (e *Environment) Define(name string, val Object, isMutable bool, typeHint string) (*VariableEntry, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := &VariableEntry{
		Name:      name,
		History:   []Object{val},
		IsMutable: isMutable,
		TypeHint:  typeHint,
	}
	e.store[name] = entry
	return entry, nil
}

// Assign updates an existing variable and records the change into its timeline history.
func (e *Environment) Assign(name string, val Object) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry, ok := e.store[name]
	if !ok {
		if e.outer != nil {
			return e.outer.Assign(name, val)
		}
		return fmt.Errorf("variabel '%s' belum didefinisikan", name)
	}

	if !entry.IsMutable {
		return fmt.Errorf("variabel '%s' bersifat tetap (immutable) dan tidak dapat diubah", name)
	}

	entry.History = append(entry.History, val)
	return nil
}

// Get retrieves the current value of a variable.
func (e *Environment) Get(name string) (Object, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	entry, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	if !ok {
		return nil, false
	}
	return entry.Current(), true
}

// GetEntry retrieves the VariableEntry to inspect its time-travel history.
func (e *Environment) GetEntry(name string) (*VariableEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	entry, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.GetEntry(name)
	}
	return entry, ok
}
