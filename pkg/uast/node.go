package uast

// NodeType identifies the specific kind of AST node.
type NodeType string

const (
	// Root and structural nodes
	NodeTypeProgram             NodeType = "Program"
	NodeTypeBlockStatement      NodeType = "BlockStatement"
	NodeTypeExpressionStatement NodeType = "ExpressionStatement"

	// Statements
	NodeTypeVariableDeclaration  NodeType = "VariableDeclaration"
	NodeTypeAssignmentStatement  NodeType = "AssignmentStatement"
	NodeTypeIfStatement          NodeType = "IfStatement"
	NodeTypeWhileStatement       NodeType = "WhileStatement"
	NodeTypeForStatement         NodeType = "ForStatement"
	NodeTypeFunctionDeclaration  NodeType = "FunctionDeclaration"
	NodeTypeReturnStatement      NodeType = "ReturnStatement"
	NodeTypeConcurrencyStatement NodeType = "ConcurrencyStatement"

	// Expressions
	NodeTypeLiteral          NodeType = "Literal"
	NodeTypeIdentifier       NodeType = "Identifier"
	NodeTypeBinaryExpression NodeType = "BinaryExpression"
	NodeTypeUnaryExpression  NodeType = "UnaryExpression"
	NodeTypeCallExpression   NodeType = "CallExpression"
	NodeTypeMemberExpression NodeType = "MemberExpression"
	NodeTypeArrayLiteral     NodeType = "ArrayLiteral"
)

// Position captures source location information for LSP, tooling, and error reporting.
type Position struct {
	File   string `json:"file,omitempty"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

// NodeMetadata holds extra information attached to any U-AST node.
// ID is useful for diffing, LSP tracking, and incremental evaluation.
// StyleHint can carry syntactic preferences or dialect notes.
type NodeMetadata struct {
	ID        string    `json:"id,omitempty"`
	Pos       *Position `json:"pos,omitempty"`
	Comments  []string  `json:"comments,omitempty"`
	StyleHint string    `json:"style_hint,omitempty"`
}

// Node represents any node in the Universal Abstract Syntax Tree.
type Node interface {
	Type() NodeType
	GetMetadata() *NodeMetadata
}

// Statement represents an executive sentence/statement in U-AST.
type Statement interface {
	Node
	statementNode()
}

// Expression represents a value-producing construct in U-AST.
type Expression interface {
	Node
	expressionNode()
}
