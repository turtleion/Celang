package uast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ToJSON serializes any U-AST node to JSON.
func ToJSON(node Node, indent bool) ([]byte, error) {
	if indent {
		return json.MarshalIndent(node, "", "  ")
	}
	return json.Marshal(node)
}

// FromJSON deserializes JSON into a *Program node.
func FromJSON(data []byte) (*Program, error) {
	var prog Program
	if err := json.Unmarshal(data, &prog); err != nil {
		return nil, fmt.Errorf("failed to parse Program from JSON: %w", err)
	}
	return &prog, nil
}

// UnmarshalExpressionJSON dynamically inspects the 'type' attribute and returns the concrete Expression.
func UnmarshalExpressionJSON(data []byte) (Expression, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	var detector struct {
		Type NodeType `json:"type"`
	}
	if err := json.Unmarshal(data, &detector); err != nil {
		return nil, fmt.Errorf("failed to determine expression node type: %w", err)
	}

	switch detector.Type {
	case NodeTypeLiteral:
		var lit Literal
		if err := json.Unmarshal(data, &lit); err != nil {
			return nil, err
		}
		return &lit, nil

	case NodeTypeIdentifier:
		var id Identifier
		if err := json.Unmarshal(data, &id); err != nil {
			return nil, err
		}
		id.NodeType = NodeTypeIdentifier
		return &id, nil

	case NodeTypeBinaryExpression:
		var bin BinaryExpression
		if err := json.Unmarshal(data, &bin); err != nil {
			return nil, err
		}
		return &bin, nil

	case NodeTypeUnaryExpression:
		var un UnaryExpression
		if err := json.Unmarshal(data, &un); err != nil {
			return nil, err
		}
		return &un, nil

	case NodeTypeCallExpression:
		var call CallExpression
		if err := json.Unmarshal(data, &call); err != nil {
			return nil, err
		}
		return &call, nil

	case NodeTypeMemberExpression:
		var mem MemberExpression
		if err := json.Unmarshal(data, &mem); err != nil {
			return nil, err
		}
		return &mem, nil

	case NodeTypeArrayLiteral:
		var arr ArrayLiteral
		if err := json.Unmarshal(data, &arr); err != nil {
			return nil, err
		}
		return &arr, nil

	case NodeTypeMapLiteral:
		var m MapLiteral
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		return &m, nil

	default:
		return nil, fmt.Errorf("unrecognized expression type: '%s'", detector.Type)
	}
}

// UnmarshalStatementJSON dynamically inspects the 'type' attribute and returns the concrete Statement.
func UnmarshalStatementJSON(data []byte) (Statement, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	var detector struct {
		Type NodeType `json:"type"`
	}
	if err := json.Unmarshal(data, &detector); err != nil {
		return nil, fmt.Errorf("failed to determine statement node type: %w", err)
	}

	switch detector.Type {
	case NodeTypeExpressionStatement:
		var exprStmt ExpressionStatement
		if err := json.Unmarshal(data, &exprStmt); err != nil {
			return nil, err
		}
		return &exprStmt, nil

	case NodeTypeVariableDeclaration:
		var varDecl VariableDeclaration
		if err := json.Unmarshal(data, &varDecl); err != nil {
			return nil, err
		}
		return &varDecl, nil

	case NodeTypeAssignmentStatement:
		var assign AssignmentStatement
		if err := json.Unmarshal(data, &assign); err != nil {
			return nil, err
		}
		return &assign, nil

	case NodeTypeBlockStatement:
		var block BlockStatement
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return &block, nil

	case NodeTypeIfStatement:
		var ifStmt IfStatement
		if err := json.Unmarshal(data, &ifStmt); err != nil {
			return nil, err
		}
		return &ifStmt, nil

	case NodeTypeSwitchStatement:
		var switchStmt SwitchStatement
		if err := json.Unmarshal(data, &switchStmt); err != nil {
			return nil, err
		}
		return &switchStmt, nil

	case NodeTypeWhileStatement:
		var whileStmt WhileStatement
		if err := json.Unmarshal(data, &whileStmt); err != nil {
			return nil, err
		}
		return &whileStmt, nil

	case NodeTypeForStatement:
		var forStmt ForStatement
		if err := json.Unmarshal(data, &forStmt); err != nil {
			return nil, err
		}
		return &forStmt, nil

	case NodeTypeForEachStatement:
		var forEachStmt ForEachStatement
		if err := json.Unmarshal(data, &forEachStmt); err != nil {
			return nil, err
		}
		return &forEachStmt, nil

	case NodeTypeFunctionDeclaration:
		var fnDecl FunctionDeclaration
		if err := json.Unmarshal(data, &fnDecl); err != nil {
			return nil, err
		}
		fnDecl.NodeType = NodeTypeFunctionDeclaration
		return &fnDecl, nil

	case NodeTypeReturnStatement:
		var retStmt ReturnStatement
		if err := json.Unmarshal(data, &retStmt); err != nil {
			return nil, err
		}
		return &retStmt, nil

	case NodeTypeConcurrencyStatement:
		var concStmt ConcurrencyStatement
		if err := json.Unmarshal(data, &concStmt); err != nil {
			return nil, err
		}
		return &concStmt, nil

	case NodeTypeImportDeclaration:
		var imp ImportDeclaration
		if err := json.Unmarshal(data, &imp); err != nil {
			return nil, err
		}
		imp.NodeType = NodeTypeImportDeclaration
		return &imp, nil

	case NodeTypeVisibilityStatement:
		var vis VisibilityStatement
		if err := json.Unmarshal(data, &vis); err != nil {
			return nil, err
		}
		vis.NodeType = NodeTypeVisibilityStatement
		return &vis, nil

	default:
		return nil, fmt.Errorf("unrecognized statement type: '%s'", detector.Type)
	}
}

