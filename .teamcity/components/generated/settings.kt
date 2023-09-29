/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 4

// specifies the default level of parallelism per-service-package
var defaultParallelism = 6

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "1.2.5"

// This represents a cron view of days of the week
const val defaultDaysOfWeek = "*"

// Cron value for any day of month
const val defaultDaysOfMonth = "*"

// Value used to make long-running builds fail due to a timeout
const val defaultBuildTimeoutDuration = 60 * 12 //12 hours in minutes

// Values that `environment` parameter is checked against,
// when deciding to change how TeamCity objects are configured
const val MAJOR_RELEASE_TESTING = "major-release-5.0.0-testing"
const val MM_UPSTREAM = "mm-upstream"
