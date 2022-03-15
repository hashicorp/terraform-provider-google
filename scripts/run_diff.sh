#!/bin/bash
set -e
set -x
if [ -z "$1" ]; then
  echo "Must provide 1 argument - name of resource to diff, e.g. 'google_compute_forwarding_rule'"
  exit 1
fi

function cleanup() {
  go mod edit -dropreplace=github.com/hashicorp/terraform-provider-clean-google
  go mod edit -droprequire=github.com/hashicorp/terraform-provider-clean-google
}

trap cleanup EXIT
if [[ -d ~/go/src/github.com/hashicorp/terraform-provider-clean-google ]]; then
  pushd ~/go/src/github.com/hashicorp/terraform-provider-clean-google
  git clean -fdx
  git reset --hard
  git checkout main
  git pull
  popd
else
  mkdir -p ~/go/src/github.com/hashicorp
  git clone https://github.com/hashicorp/terraform-provider-google ~/go/src/github.com/hashicorp/terraform-provider-clean-google
fi


go mod edit -require=github.com/hashicorp/terraform-provider-clean-google@v0.0.0
go mod edit -replace github.com/hashicorp/terraform-provider-clean-google=$(realpath ~/go/src/github.com/hashicorp/terraform-provider-clean-google)
go run scripts/diff.go --resource $1 --verbose
