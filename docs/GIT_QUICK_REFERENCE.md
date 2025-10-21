# LUMINA Git Quick Reference

**Quick commands for daily development**

---

## Daily Workflow

### Start New Feature

```bash
git checkout develop
git pull origin develop
git checkout -b feature/my-feature
# ... work ...
git add .
git commit -m "feat(scope): add new feature"
git checkout develop
git merge --no-ff feature/my-feature
git branch -d feature/my-feature
git push origin develop
```

### Commit Changes

```bash
# Stage changes
git add <files>

# Commit with conventional commits format
git commit -m "type(scope): description"

# Types: feat, fix, docs, style, refactor, perf, test, chore, ci
# Examples:
git commit -m "feat(search): add fuzzy file search"
git commit -m "fix(clipboard): prevent crash on empty selection"
git commit -m "docs(readme): update installation instructions"
```

### Check Status

```bash
git status                    # Working directory status
git log --oneline -10         # Recent commits
git diff                      # Unstaged changes
git diff --staged             # Staged changes
git branch -a                 # All branches
```

---

## Testing

```bash
go test ./...                 # Run all tests
go test -cover ./...          # With coverage
go test -race ./...           # Race detector
go test -bench=.              # Benchmarks
```

---

## Release Process

### Automated Release

```bash
./scripts/release.sh v1.1.0
# Follow prompts
# Review changes
git push origin main --tags
```

### Manual Release

```bash
# 1. Create release branch
git checkout -b release/v1.1.0 develop

# 2. Update version.go
vim version.go
# Version: "1.1.0"
# BuildDate: "2025-10-21"
# BuildPhase: "Phase 2"

# 3. Update CHANGELOG.md
vim CHANGELOG.md
# Add ## [1.1.0] section

# 4. Commit changes
git commit -am "chore(release): prepare v1.1.0"

# 5. Merge to main
git checkout main
git merge --no-ff release/v1.1.0

# 6. Tag release
git tag -a v1.1.0 -m "Release v1.1.0 - Phase 2"

# 7. Push
git push origin main --tags

# 8. Merge back to develop
git checkout develop
git merge --no-ff release/v1.1.0
git branch -d release/v1.1.0
```

---

## Emergency Fixes

### Hotfix

```bash
# 1. Create hotfix branch
git checkout -b hotfix/v1.0.2 v1.0.1-alpha

# 2. Fix bug
vim <file>
git commit -am "fix(scope): description"

# 3. Update version
vim version.go  # Bump to 1.0.2
git commit -am "chore(version): bump to 1.0.2"

# 4. Merge to main
git checkout main
git merge --no-ff hotfix/v1.0.2
git tag -a v1.0.2 -m "Hotfix: description"
git push origin main --tags

# 5. Merge to develop
git checkout develop
git merge --no-ff hotfix/v1.0.2
git branch -d hotfix/v1.0.2
```

### Undo Last Commit

```bash
# Keep changes (undo commit only)
git reset --soft HEAD~1

# Discard changes (undo commit and changes)
git reset --hard HEAD~1

# Already pushed (safe revert)
git revert <commit-hash>
git push origin <branch>
```

---

## Git Hooks

### Bypass Hooks (Use with caution!)

```bash
git commit --no-verify        # Skip pre-commit and commit-msg
git push --no-verify          # Skip pre-push
```

### Test Hooks

```bash
# Test pre-commit
git add <file>
git commit -m "test"

# Test commit-msg
git commit -m "invalid"       # Should fail
git commit -m "feat: valid"   # Should pass

# Test pre-push
git push origin <branch>
```

---

## Rollback

### Rollback to Tag

```bash
# View tags
git tag -l

# Checkout specific version
git checkout v1.0.1-alpha
go build -o lumina

# Return to latest
git checkout main
```

### Rollback Branch (DANGEROUS)

```bash
# Create backup first!
git branch backup-main-$(date +%Y%m%d)

# Rollback main to tag
git checkout main
git reset --hard v1.0.1-alpha
git push --force-with-lease origin main
```

---

## Branches

### List Branches

```bash
git branch                    # Local branches
git branch -r                 # Remote branches
git branch -a                 # All branches
```

### Delete Branches

```bash
git branch -d feature/name    # Delete local (safe)
git branch -D feature/name    # Force delete local
git push origin --delete feature/name  # Delete remote
```

### Clean Up Merged Branches

```bash
# Delete all merged branches except main/develop
git branch --merged main | grep -v "main\|develop" | xargs git branch -d
```

---

## Stashing

```bash
git stash                     # Stash changes
git stash list                # List stashes
git stash pop                 # Apply and remove latest stash
git stash apply stash@{0}     # Apply specific stash
git stash drop stash@{0}      # Remove specific stash
git stash clear               # Remove all stashes
```

---

## Viewing History

```bash
# Pretty log
git log --oneline --graph --all -20

# File history
git log -p <file>

# Search commits
git log --grep="keyword"
git log --author="name"

# Show commit
git show <commit-hash>
git show HEAD
git show HEAD~1               # Previous commit
```

---

## Tags

```bash
# List tags
git tag -l

# Create annotated tag
git tag -a v1.1.0 -m "Release message"

# Create lightweight tag
git tag v1.1.0

# Push tags
git push origin --tags
git push origin v1.1.0        # Single tag

# Delete tag
git tag -d v1.1.0             # Local
git push origin --delete v1.1.0  # Remote

# Show tag
git show v1.1.0
```

---

## Commit Message Templates

### Feature

```
feat(scope): add new feature

Detailed description of what this feature does
and why it's needed.

Addresses user need for X.
```

### Bug Fix

```
fix(scope): resolve issue with Y

Description of the bug and how it's fixed.

Fixes #123
```

### Breaking Change

```
feat(scope): major refactoring

BREAKING CHANGE: API changed from X to Y.
Migration guide: ...
```

### Documentation

```
docs(readme): update installation guide

Added instructions for Linux users.
Clarified prerequisites.
```

---

## Conventional Commit Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `docs` | Documentation | None |
| `style` | Formatting | None |
| `refactor` | Code restructuring | None |
| `perf` | Performance | Patch |
| `test` | Tests | None |
| `chore` | Maintenance | None |
| `ci` | CI/CD | None |

---

## Git Alias Shortcuts

Add to `~/.gitconfig`:

```ini
[alias]
    st = status
    co = checkout
    br = branch
    ci = commit
    cm = commit -m
    cam = commit -am
    lg = log --oneline --graph --all -20
    last = log -1 HEAD
    unstage = reset HEAD --
    undo = reset --soft HEAD~1
    amend = commit --amend --no-edit
```

Usage:
```bash
git st                        # git status
git co develop                # git checkout develop
git cm "feat: add feature"    # git commit -m "feat: add feature"
git lg                        # pretty log
git undo                      # undo last commit
```

---

## Emergency Contacts

**Stuck?** Check these first:
1. `git status` - What's the current state?
2. `git log --oneline -5` - Recent commits
3. `git reflog` - Find lost commits
4. `git stash` - Save work temporarily

**Broke something?**
1. `git reflog` - Find the commit before you broke it
2. `git reset --hard <commit-hash>` - Go back to that commit
3. Or ask for help before force pushing!

---

**Last Updated**: 2025-10-21
**See Also**: [docs/GIT_WORKFLOW.md](file:///Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/docs/GIT_WORKFLOW.md) for complete guide
