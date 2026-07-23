#!/bin/bash
# 
# get_nightly_sha.sh
#
# Thin bash wrapper to execute get_nightly_sha.py.
# Fetches the nightly-test branch and queries TeamCity for the release candidate SHA.
#
# Usage:
#   $0 [REMOTE] [--day-of-week DAY_OF_WEEK]
#
# Arguments:
#   [REMOTE]                      Optional git remote name (default: "origin").
#   --day-of-week, --day [day]    Optional day of week to filter the candidate cut by.
#                                 Defaults to "thursday". Supports full names or 3-4 letter abbreviations.
#
# Examples:
#   $0 upstream
#   $0 upstream --day-of-week wednesday
#   $0 origin --day friday
#
python3 "$(dirname "$0")/get_nightly_sha.py" "$@"

