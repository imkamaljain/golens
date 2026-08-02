package architecture

import (
	"strings"

	"github.com/imkamaljain/golens/internal/core"
)

// GraphNode represents a package in the architecture graph.
type GraphNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GraphEdge represents a dependency between two packages.
type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// ArchitectureReport holds the results of the architecture analysis.
type ArchitectureReport struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// Analyze processes the project to build an architecture dependency graph.
func Analyze(project *core.Project) *ArchitectureReport {
	report := &ArchitectureReport{
		Nodes: make([]GraphNode, 0),
		Edges: make([]GraphEdge, 0),
	}

	for id, pkg := range project.Packages {
		report.Nodes = append(report.Nodes, GraphNode{
			ID:   id,
			Name: pkg.Name,
		})

		for _, imp := range pkg.Imports {
			// Only include edges to internal project packages for the diagram
			if _, exists := project.Packages[imp]; exists {
				report.Edges = append(report.Edges, GraphEdge{
					Source: id,
					Target: imp,
				})
			} else if project.ModulePath != "" && strings.HasPrefix(imp, project.ModulePath) {
				// If it's part of the module but wasn't loaded properly for some reason
				report.Edges = append(report.Edges, GraphEdge{
					Source: id,
					Target: imp,
				})
			}
		}
	}

	return report
}
