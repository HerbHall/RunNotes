# Screenshot Capture Guide

Screenshots for the Docker Desktop Extensions Marketplace listing.

## Requirements

- **Minimum**: 3 screenshots
- **Recommended dimensions**: 2400 x 1600 pixels
- **Format**: PNG
- **Save to**: `docs/screenshots/`

## Setup

1. Build and install the extension locally:

   ```bash
   make build-extension
   echo "y" | docker extension install herbhall/runnotes:0.1.0
   ```

2. Create 3-4 sample containers with different states:

   ```bash
   docker run -d --name web-server nginx
   docker run -d --name redis-cache redis
   docker run -d --name postgres-db postgres -e POSTGRES_PASSWORD=test
   docker create --name api-staging node:22-alpine
   ```

3. Open Docker Desktop and navigate to the RunNotes extension tab.

4. Add sample notes to containers:
   - **web-server**: A markdown note with headings, bullet points, and code blocks
   - **redis-cache**: A short note with tags like `cache`, `session`
   - **postgres-db**: A pinned note with connection details

## Screenshots to Capture

### 1. Main View (`main.png`)

Show the full extension UI with:

- Container list on the left (3-4 containers visible)
- One container selected
- Note content visible in the editor on the right
- Tags visible below the editor

### 2. Markdown Preview (`markdown-preview.png`)

Same container selected but with the Preview toggle active:

- Toggle button group showing "Preview" selected
- Rendered markdown visible (headings, lists, code blocks)
- Demonstrates the markdown formatting capability

### 3. Search and Filter (`search-filter.png`)

Show the search/filter functionality:

- Type a search term in the search bar (e.g., container name or tag)
- Filtered container list showing matching results
- Or show the pin filter active with pinned containers highlighted

### 4. Export/Import (Optional) (`export-import.png`)

Show the export/import buttons in the header bar.

## Capture Tips

- Use **light mode** for primary marketplace screenshots (Docker Desktop > Settings > Appearance)
- Use Windows Snipping Tool or Snip & Sketch (`Win+Shift+S`) for precise captures
- Resize the Docker Desktop window to approximately 2400x1600 before capturing, or capture and scale
- Ensure no sensitive data is visible in container names or note content

## After Capture

1. Save screenshots as PNG files in `docs/screenshots/`
2. The Dockerfile label `com.docker.extension.screenshots` will be updated with raw GitHub URLs pointing to these files
