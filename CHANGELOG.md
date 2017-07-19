## 0.1.2 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `google_sql_database_instance`: a limited number of fields will be read during import because of [GH-114]
* `google_sql_database_instance`: `name`, `region`, `database_version`, and `master_instance_name` fields are now updated during a refresh and may display diffs

FEATURES:

* **New Resource:** `google_bigtable_instance` [GH-177]
* **New Resource:** `google_bigtable_table` [GH-177]
* **New Resource:** `google_compute_project_metadata_item` - allows management of single key/value pairs within the project metadata map [GH-176]

IMPROVEMENTS:

* compute: Add `boot_disk` property to `google_compute_instance` [GH-122]
* compute: Add `scratch_disk` property to `google_compute_instance` and deprecate `disk` [GH-123]
* compute: Add `labels` property to `google_compute_instance` [GH-150]
* compute: Add import support for `google_compute_image` [GH-194]
* compute: Add import support for `google_compute_https_health_check` [GH-213]
* container: Add timeout support ([#13203](https://github.com/hashicorp/terraform/issues/13203))
* container: Allow adding/removing zones to/from GKE clusters without recreating them [GH-152]
* project: Allow unlinking of billing account [GH-138]
* sql: Add support for importing `google_sql_database` [GH-12]
* sql: Add support for importing `google_sql_database_instance` [GH-11]
* sql: Add `charset` and `collation` properties to `google_sql_database` [GH-183]

BUG FIXES:

* compute: `compute_firewall` will no longer display a perpetual diff if `source_ranges` isn't set [GH-147]
* compute: Fix read method + test/document import for `google_compute_health_check` [GH-155]
* compute: Read named ports changes properly in `google_compute_instance_group` [GH-188]
* compute: `google_compute_image` `description` property can now be set [GH-199] 
* compute: `google_compute_target_https_proxy` will no longer display a diff if ssl certificates are referenced using only the path [GH-210]

## 0.1.1 (June 21, 2017)

BUG FIXES: 

* compute: Restrict the number of health_checks in Backend Service resources to 1. ([#145](https://github.com/terraform-providers/terraform-provider-google/issues/145))

## 0.1.0 (June 20, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `compute_disk.image`: shorthand for disk images is no longer supported, and will display a diff if used ([#1](https://github.com/terraform-providers/terraform-provider-google/issues/1))

IMPROVEMENTS:

* compute: Add support for importing `compute_backend_service` ([#40](https://github.com/terraform-providers/terraform-provider-google/issues/40))
* compute: Wait for disk resizes to complete ([#1](https://github.com/terraform-providers/terraform-provider-google/issues/1))
* compute: Support `connection_draining_timeout_sec` in `google_compute_region_backend_service` ([#101](https://github.com/terraform-providers/terraform-provider-google/issues/101))
* compute: Made `path_rule` optional in `google_compute_url_map`'s `path_matcher` block ([#118](https://github.com/terraform-providers/terraform-provider-google/issues/118))
* container: Add support for labels and tags on GKE node_config ([#7](https://github.com/terraform-providers/terraform-provider-google/issues/7))
* sql: Add an additional delay when checking for sql operations ([#15170](https://github.com/hashicorp/terraform/pull/15170))

BUG FIXES:

* compute: Changed `google_compute_instance_group_manager` `target_size` default to 0 ([#65](https://github.com/terraform-providers/terraform-provider-google/issues/65))
* storage: Represent GCS Bucket locations as uppercase in state. ([#117](https://github.com/terraform-providers/terraform-provider-google/issues/117))
