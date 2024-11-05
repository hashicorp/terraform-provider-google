/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import jetbrains.buildServer.configs.kotlin.triggers.ScheduleTrigger
import org.junit.Assert
import org.junit.Test
import projects.feature_branches.featureBranchEphemeralResources
import projects.googleCloudRootProject

class FeatureBranchEphemeralResourcesSubProject {
    @Test
    fun buildsUsingHashiCorpReposAreOnSchedule() {
        val root = googleCloudRootProject(testContextParameters())

        // Find feature branch project
        val project = getSubProject(root, featureBranchEphemeralResources)

        // All builds using the HashiCorp owned GitHub repos
        val hashiBuilds = project.buildTypes.filter { bt ->
            bt.name.contains("HashiCorp downstream")
        }

        hashiBuilds.forEach{bt ->
            Assert.assertTrue(
                "Build configuration `${bt.name}` should contain at least one trigger",
                bt.triggers.items.isNotEmpty()
            )
            // Look for at least one CRON trigger
            var found = false
            lateinit var schedulingTrigger: ScheduleTrigger
            for (item in bt.triggers.items){
                if (item.type == "schedulingTrigger") {
                    schedulingTrigger = item as ScheduleTrigger
                    found = true
                    break
                }
            }

            Assert.assertTrue(
                "Build configuration `${bt.name}` should contain a CRON/'schedulingTrigger' trigger",
                found
            )

            // Check that triggered builds are being run on the feature branch
            val isCorrectBranch: Boolean = schedulingTrigger.branchFilter == "+:refs/heads/$featureBranchEphemeralResources"

            Assert.assertTrue(
                "Build configuration `${bt.name}` is using the $featureBranchEphemeralResources branch filter",
                isCorrectBranch
            )
        }
    }

    @Test
    fun buildsUsingModularMagicianReposAreNotTriggered() {
        val root = googleCloudRootProject(testContextParameters())

        // Find feature branch project
        val project = getSubProject(root, featureBranchEphemeralResources)

        // All builds using the HashiCorp owned GitHub repos
        val magicianBuilds = project.buildTypes.filter { bt ->
            bt.name.contains("MM upstream")
        }

        magicianBuilds.forEach{bt ->
            Assert.assertTrue(
                "Build configuration `${bt.name}` should not have any triggers",
                bt.triggers.items.isEmpty()
            )
        }
    }
}
