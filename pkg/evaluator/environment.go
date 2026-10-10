package evaluator

import (
	"fmt"
	"sync"

	"celang/pkg/uast"
)

const DefaultHistoryCapacity = 50

// VariableEntry stores the state and timeline history for a variable.
// Uses an Anchor Snapshot + Bounded Ring Buffer architecture to guarantee O(1) memory usage
// even in long-running loops with millions of mutations.
type VariableEntry struct {
	Name        string
	InitialVal  Object
	PreviousVal Object
	CurrentVal  Object
	History     []Object
	HistoryCap  int
	IsMutable   bool
	TypeHint    string
	Visibility  *uast.Visibility
}

// Current returns the most recent value of the variable.
func (ve *VariableEntry) Current() Object {
	if ve.CurrentVal != nil {
		return ve.CurrentVal
	}
	if len(ve.History) > 0 {
		return ve.History[len(ve.History)-1]
	}
	return &Null{}
}

// Previous returns the value right before the current one (if available).
func (ve *VariableEntry) Previous() (Object, bool) {
	if ve.PreviousVal != nil {
		return ve.PreviousVal, true
	}
	if len(ve.History) >= 2 {
		return ve.History[len(ve.History)-2], true
	}
	return nil, false
}

// Initial returns the first value ever assigned to this variable.
func (ve *VariableEntry) Initial() Object {
	if ve.InitialVal != nil {
		return ve.InitialVal
	}
	if len(ve.History) > 0 {
		return ve.History[0]
	}
	return &Null{}
}

// ModuleRegistry caches evaluated modules across imports to prevent duplicate execution and circular loops.
type ModuleRegistry struct {
	mu      sync.RWMutex
	modules map[string]*Environment
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]*Environment),
	}
}

// Environment manages variable bindings across lexical scopes with thread-safety for goroutines.
type Environment struct {
	mu             sync.RWMutex
	store          map[string]*VariableEntry
	outer          *Environment
	FilePath       string
	ModuleRoot     string
	ModuleRegistry *ModuleRegistry
}

func NewEnvironment() *Environment {
	return &Environment{
		store:          make(map[string]*VariableEntry),
		ModuleRegistry: NewModuleRegistry(),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	if outer != nil {
		env.FilePath = outer.FilePath
		env.ModuleRoot = outer.ModuleRoot
		env.ModuleRegistry = outer.ModuleRegistry
	}
	return env
}

// Define declares a new variable in the current scope with an initial value and time-travel history.
// By default, visibility is private (VisibilityPrivate / privat!).
func (e *Environment) Define(name string, val Object, isMutable bool, typeHint string) (*VariableEntry, error) {
	return e.DefineWithVisibility(name, val, isMutable, typeHint, &uast.Visibility{Kind: uast.VisibilityPrivate})
}

// DefineWithVisibility declares a variable with an explicit visibility scope.
func (e *Environment) DefineWithVisibility(name string, val Object, isMutable bool, typeHint string, vis *uast.Visibility) (*VariableEntry, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if vis == nil {
		vis = &uast.Visibility{Kind: uast.VisibilityPrivate}
	}

	entry := &VariableEntry{
		Name:        name,
		InitialVal:  val,
		CurrentVal:  val,
		History:     []Object{val},
		HistoryCap:  DefaultHistoryCapacity,
		IsMutable:   isMutable,
		TypeHint:    typeHint,
		Visibility:  vis,
	}
	e.store[name] = entry
	return entry, nil
}

// SetVisibility updates the visibility of an existing symbol in this scope.
func (e *Environment) SetVisibility(name string, vis *uast.Visibility) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry, ok := e.store[name]
	if !ok {
		return fmt.Errorf("simbol '%s' belum didefinisikan", name)
	}
	entry.Visibility = vis
	return nil
}

// AllEntries returns a snapshot of all variables and their metadata in this scope.
func (e *Environment) AllEntries() map[string]*VariableEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()

	snapshot := make(map[string]*VariableEntry, len(e.store))
	for k, v := range e.store {
		snapshot[k] = v
	}
	return snapshot
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

	entry.PreviousVal = entry.CurrentVal
	entry.CurrentVal = val

	capLimit := entry.HistoryCap
	if capLimit <= 0 {
		capLimit = DefaultHistoryCapacity
	}

	if len(entry.History) >= capLimit {
		// Zero-allocation sliding window: drop oldest, preserve last 'capLimit' entries
		copy(entry.History, entry.History[1:])
		entry.History[len(entry.History)-1] = val
	} else {
		entry.History = append(entry.History, val)
	}

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
