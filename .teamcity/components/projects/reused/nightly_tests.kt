/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.reused

import AllNightlyTestsName
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
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
import jetbrains.buildServer.configs.kotlin.DslContext
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.triggers.finishBuildTrigger
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
    // Package builds use per-service shared-resource locks to avoid clashes with ad hoc builds.
    val packageBuildConfigs = BuildConfigurationsForPackages(allPackages, providerName, projectId, vcsRoot, sharedResources, config)

    // Create a composite build that runs all package tests
    val compositeId = replaceCharsId("${projectId}_all_tests")
    val compositeConfig = BuildType {
        id(compositeId)
        name = AllNightlyTestsName
        type = BuildTypeSettings.Type.COMPOSITE

        vcs {
            root(vcsRoot)
            cleanCheckout = true
        }

        dependencies {
            packageBuildConfigs.forEach { bc ->
                snapshot(bc) {
                    onDependencyFailure = FailureAction.ADD_PROBLEM
                    onDependencyCancel = FailureAction.ADD_PROBLEM
                }
            }
        }
    }
    compositeConfig.addTrigger(cron)

    // Create build config for sweeping the nightly test project
    var sweepersList: Map<String,Map<String,String>>
    when(providerName) {
        ProviderNameGa -> sweepersList = SweepersListGa
        ProviderNameBeta -> sweepersList = SweepersListBeta
        ProviderNameBetaDiffTest -> sweepersList = SweepersListBeta
        else -> throw Exception("Provider name not supplied when generating a nightly test subproject")
    }
    // We still allow locks in the service sweeper build configuration for adhoc triggers of services
    val serviceSweeperConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperName, sweepersList, projectId, vcsRoot, sharedResources, config)
    serviceSweeperConfig.triggers {
        finishBuildTrigger {
            buildType = "${DslContext.projectId}_${compositeId}"
            branchFilter = "+:${cron.branch}"
            successfulOnly = false
        }
    }

    // Add snapshot dependency on the composite config to run after tests finish
    serviceSweeperConfig.dependencies {
        snapshot(compositeConfig) {
            onDependencyFailure = FailureAction.IGNORE
            onDependencyCancel = FailureAction.IGNORE
        }
    }

    return Project {
        id(projectId)
        name = "Nightly Tests"
        description = "A project connected to the hashicorp/terraform-provider-${providerName} repository, where scheduled nightly tests run and users can trigger ad-hoc builds"

        buildType(compositeConfig)
        packageBuildConfigs.forEach { buildType(it) }
        buildType(serviceSweeperConfig)

        params{
            configureGoogleSpecificTestParameters(config)
        }
    }
}