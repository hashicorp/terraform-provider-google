#!/bin/bash
# 
# cut_release_branch.sh
# 
# This script securely handles the git branching and changelog generation for weekly releases.
# It checks out the verified candidate SHA, creates the base release branch, and then creates
# an isolated PR branch (changelog-*) for auditing. Finally, it fetches the authoritative 
# templates from magic-modules and runs changelog-gen to prepare the workspace for the human audit.
#
# Interface: $0 <COMMIT_SHA> <PREVIOUS_RELEASE_VERSION> <RELEASE_VERSION> <TARGET_REMOTE> <UPSTREAM_REMOTE>
# The agent dynamically determines the TARGET_REMOTE based on git remotes and user preference,
# and passes it explicitly to this script.
#
set -e

if [[ $# -ne 5 ]]; then
  echo "Usage: $0 <COMMIT_SHA> <PREVIOUS_RELEASE_VERSION> <RELEASE_VERSION> <TARGET_REMOTE> <UPSTREAM_REMOTE>"
  echo "Example: $0 8abcb41 7.40.0 7.41.0 ss origin"
  exit 1
fi

COMMIT_SHA=$1
PREVIOUS_RELEASE_VERSION=$2
RELEASE_VERSION=$3
REMOTE=$4
UPSTREAM_REMOTE=$5

if [[ -z "${GITHUB_TOKEN}" && -z "${JET_GITHUB_TOKEN}" ]]; then
  echo "AGENT_INSTRUCTION: GITHUB_TOKEN (or JET_GITHUB_TOKEN) is not set. changelog-gen requires this. Please ask the user to provide it."
  exit 1
fi

export GITHUB_TOKEN="${JET_GITHUB_TOKEN:-$GITHUB_TOKEN}"

# Check if release or changelog branches already exist locally
if git show-ref --verify --quiet refs/heads/release-${RELEASE_VERSION}; then
  echo "AGENT_INSTRUCTION: The local branch 'release-${RELEASE_VERSION}' already exists. Please delete it locally using 'git branch -D release-${RELEASE_VERSION}' and run the script again."
  exit 1
fi
if git show-ref --verify --quiet refs/heads/changelog-${RELEASE_VERSION}; then
  echo "AGENT_INSTRUCTION: The local branch 'changelog-${RELEASE_VERSION}' already exists. Please delete it locally using 'git branch -D changelog-${RELEASE_VERSION}' and run the script again."
  exit 1
fi

# Check if release or changelog branches already exist on remote
echo "Checking target remote ${REMOTE}..."
REMOTE_CHECK_OUT=$(git ls-remote --heads $REMOTE refs/heads/release-${RELEASE_VERSION} refs/heads/changelog-${RELEASE_VERSION} 2>&1)
CHECK_STATUS=$?
if [[ $CHECK_STATUS -ne 0 ]]; then
  echo "AGENT_INSTRUCTION: Failed to query target remote '${REMOTE}' via git ls-remote."
  echo "Error output: $REMOTE_CHECK_OUT"
  echo "Please check if your SSH keys are configured correctly or if the remote URL is correct."
  exit 1
fi

if echo "$REMOTE_CHECK_OUT" | grep -q "refs/heads/release-${RELEASE_VERSION}"; then
  echo "AGENT_INSTRUCTION: The remote branch 'release-${RELEASE_VERSION}' already exists on target remote '${REMOTE}'. Please delete it on the remote or resolve the conflict and run the script again."
  exit 1
fi

if echo "$REMOTE_CHECK_OUT" | grep -q "refs/heads/changelog-${RELEASE_VERSION}"; then
  echo "AGENT_INSTRUCTION: The remote branch 'changelog-${RELEASE_VERSION}' already exists on target remote '${REMOTE}'. Please delete it on the remote or resolve the conflict and run the script again."
  exit 1
fi

echo "Starting release cut for ${RELEASE_VERSION} from commit ${COMMIT_SHA}..."
echo "Using remote: ${REMOTE}"

REPO_NAME=$(basename $(git rev-parse --show-toplevel))
echo "Repository name: ${REPO_NAME}"

MODE_SUFFIX="ga"
if [[ "$REPO_NAME" == *"-beta" ]]; then
  MODE_SUFFIX="beta"
fi

# Fetch tags and commits from upstream
git fetch $UPSTREAM_REMOTE --tags
git fetch $UPSTREAM_REMOTE

echo "Calculating merge base for v${PREVIOUS_RELEASE_VERSION}..."
COMMIT_SHA_OF_LAST_RELEASE=$(git merge-base main v${PREVIOUS_RELEASE_VERSION} || git merge-base ${UPSTREAM_REMOTE}/main v${PREVIOUS_RELEASE_VERSION})

if [[ -z "${COMMIT_SHA_OF_LAST_RELEASE}" ]]; then
  echo "AGENT_INSTRUCTION: Could not determine merge-base for v${PREVIOUS_RELEASE_VERSION}. Verify the previous release version tag exists."
  exit 1
fi
echo "Merge base SHA: ${COMMIT_SHA_OF_LAST_RELEASE}"

echo "Checking out candidate SHA: ${COMMIT_SHA}"
git checkout ${COMMIT_SHA}

echo "Creating branch release-${RELEASE_VERSION}"
git checkout -b release-${RELEASE_VERSION}

echo "Pushing branch to ${REMOTE}"
git push -u $REMOTE release-${RELEASE_VERSION}

# Make PR branch for release notes to separate changelog edits from the raw release tag
echo "Creating PR branch changelog-${RELEASE_VERSION}"
git checkout -b changelog-${RELEASE_VERSION}

COMMIT_SHA_OF_LAST_COMMIT_IN_CURRENT_RELEASE=$(git rev-list -n 1 HEAD)

echo "Installing changelog-gen..."
go install github.com/paultyng/changelog-gen@master

echo "Fetching authoritative templates directly from magic-modules main via GitHub raw URLs..."
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT
curl -s -f -o "$TMP_DIR/changelog.tmpl" "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/.ci/changelog.tmpl" || {
  echo "AGENT_INSTRUCTION: Failed to download changelog.tmpl from magic-modules. Check internet connectivity."
  exit 1
}
curl -s -f -o "$TMP_DIR/release-note.tmpl" "https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/.ci/release-note.tmpl" || {
  echo "AGENT_INSTRUCTION: Failed to download release-note.tmpl from magic-modules. Check internet connectivity."
  exit 1
}

echo "Running changelog-gen..."
TMP_OUTPUT=$(mktemp)
echo "## ${RELEASE_VERSION} (Unreleased)" > "$TMP_OUTPUT"
echo "" >> "$TMP_OUTPUT"

changelog-gen -repo $REPO_NAME -branch main -owner hashicorp \
  -changelog "$TMP_DIR/changelog.tmpl" \
  -releasenote "$TMP_DIR/release-note.tmpl" \
  -no-note-label "changelog: no-release-note" \
  ${COMMIT_SHA_OF_LAST_RELEASE} ${COMMIT_SHA_OF_LAST_COMMIT_IN_CURRENT_RELEASE} >> "$TMP_OUTPUT"

echo "" >> "$TMP_OUTPUT"

# Check if CHANGELOG.md is missing the previous release notes (e.g. if previous release PR hasn't merged to main yet)
if ! grep -q "^## ${PREVIOUS_RELEASE_VERSION}" CHANGELOG.md; then
  echo "[INFO] Previous release v${PREVIOUS_RELEASE_VERSION} notes not found in CHANGELOG.md on main. Checking for unmerged previous release branch..."
  PREV_REF=""
  for CANDIDATE in "${UPSTREAM_REMOTE}/changelog-${PREVIOUS_RELEASE_VERSION}" "${UPSTREAM_REMOTE}/${PREVIOUS_RELEASE_VERSION}-changelog" "${REMOTE}/changelog-${PREVIOUS_RELEASE_VERSION}" "${REMOTE}/${PREVIOUS_RELEASE_VERSION}-changelog"; do
    if git rev-parse --verify "$CANDIDATE" >/dev/null 2>&1; then
      PREV_REF="$CANDIDATE"
      break
    fi
  done

  if [[ -n "$PREV_REF" ]]; then
    echo "Recovering v${PREVIOUS_RELEASE_VERSION} release notes from $PREV_REF..."
    PREV_CHANGELOG=$(git show "$PREV_REF:CHANGELOG.md" 2>/dev/null || true)
    if [[ -n "$PREV_CHANGELOG" ]]; then
      PREV_SECTION=$(echo "$PREV_CHANGELOG" | awk "/^## ${PREVIOUS_RELEASE_VERSION}/{flag=1; print; next} /^## /{if(flag && \$0 ~ /^## [0-9]+\\./) exit} flag")
      if [[ -n "$PREV_SECTION" ]]; then
        echo "$PREV_SECTION" >> "$TMP_OUTPUT"
        echo "" >> "$TMP_OUTPUT"
        echo "[SUCCESS] Recovered v${PREVIOUS_RELEASE_VERSION} release notes and inserted into changelog stream."
      fi
    fi
  else
    echo "[WARNING] Could not locate remote branch for previous release v${PREVIOUS_RELEASE_VERSION}."
  fi
fi

# Prepend the generated notes to the existing CHANGELOG.md, stripping existing empty header if present
if head -n 1 CHANGELOG.md | grep -q "^## ${RELEASE_VERSION} (Unreleased)"; then
  tail -n +3 CHANGELOG.md >> "$TMP_OUTPUT"
else
  cat CHANGELOG.md >> "$TMP_OUTPUT"
fi

mv "$TMP_OUTPUT" CHANGELOG.md

echo "Committing raw changelog..."
git add CHANGELOG.md
git commit -m "changelog: generate raw release notes for $RELEASE_VERSION"

echo "Pushing changelog branch to $REMOTE..."
git push -u $REMOTE changelog-${RELEASE_VERSION}

echo "Creating PR using GitHub API..."
REMOTE_URL=$(git remote get-url $REMOTE)
REPO_FULL=$(echo $REMOTE_URL | sed -e 's/.*github.com[:\/]//' -e 's/\.git$//')

PR_RESPONSE=$(curl -s -X POST -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$REPO_FULL/pulls \
  -d '{
    "title": "Release Notes for version '"$RELEASE_VERSION"' ('"$MODE_SUFFIX"')",
    "head": "changelog-'"$RELEASE_VERSION"'",
    "base": "release-'"$RELEASE_VERSION"'",
    "body": ""
  }')

PR_URL=$(echo "$PR_RESPONSE" | grep -o '"html_url": "[^"]*"' | head -n 1 | cut -d'"' -f4)

if [ -z "$PR_URL" ]; then
  echo "AGENT_INSTRUCTION: Failed to create PR. Check API response or if PR already exists."
  echo "API Response: $PR_RESPONSE"
  exit 1
fi

echo "SUCCESS: changelog-gen completed and PR created at $PR_URL."
echo "AGENT_INSTRUCTION: The workspace is now on branch 'changelog-${RELEASE_VERSION}' and PR has been created. PR URL: $PR_URL. Please provide this link to the user and ask if you should proceed to the Phase 4 Audit!"
