/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.AbsoluteId

class sweeperBuildConfigs() {

    fun preSweeperBuildConfig(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, environmentVariables: ClientConfiguration) : BuildType {
        val testPrefix = "TestAcc"
        val testTimeout = "12"
        val sweeperRegions = "us-central1"
        val sweeperRun = "" // Empty string means all sweepers run

        val configName = "Pre-Sweeper"
        val sweeperStepName = "Pre-Sweeper"

        return createBuildConfig(manualVcsRoot, preSweeperBuildConfigId, configName, sweeperStepName, parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun, path, environmentVariables)
   }

    fun postSweeperBuildConfig(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, triggerConfig: NightlyTriggerConfiguration, environmentVariables: ClientConfiguration, dependencies: ArrayList<String>) : BuildType {
        val testPrefix = "TestAcc"
        val testTimeout = "12"
        val sweeperRegions = "us-central1"
        val sweeperRun = "" // Empty string means all sweepers run

        val configName = "Post-Sweeper"
        val sweeperStepName = "Post-Sweeper"

        val build = createBuildConfig(manualVcsRoot, postSweeperBuildConfigId, configName, sweeperStepName, parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun, path, environmentVariables)

        build.addDependencies(dependencies) // Make dependent on package builds before starting
        build.addTrigger(triggerConfig)     // Post-Sweeper is triggered by cron, and dependencies result in other builds being queued

        return build
    }

    fun createBuildConfig(
        manualVcsRoot: AbsoluteId,
        configId: String,
        configName: String,
        sweeperStepName: String,
        parallelism: Int,
        testPrefix: String,
        testTimeout: String,
        sweeperRegions: String,
        sweeperRun: String,
        path: String,
        environmentVariables: ClientConfiguration,
        buildTimeout: Int = defaultBuildTimeoutDuration
        ) : BuildType {
        return BuildType {

            id(configId)

            name = configName

            vcs {
                root(rootId = manualVcsRoot)
                cleanCheckout = true
            }

            steps {
                SetGitCommitBuildId()
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunSweepers(sweeperStepName)
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                ConfigureGoogleSpecificTestParameters(environmentVariables)
                TerraformAcceptanceTestParameters(parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun)
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(path)
            }

            failureConditions {
                executionTimeoutMin = buildTimeout
            }

            // NOTE: dependencies and triggers are added by methods after the BuildType object is created
        }
    }
}
