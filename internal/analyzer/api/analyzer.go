package api

import (
	"go/ast"

	"github.com/imkamaljain/golens/internal/core"
)

// Analyze checks for complex HTTP handlers.
func Analyze(project *core.Project) []core.Finding {
	var findings []core.Finding

	for _, pkg := range project.Packages {
		if pkg.Raw == nil {
			continue
		}

		for _, file := range pkg.Raw.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					// Check if it's an HTTP handler
					// func(w http.ResponseWriter, r *http.Request)
					if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) == 2 {
						// Simple heuristic: check if body has > 20 statements
						if funcDecl.Body != nil && len(funcDecl.Body.List) > 20 {
							pos := pkg.Raw.Fset.Position(funcDecl.Pos())
							findings = append(findings, core.Finding{
								Category: "API",
								Severity: "Medium",
								Message:  "HTTP handler is too complex (>20 statements). Consider refactoring business logic into separate services.",
								File:     core.RelativePath(pos.Filename, project.Dir),
								Line:     pos.Line,
							})
						}
					}
				}
				return true
			})
		}
	}

	return findings
}
