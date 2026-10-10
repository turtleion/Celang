package uast

import (
	"encoding/json"
	"fmt"
)

// ExpressionStatement encapsulates an expression evaluated for side effects (e.g., function calls).
type ExpressionStatement struct {
	NodeType   NodeType      `json:"type"`
	Metadata   *NodeMetadata `json:"metadata,omitempty"`
	Expression Expression    `json:"expression"`
}

func (e *ExpressionStatement) Type() NodeType          { return NodeTypeExpressionStatement }
func (e *ExpressionStatement) GetMetadata() *NodeMetadata { return e.Metadata }
func (e *ExpressionStatement) statementNode()          {}

func (e *ExpressionStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType   NodeType        `json:"type"`
		Metadata   *NodeMetadata   `json:"metadata,omitempty"`
		Expression json.RawMessage `json:"expression"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	expr, err := UnmarshalExpressionJSON(raw.Expression)
	if err != nil {
		return fmt.Errorf("ExpressionStatement.expression unmarshal error: %w", err)
	}
	e.NodeType = NodeTypeExpressionStatement
	e.Metadata = raw.Metadata
	e.Expression = expr
	return nil
}

// VariableDeclaration represents defining a variable (mutable or immutable).
type VariableDeclaration struct {
	NodeType   NodeType      `json:"type"`
	Metadata   *NodeMetadata `json:"metadata,omitempty"`
	Name       string        `json:"name"`
	Value      Expression    `json:"value,omitempty"`
	IsMutable  bool          `json:"is_mutable"`
	TypeHint   string        `json:"type_hint,omitempty"`
	Visibility *Visibility   `json:"visibility,omitempty"`
}

func (v *VariableDeclaration) Type() NodeType          { return NodeTypeVariableDeclaration }
func (v *VariableDeclaration) GetMetadata() *NodeMetadata { return v.Metadata }
func (v *VariableDeclaration) statementNode()          {}

func (v *VariableDeclaration) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType   NodeType        `json:"type"`
		Metadata   *NodeMetadata   `json:"metadata,omitempty"`
		Name       string          `json:"name"`
		Value      json.RawMessage `json:"value,omitempty"`
		IsMutable  bool            `json:"is_mutable"`
		TypeHint   string          `json:"type_hint,omitempty"`
		Visibility *Visibility     `json:"visibility,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var valExpr Expression
	if len(raw.Value) > 0 && string(raw.Value) != "null" {
		var err error
		valExpr, err = UnmarshalExpressionJSON(raw.Value)
		if err != nil {
			return fmt.Errorf("VariableDeclaration.value unmarshal error: %w", err)
		}
	}
	v.NodeType = NodeTypeVariableDeclaration
	v.Metadata = raw.Metadata
	v.Name = raw.Name
	v.Value = valExpr
	v.IsMutable = raw.IsMutable
	v.TypeHint = raw.TypeHint
	v.Visibility = raw.Visibility
	return nil
}

// AssignmentStatement represents assigning or updating a variable / member.
type AssignmentStatement struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Target   Expression    `json:"target"`
	Operator string        `json:"operator"` // "=", "+=", "-=", etc.
	Value    Expression    `json:"value"`
}

func (a *AssignmentStatement) Type() NodeType          { return NodeTypeAssignmentStatement }
func (a *AssignmentStatement) GetMetadata() *NodeMetadata { return a.Metadata }
func (a *AssignmentStatement) statementNode()          {}

func (a *AssignmentStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Target   json.RawMessage `json:"target"`
		Operator string          `json:"operator"`
		Value    json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	targetExpr, err := UnmarshalExpressionJSON(raw.Target)
	if err != nil {
		return fmt.Errorf("AssignmentStatement.target unmarshal error: %w", err)
	}
	valExpr, err := UnmarshalExpressionJSON(raw.Value)
	if err != nil {
		return fmt.Errorf("AssignmentStatement.value unmarshal error: %w", err)
	}
	a.NodeType = NodeTypeAssignmentStatement
	a.Metadata = raw.Metadata
	a.Target = targetExpr
	a.Operator = raw.Operator
	a.Value = valExpr
	return nil
}

// BlockStatement groups a sequence of statements within a scope.
type BlockStatement struct {
	NodeType   NodeType      `json:"type"`
	Metadata   *NodeMetadata `json:"metadata,omitempty"`
	Statements []Statement   `json:"statements"`
}

func (b *BlockStatement) Type() NodeType          { return NodeTypeBlockStatement }
func (b *BlockStatement) GetMetadata() *NodeMetadata { return b.Metadata }
func (b *BlockStatement) statementNode()          {}

