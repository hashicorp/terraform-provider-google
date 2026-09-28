/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import AllNightlyTestsName
import AllProvidersNightlyTestsName
import DefaultBranchName
import ServiceSweeperCronName
import ServiceSweeperManualName
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
import jetbrains.buildServer.configs.kotlin.ReuseBuilds
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.SharedResources
import jetbrains.buildServer.configs.kotlin.triggers.FinishBuildTrigger
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

        // The composite gate has no sweeper steps or environment parameters.
        listOf("Project Sweeper", "Folder Sweeper").forEach { name ->
            val bt = getBuildFromProject(globalSweepersProject, name)
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
    fun globalSweepersUseFinishTriggersWithoutDependencies() {
        listOf("TeamCityTests", "Experimental_NightlyTests").forEach { projectId ->
            val root = googleCloudRootProject(testContextParameters(projectId))
            val globalSweepers = getSubProject(root, globalSweepersProjectName)
            val gate = getBuildFromProject(globalSweepers, "Nightly Sweeper Gate")
            assertEquals(BuildTypeSettings.Type.COMPOSITE, gate.type)
            assertTrue("Gate should not execute sweeper steps", gate.steps.items.isEmpty())
            assertTrue("Gate should not hold locks needed by service sweepers", gate.features.items.filterIsInstance<SharedResources>().isEmpty())
            val betaNightly = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)
            val gaNightly = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
            val sweeperBeta = getBuildFromProject(betaNightly, ServiceSweeperName)
            val sweeperGa = getBuildFromProject(gaNightly, ServiceSweeperName)
            val allProvidersComposite = getBuildFromProject(root, AllProvidersNightlyTestsName)
            assertFinishTrigger(gate, allProvidersComposite, DefaultBranchName)
            assertSnapshotDependencies(gate, listOf(sweeperGa, sweeperBeta), FailureAction.IGNORE)
            gate.dependencies.items.forEach { dependency ->
                assertEquals("Gate service sweeper dependencies should always run a new build", ReuseBuilds.NO, dependency.snapshot!!.reuseBuilds)
            }
            listOf("Project Sweeper", "Folder Sweeper").forEach { name ->
                val sweeper = getBuildFromProject(globalSweepers, name)
                assertFinishTriggerWithoutBranchFilter(sweeper, gate)
                assertTrue("Global sweeper should not have snapshot dependencies", sweeper.dependencies.items.isEmpty())
                assertSharedResourceLocks(sweeper, SharedResources {
                    lockAllValues(SharedResourceNameGa)
                    lockAllValues(SharedResourceNameBeta)
                    lockAllValues(SharedResourceNameVcr)
                })
            }

        }
    }

    private fun assertFinishTriggerWithoutBranchFilter(build: BuildType, source: BuildType) {
        assertEquals("${build.name} should have exactly one trigger", 1, build.triggers.items.size)
        val trigger = build.triggers.items.single()
        assertTrue("${build.name} should use a finish-build trigger", trigger is FinishBuildTrigger)
        trigger as FinishBuildTrigger
        assertEquals("${build.name} should watch the source build's resolved ID", source.id!!.value, trigger.buildType)
        assertTrue("${build.name} should not filter the composite gate branch", trigger.branchFilter.isNullOrEmpty())
        assertEquals("${build.name} should not require a successful source build", false, trigger.successfulOnly)
    }
}
