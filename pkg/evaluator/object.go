package evaluator

import (
	"fmt"
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
	POINTER_OBJ      ObjectType = "POINTER"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

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
	Name       string
	Parameters []uast.Parameter
	Body       *uast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.Name)
	}
	return fmt.Sprintf("fungsi %s(%s)", f.Name, strings.Join(params, ", "))
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
	var elements []string
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// Map represents a key-value dictionary.
type Map struct {
	Pairs map[string]Object
	Order []string // maintains key order for readable inspection
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var pairs []string
	for _, k := range m.Order {
		v := m.Pairs[k]
		pairs = append(pairs, fmt.Sprintf("%q: %s", k, v.Inspect()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
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
