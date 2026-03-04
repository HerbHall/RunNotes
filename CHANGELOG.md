# Changelog

All notable changes to RunNotes will be documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.1.1](https://github.com/HerbHall/RunNotes/compare/v0.1.0...v0.1.1) (2026-03-04)


### Features

* add export/import JSON for notes backup ([#8](https://github.com/HerbHall/RunNotes/issues/8)) ([#33](https://github.com/HerbHall/RunNotes/issues/33)) ([1c8a52a](https://github.com/HerbHall/RunNotes/commit/1c8a52aa6af355dab737c2fce829d3afbd7f2bf0))
* add extension icon, branding, and publishing prep ([#10](https://github.com/HerbHall/RunNotes/issues/10)) ([#37](https://github.com/HerbHall/RunNotes/issues/37)) ([a6a31d3](https://github.com/HerbHall/RunNotes/commit/a6a31d3705ccae388f379173a8b59aa6210836f5))
* add GitHub Copilot integration files ([#42](https://github.com/HerbHall/RunNotes/issues/42)) ([19660c3](https://github.com/HerbHall/RunNotes/commit/19660c3bca770c8819f20268a2973c1917aaa931))
* add markdown support with preview toggle ([#6](https://github.com/HerbHall/RunNotes/issues/6)) ([#34](https://github.com/HerbHall/RunNotes/issues/34)) ([b1c2d73](https://github.com/HerbHall/RunNotes/commit/b1c2d7359b353e871411d6cde71e980b891f410c))
* add orphaned notes detection and cleanup ([#5](https://github.com/HerbHall/RunNotes/issues/5)) ([#35](https://github.com/HerbHall/RunNotes/issues/35)) ([7b5fef2](https://github.com/HerbHall/RunNotes/commit/7b5fef25f2cf55156d8f8f998c3cda28bb1e4322))
* container lifecycle management with ID reconciliation ([#4](https://github.com/HerbHall/RunNotes/issues/4)) ([#31](https://github.com/HerbHall/RunNotes/issues/31)) ([40222d5](https://github.com/HerbHall/RunNotes/commit/40222d573085a1333411954ff1df55665ad5ddc8))
* extend search to include tags ([#7](https://github.com/HerbHall/RunNotes/issues/7)) ([#32](https://github.com/HerbHall/RunNotes/issues/32)) ([63102ee](https://github.com/HerbHall/RunNotes/commit/63102ee93138c08f17ae1c0a9959091fb9a08532))
* implement backend service with SQLite storage API ([#27](https://github.com/HerbHall/RunNotes/issues/27)) ([efecdf7](https://github.com/HerbHall/RunNotes/commit/efecdf7f1cc5fc17d5b6517c073f3c233fa42d34))
* implement container list and note editor UI ([#2](https://github.com/HerbHall/RunNotes/issues/2)) ([#30](https://github.com/HerbHall/RunNotes/issues/30)) ([c3d05e0](https://github.com/HerbHall/RunNotes/commit/c3d05e0229d571cd64fa56bae50b809d9b444895))
* initial project structure and documentation ([6460244](https://github.com/HerbHall/RunNotes/commit/646024498baf4db9129f8156725bf95f6ccfca6f))
* scaffold React frontend with Vite, MUI, and Docker theme ([#2](https://github.com/HerbHall/RunNotes/issues/2)) ([#29](https://github.com/HerbHall/RunNotes/issues/29)) ([38ae3fb](https://github.com/HerbHall/RunNotes/commit/38ae3fbb3ffaf0bcb0c031d69dce9c12b88d073e))
* support multiple notes per container ([#39](https://github.com/HerbHall/RunNotes/issues/39)) ([76bcda6](https://github.com/HerbHall/RunNotes/commit/76bcda6fadeb233a8cae502ce8732d812646c371))


### Bug Fixes

* align Go version to 1.25 and add stub UI for extension install ([#28](https://github.com/HerbHall/RunNotes/issues/28)) ([123dabf](https://github.com/HerbHall/RunNotes/commit/123dabf7da2f1f21545e28a9e4b4a30adc792ab8))
* use versioned tag in screenshot guide install command ([2b9bbd7](https://github.com/HerbHall/RunNotes/commit/2b9bbd70f6c9c913df21e16097d165ec481c7048))

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