func (b *BlockStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType   NodeType          `json:"type"`
		Metadata   *NodeMetadata     `json:"metadata,omitempty"`
		Statements []json.RawMessage `json:"statements"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	stmts := make([]Statement, len(raw.Statements))
	for i, stmtRaw := range raw.Statements {
		stmt, err := UnmarshalStatementJSON(stmtRaw)
		if err != nil {
			return fmt.Errorf("BlockStatement.statements[%d] unmarshal error: %w", i, err)
		}
		stmts[i] = stmt
	}
	b.NodeType = NodeTypeBlockStatement
	b.Metadata = raw.Metadata
	b.Statements = stmts
	return nil
}

// IfStatement represents branching logic.
type IfStatement struct {
	NodeType   NodeType        `json:"type"`
	Metadata   *NodeMetadata   `json:"metadata,omitempty"`
	Condition  Expression      `json:"condition"`
	Consequent *BlockStatement `json:"consequent"`
	Alternate  Statement       `json:"alternate,omitempty"` // BlockStatement or another IfStatement
}

func (i *IfStatement) Type() NodeType          { return NodeTypeIfStatement }
func (i *IfStatement) GetMetadata() *NodeMetadata { return i.Metadata }
func (i *IfStatement) statementNode()          {}

func (i *IfStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType   NodeType        `json:"type"`
		Metadata   *NodeMetadata   `json:"metadata,omitempty"`
		Condition  json.RawMessage `json:"condition"`
		Consequent *BlockStatement `json:"consequent"`
		Alternate  json.RawMessage `json:"alternate,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	cond, err := UnmarshalExpressionJSON(raw.Condition)
	if err != nil {
		return fmt.Errorf("IfStatement.condition unmarshal error: %w", err)
	}
	var alt Statement
	if len(raw.Alternate) > 0 && string(raw.Alternate) != "null" {
		alt, err = UnmarshalStatementJSON(raw.Alternate)
		if err != nil {
			return fmt.Errorf("IfStatement.alternate unmarshal error: %w", err)
		}
	}
	i.NodeType = NodeTypeIfStatement
	i.Metadata = raw.Metadata
	i.Condition = cond
	i.Consequent = raw.Consequent
	i.Alternate = alt
	return nil
}

// CaseClause represents a branch in a SwitchStatement.
// If Value is nil, it is the default fallback branch.
type CaseClause struct {
	NodeType NodeType    `json:"type"`
	Value    Expression  `json:"value,omitempty"`
	Body     []Statement `json:"body"`
}

func (c *CaseClause) Type() NodeType          { return NodeTypeCaseClause }
func (c *CaseClause) GetMetadata() *NodeMetadata { return nil }
func (c *CaseClause) statementNode()          {}

func (c *CaseClause) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType          `json:"type"`
		Value    json.RawMessage   `json:"value,omitempty"`
		Body     []json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var valExpr Expression
	if len(raw.Value) > 0 && string(raw.Value) != "null" {
		var err error
		valExpr, err = UnmarshalExpressionJSON(raw.Value)
		if err != nil {
			return fmt.Errorf("CaseClause.value unmarshal error: %w", err)
		}
	}
	stmts := make([]Statement, len(raw.Body))
	for i, stmtRaw := range raw.Body {
		stmt, err := UnmarshalStatementJSON(stmtRaw)
		if err != nil {
			return fmt.Errorf("CaseClause.body[%d] unmarshal error: %w", i, err)
		}
		stmts[i] = stmt
	}
	c.NodeType = NodeTypeCaseClause
	c.Value = valExpr
	c.Body = stmts
	return nil
}

// SwitchStatement represents value matching (cocok ... kasus ... lainnya).
type SwitchStatement struct {
	NodeType     NodeType      `json:"type"`
	Metadata     *NodeMetadata `json:"metadata,omitempty"`
	Discriminant Expression    `json:"discriminant"`
	Cases        []CaseClause  `json:"cases"`
}

func (s *SwitchStatement) Type() NodeType          { return NodeTypeSwitchStatement }
func (s *SwitchStatement) GetMetadata() *NodeMetadata { return s.Metadata }
func (s *SwitchStatement) statementNode()          {}

