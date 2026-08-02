package db

import (
	"go/ast"
	"strings"

	"github.com/imkamaljain/golens/internal/core"
)

// Analyze checks for DB anti-patterns like SELECT * and queries inside loops.
func Analyze(project *core.Project) []core.Finding {
	var findings []core.Finding

	for _, pkg := range project.Packages {
		if pkg.Raw == nil {
			continue
		}

		for _, file := range pkg.Raw.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for SELECT *
				if basicLit, ok := n.(*ast.BasicLit); ok {
					if basicLit.Kind.String() == "STRING" {
						val := strings.ToUpper(basicLit.Value)
						if strings.Contains(val, "SELECT *") {
							pos := pkg.Raw.Fset.Position(n.Pos())
							findings = append(findings, core.Finding{
								Category: "Database",
								Severity: "Low",
								Message:  "Using 'SELECT *' is an anti-pattern. Specify explicit columns to avoid unexpected data and reduce memory usage.",
								File:     core.RelativePath(pos.Filename, project.Dir),
								Line:     pos.Line,
							})
						}
					}
				}

				// Check for DB queries inside loops
				switch loop := n.(type) {
				case *ast.ForStmt:
					findings = append(findings, checkLoopForQueries(loop.Body, pkg, project)...)
				case *ast.RangeStmt:
					findings = append(findings, checkLoopForQueries(loop.Body, pkg, project)...)
				}

				return true
			})
		}
	}

	return findings
}

func checkLoopForQueries(body *ast.BlockStmt, pkg *core.Package, project *core.Project) []core.Finding {
	var findings []core.Finding

	if body == nil {
		return findings
	}

	ast.Inspect(body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		methodName := selExpr.Sel.Name
		if methodName == "Query" || methodName == "QueryRow" || methodName == "Exec" || methodName == "Find" {
			pos := pkg.Raw.Fset.Position(n.Pos())
			findings = append(findings, core.Finding{
				Category: "Database",
				Severity: "High",
				Message:  "Database query inside a loop (N+1 query risk). Consider batching the query outside the loop.",
				File:     core.RelativePath(pos.Filename, project.Dir),
				Line:     pos.Line,
			})
		}

		return true
	})

	return findings
}
