/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.feature_branches

import ProviderNameBeta
import ProviderNameGa
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.ServicesListBeta
import generated.ServicesListGa
import jetbrains.buildServer.configs.kotlin.Project
import replaceCharsId
import vcs_roots.HashiCorpVCSRootBeta
import vcs_roots.HashiCorpVCSRootGa
import vcs_roots.ModularMagicianVCSRootBeta
import vcs_roots.ModularMagicianVCSRootGa

const val featureBranchEphemeralResources = "FEATURE-BRANCH-ephemeral-resource"
const val EphemeralResourcesTfCoreVersion = "1.10.0-alpha20240926" // TODO - update with correct release

// featureBranchEphemeralResourcesSubProject creates a project just for testing ephemeral resources.
// We know that all ephemeral resources we're adding are part of the Resource Manager service, so we only include those builds.
// We create builds for testing the resourcemanager service:
//    - Against the GA hashicorp repo
//    - Against the GA modular-magician repo
//    - Against the Beta hashicorp repo
//    - Against the Beta modular-magician repo
// These resemble existing projects present in TeamCity, but these all use a more recent version of Terraform including
// the new ephemeral values feature.
fun featureBranchEphemeralResourcesSubProject(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId(featureBranchEphemeralResources)

    val packageName = "resourcemanager" // All ephemeral resources will be in the resourcemanager package
    val vcrConfig = getVcrAcceptanceTestConfig(allConfig) // Reused below for both MM testing build configs
    val trigger  = NightlyTriggerConfiguration(
        branch = "refs/heads/$featureBranchEphemeralResources" // triggered builds must test the feature branch
    )


    // GA
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    // How to make only build configuration to the relevant package(s)
    val resourceManagerPackageGa = ServicesListGa.getValue(packageName)

    // Enable testing using hashicorp/terraform-provider-google
    var parentId = "${projectId}_HC_GA"
    val buildConfigHashiCorpGa = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageGa.getValue("path"), "Ephemeral resources in $packageName (GA provider, HashiCorp downstream)", ProviderNameGa, parentId, HashiCorpVCSRootGa, listOf(SharedResourceNameGa), gaConfig)
    buildConfigHashiCorpGa.addTrigger(trigger)

    // Enable testing using modular-magician/terraform-provider-google
    parentId = "${projectId}_MM_GA"
    val buildConfigModularMagicianGa = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageGa.getValue("path"), "Ephemeral resources in $packageName (GA provider, MM upstream)", ProviderNameGa, parentId, ModularMagicianVCSRootGa, listOf(SharedResourceNameVcr), vcrConfig)
    // No trigger added here (MM upstream is manual only)

    // Beta
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    val resourceManagerPackageBeta = ServicesListBeta.getValue(packageName)

    // Enable testing using hashicorp/terraform-provider-google-beta
    parentId = "${projectId}_HC_BETA"
    val buildConfigHashiCorpBeta = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageBeta.getValue("path"), "Ephemeral resources in $packageName (Beta provider, HashiCorp downstream)", ProviderNameBeta, parentId, HashiCorpVCSRootBeta, listOf(SharedResourceNameBeta), betaConfig)
    buildConfigHashiCorpBeta.addTrigger(trigger)

    // Enable testing using modular-magician/terraform-provider-google-beta
    parentId = "${projectId}_MM_BETA"
    val buildConfigModularMagicianBeta = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageBeta.getValue("path"), "Ephemeral resources in $packageName (Beta provider, MM upstream)", ProviderNameBeta, parentId, ModularMagicianVCSRootBeta, listOf(SharedResourceNameVcr), vcrConfig)
    // No trigger added here (MM upstream is manual only)


    // ------

    // Make all builds use a 1.10.0-ish version of TF core
    val allBuildConfigs = listOf(buildConfigHashiCorpGa, buildConfigModularMagicianGa, buildConfigHashiCorpBeta, buildConfigModularMagicianBeta)
    allBuildConfigs.forEach{ b ->
        b.overrideTerraformCoreVersion(EphemeralResourcesTfCoreVersion)
    }

    // ------

    return Project{
        id(projectId)
        name = featureBranchEphemeralResources
        description = "Subproject for testing feature branch $featureBranchEphemeralResources"

        // Register all build configs in the project
        allBuildConfigs.forEach{ b ->
            buildType(b)
        }

        params {
            readOnlySettings()
        }
    }
}