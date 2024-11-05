/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import ProjectSweeperName
import ServiceSweeperCronName
import ServiceSweeperManualName
import ServiceSweeperName
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.triggers.ScheduleTrigger
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject

class SweeperTests {
    @Test
    fun projectSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Project sweeper project
        val projectSweeperProject = getSubProject(root, projectSweeperProjectName)

        // SKIP_PROJECT_SWEEPER should be empty so project sweepers will be run
        // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26

        projectSweeperProject.buildTypes.forEach{bt ->
            val skipProjectSweeper = bt.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to an empty value, so project sweepers are NOT skipped in the ${projectSweeperProject.name} project. Value = `${skipProjectSweeper}` ", skipProjectSweeper == "")
        }
    }

    @Test
    fun gaNightlyTestsServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        val project = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)

        // Find sweeper inside
        val sweeper = getBuildFromProject(project, ServiceSweeperName)

        // Check PACKAGE_PATH is in google (not google-beta)
        val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
        assertEquals("./google/sweeper", value)

        // SKIP_PROJECT_SWEEPER should have a value so project sweepers will be skipped
        // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26
        val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string so project sweepers are skipped in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")
    }

    @Test
    fun betaNightlyTestsServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val project = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        // Find sweeper inside
        val sweeper: BuildType = getBuildFromProject(project, ServiceSweeperName)

        // Check PACKAGE_PATH is in google-beta
        val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
        assertEquals("./google-beta/sweeper", value)

        // SKIP_PROJECT_SWEEPER should have a value so project sweepers will be skipped
        // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26
        val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string so project sweepers are skipped in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")
    }

    @Test
    fun gaMmUpstreamServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val project = getNestedProjectFromRoot(root, gaProjectName, mmUpstreamProjectName)

        // Find sweepers inside
        val cronSweeper = getBuildFromProject(project, ServiceSweeperCronName)
        val manualSweeper = getBuildFromProject(project, ServiceSweeperManualName)
        val allSweepers: ArrayList<BuildType> = arrayListOf(cronSweeper, manualSweeper)
        allSweepers.forEach{ sweeper ->
            // Check PACKAGE_PATH is in google-beta
            val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
            assertEquals("./google/sweeper", value)

            // SKIP_PROJECT_SWEEPER should have a value so project sweepers will be skipped
            // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26
            val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string so project sweepers are skipped in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")
        }
    }

    @Test
    fun betaMmUpstreamServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val project = getNestedProjectFromRoot(root, betaProjectName, mmUpstreamProjectName)

        // Find sweepers inside
        val cronSweeper = getBuildFromProject(project, ServiceSweeperCronName)
        val manualSweeper = getBuildFromProject(project, ServiceSweeperManualName)
        val allSweepers: ArrayList<BuildType> = arrayListOf(cronSweeper, manualSweeper)
        allSweepers.forEach{ sweeper ->
            // Check PACKAGE_PATH is in google-beta
            val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
            assertEquals("./google-beta/sweeper", value)

            // SKIP_PROJECT_SWEEPER should have a value so project sweepers will be skipped
            // See https://github.com/GoogleCloudPlatform/magic-modules/blob/501429790939717ca6dce76dbf4b1b82aef4e9d9/mmv1/third_party/terraform/services/resourcemanager/resource_google_project_sweeper.go#L18-L26
            val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string so project sweepers are skipped in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")
        }
    }

    @Test
    fun projectSweepersRunAfterServiceSweepers() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project's service sweeper
        val gaNightlyTests: Project = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
        val sweeperGa: BuildType = getBuildFromProject(gaNightlyTests, ServiceSweeperName)

        // Find Beta nightly test project's service sweeper
        val betaNightlyTests : Project = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)
        val sweeperBeta: BuildType = getBuildFromProject(betaNightlyTests, ServiceSweeperName)

        // Find Project sweeper project's build
        val projectSweeperProject = getSubProject(root, projectSweeperProjectName)
        val projectSweeper: BuildType = getBuildFromProject(projectSweeperProject, ProjectSweeperName)
        
        // Check only one schedule trigger is on the builds in question
        assertTrue(sweeperGa.triggers.items.size == 1)
        assertTrue(sweeperBeta.triggers.items.size == 1)
        assertTrue(projectSweeper.triggers.items.size == 1)

        // Assert that the hour value that sweeper builds are triggered at is less than the hour value that project sweeper builds are triggered at
        // i.e. sweeper builds are triggered first
        val stGa = sweeperGa.triggers.items[0] as ScheduleTrigger
        val cronGa = stGa.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        val stBeta = sweeperBeta.triggers.items[0] as ScheduleTrigger
        val cronBeta = stBeta.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        val stProject = projectSweeper.triggers.items[0] as ScheduleTrigger
        val cronProject = stProject.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        assertTrue("Service sweeper for the GA Nightly Test project should be triggered at an earlier hour than the project sweeper", cronGa.hours.toString().toInt() < cronProject.hours.toString().toInt()) // Converting nullable strings to ints
        assertTrue("Service sweeper for the Beta Nightly Test project should be triggered at an earlier hour than the project sweeper", cronBeta.hours.toString().toInt() < cronProject.hours.toString().toInt() )
    }
}
