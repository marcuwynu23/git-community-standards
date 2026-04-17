# Release Notes

## v1.0.0 - 2026-04-17

### Added
- Initial release of `git-community-standards`.
- CLI commands:
  - `git community-standards list`
  - `git community-standards apply`
  - `git community-standards apply <category>`
  - `git community-standards apply override`
  - `git community-standards apply <category> override`
  - `git community-standards --version`
  - `git community-standards -v`
- Category-based file application:
  - `root`
  - `github`
  - `issue-templates`
  - `pr-template`
- Multi-platform release workflow for linux, windows, and darwin across amd64 and arm64.
- Project tooling and docs:
  - `Makefile` with `test`, `build`, and `release`
  - CI workflows (`test.yml`, `release.yml`)
  - repository community templates and docs

### Changed
- Usage/help output standardized for `git community-standards` style.

### Fixed
- Default apply behavior now skips existing files unless override mode is requested.
- Optional upstream files handled safely with fallback paths and non-fatal skip behavior.

### Security
- No security changes in this release.
