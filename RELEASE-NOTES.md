# Release Notes

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [Unreleased]

### Added
- Go CLI tool for applying repository community standards from a template source.
- `list` command to show available categories.
- `apply` command with category support.
- `override` mode to replace existing files when explicitly requested.
- Version flags: `--version` and `-v`.
- Git-style command support: `git community-standards ...`.
- Makefile targets: `test`, `build`, `release`.
- GitHub Actions workflows:
  - `test.yml` (test + cross-platform build checks)
  - `release.yml` (tag-based multi-OS/multi-arch release assets)

### Changed
- Usage output standardized around `git community-standards` command form.
- Documentation tailored to this project and current CLI behavior.
- Release workflow now publishes archived artifacts (`.zip`) for all OS/architecture targets.
- GitHub Release body now uses the contents of `RELEASE-NOTES.md`.

### Deprecated
- None.

### Removed
- None.

### Fixed
- Safer default apply behavior by skipping existing files unless `override` is used.
- Improved resilience for optional upstream files using fallback source paths.

### Security
- None.

---

## [1.0.0] - 2026-04-17

### Added
- Initial release of the project
- Core features implemented:
  - Community standards file application CLI
  - Category-based apply support
  - Optional override behavior
  - Multi-platform release pipeline

### Changed
- N/A

### Fixed
- N/A

### Security
- N/A

---

## Release Guidelines

### Versioning
This project follows **Semantic Versioning (SemVer)**:
- **MAJOR**: incompatible API changes
- **MINOR**: backwards-compatible features
- **PATCH**: backwards-compatible bug fixes

---

## Notes

- Include links to issues or PRs when possible:
  - Example: Improve apply override handling ([#42](https://github.com/marcuwynu23/git-community-standards/pull/42))
- Highlight breaking changes clearly under a "Breaking Changes" section if needed.
- Keep entries concise and user-focused.

---

## Contributors

Thanks to everyone who contributed to this release:

- @marcuwynu23