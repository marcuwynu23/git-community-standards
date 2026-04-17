# Contributing

Thanks for your interest in improving `git-community-standards`.

We welcome bug fixes, features, documentation improvements, and tooling updates.

---

## Getting Started

1. Fork and clone the repository.
2. Create a branch from `main`.
3. Make focused changes with clear intent.
4. Run local checks before opening a PR.

---

## Development Workflow

1. Implement your change in a focused branch.
2. Keep the CLI behavior consistent across `git community-standards ...` commands.
3. Update docs (`README.md`, `RELEASE-NOTES.md`) when behavior changes.
4. Add/update tests when testable behavior changes.
5. Run local checks:

```bash
make test
make build
```

---

## Guidelines

- Keep changes small and focused.
- Prefer standard library solutions unless a dependency is clearly justified.
- Preserve backwards compatibility for existing commands when possible.
- Keep user-facing output clear and actionable.
- Ensure CI workflows stay green (`test.yml`, `release.yml`).

---

## Pull Request Process

1. Rebase/update your branch with `main`.
2. Ensure checks pass locally (`make test`, `make build`).
3. Open a PR with:
   - what changed
   - why it changed
   - how you validated it
4. If CLI output changed, include example command output in the PR description.

---

## Commit Messages

Use clear, concise, intent-first commit messages.

Examples:

```
fix race condition in worker queue
add support for custom config path
update documentation for setup
```

---

## Reporting Issues

When opening an issue, please include:

- Description of the problem
- Steps to reproduce
- Expected vs actual behavior
- OS, shell, and CLI version (`git community-standards --version`)

---

## Suggestions & Feature Requests

- Explain the problem you want to solve.
- Describe your proposed command/API behavior.
- Mention alternatives considered.

---

## Code of Conduct

Please follow `CODE_OF_CONDUCT.md` in all interactions.

---

## Notes

- Maintainers may request changes before merge.
- Not all contributions are guaranteed to be accepted, but all will be reviewed.

---

Thanks again for contributing.
