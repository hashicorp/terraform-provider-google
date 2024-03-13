/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.feature_branches

import ProviderNameBeta
import ProviderNameGa
import builds.*
import generated.PackagesListBeta
import generated.PackagesListGa
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId
import vcs_roots.ModularMagicianVCSRootBeta
import vcs_roots.ModularMagicianVCSRootGa

const val featureBranchProviderFunctionsName = "FEATURE-BRANCH-provider-functions"
const val providerFunctionsTfCoreVersion = "1.8.0-alpha20240228"

// VCS Roots specifically for pulling code from the feature branches in the downstream and upstream repos
object HashicorpVCSRootGa_featureBranchProviderFunctions: GitVcsRoot({
    name = "VCS root for the hashicorp/terraform-provider-${ProviderNameGa} repo @ refs/heads/${featureBranchProviderFunctionsName}"
    url = "https://github.com/hashicorp/terraform-provider-${ProviderNameGa}"
    branch = "refs/heads/${featureBranchProviderFunctionsName}"
    branchSpec = "" // empty as we'll access no other branches
})

object HashicorpVCSRootBeta_featureBranchProviderFunctions: GitVcsRoot({
    name = "VCS root for the hashicorp/terraform-provider-${ProviderNameBeta} repo @ refs/heads/${featureBranchProviderFunctionsName}"
    url = "https://github.com/hashicorp/terraform-provider-${ProviderNameBeta}"
    branch = "refs/heads/${featureBranchProviderFunctionsName}"
    branchSpec = "" // empty as we'll access no other branches
})

fun featureBranchProviderFunctionSubProject(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId(featureBranchProviderFunctionsName)

    val packageName = "functions" // This project will contain only builds to test this single package
    val sharedResourcesEmpty: List<String> = listOf() // No locking when testing functions
    val vcrConfig = getVcrAcceptanceTestConfig(allConfig) // Reused below for both MM testing build configs
    val trigger  = NightlyTriggerConfiguration() // Resued below for running tests against the downstream repos every night.

    var parentId: String // To be overwritten when each build config is generated below.

    // GA
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    // How to make only build configuration to the relevant package(s)
    val functionPackageGa = PackagesListGa.getValue(packageName)

    // Enable testing using hashicorp/terraform-provider-google
    parentId = "${projectId}_HC_GA"
    val buildConfigHashiCorpGa = BuildConfigurationForSinglePackage(packageName, functionPackageGa.getValue("path"), "Provider-Defined Functions (GA provider, HashiCorp downstream)", ProviderNameGa, parentId, HashicorpVCSRootGa_featureBranchProviderFunctions, sharedResourcesEmpty, gaConfig)
    buildConfigHashiCorpGa.addTrigger(trigger)

    // Enable testing using modular-magician/terraform-provider-google
    parentId = "${projectId}_MM_GA"
    val buildConfigModularMagicianGa = BuildConfigurationForSinglePackage(packageName, functionPackageGa.getValue("path"), "Provider-Defined Functions (GA provider, MM upstream)", ProviderNameGa, parentId, ModularMagicianVCSRootGa, sharedResourcesEmpty, vcrConfig)

    // Beta
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    val functionPackageBeta = PackagesListBeta.getValue("functions")

    // Enable testing using hashicorp/terraform-provider-google-beta
    parentId = "${projectId}_HC_BETA"
    val buildConfigHashiCorpBeta = BuildConfigurationForSinglePackage(packageName, functionPackageBeta.getValue("path"), "Provider-Defined Functions (Beta provider, HashiCorp downstream)", ProviderNameBeta, parentId, HashicorpVCSRootBeta_featureBranchProviderFunctions, sharedResourcesEmpty, betaConfig)
    buildConfigHashiCorpBeta.addTrigger(trigger)

    // Enable testing using modular-magician/terraform-provider-google-beta
    parentId = "${projectId}_MM_BETA"
    val buildConfigModularMagicianBeta = BuildConfigurationForSinglePackage(packageName, functionPackageBeta.getValue("path"), "Provider-Defined Functions (Beta provider, MM upstream)", ProviderNameBeta, parentId, ModularMagicianVCSRootBeta, sharedResourcesEmpty, vcrConfig)

    val allBuildConfigs = listOf(buildConfigHashiCorpGa, buildConfigModularMagicianGa, buildConfigHashiCorpBeta, buildConfigModularMagicianBeta)

    // Make these builds use a 1.8.0-ish version of TF core
    allBuildConfigs.forEach{ b ->
        b.overrideTerraformCoreVersion(providerFunctionsTfCoreVersion)
    }

    return Project{
        id(projectId)
        name = featureBranchProviderFunctionsName
        description = "Subproject for testing feature branch $featureBranchProviderFunctionsName"

        // Register feature branch-specific VCS roots in the project
        vcsRoot(HashicorpVCSRootGa_featureBranchProviderFunctions)
        vcsRoot(HashicorpVCSRootBeta_featureBranchProviderFunctions)

        // Register all build configs in the project
        allBuildConfigs.forEach{ b ->
            buildType(b)
        }

        params {
            readOnlySettings()
        }
    }
}