---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Configuration: https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/backupdr/BackupPlan.yaml
#     Template:      https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.html.markdown.tmpl
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Backup and DR Service"
description: |-
  A backup plan defines when and how to back up a resource, including the backup's schedule, retention, and location.
---

# google_backup_dr_backup_plan

A backup plan defines when and how to back up a resource, including the backup's schedule, retention, and location.


To get more information about BackupPlan, see:

* [API documentation](https://cloud.google.com/backup-disaster-recovery/docs/reference/rest)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/backup-disaster-recovery/docs)

## Example Usage - Backup Dr Backup Plan Simple


```hcl
resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                      = "us-central1"
  backup_vault_id                               = "backup-vault-simple-test"
  backup_minimum_enforced_retention_duration    = "100000s"
}

resource "google_backup_dr_backup_plan" "my-backup-plan-1" {
  location       = "us-central1"
  backup_plan_id = "backup-plan-simple-test"
  resource_type  = "compute.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id                = "rule-1"
    backup_retention_days  = 5

    standard_schedule {
      recurrence_type     = "HOURLY"
      hourly_frequency    = 6
      time_zone           = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 24
      }
    }
  }
}
```
## Example Usage - Backup Dr Backup Plan For Disk Resource


```hcl
resource "google_backup_dr_backup_vault" "my_backup_vault" {
  provider = google-beta
  location                                      = "us-central1"
  backup_vault_id                               = "backup-vault-disk-test"
  backup_minimum_enforced_retention_duration    = "100000s"
}

resource "google_backup_dr_backup_plan" "my-disk-backup-plan-1" {
  provider       = google-beta
  location       = "us-central1"
  backup_plan_id = "backup-plan-disk-test"
  resource_type  = "compute.googleapis.com/Disk"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id                = "rule-1"
    backup_retention_days  = 5

    standard_schedule {
      recurrence_type     = "HOURLY"
      hourly_frequency    = 1
      time_zone           = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 6
      }
    }
  }
}
```
## Example Usage - Backup Dr Backup Plan For Csql Resource


```hcl
resource "google_backup_dr_backup_vault" "my_backup_vault" {
  location                                      = "us-central1"
  backup_vault_id                               = "backup-vault-csql-test"
  backup_minimum_enforced_retention_duration    = "100000s"
}

resource "google_backup_dr_backup_plan" "my-csql-backup-plan-1" {
  location       = "us-central1"
  backup_plan_id = "backup-plan-csql-test"
  resource_type  = "sqladmin.googleapis.com/Instance"
  backup_vault   = google_backup_dr_backup_vault.my_backup_vault.id

  backup_rules {
    rule_id                = "rule-1"
    backup_retention_days  = 5

    standard_schedule {
      recurrence_type     = "HOURLY"
      hourly_frequency    = 6
      time_zone           = "UTC"

      backup_window {
        start_hour_of_day = 0
        end_hour_of_day   = 6
      }
    }
  }
  log_retention_days = 4
}
```

## Argument Reference

The following arguments are supported:


* `backup_vault` -
  (Required)
  Backup vault where the backups gets stored using this Backup plan.

* `resource_type` -
  (Required)
  The resource type to which the `BackupPlan` will be applied.
  Examples include, "compute.googleapis.com/Instance", "compute.googleapis.com/Disk", "sqladmin.googleapis.com/Instance" and "storage.googleapis.com/Bucket".

* `backup_rules` -
  (Required)
  The backup rules for this `BackupPlan`. There must be at least one `BackupRule` message.
  Structure is [documented below](#nested_backup_rules).

* `location` -
  (Required)
  The location for the backup plan

* `backup_plan_id` -
  (Required)
  The ID of the backup plan


* `description` -
  (Optional)
  The description allows for additional details about `BackupPlan` and its use cases to be provided.

* `log_retention_days` -
  (Optional)
  This is only applicable for CloudSql resource. Days for which logs will be stored. This value should be greater than or equal to minimum enforced log retention duration of the backup vault.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.



<a name="nested_backup_rules"></a>The `backup_rules` block supports:

* `rule_id` -
  (Required)
  The unique ID of this `BackupRule`. The `rule_id` is unique per `BackupPlan`.

* `backup_retention_days` -
  (Required)
  Configures the duration for which backup data will be kept. The value should be greater than or equal to minimum enforced retention of the backup vault.

* `standard_schedule` -
  (Required)
  StandardSchedule defines a schedule that runs within the confines of a defined window of days.
  Structure is [documented below](#nested_backup_rules_backup_rules_standard_schedule).


<a name="nested_backup_rules_backup_rules_standard_schedule"></a>The `standard_schedule` block supports:

* `recurrence_type` -
  (Required)
  RecurrenceType enumerates the applicable periodicity for the schedule.
  Possible values are: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`, `YEARLY`.

* `hourly_frequency` -
  (Optional)
  Specifies frequency for hourly backups. An hourly frequency of 2 means jobs will run every 2 hours from start time till end time defined.
  This is required for `recurrence_type`, `HOURLY` and is not applicable otherwise.

* `days_of_week` -
  (Optional)
  Specifies days of week like MONDAY or TUESDAY, on which jobs will run. This is required for `recurrence_type`, `WEEKLY` and is not applicable otherwise.
  Each value may be one of: `DAY_OF_WEEK_UNSPECIFIED`, `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY`.

* `days_of_month` -
  (Optional)
  Specifies days of months like 1, 5, or 14 on which jobs will run.

* `week_day_of_month` -
  (Optional)
  Specifies a week day of the month like FIRST SUNDAY or LAST MONDAY, on which jobs will run.
  Structure is [documented below](#nested_backup_rules_backup_rules_standard_schedule_week_day_of_month).

* `months` -
  (Optional)
  Specifies values of months
  Each value may be one of: `MONTH_UNSPECIFIED`, `JANUARY`, `FEBRUARY`, `MARCH`, `APRIL`, `MAY`, `JUNE`, `JULY`, `AUGUST`, `SEPTEMBER`, `OCTOBER`, `NOVEMBER`, `DECEMBER`.

* `time_zone` -
  (Required)
  The time zone to be used when interpreting the schedule.

* `backup_window` -
  (Optional)
  A BackupWindow defines the window of the day during which backup jobs will run. Jobs are queued at the beginning of the window and will be marked as
  `NOT_RUN` if they do not start by the end of the window.
  Structure is [documented below](#nested_backup_rules_backup_rules_standard_schedule_backup_window).


<a name="nested_backup_rules_backup_rules_standard_schedule_week_day_of_month"></a>The `week_day_of_month` block supports:

* `week_of_month` -
  (Required)
  WeekOfMonth enumerates possible weeks in the month, e.g. the first, third, or last week of the month.
  Possible values are: `WEEK_OF_MONTH_UNSPECIFIED`, `FIRST`, `SECOND`, `THIRD`, `FOURTH`, `LAST`.

* `day_of_week` -
  (Required)
  Specifies the day of the week.
  Possible values are: `DAY_OF_WEEK_UNSPECIFIED`, `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY`.

<a name="nested_backup_rules_backup_rules_standard_schedule_backup_window"></a>The `backup_window` block supports:

* `start_hour_of_day` -
  (Required)
  The hour of the day (0-23) when the window starts, for example, if the value of the start hour of the day is 6, that means the backup window starts at 6:00.

* `end_hour_of_day` -
  (Optional)
  The hour of the day (1-24) when the window ends, for example, if the value of end hour of the day is 10, that means the backup window end time is 10:00.
  The end hour of the day should be greater than the start

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/backupPlans/{{backup_plan_id}}`

* `name` -
  The name of backup plan resource created

* `backup_vault_service_account` -
  The Google Cloud Platform Service Account to be used by the BackupVault for taking backups.

* `supported_resource_types` -
  ([Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
  The list of all resource types to which the `BackupPlan` can be applied.

* `create_time` -
  When the `BackupPlan` was created.

* `update_time` -
  When the `BackupPlan` was last updated.


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 60 minutes.
- `update` - Default is 60 minutes.
- `delete` - Default is 60 minutes.

## Import


BackupPlan can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{location}}/backupPlans/{{backup_plan_id}}`
* `{{project}}/{{location}}/{{backup_plan_id}}`
* `{{location}}/{{backup_plan_id}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import BackupPlan using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{location}}/backupPlans/{{backup_plan_id}}"
  to = google_backup_dr_backup_plan.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), BackupPlan can be imported using one of the formats above. For example:

```
$ terraform import google_backup_dr_backup_plan.default projects/{{project}}/locations/{{location}}/backupPlans/{{backup_plan_id}}
$ terraform import google_backup_dr_backup_plan.default {{project}}/{{location}}/{{backup_plan_id}}
$ terraform import google_backup_dr_backup_plan.default {{location}}/{{backup_plan_id}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
