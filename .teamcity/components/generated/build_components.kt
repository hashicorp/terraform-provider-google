/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.*
import jetbrains.buildServer.configs.kotlin.buildFeatures.GolangFeature
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep
import jetbrains.buildServer.configs.kotlin.triggers.schedule

// NOTE: in time this could be pulled out into a separate Kotlin package

// The native Go test runner (which TeamCity shells out to) will fail
// the entire test suite when a single test panics, which isn't ideal.
//
// Until that changes, we'll continue to use `teamcity-go-test` to run
// each test individually

// NOTE: this file includes Extensions of Kotlin DSL classes
// See
// - BuildFeatures           https://teamcity.jetbrains.com/app/dsl-documentation/root/build-features/index.html
// - BuildSteps              https://teamcity.jetbrains.com/app/dsl-documentation/root/build-steps/index.html
// - ParametrizedWithType    https://teamcity.jetbrains.com/app/dsl-documentation/root/parametrized-with-type/index.html
// - Triggers                https://teamcity.jetbrains.com/app/dsl-documentation/root/triggers/index.html


const val useTeamCityGoTest = false

fun BuildFeatures.Golang() {
    if (useTeamCityGoTest) {
        feature(GolangFeature {
            testFormat = "json"
        })
    }
}

fun BuildSteps.ConfigureGoEnv() {
    step(ScriptBuildStep {
        name = "Configure Go version using .go-version file"
        scriptContent = "goenv install -s \$(goenv local) && goenv rehash"
    })
}

fun BuildSteps.SetGitCommitBuildId() {
    step(ScriptBuildStep {
        name = "Set build id as shortened git commit hash"
        scriptContent = """
            #!/bin/bash
            GIT_HASH=%system.build.vcs.number%
            GIT_HASH_SHORT=${'$'}{GIT_HASH:0:7}
            echo "##teamcity[buildNumber '${'$'}{GIT_HASH_SHORT}']"
        """.trimIndent()
        // ${'$'} is required to allow creating a script in TeamCity that contains
        // parts like ${GIT_HASH_SHORT} without having Kotlin syntax issues. For more info see:
        // https://youtrack.jetbrains.com/issue/KT-2425/Provide-a-way-for-escaping-the-dollar-sign-symbol-in-multiline-strings-and-string-templates
    })
}

fun BuildSteps.TagBuildToIndicatePurpose() {
    step(ScriptBuildStep {
        name = "Set build tag to indicate if build is run automatically or manually triggered"
        scriptContent = """
            #!/bin/bash
            TRIGGERED_BY_USERNAME=%teamcity.build.triggeredBy.username%

            if [[ "${'$'}TRIGGERED_BY_USERNAME" = "n/a" ]] ; then
                echo "Build was triggered as part of automated testing. We know this because the `triggeredBy.username` value was `n/a`, value: ${'$'}{TRIGGERED_BY_USERNAME}"
                TAG="nightly-test"
                echo "##teamcity[addBuildTag '${'$'}{TAG}']"
            else
                echo "Build wasn't triggered as part of automated testing. We know this because the `triggeredBy.username` value was not `n/a`, value: ${'$'}{TRIGGERED_BY_USERNAME}"
                TAG="one-off-build"
                echo "##teamcity[addBuildTag '${'$'}{TAG}']"
            fi
        """.trimIndent()
        // ${'$'} is required to allow creating a script in TeamCity that contains
        // parts like ${GIT_HASH_SHORT} without having Kotlin syntax issues. For more info see:
        // https://youtrack.jetbrains.com/issue/KT-2425/Provide-a-way-for-escaping-the-dollar-sign-symbol-in-multiline-strings-and-string-templates
    })
}

fun BuildSteps.DownloadTerraformBinary() {
    // https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip
    var terraformUrl = "https://releases.hashicorp.com/terraform/%env.TERRAFORM_CORE_VERSION%/terraform_%env.TERRAFORM_CORE_VERSION%_linux_amd64.zip"
    step(ScriptBuildStep {
        name = "Download Terraform version %s".format(defaultTerraformCoreVersion)
        scriptContent = """
        #!/bin/bash
        mkdir -p tools
        wget -O tf.zip %s
        unzip tf.zip
        mv terraform tools/
        """.format(terraformUrl).trimIndent()
    })
}

