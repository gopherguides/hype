# Repository Guidelines

## Project Structure & Module Organization

Core packages live at repo root (e.g., `element.go`, `parser.go`) and follow Go's standard package-per-file pattern, each paired with `*_test.go`. The CLI wrapper sits in `cmd/hype`, which builds the `hype` binary. Content templates and examples live under `docs/`, `mdx/`, and `slides/`, while reusable fixtures are stored in `testdata/`. Helper utilities and adapters reside in `internal/` plus feature-specific subfolders such as `atomx/` and `binding/`.

## Build, Test, and Development Commands

- `make test`: runs `go test -count 1 -race -vet=off -cover` across all packages except generated docs.
- `make build`: produces a local `hype` binary from `cmd/hype/`.
- `go run ./cmd/hype --help`: inspect CLI options without installing.
- `make docs` or `make hype`: regenerate `README.md` from `hype.md`; run after changing documentation templates.
- `go install ./cmd/hype`: install the CLI into your `$GOBIN` for reuse in other projects.

## Coding Style & Naming Conventions

Use Go 1.25+ features conservatively and keep compiler warnings at zero. Format code via `gofmt` (tabs for indentation, blank lines between logical sections) and lint with `revive` using `revive.toml`. Exported identifiers should read like `Parser`, `ExecuteError`, while private helpers stay lowerCamel. Keep package boundaries focused; avoid cyclic imports by adding shared helpers to `internal/`.

## Testing Guidelines

Unit tests mirror their source files (e.g., `figure_test.go`) and should cover happy paths plus failure parsing cases. Add integration coverage for the CLI in `cli_integration_test.go` when touching command behavior. Run `make cov` to inspect HTML coverage locally and target meaningful assertions rather than snapshot dumps. Table-driven tests are preferred; name cases with short strings describing the scenario. All new features require a failing test before implementation when practical.

## Commit & Pull Request Guidelines

Follow the existing history style: short imperative subject plus optional PR reference, e.g., `Fix non-deterministic JSON output (#46)`. Group related changes into one commit and keep generated files (like `README.md`) in the same commit when they stem from code changes. Pull requests should link issues, describe the behavior change, list verification steps (`make test`, `go run examples/...`), and include screenshots when user-facing output changes. Ensure CI (GitHub Actions `tests.yml`) passes before requesting review.
