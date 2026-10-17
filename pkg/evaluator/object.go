package evaluator

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"celang/pkg/uast"
)

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	FLOAT_OBJ        ObjectType = "FLOAT"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	STRING_OBJ       ObjectType = "STRING"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	MAP_OBJ          ObjectType = "MAP"
	POINTER_OBJ         ObjectType = "POINTER"
	STRUCT_DEF_OBJ      ObjectType = "STRUCT_DEF"
	STRUCT_INSTANCE_OBJ ObjectType = "STRUCT_INSTANCE"
	TICKET_OBJ          ObjectType = "TIKET"
	CHANNEL_OBJ         ObjectType = "SALURAN"
	ITERATOR_OBJ        ObjectType = "ITERATOR"
	ROLLBACK_SIGNAL_OBJ ObjectType = "ROLLBACK_SIGNAL"
	BREAK_SIGNAL_OBJ    ObjectType = "BREAK_SIGNAL"
	CONTINUE_SIGNAL_OBJ ObjectType = "CONTINUE_SIGNAL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Callable represents any object that can be called with arguments (e.g. bytecode CompiledFunction).
type Callable interface {
	Object
	Call(args ...Object) Object
}

// ASTNodeWrapper wraps a U-AST node as an Object for bytecode chunks.
type ASTNodeWrapper struct {
	Node uast.Node
}

func (w *ASTNodeWrapper) Type() ObjectType { return "AST_NODE" }
func (w *ASTNodeWrapper) Inspect() string   { return "<ast-node>" }

// Integer represents a 64-bit signed integer.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

// Float represents a 64-bit IEEE 754 floating point number.
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// Boolean represents true or false.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

// String represents text.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Null represents absence of value.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue wraps a value returning from a function call or block.
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

// Error represents an empathetic runtime error with helpful suggestions.
type Error struct {
	Message      string
	Pos          *uast.Position
	FriendlyHint string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Kesalahan: %s", e.Message))
	if e.Pos != nil {
		sb.WriteString(fmt.Sprintf(" [baris %d, kolom %d]", e.Pos.Line, e.Pos.Column))
	}
	if e.FriendlyHint != "" {
		sb.WriteString(fmt.Sprintf("\n💡 Saran: %s", e.FriendlyHint))
	}
	return sb.String()
}

// Function represents a user-defined function.
type Function struct {
	Name           string
	TypeParameters []uast.TypeParameter
	Parameters     []uast.Parameter
	Body           *uast.BlockStatement
	Env            *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.Name)
	}
	tpStr := ""
	if len(f.TypeParameters) > 0 {
		var tpStrs []string
		for _, tp := range f.TypeParameters {
			if tp.Constraint != "" {
				tpStrs = append(tpStrs, fmt.Sprintf("%s: %s", tp.Name, tp.Constraint))
			} else {
				tpStrs = append(tpStrs, tp.Name)
			}
		}
		tpStr = fmt.Sprintf("<%s>", strings.Join(tpStrs, ", "))
	}
	return fmt.Sprintf("fungsi %s%s(%s)", f.Name, tpStr, strings.Join(params, ", "))
}

// BuiltinFunction is a native Go function exposed to Celang.
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "fungsi bawaan (builtin)" }

// Array represents an ordered collection of objects.
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	return inspectObjectWithVisited(a, make(map[any]bool))
}

// Map represents a key-value dictionary.
type Map struct {
	Pairs map[string]Object
	Order []string // maintains key order for readable inspection
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	return inspectObjectWithVisited(m, make(map[any]bool))
}

// GetKeys returns the ordered keys of the Map, automatically reconciling with Pairs if needed.
func (m *Map) GetKeys() []string {
	if len(m.Order) != len(m.Pairs) {
		orderSet := make(map[string]bool, len(m.Order))
		var newOrder []string
		for _, k := range m.Order {
			if _, exists := m.Pairs[k]; exists && !orderSet[k] {
				orderSet[k] = true
				newOrder = append(newOrder, k)
			}
		}
		var remaining []string
		for k := range m.Pairs {
			if !orderSet[k] {
				remaining = append(remaining, k)
			}
		}
		sort.Strings(remaining)
		m.Order = append(newOrder, remaining...)
	}
	return m.Order
}

// Pointer represents a memory reference to a variable in an Environment.
type Pointer struct {
	Env     *Environment
	VarName string
}

func (p *Pointer) Type() ObjectType { return POINTER_OBJ }
func (p *Pointer) Inspect() string {
	valStr := "null"
	if val, ok := p.Env.Get(p.VarName); ok {
		valStr = val.Inspect()
	}
	return fmt.Sprintf("&%s (-> %s)", p.VarName, valStr)
}

// StructDefinition represents the blueprint of a custom data type (bentuk).
type StructDefinition struct {
	Name           string
	TypeParameters []uast.TypeParameter
	Fields         []uast.StructField
	Methods        map[string]*Function
	PrivateFields  map[string]bool
	PrivateMethods map[string]bool
	Embeds         []string
	Env            *Environment
}

