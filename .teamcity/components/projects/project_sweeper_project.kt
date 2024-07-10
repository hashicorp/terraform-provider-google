/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects

import ProjectSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.SweepersListGa
import jetbrains.buildServer.configs.kotlin.Project
import replaceCharsId
import vcs_roots.HashiCorpVCSRootGa

// projectSweeperSubProject returns a subproject that contains a sweeper for project resources
// Sweeping projects is an edge case because it doesn't respect boundaries between different testing projects GA/Beta/PR
fun projectSweeperSubProject(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId("PROJECT_SWEEPER")

    // Get config for using the GA identity (arbitrary choice as sweeper isn't confined by GA/Beta etc.)
    val gaConfig = getGaAcceptanceTestConfig(allConfig)

    // List of ALL shared resources; avoid clashing with any other running build
    val sharedResources: List<String> = listOf(SharedResourceNameGa, SharedResourceNameBeta, SharedResourceNameVcr)

    // Create build config for sweeping project resources
    // Uses the HashiCorpVCSRootGa VCS Root so that the latest sweepers in hashicorp/terraform-provider-google are used
    val serviceSweeperConfig = BuildConfigurationForProjectSweeper("N/A", ProjectSweeperName, SweepersListGa, projectId, HashiCorpVCSRootGa, sharedResources, gaConfig)
    val trigger  = NightlyTriggerConfiguration(startHour=12)
    serviceSweeperConfig.addTrigger(trigger)

    return Project{
        id(projectId)
        name = "Project Sweeper"
        description = "Subproject containing a build configuration for sweeping project resources"

        // Register build configs in the project
        buildType(serviceSweeperConfig)

        params {
            readOnlySettings()
        }
    }
}