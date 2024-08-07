# Contributing to the TeamCity configuration

## Background

The `.teamcity/` folder contains files that define what resources are present in our TeamCity environment.

We use TeamCity's Kotlin DSL to configure our TeamCity projects:
* [Kotlin DSL documentation](https://www.jetbrains.com/help/teamcity/kotlin-dsl.html)
* [Kotlin DSL reference documentation](https://teamcity.jetbrains.com/app/dsl-documentation/index.html)

For general information, [look at the TeamCity documentation](https://www.jetbrains.com/help/teamcity/teamcity-documentation.html). Whilst we use TeamCity Cloud, the TeamCity On Premise documentation is still relevant for understanding concepts.

## Environment

Note: these instructions need to be tested and improved. Please contact @SarahFrench (e.g. [open a GitHub issue](https://github.com/hashicorp/terraform-provider-google/issues/new?assignees=&labels=technical-debt&projects=&template=11_developer_productivity.md) and tag me) for help!

You will need to install:
* Java 17
    * `brew install openjdk@17`
* Maven
    * `brew install --ignore-dependencies maven`

Add the following to `~/.zshrc` and reload your terminal:

```
export JAVA_HOME=/usr/local/Cellar/openjdk@17/17.0.9/libexec/openjdk.jdk/Contents/Home
```


## Getting started

* Generate the GA provider from Magic Modules locally
* Open the provider repo and cd into the .teamcity folder
* Run `make tools` to download dependencies
* Run `make validate` to check the code for both:
    * Errors that prevent the code building
    * Logical errors in TeamCity-specific logic, e.g. the need for unique identifiers for builds.
* Run `make test` to run the automated tests defined in `.teamcity/tests`

## Rough description of the code base

Note: this is likely to go out of date, so this description is kept at a high level.

```
.teamcity/
│
├─ components/
│  ├─ builds/
│  │  # Files related to build configurations and the individual components of builds,
│  │  # e.g. how they're triggered, how parameters are set, build step definitions...
│  ├─ inputs/
│  │  # Files containing information about the packages in the providers, both GA and Beta,
│  │  # There are also files that can supply information about how those packages should be handled, 
│  │  # e.g. non-default parallelism values that should be used for specific packages.
│  ├─ projects/
│  │  # Files containing information about the projects created in the configuration.
│  │  # Files direcly inside the projects folder define the sub projects that will be created inside
│  │  # the parent project- Google, Google Beta, and Project Sweeper.
│  │  # The root_project.kt file defines the parent project where versioned settings is enabled.
│  │  ├─ reused/
│  │     # Code for dynamically generating subprojects within Google, Google Beta, and Project Sweeper projects.
│  ├─ vcs_roots/
│  │  # Definitions of VCS roots used to pull the provider code from GA/Beta versions of HashiCorp
│  │  # and Modular Magician repos.
│  ├─ constants.kt
│  │  # Global constants used in the above files
│  ├─ unique_id.kt
│     # A util function that's reused a lot - the beginnings of a utils file...
│
├─ tests/
│  # Test files
│
├─ pom.xml
│  # Dependencies
├─ settings.kts
│  # Entrypoint for the configuration as a whole. Where Context Parameters are fed into the Kotlin code as inputs.
├─ Makefile
   # Makefile is essential for testing your code changes!
# There are other misc files here, e.g. .gitignore
```

## Feature branch testing

If you want to test a feature branch on a schedule ahead of a release you can update the TeamCity configuration to include a new project specifically for testing that feature branch.

### Testing major releases

#### Adding a testing project for a major release

First, make sure that the feature branch `FEATURE-BRANCH-major-release-X.0.0` is created in the downstream TPG and TPGB repositories, where X is the major version.

See this PR as an example of adding a major release testing project: https://github.com/SarahFrench/magic-modules/pull/9/files

That PR creates a new file at `.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-X.0.0.kt` (replacing `X` with the version number). This file defines a new project that will contain all the builds run against the feature branch. See [FEATURE-BRANCH-major-release-6.0.0.kt](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt) as an example.

This file must:

* [Define a function that returns an instance of Project](https://github.com/GoogleCloudPlatform/magic-modules/blob/30ab2a2eea61cc34f439ddfe7cf840abf746ab1f/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt#L50)
* [The project should include sub projects for Google and Google Beta, and inside each use reusable code to provision nightly testing projects](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt#L59-L97)
    * Note: Including the Google and Google Beta projects is done to avoid two projects with the name "Nightly Tests" existing side-by-side. Names have to be unique within the scope of a containing project.
* For VCS roots, [create 2 new roots that default to using the feature branch from the TPG and TPGB repos](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt#L22-L38).
   * By using a new VCS root, instead of importing [existing roots](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/.teamcity/components/vcs_roots/vcs_roots.kt#L14-L32), you can ensure that the new project will default to using the feature branch even when builds are triggered manually.
   * Be aware of how ["logical branch names"](https://www.jetbrains.com/help/teamcity/working-with-feature-branches.html#Logical+Branch+Name) may interact with the branch filter on the cron trigger. Always manually test your configuration is able to launch builds using that trigger!
* For both nightly test projects you need to pass in a argument describing the CRON trigger that will be used to trigger builds.
    * When testing major release feature branches we tend to:
        * Trigger tests on the GA version of the [feature branch on Thursdays](https://github.com/GoogleCloudPlatform/magic-modules/blob/30ab2a2eea61cc34f439ddfe7cf840abf746ab1f/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt#L72)
        * Trigger tests on the Beta version of the [feature branch on Fridays](https://github.com/GoogleCloudPlatform/magic-modules/blob/30ab2a2eea61cc34f439ddfe7cf840abf746ab1f/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt#L92)
        * The non-feature branch projects will need to be updated to run on all days except these!
    * You'll also need to [pass the feature branch name into the CRON trigger config class](https://github.com/GoogleCloudPlatform/magic-modules/blob/2778e6b73d802c6709d10d56fc3b8a3891168e6e/mmv1/third_party/terraform/.teamcity/components/projects/feature_branches/FEATURE-BRANCH-major-release-6.0.0.kt#L71). Note that the string needs to start with `refs/heads/`.
* Don't forget to update the files that define the long-lived nightly test projects, making their CRON schedules the opposite of what's described for the feature branch testing projects:


```diff
// .teamcity/components/projects/google_ga_subproject.kt

- subProject(nightlyTests(gaId, ProviderNameGa, HashiCorpVCSRootGa, gaConfig))
+ subProject(nightlyTests(gaId, ProviderNameGa, HashiCorpVCSRootGa, gaConfig, NightlyTriggerConfiguration(daysOfWeek="1-4,6-7"))) // All nights except Thursday (5) for GA; feature branch testing happens on Thursdays and TeamCity numbers days Sun=1...Sat=7

// .teamcity/components/projects/google_beta_subproject.kt

- subProject(nightlyTests(betaId, ProviderNameBeta, HashiCorpVCSRootBeta, betaConfig))
+ subProject(nightlyTests(betaId, ProviderNameBeta, HashiCorpVCSRootBeta, betaConfig, NightlyTriggerConfiguration(daysOfWeek="1-5,7"))) // All nights except Friday (6) for Beta; feature branch testing happens on Fridays and TeamCity numbers days Sun=1...Sat=7
```

Finally, you need to register the new feature branch testing project in [the root "Google Cloud" project](https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud). Update [root_project.kt](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/.teamcity/components/projects/root_project.kt) to import code from your feature branch specific file, and [register a new sub project in the root project](https://github.com/GoogleCloudPlatform/magic-modules/blob/fd1a2272507d09214cf225b2ac05dfb363d3fb98/mmv1/third_party/terraform/.teamcity/components/projects/root_project.kt#L67-L68):

```java
Project {
    description = "Contains all testing projects for the GA and Beta versions of the Google provider."

    ...

    // Feature branch testing
    subProject(featureBranchMajorRelease600_Project(allConfig)) // FEATURE-BRANCH-major-release-6.0.0

}
```

Don't forget to check that the code builds by running `make test` in the `.teamcity` folder while you're making these changes.

To test your changes before merging a PR in magic-modules to update the TeamCity configuration, you can create a new TeamCity project that pulls its config from [the `auto-pr-N` branch for your PR in the modular-magician/terraform-provider-google repo](https://github.com/GoogleCloudPlatform/magic-modules/pull/11104#issuecomment-2206785710). See [USE_CONFIG_WITH_TEAMCITY.md](./USE_CONFIG_WITH_TEAMCITY.md) for details on creating new projects.

#### Removing a test project, after a major release

Once a major release is out you can safely delete the new project in TeamCity, and return the cron schedules to normal (i.e. every day of the week testing main branch).

Here is PR illustrating those changes: https://github.com/SarahFrench/magic-modules/pull/8

### Other feature branches

You can do the above for feature branches that aren't major releases! However depending on the feature it may make sense to limit what builds run in that project; major releases need us to test all services, but is that also true for your feature branch?

An example you can look at is this [PR](https://github.com/GoogleCloudPlatform/magic-modules/pull/10088) that creates a project within TeamCity for testing provider-defined functions. The project only ran builds for running acceptance tests for provider-defined functions and we did not need to be mindful of cron schedules as the provider-defined function acceptance tests did not have the ability to conflict with nightly tests.
