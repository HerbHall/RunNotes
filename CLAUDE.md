# RunNotes

Docker Desktop extension for attaching notes and annotations to containers.

## Tech Stack

- **Frontend**: React + Material UI (Docker Desktop extension standard)
- **Backend**: Go service running in Desktop VM (pure-Go SQLite, single binary)
- **Storage**: SQLite on Docker volume (see [docs/DATA_MODEL.md](docs/DATA_MODEL.md))
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

## CI and Linting

- **Markdown lint**: `.github/workflows/lint.yml` runs `markdownlint-cli2` on all `**/*.md` files
- **Dockerfile lint**: Same workflow runs `hadolint` on `Dockerfile`
- **hadolint config**: `.hadolint.yaml` ignores DL3048 (vendor labels) and DL3045 (COPY without WORKDIR) — standard Docker extension patterns
- **markdownlint config**: `.markdownlint.json` at project root (overrides DevSpace parent)
- **VERSION file**: Contains semver string (`0.1.0`), used by Makefile and metadata

## Phased Roadmap

Development follows a 5-phase plan with research, implementation, and gate issues at each phase. See GitHub issues for the full backlog. Phases: Backend (1) -> Frontend (2) -> Lifecycle (3) -> Enhancements (4) -> Ship (5).

## File Layout

```text
RunNotes/
├── .github/workflows/   - CI workflows (lint.yml)
├── .hadolint.yaml       - Hadolint config for Dockerfile linting
├── .markdownlint.json   - Markdownlint config
├── docs/                - Architecture docs, data model, feasibility research
├── ui/                  - React frontend (scaffold when starting Phase 2)
├── backend/             - Go backend service (scaffold when starting Phase 1)
├── metadata.json        - Extension metadata
├── Dockerfile           - Extension image
├── Makefile             - Build targets
├── VERSION              - Semver version string
├── CLAUDE.md            - Project context for Claude Code
├── HANDOFF.md           - Claude Code handoff notes
├── CONTRIBUTING.md      - Contribution guidelines
├── CHANGELOG.md         - Keep-a-Changelog format
└── README.md            - Project overview
```