// RunSweepers runs sweepers, and relies on set build configuration parameters
fun BuildSteps.RunSweepers(sweeperStepName : String) {
    step(ScriptBuildStep{
        name = sweeperStepName
        scriptContent = "go test -v \"%PACKAGE_PATH%\" -sweep=\"%SWEEPER_REGIONS%\" -sweep-allow-failures -sweep-run=\"%SWEEP_RUN%\" -timeout 30m"
    })
}

// RunAcceptanceTests runs tests for a given directory, using either:
// - TeamCity's test runner - stops remaining tests after a failure
// - jen20/teamcity-go-test - allows tests to continue after a failure, and requires a test binary
fun BuildSteps.RunAcceptanceTests() {
    if (useTeamCityGoTest) {
        step(ScriptBuildStep {
            name = "Run Tests"
            scriptContent = "go test -v \"%PACKAGE_PATH%\" -timeout=\"%TIMEOUT%h\" -test.parallel=\"%PARALLELISM%\" -run=\"%TEST_PREFIX%\" -json"
        })
    } else {
        step(ScriptBuildStep {
            name = "Compile Test Binary"
            scriptContent = "go test -c -o test-binary"
            workingDir = "%PACKAGE_PATH%"
        })

        step(ScriptBuildStep {
            name = "Run via jen20/teamcity-go-test"
            scriptContent = "./test-binary -test.list=\"%TEST_PREFIX%\" | teamcity-go-test -test ./test-binary -parallelism \"%PARALLELISM%\" -timeout \"%TIMEOUT%h\""
            workingDir = "%PACKAGE_PATH%"
        })
    }
}

fun ParametrizedWithType.TerraformAcceptanceTestParameters(parallelism : Int, prefix : String, timeout: String, sweeperRegions: String, sweepRun: String) {
    text("PARALLELISM", "%d".format(parallelism))
    text("TEST_PREFIX", prefix)
    text("TIMEOUT", timeout)
    text("SWEEPER_REGIONS", sweeperRegions)
    text("SWEEP_RUN", sweepRun)
}

fun ParametrizedWithType.ReadOnlySettings() {
    hiddenVariable("teamcity.ui.settings.readOnly", "true", "Requires build configurations be edited via Kotlin")
}

fun ParametrizedWithType.TerraformAcceptanceTestsFlag() {
    hiddenVariable("env.TF_ACC", "1", "Set to a value to run the Acceptance Tests")
}

fun ParametrizedWithType.TerraformCoreBinaryTesting() {
    text("env.TERRAFORM_CORE_VERSION", defaultTerraformCoreVersion, "The version of Terraform Core which should be used for testing")
    hiddenVariable("env.TF_ACC_TERRAFORM_PATH", "%system.teamcity.build.checkoutDir%/tools/terraform", "The path where the Terraform Binary is located")
}

fun ParametrizedWithType.TerraformShouldPanicForSchemaErrors() {
    hiddenVariable("env.TF_SCHEMA_PANIC_ON_ERROR", "1", "Panic if unknown/unmatched fields are set into the state")
}

fun ParametrizedWithType.WorkingDirectory(path : String) {
    text("PACKAGE_PATH", path, "", "The path at which to run - automatically updated", ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenVariable(name: String, value: String, description: String) {
    text(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenPasswordVariable(name: String, value: String, description: String) {
    password(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun Triggers.RunNightly(config: NightlyTriggerConfiguration) {
    val filter = "+:" + config.branchRef // e.g. "+:refs/heads/main"

    schedule{
        enabled = config.nightlyTestsEnabled
        branchFilter = filter
        triggerBuild = always() // Run build even if no new commits/pending changes
        withPendingChangesOnly = false
        enforceCleanCheckout = true

        schedulingPolicy = cron {
            hours = config.startHour.toString()
            timezone = "SERVER"

            dayOfWeek = config.daysOfWeek
            dayOfMonth = config.daysOfMonth
        }
    }
}

fun BuildType.addTrigger(triggerConfig: NightlyTriggerConfiguration){
    triggers {
        RunNightly(triggerConfig)
    }
}
