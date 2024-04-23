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
        val project = googleCloudRootProject(testContextParameters())
        project.buildTypes.forEach { bt ->
            assertTrue("Build '${bt.id}' doesn't use clean checkout", bt.vcs.cleanCheckout)
        }
    }
}
