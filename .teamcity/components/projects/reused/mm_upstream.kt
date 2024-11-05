/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.reused

import MMUpstreamProjectId
import ProviderNameBeta
import ProviderNameGa
import ServiceSweeperCronName
import ServiceSweeperManualName
import SharedResourceNameVcr
import builds.*
import generated.*
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun mmUpstream(parentProject: String, providerName: String, vcsRoot: GitVcsRoot, cronSweeperVcsRoot: GitVcsRoot, config: AccTestConfiguration, cron: NightlyTriggerConfiguration): Project {

    // Create unique ID for the dynamically-created project
    var projectId = "${parentProject}_${MMUpstreamProjectId}"
    projectId = replaceCharsId(projectId)

    // Shared resource allows ad hoc builds and sweeper builds to not clash
    var sharedResources: List<String> = listOf(SharedResourceNameVcr)

    // Create build configs for each package defined in packages.kt and services_ga.kt/services_beta.kt files
    val allPackages = getAllPackageInProviderVersion(providerName)
    val packageBuildConfigs = BuildConfigurationsForPackages(allPackages, providerName, projectId, vcsRoot, sharedResources, config)

    // Create build config for sweeping the VCR test project - everything except projects
    var sweepersList: Map<String,Map<String,String>>
    when(providerName) {
        ProviderNameGa -> sweepersList = SweepersListGa
        ProviderNameBeta -> sweepersList = SweepersListBeta
        else -> throw Exception("Provider name not supplied when generating a nightly test subproject")
    }

    // This build is for manually-initiated runs of sweepers, to test changes to sweepers from the upstream repo
    val serviceSweeperManualConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperManualName, sweepersList, projectId, vcsRoot, sharedResources, config)

    // This build runs on a schedule to do actual sweeping of the VCR project, using the downstream repo's code
    val serviceSweeperCronConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperCronName, sweepersList, projectId, cronSweeperVcsRoot, sharedResources, config)
    serviceSweeperCronConfig.addTrigger(cron)

    return Project {
        id(projectId)
        name = "Upstream MM Testing"
        description = "A project connected to the modular-magician/terraform-provider-${providerName} repository, to let users trigger ad-hoc builds against branches for PRs"

        // Register build configs in the project
        packageBuildConfigs.forEach { buildConfiguration: BuildType ->
            buildType(buildConfiguration)
        }
        buildType(serviceSweeperManualConfig)
        buildType(serviceSweeperCronConfig)

        params{
            configureGoogleSpecificTestParameters(config)
        }
    }
}

fun getAllPackageInProviderVersion(providerName: String): Map<String, Map<String,String>> {
    var allPackages: Map<String, Map<String, String>> = mapOf()
    if (providerName == ProviderNameGa){
        allPackages = PackagesListGa + ServicesListGa
    }
    if (providerName == ProviderNameBeta){
        allPackages = PackagesListBeta + ServicesListBeta
    }
    return allPackages
}