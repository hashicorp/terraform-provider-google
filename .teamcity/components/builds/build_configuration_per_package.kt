/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package builds

import DefaultBuildTimeoutDuration
import DefaultParallelism
import generated.ServiceParallelism
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.sharedResources
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

// BuildConfigurationsForPackages accepts a map containing details of multiple packages in a provider and returns a list of build configurations for them all.
// Intended to be used in projects where we're testing all packages, e.g. the nightly test projects
fun BuildConfigurationsForPackages(packages: Map<String, Map<String, String>>, providerName: String, parentProjectName: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration): List<BuildType> {
    val list = ArrayList<BuildType>()

    // Create build configurations for all packages, except sweeper
    packages.forEach { (packageName, info) ->
        val path: String = info.getValue("path").toString()
        val displayName: String = info.getValue("displayName").toString()

        val pkg = PackageDetails(packageName, displayName, providerName, parentProjectName)
        val buildConfig = pkg.buildConfiguration(path, vcsRoot, sharedResources, environmentVariables)
        list.add(buildConfig)
    }

    return list
}

// BuildConfigurationForSinglePackage accepts details of a single package in a provider and returns a build configuration for it
// Intended to be used in short-lived projects where we're testing specific packages, e.g. feature branch testing
fun BuildConfigurationForSinglePackage(packageName: String, packagePath: String, packageDisplayName: String, providerName: String, parentProjectName: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration): BuildType{
    val pkg = PackageDetails(packageName, packageDisplayName, providerName, parentProjectName)
    return pkg.buildConfiguration(packagePath, vcsRoot, sharedResources, environmentVariables)
}

class PackageDetails(private val packageName: String, private val displayName: String, private val providerName: String, private val parentProjectName: String) {

    // buildConfiguration returns a BuildType for a service package
    // For BuildType docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/build-type/index.html
    fun buildConfiguration(path: String, vcsRoot: GitVcsRoot, sharedResources: List<String>, environmentVariables: AccTestConfiguration, buildTimeout: Int = DefaultBuildTimeoutDuration): BuildType {

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
                runAcceptanceTests()
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

    private fun uniqueID(): String {
        // Replacing chars can be necessary, due to limitations on IDs
        // "ID should start with a latin letter and contain only latin letters, digits and underscores (at most 225 characters)."
        var id = "%s_%s_PACKAGE_%s".format(this.parentProjectName, this.providerName, this.packageName)
        return replaceCharsId(id)
    }
}
