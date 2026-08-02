package core

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

// Project represents a Go project being analyzed.
type Project struct {
	Dir        string
	ModulePath string
	GoVersion  string
	Packages   map[string]*Package
}

// Package represents a Go package within the project.
type Package struct {
	ID      string
	Name    string
	Path    string
	Imports []string
	Raw     *packages.Package
}

// Finding represents an issue discovered by an analyzer.
type Finding struct {
	Category string `json:"category"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Report represents the final analysis report.
type Report struct {
	Project  *Project
	Findings []Finding `json:"findings"`
}

// RelativePath strips the project directory from the absolute file path.
func RelativePath(filePath, projectDir string) string {
	if strings.HasPrefix(filePath, projectDir) {
		filePath = strings.TrimPrefix(filePath, projectDir)
		filePath = strings.TrimPrefix(filePath, "/")
		filePath = strings.TrimPrefix(filePath, "\\")
	}
	return filePath
}

