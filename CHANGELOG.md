# Changelog

All notable changes to RunNotes will be documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.1.0] - 2026-03-02

### Added

- Container notes with tag-based organization
- Markdown support with live preview toggle
- Name-based persistence (survives docker-compose down/up cycles)
- Container ID auto-reconciliation when containers are recreated
- Search by container name, note content, or tags
- Pin/unpin notes for quick access
- Orphaned note detection and cleanup dialog
- Export/import notes as JSON backup
- SQLite storage on Docker volume for persistence across restarts
