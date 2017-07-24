#!/bin/bash

# This script rewrites [GH-nnnn]-style references in the CHANGELOG.md file to
# be Markdown links to the given github issues.
#
# This is run during releases so that the issue references in all of the
# released items are presented as clickable links, but we can just use the
# easy [GH-nnnn] shorthand for quickly adding items to the "Unrelease" section
# while merging things between releases.

set -e

if [[ ! -f CHANGELOG.md ]]; then
  echo "ERROR: CHANGELOG.md not found in pwd."
  echo "Please run this from the root of the terraform provider repository"
  exit 1
fi

if [[ `uname` == "Darwin" ]]; then
  echo "Using BSD sed"
  SED="sed -i.bak -E -e"
else
  echo "Using GNU sed"
  SED="sed -i.bak -r -e"
fi

PROVIDER_URL="https:\/\/github.com\/terraform-providers\/terraform-provider-google\/issues"

$SED "s/GH-([0-9]+)/\[#\1\]\($PROVIDER_URL\/\1\)/g" -e 's/\[\[#(.+)([0-9])\)]$/(\[#\1\2))/g' CHANGELOG.md

rm CHANGELOG.md.bak
