# RunNotes -- Copilot Instructions

Docker Desktop extension for managing notes and documentation for Docker containers and images.

## Tech Stack

- **Frontend**: React 18, TypeScript, MUI v5, Vite
- **Extension SDK**: @docker/extension-api-client, @docker/docker-mui-theme
- **Backend**: Go (Docker Desktop extension VM service)
- **Container**: Docker (multi-arch linux/amd64, linux/arm64)

## Project Structure

```text
RunNotes/
├── cmd/              - Go backend entry point (extension VM service)
├── internal/         - Go backend packages
├── ui/               - React frontend source (components, hooks, API)
├── docs/             - Documentation and screenshots
├── scripts/          - Build and utility scripts
├── Dockerfile        - Multi-stage build (Go backend + React frontend)
├── metadata.json     - Docker Desktop extension metadata
├── .github/          - CI workflows and Copilot config
└── CLAUDE.md         - Claude Code instructions
```

## Code Style

- Conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
- Co-author tag: `Co-Authored-By: GitHub Copilot <noreply@github.com>`
- Functional React components with typed props interfaces
- MUI v5 patterns (`InputProps`, not `slotProps.input` for TextField adornments)
- Vitest + Testing Library for frontend tests

## Coding Guidelines

- Fix errors immediately -- never classify them as pre-existing
- Build, test, and lint must pass before any commit
- Never skip hooks (`--no-verify`) or force-push main
- Use `unknown` with type guards instead of `any` in TypeScript
- Remove unused code completely; no backwards-compatibility hacks

## Available Resources

```bash
cd ui && npm run build    # Build frontend
cd ui && npm test         # Run frontend tests
cd ui && npm run lint     # Lint frontend
cd ui && npx tsc --noEmit # TypeScript type check
docker build -t runnotes .  # Build Docker image
docker extension validate runnotes  # Validate extension metadata
```

## Do NOT

- Use `any` in TypeScript or suppress TypeScript errors with `as unknown`
- Commit generated files without regenerating them first
- Add dependencies without updating the lock file (`npm install`)
- Store secrets, tokens, or credentials in code or config files
- Mark work as complete when known errors remain
- Use MUI v6 patterns (this project uses MUI v5 via @docker/docker-mui-theme)
