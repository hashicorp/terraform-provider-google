/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import org.junit.Assert.assertTrue
import org.junit.Test
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert
import projects.googleCloudRootProject

class NightlyTestProjectsTests {
    @Test
    fun allBuildsShouldHaveTrigger() {
        val project = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        var gaNightlyTestProject = getSubProject(project, gaProjectName, nightlyTestsProjectName)

        // Find Beta nightly test project
        var betaNightlyTestProject = getSubProject(project, betaProjectName, nightlyTestsProjectName)

        // Make assertions about builds in both nightly test projects
        (gaNightlyTestProject.buildTypes + betaNightlyTestProject.buildTypes).forEach{bt ->
            assertTrue("Build configuration `${bt.name}` contains at least one trigger", bt.triggers.items.isNotEmpty())
             // Look for at least one CRON trigger
            var found: Boolean = false
            for (item in bt.triggers.items){
                if (item.type == "schedulingTrigger") {
                    found = true
                    break
                }
            }
            assertTrue("Build configuration `${bt.name}` contains a CRON trigger", found)
        }
    }
}
