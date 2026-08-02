package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/imkamaljain/golens/internal/analyzer/architecture"
	"github.com/imkamaljain/golens/internal/core"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoLens Report</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.28.1/cytoscape.min.js"></script>
    <script src="https://unpkg.com/lucide@latest"></script>
    <style>
        body { font-family: 'Inter', sans-serif; background-color: #f3f4f6; }
        #cy { width: 100%; height: 500px; }
        .severity-high { background-color: #fee2e2; color: #991b1b; border-left-color: #ef4444; }
        .severity-medium { background-color: #fef9c3; color: #854d0e; border-left-color: #eab308; }
        .severity-low { background-color: #e0f2fe; color: #075985; border-left-color: #0ea5e9; }
        .health-excellent { color: #16a34a; }
        .health-fair { color: #ca8a04; }
        .health-poor { color: #dc2626; }
    </style>
</head>
<body class="text-gray-800">

<div class="min-h-screen flex flex-col">
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-10">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
                <div class="bg-blue-600 p-2 rounded-lg shadow-sm">
                    <i data-lucide="stethoscope" class="text-white w-6 h-6"></i>
                </div>
                <h1 class="text-2xl font-bold text-gray-900 tracking-tight">GoLens</h1>
            </div>
            <div class="text-sm font-medium text-gray-500 bg-gray-100 px-3 py-1 rounded-full border border-gray-200">
                Go {{.ReportData.Project.GoVersion}} &bull; {{.ReportData.Project.ModulePath}}
            </div>
        </div>
    </header>

    <!-- Main Content -->
    <main class="flex-grow max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full space-y-8">
        
        <!-- Scorecards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 flex flex-col items-center justify-center transition hover:shadow-md">
                <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Health Score</h3>
                <div class="text-5xl font-bold {{if ge .HealthScore 90}}health-excellent{{else if ge .HealthScore 70}}health-fair{{else}}health-poor{{end}}">
                    {{.HealthScore}}
                </div>
                <p class="text-xs text-gray-400 mt-2">out of 100</p>
            </div>
            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 flex flex-col items-center justify-center transition hover:shadow-md">
                <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Packages Analyzed</h3>
                <div class="text-5xl font-bold text-blue-600">
                    {{.TotalPackages}}
                </div>
                <p class="text-xs text-gray-400 mt-2">modules scanned</p>
            </div>
            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 flex flex-col items-center justify-center transition hover:shadow-md">
                <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-2">Issues Found</h3>
                <div class="text-5xl font-bold {{if eq .TotalIssues 0}}text-green-500{{else}}text-red-500{{end}}">
                    {{.TotalIssues}}
                </div>
                <p class="text-xs text-gray-400 mt-2">requires attention</p>
            </div>
        </div>

        <!-- Architecture Graph -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200 bg-gray-50 flex items-center space-x-2">
                <i data-lucide="share-2" class="w-5 h-5 text-gray-500"></i>
                <h2 class="text-lg font-semibold text-gray-800">Architecture Dependency Graph</h2>
            </div>
            <div id="cy" class="bg-slate-50 cursor-grab active:cursor-grabbing"></div>
        </div>

        <!-- Findings -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200 bg-gray-50 flex items-center space-x-2">
                <i data-lucide="alert-triangle" class="w-5 h-5 text-gray-500"></i>
                <h2 class="text-lg font-semibold text-gray-800">Diagnostic Findings</h2>
            </div>
            <div class="p-6">
                {{if eq .TotalIssues 0}}
                    <div class="text-center py-12">
                        <i data-lucide="check-circle" class="w-16 h-16 text-green-500 mx-auto mb-4 opacity-50"></i>
                        <h3 class="text-lg font-medium text-gray-900">Perfect Health!</h3>
                        <p class="text-gray-500 mt-1">No anti-patterns or issues were detected.</p>
                    </div>
                {{else}}
                    <div class="space-y-4">
                        {{range .ReportData.Findings}}
                            <div class="border-l-4 rounded-r-lg p-4 shadow-sm transition hover:shadow-md {{if eq .Severity "High"}}severity-high{{else if eq .Severity "Medium"}}severity-medium{{else}}severity-low{{end}}">
                                <div class="flex justify-between items-start">
                                    <div>
                                        <div class="flex items-center space-x-2 mb-1">
                                            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-white bg-opacity-50">
                                                {{.Category}}
                                            </span>
                                            <span class="text-xs font-bold uppercase">{{.Severity}}</span>
                                        </div>
                                        <h4 class="text-base font-semibold mb-1">{{.Message}}</h4>
                                    </div>
                                </div>
                                <div class="mt-2 flex items-center text-sm font-mono opacity-80">
                                    <i data-lucide="file-code" class="w-4 h-4 mr-1"></i>
                                    {{.File}}:{{.Line}}
                                </div>
                            </div>
                        {{end}}
                    </div>
                {{end}}
            </div>
        </div>
    </main>

    <footer class="bg-white border-t border-gray-200 py-6 mt-auto">
        <div class="max-w-7xl mx-auto px-4 text-center text-sm text-gray-500">
            Generated by GoLens &bull; AI-Powered Performance Intelligence
        </div>
    </footer>
</div>

<script>
    // Initialize icons
    lucide.createIcons();

    // Render Architecture Graph
    const rawData = {{.JSONData}};
    const elements = [];
    
    rawData.nodes.forEach(n => {
        elements.push({ data: { id: n.id, name: n.name } });
    });
    
    rawData.edges.forEach(e => {
        elements.push({ data: { source: e.source, target: e.target } });
    });

    var cy = cytoscape({
        container: document.getElementById('cy'),
        elements: elements,
        style: [
            {
                selector: 'node',
                style: {
                    'background-color': '#2563eb', // blue-600
                    'label': 'data(name)',
                    'color': '#1f2937', // gray-800
                    'text-valign': 'bottom',
                    'text-halign': 'center',
                    'text-margin-y': 6,
                    'font-size': '12px',
                    'font-family': 'Inter',
                    'font-weight': '600',
                    'width': 30,
                    'height': 30,
                    'border-width': 2,
                    'border-color': '#bfdbfe' // blue-200
                }
            },
            {
                selector: 'edge',
                style: {
                    'width': 2,
                    'line-color': '#cbd5e1', // slate-300
                    'target-arrow-color': '#cbd5e1',
                    'target-arrow-shape': 'triangle',
                    'curve-style': 'bezier'
                }
            }
        ],
        layout: {
            name: 'cose',
            padding: 50,
            nodeRepulsion: 400000,
            idealEdgeLength: 100,
            edgeElasticity: 100,
            nestingFactor: 5,
            gravity: 250,
            numIter: 1000,
            initialTemp: 200,
            coolingFactor: 0.95,
            minTemp: 1.0,
            animate: false
        }
    });
</script>
</body>
</html>`

type TemplateData struct {
	JSONData      string
	ReportData    *core.Report
	HealthScore   int
	TotalPackages int
	TotalIssues   int
}

// GenerateHTML renders the architecture report to an HTML file.
func GenerateHTML(archReport *architecture.ArchitectureReport, fullReport *core.Report, outputDir string) error {
	jsonData, err := json.Marshal(archReport)
	if err != nil {
		return fmt.Errorf("failed to serialize report data: %w", err)
	}

	// Calculate metrics
	totalPackages := len(fullReport.Project.Packages)
	totalIssues := len(fullReport.Findings)

	// Calculate Health Score (start 100, subtract points based on severity)
	healthScore := 100
	for _, f := range fullReport.Findings {
		switch f.Severity {
		case "High":
			healthScore -= 10
		case "Medium":
			healthScore -= 5
		case "Low":
			healthScore -= 2
		}
	}
	if healthScore < 0 {
		healthScore = 0
	}

	outPath := filepath.Join(outputDir, "golens-report.html")
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := TemplateData{
		JSONData:      string(jsonData),
		ReportData:    fullReport,
		HealthScore:   healthScore,
		TotalPackages: totalPackages,
		TotalIssues:   totalIssues,
	}

	return tmpl.Execute(f, data)
}
