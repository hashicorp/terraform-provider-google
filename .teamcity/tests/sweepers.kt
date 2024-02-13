/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package tests

import ServiceSweeperName
import jetbrains.buildServer.configs.kotlin.BuildType
import org.junit.Assert.assertTrue
import org.junit.Test
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert
import projects.googleCloudRootProject

class SweeperTests {
    @Test
    fun projectSweeperProjectDoesNotSkipProjectSweep() {
        val project = googleCloudRootProject(testContextParameters())

        // Find Project sweeper project
        val projectSweeperProject: Project? =  project.subProjects.find { p->  p.name == projectSweeperProjectName}
        if (projectSweeperProject == null) {
            Assert.fail("Could not find the Project Sweeper project")
        }

        // For the project sweeper to be skipped, SKIP_PROJECT_SWEEPER needs a value
        // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26

        projectSweeperProject!!.buildTypes.forEach{bt ->
            val value = bt.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER is set to an empty value, so project sweepers are NOT skipped. Value = `${value}` ", value == "")
        }
    }

    @Test
    fun gaNightlyProjectServiceSweeperSkipsProjectSweep() {
        val project = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        val gaProject: Project? =  project.subProjects.find { p->  p.name == gaProjectName}
        if (gaProject == null) {
            Assert.fail("Could not find the Google (GA) project")
        }
        val gaNightlyTestProject: Project? = gaProject!!.subProjects.find { p->  p.name == nightlyTestsProjectName}
        if (gaNightlyTestProject == null) {
            Assert.fail("Could not find the Google (GA) Nightly Test project")
        }

        // Find sweeper inside
        val sweeper: BuildType? = gaNightlyTestProject!!.buildTypes.find { p-> p.name == ServiceSweeperName}
        if (sweeper == null) {
            Assert.fail("Could not find the sweeper build in the Google (GA) Nightly Test project")
        }

        // For the project sweeper to be skipped, SKIP_PROJECT_SWEEPER needs a value
        // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26

        val value = sweeper!!.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER is set to a non-empty string, so project sweepers are skipped. Value = `${value}` ", value != "")
    }

    @Test
    fun betaNightlyProjectServiceSweeperSkipsProjectSweep() {
        val project = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val betaProject: Project? =  project.subProjects.find { p->  p.name == betaProjectName}
        if (betaProject == null) {
            Assert.fail("Could not find the Google (GA) project")
        }
        val betaNightlyTestProject: Project? = betaProject!!.subProjects.find { p->  p.name == nightlyTestsProjectName}
        if (betaNightlyTestProject == null) {
            Assert.fail("Could not find the Google (GA) Nightly Test project")
        }

        // Find sweeper inside
        val sweeper: BuildType? = betaNightlyTestProject!!.buildTypes.find { p-> p.name == ServiceSweeperName}
        if (sweeper == null) {
            Assert.fail("Could not find the sweeper build in the Google (GA) Nightly Test project")
        }

        // For the project sweeper to be skipped, SKIP_PROJECT_SWEEPER needs a value
        // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26

        val value = sweeper!!.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER is set to a non-empty string, so project sweepers are skipped. Value = `${value}` ", value != "")
    }
}
