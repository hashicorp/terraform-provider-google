## 1.0.2 (Unreleased)

FEATURES:
* **New Resource:** `google_logging_folder_sink` ([#470](https://github.com/terraform-providers/terraform-provider-google/pull/470))
* **New Resource:** `google_organization_policy` ([#523](https://github.com/terraform-providers/terraform-provider-google/pull/523))
* **New Resource:** `google_compute_target_tcp_proxy` ([#528](https://github.com/terraform-providers/terraform-provider-google/pull/528))
* **New Resource:** `google_compute_region_autoscaler` ([#544](https://github.com/terraform-providers/terraform-provider-google/pull/544))
* **New Resources:** `google_compute_shared_vpc_host_project` and `google_compute_shared_vpc_service_project` ([#544](https://github.com/terraform-providers/terraform-provider-google/pull/572))

IMPROVEMENTS:
* compute: Generate network link without calling network API in `google_compute_subnetwork` [GH-527]
* compute: Generate network link without calling network API in `google_compute_vpn_gateway` and `google_compute_router` [GH-527]
* compute: Add import support to `google_compute_target_tcp_proxy` [GH-534]
* compute: Add labels support to `google_compute_instance_template` ([#17](https://github.com/terraform-providers/terraform-provider-google/issues/17))
* compute: `google_vpn_tunnel` - Mark 'shared_secret' as sensitive [GH-561]
* container: Allow disabling of Kubernetes Dashboard via `kubernetes_dashboard` addon ([#433](https://github.com/terraform-providers/terraform-provider-google/issues/433))
* container: Merge the schemas and logic for the node pool resource and the node pool field in the cluster to aid in maintainability [GH-489]
* container: Add master_version to container cluster [GH-538]
* sql: Add new retry wrapper fn, retry sql database instance operations that commonly 503 [GH-417]
* pubsub: `push_config` field for a `google_pubsub_subscription` is not updateable [GH-512]

BUG FIXES:
* compute: Fix bug in `google_compute_instance` preventing the `assigned_nat_ip` field from ever getting assigned [GH-536]
* compute: Fix bug in `google_compute_firewall` causing the beta APIs even if no beta features are used [GH-500]
* compute: Fix bug in `google_network_peering` preventing creating a peering for a network outside the provider default project [GH-496]
* compute: Fix BackendService group hash when instance groups use beta features [GH-522]
* compute: Make `disk.device_name` computed in `google_compute_instance_template` [GH-566]
* dns: Error out if DNS zone is not found [GH-560]
* container: Fix crash when creating node pools with `name_prefix` or no name [GH-531]
* container: Fix cluster version upgrades [GH-577]

## 1.0.1 (October 02, 2017)

BUG FIXES:
* compute: Fix bug that prevented the state migration for `google_compute_instance` from updating to use attached_disk, boot_disk, and scratch_disk. ([#511](https://github.com/terraform-providers/terraform-provider-google/issues/511))
* compute: Fix bug causing a crash if the API returns an error on `google_compute_instance` creation [GH-556]

## 1.0.0 (October 02, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:
* compute: A state migration was added to convert `google_compute_instance.disk` fields into the correct one of `attached_disk`, `boot_disk`, or `scratch_disk`. This will lead to plan-time diffs for anyone still using the `disk` field. Please verify its results carefully and update configs appropriately.
* container: `google_container_cluster.node_pool.initial_node_count` is now deprecated. Please replace with `google_container_cluster.node_pool.node_count` instead. ([#331](https://github.com/terraform-providers/terraform-provider-google/issues/331))
* storage: `google_storage_bucket_acl` now sets the bucket ACL to whatever is in the config, correcting any drift. This means any permissions set automatically by GCP (e.g., project-viewers-\* policies, etc.) will be removed unless they're added to your config. Also, the `OWNER:project-owners-{project-id}` will never be deleted, as the API won't allow it. This is now correctly handled, and it is removed from state without being deleted in the API. ([#358](https://github.com/terraform-providers/terraform-provider-google/issues/358)] [[#439](https://github.com/terraform-providers/terraform-provider-google/issues/439))

FEATURES:
* **New Data Source:** `google_client_config` ([#385](https://github.com/terraform-providers/terraform-provider-google/issues/385))
* **New Resource:** `google_compute_region_instance_group_manager` ([#394](https://github.com/terraform-providers/terraform-provider-google/issues/394))
* **New Resource:** `google_folder` ([#416](https://github.com/terraform-providers/terraform-provider-google/issues/416))
* **New Resource:** `google_folder_iam_policy` ([#447](https://github.com/terraform-providers/terraform-provider-google/issues/447))
* **New Resource:** `google_logging_project_sink` ([#432](https://github.com/terraform-providers/terraform-provider-google/issues/432))
* **New Resource:** `google_logging_billing_account_sink` ([#457](https://github.com/terraform-providers/terraform-provider-google/issues/457))

IMPROVEMENTS:
* bigquery: Support Bigquery Views ([#230](https://github.com/terraform-providers/terraform-provider-google/issues/230))
* container: Add import support for `google_container_cluster` ([#391](https://github.com/terraform-providers/terraform-provider-google/issues/391))
* container: Add support for resizing a node pool defined in `google_container_cluster` ([#331](https://github.com/terraform-providers/terraform-provider-google/issues/331))
* container: Allow updating `google_container_cluster.logging_service` ([#343](https://github.com/terraform-providers/terraform-provider-google/issues/343))
* container: Add support for 'node_config.preemptible' field on `google_container_cluster` ([#341](https://github.com/terraform-providers/terraform-provider-google/issues/341))
* container: Allow min node counts of 0 for node pool autoscaling ([#468](https://github.com/terraform-providers/terraform-provider-google/issues/468))
* compute: Add support for 'labels' field on `google_compute_image` ([#339](https://github.com/terraform-providers/terraform-provider-google/issues/339))
* compute: Add support for 'labels' field on `google_compute_disk` ([#344](https://github.com/terraform-providers/terraform-provider-google/issues/344))
* compute: Add support for `labels` field on `google_compute_global_forwarding_rule` ([#354](https://github.com/terraform-providers/terraform-provider-google/issues/354))
* compute: Add support for 'guest_accelerators' (GPU) on `google_compute_instance` ([#330](https://github.com/terraform-providers/terraform-provider-google/issues/330))
* compute: Add support for 'priority' field on `google_compute_firewall` ([#342](https://github.com/terraform-providers/terraform-provider-google/issues/342))
* compute: `google_compute_firewall` network field now supports self_link in addition of name ([#477](https://github.com/terraform-providers/terraform-provider-google/issues/477))
* compute: Add support for 'min_cpu_platform' in `google_compute_instance` ([#349](https://github.com/terraform-providers/terraform-provider-google/issues/349))
* compute: Add support for 'alias_ip_range' in `google_compute_instance` ([#375](https://github.com/terraform-providers/terraform-provider-google/issues/375))
* compute: Add support for computed field 'instance_id' in `google_compute_instance` ([#427](https://github.com/terraform-providers/terraform-provider-google/issues/427))
* compute: Improve import for `google_compute_address` to support multiple id formats. ([#378](https://github.com/terraform-providers/terraform-provider-google/issues/378))
* compute: Add state migration from `disk` to boot_disk/scratch_disk/attached_disk ([#329](https://github.com/terraform-providers/terraform-provider-google/issues/329))
* compute: Mark certificate as sensitive within `google_compute_ssl_certificate` ([#490](https://github.com/terraform-providers/terraform-provider-google/issues/490))
* project: Add support for 'labels' field on `google_project` ([#383](https://github.com/terraform-providers/terraform-provider-google/issues/383))
* project: Move a `google_project` in and out of a folder ([#438](https://github.com/terraform-providers/terraform-provider-google/issues/438))
* pubsub: Add import support for `google_pubsub_topic`. ([#392](https://github.com/terraform-providers/terraform-provider-google/issues/392))
* pubsub: Add import support for `google_pubsub_subscription`. ([#456](https://github.com/terraform-providers/terraform-provider-google/issues/456))
* sql: Add support for `connection_name` in `google_sql_database_instance` ([#387](https://github.com/terraform-providers/terraform-provider-google/issues/387))
* storage: Add support for versioning in `google_storage_bucket` ([#381](https://github.com/terraform-providers/terraform-provider-google/issues/381))

BUG FIXES:
* compute/sql: Fix a few instances where we read the project from the provider config and not using the helper function ([#469](https://github.com/terraform-providers/terraform-provider-google/issues/469))
* compute: Fix bug with CSEK where the key stored in state might be associated with the wrong disk ([#327](https://github.com/terraform-providers/terraform-provider-google/issues/327))
* compute: Fix bug where 'session_affinity' would get reset on `google_compute_backend_service` resource ([#348](https://github.com/terraform-providers/terraform-provider-google/issues/348))
* sql: Fixed bug where ip_address elements were offset incorrectly ([#352](https://github.com/terraform-providers/terraform-provider-google/issues/352))
* sql: Fixed bug where default user on replica would cause an incorrect delete api call ([#347](https://github.com/terraform-providers/terraform-provider-google/issues/347))
* project: Fixed bug where deleting a project outside Terraform would cause `google_project` to fail. ([#466](https://github.com/terraform-providers/terraform-provider-google/issues/466))
* pubsub: Fixed bug where `google_pubsub_subscription` did not read its state from the API. ([#456](https://github.com/terraform-providers/terraform-provider-google/issues/456))

## 0.1.3 (August 17, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:
* bigtable: `num_nodes` in `google_bigtable_instance` no longer defaults to `3`; if you used that default, it will need to be explicitly set. ([#313](https://github.com/terraform-providers/terraform-provider-google/issues/313))
* compute: `automatic_restart` and `on_host_maintenance` have been removed from `google_compute_instance_template`. Use `scheduling.automatic_restart` or `scheduling.on_host_maintenance` instead. ([#224](https://github.com/terraform-providers/terraform-provider-google/issues/224))

FEATURES:
* **New Data Source:** `google_compute_instance_group` ([#267](https://github.com/terraform-providers/terraform-provider-google/issues/267))
* **New Data Source:** `google_dns_managed_zone` ([#268](https://github.com/terraform-providers/terraform-provider-google/issues/268))
* **New Resource:** `google_compute_project_metadata_item` - allows management of single key/value pairs within the project metadata map ([#176](https://github.com/terraform-providers/terraform-provider-google/issues/176))
* **New Resource:** `google_project_iam_binding` - allows fine-grained control of a project's IAM policy, controlling only a single binding. ([#171](https://github.com/terraform-providers/terraform-provider-google/issues/171))
* **New Resource:** `google_project_iam_member` - allows fine-grained control of a project's IAM policy, controlling only a single member in a binding. ([#171](https://github.com/terraform-providers/terraform-provider-google/issues/171))
* **New Resource:** `google_compute_network_peering` ([#259](https://github.com/terraform-providers/terraform-provider-google/issues/259))
* **New Resource:** `google_runtimeconfig_config` - allows creating, updating and deleting Google RuntimeConfig resources ([#315](https://github.com/terraform-providers/terraform-provider-google/issues/315))
* **New Resource:** `google_runtimeconfig_variable` - allows creating, updating, and deleting Google RuntimeConfig variables ([#315](https://github.com/terraform-providers/terraform-provider-google/issues/315))
* **New Resource:** `google_sourcerepo_repository` - allows creating and deleting Google Source Repositories ([#256](https://github.com/terraform-providers/terraform-provider-google/issues/256))
* **New Resource:** `google_spanner_instance` - allows creating, updating and deleting Google Spanner Instance ([#270](https://github.com/terraform-providers/terraform-provider-google/issues/270))
* **New Resource:** `google_spanner_database` - allows creating, updating and deleting Google Spanner Database ([#271](https://github.com/terraform-providers/terraform-provider-google/issues/271))

IMPROVEMENTS:
* bigtable: Add support for `instance_type` to `google_bigtable_instance`. ([#313](https://github.com/terraform-providers/terraform-provider-google/issues/313))
* compute: Add import support for `google_compute_subnetwork` ([#227](https://github.com/terraform-providers/terraform-provider-google/issues/227))
* compute: Add import support for `google_container_node_pool` ([#284](https://github.com/terraform-providers/terraform-provider-google/issues/284))
* compute: Change google_container_node_pool ID format to zone/cluster/name to remove artificial restriction on node pool name across clusters ([#304](https://github.com/terraform-providers/terraform-provider-google/issues/304))
* compute: Add support for `auto_healing_policies` to `google_compute_instance_group_manager` ([#249](https://github.com/terraform-providers/terraform-provider-google/issues/249))
* compute: Add support for `ip_version` to `google_compute_global_forwarding_rule` ([#265](https://github.com/terraform-providers/terraform-provider-google/issues/265))
* compute: Add support for `ip_version` to `google_compute_global_address` ([#250](https://github.com/terraform-providers/terraform-provider-google/issues/250))
* compute: Add support for `subnetwork` as a self_link to `google_compute_instance`. ([#290](https://github.com/terraform-providers/terraform-provider-google/issues/290))
* compute: Add support for `secondary_ip_range` to `google_compute_subnetwork`. ([#310](https://github.com/terraform-providers/terraform-provider-google/issues/310))
* compute: Add support for multiple `network_interface`'s to `google_compute_instance`. ([#289](https://github.com/terraform-providers/terraform-provider-google/issues/289))
* compute: Add support for `denied` to `google_compute_firewall` ([#282](https://github.com/terraform-providers/terraform-provider-google/issues/282))
* compute: Add support for egress traffic using `direction` to `google_compute_firewall` ([#306](https://github.com/terraform-providers/terraform-provider-google/issues/306))
* compute: When disks are created from snapshots, both snapshot names and URLs may be used ([#238](https://github.com/terraform-providers/terraform-provider-google/issues/238))
* container: Add support for node pool autoscaling ([#157](https://github.com/terraform-providers/terraform-provider-google/issues/157))
* container: Add NodeConfig support on `google_container_node_pool` ([#184](https://github.com/terraform-providers/terraform-provider-google/issues/184))
* container: Add support for legacyAbac to `google_container_cluster` ([#261](https://github.com/terraform-providers/terraform-provider-google/issues/261))
* container: Allow configuring node_config of node_pools specified in `google_container_cluster` ([#299](https://github.com/terraform-providers/terraform-provider-google/issues/299))
* sql: Persist state from the API for `google_sql_database_instance` regardless of what attributes the user has set ([#208](https://github.com/terraform-providers/terraform-provider-google/issues/208))
* storage: Buckets now can have lifecycle properties ([#6](https://github.com/terraform-providers/terraform-provider-google/pull/6))

BUG FIXES:
* bigquery: Fix type panic on expiration_time ([#209](https://github.com/terraform-providers/terraform-provider-google/issues/209))
* compute: Marked 'private_key' as sensitive ([#220](https://github.com/terraform-providers/terraform-provider-google/pull/220))
* compute: Fix disk type "Malformed URL" error on `google_compute_instance` boot disks ([#275](https://github.com/terraform-providers/terraform-provider-google/issues/275))
* compute: Refresh `google_compute_autoscaler` using the `zone` set in state instead of scanning for the first one with a matching name in the provider region. ([#193](https://github.com/terraform-providers/terraform-provider-google/issues/193))
* compute: `google_compute_instance` reads `scheduling` fields from GCP ([#237](https://github.com/terraform-providers/terraform-provider-google/issues/237))
* compute: Fix bug where `scheduling.automatic_restart` set to false on `google_compute_instance_template` would force recreate ([#224](https://github.com/terraform-providers/terraform-provider-google/issues/224))
* container: Fix error if `google_container_node_pool` deleted out of band ([#293](https://github.com/terraform-providers/terraform-provider-google/issues/293))
* container: Fail when both name and name_prefix are set for node_pool in `google_container_cluster` ([#296](https://github.com/terraform-providers/terraform-provider-google/issues/296))
* container: Allow upgrading GKE versions and provide better error message handling ([#291](https://github.com/terraform-providers/terraform-provider-google/issues/291))

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
