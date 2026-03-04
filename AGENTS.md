<!--
  Scope: AGENTS.md guides the Copilot coding agent and Copilot Chat.
  For code completion and code review patterns, see .github/copilot-instructions.md
  and .github/instructions/*.instructions.md
  For Claude Code, see CLAUDE.md
-->

# RunNotes

Docker Desktop extension for managing notes and documentation for Docker containers and images.

## Tech Stack

- **Frontend**: React 18, TypeScript, MUI v5, Vite
- **Extension SDK**: @docker/extension-api-client, @docker/docker-mui-theme
- **Backend**: Go (Docker Desktop extension VM service)
- **Container**: Docker (multi-arch linux/amd64, linux/arm64)

## Build and Test Commands

```bash
# Frontend build
cd ui && npm run build

# Frontend tests
cd ui && npm test

# Frontend lint
cd ui && npm run lint

# TypeScript type check
cd ui && npx tsc --noEmit

# Docker build
docker build -t runnotes .

# Extension validation
docker extension validate runnotes

# Full verification (run before any PR)
cd ui && npx tsc --noEmit && npm run lint && npm run build
```

## Project Structure

```text
RunNotes/
├── cmd/              - Go backend entry point (extension VM service)
├── internal/         - Go backend packages
├── ui/               - React frontend source
│   ├── src/          - Components, hooks, API layer
│   └── package.json  - npm dependencies
├── docs/             - Documentation and screenshots
├── scripts/          - Build and utility scripts
├── Dockerfile        - Multi-stage build (Go backend + React frontend)
├── docker-compose.yaml
├── metadata.json     - Docker Desktop extension metadata
├── docker.svg        - Extension icon
├── Makefile          - Build automation
├── CLAUDE.md         - Claude Code instructions
└── .github/          - CI workflows and Copilot config
```

## Workflow Rules

### Always Do

- Create a feature branch for every change (`feature/issue-NNN-description`)
- Use conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
- Run build, test, and lint before opening a PR
- Define props interfaces above components; use functional components
- Use `InputProps` (not `slotProps.input`) for MUI v5 TextField adornments
- Fix every error you find, regardless of who introduced it

### Ask First

- Adding new dependencies (check if existing packages cover the need)
- Architectural changes (new packages, major interface changes)
- Changes to the Docker Desktop extension metadata or Dockerfile labels
- Changes to CI/CD workflows
- Removing or renaming exported components or API functions

### Never Do

- Commit directly to `main` -- always use feature branches
- Skip tests or lint checks -- even for "small changes"
- Use `--no-verify` or `--force` flags
- Commit secrets, credentials, or API keys
- Add TODO comments without a linked issue number
- Mark work as complete when build, test, or lint failures remain
- Use `any` in TypeScript -- use `unknown` with type guards instead

## Core Principles

These are unconditional -- no optimization or time pressure overrides them:

1. **Quality**: Once found, always fix, never leave. There is no "pre-existing" error.
2. **Verification**: Build, test, and lint must pass before any commit.
3. **Safety**: Never force-push `main`. Never skip hooks. Never commit secrets.
4. **Honesty**: Never mark work as complete when it is not.

## Testing Conventions

```tsx
// Vitest + Testing Library
// Use descriptive test names that explain the expected behavior
describe('NoteEditor', () => {
    it('renders the note title from props', () => {
        render(<NoteEditor title="My Note" onChange={vi.fn()} />)
        expect(screen.getByDisplayValue('My Note')).toBeInTheDocument()
    })

    it('calls onChange when the title is edited', async () => {
        const onChange = vi.fn()
        render(<NoteEditor title="" onChange={onChange} />)
        await userEvent.type(screen.getByRole('textbox'), 'New Title')
        expect(onChange).toHaveBeenCalled()
    })
})
```

## Commit Format

```text
feat: add container note search

Implements full-text search across container notes with debounced input.

Closes #42
Co-Authored-By: GitHub Copilot <copilot@github.com>
```

Types: `feat` (new feature), `fix` (bug fix), `refactor` (no behavior change),
`docs` (documentation only), `test` (tests only), `chore` (build/tooling).
