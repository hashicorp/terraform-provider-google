/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is copied from mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.*
import jetbrains.buildServer.configs.kotlin.AbsoluteId

class packageDetails(packageName: String, displayName: String, providerName: String) {
    val packageName = packageName
    val displayName = displayName
    val providerName = providerName

    // buildConfiguration returns a BuildType for a service package
    // For BuildType docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/build-type/index.html
    fun buildConfiguration(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, environmentVariables: ClientConfiguration, buildTimeout: Int = defaultBuildTimeoutDuration) : BuildType {

        val testPrefix: String = "TestAcc"
        val testTimeout: String = "12"
        val sweeperRegions: String = "" // Not used
        val sweeperRun: String = "" // Not used

        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID())

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(rootId = manualVcsRoot)
                cleanCheckout = true
            }

            steps {
                SetGitCommitBuildId()
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunAcceptanceTests()
            }

            features {
                Golang()
            }

            params {
                ConfigureGoogleSpecificTestParameters(environmentVariables)
                // TODO(SarahFrench) Split TerraformAcceptanceTestParameters function into 2: one that's used for all tests/sweeper commands, and one that's specific to sweepers
                // We shouldn't be adding sweeper-specific parameters to non-sweeper builds
                TerraformAcceptanceTestParameters(parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun)
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(path)
            }

            failureConditions {
                errorMessage = true
                executionTimeoutMin = buildTimeout
            }

            // Dependencies are not set here; instead, the `sequential` block in the Project instance creates dependencies between builds
            // Triggers are not set here; the pre-sweeper at the start of the `sequential` block has a cron trigger.
        }
    }

    fun uniqueID() : String {
        // Replacing chars can be necessary, due to limitations on IDs
        // "ID should start with a latin letter and contain only latin letters, digits and underscores (at most 225 characters)." 
        var pv = this.providerName.replace("-", "").toUpperCase()
        var pkg = this.packageName.toUpperCase()

        return "%s_PACKAGE_%s".format(pv, pkg)
    }
}
