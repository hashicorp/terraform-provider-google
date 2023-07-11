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
        name = "Configure Go Version"
        scriptContent = "goenv install -s \$(goenv local) && goenv rehash"
    })
}

fun BuildSteps.DownloadTerraformBinary() {
    // https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip
    var terraformUrl = "https://releases.hashicorp.com/terraform/%env.TERRAFORM_CORE_VERSION%/terraform_%env.TERRAFORM_CORE_VERSION%_linux_amd64.zip"
    step(ScriptBuildStep {
        name = "Download Terraform Core v%env.TERRAFORM_CORE_VERSION%.."
        scriptContent = "mkdir -p tools && wget -O tf.zip %s && unzip tf.zip && mv terraform tools/".format(terraformUrl)
    })
}

fun servicePath(path : String, packageName: String) : String {
    return "./%s/%s".format(path, packageName)
}

fun BuildSteps.RunAcceptanceTests(path : String, packageName: String) {
    var packagePath = servicePath(path, packageName)
    var withTestsDirectoryPath = "##teamcity[setParameter name='PACKAGE_PATH' value='%s/tests']".format(packagePath)

    // some packages use a ./tests folder, others don't - conditionally append that if needed
    step(ScriptBuildStep {
        name          = "Determine Working Directory for this Package"
        scriptContent = "if [ -d \"%s/tests\" ]; then echo \"%s\"; fi".format(packagePath, withTestsDirectoryPath)
    })

    step(ScriptBuildStep{
        name = "Pre-Sweeper"
        scriptContent = "go test -v \"%PACKAGE_PATH%\" -sweep=\"%SWEEPER_REGIONS%\"  -sweep-allow-failures -sweep-run=\"%SWEEP_RUN%\" -timeout 30m"
    })

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
            // ./test-binary -test.list=TestAccComputeRegionDisk_ | teamcity-go-test -test ./test-binary -timeout 1s
            name = "Run via jen20/teamcity-go-test"
            scriptContent = "./test-binary -test.list=\"%TEST_PREFIX%\" | teamcity-go-test -test ./test-binary -parallelism \"%PARALLELISM%\" -timeout \"%TIMEOUT%h\""
            workingDir = "%PACKAGE_PATH%"
        })
    }

    step(ScriptBuildStep{
        name = "Post-Sweeper"
        scriptContent = "go test -v \"%PACKAGE_PATH%\" -sweep=\"%SWEEPER_REGIONS%\"  -sweep-allow-failures -sweep-run=\"%SWEEP_RUN%\" -timeout 30m"
    })
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

fun ParametrizedWithType.WorkingDirectory(path : String, packageName: String) {
    text("PACKAGE_PATH", servicePath(path, packageName), "", "The path at which to run - automatically updated", ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenVariable(name: String, value: String, description: String) {
    text(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun ParametrizedWithType.hiddenPasswordVariable(name: String, value: String, description: String) {
    password(name, value, "", description, ParameterDisplay.HIDDEN)
}

fun Triggers.RunNightly(nightlyTestsEnabled: Boolean, startHour: Int, daysOfWeek: String, daysOfMonth: String) {
    schedule{
        enabled = nightlyTestsEnabled
        branchFilter = "+:refs/heads/main"

        schedulingPolicy = cron {
            hours = startHour.toString()
            timezone = "SERVER"

            dayOfWeek = daysOfWeek
            dayOfMonth = daysOfMonth
        }
    }
}