func (s *SwitchStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType     NodeType        `json:"type"`
		Metadata     *NodeMetadata   `json:"metadata,omitempty"`
		Discriminant json.RawMessage `json:"discriminant"`
		Cases        []CaseClause    `json:"cases"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	discExpr, err := UnmarshalExpressionJSON(raw.Discriminant)
	if err != nil {
		return fmt.Errorf("SwitchStatement.discriminant unmarshal error: %w", err)
	}
	if raw.Cases == nil {
		raw.Cases = []CaseClause{}
	}
	s.NodeType = NodeTypeSwitchStatement
	s.Metadata = raw.Metadata
	s.Discriminant = discExpr
	s.Cases = raw.Cases
	return nil
}

// WhileStatement represents repetition while condition is true.
type WhileStatement struct {
	NodeType  NodeType        `json:"type"`
	Metadata  *NodeMetadata   `json:"metadata,omitempty"`
	Condition Expression      `json:"condition"`
	Body      *BlockStatement `json:"body"`
}

func (w *WhileStatement) Type() NodeType          { return NodeTypeWhileStatement }
func (w *WhileStatement) GetMetadata() *NodeMetadata { return w.Metadata }
func (w *WhileStatement) statementNode()          {}

func (w *WhileStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType  NodeType        `json:"type"`
		Metadata  *NodeMetadata   `json:"metadata,omitempty"`
		Condition json.RawMessage `json:"condition"`
		Body      *BlockStatement `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	cond, err := UnmarshalExpressionJSON(raw.Condition)
	if err != nil {
		return fmt.Errorf("WhileStatement.condition unmarshal error: %w", err)
	}
	w.NodeType = NodeTypeWhileStatement
	w.Metadata = raw.Metadata
	w.Condition = cond
	w.Body = raw.Body
	return nil
}

// ForStatement represents bounded iteration.
type ForStatement struct {
	NodeType  NodeType        `json:"type"`
	Metadata  *NodeMetadata   `json:"metadata,omitempty"`
	Init      Statement       `json:"init,omitempty"`
	Condition Expression      `json:"condition,omitempty"`
	Update    Statement       `json:"update,omitempty"`
	Body      *BlockStatement `json:"body"`
}

func (f *ForStatement) Type() NodeType          { return NodeTypeForStatement }
func (f *ForStatement) GetMetadata() *NodeMetadata { return f.Metadata }
func (f *ForStatement) statementNode()          {}

func (f *ForStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType  NodeType        `json:"type"`
		Metadata  *NodeMetadata   `json:"metadata,omitempty"`
		Init      json.RawMessage `json:"init,omitempty"`
		Condition json.RawMessage `json:"condition,omitempty"`
		Update    json.RawMessage `json:"update,omitempty"`
		Body      *BlockStatement `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var initStmt Statement
	if len(raw.Init) > 0 && string(raw.Init) != "null" {
		var err error
		initStmt, err = UnmarshalStatementJSON(raw.Init)
		if err != nil {
			return fmt.Errorf("ForStatement.init unmarshal error: %w", err)
		}
	}
	var condExpr Expression
	if len(raw.Condition) > 0 && string(raw.Condition) != "null" {
		var err error
		condExpr, err = UnmarshalExpressionJSON(raw.Condition)
		if err != nil {
			return fmt.Errorf("ForStatement.condition unmarshal error: %w", err)
		}
	}
	var updateStmt Statement
	if len(raw.Update) > 0 && string(raw.Update) != "null" {
		var err error
		updateStmt, err = UnmarshalStatementJSON(raw.Update)
		if err != nil {
			return fmt.Errorf("ForStatement.update unmarshal error: %w", err)
		}
	}
	f.NodeType = NodeTypeForStatement
	f.Metadata = raw.Metadata
	f.Init = initStmt
	f.Condition = condExpr
	f.Update = updateStmt
	f.Body = raw.Body
	return nil
}

// ForEachStatement represents iterating through a collection (e.g. untuk item dalam daftar { ... }).
type ForEachStatement struct {
	NodeType NodeType        `json:"type"`
	Metadata *NodeMetadata   `json:"metadata,omitempty"`
	Item     string          `json:"item"`
	Iterable Expression      `json:"iterable"`
	Body     *BlockStatement `json:"body"`
}

func (f *ForEachStatement) Type() NodeType          { return NodeTypeForEachStatement }
func (f *ForEachStatement) GetMetadata() *NodeMetadata { return f.Metadata }
func (f *ForEachStatement) statementNode()          {}

func (f *ForEachStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Item     string          `json:"item"`
		Iterable json.RawMessage `json:"iterable"`
		Body     *BlockStatement `json:"body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	iterableExpr, err := UnmarshalExpressionJSON(raw.Iterable)
	if err != nil {
		return fmt.Errorf("ForEachStatement.iterable unmarshal error: %w", err)
	}
	f.NodeType = NodeTypeForEachStatement
	f.Metadata = raw.Metadata
	f.Item = raw.Item
	f.Iterable = iterableExpr
	f.Body = raw.Body
	return nil
}

// Parameter defines a function parameter with optional type hint and default value.
type Parameter struct {
	Name         string     `json:"name"`
	TypeHint     string     `json:"type_hint,omitempty"`
	DefaultValue Expression `json:"default_value,omitempty"`
}

