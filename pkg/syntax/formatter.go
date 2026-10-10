package syntax

import (
	"fmt"
	"strings"

	"celang/pkg/uast"
)

// Formatter converts a U-AST into source code formatted in a chosen StyleProfile.
// This is the Reverse Parser (AST -> Text) component of the Syntax Styling Engine.
type Formatter struct {
	profile *StyleProfile
}

func NewFormatter(profile *StyleProfile) *Formatter {
	if profile == nil {
		profile = IndonesianProfile()
	}
	return &Formatter{profile: profile}
}

// Format converts an entire U-AST Program into formatted text.
func (f *Formatter) Format(prog *uast.Program) string {
	if prog == nil {
		return ""
	}

	var sb strings.Builder
	for i, stmt := range prog.Body {
		if i > 0 {
			// Add blank line before major blocks for readability
			if isMajorBlock(stmt) {
				sb.WriteString("\n")
			}
		}
		sb.WriteString(f.formatStatement(stmt, 0))
		sb.WriteString("\n")
	}

	return strings.TrimRight(sb.String(), "\n") + "\n"
}

func (f *Formatter) formatStatement(stmt uast.Statement, indent int) string {
	prefix := strings.Repeat(" ", indent*f.profile.IndentWidth)

	switch s := stmt.(type) {
	case *uast.VariableDeclaration:
		var sb strings.Builder
		sb.WriteString(prefix)

		if s.IsMutable {
			if f.profile.VarDecl != "" {
				sb.WriteString(f.profile.VarDecl)
				sb.WriteString(" ")
			}
		} else {
			sb.WriteString(f.profile.ConstDecl)
			sb.WriteString(" ")
		}

		sb.WriteString(s.Name)

		if s.TypeHint != "" {
			sb.WriteString(" ")
			sb.WriteString(f.profile.TypeAssign)
			sb.WriteString(" ")
			sb.WriteString(f.formatTypeHint(s.TypeHint))
		}

		if s.Value != nil {
			sb.WriteString(" = ")
			sb.WriteString(f.formatExpression(s.Value))
		}

		return sb.String()

	case *uast.AssignmentStatement:
		targetStr := f.formatExpression(s.Target)
		valStr := f.formatExpression(s.Value)
		return fmt.Sprintf("%s%s %s %s", prefix, targetStr, s.Operator, valStr)

	case *uast.IfStatement:
		var sb strings.Builder
		sb.WriteString(prefix)
		sb.WriteString(f.profile.If)
		sb.WriteString(" ")
		sb.WriteString(f.formatExpression(s.Condition))
		sb.WriteString(" {\n")

		if s.Consequent != nil {
			for _, child := range s.Consequent.Statements {
				sb.WriteString(f.formatStatement(child, indent+1))
				sb.WriteString("\n")
			}
		}

		currAlt := s.Alternate
		for currAlt != nil {
			if elseIf, ok := currAlt.(*uast.IfStatement); ok {
				sb.WriteString(prefix)
				sb.WriteString("} ")
				sb.WriteString(f.profile.ElseIf)
				sb.WriteString(" ")
				sb.WriteString(f.formatExpression(elseIf.Condition))
				sb.WriteString(" {\n")

				if elseIf.Consequent != nil {
					for _, child := range elseIf.Consequent.Statements {
						sb.WriteString(f.formatStatement(child, indent+1))
						sb.WriteString("\n")
					}
				}
				currAlt = elseIf.Alternate
			} else if block, ok := currAlt.(*uast.BlockStatement); ok {
				sb.WriteString(prefix)
				sb.WriteString("} ")
				sb.WriteString(f.profile.Else)
				sb.WriteString(" {\n")

				for _, child := range block.Statements {
					sb.WriteString(f.formatStatement(child, indent+1))
					sb.WriteString("\n")
				}
				currAlt = nil
			} else {
				break
			}
		}

		sb.WriteString(prefix)
		sb.WriteString("}")
		return sb.String()

	case *uast.SwitchStatement:
		var sb strings.Builder
		sb.WriteString(prefix)
		sb.WriteString(f.profile.Switch)
		if f.profile.Name == "javascript" {
			sb.WriteString(" (")
			sb.WriteString(f.formatExpression(s.Discriminant))
			sb.WriteString(") {\n")
		} else {
			sb.WriteString(" ")
			sb.WriteString(f.formatExpression(s.Discriminant))
			sb.WriteString(" {\n")
		}

		for _, c := range s.Cases {
			casePrefix := strings.Repeat(" ", (indent+1)*f.profile.IndentWidth)
			if c.Value != nil {
				sb.WriteString(casePrefix)
				sb.WriteString(f.profile.Case)
				sb.WriteString(" ")
				sb.WriteString(f.formatExpression(c.Value))
				sb.WriteString(":\n")
			} else {
				sb.WriteString(casePrefix)
				sb.WriteString(f.profile.Default)
				sb.WriteString(":\n")
			}
			for _, child := range c.Body {
				sb.WriteString(f.formatStatement(child, indent+2))
				sb.WriteString("\n")
			}
		}

		sb.WriteString(prefix)
		sb.WriteString("}")
		return sb.String()

	case *uast.WhileStatement:
		var sb strings.Builder
		sb.WriteString(prefix)
		sb.WriteString(f.profile.While)
		sb.WriteString(" ")
		sb.WriteString(f.formatExpression(s.Condition))
		sb.WriteString(" {\n")

		if s.Body != nil {
			for _, child := range s.Body.Statements {
				sb.WriteString(f.formatStatement(child, indent+1))
				sb.WriteString("\n")
			}
		}

		sb.WriteString(prefix)
		sb.WriteString("}")
		return sb.String()

	case *uast.ForEachStatement:
		var sb strings.Builder
		if f.profile.Name == "javascript" {
			sb.WriteString(fmt.Sprintf("%sfor (let %s of %s) {\n", prefix, s.Item, f.formatExpression(s.Iterable)))
		} else {
			sb.WriteString(prefix)
			sb.WriteString(f.profile.For)
			sb.WriteString(" ")
			sb.WriteString(s.Item)
			sb.WriteString(" ")
			sb.WriteString(f.profile.In)
			sb.WriteString(" ")
			sb.WriteString(f.formatExpression(s.Iterable))
			sb.WriteString(" {\n")
		}

		if s.Body != nil {
			for _, child := range s.Body.Statements {
				sb.WriteString(f.formatStatement(child, indent+1))
				sb.WriteString("\n")
			}
		}

		sb.WriteString(prefix)
		sb.WriteString("}")
		return sb.String()

	case *uast.FunctionDeclaration:
		var sb strings.Builder
		sb.WriteString(prefix)
		sb.WriteString(f.profile.Func)
		sb.WriteString(" ")
		sb.WriteString(s.Name)
		sb.WriteString("(")

		params := make([]string, len(s.Parameters))
		for i, p := range s.Parameters {
			pStr := p.Name
			if p.TypeHint != "" {
				pStr += " " + f.profile.TypeAssign + " " + f.formatTypeHint(p.TypeHint)
			}
			if p.DefaultValue != nil {
				pStr += " = " + f.formatExpression(p.DefaultValue)
			}
			params[i] = pStr
		}
		sb.WriteString(strings.Join(params, ", "))
		sb.WriteString(")")

		if s.ReturnType != "" && f.profile.ReturnArrow != "" {
			sb.WriteString(" ")
			sb.WriteString(f.profile.ReturnArrow)
			sb.WriteString(" ")
			sb.WriteString(f.formatTypeHint(s.ReturnType))
		}

		sb.WriteString(" {\n")
		if s.Body != nil {
			for _, child := range s.Body.Statements {
				sb.WriteString(f.formatStatement(child, indent+1))
				sb.WriteString("\n")
			}
		}
		sb.WriteString(prefix)
		sb.WriteString("}")
		return sb.String()

	case *uast.ReturnStatement:
		if s.Value != nil {
			return fmt.Sprintf("%s%s %s", prefix, f.profile.Return, f.formatExpression(s.Value))
		}
		return fmt.Sprintf("%s%s", prefix, f.profile.Return)

	case *uast.ConcurrencyStatement:
		var sb strings.Builder
		modeKeyword := f.profile.Async
		if s.Mode == "bersamaan" || s.Mode == "parallel" {
			modeKeyword = f.profile.Parallel
		}

		sb.WriteString(prefix)
		sb.WriteString(modeKeyword)
		sb.WriteString(" {\n")

		for _, task := range s.Tasks {
			sb.WriteString(f.formatStatement(task, indent+1))
			sb.WriteString("\n")
		}

		sb.WriteString(prefix)
		sb.WriteString("}")
		return sb.String()

	case *uast.ExpressionStatement:
		return prefix + f.formatExpression(s.Expression)

	case *uast.BlockStatement:
		var sb strings.Builder
		sb.WriteString(prefix + "{\n")
		for _, child := range s.Statements {
			sb.WriteString(f.formatStatement(child, indent+1))
			sb.WriteString("\n")
		}
		sb.WriteString(prefix + "}")
		return sb.String()
	}

	return prefix
}

