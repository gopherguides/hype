# Contributing to Hype

Thank you for your interest in contributing to Hype!

## Development Setup

1. **Prerequisites**
   - Go 1.21 or later
   - Git

2. **Clone the repository**
   ```bash
   git clone https://github.com/gopherguides/hype.git
   cd hype
   ```

3. **Build and install locally**
   ```bash
   go install ./cmd/hype
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## Branch Naming

Use these prefixes for your branches:

- `feat/` - New features (e.g., `feat/add-pdf-export`)
- `fix/` - Bug fixes (e.g., `fix/broken-include-paths`)
- `docs/` - Documentation changes (e.g., `docs/update-readme`)
- `refactor/` - Code refactoring (e.g., `refactor/simplify-parser`)

## Pull Request Workflow

1. **Create a feature branch from main**
   ```bash
   git checkout main
   git pull origin main
   git checkout -b feat/your-feature
   ```

2. **Make your changes**
   - Write clear, concise code
   - Add tests for new functionality
   - Update documentation if needed

3. **Commit your changes**
   - Use conventional commit format: `type(scope): description`
   - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
   - Example: `feat(parser): add support for custom attributes`

4. **Push and create a PR**
   ```bash
   git push -u origin feat/your-feature
   gh pr create
   ```

5. **Address review feedback**
   - Make requested changes
   - Push additional commits
   - Request re-review when ready

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Keep functions focused and small

## Testing

- Write table-driven tests when possible
- Test both success and error cases
- Use `testdata/` directories for test fixtures

## Documentation

The README is generated from `hype.md` using hype itself:

```bash
hype export -format=markdown -f hype.md > README.md
```

If you modify documentation in `docs/` or `hype.md`, regenerate the README.

## Questions?

Open an issue if you have questions or need help getting started.
