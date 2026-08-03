# User Guide (template)

> This is the *user guide* template shipped under `docs/templates/general/`.
> The canonical, up-to-date guide for the `git-community-standards` tool lives
> at [USER-GUIDE.md](USER-GUIDE.md) in the repository root.

This document explains how to **use** a project that adopts these community
standards. It is intended to be copied into the *consuming* repository (for
example, as `USER-GUIDE.md` at its root) so contributors know how to install,
build, test, and ship changes.

---

## Who is this for?

This guide is for **contributors and maintainers** of this project. If you just
want to use the project as a library or binary, see the top-level `README.md`.

---

## Prerequisites

| Tool | Minimum version |
| --- | --- |
| Go | 1.23+ |
| GNU Make | any recent (optional) |
| Git | 2.x |

---

## Building

### Local build

```bash
make build
# binary: ./bin/<binary-name>
```

### Build from source (without Make)

```bash
go build -o <binary-name> .
```

On Windows the binary is `<binary-name>.exe`.

### Release binaries

```bash
make release
```

This produces multi-architecture binaries in `./dist/`.

---

## Testing

```bash
make test
# or directly:
go test ./...
```

Tests run on linux, windows, and macOS in CI (see `.github/workflows/test.yml`).

---

## Running

```bash
./bin/<binary-name> list        # show available templates
./bin/<binary-name> --version   # print the version
```

> Replace `<binary-name>` with the actual project binary name.

---

## Development workflow

1. Fork the repository.
2. Clone your fork locally.
3. Create a feature branch: `git switch -c feature/<short-description>`.
4. Make changes, then run `make test` and `make build` locally.
5. Commit using [Conventional Commits](https://www.conventionalcommits.com/):
   `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`.
6. Open a pull request describing what changed and why.

---

## Committing

Follow these rules:

- Use lowercase type prefixes: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`.
- Keep the subject line under 72 characters.
- Wrap the body at 72 characters.
- Reference issues: `Closes #123`, `Fixes #456`.

---

## Releasing

- Releases are automated via `.github/workflows/release.yml`.
- Tag a commit with a version (for example `v1.2.0`) and push the tag:
  ```bash
  git tag v1.2.0
  git push origin v1.2.0
  ```
- Release notes are generated from commits since the last release.

---

## Style and conventions

- Go code should pass `go vet ./...` and `gofmt` checks.
- Run `go test ./...` before pushing.
- Document new commands and flags in the relevant help/usage output.

---

## Getting help

- Open an issue for bugs or feature requests (use the provided templates).
- Start a discussion for ideas, questions, or feedback.
