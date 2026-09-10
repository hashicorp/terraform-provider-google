# How to Use TeamCity

The testing environment in TeamCity can be found at https://hashicorp.teamcity.com/

Contents:
* [Projects in TeamCity](#projects-in-teamcity)
* [How to do common tasks in TeamCity](#how-to-do-common-tasks-in-teamcity)
    * [Looking at nightly test results](#looking-at-nightly-test-results)
    * [Ad-hoc testing of branches in the downstream repositories](#ad-hoc-testing-of-branches-in-the-downstream-repositories)
    * [Ad-hoc testing of branches in the upstream repositories, while reviewing a PR](#ad-hoc-testing-of-branches-in-the-upstream-repositories-while-reviewing-a-pr)



## Projects in TeamCity

This is the hierarchy of projects in TeamCity currently:

```
Google Cloud/
│
├─ Google/
│  ├─ Nightly Tests
│  ├─ Upstream MM Testing
│
├─ Google Beta/
│  ├─ Nightly Tests
│  ├─ Upstream MM Testing
│  ├─ VCR Recording
│
├─ Project Sweeper/
```

* Projects within `Google` are for testing the GA provider.
* Projects within `Google Beta` are for testing the GA provider.
* The `Project Sweeper` project contains only the sweeper for `google_project` resources.
   * This sweeper's effects aren't confined to a single GCP project, so requires special control to ensure it doesn't interfere with other builds' acceptance tests.

## How to do common tasks in TeamCity

### Looking at nightly test results

What is the `nightly-test` branch?:
* A branch created by the [`TeamCity - nightly` Github Action](https://github.com/hashicorp/terraform-provider-google/actions/workflows/teamcity-nightly-workflow.yaml) as a way to select a recent commit on main that all tests should run on
* The branch is recreated each night
* If we use the main branch directly instead: Services are separated in their own builds and are sent into an agent queue. Depending on the timing of the queue a build could potentially have a different commit from the rest if `main` received a commit prior to exiting.

The tests interact fully with Google APIs and are used to identify any breaking changes introduced by either a recent PR or a change in the API itself.

You can find the builds for nightly tests at:

* [Google > Nightly Tests](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS?branch=refs%2Fheads%2Fnightly-test&mode=builds#all-projects)
* [Google Beta > Nightly Tests](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_BETA_NIGHTLYTESTS?branch=refs%2Fheads%2Fnightly-test&mode=builds#all-projects)

These projects contain a build configuration per service package.

To view all the failed tests for a given commit:

* Navigate to the `Change Log` tab on the relevant Nightly Tests page.
* Search for the commit used in a given night's tests and click on the commit hash.
* Scroll down the page and click the `Problems & Tests` tab
* You will see a list of tests that have failed across all builds using that commit.
    * Bolded test names are tests that were new failures.
    * Strikethrough test names are tests that have since passed in subsequent builds.
    * Plain text means that the test has failed previously and hasn't since passed. 
 

### Ad-hoc testing of branches in the downstream repositories

In preparation for a release you may need to run tests on a release branch present in the downstream hashicorp/terraform-provider-google(-beta) repos.

To do this you should:

* Navigate to [`Google > Nightly Tests`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS?mode=builds#all-projects) or [`Google Beta > Nightly Tests`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_BETA_NIGHTLYTESTS#all-projects) depending on which provider you want to test
* Click on the build configuration for the service you want to test, e.g. [`Accessapproval - Acceptance Tests`](https://hashicorp.teamcity.com/buildConfiguration/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS_GOOGLE_PACKAGE_ACCESSAPPROVAL?branch=%3Cdefault%3E&buildTypeTab=overview&mode=builds)
* On the page for the service's build configuation, look in the top right of the page for two buttons, one with the text "Run" and the other containing "...". Click on the righthand-side "..." button to launch the Custom Build modal.
* Launching a [Custom Build](https://www.jetbrains.com/help/teamcity/running-custom-build.html):
  * Select which branch to test
     * By default the `Nightly Tests` projects will use the `nightly-test` branch.
     * To change to a different branch, open the `Changes` tab, and select the branch from the `Build branch` dropdown. For more info [see the documentation on build branches](https://www.jetbrains.com/help/teamcity/running-custom-build.html#Build+Branch).
  * Change parameters for the Custom Build you're preparing to launch
    * Click the `Parameters` tab
    * Edit the value of configuration parameters or environment variables as needed.
      * `TEST_PREFIX` - change this to restrict which tests are run
      * `PARALLELISM` - controls the number of tests run in parallel
      * `The version of Terraform Core which should be...` - change this to a valid version number of Terraform to be downloaded and used in the build.
* When you've finished customising your build, click "Run Build" at the bottom of the modal.



### Ad-hoc testing of branches in the upstream repositories, while reviewing a PR

When reviewing a PR you may need to run acceptance tests using the code shown in the Diff Report, present in branches in the modular-magician/terraform-provider-google(-beta) repositories.

To do this you should navigate to either of these projects and run a custom build:
* [`Google > Upstream MM Testing`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_MMUPSTREAMTESTS?branch=refs%2Fheads%2Fnightly-test&mode=builds#all-projects)
* [`Google Beta > Upstream MM Testing`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_BETA_MMUPSTREAMTESTS?branch=refs%2Fheads%2Fnightly-test&mode=builds#all-projects)

Builds in these projects will test the code present in the Modular Magician's forks and will run tests against the VCR testing project in GCP.

See the section above about how to customise and run a Custom Build.


### Triggering VCR tests to record new cassettes

Sometimes VCR cassettes need to be re-recorded by manual intervention, for example if a VCR test is failing across all PRs due to a bad cassette that isn't being replaced. Our VCR tests on PRs only use the Beta provider, so the only place to record VCR cassettes in TeamCity is:

* [`Google Beta > VCR  Recording`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_BETA_VCRRECORDING#all-projects)

There are two build configurations in this project: one for using hashicorp/terraform-provider-google-beta to run the tests, and the other for using modular-magician/terraform-provider-google-beta. Make sure to use the correct one for your use case.

The VCR recording feature allows you to run tests across multiple service packages at once, if needed. The tests are run using the same command as the [Makefile's `testacc` target](https://github.com/hashicorp/terraform-provider-google/blob/6f7a4648aef25bce130817c38556dabbe8265bc3/GNUmakefile#L17-L18), so the values you need to set in custom builds should be familiar:

* TEST - Controls which folders are scanned for tests to run. This defaults to `./google-beta/services/...` but you can make the build faster by specifying a single service package like `./google-beta/services/pubsub`, if possible.
* TESTARGS - Controls which tests are scanned for. The value defaults to `-run=%TEST_PREFIX%`, where TEST_PREFIX is `TestAcc`. You can change the value of either TEST_PREFIX or TESTARGs to achieve the same outcome.
    * When running a list of tests I recommend editing `TESTARGS` directly, e.g. changing the value to `-run=(TestAccTest1|TestAccTest2|etc...)`
* VCR_MODE - this defaults to `RECORDING`, but you can change it to `REPLAYING`.

NOTE: VCR_PATH is already set and doesn't need to be altered.

In RECORDING mode the build will run the acceptance tests and then push the recorded cassettes to a GCS Bucket. The target bucket is controlled by Context Parameters set in TeamCity.

In REPLAYING mode the build will download VCR cassettes from a GCS Bucket and run the acceptance tests using that input.

## Sweepers

### Sweeping the Nightly Test Projects

The Service Sweeper builds in [`Google > Nightly Tests`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS?mode=builds#all-projects) and [`Google Beta > Nightly Tests`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_BETA_NIGHTLYTESTS#all-projects) run every night via CRON. They are designed to not run until there are no builds testing any services in the GA or Beta nightly test GCP projects. No acceptance testing builds will start until the sweeper stops.

### Sweeping the VCR Project

The Service Sweeper builds in [`Google > Upstream MM Testing`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_MMUPSTREAMTESTS#all-projects) and [`Google Beta > Upstream MM Testing`](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_BETA_MMUPSTREAMTESTS#all-projects) run every night via CRON. They are redundant as both sweep the VCR project, but I've left them both in. They are designed to not run until there are no builds testing any services in the VCR test GCP project. No acceptance testing builds will start until the sweeper stops.

### Sweeping `google_project` Resources

When testing the GA and Beta providers we can run tests in parallel because those tests use separate GCP projects. This creates a boundary between the two test suites and ensures they don't clash. However if an acceptance test provisions `google_project` resources in the process then there is no longer a clear GA/Beta boundary based on which host project is in use. This makes sweeping up these resources tough, as there's potential to disrupt any other running build.

The `google_project` resource can be swept up safely if there are no other ongoing builds testing the GA/Beta Google providers. The `Project Sweeper` project contains a special build configuration for this sweeper that locks access to the GA/Beta/VCR GCP projects while it runs. This means the buid must wait for all other builds to stop before it starts, and while it is running no other Google-related builds can leave the queue and start running.
