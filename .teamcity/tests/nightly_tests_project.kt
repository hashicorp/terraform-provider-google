/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

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
        var gaProject: Project? =  project.subProjects.find { p->  p.name == gaProjectName}
        if (gaProject == null) {
            Assert.fail("Could not find the Google (GA) project")
        }
        var gaNightlyTestProject: Project? = gaProject!!.subProjects.find { p->  p.name == nightlyTestsProjectName}
        if (gaNightlyTestProject == null) {
            Assert.fail("Could not find the Google (GA) Nightly Test project")
        }

        // Find Beta nightly test project
        var betaProject: Project? =  project.subProjects.find { p->  p.name == betaProjectName}
        if (betaProject == null) {
            Assert.fail("Could not find the Google (Beta) project")
        }
        var betaNightlyTestProject: Project? = betaProject!!.subProjects.find { p->  p.name == nightlyTestsProjectName}
        if (betaNightlyTestProject == null) {
            Assert.fail("Could not find the Google (GA) Nightly Test project")
        }

        (gaNightlyTestProject!!.buildTypes + betaNightlyTestProject!!.buildTypes).forEach{bt ->
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
