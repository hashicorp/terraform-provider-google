/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package projects.feature_branches

import ProviderNameBeta
import ProviderNameGa
import builds.*
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import projects.reused.nightlyTests
import replaceCharsId

const val branchName = "FEATURE-BRANCH-major-release-6.0.0"

// VCS Roots specifically for pulling code from the feature branches in the downstream repos

object HashicorpVCSRootGa_featureBranchMajorRelease600: GitVcsRoot({
    name = "VCS root for the hashicorp/terraform-provider-${ProviderNameGa} repo @ refs/heads/${branchName}"
    url = "https://github.com/hashicorp/terraform-provider-${ProviderNameGa}"
    branch = "refs/heads/${branchName}"
    branchSpec = """
        +:(refs/heads/*)
        -:refs/pulls/*
    """.trimIndent()
})

object HashicorpVCSRootBeta_featureBranchMajorRelease600: GitVcsRoot({
    name = "VCS root for the hashicorp/terraform-provider-${ProviderNameBeta} repo @ refs/heads/${branchName}"
    url = "https://github.com/hashicorp/terraform-provider-${ProviderNameBeta}"
    branch = "refs/heads/${branchName}"
    branchSpec = """
        +:(refs/heads/*)
        -:refs/pulls/*
    """.trimIndent()
})

fun featureBranchMajorRelease600_Project(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId(branchName)
    val gaProjectId = replaceCharsId(projectId + "_GA")
    val betaProjectId= replaceCharsId(projectId + "_BETA")

    // Get config for using the GA and Beta identities
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)

    return Project{
        id(projectId)
        name = "6.0.0 Major Release Testing"
        description = "Subproject for testing feature branch $branchName"

        // Register feature branch-specific VCS roots in the project
        vcsRoot(HashicorpVCSRootGa_featureBranchMajorRelease600)
        vcsRoot(HashicorpVCSRootBeta_featureBranchMajorRelease600)

        // Nested Nightly Test project that uses hashicorp/terraform-provider-google
        subProject(
            Project{
                id(gaProjectId)
                name = "Google"
                subProject(
                    nightlyTests(
                        gaProjectId,
                        ProviderNameGa,
                        HashicorpVCSRootGa_featureBranchMajorRelease600,
                        gaConfig,
                        NightlyTriggerConfiguration(
                            branch = "refs/heads/${branchName}", // Make triggered builds use the feature branch
                            daysOfWeek = "5"     // Thursday for GA, TeamCity numbers days Sun=1...Sat=7
                        ), 
                    )
                )
            }
        )

        // Nested Nightly Test project that uses hashicorp/terraform-provider-google-beta
        subProject(
            Project {
                id(betaProjectId)
                name = "Google Beta"
                subProject(
                    nightlyTests(
                        betaProjectId,
                        ProviderNameBeta,
                        HashicorpVCSRootBeta_featureBranchMajorRelease600,
                        betaConfig,
                        NightlyTriggerConfiguration(
                            branch = "refs/heads/${branchName}", // Make triggered builds use the feature branch
                            daysOfWeek="6"       // Friday for Beta, TeamCity numbers days Sun=1...Sat=7
                        ),
                    )
                )
            }
        )

        params {
            readOnlySettings()
        }
    }
}