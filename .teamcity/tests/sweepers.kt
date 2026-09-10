/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import GlobalSweepersProjectName
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
    fun globalSweepersConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Global sweepers project
        val globalSweepersProject = getSubProject(root, globalSweepersProjectName)

        globalSweepersProject.buildTypes.forEach{bt ->
            val skipProjectSweeper = bt.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to an empty value in the ${globalSweepersProject.name} project. Value = `${skipProjectSweeper}` ", skipProjectSweeper == "")

            val skipFolderSweeper = bt.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
            assertTrue("env.SKIP_FOLDER_SWEEPER should be set to an empty value in the ${globalSweepersProject.name} project. Value = `${skipFolderSweeper}` ", skipFolderSweeper == "")
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

        // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
        val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

        val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
        assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
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

        // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
        val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

        val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
        assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
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

            // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
            val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

            val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
            assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
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

            // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
            val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

            val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
            assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
        }
    }

    @Test
    fun globalSweepersRunAfterServiceSweepers() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project's service sweeper
        val gaNightlyTests: Project = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
        val sweeperGa: BuildType = getBuildFromProject(gaNightlyTests, ServiceSweeperName)

        // Find Beta nightly test project's service sweeper
        val betaNightlyTests : Project = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)
        val sweeperBeta: BuildType = getBuildFromProject(betaNightlyTests, ServiceSweeperName)

        // Find Global sweepers project's builds
        val globalSweepersProject = getSubProject(root, globalSweepersProjectName)
        val projectSweeper: BuildType = getBuildFromProject(globalSweepersProject, "Project Sweeper")
        val folderSweeper: BuildType = getBuildFromProject(globalSweepersProject, "Folder Sweeper")
        
        // Check only one schedule trigger is on the builds in question
        assertTrue(sweeperGa.triggers.items.size == 1)
        assertTrue(sweeperBeta.triggers.items.size == 1)
        assertTrue(projectSweeper.triggers.items.size == 1)
        assertTrue(folderSweeper.triggers.items.size == 1)

        // Assert that the hour value that sweeper builds are triggered at is less than the hour value that project/folder sweeper builds are triggered at
        val stGa = sweeperGa.triggers.items[0] as ScheduleTrigger
        val cronGa = stGa.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        val stBeta = sweeperBeta.triggers.items[0] as ScheduleTrigger
        val cronBeta = stBeta.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        
        val stProject = projectSweeper.triggers.items[0] as ScheduleTrigger
        val cronProject = stProject.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        
        val stFolder = folderSweeper.triggers.items[0] as ScheduleTrigger
        val cronFolder = stFolder.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron
        
        assertTrue("Service sweeper for the GA Nightly Test project should be triggered at an earlier hour than the project sweeper", cronGa.hours.toString().toInt() < cronProject.hours.toString().toInt())
        assertTrue("Service sweeper for the Beta Nightly Test project should be triggered at an earlier hour than the project sweeper", cronBeta.hours.toString().toInt() < cronProject.hours.toString().toInt())
        
        assertTrue("Service sweeper for the GA Nightly Test project should be triggered at an earlier hour than the folder sweeper", cronGa.hours.toString().toInt() < cronFolder.hours.toString().toInt())
        assertTrue("Service sweeper for the Beta Nightly Test project should be triggered at an earlier hour than the folder sweeper", cronBeta.hours.toString().toInt() < cronFolder.hours.toString().toInt())
    }
}
