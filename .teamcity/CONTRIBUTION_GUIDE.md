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


## Getting started

* Generate the GA provider from Magic Modules locally
* Open the provider repo and cd into the .teamcity folder
* Run `make tools` to download dependencies
* Run `make validate` to check the code for both:
    * Errors that prevent the code building
    * Logical errors in TeamCity-specific logic, e.g. the need for unique identifiers for builds.
* Run `make tests` to run the automated tests defined in `.teamcity/tests`

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
