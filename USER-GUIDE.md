# User Guide

`git-community-standards` is a command-line tool that bootstraps or refreshes
the standard community files (README, CONTRIBUTING, CODE_OF_CONDUCT, pull
request / issue templates, security policy, and more) for a repository from a
shared template source hosted on GitHub.

---

## Installation

### Build from source

```bash
git clone https://github.com/marcuwynu23/git-community-standards.git
cd git-community-standards
go build -o git-community-standards .
```

On Windows, build `git-community-standards.exe` or run `make build`.

### Install via the Makefile

```bash
make build     # builds bin/git-community-standards(.exe)
```

### Pre-built releases

Run the [release workflow](https://github.com/marcuwynu23/git-community-standards/actions)
to download binaries for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64,
windows/amd64, and windows/arm64.

---

## Usage

The tool is designed to be invoked as a `git` subcommand alias, but it also
works when run as a standalone binary.

```bash
git community-standards list
git community-standards apply
git community-standards apply <github|gitlab|bitbucket|none>
git community-standards apply override
git community-standards apply <github|gitlab|bitbucket> override
git community-standards --version
git community-standards -v
```

> If `git community-standards` is not configured as a git subcommand, invoke
> the binary directly (`git-community-standards`) with the same arguments.

---

## Commands

| Command | Description |
| --- | --- |
| `list` | Show the general docs and the available platform template sets. |
| `apply` | Apply the general community docs only. |
| `apply <platform>` | Apply general docs plus the given platform templates. |
| `apply none` / `apply general` | Apply general community docs only (explicit). |
| `apply override` | Replace existing files instead of skipping them. |
| `--version` / `-v` | Print the tool version. |

### Platforms

| Platform | Target directory | Contents |
| --- | --- | --- |
| `none` (default) | repository root | README, LICENSE, CONTRIBUTING, CODE_OF_CONDUCT, RELEASE-NOTES, SECURITY |
| `github` | `.github/` | `FUNDING.yml`, `ISSUE_TEMPLATE/bug_report.md`, `ISSUE_TEMPLATE/feature_request.md`, `PULL_REQUEST_TEMPLATE.md` |
| `gitlab` | `.gitlab/` | `issue_templates/bug_report.md`, `issue_templates/feature_request.md`, `merge_request_templates/default.md` |
| `bitbucket` | `.bitbucket/` | `ISSUE_TEMPLATE/bug_report.md`, `ISSUE_TEMPLATE/feature_request.md`, `PULL_REQUEST_TEMPLATE.md` |

#### GitHub

```bash
git community-standards apply github
```

Adds GitHub-specific templates under `.github`.

#### GitLab

```bash
git community-standards apply gitlab
```

Adds GitLab-specific templates under `.gitlab` (issue templates and a merge
request template). GitLab uses *Merge Request* terminology, so the pull-request
template file is named `MR_TEMPLATE.md` and is written to
`.gitlab/merge_request_templates/default.md`.

#### Bitbucket

```bash
git community-standards apply bitbucket
```

Adds Bitbucket-specific templates under `.bitbucket`.

---

## Template priority and fallbacks

Each file spec lists a primary `RemotePath` plus one or more `Fallbacks`. The
tool downloads the first source that returns successfully. This keeps the tool
forward-compatible with older template layouts.

Examples:
- `LICENSE` is fetched from `general/LICENSE.md`, then `general/LICENSE`,
  `general/LICENCE`, and `general/LICENCE.md`.
- `SECURITY.md` is fetched from `general/SECURITY.md`, falling back to
  `general/SECURITY_POLICY.md`.

---

## Override behavior

By default the tool **skips** any file that already exists in your repository
so it never clobbers existing work:

```
Skipped README.md (already exists). Use `apply override` to replace it.
```

To force-replace existing files, append `override`:

```bash
git community-standards apply github override
```

---

## Output and colors

Output is colorized when written to an interactive terminal. Color is
automatically disabled when:
- stdout is piped or redirected,
- the `NO_COLOR` environment variable is set, or
- the `TERM` variable is empty or `dumb`.

All error messages are written to **stderr**, so they stay visible when stdout
is piped.

---

## Project layout

```
.
├── main.go                 # thin entry point (`package main`)
├── internal/app/           # core logic (testable library code)
│   ├── app.go
│   ├── color.go
│   ├── filelists.go
│   └── app_test.go
└── docs/templates/
    ├── general/            # platform-independent community docs
    ├── github/             # GitHub-specific templates
    ├── gitlab/             # GitLab-specific templates
    └── bitbucket/          # Bitbucket-specific templates
```

---

## Development

### Requirements

- Go 1.23+
- GNU Make (optional, for convenience targets)

### Local checks

```bash
make test    # run unit tests
make build   # build local binary into ./bin
```

### Cross-platform build check

```bash
make test           # unit tests (linux, windows, macos)
make release         # produces multi-arch binaries in ./dist
```

See the [Makefile](Makefile) for available targets.

---

## Troubleshooting

**Files are skipped unexpectedly.** Ensure you are running from the repository
root and re-run with `override` to replace existing files.

**Downloads fail.** The tool fetches templates from
`https://raw.githubusercontent.com/marcuwynu23/git-community-standards/refs/heads/main/docs/templates`.
Check network access and that you are not blocked by a firewall or proxy.

**No color in output.** Color only enables for interactive terminals. Confirm
stdout is a TTY, or unset `NO_COLOR` and set a real `TERM`.
