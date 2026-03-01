# RunNotes

A Docker Desktop extension for attaching notes and annotations to containers. Never forget why a container exists, what config quirks it has, or what you were testing.

## The Problem

Docker Desktop shows container name, image, status, ports, and logs — but there's no built-in way to capture *why* a container exists, what it's for, or what you were doing with it. RunNotes fills that gap.

## Features (Planned)

- **Container notes** — Attach text/markdown notes to any container
- **Name-based persistence** — Notes survive container recreation (keyed by name, not ephemeral ID)
- **Search** — Find notes across all containers
- **Orphan management** — View and clean up notes for containers that no longer exist
- **Export** — Back up your notes as JSON

## Status

🚧 **Pre-development** — Architecture and planning phase.

## Architecture

RunNotes is a Docker Desktop extension with three components:

- **Frontend** — React UI tab in Docker Desktop (Material UI to match DD's look)
- **Backend** — Lightweight service running in the Desktop VM (SQLite or JSON storage)
- **Host binary** — Optional, for host-side data persistence

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for full technical details.

## Development

### Prerequisites

- Docker Desktop 4.8.0+ with Extensions enabled
- Node.js 18+ (for frontend development)
- Docker Extensions CLI (`docker extension` command)

### Getting Started

```powershell
# Clone
git clone https://github.com/HerbHall/RunNotes.git
cd RunNotes

# Build and install locally
make build-extension
make install-extension
```

### Project Structure

```text
RunNotes/
├── docs/                - Architecture and research documentation
├── ui/                  - React frontend (Docker Desktop extension tab)
├── backend/             - Backend service (runs in Desktop VM)
├── metadata.json        - Docker extension metadata
├── Dockerfile           - Extension image build
├── Makefile             - Build automation
├── CLAUDE.md            - Project context for Claude Code
├── CONTRIBUTING.md      - Contribution guidelines
└── README.md            - This file
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE) — Copyright (c) 2026 Herb Hall
