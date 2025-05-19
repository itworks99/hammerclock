# Hammerclock Development Guidelines

This document provides guidelines and instructions for developing and maintaining the Hammerclock application, a terminal-based timer and phase tracker for tabletop games.

Hammerclock is a Go-based terminal user interface (TUI) application built with the tview framework, following the Model-View-Update (MVU) architecture. The code should be idiomatic, modular, testable, and aligned with modern Go and MVU best practices.

## Role Expectations

- Enforce strict MVU separation: **Model**, **Update**, and **View** logic must remain isolated.
- Help write idiomatic, maintainable, and testable Go code.
- Favor **small, focused interfaces** and **explicit dependency injection**.
- Ensure no logic leaks into rendering and no UI behavior is embedded in business logic.

## MVU Guidelines

- **Model**: Clean, plain Go struct representing app state. Model is defined in `internal/hammerclock/common/types.go`.
- **Update**: A pure function: `Update(msg, model) → (model, cmd)`. No side effects outside commands. Defined in `internal/hammerclock/update.go`.
- **View**: A component that renders UI based on the model. No state mutations or side effects. Defined in `internal/hammerclock/view.go`.
- **Commands**: Functions that return messages to be processed by the update function. Used for side effects.
- Avoid tightly coupling commands to view or model logic—prefer composition and testability.

## Current Project Structure

- `/cmd/hammerclock/main.go` - App entry point
- `/internal/hammerclock/`
  - `model.go` - Model initialization
  - `update.go` - Update logic
  - `view.go` - View rendering
  - `/common/`
    - `messages.go` - Message type definitions
    - `types.go` - Core type definitions including Model
  - `/config/` - Application configuration
  - `/logging/` - Game session logging
  - `/options/` - User options management
  - `/palette/` - Color theme definitions
  - `/rules/` - Game rule definitions
  - `/ui/` - UI components
- `/test/` - Test files
- `ARCHITECTURE.MD` - Architecture documentation
- `CONTRIBUTING.MD` - Contribution guidelines
- `README.MD` - Project documentation

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

## Logging and Observability

Hammerclock includes a built-in logging system that:

- Records game events and player actions
- Writes logs to a CSV file (`logs.csv`)
- Uses a buffered channel for non-blocking logging
- Supports both in-memory action logs (for UI display) and persistent logs (for post-game analysis)

When adding new log entries:

```go
// Add a log entry for a specific player
logging.AddLogEntry(player, model, "Player %s completed phase %s", player.Name, phaseName)
```

For application observability:

- Consider adding metrics for startup time, render latency, or error counts
- Structured logging would be a useful future enhancement

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

- Go 1.18.0 or later
- The tview library for terminal UI (already included in go.mod)

### Building the Application

1. Clone the repository
2. Navigate to the project root directory
3. Build the application using:

   ```powershell
   go build -o bin/hammerclock.exe cmd/hammerclock/main.go
   ```

4. Run the application:

   ```powershell
   ./bin/hammerclock.exe
   ```

### Live build and reload

For live reload during development, use the `air` tool:

1. Install `air`:

   ```powershell
   go install github.com/cosmtrek/air@latest
   ```

2. Create a `.air.toml` configuration file in the project root:

   ```toml
   root = "."
   tmp_dir = "tmp"

   [build]
   cmd = "go build -o ./tmp/main.exe ./cmd/hammerclock"
   bin = "tmp/main.exe"
   include_ext = ["go", "json"]
   exclude_dir = ["vendor", "bin"]
   ```

3. Run `air` in the project root directory:

   ```powershell
   air
   ```

### Terminal UI Implementation

The application uses the tview library for terminal UI. Key aspects:

- Screen initialization and cleanup in the main function (`cmd/hammerclock/main.go`)
- Event handling through the message system (`internal/hammerclock/update.go`)
- Time tracking using Go's ticker functionality
- UI components organized in `internal/hammerclock/ui/`

### Message System

When creating new message types:

1. Define the message struct in `internal/hammerclock/common/messages.go`
2. Add a handler in the `Update` function in `internal/hammerclock/update.go`
3. Implement the handler function following the MVU pattern

Example:

```go
// In messages.go
type NewFeatureMsg struct {
    Parameter string
}

// In update.go
func Update(msg common.Message, model common.Model) (common.Model, Command) {
    switch msg := msg.(type) {
    // ...existing message handlers...
    case *common.NewFeatureMsg:
        return handleNewFeature(msg, model)
    }
}

func handleNewFeature(msg *common.NewFeatureMsg, model common.Model) (common.Model, Command) {
    newModel := model // Create a copy
    // Update model based on message
    return newModel, noCommand
}
```

### Adding New Features

When adding new features:

1. For UI changes:
   - Add new components in the `ui` directory
   - Update the view rendering in `view.go`
2. For game logic:
   - Define new message types
   - Implement handlers in `update.go`
   - Update the model as needed
3. For configuration changes:
   - Update structures in `options/options.go`
   - Ensure backward compatibility with existing config files
4. Add appropriate tests for the new functionality

### Debugging

- Use the logging system for debug output (`internal/hammerclock/logging/logging.go`)
- Add log entries to track state changes and user actions
- Use `fmt.Printf()` for temporary debug output during development
- Test UI changes incrementally to avoid breaking the interface

## UI Components

Hammerclock uses several UI components organized in the `/internal/hammerclock/ui/` directory:

| Component    | Purpose                            | File              |
| ------------ | ---------------------------------- | ----------------- |
| AboutPanel   | Displays application information   | `AboutPanel.go`   |
| Clock        | Displays and manages time tracking | `clock.go`        |
| LogPanel     | Shows player action logs           | `LogPanel.go`     |
| MenuBar      | Provides navigation controls       | `MenuBar.go`      |
| OptionsPanel | Manages game settings              | `OptionsPanel.go` |
| PlayerPanel  | Displays player information        | `PlayerPanel.go`  |
| StatusPanel  | Shows game status                  | `StatusPanel.go`  |

When creating or modifying UI components:

1. Follow the existing component patterns
2. Keep drawing logic separate from state management
3. Update the View's Render method to use your component

## Color Palettes

Hammerclock supports multiple color palettes, defined in `internal/hammerclock/palette/palette.go`:

- `k9s`: Default color scheme based on the K9s terminal UI application
- `dracula`: Dark theme with purple accents
- `monokai`: Dark theme with vibrant colors
- `warhammer`: Theme inspired by Warhammer 40K colors
- `killteam`: Theme inspired by Kill Team colors

When adding a new color palette:

1. Define a new ColorPalette struct in the palette package
2. Add the palette name to the `ColorPalettes()` function
3. Add a case for your palette in the `ColorPaletteByName()` function

## Future Improvements

- Introduce more granular message types for specific state changes
- Implement deeper immutability for nested state objects
- Add comprehensive test coverage with unit and integration tests
- Enhance logging with filtering and search capabilities
- Support for exporting logs in multiple formats
- Add custom timers and time limits for competitive play
- Implement network play support for remote players
- Create a headless mode for tournament use
