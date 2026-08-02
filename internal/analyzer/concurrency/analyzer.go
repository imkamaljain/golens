package concurrency

import (
	"go/ast"
	"strings"

	"github.com/imkamaljain/golens/internal/core"
)

// Analyze checks for potential concurrency issues like unbounded goroutines in loops.
func Analyze(project *core.Project) []core.Finding {
	var findings []core.Finding

	for _, pkg := range project.Packages {
		if pkg.Raw == nil {
			continue
		}

		for _, file := range pkg.Raw.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				// We're looking for `go` statements inside loops
				switch loop := n.(type) {
				case *ast.ForStmt:
					findings = append(findings, checkLoopForGoroutines(loop.Body, pkg, project)...)
				case *ast.RangeStmt:
					findings = append(findings, checkLoopForGoroutines(loop.Body, pkg, project)...)
				}
				return true
			})
		}
	}

	return findings
}

func checkLoopForGoroutines(body *ast.BlockStmt, pkg *core.Package, project *core.Project) []core.Finding {
	var findings []core.Finding

	if body == nil {
		return findings
	}

	ast.Inspect(body, func(n ast.Node) bool {
		if goStmt, ok := n.(*ast.GoStmt); ok {
			pos := pkg.Raw.Fset.Position(goStmt.Pos())
			fileName := pos.Filename
			if strings.HasPrefix(fileName, project.Dir) {
				fileName = strings.TrimPrefix(fileName, project.Dir)
				fileName = strings.TrimPrefix(fileName, "/")
				fileName = strings.TrimPrefix(fileName, "\\")
			}

			findings = append(findings, core.Finding{
				Category: "Concurrency",
				Severity: "High",
				Message:  "Goroutine launched inside a loop. This may lead to unbounded concurrency or resource exhaustion.",
				File:     fileName,
				Line:     pos.Line,
			})
		}
		return true
	})

	return findings
}
