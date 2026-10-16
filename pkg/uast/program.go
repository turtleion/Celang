package uast

import (
	"encoding/json"
	"fmt"
)

// Program is the root node of a Celang U-AST compilation unit.
// This is the structure that is stored as data (.ast / .json) in version control.
type Program struct {
	NodeType    NodeType      `json:"type"`
	Version     string        `json:"version"`
	Metadata    *NodeMetadata `json:"metadata,omitempty"`
	Body        []Statement   `json:"body"`
	IsUniversal bool          `json:"is_universal,omitempty"`
}

func (p *Program) Type() NodeType          { return NodeTypeProgram }
func (p *Program) GetMetadata() *NodeMetadata { return p.Metadata }

func (p *Program) UnmarshalJSON(data []byte) error {
	var raw struct {
		NodeType    NodeType          `json:"type"`
		Version     string            `json:"version"`
		Metadata    *NodeMetadata     `json:"metadata,omitempty"`
		Body        []json.RawMessage `json:"body"`
		IsUniversal bool              `json:"is_universal,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	stmts := make([]Statement, len(raw.Body))
	for i, stmtRaw := range raw.Body {
		stmt, err := UnmarshalStatementJSON(stmtRaw)
		if err != nil {
			return fmt.Errorf("Program.body[%d] unmarshal error: %w", i, err)
		}
		stmts[i] = stmt
	}
	p.NodeType = NodeTypeProgram
	p.Version = raw.Version
	p.Metadata = raw.Metadata
	p.Body = stmts
	p.IsUniversal = raw.IsUniversal
	return nil
}
