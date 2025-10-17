/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package builds

import ArtifactRules
import DefaultBuildTimeoutDuration
import DefaultParallelism
import generated.ServiceParallelism
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.failureConditions.BuildFailureOnText
import jetbrains.buildServer.configs.kotlin.failureConditions.failOnText
import jetbrains.buildServer.configs.kotlin.sharedResources
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

// BuildConfigurationsForPackages accepts a map containing details of multiple packages in a provider and returns a list of build configurations for them all.
// Intended to be used in projects where we're testing all packages, e.g. the nightly test projects
fun BuildConfigurationsForPackages(packages: Map<String, Map<String, String>>, providerName: String, parentProjectName: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration, testPrefix: String = "TestAcc", releaseDiffTest: Boolean = false): List<BuildType> {
    val list = ArrayList<BuildType>()

    // Create build configurations for all packages, except sweeper
    packages.forEach { (packageName, info) ->
        val path: String = info.getValue("path").toString()
        val displayName: String = info.getValue("displayName").toString()

        val pkg = PackageDetails(packageName, displayName, providerName, parentProjectName, releaseDiffTest)
        val buildConfig = pkg.buildConfiguration(path, vcsRoot, sharedResources, environmentVariables, testPrefix = testPrefix, releaseDiffTest = releaseDiffTest)
        list.add(buildConfig)
    }

    return list
}

// BuildConfigurationForSinglePackage accepts details of a single package in a provider and returns a build configuration for it
// Intended to be used in short-lived projects where we're testing specific packages, e.g. feature branch testing
fun BuildConfigurationForSinglePackage(packageName: String, packagePath: String, packageDisplayName: String, providerName: String, parentProjectName: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration, testPrefix: String = "TestAcc", releaseDiffTest: Boolean = false): BuildType{
    val pkg = PackageDetails(packageName, packageDisplayName, providerName, parentProjectName, releaseDiffTest = releaseDiffTest)
    return pkg.buildConfiguration(packagePath, vcsRoot, sharedResources, environmentVariables, testPrefix = testPrefix, releaseDiffTest = releaseDiffTest)
}

class PackageDetails(private val packageName: String, private val displayName: String, private val providerName: String, private val parentProjectName: String, private val releaseDiffTest: Boolean) {

    // buildConfiguration returns a BuildType for a service package
    // For BuildType docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/build-type/index.html
    fun buildConfiguration(path: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration, buildTimeout: Int = DefaultBuildTimeoutDuration, testPrefix: String, releaseDiffTest: Boolean): BuildType {
        val testPrefix = "TestAcc"
        val testTimeout = "12"

        var parallelism = DefaultParallelism
        if (ServiceParallelism.containsKey(packageName)){
            parallelism = ServiceParallelism.getValue(packageName)
        }

        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID())

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(vcsRoot)
                cleanCheckout = true
            }

            steps {
                setGitCommitBuildId()
                tagBuildToIndicateTriggerMethod()
                configureGoEnv()
                downloadTerraformBinary()
                if (releaseDiffTest) {
                    runDiffTests()
                } else {
                    runAcceptanceTests()
                }
                saveArtifactsToGCS()
                archiveArtifactsIfOverLimit() // Must be after push to GCS step, as this step impacts debug log files
            }

            features {
                golang()
                if (sharedResources.isNotEmpty()) {
                    sharedResources {
                        // When the build runs, it locks the value(s) below
                        sharedResources.forEach { sr ->
                            lockSpecificValue(sr, packageName)
                        }
                    }
                }
            }

            params {
                configureGoogleSpecificTestParameters(environmentVariables)
                acceptanceTestBuildParams(parallelism, testPrefix, testTimeout, releaseDiffTest)
                terraformLoggingParameters(environmentVariables, providerName)
                terraformCoreBinaryTesting()
                terraformShouldPanicForSchemaErrors()
                readOnlySettings()
                workingDirectory(path)
            }

            artifactRules = ArtifactRules

            failureConditions {
                errorMessage = true
                executionTimeoutMin = buildTimeout

                // Stop builds if the branch does not exist
                failOnText {
                  conditionType = BuildFailureOnText.ConditionType.CONTAINS
                  pattern = "which does not correspond to any branch monitored by the build VCS roots"
                  failureMessage = "Error: The branch %teamcity.build.branch% does not exist"
                  reverse = false
                  stopBuildOnFailure = true
                }
            }

        }
    }

    private fun uniqueID(): String {
        // Replacing chars can be necessary, due to limitations on IDs
        // "ID should start with a latin letter and contain only latin letters, digits and underscores (at most 225 characters)."
        var id = "%s_%s_PACKAGE_%s".format(this.parentProjectName, this.providerName, this.packageName)
        return replaceCharsId(id)
    }
}