func (s *StructDefinition) Type() ObjectType { return STRUCT_DEF_OBJ }
func (s *StructDefinition) Inspect() string {
	if len(s.TypeParameters) > 0 {
		var tpStrs []string
		for _, tp := range s.TypeParameters {
			if tp.Constraint != "" {
				tpStrs = append(tpStrs, fmt.Sprintf("%s: %s", tp.Name, tp.Constraint))
			} else {
				tpStrs = append(tpStrs, tp.Name)
			}
		}
		return fmt.Sprintf("bentuk %s<%s>", s.Name, strings.Join(tpStrs, ", "))
	}
	return fmt.Sprintf("bentuk %s", s.Name)
}

// StructInstance represents an instantiated struct object with fields and methods.
type StructInstance struct {
	Def           *StructDefinition
	TypeArguments []string
	Fields        map[string]Object
	Env           *Environment
}

func (s *StructInstance) Type() ObjectType { return STRUCT_INSTANCE_OBJ }
func (s *StructInstance) Inspect() string {
	return inspectObjectWithVisited(s, make(map[any]bool))
}

func inspectObjectWithVisited(obj Object, visited map[any]bool) string {
	if obj == nil {
		return "null"
	}
	switch o := obj.(type) {
	case *StructInstance:
		name := o.Def.Name
		if len(o.TypeArguments) > 0 {
			name += fmt.Sprintf("<%s>", strings.Join(o.TypeArguments, ", "))
		}
		if visited[o] {
			return fmt.Sprintf("%s{...}", name)
		}
		visited[o] = true
		defer delete(visited, o)
		var parts []string
		for _, f := range o.Def.Fields {
			val := o.Fields[f.Name]
			if val == nil {
				val = &Null{}
			}
			parts = append(parts, fmt.Sprintf("%s: %s", f.Name, inspectObjectWithVisited(val, visited)))
		}
		return fmt.Sprintf("%s{%s}", name, strings.Join(parts, ", "))

	case *Array:
		if visited[o] {
			return "[...]"
		}
		visited[o] = true
		defer delete(visited, o)
		var elements []string
		for _, el := range o.Elements {
			elements = append(elements, inspectObjectWithVisited(el, visited))
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case *Map:
		if visited[o] {
			return "{...}"
		}
		visited[o] = true
		defer delete(visited, o)
		var pairs []string
		for _, k := range o.GetKeys() {
			v := o.Pairs[k]
			pairs = append(pairs, fmt.Sprintf("%q: %s", k, inspectObjectWithVisited(v, visited)))
		}
		return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))

	default:
		return o.Inspect()
	}
}

// IteratorObject wraps stateful iteration for arrays, strings, maps, and channels.
type IteratorObject struct {
	NextFunc func() (Object, bool)
}

func (it *IteratorObject) Type() ObjectType { return ITERATOR_OBJ }
func (it *IteratorObject) Inspect() string   { return "<iterator>" }
func (it *IteratorObject) Next() (Object, bool) {
	if it.NextFunc == nil {
		return nil, false
	}
	return it.NextFunc()
}

// NewIterator creates an iterator for any iterable Celang object.
func NewIterator(obj Object) (*IteratorObject, error) {
	switch o := obj.(type) {
	case *Array:
		idx := 0
		return &IteratorObject{
			NextFunc: func() (Object, bool) {
				if idx >= len(o.Elements) {
					return nil, false
				}
				val := o.Elements[idx]
				idx++
				return val, true
			},
		}, nil
	case *String:
		runes := []rune(o.Value)
		idx := 0
		return &IteratorObject{
			NextFunc: func() (Object, bool) {
				if idx >= len(runes) {
					return nil, false
				}
				val := &String{Value: string(runes[idx])}
				idx++
				return val, true
			},
		}, nil
	case *Map:
		idx := 0
		keys := o.GetKeys()
		return &IteratorObject{
			NextFunc: func() (Object, bool) {
				if idx >= len(keys) {
					return nil, false
				}
				val := &String{Value: keys[idx]}
				idx++
				return val, true
			},
		}, nil
	case *Channel:
		return &IteratorObject{
			NextFunc: func() (Object, bool) {
				val, ok := <-o.Ch
				if !ok {
					return nil, false
				}
				return val, true
			},
		}, nil
	default:
		return nil, fmt.Errorf("objek tipe %s tidak dapat diiterasi", obj.Type())
	}
}

// RollbackSignal is generated by 'batalkan' to unwind the execution to the enclosing transaction.
type RollbackSignal struct {
	Reason Object
}

func (r *RollbackSignal) Type() ObjectType { return ROLLBACK_SIGNAL_OBJ }
func (r *RollbackSignal) Inspect() string {
	if r.Reason != nil {
		return r.Reason.Inspect()
	}
	return "transaksi dibatalkan"
}

// BreakSignal is generated by 'berhenti' / 'break' to exit loops.
type BreakSignal struct{}

func (b *BreakSignal) Type() ObjectType { return BREAK_SIGNAL_OBJ }
func (b *BreakSignal) Inspect() string  { return "berhenti" }

// ContinueSignal is generated by 'lanjut' / 'continue' to skip to next loop iteration.
type ContinueSignal struct{}

func (c *ContinueSignal) Type() ObjectType { return CONTINUE_SIGNAL_OBJ }
func (c *ContinueSignal) Inspect() string  { return "lanjut" }

var (
	BREAK_SIGNAL    = &BreakSignal{}
	CONTINUE_SIGNAL = &ContinueSignal{}
)


