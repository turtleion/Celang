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
	Name         string
	InitialVal   Object
	PreviousVal  Object
	CurrentVal   Object
	History      []Object
	HistoryCap   int
	IsMutable    bool
	TypeHint     string
	Visibility   *uast.Visibility
	AliasedEntry *VariableEntry
}

// Resolve unwraps any reference chain and returns the root VariableEntry.
func (ve *VariableEntry) Resolve() *VariableEntry {
	curr := ve
	for curr.AliasedEntry != nil {
		curr = curr.AliasedEntry
	}
	return curr
}

// Current returns the most recent value of the variable.
func (ve *VariableEntry) Current() Object {
	if ve.AliasedEntry != nil {
		return ve.AliasedEntry.Current()
	}
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
	if ve.AliasedEntry != nil {
		return ve.AliasedEntry.Previous()
	}
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
	if ve.AliasedEntry != nil {
		return ve.AliasedEntry.Initial()
	}
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
	typeAliases    map[string]string
	outer          *Environment
	FilePath       string
	ModuleRoot     string
	ModuleRegistry *ModuleRegistry
	defers         []uast.Statement
}

func NewEnvironment() *Environment {
	return &Environment{
		store:          make(map[string]*VariableEntry),
		typeAliases:    make(map[string]string),
		ModuleRegistry: NewModuleRegistry(),
		defers:         make([]uast.Statement, 0),
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

// PushDefer adds a statement to be executed at scope exit (LIFO).
func (e *Environment) PushDefer(stmt uast.Statement) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.defers = append(e.defers, stmt)
}

// PopAllDefers extracts all deferred statements in LIFO (reverse) order and clears the stack.
func (e *Environment) PopAllDefers() []uast.Statement {
	e.mu.Lock()
	defer e.mu.Unlock()
	n := len(e.defers)
	if n == 0 {
		return nil
	}
	reversed := make([]uast.Statement, n)
	for i := 0; i < n; i++ {
		reversed[i] = e.defers[n-1-i]
	}
	e.defers = nil
	return reversed
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

// DefineTypeAlias registers a custom type alias in the current environment.
func (e *Environment) DefineTypeAlias(name, baseType string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.typeAliases == nil {
		e.typeAliases = make(map[string]string)
	}
	e.typeAliases[name] = baseType
}

// GetTypeAlias retrieves a type alias if defined in this or outer environments.
func (e *Environment) GetTypeAlias(name string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.typeAliases != nil {
		if base, ok := e.typeAliases[name]; ok {
			return base, true
		}
	}
	if e.outer != nil {
		return e.outer.GetTypeAlias(name)
	}
	return "", false
}

// DefineReference creates a transparent reference alias pointing to a target VariableEntry.
func (e *Environment) DefineReference(aliasName string, target *VariableEntry) (*VariableEntry, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if target == nil {
		return nil, fmt.Errorf("target rujuk tidak boleh nil")
	}

	target = target.Resolve()
	entry := &VariableEntry{
		Name:         aliasName,
		AliasedEntry: target,
		IsMutable:    true,
	}
	e.store[aliasName] = entry
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

	target := entry.Resolve()
	if !target.IsMutable {
		return fmt.Errorf("variabel '%s' bersifat tetap (immutable) dan tidak dapat diubah", target.Name)
	}

	target.PreviousVal = target.CurrentVal
	target.CurrentVal = val

	capLimit := target.HistoryCap
	if capLimit <= 0 {
		capLimit = DefaultHistoryCapacity
	}

	if len(target.History) >= capLimit {
		// Zero-allocation sliding window: drop oldest, preserve last 'capLimit' entries
		copy(target.History, target.History[1:])
		target.History[len(target.History)-1] = val
	} else {
		target.History = append(target.History, val)
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
	return entry.Resolve().Current(), true
}

// GetEntry retrieves the VariableEntry to inspect its time-travel history.
func (e *Environment) GetEntry(name string) (*VariableEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	entry, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.GetEntry(name)
	}
	if !ok {
		return nil, false
	}
	return entry.Resolve(), true
}
