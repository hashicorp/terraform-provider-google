/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package projects.reused

import NightlyTestsProjectId
import ProviderNameBeta
import ProviderNameGa
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import builds.*
import generated.PackagesList
import generated.SweepersList
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun nightlyTests(parentProject:String, providerName: String, vcsRoot: GitVcsRoot, config: AccTestConfiguration): Project {

    // Create unique ID for the dynamically-created project
    var projectId = "${parentProject}_${NightlyTestsProjectId}"
    projectId = replaceCharsId(projectId)

    // Nightly test projects run all acceptance tests overnight
    // Here we ensure the project uses the appropriate Shared Resource to ensure no clashes between builds and/or sweepers
    var sharedResources: ArrayList<String>
    when(providerName) {
        ProviderNameGa -> sharedResources = arrayListOf(SharedResourceNameGa)
        ProviderNameBeta -> sharedResources = arrayListOf(SharedResourceNameBeta)
        else -> throw Exception("Provider name not supplied when generating a nightly test subproject")
    }

    // Create build configs to run acceptance tests for each package defined in packages.kt and services.kt files
    val allPackages = getAllPackageInProviderVersion(providerName)
    val packageBuildConfigs = BuildConfigurationsForPackages(allPackages, providerName, projectId, vcsRoot, sharedResources, config)
    val accTestTrigger  = NightlyTriggerConfiguration()
    packageBuildConfigs.forEach { buildConfiguration ->
        buildConfiguration.addTrigger(accTestTrigger)
    }

    // Create build config for sweeping the nightly test project
    val serviceSweeperConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperName, SweepersList, projectId, vcsRoot, sharedResources, config)
    val sweeperTrigger  = NightlyTriggerConfiguration(startHour=12)  // Override hour
    serviceSweeperConfig.addTrigger(sweeperTrigger)

    return Project {
        id(projectId)
        name = "Nightly Tests"
        description = "A project connected to the hashicorp/terraform-provider-${providerName} repository, where scheduled nightly tests run and users can trigger ad-hoc builds"

        // Register build configs in the project
        packageBuildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
        buildType(serviceSweeperConfig)

        params{
            configureGoogleSpecificTestParameters(config)
        }
    }
}