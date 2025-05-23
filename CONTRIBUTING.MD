# Contributing to hammerclock

This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming environment for everyone.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git

### Setting Up the Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:

   ```
   git clone https://github.com/your-username/hammerclock.git
   cd hammerclock
   ```

3. Add the original repository as an upstream remote:

   ```
   git remote add upstream https://github.com/itworks99/hammerclock.git
   ```

4. Install dependencies:

   ```
   go mod download
   ```

5. Build the project:

   ```
   go build -o bin/hammerclock.exe cmd/hammerclock/main.go
   ```

## Project Structure

The project follows this structure:

- `/cmd/hammerclock/main.go` - Application entry point
- `/internal/hammerclock/` - Core application code
  - `model.go` - Model initialization
  - `update.go` - Update logic
  - `view.go` - View rendering
  - `/common/` - Shared types and messages
  - `/config/` - Application configuration
  - `/logging/` - Game session logging
  - `/options/` - User options management
  - `/palette/` - Color theme definitions
  - `/ui/` - UI components

See [ARCHITECTURE.MD](ARCHITECTURE.MD) for more details about the application design.

## Development Workflow

### Branching Strategy

- `main`: Production-ready code
- `develop`: Development branch for integrating features
- Feature branches: Create from `develop` with the naming convention `feature/feature-name`
- Bug fix branches: Create from `develop` with the naming convention `fix/bug-name`

### Commits

- Use clear, descriptive commit messages
- Begin commit messages with a verb in the imperative mood (e.g., "Add", "Fix", "Update")
- Reference issue numbers in commit messages when applicable

### Pull Requests

1. Create a new branch from `develop`
2. Make your changes in small, logical commits
3. Push your branch to your fork
4. Create a pull request to the `develop` branch of the original repository
5. Include a clear description of the changes and any relevant issue numbers

## Coding Guidelines

### Architecture

See [ARCHITECTURE.MD](ARCHITECTURE.MD) for more details about the application design.

### Style Guide

- Follow standard Go style conventions
- Use `gofmt` to format your code
- Keep functions small and focused
- Write descriptive variable and function names
- Add comments for complex logic

### Documentation

- Update documentation when making significant changes
- Document public functions, types, and methods
- Include examples where appropriate
- Update the README.md when adding new features or changing existing functionality

## Testing

- Write tests for new features and bug fixes
- Run tests before creating pull requests:

  ```
  go test ./...
  ```

- Make sure deadlock tests pass:

  ```
  go test -v ./test/
  ```

## Building and Running

To build and run the application:

```
go build -o hammerclock.exe cmd/hammerclock/main.go
./hammerclock.exe
```

For development with live reloading, you can use the `air` tool:

1. Install air:

   ```
   go install github.com/cosmtrek/air@latest
   ```

2. Create a `.air.toml` configuration file in the project root with:

   ```toml
   root = "."
   tmp_dir = "tmp"

   [build]
   cmd = "go build -o ./tmp/main.exe ./cmd/hammerclock"
   bin = "tmp/main.exe"
   include_ext = ["go", "json"]
   exclude_dir = ["vendor", "bin"]
   ```

3. Run air:

   ```
   air
   ```

