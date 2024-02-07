// This file is controlled by MMv1, any changes made here will be overwritten

package projects.reused

import MMUpstreamProjectId
import ProviderNameBeta
import ProviderNameGa
import ServiceSweeperName
import SharedResourceNameVcr
import builds.*
import generated.PackagesList
import generated.ServicesListGa
import generated.ServicesListBeta
import generated.SweepersList
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun mmUpstream(parentProject: String, providerName: String, vcsRoot: GitVcsRoot, config: AccTestConfiguration): Project {

    // Create unique ID for the dynamically-created project
    var projectId = "${parentProject}_${MMUpstreamProjectId}"
    projectId = replaceCharsId(projectId)

    // Shared resource allows ad hoc builds and sweeper builds to not clash
    var sharedResources: List<String> = listOf(SharedResourceNameVcr)

    // Create build configs for each package defined in packages.kt and services_ga.kt/services_beta.kt files
    val allPackages = getAllPackageInProviderVersion(providerName)
    val packageBuildConfigs = BuildConfigurationsForPackages(allPackages, providerName, projectId, vcsRoot, sharedResources, config)

    // Create build config for sweeping the nightly test project - everything except projects
    val serviceSweeperConfig = BuildConfigurationForSweeper(providerName, ServiceSweeperName, SweepersList, projectId, vcsRoot, sharedResources, config)
    val trigger  = NightlyTriggerConfiguration()
    serviceSweeperConfig.addTrigger(trigger) // Only the sweeper is on a schedule in this project

    return Project {
        id(projectId)
        name = "MM Upstream Testing"
        description = "A project connected to the modular-magician/terraform-provider-${providerName} repository, to let users trigger ad-hoc builds against branches for PRs"

        // Register build configs in the project
        packageBuildConfigs.forEach { buildConfiguration: BuildType ->
            buildType(buildConfiguration)
        }
        buildType(serviceSweeperConfig)

        params{
            configureGoogleSpecificTestParameters(config)
        }
    }
}

fun getAllPackageInProviderVersion(providerName: String): Map<String, Map<String,String>> {
    var allPackages: Map<String, Map<String, String>> = mapOf()
    if (providerName == ProviderNameGa){
        allPackages = PackagesList + ServicesListGa
    }
    if (providerName == ProviderNameBeta){
        allPackages = PackagesList + ServicesListBeta
    }
    return allPackages
}