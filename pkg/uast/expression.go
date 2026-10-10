package uast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type LiteralKind string

const (
	LiteralKindInt    LiteralKind = "int"
	LiteralKindFloat  LiteralKind = "float"
	LiteralKindString LiteralKind = "string"
	LiteralKindBool   LiteralKind = "bool"
	LiteralKindNull   LiteralKind = "null"
)

// Literal represents a constant value (e.g. 42, 3.14, "halo", true, null).
type Literal struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Kind     LiteralKind   `json:"kind"`
	Raw      string        `json:"raw"`
	Value    any           `json:"value"`
}

func (l *Literal) Type() NodeType          { return NodeTypeLiteral }
func (l *Literal) GetMetadata() *NodeMetadata { return l.Metadata }
func (l *Literal) expressionNode()         {}

func (l *Literal) UnmarshalJSON(data []byte) error {
	type rawLit Literal
	var raw rawLit
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*l = Literal(raw)
	l.NodeType = NodeTypeLiteral

	// Parse Value strongly based on Kind and Raw
	switch l.Kind {
	case LiteralKindInt:
		if l.Raw != "" {
			if v, err := strconv.ParseInt(l.Raw, 10, 64); err == nil {
				l.Value = v
				return nil
			}
		}
		if f, ok := l.Value.(float64); ok {
			l.Value = int64(f)
		}
	case LiteralKindFloat:
		if l.Raw != "" {
			if v, err := strconv.ParseFloat(l.Raw, 64); err == nil {
				l.Value = v
				return nil
			}
		}
		if f, ok := l.Value.(float64); ok {
			l.Value = f
		}
	case LiteralKindBool:
		if l.Raw != "" {
			if v, err := strconv.ParseBool(l.Raw); err == nil {
				l.Value = v
				return nil
			}
		}
	case LiteralKindNull:
		l.Value = nil
	case LiteralKindString:
		if s, ok := l.Value.(string); ok {
			l.Value = s
		}
	}
	return nil
}

// Identifier represents a named reference to a variable, function, or binding.
type Identifier struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Name     string        `json:"name"`
}

func (i *Identifier) Type() NodeType          { return NodeTypeIdentifier }
func (i *Identifier) GetMetadata() *NodeMetadata { return i.Metadata }
func (i *Identifier) expressionNode()         {}

// BinaryExpression represents binary operations (e.g., a + b, x == y, count > 5).
type BinaryExpression struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Operator string        `json:"operator"`
	Left     Expression    `json:"left"`
	Right    Expression    `json:"right"`
}

func (b *BinaryExpression) Type() NodeType          { return NodeTypeBinaryExpression }
func (b *BinaryExpression) GetMetadata() *NodeMetadata { return b.Metadata }
func (b *BinaryExpression) expressionNode()         {}

func (b *BinaryExpression) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Operator string          `json:"operator"`
		Left     json.RawMessage `json:"left"`
		Right    json.RawMessage `json:"right"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	leftExpr, err := UnmarshalExpressionJSON(raw.Left)
	if err != nil {
		return fmt.Errorf("BinaryExpression.left unmarshal error: %w", err)
	}
	rightExpr, err := UnmarshalExpressionJSON(raw.Right)
	if err != nil {
		return fmt.Errorf("BinaryExpression.right unmarshal error: %w", err)
	}
	b.NodeType = NodeTypeBinaryExpression
	b.Metadata = raw.Metadata
	b.Operator = raw.Operator
	b.Left = leftExpr
	b.Right = rightExpr
	return nil
}

// UnaryExpression represents unary operations (e.g., -x, !active, not done).
type UnaryExpression struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Operator string        `json:"operator"`
	Argument Expression    `json:"argument"`
	IsPrefix bool          `json:"is_prefix"`
}

func (u *UnaryExpression) Type() NodeType          { return NodeTypeUnaryExpression }
func (u *UnaryExpression) GetMetadata() *NodeMetadata { return u.Metadata }
func (u *UnaryExpression) expressionNode()         {}

func (u *UnaryExpression) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Operator string          `json:"operator"`
		Argument json.RawMessage `json:"argument"`
		IsPrefix bool            `json:"is_prefix"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	argExpr, err := UnmarshalExpressionJSON(raw.Argument)
	if err != nil {
		return fmt.Errorf("UnaryExpression.argument unmarshal error: %w", err)
	}
	u.NodeType = NodeTypeUnaryExpression
	u.Metadata = raw.Metadata
	u.Operator = raw.Operator
	u.Argument = argExpr
	u.IsPrefix = raw.IsPrefix
	return nil
}

