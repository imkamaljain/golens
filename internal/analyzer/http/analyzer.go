package http

import (
	"go/ast"
	"strings"

	"github.com/imkamaljain/golens/internal/core"
)

// Analyze checks for common HTTP server misconfigurations (like using http.ListenAndServe without timeouts).
func Analyze(project *core.Project) []core.Finding {
	var findings []core.Finding

	for _, pkg := range project.Packages {
		if pkg.Raw == nil {
			continue
		}

		for _, file := range pkg.Raw.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selExpr.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == "http" && selExpr.Sel.Name == "ListenAndServe" {
					pos := pkg.Raw.Fset.Position(n.Pos())
					fileName := pos.Filename
					if strings.HasPrefix(fileName, project.Dir) {
						fileName = strings.TrimPrefix(fileName, project.Dir)
						fileName = strings.TrimPrefix(fileName, "/")
						fileName = strings.TrimPrefix(fileName, "\\")
					}

					findings = append(findings, core.Finding{
						Category: "HTTP",
						Severity: "Medium",
						Message:  "Using http.ListenAndServe directly does not allow setting timeouts. Use a custom http.Server instead to prevent slowloris attacks.",
						File:     fileName,
						Line:     pos.Line,
					})
				}

				return true
			})
		}
	}

	return findings
}
