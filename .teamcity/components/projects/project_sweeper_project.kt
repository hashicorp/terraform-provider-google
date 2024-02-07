// This file is controlled by MMv1, any changes made here will be overwritten

package projects

import ProjectSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.SweepersList
import jetbrains.buildServer.configs.kotlin.Project
import replaceCharsId
import vcs_roots.HashiCorpVCSRootGa

// projectSweeperSubProject returns a subproject that contains a sweeper for project resources
// Sweeping projects is an edge case because it doesn't respect boundaries between different testing projects GA/Beta/PR
fun projectSweeperSubProject(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId("PROJECT_SWEEPER")

    // Get config for using the GA identity (arbitrary choice as sweeper isn't confined by GA/Beta etc)
    val gaConfig = getGaAcceptanceTestConfig(allConfig)

    // List of ALL shared resources; avoid clashing with any other running build
    val sharedResources: List<String> = listOf(SharedResourceNameGa, SharedResourceNameBeta, SharedResourceNameVcr)

    // Create build config for sweeping project resources
    // Uses the HashiCorpVCSRootGa VCS Root so that the latest sweepers in hashicorp/terraform-provider-google are used
    val serviceSweeperConfig = BuildConfigurationForSweeper("N/A", ProjectSweeperName, SweepersList, projectId, HashiCorpVCSRootGa, sharedResources, gaConfig)
    serviceSweeperConfig.enableProjectSweep()
    val trigger  = NightlyTriggerConfiguration()
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