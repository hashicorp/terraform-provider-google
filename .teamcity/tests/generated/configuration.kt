/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

package tests

import Google
import org.junit.Assert.assertTrue
import org.junit.Test
import useTeamCityGoTest

class ConfigurationTests {
    @Test
    fun buildShouldFailOnError() {
        val project = Google("default", testVcsRootId(), "refs/heads/main", testConfiguration())
        project.buildTypes.forEach { bt ->
            assertTrue("Build '${bt.id}' should fail on errors!", bt.failureConditions.errorMessage)
        }
    }

    @Test
    fun buildShouldHaveGoTestFeature() {
        val project = Google("default",  testVcsRootId(), "refs/heads/main",testConfiguration())
        project.buildTypes.forEach{ bt ->
            var exists = false
            bt.features.items.forEach { f ->
                if (f.type == "golang") {
                    exists = true
                }
            }

            if (useTeamCityGoTest) {
                assertTrue("Build %s doesn't have Go Test Json enabled".format(bt.name), exists)
            }
        }
    }

    // Commented out because it's not true that all builds have triggers now.
    // Once I have the ability to run tests I'll address this - writing new tests for the new config
    // @Test
    // fun buildShouldHaveTrigger() {
    //     val project = Google("default",  testVcsRootId(), "refs/heads/main", testConfiguration())
    //     var exists = false
    //     project.buildTypes.forEach{ bt ->
    //         bt.triggers.items.forEach { t ->
    //             if (t.type == "schedulingTrigger") {
    //                 exists = true
    //             }
    //         }
    //     }
    //     assertTrue("The Build Configuration should have a Trigger", exists)
    // }
}
