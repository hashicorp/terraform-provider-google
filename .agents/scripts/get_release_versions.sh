#!/bin/bash
# 
# get_release_versions.sh
#
# Thin bash wrapper to execute get_release_versions.py.
# Fetches tags from the upstream remote and calculates semver versions.
#
# Usage:
#   $0 [REMOTE]
#
# Arguments:
#   [REMOTE]    Optional git remote name (default: "origin").
#
# Examples:
#   $0 origin
#   $0 upstream
#

python3 "$(dirname "$0")/get_release_versions.py" "$@"
