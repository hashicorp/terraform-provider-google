/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package builds

import ArtifactRules
import DefaultBuildTimeoutDuration
import DefaultParallelism
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.failureConditions.BuildFailureOnText
import jetbrains.buildServer.configs.kotlin.failureConditions.failOnText
import jetbrains.buildServer.configs.kotlin.sharedResources
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot

class VcrDetails(private val providerName: String, private val buildId: String, private val vcsRoot: GitVcsRoot, private val sharedResources: List<String>) {

    fun vcrBuildConfig(
        environmentVariables: AccTestConfiguration,
    ): BuildType {

        // These hardcoded values affect the test runner's behaviour
        val testPrefix = "TestAcc"
        val testTimeout = "12"
        val parallelism = DefaultParallelism
        val buildTimeout: Int = DefaultBuildTimeoutDuration
        val releaseDiffTest = false

        // Path is just ./google(-beta) here, whereas nightly test builds use paths like ./google/something/specific
        // This helps VCR testing builds to run tests across multiple packages
        val path = "./${providerName}"

        val repo = vcsRoot.url!!.replace("https://github.com/", "")

        return BuildType {
            id(buildId)

            name = "VCR Recording - Using ${repo}"

            vcs {
                root(vcsRoot)
                cleanCheckout = true
            }

            steps {
                setGitCommitBuildId()
                checkVcrEnvironmentVariables()
                tagBuildToIndicateVcrMode()
                configureGoEnv()
                downloadTerraformBinary()
                runVcrTestRecordingSetup()
                runVcrAcceptanceTests()
                runVcrTestRecordingSaveCassettes()
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
                vcrEnvironmentVariables(environmentVariables, providerName)
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
}