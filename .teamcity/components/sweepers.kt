/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.AbsoluteId

class sweeperDetails() {

    fun sweeperBuildConfig(
        sweeperName: String,
        path: String,
        providerName: String,
        manualVcsRoot: AbsoluteId,
        parallelism: Int,
        environmentVariables: ClientConfiguration,
        buildTimeout: Int = defaultBuildTimeoutDuration
        ) : BuildType {

        // These hardcoded values affect the sweeper CLI command's behaviour
        val testPrefix: String = "TestAcc"
        val testTimeout: String = "12"
        val sweeperRegions: String = "us-central1"
        val sweeperRun: String = "" // Empty string means all sweepers run
        
        return BuildType {

            id(createID(sweeperName))

            name = sweeperName

            vcs {
                root(rootId = manualVcsRoot)
                cleanCheckout = true
            }

            steps {
                SetGitCommitBuildId()
                TagBuildToIndicatePurpose()
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunSweepers(sweeperName)
            }

            features {
                Golang()
            }

            params {
                ConfigureGoogleSpecificTestParameters(environmentVariables)
                TerraformAcceptanceTestParameters(parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun)
                TerraformLoggingParameters()
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(path)
            }

            artifactRules = "%teamcity.build.checkoutDir%/debug*.txt"

            failureConditions {
                errorMessage = true
                executionTimeoutMin = buildTimeout
            }

            // NOTE: dependencies and triggers are added by methods after the BuildType object is created
        }
    }

    fun createID(name: String) : String {
        // Replacing chars can be necessary, due to limitations on IDs
        // "ID should start with a latin letter and contain only latin letters, digits and underscores (at most 225 characters)." 
        return name.replace("-", "_").toUpperCase()
    }
}
