/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.reused

import NightlyTestsProjectId
import ProviderNameGa
import ProviderNameBeta
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import builds.*
import generated.SweepersListBeta
import generated.SweepersListGa
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun weeklyDiffTests(parentProject:String, providerName: String, vcsRoot: GitVcsRoot, config: AccTestConfiguration, cron: NightlyTriggerConfiguration): Project {

    var projectId = "${parentProject}_${NightlyTestsProjectId}"
    projectId = replaceCharsId(projectId)

    // Nightly test projects run all acceptance tests overnight
    // Here we ensure the project uses the appropriate Shared Resource to ensure no clashes between builds and/or sweepers
    var sharedResources: ArrayList<String>
    when(providerName) {
        ProviderNameGa -> sharedResources = arrayListOf(SharedResourceNameGa)
        ProviderNameBeta -> sharedResources = arrayListOf(SharedResourceNameBeta)
        else -> throw Exception("Provider name not supplied when generating a weekly diff test subproject")
    }

    // Create build configs to run acceptance tests for each package defined in packages.kt and services.kt files
    // and add cron trigger to them all
    val allPackages = getAllPackageInProviderVersion(providerName)
    val packageBuildConfigs = BuildConfigurationsForPackages(allPackages, providerName, projectId, vcsRoot, sharedResources, config, releaseDiffTest = true)
    packageBuildConfigs.forEach { buildConfiguration ->
        buildConfiguration.addTrigger(cron)
    }

    // Create build config for sweeping the nightly test project
    var sweepersList: Map<String,Map<String,String>>
    when(providerName) {
        ProviderNameGa -> sweepersList = SweepersListGa
        ProviderNameBeta -> sweepersList = SweepersListBeta
        else -> throw Exception("Provider name not supplied when generating a weekly diff test subproject")
    }
    val serviceSweeperConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperName, sweepersList, projectId, vcsRoot, sharedResources, config)
    val sweeperCron = cron.clone()
    sweeperCron.startHour += 5  // Ensure triggered after the package test builds are triggered
    serviceSweeperConfig.addTrigger(sweeperCron)

    return Project {
        id(projectId)
        name = "Weekly Diff Tests"
        description = "A project connected to the hashicorp/terraform-provider-${providerName} repository, where scheduled weekly diff tests run and users can trigger ad-hoc builds"

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