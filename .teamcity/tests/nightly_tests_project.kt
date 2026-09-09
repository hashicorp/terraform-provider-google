/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import AllNightlyTestsName
import DefaultBranchName
import ProviderNameBeta
import ProviderNameGa
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import builds.NightlyTriggerConfiguration
import builds.getGaAcceptanceTestConfig
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.SharedResources
import jetbrains.buildServer.configs.kotlin.triggers.ScheduleTrigger
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject
import projects.reused.getAllPackageInProviderVersion
import projects.reused.nightlyTests
import vcs_roots.HashiCorpVCSRootGa

class NightlyTestProjectsTests {
    @Test
    fun onlyCompositeShouldHaveNightlySchedule() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        var gaNightlyTestProject = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)

        // Find Beta nightly test project
        var betaNightlyTestProject = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        listOf(gaNightlyTestProject, betaNightlyTestProject).forEach { project ->
            val composite = getBuildFromProject(project, AllNightlyTestsName)
            assertEquals("Composite should have one schedule trigger", 1, composite.triggers.items.size)
            val trigger = composite.triggers.items.single()
            assertTrue("Composite should use a schedule trigger", trigger is ScheduleTrigger)
            trigger as ScheduleTrigger
            assertTrue("Composite should use CRON scheduling", trigger.schedulingPolicy is ScheduleTrigger.SchedulingPolicy.Cron)
            assertEquals("Composite should select the nightly branch", "+:$DefaultBranchName", trigger.branchFilter)
            assertEquals("Nightly runs should not require pending changes", false, trigger.withPendingChangesOnly)

            val sweeper = getBuildFromProject(project, ServiceSweeperName)
            assertFinishTrigger(sweeper, composite, DefaultBranchName)
            project.buildTypes.filter { it != composite && it != sweeper }.forEach { build ->
                assertTrue("Package build `${build.name}` should have no independent trigger", build.triggers.items.isEmpty())
            }
        }
    }

    @Test
    fun nightlyTestsShouldHaveCompositeAllTestsBuild() {
        val root = googleCloudRootProject(testContextParameters())

        var gaNightlyTestProject = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
        var betaNightlyTestProject = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        listOf(gaNightlyTestProject, betaNightlyTestProject).forEach { project ->
            val composite = getBuildFromProject(project, AllNightlyTestsName)
            assertEquals("Build configuration `${composite.name}` should be a COMPOSITE build", BuildTypeSettings.Type.COMPOSITE, composite.type)

            val packageBuilds = project.buildTypes.filter { bt ->
                bt.name != ServiceSweeperName && bt.name != AllNightlyTestsName
            }
            assertTrue("Nightly test project `${project.name}` should have package test builds", packageBuilds.isNotEmpty())
            assertSnapshotDependencies(composite, packageBuilds, FailureAction.ADD_PROBLEM)

            val sweeper = getBuildFromProject(project, ServiceSweeperName)
            assertSnapshotDependencies(sweeper, listOf(composite), FailureAction.IGNORE)
        }
    }

    @Test
    fun serviceSweeperShouldFollowCustomNightlyBranch() {
        val config = getGaAcceptanceTestConfig(testContextParameters())
        val cron = NightlyTriggerConfiguration(
            branch = "refs/heads/experimental-nightly",
            nightlyTestsEnabled = false,
            startHour = 7
        )
        val project = nightlyTests("EXPERIMENTAL", ProviderNameGa, HashiCorpVCSRootGa, config, cron)
        val composite = getBuildFromProject(project, AllNightlyTestsName)
        val schedule = composite.triggers.items.single() as ScheduleTrigger
        assertEquals(false, schedule.enabled)
        assertEquals("+:${cron.branch}", schedule.branchFilter)
        assertEquals("7", (schedule.schedulingPolicy as ScheduleTrigger.SchedulingPolicy.Cron).hours)
        assertFinishTrigger(getBuildFromProject(project, ServiceSweeperName), composite, cron.branch)
    }

    @Test
    fun nightlyPackagesAndSweepersShouldRetainProviderLocks() {
        val root = googleCloudRootProject(testContextParameters())
        listOf(
            Triple(gaProjectName, ProviderNameGa, SharedResourceNameGa),
            Triple(betaProjectName, ProviderNameBeta, SharedResourceNameBeta)
        ).forEach { (name, provider, resource) ->
            val project = getNestedProjectFromRoot(root, name, nightlyTestsProjectName)
            val composite = getBuildFromProject(project, AllNightlyTestsName)
            val sweeper = getBuildFromProject(project, ServiceSweeperName)
            assertTrue("Composite should not hold locks needed by package builds", composite.features.items.filterIsInstance<SharedResources>().isEmpty())
            assertSharedResourceLocks(sweeper, SharedResources { lockAllValues(resource) })
            project.buildTypes.filter { it != composite && it != sweeper }.forEach { build ->
                val path = build.params.findRawParam("PACKAGE_PATH")!!.value
                val packageName = getAllPackageInProviderVersion(provider).entries.single { it.value.getValue("path") == path }.key
                assertSharedResourceLocks(build, SharedResources { lockSpecificValue(resource, packageName) })
            }
        }
    }
}
