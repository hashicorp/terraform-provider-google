## 5.45.2 (Feb 10, 2025)

NOTES:
* `5.45.2` contains no changes from `5.45.1`. This release is being made to ensure that the version numbers of the `google` and `google-beta` provider releases remain aligned, as [`google-beta`'s `5.45.2` release](https://github.com/hashicorp/terraform-provider-google-beta/releases/tag/v5.45.2) contains a beta-only change.

## 6.19.0 (Feb 3, 2025)
DEPRECATIONS:
* beyondcorp: deprecated `location` on `google_beyondcorp_security_gateway`. The only valid value is `global`, which is now also the default value. The field will be removed in a future major release. ([#21006](https://github.com/hashicorp/terraform-provider-google/pull/21006))

FEATURES:
* **New Data Source:** `google_parameter_manager_parameter_version` ([#21055](https://github.com/hashicorp/terraform-provider-google/pull/21055))
* **New Data Source:** `google_parameter_manager_parameters` ([#21043](https://github.com/hashicorp/terraform-provider-google/pull/21043))
* **New Data Source:** `google_parameter_manager_regional_parameter_version` ([#21073](https://github.com/hashicorp/terraform-provider-google/pull/21073))
* **New Resource:** `google_beyondcorp_security_gateway_iam_binding` ([#21078](https://github.com/hashicorp/terraform-provider-google/pull/21078))
* **New Resource:** `google_beyondcorp_security_gateway_iam_member` ([#21078](https://github.com/hashicorp/terraform-provider-google/pull/21078))
* **New Resource:** `google_beyondcorp_security_gateway_iam_policy` ([#21078](https://github.com/hashicorp/terraform-provider-google/pull/21078))

IMPROVEMENTS:
* accesscontextmanager: added `etag` to `google_access_context_manager_service_perimeter_dry_run_resource` to prevent overriding list of resources ([#21005](https://github.com/hashicorp/terraform-provider-google/pull/21005))
* compute: allowed parallelization of `google_compute_(region_)per_instance_config` by not locking on the parent resource, but including instance name. ([#21001](https://github.com/hashicorp/terraform-provider-google/pull/21001))
* compute: added `network_profile` field to `google_compute_network` resource. ([#21027](https://github.com/hashicorp/terraform-provider-google/pull/21027))
* compute: added `zero_advertised_route_priority` field to `google_compute_router_peer` ([#21024](https://github.com/hashicorp/terraform-provider-google/pull/21024))
* container: added `max_run_duration` to `node_config` in `google_container_cluster` and `google_container_node_pool` ([#21071](https://github.com/hashicorp/terraform-provider-google/pull/21071))
* dataproc: added `encryption_config` to `google_dataproc_workflow_template` ([#21077](https://github.com/hashicorp/terraform-provider-google/pull/21077))
* gkehub2: added support for `fleet_default_member_config.config_management.config_sync.metrics_gcp_service_account_email` field to `google_gke_hub_feature` resource ([#21042](https://github.com/hashicorp/terraform-provider-google/pull/21042))
* iam: added `prefix` and `regex` fields to `google_service_accounts` data source ([#21020](https://github.com/hashicorp/terraform-provider-google/pull/21020))
* pubsub: added `ingestion_data_source_settings.aws_msk` and `ingestion_data_source_settings.confluent_cloud` fields to `google_pubsub_topic` resource ([#20999](https://github.com/hashicorp/terraform-provider-google/pull/20999))
* spanner: added `encryption_config` field to  `google_spanner_backup_schedule` ([#21067](https://github.com/hashicorp/terraform-provider-google/pull/21067))
* workflows: added `tags` and `workflow_tags` fields to `google_workflows_workflow` resource ([#21053](https://github.com/hashicorp/terraform-provider-google/pull/21053))

BUG FIXES:
* alloydb: marked `google_alloydb_user.password` as sensitive ([#21014](https://github.com/hashicorp/terraform-provider-google/pull/21014))
* beyondcorp: corrected `location` to always be global in `google_beyondcorp_security_gateway` ([#21006](https://github.com/hashicorp/terraform-provider-google/pull/21006))
* cloudquotas: removed validation for `parent` in `google_cloud_quotas_quota_adjuster_settings` ([#21054](https://github.com/hashicorp/terraform-provider-google/pull/21054))
* compute: made `google_compute_router_peer.advertised_route_priority` use server-side default if unset. To set the value to `0` you must also set `zero_advertised_route_priority = true`. ([#21024](https://github.com/hashicorp/terraform-provider-google/pull/21024))
* container: fixed a diff caused by server-side set values for `node_config.resource_labels` ([#21082](https://github.com/hashicorp/terraform-provider-google/pull/21082))
* container: marked `cluster_autoscaling.resource_limits.maximum` as required, as requests would fail if it was not set ([#21051](https://github.com/hashicorp/terraform-provider-google/pull/21051))
* firestore: fixed error preventing deletion of wildcard `google_firestore_field` resources ([#21034](https://github.com/hashicorp/terraform-provider-google/pull/21034))
* netapp: fixed an issue where a diff on `zone` would be found if it was unspecified in `google_netapp_storage_pool` ([#21060](https://github.com/hashicorp/terraform-provider-google/pull/21060))
* networksecurity: fixed sporadic-diff in `google_network_security_security_profile` ([#21070](https://github.com/hashicorp/terraform-provider-google/pull/21070))
* spanner: fixed bug with `google_spanner_instance.force_destroy` not setting `billing_project` value correctly ([#21023](https://github.com/hashicorp/terraform-provider-google/pull/21023))
* storage: fixed an issue where plans with a dependency on the `content` field in the `google_storage_bucket_object_content` data source could erroneously fail ([#21074](https://github.com/hashicorp/terraform-provider-google/pull/21074))

## 5.45.1 (January 29, 2025)

NOTES:
* 5.45.1 is a backport release, responding to a new GKE label being applied that can cause unwanted diffs in node pools. The changes in this release will be available in 6.18.1 and users upgrading to 6.X should upgrade to that version or higher.

BUG FIXES:
* container: fixed a diff caused by server-side set values for `node_config.resource_labels` ([#21082](https://github.com/hashicorp/terraform-provider-google/pull/21082))

## 5.45.0 (November 11, 2024)

NOTES:
* 5.45.0 is a backport release, responding to a new Spanner feature that may result in creation of unwanted backups for users. The changes in this release will be available in 6.11.0 and users upgrading to 6.X should upgrade to that version or higher.

IMPROVEMENTS:
* spanner: added `default_backup_schedule_type` field to  `google_spanner_instance` ([#20213](https://github.com/hashicorp/terraform-provider-google/pull/20213))

## 5.44.2 (October 14, 2024)

Notes:
* 5.44.2 is a backport release, responding to a GKE rollout that created permadiffs for many users. The changes in this release will be available in 6.7.0 and users upgrading to 6.X should upgrade to that version or higher.

IMPROVEMENTS:
* container: `google_container_cluster` will now accept server-specified values for `node_pool_auto_config.0.node_kubelet_config` when it is not defined in configuration and will not detect drift. Note that this means that removing the value from configuration will now preserve old settings instead of reverting the old settings. ([#19817](https://github.com/hashicorp/terraform-provider-google/pull/19817))

BUG FIXES:
* container: fixed a diff triggered by a new API-side default value for `node_config.0.kubelet_config.0.insecure_kubelet_readonly_port_enabled`. Terraform will now accept server-specified values for `node_config.0.kubelet_config` when it is not defined in configuration and will not detect drift. Note that this means that removing the value from configuration will now preserve old settings instead of reverting the old settings. ([#19817](https://github.com/hashicorp/terraform-provider-google/pull/19817))

## 5.44.1 (September 23, 2024)
NOTES:
* 5.44.1 is a backport release, intended to pull in critical container improvements and fixes for issues introduced in 5.44.0

IMPROVEMENTS:
* container: added in-place update support for `gcfs_config` in in `google_container_cluster` and `google_container_node_pool` ([#19365](https://github.com/hashicorp/terraform-provider-google/pull/19365)) ([#19512](https://github.com/hashicorp/terraform-provider-google/pull/19512))

BUG FIXES:
* container: fixed a permadiff on `gcfs_config` in `google_container_cluster` and `google_container_node_pool` ([#19512](https://github.com/hashicorp/terraform-provider-google/pull/19512))
* container: fixed a bug where specifying `node_pool_defaults.node_config_defaults` with `enable_autopilot = true` will cause `google_container_cluster` resource creation failure. ([#19543](https://github.com/hashicorp/terraform-provider-google/pull/19543))

## 5.44.0 (September 9, 2024)

NOTES:
* 5.44.0 is a backport release, intended to pull in critical container improvements from 6.2.0

IMPROVEMENTS:
* container: added `insecure_kubelet_readonly_port_enabled` to `node_pool.node_config.kubelet_config` and `node_config.kubelet_config` in `google_container_node_pool` resource. ([#19312](https://github.com/hashicorp/terraform-provider-google/pull/19312))
* container: added `insecure_kubelet_readonly_port_enabled` to `node_pool_defaults.node_config_defaults`, `node_pool.node_config.kubelet_config`, and `node_config.kubelet_config` in `google_container_cluster` resource. ([#19312](https://github.com/hashicorp/terraform-provider-google/pull/19312))
* container: added `node_pool_auto_config.node_kublet_config.insecure_kubelet_readonly_port_enabled` field to `google_container_cluster`. ([#19320](https://github.com/hashicorp/terraform-provider-google/pull/19320))

## 5.43.1 (August 30, 2024)

NOTES:
* 5.43.1 is a backport release, and some changes will not appear in 6.X series releases until 6.1.0

BUG FIXES:
* pubsub: fixed a validation bug that didn't allow empty filter definitions for `google_pubsub_subscription` resources ([#19284](https://github.com/hashicorp/terraform-provider-google/pull/19284))

## 5.43.0 (August 26, 2024)

DEPRECATIONS:
* storage: deprecated `lifecycle_rule.condition.no_age` field in `google_storage_bucket`. Use the new `lifecycle_rule.condition.send_age_if_zero` field instead. ([#19172](https://github.com/hashicorp/terraform-provider-google/pull/19172))

FEATURES:
* **New Resource:** `google_kms_ekm_connection_iam_binding` ([#19132](https://github.com/hashicorp/terraform-provider-google/pull/19132))
* **New Resource:** `google_kms_ekm_connection_iam_member` ([#19132](https://github.com/hashicorp/terraform-provider-google/pull/19132))
* **New Resource:** `google_kms_ekm_connection_iam_policy` ([#19132](https://github.com/hashicorp/terraform-provider-google/pull/19132))
* **New Resource:** `google_scc_v2_organization_scc_big_query_exports` ([#19184](https://github.com/hashicorp/terraform-provider-google/pull/19184))

IMPROVEMENTS:
* compute: added `label_fingerprint` field to `google_compute_global_address` resource ([#19204](https://github.com/hashicorp/terraform-provider-google/pull/19204))
* compute: exposed service side id as new output field `forwarding_rule_id` on resource `google_compute_forwarding_rule` ([#19139](https://github.com/hashicorp/terraform-provider-google/pull/19139))
* container: added EXTENDED as a valid option for `release_channel` field in `google_container_cluster` resource ([#19141](https://github.com/hashicorp/terraform-provider-google/pull/19141))
* logging: changed `enable_analytics` parsing to "no preference" in analytics if omitted, instead of explicitly disabling analytics in `google_logging_project_bucket_config` ([#19126](https://github.com/hashicorp/terraform-provider-google/pull/19126))
* pusbub: added validation to `filter` field in resource `google_pubsub_subscription` ([#19131](https://github.com/hashicorp/terraform-provider-google/pull/19131))
* resourcemanager: added `default_labels` field to `google_client_config` data source ([#19170](https://github.com/hashicorp/terraform-provider-google/pull/19170))
* vmwareengine: added PC undelete support in `google_vmwareengine_private_cloud` ([#19192](https://github.com/hashicorp/terraform-provider-google/pull/19192))

BUG FIXES:
* alloydb: fixed a permadiff on `psc_instance_config` in `google_alloydb_instance` resource ([#19143](https://github.com/hashicorp/terraform-provider-google/pull/19143))
* compute: fixed a malformed URL that affected updating the `server_tls_policy` property on `google_compute_target_https_proxy` resources ([#19164](https://github.com/hashicorp/terraform-provider-google/pull/19164))
* compute: fixed bug where the `labels` field could not be updated on `google_compute_global_address` ([#19204](https://github.com/hashicorp/terraform-provider-google/pull/19204))
* compute: fixed force diff replacement logic for `network_ip` on resource `google_compute_instance` ([#19135](https://github.com/hashicorp/terraform-provider-google/pull/19135))

## 5.42.0 (August 19, 2024)
DEPRECATIONS:
* compute: setting `google_compute_subnetwork.secondary_ip_range = []` to explicitly set a list of empty objects is deprecated and will produce an error in the upcoming major release. Use `send_secondary_ip_range_if_empty` while removing `secondary_ip_range` from config instead. ([#19122](https://github.com/hashicorp/terraform-provider-google/pull/19122))

FEATURES:
* **New Data Source:** `google_artifact_registry_locations` ([#19047](https://github.com/hashicorp/terraform-provider-google/pull/19047))
* **New Data Source:** `google_cloud_identity_transitive_group_memberships` ([#19038](https://github.com/hashicorp/terraform-provider-google/pull/19038))
* **New Resource:** `google_discovery_engine_schema` ([#19124](https://github.com/hashicorp/terraform-provider-google/pull/19124))
* **New Resource:** `google_scc_folder_notification_config` ([#19057](https://github.com/hashicorp/terraform-provider-google/pull/19057))
* **New Resource:** `google_scc_v2_folder_notification_config` ([#19055](https://github.com/hashicorp/terraform-provider-google/pull/19055))
* **New Resource:** `google_vertex_ai_index_endpoint_deployed_index` ([#19061](https://github.com/hashicorp/terraform-provider-google/pull/19061))

IMPROVEMENTS:
* clouddeploy: added `serial_pipeline.stages.strategy.canary.runtime_config.kubernetes.gateway_service_mesh.pod_selector_label` and `serial_pipeline.stages.strategy.canary.runtime_config.kubernetes.service_networking.pod_selector_label` fields to `google_clouddeploy_delivery_pipeline` resource ([#19100](https://github.com/hashicorp/terraform-provider-google/pull/19100))
* compute: added `send_secondary_ip_range_if_empty` to `google_compute_subnetwork` ([#19122](https://github.com/hashicorp/terraform-provider-google/pull/19122))
* discoveryengine: added `skip_default_schema_creation` field to `google_data_store` resource ([#19017](https://github.com/hashicorp/terraform-provider-google/pull/19017))
* dns: changed `load_balancer_type` field from required to optional in `google_dns_record_set` ([#19050](https://github.com/hashicorp/terraform-provider-google/pull/19050))
* firestore: added `cmek_config` field to `google_firestore_database` resource ([#19107](https://github.com/hashicorp/terraform-provider-google/pull/19107))
* servicenetworking: added `update_on_creation_fail` field to `google_service_networking_connection` resource. When it is set to true, enforce an update of the reserved peering ranges on the existing service networking connection in case of a new connection creation failure. ([#19035](https://github.com/hashicorp/terraform-provider-google/pull/19035))
* sql: added `server_ca_mode` field to `google_sql_database_instance` resource ([#18998](https://github.com/hashicorp/terraform-provider-google/pull/18998))

BUG FIXES:
* bigquery: made `google_bigquery_dataset_iam_member` non-authoritative. To remove a bigquery dataset iam member, use an authoritative resource like `google_bigquery_dataset_iam_policy` ([#19121](https://github.com/hashicorp/terraform-provider-google/pull/19121))
* cloudfunctions2: fixed a "Provider produced inconsistent final plan" bug affecting the `service_config.environment_variables` field in `google_cloudfunctions2_function` resource ([#19024](https://github.com/hashicorp/terraform-provider-google/pull/19024))
* cloudfunctions2: fixed a permadiff on `storage_source.generation` in `google_cloudfunctions2_function` resource ([#19031](https://github.com/hashicorp/terraform-provider-google/pull/19031))
* compute: fixed issue where sub-resources managed by `google_compute_forwarding_rule` prevented resource deletion ([#19117](https://github.com/hashicorp/terraform-provider-google/pull/19117))
* logging: changed `google_logging_project_bucket_config.enable_analytics` behavior to set "no preference" in analytics if omitted, instead of explicitly disabling analytics. ([#19126](https://github.com/hashicorp/terraform-provider-google/pull/19126))
* workbench: fixed a bug with `google_workbench_instance` metadata drifting when using custom containers. ([#19119](https://github.com/hashicorp/terraform-provider-google/pull/19119))

## 5.41.0 (August 13, 2024)

DEPRECATIONS:
* resourcemanager: deprecated `skip_delete` field in the `google_project` resource. Use `deletion_policy` instead. ([#18867](https://github.com/hashicorp/terraform-provider-google/pull/18867))

FEATURES:
* **New Data Source:** `google_logging_log_view_iam_policy` ([#18990](https://github.com/hashicorp/terraform-provider-google/pull/18990))
* **New Data Source:** `google_scc_v2_organization_source_iam_policy` ([#19004](https://github.com/hashicorp/terraform-provider-google/pull/19004))
* **New Resource:** `google_access_context_manager_service_perimeter_dry_run_egress_policy` ([#18994](https://github.com/hashicorp/terraform-provider-google/pull/18994))
* **New Resource:** `google_access_context_manager_service_perimeter_dry_run_ingress_policy` ([#18994](https://github.com/hashicorp/terraform-provider-google/pull/18994))
* **New Resource:** `google_scc_v2_folder_mute_config` ([#18924](https://github.com/hashicorp/terraform-provider-google/pull/18924))
* **New Resource:** `google_scc_v2_project_mute_config` ([#18993](https://github.com/hashicorp/terraform-provider-google/pull/18993))
* **New Resource:** `google_scc_v2_project_notification_config` ([#19008](https://github.com/hashicorp/terraform-provider-google/pull/19008))
* **New Resource:** `google_scc_v2_organization_source` ([#19004](https://github.com/hashicorp/terraform-provider-google/pull/19004))
* **New Resource:** `google_scc_v2_organization_source_iam_binding` ([#19004](https://github.com/hashicorp/terraform-provider-google/pull/19004))
* **New Resource:** `google_scc_v2_organization_source_iam_member` ([#19004](https://github.com/hashicorp/terraform-provider-google/pull/19004))
* **New Resource:** `google_scc_v2_organization_source_iam_policy` ([#19004](https://github.com/hashicorp/terraform-provider-google/pull/19004))
* **New Resource:** `google_logging_log_view_iam_binding` ([#18990](https://github.com/hashicorp/terraform-provider-google/pull/18990))
* **New Resource:** `google_logging_log_view_iam_member` ([#18990](https://github.com/hashicorp/terraform-provider-google/pull/18990))
* **New Resource:** `google_logging_log_view_iam_policy` ([#18990](https://github.com/hashicorp/terraform-provider-google/pull/18990))

IMPROVEMENTS:
* clouddeploy: added `gke.proxy_url` field to `google_clouddeploy_target` ([#19016](https://github.com/hashicorp/terraform-provider-google/pull/19016))
* cloudrunv2: added field `binary_authorization.policy` to resource `google_cloud_run_v2_job` and resource `google_cloud_run_v2_service` to support named binary authorization policy. ([#18995](https://github.com/hashicorp/terraform-provider-google/pull/18995))
* compute: added `source_regions` field to `google_compute_healthcheck` resource ([#19006](https://github.com/hashicorp/terraform-provider-google/pull/19006))
* compute: added update-in-place support for the `google_compute_target_https_proxy.server_tls_policy` field ([#18996](https://github.com/hashicorp/terraform-provider-google/pull/18996))
* compute: added update-in-place support for the `google_compute_region_target_https_proxy.server_tls_policy` field ([#19007](https://github.com/hashicorp/terraform-provider-google/pull/19007))
* container: added `auto_provisioning_locations` field to `google_container_cluster` ([#18928](https://github.com/hashicorp/terraform-provider-google/pull/18928))
* dataform: added `kms_key_name` field to `google_dataform_repository` resource ([#18947](https://github.com/hashicorp/terraform-provider-google/pull/18947))
* discoveryengine: added `skip_default_schema_creation` field to `google_discovery_engine_data_store` resource ([#19017](https://github.com/hashicorp/terraform-provider-google/pull/19017))
* gkehub: added `configmanagement.management` and `configmanagement.config_sync.enabled` fields to `google_gkehub_feature_membership` ([#19016](https://github.com/hashicorp/terraform-provider-google/pull/19016))
* gkehub: added `management` field to `google_gke_hub_feature.fleet_default_member_config.configmanagement` ([#18963](https://github.com/hashicorp/terraform-provider-google/pull/18963))
* resourcemanager: added `deletion_policy` field to the `google_project` resource. Setting `deletion_policy` to `PREVENT` will protect the project against any destroy actions caused by a terraform apply or terraform destroy. Setting `deletion_policy` to `ABANDON` allows the resource to be abandoned rather than deleted and it behaves the same with `skip_delete = true`. Default value is `DELETE`. `skip_delete = true` takes precedence over `deletion_policy = "DELETE"`.
* storage: added `force_destroy` field to `google_storage_managed_folder` resource ([#18973](https://github.com/hashicorp/terraform-provider-google/pull/18973))
* storage: added `generation` field to `google_storage_bucket_object` resource ([#18971](https://github.com/hashicorp/terraform-provider-google/pull/18971))

BUG FIXES:
* compute: fixed `google_compute_instance.alias_ip_range` update behavior to avoid temporarily deleting unchanged alias IP ranges ([#19015](https://github.com/hashicorp/terraform-provider-google/pull/19015))
* compute: fixed the bug that creation of PSC forwarding rules fails in `google_compute_forwarding_rule` resource when provider default labels are set ([#18984](https://github.com/hashicorp/terraform-provider-google/pull/18984))
* sql: fixed a perma-diff in `settings.insights_config` in `google_sql_database_instance` ([#18962](https://github.com/hashicorp/terraform-provider-google/pull/18962))




## 5.40.0 (August 5, 2024)

IMPROVEMENTS:
* bigquery: added support for value `DELTA_LAKE` to `source_format` in `google_bigquery_table` resource ([#18915](https://github.com/hashicorp/terraform-provider-google/pull/18915))
* compute: added `access_mode` field to `google_compute_disk` resource ([#18857](https://github.com/hashicorp/terraform-provider-google/pull/18857))
* compute: added `stack_type`, and `gateway_ip_version` fields to `google_compute_router` resource ([#18839](https://github.com/hashicorp/terraform-provider-google/pull/18839))
* container: added field `ray_operator_config` for `resource_container_cluster` ([#18825](https://github.com/hashicorp/terraform-provider-google/pull/18825))
* container: promoted `additional_node_network_configs` and `additional_pod_network_configs` fields to GA in the `google_container_node_pool` resource ([#18842](https://github.com/hashicorp/terraform-provider-google/pull/18842))
* container: promoted `enable_multi_networking` to GA in the `google_container_cluster` resource ([#18842](https://github.com/hashicorp/terraform-provider-google/pull/18842))
* monitoring: updated `goal` field to accept a max threshold of up to 0.9999 in `google_monitoring_slo` resource to 0.9999 ([#18845](https://github.com/hashicorp/terraform-provider-google/pull/18845))
* networkconnectivity: added `export_psc` field to `google_network_connectivity_hub` resource ([#18866](https://github.com/hashicorp/terraform-provider-google/pull/18866))
* sql: added `enable_dataplex_integration` field to `google_sql_database_instance` resource ([#18852](https://github.com/hashicorp/terraform-provider-google/pull/18852))

BUG FIXES:
* bigquery: fixed a permadiff when handling "assets" in `params` in the `google_bigquery_data_transfer_config` resource ([#18898](https://github.com/hashicorp/terraform-provider-google/pull/18898))
* bigquery: fixed an issue preventing certain keys in `params` from being assigned values in `google_bigquery_data_transfer_config` ([#18888](https://github.com/hashicorp/terraform-provider-google/pull/18888))
* compute: fixed perma-diff of `advertised_ip_ranges` field in `google_compute_router` resource ([#18869](https://github.com/hashicorp/terraform-provider-google/pull/18869))
* container: fixed perma-diff on `node_config.guest_accelerator.gpu_driver_installation_config` field in GKE 1.30+ in `google_container_node_pool` resource ([#18835](https://github.com/hashicorp/terraform-provider-google/pull/18835))
* sql: fixed a perma-diff in `settings.insights_config` in `google_sql_database_instance` ([#18962](https://github.com/hashicorp/terraform-provider-google/pull/18962))

## v5.39.1 (July 30th, 2024)

BUG FIXES:
* datastream: fixed a breaking change in 5.39.0 `google_datastream_stream` that made one of `destination_config.bigquery_destination_config.merge` or `destination_config.bigquery_destination_config.append_only` required ([#18903](https://github.com/hashicorp/terraform-provider-google/pull/18903))

## 5.39.0 (July 29th, 2024)

NOTES:
* networkconnectivity: migrated `google_network_connectivity_hub` from DCL to MMv1 ([#18724](https://github.com/hashicorp/terraform-provider-google/pull/18724))
* networkconnectivity: migrated `google_network_connectivity_spoke` from DCL to MMv1 ([#18779](https://github.com/hashicorp/terraform-provider-google/pull/18779))

DEPRECATIONS:
* bigquery: deprecated `allow_resource_tags_on_deletion` in `google_bigquery_table`. ([#18811](https://github.com/hashicorp/terraform-provider-google/pull/18811))
* bigqueryreservation: deprecated `multi_region_auxiliary` on `google_bigquery_reservation`. ([#18803](https://github.com/hashicorp/terraform-provider-google/pull/18803))
* datastore: deprecated the resource `google_datastore_index`. Use the `google_firestore_index` resource instead. ([#18781](https://github.com/hashicorp/terraform-provider-google/pull/18781))

FEATURES:
* **New Resource:** `google_apigee_environment_keyvaluemaps_entries` ([#18707](https://github.com/hashicorp/terraform-provider-google/pull/18707))
* **New Resource:** `google_apigee_environment_keyvaluemaps` ([#18707](https://github.com/hashicorp/terraform-provider-google/pull/18707))
* **New Resource:** `google_compute_resize_request` ([#18725](https://github.com/hashicorp/terraform-provider-google/pull/18725))
* **New Resource:** `google_compute_router_route_policy` ([#18759](https://github.com/hashicorp/terraform-provider-google/pull/18759))
* **New Resource:** `google_scc_v2_organization_mute_config` ([#18752](https://github.com/hashicorp/terraform-provider-google/pull/18752))

IMPROVEMENTS:
* alloydb: added `observability_config` field to `google_alloydb_instance` resource ([#18743](https://github.com/hashicorp/terraform-provider-google/pull/18743))
* bigquery: added `resource_tags` field to `google_bigquery_dataset` resource (ga) ([#18711](https://github.com/hashicorp/terraform-provider-google/pull/18711))
* bigquery: added `resource_tags` field to `google_bigquery_table` resource ([#18741](https://github.com/hashicorp/terraform-provider-google/pull/18741))
* bigtable: added `data_boost_isolation_read_only` and `data_boost_isolation_read_only.compute_billing_owner` fields to `google_bigtable_app_profile` resource ([#18819](https://github.com/hashicorp/terraform-provider-google/pull/18819))
* cloudfunctions: added `build_service_account` field to `google_cloudfunctions_function` resource ([#18702](https://github.com/hashicorp/terraform-provider-google/pull/18702))
* compute: added `aws_v4_authentication` fields to `google_compute_backend_service` resource ([#18796](https://github.com/hashicorp/terraform-provider-google/pull/18796))
* compute: added `custom_learned_ip_ranges` and `custom_learned_route_priority` fields to `google_compute_router_peer` resource ([#18727](https://github.com/hashicorp/terraform-provider-google/pull/18727))
* compute: added `export_policies` and `import_policies` fields  to `google_compute_router_peer` resource ([#18759](https://github.com/hashicorp/terraform-provider-google/pull/18759))
* compute: added `shared_secret` field to `google_compute_public_advertised_prefix` resource ([#18786](https://github.com/hashicorp/terraform-provider-google/pull/18786))
* compute: added `storage_pool` under `boot_disk.initialize_params` to `google_compute_instance` resource ([#18817](https://github.com/hashicorp/terraform-provider-google/pull/18817))
* compute: changed `target_service` field on the `google_compute_service_attachment` resource to accept a `ForwardingRule` or `Gateway` URL. ([#18742](https://github.com/hashicorp/terraform-provider-google/pull/18742))
* container: added field `ray_operator_config` for `google_container_cluster` ([#18825](https://github.com/hashicorp/terraform-provider-google/pull/18825))
* datastream: added `merge` and `append_only` fields to `google_datastream_stream` resource ([#18726](https://github.com/hashicorp/terraform-provider-google/pull/18726))
* datastream: promoted `source_config.sql_server_source_config` and `backfill_all.sql_server_excluded_objects` fields in `google_datastream_stream` resource from beta to GA ([#18732](https://github.com/hashicorp/terraform-provider-google/pull/18732))
* datastream: promoted `sql_server_profile` field in `google_datastream_connection_profile` resource from beta to GA ([#18732](https://github.com/hashicorp/terraform-provider-google/pull/18732))
* dlp: added `cloud_storage_target` field to `google_data_loss_prevention_discovery_config` resource ([#18740](https://github.com/hashicorp/terraform-provider-google/pull/18740))
* resourcemanager: added `check_if_service_has_usage_on_destroy` field to `google_project_service` resource ([#18753](https://github.com/hashicorp/terraform-provider-google/pull/18753))
* resourcemanager: added the `member` property to `google_project_service_identity` ([#18695](https://github.com/hashicorp/terraform-provider-google/pull/18695))
* vmwareengine: added `deletion_delay_hours` field to `google_vmwareengine_private_cloud` resource ([#18698](https://github.com/hashicorp/terraform-provider-google/pull/18698))
* vmwareengine: supported type change from `TIME_LIMITED` to `STANDARD` for multi-node `google_vmwareengine_private_cloud` resource ([#18698](https://github.com/hashicorp/terraform-provider-google/pull/18698))
* workbench: added `access_configs` to `google_workbench_instance` resource ([#18737](https://github.com/hashicorp/terraform-provider-google/pull/18737))

BUG FIXES:
* compute: fixed perma-diff for `interconnect_type` being `DEDICATED` in `google_compute_interconnect` resource ([#18761](https://github.com/hashicorp/terraform-provider-google/pull/18761))
* dialogflowcx: fixed intermittent issues with retrieving resource state soon after creating `google_dialogflow_cx_security_settings` resources ([#18792](https://github.com/hashicorp/terraform-provider-google/pull/18792))
* firestore: fixed missing import of `field` for `google_firestore_field`. ([#18771](https://github.com/hashicorp/terraform-provider-google/pull/18771))
* firestore: fixed bug where fields `database`, `collection`, `document_id`, and `field` could not be updated on `google_firestore_document` and `google_firestore_field` resources. ([#18821](https://github.com/hashicorp/terraform-provider-google/pull/18821))
* netapp: made the `smb_settings` field on the `google_netapp_volume` resource default to the value returned from the API. This solves permadiffs when the field is unset. ([#18790](https://github.com/hashicorp/terraform-provider-google/pull/18790))
* networksecurity: added recreate functionality on update for `client_validation_mode` and `client_validation_trust_config` in `google_network_security_server_tls_policy` ([#18769](https://github.com/hashicorp/terraform-provider-google/pull/18769))

## 5.38.0 (July 15, 2024)

FEATURES:
* **New Data Source:** `google_gke_hub_membership_binding` ([#18680](https://github.com/hashicorp/terraform-provider-google/pull/18680))
* **New Data Source:** `google_site_verification_token` ([#18688](https://github.com/hashicorp/terraform-provider-google/pull/18688))
* **New Resource:** `google_scc_project_notification_config` ([#18682](https://github.com/hashicorp/terraform-provider-google/pull/18682))

IMPROVEMENTS:
* compute: promoted `labels` field on `google_compute_global_address` resource from beta to GA ([#18646](https://github.com/hashicorp/terraform-provider-google/pull/18646))
* compute: made the `google_compute_resource_policy` resource updatable in-place ([#18673](https://github.com/hashicorp/terraform-provider-google/pull/18673))
* privilegedaccessmanager: promoted `google_privileged_access_manager_entitlement` resource from beta to GA ([#18686](https://github.com/hashicorp/terraform-provider-google/pull/18686))
* vertexai: added `project_number` field to `google_vertex_ai_feature_online_store_featureview` resource ([#18637](https://github.com/hashicorp/terraform-provider-google/pull/18637))

BUG FIXES:
* cloudfunctions2: fixed permadiffs on `service_config.environment_variables` field in `google_cloudfunctions2_function` resource ([#18651](https://github.com/hashicorp/terraform-provider-google/pull/18651))

## 5.37.0 (July 8, 2024)

FEATURES:
* **New Data Source:** `google_kms_crypto_keys` ([#18605](https://github.com/hashicorp/terraform-provider-google/pull/18605))
* **New Data Source:** `google_kms_key_rings` ([#18611](https://github.com/hashicorp/terraform-provider-google/pull/18611))
* **New Resource:** `google_scc_v2_organization_notification_config` ([#18594](https://github.com/hashicorp/terraform-provider-google/pull/18594))
* **New Resource:** `google_secure_source_manager_repository` ([#18576](https://github.com/hashicorp/terraform-provider-google/pull/18576))
* **New Resource:** `google_storage_managed_folder_iam` ([#18555](https://github.com/hashicorp/terraform-provider-google/pull/18555))
* **New Resource:** `google_storage_managed_folder` ([#18555](https://github.com/hashicorp/terraform-provider-google/pull/18555))

IMPROVEMENTS:
* certificatemanager: added `allowlisted_certificates` field to `google_certificate_manager_trust_config` resource ([#18587](https://github.com/hashicorp/terraform-provider-google/pull/18587))
* compute: added `max_run_duration` and `on_instance_stop_action` fields to `google_compute_instance`, `google_compute_instance_template`, and `google_compute_instance_from_machine_image` resources ([#18623](https://github.com/hashicorp/terraform-provider-google/pull/18623))
* dataplex: added `sql_assertion` field to `google_dataplex_datascan` resource ([#18559](https://github.com/hashicorp/terraform-provider-google/pull/18559))
* gkehub: added `fleet_default_member_config.configmanagement.config_sync.enabled` field to `google_gke_hub_feature` resource ([#18582](https://github.com/hashicorp/terraform-provider-google/pull/18582))
* netapp: added `zone` and `replica_zone` field to `google_netapp_storage_pool` resource ([#18609](https://github.com/hashicorp/terraform-provider-google/pull/18609))
* vertexai: added `project_number` field to `google_vertex_ai_feature_online_store_featureview` resource ([#18637](https://github.com/hashicorp/terraform-provider-google/pull/18637))
* workstations: added `host.gce_instance.vm_tags` field to `google_workstations_workstation_config` resource ([#18588](https://github.com/hashicorp/terraform-provider-google/pull/18588))

BUG FIXES:
* compute: fixed a bug preventing the creation of `google_compute_autoscaler` and `google_compute_region_autoscaler` resources if both `autoscaling_policy.max_replicas` and `autoscaling_policy.min_replicas` were configured as zero. ([#18607](https://github.com/hashicorp/terraform-provider-google/pull/18607))
* resourcemanager: mitigated eventual consistency issues by adding a 10s wait after `google_service_account_key` resource creation ([#18566](https://github.com/hashicorp/terraform-provider-google/pull/18566))
* vertexai: fixed issue where updating "metadata" field could fail in `google_vertex_ai_index` resource ([#18632](https://github.com/hashicorp/terraform-provider-google/pull/18632))

## 5.36.0 (July 1, 2024)

FEATURES:
* **New Resource:** `google_storage_managed_folder_iam` ([#18555](https://github.com/hashicorp/terraform-provider-google/pull/18555))
* **New Resource:** `google_storage_managed_folder` ([#18555](https://github.com/hashicorp/terraform-provider-google/pull/18555))

IMPROVEMENTS:
* bigtable: added `ignore_warnings` field to `google_bigtable_gc_policy` resource ([#18492](https://github.com/hashicorp/terraform-provider-google/pull/18492))
* cloudfunctions2: added `build_config.automatic_update_policy` and `build_config.on_deploy_update_policy` fields to `google_cloudfunctions2_function` resource ([#18540](https://github.com/hashicorp/terraform-provider-google/pull/18540))
* compute: added `confidential_instance_config.confidential_instance_type` field to `google_compute_instance`, `google_compute_instance_template`, and `google_compute_region_instance_template` resources ([#18554](https://github.com/hashicorp/terraform-provider-google/pull/18554))
* compute: added `custom_error_response_policy` and `default_custom_error_response_policy` fields to `google_compute_url_map` resource ([#18511](https://github.com/hashicorp/terraform-provider-google/pull/18511))
* compute: added `tls_early_data` field to `google_compute_target_https_proxy` resource ([#18512](https://github.com/hashicorp/terraform-provider-google/pull/18512))
* compute: promoted `google_compute_network_attachment` resource from beta to GA ([#18494](https://github.com/hashicorp/terraform-provider-google/pull/18494))
* datafusion: added `connection_type` and `private_service_connect_config` fields to `google_data_fusion_instance` resource ([#18525](https://github.com/hashicorp/terraform-provider-google/pull/18525))
* healthcare: added `encryption_spec` field to `google_healthcare_dataset` resource ([#18528](https://github.com/hashicorp/terraform-provider-google/pull/18528))
* monitoring: added `links` field to `google_monitoring_alert_policy` resource ([#18549](https://github.com/hashicorp/terraform-provider-google/pull/18549))
* vertexai: added update support for `big_query.entity_id_columns` field on `google_vertex_ai_feature_group` resource ([#18493](https://github.com/hashicorp/terraform-provider-google/pull/18493))
* vertexai: promoted `dedicated_serving_endpoint` field on `google_vertex_ai_feature_online_store` resource from beta to GA ([#18513](https://github.com/hashicorp/terraform-provider-google/pull/18513))

BUG FIXES:
* accesscontextmanager: fixed perma-diff caused by ordering of `service_perimeters` in `google_access_context_manager_service_perimeters` resource ([#18520](https://github.com/hashicorp/terraform-provider-google/pull/18520))
* compute: fixed a crash in `google_compute_reservation` resource when `share_settings` field has changes ([#18498](https://github.com/hashicorp/terraform-provider-google/pull/18498))
* compute: fixed issue in `google_compute_instance` resource where `service_account` is not set when specifying `service_account.email` and no `service_account.scopes` ([#18521](https://github.com/hashicorp/terraform-provider-google/pull/18521))
* gkehub2: fixed `google_gke_hub_feature` resource to allow `fleet_default_member_config` field to be unset ([#18487](https://github.com/hashicorp/terraform-provider-google/pull/18487))
* identityplatform: fixed perma-diff on `google_identity_platform_config` resource when `sms_region_config` is not set ([#18537](https://github.com/hashicorp/terraform-provider-google/pull/18537))
* logging: fixed perma-diff on `index_configs` in `google_logging_organization_bucket_config` resource ([#18501](https://github.com/hashicorp/terraform-provider-google/pull/18501))

## 5.35.0 (June 24, 2024)

FEATURES:
* **New Data Source:** `google_artifact_registry_docker_image` ([#18446](https://github.com/hashicorp/terraform-provider-google/pull/18446))
* **New Resource:** `google_service_networking_vpc_service_controls` ([#18448](https://github.com/hashicorp/terraform-provider-google/pull/18448))

IMPROVEMENTS:
* billingbudget: added `enable_project_level_recipients` field to `google_billing_budget` resource ([#18437](https://github.com/hashicorp/terraform-provider-google/pull/18437))
* compute: added `action_token_site_keys` and `session_token_site_keys` fields to `google_compute_security_policy` and `google_compute_security_policy_rule` resources ([#18414](https://github.com/hashicorp/terraform-provider-google/pull/18414))
* gkehub2: added `ENTERPRISE` option to `security_posture_config` field on `google_gke_hub_fleet` resource ([#18440](https://github.com/hashicorp/terraform-provider-google/pull/18440))
* pubsub: added `bigquery_config.service_account_email` field to `google_pubsub_subscription` resource ([#18444](https://github.com/hashicorp/terraform-provider-google/pull/18444))
* redis: added `maintenance_version` field to `google_redis_instance` resource ([#18424](https://github.com/hashicorp/terraform-provider-google/pull/18424))
* storage: changed update behavior in `google_storage_bucket_object` to no longer delete to avoid object deletion on content update ([#18479](https://github.com/hashicorp/terraform-provider-google/pull/18479))
* sql: added support for more MySQL values in `type` field of `google_sql_user` resource ([#18452](https://github.com/hashicorp/terraform-provider-google/pull/18452))
* sql: increased timeouts on `google_sql_database_instance` to 90m to account for longer-running actions such as creation through cloning ([#18458](https://github.com/hashicorp/terraform-provider-google/pull/18458))
* workbench: added update support to `gce_setup.boot_disk` and `gce_setup.data_disks` fields in `google_workbench_instance` resource ([#18482](https://github.com/hashicorp/terraform-provider-google/pull/18482))

BUG FIXES:
* compute: updated `google_compute_instance` to force reboot if `min_node_cpus` is updated ([#18420](https://github.com/hashicorp/terraform-provider-google/pull/18420))
* compute: fixed `description` field in `google_compute_firewall` to support empty/null values on update ([#18478](https://github.com/hashicorp/terraform-provider-google/pull/18478))
* compute: fixed perma-diff on `google_compute_disk` for Ubuntu amd64 canonical LTS images ([#18418](https://github.com/hashicorp/terraform-provider-google/pull/18418))
* storage: fixed lowercased `custom_placement_config` values in `google_storage_bucket` causing perma-destroy ([#18456](https://github.com/hashicorp/terraform-provider-google/pull/18456))
* workbench: fixed issue where instance was not starting after an update in `google_workbench_instance` resource ([#18464](https://github.com/hashicorp/terraform-provider-google/pull/18464))
* workbench: fixed perma-diff caused by empty `accelerator_configs` in `google_workbench_instance` resource ([#18464](https://github.com/hashicorp/terraform-provider-google/pull/18464))

## 5.34.0 (June 17, 2024)

NOTES:
* compute: Updated field description of `connection_draining_timeout_sec`, `balancing_mode` and `outlier_detection` in `google_compute_region_backend_service` and `google_compute_backend_service`  to inform that default values will be changed in 6.0.0 ([#18399](https://github.com/hashicorp/terraform-provider-google/pull/18399))

FEATURES:
* **New Resource:** `google_netapp_backup` ([#18357](https://github.com/hashicorp/terraform-provider-google/pull/18357))
* **New Resource:** `google_network_services_service_lb_policies` ([#18326](https://github.com/hashicorp/terraform-provider-google/pull/18326))
* **New Resource:** `google_scc_management_folder_security_health_analytics_custom_module` ([#18360](https://github.com/hashicorp/terraform-provider-google/pull/18360))
* **New Resource:** `google_scc_management_project_security_health_analytics_custom_module` ([#18369](https://github.com/hashicorp/terraform-provider-google/pull/18369))
* **New Resource:** `google_scc_management_organization_security_health_analytics_custom_module` ([#18374](https://github.com/hashicorp/terraform-provider-google/pull/18374))

IMPROVEMENTS:
* alloydb: changed the resource `google_alloydb_instance` to be created directly with public IP enabled instead of creating the resource with public IP disabled and then enabling it ([#18344](https://github.com/hashicorp/terraform-provider-google/pull/18344))
* bigtable: added `automated_backup_configuration` field to `google_bigtable_table` resource ([#18335](https://github.com/hashicorp/terraform-provider-google/pull/18335))
* cloudbuildv2: added support for connecting to Bitbucket Data Center and Bitbucket Cloud with the `bitbucket_data_center_config` and `bitbucket_cloud_config` fields in `google_cloudbuildv2_connection` ([#18375](https://github.com/hashicorp/terraform-provider-google/pull/18375))
* compute: added update support to `ssl_policy` field in `google_compute_region_target_https_proxy` resource ([#18361](https://github.com/hashicorp/terraform-provider-google/pull/18361))
* compute: removed enum validation on `guest_os_features.type` in `google_compute_disk` to allow for new features to be used without provider update ([#18331](https://github.com/hashicorp/terraform-provider-google/pull/18331))
* compute: updated documentation of google_compute_target_https_proxy and google_compute_region_target_https_proxy ([#18358](https://github.com/hashicorp/terraform-provider-google/pull/18358))
* container: added support for `security_posture_config.mode` value "ENTERPRISE" in `resource_container_cluster` ([#18334](https://github.com/hashicorp/terraform-provider-google/pull/18334))
* discoveryengine: added `document_processing_config` field to `google_discovery_engine_data_store` resource ([#18350](https://github.com/hashicorp/terraform-provider-google/pull/18350))
* edgecontainer: added 'maintenance_exclusions' field to 'google_edgecontainer_cluster' resource ([#18370](https://github.com/hashicorp/terraform-provider-google/pull/18370))
* gkehub: added `prevent_drift` field to ConfigManagement `fleet_default_member_config` ([#18330](https://github.com/hashicorp/terraform-provider-google/pull/18330))
* netapp: added `administrators` field to `google_netapp_active_directory` resource ([#18333](https://github.com/hashicorp/terraform-provider-google/pull/18333))
* vertexai: promoted `optimized` field to GA for `google_vertex_ai_feature_online_store` resource ([#18348](https://github.com/hashicorp/terraform-provider-google/pull/18348))
* workbench: updated the metadata keys managed by the backend. ([#18367](https://github.com/hashicorp/terraform-provider-google/pull/18367))

BUG FIXES:
* compute: fixed an issue where `google_compute_instance_group_manager` with a pending operation was incorrectly removed due to the operation no longer being present in the backend ([#18380](https://github.com/hashicorp/terraform-provider-google/pull/18380))
* compute: fixed issue where users could not create `google_compute_security_policy` resources with `layer_7_ddos_defense_config` explicitly disabled ([#18345](https://github.com/hashicorp/terraform-provider-google/pull/18345))
* workbench: fixed a bug in the `google_workbench_instance` resource where specifying a network in some scenarios would cause instance creation to fail ([#18404](https://github.com/hashicorp/terraform-provider-google/pull/18404)

## 5.33.0 (June 10, 2024)

DEPRECATIONS:
* healthcare: deprecated `notification_config` in `google_healthcare_fhir_store` resource. Use `notification_configs` instead. ([#18306](https://github.com/hashicorp/terraform-provider-google/pull/18306))

FEATURES:
* **New Data Source:** `google_compute_security_policy` ([#18316](https://github.com/hashicorp/terraform-provider-google/pull/18316))
* **New Resource:** `google_compute_project_cloud_armor_tier` ([#18319](https://github.com/hashicorp/terraform-provider-google/pull/18319))
* **New Resource:** `google_network_services_service_lb_policies` ([#18326](https://github.com/hashicorp/terraform-provider-google/pull/18326))
* **New Resource:** `google_scc_management_organization_event_threat_detection_custom_module` ([#18317](https://github.com/hashicorp/terraform-provider-google/pull/18317))
* **New Resource:** `google_spanner_instance_config` ([#18322](https://github.com/hashicorp/terraform-provider-google/pull/18322))

IMPROVEMENTS:
* appengine: added `flexible_runtime_settings` field to `google_app_engine_flexible_app_version` resource ([#18325](https://github.com/hashicorp/terraform-provider-google/pull/18325))
* bigtable: added `force_destroy` field to `google_bigtable_instance` resource. This will force delete any backups present in the instance and allow the instance to be deleted. ([#18291](https://github.com/hashicorp/terraform-provider-google/pull/18291))
* clouddeploy: added `execution_configs.verbose` field to `google_clouddeploy_target` resource ([#18292](https://github.com/hashicorp/terraform-provider-google/pull/18292))
* compute: added `storage_pool` field to `google_compute_disk` resource ([#18273](https://github.com/hashicorp/terraform-provider-google/pull/18273))
* dlp: added `secrets_discovery_target`, `cloud_sql_target.filter.database_resource_reference`, and `big_query_target.filter.table_reference` fields to `google_data_loss_prevention_discovery_config` resource ([#18324](https://github.com/hashicorp/terraform-provider-google/pull/18324))
* gkebackup: added `backup_schedule.backup_config.permissive_mode` field to `google_gke_backup_backup_plan` resource ([#18266](https://github.com/hashicorp/terraform-provider-google/pull/18266))
* gkebackup: added `restore_config.restore_order` field to `google_gke_backup_restore_plan` resource ([#18266](https://github.com/hashicorp/terraform-provider-google/pull/18266))
* gkebackup: added `restore_config.volume_data_restore_policy_bindings` field to `google_gke_backup_restore_plan` resource ([#18266](https://github.com/hashicorp/terraform-provider-google/pull/18266))
* gkebackup: added new enum values `MERGE_SKIP_ON_CONFLICT`, `MERGE_REPLACE_VOLUME_ON_CONFLICT` and `MERGE_REPLACE_ON_CONFLICT` to field `restore_config.namespaced_resource_restore_mode` in `google_gke_backup_restore_plan` resource ([#18266](https://github.com/hashicorp/terraform-provider-google/pull/18266))
* healthcare: added `notification_config.send_for_bulk_import` field to `google_healthcare_dicom_store` resource ([#18320](https://github.com/hashicorp/terraform-provider-google/pull/18320))
* healthcare: added `notification_configs` field to `google_healthcare_fhir_store` resource ([#18306](https://github.com/hashicorp/terraform-provider-google/pull/18306))
* integrationconnectors: added `endpoint_global_access` field to `google_integration_connectors_endpoint_attachment` resource ([#18293](https://github.com/hashicorp/terraform-provider-google/pull/18293))
* netapp: added `backup_config` field to `google_netapp_volume` resource ([#18286](https://github.com/hashicorp/terraform-provider-google/pull/18286))
* redis: added `zone_distribution_config` field to `google_redis_cluster` resource ([#18307](https://github.com/hashicorp/terraform-provider-google/pull/18307))
* resourcemanager: added support for `range_type = "default-domains-netblocks"` in `google_netblock_ip_ranges` data source ([#18290](https://github.com/hashicorp/terraform-provider-google/pull/18290))
* secretmanager: added support for IAM conditions in `google_secret_manager_secret_iam_*` resources ([#18294](https://github.com/hashicorp/terraform-provider-google/pull/18294))
* workstations: added `boot_disk_size_gb`, `enable_nested_virtualization`, and `pool_size` to `host.gce_instance.boost_configs` in `google_workstations_workstation_config` resource ([#18310](https://github.com/hashicorp/terraform-provider-google/pull/18310))

BUG FIXES:
* container: fixed `google_container_node_pool` crash if `node_config.secondary_boot_disks.mode` is not set ([#18323](https://github.com/hashicorp/terraform-provider-google/pull/18323))
* dlp: removed `required` on `inspect_config.limits.max_findings_per_info_type.info_type` field to allow the use of default limit by not setting this field in `google_data_loss_prevention_inspect_template` resource ([#18285](https://github.com/hashicorp/terraform-provider-google/pull/18285))
* provider: fixed application default credential and access token authorization when `universe_domain` is set ([#18272](https://github.com/hashicorp/terraform-provider-google/pull/18272))


## 5.32.0 (June 3, 2024)

NOTES:
* privateca: converted `google_privateca_certificate_template` to now use the MMv1 engine instead of DCL ([#18224](https://github.com/hashicorp/terraform-provider-google/pull/18224))

FEATURES:
* **New Resource:** `google_dataplex_entry_type` ([#18229](https://github.com/hashicorp/terraform-provider-google/pull/18229))
* **New Resource:** `google_logging_log_view_iam_binding` ([#18243](https://github.com/hashicorp/terraform-provider-google/pull/18243))
* **New Resource:** `google_logging_log_view_iam_member` ([#18243](https://github.com/hashicorp/terraform-provider-google/pull/18243))
* **New Resource:** `google_logging_log_view_iam_policy` ([#18243](https://github.com/hashicorp/terraform-provider-google/pull/18243))

IMPROVEMENTS:
* alloydb: added `psc_config` field to `google_alloydb_cluster` resource ([#18263](https://github.com/hashicorp/terraform-provider-google/pull/18263))
* alloydb: added `psc_instance_config` field to `google_alloydb_instance` resource ([#18263](https://github.com/hashicorp/terraform-provider-google/pull/18263))
* cloudrunv2: added `default_uri_disabled` field to resource `google_cloud_run_v2_service` resource ([#18246](https://github.com/hashicorp/terraform-provider-google/pull/18246))
* compute: added `NONE` to acceptable options for `update_policy.minimal_action` field in `google_compute_instance_group_manager` resource ([#18236](https://github.com/hashicorp/terraform-provider-google/pull/18236))
* looker: increased validation length of `name` to `google_looker_instance` resource ([#18244](https://github.com/hashicorp/terraform-provider-google/pull/18244))
* sql: updated support for a new value `week5` in field `setting.maintenance_window.update_track` in `google_sql_database_instance` resource ([#18223](https://github.com/hashicorp/terraform-provider-google/pull/18223))

BUG FIXES:
* cloudrunv2: added validation for `timeout` field to `google_cloud_run_v2_job` and `google_cloud_run_v2_service` resources ([#18260](https://github.com/hashicorp/terraform-provider-google/pull/18260))
* compute: fixed permadiff in ordering of `advertised_ip_ranges.range` field on `google_compute_router` resource ([#18228](https://github.com/hashicorp/terraform-provider-google/pull/18228))
* iam: added a 10 second sleep when creating a 'google_service_account' resource to reduce eventual consistency errors([#18261](https://github.com/hashicorp/terraform-provider-google/pull/18261))
* storage: fixed `google_storage_bucket.lifecycle_rule.condition` block fields  `days_since_noncurrent_time` and `days_since_custom_time`  and `num_newer_versions` were not working for 0 value ([#18231](https://github.com/hashicorp/terraform-provider-google/pull/18231))

## 5.31.0 (May 28, 2024)

FEATURES:
* **New Data Source:** `google_compute_subnetworks` ([#18159](https://github.com/hashicorp/terraform-provider-google/pull/18159))
* **New Resource:** `google_dataplex_aspect_type` ([#18201](https://github.com/hashicorp/terraform-provider-google/pull/18201))
* **New Resource:** `google_dataplex_entry_group` ([#18188](https://github.com/hashicorp/terraform-provider-google/pull/18188))
* **New Resource:** `google_kms_autokey_config` ([#18179](https://github.com/hashicorp/terraform-provider-google/pull/18179))
* **New Resource:** `google_kms_key_handle` ([#18179](https://github.com/hashicorp/terraform-provider-google/pull/18179))
* **New Resource:** `google_network_services_lb_route_extension` ([#18195](https://github.com/hashicorp/terraform-provider-google/pull/18195))

IMPROVEMENTS:
* appengine: added field `instance_ip_mode` to resource `google_app_engine_flexible_app_version` resource (beta) ([#18168](https://github.com/hashicorp/terraform-provider-google/pull/18168))
* bigquery: added `external_data_configuration.bigtable_options` to `google_bigquery_table` ([#18181](https://github.com/hashicorp/terraform-provider-google/pull/18181))
* composer: added support for importing `google_composer_user_workloads_secret` via the "{{environment}}/{{name}}" format. ([#7390](https://github.com/hashicorp/terraform-provider-google-beta/pull/7390))
* composer: improved timeouts for `google_composer_user_workloads_secret`. ([#7390](https://github.com/hashicorp/terraform-provider-google-beta/pull/7390))
* compute: added `TLS_JA3_FINGERPRINT` and `USER_IP` options in field `rate_limit_options.enforce_on_key` to `google_compute_security_policy` resource ([#18167](https://github.com/hashicorp/terraform-provider-google/pull/18167))
* compute: added 'rateLimitOptions' field to 'google_compute_security_policy_rule' resource ([#18167](https://github.com/hashicorp/terraform-provider-google/pull/18167))
* compute: changed `google_compute_region_ssl_policy`'s `region` field to optional and allow to be inferred from environment ([#18178](https://github.com/hashicorp/terraform-provider-google/pull/18178))
* compute: added `subnet_length` field to `google_compute_interconnect_attachment` resource ([#18187](https://github.com/hashicorp/terraform-provider-google/pull/18187))
* container: added `containerd_config` field and subfields to `google_container_cluster` and `google_container_node_pool` resources, to allow those resources to access private image registries. ([#18160](https://github.com/hashicorp/terraform-provider-google/pull/18160))
* container: allowed both `enable_autopilot` and `workload_identity_config` to be set in `google_container_cluster` resource. ([#18166](https://github.com/hashicorp/terraform-provider-google/pull/18166))
* datastream: added `create_without_validation` field to `google_datastream_connection_profile`, `google_datastream_private_connection` and `google_datastream_stream` resources ([#18176](https://github.com/hashicorp/terraform-provider-google/pull/18176))
* network-security: added `trust_config`, `min_tls_version`, `tls_feature_profile` and `custom_tls_features` fields to `google_network_security_tls_inspection_policy` resource ([#18139](https://github.com/hashicorp/terraform-provider-google/pull/18139))
* networkservices: made field `load_balancing_scheme` immutable in resource `google_network_services_lb_traffic_extension`, as in-place updating is always failing ([#18195](https://github.com/hashicorp/terraform-provider-google/pull/18195))
* networkservices: made required fields `extension_chains.extensions.authority ` and `extension_chains.extensions.timeout` optional in resource `google_network_services_lb_traffic_extension` ([#18195](https://github.com/hashicorp/terraform-provider-google/pull/18195))
* networkservices: removed unsupported load balancing scheme `LOAD_BALANCING_SCHEME_UNSPECIFIED` from the field `load_balancing_scheme` in resource `google_network_services_lb_traffic_extension` ([#18195](https://github.com/hashicorp/terraform-provider-google/pull/18195))
* pubsub: added `cloud_storage_config.filename_datetime_format` field to `google_pubsub_subscription` resource ([#18180](https://github.com/hashicorp/terraform-provider-google/pull/18180))
* tpu: added `type` of `accelerator_config` to `google_tpu_v2_vm` resource ([#18148](https://github.com/hashicorp/terraform-provider-google/pull/18148))

BUG FIXES:
* monitoring: fixed a permadiff with `monitored_resource.labels` property in the `google_monitoring_uptime_check_config` resource ([#18174](https://github.com/hashicorp/terraform-provider-google/pull/18174))
* storage: fixed a bug where field `autoclass` block is generating permadiff whenever the block is removed from the config  in `google_storage_bucket` resource ([#18197](https://github.com/hashicorp/terraform-provider-google/pull/18197))
* storagetransfer: fixed a permadiff with `transfer_spec.0.aws_s3_data_source.0.aws_access_key` `resource_storage_transfer_job` ([#18190](https://github.com/hashicorp/terraform-provider-google/pull/18190))

## 5.30.0 (May 20, 2024)

FEATURES:
* **New Data Source:** `google_cloud_asset_resources_search_all` ([#18129](https://github.com/hashicorp/terraform-provider-google/pull/18129))
* **New Resource:** `google_compute_interconnect` ([#18064](https://github.com/hashicorp/terraform-provider-google/pull/18064))
* **New Resource:** `google_network_services_lb_traffic_extension` ([#18138](https://github.com/hashicorp/terraform-provider-google/pull/18138))

IMPROVEMENTS:
* compute:  added `kms_key_name` field to `google_bigquery_connection` resource ([#18057](https://github.com/hashicorp/terraform-provider-google/pull/18057))
* compute: added `auto_network_tier` field to `google_compute_router_nat` resource ([#18055](https://github.com/hashicorp/terraform-provider-google/pull/18055))
* compute: promoted `enable_ipv4`, `ipv4_nexthop_address` and `peer_ipv4_nexthop_address` fields in `google_compute_router_peer` resource to GA ([#18056](https://github.com/hashicorp/terraform-provider-google/pull/18056))
* compute: promoted `identifier_range` field in `google_compute_router` resource to GA ([#18056](https://github.com/hashicorp/terraform-provider-google/pull/18056))
* compute: promoted `ip_version` field in `google_compute_router_interface` resource to GA ([#18056](https://github.com/hashicorp/terraform-provider-google/pull/18056))
* container: added `KUBELET` and `CADVISOR` options to `monitoring_config.enable_components` in `google_container_cluster` resource ([#18090](https://github.com/hashicorp/terraform-provider-google/pull/18090))
* dataproc: added `local_ssd_interface` to `google_dataproc_cluster` resource ([#18137](https://github.com/hashicorp/terraform-provider-google/pull/18137))
* dataprocmetastore: promoted `google_dataproc_metastore_federation` to GA ([#18084](https://github.com/hashicorp/terraform-provider-google/pull/18084))
* dlp: added `cloud_sql_target` field to `google_data_loss_prevention_discovery_config` resource ([#18063](https://github.com/hashicorp/terraform-provider-google/pull/18063))
* netapp: added `FLEX` value to field `service_level` in `google_netapp_storage_pool` resource ([#18088](https://github.com/hashicorp/terraform-provider-google/pull/18088))
* networksecurity: added `trust_config`, `min_tls_version`, `tls_feature_profile` and `custom_tls_features` fields to `google_network_security_tls_inspection_policy` resource ([#18139](https://github.com/hashicorp/terraform-provider-google/pull/18139))
* networkservices: supported in-place update for `gateway_security_policy` and `certificate_urls` fields in `google_network_services_gateway` resource ([#18082](https://github.com/hashicorp/terraform-provider-google/pull/18082))

BUG FIXES:
* compute: fixed a perma-diff on `machine_type` field in `google_compute_instance` resource ([#18071](https://github.com/hashicorp/terraform-provider-google/pull/18071))
* compute: fixed a perma-diff on `type` field in `google_compute_disk` resource ([#18071](https://github.com/hashicorp/terraform-provider-google/pull/18071))
* storage: fixed update issue for `lifecycle_rule.condition.custom_time_before` and `lifecycle_rule.condition.noncurrent_time_before` in `google_storage_bucket` resource ([#18127](https://github.com/hashicorp/terraform-provider-google/pull/18127))

## 5.29.1 (May 14, 2024)

BREAKING CHANGES:
* compute: removed `secondary_ip_range.reserved_internal_range` field from `google_compute_subnetwork` ([18133](https://github.com/hashicorp/terraform-provider-google/pull/18133))

## 5.29.0 (May 13, 2024)

NOTES:
* compute: added documentation for `md5_authentication_key` field in `google_compute_router_peer` resource. The field was introduced in [v5.12.0](https://github.com/hashicorp/terraform-provider-google/releases/tag/v5.12.0), but documentation was unintentionally omitted at that time. ([#17991](https://github.com/hashicorp/terraform-provider-google/pull/17991))

FEATURES:
* **New Resource:** `google_bigtable_authorized_view` ([#18006](https://github.com/hashicorp/terraform-provider-google/pull/18006))
* **New Resource:** `google_integration_connectors_managed_zone` ([#18029](https://github.com/hashicorp/terraform-provider-google/pull/18029))
* **New Resource:** `google_network_connectivity_regional_endpoint` ([#18014](https://github.com/hashicorp/terraform-provider-google/pull/18014))
* **New Resource:** `google_network_security_security_profile` ([#18025](https://github.com/hashicorp/terraform-provider-google/pull/18025))
* **New Resource:** `google_network_security_security_profile_group` ([#18025](https://github.com/hashicorp/terraform-provider-google/pull/18025))
* **New Resource:** `google_network_security_firewall_endpoint` ([#18025](https://github.com/hashicorp/terraform-provider-google/pull/18025))
* **New Resource:** `google_network_security_firewall_endpoint_association` ([#18025](https://github.com/hashicorp/terraform-provider-google/pull/18025))

IMPROVEMENTS:
* clouddeploy: added `custom_target` field to  `google_clouddeploy_target` resource ([#18000](https://github.com/hashicorp/terraform-provider-google/pull/18000))
* clouddeploy: added `google_cloud_build_repo` to `custom_target_type` resource ([#18040](https://github.com/hashicorp/terraform-provider-google/pull/18040))
* compute: added `preconfigured_waf_config` field to `google_compute_region_security_policy_rule` resource; ([#18039](https://github.com/hashicorp/terraform-provider-google/pull/18039))
* compute: added `rate_limit_options` field to `google_compute_region_security_policy_rule` resource; ([#18039](https://github.com/hashicorp/terraform-provider-google/pull/18039))
* compute: added `security_profile_group`, `tls_inspect` to `google_compute_firewall_policy_rule` ([#18000](https://github.com/hashicorp/terraform-provider-google/pull/18000))
* compute: added `security_profile_group`, `tls_inspect` to `google_compute_network_firewall_policy_rule` ([#18000](https://github.com/hashicorp/terraform-provider-google/pull/18000))
* compute: added fields `reserved_internal_range` and `secondary_ip_ranges.reserved_internal_range` to `google_compute_subnetwork` resource ([#18026](https://github.com/hashicorp/terraform-provider-google/pull/18026))
* container: added `dns_config.additive_vpc_scope_dns_domain` field to `google_container_cluster` resource ([#18031](https://github.com/hashicorp/terraform-provider-google/pull/18031))
* container: added `enable_nested_virtualization` field to `google_container_node_pool` and `google_container_cluster` resource. ([#18015](https://github.com/hashicorp/terraform-provider-google/pull/18015))
* iam: added `extra_attributes_oauth2_client` field to `google_iam_workforce_pool_provider` resource ([#18027](https://github.com/hashicorp/terraform-provider-google/pull/18027))
* privateca: added `maximum_lifetime` field to  `google_privateca_certificate_template` resource ([#18000](https://github.com/hashicorp/terraform-provider-google/pull/18000))

## 5.28.0 (May 6, 2024)

DEPRECATIONS:
* integrations: deprecated `create_sample_workflows` and `provision_gmek` fields in `google_integrations_client`.  ([#17945](https://github.com/hashicorp/terraform-provider-google/pull/17945))

FEATURES:
* **New Data Source:** `google_storage_buckets` ([#17960](https://github.com/hashicorp/terraform-provider-google/pull/17960))
* **New Resource:** `google_compute_security_policy_rule` ([#17937](https://github.com/hashicorp/terraform-provider-google/pull/17937))

IMPROVEMENTS:
* alloydb: added `maintenance_update_policy` field to `google_alloydb_cluster` resource ([#17954](https://github.com/hashicorp/terraform-provider-google/pull/17954))
* bigquery: promoted `external_dataset_reference` in `google_bigquery_dataset` to GA ([#17944](https://github.com/hashicorp/terraform-provider-google/pull/17944))
* composer: promoted `config.software_config.image_version` in-place update to GA in resource `google_composer_environment` ([#17986](https://github.com/hashicorp/terraform-provider-google/pull/17986))
* container: added `node_config.secondary_boot_disks` field to `google_container_node_pool` ([#17962](https://github.com/hashicorp/terraform-provider-google/pull/17962))
* integrations: added `create_sample_integrations` field to `google_integrations_client`, replacing deprecated field `create_sample_workflows`. ([#17945](https://github.com/hashicorp/terraform-provider-google/pull/17945))
* redis: added `redis_configs` field to `google_redis_cluster` resource ([#17956](https://github.com/hashicorp/terraform-provider-google/pull/17956))

BUG FIXES:
* dns: fixed bug where the deletion of `google_dns_managed_zone` resources was blocked by any associated SOA-type `google_dns_record_set` resources ([#17989](https://github.com/hashicorp/terraform-provider-google/pull/17989))
* storage: fixed an issue where `google_storage_bucket_object` and `google_storage_bucket_objects` data sources would ignore custom endpoints ([#17952](https://github.com/hashicorp/terraform-provider-google/pull/17952))

## 5.27.0 (Apr 30, 2024)

FEATURES:
* **New Data Source:** `google_storage_bucket_objects` ([#17920](https://github.com/hashicorp/terraform-provider-google/pull/17920))
* **New Resource:** `google_compute_security_policy_rule` ([#17937](https://github.com/hashicorp/terraform-provider-google/pull/17937))
* **New Resource:** `google_data_loss_prevention_discovery_config` ([#17887](https://github.com/hashicorp/terraform-provider-google/pull/17887))
* **New Resource:** `google_integrations_auth_config` ([#17917](https://github.com/hashicorp/terraform-provider-google/pull/17917))
* **New Resource:** `google_network_connectivity_internal_range` ([#17909](https://github.com/hashicorp/terraform-provider-google/pull/17909))

IMPROVEMENTS:
* alloydb: added `network_config` field to `google_alloydb_instance` resource ([#17921](https://github.com/hashicorp/terraform-provider-google/pull/17921))
* alloydb: added `public_ip_address` field  to `google_alloydb_instance` resource ([#17921](https://github.com/hashicorp/terraform-provider-google/pull/17921))
* apigee: added `forward_proxy_uri` field to `google_apigee_environment` resource ([#17902](https://github.com/hashicorp/terraform-provider-google/pull/17902))
* bigquerydatapolicy: added `data_masking_policy.routine` field to `google_bigquery_data_policy` resource ([#17885](https://github.com/hashicorp/terraform-provider-google/pull/17885))
* compute: added `server_tls_policy` field to `google_compute_region_target_https_proxy` resource ([#17934](https://github.com/hashicorp/terraform-provider-google/pull/17934))
* logging: added `intercept_children` field to `google_logging_organization_sink` and `google_logging_folder_sink` resources ([#17932](https://github.com/hashicorp/terraform-provider-google/pull/17932))
* monitoring: added `service_agent_authentication` field to `google_monitoring_uptime_check_config` resource ([#17929](https://github.com/hashicorp/terraform-provider-google/pull/17929))
* privateca: added `subject_key_id` field to `google_privateca_certificate` and `google_privateca_certificate_authority` resources ([#17923](https://github.com/hashicorp/terraform-provider-google/pull/17923))
* secretmanager: added `version_destroy_ttl` field to `google_secret_manager_secret` resource ([#17888](https://github.com/hashicorp/terraform-provider-google/pull/17888))

BUG FIXES:
* appengine: added suppression for a diff in `google_app_engine_standard_app_version.automatic_scaling` when the block is unset in configuration ([#17905](https://github.com/hashicorp/terraform-provider-google/pull/17905))
* sql: fixed issues with updating the `enable_google_ml_integration` field in `google_sql_database_instance` resource ([#17878](https://github.com/hashicorp/terraform-provider-google/pull/17878))

## 5.26.0 (Apr 22, 2024)

FEATURES:
* **New Resource:** `google_project_iam_member_remove` ([#17871](https://github.com/hashicorp/terraform-provider-google/pull/17871))

IMPROVEMENTS:
* apigee: added support for `api_consumer_data_location`, `api_consumer_data_encryption_key_name`, and `control_plane_encryption_key_name` in `google_apigee_organization` ([#17874](https://github.com/hashicorp/terraform-provider-google/pull/17874))
* artifactregistry: added `remote_repository_config.<facade>_repository.custom_repository.uri` field to `google_artifact_registry_repository` resource. ([#17840](https://github.com/hashicorp/terraform-provider-google/pull/17840))
* bigquery: added `resource_tags` field to `google_bigquery_table` resource ([#17876](https://github.com/hashicorp/terraform-provider-google/pull/17876))
* billing: added `ownership_scope` field to `google_billing_budget` resource ([#17868](https://github.com/hashicorp/terraform-provider-google/pull/17868))
* cloudfunctions2: added `build_config.service_account` field to `google_cloudfunctions2_function` resource ([#17841](https://github.com/hashicorp/terraform-provider-google/pull/17841))
* resourcemanager: added the field `api_method` to datasource `google_active_folder` so you can use either `SEARCH` or `LIST` to find your folder ([#17877](https://github.com/hashicorp/terraform-provider-google/pull/17877))
* storage: added labels validation to `google_storage_bucket` resource ([#17806](https://github.com/hashicorp/terraform-provider-google/pull/17806))

BUG FIXES:
* apigee: fixed permadiff in ordering of `google_apigee_organization.properties.property`. ([#17850](https://github.com/hashicorp/terraform-provider-google/pull/17850))
* cloudrun: fixed the bug that computed `metadata.0.labels` and `metadata.0.annotations` fields don't appear in terraform plan when creating resource `google_cloud_run_service` and `google_cloud_run_domain_mapping` ([#17815](https://github.com/hashicorp/terraform-provider-google/pull/17815))
* dns: fixed bug where some methods of authentication didn't work when using `dns` data sources ([#17847](https://github.com/hashicorp/terraform-provider-google/pull/17847))
* iam: fixed a bug that prevented setting `create_ignore_already_exists` on existing resources in `google_service_account`. ([#17856](https://github.com/hashicorp/terraform-provider-google/pull/17856))
* sql: fixed issues with updating the `enable_google_ml_integration` field in `google_sql_database_instance` resource ([#17878](https://github.com/hashicorp/terraform-provider-google/pull/17878))
* storage: added validation to `name` field in `google_storage_bucket` resource ([#17858](https://github.com/hashicorp/terraform-provider-google/pull/17858))
* vmwareengine: fixed stretched cluster creation in `google_vmwareengine_private_cloud` ([#17875](https://github.com/hashicorp/terraform-provider-google/pull/17875))

## 5.25.0 (Apr 15, 2024)

FEATURES:
* **New Data Source:** `google_tags_tag_keys` ([#17782](https://github.com/hashicorp/terraform-provider-google/pull/17782))
* **New Data Source:** `google_tags_tag_values` ([#17782](https://github.com/hashicorp/terraform-provider-google/pull/17782))

IMPROVEMENTS:
* bigquery: added in-place schema column drop support for `google_bigquery_table` resource ([#17777](https://github.com/hashicorp/terraform-provider-google/pull/17777))
* compute: added `endpoint_types` field to `google_compute_router_nat` resource ([#17771](https://github.com/hashicorp/terraform-provider-google/pull/17771))
* compute: increased timeouts from 8 minutes to 20 minutes for `google_compute_security_policy` resource ([#17793](https://github.com/hashicorp/terraform-provider-google/pull/17793))
* compute: promoted `google_compute_instance_settings` to GA ([#17781](https://github.com/hashicorp/terraform-provider-google/pull/17781))
* container: added `stateful_ha_config` field to `google_container_cluster` resource ([#17796](https://github.com/hashicorp/terraform-provider-google/pull/17796))
* firestore: added `vector_config` field to `google_firestore_index` resource ([#17758](https://github.com/hashicorp/terraform-provider-google/pull/17758))
* gkebackup: added `backup_schedule.rpo_config` field to `google_gke_backup_backup_plan` resource ([#17805](https://github.com/hashicorp/terraform-provider-google/pull/17805))
* networksecurity: added `disabled` field to `google_network_security_firewall_endpoint_association` resource; ([#17762](https://github.com/hashicorp/terraform-provider-google/pull/17762))
* sql: added `enable_google_ml_integration` field to `google_sql_database_instance` resource ([#17798](https://github.com/hashicorp/terraform-provider-google/pull/17798))
* storage: added labels validation to `google_storage_bucket` resource ([#17806](https://github.com/hashicorp/terraform-provider-google/pull/17806))
* vmwareengine: added `preferred_zone` and `secondary_zone` fields to `google_vmwareengine_private_cloud` resource ([#17803](https://github.com/hashicorp/terraform-provider-google/pull/17803))

BUG FIXES:
* networksecurity: fixed an issue where `google_network_security_firewall_endpoint_association` resources could not be created due to a bad parameter ([#17762](https://github.com/hashicorp/terraform-provider-google/pull/17762))
* privateca: fixed permission issue by specifying signer certs chain when activating a sub-CA across regions for `google_privateca_certificate_authority` resource ([#17783](https://github.com/hashicorp/terraform-provider-google/pull/17783))

## 5.24.0 (Apr 8, 2024)

IMPROVEMENTS:
* container: added `enable_cilium_clusterwide_network_policy` field to `google_container_cluster` resource ([#17738](https://github.com/hashicorp/terraform-provider-google/pull/17738))
* container: added `node_pool_auto_config.resource_manager_tags` field to `google_container_cluster` resource ([#17715](https://github.com/hashicorp/terraform-provider-google/pull/17715))
* gkeonprem: added `disable_bundled_ingress` field to `google_gkeonprem_vmware_cluster` resource ([#17718](https://github.com/hashicorp/terraform-provider-google/pull/17718))
* redis: added `node_type` and `precise_size_gb` fields to `google_redis_cluster` ([#17742](https://github.com/hashicorp/terraform-provider-google/pull/17742))
* storage: added `project_number` attribute to `google_storage_bucket` resource and data source ([#17719](https://github.com/hashicorp/terraform-provider-google/pull/17719))
* storage: added ability to provide `project` argument to `google_storage_bucket` data source. This will not impact reading the resource's data, instead this helps users avoid calls to the Compute API within the data source. ([#17719](https://github.com/hashicorp/terraform-provider-google/pull/17719))

BUG FIXES:
* appengine: fixed a crash in `google_app_engine_flexible_app_version` due to the `deployment` field not being returned by the API ([#17744](https://github.com/hashicorp/terraform-provider-google/pull/17744))
* bigquery: fixed a crash when `google_bigquery_table` had a `primary_key.columns` entry set to `""` ([#17721](https://github.com/hashicorp/terraform-provider-google/pull/17721))
* compute: fixed update scenarios on`google_compute_region_target_https_proxy` and `google_compute_target_https_proxy` resources. ([#17733](https://github.com/hashicorp/terraform-provider-google/pull/17733))

## 5.23.0 (Apr 1, 2024)

NOTES:
* provider: introduced support for [provider-defined functions](https://developer.hashicorp.com/terraform/plugin/framework/functions). This feature is in Terraform v1.8.0+. ([#17694](https://github.com/hashicorp/terraform-provider-google/pull/17694))

DEPRECATIONS:
* kms: deprecated `attestation.external_protection_level_options` in favor of `external_protection_level_options` in `google_kms_crypto_key_version` ([#17704](https://github.com/hashicorp/terraform-provider-google/pull/17704))

FEATURES:
* **New Data Source:** `google_apphub_application` ([#17679](https://github.com/hashicorp/terraform-provider-google/pull/17679))
* **New Resource:** `google_cloud_quotas_quota_preference` ([#17637](https://github.com/hashicorp/terraform-provider-google/pull/17637))
* **New Resource:** `google_vertex_ai_deployment_resource_pool` ([#17707](https://github.com/hashicorp/terraform-provider-google/pull/17707))
* **New Resource:** `google_integrations_client` ([#17640](https://github.com/hashicorp/terraform-provider-google/pull/17640))

IMPROVEMENTS:
* bigquery: added `dataGovernanceType` to `google_bigquery_routine` resource ([#17689](https://github.com/hashicorp/terraform-provider-google/pull/17689))
* bigquery: added support for `external_data_configuration.json_extension` to `google_bigquery_table` ([#17663](https://github.com/hashicorp/terraform-provider-google/pull/17663))
* compute: added `cloud_router_ipv6_address`, `customer_router_ipv6_address` fields to `google_compute_interconnect_attachment` resource ([#17692](https://github.com/hashicorp/terraform-provider-google/pull/17692))
* compute: added `generated_id` field to `google_compute_region_backend_service` resource ([#17639](https://github.com/hashicorp/terraform-provider-google/pull/17639))
* integrations: added deletion support for `google_integrations_client` resource ([#17678](https://github.com/hashicorp/terraform-provider-google/pull/17678))
* kms: added `crypto_key_backend` field to `google_kms_crypto_key` resource ([#17704](https://github.com/hashicorp/terraform-provider-google/pull/17704))
* metastore: added `scheduled_backup` field to `google_dataproc_metastore_service` resource ([#17673](https://github.com/hashicorp/terraform-provider-google/pull/17673))
* provider: added provider-defined function `name_from_id` for retrieving the short-form name of a resource from its self link or id ([#17694](https://github.com/hashicorp/terraform-provider-google/pull/17694))
* provider: added provider-defined function `project_from_id` for retrieving the project id from a resource's self link or id ([#17694](https://github.com/hashicorp/terraform-provider-google/pull/17694))
* provider: added provider-defined function `region_from_zone` for deriving a region from a zone's name ([#17694](https://github.com/hashicorp/terraform-provider-google/pull/17694))
* provider: added provider-defined functions `location_from_id`, `region_from_id`, and `zone_from_id` for retrieving the location/region/zone names from a resource's self link or id ([#17694](https://github.com/hashicorp/terraform-provider-google/pull/17694))

BUG FIXES:
* cloudrunv2: fixed Terraform state inconsistency when resource `google_cloud_run_v2_job` creation fails ([#17711](https://github.com/hashicorp/terraform-provider-google/pull/17711))
* cloudrunv2: fixed Terraform state inconsistency when resource `google_cloud_run_v2_service` creation fails ([#17711](https://github.com/hashicorp/terraform-provider-google/pull/17711))
* container: fixed `google_container_cluster` permadiff when `master_ipv4_cidr_block` is set for a private flexible cluster ([#17687](https://github.com/hashicorp/terraform-provider-google/pull/17687))
* dataflow: fixed an issue where the provider would crash when `enableStreamingEngine` is set as a `parameter` value in `google_dataflow_flex_template_job` ([#17712](https://github.com/hashicorp/terraform-provider-google/pull/17712))
* kms: added top-level `external_protection_level_options` field in `google_kms_crypto_key_version` resource ([#17704](https://github.com/hashicorp/terraform-provider-google/pull/17704))

## 5.22.0 (Mar 26, 2024)

BREAKING CHANGES:
* networksecurity: added required field `billing_project_id` to `google_network_security_firewall_endpoint` resource. Any configuration without `billing_project_id` specified will cause resource creation fail (beta) ([#17630](https://github.com/hashicorp/terraform-provider-google/pull/17630))

FEATURES:
* **New Data Source:** `google_cloud_quotas_quota_info` ([#17564](https://github.com/hashicorp/terraform-provider-google/pull/17564))
* **New Data Source:** `google_cloud_quotas_quota_infos` ([#17617](https://github.com/hashicorp/terraform-provider-google/pull/17617))
* **New Resource:** `google_access_context_manager_service_perimeter_dry_run_resource` ([#17614](https://github.com/hashicorp/terraform-provider-google/pull/17614))

IMPROVEMENTS:
* accesscontextmanager: supported managing service perimeter dry run resources outside the perimeter via new resource `google_access_context_manager_service_perimeter_dry_run_resource` ([#17614](https://github.com/hashicorp/terraform-provider-google/pull/17614))
* cloudrunv2: added plan-time validation to restrict number of ports to 1 in `google_cloud_run_v2_service` ([#17594](https://github.com/hashicorp/terraform-provider-google/pull/17594))
* composer: added field `count` to validate number of DAG processors in `google_composer_environment` ([#17625](https://github.com/hashicorp/terraform-provider-google/pull/17625))
* compute: added enumeration value `SEV_LIVE_MIGRATABLE_V2` for the `guest_os_features` of `google_compute_disk` ([#17629](https://github.com/hashicorp/terraform-provider-google/pull/17629))
* compute: added `status.all_instances_config.revision` field to `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#17595](https://github.com/hashicorp/terraform-provider-google/pull/17595))
* compute: added field `path_template_match` to resource `google_compute_region_url_map` ([#17571](https://github.com/hashicorp/terraform-provider-google/pull/17571))
* compute: added field `path_template_rewrite` to resource `google_compute_region_url_map` ([#17571](https://github.com/hashicorp/terraform-provider-google/pull/17571))
* pubsub: added `ingestion_data_source_settings` field to `google_pubsub_topic` resource ([#17604](https://github.com/hashicorp/terraform-provider-google/pull/17604))
* storage: added 'soft_delete_policy' to 'google_storage_bucket' resource ([#17624](https://github.com/hashicorp/terraform-provider-google/pull/17624))

BUG FIXES:
* accesscontextmanager: fixed an issue with `access_context_manager_service_perimeter_ingress_policy` and `access_context_manager_service_perimeter_egress_policy` where updates could not be applied after initial creation. Any updates applied to these resources will now involve their recreation. To ensure that new policies are added before old ones are removed, add a `lifecycle` block with `create_before_destroy = true` to your resource configuration alongside other updates. ([#17596](https://github.com/hashicorp/terraform-provider-google/pull/17596))
* firebase: made the `google_firebase_android_app` resource's `package_name` field required and immutable. This prevents API errors encountered by users who attempted to update or leave that field unset in their configurations. ([#17585](https://github.com/hashicorp/terraform-provider-google/pull/17585))
* spanner: removed validation function for the field `version_retention_period` in the resource `google_spanner_database` and directly returned error from backend ([#17621](https://github.com/hashicorp/terraform-provider-google/pull/17621))

## 5.21.0 (Mar 18, 2024)

FEATURES:
* **New Data Source:** `google_apphub_discovered_service` ([#17548](https://github.com/hashicorp/terraform-provider-google/pull/17548))
* **New Data Source:** `google_apphub_discovered_workload` ([#17553](https://github.com/hashicorp/terraform-provider-google/pull/17553))
* **New Data Source:** `google_cloud_quotas_quota_info` ([#17564](https://github.com/hashicorp/terraform-provider-google/pull/17564))
* **New Resource:** `google_apphub_workload` ([#17561](https://github.com/hashicorp/terraform-provider-google/pull/17561))
* **New Resource:** `google_firebase_app_check_device_check_config` ([#17517](https://github.com/hashicorp/terraform-provider-google/pull/17517))
* **New Resource:** `google_iap_tunnel_dest_group` ([#17533](https://github.com/hashicorp/terraform-provider-google/pull/17533))
* **New Resource:** `google_kms_ekm_connection` ([#17512](https://github.com/hashicorp/terraform-provider-google/pull/17512))
* **New Resource:** `google_apphub_application` ([#17499](https://github.com/hashicorp/terraform-provider-google/pull/17499))
* **New Resource:** `google_apphub_service` ([#17562](https://github.com/hashicorp/terraform-provider-google/pull/17562))
* **New Resource:** `google_apphub_service_project_attachment` ([#17536](https://github.com/hashicorp/terraform-provider-google/pull/17536))
* **New Resource:** `google_network_security_firewall_endpoint_association` ([#17540](https://github.com/hashicorp/terraform-provider-google/pull/17540))

IMPROVEMENTS:
* cloudrunv2: added support for `scaling.min_instance_count` in `google_cloud_run_v2_service`. ([#17501](https://github.com/hashicorp/terraform-provider-google/pull/17501))
* compute: added `metric.single_instance_assignment` and `metric.filter` to `google_compute_region_autoscaler` ([#17519](https://github.com/hashicorp/terraform-provider-google/pull/17519))
* container: added `queued_provisioning` to `google_container_node_pool` ([#17549](https://github.com/hashicorp/terraform-provider-google/pull/17549))
* gkeonprem: allowed `vcenter_network` to be set in `google_gkeonprem_vmware_cluster`, previously it was output-only ([#17505](https://github.com/hashicorp/terraform-provider-google/pull/17505))
* workstations: added support for `ephemeral_directories` in `google_workstations_workstation_config` ([#17515](https://github.com/hashicorp/terraform-provider-google/pull/17515))

BUG FIXES:
* compute: allowed sending empty values for `SERVERLESS` in `google_compute_region_network_endpoint_group` resource ([#17500](https://github.com/hashicorp/terraform-provider-google/pull/17500))
* notebooks: fixed an issue where default tags would cause a diff recreating `google_notebooks_instance` resources ([#17559](https://github.com/hashicorp/terraform-provider-google/pull/17559))
* storage: fixed an issue where two or more lifecycle rules with different values of `no_age` field always generates change in `google_storage_bucket` resource. ([#17513](https://github.com/hashicorp/terraform-provider-google/pull/17513))

## 5.20.0 (Mar 11, 2024)

FEATURES:
* **New Resource:** `google_clouddeploy_custom_target_type_iam_*` ([#17445](https://github.com/hashicorp/terraform-provider-google/pull/17445))

IMPROVEMENTS:
* certificatemanager: added `type` field to `google_certificate_manager_dns_authorization` resource ([#17459](https://github.com/hashicorp/terraform-provider-google/pull/17459))
* compute: added the `network_url` attribute to the `consumer_accept_list`-block of the `google_compute_service_attachment` resource ([#17492](https://github.com/hashicorp/terraform-provider-google/pull/17492))
* gkehub: added support for `policycontroller.policy_controller_hub_config.policy_content.bundles` and 
`policycontroller.policy_controller_hub_config.deployment_configs` fields to `google_gke_hub_feature_membership` ([#17483](https://github.com/hashicorp/terraform-provider-google/pull/17483))

BUG FIXES:
* artifactregistry: fixed permadiff when `google_artifact_repository.docker_config` field is unset ([#17484](https://github.com/hashicorp/terraform-provider-google/pull/17484))
* bigquery: corrected plan-time validation on `google_bigquery_dataset.dataset_id` ([#17449](https://github.com/hashicorp/terraform-provider-google/pull/17449))
* kms: fixed issue where `google_kms_crypto_key_version.attestation.cert_chains` properties were incorrectly set to type string ([#17486](https://github.com/hashicorp/terraform-provider-google/pull/17486))

## 5.19.0 (Mar 4, 2024)

FEATURES:
* **New Resource:** `google_clouddeploy_automation`([#17427](https://github.com/hashicorp/terraform-provider-google/pull/17427))
* **New Resource:** `google_clouddeploy_target_iam_*` ([#17368](https://github.com/hashicorp/terraform-provider-google/pull/17368))

IMPROVEMENTS:
* bigquery: added `remote_function_options` field to `google_bigquery_routine` resource ([#17382](https://github.com/hashicorp/terraform-provider-google/pull/17382))
* certificatemanager: added `location` field to `google_certificate_manager_dns_authorization` resource ([#17358](https://github.com/hashicorp/terraform-provider-google/pull/17358))
* composer: added validations for composer 2/3 only fields in `google_composer_environment` ([#17361](https://github.com/hashicorp/terraform-provider-google/pull/17361))
* compute: added `certificate_manager_certificates` field to `google_compute_region_target_https_proxy` resource ([#17365](https://github.com/hashicorp/terraform-provider-google/pull/17365))
* compute: promoted `all_instances_config` field in resources `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` to GA ([#17414](https://github.com/hashicorp/terraform-provider-google/pull/17414))
* container: promoted `enable_confidential_storage` from `node_config` in `google_container_cluster` and `google_container_node_pool` to GA ([#17367](https://github.com/hashicorp/terraform-provider-google/pull/17367))
* gkehub2: added `namespace_labels` field to `google_gke_hub_scope` resource ([#17421](https://github.com/hashicorp/terraform-provider-google/pull/17421))

BUG FIXES:
* resourcemanager: added a retry to deleting the default network when `auto_create_network` is false in `google_project` ([#17419](https://github.com/hashicorp/terraform-provider-google/pull/17419))

## 5.18.0 (Feb 26, 2024)

BREAKING CHANGES:
* securityposture: marked `policy_sets` and `policy_sets.policies` required in `google_securityposture_posture`. API validation already enforced this, so no resources could be provisioned without these ([#17303](https://github.com/hashicorp/terraform-provider-google/pull/17303))

FEATURES:
* **New Data Source:** `google_compute_forwarding_rules` ([#17342](https://github.com/hashicorp/terraform-provider-google/pull/17342))
* **New Resource:** `google_firebase_app_check_app_attest_config` ([#17279](https://github.com/hashicorp/terraform-provider-google/pull/17279))
* **New Resource:** `google_firebase_app_check_play_integrity_config` ([#17279](https://github.com/hashicorp/terraform-provider-google/pull/17279))
* **New Resource:** `google_firebase_app_check_recaptcha_enterprise_config` ([#17327](https://github.com/hashicorp/terraform-provider-google/pull/17327))
* **New Resource:** `google_firebase_app_check_recaptcha_v3_config` ([#17327](https://github.com/hashicorp/terraform-provider-google/pull/17327))
* **New Resource:** `google_migration_center_preference_set` ([#17291](https://github.com/hashicorp/terraform-provider-google/pull/17291))
* **New Resource:** `google_netapp_volume_replication` ([#17348](https://github.com/hashicorp/terraform-provider-google/pull/17348))

IMPROVEMENTS:
* cloudfunctions: added output-only `version_id` field on `google_cloudfunctions_function` ([#17273](https://github.com/hashicorp/terraform-provider-google/pull/17273))
* composer: supported patch versions of airflow on `google_composer_environment` ([#17345](https://github.com/hashicorp/terraform-provider-google/pull/17345))
* compute: supported updating `network_interface.stack_type` field on `google_compute_instance` resource. ([#17295](https://github.com/hashicorp/terraform-provider-google/pull/17295))
* container: added `node_config.resource_manager_tags` field to `google_container_cluster` resource ([#17346](https://github.com/hashicorp/terraform-provider-google/pull/17346))
* container: added `node_config.resource_manager_tags` field to `google_container_node_pool` resource ([#17346](https://github.com/hashicorp/terraform-provider-google/pull/17346))
* container: added output-only fields `membership_id` and  `membership_location` under `fleet` in `google_container_cluster` resource ([#17305](https://github.com/hashicorp/terraform-provider-google/pull/17305))
* looker: added `custom_domain` field to `google_looker_instance ` resource ([#17301](https://github.com/hashicorp/terraform-provider-google/pull/17301))
* netapp: added field `restore_parameters` and output-only fields `state`, `state_details` and `create_time` to `google_netapp_volume` resource ([#17293](https://github.com/hashicorp/terraform-provider-google/pull/17293))
* workbench: added `container_image` field to `google_workbench_instance` resource ([#17326](https://github.com/hashicorp/terraform-provider-google/pull/17326))
* workbench: added `shielded_instance_config` field to `google_workbench_instance` resource ([#17306](https://github.com/hashicorp/terraform-provider-google/pull/17306))

BUG FIXES:
* bigquery: allowed users to set permissions for `principal`/`principalSets` (`iamMember`) in `google_bigquery_dataset_iam_member`. ([#17292](https://github.com/hashicorp/terraform-provider-google/pull/17292))
* cloudfunctions2: fixed an issue where not specifying `event_config.trigger_region` in `google_cloudfunctions2_function` resulted in a permanent diff. The field now pulls a default value from the API when unset. ([#17328](https://github.com/hashicorp/terraform-provider-google/pull/17328))
* compute: fixed issue where changes only in `stateful_(internal|external)_ip` would not trigger an update for `google_compute_(region_)instance_group_manager` ([#17297](https://github.com/hashicorp/terraform-provider-google/pull/17297))
* compute: fixed perma-diff on `min_ports_per_vm` in `google_compute_router_nat` when the field is unset by making the field default to the API-set value ([#17337](https://github.com/hashicorp/terraform-provider-google/pull/17337))
* dataflow: fixed crash in `google_dataflox_job` to return an error instead if a job's Environment field is nil when reading job information ([#17344](https://github.com/hashicorp/terraform-provider-google/pull/17344))
* notebooks: changed `tag` field to default to the API's value if not specified in `google_notebooks_instance` ([#17323](https://github.com/hashicorp/terraform-provider-google/pull/17323))

## 5.17.0 (Feb 20, 2024)

NOTES:
* cloudbuildv2: changed underlying actuation engine for `google_cloudbuildv2_connection`, there should be no user-facing impact ([#17222](https://github.com/hashicorp/terraform-provider-google/pull/17222))

DEPRECATIONS:
* container: deprecated support for `relay_mode` field in `google_container_cluster.monitoring_config.advanced_datapath_observability_config` in favor of `enable_relay` field, `relay_mode` field will be removed in a future major release ([#17262](https://github.com/hashicorp/terraform-provider-google/pull/17262))

FEATURES:
* **New Resource:** `google_firebase_app_check_debug_token` ([#17242](https://github.com/hashicorp/terraform-provider-google/pull/17242))
* **New Resource:** `google_clouddeploy_custom_target_type` ([#17254](https://github.com/hashicorp/terraform-provider-google/pull/17254))

IMPROVEMENTS:
* cloudasset: allowed overriding the billing project for the `google_cloud_asset_resources_search_all` datasource
* clouddeploy: added support for `canary_revision_tags`, `prior_revision_tags`, `stable_revision_tags`, and `stable_cutback_duration` to `google_clouddeploy_delivery_pipeline`
* cloudfunctions: expose `version_id` on `google_cloudfunctions_function` ([#17273](https://github.com/hashicorp/terraform-provider-google/pull/17273))
* compute: promoted `user_ip_request_headers` field on `google_compute_security_policy` resource to GA ([#17271](https://github.com/hashicorp/terraform-provider-google/pull/17271))
* container: added support for `enable_relay` field to `google_container_cluster.monitoring_config.advanced_datapath_observability_config` ([#17262](https://github.com/hashicorp/terraform-provider-google/pull/17262))
* eventarc: added support for `http_endpoint.uri` and `network_config.network_attachment` to `google_eventarc_trigger` ([#17237](https://github.com/hashicorp/terraform-provider-google/pull/17237))
* healthcare: added `reject_duplicate_message` field to `google_healthcare_hl7_v2_store ` resource ([#17267](https://github.com/hashicorp/terraform-provider-google/pull/17267))
* identityplatform: added `client`, `permissions`, `monitoring` and `mfa` fields to `google_identity_platform_config` ([#17225](https://github.com/hashicorp/terraform-provider-google/pull/17225))
* notebooks: added `desired_state` field to `google_notebooks_instance` ([#17268](https://github.com/hashicorp/terraform-provider-google/pull/17268))
* vertexai: added `feature_registry_source` field to `google_vertex_ai_feature_online_store_featureview` resource ([#17264](https://github.com/hashicorp/terraform-provider-google/pull/17264))
* workbench: added `desired_state` field to `google_workbench_instance` resource ([#17270](https://github.com/hashicorp/terraform-provider-google/pull/17270))

BUG FIXES:
* compute: made `resource_manager_tags` updatable on `google_compute_instance_template` and `google_compute_region_instance_template` ([#17256](https://github.com/hashicorp/terraform-provider-google/pull/17256))
* notebooks: prevented recreation of `google_notebooks_instance` when `kms_key` or `service_account_scopes` are changed server-side ([#17232](https://github.com/hashicorp/terraform-provider-google/pull/17232))

## 5.16.0 (Feb 12, 2024)

FEATURES:
* **New Resource:** `google_clouddeploy_delivery_pipeline_iam_*` ([#17180](https://github.com/hashicorp/terraform-provider-google/pull/17180))
* **New Resource:** `google_compute_instance_group_membership` ([#17188](https://github.com/hashicorp/terraform-provider-google/pull/17188))
* **New Resource:** `google_discovery_engine_search_engine` ([#17146](https://github.com/hashicorp/terraform-provider-google/pull/17146))
* **New Resource:** `google_firebase_app_check_service_config` ([#17155](https://github.com/hashicorp/terraform-provider-google/pull/17155))

IMPROVEMENTS:
* bigquery: promoted `table_replication_info` field on `resource_bigquery_table` resource to GA ([#17181](https://github.com/hashicorp/terraform-provider-google/pull/17181))
* networksecurity: removed unused custom code from `google_network_security_address_group` ([#17183](https://github.com/hashicorp/terraform-provider-google/pull/17183))
* provider: added an optional provider level label `goog-terraform-provisioned` to identify resources that were created by Terraform when viewing/editing these resources in other tools. ([#17170](https://github.com/hashicorp/terraform-provider-google/pull/17170))

## 5.15.0 (Feb 5, 2024)

FEATURES:
* **New Data Source:** `google_compute_machine_types` ([#17107](https://github.com/hashicorp/terraform-provider-google/pull/17107))
* **New Resource:** `google_blockchain_node_engine_blockchain_nodes` ([#17096](https://github.com/hashicorp/terraform-provider-google/pull/17096))
* **New Resource:** `google_compute_region_network_endpoint` ([#17137](https://github.com/hashicorp/terraform-provider-google/pull/17137))
* **New Resource:** `google_discovery_engine_chat_engine` ([#17145](https://github.com/hashicorp/terraform-provider-google/pull/17145))
* **New Resource:** `google_discovery_engine_search_engine` ([#17146](https://github.com/hashicorp/terraform-provider-google/pull/17146))
* **New Resource:** `google_netapp_volume_snapshot` ([#17138](https://github.com/hashicorp/terraform-provider-google/pull/17138))

IMPROVEMENTS:
* compute: added `INTERNET_IP_PORT` and `INTERNET_FQDN_PORT` options for the `google_compute_region_network_endpoint_group` resource. ([#17137](https://github.com/hashicorp/terraform-provider-google/pull/17137))
* compute: added `creation_timestamp` to `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager`. ([#17110](https://github.com/hashicorp/terraform-provider-google/pull/17110))
* compute: added `disk_id` attribute to `google_compute_disk` resource ([#17112](https://github.com/hashicorp/terraform-provider-google/pull/17112))
* compute: added `stack_type` attribute for `google_compute_interconnect_attachment` resource. ([#17139](https://github.com/hashicorp/terraform-provider-google/pull/17139))
* compute: updated the `google_compute_security_policy` resource's `json_parsing` field to accept the value `STANDARD_WITH_GRAPHQL` ([#17097](https://github.com/hashicorp/terraform-provider-google/pull/17097))
* memcache: added `reserved_ip_range_id` field to `google_memcache_instance` resource ([#17101](https://github.com/hashicorp/terraform-provider-google/pull/17101))
* netapp: added `deletion_policy` field to `google_netapp_volume` resource ([#17111](https://github.com/hashicorp/terraform-provider-google/pull/17111))

BUG FIXES:
* alloydb: fixed an issue where `database_flags` in secondary `google_alloydb_instance` resources would cause a diff, as they are copied from the primary ([#17128](https://github.com/hashicorp/terraform-provider-google/pull/17128))
* filestore: made `google_filestore_instance.source_backup` field configurable ([#17099](https://github.com/hashicorp/terraform-provider-google/pull/17099))
* vmwareengine: fixed a bug to prevent recreation of existing [`google_vmwareengine_private_cloud`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/vmwareengine_private_cloud) resources when upgrading provider version from <5.10.0 ([#17135](https://github.com/hashicorp/terraform-provider-google/pull/17135)

## 5.14.0 (Jan 29, 2024)

FEATURES:
* **New Resource:** `google_discovery_engine_data_store` ([#17084](https://github.com/hashicorp/terraform-provider-google/pull/17084))
* **New Resource:** `google_securityposture_posture_deployment` ([#17085](https://github.com/hashicorp/terraform-provider-google/pull/17085))
* **New Resource:** `google_securityposture_posture` ([#17079](https://github.com/hashicorp/terraform-provider-google/pull/17079))

IMPROVEMENTS:
* artifactregistry: promoted `cleanup_policies` and `cleanup_policy_dry_run` fields to GA for `google_artifactregistry_repository` resource ([#17074](https://github.com/hashicorp/terraform-provider-google/pull/17074))
* composer: added `data_retention_config` field to `google_composer_environment` resource ([#17050](https://github.com/hashicorp/terraform-provider-google/pull/17050))
* logging: updated the `google_logging_project_bucket_config` resource to be created using the asynchronous create method ([#17067](https://github.com/hashicorp/terraform-provider-google/pull/17067))
* pubsub: added `use_table_schema` field to `google_pubsub_subscription` resource ([#17054](https://github.com/hashicorp/terraform-provider-google/pull/17054))
* workflows: added `call_log_level` field to `google_workflows_workflow` resource ([#17051](https://github.com/hashicorp/terraform-provider-google/pull/17051))

BUG FIXES:
* cloudfunctions2: fixed permadiff when `build_config.docker_repository` field is not specified on `google_cloudfunctions2_function` resource ([#17072](https://github.com/hashicorp/terraform-provider-google/pull/17072))
* compute: fixed error when `iap` field is unset for `google_compute_region_backend_service` resource ([#17071](https://github.com/hashicorp/terraform-provider-google/pull/17071))
* eventarc: fixed error when setting `destination.cloud_function` field on `google_eventarc_trigger` resource by making it output-only ([#17052](https://github.com/hashicorp/terraform-provider-google/pull/17052))


## 5.13.0 (Jan 22, 2024)

NOTES:
* cloudbuildv2: changed underlying actuation engine for `google_cloudbuildv2_repository`, there should be no user-facing impact ([#16969](https://github.com/hashicorp/terraform-provider-google/pull/16969))
* provider: added support for in-place update for `labels` and `terraform_labels` fields in immutable resources ([#17016](https://github.com/hashicorp/terraform-provider-google/pull/17016))

FEATURES:
* **New Resource:** `google_netapp_backup_policy` ([#16962](https://github.com/hashicorp/terraform-provider-google/pull/16962))
* **New Resource:** `google_netapp_volume` ([#16990](https://github.com/hashicorp/terraform-provider-google/pull/16990))
* **New Resource:** `google_network_security_address_group_iam_*` ([#17013](https://github.com/hashicorp/terraform-provider-google/pull/17013))
* **New Resource:** `google_vertex_ai_feature_group_feature` ([#17015](https://github.com/hashicorp/terraform-provider-google/pull/17015))

IMPROVEMENTS:
* alloydb: allowed `database_version` as an input on `google_alloydb_cluster` resource ([#16967](https://github.com/hashicorp/terraform-provider-google/pull/16967))
* bigquery: added `spark_options` field to `google_bigquery_routine` resource ([#17028](https://github.com/hashicorp/terraform-provider-google/pull/17028))
* cloudrunv2: added `nfs` and `gcs` fields to `google_cloud_run_v2_service.template.volumes` ([#16972](https://github.com/hashicorp/terraform-provider-google/pull/16972))
* cloudrunv2: added `tcp_socket` field to `google_cloud_run_v2.template.containers.liveness_probe` ([#16972](https://github.com/hashicorp/terraform-provider-google/pull/16972))
* compute: added `enable_confidential_compute` field to `google_compute_instance.boot_disk.initialize_params` ([#16968](https://github.com/hashicorp/terraform-provider-google/pull/16968))
* compute: added `enable_confidential_compute` field to `google_compute_disk` resource ([#16968](https://github.com/hashicorp/terraform-provider-google/pull/16968))
* gkehub2: added `clusterupgrade` field to `google_gke_hub_feature` resource ([#16951](https://github.com/hashicorp/terraform-provider-google/pull/16951))
* notebooks: allowed `machine_type` and `accelerator_config` to be updatable on `google_notebooks_runtime` resource ([#16993](https://github.com/hashicorp/terraform-provider-google/pull/16993))

BUG FIXES:
* compute: fixed the bug that `max_ttl` is sent in API calls even it is removed from configuration when changing cache_mode to FORCE_CACHE_ALL in `google_compute_backend_bucket` resource ([#16976](https://github.com/hashicorp/terraform-provider-google/pull/16976))
* networkservices: fixed a perma-diff on `addresses` field in `google_network_services_gateway` resource ([#17035](https://github.com/hashicorp/terraform-provider-google/pull/17035))
* provider: fixed `universe_domain` behavior to correctly throw an error when explicitly configured `universe_domain` values did not match credentials assumed to be in the default universe ([#17014](https://github.com/hashicorp/terraform-provider-google/pull/17014))
* spanner: fixed error when adding `autoscaling_config` to an existing `google_spanner_instance` resource ([#17033](https://github.com/hashicorp/terraform-provider-google/pull/17033))

## 5.12.0 (Jan 16, 2024)

FEATURES:
* **New Data Source:** `google_dns_managed_zones` ([#16949](https://github.com/hashicorp/terraform-provider-google/pull/16949))
* **New Data Source:** `google_filestore_instance` ([#16931](https://github.com/hashicorp/terraform-provider-google/pull/16931))
* **New Data Source:** `google_vmwareengine_external_access_rule` ([#16912](https://github.com/hashicorp/terraform-provider-google/pull/16912))
* **New Resource:** `google_clouddomains_registration` ([#16947](https://github.com/hashicorp/terraform-provider-google/pull/16947))
* **New Resource:** `google_netapp_kmsconfig` ([#16945](https://github.com/hashicorp/terraform-provider-google/pull/16945))
* **New Resource:** `google_vertex_ai_feature_online_store_featureview` ([#16930](https://github.com/hashicorp/terraform-provider-google/pull/16930))
* **New Resource:** `google_vmwareengine_external_access_rule` ([#16912](https://github.com/hashicorp/terraform-provider-google/pull/16912))

IMPROVEMENTS:
* compute: added `md5_authentication_key` field to `google_compute_router_peer` resource ([#16923](https://github.com/hashicorp/terraform-provider-google/pull/16923))
* compute: added in-place update support to `params.resource_manager_tags` field in `google_compute_instance` resource ([#16942](https://github.com/hashicorp/terraform-provider-google/pull/16942))
* compute: added in-place update support to `description` field in `google_compute_instance` resource ([#16900](https://github.com/hashicorp/terraform-provider-google/pull/16900))
* gkehub: added `policycontroller` field to `google_gke_hub_feature_membership` resource ([#16916](https://github.com/hashicorp/terraform-provider-google/pull/16916))
* gkehub2: added `clusterupgrade` field to `google_gke_hub_feature` resource ([#16951](https://github.com/hashicorp/terraform-provider-google/pull/16951))
* gkeonprem: added in-place update support to `vsphere_config` field and added `host_groups` field in `google_gkeonprem_vmware_node_pool` resource ([#16896](https://github.com/hashicorp/terraform-provider-google/pull/16896))
* iam: added `create_ignore_already_exists` field to `google_service_account` resource. If `ignore_create_already_exists` is set to true, resource creation would succeed when response error is 409 `ALREADY_EXISTS`. ([#16927](https://github.com/hashicorp/terraform-provider-google/pull/16927))
* servicenetworking: added field `deletion_policy` to `google_service_networking_connection` ([#16944](https://github.com/hashicorp/terraform-provider-google/pull/16944))
* sql: set `replica_configuration`, `ca_cert`, and `server_ca_cert` fields to be sensitive in `google_sql_instance` and `google_sql_ssl_cert` resources ([#16932](https://github.com/hashicorp/terraform-provider-google/pull/16932))

BUG FIXES:
* bigquery: fixed perma-diff of `encryption_configuration` when API returns an empty object on `google_bigquery_table` resource ([#16926](https://github.com/hashicorp/terraform-provider-google/pull/16926))
* compute: fixed an issue where the provider would `wait_for_instances` if set before deleting on `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` resources ([#16943](https://github.com/hashicorp/terraform-provider-google/pull/16943))
* compute: fixed perma-diff that reordered `stateful_external_ip` and `stateful_internal_ip` blocks on `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` resources ([#16910](https://github.com/hashicorp/terraform-provider-google/pull/16910))
* datapipeline: fixed perma-diff of `scheduler_service_account_email` when it's not explicitly specified in `google_data_pipeline_pipeline` resource ([#16917](https://github.com/hashicorp/terraform-provider-google/pull/16917))
* edgecontainer: fixed resource import on `google_edgecontainer_vpn_connection` resource ([#16948](https://github.com/hashicorp/terraform-provider-google/pull/16948))
* servicemanagement: fixed an issue where an inconsistent plan would be created when certain fields such as `openapi_config`, `grpc_config`, and `protoc_output_base64`, had computed values in `google_endpoints_service` resource ([#16946](https://github.com/hashicorp/terraform-provider-google/pull/16946))
* storage: fixed an issue where retry timeout wasn't being utilized when creating `google_storage_bucket` resource ([#16902](https://github.com/hashicorp/terraform-provider-google/pull/16902))

## 5.11.0 (Jan 08, 2024)

NOTES:
* compute: changed underlying actuation engine for `google_network_firewall_policy` and `google_region_network_firewall_policy`, there should be no user-facing impact ([#16837](https://github.com/hashicorp/terraform-provider-google/pull/16837))

DEPRECATIONS:
* gkehub2: deprecated field `configmanagement.config_sync.oci.version` in `google_gke_hub_feature` resource ([#16818](https://github.com/hashicorp/terraform-provider-google/pull/16818))

FEATURES:
* **New Data Source:** `google_compute_reservation` ([#16860](https://github.com/hashicorp/terraform-provider-google/pull/16860))
* **New Resource:** `google_integration_connectors_endpoint_attachment` ([#16822](https://github.com/hashicorp/terraform-provider-google/pull/16822))
* **New Resource:** `google_logging_folder_settings` ([#16800](https://github.com/hashicorp/terraform-provider-google/pull/16800))
* **New Resource:** `google_logging_organization_settings` ([#16800](https://github.com/hashicorp/terraform-provider-google/pull/16800))
* **New Resource:** `google_netapp_active_directory` ([#16844](https://github.com/hashicorp/terraform-provider-google/pull/16844))
* **New Resource:** `google_vertex_ai_feature_online_store` ([#16840](https://github.com/hashicorp/terraform-provider-google/pull/16840))
* **New Resource:** `google_vertex_ai_feature_group` ([#16842](https://github.com/hashicorp/terraform-provider-google/pull/16842))
* **New Resource:** `google_netapp_backup_vault` ([#16876](https://github.com/hashicorp/terraform-provider-google/pull/16876))

IMPROVEMENTS:
* bigqueryanalyticshub: added `restricted_export_config` field to `google_bigquery_analytics_hub_listing ` resource ([#16850](https://github.com/hashicorp/terraform-provider-google/pull/16850))
* composer: added support for `composer_internal_ipv4_cidr_block` field to `google_composer_environment` ([#16815](https://github.com/hashicorp/terraform-provider-google/pull/16815))
* compute: added `provisioned_iops`and `provisioned_throughput` fields under `boot_disk.initialize_params` to `google_compute_instance` resource ([#16871](https://github.com/hashicorp/terraform-provider-google/pull/16871))
* compute: added `resource_manager_tags` and `disk.resource_manager_tags` for `google_compute_instance_template` ([#16889](https://github.com/hashicorp/terraform-provider-google/pull/16889))
* compute: added `resource_manager_tags` and `disk.resource_manager_tags` for `google_compute_region_instance_template` ([#16889](https://github.com/hashicorp/terraform-provider-google/pull/16889))
* dataproc: added `auxiliary_node_groups` field to `google_dataproc_cluster` resource ([#16798](https://github.com/hashicorp/terraform-provider-google/pull/16798))
* edgecontainer: increased default timeout on `google_edgecontainer_cluster`, `google_edgecontainer_node_pool` to 480m from 60m ([#16886](https://github.com/hashicorp/terraform-provider-google/pull/16886))
* gkehub2: added field `version` under `configmanagement` in `google_gke_hub_feature` resource ([#16818](https://github.com/hashicorp/terraform-provider-google/pull/16818))
* kms: added output-only field `primary` to `google_kms_crypto_key` ([#16845](https://github.com/hashicorp/terraform-provider-google/pull/16845))
* metastore: added `endpoint_protocol`, `metadata_integration`, and `auxiliary_versions` to `google_dataproc_metastore_service` ([#16823](https://github.com/hashicorp/terraform-provider-google/pull/16823))
* sql: added support for IAM GROUP authentication in the `type` field of `google_sql_user` ([#16853](https://github.com/hashicorp/terraform-provider-google/pull/16853))
* storagetransfer: made `name` field settable on `google_storage_transfer_job` ([#16838](https://github.com/hashicorp/terraform-provider-google/pull/16838))

BUG FIXES:
* container: added check that `node_version` and `min_master_version` are the same on create of `google_container_cluster`, when running terraform plan ([#16817](https://github.com/hashicorp/terraform-provider-google/pull/16817))
* container: fixed a bug where disabling PDCSI addon `gce_persistent_disk_csi_driver_config` during creation will result in permadiff in `google_container_cluster` resource ([#16794](https://github.com/hashicorp/terraform-provider-google/pull/16794))
* container: fixed an issue in which migrating from the deprecated Binauthz enablement bool to the new evaluation mode enum inadvertently caused two cluster update events, instead of none. ([#16851](https://github.com/hashicorp/terraform-provider-google/pull/16851))
* containerattached: fixed crash when updating a cluster to remove `admin_users` or `admin_groups` in `google_container_attached_cluster` ([#16852](https://github.com/hashicorp/terraform-provider-google/pull/16852))
* dialogflowcx: fixed a permadiff in the `git_integration_settings` field of `google_diagflow_cx_agent` ([#16803](https://github.com/hashicorp/terraform-provider-google/pull/16803))
* monitoring: fixed the index out of range crash in `dashboard_json` for the resource `google_monitoring_dashboard` ([#16792](https://github.com/hashicorp/terraform-provider-google/pull/16792))

## 5.10.0 (Dec 18, 2023)

FEATURES:
* **New Data Source:** `google_compute_region_disk` ([#16732](https://github.com/hashicorp/terraform-provider-google/pull/16732))
* **New Data Source:** `google_vmwareengine_external_address` ([#16698](https://github.com/hashicorp/terraform-provider-google/pull/16698))
* **New Data Source:** `google_vmwareengine_subnet` ([#16700](https://github.com/hashicorp/terraform-provider-google/pull/16700))
* **New Data Source:** `google_vmwareengine_vcenter_credentials` ([#16709](https://github.com/hashicorp/terraform-provider-google/pull/16709))
* **New Resource:** `google_vmwareengine_cluster` ([#16757](https://github.com/hashicorp/terraform-provider-google/pull/16757))
* **New Resource:** `google_vmwareengine_external_address` ([#16698](https://github.com/hashicorp/terraform-provider-google/pull/16698))
* **New Resource:** `google_vmwareengine_subnet` ([#16700](https://github.com/hashicorp/terraform-provider-google/pull/16700))
* **New Resource:** `google_workbench_instance` ([#16773](https://github.com/hashicorp/terraform-provider-google/pull/16773))
* **New Resource:** `google_workbench_instance_iam_*` ([#16773](https://github.com/hashicorp/terraform-provider-google/pull/16773))

IMPROVEMENTS:
* compute: added `numeric_id` field to `google_compute_network` resource ([#16712](https://github.com/hashicorp/terraform-provider-google/pull/16712))
* compute: added `remove_instance_on_destroy` option to `google_compute_per_instance_config` resource ([#16729](https://github.com/hashicorp/terraform-provider-google/pull/16729))
* compute: added `remove_instance_on_destroy` option to `google_compute_region_per_instance_config` resource ([#16729](https://github.com/hashicorp/terraform-provider-google/pull/16729))
* container: added `network_performance_config` field to `google_container_node_pool` resource to support GKE tier 1 networking ([#16688](https://github.com/hashicorp/terraform-provider-google/pull/16688))
* container: added support for in-place update for `machine_type`/`disk_type`/`disk_size_gb` in `google_container_node_pool` resource ([#16724](https://github.com/hashicorp/terraform-provider-google/pull/16724))
* containerazure: added `config.labels` to `google_container_azure_node_pool` ([#16754](https://github.com/hashicorp/terraform-provider-google/pull/16754))
* dataform: added `display_name`, `labels` and `npmrc_environment_variables_secret_version` fields to `google_dataform_repository` resource ([#16733](https://github.com/hashicorp/terraform-provider-google/pull/16733))
* monitoring: added `severity` field to `google_monitoring_alert_policy` resource ([#16775](https://github.com/hashicorp/terraform-provider-google/pull/16775))
* notebooks: added support for `labels` to `google_notebooks_runtime` ([#16783](https://github.com/hashicorp/terraform-provider-google/pull/16783))
* recaptchaenterprise: added `waf_settings` to `google_recaptcha_enterprise_key` ([#16754](https://github.com/hashicorp/terraform-provider-google/pull/16754))
* securesourcemanager: added `host_config`, `state_note`, `kms_key`, and `private_config` fields to `google_secure_source_manager_instance` resource ([#16731](https://github.com/hashicorp/terraform-provider-google/pull/16731))
* spanner: added `autoscaling_config.max_nodes` and `autoscaling_config.min_nodes` to `google_spanner_instance` ([#16786](https://github.com/hashicorp/terraform-provider-google/pull/16786))
* storage: added `rpo` field to `google_storage_bucket` resource ([#16756](https://github.com/hashicorp/terraform-provider-google/pull/16756))
* vmwareengine: added `type` field to `google_vmwareengine_private_cloud` resource ([#16781](https://github.com/hashicorp/terraform-provider-google/pull/16781))
* workloadidentity: added `saml` block to `google_iam_workload_identity_pool_provider` resource ([#16710](https://github.com/hashicorp/terraform-provider-google/pull/16710))

BUG FIXES:
* logging: fixed an issue where value change of `unique_writer_identity` on `google_logging_project_sink` does not trigger diff on dependent's usages of `writer_identity`  ([#16776](https://github.com/hashicorp/terraform-provider-google/pull/16776))

## 5.9.0 (Dec 11, 2023)

FEATURES:
* **New Data Source:** `google_logging_folder_settings` ([#16658](https://github.com/hashicorp/terraform-provider-google/pull/16658))
* **New Data Source:** `google_logging_organization_settings` ([#16658](https://github.com/hashicorp/terraform-provider-google/pull/16658))
* **New Data Source:** `google_logging_project_settings` ([#16658](https://github.com/hashicorp/terraform-provider-google/pull/16658))
* **New Data Source:** `google_vmwareengine_network_policy` ([#16639](https://github.com/hashicorp/terraform-provider-google/pull/16639))
* **New Data Source:** `google_vmwareengine_nsx_credentials` ([#16669](https://github.com/hashicorp/terraform-provider-google/pull/16669))
* **New Resource:** `google_scc_event_threat_detection_custom_module` ([#16649](https://github.com/hashicorp/terraform-provider-google/pull/16649))
* **New Resource:** `google_secure_source_manager_instance` ([#16637](https://github.com/hashicorp/terraform-provider-google/pull/16637))
* **New Resource:** `google_vmwareengine_network_policy` ([#16639](https://github.com/hashicorp/terraform-provider-google/pull/16639))

IMPROVEMENTS:
* bigqueryconnection: added `spark` support to `google_bigquery_connection` resource ([#16677](https://github.com/hashicorp/terraform-provider-google/pull/16677))
* cloudidentity: added `expiry_detail` field to `google_cloud_identity_group_membership` resource ([#16643](https://github.com/hashicorp/terraform-provider-google/pull/16643))
* container: added `autoscaling_profile` field in the `cluster_autoscaling` block in `google_container_cluster` resource ([#16653](https://github.com/hashicorp/terraform-provider-google/pull/16653))
* gkehub: added `default_cluster_config` field to `google_gke_hub_fleet` resource  ([#16630](https://github.com/hashicorp/terraform-provider-google/pull/16630))
* gkehub: added `binary_authorization_config` field to `google_gke_hub_fleet` resource ([#16674](https://github.com/hashicorp/terraform-provider-google/pull/16674))
* sql: added support for in-place updates to the `edition` field in `google_sql_database_instance` resource ([#16629](https://github.com/hashicorp/terraform-provider-google/pull/16629))

BUG FIXES:
* artifactregistry: fixed permadiff due to unsorted `virtual_repository_config` array in `google_artifact_registry_repository` ([#16646](https://github.com/hashicorp/terraform-provider-google/pull/16646))
* container: made `dns_config` field updatable on `google_container_cluster` resource ([#16652](https://github.com/hashicorp/terraform-provider-google/pull/16652))
* dlp: added conflicting field validation in the `storage_config.timespan_config` block in `data_loss_prevention_job_trigger` resource ([#16628](https://github.com/hashicorp/terraform-provider-google/pull/16628))
* dlp: updated the `storage_config.timespan_config.timestamp_field` field in `data_loss_prevention_job_trigger` to be optional ([#16628](https://github.com/hashicorp/terraform-provider-google/pull/16628))
* firestore: added retries during creation of `google_firestore_index` resources to address retryable 409 code API errors ("Please retry, underlying data changed", and "Aborted due to cross-transaction contention") ([#16618](https://github.com/hashicorp/terraform-provider-google/pull/16618), [#16670](https://github.com/hashicorp/terraform-provider-google/pull/16670))
* storage: fixed unexpected `lifecycle_rule` conditions being added for `google_storage_bucket` ([#16683](https://github.com/hashicorp/terraform-provider-google/pull/16683))

## 5.8.0 (Dec 4, 2023)

FEATURES:
* **New Data Source:** `google_vmwareengine_network_peering` ([#16616](https://github.com/hashicorp/terraform-provider-google/pull/16616))
* **New Resource:** `google_migration_center_group` ([#16549](https://github.com/hashicorp/terraform-provider-google/pull/16549))
* **New Resource:** `google_netapp_storage_pool` ([#16573](https://github.com/hashicorp/terraform-provider-google/pull/16573))
* **New Resource:** `google_vmwareengine_network` (ga) ([#16583](https://github.com/hashicorp/terraform-provider-google/pull/16583))
* **New Resource:** `google_vmwareengine_network_peering` ([#16616](https://github.com/hashicorp/terraform-provider-google/pull/16616))

IMPROVEMENTS:
* artifactregistry: added `remote_repository_config.upstream_credentials` field to `google_artifact_registry_repository` resource ([#16562](https://github.com/hashicorp/terraform-provider-google/pull/16562))
* cloudbuild: added fields `build.artifacts.maven_artifacts`, `build.artifacts.npm_packages `, and `build.artifacts.python_packages ` to resource `google_cloudbuild_trigger` ([#16543](https://github.com/hashicorp/terraform-provider-google/pull/16543))
* cloudrunv2: promoted field `depends_on` in `google_cloud_run_v2_service` to GA ([#16577](https://github.com/hashicorp/terraform-provider-google/pull/16577))
* composer: added `database_config.zone` field in `google_composer_environment` ([#16551](https://github.com/hashicorp/terraform-provider-google/pull/16551))
* compute: added field `service_directory_registrations` to resource `google_compute_global_forwarding_rule` ([#16581](https://github.com/hashicorp/terraform-provider-google/pull/16581))
* firestore: added virtual field `deletion_policy` to `google_firestore_database` ([#16576](https://github.com/hashicorp/terraform-provider-google/pull/16576))
* firestore: enabled database deletion upon destroy for `google_firestore_database` ([#16576](https://github.com/hashicorp/terraform-provider-google/pull/16576))
* gkehub2: added `policycontroller` field to `fleet_default_member_config` in `google_gke_hub_feature` ([#16542](https://github.com/hashicorp/terraform-provider-google/pull/16542))
* iam: added `allowed_services`, `disable_programmatic_signin` fields to `google_iam_workforce_pool` resource ([#16580](https://github.com/hashicorp/terraform-provider-google/pull/16580))
* vmwareengine: added `STANDARD` type support to `google_vmwareengine_network` resource ([#16583](https://github.com/hashicorp/terraform-provider-google/pull/16583))
* vmwareengine: promoted `google_vmwareengine_private_cloud` resource to GA ([#16613](https://github.com/hashicorp/terraform-provider-google/pull/16613))

BUG FIXES:
* compute: fixed a permadiff caused by issues with ipv6 diff suppression in `google_compute_forwarding_rule` and `google_compute_global_forwarding_rule` ([#16550](https://github.com/hashicorp/terraform-provider-google/pull/16550))
* firestore: fixed an issue where `google_firestore_database` could be deleted when `delete_protection_state` was `DELETE_PROTECTION_ENABLED` ([#16576](https://github.com/hashicorp/terraform-provider-google/pull/16576))
* firestore: made resource creation retry for 409 errors with the text "Aborted due to cross-transaction contention" in `google_firestore_index ` ([#16618](https://github.com/hashicorp/terraform-provider-google/pull/16618))

## 5.7.0 (Nov 20, 2023)

DEPRECATIONS:
* gkehub: deprecated `config_management.binauthz` in `google_gke_hub_feature_membership` ([#16536](https://github.com/hashicorp/terraform-provider-google/pull/16536))

IMPROVEMENTS:
* bigtable: added `standard_isolation` and `standard_isolation.priority` fields to `google_bigtable_app_profile` resource ([#16485](https://github.com/hashicorp/terraform-provider-google/pull/16485))
* cloudrunv2: promoted `custom_audiences` field to GA on `google_cloud_run_v2_service` resource ([#16510](https://github.com/hashicorp/terraform-provider-google/pull/16510))
* compute: promoted `labels` field to GA on `google_compute_vpn_tunnel` resource ([#16508](https://github.com/hashicorp/terraform-provider-google/pull/16508))
* containerattached: added `proxy_config` field to `google_container_attached_cluster` resource ([#16524](https://github.com/hashicorp/terraform-provider-google/pull/16524))
* gkehub: added `membership_location` field to `google_gke_hub_feature_membership` resource ([#16536](https://github.com/hashicorp/terraform-provider-google/pull/16536))
* logging: made the change to aqcuire and update the `google_logging_project_sink` resource that already exists at the desired location. These logging buckets cannot be removed so deleting this resource will remove the bucket config from your terraform state but will leave the logging bucket unchanged. ([#16513](https://github.com/hashicorp/terraform-provider-google/pull/16513))
* memcache: added `MEMCACHE_1_6_15` as a possible value for `memcache_version` in `google_memcache_instance` resource ([#16531](https://github.com/hashicorp/terraform-provider-google/pull/16531))
* monitoring: added error message to delete Alert Policies first on 400 response when deleting `google_monitoring_uptime_check_config` resource ([#16535](https://github.com/hashicorp/terraform-provider-google/pull/16535))
* spanner: added `autoscaling_config` field to `google_spanner_instance` resource ([#16473](https://github.com/hashicorp/terraform-provider-google/pull/16473))
* workflows: promoted `user_env_vars` field to GA on `google_workflows_workflow` resource ([#16477](https://github.com/hashicorp/terraform-provider-google/pull/16477))

BUG FIXES:
* compute: changed `external_ipv6_prefix` field to not be output only in `google_compute_subnetwork` resource ([#16480](https://github.com/hashicorp/terraform-provider-google/pull/16480))
* compute: fixed issue where `google_compute_attached_disk` would produce an error for certain zone configs ([#16484](https://github.com/hashicorp/terraform-provider-google/pull/16484))
* edgecontainer: fixed update method of `google_edgecontainer_cluster` resource ([#16490](https://github.com/hashicorp/terraform-provider-google/pull/16490))
* provider: fixed an issue where universe domains would not overwrite API endpoints ([#16521](https://github.com/hashicorp/terraform-provider-google/pull/16521))
* resourcemanager: made `data_source_google_project_service` no longer return an error when the service is not enabled ([#16525](https://github.com/hashicorp/terraform-provider-google/pull/16525))
* sql: `ssl_mode` field is not stored in terraform state if it has never been used in `google_sql_database_instance` resource ([#16486](https://github.com/hashicorp/terraform-provider-google/pull/16486))
  
NOTES:
* dataproc: backfilled `terraform_labels` field for resource `google_dataproc_workflow_template`, so resource recreation won't happen during provider upgrade from `4.x` to `5.7` ([#16517](https://github.com/hashicorp/terraform-provider-google/pull/16517))
* * provider: backfilled `terraform_labels` field for some immutable resources, so resource recreation won't happen during provider upgrade from `4.X` to `5.7` ([#16518](https://github.com/hashicorp/terraform-provider-google/pull/16518))

## 5.6.0 (Nov 13, 2023)

FEATURES:
* **New Resource:** `google_integration_connectors_connection` ([#16468](https://github.com/hashicorp/terraform-provider-google/pull/16468))

IMPROVEMENTS:
* assuredworkloads: added `enable_sovereign_controls`, `partner`, `partner_permissions`, `violation_notifications_enabled`, and several other output-only fields to `google_assured_workloads_workloads` ([#16433](https://github.com/hashicorp/terraform-provider-google/pull/16433))
* composer: added `storage_config` to `google_composer_environment` ([#16455](https://github.com/hashicorp/terraform-provider-google/pull/16455))
* container: added `fleet` field to `google_container_cluster` resource ([#16466](https://github.com/hashicorp/terraform-provider-google/pull/16466))
* containeraws: added `admin_groups` to `google_container_aws_cluster` ([#16433](https://github.com/hashicorp/terraform-provider-google/pull/16433))
* containerazure: added `admin_groups` to `google_container_azure_cluster` ([#16433](https://github.com/hashicorp/terraform-provider-google/pull/16433))
* dataproc: added support for `instance_flexibility_policy` in `google_dataproc_cluster` ([#16417](https://github.com/hashicorp/terraform-provider-google/pull/16417))
* dialogflowcx: added `is_default_start_flow` field to `google_dialogflow_cx_flow` resource to allow management of default flow resources via Terraform ([#16441](https://github.com/hashicorp/terraform-provider-google/pull/16441))
* dialogflowcx: added `is_default_welcome_intent` and `is_default_negative_intent` fields to `google_dialogflow_cx_intent` resource to allow management of default intent resources via Terraform ([#16441](https://github.com/hashicorp/terraform-provider-google/pull/16441))
* * gkehub: added `fleet_default_member_config` field to `google_gke_hub_feature` resource ([#16457](https://github.com/hashicorp/terraform-provider-google/pull/16457))
* gkehub: added `metrics_gcp_service_account_email` to `google_gke_hub_feature_membership` ([#16433](https://github.com/hashicorp/terraform-provider-google/pull/16433))
* logging: added `index_configs` field to `logging_bucket_config` resource ([#16437](https://github.com/hashicorp/terraform-provider-google/pull/16437))
* logging: added `index_configs` field to `logging_project_bucket_config` resource ([#16437](https://github.com/hashicorp/terraform-provider-google/pull/16437))
* monitoring: added `pings_count`, `user_labels`, and `custom_content_type` fields to `google_monitoring_uptime_check_config` resource ([#16420](https://github.com/hashicorp/terraform-provider-google/pull/16420))
* spanner: added `autoscaling_config` field to  `google_spanner_instance` ([#16473](https://github.com/hashicorp/terraform-provider-google/pull/16473))
* sql: added `ssl_mode` field to `google_sql_database_instance` resource ([#16394](https://github.com/hashicorp/terraform-provider-google/pull/16394))
* vertexai: added `private_service_connect_config` to `google_vertex_ai_index_endpoint` ([#16471](https://github.com/hashicorp/terraform-provider-google/pull/16471))
* workstations: added `domain_config` field to resource `google_workstations_workstation_cluster` (beta) ([#16464](https://github.com/hashicorp/terraform-provider-google/pull/16464))

BUG FIXES:
* assuredworkloads: made the `violation_notifications_enabled` field on the `google_assured_workloads_workload` resource default to values returned from the API when unset in a users configuration ([#16465](https://github.com/hashicorp/terraform-provider-google/pull/16465))
* provider: made `terraform_labels` immutable in immutable resources to not block the upgrade. This will create a Terraform plan that recreates the resource on `4.X` -> `5.6.0` upgrade for affected resources. A mitigation to backfill the values during the upgrade is planned, and will release resource-by-resource. ([#16469](https://github.com/hashicorp/terraform-provider-google/pull/16469))

## 5.5.0 (Nov 06, 2023)

FEATURES:
* **New Data Source:** `google_bigquery_dataset` ([#16368](https://github.com/hashicorp/terraform-provider-google/pull/16368))

IMPROVEMENTS:
* alloydb: added `SECONDARY` as an option for `instance_type` field in `google_alloydb_instance` resource, to support creation of secondary instance inside a secondary cluster. ([#16398](https://github.com/hashicorp/terraform-provider-google/pull/16398))
* alloydb: added `deletion_policy` field to `google_alloydb_cluster` resource, to allow force-destroying instances along with their cluster. This is necessary to delete secondary instances, which cannot be deleted otherwise. ([#16398](https://github.com/hashicorp/terraform-provider-google/pull/16398))
* alloydb: added support to promote `google_alloydb_cluster` resources from secondary to primary ([#16413](https://github.com/hashicorp/terraform-provider-google/pull/16413))
* alloydb: increased default timeout on `google_alloydb_instance` to 120m from 40m ([#16398](https://github.com/hashicorp/terraform-provider-google/pull/16398))
* dataproc: added `instance_flexibility_policy` field ro `google_dataproc_cluster` resource ([#16417](https://github.com/hashicorp/terraform-provider-google/pull/16417))
* monitoring: added `subject` field to `google_monitoring_alert_policy` resource ([#16414](https://github.com/hashicorp/terraform-provider-google/pull/16414))
* storage: added `enable_object_retention` field to `google_storage_bucket` resource ([#16412](https://github.com/hashicorp/terraform-provider-google/pull/16412))
* storage: added `retention` field to `google_storage_bucket_object` resource ([#16412](https://github.com/hashicorp/terraform-provider-google/pull/16412))

BUG FIXES:
* firestore: fixed an issue with creation of multiple `google_firestore_field` resources ([#16372](https://github.com/hashicorp/terraform-provider-google/pull/16372))

## 5.4.0 (Oct 30, 2023)

DEPRECATIONS:
* bigquery: deprecated `cloud_spanner.use_serverless_analytics` on `google_bigquery_connection`. Use `cloud_spanner.use_data_boost` instead. ([#16310](https://github.com/hashicorp/terraform-provider-google/pull/16310))

NOTES:
* provider: added `universe_domain` attribute as a provider attribute ([#16323](https://github.com/hashicorp/terraform-provider-google/pull/16323))

BREAKING CHANGES:
* cloudrunv2: marked `location` field as required in resource `google_cloud_run_v2_job`. Any configuration without `location` specified will cause resource creation fail ([#16311](https://github.com/hashicorp/terraform-provider-google/pull/16311))
* cloudrunv2: marked `location` field as required in resource `google_cloud_run_v2_service`. Any configuration without `location` specified will cause resource creation fail ([#16311](https://github.com/hashicorp/terraform-provider-google/pull/16311))

FEATURES:
* **New Data Source:** `google_cloud_identity_group_lookup` ([#16296](https://github.com/hashicorp/terraform-provider-google/pull/16296))
* **New Resource:** `google_network_connectivity_policy_based_route` ([#16326](https://github.com/hashicorp/terraform-provider-google/pull/16326))
* **New Resource:** `google_pubsub_schema_iam_*` ([#16301](https://github.com/hashicorp/terraform-provider-google/pull/16301))

IMPROVEMENTS:
* accesscontextmanager: added support for specifying `vpc_network_sources` to `google_access_context_manager_access_levels`, `google_access_context_manager_access_level`, and `google_access_context_manager_access_level_condition` ([#16327](https://github.com/hashicorp/terraform-provider-google/pull/16327))
* apigee: added support for `type` in `google_apigee_environment` ([#16349](https://github.com/hashicorp/terraform-provider-google/pull/16349))
* bigquery: added `cloud_spanner.database_role`, `cloud_spanner.use_data_boost`, and `cloud_spanner.max_parallelism` fields to `google_bigquery_connection` ([#16310](https://github.com/hashicorp/terraform-provider-google/pull/16310))
* bigquery: added support for `iam_member` to `google_bigquery_dataset.access` ([#16322](https://github.com/hashicorp/terraform-provider-google/pull/16322))
* container: promoted field `identity_service_config` in `google_container_cluster` to GA ([#16305](https://github.com/hashicorp/terraform-provider-google/pull/16305))
* container: added update support for `google_container_node_pool.node_config.taint` ([#16306](https://github.com/hashicorp/terraform-provider-google/pull/16306))
* containerattached: added `admin_groups` field to `google_container_attached_cluster` resource ([#16307](https://github.com/hashicorp/terraform-provider-google/pull/16307))
* dialogflowcx: added `advanced_settings` field to `google_dialogflow_cx_flow` resource ([#16315](https://github.com/hashicorp/terraform-provider-google/pull/16315))
* dialogflowcx: added `advanced_settings` fields to `google_dialogflow_cx_page` resource ([#16315](https://github.com/hashicorp/terraform-provider-google/pull/16315))
* dialogflowcx: added `advanced_settings`, `text_to_speech_settings`, `git_integration_settings` fields to `google_dialogflow_cx_agent` resource ([#16315](https://github.com/hashicorp/terraform-provider-google/pull/16315))

BUG FIXES:
* bigquery: fixed a bug when updating a `google_bigquery_dataset` that contained an `iamMember` access rule added out of band with Terraform ([#16322](https://github.com/hashicorp/terraform-provider-google/pull/16322))
* bigqueryreservation: fixed bug of incorrect resource recreation when `capacity_commitment_id` is unspecified in resource `google_bigquery_capacity_commitment` ([#16320](https://github.com/hashicorp/terraform-provider-google/pull/16320))
* cloudrunv2: made `annotations` field on the `google_cloud_run_v2_job` data source include all annotations present on the resource in GCP ([#16300](https://github.com/hashicorp/terraform-provider-google/pull/16300))
* cloudrunv2: made `annotations` field on the `google_cloud_run_v2_service` data source include all annotations present on the resource in GCP ([#16300](https://github.com/hashicorp/terraform-provider-google/pull/16300))
* cloudrunv2: made `labels` and `terraform labels` fields on the `google_cloud_run_v2_job` data source include all annotations present on the resource in GCP ([#16300](https://github.com/hashicorp/terraform-provider-google/pull/16300))
* cloudrunv2: made `labels` and `terraform labels` fields on the `google_cloud_run_v2_service` data source include all annotations present on the resource in GCP ([#16300](https://github.com/hashicorp/terraform-provider-google/pull/16300))
* edgecontainer: fixed an issue where the update endpoint for `google_edgecontainer_cluster` was incorrect. ([#16347](https://github.com/hashicorp/terraform-provider-google/pull/16347))
* redis: allow `replica_count` to be set to zero in the `google_redis_cluster` resource ([#16302](https://github.com/hashicorp/terraform-provider-google/pull/16302))

## 5.3.0 (Oct 23, 2023)

DEPRECATIONS:
* bigquery: deprecated `time_partitioning.require_partition_filter` in favor of new top level field `require_partition_filter` in resource `google_bigquery_table` ([#16238](https://github.com/hashicorp/terraform-provider-google/pull/16238))

FEATURES:
* **New Data Source:** `google_cloud_run_v2_job` ([#16260](https://github.com/hashicorp/terraform-provider-google/pull/16260))
* **New Data Source:** `google_cloud_run_v2_service` ([#16290](https://github.com/hashicorp/terraform-provider-google/pull/16290))
* **New Data Source:** `google_compute_networks` ([#16240](https://github.com/hashicorp/terraform-provider-google/pull/16240))
* **New Resource:** `google_org_policy_custom_constraint` ([#16220](https://github.com/hashicorp/terraform-provider-google/pull/16220))

IMPROVEMENTS:
* cloudidentity: added `additional_group_keys` attribute to `google_cloud_identity_group` resource ([#16250](https://github.com/hashicorp/terraform-provider-google/pull/16250))
* composer: promoted `config.0.workloads_config.0.triggerer` to GA in resource `google_composer_environment` ([#16218](https://github.com/hashicorp/terraform-provider-google/pull/16218))
* compute: added `internal_ipv6_range` to `google_compute_network` data source and `internal_ipv6_prefix` field to `google_compute_subnetwork` data source ([#16267](https://github.com/hashicorp/terraform-provider-google/pull/16267))
* container: added support for `security_posture_config.vulnerability_mode` value `VULNERABILITY_ENTERPRISE`in `google_container_cluster` ([#16283](https://github.com/hashicorp/terraform-provider-google/pull/16283))
* dataform: added `ssh_authentication_config` and `service_account` to `google_dataform_repository` resource ([#16205](https://github.com/hashicorp/terraform-provider-google/pull/16205))
* dataproc: added `min_num_instances` field to `google_dataproc_cluster` resource ([#16249](https://github.com/hashicorp/terraform-provider-google/pull/16249))
* gkeonprem: promoted `google_gkeonprem_bare_metal_admin_cluster`, `google_gkeonprem_bare_metal_cluster`, and `google_gkeonprem_bare_metal_node_pool` resources to GA ([#16237](https://github.com/hashicorp/terraform-provider-google/pull/16237))
* gkeonprem: promoted `google_gkeonprem_vmware_cluster` and `google_gkeonprem_vmware_node_pool` resources to GA ([#16237](https://github.com/hashicorp/terraform-provider-google/pull/16237))
* logging: added `custom_writer_identity` field to `google_logging_project_sink` ([#16216](https://github.com/hashicorp/terraform-provider-google/pull/16216))
* secretmanager: made `ttl` field mutable in `google_secret_manager_secret` ([#16285](https://github.com/hashicorp/terraform-provider-google/pull/16285))
* storage: added `terminal_storage_class` to the `autoclass` field in `google_storage_bucket` resource ([#16282](https://github.com/hashicorp/terraform-provider-google/pull/16282))

BUG FIXES:
* bigquerydatatransfer: fixed an error when updating `google_bigquery_data_transfer_config` related to incorrect update masks ([#16269](https://github.com/hashicorp/terraform-provider-google/pull/16269))
* compute: fixed an error during the deletion when post was set to 0 on `google_compute_global_network_endpoint` ([#16286](https://github.com/hashicorp/terraform-provider-google/pull/16286))
* compute: fixed an issue with TTLs being sent for `google_compute_backend_service` when `cache_mode` is set to `USE_ORIGIN_HEADERS` ([#16245](https://github.com/hashicorp/terraform-provider-google/pull/16245))
* container: fixed an issue where empty `autoscaling` block would crash the provider for `google_container_node_pool` ([#16212](https://github.com/hashicorp/terraform-provider-google/pull/16212))
* dataflow: fixed a bug where resource updates returns an error if only `labels` has changes for batch `google_dataflow_job` and `google_dataflow_flex_template_job` ([#16248](https://github.com/hashicorp/terraform-provider-google/pull/16248))
* dialogflowcx: fixed updating `google_dialogflow_cx_version`; updates will no longer time out. ([#16214](https://github.com/hashicorp/terraform-provider-google/pull/16214))
* sql: fixed a bug where adding the `edition` field to a `google_sql_database_instance` resource that already existed and used ENTERPRISE edition resulted in a permant diff in plans ([#16215](https://github.com/hashicorp/terraform-provider-google/pull/16215))
* sql: removed host validation to support IP address and DNS address in host in `google_sql_source_representation_instance` resource ([#16235](https://github.com/hashicorp/terraform-provider-google/pull/16235))

## 5.2.0 (Oct 16, 2023)

FEATURES:
* **New Data Source:** `google_secret_manager_secrets` ([#16182](https://github.com/hashicorp/terraform-provider-google/pull/16182))
* **New Resource:** `google_alloydb_user` ([#16141](https://github.com/hashicorp/terraform-provider-google/pull/16141))
* **New Resource:** `google_firestore_backup_schedule` ([#16186](https://github.com/hashicorp/terraform-provider-google/pull/16186))
* **New Resource:** `google_redis_cluster` ([#16203](https://github.com/hashicorp/terraform-provider-google/pull/16203))

IMPROVEMENTS:
* alloydb: added `cluster_type` and `secondary_config` fields to support secondary clusters in `google_alloydb_cluster` resource. ([#16197](https://github.com/hashicorp/terraform-provider-google/pull/16197))
* compute: added `recreate_closed_psc` flag to support recreating the PSC Consumer forwarding rule if the `psc_connection_status` is closed on `google_compute_forwarding_rule`. ([#16188](https://github.com/hashicorp/terraform-provider-google/pull/16188))
* compute: added `INTERNET_IP_PORT`, `INTERNET_FQDN_PORT`, `SERVERLESS`, and `PRIVATE_SERVICE_CONNECT` as acceptable values for the `network_endpoint_type` field for the `resource_compute_network_endpoint_group` resource ([#16194](https://github.com/hashicorp/terraform-provider-google/pull/16194))
* compute: added `SEV_LIVE_MIGRATABLE_V2` to `guest_os_features` enum on `google_compute_image` resource. ([#16187](https://github.com/hashicorp/terraform-provider-google/pull/16187))
* compute: added `allow_subnet_cidr_routes_overlap` field to `google_compute_subnetwork` resource ([#16116](https://github.com/hashicorp/terraform-provider-google/pull/16116))
* compute: promoted `labels`, `effective_labels`, `terraform_labels`, and `label_fingerprint` fields in `google_compute_address` to GA ([#16120](https://github.com/hashicorp/terraform-provider-google/pull/16120))
* compute: promoted `internal_ip` and `external_ip` fields in resources `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` to GA ([#16140](https://github.com/hashicorp/terraform-provider-google/pull/16140))
* compute: promoted `internal_ip` and `external_ip` fields in resources `google_compute_per_instance_config` and `google_compute_region_per_instance_config` to GA ([#16140](https://github.com/hashicorp/terraform-provider-google/pull/16140))
* iamworkforcepool: promoted field `oidc.jwks_json` in resource `google_iam_workforce_pool` to GA ([#16199](https://github.com/hashicorp/terraform-provider-google/pull/16199))

BUG FIXES:
* alloydb: added `client_connection_config` field to `google_alloydb_instance` resource ([#16202](https://github.com/hashicorp/terraform-provider-google/pull/16202))
* bigquery: removed mutual exclusivity checks for `view`, `materialized_view`, and `schema` for the `google_bigquery_table` resource ([#16193](https://github.com/hashicorp/terraform-provider-google/pull/16193))
* compute: added `certificate_manager_certificates` field to `google_compute_target_https_proxy` resource ([#16179](https://github.com/hashicorp/terraform-provider-google/pull/16179))
* compute: fixed an issue where external `google_compute_global_address` can't be created when `network_tier` in `google_compute_project_default_network_tier` is set to `STANDARD`  ([#16144](https://github.com/hashicorp/terraform-provider-google/pull/16144))
* compute: fixed a false permadiff on `ip_address` when it is set to ipv6 on `google_compute_forwarding_rule` ([#16115](https://github.com/hashicorp/terraform-provider-google/pull/16115))
* provider: fixed a bug where an update request was sent to services when updateMask is empty ([#16111](https://github.com/hashicorp/terraform-provider-google/pull/16111))

## 5.1.0 (Oct 9, 2023)

FEATURES:
* **New Resource:** `google_database_migration_service_private_connection` ([#16104](https://github.com/hashicorp/terraform-provider-google/pull/16104))
* **New Resource:** `google_edgecontainer_cluster` ([#16055](https://github.com/hashicorp/terraform-provider-google/pull/16055))
* **New Resource:** `google_edgecontainer_node_pool` ([#16055](https://github.com/hashicorp/terraform-provider-google/pull/16055))
* **New Resource:** `google_edgecontainer_vpn_connection` ([#16055](https://github.com/hashicorp/terraform-provider-google/pull/16055))
* **New Resource:** `google_firebase_hosting_custom_domain` ([#16062](https://github.com/hashicorp/terraform-provider-google/pull/16062))
* **New Resource:** `google_gke_hub_fleet` ([#16072](https://github.com/hashicorp/terraform-provider-google/pull/16072))

IMPROVEMENTS:
* compute: added `device_name` field to `scratch_disk` block of `google_compute_instance` resource ([#16049](https://github.com/hashicorp/terraform-provider-google/pull/16049))
* container: added `node_config.linux_node_config.cgroup_mode` field to `google_container_node_pool` ([#16103](https://github.com/hashicorp/terraform-provider-google/pull/16103))
* databasemigrationservice: added support for `oracle` profiles to `google_database_migration_service_connection_profile` ([#16087](https://github.com/hashicorp/terraform-provider-google/pull/16087))
* firestore: added `api_scope` field to `google_firestore_index` resource ([#16085](https://github.com/hashicorp/terraform-provider-google/pull/16085))
* gkehub: added `location` field to `google_gke_hub_membership_iam_*` resources ([#16105](https://github.com/hashicorp/terraform-provider-google/pull/16105))
* gkehub: added `location` field to `google_gke_hub_membership` resource ([#16105](https://github.com/hashicorp/terraform-provider-google/pull/16105))
* gkeonprem: added update-in-place support for `vcenter` fields in `google_gkeonprem_vmware_cluster` ([#16073](https://github.com/hashicorp/terraform-provider-google/pull/16073))
* identityplatform: added `sms_region_config` to the resource `google_identity_platform_config` ([#16044](https://github.com/hashicorp/terraform-provider-google/pull/16044))

BUG FIXES:
* dns: fixed record set configuration parsing in `google_dns_record_set` ([#16042](https://github.com/hashicorp/terraform-provider-google/pull/16042))
* provider: fixed an issue where the plugin-framework implementation of the provider handled default region values that were self-links differently to the SDK implementation. This issue is not believed to have affected users because of downstream functions that turn self links into region names. ([#16100](https://github.com/hashicorp/terraform-provider-google/pull/16100))
* provider: fixed a bug that caused update requests to be sent for resources with a `terraform_labels` field even if no fields were updated ([#16111](https://github.com/hashicorp/terraform-provider-google/pull/16111))

## 5.0.0 (Oct 2, 2023)

KNOWN ISSUES:

* Updating some resources post-upgrade results in an error like "The update_mask in the Update{{Resource}}Request must be set". This should be resolved in `5.1.0`, see https://github.com/hashicorp/terraform-provider-google/issues/16091 for details.

[Terraform Google Provider 5.0.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_5_upgrade)

NOTES:
* provider: some provider default values are now shown at plan-time ([#15707](https://github.com/hashicorp/terraform-provider-google/pull/15707))

LABELS REWORK:
* provider: default labels configured on the provider through the new `default_labels` field are now supported. The default labels configured on the provider will be applied to all of the resources with standard `labels` field.
* provider: resources with labels - three label-related fields are now in all of the resources with standard `labels` field. `labels` field is non-authoritative and only manages the labels defined by the users on the resource through Terraform. The new output-only `terraform_labels` field merges the labels defined by the users on the resource through Terraform and the default labels configured on the provider. The new output-only `effective_labels` field lists all of labels present on the resource in GCP, including the labels configured through Terraform, the system, and other clients.
* provider: resources with annotations - two annotation-related fields are now in all of the resources with standard `annotations` field. The `annotations` field is non-authoritative and only manages the annotations defined by the users on the resource through Terraform. The new output-only `effective_annotations` field lists all of annotations present on the resource in GCP, including the annotations configured through Terraform, the system, and other clients.
* provider: datasources with labels - three fields `labels`, `terraform_labels`, and `effective_labels` are now present in most resource-based datasources. All three fields have all of labels present on the resource in GCP including the labels configured through Terraform, the system, and other clients, equivalent to `effective_labels` on the resource.
* provider: datasources with annotations - both `annotations` and `effective_annotations` are now present in most resource-based datasources. Both fields have all of annotations present on the resource in GCP including the annotations configured through Terraform, the system, and other clients, equivalent to `effective_annotations` on the resource.

BREAKING CHANGES:
* provider: added provider-level validation so these fields are not set as empty strings in a user's config: `credentials`, `access_token`, `impersonate_service_account`, `project`, `billing_project`, `region`, `zone` ([#15968](https://github.com/hashicorp/terraform-provider-google/pull/15968))
* provider: fixed many import functions throughout the provider that matched a subset of the provided input when possible. Now, the GCP resource id supplied to "terraform import" must match exactly. ([#15977](https://github.com/hashicorp/terraform-provider-google/pull/15977))
* provider: made data sources return errors on 404s when applicable instead of silently failing ([#15799](https://github.com/hashicorp/terraform-provider-google/pull/15799))
* provider: made empty strings in the provider configuration block no longer be ignored when configuring the provider([#15968](https://github.com/hashicorp/terraform-provider-google/pull/15968))
* accesscontextmanager: changed multiple array fields to sets where appropriate to prevent duplicates and fix diffs caused by server side reordering. ([#15756](https://github.com/hashicorp/terraform-provider-google/pull/15756))
* bigquery: added more input validations for `google_bigquery_table` schema ([#15338](https://github.com/hashicorp/terraform-provider-google/pull/15338))
* bigquery: made `routine_type` required for `google_bigquery_routine` ([#15517](https://github.com/hashicorp/terraform-provider-google/pull/15517))
* cloudfunction2: made `location` required on `google_cloudfunctions2_function` ([#15830](https://github.com/hashicorp/terraform-provider-google/pull/15830))
* cloudiot: removed deprecated datasource `google_cloudiot_registry_iam_policy` ([#15739](https://github.com/hashicorp/terraform-provider-google/pull/15739))
* cloudiot: removed deprecated resource `google_cloudiot_device` ([#15739](https://github.com/hashicorp/terraform-provider-google/pull/15739))
* cloudiot: removed deprecated resource  `google_cloudiot_registry` ([#15739](https://github.com/hashicorp/terraform-provider-google/pull/15739))
* cloudiot: removed deprecated resource `google_cloudiot_registry_iam_*` ([#15739](https://github.com/hashicorp/terraform-provider-google/pull/15739))
* cloudrunv2: removed deprecated field `liveness_probe.tcp_socket` from `google_cloud_run_v2_service` resource. ([#15430](https://github.com/hashicorp/terraform-provider-google/pull/15430))
* cloudrunv2: removed deprecated fields `startup_probe` and `liveness_probe` from `google_cloud_run_v2_job` resource. ([#15430](https://github.com/hashicorp/terraform-provider-google/pull/15430))
* cloudrunv2: retyped `volumes.cloud_sql_instance.instances` to SET from ARRAY for `google_cloud_run_v2_service` ([#15831](https://github.com/hashicorp/terraform-provider-google/pull/15831))
* compute: made `google_compute_node_group` require one of `initial_size` or `autoscaling_policy` fields configured upon resource creation ([#16006](https://github.com/hashicorp/terraform-provider-google/pull/16006))
* compute: made `size` in `google_compute_node_group` an output only field. ([#16006](https://github.com/hashicorp/terraform-provider-google/pull/16006))
* compute: removed default value for `rule.rate_limit_options.encorce_on_key` on resource `google_compute_security_policy` ([#15681](https://github.com/hashicorp/terraform-provider-google/pull/15681))
* compute: retyped `consumer_accept_lists` to a SET from an ARRAY type for `google_compute_service_attachment` ([#15985](https://github.com/hashicorp/terraform-provider-google/pull/15985))
* container: added `deletion_protection` to `google_container_cluster` which is enabled to `true` by default. When enabled, this field prevents Terraform from deleting the resource. ([#16013](https://github.com/hashicorp/terraform-provider-google/pull/16013))
* container: changed `management.auto_repair` and `management.auto_upgrade` defaults to true in `google_container_node_pool` ([#15931](https://github.com/hashicorp/terraform-provider-google/pull/15931))
* container: changed `networking_mode` default to `VPC_NATIVE` for newly created `google_container_cluster` resources ([#6402](https://github.com/hashicorp/terraform-provider-google-beta/pull/6402))
* container: removed `enable_binary_authorization` in `google_container_cluster` ([#15868](https://github.com/hashicorp/terraform-provider-google/pull/15868))
* container: removed default for `logging_variant` in `google_container_node_pool` ([#15931](https://github.com/hashicorp/terraform-provider-google/pull/15931))
* container: removed default value in `network_policy.provider` in `google_container_cluster` ([#15920](https://github.com/hashicorp/terraform-provider-google/pull/15920))
* container: removed the behaviour that `google_container_cluster` will delete the cluster if it's created in an error state. Instead, it will mark the cluster as tainted, allowing manual inspection and intervention. To proceed with deletion, run another `terraform apply`. ([#15887](https://github.com/hashicorp/terraform-provider-google/pull/15887))
* container: reworked the `taint` field in `google_container_cluster` and `google_container_node_pool` to only manage a subset of taint keys based on those already in state. Most existing resources are unaffected, unless they use `sandbox_config`- see upgrade guide for details. ([#15959](https://github.com/hashicorp/terraform-provider-google/pull/15959))
* dataplex: removed `data_profile_result` and `data_quality_result` from `google_dataplex_scan` ([#15505](https://github.com/hashicorp/terraform-provider-google/pull/15505))
* firebase: changed `deletion_policy` default to `DELETE` for `google_firebase_web_app`. ([#15406](https://github.com/hashicorp/terraform-provider-google/pull/15406))
* firebase: removed `google_firebase_project_location` ([#15764](https://github.com/hashicorp/terraform-provider-google/pull/15764))
* gameservices: removed Terraform support for `gameservices` ([#15558](https://github.com/hashicorp/terraform-provider-google/pull/15558))
* logging: changed the default value of `unique_writer_identity` from `false` to `true` in `google_logging_project_sink`. ([#15743](https://github.com/hashicorp/terraform-provider-google/pull/15743))
* logging: made `growth_factor`, `num_finite_buckets`, and `scale` required for `google_logging_metric` ([#15680](https://github.com/hashicorp/terraform-provider-google/pull/15680))
* looker: removed `LOOKER_MODELER` as a possible value in `google_looker_instance.platform_edition` ([#15956](https://github.com/hashicorp/terraform-provider-google/pull/15956))
* monitoring: fixed perma-diffs in `google_monitoring_dashboard.dashboard_json` by suppressing values returned by the API that are not in configuration ([#16014](https://github.com/hashicorp/terraform-provider-google/pull/16014))
* monitoring: made `labels` immutable in `google_monitoring_metric_descriptor` ([#15988](https://github.com/hashicorp/terraform-provider-google/pull/15988))
* privateca: removed deprecated fields `config_values`, `pem_certificates` from `google_privateca_certificate` ([#15537](https://github.com/hashicorp/terraform-provider-google/pull/15537))
* secretmanager: removed `automatic` field in `google_secret_manager_secret` resource ([#15859](https://github.com/hashicorp/terraform-provider-google/pull/15859))
* servicenetworking: used Create instead of Patch to create `google_service_networking_connection` ([#15761](https://github.com/hashicorp/terraform-provider-google/pull/15761))
* servicenetworking: used the `deleteConnection` method to delete the resource `google_service_networking_connection` ([#15934](https://github.com/hashicorp/terraform-provider-google/pull/15934))

FEATURES:
* **New Resource:** `google_scc_folder_custom_module` ([#15979](https://github.com/hashicorp/terraform-provider-google/pull/15979))
* **New Resource:** `google_scc_organization_custom_module` ([#16012](https://github.com/hashicorp/terraform-provider-google/pull/16012))

IMPROVEMENTS:
* alloydb: added additional fields to `google_alloydb_instance` and `google_alloydb_backup` ([#15973](https://github.com/hashicorp/terraform-provider-google/pull/15974))
* artifactregistry: added support for remote APT and YUM repositories to `google_artifact_registry_repository` ([#15973](https://github.com/hashicorp/terraform-provider-google/pull/15973))
* baremetal: made delete a noop for the resource `google_bare_metal_admin_cluster` to better align with actual behavior ([#16010](https://github.com/hashicorp/terraform-provider-google/pull/16010))
* bigtable: added `state` output attribute to `google_bigtable_instance` clusters ([#15961](https://github.com/hashicorp/terraform-provider-google/pull/15961))
* compute: made `google_compute_node_group` mutable ([#16006](https://github.com/hashicorp/terraform-provider-google/pull/16006))
* container: added the `effective_taints` attribute to `google_container_cluster` and `google_container_node_pool`, outputting all known taint values ([#15959](https://github.com/hashicorp/terraform-provider-google/pull/15959))
* container: allowed setting `addons_config.gcs_fuse_csi_driver_config` on `google_container_cluster` with `enable_autopilot: true`. ([#15996](https://github.com/hashicorp/terraform-provider-google/pull/15996))
* containeraws: added `binary_authorization` to `google_container_aws_cluster` ([#15989](https://github.com/hashicorp/terraform-provider-google/pull/15989))
* containeraws: added `update_settings` to `google_container_aws_node_pool` ([#15989](https://github.com/hashicorp/terraform-provider-google/pull/15989))
* google_compute_instance ([#15933](https://github.com/hashicorp/terraform-provider-google/pull/15933))
* osconfig: added `week_day_of_month.day_offset` field to the `google_os_config_patch_deployment` resource ([#15997](https://github.com/hashicorp/terraform-provider-google/pull/15997))
* secretmanager: allowed update for `rotation.rotation_period` field in `google_secret_manager_secret` resource ([#15952](https://github.com/hashicorp/terraform-provider-google/pull/15952))
* sql: added `preferred_zone` field to `google_sql_database_instance` resource ([#15971](https://github.com/hashicorp/terraform-provider-google/pull/15971))
* storagetransfer: added `event_stream` field to `google_storage_transfer_job` resource ([#16004](https://github.com/hashicorp/terraform-provider-google/pull/16004))

BUG FIXES:
* bigquery: fixed diff suppression in `external_data_configuration.connection_id` in `google_bigquery_table` ([#15983](https://github.com/hashicorp/terraform-provider-google/pull/15983))
* bigquery: fixed view and materialized view creation when schema is specified in `google_bigquery_table` ([#15442](https://github.com/hashicorp/terraform-provider-google/pull/15442))
* bigtable: avoided re-creation of `google_bigtable_instance` when cluster is still updating and storage type changed ([#15961](https://github.com/hashicorp/terraform-provider-google/pull/15961))
* bigtable: fixed a bug where dynamically created clusters would incorrectly run into duplication error in `google_bigtable_instance` ([#15940](https://github.com/hashicorp/terraform-provider-google/pull/15940))
* compute: removed the default value for field `reconcile_connections ` in resource `google_compute_service_attachment`, the field will now default to a value returned by the API when not set in configuration ([#15919](https://github.com/hashicorp/terraform-provider-google/pull/15919))
* compute: replaced incorrect default value for `enable_endpoint_independent_mapping` with APIs default in resource `google_compute_router_nat` ([#15478](https://github.com/hashicorp/terraform-provider-google/pull/15478))
* container: fixed an issue in `google_container_node_pool` where empty `linux_node_config.sysctls` would crash the provider ([#15941](https://github.com/hashicorp/terraform-provider-google/pull/15941))
* dataflow: fixed issue causing error message when max_workers and num_workers were supplied via parameters in `google_dataflow_flex_template_job` ([#15976](https://github.com/hashicorp/terraform-provider-google/pull/15976))
* dataflow: fixed max_workers read value permanently displaying as 0 in `google_dataflow_flex_template_job` ([#15976](https://github.com/hashicorp/terraform-provider-google/pull/15976))
* dataflow: fixed permadiff when SdkPipeline values are supplied via parameters in `google_dataflow_flex_template_job` ([#15976](https://github.com/hashicorp/terraform-provider-google/pull/15976))
* identityplayform: fixed a potential perma-diff for `sign_in` in `google_identity_platform_config` resource ([#15907](https://github.com/hashicorp/terraform-provider-google/pull/15907))
* firebase: made `google_firebase_rules.release` immutable ([#15989](https://github.com/hashicorp/terraform-provider-google/pull/15989))
* monitoring: fixed an issue where `metadata` was not able to be updated in `google_monitoring_metric_descriptor` ([#16014](https://github.com/hashicorp/terraform-provider-google/pull/16014))
* monitoring: fixed bug where importing `google_monitoring_notification_channel` failed when no default project was supplied in provider configuration or through environment variables ([#15929](https://github.com/hashicorp/terraform-provider-google/pull/15929))
* secretmanager: fixed an issue in `google_secretmanager_secret` where replacing `replication.automatic` with `replication.auto` would destroy and recreate the resource ([#15922](https://github.com/hashicorp/terraform-provider-google/pull/15922))
* sql: fixed diffs when re-ordering existing `database_flags` in `google_sql_database_instance` ([#15678](https://github.com/hashicorp/terraform-provider-google/pull/15678))
* tags: fixed import failure on `google_tags_tag_binding` ([#16005](https://github.com/hashicorp/terraform-provider-google/pull/16005))
* vertexai: made `contents_delta_uri` a required field in `google_vertex_ai_index` as omitting it would result in an error ([#15992](https://github.com/hashicorp/terraform-provider-google/pull/15992))