// CallExpression represents invoking a function/routine (e.g. cetak("Halo"), hitung(1, 2)).
type CallExpression struct {
	NodeType  NodeType      `json:"type"`
	Metadata  *NodeMetadata `json:"metadata,omitempty"`
	Callee    Expression    `json:"callee"`
	Arguments []Expression  `json:"arguments"`
}

func (c *CallExpression) Type() NodeType          { return NodeTypeCallExpression }
func (c *CallExpression) GetMetadata() *NodeMetadata { return c.Metadata }
func (c *CallExpression) expressionNode()         {}

func (c *CallExpression) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType  NodeType          `json:"type"`
		Metadata  *NodeMetadata     `json:"metadata,omitempty"`
		Callee    json.RawMessage   `json:"callee"`
		Arguments []json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	calleeExpr, err := UnmarshalExpressionJSON(raw.Callee)
	if err != nil {
		return fmt.Errorf("CallExpression.callee unmarshal error: %w", err)
	}
	args := make([]Expression, len(raw.Arguments))
	for i, argRaw := range raw.Arguments {
		expr, err := UnmarshalExpressionJSON(argRaw)
		if err != nil {
			return fmt.Errorf("CallExpression.arguments[%d] unmarshal error: %w", i, err)
		}
		args[i] = expr
	}
	c.NodeType = NodeTypeCallExpression
	c.Metadata = raw.Metadata
	c.Callee = calleeExpr
	c.Arguments = args
	return nil
}

// MemberExpression represents accessing properties (e.g., x.sebelumnya, user.nama, data[0]).
// This is essential for features like Time-Travel Variables (.sebelumnya, .awal).
type MemberExpression struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Object   Expression    `json:"object"`
	Property string        `json:"property"`
	Computed bool          `json:"computed"` // true if obj[prop]
}

func (m *MemberExpression) Type() NodeType          { return NodeTypeMemberExpression }
func (m *MemberExpression) GetMetadata() *NodeMetadata { return m.Metadata }
func (m *MemberExpression) expressionNode()         {}

func (m *MemberExpression) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Object   json.RawMessage `json:"object"`
		Property string          `json:"property"`
		Computed bool            `json:"computed"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	objExpr, err := UnmarshalExpressionJSON(raw.Object)
	if err != nil {
		return fmt.Errorf("MemberExpression.object unmarshal error: %w", err)
	}
	m.NodeType = NodeTypeMemberExpression
	m.Metadata = raw.Metadata
	m.Object = objExpr
	m.Property = raw.Property
	m.Computed = raw.Computed
	return nil
}

// ArrayLiteral represents lists or arrays of expressions (e.g., [1, 2, 3]).
type ArrayLiteral struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Elements []Expression  `json:"elements"`
}

func (a *ArrayLiteral) Type() NodeType          { return NodeTypeArrayLiteral }
func (a *ArrayLiteral) GetMetadata() *NodeMetadata { return a.Metadata }
func (a *ArrayLiteral) expressionNode()         {}

func (a *ArrayLiteral) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType          `json:"type"`
		Metadata *NodeMetadata     `json:"metadata,omitempty"`
		Elements []json.RawMessage `json:"elements"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	elems := make([]Expression, len(raw.Elements))
	for i, elemRaw := range raw.Elements {
		expr, err := UnmarshalExpressionJSON(elemRaw)
		if err != nil {
			return fmt.Errorf("ArrayLiteral.elements[%d] unmarshal error: %w", i, err)
		}
		elems[i] = expr
	}
	a.NodeType = NodeTypeArrayLiteral
	a.Metadata = raw.Metadata
	a.Elements = elems
	return nil
}

// MapProperty represents a key-value entry in a Map literal.
type MapProperty struct {
	Key   Expression `json:"key"`
	Value Expression `json:"value"`
}

