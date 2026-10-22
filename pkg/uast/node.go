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
	NodeTypeSwitchStatement      NodeType = "SwitchStatement"
	NodeTypeCaseClause           NodeType = "CaseClause"
	NodeTypeWhileStatement       NodeType = "WhileStatement"
	NodeTypeForStatement         NodeType = "ForStatement"
	NodeTypeForEachStatement     NodeType = "ForEachStatement"
	NodeTypeFunctionDeclaration  NodeType = "FunctionDeclaration"
	NodeTypeReturnStatement      NodeType = "ReturnStatement"
	NodeTypeBreakStatement       NodeType = "BreakStatement"
	NodeTypeContinueStatement    NodeType = "ContinueStatement"
	NodeTypeConcurrencyStatement NodeType = "ConcurrencyStatement"
	NodeTypeImportDeclaration    NodeType = "ImportDeclaration"
	NodeTypeVisibilityStatement  NodeType = "VisibilityStatement"
	NodeTypeTryCatchStatement    NodeType = "TryCatchStatement"
	NodeTypeCatchClause          NodeType = "CatchClause"
	NodeTypeThrowStatement       NodeType = "ThrowStatement"
	NodeTypeStyleBlockStatement  NodeType = "StyleBlockStatement"
	NodeTypeDestructureDeclaration NodeType = "DestructureDeclaration"
	NodeTypeDestructureAssignment  NodeType = "DestructureAssignment"
	NodeTypeStructDeclaration      NodeType = "StructDeclaration"
	NodeTypeReferenceDeclaration   NodeType = "ReferenceDeclaration"
	NodeTypeTypeAliasDeclaration   NodeType = "TypeAliasDeclaration"
	NodeTypeDeferStatement         NodeType = "DeferStatement"
	NodeTypeTransactionStatement   NodeType = "TransactionStatement"
	NodeTypeRollbackStatement      NodeType = "RollbackStatement"
	NodeTypeReactiveDeclaration    NodeType = "ReactiveDeclaration"
	NodeTypeWatchStatement         NodeType = "WatchStatement"
	NodeTypeTestStatement          NodeType = "TestStatement"

	// Expressions
	NodeTypeLiteral          NodeType = "Literal"
	NodeTypeIdentifier       NodeType = "Identifier"
	NodeTypeBinaryExpression NodeType = "BinaryExpression"
	NodeTypeUnaryExpression  NodeType = "UnaryExpression"
	NodeTypeCallExpression   NodeType = "CallExpression"
	NodeTypeMemberExpression NodeType = "MemberExpression"
	NodeTypeArrayLiteral        NodeType = "ArrayLiteral"
	NodeTypeMapLiteral          NodeType = "MapLiteral"
	NodeTypePipelineExpression  NodeType = "PipelineExpression"
	NodeTypeFunctionExpression  NodeType = "FunctionExpression"
	NodeTypeTemplateLiteral     NodeType = "TemplateLiteral"
	NodeTypeCastExpression       NodeType = "CastExpression"
	NodeTypeSafeNavigationExpression NodeType = "SafeNavigationExpression"
	NodeTypeNullCoalesceExpression   NodeType = "NullCoalesceExpression"
	NodeTypeRangeExpression          NodeType = "RangeExpression"
	NodeTypeTitipExpression          NodeType = "TitipExpression"
	NodeTypeJemputExpression         NodeType = "JemputExpression"
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

// VisibilityKind defines the scope of access for a symbol.
type VisibilityKind string

const (
	VisibilityPrivate   VisibilityKind = "private"   // privat! (default: file-scoped)
	VisibilityModule    VisibilityKind = "module"    // privat-modul! (module-scoped)
	VisibilityPublic    VisibilityKind = "public"    // publik! (unrestricted)
	VisibilityAllowlist VisibilityKind = "allowlist" // privat-["./target"]! (specific targets)
)

// Visibility defines the access control attached to a declaration or export.
type Visibility struct {
	Kind    VisibilityKind `json:"kind"`
	Targets []string       `json:"targets,omitempty"` // For allowlist: relative file paths or module names
}
