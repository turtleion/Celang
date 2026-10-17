package evaluator

import (
	"fmt"
	"sync"
	"sync/atomic"

	"celang/pkg/uast"
)

const DefaultHistoryCapacity = 50

// VariableEntry stores the state and timeline history for a variable.
// Uses an Anchor Snapshot + Bounded Ring Buffer architecture to guarantee O(1) memory usage
// even in long-running loops with millions of mutations.
// ReactiveWatcher represents an observer listening to changes on a variable (pantau).
type ReactiveWatcher struct {
	TargetName string
	NewValVar  string
	OldValVar  string
	Body       *uast.BlockStatement
	Env        *Environment
}

type VariableEntry struct {
	Name         string
	InitialVal   Object
	PreviousVal  Object
	CurrentVal   Object
	History      []Object
	HistoryCap   int
	IsMutable    bool
	IsConstant   bool
	TypeHint     string
	Visibility   *uast.Visibility
	AliasedEntry *VariableEntry

	// Reactive & Observer extensions
	IsReactive   bool
	ReactiveExpr uast.Expression
	ReactiveEnv  *Environment
	Dependents   []*VariableEntry
	Watchers     []*ReactiveWatcher
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
	store           map[string]*VariableEntry
	typeAliases     map[string]string
	outer           *Environment
	FilePath        string
	ModuleRoot      string
	ModuleRegistry  *ModuleRegistry
	defers          []uast.Statement
	savepoints      []*transactionFrame
	reactivityDepth int64
	activeTracker   *VariableEntry
	IsUniversal     bool
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
		env.IsUniversal = outer.IsUniversal
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

	// Track newly defined variable in active transaction frame
	curr := e
	for curr != nil {
		if curr != e {
			curr.mu.Lock()
		}
		hasSavepoints := len(curr.savepoints) > 0
		if hasSavepoints {
			top := curr.savepoints[len(curr.savepoints)-1]
			top.created = append(top.created, name)
		}
		next := curr.outer
		if curr != e {
			curr.mu.Unlock()
		}
		if hasSavepoints {
			break
		}
		curr = next
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

func (e *Environment) getRoot() *Environment {
	curr := e
	for curr.outer != nil {
		curr = curr.outer
	}
	return curr
}

func (e *Environment) getReactivityDepth() int {
	root := e.getRoot()
	return int(atomic.LoadInt64(&root.reactivityDepth))
}

func (e *Environment) incReactivityDepth() {
	root := e.getRoot()
	atomic.AddInt64(&root.reactivityDepth, 1)
}

func (e *Environment) decReactivityDepth() {
	root := e.getRoot()
	atomic.AddInt64(&root.reactivityDepth, -1)
}

// Assign updates an existing variable and records the change into its timeline history.
func (e *Environment) Assign(name string, val Object) error {
	e.mu.Lock()

	// Guard against infinite recursive watcher loops
	if e.getReactivityDepth() >= 50 {
		e.mu.Unlock()
		return fmt.Errorf("terdeteksi rekursi pengamat melebihi batas (maksimum 50 siklus pengamat)")
	}

	entry, ok := e.store[name]
	if !ok {
		e.mu.Unlock()
		if e.outer != nil {
			return e.outer.Assign(name, val)
		}
		return fmt.Errorf("variabel '%s' belum didefinisikan", name)
	}

	target := entry.Resolve()
	if !target.IsMutable {
		e.mu.Unlock()
		return fmt.Errorf("variabel '%s' bersifat tetap (immutable) dan tidak dapat diubah", target.Name)
	}

	// Record savepoint for active transaction rollback
	e.recordSavepoint(target)

	oldVal := target.CurrentVal
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

	hasReactivity := len(target.Dependents) > 0 || len(target.Watchers) > 0
	e.mu.Unlock()

	// Propagate Reactivity and Watchers outside mutex lock
	if hasReactivity {
		if errObj := e.triggerReactivity(target, val, oldVal, make(map[*VariableEntry]bool)); errObj != nil && isError(errObj) {
			return fmt.Errorf("%s", errObj.Inspect())
		}
	}

	return nil
}

// Get retrieves the current value of a variable.
func (e *Environment) Get(name string) (Object, bool) {
	entry, ok := e.GetEntry(name)
	if !ok || entry == nil {
		return nil, false
	}
	return entry.Resolve().Current(), true
}

// GetEntry retrieves the VariableEntry to inspect its time-travel history or track reactive dependencies.
func (e *Environment) GetEntry(name string) (*VariableEntry, bool) {
	e.mu.RLock()
	entry, ok := e.store[name]
	tracker := e.getActiveTracker()
	e.mu.RUnlock()

	if ok {
		root := entry.Resolve()
		if tracker != nil && tracker != root {
			e.mu.Lock()
			found := false
			for _, dep := range root.Dependents {
				if dep == tracker {
					found = true
					break
				}
			}
			if !found {
				root.Dependents = append(root.Dependents, tracker)
			}
			e.mu.Unlock()
		}
		return root, true
	}

	if e.outer != nil {
		return e.outer.GetEntry(name)
	}
	return nil, false
}

// GetLocalEntry retrieves the VariableEntry defined in this scope only (does not traverse outer scopes).
func (e *Environment) GetLocalEntry(name string) (*VariableEntry, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	entry, ok := e.store[name]
	if !ok {
		return nil, false
	}
	return entry.Resolve(), true
}

// SetActiveTracker sets the reactive entry that is currently collecting dependencies.
func (e *Environment) SetActiveTracker(entry *VariableEntry) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.activeTracker = entry
}

func (e *Environment) getActiveTracker() *VariableEntry {
	curr := e
	for curr != nil {
		if curr.activeTracker != nil {
			return curr.activeTracker
		}
		curr = curr.outer
	}
	return nil
}

type transactionSavepoint struct {
	val        Object
	prevVal    Object
	historyLen int
}

type transactionFrame struct {
	mutations map[*VariableEntry]transactionSavepoint
	created   []string
}

func (e *Environment) recordSavepoint(target *VariableEntry) {
	curr := e
	for curr != nil {
		for i := range curr.savepoints {
			frame := curr.savepoints[i]
			if _, exists := frame.mutations[target]; !exists {
				frame.mutations[target] = transactionSavepoint{
					val:        deepClone(target.CurrentVal),
					prevVal:    deepClone(target.PreviousVal),
					historyLen: len(target.History),
				}
			}
		}
		curr = curr.outer
	}
}

// RecordSavepointForVar records a savepoint for an existing variable before its internal fields/elements are mutated.
func (e *Environment) RecordSavepointForVar(name string) {
	entry, ok := e.GetEntry(name)
	if ok && entry != nil {
		e.recordSavepoint(entry.Resolve())
	}
}

// BeginTransaction starts a new transactional frame where all variable mutations can be rolled back.
func (e *Environment) BeginTransaction() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.savepoints = append(e.savepoints, &transactionFrame{
		mutations: make(map[*VariableEntry]transactionSavepoint),
		created:   make([]string, 0),
	})
}

