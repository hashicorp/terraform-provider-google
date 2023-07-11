// this file is copied from mmv1, any changes made here will be overwritten

// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 4

// specifies the default level of parallelism per-service-package
var defaultParallelism = 12

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "1.2.5"

// This represents a cron view of days of the week, Monday - Friday.
const val defaultDaysOfWeek = "*"

// Cron value for any day of month
const val defaultDaysOfMonth = "*"
