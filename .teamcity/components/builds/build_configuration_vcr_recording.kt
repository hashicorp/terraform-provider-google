// This file is controlled by MMv1, any changes made here will be overwritten

package builds

import DefaultBuildTimeoutDuration
import DefaultParallelism
import jetbrains.buildServer.configs.kotlin.BuildType
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
                acceptanceTestBuildParams(parallelism, testPrefix, testTimeout)
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
}