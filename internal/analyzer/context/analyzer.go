package context

import (
	"go/ast"
	"strings"

	"github.com/imkamaljain/golens/internal/core"
)

// Analyze checks for context.Background() or context.TODO() misuse.
func Analyze(project *core.Project) []core.Finding {
	var findings []core.Finding

	for _, pkg := range project.Packages {
		if pkg.Raw == nil {
			continue
		}

		for _, file := range pkg.Raw.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				// We want to skip inspecting inside `main` or `init` as it's often OK to use context.Background there.
				if fd, ok := n.(*ast.FuncDecl); ok {
					if fd.Name.Name == "main" || fd.Name.Name == "init" {
						return false
					}
				}

				// Look for CallExpr
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Look for SelectorExpr like `context.Background`
				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selExpr.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == "context" {
					pos := pkg.Raw.Fset.Position(n.Pos())
					// Strip the dir path for cleaner output if desired, here just use what Fset gives
					fileName := pos.Filename
					if strings.HasPrefix(fileName, project.Dir) {
						fileName = strings.TrimPrefix(fileName, project.Dir)
						fileName = strings.TrimPrefix(fileName, "/")
						fileName = strings.TrimPrefix(fileName, "\\")
					}

					if selExpr.Sel.Name == "Background" {
						findings = append(findings, core.Finding{
							Category: "Context",
							Severity: "Medium",
							Message:  "context.Background() used outside of main/init. Contexts should typically be passed in.",
							File:     fileName,
							Line:     pos.Line,
						})
					} else if selExpr.Sel.Name == "TODO" {
						findings = append(findings, core.Finding{
							Category: "Context",
							Severity: "Low",
							Message:  "context.TODO() used. This should be replaced with a real context.",
							File:     fileName,
							Line:     pos.Line,
						})
					}
				}

				return true
			})
		}
	}

	return findings
}
