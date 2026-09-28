/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import AllProvidersNightlyTestsName
import DefaultBranchName
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject
import vcs_roots.HashiCorpVCSRootBeta
import vcs_roots.HashiCorpVCSRootBetaNightly
import vcs_roots.HashiCorpVCSRootGa
import vcs_roots.HashiCorpVCSRootGaNightly

class VcsTests {
    @Test
    fun nightlyTestsAndSweepersShouldDefaultToNightlyBranch() {
        val root = googleCloudRootProject(testContextParameters())
        listOf(
            gaProjectName to HashiCorpVCSRootGaNightly,
            betaProjectName to HashiCorpVCSRootBetaNightly
        ).forEach { (provider, vcsRoot) ->
            assertTrue("Nightly VCS root should be registered", root.roots.contains(vcsRoot))
            assertEquals(DefaultBranchName, vcsRoot.branch)
            val project = getNestedProjectFromRoot(root, provider, nightlyTestsProjectName)
            project.buildTypes.forEach { build ->
                assertEquals("${build.name} should use its provider's nightly VCS root", vcsRoot.id, build.vcs.entries.single().root.id)
            }
        }
        val composite = getBuildFromProject(root, AllProvidersNightlyTestsName)
        assertEquals(HashiCorpVCSRootGaNightly.id, composite.vcs.entries.single().root.id)
        val globalSweepers = getSubProject(root, globalSweepersProjectName)
        assertTrue("Gate should be branchless", getBuildFromProject(globalSweepers, "Nightly Sweeper Gate").vcs.entries.isEmpty())
        listOf("Project Sweeper", "Folder Sweeper").forEach { name ->
            val sweeper = getBuildFromProject(globalSweepers, name)
            assertEquals(HashiCorpVCSRootGaNightly.id, sweeper.vcs.entries.single().root.id)
        }
        assertEquals("refs/heads/main", HashiCorpVCSRootGa.branch)
        assertEquals("refs/heads/main", HashiCorpVCSRootBeta.branch)
    }

    @Test
    fun buildsHaveCleanCheckOut() {
        val root = googleCloudRootProject(testContextParameters())

        val gaProject = getSubProject(root, gaProjectName)
        val betaProject = getSubProject(root, betaProjectName)
        val globalSweepersProject = getSubProject(root, globalSweepersProjectName)

        val allProjects = arrayListOf(gaProject, betaProject, globalSweepersProject)

        allProjects.forEach { p ->
            p.subProjects.forEach { sp->
                // Test is created on assumption of project structure having max 2 layers of nested project (Root > Project A > Project B)
                assertTrue("TeamCity configuration is nested deeper than this test checks; test should be rewritten", sp.subProjects.size == 0)

                sp.buildTypes.forEach{ bt ->
                    assertTrue("Build '${bt.id}' should use clean checkout", bt.vcs.cleanCheckout)
                }
            }
        }
    }
}
