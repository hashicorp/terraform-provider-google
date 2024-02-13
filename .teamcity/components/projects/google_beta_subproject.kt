/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package projects

import ProviderNameBeta
import builds.AllContextParameters
import builds.getBetaAcceptanceTestConfig
import builds.getVcrAcceptanceTestConfig
import builds.readOnlySettings
import jetbrains.buildServer.configs.kotlin.Project
import projects.reused.mmUpstream
import projects.reused.nightlyTests
import projects.reused.vcrRecording
import replaceCharsId
import vcs_roots.HashiCorpVCSRootBeta
import vcs_roots.ModularMagicianVCSRootBeta

// googleSubProjectBeta returns a subproject that is used for testing terraform-provider-google-beta (Beta)
fun googleSubProjectBeta(allConfig: AllContextParameters): Project {

    var betaId = replaceCharsId("GOOGLE_BETA")

    // Get config for using the Beta and VCR identities
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    val vcrConfig = getVcrAcceptanceTestConfig(allConfig)

    return Project{
        id(betaId)
        name = "Google Beta"
        description = "Subproject containing builds for testing the Beta version of the Google provider"

        // Nightly Test project that uses hashicorp/terraform-provider-google-beta
        subProject(nightlyTests(betaId, ProviderNameBeta, HashiCorpVCSRootBeta, betaConfig))

        // MM Upstream project that uses modular-magician/terraform-provider-google-beta
        subProject(mmUpstream(betaId, ProviderNameBeta, ModularMagicianVCSRootBeta, vcrConfig))

        // VCR recording project that allows VCR recordings to be made using hashicorp/terraform-provider-google-beta OR modular-magician/terraform-provider-google-beta
        // This is only present for the Beta provider, as only TPGB VCR recordings are used.
        subProject(vcrRecording(betaId, ProviderNameBeta, HashiCorpVCSRootBeta, ModularMagicianVCSRootBeta, vcrConfig))

        params {
            readOnlySettings()
        }
    }
}