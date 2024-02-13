/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package projects.reused

import SharedResourceNameVcr
import VcrRecordingProjectId
import builds.*
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun vcrRecording(parentProject:String, providerName: String, hashicorpVcsRoot: GitVcsRoot, modularMagicianVcsRoot: GitVcsRoot, config: AccTestConfiguration): Project {

    // Create unique ID for the dynamically-created project
    var projectId = "${parentProject}_${VcrRecordingProjectId}"
    projectId = replaceCharsId(projectId)

    val buildIdHashiCorp = replaceCharsId("${providerName}_HASHICORP_VCR")
    val buildIdModularMagician = replaceCharsId("${providerName}_MODMAGICIAN_VCR")

    // Shared resource allows VCR recording process to not clash with acceptance test or sweeper
    var sharedResources: List<String> = listOf(SharedResourceNameVcr)

    // Create the build config for hashicorp/terraform-provider-google
    var hcVcr = VcrDetails(providerName, buildIdHashiCorp, hashicorpVcsRoot, sharedResources)
    var hcBuildConfig = hcVcr.vcrBuildConfig(config)

    // Create the build config for modular-magician/terraform-provider-google
    var mmVcr = VcrDetails(providerName, buildIdModularMagician, modularMagicianVcsRoot, sharedResources)
    var mmBuildConfig = mmVcr.vcrBuildConfig(config)

    // Return VCR project with both build configs
    return Project {
        id(projectId)
        name = "VCR Recording"
        description = "A project connected to both modular-magician and hashicorp hosted terraform-provider-google-beta repositories, where users can trigger ad-hoc tests to re-record VCR cassettes"

        buildType(
            hcBuildConfig
        )
        buildType(
            mmBuildConfig
        )
    }
}