## 0.1.3 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

FEATURES:
* **New Data Source:** `google_dns_managed_zone` [GH-268]
* **New Resource:** `google_compute_project_metadata_item` - allows management of single key/value pairs within the project metadata map [GH-176]
* **New Resource:** `google_project_iam_binding` - allows fine-grained control of a project's IAM policy, controlling only a single binding. [GH-171]
* **New Resource:** `google_project_iam_member` - allows fine-grained control of a project's IAM policy, controlling only a single member in a binding. [GH-171]
* **New Resource:** `google_compute_network_peering` ([#259](https://github.com/terraform-providers/terraform-provider-google/issues/259))

IMPROVEMENTS:
* compute: Add import support for `google_compute_subnetwork` ([#227](https://github.com/terraform-providers/terraform-provider-google/issues/227))
* container: Add support for node pool autoscaling [GH-157]
* container: Add NodeConfig support on `google_container_node_pool` [GH-184]
* container: Add support for legacyAbac to `google_container_cluster` [GH-261]
* sql: Persist state from the API for `google_sql_database_instance` regardless of what attributes the user has set [GH-208]
* storage: Buckets now can have lifecycle properties ([#6](https://github.com/terraform-providers/terraform-provider-google/pull/6))

BUG FIXES:
* bigquery: Fix type panic on expiration_time [GH-209]
* compute: Marked 'private_key' as sensitive ([#220](https://github.com/terraform-providers/terraform-provider-google/pull/220))
* compute: Fix disk type "Malformed URL" error on `google_compute_instance` boot disks [GH-275]

## 0.1.2 (July 20, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `google_sql_database_instance`: a limited number of fields will be read during import because of ([#114](https://github.com/terraform-providers/terraform-provider-google/issues/114))
* `google_sql_database_instance`: `name`, `region`, `database_version`, and `master_instance_name` fields are now updated during a refresh and may display diffs

FEATURES:

* **New Resource:** `google_bigtable_instance` ([#177](https://github.com/terraform-providers/terraform-provider-google/issues/177))
* **New Resource:** `google_bigtable_table` ([#177](https://github.com/terraform-providers/terraform-provider-google/issues/177))

IMPROVEMENTS:

* compute: Add `boot_disk` property to `google_compute_instance` ([#122](https://github.com/terraform-providers/terraform-provider-google/issues/122))
* compute: Add `scratch_disk` property to `google_compute_instance` and deprecate `disk` ([#123](https://github.com/terraform-providers/terraform-provider-google/issues/123))
* compute: Add `labels` property to `google_compute_instance` ([#150](https://github.com/terraform-providers/terraform-provider-google/issues/150))
* compute: Add import support for `google_compute_image` ([#194](https://github.com/terraform-providers/terraform-provider-google/issues/194))
* compute: Add import support for `google_compute_https_health_check` ([#213](https://github.com/terraform-providers/terraform-provider-google/issues/213))
* compute: Add import support for `google_compute_instance_group` ([#201](https://github.com/terraform-providers/terraform-provider-google/issues/201))
* container: Add timeout support ([#13203](https://github.com/hashicorp/terraform/issues/13203))
* container: Allow adding/removing zones to/from GKE clusters without recreating them ([#152](https://github.com/terraform-providers/terraform-provider-google/issues/152))
* project: Allow unlinking of billing account ([#138](https://github.com/terraform-providers/terraform-provider-google/issues/138))
* sql: Add support for importing `google_sql_database` ([#12](https://github.com/terraform-providers/terraform-provider-google/issues/12))
* sql: Add support for importing `google_sql_database_instance` ([#11](https://github.com/terraform-providers/terraform-provider-google/issues/11))
* sql: Add `charset` and `collation` properties to `google_sql_database` ([#183](https://github.com/terraform-providers/terraform-provider-google/issues/183))

BUG FIXES:

* compute: `compute_firewall` will no longer display a perpetual diff if `source_ranges` isn't set ([#147](https://github.com/terraform-providers/terraform-provider-google/issues/147))
* compute: Fix read method + test/document import for `google_compute_health_check` ([#155](https://github.com/terraform-providers/terraform-provider-google/issues/155))
* compute: Read named ports changes properly in `google_compute_instance_group` ([#188](https://github.com/terraform-providers/terraform-provider-google/issues/188))
* compute: `google_compute_image` `description` property can now be set [[#199](https://github.com/terraform-providers/terraform-provider-google/issues/199)] 
* compute: `google_compute_target_https_proxy` will no longer display a diff if ssl certificates are referenced using only the path ([#210](https://github.com/terraform-providers/terraform-provider-google/issues/210))

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
