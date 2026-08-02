package graphql

import (
	"go/ast"
	"strings"

	"github.com/imkamaljain/golens/internal/core"
)

// Analyze checks for missing dataloaders in GraphQL resolvers (assuming 99designs/gqlgen).
func Analyze(project *core.Project) []core.Finding {
	var findings []core.Finding

	// Quick check: does any package import gqlgen?
	usesGqlgen := false
	for _, pkg := range project.Packages {
		for _, imp := range pkg.Imports {
			if strings.Contains(imp, "99designs/gqlgen") {
				usesGqlgen = true
				break
			}
		}
	}

	if !usesGqlgen {
		// If not using gqlgen, skip this analysis
		return findings
	}

	for _, pkg := range project.Packages {
		if pkg.Raw == nil {
			continue
		}

		for _, file := range pkg.Raw.Syntax {
			// Check if file is likely a resolver
			pos := pkg.Raw.Fset.Position(file.Pos())
			if !strings.Contains(pos.Filename, "resolver") {
				continue
			}

			// Check if dataloader is imported
			hasDataloader := false
			for _, imp := range file.Imports {
				if strings.Contains(imp.Path.Value, "dataloader") {
					hasDataloader = true
					break
				}
			}

			// If it's a resolver and has no dataloader, warn on any DB-like call inside a method
			if !hasDataloader {
				ast.Inspect(file, func(n ast.Node) bool {
					if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
						// It's a method (likely a resolver method)
						ast.Inspect(funcDecl.Body, func(innerNode ast.Node) bool {
							callExpr, ok := innerNode.(*ast.CallExpr)
							if !ok {
								return true
							}

							selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
							if !ok {
								return true
							}

							methodName := selExpr.Sel.Name
							if methodName == "Query" || methodName == "QueryRow" || methodName == "Find" {
								innerPos := pkg.Raw.Fset.Position(innerNode.Pos())
								findings = append(findings, core.Finding{
									Category: "GraphQL",
									Severity: "High",
									Message:  "Direct database call found in GraphQL resolver without a dataloader. This causes N+1 query performance issues.",
									File:     core.RelativePath(innerPos.Filename, project.Dir),
									Line:     innerPos.Line,
								})
							}
							return true
						})
					}
					return true
				})
			}
		}
	}

	return findings
}
