#!/usr/bin/env python3
import sys
import subprocess
import argparse
import json

def run_git(args):
    result = subprocess.run(["git"] + args, capture_output=True, text=True, check=True)
    return result.stdout.strip()

def parse_semver(version_str):
    # Strip leading 'v' if present
    val = version_str.lstrip('v')
    parts = val.split('.')
    if len(parts) < 3:
        return None
    try:
        major = int(parts[0])
        minor = int(parts[1])
        # Patch might contain pre-release like 3.0.0-beta.1, handle it
        patch_part = parts[2].split('-')[0]
        patch = int(patch_part)
        is_prerelease = '-' in parts[2]
        return (major, minor, patch, is_prerelease)
    except ValueError:
        return None

def main():
    parser = argparse.ArgumentParser(description="Determine the previous and next release versions.")
    parser.add_argument("remote", nargs="?", default="origin", help="Git remote to fetch tags from")
    args = parser.parse_args()

    remote = args.remote

    # 1. Fetch tags from the upstream remote
    print(f"Fetching tags from remote '{remote}'...", file=sys.stderr)
    try:
        subprocess.run(["git", "fetch", remote, "--tags", "--quiet"], check=True)
    except subprocess.CalledProcessError:
        print(f"AGENT_INSTRUCTION: Could not fetch tags from remote '{remote}'. Verify the remote exists and you have access.", file=sys.stderr)
        sys.exit(1)

    # 2. Get tags
    try:
        tags_raw = run_git(["tag", "-l", "v*.*.*"])
    except Exception as e:
        print(f"AGENT_INSTRUCTION: Error listing git tags: {str(e)}", file=sys.stderr)
        sys.exit(1)

    tags = tags_raw.splitlines()
    parsed_versions = []
    for t in tags:
        parsed = parse_semver(t)
        if parsed and not parsed[3]: # Filter out pre-releases
            parsed_versions.append((parsed[0], parsed[1], parsed[2], t))

    if not parsed_versions:
        print("AGENT_INSTRUCTION: No valid stable release tags found. Verify the repository has tags in 'vX.Y.Z' format.", file=sys.stderr)
        sys.exit(1)

    # Sort versions descending
    parsed_versions.sort(key=lambda x: (x[0], x[1], x[2]), reverse=True)
    latest_tuple = parsed_versions[0]
    previous_version = latest_tuple[3].lstrip('v')

    # Calculate next minor version
    next_major = latest_tuple[0]
    next_minor = latest_tuple[1] + 1
    next_patch = 0
    next_version = f"{next_major}.{next_minor}.{next_patch}"

    # Print results as JSON to stdout
    result = {
        "previous_version": previous_version,
        "next_version": next_version
    }
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
