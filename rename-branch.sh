#!/bin/bash

# Ensure an argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <branch-name>"
    exit 1
fi

# Prefixes for the old and new branch names
OLD_PREFIX="feat"
NEW_PREFIX="closed"

# Construct the old and new branch names
OLD_BRANCH_NAME="${OLD_PREFIX}/$1"
NEW_BRANCH_NAME="${NEW_PREFIX}/$1"

# Step 1: Rename the local branch (assumes you are on a different branch already)
git branch -m $OLD_BRANCH_NAME $NEW_BRANCH_NAME

# Step 2: Delete the old branch from remote
git push origin --delete $OLD_BRANCH_NAME

# Step 3: Push the new branch name to remote and set upstream
git push origin -u $NEW_BRANCH_NAME

echo "Branch has been renamed from $OLD_BRANCH_NAME to $NEW_BRANCH_NAME and pushed to remote."
