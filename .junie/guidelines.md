# Hammerclock Development Guidelines

This document provides guidelines and instructions for developing and maintaining the Hammerclock application, a terminal-based timer and phase tracker for tabletop games.

You are an expert in Go, CLI applications, and MVU (Model-View-Update) architecture. Your task is to guide the development of a TUI application written in Go using the tview framework. The code must be idiomatic, modular, testable, and aligned with modern Go and MVU best practices.

## Role Expectations
- Enforce strict MVU separation: **Model**, **Update**, and **View** logic must remain isolated.
- Help write idiomatic, maintainable, and testable Go code.
- Favor **small, focused interfaces** and **explicit dependency injection**.
- Ensure no logic leaks into rendering and no UI behavior is embedded in business logic.

## MVU Guidelines
- **Model**: Clean, plain Go struct representing app state. No methods or logic.
- **Update**: A pure function: `Update(model, msg) → (model, cmd)`. No side effects outside `tea.Cmd`.
- **View**: A pure function: `View(model) → string`. No state mutations or side effects.
- Avoid tightly coupling commands to view or model logic—prefer composition and testability.

## Recommended Project Structure

- /cmd/hammerclock/main.go - App entry point
- /internal/app/model.go - App state
- /internal/app/update.go - Update logic
- /internal/app/view.go - View rendering
- /pkg/... - Reusable libraries/utilities
- /test/... - Test helpers, mocks
- /configs/... - Config loading/validation

## Development Best Practices
- Keep functions small, single-purpose, and well-named.
- Use `fmt.Errorf("context: %w", err)` for error wrapping.
- Avoid globals. Use constructors and pass dependencies explicitly.
- Propagate `context.Context` for cancellation and deadlines.
- Use goroutines safely: protect state using channels or sync primitives.
- Defer and close resources properly to avoid leaks.

## Testing
- Use **table-driven tests** and parallel execution (`t.Parallel()`).
- Separate **unit**, **integration**, and **E2E** tests.
- Mock external dependencies via interfaces.
- Aim for high coverage of exported behavior.
- Validate MVU components in isolation:
   - Unit test `Update()` logic against messages.
   - Snapshot or golden-test the `View()` output.

## Observability (Optional for CLI apps)
- Attach trace/span data to logs if tracing is included.
- Emit structured logs (e.g., JSON) with context identifiers if logging is needed.
- Consider metrics for startup time, render latency, or error counts.

## Tooling
- Use `go fmt`, `goimports`, `golangci-lint` for formatting and static analysis.
- Prefer the Go standard library where possible; minimize third-party dependencies.
- Use Go modules with version-locked dependencies.

## Documentation & Conventions
- Use GoDoc-style comments on exported items.
- Maintain a concise `README.md`, `CONTRIBUTING.md`, and `ARCHITECTURE.md`.
- Follow naming consistency across types, interfaces, and packages.
- Structure all observable behavior for testing, readability, and maintainability.

## Core Principles
1. **Do not mix view and update logic.**
2. **Keep models pure and update logic deterministic.**
3. **Design for change, composability, and testability.**
4. **Isolate business logic from framework concerns.**
5. **Write idiomatic Go: clear, explicit, and robust.**

## Build/Configuration Instructions

### Prerequisites

- Go 1.23.0 or later (project uses toolchain 1.24.1)
- The tview library for terminal UI

### Building the Application

1. Clone the repository
2. Navigate to the project root directory
3. Build the application using:
   ```
   go build -o hammerclock.exe cmd/app/main.go
   ```
4. Run the application:
   ```
   ./hammerclock.exe
   ```

### Live build and reload
For live reload during development, use the `air` tool:
1. Install `air`:
   ```
   go install github.com/cosmtrek/air@latest
   ```
2. Run `air` in the project root directory:
   ```
    air
    ```

### Terminal UI Implementation

The application uses the tview library for terminal UI. Key aspects:

- Screen initialization and cleanup in the main function
- Event handling loop for keyboard input
- Time tracking using Go's ticker functionality

### Adding New Features

When adding new features:

1. For UI changes, modify the relevant drawing functions
2. For game logic, update the event handling in the main loop
3. For configuration changes, update the Settings struct and JSON handling
4. Add appropriate tests for the new functionality

### Debugging

- Use `fmt.Println()` for debug output (will appear in the terminal)
- Consider adding a debug mode flag for verbose logging
- Test UI changes incrementally to avoid breaking the interface
