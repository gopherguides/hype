# Release Process

This document describes how to create a new release for the Hype project using GoReleaser.

## Prerequisites

### 1. Install GoReleaser

```bash
# Install via Homebrew (recommended on macOS)
brew install goreleaser

# Or install via Go
go install github.com/goreleaser/goreleaser/v2@latest
```

### 2. GitHub Token Setup

You need a GitHub Personal Access Token with appropriate permissions:

1. Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Create a token with `repo` scope (for private repos) or `public_repo` scope (for public repos)
3. Export it locally:

```bash
export GITHUB_TOKEN=your_token_here
```

Tip: Add this to your `~/.zshrc` or `~/.bashrc` for persistence:

```bash
echo 'export GITHUB_TOKEN=your_token_here' >> ~/.zshrc
```

## Release Methods

### Method 1: Manual Release (Current Process)

Use this method for immediate releases or if you prefer manual control:

1. **Ensure you're on main branch and up to date:**

    ```bash
    git checkout main
    git pull origin main
    ```

2. **Verify all changes are committed:**

    ```bash
    git status
    # Should show "nothing to commit, working tree clean"
    ```

3. **Create and push the tag:**

    ```bash
    # Replace v0.3.2 with your desired version
    git tag v0.3.2
    git push origin v0.3.2
    ```

4. **Run GoReleaser:**

    ```bash
    goreleaser release --clean
    ```

### Method 2: Automated Release (Recommended for Future)

If you want automated releases via GitHub Actions:

1. **One-time setup - Commit the release workflow:**

   ```bash
   git add .github/workflows/release.yml
   git commit -m "Add automated release workflow"
   git push origin main
   ```

2. **For future releases, simply push a tag:**

   ```bash
   # Replace v0.3.3 with your desired version
   git tag v0.3.3
   git push origin v0.3.3
   ```
    
    The GitHub Action will automatically run GoReleaser and create the release.

## What Gets Created

After a successful release, you'll have:

- A new release at `https://github.com/gopherguides/hype/releases/tag/v{version}`
- Binaries for:
  - Linux (x86_64, ARM64)
  - Windows (x86_64)
  - macOS (x86_64, ARM64)
- Automatically generated changelog
- Downloadable archives for each platform (tar.gz for Unix, zip for Windows)

## Version Numbering

Follow semantic versioning (semver):

- `v1.0.0` - Major version (breaking changes)
- `v0.1.0` - Minor version (new features, backwards compatible)
- `v0.0.1` - Patch version (bug fixes, backwards compatible)

## Current Release Configuration

The release is configured via `.goreleaser.yaml`:

- Main binary: `./cmd/hype/main.go`
- Builds for: Linux, Windows, macOS
- Archives: tar.gz (Unix), zip (Windows)
- Hooks: `go mod tidy` and `go generate ./...` before build

## Troubleshooting

### GoReleaser Command Not Found

```bash
# Reinstall GoReleaser
brew install goreleaser
# Or
go install github.com/goreleaser/goreleaser/v2@latest
```

### Invalid GitHub Token

```bash
# Verify token is set
echo $GITHUB_TOKEN
# If empty, export it again
export GITHUB_TOKEN=your_token_here
```

### Tag Already Exists

```bash
# Delete local tag
git tag -d v0.3.2
# Delete remote tag
git push origin --delete v0.3.2
# Recreate tag
git tag v0.3.2
git push origin v0.3.2
```

### Build Failures

```bash
# Test the build locally first
goreleaser build --snapshot --clean
```

### Permission Denied

Ensure your GitHub token has the correct permissions:

- `public_repo` for public repositories
- `repo` for private repositories

## Quick Reference

```bash
# Complete release process (manual)
git checkout main && git pull origin main
git tag v0.3.2
git push origin v0.3.2
goreleaser release --clean

# Complete release process (automated, after workflow is set up)
git checkout main && git pull origin main
git tag v0.3.2
git push origin v0.3.2
# GitHub Actions handles the rest
```