func (f *Formatter) formatExpression(expr uast.Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *uast.Literal:
		switch e.Kind {
		case uast.LiteralKindInt, uast.LiteralKindFloat:
			return e.Raw
		case uast.LiteralKindString:
			return e.Raw
		case uast.LiteralKindBool:
			if b, ok := e.Value.(bool); ok && b {
				return f.profile.True
			}
			return f.profile.False
		case uast.LiteralKindNull:
			return f.profile.Null
		}
		return "null"

	case *uast.Identifier:
		return e.Name

	case *uast.BinaryExpression:
		op := e.Operator
		// Adapt logical operator to active profile
		switch op {
		case "&&", "dan", "and":
			op = f.profile.And
		case "||", "atau", "or":
			op = f.profile.Or
		}
		return fmt.Sprintf("%s %s %s", f.formatExpression(e.Left), op, f.formatExpression(e.Right))

	case *uast.UnaryExpression:
		op := e.Operator
		if op == "!" || op == "bukan" || op == "tidak" || op == "not" {
			if f.profile.Not == "!" {
				return fmt.Sprintf("!%s", f.formatExpression(e.Argument))
			}
			return fmt.Sprintf("%s %s", f.profile.Not, f.formatExpression(e.Argument))
		}
		return fmt.Sprintf("%s%s", op, f.formatExpression(e.Argument))

	case *uast.CallExpression:
		calleeStr := f.formatExpression(e.Callee)
		// Adapt builtins: cetak <-> print
		if calleeStr == "cetak" && (f.profile.Name == "pythonic" || f.profile.Name == "classic") {
			calleeStr = "print"
		} else if calleeStr == "print" && f.profile.Name == "indonesian" {
			calleeStr = "cetak"
		}

		args := make([]string, len(e.Arguments))
		for i, a := range e.Arguments {
			args[i] = f.formatExpression(a)
		}
		return fmt.Sprintf("%s(%s)", calleeStr, strings.Join(args, ", "))

	case *uast.MemberExpression:
		objStr := f.formatExpression(e.Object)
		if e.Computed {
			return fmt.Sprintf("%s[%s]", objStr, e.Property)
		}
		return fmt.Sprintf("%s.%s", objStr, e.Property)

	case *uast.ArrayLiteral:
		elements := make([]string, len(e.Elements))
		for i, el := range e.Elements {
			elements[i] = f.formatExpression(el)
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case *uast.MapLiteral:
		if len(e.Properties) == 0 {
			return "{}"
		}
		pairs := make([]string, len(e.Properties))
		for i, p := range e.Properties {
			keyStr := f.formatExpression(p.Key)
			valStr := f.formatExpression(p.Value)
			pairs[i] = fmt.Sprintf("%s: %s", keyStr, valStr)
		}
		return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
	}

	return ""
}

func (f *Formatter) formatTypeHint(th string) string {
	if f.profile.Name == "indonesian" {
		th = strings.ReplaceAll(th, "|", " atau ")
	} else {
		th = strings.ReplaceAll(th, " atau ", "|")
	}
	return th
}

func isMajorBlock(stmt uast.Statement) bool {
	switch stmt.(type) {
	case *uast.FunctionDeclaration, *uast.IfStatement, *uast.WhileStatement, *uast.ForEachStatement, *uast.ConcurrencyStatement:
		return true
	default:
		return false
	}
}
