/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.reused

import NightlyTestsProjectId
import ProviderNameBeta
import ProviderNameGa
import ProviderNameBetaDiffTest
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import builds.*
import generated.SweepersListBeta
import generated.SweepersListGa
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun nightlyTests(parentProject:String, providerName: String, vcsRoot: GitVcsRoot, config: AccTestConfiguration, cron: NightlyTriggerConfiguration): Project {

    // Create unique ID for the dynamically-created project
    var projectId = "${parentProject}_${NightlyTestsProjectId}"
    projectId = replaceCharsId(projectId)

    // Nightly test projects run all acceptance tests overnight
    // Here we ensure the project uses the appropriate Shared Resource to ensure no clashes between builds and/or sweepers
    var sharedResources: ArrayList<String>
    when(providerName) {
        ProviderNameGa -> sharedResources = arrayListOf(SharedResourceNameGa)
        ProviderNameBeta -> sharedResources = arrayListOf(SharedResourceNameBeta)
        ProviderNameBetaDiffTest -> sharedResources = arrayListOf(SharedResourceNameBeta)
        else -> throw Exception("Provider name not supplied when generating a nightly test subproject")
    }

    // Create build configs to run acceptance tests for each package defined in packages.kt and services.kt files
    val allPackages = getAllPackageInProviderVersion(providerName)
    val packageBuildConfigs = BuildConfigurationsForPackages(allPackages, providerName, projectId, vcsRoot, sharedResources, config)

    // Create a composite build that runs all package tests
    val compositeConfig = jetbrains.buildServer.configs.kotlin.BuildType {
        id(replaceCharsId("${projectId}_all_tests"))
        name = "All Nightly Tests"
        type = jetbrains.buildServer.configs.kotlin.BuildTypeSettings.Type.COMPOSITE

        vcs {
            root(vcsRoot)
        }

        dependencies {
            packageBuildConfigs.forEach { bc ->
                snapshot(bc) {
                    onDependencyFailure = jetbrains.buildServer.configs.kotlin.FailureAction.ADD_PROBLEM
                    onDependencyCancel = jetbrains.buildServer.configs.kotlin.FailureAction.ADD_PROBLEM
                }
            }
        }
    }

    // Create build config for sweeping the nightly test project
    var sweepersList: Map<String,Map<String,String>>
    when(providerName) {
        ProviderNameGa -> sweepersList = SweepersListGa
        ProviderNameBeta -> sweepersList = SweepersListBeta
        ProviderNameBetaDiffTest -> sweepersList = SweepersListBeta
        else -> throw Exception("Provider name not supplied when generating a nightly test subproject")
    }
    val serviceSweeperConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperName, sweepersList, projectId, vcsRoot, sharedResources, config)
    
    // Add snapshot dependency on the composite config to run after tests finish
    serviceSweeperConfig.dependencies {
        snapshot(compositeConfig) {
            onDependencyFailure = jetbrains.buildServer.configs.kotlin.FailureAction.IGNORE
            onDependencyCancel = jetbrains.buildServer.configs.kotlin.FailureAction.IGNORE
        }
    }
    
    // Trigger the sweeper, which will recursively trigger the composite build and all package tests
    serviceSweeperConfig.addTrigger(cron)

    return Project {
        id(projectId)
        name = "Nightly Tests"
        description = "A project connected to the hashicorp/terraform-provider-${providerName} repository, where scheduled nightly tests run and users can trigger ad-hoc builds"

        // Register build configs in the project
        packageBuildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
        buildType(compositeConfig)
        buildType(serviceSweeperConfig)

        params{
            configureGoogleSpecificTestParameters(config)
        }
    }
}