// -------------------------------------------------------------
// Helper Constructors for Clean AST Construction
// -------------------------------------------------------------

func NewProgram(body ...Statement) *Program {
	if body == nil {
		body = []Statement{}
	}
	return &Program{
		NodeType: NodeTypeProgram,
		Version:  "1.0",
		Body:     body,
	}
}

func NewIntLiteral(val int64) *Literal {
	return &Literal{
		NodeType: NodeTypeLiteral,
		Kind:     LiteralKindInt,
		Raw:      strconv.FormatInt(val, 10),
		Value:    val,
	}
}

func NewFloatLiteral(val float64) *Literal {
	return &Literal{
		NodeType: NodeTypeLiteral,
		Kind:     LiteralKindFloat,
		Raw:      strconv.FormatFloat(val, 'f', -1, 64),
		Value:    val,
	}
}

func NewStringLiteral(val string) *Literal {
	return &Literal{
		NodeType: NodeTypeLiteral,
		Kind:     LiteralKindString,
		Raw:      strconv.Quote(val),
		Value:    val,
	}
}

func NewBoolLiteral(val bool) *Literal {
	return &Literal{
		NodeType: NodeTypeLiteral,
		Kind:     LiteralKindBool,
		Raw:      strconv.FormatBool(val),
		Value:    val,
	}
}

func NewNullLiteral() *Literal {
	return &Literal{
		NodeType: NodeTypeLiteral,
		Kind:     LiteralKindNull,
		Raw:      "null",
		Value:    nil,
	}
}

func NewIdentifier(name string) *Identifier {
	return &Identifier{
		NodeType: NodeTypeIdentifier,
		Name:     name,
	}
}

