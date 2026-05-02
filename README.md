# GitHub Actions Manager

A premium Terminal User Interface (TUI) for managing self-hosted GitHub Actions runners on Linux. Built with Go and the Bubble Tea framework, this tool provides a streamlined experience for installing, configuring, and managing the lifecycle of your runners.

## Key Features

- **Effortless Installation**: Interactive flow to download and configure runners with live progress tracking.
- **Service Management**: Easily install, start, stop, and uninstall runners as systemd services.
- **Robust Removal**: Choose between standard removal (token-based) or Force Mode for deep cleanup.
- **Premium UI/UX**: Beautiful, modular TUI design with smooth transitions and intelligent navigation.
- **Developer Friendly**: Highly modular architecture (MVU pattern), zero duplicate logic, and fully linted codebase.

## Quick Start

### Prerequisites

- Go 1.21 or later
- golangci-lint (optional, for development)
- Linux environment (required for runner services)

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/AzPepoze/github-actions-manager.git
cd github-actions-manager
make build
```

The binary will be available at ./bin/github-actions-manager.

### Usage

Run the manager:

```bash
make start
```

## Development

This project uses a modern Go stack:
- **Framework**: Bubble Tea (MVU architecture)
- **Styling**: Lip Gloss
- **Components**: Bubbles

### Makefile Commands

- `make build`: Compiles the binary to ./bin.
- `make start`: Builds and runs the application.
- `make lint`: Runs golangci-lint for code quality checks.
- `make clean`: Removes build artifacts.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---
Built by AzPepoze
