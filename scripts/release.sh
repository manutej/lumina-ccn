#!/bin/bash
# LUMINA Release Automation Script
# Usage: ./scripts/release.sh [version]
# Example: ./scripts/release.sh v1.1.0

set -e  # Exit on error

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get version from argument or prompt
if [ -z "$1" ]; then
    echo -e "${YELLOW}Enter version (e.g., v1.1.0):${NC}"
    read VERSION
else
    VERSION=$1
fi

# Validate version format
if ! echo "$VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[a-z0-9]+)?$'; then
    echo -e "${RED}Error: Invalid version format. Use vX.Y.Z or vX.Y.Z-suffix${NC}"
    echo "Examples: v1.1.0, v1.0.2-alpha, v2.0.0-beta.1"
    exit 1
fi

echo -e "${BLUE}=== LUMINA Release Automation ===${NC}"
echo -e "Version: ${GREEN}$VERSION${NC}"
echo ""

# 1. Check working directory is clean
echo -e "${BLUE}[1/10]${NC} Checking working directory..."
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Working directory is not clean${NC}"
    git status --short
    echo ""
    echo "Commit or stash your changes before releasing."
    exit 1
fi
echo -e "${GREEN}✓ Working directory is clean${NC}"
echo ""

# 2. Check we're on develop or release branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo -e "${BLUE}[2/10]${NC} Checking current branch..."
if [[ "$CURRENT_BRANCH" != "develop" && "$CURRENT_BRANCH" != release/* ]]; then
    echo -e "${YELLOW}Warning: Not on develop or release branch (current: $CURRENT_BRANCH)${NC}"
    echo "Create release from develop branch? (y/n)"
    read -r response
    if [[ "$response" != "y" ]]; then
        exit 1
    fi
    git checkout develop
fi
echo -e "${GREEN}✓ On branch: $CURRENT_BRANCH${NC}"
echo ""

# 3. Pull latest changes
echo -e "${BLUE}[3/10]${NC} Pulling latest changes..."
git pull origin $CURRENT_BRANCH || echo -e "${YELLOW}No remote configured, skipping pull${NC}"
echo -e "${GREEN}✓ Up to date${NC}"
echo ""

# 4. Run tests
echo -e "${BLUE}[4/10]${NC} Running tests..."
if ! go test ./... ; then
    echo -e "${RED}Error: Tests failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ All tests passed${NC}"
echo ""

# 5. Run linting
echo -e "${BLUE}[5/10]${NC} Running go vet..."
if ! go vet ./... ; then
    echo -e "${RED}Error: go vet failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Linting passed${NC}"
echo ""

# 6. Build binary
echo -e "${BLUE}[6/10]${NC} Building binary..."
if ! go build -o lumina ; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# 7. Update version.go
echo -e "${BLUE}[7/10]${NC} Updating version.go..."
VERSION_NUMBER=${VERSION#v}  # Remove 'v' prefix
BUILD_DATE=$(date +%Y-%m-%d)

# Determine phase based on version
if [[ "$VERSION" == *"alpha"* ]]; then
    PHASE="Alpha"
elif [[ "$VERSION" == *"beta"* ]]; then
    PHASE="Beta"
else
    # Parse version to determine phase
    MAJOR=$(echo $VERSION_NUMBER | cut -d. -f1)
    MINOR=$(echo $VERSION_NUMBER | cut -d. -f2)
    if [ "$MINOR" -eq 0 ]; then
        PHASE="Phase $MAJOR"
    else
        PHASE="Phase $MAJOR.$MINOR"
    fi
fi

cat > version.go << EOF
package main

const (
    Version   = "$VERSION_NUMBER"
    BuildDate = "$BUILD_DATE"
    BuildPhase = "$PHASE"
)
EOF

git add version.go
echo -e "${GREEN}✓ version.go updated${NC}"
echo ""

# 8. Update CHANGELOG.md
echo -e "${BLUE}[8/10]${NC} Preparing CHANGELOG.md..."
echo -e "${YELLOW}Please update CHANGELOG.md with release notes for $VERSION${NC}"
echo "Press Enter when ready..."
read

# Verify CHANGELOG was updated
if ! grep -q "## \[$VERSION_NUMBER\]" CHANGELOG.md; then
    echo -e "${RED}Error: CHANGELOG.md not updated with [$VERSION_NUMBER]${NC}"
    echo "Add a section for this version before continuing."
    exit 1
fi

git add CHANGELOG.md
echo -e "${GREEN}✓ CHANGELOG.md ready${NC}"
echo ""

# 9. Create release commit
echo -e "${BLUE}[9/10]${NC} Creating release commit..."
git commit -m "chore(release): prepare $VERSION

- Update version to $VERSION_NUMBER
- Update CHANGELOG.md
- Build date: $BUILD_DATE"

echo -e "${GREEN}✓ Release commit created${NC}"
echo ""

# 10. Create and push tag
echo -e "${BLUE}[10/10]${NC} Creating git tag..."

TAG_MESSAGE="Release $VERSION - $PHASE

$(awk "/## \[$VERSION_NUMBER\]/,/## \[/" CHANGELOG.md | sed '$d' | tail -n +2)"

git tag -a "$VERSION" -m "$TAG_MESSAGE"
echo -e "${GREEN}✓ Tag $VERSION created${NC}"
echo ""

# Summary
echo -e "${GREEN}=== Release Prepared Successfully ===${NC}"
echo ""
echo "Next steps:"
echo ""
echo "1. Review the changes:"
echo -e "   ${BLUE}git show HEAD${NC}"
echo -e "   ${BLUE}git show $VERSION${NC}"
echo ""
echo "2. If everything looks good, push to remote:"
echo -e "   ${BLUE}git push origin $CURRENT_BRANCH${NC}"
echo -e "   ${BLUE}git push origin --tags${NC}"
echo ""
echo "3. If on a release branch, merge to main:"
echo -e "   ${BLUE}git checkout main${NC}"
echo -e "   ${BLUE}git merge --no-ff $CURRENT_BRANCH${NC}"
echo -e "   ${BLUE}git push origin main --tags${NC}"
echo ""
echo "4. GitHub Actions will automatically:"
echo "   - Run tests on all platforms"
echo "   - Build binaries"
echo "   - Create GitHub release"
echo "   - Upload release artifacts"
echo ""
echo -e "${YELLOW}To undo this release preparation:${NC}"
echo -e "   ${BLUE}git reset --hard HEAD~1${NC}"
echo -e "   ${BLUE}git tag -d $VERSION${NC}"
echo ""
