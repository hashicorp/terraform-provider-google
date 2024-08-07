/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package builds

import DefaultBranchName
import DefaultDaysOfMonth
import DefaultDaysOfWeek
import DefaultStartHour
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Triggers
import jetbrains.buildServer.configs.kotlin.triggers.schedule

class NightlyTriggerConfiguration(
    val branch: String = DefaultBranchName,
    val nightlyTestsEnabled: Boolean = true,
    val startHour: Int = DefaultStartHour,
    val daysOfWeek: String = DefaultDaysOfWeek,
    val daysOfMonth: String = DefaultDaysOfMonth
)

fun Triggers.runNightly(config: NightlyTriggerConfiguration) {

    schedule{
        enabled = config.nightlyTestsEnabled
        branchFilter = "+:refs/heads/nightly-test"
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
