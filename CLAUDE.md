# RunNotes

Docker Desktop extension for attaching notes and annotations to containers.

## Tech Stack

- **Frontend**: React + Material UI (Docker Desktop extension standard)
- **Backend**: Go + `modernc.org/sqlite` (pure Go, no CGO) running in Desktop VM
- **Storage**: SQLite on Docker volume (see [docs/DATA_MODEL.md](docs/DATA_MODEL.md))
- **Migrations**: Goose v3 with `go:embed` (baked into binary)
- **Build**: Docker Extensions CLI + Makefile
- **Platform**: Docker Desktop 4.8.0+ (Windows, Mac, Linux)

## Key Design Decisions

- Notes are keyed by **container name** (primary) and **container ID** (secondary) so they survive `docker-compose down && up` cycles
- Storage lives in a Docker volume attached to the backend container, persisting across Desktop restarts
- Extension UI reinitializes on every tab switch — all state must come from backend, never React state alone
- Extension framework only supports Linux containers, even on Windows/Mac

## Project Conventions

- Commit messages: conventional commits (`feat:`, `fix:`, `docs:`, `chore:`)
- Co-authored commits with Claude: `Co-Authored-By: Claude <noreply@anthropic.com>`
- Issues track all work; PRs reference issue numbers
- PowerShell is the primary scripting shell on Windows

## Build Commands

- `make test` — Run Go tests
- `make lint` — Run golangci-lint
- `make build-backend` — Build backend binary (local arch)
- `make dev-backend` — Run backend in dev mode (TCP :3001)
- `make cross-check` — Verify cross-compilation for linux/amd64 and linux/arm64
- `make build-extension` — Build Docker extension image
- `make install-extension` — Install extension into Docker Desktop

## CI and Linting

- **CI workflow**: `.github/workflows/lint.yml` — Markdown lint, Dockerfile lint, Go test, Go lint
- **Go test**: `go test -race ./...` on Linux CI runner
- **Go lint**: `golangci-lint-action@v7` with `v2.10.1`
- **Markdown lint**: `markdownlint-cli2` on all `**/*.md` files
- **Dockerfile lint**: `hadolint` on `Dockerfile`
- **hadolint config**: `.hadolint.yaml` ignores DL3048 (vendor labels) and DL3045 (COPY without WORKDIR) — standard Docker extension patterns
- **markdownlint config**: `.markdownlint.json` at project root (overrides DevSpace parent)
- **VERSION file**: Contains semver string (`0.1.0`), used by Makefile and metadata

## Phased Roadmap

Development follows a 5-phase plan with research, implementation, and gate issues at each phase. See GitHub issues for the full backlog. Phases: Backend (1) -> Frontend (2) -> Lifecycle (3) -> Enhancements (4) -> Ship (5).

## File Layout

```text
RunNotes/
├── cmd/backend/         - Go entry point (main.go)
├── internal/
│   ├── database/        - SQLite open, PRAGMAs, Goose migrations
│   ├── handler/         - HTTP handlers for note CRUD
│   ├── models/          - Note, CreateNoteRequest, UpdateNoteRequest
│   └── store/           - NoteStore data access layer
├── .github/workflows/   - CI workflow (lint, test, build)
├── docs/                - Architecture, data model, feasibility
├── ui/                  - React frontend (Phase 2)
├── go.mod / go.sum      - Go module definition
├── docker-compose.yaml  - VM service with volume mount
├── metadata.json        - Extension metadata (composefile pattern)
├── Dockerfile           - Multi-stage: Go builder + Alpine final
├── Makefile             - Build, test, lint, dev targets
└── VERSION              - Semver version string (0.1.0)
```
