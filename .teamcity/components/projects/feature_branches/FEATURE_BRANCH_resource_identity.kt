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
import components.projects.feature_branches.getServicesList
import DefaultStartHour

const val featureBranchResourceIdentity = "FEATURE-BRANCH-resource-identity"
const val ResourceIdentityTfCoreVersion = "1.12.0-beta2"

fun featureBranchResourceIdentitySubProject(allConfig: AllContextParameters): Project {

    val trigger  = NightlyTriggerConfiguration(
        branch = "refs/heads/$featureBranchResourceIdentity", // triggered builds must test the feature branch
        startHour = DefaultStartHour + 6,
        nightlyTestsEnabled = false
    )
    val vcrConfig = getVcrAcceptanceTestConfig(allConfig) // Reused below for both MM testing build configs
    val servicesToTest = arrayOf("secretmanager", "resourcemanager")

    // GA
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    // These are the packages that have resources that will use write-only attributes
    var ServicesListWriteOnlyGA = getServicesList(servicesToTest, "GA")

    val buildConfigsGa = BuildConfigurationsForPackages(ServicesListWriteOnlyGA, ProviderNameGa, "ResourceIdentityGa - HC", HashiCorpVCSRootGa, listOf(SharedResourceNameGa), gaConfig, "TestAcc.*ResourceIdentity")
    buildConfigsGa.forEach{ builds ->
        builds.addTrigger(trigger)
    }

    var ServicesListWriteOnlyGaMM = getServicesList(servicesToTest, "GA-MM")
    val buildConfigsMMGa = BuildConfigurationsForPackages(ServicesListWriteOnlyGaMM, ProviderNameGa, "ResourceIdentityGa - MM", ModularMagicianVCSRootGa, listOf(SharedResourceNameGa), vcrConfig, "TestAcc.*ResourceIdentity")

    // Beta
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    var ServicesListWriteOnlyBeta = getServicesList(servicesToTest, "Beta")
    val buildConfigsBeta = BuildConfigurationsForPackages(ServicesListWriteOnlyBeta, ProviderNameBeta, "ResourceIdentityBeta - HC", HashiCorpVCSRootBeta, listOf(SharedResourceNameBeta), betaConfig, "TestAcc.*ResourceIdentity")
    buildConfigsBeta.forEach{ builds ->
        builds.addTrigger(trigger)
    }

    var ServicesListWriteOnlyBetaMM = getServicesList(servicesToTest, "Beta-MM")
    val buildConfigsMMBeta = BuildConfigurationsForPackages(ServicesListWriteOnlyBetaMM, ProviderNameBeta, "ResourceIdentityBeta - MM", ModularMagicianVCSRootBeta, listOf(SharedResourceNameBeta), vcrConfig, "TestAcc.*ResourceIdentity")

    // Make all builds use a 1.12.0-ish version of TF core
    val allBuildConfigs = buildConfigsGa + buildConfigsBeta + buildConfigsMMGa + buildConfigsMMBeta
    allBuildConfigs.forEach{ builds ->
        builds.overrideTerraformCoreVersion(ResourceIdentityTfCoreVersion)
    }

    // ------

    return Project{
        id("FEATURE_BRANCH_resource_identity")
        name = featureBranchResourceIdentity
        description = "Subproject for testing feature branch $featureBranchResourceIdentity"

        // Register all build configs in the project
        allBuildConfigs.forEach{ builds ->
            buildType(builds)
        }

        params {
            readOnlySettings()
        }
    }
}