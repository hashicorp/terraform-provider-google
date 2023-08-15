/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is copied from mmv1, any changes made here will be overwritten

package tests

import org.junit.Assert.assertTrue
import org.junit.Assert.assertFalse
import org.junit.Test
import ShouldAddTrigger
import MM_UPSTREAM
import MAJOR_RELEASE_TESTING

class HelperTests {
    @Test
    fun funShouldAddTrigger_random_string() {
        val environment = "foobar"
        assertTrue("Cron triggers should be added to projects with a random environment value" , ShouldAddTrigger(environment))
    }

    @Test
    fun funShouldAddTrigger_MAJOR_RELEASE_TESTING() {
        val environment = MAJOR_RELEASE_TESTING
        assertTrue("Cron triggers should be added to projects used for testing the 5.0.0 major release" , ShouldAddTrigger(environment))
    }

    @Test
    fun funShouldAddTrigger_MM_UPSTREAM() {
        val environment = MM_UPSTREAM
        assertFalse("Cron triggers should NOT be added to projects using the MM upstream repo" , ShouldAddTrigger(environment))
    }
}