func (p *MapProperty) UnmarshalJSON(data []byte) error {
	var raw struct {
		Key   json.RawMessage `json:"key"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	keyExpr, err := UnmarshalExpressionJSON(raw.Key)
	if err != nil {
		return fmt.Errorf("MapProperty.key unmarshal error: %w", err)
	}
	valExpr, err := UnmarshalExpressionJSON(raw.Value)
	if err != nil {
		return fmt.Errorf("MapProperty.value unmarshal error: %w", err)
	}
	p.Key = keyExpr
	p.Value = valExpr
	return nil
}

// MapLiteral represents dictionary/hash map structures (e.g., {"kunci": "nilai", "poin" itu 100}).
type MapLiteral struct {
	NodeType   NodeType      `json:"type"`
	Metadata   *NodeMetadata `json:"metadata,omitempty"`
	Properties []MapProperty `json:"properties"`
}

func (m *MapLiteral) Type() NodeType          { return NodeTypeMapLiteral }
func (m *MapLiteral) GetMetadata() *NodeMetadata { return m.Metadata }
func (m *MapLiteral) expressionNode()         {}

func (m *MapLiteral) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType   NodeType      `json:"type"`
		Metadata   *NodeMetadata `json:"metadata,omitempty"`
		Properties []MapProperty `json:"properties"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw.Properties == nil {
		raw.Properties = []MapProperty{}
	}
	m.NodeType = NodeTypeMapLiteral
	m.Metadata = raw.Metadata
	m.Properties = raw.Properties
	return nil
}

// PipelineExpression represents piping a value into a function (left |> right).
type PipelineExpression struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Left     Expression    `json:"left"`
	Right    Expression    `json:"right"`
}

func (p *PipelineExpression) Type() NodeType          { return NodeTypePipelineExpression }
func (p *PipelineExpression) GetMetadata() *NodeMetadata { return p.Metadata }
func (p *PipelineExpression) expressionNode()         {}

func (p *PipelineExpression) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Left     json.RawMessage `json:"left"`
		Right    json.RawMessage `json:"right"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	leftExpr, err := UnmarshalExpressionJSON(raw.Left)
	if err != nil {
		return fmt.Errorf("PipelineExpression.left unmarshal error: %w", err)
	}
	rightExpr, err := UnmarshalExpressionJSON(raw.Right)
	if err != nil {
		return fmt.Errorf("PipelineExpression.right unmarshal error: %w", err)
	}
	p.NodeType = NodeTypePipelineExpression
	p.Metadata = raw.Metadata
	p.Left = leftExpr
	p.Right = rightExpr
	return nil
}

// FunctionExpression represents an anonymous function literal (e.g., fungsi(x) { kembalikan x * 2 }).
type FunctionExpression struct {
	NodeType   NodeType        `json:"type"`
	Metadata   *NodeMetadata   `json:"metadata,omitempty"`
	Parameters []Parameter     `json:"parameters"`
	ReturnType string          `json:"return_type,omitempty"`
	Body       *BlockStatement `json:"body"`
}

func (f *FunctionExpression) Type() NodeType          { return NodeTypeFunctionExpression }
func (f *FunctionExpression) GetMetadata() *NodeMetadata { return f.Metadata }
func (f *FunctionExpression) expressionNode()         {}

func (f *FunctionExpression) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType   NodeType        `json:"type"`
		Metadata   *NodeMetadata   `json:"metadata,omitempty"`
		Parameters []Parameter     `json:"parameters"`
		ReturnType string          `json:"return_type,omitempty"`
		Body       *BlockStatement `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	f.NodeType = NodeTypeFunctionExpression
	f.Metadata = raw.Metadata
	f.Parameters = raw.Parameters
	f.ReturnType = raw.ReturnType
	f.Body = raw.Body
	return nil
}

// TemplateLiteral represents an interpolated string composed of literal string parts and expressions.
type TemplateLiteral struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Parts    []Expression  `json:"parts"`
}

func (t *TemplateLiteral) Type() NodeType          { return NodeTypeTemplateLiteral }
func (t *TemplateLiteral) GetMetadata() *NodeMetadata { return t.Metadata }
func (t *TemplateLiteral) expressionNode()         {}

func (t *TemplateLiteral) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType          `json:"type"`
		Metadata *NodeMetadata     `json:"metadata,omitempty"`
		Parts    []json.RawMessage `json:"parts"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.NodeType = NodeTypeTemplateLiteral
	t.Metadata = raw.Metadata
	t.Parts = make([]Expression, len(raw.Parts))
	for i, rawPart := range raw.Parts {
		expr, err := UnmarshalExpressionJSON(rawPart)
		if err != nil {
			return fmt.Errorf("TemplateLiteral part %d unmarshal error: %w", i, err)
		}
		t.Parts[i] = expr
	}
	return nil
}


