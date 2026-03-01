# RunNotes Feasibility Assessment

## Summary

**Verdict**: Highly feasible. The Docker Desktop Extensions SDK provides all necessary capabilities and there is no existing extension filling this niche.

## Technical Feasibility

The Docker Desktop Extensions SDK supports:

- **React frontend** with a dedicated tab in the sidebar
- **Backend services** running in the Desktop VM with Docker volume persistence
- **Host executables** for platform-specific operations
- **Docker CLI access** from the frontend via `ddClient.docker.cli`
- **Socket-based communication** between frontend and backend (no port conflicts)

The SDK is mature (GA since Docker Desktop 4.8.0) and well-documented with TypeScript support.

## Market Gap

Docker Desktop currently provides container name, image, status, ports, and logs — but no mechanism for user-authored annotations. No existing extension in the Docker Hub marketplace addresses this use case. The closest tools are monitoring/metrics extensions, which serve a different purpose.

## Effort Estimate

| Tier | Scope | Effort |
|------|-------|--------|
| MVP | Text notes per container, local volume storage, list view | 1-2 days |
| Solid | Markdown, search, name-based persistence, export | 3-5 days |
| Full | Tags, Compose-project notes, orphan management, shared notes | 1-2 weeks |
