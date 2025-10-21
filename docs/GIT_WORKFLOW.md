# LUMINA Git Workflow Guide

## Overview

LUMINA uses a simplified GitHub Flow with release branches for structured phase-based development.

## Branch Strategy

### Branch Types

1. **main** - Stable production releases
   - Always deployable
   - Protected branch (requires PR for changes)
   - Tagged with semantic versions (v1.0.0, v1.1.0, etc.)

2. **develop** - Integration branch for next version
   - All feature branches merge here first
   - Testing ground for combined features
   - Synced to main via release branches

3. **feature/{name}** - Individual feature development
   - Created from: `develop`
   - Merged to: `develop`
   - Naming: `feature/toc-integration`, `feature/fast-search`
   - Lifespan: 1-3 days (keep small and focused)

4. **hotfix/{version}** - Critical production fixes
   - Created from: `main`
   - Merged to: `main` AND `develop`
   - Naming: `hotfix/v1.0.2`, `hotfix/clipboard-crash`
   - Use for: Security issues, critical bugs in production

5. **release/{version}** - Release preparation
   - Created from: `develop`
   - Merged to: `main` (then tagged) AND `develop`
   - Naming: `release/v1.1.0`, `release/v2.0.0`
   - Use for: Version bumps, CHANGELOG updates, final QA

## Workflow Examples

### Phase 2 Development Workflow

```bash
# Start Phase 2 Development
git checkout main
git checkout -b develop

# Feature 1: Fast Search
git checkout -b feature/fast-search develop
# ... implement search.go, add tests ...
git add search.go search_test.go
git commit -m "feat(search): implement fuzzy file search"
git checkout develop
git merge --no-ff feature/fast-search
git branch -d feature/fast-search

# Feature 2: Content Search
git checkout -b feature/content-search develop
# ... implement ...
git commit -m "feat(search): add content search with ripgrep"
git checkout develop
git merge --no-ff feature/content-search
git branch -d feature/content-search

# When all Phase 2 features complete
git checkout -b release/v1.1.0 develop
# Update version.go, CHANGELOG.md
git commit -am "chore(release): prepare v1.1.0"
git checkout main
git merge --no-ff release/v1.1.0
git tag -a v1.1.0 -m "Phase 2: Fast Search & Content Discovery"
git push origin main --tags
git checkout develop
git merge --no-ff release/v1.1.0
git branch -d release/v1.1.0
```

### Hotfix Workflow

```bash
# Critical bug found in v1.0.1-alpha
git checkout -b hotfix/v1.0.2 v1.0.1-alpha
# ... fix clipboard crash ...
git commit -am "fix(clipboard): prevent crash on empty selection"
# Update version.go to 1.0.2
git commit -am "chore(version): bump to 1.0.2"

# Merge to main
git checkout main
git merge --no-ff hotfix/v1.0.2
git tag -a v1.0.2 -m "Hotfix: Clipboard crash fix"
git push origin main --tags

# Merge to develop
git checkout develop
git merge --no-ff hotfix/v1.0.2
git branch -d hotfix/v1.0.2
```

## Commit Message Standards

### Format

Follow Conventional Commits specification:

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, no logic change)
- **refactor**: Code refactoring (no feature or bug fix)
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (dependencies, version bumps)
- **ci**: CI/CD changes

### Scope (Optional)

Component being modified: `search`, `clipboard`, `toc`, `keybindings`, `ui`

### Examples

```bash
# Good commit messages
git commit -m "feat(search): implement fuzzy file search with scoring algorithm"
git commit -m "fix(clipboard): prevent crash when copying empty selection"
git commit -m "docs(readme): add installation instructions for Linux"
git commit -m "test(toc): add unit tests for heading extraction"
git commit -m "refactor(keybindings): replace switch with action dispatch pattern"
git commit -m "perf(rendering): optimize markdown rendering for large files"

# Multi-line commit
git commit -m "feat(search): add content search with ripgrep

Implements full-text search across markdown files using ripgrep.
Includes fuzzy matching and result highlighting.

Closes #42"
```

## Tagging Strategy

### Semantic Versioning

Format: `vMAJOR.MINOR.PATCH[-PRERELEASE]`

- **MAJOR**: Breaking changes (v2.0.0)
- **MINOR**: New features, backward compatible (v1.1.0)
- **PATCH**: Bug fixes, backward compatible (v1.0.2)
- **PRERELEASE**: Alpha, beta, rc (v1.0.0-alpha, v1.0.0-beta.1)

### Creating Tags

```bash
# Annotated tag (recommended)
git tag -a v1.1.0 -m "Phase 2: Fast Search & Content Discovery"

# Tag with detailed message
git tag -a v1.1.0 -m "$(cat <<EOF
Release v1.1.0 - Phase 2

Features:
- Fast file search with fuzzy matching
- Content search with ripgrep
- File watching and auto-reload
- Enhanced vim keybindings

Performance:
- 50% faster file tree rendering
- Search results in <100ms

Breaking Changes: None
EOF
)"

# Push tags
git push origin --tags

# List tags
git tag -l

# Show tag details
git show v1.1.0
```

## Release Checklist

Before creating a release:

- [ ] All features implemented and tested
- [ ] Unit tests passing (`go test ./...`)
- [ ] Build successful (`go build -o lumina`)
- [ ] Manual testing complete
- [ ] Documentation updated (README.md, docs/)
- [ ] CHANGELOG.md updated with all changes
- [ ] Version bumped in `version.go`
- [ ] No uncommitted changes (`git status`)
- [ ] Release branch merged to main
- [ ] Git tag created and pushed
- [ ] GitHub release created with binaries
- [ ] Release branch merged back to develop

## Rollback Procedures

### Rollback to Previous Version

```bash
# Check available versions
git tag -l

# Rollback to specific version (soft - creates new branch)
git checkout -b rollback-v1.0.1 v1.0.1-alpha
go build -o lumina

# Rollback main to previous version (destructive)
# Create backup first!
git branch backup-main-$(date +%Y%m%d-%H%M%S)
git checkout main
git reset --hard v1.0.1-alpha
git push --force-with-lease origin main  # Use with caution!
```

### Undo Last Commit (Not Pushed)

```bash
# Keep changes, undo commit
git reset --soft HEAD~1

# Discard changes, undo commit
git reset --hard HEAD~1
```

### Revert Commit (Already Pushed)

```bash
# Safe revert - creates new commit
git revert <commit-hash>
git push origin main
```

## Git Hooks

See `.git/hooks/` for automated quality checks:

- `pre-commit`: Run tests, linting, formatting checks
- `commit-msg`: Validate commit message format
- `pre-push`: Run full test suite, build verification

## Branch Protection Rules

When GitHub remote is configured:

### main branch
- Require pull request reviews (1 reviewer)
- Require status checks (tests, build)
- No direct pushes
- Require linear history

### develop branch
- Require status checks
- Allow direct pushes (for quick fixes)
- No force pushes

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Git Best Practices](https://git-scm.com/book/en/v2)

---

**Last Updated**: 2025-10-21
