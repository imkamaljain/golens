package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/imkamaljain/golens/internal/analyzer/api"
	"github.com/imkamaljain/golens/internal/analyzer/architecture"
	"github.com/imkamaljain/golens/internal/analyzer/concurrency"
	"github.com/imkamaljain/golens/internal/analyzer/context"
	"github.com/imkamaljain/golens/internal/analyzer/db"
	"github.com/imkamaljain/golens/internal/analyzer/graphql"
	httparch "github.com/imkamaljain/golens/internal/analyzer/http"
	"github.com/imkamaljain/golens/internal/core"
	"github.com/imkamaljain/golens/internal/report"
	"github.com/imkamaljain/golens/internal/scanner"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [path]",
	Short: "Inspect a Go project for architectural and structural issues",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Scanning project at: %s\n", absPath)

		proj, err := scanner.ScanProject(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Detected Go module: %s (Go %s)\n", proj.ModulePath, proj.GoVersion)
		fmt.Printf("Found %d packages.\n", len(proj.Packages))

		fmt.Println("Analyzing architecture...")
		archReport := architecture.Analyze(proj)
		fmt.Printf("Generated graph with %d nodes and %d edges.\n", len(archReport.Nodes), len(archReport.Edges))

		fmt.Println("Running context analyzer...")
		ctxFindings := context.Analyze(proj)

		fmt.Println("Running concurrency analyzer...")
		concFindings := concurrency.Analyze(proj)

		fmt.Println("Running HTTP analyzer...")
		httpFindings := httparch.Analyze(proj)

		fmt.Println("Running Database analyzer...")
		dbFindings := db.Analyze(proj)

		fmt.Println("Running GraphQL analyzer...")
		gqlFindings := graphql.Analyze(proj)

		fmt.Println("Running API analyzer...")
		apiFindings := api.Analyze(proj)

		var allFindings []core.Finding
		allFindings = append(allFindings, ctxFindings...)
		allFindings = append(allFindings, concFindings...)
		allFindings = append(allFindings, httpFindings...)
		allFindings = append(allFindings, dbFindings...)
		allFindings = append(allFindings, gqlFindings...)
		allFindings = append(allFindings, apiFindings...)
		fmt.Printf("Found %d total issues.\n", len(allFindings))

		finalReport := &core.Report{
			Project:  proj,
			Findings: allFindings,
		}

		fmt.Println("Generating HTML report...")
		err = report.GenerateHTML(archReport, finalReport, ".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate report: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Done! Report saved to golens-report.html")
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
