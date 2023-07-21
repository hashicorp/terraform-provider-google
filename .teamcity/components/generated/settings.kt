// this file is auto-generated with mmv1, any changes made here will be overwritten

// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 4

// specifies the default level of parallelism per-service-package
var defaultParallelism = 12

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "1.2.5"

// This represents a cron view of days of the week
const val defaultDaysOfWeek = "1-3,5-7" // All nights except Thursday for GA; feature branch testing happens on Thursdays

// Cron value for any day of month
const val defaultDaysOfMonth = "*"

// Values that `environment` parameter is checked against,
// when deciding to change how TeamCity objects are configured
const val MAJOR_RELEASE_TESTING = "major-release-5.0.0-testing"
