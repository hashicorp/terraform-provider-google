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

const val featureBranchEphemeralWriteOnly = "FEATURE-BRANCH-ephemeral-write-only"
const val EphemeralWriteOnlyTfCoreVersion = "1.11.0-alpha20241211"

fun featureBranchEphemeralWriteOnlySubProject(allConfig: AllContextParameters): Project {

    val trigger  = NightlyTriggerConfiguration(
        branch = "refs/heads/$featureBranchEphemeralWriteOnly", // triggered builds must test the feature branch
        startHour = DefaultStartHour + 6,
    )
    val vcrConfig = getVcrAcceptanceTestConfig(allConfig) // Reused below for both MM testing build configs

    // GA
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    // These are the packages that have resources that will use write-only attributes
    var ServicesListWriteOnlyGA = getServicesList(arrayOf("compute", "secretmanager", "sql", "bigquerydatatransfer"), "GA")

    val buildConfigsGa = BuildConfigurationsForPackages(ServicesListWriteOnlyGA, ProviderNameGa, "EphemeralWriteOnlyGa - HC", HashiCorpVCSRootGa, listOf(SharedResourceNameGa), gaConfig, "TestAcc.*Ephemeral")
    buildConfigsGa.forEach{ builds ->
        builds.addTrigger(trigger)
    }

    var ServicesListWriteOnlyGaMM = getServicesList(arrayOf("compute", "secretmanager", "sql", "bigquerydatatransfer"), "GA-MM")
    val buildConfigsMMGa = BuildConfigurationsForPackages(ServicesListWriteOnlyGaMM, ProviderNameGa, "EphemeralWriteOnlyGa - MM", ModularMagicianVCSRootGa, listOf(SharedResourceNameGa), vcrConfig, "TestAcc.*Ephemeral")

    // Beta
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    var ServicesListWriteOnlyBeta = getServicesList(arrayOf("compute", "secretmanager", "sql", "bigquerydatatransfer"), "Beta")
    val buildConfigsBeta = BuildConfigurationsForPackages(ServicesListWriteOnlyBeta, ProviderNameBeta, "EphemeralWriteOnlyBeta - HC", HashiCorpVCSRootBeta, listOf(SharedResourceNameBeta), betaConfig, "TestAcc.*Ephemeral")
    buildConfigsBeta.forEach{ builds ->
        builds.addTrigger(trigger)
    }

    var ServicesListWriteOnlyBetaMM = getServicesList(arrayOf("compute", "secretmanager", "sql", "bigquerydatatransfer"), "Beta-MM")
    val buildConfigsMMBeta = BuildConfigurationsForPackages(ServicesListWriteOnlyBetaMM, ProviderNameBeta, "EphemeralWriteOnlyBeta - MM", ModularMagicianVCSRootBeta, listOf(SharedResourceNameBeta), vcrConfig, "TestAcc.*Ephemeral")

    // Make all builds use a 1.11.0-ish version of TF core
    val allBuildConfigs = buildConfigsGa + buildConfigsBeta + buildConfigsMMGa + buildConfigsMMBeta
    allBuildConfigs.forEach{ builds ->
        builds.overrideTerraformCoreVersion(EphemeralWriteOnlyTfCoreVersion)
    }

    // ------

    return Project{
        id("FEATURE_BRANCH_ephemeral_write_only")
        name = featureBranchEphemeralWriteOnly
        description = "Subproject for testing feature branch $featureBranchEphemeralWriteOnly"

        // Register all build configs in the project
        allBuildConfigs.forEach{ builds ->
            buildType(builds)
        }

        params {
            readOnlySettings()
        }
    }
}