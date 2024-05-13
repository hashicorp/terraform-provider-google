/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject

class VcsTests {
    @Test
    fun buildsHaveCleanCheckOut() {
        val root = googleCloudRootProject(testContextParameters())

        val gaProject = getSubProject(root, gaProjectName)
        val betaProject = getSubProject(root, betaProjectName)
        val projectSweeperProject = getSubProject(root, betaProjectName)

        val allProjects = arrayListOf(gaProject, betaProject, projectSweeperProject)

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
