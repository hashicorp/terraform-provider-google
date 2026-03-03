#!/usr/bin/env bash
# scripts/get_unit_test_pkgs.sh

# Ask Go for BOTH the Import Path and the exact Directory Path on disk, separated by a pipe '|'
go list -e -f '{{.ImportPath}}|{{.Dir}}' ./... | grep -v "/scripts" | while IFS='|' read -r pkg dir; do

    # 1. For NON-SERVICE packages:
    if [[ "$pkg" != *"/google/services/"* ]]; then
        # Check if any file in this directory has a test.
        if grep -q "^func Test" "$dir"/*_test.go 2>/dev/null; then
            echo "$pkg"
        fi
        continue
    fi

    # 2. For SERVICE packages: 
    # Logic: If line starts with "func Test" AND does NOT start with "func TestAcc", exit with success immediately.
    if awk '/^func Test/ && !/^func TestAcc/ { found=1; exit } END { exit !found }' "$dir"/*_test.go 2>/dev/null; then
        echo "$pkg"
    fi

done