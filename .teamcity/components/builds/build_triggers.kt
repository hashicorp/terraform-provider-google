/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package builds

import DefaultDaysOfMonth
import DefaultDaysOfWeek
import DefaultStartHour
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Triggers
import jetbrains.buildServer.configs.kotlin.triggers.schedule

class NightlyTriggerConfiguration(
    val nightlyTestsEnabled: Boolean = true,
    val startHour: Int = DefaultStartHour,
    val daysOfWeek: String = DefaultDaysOfWeek,
    val daysOfMonth: String = DefaultDaysOfMonth
)

fun Triggers.runNightly(config: NightlyTriggerConfiguration) {
    val filter = "+:refs/heads/main"

    schedule{
        enabled = config.nightlyTestsEnabled
        branchFilter = filter
        triggerBuild = always() // Run build even if no new commits/pending changes
        withPendingChangesOnly = false
        enforceCleanCheckout = true

        schedulingPolicy = cron {
            hours = config.startHour.toString()
            timezone = "SERVER"

            dayOfWeek = config.daysOfWeek
            dayOfMonth = config.daysOfMonth
        }
    }
}

// BuildType.addTrigger enables adding a CRON trigger after a build configuration has been initialised
fun BuildType.addTrigger(triggerConfig: NightlyTriggerConfiguration){
    triggers {
        runNightly(triggerConfig)
    }
}
