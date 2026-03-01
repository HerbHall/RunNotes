# RunNotes Architecture

## Overview

RunNotes is a Docker Desktop extension that allows users to attach notes and annotations to containers. It follows the standard Docker Desktop extension architecture with a React frontend, an optional backend service, and host-side persistence.

## Components

### Frontend (React UI)

The extension adds a tab to Docker Desktop's sidebar. The UI provides:

- A list of containers (via `docker ps` through the Extensions SDK)
- A note editor panel for the selected container
- Search across all notes
- An orphaned notes view for cleanup

The frontend uses React with Material UI to match Docker Desktop's native look and feel. The Extensions SDK (`@docker/extension-api-client`) provides access to Docker CLI commands and backend communication.

**Key constraint**: Docker Desktop reinitializes the extension UI every time the user navigates to the tab. All state must be fetched from the backend on mount — React state alone cannot persist.

### Backend (VM Service)

A Go service running inside the Docker Desktop VM. Handles:

- CRUD operations for notes
- SQLite storage management
- Container name/ID resolution

Communication between frontend and backend uses the Extensions SDK socket/named-pipe mechanism (not HTTP ports), avoiding port collisions with host applications.

### Storage Design

Notes are stored with a dual-key system:

- **Primary key**: Container name (or Compose project + service name)
- **Secondary key**: Container ID

This design ensures notes survive container recreation (`docker-compose down && up`), which is the most common workflow. When a container is recreated with the same name, its notes carry over automatically.

Storage uses **SQLite** — reliable, supports querying, single-file database. The database resides on a Docker volume mounted to the backend container, persisting across Docker Desktop restarts.

### Data Model

See [DATA_MODEL.md](DATA_MODEL.md) for the full schema definition, Go backend types, and TypeScript frontend types.

## Extension Metadata

The `metadata.json` file defines the extension for Docker Desktop:

- `ui.dashboard-tab` — Registers the RunNotes tab
- `vm.image` — Backend service image
- `host.binaries` — Optional host-side executables (for data export)

## Constraints

- Extension framework only supports **Linux containers** (even on Windows/Mac)
- Multi-arch support required: `linux/amd64` and `linux/arm64` minimum
- Extensions are distributed as Docker images via Docker Hub
- No native Windows container support in the extension framework

## References

- [Docker Extensions SDK](https://www.docker.com/developers/sdk/)