// CommitTransaction confirms all mutations within the active transaction frame.
func (e *Environment) CommitTransaction() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.savepoints) == 0 {
		return
	}
	e.savepoints = e.savepoints[:len(e.savepoints)-1]
}

// RollbackTransaction reverts all mutations made during the current transaction frame.
func (e *Environment) RollbackTransaction() {
	e.mu.Lock()
	if len(e.savepoints) == 0 {
		e.mu.Unlock()
		return
	}
	top := e.savepoints[len(e.savepoints)-1]
	e.savepoints = e.savepoints[:len(e.savepoints)-1]

	// 1. Remove variables declared in this transaction frame
	for _, createdName := range top.created {
		delete(e.store, createdName)
	}

	// 2. Restore mutated variables to their pre-transaction snapshot
	type rollItem struct {
		entry *VariableEntry
		newV  Object
		oldV  Object
	}
	var toTrigger []rollItem

	for entry, sp := range top.mutations {
		oldVal := entry.CurrentVal
		entry.CurrentVal = sp.val
		entry.PreviousVal = sp.prevVal
		if sp.historyLen <= len(entry.History) {
			entry.History = entry.History[:sp.historyLen]
		}
		if len(entry.Dependents) > 0 || len(entry.Watchers) > 0 {
			toTrigger = append(toTrigger, rollItem{entry: entry, newV: sp.val, oldV: oldVal})
		}
	}
	e.mu.Unlock()

	for _, item := range toTrigger {
		e.triggerReactivity(item.entry, item.newV, item.oldV, make(map[*VariableEntry]bool))
	}
}

