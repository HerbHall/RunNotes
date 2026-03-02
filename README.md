# RunNotes

A Docker Desktop extension for attaching notes and annotations to containers. Never forget why a container exists, what config quirks it has, or what you were testing.

## The Problem

Docker Desktop shows container name, image, status, ports, and logs — but there's no built-in way to capture *why* a container exists, what it's for, or what you were doing with it. RunNotes fills that gap.

## Features

- **Container notes** — Attach rich text notes to any container with tag-based organization
- **Markdown support** — Write notes in Markdown with live preview toggle
- **Name-based persistence** — Notes survive `docker-compose down && up` cycles (keyed by name, not ephemeral ID)
- **Lifecycle management** — Container IDs auto-reconcile when containers are recreated
- **Search and filter** — Find containers by name, content, or tags; filter by pinned notes
- **Pin notes** — Pin important notes for quick access
- **Orphan management** — Detect and clean up notes for containers that no longer exist
- **Export/Import** — Back up all notes as JSON, restore from backup

## Status

**In Development** — Feature-complete with container notes, markdown, search, export/import, and lifecycle management.

## Architecture

RunNotes is a Docker Desktop extension with three components:

- **Frontend** — React UI tab in Docker Desktop (Material UI to match DD's look)
- **Backend** — Go service running in the Desktop VM (SQLite storage)
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
├── cmd/backend/         - Go entry point
├── internal/            - Backend packages (database, handler, models, store)
├── docs/                - Architecture and research documentation
├── ui/                  - React frontend (Vite + MUI + TypeScript)
├── metadata.json        - Docker extension metadata
├── docker-compose.yaml  - VM service definition
├── Dockerfile           - Multi-stage extension image build
├── Makefile             - Build automation
└── README.md            - This file
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE) — Copyright (c) 2026 Herb Hall
