# git-community-standards

CLI tool to apply standard community files to a GitHub repository from a shared template source.

## Overview

`git-community-standards` helps keep repository governance docs consistent by downloading and writing common files such as:

- `README.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `RELEASE-NOTES.md`
- `LICENSE` (when available in source)
- `SECURITY.md` (when available in source)
- `.github/FUNDING.yml`
- `.github/ISSUE_TEMPLATE/*`
- `.github/PULL_REQUEST_TEMPLATE.md`

The tool supports safe default behavior (skip existing files) and explicit override mode.

## Installation

### Build from source

```bash
git clone https://github.com/marcuwynu23/git-community-standards.git
cd git-community-standards
go build -o git-community-standards .
```

On Windows, build `git-community-standards.exe`.

### Development build with Makefile

```bash
make build
```

## Usage

This project is designed to be used in git-style form:

```bash
git community-standards list
git community-standards apply
git community-standards apply <category>
git community-standards apply override
git community-standards apply <category> override
git community-standards --version
git community-standards -v
```

If you run the binary directly, the same command arguments work without the leading `git`.

### Commands

- `list`: show available categories.
- `apply`: apply all categories.
- `apply <category>`: apply one category.
- `override`: replace existing files instead of skipping them.

### Categories

- `root`
- `github`
- `issue-templates`
- `pr-template`

## Development

### Requirements

- Go 1.23+
- GNU Make (optional, for convenience commands)

### Local checks

```bash
make test
make build
```

### Release builds

```bash
make release
```

This produces binaries in `dist/` for:

- linux/amd64, linux/arm64
- darwin/amd64, darwin/arm64
- windows/amd64, windows/arm64

## CI and Releases

- `test.yml` validates tests and cross-platform build checks.
- `release.yml` runs on tags (`v*`), depends on test workflow, and publishes release assets for all OS/architecture targets.

## Contributing

Contributions are welcome. See `CONTRIBUTING.md` for contribution workflow and standards.

## License

See `LICENSE`.
