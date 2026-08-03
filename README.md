<div align="center">

# git-community-standards

CLI tool to apply standard community files to a repository from a shared template source, with templates for GitHub, GitLab, and Bitbucket.

[![Test](https://img.shields.io/github/actions/workflow/status/marcuwynu23/git-community-standards/test.yml?branch=main&label=test)](https://github.com/marcuwynu23/git-community-standards/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/actions/workflow/status/marcuwynu23/git-community-standards/release.yml?label=release)](https://github.com/marcuwynu23/git-community-standards/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/v/release/marcuwynu23/git-community-standards)](https://github.com/marcuwynu23/git-community-standards/releases)
[![License](https://img.shields.io/github/license/marcuwynu23/git-community-standards)](./LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.23%2B-00ADD8)](https://go.dev/)

</div>

## Overview

`git-community-standards` helps keep repository governance docs consistent by downloading and writing common files:

**General docs (always applied):**
- `README.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `RELEASE-NOTES.md`
- `LICENSE` (when available in source)
- `SECURITY.md` (when available in source)

**Platform templates (choose `github`, `gitlab`, or `bitbucket`):**
- GitHub: `.github/FUNDING.yml`, `.github/ISSUE_TEMPLATE/*`, `.github/PULL_REQUEST_TEMPLATE.md`
- GitLab: `.gitlab/issue_templates/*`, `.gitlab/merge_request_templates/default.md`
- Bitbucket: `.bitbucket/ISSUE_TEMPLATE/*`, `.bitbucket/PULL_REQUEST_TEMPLATE.md`

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
git community-standards apply <github|gitlab|bitbucket|none>
git community-standards apply override
git community-standards apply <github|gitlab|bitbucket> override
git community-standards --version
git community-standards -v
```

If you run the binary directly, the same command arguments work without the leading `git`.

Output is colorized when written to an interactive terminal. Piping, redirecting, setting `NO_COLOR`, or the `dumb` `TERM` value disables color automatically; errors always go to stderr.

### Commands

- `list`: show available platforms and general docs.
- `apply`: apply the general community docs only.
- `apply <platform>`: apply general docs plus the given platform templates (`github`, `gitlab`, or `bitbucket`).
- `apply none` / `apply general`: apply general community docs only.
- `override`: replace existing files instead of skipping them.

### Platforms

- `none` (default): general community docs only.
- `github`: GitHub templates under `.github` (FUNDING, issue templates, pull request template).
- `gitlab`: GitLab templates under `.gitlab` (issue templates, merge request template).
- `bitbucket`: Bitbucket templates (issue templates and pull request template).

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