// DefineReactive registers a reactive derived variable and tracks its dependencies.
func (e *Environment) DefineReactive(name string, expr uast.Expression, evalEnv *Environment) (*VariableEntry, error) {
	e.mu.Lock()
	entry := &VariableEntry{
		Name:         name,
		IsMutable:    true,
		IsReactive:   true,
		ReactiveExpr: expr,
		ReactiveEnv:  evalEnv,
		HistoryCap:   DefaultHistoryCapacity,
		Visibility:   &uast.Visibility{Kind: uast.VisibilityPrivate},
	}
	e.store[name] = entry
	e.mu.Unlock()

	evalEnv.SetActiveTracker(entry)
	initialVal := Eval(expr, evalEnv)
	evalEnv.SetActiveTracker(nil)

	if isError(initialVal) {
		return nil, fmt.Errorf("%s", initialVal.Inspect())
	}

	e.mu.Lock()
	entry.InitialVal = initialVal
	entry.CurrentVal = initialVal
	entry.History = []Object{initialVal}
	e.mu.Unlock()

	return entry, nil
}

func (e *Environment) triggerReactivity(target *VariableEntry, newVal, oldVal Object, visited map[*VariableEntry]bool) Object {
	if visited[target] {
		return nil
	}
	visited[target] = true

	e.incReactivityDepth()
	defer e.decReactivityDepth()

	var watchers []*ReactiveWatcher
	if len(target.Watchers) > 0 {
		watchers = append([]*ReactiveWatcher{}, target.Watchers...)
	}
	var dependents []*VariableEntry
	if len(target.Dependents) > 0 {
		dependents = append([]*VariableEntry{}, target.Dependents...)
	}

	// 1. Notify watchers (pantau)
	for _, w := range watchers {
		if w.Body != nil && w.Env != nil {
			wEnv := NewEnclosedEnvironment(w.Env)
			nVal := newVal
			if nVal == nil {
				nVal = NULL
			}
			oVal := oldVal
			if oVal == nil {
				oVal = NULL
			}
			wEnv.Define(w.NewValVar, nVal, false, "")
			wEnv.Define(w.OldValVar, oVal, false, "")
			res := Eval(w.Body, wEnv)
			if isError(res) {
				return res
			}
		}
	}

	// 2. Update dependent derived variables (reaktif)
	for _, dep := range dependents {
		if dep.ReactiveExpr != nil && dep.ReactiveEnv != nil {
			depOld := dep.CurrentVal
			dep.ReactiveEnv.SetActiveTracker(dep)
			depNew := Eval(dep.ReactiveExpr, dep.ReactiveEnv)
			dep.ReactiveEnv.SetActiveTracker(nil)
			if isError(depNew) {
				return depNew
			}
			dep.PreviousVal = depOld
			dep.CurrentVal = depNew

			capLimit := dep.HistoryCap
			if capLimit <= 0 {
				capLimit = DefaultHistoryCapacity
			}
			if len(dep.History) >= capLimit {
				copy(dep.History, dep.History[1:])
				dep.History[len(dep.History)-1] = depNew
			} else {
				dep.History = append(dep.History, depNew)
			}
			if errRes := e.triggerReactivity(dep, depNew, depOld, visited); errRes != nil && isError(errRes) {
				return errRes
			}
		}
	}
	return nil
}
