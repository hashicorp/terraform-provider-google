## 2.7.0 (Unreleased)

BACKWARDS INCOMPATIBILITIES:
* cloudfunctions: `google_cloudfunctions_function.runtime` now has an explicit default value of `nodejs6`. Users who have a different value set in the API but the value undefined in their config will see a diff. [GH-3605]

FEATURES: 
* **New Resources**: `google_compute_instance_iam_binding`, `google_compute_instance_iam_member`, and `google_compute_instance_iam_policy` are now available. ([#3551](https://github.com/terraform-providers/terraform-provider-google/pull/3551))

BUG FIXES:
* cloudfunctions: `google_cloudfunctions_function.runtime` now has an explicit default value of `nodejs6`. [GH-3605]
* monitoring: patching `google_monitoring_alert_policy` is more likely to succeed [GH-3587]

## 2.6.0 (May 07, 2019)

KNOWN ISSUES:
* cloudfunctions: `google_cloudfunctions_function`s without a `runtime` set will fail to create due to an upstream API change. You can work around this by setting an explicit `runtime` in `2.X` series releases.

DEPRECATIONS:
* monitoring: `google_monitoring_alert_policy` `labels` was deprecated, as the field was never used and it was typed incorrectly. ([#3494](https://github.com/terraform-providers/terraform-provider-google/issues/3494))

FEATURES: 
* **New Datasource**: `google_compute_node_types` for sole-tenant node types is now available. ([#3446](https://github.com/terraform-providers/terraform-provider-google/pull/3446))
* **New Resource**: `google_compute_node_group` for sole-tenant nodes is now available. ([#3514](https://github.com/terraform-providers/terraform-provider-google/pull/3514))
* **New Resource**: `google_compute_node_template` for sole-tenant nodes is now available. ([#3446](https://github.com/terraform-providers/terraform-provider-google/pull/3446))
* **New Resource**: `google_filestore_instance` is now available at GA. ([#3522](https://github.com/terraform-providers/terraform-provider-google/issues/3522))
* **New Resource**: `google_firestore_index` is now available to configure composite indexes on Firestore. ([#3484](https://github.com/terraform-providers/terraform-provider-google/issues/3484))
* **New Resource**: `google_logging_metric` is now available to configure Stackdriver logs-based metrics. ([#1702](https://github.com/GoogleCloudPlatform/magic-modules/pull/1702))
* **New Resources**: `google_compute_subnetwork_iam_binding`, `google_compute_subnetwork_iam_member`, and `google_compute_subnetwork_iam_policy` are now available at GA. ([#3541](https://github.com/terraform-providers/terraform-provider-google/issues/3541))

ENHANCEMENTS:
* dataflow: `google_dataflow_job`'s `network` and `subnetwork` can be configured. ([#3476](https://github.com/terraform-providers/terraform-provider-google/issues/3476))
* monitoring: `google_monitoring_alert_policy` `user_labels` support was added. ([#3494](https://github.com/terraform-providers/terraform-provider-google/issues/3494))
* compute: `google_compute_instance` and `google_compute_instance_template` now support node affinities for scheduling on sole tenant nodes [#3553](https://github.com/terraform-providers/terraform-provider-google/pull/3553)
* compute: `google_compute_region_backend_service` is now generated with Magic Modules, adding configurable timeouts, multiple import formats, `creation_timestamp` output. ([#3521](https://github.com/terraform-providers/terraform-provider-google/pull/3521))
* pubsub: `google_pubsub_subscription` now supports setting an `expiration_policy`. ([#1703](https://github.com/GoogleCloudPlatform/magic-modules/pull/1703))

BUG FIXES:
* bigquery: `google_bigquery_table` will work with a larger range of projects id formats. ([#3486](https://github.com/terraform-providers/terraform-provider-google/issues/3486))
* cloudfunctions: `google_cloudfunctions_fucntion` no longer restricts an outdated list of `region`s ([#3530](https://github.com/terraform-providers/terraform-provider-google/issues/3530))
* compute: `google_compute_instance` now retries updating metadata when fingerprints are mismatched. ([#3372](https://github.com/terraform-providers/terraform-provider-google/issues/3372))
* compute: `google_compute_subnetwork.secondary_ip_ranges` doesn't cause a diff on out of band changes, allows updating to empty list of ranges. ([#3496](https://github.com/terraform-providers/terraform-provider-google/issues/3496))
* container: `google_container_cluster` setting networks / subnetworks by name works with `location`. ([#3492](https://github.com/terraform-providers/terraform-provider-google/issues/3492))
* container: `google_container_cluster` removed an overly restrictive validation restricting `node_pool` and `remove_default_node_pool` being specified at the same time. ([#3497](https://github.com/terraform-providers/terraform-provider-google/issues/3497))
* storage: `data.google_storage_bucket_object` now correctly URL encodes the slashes in a file name ([#1613](https://github.com/terraform-providers/terraform-provider-google/issues/1613))

## 2.5.1 (April 22, 2019)

BUG FIXES:
* compute: `google_compute_backend_service` handles empty/nil `iap` block created by previous providers properly. ([#3459](https://github.com/terraform-providers/terraform-provider-google/issues/3459))
* compute: `google_compute_backend_service` allows multiple instance types in `backends.group` again. ([#3463](https://github.com/terraform-providers/terraform-provider-google/issues/3463))
* dns: `google_dns_managed_zone` does not permadiff when visiblity is set to default and returned as empty from API ([#3459](https://github.com/terraform-providers/terraform-provider-google/issues/3461))
* google_projects: Datasource `google_projects` now handles paginated results from listing projects ([#3464](https://github.com/terraform-providers/terraform-provider-google/pull/3464))
* google_project_iam: `google_project_iam_policy/member/binding` now attempts to retry for read-only operations as well as retrying read-write operations ([#3455](https://github.com/terraform-providers/terraform-provider-google/pull/3455))
* kms: `google_kms_crypto_key.rotation_period` now can be an empty string to allow for unset behavior in modules ([#3468](https://github.com/terraform-providers/terraform-provider-google/pull/3468))

## 2.5.0 (April 18, 2019)

KNOWN ISSUES:
* compute: `google_compute_subnetwork` will fail to reorder `secondary_ip_range` values at apply time
* compute: `google_compute_subnetwork`s used with a VPC-native GKE cluster will have a diff if that cluster creates secondary ranges automatically.

BACKWARDS INCOMPATIBILITIES:
* all: This is the first release to use the 0.12 SDK required for Terraform 0.12 support. Some provider behaviour may have changed as a result of changes made by the new SDK version.
* compute: `google_compute_instance_group` will not reconcile instances recreated within the same `terraform apply` due to underlying `0.12` SDK changes in the provider. ([#616](https://github.com/terraform-providers/terraform-provider-google/issues/616))
* compute: `google_compute_subnetwork` will have a diff if `secondary_ip_range` values defined in config don't exactly match real state; if so, they will need to be reconciled. ([#3432](https://github.com/terraform-providers/terraform-provider-google/issues/3432))
* container: `google_container_cluster` will have a diff if `master_authorized_networks.cidr_blocks` defined in config doesn't exactly match the real state; if so, it will need to be reconciled. ([#3427](https://github.com/terraform-providers/terraform-provider-google/issues/3427))


BUG FIXES:
* container: `google_container_cluster` catch out of band changes to `master_authorized_networks.cidr_blocks`. ([#3427](https://github.com/terraform-providers/terraform-provider-google/issues/3427))

## 2.4.1 (April 30, 2019)

NOTES:
This 2.4.1 release is a bugfix release for 2.4.0. It backports the fixes applied in the 2.5.1 release to the 2.4.0 series.

BUG FIXES:
* compute: `google_compute_backend_service` handles empty/nil `iap` block created by previous providers properly. ([#3459](https://github.com/terraform-providers/terraform-provider-google/issues/3459))
* compute: `google_compute_backend_service` allows multiple instance types in `backends.group` again. ([#3463](https://github.com/terraform-providers/terraform-provider-google/issues/3463))
* dns: `google_dns_managed_zone` does not permadiff when visiblity is set to default and returned as empty from API ([#3459](https://github.com/terraform-providers/terraform-provider-google/issues/3461))

## 2.4.0 (April 15, 2019)

KNOWN ISSUES:

* compute: `google_compute_backend_service` resources created with past provider versions won't work with `2.4.0`. You can pin your provider version or manually delete them and recreate them until this is resolved. (https://github.com/terraform-providers/terraform-provider-google/issues/3441)
* dns: `google_dns_managed_zone.visibility` will cause a diff if set to `public`. Setting it to `""` (defaulting to public) will work around this. (https://github.com/terraform-providers/terraform-provider-google/issues/3435)

FEATURES:
* **New Resource**: `google_access_context_manager_access_policy` is now available at GA. ([#3358](https://github.com/terraform-providers/terraform-provider-google/issues/3358))
* **New Resource**: `google_access_context_manager_access_level` is now available at GA. ([#3358](https://github.com/terraform-providers/terraform-provider-google/issues/3358))
* **New Resource**: `google_access_context_manager_service_perimeter` is now available at GA. ([#3358](https://github.com/terraform-providers/terraform-provider-google/issues/3358))
* **New Datasource**: `google_service_account_access_token` is now available. ([#3357](https://github.com/terraform-providers/terraform-provider-google/issues/3357))

ENHANCEMENTS:
* compute: `google_compute_backend_service` is now generated with Magic Modules, adding configurable timeouts, multiple import formats, `creation_timestamp` output. ([#3345](https://github.com/terraform-providers/terraform-provider-google/issues/3345))
* compute: `google_compute_backend_service` now supports `load_balancing_scheme` and `cdn_policy.signed_url_cache_max_age_sec`. ([#3375](https://github.com/terraform-providers/terraform-provider-google/issues/3375))
* compute: `google_compute_network` now supports `delete_default_routes_on_create` to delete pre-created routes at network creation time. ([#3391](https://github.com/terraform-providers/terraform-provider-google/issues/3391))
* dns: `google_dns_managed_zone.private_visibility_config`, part of private DNS, is now generally available. ([#3352](https://github.com/terraform-providers/terraform-provider-google/issues/3352))

BUG FIXES:
* container: `google_container_cluster` will ignore out of band changes on `node_ipv4_cidr_block`. ([#3319](https://github.com/terraform-providers/terraform-provider-google/issues/3319))
* container: `google_container_cluster` will now reject config with both `node_pool` and `remove_default_node_pool` defined ([#3422](https://github.com/terraform-providers/terraform-provider-google/issues/3422))
* container: `google_container_cluster` will allow >20 `cidr_blocks` in `master_authorized_networks_config`. ([#3397](https://github.com/terraform-providers/terraform-provider-google/issues/3397))
* netblock: `data.google_netblock_ip_ranges.cidr_blocks` will better handle ipv6 input. ([#3390](https://github.com/terraform-providers/terraform-provider-google/issues/3390))
* sql: `google_sql_database_instance` will retry reads during Terraform refreshes if it hits a rate limit. ([#3366](https://github.com/terraform-providers/terraform-provider-google/issues/3366))

## 2.3.0 (March 26, 2019)

DEPRECATIONS:
* container: `google_container_cluster` `zone` and `region` fields are deprecated in favour of `location`, `additional_zones` in favour of `node_locations`. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* container: `google_container_node_pool` `zone` and `region` fields are deprecated in favour of `location`. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* container: `data.google_container_cluster` `zone` and `region` fields are deprecated in favour of `location`. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* container: `google_container_engine_versions` `zone` and `region` fields are deprecated in favour of `location`. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))

FEATURES:
* **New Datasource**: `google_*_organization_policy` Adding datasources for folder and project org policy ([#3137](https://github.com/terraform-providers/terraform-provider-google/issues/3137))

ENHANCEMENTS:
* compute: `google_compute_disk`, `google_compute_region_disk` now support `physical_block_size_bytes` ([#526](https://github.com/terraform-providers/terraform-provider-google/issues/526))
* compute: `google_compute_forwarding_rule` supports specifying `all_ports` for internal load balancing. ([#3309](https://github.com/terraform-providers/terraform-provider-google/issues/3309))
* compute: `google_compute_vpn_tunnel` will properly apply labels. ([#3277](https://github.com/terraform-providers/terraform-provider-google/issues/3277))
* container: `google_container_cluster` adds a unified `location` field for regions and zones, `node_locations` to manage extra zones for multi-zonal clusters and specific zones for regional clusters. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* container: `google_container_node_pool` adds a unified `location` field for regions and zones. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* container: `data.google_container_cluster` adds a unified `location` field for regions and zones. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* container: `google_container_engine_versions` adds a unified `location` field for regions and zones. ([#3114](https://github.com/terraform-providers/terraform-provider-google/issues/3114))
* dataflow: `google_dataflow_job` has support for custom service accounts with `service_account_email`. ([#3238](https://github.com/terraform-providers/terraform-provider-google/issues/3238))
* monitoring: `google_monitoring_uptime_check_config` Add a computed field for uptime check id ([#3138](https://github.com/terraform-providers/terraform-provider-google/issues/3138))
* resourcemanager: `google_*_organization_policy` Add import support for folder and project organization_policies ([#3218](https://github.com/terraform-providers/terraform-provider-google/issues/3218))
* sql: `google_sql_ssl_cert` Allow project to be specified at resource level ([#3235](https://github.com/terraform-providers/terraform-provider-google/issues/3235))
* storage: `google_storage_bucket` Change storage bucket import logic to avoid calls to compute api ([#3244](https://github.com/terraform-providers/terraform-provider-google/issues/3244))
* storage: `google_storage_bucket.storage_class` supports updating. ([#3297](https://github.com/terraform-providers/terraform-provider-google/issues/3297))
* various: Some import formats that previously failed will now work as documented. ([#3283](https://github.com/terraform-providers/terraform-provider-google/issues/3283))

BUG FIXES:
* compute: `google_compute_disk` will properly detach instances again. ([#3269](https://github.com/terraform-providers/terraform-provider-google/issues/3269))
* container: `google_container_cluster`, `google_container_node_pool` properly suppress new GKE `1.12` `metadata` values. ([#3233](https://github.com/terraform-providers/terraform-provider-google/issues/3233))
* container: `google_container_cluster` properly collects service-level errors from the API ([#2941](https://github.com/terraform-providers/terraform-provider-google/issues/2941))
* monitoring: `google_monitoring_uptime_check_config` Change all fields for monitored resource to force recreation ([#3132](https://github.com/terraform-providers/terraform-provider-google/issues/3132))
* various: Retry only 409 concurrent operation errors and not naming conflicts ([#3285](https://github.com/terraform-providers/terraform-provider-google/issues/3285))

## 2.2.0 (March 12, 2019)

KNOWN ISSUES:

* compute: `google_compute_disk` is unable to detach instances at deletion time.

---

FEATURES:
* **New Datasource**: `data.google_projects` for retrieving a list of projects based on a filter. ([#3178](https://github.com/terraform-providers/terraform-provider-google/issues/3178))
* **New Resource**: `google_tpu_node` for Cloud TPU Nodes ([#3179](https://github.com/terraform-providers/terraform-provider-google/issues/3179))

ENHANCEMENTS:
* compute: `google_compute_disk` and `google_compute_region_disk` will now detach themselves from a more up to date set of users at delete time. ([#3154](https://github.com/terraform-providers/terraform-provider-google/issues/3154))
* compute: `google_compute_network` is now generated by Magic Modules, supporting configurable timeouts and more import formats. ([#3203](https://github.com/terraform-providers/terraform-provider-google/issues/3203))
* compute: `google_compute_firewall` will validate the maximum size of service account lists at plan time. ([#3201](https://github.com/terraform-providers/terraform-provider-google/issues/3201))
* container: `google_container_cluster` can now disable VPC Native clusters with `ip_allocation_policy.use_ip_aliases` ([#3174](https://github.com/terraform-providers/terraform-provider-google/issues/3174))
* container: `data.google_container_engine_versions` supports `version_prefix` to allow fuzzy version matching. Using this field, Terraform can match the latest version of a major, minor, or patch release. ([#3199](https://github.com/terraform-providers/terraform-provider-google/issues/3199))
* pubsub: `google_pubsub_subscription` now supports configuring `message_retention_duration` and `retain_acked_messages`. ([#3193](https://github.com/terraform-providers/terraform-provider-google/issues/3193))

BUG FIXES:
* app_engine: `google_app_engine_application` correctly outputs `gcr_domain`.  ([#3149](https://github.com/terraform-providers/terraform-provider-google/issues/3149))
* compute: `data.google_compute_subnetwork` outputs the `self_link` field again. ([#3156](https://github.com/terraform-providers/terraform-provider-google/issues/3156))
* compute: `google_compute_attached_disk` is now removed from state if the instance was removed. ([#3183](https://github.com/terraform-providers/terraform-provider-google/issues/3183))
* container: `google_container_cluster` private_cluster_config now has a diff suppress to prevent a permadiff for and allows for empty `master_ipv4_cidr_block`  ([#460](https://github.com/terraform-providers/terraform-provider-google/issues/460))
* container: `google_container_cluster` import behavior fixed/documented for TF-state-only fields (`remove_default_node_pool`, `min_master_version`) ([#3146](https://github.com/terraform-providers/terraform-provider-google/issues/3146)][[#3169](https://github.com/terraform-providers/terraform-provider-google/issues/3169)][[#3180](https://github.com/terraform-providers/terraform-provider-google/issues/3180))
* storagetransfer: `google_storage_transfer_job` will no longer crash when accessing nil dates. ([#3185](https://github.com/terraform-providers/terraform-provider-google/issues/3185))

## 2.1.0 (February 26, 2019)

FEATURES:
* **New Datasource**: `google_client_openid_userinfo` for retrieving the `email` used to authenticate with GCP. ([#3103](https://github.com/terraform-providers/terraform-provider-google/issues/3103))

ENHANCEMENTS:
* compute: `data.google_compute_subnetwork` can now be addressed by `self_link` as an alternative to the existing `name`/`region`/`project` fields. ([#3040](https://github.com/terraform-providers/terraform-provider-google/issues/3040))
* pubsub: `google_pubsub_topic` is now generated using Magic Modules, adding Open in Cloud Shell examples, configurable timeouts, and the `labels` field. ([#3043](https://github.com/terraform-providers/terraform-provider-google/issues/3043))
* pubsub: `google_pubsub_subscription` is now generated using Magic Modules, adding Open in Cloud Shell examples, configurable timeouts, update support, and the `labels` field. ([#3043](https://github.com/terraform-providers/terraform-provider-google/issues/3043))
* sql: `google_sql_database_instance` now provides `public_ip_address` and `private_ip_address` outputs of the first public and private IP of the instance respectively. ([#3091](https://github.com/terraform-providers/terraform-provider-google/issues/3091))


BUG FIXES:
* sql: `google_sql_database_instance` allows the empty string to be set for `private_network`. ([#3091](https://github.com/terraform-providers/terraform-provider-google/issues/3091))

## 2.0.0 (February 12, 2019)

BACKWARDS INCOMPATIBILITIES:
* bigtable: `google_bigtable_instance.cluster.num_nodes` will fail at plan time if `DEVELOPMENT` instances have `num_nodes = "0"` set explicitly. If it has been set, unset the field. ([#2401](https://github.com/terraform-providers/terraform-provider-google/issues/2401))
* cloudbuild: `google_cloudbuild_trigger.build.step.args` is now a list instead of space separated strings. ([#2790](https://github.com/terraform-providers/terraform-provider-google/issues/2790))
* cloudfunctions: `google_cloudfunctions_function.retry_on_failure` has been removed. Use `event_trigger.failure_policy.retry` instead. ([#2392](https://github.com/terraform-providers/terraform-provider-google/issues/2392))
* composer: `google_composer_environment.node_config.zone` is now `Required`. ([#2967](https://github.com/terraform-providers/terraform-provider-google/issues/2967))
* compute: `google_compute_instance`, `google_compute_instance_from_template` `metadata` field is now authoritative and will remove values not explicitly set in config. ([#2208](https://github.com/terraform-providers/terraform-provider-google/issues/2208))
* compute: `google_compute_project_metadata` resource is now authoritative and will remove values not explicitly set in config. ([#2205](https://github.com/terraform-providers/terraform-provider-google/issues/2205))
* compute: `google_compute_url_map` resource is now authoritative and will remove values not explicitly set in config. ([#2245](https://github.com/terraform-providers/terraform-provider-google/issues/2245))
* compute: `google_compute_global_forwarding_rule.labels` is removed from the `google` provider and must be used in the `google-beta` provider. ([#2399](https://github.com/terraform-providers/terraform-provider-google/issues/2399))
* compute: `google_compute_subnetwork_iam_binding`, `google_compute_subnetwork_iam_member`, `google_compute_subnetwork_iam_policy` are removed from the `google` provider and must be used in the `google-beta` provider. ([#2398](https://github.com/terraform-providers/terraform-provider-google/issues/2398))
* compute: `google_compute_backend_service.custom_request_headers` is removed from the `google` provider and must be used in the `google-beta` provider. ([#2405](https://github.com/terraform-providers/terraform-provider-google/issues/2405))
* compute: `google_compute_snapshot.snapshot_encryption_key_raw`, `google_compute_snapshot.snapshot_encryption_key_sha256`, `google_compute_snapshot.source_disk_encryption_key_raw`, `google_compute_snapshot.source_disk_encryption_key_sha256` fields are now removed. Use `google_compute_snapshot.snapshot_encryption_key.0.raw_key`, `google_compute_snapshot.snapshot_encryption_key.0.sha256`, `google_compute_snapshot.source_disk_encryption_key.0.raw_key`, `google_compute_snapshot.source_disk_encryption_key.0.sha256` instead. ([#2572](https://github.com/terraform-providers/terraform-provider-google/issues/2572)][[#2624](https://github.com/terraform-providers/terraform-provider-google/issues/2624))
* container: `google_container_node_pool.max_pods_per_node` is removed from the `google` provider and must be used in the `google-beta` provider. ([#2391](https://github.com/terraform-providers/terraform-provider-google/issues/2391))
* compute: `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` have had their `version`, `auto_healing_policies`, and `rolling_update_policy` fields removed from the `google` provider. They must be used in the `google-beta` provider. `rolling_update_policy` was renamed `update_policy` in that provider. ([#2392](https://github.com/terraform-providers/terraform-provider-google/issues/2392))
* compute: `google_compute_instance_group_manager` is no longer imported by the provider-level region. Set the appropriate provider-level zone instead. ([#2693](https://github.com/terraform-providers/terraform-provider-google/issues/2693))
* compute: `google_compute_region_instance_group_manager.update_strategy` in the `google-beta` provider has been removed. ([#2594](https://github.com/terraform-providers/terraform-provider-google/issues/2594))
* compute: `google_compute_instance`, `google_compute_instance_template`, `google_compute_instance_from_template` have had the `network_interface.address` field removed. ([#2595](https://github.com/terraform-providers/terraform-provider-google/issues/2595))
* compute: `google_compute_disk` is no longer imported by the provider-level region. Set the appropriate provider-level zone instead. ([#2694](https://github.com/terraform-providers/terraform-provider-google/issues/2694))
* compute: `google_compute_router_nat.subnetwork.source_ip_ranges_to_nat` is now Required inside `subnetwork` blocks. ([#2749](https://github.com/terraform-providers/terraform-provider-google/issues/2749))
* compute: `google_compute_ssl_certificate`'s `private_key` field is no longer stored in state in cleartext; it is now SHA256 encoded. ([#2976](https://github.com/terraform-providers/terraform-provider-google/issues/2976))
* container: `google_container_cluster` fields (`private_cluster`, `master_ipv4_cidr_block`) are removed. Use `private_cluster_config` and `private_cluster_config.master_ipv4_cidr_block` instead. ([#2395](https://github.com/terraform-providers/terraform-provider-google/issues/2395))
* container: `google_container_cluster` fields (`enable_binary_authorization`, `enable_tpu`, `pod_security_policy_config`) are removed from the `google` provider and must be used in the `google-beta` provider. ([#2395](https://github.com/terraform-providers/terraform-provider-google/issues/2395))
* container: `google_container_cluster.node_config` fields (`taints`, `workload_metadata_config`) are removed from the `google` provider and must be used in the `google-beta` provider. ([#2601](https://github.com/terraform-providers/terraform-provider-google/issues/2601))
* container: `google_container_node_pool.node_config` fields (`taints`, `workload_metadata_config`) are removed from the `google` provider and must be used in the `google-beta` provider. ([#2601](https://github.com/terraform-providers/terraform-provider-google/issues/2601))
* container: `google_container_node_pool`'s `name_prefix` field has been restored and is no longer deprecated. ([#2975](https://github.com/terraform-providers/terraform-provider-google/issues/2975))
* sql: `google_sql_database_instance` resource is now authoritative and will remove values not explicitly set in config. ([#2203](https://github.com/terraform-providers/terraform-provider-google/issues/2203))
* bigtable: `google_bigtable_instance` `zone` field is no longer inferred from the provider.
* endpoints: `google_endpoints_service.protoc_output` was removed. Use `google_endpoints_service.protoc_output_base64` instead. ([#2396](https://github.com/terraform-providers/terraform-provider-google/issues/2396))
* resourcemanager: `google_project_iam_policy` is now authoritative and will remove values not explicitly set in config. Several fields were removed that made it authoritative: `authoritative`, `restore_policy`, and `disable_project`. This resource is very dangerous! Ensure you are not using the removed fields (`authoritative`, `restore_policy`, `disable_project`). ([#2315](https://github.com/terraform-providers/terraform-provider-google/issues/2315))
* resourcemanager: Datasource `google_service_account_key.service_account_id` has been removed. Use the `name` field instead. ([#2397](https://github.com/terraform-providers/terraform-provider-google/issues/2397))
* resourcemanager: `google_project.app_engine` has been removed. Use the `google_app_engine_application` resource instead. ([#2386](https://github.com/terraform-providers/terraform-provider-google/issues/2386))
* resourcemanager: `google_organization_custom_role.deleted` is now an output-only attribute. Use `terraform destroy`, or remove the resource from your config instead. ([#2596](https://github.com/terraform-providers/terraform-provider-google/issues/2596))
* resourcemanager: `google_project_custom_role.deleted` is now an output-only attribute. Use `terraform destroy`, or remove the resource from your config instead. ([#2619](https://github.com/terraform-providers/terraform-provider-google/issues/2619))
* serviceusage: `google_project_service` will now error instead of silently disabling dependent services if `disable_dependent_services` is unset. ([#2938](https://github.com/terraform-providers/terraform-provider-google/issues/2938))
* storage: `google_storage_object_acl.role_entity` is now authoritative and will remove values not explicitly set in config. Use `google_storage_object_access_control` for fine-grained management. ([#2316](https://github.com/terraform-providers/terraform-provider-google/issues/2316))
* storage: `google_storage_default_object_acl.role_entity` is now authoritative and will remove values not explicitly set in config. ([#2345](https://github.com/terraform-providers/terraform-provider-google/issues/2345))
* iam: `google_*_iam_binding` Change all IAM bindings to be authoritative ([#2764](https://github.com/terraform-providers/terraform-provider-google/issues/2764))

FEATURES:
* **New Resource**: `google_access_context_manager_access_policy` for managing the container for an organization's access levels. ([`google-beta`#96](https://github.com/terraform-providers/terraform-provider-google-beta/pull/96))
* **New Resource**: `google_access_context_manager_access_level` for managing an organization's access levels. ([`google-beta`#149](https://github.com/terraform-providers/terraform-provider-google-beta/pull/149))
* **New Resource**: `google_access_context_manager_service_perimeter` for managing service perimeters in an access policy. ([`google-beta`#246](https://github.com/terraform-providers/terraform-provider-google-beta/pull/246))
* **New Resource**: `google_storage_transfer_job` for managing recurring storage transfers with Google Cloud Storage. ([#2707](https://github.com/terraform-providers/terraform-provider-google/issues/2707))
* **New Datasource**: `google_storage_transfer_project_service_account` data source for retrieving the Storage Transfer service account for a project ([#2692](https://github.com/terraform-providers/terraform-provider-google/issues/2692))
* **New Resource**: `google_app_engine_firewall_rule` ([#2738](https://github.com/terraform-providers/terraform-provider-google/issues/2738)][[#2849](https://github.com/terraform-providers/terraform-provider-google/issues/2849))
* **New Resource**: `google_project_iam_audit_config` ([#2731](https://github.com/terraform-providers/terraform-provider-google/issues/2731))
* **New Datasource**: `google_kms_crypto_key` data source for an externally managed KMS crypto key ([#2891](https://github.com/terraform-providers/terraform-provider-google/issues/2891))
* **New Datasource**: `google_kms_key_ring` ([#2891](https://github.com/terraform-providers/terraform-provider-google/issues/2891))

ENHANCEMENTS:
* provider: Add `access_token` config option to allow Terraform to authenticate using short-lived Google OAuth 2.0 access token ([#2838](https://github.com/terraform-providers/terraform-provider-google/issues/2838))
* bigquery: Add `default_partition_expiration_ms` field to `google_bigquery_dataset` resource. ([#2287](https://github.com/terraform-providers/terraform-provider-google/issues/2287))
* bigquery: Add `delete_contents_on_destroy` field to `google_bigquery_dataset` resource. ([#2986](https://github.com/terraform-providers/terraform-provider-google/issues/2986))
* bigquery: Add `time_partitioning.require_partition_filter` to `google_bigquery_table` resource. ([#2815](https://github.com/terraform-providers/terraform-provider-google/issues/2815))
* bigquery: Allow more BigQuery regions ([#2566](https://github.com/terraform-providers/terraform-provider-google/issues/2566))
* bigtable: Add `column_family` at create time to `google_bigtable_table`. ([#2228](https://github.com/terraform-providers/terraform-provider-google/issues/2228))
* bigtable: Add multi-zone (inside one region) replication to `google_bigtable_instance`. ([#2313](https://github.com/terraform-providers/terraform-provider-google/issues/2313)] [[#2289](https://github.com/terraform-providers/terraform-provider-google/issues/2289))
* cloudbuild: `google_cloudbuild_trigger` is now autogenerated, adding more configurable timeouts, import support, and the `disabled` field. `ignored_files`, `included_files` are now updatable. ([#2790](https://github.com/terraform-providers/terraform-provider-google/issues/2790)] [[#2871](https://github.com/terraform-providers/terraform-provider-google/issues/2871))
* cloudfunctions: ` google_cloudfunctions_function` now has souce repo support ([#2650](https://github.com/terraform-providers/terraform-provider-google/issues/2650))
* cloudfunctions: `google_cloudfunctions_function` now supports `service_account_email` for self-provided service accounts. ([#2947](https://github.com/terraform-providers/terraform-provider-google/issues/2947))
* compute: `google_compute_forwarding_rule` supports specifying `all_ports` for internal load balancing. ([`google-beta`#297](https://github.com/terraform-providers/terraform-provider-google-beta/pull/297))
* compute: `google_compute_image` is now autogenerated and supports multiple import formats, and `size_gb` attribute. ([#2769](https://github.com/terraform-providers/terraform-provider-google/issues/2769))
* compute: `google_compute_url_map` resource is now autogenerated and supports multiple import formats. ([#2245](https://github.com/terraform-providers/terraform-provider-google/issues/2245))
* compute: Add `name`, `unique_id`, and `display_name` properties to `data.google_compute_default_service_account` ([#2778](https://github.com/terraform-providers/terraform-provider-google/issues/2778))
* compute: `google_compute_disk` Add support for KMS encryption to compute disk ([#2884](https://github.com/terraform-providers/terraform-provider-google/issues/2884))
* compute: Add support for PARTNER interconnects. ([#2959](https://github.com/terraform-providers/terraform-provider-google/issues/2959))
* dataproc: Add `accelerators` support to `google_dataproc_cluster` to allow using GPU accelerators. ([#2411](https://github.com/terraform-providers/terraform-provider-google/issues/2411))
* dataproc: `google_dataproc_cluster` Add support for KMS encryption to dataproc cluster ([#2840](https://github.com/terraform-providers/terraform-provider-google/issues/2840))
* project: The google_iam_policy data source now supports Audit Configs ([#2687](https://github.com/terraform-providers/terraform-provider-google/issues/2687))
* kms: Add support for `protection_level` to `google_kms_crypto_key` ([#2751](https://github.com/terraform-providers/terraform-provider-google/issues/2751))
* resourcemanager: add `inherit_from_parent` to all org policy resources ([#2653](https://github.com/terraform-providers/terraform-provider-google/issues/2653))
* serviceusage: `google_project_service` now supports `disable_dependent_services` to control whether services can disable services that depend on them at disable-time. ([#2938](https://github.com/terraform-providers/terraform-provider-google/issues/2938))
* sourcerepo: `google_sourcerepo_repository` is now autogenerated, adding configurable timeouts. ([#2797](https://github.com/terraform-providers/terraform-provider-google/issues/2797))
* storage: `google_storage_object_acl` can more easily swap between `role_entity` and `predefined_acl` ACL definitions. ([#2316](https://github.com/terraform-providers/terraform-provider-google/issues/2316))
* storage: `google_storage_bucket` has support for `requester_pays` ([#2580](https://github.com/terraform-providers/terraform-provider-google/issues/2580))
* storage: `google_storage_bucket_object` exports `output_name` for interpolations on `name`, allowing you to trigger reapplication of `google_storage_object_acl` on recreated objects. ([#2914](https://github.com/terraform-providers/terraform-provider-google/issues/2914))
* storage: During a force destroy, `google_storage_bucket` will delete objects in parallel instead of serially. ([#2944](https://github.com/terraform-providers/terraform-provider-google/issues/2944))
* spanner: `google_spanner_database` is autogenerated and supports timeouts. ([#2812](https://github.com/terraform-providers/terraform-provider-google/issues/2812))
* spanner: `google_spanner_instance` is autogenerated and supports timeouts. ([#2892](https://github.com/terraform-providers/terraform-provider-google/issues/2892))

BUG FIXES:

* cloudbuild: allow `google_cloudbuild_trigger.trigger_template.project` to not be set ([#2655](https://github.com/terraform-providers/terraform-provider-google/issues/2655))
* cloudbuild: fix update so it doesn't error every time ([#2743](https://github.com/terraform-providers/terraform-provider-google/issues/2743))
* cloudfunctions: No longer over-validate project ids in `google_cloudfunctions_function` ([#2780](https://github.com/terraform-providers/terraform-provider-google/issues/2780))
* compute: attached_disk now supports region disks ([#2441](https://github.com/terraform-providers/terraform-provider-google/issues/2441))
* compute: extract vpn tunnel region/project from vpn gateway ([#2640](https://github.com/terraform-providers/terraform-provider-google/issues/2640))
* compute: send instance scheduling block with automaticrestart true if there is none in cfg ([#2638](https://github.com/terraform-providers/terraform-provider-google/issues/2638))
* compute: fix disk behaivor in compute_instance_from_template ([#2695](https://github.com/terraform-providers/terraform-provider-google/issues/2695))
* compute: add diffsuppress for region_autoscaler.target so it can be used with both versions of the provider ([#2770](https://github.com/terraform-providers/terraform-provider-google/issues/2770))
* compute: fix ID for inferring project for old compute_project_metadata states ([#2844](https://github.com/terraform-providers/terraform-provider-google/issues/2844))
* compute: `google_compute_backend_service` will send the correct `iap` block values during updates ([#2978](https://github.com/terraform-providers/terraform-provider-google/issues/2978))
* container: fix failure when updating node versions ([#2872](https://github.com/terraform-providers/terraform-provider-google/issues/2872))
* dataproc: convert dataproc_cluster.cluster_config.gce_cluster_config.tags into a set ([#2633](https://github.com/terraform-providers/terraform-provider-google/issues/2633))
* iam: fix permadiff when stage is ALPHA ([#2370](https://github.com/terraform-providers/terraform-provider-google/issues/2370))
* iam: add another retry if iam read returns nil ([#2629](https://github.com/terraform-providers/terraform-provider-google/issues/2629))
* monitoring: `uptime_check_config` can now be updated and won't error when changing duration. ([#2786](https://github.com/terraform-providers/terraform-provider-google/issues/2786))
* runtimeconfig: allow more characters in runtimeconfig name ([#2643](https://github.com/terraform-providers/terraform-provider-google/issues/2643))
* sql: send maintenance_window.hour even if it's zero, since that's a valid value ([#2630](https://github.com/terraform-providers/terraform-provider-google/issues/2630))
* sql: allow cross-project imports for sql user ([#2632](https://github.com/terraform-providers/terraform-provider-google/issues/2632))
* sql: mark region as computed in sql db instance since we use getregion ([#2635](https://github.com/terraform-providers/terraform-provider-google/issues/2635))
* sql: `google_sql_database_instance` Stop SQL instances from reporting failing to destroy ([#2811](https://github.com/terraform-providers/terraform-provider-google/issues/2811))

## 1.20.0 (December 14, 2018)

DEPRECATIONS:
* **Deprecated `google_compute_snapshot`'s top-level encryption fields.** ([#2572](https://github.com/terraform-providers/terraform-provider-google/issues/2572))

FEATURES:
* **New Resource**: `google_storage_object_access_control` for fine-grained management of ACLs on Google Cloud Storage objects ([#2256](https://github.com/terraform-providers/terraform-provider-google/issues/2256))
* **New Resource**: `google_storage_default_object_access_control` for fine-grained management of default object ACLs on Google Cloud Storage buckets ([#2358](https://github.com/terraform-providers/terraform-provider-google/issues/2358))
* **New Resource**: `google_sql_ssl_cert` for Google Cloud SQL client SSL certificates. ([#2290](https://github.com/terraform-providers/terraform-provider-google/issues/2290))
* **New Resource**: `google_monitoring_notification_channel` ([#2452](https://github.com/terraform-providers/terraform-provider-google/issues/2452))
* **New Resource**: `google_compute_router_nat` ([#2576](https://github.com/terraform-providers/terraform-provider-google/issues/2576))
* **New Resource**: `google_monitoring_group` ([#2451](https://github.com/terraform-providers/terraform-provider-google/issues/2451))
* **New Resource**: `google_billing_account_iam_binding`, `google_billing_account_iam_member`, `google_billing_account_iam_policy` for managing Billing Account IAM policies, including managing Billing Account users. ([#2143](https://github.com/terraform-providers/terraform-provider-google/issues/2143))
* **New Datasource**: `google_iam_role` datasource to be able to read an IAM role's permissions. ([#2482](https://github.com/terraform-providers/terraform-provider-google/issues/2482))

ENHANCEMENTS:
* cloudbuild: Added Update support for `google_cloudbuild_trigger`.  ([#2121](https://github.com/terraform-providers/terraform-provider-google/issues/2121))
* cloudfunctions: Add `runtime` support to `google_cloudfunctions_function` ([#2340](https://github.com/terraform-providers/terraform-provider-google/issues/2340))
* cloudfunctions: Add new-style Storage and Pub/Sub trigger support to `google_cloudfunctions_function` ([#2412](https://github.com/terraform-providers/terraform-provider-google/issues/2412))
* compute: `google_compute_health_check` supports for content-based load balancing (`response` field) in HTTP(S) checks. ([#2550](https://github.com/terraform-providers/terraform-provider-google/issues/2550))
* container: regional and private clusters are in GA now ([#2364](https://github.com/terraform-providers/terraform-provider-google/issues/2364))
* iam: `google_service_accounts` now supports multiple import formats. ([#2261](https://github.com/terraform-providers/terraform-provider-google/issues/2261))
* sql: add support for private IP for SQL instances. ([#2662](https://github.com/terraform-providers/terraform-provider-google/issues/2662))

BUG FIXES:
* bigquery: added australia and europe regions to the validate function ([#2333](https://github.com/terraform-providers/terraform-provider-google/issues/2333))
* compute: `google_compute_disk.snapshot`, `google_compute_region_disk.snapshot` properly allow partial URIs. ([#2450](https://github.com/terraform-providers/terraform-provider-google/issues/2450))
* compute: The `google_compute_instance` datasource can now be addressed by `self_link`. ([#2874](https://github.com/terraform-providers/terraform-provider-google/issues/2874))
* compute: `google_compute_image.licenses` elements properly allow partial URIs / versioned self links. ([#3018](https://github.com/terraform-providers/terraform-provider-google/issues/3018))
* compute: `google_compute_project_metadata` can now be imported from a project other than the one specified in your config. ([#3018](https://github.com/terraform-providers/terraform-provider-google/issues/3018))
* pubsub: fix issue where not all attributes were saved in state ([#2469](https://github.com/terraform-providers/terraform-provider-google/issues/2469))


## 1.19.1 (October 12, 2018)

BUG FIXES:

* all: fix deprecation links in resources ([#2197](https://github.com/terraform-providers/terraform-provider-google/issues/2197)] [[#2196](https://github.com/terraform-providers/terraform-provider-google/issues/2196))
* all: fix panics caused by including empty blocks with lists ([#2229](https://github.com/terraform-providers/terraform-provider-google/issues/2229)] [[#2233](https://github.com/terraform-providers/terraform-provider-google/issues/2233)] [[#2239](https://github.com/terraform-providers/terraform-provider-google/issues/2239))
* compute: allow instance templates to have disks with no source image set ([#2218](https://github.com/terraform-providers/terraform-provider-google/issues/2218))
* project: fix plan output when app engine api is not enabled ([#2204](https://github.com/terraform-providers/terraform-provider-google/issues/2204))

## 1.19.0 (October 08, 2018)

BACKWARDS INCOMPATIBILITIES:
* all: beta fields have been deprecated in favor of the new `google-beta` provider. See https://terraform.io/docs/providers/google/provider_versions.html for more info. ([#2152](https://github.com/terraform-providers/terraform-provider-google/issues/2152)] [[#2142](https://github.com/terraform-providers/terraform-provider-google/issues/2142))
* bigtable: `google_bigtable_instance` deprecated the `cluster_id`, `zone`, `num_nodes`, and `storage_type` fields, creating a `cluster` block containing those fields instead. ([#2161](https://github.com/terraform-providers/terraform-provider-google/issues/2161))
* cloudfunctions: `google_cloudfunctions_function` and `datasource_google_cloudfunctions_function` deprecated `trigger_bucket` and `trigger_topic` in favor of the new `event_trigger` field, and deprecated `retry_on_failure` in favor of the `event_trigger.failure_policy.retry` field. ([#2158](https://github.com/terraform-providers/terraform-provider-google/issues/2158))
* compute: `google_compute_instance`, `google_compute_instance_template`, `google_compute_instance_from_template` have had the `network_interface.address` field deprecated and the `network_interface.network_ip` field undeprecated to better match the API. Terraform configurations should migrate from `network_interface.address` to `network_interface.network_ip`. ([#2096](https://github.com/terraform-providers/terraform-provider-google/issues/2096))
* compute: `google_compute_instance`, `google_compute_instance_from_template` have had the `network_interface.0.access_config.0.assigned_nat_ip` field deprecated. Please use `network_interface.0.access_config.0.nat_ip` instead.
* compute: `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` have had their `version`, `auto_healing_policies`, and `rolling_update_policy` fields deprecated. `google_compute_instance_group_manager` also now accepts `REPLACE` for `update_strategy`, which is an alias for `RESTART`, and is preferred. ([#2156](https://github.com/terraform-providers/terraform-provider-google/issues/2156))
* project: `google_project`'s `app_engine` sub-block has been deprecated. Please use the `google_app_engine_app` resource instead. Changing between the two should not force project re-creation. ([#2147](https://github.com/terraform-providers/terraform-provider-google/issues/2147))
* project: `google_project_iam_policy`'s `restore_policy` field is now deprecated ([#2186](https://github.com/terraform-providers/terraform-provider-google/issues/2186))

FEATURES: 
* **New Datasource**: `google_compute_instance` ([#1906](https://github.com/terraform-providers/terraform-provider-google/issues/1906))
* **New Resource**: `google_compute_interconnect_attachment` ([#1140](https://github.com/terraform-providers/terraform-provider-google/issues/1140))
* **New Resource**: `google_filestore_instance` ([#2088](https://github.com/terraform-providers/terraform-provider-google/issues/2088))
* **New Resource**: `google_app_engine_application` ([#2147](https://github.com/terraform-providers/terraform-provider-google/issues/2147))

ENHANCEMENTS:
* container: Add `enable_tpu` flag to google_container_cluster ([#1974](https://github.com/terraform-providers/terraform-provider-google/issues/1974))
* dns: `google_dns_managed_zone` is now importable ([#1944](https://github.com/terraform-providers/terraform-provider-google/issues/1944))
* dns: `google_dns_managed_zone` is now entirely GA ([#2154](https://github.com/terraform-providers/terraform-provider-google/issues/2154))
* runtimeconfig: `google_runtimeconfig_config` and `google_runtimeconfig_variable` are now importable. ([#2054](https://github.com/terraform-providers/terraform-provider-google/issues/2054))
* services: containeranalysis.googleapis.com can now be enabled ([#2095](https://github.com/terraform-providers/terraform-provider-google/issues/2095))

BUG FIXES:
* compute: fix instance template interaction with regional disk self links ([#2138](https://github.com/terraform-providers/terraform-provider-google/issues/2138))
* compute: fix diff when using image shorthands for instance templates ([#1995](https://github.com/terraform-providers/terraform-provider-google/issues/1995))
* compute: fix error when reading instance templates created from disks and referenced by name instead of self_link ([#2153](https://github.com/terraform-providers/terraform-provider-google/issues/2153))
* container: Make max_pods_per_node ForceNew ([#2139](https://github.com/terraform-providers/terraform-provider-google/issues/2139))
* services: make google_project_service more resilient to projects being deleted ([#2090](https://github.com/terraform-providers/terraform-provider-google/issues/2090))
* sql: retry failed sql calls ([#2174](https://github.com/terraform-providers/terraform-provider-google/issues/2174))

## 1.18.0 (September 17, 2018)

BACKWARDS INCOMPATIBILITIES:
* compute: instance templates used to not set any disks in the template in state unless they were in the config, as well. It also only stored the image name in state. Both of these were bugs, and have been fixed. They should not cause any disruption. If you were interpolating an image name from a disk in an instance template, you'll need to update your config to strip out everything before the last `/`. If you imported an instance template, and did not add all the disks in the template to your config, you'll see a diff; add those disks to your config, and it will go away. Those are the only two instances where this change should effect you. We apologise for the inconvenience. ([#1916](https://github.com/terraform-providers/terraform-provider-google/issues/1916))
* iam: `google_*_custom_roles` now treats `delete` as deprecated - to actually delete roles, remove from config.  
* provider: This is the first release tested against and built with Go 1.11, which required go fmt changes to the code. If you are building a custom version of this provider or running tests using the repository Make targets (e.g. make build) when using a previous version of Go, you will receive errors. You can use the underlying go commands (e.g. go build) to workaround the go fmt check in the Make targets until you are able to upgrade Go.

FEATURES: 
* **New Resource**: `google_compute_attached_disk` ([#1585](https://github.com/terraform-providers/terraform-provider-google/issues/1585))
* **New Resource**: `google_composer_environment` ([#2001](https://github.com/terraform-providers/terraform-provider-google/issues/2001))

IMPROVEMENTS:
* bigquery: Add Support For BigQuery Access Control ([#1931](https://github.com/terraform-providers/terraform-provider-google/issues/1931))
* compute: `google_compute_health_check` is autogenerated, exposing the `type` attribute and accepting more import formats. ([#1941](https://github.com/terraform-providers/terraform-provider-google/issues/1941))
* compute: `google_compute_ssl_certificate` is autogenerated, exposing the `creation_timestamp` attribute and accepting more import formats. Note: `certificate_id` was changed to an int from a string. This should have no effect on backwards compatibility, but please report a bug if you have any issues! ([#2015](https://github.com/terraform-providers/terraform-provider-google/issues/2015))
* container: Addition of create_subnetwork and other fields relevant for Alias IPs ([#1921](https://github.com/terraform-providers/terraform-provider-google/issues/1921))
* dataflow: Add region choice to dataflow jobs ([#1979](https://github.com/terraform-providers/terraform-provider-google/issues/1979))
* logging: Add import support for `google_logging_organization_sink`, `google_logging_folder_sink`, `google_logging_billing_account_sink` ([#1860](https://github.com/terraform-providers/terraform-provider-google/issues/1860))
* logging: Sending a default update mask for all logging sinks to prevent future breakages ([#1991](https://github.com/terraform-providers/terraform-provider-google/issues/1991))
* dns: Adding support for labels to managed DNS ([#1803](https://github.com/terraform-providers/terraform-provider-google/issues/1803))
* container: Add support for `max_pods_per_node` for private clusters. ([#2038](https://github.com/terraform-providers/terraform-provider-google/issues/2038))

BUG FIXES:
* compute: Store google_compute_vpn_tunnel.router as a self_link to avoid permadiffs. ([#2003](https://github.com/terraform-providers/terraform-provider-google/issues/2003))
* iam: Prevent error when attempting to recreate recently soft-deleted `google_(project|organization)_iam_custom_role`. Instead, roles that are able to be undeleted will be undeleted-updated, as long as they were deleted within 7 days. ([#1681](https://github.com/terraform-providers/terraform-provider-google/issues/1681))
* project: make validation for project id less restrictive ([#1878](https://github.com/terraform-providers/terraform-provider-google/issues/1878))

## 1.17.1 (August 22, 2018)

BUG FIXES:
* container: fix panic on gke binauth ([#1924](https://github.com/terraform-providers/terraform-provider-google/issues/1924))

## 1.17.0 (August 22, 2018)

FEATURES:
* **New Datasource**: `google_project_services` ([#1822](https://github.com/terraform-providers/terraform-provider-google/issues/1822))
* **New Resource**: `google_compute_region_disk` ([#1755](https://github.com/terraform-providers/terraform-provider-google/issues/1755))
* **New Resource**: `google_binary_authorization_attestor` ([#1885](https://github.com/terraform-providers/terraform-provider-google/issues/1885))
* **New Resource**: `google_binary_authorization_policy` ([#1885](https://github.com/terraform-providers/terraform-provider-google/issues/1885))
* **New Resource**: `google_container_analysis_note` ([#1885](https://github.com/terraform-providers/terraform-provider-google/issues/1885))

IMPROVEMENTS:
* cloudfunctions: Add support for updating function code in place ([#1781](https://github.com/terraform-providers/terraform-provider-google/issues/1781))
* cloudbuild: Add support for substitutions in triggers ([#1810](https://github.com/terraform-providers/terraform-provider-google/issues/1810))
* compute: Bring regional instance groups up to par with zonal instance groups. ([#1809](https://github.com/terraform-providers/terraform-provider-google/issues/1809))
* compute: Add labels to Address and GlobalAddress. ([#1811](https://github.com/terraform-providers/terraform-provider-google/issues/1811))
* container: allow updating node image types ([#1843](https://github.com/terraform-providers/terraform-provider-google/issues/1843))
* container: Add support for binary authorization in GKE ([#1884](https://github.com/terraform-providers/terraform-provider-google/issues/1884))
* compute: Allow update of master auth on GKE container cluster. ([#1873](https://github.com/terraform-providers/terraform-provider-google/issues/1873))
* compute: Add support for `boot_disk_type` to `google_dataproc_cluster`. ([#1855](https://github.com/terraform-providers/terraform-provider-google/issues/1855))
* compute: Generate resource_compute_firewall in magic-modules. Make more fields updatable by using PATCH instead of PUT. ([#1907](https://github.com/terraform-providers/terraform-provider-google/issues/1907))
* storage: Add user_project support to `google_storage_project_service_account` data source ([#1913](https://github.com/terraform-providers/terraform-provider-google/issues/1913))

BUG FIXES:
* project: Fix bug where app engine wasn't getting enabled on projects that had billing enabled ([#1795](https://github.com/terraform-providers/terraform-provider-google/issues/1795))
* redis: Allow authorized network to be a name or self link ([#1782](https://github.com/terraform-providers/terraform-provider-google/issues/1782))
* sql: lock on master name when creating replicas ([#1798](https://github.com/terraform-providers/terraform-provider-google/issues/1798))
* storage: allow all role-entity pairs to be unordered ([#1787](https://github.com/terraform-providers/terraform-provider-google/issues/1787))
* compute: allow switching from a daily `ubuntu-minimal` build to `ubuntu-minimal-lts` instead of only `ubuntu`. ([#1870](https://github.com/terraform-providers/terraform-provider-google/issues/1870))
* kms: allow project ids with colons ([#1865](https://github.com/terraform-providers/terraform-provider-google/issues/1865))
* compute: allow project iam policy import with a resource that doesn't match provider project. ([#1875](https://github.com/terraform-providers/terraform-provider-google/issues/1875))
* compute: Ensure regional container clusters update correctly.  ([#1887](https://github.com/terraform-providers/terraform-provider-google/issues/1887))

## 1.16.2 (July 18, 2018)

BUG FIXES:
* compute: use patch instead of put to update router ([#1780](https://github.com/terraform-providers/terraform-provider-google/issues/1780))
* compute: allow a lot more fields in `google_compute_firewall` to be updated to their empty value ([#1784](https://github.com/terraform-providers/terraform-provider-google/issues/1784))
* compute: allow setting instance scheduling booleans on `google_compute_instance` to false ([#1779](https://github.com/terraform-providers/terraform-provider-google/issues/1779))
* compute: ensure router peers and interfaces are always removed.  ([#1877](https://github.com/terraform-providers/terraform-provider-google/issues/1877))

## 1.16.1 (July 16, 2018)

BUG FIXES:
* container: Fix crash when updating resource labels on a cluster ([#1769](https://github.com/terraform-providers/terraform-provider-google/issues/1769))

## 1.16.0 (July 12, 2018)

FEATURES:
* **New Resource**: `compute_instance_from_template` ([#1652](https://github.com/terraform-providers/terraform-provider-google/issues/1652))

IMPROVEMENTS:
* compute: Autogenerate `google_compute_forwarding_rule`, adding labels, service labels, and service name attribute.
* compute: add `quic_override` to `google_compute_target_https_proxy` ([#1718](https://github.com/terraform-providers/terraform-provider-google/issues/1718))
* compute: add support for licenses to `compute_image` ([#1717](https://github.com/terraform-providers/terraform-provider-google/issues/1717))
* compute: Autogenerate router resource. Also adds update support and a few new fields (advertise_mode, advertised_groups, advertised_ip_ranges). ([#1723](https://github.com/terraform-providers/terraform-provider-google/issues/1723))
* container: add ability to configure resource labels on `google_container_cluster` ([#1663](https://github.com/terraform-providers/terraform-provider-google/issues/1663))
* container: increase max number of `master_authorized_networks` to 20 ([#1733](https://github.com/terraform-providers/terraform-provider-google/issues/1733))
* container: support specifying `disk_type` for `node_config` ([#1665](https://github.com/terraform-providers/terraform-provider-google/issues/1665))
* project: correctly paginate when more than 50 services are enabled ([#1737](https://github.com/terraform-providers/terraform-provider-google/issues/1737))
* redis: Support Redis Configuration ([#1706](https://github.com/terraform-providers/terraform-provider-google/issues/1706))

BUG FIXES:
* all: Fix retries for wrapped errors ([#1760](https://github.com/terraform-providers/terraform-provider-google/issues/1760))
* iot: Retry creation of Cloud IoT registry ([#1713](https://github.com/terraform-providers/terraform-provider-google/issues/1713))
* project: ignore stackdriverprovisioning service, so it doesn't permadiff ([#1763](https://github.com/terraform-providers/terraform-provider-google/issues/1763))

## 1.15.0 (June 25, 2018)

FEATURES:

IMPROVEMENTS:
* compute: Autogenerate `compute_subnetwork` ([#1661](https://github.com/terraform-providers/terraform-provider-google/issues/1661))
* container: Allow specifying project when importing container_node_pool ([#1653](https://github.com/terraform-providers/terraform-provider-google/issues/1653))
* dns: Add update support for `dns_managed_zone` ([#1617](https://github.com/terraform-providers/terraform-provider-google/issues/1617))
* project: App Engine application fields can now be updated in-place where possible ([#1621](https://github.com/terraform-providers/terraform-provider-google/issues/1621))
* storage: Add `project` field for GCS service account data source ([#1677](https://github.com/terraform-providers/terraform-provider-google/issues/1677))
* sql: Attempting to shrink an `sql_database_instance`'s disk size will now force recreation of the resource ([#1684](https://github.com/terraform-providers/terraform-provider-google/issues/1684))

BUG FIXES:
* all: Check for done operations before waiting on them. This fixes a 403 we were getting when trying to enable already-enabled services. ([#1632](https://github.com/terraform-providers/terraform-provider-google/issues/1632))
* bigquery: add error checking for bigquery dataset id ([#1638](https://github.com/terraform-providers/terraform-provider-google/issues/1638))
* compute: Store v1 `self_link` for `(sub)?network` in `google_compute_instance` ([#1629](https://github.com/terraform-providers/terraform-provider-google/issues/1629))
* compute: `zone` field in `google_compute_disk` should be optional ([#1631](https://github.com/terraform-providers/terraform-provider-google/issues/1631))
* compute: name_prefix is no longer deprecated for SSL certificates ([#1622](https://github.com/terraform-providers/terraform-provider-google/issues/1622))
* compute: for global address ip_version, IPV4 and empty are equivalent. ([#1639](https://github.com/terraform-providers/terraform-provider-google/issues/1639))
* compute: fix default service account data source to actually set the email and project ([#1690](https://github.com/terraform-providers/terraform-provider-google/issues/1690))
* container: fix permadiff on `container_cluster`'s `pod_security_policy_config` ([#1670](https://github.com/terraform-providers/terraform-provider-google/issues/1670))
* container: removing sub-blocks of `container_cluster` like maintenance windows will now delete them from the API ([#1685](https://github.com/terraform-providers/terraform-provider-google/issues/1685))
* container: retry node pool writes on failed precondition ([#1660](https://github.com/terraform-providers/terraform-provider-google/issues/1660))
* iam: Fixes issue with consecutive whitespace ([#1625](https://github.com/terraform-providers/terraform-provider-google/issues/1625))
* iam: use same mutex for project_iam_policy as the other project_iam resources ([#1645](https://github.com/terraform-providers/terraform-provider-google/issues/1645))
* iam: don't error if service account key is already gone on delete ([#1659](https://github.com/terraform-providers/terraform-provider-google/issues/1659))
* iam: Fix bug in v1.14 where service_account_key needed project set ([#1664](https://github.com/terraform-providers/terraform-provider-google/issues/1664))
* iot: fix updatemask so updates actually work ([#1640](https://github.com/terraform-providers/terraform-provider-google/issues/1640))
* storage: fix a permadiff in bucket ACL role entities ([#1692](https://github.com/terraform-providers/terraform-provider-google/issues/1692))

## 1.14.0 (June 07, 2018)

FEATURES:
* **New Datasource**: `google_service_account` ([#1535](https://github.com/terraform-providers/terraform-provider-google/issues/1535))
* **New Datasource**: `google_service_account_key` ([#1535](https://github.com/terraform-providers/terraform-provider-google/issues/1535))
* **New Datasource**: `google_netblock_ip_ranges` ([#1580](https://github.com/terraform-providers/terraform-provider-google/issues/1580))
* **New Datasource**: `google_compute_regions` ([#1603](https://github.com/terraform-providers/terraform-provider-google/issues/1603))

IMPROVEMENTS:
* compute: As part of migrating `google_compute_disk` to be autogenerated, enabled encrypted source snapshot & images. [[#1521](https://github.com/terraform-providers/terraform-provider-google/issues/1521)].
* compute: Accept subnetwork name only in `google_forwarding_rule` ([#1552](https://github.com/terraform-providers/terraform-provider-google/issues/1552))
* compute: Add disabled property to `google_compute_firewall` ([#1536](https://github.com/terraform-providers/terraform-provider-google/issues/1536))
* compute: Add support for custom request headers in `google_compute_backend_service` ([#1537](https://github.com/terraform-providers/terraform-provider-google/issues/1537))
* compute: Add support for `ssl_policy` to `google_compute_target_ssl_proxy` ([#1568](https://github.com/terraform-providers/terraform-provider-google/issues/1568))
* compute: Add support for `version`s in instance group manager ([#1499](https://github.com/terraform-providers/terraform-provider-google/issues/1499))
* compute: Add support for `network_tier` to address, instance and instance_template ([#1530](https://github.com/terraform-providers/terraform-provider-google/issues/1530))
* cloudbuild: Use the project defined in `trigger_template` when creating a `google_cloudbuild_trigger` ([#1556](https://github.com/terraform-providers/terraform-provider-google/issues/1556))
* cloudbuild: Support configuration file in repository for `google_cloudbuild_trigger` ([#1557](https://github.com/terraform-providers/terraform-provider-google/issues/1557))
* kms: Add basic update for `google_kms_crypto_key` resource ([#1511](https://github.com/terraform-providers/terraform-provider-google/issues/1511))
* project: Use default provider project for `google_project_services` if project field is empty ([#1553](https://github.com/terraform-providers/terraform-provider-google/issues/1553))
* project: Added support for restoring default organization policies ([#1477](https://github.com/terraform-providers/terraform-provider-google/issues/1477))
* project: Handle spurious Cloud API errors and performance issues for `google_project_service(s)` ([#1565](https://github.com/terraform-providers/terraform-provider-google/issues/1565))
* redis: Add update support for Redis Instances ([#1590](https://github.com/terraform-providers/terraform-provider-google/issues/1590))
* sql: Add labels support in `sql_database_instance` ([#1567](https://github.com/terraform-providers/terraform-provider-google/issues/1567))

BUG FIXES:
* dns: Suppress diff for ipv6 address in `google_dns_record_set` ([#1551](https://github.com/terraform-providers/terraform-provider-google/issues/1551))
* storage: Support removing a label in `google_storage_bucket` ([#1550](https://github.com/terraform-providers/terraform-provider-google/issues/1550))
* compute: Fix perpetual diff caused by the `google_instance_group` self_link in `google_regional_instance_group_manager` ([#1549](https://github.com/terraform-providers/terraform-provider-google/issues/1549))
* project: Retry while listing enabled services ([#1573](https://github.com/terraform-providers/terraform-provider-google/issues/1573))
* redis: Allow self links for redis authorized network ([#1599](https://github.com/terraform-providers/terraform-provider-google/issues/1599))

## 1.13.0 (May 24, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:
* `google_project_service`/`google_project_services` now use the [Service Usage API](https://cloud.google.com/service-usage). Users of those resources will need to enable the API at https://console.cloud.google.com/apis/api/serviceusage.googleapis.com.
* If you have a `google_project` resource where App Engine is enabled in the project, add an `app_engine` [block](https://www.terraform.io/docs/providers/google/r/google_project.html#app_engine) to your resource before running Terraform after upgrading to this version, or hold off on upgrading for now. See [#1561](https://github.com/terraform-providers/terraform-provider-google/issues/1561), which has more details and an ongoing investigation of other potential fixes.

FEATURES:
* **New Resource**: `google_cloudbuild_trigger`. ([#1357](https://github.com/terraform-providers/terraform-provider-google/issues/1357))
* **New Resource**: `google_storage_bucket_iam_policy` ([#1190](https://github.com/terraform-providers/terraform-provider-google/issues/1190))
* **New Resource**: `google_resource_manager_lien` ([#1484](https://github.com/terraform-providers/terraform-provider-google/issues/1484))
* **New Resource**: `google_logging_billing_account_exclusion` ([#990](https://github.com/terraform-providers/terraform-provider-google/issues/990))
* **New Resource**: `google_logging_folder_exclusion` ([#990](https://github.com/terraform-providers/terraform-provider-google/issues/990))
* **New Resource**: `google_logging_organization_exclusion` ([#990](https://github.com/terraform-providers/terraform-provider-google/issues/990))
* **New Resource**: `google_logging_project_exclusion` ([#990](https://github.com/terraform-providers/terraform-provider-google/issues/990))
* **New Resource**: `google_redis_instance` ([#1485](https://github.com/terraform-providers/terraform-provider-google/issues/1485))
* App Engine applications can now be managed using the `app_engine` field in `google_project` ([#1503](https://github.com/terraform-providers/terraform-provider-google/issues/1503))

IMPROVEMENTS:
* cloudfunctions: add ability to retry cloud functions on failure ([#1452](https://github.com/terraform-providers/terraform-provider-google/issues/1452))
* container: Add support for regional cluster in `google_container` datasource ([#1441](https://github.com/terraform-providers/terraform-provider-google/issues/1441))
* container: Add GKE Shared VPC support ([#1528](https://github.com/terraform-providers/terraform-provider-google/issues/1528))
* compute: autogenerate `google_compute_ssl_policy` ([#1478](https://github.com/terraform-providers/terraform-provider-google/issues/1478))
* compute: add support for `ssl_policy` to `google_target_https_proxy` ([#1466](https://github.com/terraform-providers/terraform-provider-google/issues/1466))
* project: Added name and project_id plan-time validations ([#1519](https://github.com/terraform-providers/terraform-provider-google/issues/1519))

BUG FIXES:
* compute: Compare region_backend_service.backend[].group as a relative path ([#1487](https://github.com/terraform-providers/terraform-provider-google/issues/1487))
* compute: Fixed `region_backend_service` to calc hash using relative path ([#1491](https://github.com/terraform-providers/terraform-provider-google/issues/1491))
* sql: Fix panic on empty maintenance window ([#1507](https://github.com/terraform-providers/terraform-provider-google/issues/1507))

## 1.12.0 (May 04, 2018)
FEATURES:
* spanner: New resources to manage IAM for Spanner Databases: google_spanner_database_iam_binding, google_spanner_database_iam_member, and google_spanner_database_iam_policy ([#1386](https://github.com/terraform-providers/terraform-provider-google/issues/1386))
* spanner: New resources to manage IAM for Spanner Instances: google_spanner_instance_iam_binding, google_spanner_instance_iam_member, and google_spanner_instance_iam_policy ([#1387](https://github.com/terraform-providers/terraform-provider-google/issues/1387))

IMPROVEMENTS:
* compute: Autogenerate `google_vpn_gateway` ([#1409](https://github.com/terraform-providers/terraform-provider-google/issues/1409))
* compute: add `enable_flow_logs` field to subnetwork ([#1385](https://github.com/terraform-providers/terraform-provider-google/issues/1385))
* project: Don't fail if `folder_id` and `org_id` are set but one is empty for `google_project` ([#1425](https://github.com/terraform-providers/terraform-provider-google/issues/1425))

BUG FIXES:
* compute: Always parse fixed64 string to int64 even on 32 bits platform to prevent out-of-range crash. ([#1429](https://github.com/terraform-providers/terraform-provider-google/issues/1429))

## 1.11.0 (May 01, 2018)

IMPROVEMENTS:
* compute: Add `public_ptr_domain_name` to `google_compute_instance`.  ([#1349](https://github.com/terraform-providers/terraform-provider-google/issues/1349))
* compute: Autogenerate `google_compute_global_address`. ([#1379](https://github.com/terraform-providers/terraform-provider-google/issues/1379))
* compute: Autogenerate `google_compute_target_http_proxy`. ([#1391](https://github.com/terraform-providers/terraform-provider-google/issues/1391))
* compute: Autogenerate `google_compute_target_http_proxy`. ([#1373](https://github.com/terraform-providers/terraform-provider-google/issues/1373))
* compute: Simplify autogenerated code for `google_compute_target_http_proxy` and `google_compute_target_ssl_proxy`. ([#1395](https://github.com/terraform-providers/terraform-provider-google/issues/1395))
* compute: Use partial state setting in `google_compute_target_http_proxy` and `google_compute_target_ssl_proxy` to better handle mid-update errors. ([#1392](https://github.com/terraform-providers/terraform-provider-google/issues/1392))
* compute: Use the v1 API for `google_compute_address` ([#1384](https://github.com/terraform-providers/terraform-provider-google/issues/1384))
* compute: Properly detect when `public_ptr_domain_name` isn't set. ([#1383](https://github.com/terraform-providers/terraform-provider-google/issues/1383))
* compute: Use the v1 API for `google_compute_ssl_policy` ([#1368](https://github.com/terraform-providers/terraform-provider-google/issues/1368))
* container: Add `issue_client_certificate` to `google_container_cluster`. ([#1396](https://github.com/terraform-providers/terraform-provider-google/issues/1396))
* container: Support regional clusters for node pools. ([#1320](https://github.com/terraform-providers/terraform-provider-google/issues/1320))
* all: List of resources is now partially auto-generated ([#1397](https://github.com/terraform-providers/terraform-provider-google/issues/1397)] [[#1402](https://github.com/terraform-providers/terraform-provider-google/issues/1402))

BUG FIXES:
* iam: expand the validation for service accounts to include App Engine and compute default service accounts ([#1390](https://github.com/terraform-providers/terraform-provider-google/issues/1390))
* sql: Increase timeouts ([#1381](https://github.com/terraform-providers/terraform-provider-google/issues/1381))
* website: fix broken layouts ([#1405](https://github.com/terraform-providers/terraform-provider-google/issues/1405))

## 1.10.0 (April 20, 2018)

FEATURES:
* **New Data Source** `google_folder` ([#1280](https://github.com/terraform-providers/terraform-provider-google/issues/1280))
* **New Resource** `google_compute_subnetwork_iam_binding` ([#1305](https://github.com/terraform-providers/terraform-provider-google/issues/1305))
* **New Resource** `google_compute_subnetwork_iam_member` ([#1305](https://github.com/terraform-providers/terraform-provider-google/issues/1305))
* **New Resource** `google_compute_subnetwork_iam_policy` ([#1305](https://github.com/terraform-providers/terraform-provider-google/issues/1305))

IMPROVEMENTS:
* compute: Add timeouts to `google_compute_snapshot` ([#1309](https://github.com/terraform-providers/terraform-provider-google/issues/1309))
* compute: un-deprecate name_prefix for instance templates ([#1328](https://github.com/terraform-providers/terraform-provider-google/issues/1328))
* compute: Add `default_cluster_version` field to `data_source_google_container_engine_versions`. ([#1355](https://github.com/terraform-providers/terraform-provider-google/issues/1355))
* compute: Add `max_connections` and `max_connections_per_instance` to `resource_compute_backend_service` ([#1353](https://github.com/terraform-providers/terraform-provider-google/issues/1353))
* all: Maintain parity with GCP Console UI by allowing removal of default project networks.  ([#1316](https://github.com/terraform-providers/terraform-provider-google/issues/1316))
* all: Use standard user-agent header ([#1332](https://github.com/terraform-providers/terraform-provider-google/issues/1332))

BUG FIXES:
* compute: fix error introduced when attached disks are deleted out of band ([#1301](https://github.com/terraform-providers/terraform-provider-google/issues/1301))
* container: Use correct project id regex in `google_container_cluster` ([#1311](https://github.com/terraform-providers/terraform-provider-google/issues/1311))
* folder: Escape the display name in active folder data source (in case of spaces, etc) ([#1261](https://github.com/terraform-providers/terraform-provider-google/issues/1261))
* project: Fix auto-delete default network in google_project ([#1336](https://github.com/terraform-providers/terraform-provider-google/issues/1336))

## 1.9.0 (April 05, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:
* `name_prefix` is now deprecated in all resources that support it ([#1035](https://github.com/terraform-providers/terraform-provider-google/issues/1035))

FEATURES:
* **New Data Source** `google_compute_ssl_policy` ([#1247](https://github.com/terraform-providers/terraform-provider-google/issues/1247))
* **New Resource** `google_compute_security_policy` ([#1242](https://github.com/terraform-providers/terraform-provider-google/issues/1242))
* **New Resource** `google_compute_ssl_policy` ([#1247](https://github.com/terraform-providers/terraform-provider-google/issues/1247))
* **New Resource** `google_project_organization_policy` ([#1226](https://github.com/terraform-providers/terraform-provider-google/issues/1226))

IMPROVEMENTS:
* all: Read `GOOGLE_CLOUD_PROJECT` environment variable also ([#1271](https://github.com/terraform-providers/terraform-provider-google/issues/1271))
* bigquery: Add time partitioning field to `google_bigquery_table` resource ([#1240](https://github.com/terraform-providers/terraform-provider-google/issues/1240))
* config: Add OAuth access token to `google_client_config` data source [[#1277](https://github.com/terraform-providers/terraform-provider-google/issues/1277)] 
* compute: Add `wait_for_instances` field to `google_compute_instance_group_manager` and self_link option to the `google_compute_instance_group` data source ([#1222](https://github.com/terraform-providers/terraform-provider-google/issues/1222))
* compute: add support for security policies in backend services ([#1243](https://github.com/terraform-providers/terraform-provider-google/issues/1243))
* compute: regional instance group managers now support rolling updates ([#1260](https://github.com/terraform-providers/terraform-provider-google/issues/1260))
* container: add ability to delete the default node pool ([#1245](https://github.com/terraform-providers/terraform-provider-google/issues/1245))
* container: Add update support for pod security policy ([#1195](https://github.com/terraform-providers/terraform-provider-google/issues/1195))
* container: Add gke node taints ([#1264](https://github.com/terraform-providers/terraform-provider-google/issues/1264))
* container: Add support for node pool versions ([#1266](https://github.com/terraform-providers/terraform-provider-google/issues/1266))
* container: Add support for private clusters ([#1250](https://github.com/terraform-providers/terraform-provider-google/issues/1250))
* container: Updates container_cluster to set `enable_legacy_abac` to false by default ([#1281](https://github.com/terraform-providers/terraform-provider-google/issues/1281))
* container: Add support for regional GKE clusters in `google_container_cluster` ([#1181](https://github.com/terraform-providers/terraform-provider-google/issues/1181))
* iam: allow setting service account email as id for service account keys ([#1256](https://github.com/terraform-providers/terraform-provider-google/issues/1256))
* sql: add custom timeouts support for sql database instance ([#1288](https://github.com/terraform-providers/terraform-provider-google/issues/1288))
* sql: Retry on 429 and 503 errors on sql admin operation ([#1212](https://github.com/terraform-providers/terraform-provider-google/issues/1212))
* project: Add disable_on_destroy flag to `google_project_services` ([#1293](https://github.com/terraform-providers/terraform-provider-google/issues/1293))

BUG FIXES:
* compute: fix panic when setting empty iap block ([#1232](https://github.com/terraform-providers/terraform-provider-google/issues/1232))
* compute: protect against an instance getting deleted by an igm while the disk is being detached ([#1241](https://github.com/terraform-providers/terraform-provider-google/issues/1241))
* compute: Add DiffSuppress for URL maps on Target HTTP(S) Proxies ([#1263](https://github.com/terraform-providers/terraform-provider-google/issues/1263))
* storage: Set force_destroy when importing storage buckets ([#1223](https://github.com/terraform-providers/terraform-provider-google/issues/1223))
* storage: Delete all object version when deleting all objects in a bucket ([#1285](https://github.com/terraform-providers/terraform-provider-google/issues/1285))

## 1.8.0 (March 19, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:
* `google_dataproc_cluster.delete_autogen_bucket` is now deprecated ([#1171](https://github.com/terraform-providers/terraform-provider-google/issues/1171))

FEATURES:
* **New Resource** `google_organization_iam_policy` (see docs for caveats) ([#1196](https://github.com/terraform-providers/terraform-provider-google/issues/1196))

IMPROVEMENTS:
* container: un-deprecate `google_container_node_pool.initial_node_count` ([#1176](https://github.com/terraform-providers/terraform-provider-google/issues/1176))
* container: Add support for pod security policy ([#1192](https://github.com/terraform-providers/terraform-provider-google/issues/1192))
* container: Add support for GKE metadata concealment ([#1199](https://github.com/terraform-providers/terraform-provider-google/issues/1199))
* container: Add support for GKE network policy config addon. ([#1200](https://github.com/terraform-providers/terraform-provider-google/issues/1200))
* container: Add support for `instance_group_urls` in `google_container_node_pool` ([#1207](https://github.com/terraform-providers/terraform-provider-google/issues/1207))
* compute: Rolling update support for instance group manager ([#1137](https://github.com/terraform-providers/terraform-provider-google/issues/1137))
* compute: Add `cdn_policy` field to backend service ([#1208](https://github.com/terraform-providers/terraform-provider-google/issues/1208))
* compute: Add support for deletion protection. ([#1205](https://github.com/terraform-providers/terraform-provider-google/issues/1205))
* all: IAM resources now wait for propagation before reporting created. ([#1197](https://github.com/terraform-providers/terraform-provider-google/issues/1197))

BUG FIXES:
* compute: Properly set `image_id` field on `data_google_compute_image` in state ([#1217](https://github.com/terraform-providers/terraform-provider-google/issues/1217))
* compute: Properly set `project` field on `google_compute_project_metadata` in state ([#1217](https://github.com/terraform-providers/terraform-provider-google/issues/1217))
* dataproc: Properly set `cluster_config.0.initialization_action` on `google_dataproc_cluster` in state ([#1217](https://github.com/terraform-providers/terraform-provider-google/issues/1217))

## 1.7.0 (March 12, 2018)

Features:
* **New Data Source** `google_compute_forwarding_rule` ([#1078](https://github.com/terraform-providers/terraform-provider-google/issues/1078))
* **New Data Source** `google_compute_vpn_gateway` ([#1071](https://github.com/terraform-providers/terraform-provider-google/issues/1071))
* **New Data Source** `google_project` ([#1111](https://github.com/terraform-providers/terraform-provider-google/issues/1111))
* **New Data Source** `google_compute_backend_service` ([#1150](https://github.com/terraform-providers/terraform-provider-google/issues/1150))
* **New Data Source** `google_storage_project_service_account` ([#1110](https://github.com/terraform-providers/terraform-provider-google/issues/1110))
* **New Data Source** `google_compute_default_service_account` ([#1119](https://github.com/terraform-providers/terraform-provider-google/issues/1119))
* **New Resource** `google_folder_iam_binding` ([#1076](https://github.com/terraform-providers/terraform-provider-google/issues/1076))
* **New Resource** `google_folder_iam_member` ([#1076](https://github.com/terraform-providers/terraform-provider-google/issues/1076))
* **New Resource** `google_project_usage_export_bucket` ([#1080](https://github.com/terraform-providers/terraform-provider-google/issues/1080))

IMPROVEMENTS:
* compute: add support for updating alias ips in instances ([#1084](https://github.com/terraform-providers/terraform-provider-google/issues/1084))
* compute: allow setting a route resource's `description` attribute ([#1088](https://github.com/terraform-providers/terraform-provider-google/issues/1088))
* compute: allow lowercase ip protocols in forwarding rules ([#1118](https://github.com/terraform-providers/terraform-provider-google/issues/1118))
* compute: `google_compute_zones` datasource accepts a `project` parameter ([#1122](https://github.com/terraform-providers/terraform-provider-google/issues/1122))
* compute: Support `distributionPolicy` when creating regional instance group managers. ([#1092](https://github.com/terraform-providers/terraform-provider-google/issues/1092))
* compute: Timeout customization for `google_compute_backend_bucket`, `google_compute_http_health_check`, and `google_compute_https_health_check` ([#1177](https://github.com/terraform-providers/terraform-provider-google/issues/1177))
* container: Fail if the ip_allocation_policy doesn't specify secondary range names ([#1065](https://github.com/terraform-providers/terraform-provider-google/issues/1065))
* container: Allow specifying accelerators in cluster node_config. ([#1115](https://github.com/terraform-providers/terraform-provider-google/issues/1115))
* pubsub: Add project field to iam pubsub topic resources ([#1154](https://github.com/terraform-providers/terraform-provider-google/issues/1154))
* sql: Support multiple users with the same name for different host for 1st gen SQL instances. ([#1066](https://github.com/terraform-providers/terraform-provider-google/issues/1066))
* sql: Add SQL DB Instance attribute `first_ip_address` ([#1050](https://github.com/terraform-providers/terraform-provider-google/issues/1050))

BUG FIXES:
* compute: Don't store disk in state if it didn't create ([#1129](https://github.com/terraform-providers/terraform-provider-google/issues/1129))
* compute: Check set equality for service account scope changes ([#1130](https://github.com/terraform-providers/terraform-provider-google/issues/1130))
* compute: Disk now accepts project id with '.' and ':' ([#1145](https://github.com/terraform-providers/terraform-provider-google/issues/1145))
* dataproc: fix typos in pyspark dataproc job resource that led to args not working ([#1120](https://github.com/terraform-providers/terraform-provider-google/issues/1120))
* dns: fix perpetual diffs when names aren't all uppercase or if TXT records aren't quoted ([#1141](https://github.com/terraform-providers/terraform-provider-google/issues/1141))
* spanner: Accepts project id with '.' and ':' ([#1151](https://github.com/terraform-providers/terraform-provider-google/issues/1151))

## 1.6.0 (February 09, 2018)

Features:
* **New Resource** `google_cloudiot_registry` ([#970](https://github.com/terraform-providers/terraform-provider-google/issues/970))
* **New Resource** `google_endpoints_service` ([#933](https://github.com/terraform-providers/terraform-provider-google/issues/933))
* **New Resource** `google_storage_default_object_acl` ([#992](https://github.com/terraform-providers/terraform-provider-google/issues/992))
* **New Resource** `google_storage_notification` ([#1033](https://github.com/terraform-providers/terraform-provider-google/issues/1033))

IMPROVEMENTS:
* compute: Suppress diff if `guest_accelerators` count is 0 in `google_compute_instance` and `google_compute_instance_template` ([#866](https://github.com/terraform-providers/terraform-provider-google/issues/866))
* compute: Add update support for machine type, min cpu platform, and service accounts ([#1005](https://github.com/terraform-providers/terraform-provider-google/issues/1005))
* compute: Add import support for google_compute_shared_vpc_host_project/google_compute_shared_vpc_service_project resources ([#1004](https://github.com/terraform-providers/terraform-provider-google/issues/1004))
* compute: Make route priority optional since Compute has a default value. ([#1009](https://github.com/terraform-providers/terraform-provider-google/issues/1009))
* container: Suppress diff for empty/default provider in `google_container_cluster` network policy [#1031](https://github.com/terraform-providers/terraform-provider-google/issues/1031)
* container: Return an error if name and name prefix are specified in node pool ([#1062](https://github.com/terraform-providers/terraform-provider-google/issues/1062))
* sql: Support for PostgreSQL high availability ([#1001](https://github.com/terraform-providers/terraform-provider-google/issues/1001))
* sql: Support for ServerCaCert in Cloud SQL instance. (Related to [#635](https://github.com/terraform-providers/terraform-provider-google/issues/635))
* storage: Add support for setting bucket's logging config ([#946](https://github.com/terraform-providers/terraform-provider-google/issues/946))


BUG FIXES:

* project: Fix crash when errors are encountered updating a `google_project` ([#1016](https://github.com/terraform-providers/terraform-provider-google/issues/1016))
* logging: Set project during import for `google_logging_project_sink` to avoid recreation ([#1018](https://github.com/terraform-providers/terraform-provider-google/issues/1018))
* compute: Suppress diff on image field when referring to unconventional public image family naming pattern ([#1024](https://github.com/terraform-providers/terraform-provider-google/issues/1024))
* compute: Backend service backed by a group couldn't be created or updated because both max_rate and max_rate_per_instance would always be set to zero and they can't be both set. ([#1051](https://github.com/terraform-providers/terraform-provider-google/issues/1051))
* container: Fix perpetual diff in `google_container_cluster` if the subnetwork field is not specified ([#1061](https://github.com/terraform-providers/terraform-provider-google/issues/1061))

## 1.5.0 (January 18, 2018)

FEATURES:
* **New Resource:** `google_cloudfunctions_function` ([#899](https://github.com/terraform-providers/terraform-provider-google/issues/899))
* **New Resource:** `google_logging_organization_sink` ([#923](https://github.com/terraform-providers/terraform-provider-google/issues/923))
* **New Resource:** `google_service_account_iam_binding` ([#840](https://github.com/terraform-providers/terraform-provider-google/issues/840))
* **New Resource:** `google_service_account_iam_member` ([#840](https://github.com/terraform-providers/terraform-provider-google/issues/840))
* **New Resource:** `google_service_account_iam_policy` ([#840](https://github.com/terraform-providers/terraform-provider-google/issues/840))
* **New Resource:** `google_pubsub_topic_iam_binding` ([#875](https://github.com/terraform-providers/terraform-provider-google/issues/875))
* **New Resource:** `google_pubsub_topic_iam_member` ([#875](https://github.com/terraform-providers/terraform-provider-google/issues/875))
* **New Resource:** `google_pubsub_topic_iam_policy` ([#875](https://github.com/terraform-providers/terraform-provider-google/issues/875))
* **New Resource:** `google_dataflow_job` ([#855](https://github.com/terraform-providers/terraform-provider-google/issues/855))
* **New Data Source:** `google_compute_region_instance_group` ([#851](https://github.com/terraform-providers/terraform-provider-google/issues/851))
* **New Data Source:** `google_container_cluster` ([#740](https://github.com/terraform-providers/terraform-provider-google/issues/740))
* **New Data Source:** `google_kms_secret` ([#741](https://github.com/terraform-providers/terraform-provider-google/issues/741))
* **New Data Source:** `google_billing_account`([#889](https://github.com/terraform-providers/terraform-provider-google/issues/889))
* **New Data Source:** `google_organization` ([#887](https://github.com/terraform-providers/terraform-provider-google/issues/887))
* **New Data Source:** `google_container_registry_repository` ([#954](https://github.com/terraform-providers/terraform-provider-google/issues/954))
* **New Data Source:** `google_container_registry_image` ([#954](https://github.com/terraform-providers/terraform-provider-google/issues/954))

IMPROVEMENTS:
* iam: Add support for import of IAM resources (project, folder, organizations, crypto keys, and key rings).  ([#835](https://github.com/terraform-providers/terraform-provider-google/issues/835))
* compute: Add support for routing mode in compute network. ([#838](https://github.com/terraform-providers/terraform-provider-google/issues/838))
* compute: Add configurable create/update/delete timeouts to `google_compute_instance` ([#856](https://github.com/terraform-providers/terraform-provider-google/issues/856))
* compute: Add configurable create/update/delete timeouts to `google_compute_subnetwork` ([#871](https://github.com/terraform-providers/terraform-provider-google/issues/871))
* compute: Add update support for `routing_mode` in `google_compute_network` ([#857](https://github.com/terraform-providers/terraform-provider-google/issues/857))
* compute: Add import support for `google_compute_instance` ([#873](https://github.com/terraform-providers/terraform-provider-google/issues/873))
* compute: More descriptive error message for health check not found in `google_compute_target_pool` ([#883](https://github.com/terraform-providers/terraform-provider-google/issues/883))
* compute: Add `disable_on_destroy` (default true) for `google_project_service`. ([#965](https://github.com/terraform-providers/terraform-provider-google/issues/965))
* compute: Add update support for subnetwork IP CIDR range expansion ([#945](https://github.com/terraform-providers/terraform-provider-google/issues/945))
* compute: Read boot disk initialization params from API in `google_compute_instance` ([#948](https://github.com/terraform-providers/terraform-provider-google/issues/948))
* container: Ensure operations on a cluster are applied serially ([#937](https://github.com/terraform-providers/terraform-provider-google/issues/937))
* container: Don't recreate container_cluster when maintenance_window changes ([#893](https://github.com/terraform-providers/terraform-provider-google/issues/893))
* dataproc: Add "internal IP only" support for Dataproc clusters ([#837](https://github.com/terraform-providers/terraform-provider-google/issues/837))
* dataproc: Support `self_link` from a different project in dataproc network and subnetwork fields ([#935](https://github.com/terraform-providers/terraform-provider-google/issues/935))
* sourcerepo: Export new `url` field for `google_sourcerepo_repository` ([#943](https://github.com/terraform-providers/terraform-provider-google/issues/943))
* folder: Support more format for `folder` field in `google_folder_organization_policy` ([#963](https://github.com/terraform-providers/terraform-provider-google/issues/963))
* dns: Add import support to `google_dns_record_set` ([#895](https://github.com/terraform-providers/terraform-provider-google/issues/895))
* all: Make provider-wide region optional ([#916](https://github.com/terraform-providers/terraform-provider-google/issues/916))
* all: Infers region from zone schema before using the provider-level region ([#938](https://github.com/terraform-providers/terraform-provider-google/issues/938))
* all: Upgrade terraform core to v0.11.2 ([#940](https://github.com/terraform-providers/terraform-provider-google/issues/940))

BUG FIXES:
* compute: Suppress diff for equivalent value in `google_compute_disk` image field ([#884](https://github.com/terraform-providers/terraform-provider-google/issues/884))
* compute: Read IAP settings properly in `google_compute_backend_service` ([#907](https://github.com/terraform-providers/terraform-provider-google/issues/907))
* compute: Fix bug causing a crash when specifying unknown network in `google_compute_network_peering` ([#918](https://github.com/terraform-providers/terraform-provider-google/issues/918))
* compute: Fix failing update when changing `google_compute_health_check` type ([#944](https://github.com/terraform-providers/terraform-provider-google/issues/944))
* compute: Fix bug blocking `google_compute_autoscaler` from containing multiple metrics. ([#966](https://github.com/terraform-providers/terraform-provider-google/issues/966))
* container: Set default scopes when creating GKE clusters/node pools ([#924](https://github.com/terraform-providers/terraform-provider-google/issues/924))
* storage: Fix bug blocking the update of a storage object if its content is dynamic/interpolated ([#848](https://github.com/terraform-providers/terraform-provider-google/issues/848))
* storage: Fix bug preventing the removal of lifecycle rules for a `google_storage_bucket` ([#850](https://github.com/terraform-providers/terraform-provider-google/issues/850))
* all: Fix bug causing a perpetual diff when using provider-default zone ([#914](https://github.com/terraform-providers/terraform-provider-google/issues/914))

## 1.4.0 (December 11, 2017)

FEATURES:
* **New Data Source:** `google_compute_image` ([#128](https://github.com/terraform-providers/terraform-provider-google/issues/128))
* **New Resource:** `google_storage_bucket_iam_binding` ([#822](https://github.com/terraform-providers/terraform-provider-google/issues/822))
* **New Resource:** `google_storage_bucket_iam_member` ([#822](https://github.com/terraform-providers/terraform-provider-google/issues/822))

IMPROVEMENTS:

* all: Add support for `zone` at the provider level, as a default for all zonal resources.  ([#816](https://github.com/terraform-providers/terraform-provider-google/issues/816))
* compute: Add support for `min_cpu_platform` to `google_compute_instance_template` ([#808](https://github.com/terraform-providers/terraform-provider-google/issues/808))
* compute: Add example for Shared VPC (aka cross-project networking, or XPN). ([#810](https://github.com/terraform-providers/terraform-provider-google/issues/810))

BUG FIXES:

* all: Fix bug that disallowed using file paths for credentials ([#832](https://github.com/terraform-providers/terraform-provider-google/issues/832))
* dns: Fix bug that broke NS records on subdomains ([#807](https://github.com/terraform-providers/terraform-provider-google/issues/807))
* bigquery: Fix bug causing a crash if the import id was invalid ([#828](https://github.com/terraform-providers/terraform-provider-google/issues/828))

## 1.3.0 (November 30, 2017)

FEATURES:
* **New Resource:** `google_folder_organization_policy` ([#747](https://github.com/terraform-providers/terraform-provider-google/issues/747))
* **New Resource:** `google_kms_key_ring_iam_binding` ([#781](https://github.com/terraform-providers/terraform-provider-google/issues/781))
* **New Resource:** `google_kms_key_ring_iam_member` ([#781](https://github.com/terraform-providers/terraform-provider-google/issues/781))
* **New Resource:** `google_kms_crypto_key_iam_binding` ([#781](https://github.com/terraform-providers/terraform-provider-google/issues/781))
* **New Resource:** `google_kms_crypto_key_iam_member` ([#781](https://github.com/terraform-providers/terraform-provider-google/issues/781))
* **New Resource:** `google_project_custom_iam_role` ([#709](https://github.com/terraform-providers/terraform-provider-google/issues/709))
* **New Resource:** `google_organization_custom_iam_role` ([#735](https://github.com/terraform-providers/terraform-provider-google/issues/735))
* **New Resource:** `google_organization_iam_binding` ([#775](https://github.com/terraform-providers/terraform-provider-google/issues/775))
* **New Resource:** `google_organization_iam_member` ([#775](https://github.com/terraform-providers/terraform-provider-google/issues/775))
* **New Resource:** `google_dataproc_job` ([#253](https://github.com/terraform-providers/terraform-provider-google/issues/253))
* **New Data Source:** `google_active_folder` ([#738](https://github.com/terraform-providers/terraform-provider-google/issues/738))
* **New Data Source:** `google_compute_address` ([#748](https://github.com/terraform-providers/terraform-provider-google/issues/748))
* **New Data Source:** `google_compute_global_address` ([#759](https://github.com/terraform-providers/terraform-provider-google/issues/759))

IMPROVEMENTS:
* compute: Add import support for `google_compute_ssl_certificates` ([#678](https://github.com/terraform-providers/terraform-provider-google/issues/678))
* compute: Add import support for `google_compute_target_http_proxy` ([#678](https://github.com/terraform-providers/terraform-provider-google/issues/678))
* compute: Add import support for `google_compute_target_https_proxy` ([#678](https://github.com/terraform-providers/terraform-provider-google/issues/678))
* compute: Add partial import support for `google_compute_url_map` ([#678](https://github.com/terraform-providers/terraform-provider-google/issues/678))
* compute: Add import support for `google_compute_backend_bucket` ([#736](https://github.com/terraform-providers/terraform-provider-google/issues/736))
* compute: Add configurable timeouts for disks ([#717](https://github.com/terraform-providers/terraform-provider-google/issues/717))
* compute: Use v1 API now that all beta features are in GA for `google_compute_firewall` [[#768](https://github.com/terraform-providers/terraform-provider-google/issues/768)] 
* compute: Add Alias IP and Guest Accelerator support to Instance Templates ([#639](https://github.com/terraform-providers/terraform-provider-google/issues/639))
* container: Relax diff on `daily_maintenance_window.start_time` for `google_container_cluster` ([#726](https://github.com/terraform-providers/terraform-provider-google/issues/726))
* container: Allow node pools with size 0 ([#752](https://github.com/terraform-providers/terraform-provider-google/issues/752))
* container: Add support for `google_container_node_pool` management ([#669](https://github.com/terraform-providers/terraform-provider-google/issues/669))
* container: Add container cluster network policy ([#630](https://github.com/terraform-providers/terraform-provider-google/issues/630))
* container: add support for ip aliasing in `google_container_cluster` ([#654](https://github.com/terraform-providers/terraform-provider-google/issues/654))
* kms: Adds support for creating KMS CryptoKeys resources ([#692](https://github.com/terraform-providers/terraform-provider-google/issues/692))
* project: Add validation for `account_id` in `google_service_account` ([#793](https://github.com/terraform-providers/terraform-provider-google/issues/793))
* storage: Detect file changes in `google_storage_bucket_object` when using source field ([#789](https://github.com/terraform-providers/terraform-provider-google/issues/789))
* all: Consistently store the project and region fields value in state. ([#784](https://github.com/terraform-providers/terraform-provider-google/issues/784))

BUG FIXES:
* bigquery: Set UseLegacySql to true for compatibility with the BigQuery API ([#724](https://github.com/terraform-providers/terraform-provider-google/issues/724))
* compute: Fix perpetual diff with `next_hop_instance` field in `google_compute_route` ([#716](https://github.com/terraform-providers/terraform-provider-google/issues/716))
* compute: Restore the `ipv4_range` field to `google_compute_network` to support legacy VPCs ([#805](https://github.com/terraform-providers/terraform-provider-google/issues/805))
* project: Fix timeout issue with project services ([#737](https://github.com/terraform-providers/terraform-provider-google/issues/737))
* sql: Fix perpetual diff with `authorized_networks` field in `google_sql_database_instance` ([#733](https://github.com/terraform-providers/terraform-provider-google/issues/733))
* sql: give disk_autoresize a default in `google_sql_database_instance` ([#806](https://github.com/terraform-providers/terraform-provider-google/issues/806))

## 1.2.0 (November 09, 2017)

FEATURES:

* **New Resource:** `google_service_account_key` ([#472](https://github.com/terraform-providers/terraform-provider-google/issues/472))
* **New Resource:** `google_kms_key_ring` ([#518](https://github.com/terraform-providers/terraform-provider-google/issues/518))
* **New Resource:** `google_dataproc_cluster` ([#252](https://github.com/terraform-providers/terraform-provider-google/issues/252))
* **New Resource:** `google_project_service` ([#668](https://github.com/terraform-providers/terraform-provider-google/issues/668))

IMPROVEMENTS:
* compute: Add import support for `google_compute_global_forwarding_rule` ([#653](https://github.com/terraform-providers/terraform-provider-google/issues/653))
* compute: Add IAP support for backend services ([#471](https://github.com/terraform-providers/terraform-provider-google/issues/471))
* compute: Allow attaching and detaching disks from instances ([#636](https://github.com/terraform-providers/terraform-provider-google/issues/636))
* compute: Add support for source/target service accounts to `google_compute_firewall` ([#681](https://github.com/terraform-providers/terraform-provider-google/issues/681))
* compute: Add `secondary_ip_range` support to `google_compute_subnetwork` data source ([#687](https://github.com/terraform-providers/terraform-provider-google/issues/687))
* compute: Add support for internal address (beta feature) in `google_compute_address` ([#594](https://github.com/terraform-providers/terraform-provider-google/issues/594))
* compute: Add support to `google_compute_target_pool` for health checks self_link ([#702](https://github.com/terraform-providers/terraform-provider-google/issues/702))
* container: Add support for CPU Platform in `google_container_node_pool` and `google_container_cluster` ([#622](https://github.com/terraform-providers/terraform-provider-google/issues/622))
* container: Add support for Kubernetes alpha features ([#646](https://github.com/terraform-providers/terraform-provider-google/issues/646))
* container: Add support for master authorized networks in `google_container_cluster` ([#626](https://github.com/terraform-providers/terraform-provider-google/issues/626))
* container: Add support for maintenance window on `google_container_cluster` ([#670](https://github.com/terraform-providers/terraform-provider-google/issues/670))
* logging: Make `google_logging_project_sink` resource importable ([#688](https://github.com/terraform-providers/terraform-provider-google/issues/688))
* project: Make `google_service_account` resource importable ([#606](https://github.com/terraform-providers/terraform-provider-google/issues/606))
* project: Project is optional and default to the provider value in `google_project_iam_policy` ([#691](https://github.com/terraform-providers/terraform-provider-google/issues/691))
* pubsub: Create a `google_pubsub_subscription` for a topic in a different project ([#640](https://github.com/terraform-providers/terraform-provider-google/issues/640))
* storage: Add labels to `google_storage_bucket` ([#652](https://github.com/terraform-providers/terraform-provider-google/issues/652))

BUG FIXES:
* compute: Increase timeout for deleting networks ([#662](https://github.com/terraform-providers/terraform-provider-google/issues/662))
* compute: Fix disk migration bug with empty `initialize_params` block ([#664](https://github.com/terraform-providers/terraform-provider-google/issues/664))
* compute: Update `google_compute_target_pool` to no longer have a plan/apply loop with instance URLs ([#666](https://github.com/terraform-providers/terraform-provider-google/issues/666))
* container: `google_container_cluster.node_config.oauth_scopes` no longer need to be set alphabetically ([#506](https://github.com/terraform-providers/terraform-provider-google/issues/506))
* dns: `google_dns_record_set` can now manage NS records ([#359](https://github.com/terraform-providers/terraform-provider-google/issues/359))
* project: Set valid default `public_key_type` for `google_service_account_key` ([#686](https://github.com/terraform-providers/terraform-provider-google/issues/686))

## 1.1.1 (October 24, 2017)

FEATURES:

* **New Resource:** `google_compute_target_ssl_proxy` ([#569](https://github.com/terraform-providers/terraform-provider-google/issues/569))
* **New Data Source:** `google_compute_lb_ip_ranges` ([#567](https://github.com/terraform-providers/terraform-provider-google/issues/567))

IMPROVEMENTS:
* compute: Make `boot_disk` required; remove checks around expected number of disks ([#600](https://github.com/terraform-providers/terraform-provider-google/issues/600))
* compute: Allow setting boot and attached disk sources by name or self link ([#605](https://github.com/terraform-providers/terraform-provider-google/issues/605))
* container: Allow updating `google_container_cluster.monitoring_service` ([#598](https://github.com/terraform-providers/terraform-provider-google/issues/598))
* container: Allow updating `google_container_cluster.addons_config` ([#597](https://github.com/terraform-providers/terraform-provider-google/issues/597))
* project: Make `google_project_services` resource importable ([#601](https://github.com/terraform-providers/terraform-provider-google/issues/601))

BUG FIXES:
* compute: Fix import functionality in `google_compute_route` ([#565](https://github.com/terraform-providers/terraform-provider-google/issues/565))
* compute: Migrate boot disk initialize params ([#592](https://github.com/terraform-providers/terraform-provider-google/issues/592))

## 1.1.0 (October 12, 2017)

FEATURES:
* **New Resource:** `google_logging_folder_sink` ([#470](https://github.com/terraform-providers/terraform-provider-google/pull/470))
* **New Resource:** `google_organization_policy` ([#523](https://github.com/terraform-providers/terraform-provider-google/pull/523))
* **New Resource:** `google_compute_target_tcp_proxy` ([#528](https://github.com/terraform-providers/terraform-provider-google/pull/528))
* **New Resource:** `google_compute_region_autoscaler` ([#544](https://github.com/terraform-providers/terraform-provider-google/pull/544))
* **New Resources:** `google_compute_shared_vpc_host_project` and `google_compute_shared_vpc_service_project` ([#544](https://github.com/terraform-providers/terraform-provider-google/pull/572))

IMPROVEMENTS:
* compute: Generate network link without calling network API in `google_compute_subnetwork` ([#527](https://github.com/terraform-providers/terraform-provider-google/issues/527))
* compute: Generate network link without calling network API in `google_compute_vpn_gateway` and `google_compute_router` ([#527](https://github.com/terraform-providers/terraform-provider-google/issues/527))
* compute: Add import support to `google_compute_target_tcp_proxy` ([#534](https://github.com/terraform-providers/terraform-provider-google/issues/534))
* compute: Add labels support to `google_compute_instance_template` ([#17](https://github.com/terraform-providers/terraform-provider-google/issues/17))
* compute: `google_vpn_tunnel` - Mark 'shared_secret' as sensitive ([#561](https://github.com/terraform-providers/terraform-provider-google/issues/561))
* container: Allow disabling of Kubernetes Dashboard via `kubernetes_dashboard` addon ([#433](https://github.com/terraform-providers/terraform-provider-google/issues/433))
* container: Merge the schemas and logic for the node pool resource and the node pool field in the cluster to aid in maintainability ([#489](https://github.com/terraform-providers/terraform-provider-google/issues/489))
* container: Add master_version to container cluster ([#538](https://github.com/terraform-providers/terraform-provider-google/issues/538))
* sql: Add new retry wrapper fn, retry sql database instance operations that commonly 503 ([#417](https://github.com/terraform-providers/terraform-provider-google/issues/417))
* pubsub: `push_config` field for a `google_pubsub_subscription` is not updateable ([#512](https://github.com/terraform-providers/terraform-provider-google/issues/512))

BUG FIXES:
* compute: Fix bug in `google_compute_instance` preventing the `assigned_nat_ip` field from ever getting assigned ([#536](https://github.com/terraform-providers/terraform-provider-google/issues/536))
* compute: Fix bug in `google_compute_firewall` causing the beta APIs even if no beta features are used ([#500](https://github.com/terraform-providers/terraform-provider-google/issues/500))
* compute: Fix bug in `google_network_peering` preventing creating a peering for a network outside the provider default project ([#496](https://github.com/terraform-providers/terraform-provider-google/issues/496))
* compute: Fix BackendService group hash when instance groups use beta features ([#522](https://github.com/terraform-providers/terraform-provider-google/issues/522))
* compute: Make `disk.device_name` computed in `google_compute_instance_template` ([#566](https://github.com/terraform-providers/terraform-provider-google/issues/566))
* dns: Error out if DNS zone is not found ([#560](https://github.com/terraform-providers/terraform-provider-google/issues/560))
* container: Fix crash when creating node pools with `name_prefix` or no name ([#531](https://github.com/terraform-providers/terraform-provider-google/issues/531))
* container: Fix cluster version upgrades ([#577](https://github.com/terraform-providers/terraform-provider-google/issues/577))

## 1.0.1 (October 02, 2017)

BUG FIXES:
* compute: Fix bug that prevented the state migration for `google_compute_instance` from updating to use attached_disk, boot_disk, and scratch_disk. ([#511](https://github.com/terraform-providers/terraform-provider-google/issues/511))
* compute: Fix bug causing a crash if the API returns an error on `google_compute_instance` creation ([#556](https://github.com/terraform-providers/terraform-provider-google/issues/556))

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
