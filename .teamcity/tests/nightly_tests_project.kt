/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import jetbrains.buildServer.configs.kotlin.triggers.ScheduleTrigger
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject

class NightlyTestProjectsTests {
    @Test
    fun allBuildsShouldHaveTrigger() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        var gaNightlyTestProject = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)

        // Find Beta nightly test project
        var betaNightlyTestProject = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        // Make assertions about builds in both nightly test projects
        (gaNightlyTestProject.buildTypes + betaNightlyTestProject.buildTypes).forEach{bt ->
            assertTrue("Build configuration `${bt.name}` should contain at least one trigger", bt.triggers.items.isNotEmpty())
             // Look for at least one CRON trigger
            var found: Boolean = false
            lateinit var schedulingTrigger: ScheduleTrigger
            for (item in bt.triggers.items){
                if (item.type == "schedulingTrigger") {
                    schedulingTrigger = item as ScheduleTrigger
                    found = true
                    break
                }
            }

            assertTrue("Build configuration `${bt.name}` should contain a CRON/'schedulingTrigger' trigger", found)

            // Check that nightly test is being ran on the nightly-test branch
            var isNightlyTestBranch: Boolean = false
            if (schedulingTrigger.branchFilter == "+:refs/heads/nightly-test"){
                isNightlyTestBranch = true
            }
            assertTrue("Build configuration `${bt.name}` is using the nightly-test branch filter;", isNightlyTestBranch)
        }
    }
}
