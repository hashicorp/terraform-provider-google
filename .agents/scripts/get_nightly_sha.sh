#!/bin/bash
# Thin bash wrapper to execute get_nightly_sha.py
python3 "$(dirname "$0")/get_nightly_sha.py" "$@"