func (p *Parameter) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name         string          `json:"name"`
		TypeHint     string          `json:"type_hint,omitempty"`
		DefaultValue json.RawMessage `json:"default_value,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Name = raw.Name
	p.TypeHint = raw.TypeHint
	if len(raw.DefaultValue) > 0 && string(raw.DefaultValue) != "null" {
		defVal, err := UnmarshalExpressionJSON(raw.DefaultValue)
		if err != nil {
			return fmt.Errorf("Parameter.default_value unmarshal error: %w", err)
		}
		p.DefaultValue = defVal
	}
	return nil
}

// FunctionDeclaration represents a named function definition.
type FunctionDeclaration struct {
	NodeType   NodeType        `json:"type"`
	Metadata   *NodeMetadata   `json:"metadata,omitempty"`
	Name       string          `json:"name"`
	Parameters []Parameter     `json:"parameters"`
	Body       *BlockStatement `json:"body"`
	ReturnType string          `json:"return_type,omitempty"`
	Visibility *Visibility     `json:"visibility,omitempty"`
}

func (f *FunctionDeclaration) Type() NodeType          { return NodeTypeFunctionDeclaration }
func (f *FunctionDeclaration) GetMetadata() *NodeMetadata { return f.Metadata }
func (f *FunctionDeclaration) statementNode()          {}

// ReturnStatement exits a function and optionally produces a value.
type ReturnStatement struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Value    Expression    `json:"value,omitempty"`
}

func (r *ReturnStatement) Type() NodeType          { return NodeTypeReturnStatement }
func (r *ReturnStatement) GetMetadata() *NodeMetadata { return r.Metadata }
func (r *ReturnStatement) statementNode()          {}

func (r *ReturnStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType        `json:"type"`
		Metadata *NodeMetadata   `json:"metadata,omitempty"`
		Value    json.RawMessage `json:"value,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var val Expression
	if len(raw.Value) > 0 && string(raw.Value) != "null" {
		var err error
		val, err = UnmarshalExpressionJSON(raw.Value)
		if err != nil {
			return fmt.Errorf("ReturnStatement.value unmarshal error: %w", err)
		}
	}
	r.NodeType = NodeTypeReturnStatement
	r.Metadata = raw.Metadata
	r.Value = val
	return nil
}

// ConcurrencyStatement represents concurrent execution block (e.g. 'sambil', 'bersamaan dengan').
type ConcurrencyStatement struct {
	NodeType NodeType      `json:"type"`
	Metadata *NodeMetadata `json:"metadata,omitempty"`
	Mode     string        `json:"mode"` // e.g. "sambil" (background goroutine) or "bersamaan" (fan-out wait group)
	Tasks    []Statement   `json:"tasks"`
}

func (c *ConcurrencyStatement) Type() NodeType          { return NodeTypeConcurrencyStatement }
func (c *ConcurrencyStatement) GetMetadata() *NodeMetadata { return c.Metadata }
func (c *ConcurrencyStatement) statementNode()          {}

func (c *ConcurrencyStatement) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType NodeType          `json:"type"`
		Metadata *NodeMetadata     `json:"metadata,omitempty"`
		Mode     string            `json:"mode"`
		Tasks    []json.RawMessage `json:"tasks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	tasks := make([]Statement, len(raw.Tasks))
	for i, taskRaw := range raw.Tasks {
		stmt, err := UnmarshalStatementJSON(taskRaw)
		if err != nil {
			return fmt.Errorf("ConcurrencyStatement.tasks[%d] unmarshal error: %w", i, err)
		}
		tasks[i] = stmt
	}
	c.NodeType = NodeTypeConcurrencyStatement
	c.Metadata = raw.Metadata
	c.Mode = raw.Mode
	c.Tasks = tasks
	return nil
}

// ImportSpecifier represents an imported identifier or alias.
type ImportSpecifier struct {
	Imported string `json:"imported"`
	Local    string `json:"local"`
	IsAll    bool   `json:"is_all,omitempty"`
}

// ImportDeclaration represents importing symbols from another module/file.
type ImportDeclaration struct {
	NodeType   NodeType          `json:"type"`
	Metadata   *NodeMetadata     `json:"metadata,omitempty"`
	Source     string            `json:"source"`
	Specifiers []ImportSpecifier `json:"specifiers"`
}

func (i *ImportDeclaration) Type() NodeType          { return NodeTypeImportDeclaration }
func (i *ImportDeclaration) GetMetadata() *NodeMetadata { return i.Metadata }
func (i *ImportDeclaration) statementNode()          {}

// VisibilityStatement represents a standalone export/visibility declaration (e.g. publik! { kali, bagi }).
type VisibilityStatement struct {
	NodeType    NodeType      `json:"type"`
	Metadata    *NodeMetadata `json:"metadata,omitempty"`
	Visibility  Visibility    `json:"visibility"`
	Identifiers []string      `json:"identifiers"`
}

func (v *VisibilityStatement) Type() NodeType          { return NodeTypeVisibilityStatement }
func (v *VisibilityStatement) GetMetadata() *NodeMetadata { return v.Metadata }
func (v *VisibilityStatement) statementNode()          {}
