package scanner

import (
	"fmt"
	"strings"

	"github.com/imkamaljain/golens/internal/core"
	"golang.org/x/tools/go/packages"
)

// ScanProject analyzes the Go project at the given directory path.
func ScanProject(dir string) (*core.Project, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedModule | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  dir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in %s", dir)
	}

	project := &core.Project{
		Dir:      dir,
		Packages: make(map[string]*core.Package),
	}

	for _, p := range pkgs {
		if p.Module != nil {
			project.ModulePath = p.Module.Path
			project.GoVersion = p.Module.GoVersion
		}

		// Only include project's own packages, ignore standard library and external deps
		// For simplicity in v0.1, we consider everything returned by "./..." as a project package.
		if !strings.HasPrefix(p.ID, project.ModulePath) && project.ModulePath != "" {
			continue
		}

		var imports []string
		for impPath := range p.Imports {
			imports = append(imports, impPath)
		}

		project.Packages[p.ID] = &core.Package{
			ID:      p.ID,
			Name:    p.Name,
			Path:    p.PkgPath,
			Imports: imports,
			Raw:     p,
		}
	}

	return project, nil
}
