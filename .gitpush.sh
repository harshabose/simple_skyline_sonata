#!/bin/bash

# gitpush.sh - Automatically commit and push changes in all submodules and main repo
# Usage: ./gitpush.sh [commit-message]

# Colors for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Set main project directory - adjust this if script is not in project root
MAIN_DIR="$(pwd)"
DEPENDENCIES_DIR="$MAIN_DIR/dependencies"
TIMESTAMP=$(date +"%d-%m-%Y %H:%M")

# Get commit message from argument or use default
COMMIT_MSG="${1:-general commit $TIMESTAMP}"

# Simple debug function
debug() {
    echo -e "${YELLOW}DEBUG: $1${NC}"
}

# Function to commit and push changes if they exist
commit_and_push() {
    local repo_dir="$1"
    local repo_name="$2"

    echo -e "${BLUE}Processing $repo_name at $repo_dir${NC}"

    # Check if repo has changes
    cd "$repo_dir"
    if git status --porcelain | grep -q .; then
        echo -e "${YELLOW}Changes detected in $repo_name${NC}"

        # Stage all changes
        git add .

        # Create commit message
        local msg="$repo_name $COMMIT_MSG"

        # Commit changes
        git commit -m "$msg"

        # Push to remote
        echo -e "${GREEN}Pushing changes to $repo_name${NC}"
        git push origin master

        echo -e "${GREEN}✓ $repo_name successfully updated${NC}"
    else
        echo -e "${GREEN}No changes in $repo_name${NC}"
    fi
}

# Main execution
echo -e "${BLUE}===== Starting Git automation script =====${NC}"
echo -e "${BLUE}Timestamp: $TIMESTAMP${NC}"

# Check if dependencies directory exists
if [ ! -d "$DEPENDENCIES_DIR" ]; then
    echo -e "${RED}Error: Dependencies directory not found at $DEPENDENCIES_DIR${NC}"
    exit 1
fi

# Process all direct submodules in dependencies directory
echo -e "${BLUE}Processing submodules...${NC}"
echo -e "${BLUE}Dependencies directory: $DEPENDENCIES_DIR${NC}"

# List the direct subdirectories to process
for SUBMODULE in "$DEPENDENCIES_DIR"/*; do
    if [ -d "$SUBMODULE" ] && [ -d "$SUBMODULE/.git" ]; then
        SUBMODULE_NAME=$(basename "$SUBMODULE")
        echo -e "${BLUE}Found submodule: $SUBMODULE_NAME${NC}"
        commit_and_push "$SUBMODULE" "$SUBMODULE_NAME"
    elif [ -d "$SUBMODULE" ]; then
        SUBMODULE_NAME=$(basename "$SUBMODULE")
        echo -e "${BLUE}Checking subdirectory: $SUBMODULE_NAME${NC}"

        # Check for nested git repositories
        for NESTED_REPO in "$SUBMODULE"/*; do
            if [ -d "$NESTED_REPO" ] && [ -d "$NESTED_REPO/.git" ]; then
                NESTED_NAME="$SUBMODULE_NAME/$(basename "$NESTED_REPO")"
                echo -e "${BLUE}Found nested submodule: $NESTED_NAME${NC}"
                commit_and_push "$NESTED_REPO" "$NESTED_NAME"
            fi
        done
    fi
done

# Finally, commit changes in the main repo
echo -e "${BLUE}Processing main repository...${NC}"
# shellcheck disable=SC2164
cd "$MAIN_DIR"

# Check if main repo has changes, including potentially updated submodule references
if git status --porcelain | grep -q .; then
    echo -e "${YELLOW}Changes detected in main repository${NC}"

    # Stage all changes
    git add .

    # Commit and push
    git commit -m "$COMMIT_MSG"
    echo -e "${GREEN}Pushing changes to main repository${NC}"
    git push origin master

    echo -e "${GREEN}✓ Main repository successfully updated${NC}"
else
    echo -e "${GREEN}No changes in main repository${NC}"
fi

echo -e "${GREEN}===== Git automation completed successfully =====${NC}"