package uast_test

import (
	"encoding/json"
	"testing"

	"celang/pkg/uast"
)

func TestUAST_JSONRoundTrip(t *testing.T) {
	// Construct a rich Celang U-AST program:
	//
	// let x = 10
	// x = 25
	// let prev = x.sebelumnya
	// if prev > 5 {
	//     cetak(prev + " ok")
	// }
	// sambil {
	//     hitung()
	// }
	originalProgram := uast.NewProgram(
		uast.NewVarDecl("x", uast.NewIntLiteral(10), true, "int"),
		uast.NewAssign(uast.NewIdentifier("x"), "=", uast.NewIntLiteral(25)),
		uast.NewVarDecl(
			"prev",
			uast.NewMemberExpr(uast.NewIdentifier("x"), "sebelumnya", false),
			false,
			"int",
		),
		uast.NewIf(
			uast.NewBinaryExpr(uast.NewIdentifier("prev"), ">", uast.NewIntLiteral(5)),
			uast.NewBlock(
				uast.NewExprStmt(
					uast.NewCallExpr(
						uast.NewIdentifier("cetak"),
						uast.NewBinaryExpr(uast.NewIdentifier("prev"), "+", uast.NewStringLiteral(" ok")),
					),
				),
			),
			nil,
		),
		uast.NewConcurrencyStmt(
			"sambil",
			uast.NewExprStmt(uast.NewCallExpr(uast.NewIdentifier("hitung"))),
		),
	)

	// Step 1: Serialize to JSON
	jsonData, err := uast.ToJSON(originalProgram, true)
	if err != nil {
		t.Fatalf("Failed to serialize U-AST to JSON: %v", err)
	}

	t.Logf("Generated U-AST JSON:\n%s", string(jsonData))

	// Step 2: Deserialize back from JSON
	restoredProgram, err := uast.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize U-AST from JSON: %v", err)
	}

	// Step 3: Re-serialize restored program and verify identical JSON
	reSerialized, err := uast.ToJSON(restoredProgram, true)
	if err != nil {
		t.Fatalf("Failed to re-serialize restored U-AST: %v", err)
	}

	if string(jsonData) != string(reSerialized) {
		t.Errorf("Roundtrip mismatch!\nOriginal:\n%s\n\nRestored:\n%s", string(jsonData), string(reSerialized))
	}

	// Verify statements count
	if len(restoredProgram.Body) != len(originalProgram.Body) {
		t.Fatalf("Expected %d statements, got %d", len(originalProgram.Body), len(restoredProgram.Body))
	}

	// Verify specific node details
	varDecl, ok := restoredProgram.Body[0].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Statement 0 is not *VariableDeclaration: %T", restoredProgram.Body[0])
	}
	if varDecl.Name != "x" || !varDecl.IsMutable {
		t.Errorf("Unexpected varDecl: %+v", varDecl)
	}
	litVal, ok := varDecl.Value.(*uast.Literal)
	if !ok || litVal.Value != int64(10) {
		t.Errorf("Unexpected literal value in varDecl: %+v", varDecl.Value)
	}
}

func TestUAST_DirectDeserialization(t *testing.T) {
	rawJSON := `{
		"type": "Program",
		"version": "1.0",
		"body": [
			{
				"type": "VariableDeclaration",
				"name": "pesan",
				"is_mutable": false,
				"value": {
					"type": "Literal",
					"kind": "string",
					"raw": "\"Halo Celang!\"",
					"value": "Halo Celang!"
				}
			},
			{
				"type": "ExpressionStatement",
				"expression": {
					"type": "CallExpression",
					"callee": {
						"type": "Identifier",
						"name": "cetak"
					},
					"arguments": [
						{
							"type": "Identifier",
							"name": "pesan"
						}
					]
				}
			}
		]
	}`

	var prog uast.Program
	err := json.Unmarshal([]byte(rawJSON), &prog)
	if err != nil {
		t.Fatalf("Direct unmarshal failed: %v", err)
	}

	if len(prog.Body) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(prog.Body))
	}

	exprStmt, ok := prog.Body[1].(*uast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected *ExpressionStatement, got %T", prog.Body[1])
	}

	callExpr, ok := exprStmt.Expression.(*uast.CallExpression)
	if !ok {
		t.Fatalf("Expected *CallExpression, got %T", exprStmt.Expression)
	}

	calleeId, ok := callExpr.Callee.(*uast.Identifier)
	if !ok || calleeId.Name != "cetak" {
		t.Fatalf("Expected callee 'cetak', got %+v", callExpr.Callee)
	}
}
