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

## File Layout

```text
RunNotes/
├── docs/                - Architecture docs, feasibility research
├── ui/                  - React frontend
├── backend/             - Backend service
├── metadata.json        - Extension metadata
├── Dockerfile           - Extension image
└── Makefile             - Build targets
```