func NewBinaryExpr(left Expression, op string, right Expression) *BinaryExpression {
	return &BinaryExpression{
		NodeType: NodeTypeBinaryExpression,
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

func NewUnaryExpr(op string, arg Expression, isPrefix bool) *UnaryExpression {
	return &UnaryExpression{
		NodeType: NodeTypeUnaryExpression,
		Operator: op,
		Argument: arg,
		IsPrefix: isPrefix,
	}
}

func NewCallExpr(callee Expression, args ...Expression) *CallExpression {
	if args == nil {
		args = []Expression{}
	}
	return &CallExpression{
		NodeType:  NodeTypeCallExpression,
		Callee:    callee,
		Arguments: args,
	}
}

func NewMemberExpr(obj Expression, prop string, computed bool) *MemberExpression {
	return &MemberExpression{
		NodeType: NodeTypeMemberExpression,
		Object:   obj,
		Property: prop,
		Computed: computed,
	}
}

func NewArrayLiteral(elements ...Expression) *ArrayLiteral {
	if elements == nil {
		elements = []Expression{}
	}
	return &ArrayLiteral{
		NodeType: NodeTypeArrayLiteral,
		Elements: elements,
	}
}

func NewMapLiteral(properties ...MapProperty) *MapLiteral {
	if properties == nil {
		properties = []MapProperty{}
	}
	return &MapLiteral{
		NodeType:   NodeTypeMapLiteral,
		Properties: properties,
	}
}

func NewMapProperty(key, value Expression) MapProperty {
	return MapProperty{
		Key:   key,
		Value: value,
	}
}

func NewVarDecl(name string, val Expression, isMutable bool, typeHint string) *VariableDeclaration {
	return &VariableDeclaration{
		NodeType:  NodeTypeVariableDeclaration,
		Name:      name,
		Value:     val,
		IsMutable: isMutable,
		TypeHint:  typeHint,
	}
}

func NewAssign(target Expression, op string, val Expression) *AssignmentStatement {
	return &AssignmentStatement{
		NodeType: NodeTypeAssignmentStatement,
		Target:   target,
		Operator: op,
		Value:    val,
	}
}

func NewExprStmt(expr Expression) *ExpressionStatement {
	return &ExpressionStatement{
		NodeType:   NodeTypeExpressionStatement,
		Expression: expr,
	}
}

func NewBlock(stmts ...Statement) *BlockStatement {
	if stmts == nil {
		stmts = []Statement{}
	}
	return &BlockStatement{
		NodeType:   NodeTypeBlockStatement,
		Statements: stmts,
	}
}

func NewIf(cond Expression, cons *BlockStatement, alt Statement) *IfStatement {
	return &IfStatement{
		NodeType:   NodeTypeIfStatement,
		Condition:  cond,
		Consequent: cons,
		Alternate:  alt,
	}
}

func NewSwitch(discriminant Expression, cases ...CaseClause) *SwitchStatement {
	if cases == nil {
		cases = []CaseClause{}
	}
	return &SwitchStatement{
		NodeType:     NodeTypeSwitchStatement,
		Discriminant: discriminant,
		Cases:        cases,
	}
}

func NewCase(val Expression, body ...Statement) CaseClause {
	if body == nil {
		body = []Statement{}
	}
	return CaseClause{
		NodeType: NodeTypeCaseClause,
		Value:    val,
		Body:     body,
	}
}

func NewWhile(cond Expression, body *BlockStatement) *WhileStatement {
	return &WhileStatement{
		NodeType:  NodeTypeWhileStatement,
		Condition: cond,
		Body:      body,
	}
}

func NewFor(init Statement, cond Expression, update Statement, body *BlockStatement) *ForStatement {
	return &ForStatement{
		NodeType:  NodeTypeForStatement,
		Init:      init,
		Condition: cond,
		Update:    update,
		Body:      body,
	}
}

func NewForEach(item string, iterable Expression, body *BlockStatement) *ForEachStatement {
	return &ForEachStatement{
		NodeType: NodeTypeForEachStatement,
		Item:     item,
		Iterable: iterable,
		Body:     body,
	}
}

func NewFuncDecl(name string, params []Parameter, body *BlockStatement, retType string) *FunctionDeclaration {
	if params == nil {
		params = []Parameter{}
	}
	return &FunctionDeclaration{
		NodeType:   NodeTypeFunctionDeclaration,
		Name:       name,
		Parameters: params,
		Body:       body,
		ReturnType: retType,
	}
}

func NewReturn(val Expression) *ReturnStatement {
	return &ReturnStatement{
		NodeType: NodeTypeReturnStatement,
		Value:    val,
	}
}

func NewConcurrencyStmt(mode string, tasks ...Statement) *ConcurrencyStatement {
	if tasks == nil {
		tasks = []Statement{}
	}
	return &ConcurrencyStatement{
		NodeType: NodeTypeConcurrencyStatement,
		Mode:     mode,
		Tasks:    tasks,
	}
}
