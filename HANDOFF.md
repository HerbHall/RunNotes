# RunNotes — Claude Code Handoff

## What This Is

RunNotes is a Docker Desktop extension that lets users attach notes/annotations to containers. It fills a real gap — Docker Desktop shows container name, image, status, ports, and logs but has no way to capture *why* a container exists or what you were doing with it.

**Owner**: Herb Hall (github.com/HerbHall)
**License**: MIT
**Status**: Pre-development — scaffold exists, no code yet

---

## What Already Exists (D:\devspace\RunNotes)

The project scaffold is mostly in place:

- ✅ `CLAUDE.md` — Project context
- ✅ `README.md` — Full project overview
- ✅ `CONTRIBUTING.md` — Contribution guidelines
- ✅ `CHANGELOG.md` — Keep-a-Changelog format
- ✅ `LICENSE` — MIT
- ✅ `VERSION` — 0.1.0
- ✅ `.gitignore` — Comprehensive
- ✅ `.editorconfig` — Workspace standard
- ✅ `metadata.json` — Docker extension metadata (basic)
- ✅ `Dockerfile` — Extension image stub (labels set, stages TODO)
- ✅ `Makefile` — Build/install/push targets
- ✅ `docker.svg` — Placeholder icon
- ✅ `docs/ARCHITECTURE.md` — Full architecture design
- ✅ `docs/FEASIBILITY.md` — Feasibility assessment
- ✅ `create-issues.bat` — 10 GitHub issues ready to create
- ✅ `.git` initialized
- ✅ `.github/` directory exists (empty)

---

## What Still Needs To Be Done

### 1. GitHub Repository (FIRST PRIORITY)

Create the remote repo and push:

```powershell
# In CMD (not PowerShell — avoid bracket escaping issues with gh)
cd /d D:\devspace\RunNotes
cmd /c "gh repo create HerbHall/RunNotes --public --source=. --remote=origin --description "Docker Desktop extension for attaching notes to containers""
git add -A
git commit -m "chore: initial project scaffold"
git push -u origin main
```

Herb has `gh` CLI installed and authenticated as HerbHall. Use `cmd /c` wrapper in PowerShell to avoid bracket escaping issues.

### 2. Create GitHub Issues

Run the existing batch file:

```powershell
cmd /c "D:\devspace\RunNotes\create-issues.bat"
```

This creates 10 issues covering the full backlog:

1. Data model design (mvp)
2. React UI — container list + note editor (mvp)
3. Backend service — storage API (mvp)
4. Container identity/lifecycle management (mvp)
5. Orphaned notes detection (enhancement)
6. Markdown support (enhancement)
7. Search across notes (enhancement)
8. Export/import JSON (enhancement)
9. Docker Hub publishing (chore)
10. Icon and branding (docs)

**NOTE**: The batch file uses labels `feat,mvp` etc. These labels need to exist on the repo first. Create them:

```powershell
cmd /c "gh label create mvp --color 0E8A16 --description "Minimum viable product" --repo HerbHall/RunNotes"
cmd /c "gh label create feat --color 1D76DB --description "New feature" --repo HerbHall/RunNotes"
cmd /c "gh label create enhancement --color A2EEEF --description "Enhancement" --repo HerbHall/RunNotes"
cmd /c "gh label create chore --color FBCA04 --description "Maintenance task" --repo HerbHall/RunNotes"
```

### 3. Missing Files (Match DevKit Patterns)

These exist in DevKit but not RunNotes:

- `.markdownlint.json` — Copy from `D:\devspace\.markdownlint.json`
- `.github/workflows/lint.yml` — Basic CI (markdownlint, dockerfile lint)
- `METHODOLOGY.md` — Optional, only if project methodology needs documenting

### 4. Source Directories (Not Yet Created)

The actual code directories don't exist yet:

```text
ui/           — React frontend (create when starting issue #2)
backend/      — Go or Node backend service (create when starting issue #3)
```

Don't create these empty — they should be scaffolded when development begins.

### 5. Metadata.json Enhancements

Current `metadata.json` is minimal. When development starts, add:

- `com.docker.desktop.extension.api.version` (already in Dockerfile labels, should also be validated)
- Consider adding a `compose.yaml` for the VM service if it needs volume mounts

### 6. VS Code Workspace File

Create `D:\devspace\runnotes.code-workspace` to match the DevKit pattern:

```json
{
  "folders": [
    { "path": "RunNotes" }
  ],
  "settings": {}
}
```

---

## Key Architecture Decisions

These are settled from the research phase — don't revisit:

- **Notes keyed by container NAME** (not ID) as primary key, so notes survive `docker-compose down && up`
- **Container ID as secondary key** for disambiguation
- **Storage in Docker volume** attached to backend container (survives DD restarts)
- **SQLite preferred** over JSON for storage (supports querying, single-file)
- **Socket communication** between frontend and backend (Extensions SDK standard, no port conflicts)
- **React + Material UI** frontend to match Docker Desktop look and feel
- **Extension framework only supports Linux containers** even on Windows — this is fine, backend runs in DD VM
- **Multi-arch required**: linux/amd64 + linux/arm64

## Name Research

"RunNotes" was chosen after research confirmed:

- No Docker CLI command or concept named "runnotes"
- No existing Docker extension using this name
- The name "Manifest" was considered first but rejected — `docker manifest` is a core Docker CLI command AND there's already a Docker Desktop extension called "Manifest" (memory metrics tool by oslabs-beta)
- "RunNotes" connects to `docker run` + notes — immediately clear to Docker users

## Docker Desktop Extension SDK Key Facts

- Extensions are distributed as Docker images via Docker Hub
- Architecture: React frontend tab + optional backend VM service + optional host binaries
- Frontend can invoke `docker` CLI commands via `ddClient.docker.cli`
- Frontend-to-backend communication via socket/named pipe (not HTTP)
- UI reinitializes EVERY time user navigates to the extension tab — state must come from backend
- `docker extension init` scaffolds a working starter project
- Image: `herbhall/runnotes` on Docker Hub

## Herb's Preferences (Important)

- **Green/earthy colors** for branding — dislikes blue
- **PowerShell is primary shell** on Windows
- **Use `cmd /c` wrapper** for `gh` commands in PowerShell (bracket escaping)
- **Conventional commits**: `feat:`, `fix:`, `docs:`, `chore:`
- **Co-authored commits**: `Co-Authored-By: Claude <noreply@anthropic.com>`
- **Executes steps immediately as he reads them** — put prerequisites BEFORE action steps
- **Visual/point-and-click preference** over CLI memorization
- **Self-hosted solutions** preferred over cloud-dependent alternatives

---

## Suggested Session Plan

1. Create GitHub repo + push scaffold
2. Create labels + run create-issues.bat
3. Copy `.markdownlint.json` from devspace root
4. Create `.github/workflows/lint.yml` (basic CI)
5. Create `runnotes.code-workspace`
6. Commit and push all additions
7. Begin MVP development (issues #1-4) when ready
