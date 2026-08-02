# GoLens

> AI-powered performance intelligence and architectural analysis for Go applications.

GoLens is a modern, comprehensive health diagnostics platform for scaling Go backend services. It moves beyond standard linting by analyzing your Abstract Syntax Trees (AST) to expose hidden architectural anti-patterns, concurrency leaks, and database inefficiencies and visualizes everything in a premium, interactive HTML dashboard.

---

## 🚀 Features

GoLens currently ships with seven powerful static analyzers:

*   **Architecture Analyzer**: Automatically maps out the dependency graph of all your internal packages.
*   **Context Analyzer**: Detects the misuse of `context.Background()` or `context.TODO()` buried deep within business logic instead of being properly propagated.
*   **Concurrency Analyzer**: Spots dangerous unbounded goroutines launched inside `for` or `range` loops that can cause massive memory leaks.
*   **Database Analyzer**: Exposes the `SELECT *` anti-pattern and detects database queries executed inside loops (the classic N+1 query problem).
*   **GraphQL Analyzer**: Identifies direct database calls inside `gqlgen` resolvers that are missing `dataloader` implementations, preventing nested query explosion.
*   **API Analyzer**: Detects overly complex HTTP handlers that should be refactored into smaller, testable services.
*   **HTTP Analyzer**: Flags the use of bare `http.ListenAndServe`, which is vulnerable to slowloris attacks due to a lack of explicit timeouts.

---

## 🛠️ Installation

You can run GoLens directly from source:

```bash
git clone https://github.com/imkamaljain/golens.git
cd golens
go build -o golens ./cmd/golens
```

*(Once published, you will be able to install it globally via `go install github.com/imkamaljain/golens/cmd/golens@latest`)*

---

## 💻 Usage

To inspect a Go project, navigate to the root directory of your target service and run:

```bash
golens inspect .
```

### The Interactive Dashboard
GoLens will analyze your entire module in seconds and generate a `golens-report.html` file in the same directory. 

Open this file in any web browser to view your **Performance Dashboard**, featuring:
- An Executive **Health Score** and package summary.
- A **Cytoscape.js** interactive map of your architecture.
- A beautifully styled data grid of all identified **Findings**, color-coded by severity (High, Medium, Low) with exact file paths and line numbers.
