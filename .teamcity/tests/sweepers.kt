/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import AllNightlyTestsName
import DefaultBranchName
import ServiceSweeperCronName
import ServiceSweeperManualName
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.SharedResources
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
    fun globalSweepersDependOnAllNightlyTests() {
        listOf("TeamCityTests", "Experimental_NightlyTests").forEach { projectId ->
            val root = googleCloudRootProject(testContextParameters(projectId))
            val gaNightly = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
            val betaNightly = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)
            val gaComposite = getBuildFromProject(gaNightly, AllNightlyTestsName)
            val betaComposite = getBuildFromProject(betaNightly, AllNightlyTestsName)
            val sweeperGa = getBuildFromProject(gaNightly, ServiceSweeperName)
            val sweeperBeta = getBuildFromProject(betaNightly, ServiceSweeperName)
            assertFinishTrigger(sweeperGa, gaComposite, DefaultBranchName)
            assertFinishTrigger(sweeperBeta, betaComposite, DefaultBranchName)

            val globalSweepers = getSubProject(root, globalSweepersProjectName)
            listOf("Project Sweeper", "Folder Sweeper").forEach { name ->
                val sweeper = getBuildFromProject(globalSweepers, name)
                assertFinishTrigger(sweeper, sweeperGa, DefaultBranchName)
                assertSnapshotDependencies(sweeper, listOf(gaComposite, betaComposite, sweeperGa, sweeperBeta), FailureAction.IGNORE)
                assertSharedResourceLocks(sweeper, SharedResources {
                    lockAllValues(SharedResourceNameGa)
                    lockAllValues(SharedResourceNameBeta)
                    lockAllValues(SharedResourceNameVcr)
                })
            }
        }
    }
}
