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
        val id = "PRE_SWEEPER_TESTS"

        val buildConfig = createBuildConfig(manualVcsRoot, id, configName, sweeperStepName, parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun, path, environmentVariables)
        return buildConfig
   }

    fun postSweeperBuildConfig(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, environmentVariables: ClientConfiguration) : BuildType {
        val testPrefix = "TestAcc"
        val testTimeout = "12"
        val sweeperRegions = "us-central1"
        val sweeperRun = "" // Empty string means all sweepers run

        val configName = "Post-Sweeper"
        val sweeperStepName = "Post-Sweeper"
        val id = "POST_SWEEPER_TESTS"

        return createBuildConfig(manualVcsRoot, id, configName, sweeperStepName, parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun, path, environmentVariables)
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
        }
    }
}
