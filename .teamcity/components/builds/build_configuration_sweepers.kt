/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package builds

import DefaultBuildTimeoutDuration
import DefaultParallelism
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.sharedResources
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId


fun BuildConfigurationForServiceSweeper(providerName: String, sweeperName: String, packages: Map<String, Map<String, String>>, parentProjectName: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration): BuildType {
    val sweeperPackage: Map<String, String> = packages.getValue("sweeper")
    val sweeperPath: String = sweeperPackage.getValue("path").toString()

    val sweeperRun = "" // Empty string means all sweepers run
    val sweeperRegions = "us-central1"

    val s = SweeperDetails(sweeperName, parentProjectName, providerName, sweeperRun, sweeperRegions)

    val bc = s.sweeperBuildConfig(sweeperPath, vcsRoot, sharedResources, DefaultParallelism, environmentVariables)
    bc.disableProjectSweep()
    return bc
}

fun BuildConfigurationForProjectSweeper(providerName: String, sweeperName: String, packages: Map<String, Map<String, String>>, parentProjectName: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration): BuildType {
    val sweeperPackage: Map<String, String> = packages.getValue("sweeper")
    val sweeperPath: String = sweeperPackage.getValue("path").toString()

    val sweeperRun = "GoogleProject" // Name from .google/services/resourcemanager/resource_google_project_sweeper.go
    val sweeperRegions = "us-central1" // A value needs to be present, despite projects not being regional resources

    val s = SweeperDetails(sweeperName, parentProjectName, providerName, sweeperRun, sweeperRegions)

    val bc = s.sweeperBuildConfig(sweeperPath, vcsRoot, sharedResources, DefaultParallelism, environmentVariables)
    bc.enableProjectSweep()
    return bc
}

class SweeperDetails(private val sweeperName: String, private val parentProjectName: String, private val providerName: String, private val sweeperRun: String, private val sweeperRegions: String) {

    fun sweeperBuildConfig(
        path: String,
        vcsRoot: GitVcsRoot,
        sharedResources: List<String>,
        parallelism: Int,
        environmentVariables: AccTestConfiguration,
        buildTimeout: Int = DefaultBuildTimeoutDuration
    ): BuildType {

        // These hardcoded values affect the sweeper CLI command's behaviour
        val testPrefix = "TestAcc"
        val testTimeout = "12"

        return BuildType {

            id(uniqueID())

            name = sweeperName

            vcs {
                root(vcsRoot)
                cleanCheckout = true
            }

            steps {
                setGitCommitBuildId()
                tagBuildToIndicateTriggerMethod()
                configureGoEnv()
                downloadTerraformBinary()
                runSweepers(sweeperName)
            }

            features {
                golang()
                if (sharedResources.isNotEmpty()) {
                    sharedResources {
                        // When the build runs, it locks the value(s) below
                        sharedResources.forEach { sr ->
                            lockAllValues(sr)
                        }
                    }
                }
            }

            params {
                configureGoogleSpecificTestParameters(environmentVariables)
                acceptanceTestBuildParams(parallelism, testPrefix, testTimeout)
                sweeperParameters(sweeperRegions, sweeperRun)
                terraformLoggingParameters(providerName)
                terraformCoreBinaryTesting()
                terraformShouldPanicForSchemaErrors()
                readOnlySettings()
                workingDirectory(path)
            }

            artifactRules = "%teamcity.build.checkoutDir%/debug*.txt"

            failureConditions {
                errorMessage = true
                executionTimeoutMin = buildTimeout
            }

        }
    }

    private fun uniqueID(): String {
        // Replacing chars can be necessary, due to limitations on IDs
        // "ID should start with a latin letter and contain only latin letters, digits and underscores (at most 225 characters)."
        var id = "%s_%s".format(this.parentProjectName, this.sweeperName)
        return replaceCharsId(id)
    }
}
