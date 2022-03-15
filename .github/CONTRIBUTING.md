# Contributing to Terraform - Google Provider

For a set of general guidelines, see the [CONTRIBUTING.md](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md) page in the main Terraform repository.

The following are certain Google Provider-specific things to be aware of when contributing.

## Go

See the [.go-version](https://github.com/hashicorp/terraform-provider-google/blob/master/.go-version) file for which version of Go to use while developing the provider. You can manage it automatically using [`goenv`](https://github.com/syndbg/goenv).

We aim to make the Google Provider a good steward of Go practices. See https://github.com/golang/go/wiki/CodeReviewComments for common Go mistakes that you should attempt to avoid.

## Generated Resources

We maintain 2 different versions of the Google Terraform provider; the [`google` provider](https://github.com/hashicorp/terraform-provider-google) and the [`google-beta` provider](https://github.com/hashicorp/terraform-provider-google-beta). The `google` provider supports GA ([general availability](https://cloud.google.com/terms/launch-stages)) features, and `google-beta` supports beta features.

We are using code generation tool called [Magic Modules](https://github.com/googleCloudPlatform/magic-modules/) that uses a shared code base to generate both providers. Some Terraform resources are fully generated, whereas some resources are hand written and located in [the third_party/terraform/ folder in magic modules](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/third_party/terraform/resources). Generated resources will have a prominent header at the top of the file identifying them. Hand written resources have a .go or .go.erb extension but will eventually be migrated into the code generation tool with the goal of having all resources fully generated.

For more details on Magic Modules please visit [the readme](https://github.com/GoogleCloudPlatform/magic-modules). For feature requests or bugs regarding those resources, please continue to file issues in the [terraform-provider-google issue tracker](https://github.com/hashicorp/terraform-provider-google/issues). PRs changing those resources directly will not be accepted, and they must be edited through Magic Modules.

## Beta vs GA providers

Fields that are only available in beta versions of the Google Cloud Platform API will need to be added only to the `google-beta` provider and excluded from the `google` provider. For generated versions, you can apply `min_version: beta`to a specific field or resource. For handwriten resources, you will need to use templating to exclude beta fields/resources from generating into the ga version. Look for `*.erb` files in [resources](https://github.com/GoogleCloudPlatform/magic-modules/tree/master/mmv1/third_party/terraform/resources) for examples.

## Tests

### Running Tests

Configuring tests is similar to configuring the provider; see the [Provider Configuration Reference](https://www.terraform.io/docs/providers/google/provider_reference.html#configuration-reference) for more details. Tests generally assume the following environment variables must be set in order to run tests:

```
GOOGLE_PROJECT
GOOGLE_CREDENTIALS|GOOGLE_CLOUD_KEYFILE_JSON|GCLOUD_KEYFILE_JSON|GOOGLE_USE_DEFAULT_CREDENTIALS
GOOGLE_REGION
GOOGLE_ZONE
```

Note that the credentials you provide must be granted wide permissions on the specified project. These tests provision real resources, and require permission in order to do so. Most developers on the team grant their test service account `roles/editor` or `roles/owner` on their project. Additionally, to ensure that your tests are performed in a region and zone with wide support for GCP features, `GOOGLE_REGION` should be set to `us-central1` and `GOOGLE_ZONE` to `us-central1-a`.

Additional variable may be required for other tests, and should get flagged when running them by Go skipping the test and flagging in the output it was skipped, with a skip message explaining why. The most typical extra values required are those required for project creation:

```
GOOGLE_ORG
GOOGLE_BILLING_ACCOUNT
```

When running tests, specify which to run using `TESTARGS`, such as:

```
TF_LOG=TRACE make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
```

The `TESTARGS` variable is regexp-like, so multiple tests can be run in parallel by specifying a common substring of those tests (for example, `TestAccContainerNodePool` to run all node pool tests). There are 1500+ tests, and running all of them takes over 6 hours and requires a lot of GCP quota.

Note: `TF_LOG=TRACE` is optional; it [enables verbose logging](https://www.terraform.io/docs/internals/debugging.html) during tests, including all API request/response cycles. `> output.log` redirects the test output to a file for analysis, which is useful because `TRACE` logging can be extremely verbose.

### Ensuring no plan-time difference to latest provider release (optional)

In a case where you are editing an existing field you might want to ensure the resource you are modifying doesn't result in a diff to existing deployments. You can run set the [environment variable](https://github.com/GoogleCloudPlatform/magic-modules/blob/a30da2040ca7b8bd37186d8521a911e7469da632/mmv1/third_party/terraform/utils/provider_test.go.erb#L284-L286) `RELEASE_DIFF` before running a test. This will append plan only steps using the latest released/published provider (`google` or `google-beta`) after all configuration deployments to ensure uniformity.

```
export RELEASE_DIFF=true
TF_LOG=TRACE make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
```

### Writing Tests

Tests should confirm that a resource can be created, and fields or resources that support update should have tests for update. Configs need to include all dependencies of the resource in a valid configuration, and should generally contain exactly one instance of the resource being tested (w/ exceptions if a pair of resources is a regular part of the workflow, such as network peerings).

Resources that are importable (~90% of resources) should have a test that confirms that every field is importable. This is added as an extra TestStep with the following format, substituting the appropriate `ResourceName` value based on the test config:

```
resource.TestStep{
	ResourceName:      "google_compute_backend_service.foobar",
	ImportState:       true,
	ImportStateVerify: true,
},
```

### Sweepers

Running provider tests often can lead to dangling test resources caused by test failures. Terraform has a capability to run [Sweepers](https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html) which can go through and delete resources. In TPG, sweepers mainly:
1. List every resource in a project of a specific kind
2. Iterate through the list and determine if a resource is [sweepable](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/third_party/terraform/utils/gcp_sweeper_test.go#L46)
3. If sweepable, delete the resource

Sweepers run by using the `-sweep` and `-sweep-run` `TESTARGS` flags:

```
TF_LOG=TRACE make testacc TEST=./google TESTARGS='-sweep=us-central1 -sweep-run=<sweeper-name-here>' > output.log
```

## Instructing terraform to use a local copy of the provider

Note that these instructions apply to `0.13+`. For prior Terraform versions, look at past versions of this page for instructions.

### Using released terraform binary with local provider binary

Setup:
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google/5.0.0/darwin_amd64
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google-beta/5.0.0/darwin_amd64
ln -s $GOPATH/bin/terraform-provider-google ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google/5.0.0/darwin_amd64/terraform-provider-google_v5.0.0
ln -s $GOPATH/bin/terraform-provider-google-beta ~/.terraform.d/plugins/registry.terraform.io/hashicorp/google-beta/5.0.0/darwin_amd64/terraform-provider-google-beta_v5.0.0
```

Once this setup is complete, terraform will automatically use the binaries generated by the `make build` commands in the `terraform-provider-google` and `terraform-provider-google-beta` repositories instead of downloading the latest versions. To undo this, you can run:

```bash
rm -rf ~/.terraform.d/plugins/registry.terraform.io/hashicorp/
```

For more information, check out Hashicorp's documentation on the [0.13+ filesystem layout](https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers).

If multiple versions are available in a plugin directory (for example after `terraform providers mirror` is used), Terraform will pick the most up-to-date provider version within version constraints. As such, we recommend using a version that is several major versions ahead for your local copy of the provider, such as `5.0.0`.

### Building terraform binary from source

If you build the terraform binary from source, it will automatically use the binaries generated by the `make build` commands in the `terraform-provider-google` and `terraform-provider-google-beta` repositories instead of downloading the latest versions.

### FAQ

* I've followed these instructions and a preexisting config is still using the wrong provider version.
  * If you've already used a release version of a provider in a given directory by running `terraform init`, Terraform will not use the locally built copy; remove the release version from the `./.terraform/` to start using your locally built copy.

# Maintainer-specific information

## Reviewing / Merging Code

When reviewing/merging code, roughly follow the guidelines set in the
[Maintainer's Etiquette](https://github.com/hashicorp/terraform/blob/master/docs/maintainer-etiquette.md)
guide. One caveat is that they're fairly old and apply primarily to HashiCorp employees, but the general guidance about merging / changelogs is still relevant.

## Upstreaming community PRs to Magic Modules

We recommend that the majority of contributors contribute directly to [Magic Modules](https://github.com/googleCloudPlatform/magic-modules/) rather than directly to this repo, as we're much more able to run automated tests against PRs in that repository. If a contributor is unable to contribute to that repository but can contribute here, we're able to semi-automatically upstream changes with the following after conducting a primary review in this repo:

When contributors update handwritten files, we've got a couple bash fns to make the process simpler. Define the following in
your `.bashrc` or `.bash_profile`.

```bash
function tpgpatch1 {
  pr_username=$(echo $1 | cut -d ':' -f1)
  feature_branch=$(echo $1 | cut -d ':' -f2)
  git remote add $pr_username git@github.com:$pr_username/${PWD##*/}
  git fetch $pr_username
  git checkout $pr_username/$feature_branch
  git format-patch $(git merge-base HEAD master)
}

function tpgpatch2 {
  for patch in $GOPATH/src/github.com/hashicorp/terraform-provider-google*/*.patch; do
    echo "checking ${patch}"
        if git apply --stat $patch | grep "google/"; then
                git am -3 -i $patch -p2 --directory=third_party/terraform/resources/ --include="*.go"
        fi
        if git apply --stat $patch | grep "google-beta/"; then
                git am -3 -i $patch -p2 --directory=third_party/terraform/resources/ --include="*.go"
        fi
        if git apply --stat $patch | grep "markdown"; then
                git am -3 -i $patch --directory=third_party/terraform/ --include="*website/*"
        fi
  done
}
```

With those functions defined:

1. Check out both the provider and MM repo to `master`, committing/stashing any local changes
1. In the MM repo, run `git checkout -b {{branch}}` to create a branch for your upstreaming PR
1. Click the clipboard button next to the `author:branch` indicator on the PR to copy it.
1. Run `tpgpatch1 author:branch` from the provider repo
1. Run `tpgpatch2` from the MM repo
1. Remove the patch files from the provider repo

At this point, you should be checked out to a branch with the changes to handwritten files included in the MM repo. For
generated files and most compiled files, you'll need to perform the upstreaming manually. After getting your local branch
ready:

1. Open a PR in MM with a complete release note, and edit the PR in this repo to include it
1. Assign a reviewer- generally they'll rubberstamp the change, since it's already been approved
    * You can ignore CLA notices for `third_party`-only changes, it's not subject to it.
1. Once approved, merge the PR here followed by the MM PR. The Magician will correct any deltas between the original code & the MM change, generally just minor formatting such as whitespace.
