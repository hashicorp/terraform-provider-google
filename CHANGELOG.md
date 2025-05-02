## 6.33.0 (Unreleased)

DEPRECATIONS:
* tpu: deprecated `google_tpu_node` resource. `google_tpu_node` is deprecated and will be removed in a future major release. Use `google_tpu_v2_vm` instead. ([#22552](https://github.com/hashicorp/terraform-provider-google/pull/22552))

FEATURES:
* **New Resource:** `google_apigee_security_profile_v2` ([#22524](https://github.com/hashicorp/terraform-provider-google/pull/22524))
* **New Resource:** `google_resource_manager_capability` (beta) ([#22582](https://github.com/hashicorp/terraform-provider-google/pull/22582))

IMPROVEMENTS:
* bigtable: added `cluster.node_scaling_factor` field to `google_bigtable_instance` resource ([#22560](https://github.com/hashicorp/terraform-provider-google/pull/22560))
* cloudrunv2: added `scaling_mode` and `manual_instance_count` fields to `google_cloud_run_v2_service` resource ([#22561](https://github.com/hashicorp/terraform-provider-google/pull/22561))
* container: added `flex_start` to `node_config` in `google_container_cluster` and `google_container_node_pool` (ga revert) ([#22542](https://github.com/hashicorp/terraform-provider-google/pull/22542))
* networkconnectivity: added `state_reason` field to `spoke` resource ([#22525](https://github.com/hashicorp/terraform-provider-google/pull/22525))
* sql: added `connection_pool_config` field. ([#22583](https://github.com/hashicorp/terraform-provider-google/pull/22583))
* vpcaccess: changed fields `min_instances`, `max_instances`, `machine_type` to allow update `google_vpc_access_connector` without without recreation. ([#22572](https://github.com/hashicorp/terraform-provider-google/pull/22572))

BUG FIXES:
* compute: fixed the bug when validating the subnetwork project in `google_compute_instance` resource ([#22571](https://github.com/hashicorp/terraform-provider-google/pull/22571))
* workbench: fixed a permadiff on `metadata` of `instance-region` in `google_workbench_instance` resource ([#22553](https://github.com/hashicorp/terraform-provider-google/pull/22553))

## 6.32.0 (Apr 25, 2025)

NOTES:
* `6.32.0` contains no changes from `6.31.1`. This release is being made to ensure that the version numbers of the `google` and `google-beta` provider releases remain aligned, as [`google-beta`'s `6.32.0` release](https://github.com/hashicorp/terraform-provider-google-beta/releases/tag/v6.32.0) contains a beta-only change.


## 6.31.1 (Apr 25, 2025)

BUG FIXES:
* storage: removed extra permission (storage.anywhereCaches.list) required for destroying a `resource_storage_bucket` ([#22442](https://github.com/hashicorp/terraform-provider-google/pull/22442))

## 6.31.0 (Apr 22, 2025)

DEPRECATIONS:
* integrations: deprecated `run_as_service_account` field in `google_integrations_client` resource ([#22312](https://github.com/hashicorp/terraform-provider-google/pull/22312))

FEATURES:
* **New Resource:** `google_compute_resource_policy_attachment` ([#22400](https://github.com/hashicorp/terraform-provider-google/pull/22400))
* **New Resource:** `google_compute_storage_pool` ([#22343](https://github.com/hashicorp/terraform-provider-google/pull/22343))
* **New Resource:** `google_gke_backup_backup_channel` ([#22393](https://github.com/hashicorp/terraform-provider-google/pull/22393))
* **New Resource:** `google_gke_backup_restore_channel` ([#22393](https://github.com/hashicorp/terraform-provider-google/pull/22393))
* **New Resource:** `google_iap_web_cloud_run_service_iam_binding` ([#22399](https://github.com/hashicorp/terraform-provider-google/pull/22399))
* **New Resource:** `google_iap_web_cloud_run_service_iam_member` ([#22399](https://github.com/hashicorp/terraform-provider-google/pull/22399))
* **New Resource:** `google_iap_web_cloud_run_service_iam_policy` ([#22399](https://github.com/hashicorp/terraform-provider-google/pull/22399))
* **New Resource:** `google_storage_batch_operations_job` ([#22333](https://github.com/hashicorp/terraform-provider-google/pull/22333))

IMPROVEMENTS:
* accesscontextmanager: added `scoped_access_settings` field to `gcp_user_access_binding` resource ([#22308](https://github.com/hashicorp/terraform-provider-google/pull/22308))
* alloydb: added `machine_type` field to `google_alloydb_instance` resource ([#22352](https://github.com/hashicorp/terraform-provider-google/pull/22352))
* artifactregistry: added `DEBIAN_SNAPSHOT` enum value to `repository_base` in `google_artifact_registry_repository` ([#22315](https://github.com/hashicorp/terraform-provider-google/pull/22315))
* bigquery: added `external_catalog_dataset_options` fields to `google_bigquery_dataset` resource ([#22377](https://github.com/hashicorp/terraform-provider-google/pull/22377))
* compute: added `log_config.optional_mode`, `log_config.optional_fields`, `backend.preference`, `max_stream_duration` and `cdn_policy.request_coalescing` fields to `google_compute_backend_service` resource ([#22391](https://github.com/hashicorp/terraform-provider-google/pull/22391))
* container: added support for updating the `confidential_nodes` field in `google_container_node_pool` ([#22363](https://github.com/hashicorp/terraform-provider-google/pull/22363))
* discoveryengine: added `allow_cross_region` field to `google_discovery_engine_chat_engine` resource ([#22336](https://github.com/hashicorp/terraform-provider-google/pull/22336))
* gkehub: added `configmanagement.config_sync.deployment_overrides` field to `google_gke_hub_feature_membership` resource ([#22403](https://github.com/hashicorp/terraform-provider-google/pull/22403))
* kms: added new enum values for `import_method` field in  `google_kms_key_ring_import_job` resource ([#22314](https://github.com/hashicorp/terraform-provider-google/pull/22314))
* metastore: added `tags` field to `google_dataproc_metastore_service` resource to allow setting tags for services at creation time ([#22313](https://github.com/hashicorp/terraform-provider-google/pull/22313))
* monitoring: added `log_check_failures` to `google_monitoring_uptime_check_config` ([#22351](https://github.com/hashicorp/terraform-provider-google/pull/22351))
* networkconnectivity: added IPv6 support to `google_network_connectivity_internal_range` resource ([#22401](https://github.com/hashicorp/terraform-provider-google/pull/22401))
* networkconnectivity: added `exclude_cidr_ranges` field to `google_network_connectivity_internal_range` resource ([#22332](https://github.com/hashicorp/terraform-provider-google/pull/22332))
* privateca: added `backdate_duration` field to the `google_privateca_ca_pool` resource to add support for backdating the `not_before_time` of certificates ([#22380](https://github.com/hashicorp/terraform-provider-google/pull/22380))
* redis: added `tags` field to `google_redis_instance` ([#22337](https://github.com/hashicorp/terraform-provider-google/pull/22337))
* sql: added `custom_subject_alternative_names` field to `instances` resource ([#22357](https://github.com/hashicorp/terraform-provider-google/pull/22357))
* sql: added `data_disk_provisioned_iops` and `data_disk_provisioned_throughput` fields to `google_sql_database_instance` resource ([#22398](https://github.com/hashicorp/terraform-provider-google/pull/22398))
* sql: added `retain_backups_on_delete` field to `google_sql_database_instance` resource ([#22334](https://github.com/hashicorp/terraform-provider-google/pull/22334))

BUG FIXES:
* colab: fixed perma-diff in `google_colab_runtime_template` caused by not returning default values. ([#22338](https://github.com/hashicorp/terraform-provider-google/pull/22338))
* discoveryengine: fixed `google_discovery_engine_target_site` operations to allow for enough time to index before timing out ([#22358](https://github.com/hashicorp/terraform-provider-google/pull/22358))
* compute: fixed perma-diff in `google_compute_network_firewall_policy_rule` when `security_profile_group` starts with `//` ([#22402](https://github.com/hashicorp/terraform-provider-google/pull/22402))
* healthcare: made `google_healthcare_pipeline_job` wait for creation and update operation to complete ([#22339](https://github.com/hashicorp/terraform-provider-google/pull/22339))
* identityplatform: fixed perma-diff in `google_identity_platform_config` when fields in `blocking_functions.forward_inbound_credentials` are set to `false` ([#22384](https://github.com/hashicorp/terraform-provider-google/pull/22384))
* sql: added diff suppression for some version changes to`google_sql_database_instance`. Diffs for `database_version` for MySQL 8.0 will be suppressed when the version is updated by auto version upgrade.([#22356](https://github.com/hashicorp/terraform-provider-google/pull/22356))
* sql: fixed the issue of shortened version of failover_dr_replica_name causes unnecessary diff in `google_sql_database_instance` ([#22319](https://github.com/hashicorp/terraform-provider-google/pull/22319))

## 6.30.0 (Apr 15, 2025)

FEATURES:
* **New Resource:** `google_developer_connect_account_connector` ([#22270](https://github.com/hashicorp/terraform-provider-google/pull/22270))
* **New Resource:** `google_vertex_ai_feature_group_iam_*` ([#22260](https://github.com/hashicorp/terraform-provider-google/pull/22260))
* **New Resource:** `google_vertex_ai_feature_online_store_iam_*` ([#22260](https://github.com/hashicorp/terraform-provider-google/pull/22260))
* **New Resource:** `google_vertex_ai_feature_online_store_featureview_iam_*` ([#22260](https://github.com/hashicorp/terraform-provider-google/pull/22260))

IMPROVEMENTS:
* bigquery: added `external_catalog_table_options` and `schema_foreign_type_info` fields to  `google_bigquery_table` resource ([#22302](https://github.com/hashicorp/terraform-provider-google/pull/22302))
* cloudrunv2: added `iap_enabled` field to `google_cloud_run_v2_service` resource ([#22301](https://github.com/hashicorp/terraform-provider-google/pull/22301))
* compute: added `source_disk_encryption_key.kms_key_self_link` and `source_disk_encryption_key.rsa_encrypted_key` fields to `google_compute_snapshot` resource ([#22247](https://github.com/hashicorp/terraform-provider-google/pull/22247))
* compute: added `source_disk_encryption_key`, `source_image_encryption_key` and `source_snapshot_encryption_key` fields to `google_compute_image` resource ([#22247](https://github.com/hashicorp/terraform-provider-google/pull/22247))
* compute: added `type`, `source_nat_active_ranges` and `source_nat_drain_ranges` fields to `google_compute_router_nat` resource ([#22282](https://github.com/hashicorp/terraform-provider-google/pull/22282))
* databasemigrationservice: allowed setting `ssl.type` in `google_database_migration_service_connection_profile` resource ([#22268](https://github.com/hashicorp/terraform-provider-google/pull/22268))
* firestore: added `MONGODB_COMPATIBLE_API` enum option to `api_scope` field in `google_firestore_index` resource ([#22287](https://github.com/hashicorp/terraform-provider-google/pull/22287))
* firestore: added `database_edition` field to `google_firestore_database` resource ([#22287](https://github.com/hashicorp/terraform-provider-google/pull/22287))
* firestore: added `density` and `multikey` fields to `google_firestore_index` resource ([#22287](https://github.com/hashicorp/terraform-provider-google/pull/22287))
* memorystore: added `managed_backup_source` and `gcs_source` fields to `google_memorystore_instance` resource ([#22295](https://github.com/hashicorp/terraform-provider-google/pull/22295))
* monitoring: added `password_wo` write-only field and `password_wo_version` field to `google_monitoring_uptime_check_config` resource ([#22242](https://github.com/hashicorp/terraform-provider-google/pull/22242))
* redis: added `managed_backup_source` and `gcs_source` fields to `google_redis_cluster` resource ([#22277](https://github.com/hashicorp/terraform-provider-google/pull/22277))
* storage: added support for deleting pending caches present on bucket when setting `force_destory` to true in `google_storage_bucket` resource ([#22262](https://github.com/hashicorp/terraform-provider-google/pull/22262))
* storagecontrol: added `trial_config` field to `google_storage_control_folder_intelligence_config` resource ([#22236](https://github.com/hashicorp/terraform-provider-google/pull/22236))
* storagecontrol: added `trial_config` field to `google_storage_control_organization_intelligence_config` resource ([#22236](https://github.com/hashicorp/terraform-provider-google/pull/22236))
* storagecontrol: added `trial_config` field to `google_storage_control_project_intelligence_config` resource ([#22236](https://github.com/hashicorp/terraform-provider-google/pull/22236))

BUG FIXES:
* container: fixed perma-diff in `fleet` field when the `fleet.project` field being added is null or empty in `google_container_cluster` resource ([#22240](https://github.com/hashicorp/terraform-provider-google/pull/22240))
* pubsub: fixed perma-diff by changing `allowed_persistence_regions` field to set in `google_pubsub_topic` resource ([#22273](https://github.com/hashicorp/terraform-provider-google/pull/22273))

## 6.29.0 (Apr 8, 2025)

FEATURES:
* **New Resource:** `google_apigee_control_plane_access` ([#22209](https://github.com/hashicorp/terraform-provider-google/pull/22209))
* **New Resource:** `google_clouddeploy_deploy_policy` ([#22190](https://github.com/hashicorp/terraform-provider-google/pull/22190))
* **New Resource:** `google_gemini_code_tools_setting_binding` ([#22226](https://github.com/hashicorp/terraform-provider-google/pull/22226))
* **New Resource:** `google_gemini_code_tools_setting` ([#22203](https://github.com/hashicorp/terraform-provider-google/pull/22203))
* **New Resource:** `google_os_config_v2_policy_orchestrator_for_organization` ([#22192](https://github.com/hashicorp/terraform-provider-google/pull/22192))

IMPROVEMENTS:
* accesscontextmanager: added `session_settings` field to `gcp_user_access_binding` resource ([#22227](https://github.com/hashicorp/terraform-provider-google/pull/22227))
* cloudedeploy: added `timed_promote_release_rule` and `repair_rollout_rule` fields to `google_clouddeploy_automation` resource ([#22190](https://github.com/hashicorp/terraform-provider-google/pull/22190))
* compute: added `group_placement_policy.0.tpu_topology` field to `google_compute_resource_policy` resource ([#22201](https://github.com/hashicorp/terraform-provider-google/pull/22201))
* datastream: added support for creating streams for Salesforce source in `google_datastream_stream` ([#22205](https://github.com/hashicorp/terraform-provider-google/pull/22205))
* gkehub: enabled partial results to be returned when a cloud region is unreachable in `google_gke_hub_feature ` ([#22218](https://github.com/hashicorp/terraform-provider-google/pull/22218))
* gkeonprem: added `enable_advanced_cluster` field to `google_gkeonprem_vmware_admin_cluster` resource ([#22188](https://github.com/hashicorp/terraform-provider-google/pull/22188))
* gkeonprem: added `enable_advanced_cluster` field to `google_gkeonprem_vmware_cluster` resource ([#22188](https://github.com/hashicorp/terraform-provider-google/pull/22188))
* memorystore: added `automated_backup_config` field to `google_memorystore_instance` resource, ([#22208](https://github.com/hashicorp/terraform-provider-google/pull/22208))
* netapp: added `tiering_policy` to `google_netapp_volume_replication` resource ([#22223](https://github.com/hashicorp/terraform-provider-google/pull/22223))
* parametermanagerregional: added `kms_key_version` field to `google_parameter_manager_regional_parameter_version` resource and datasource ([#22213](https://github.com/hashicorp/terraform-provider-google/pull/22213))
* parametermanagerregional: added `kms_key` field to `google_parameter_manager_regional_parameter` resource and `google_parameter_manager_regional_parameters` datasource ([#22213](https://github.com/hashicorp/terraform-provider-google/pull/22213))
* redis: added `automated_backup_config` field to `google_redis_cluster` ([#22117](https://github.com/hashicorp/terraform-provider-google/pull/22117))
* storage: added `md5hexhash` field in `google_storage_bucket_object` ([#22229](https://github.com/hashicorp/terraform-provider-google/pull/22229))
* workbench: added `confidential_instance_config` field to `google_workbench_instance` resource ([#22178](https://github.com/hashicorp/terraform-provider-google/pull/22178))

BUG FIXES:
* colab: fixed an issue where `google_colab_*` resources incorrectly required a provider-level region matching the resource location ([#22217](https://github.com/hashicorp/terraform-provider-google/pull/22217))
* datastream: updated `private_key`to be mutable in `google_datastream_connection_profile` resource. ([#22179](https://github.com/hashicorp/terraform-provider-google/pull/22179))
  
## 6.28.0 (Apr 1, 2025)

DEPRECATIONS:
* compute: deprecated `enable_flow_logs` in favor of `log_config` on `google_compute_subnetwork` resource.  If `log_config` is present, flow logs are enabled, and `enable_flow_logs` can be safely removed. ([#22111](https://github.com/hashicorp/terraform-provider-google/pull/22111))
* containerregistry: Deprecated `google_container_registry` resource, and `google_container_registry_image` and `google_container_registry_repository` data sources. Use `google_artifact_registry_repository` instead. ([#22071](https://github.com/hashicorp/terraform-provider-google/pull/22071))

FEATURES:
* **New Data Source:** `google_compute_region_backend_service` ([#21986](https://github.com/hashicorp/terraform-provider-google/pull/21986))
* **New Data Source:** `google_organization_iam_custom_roles` ([#22035](https://github.com/hashicorp/terraform-provider-google/pull/22035))
* **New Data Source:** `google_parameter_manager_parameter_version_render` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_parameter_version` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_parameter` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_parameters` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_regional_parameter_version_render` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_regional_parameter_version` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_regional_parameter` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_parameter_manager_regional_parameters` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Data Source:** `google_storage_control_folder_intelligence_config` ([#22077](https://github.com/hashicorp/terraform-provider-google/pull/22077))
* **New Data Source:** `google_storage_control_organization_intelligence_config` ([#22077](https://github.com/hashicorp/terraform-provider-google/pull/22077))
* **New Data Source:** `google_storage_control_project_intelligence_config` ([#22077](https://github.com/hashicorp/terraform-provider-google/pull/22077))
* **New Resource:** `google_apigee_dns_zone` ([#21992](https://github.com/hashicorp/terraform-provider-google/pull/21992))
* **New Resource:** `google_chronicle_data_access_scope` ([#21982](https://github.com/hashicorp/terraform-provider-google/pull/21982))
* **New Resource:** `google_chronicle_referencelist` ([#22090](https://github.com/hashicorp/terraform-provider-google/pull/22090))
* **New Resource:** `google_chronicle_retrohunt` ([#22092](https://github.com/hashicorp/terraform-provider-google/pull/22092))
* **New Resource:** `google_chronicle_rule` ([#22089](https://github.com/hashicorp/terraform-provider-google/pull/22089))
* **New Resource:** `google_chronicle_rule_deployment` ([#22093](https://github.com/hashicorp/terraform-provider-google/pull/22093))
* **New Resource:** `google_chronicle_watchlist` ([#21989](https://github.com/hashicorp/terraform-provider-google/pull/21989))
* **New Resource:** `google_dataproc_metastore_database_iam_*` resources ([#21985](https://github.com/hashicorp/terraform-provider-google/pull/21985))
* **New Resource:** `google_dataproc_metastore_table_iam_*` ([#22064](https://github.com/hashicorp/terraform-provider-google/pull/22064))
* **New Resource:** `google_discovery_engine_sitemap` ([#21976](https://github.com/hashicorp/terraform-provider-google/pull/21976))
* **New Resource:** `google_eventarc_enrollment` ([#22028](https://github.com/hashicorp/terraform-provider-google/pull/22028))
* **New Resource:** `google_firebase_app_hosting_build` ([#22063](https://github.com/hashicorp/terraform-provider-google/pull/22063))
* **New Resource:** `google_memorystore_instance_desired_user_created_endpoints` ([#22073](https://github.com/hashicorp/terraform-provider-google/pull/22073))
* **New Resource:** `google_parameter_manager_parameter_version` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Resource:** `google_parameter_manager_parameter` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Resource:** `google_parameter_manager_regional_parameter_version` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Resource:** `google_parameter_manager_regional_parameter` ([#22099](https://github.com/hashicorp/terraform-provider-google/pull/22099))
* **New Resource:** `google_storage_control_folder_intelligence_config` ([#22061](https://github.com/hashicorp/terraform-provider-google/pull/22061))
* **New Resource:** `google_storage_control_organization_intelligence_config` ([#21987](https://github.com/hashicorp/terraform-provider-google/pull/21987))

IMPROVEMENTS:
* accesscontextmanager: added `roles` field to ingress and egress policies of `google_access_context_manager_service_perimeter*` resources ([#22086](https://github.com/hashicorp/terraform-provider-google/pull/22086))
* cloudfunctions2: added `binary_authorization_policy` field to `google_cloudfunctions2_function` resource ([#22070](https://github.com/hashicorp/terraform-provider-google/pull/22070))
* cloudrun: promoted `node_selector` field in `google_cloud_run_service` resource to GA ([#22054](https://github.com/hashicorp/terraform-provider-google/pull/22054))
* cloudrunv2: added `gpu_zonal_redundancy_disabled` field to `google_cloud_run_v2_service` resource ([#22054](https://github.com/hashicorp/terraform-provider-google/pull/22054))
* cloudrunv2: promoted `node_selector` field in  `google_cloud_run_v2_service` resource to GA ([#22054](https://github.com/hashicorp/terraform-provider-google/pull/22054))
* compute: added `md5_authentication_keys` field to `google_compute_router` resource ([#22101](https://github.com/hashicorp/terraform-provider-google/pull/22101))
* compute: added `EXTERNAL_IPV6_SUBNETWORK_CREATION` as a supported value for the `mode` field in `google_compute_public_delegated_prefix` resource ([#22037](https://github.com/hashicorp/terraform-provider-google/pull/22037))
* compute: added `external_ipv6_prefix`, `stack_type`, and `ipv6_access_type` fields to `google_compute_subnetwork` data source ([#22085](https://github.com/hashicorp/terraform-provider-google/pull/22085))
* compute: added several `boot_disk`, `attached_disk`, and `instance_encryption_key` fields to `google_compute_instance` and `google_compute_instance_template` resources ([#22096](https://github.com/hashicorp/terraform-provider-google/pull/22096))
* compute: added `image_encryption_key.raw_key` and `image_encryption_key.rsa_encrypted_key` fields to `google_compute_image` resource ([#22096](https://github.com/hashicorp/terraform-provider-google/pull/22096))
* compute: added `snapshot_encryption_key.rsa_encrypted_key` field to `google_compute_snapshot` resource ([#22096](https://github.com/hashicorp/terraform-provider-google/pull/22096))
* container: added `auto_monitoring_config` field to `google_container_cluster` resource ([#21970](https://github.com/hashicorp/terraform-provider-google/pull/21970))
* container: added `disable_l4_lb_firewall_reconciliation` field to `google_container_cluster` resource ([#22065](https://github.com/hashicorp/terraform-provider-google/pull/22065))
* datafusion: added `tags` field to `google_data_fusion_instance` resource to allow setting tags for instances at creation time ([#21977](https://github.com/hashicorp/terraform-provider-google/pull/21977))
* datastream: added `blmt_config` field to `bigquery_destination_config` resource to enable support for BigLake Managed Tables streams ([#22109](https://github.com/hashicorp/terraform-provider-google/pull/22109))
* datastream: added `secret_manager_stored_password` field to `google_datastream_connection_profile` resource ([#22046](https://github.com/hashicorp/terraform-provider-google/pull/22046))
* identityplatform: added `disabled_user_signup` and `disabled_user_deletion` to `google_identity_platform_tenant` resource ([#21983](https://github.com/hashicorp/terraform-provider-google/pull/21983))
* memorystore: added `psc_attachment_details` field to `google_memorystore_instance` resource, to enable use of the fine-grained resource `google_memorystore_instance_desired_user_created_connections` ([#22073](https://github.com/hashicorp/terraform-provider-google/pull/22073))
* memorystore: added the `cross_cluster_replication_config` field to the `google_redis_cluster` resource ([#22097](https://github.com/hashicorp/terraform-provider-google/pull/22097))
* metastore: added `deletion_protection` field to `google_dataproc_metastore_federation` resource ([#22106](https://github.com/hashicorp/terraform-provider-google/pull/22106))
* networksecurity: added `antivirus_overrides` field to `google_network_security_security_profile` resource ([#22060](https://github.com/hashicorp/terraform-provider-google/pull/22060))
* networksecurity: added `connected_deployment_groups` and `associations` fields to `google_network_security_mirroring_endpoint_group` resource ([#21974](https://github.com/hashicorp/terraform-provider-google/pull/21974))
* networksecurity: added `locations` field to `google_network_security_mirroring_deployment_group` resource ([#21975](https://github.com/hashicorp/terraform-provider-google/pull/21975))
* networksecurity: added `locations` field to `google_network_security_mirroring_endpoint_group_association` resource ([#21971](https://github.com/hashicorp/terraform-provider-google/pull/21971))
* parametermanager: added `kms_key_version` field to `google_parameter_manager_parameter_version` resource and datasource ([#22058](https://github.com/hashicorp/terraform-provider-google/pull/22058))
* parametermanager: added `kms_key` field to `google_parameter_manager_parameter` resource and `google_parameter_manager_parameters` datasource ([#22058](https://github.com/hashicorp/terraform-provider-google/pull/22058))
* provider: added `external_credentials` block in `provider` ([#22081](https://github.com/hashicorp/terraform-provider-google/pull/22081))
* redis: added `automated_backup_config` field to `google_redis_cluster` resource ([#22117](https://github.com/hashicorp/terraform-provider-google/pull/22117))
* storage: added `content_base64` field in `google_storage_bucket_object_content` datasource ([#22051](https://github.com/hashicorp/terraform-provider-google/pull/22051))

BUG FIXES:
* alloydb: added a mutex to `google_alloydb_cluster` to prevent conflicts among multiple cluster operations ([#21972](https://github.com/hashicorp/terraform-provider-google/pull/21972))
* artifactregistry: fixed type assertion panic in `google_artifact_registry_repository` resource ([#22100](https://github.com/hashicorp/terraform-provider-google/pull/22100))
* bigtable: fixed `automated_backup_policy` field for `google_bigtable_table` resource ([#22034](https://github.com/hashicorp/terraform-provider-google/pull/22034))
* cloudrunv2: fixed the diffs for unchanged `template.template.containers.env` in `google_cloud_run_v2_job` resource ([#22115](https://github.com/hashicorp/terraform-provider-google/pull/22115))
* compute: fixed a regression in `google_compute_subnetwork` where setting `log_config` would not enable flow logs without `enable_flow_logs` also being set to true. To enable or disable flow logs, please use `log_config`. `enable_flow_logs` is now deprecated and will be removed in the next major release. ([#22111](https://github.com/hashicorp/terraform-provider-google/pull/22111))
* compute: fixed unable to update the `preview` field for `google_compute_security_policy_rule` resource ([#21984](https://github.com/hashicorp/terraform-provider-google/pull/21984))
* orgpolicy: fix permadiff in `google_org_policy_policy` when multiple rules are present ([#21981](https://github.com/hashicorp/terraform-provider-google/pull/21981))
* resourcemanager: increased page size for list services api to help any teams hitting `ListEnabledRequestsPerMinutePerProject` quota issues ([#22050](https://github.com/hashicorp/terraform-provider-google/pull/22050))
* spanner: fixed issue with applying changes in provider `default_labels` on `google_spanner_instance` resource ([#22036](https://github.com/hashicorp/terraform-provider-google/pull/22036))
* storage: fixed `google_storage_anywhere_cache` to cancel long-running operations after create and update requests timeout ([#22031](https://github.com/hashicorp/terraform-provider-google/pull/22031))
* workbench: fixed metadata permadiff in `google_workbench_instance` resource ([#22056](https://github.com/hashicorp/terraform-provider-google/pull/22056))

## 6.27.0 (Mar 25, 2025)

FEATURES:
* **New Data Source:** `google_compute_images` ([#21872](https://github.com/hashicorp/terraform-provider-google/pull/21872))
* **New Data Source:** `google_organization_iam_custom_role` ([#21922](https://github.com/hashicorp/terraform-provider-google/pull/21922))
* **New Resource:** `google_lustre_instance` ([#21963](https://github.com/hashicorp/terraform-provider-google/pull/21963))
* **New Resource:** `google_os_config_v2_policy_orchestrator` ([#21930](https://github.com/hashicorp/terraform-provider-google/pull/21930))
* **New Resource:** `google_storage_control_project_intelligence_config` ([#21902](https://github.com/hashicorp/terraform-provider-google/pull/21902))
* **New Resource:** `google_chronicle_data_access_label` ([#21956](https://github.com/hashicorp/terraform-provider-google/pull/21956))
* **New Resource:** `google_compute_router_route_policy` ([#21945](https://github.com/hashicorp/terraform-provider-google/pull/21945))

IMPROVEMENTS:
* bigquery: added `secondary_location` and `replication_status` fields to support managed disaster recovery feature in `google_bigquery_reservation` ([#21920](https://github.com/hashicorp/terraform-provider-google/pull/21920))
* clouddeploy: added `dns_endpoint` field to to `google_clouddeploy_target` resource ([#21868](https://github.com/hashicorp/terraform-provider-google/pull/21868))
* compute: added `shielded_instance_initial_state` structure to `google_compute_image` resource ([#21937](https://github.com/hashicorp/terraform-provider-google/pull/21937))
* compute: added `LINK_TYPE_ETHERNET_400G_LR4` enum value to `link_type` field in `google_compute_interconnect` resource ([#21903](https://github.com/hashicorp/terraform-provider-google/pull/21903))
* compute: added `architecture` and `guest_os_features` to `google_compute_instance` ([#21875](https://github.com/hashicorp/terraform-provider-google/pull/21875))
* compute: added `workload_policy.type`, `workload_policy.max_topology_distance` and `workload_policy.accelerator_topology` fields to `google_compute_resource_policy` resource ([#21961](https://github.com/hashicorp/terraform-provider-google/pull/21961))
* container: added `ip_endpoints_config` field to `google_container_cluster` resource ([#21959](https://github.com/hashicorp/terraform-provider-google/pull/21959))
* container: added `node_config.windows_node_config` field to `google_container_node_pool` resource. ([#21876](https://github.com/hashicorp/terraform-provider-google/pull/21876))
* container: added `pod_autoscaling` field to `google_container_cluster` resource ([#21919](https://github.com/hashicorp/terraform-provider-google/pull/21919))
* memorystore: added the `maintenance_policy` field to the `google_memorystore_instance` resource ([#21957](https://github.com/hashicorp/terraform-provider-google/pull/21957))
* memorystore: enabled update support for `node_type` field in `google_memorystore_instance` resource ([#21899](https://github.com/hashicorp/terraform-provider-google/pull/21899))
* metastore: promoted `scaling_config` field of `google_dataproc_metastore_service` resource to GA ([#21877](https://github.com/hashicorp/terraform-provider-google/pull/21877))
* networksecurity: added `connected_deployment_group` and `associations` fields to `google_network_security_intercept_endpoint_group` resource ([#21940](https://github.com/hashicorp/terraform-provider-google/pull/21940))
* networksecurity: added `locations` field to `google_network_security_intercept_deployment_group` resource ([#21923](https://github.com/hashicorp/terraform-provider-google/pull/21923))
* networksecurity: added `locations` field to `google_network_security_intercept_endpoint_group_association` resource ([#21962](https://github.com/hashicorp/terraform-provider-google/pull/21962))
* redis: added update support for `google_redis_cluster` `node_type` ([#21870](https://github.com/hashicorp/terraform-provider-google/pull/21870))
* storage: added metadata_options in `google_storage_transfer_job` ([#21897](https://github.com/hashicorp/terraform-provider-google/pull/21897))

BUG FIXES:
* bigqueryanalyticshub: fixed a bug in `google_bigquery_analytics_hub_listing_subscription` where a subscription using a different project than the dataset would not work ([#21958](https://github.com/hashicorp/terraform-provider-google/pull/21958))
* cloudrun: fixed the perma-diffs for unchanged `template.spec.containers.env` in `google_cloud_run_service` resource ([#21916](https://github.com/hashicorp/terraform-provider-google/pull/21916))
* cloudrunv2: fixed the perma-diffs for unchanged `template.containers.env` in `google_cloud_run_v2_service` resource ([#21916](https://github.com/hashicorp/terraform-provider-google/pull/21916))
* compute: fixed the issue that user can't use regional disk in `google_compute_instance_template` ([#21901](https://github.com/hashicorp/terraform-provider-google/pull/21901))
* dataflow: fixed a permadiff on `template_gcs_path` in `google_dataflow_job` resource ([#21894](https://github.com/hashicorp/terraform-provider-google/pull/21894))
* storage: lowered the minimum required items for `custom_placement_config.data_locations` from 2 to 1, and removed the Terraform-enforced maximum item limit for the field in `google_storage_bucket` ([#21878](https://github.com/hashicorp/terraform-provider-google/pull/21878))

## 6.26.0 (Mar 18, 2025)

FEATURES:
* **New Data Source:** `google_project_iam_custom_role` ([#21866](https://github.com/hashicorp/terraform-provider-google/pull/21866))
* **New Data Source:** `google_project_iam_custom_roles` ([#21813](https://github.com/hashicorp/terraform-provider-google/pull/21813))
* **New Resource:** `google_eventarc_pipeline` ([#21761](https://github.com/hashicorp/terraform-provider-google/pull/21761))
* **New Resource:** `google_firebase_app_hosting_backend` ([#21840](https://github.com/hashicorp/terraform-provider-google/pull/21840))
* **New Resource:** `google_network_security_mirroring_deployment` ([#21853](https://github.com/hashicorp/terraform-provider-google/pull/21853))
* **New Resource:** `google_network_security_mirroring_deployment_group` ([#21853](https://github.com/hashicorp/terraform-provider-google/pull/21853))
* **New Resource:** `google_network_security_mirroring_endpoint_group_association` ([#21853](https://github.com/hashicorp/terraform-provider-google/pull/21853))
* **New Resource:** `google_network_security_mirroring_endpoint_group` ([#21853](https://github.com/hashicorp/terraform-provider-google/pull/21853))

IMPROVEMENTS:
* alloydb: added `psc_config` field to ``google_alloydb_cluster` resource ([#21863](https://github.com/hashicorp/terraform-provider-google/pull/21863))
* bigquery: added `table_metadata_view` query param to `google_bigquery_table` ([#21838](https://github.com/hashicorp/terraform-provider-google/pull/21838))
* clouddeploy: added `dns_endpoint` field to to `google_clouddeploy_target` resource ([#21868](https://github.com/hashicorp/terraform-provider-google/pull/21868))
* compute: added `UNRESTRICTED` option to the `tls_early_data` field in the `google_compute_target_https_proxy` resource ([#21821](https://github.com/hashicorp/terraform-provider-google/pull/21821))
* compute: added `enable_flow_logs` and `state` fields to `google_compute_subnetwork` resource ([#21851](https://github.com/hashicorp/terraform-provider-google/pull/21851))
* compute: promoted fields `single_instance_assignment` and `filter` to GA for `google_compute_autoscaler` resource ([#21760](https://github.com/hashicorp/terraform-provider-google/pull/21760))
* container: added additional value `KCP_HPA` for `logging_config.enable_components` field in `google_container_cluster` resource ([#21836](https://github.com/hashicorp/terraform-provider-google/pull/21836))
* dataform: added `deletion_policy` field to `google_dataform_repository` resource. Default value is `DELETE`. Setting `deletion_policy` to `FORCE` will delete any child resources of this repository as well. ([#21864](https://github.com/hashicorp/terraform-provider-google/pull/21864))
* memorystore: added update support for `engine_version` field in `google_memorystore_instance` resource ([#21843](https://github.com/hashicorp/terraform-provider-google/pull/21843))
* metastore: added `create_time` and `update_time` fields to `google_dataproc_metastore_federation` resource ([#21824](https://github.com/hashicorp/terraform-provider-google/pull/21824))
* metastore: added `create_time` and `update_time` fields to `google_dataproc_metastore_service` resource ([#21817](https://github.com/hashicorp/terraform-provider-google/pull/21817))
* networksecurity: added `not_operations` field to `google_network_security_authz_policy` resource ([#21785](https://github.com/hashicorp/terraform-provider-google/pull/21785))
* networkservices: added `ip_version` and `envoy_headers` fields to `google_network_services_gateway` resource ([#21788](https://github.com/hashicorp/terraform-provider-google/pull/21788))
* sql: increased `settings.insights_config.query_string_length` and `settings.insights_config.query_string_length` limits for Enterprise Plus edition `sql_database_instance` resource. ([#21848](https://github.com/hashicorp/terraform-provider-google/pull/21848))
* storageinsights: added `parquet_options` field to `google_storage_insights_report_config` resource ([#21816](https://github.com/hashicorp/terraform-provider-google/pull/21816))
* workflows: added `execution_history_level` field to `google_workflows_workflow` resource ([#21782](https://github.com/hashicorp/terraform-provider-google/pull/21782))

BUG FIXES:
* accesscontextmanager: fixed panic on empty `access_policies` in `google_access_context_manager_access_policy` ([#21845](https://github.com/hashicorp/terraform-provider-google/pull/21845))
* compute: adjusted mapped image names that were preventing usage of `fedora-coreos` in `google_compute_image` resource ([#21787](https://github.com/hashicorp/terraform-provider-google/pull/21787))
* container: re-added `DNS_SCOPE_UNSPECIFIED` value to the `dns_config.cluster_dns_scope` field in `google_container_cluster` resource and suppressed diffs between `DNS_SCOPE_UNSPECIFIED` in config and empty/null in state ([#21861](https://github.com/hashicorp/terraform-provider-google/pull/21861))
* discoveryengine: changed field `dataStoreIds` to mutable in `google_discovery_engine_search_engine` ([#21759](https://github.com/hashicorp/terraform-provider-google/pull/21759))
* networksecurity: `min_tls_version` and `tls_feature_profile` fields updated to use the server assigned default and prevent a permadiff in `google_network_security_tls_inspection_policy` resource. ([#21788](https://github.com/hashicorp/terraform-provider-google/pull/21788))
* oslogin: added a wait after creating `google_os_login_ssh_public_key` to allow propagation ([#21860](https://github.com/hashicorp/terraform-provider-google/pull/21860))
* spanner: fixed issue with disabling autoscaling in `google_spanner_instance` ([#21852](https://github.com/hashicorp/terraform-provider-google/pull/21852))

## 6.25.0 (Mar 11, 2025)

NOTES:
* eventarc: `google_eventarc_channel` now uses MMv1 engine instead of DCL. ([#21728](https://github.com/hashicorp/terraform-provider-google/pull/21728))
* workbench: increased create timeout for `google_workbench_instance` to 40mins. ([#21700](https://github.com/hashicorp/terraform-provider-google/pull/21700))

FEATURES:
* **New Data Source:** `google_compute_region_ssl_policy` ([#21633](https://github.com/hashicorp/terraform-provider-google/pull/21633))
* **New Resource:** `google_eventarc_google_api_source` ([#21732](https://github.com/hashicorp/terraform-provider-google/pull/21732))
* **New Resource:** `google_iam_oauth_client_credential` ([#21731](https://github.com/hashicorp/terraform-provider-google/pull/21731))
* **New Resource:** `google_iam_oauth_client` ([#21660](https://github.com/hashicorp/terraform-provider-google/pull/21660))
* **New Resource:** `network_services_endpoint_policy` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))
* **New Resource:** `network_services_grpc_route` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))
* **New Resource:** `network_services_http_route` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))
* **New Resource:** `network_services_mesh` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))
* **New Resource:** `network_services_service_binding` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))
* **New Resource:** `network_services_tcp_route` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))
* **New Resource:** `network_services_tls_route` ([#21676](https://github.com/hashicorp/terraform-provider-google/pull/21676))

IMPROVEMENTS:
* alloydb: added `psc_instance_config.psc_interface_configs` field to `google_alloydb_instance` resource ([#21701](https://github.com/hashicorp/terraform-provider-google/pull/21701))
* compute: added `create_snapshot_before_destroy` to `google_compute_disk` and `google_compute_region_disk` to enable creating a snapshot before disk deletion ([#21636](https://github.com/hashicorp/terraform-provider-google/pull/21636))
* compute: added `ip_collection` and `ipv6_gce_endpoint` fields to `google_compute_subnetwork` resource ([#21730](https://github.com/hashicorp/terraform-provider-google/pull/21730))
* compute: added `log_config.optional_mode` and `log_config.optional_fields` fields to `google_compute_region_backend_service` resource ([#21722](https://github.com/hashicorp/terraform-provider-google/pull/21722))
* compute: added `rsa_encrypted_key` to `google_compute_region_disk` ([#21636](https://github.com/hashicorp/terraform-provider-google/pull/21636))
* compute: added `scheduling.termination_time` field to `google_compute_instance`, `google_compute_instance_from_machine_image`, `google_compute_instance_from_template`, `google_compute_instance_template`, and `google_compute_region_instance_template` resources ([#21717](https://github.com/hashicorp/terraform-provider-google/pull/21717))
* compute: added update support for 'purpose' field in `google_compute_subnetwork` resource ([#21729](https://github.com/hashicorp/terraform-provider-google/pull/21729))
* compute: added update support for `firewall_policy` in `google_compute_firewall_policy_association` resource. It is recommended to only perform this operation in combination with a protective lifecycle tag such as "create_before_destroy" or "prevent_destroy" on your previous `firewall_policy` resource in order to prevent situations where a target attachment has no associated policy. ([#21735](https://github.com/hashicorp/terraform-provider-google/pull/21735))
* container: added "JOBSET" as a supported value for `enable_components` in `google_container_cluster` resource ([#21657](https://github.com/hashicorp/terraform-provider-google/pull/21657))
* firebasedataconnect: added `deletion_policy` field to `google_firebase_data_connect_service` resource ([#21736](https://github.com/hashicorp/terraform-provider-google/pull/21736))
* networksecurity: added `description` field to `google_network_security_intercept_deployment`, `google_network_security_intercept_deployment_group`, `google_network_security_intercept_endpoint_group` resources ([#21711](https://github.com/hashicorp/terraform-provider-google/pull/21711))
* networksecurity: added `description` field to `google_network_security_mirroring_deployment`, `google_network_security_mirroring_deployment_group`, `google_network_security_mirroring_endpoint_group` resources ([#21714](https://github.com/hashicorp/terraform-provider-google/pull/21714))
* tpuv2: added `spot` field to `google_tpu_v2_vm` resource ([#21716](https://github.com/hashicorp/terraform-provider-google/pull/21716))
* workstations: added `tags` field to `google_workstations_workstation_cluster` resource ([#21635](https://github.com/hashicorp/terraform-provider-google/pull/21635))

BUG FIXES:
* backupdr: added missing `SUNDAY` option to `days_of_week` field in `google_backup_dr_backup_plan` resource ([#21640](https://github.com/hashicorp/terraform-provider-google/pull/21640))
* compute: fixed `network_interface.internal_ipv6_prefix_length` not being set or read in Terraform state in `google_compute_instance` resource ([#21638](https://github.com/hashicorp/terraform-provider-google/pull/21638))
* compute: fixed bug in `google_compute_router_nat` where `max_ports_per_vm` couldn't be unset once set. ([#21721](https://github.com/hashicorp/terraform-provider-google/pull/21721))
* container: fixed perma-diff in `google_container_cluster` when `cluster_dns_scope` is unspecified ([#21637](https://github.com/hashicorp/terraform-provider-google/pull/21637))
* networksecurity: added wait time on `google_network_security_gateway_security_policy_rule` resource when creating and deleting to prevent race conditions ([#21643](https://github.com/hashicorp/terraform-provider-google/pull/21643))

## 6.24.0 (Mar 3, 2025)

NOTES:
* gemini: removed unsupported value `GEMINI_CLOUD_ASSIST` for field `product` in `google_gemini_logging_setting_binding` resource ([#21630](https://github.com/hashicorp/terraform-provider-google/pull/21630))
* iam: added member value to the error message when member validation fails for google_project_iam_* ([#21586](https://github.com/hashicorp/terraform-provider-google/pull/21586))

DEPRECATIONS:
* datacatalog: deprecated `google_data_catalog_entry` and `google_data_catalog_tag` resources. For steps to transition your Data Catalog users, workloads, and content to Dataplex Catalog, see https://cloud.google.com/dataplex/docs/transition-to-dataplex-catalog. ([#21541](https://github.com/hashicorp/terraform-provider-google/pull/21541))
* notebooks: deprecated non-functional `google_notebooks_location` resource ([#21517](https://github.com/hashicorp/terraform-provider-google/pull/21517))

FEATURES:
* **New Data Source:** `google_memorystore_instance` ([#21579](https://github.com/hashicorp/terraform-provider-google/pull/21579))
* **New Resource:** `google_apihub_host_project_registration` ([#21607](https://github.com/hashicorp/terraform-provider-google/pull/21607))
* **New Resource:** `google_compute_instant_snapshot` ([#21598](https://github.com/hashicorp/terraform-provider-google/pull/21598))
* **New Resource:** `google_eventarc_message_bus` ([#21611](https://github.com/hashicorp/terraform-provider-google/pull/21611))
* **New Resource:** `google_gemini_data_sharing_with_google_setting_binding` (GA) ([#21629](https://github.com/hashicorp/terraform-provider-google/pull/21629))
* **New Resource:** `google_gemini_gcp_enablement_setting_binding` (GA) ([#21587](https://github.com/hashicorp/terraform-provider-google/pull/21587))
* **New Resource:** `google_gemini_gemini_gcp_enablement_setting_binding` ([#21540](https://github.com/hashicorp/terraform-provider-google/pull/21540))
* **New Resource:** `google_storage_anywhere_cache` ([#21537](https://github.com/hashicorp/terraform-provider-google/pull/21537))

IMPROVEMENTS:
* alloydb: added ability to upgrade major version in `google_alloydb_cluster` with `database_version` ([#21582](https://github.com/hashicorp/terraform-provider-google/pull/21582))
* compute: added `creation_timestamp`, `next_hop_peering`, ` warnings.code`, `warnings.message`, `warnings.data.key`, `warnings.data.value`, `next_hop_hub`, `route_type`, `as_paths.path_segment_type`, `as_paths.as_lists` and `route_status`  fields to `google_compute_route` resource ([#21534](https://github.com/hashicorp/terraform-provider-google/pull/21534))
* compute: added `max_stream_duration` field to `google_compute_url_map` resource ([#21535](https://github.com/hashicorp/terraform-provider-google/pull/21535))
* compute: added `network_interface.network_attachment` field to `google_compute_instance` resource (ga) ([#21606](https://github.com/hashicorp/terraform-provider-google/pull/21606))
* compute: added `network_interface.network_attachment` to `google_compute_instance` data source (ga) ([#21606](https://github.com/hashicorp/terraform-provider-google/pull/21606))
* compute: added fields `architecture`, `source_instant_snapshot`, `source_storage_object`, `resource_manager_tags`  to `google_compute_disk`. ([#21598](https://github.com/hashicorp/terraform-provider-google/pull/21598))
* container: added enum  value `UPGRADE_INFO_EVENT` for GKE notification filter in `google_container_cluster` resource ([#21609](https://github.com/hashicorp/terraform-provider-google/pull/21609))
* iam: added `AZURE_AD_GROUPS_ID` field to `google_iam_workforce_pool_provider.extra_attributes_oauth2_client.attributes_type` resource ([#21624](https://github.com/hashicorp/terraform-provider-google/pull/21624))
* networkconnectivity: added `policy_mode` field to `google_network_connectivity_hub` resource ([#21589](https://github.com/hashicorp/terraform-provider-google/pull/21589))
* networkservices: added `location` field to `google_network_services_grpc_route` resource ([#21621](https://github.com/hashicorp/terraform-provider-google/pull/21621))
* storagetransfer: added `logging_config` field to `google_storage_transfer_job` resource ([#21523](https://github.com/hashicorp/terraform-provider-google/pull/21523))

BUG FIXES:
* bigquery: updated the `max_staleness` field in `google_bigquery_table` to be a computed field ([#21596](https://github.com/hashicorp/terraform-provider-google/pull/21596))
* chronicle: fixed an error during resource creation with certain `run_frequency` configurations in `google_chronicle_rule_deployment` ([#21610](https://github.com/hashicorp/terraform-provider-google/pull/21610))
* discoveryengine: fixed bug preventing creation of `google_discovery_engine_target_site` resources ([#21628](https://github.com/hashicorp/terraform-provider-google/pull/21628))
* eventarc: fixed an issue where `google_eventarc_trigger` creation failed due to the region could not be parsed from the trigger's name ([#21528](https://github.com/hashicorp/terraform-provider-google/pull/21528))
* publicca: encode b64_mac_key in base64url, not in base64 ([#21612](https://github.com/hashicorp/terraform-provider-google/pull/21612))
* storage: fixed a 412 error returned on some `google_storage_bucket_iam_policy` deletions ([#21626](https://github.com/hashicorp/terraform-provider-google/pull/21626))

## 6.23.0 (Feb 26, 2025)

NOTES:
* The `google_sql_user` resource now supports `password_wo` [write-only arguments](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)
* The `google_bigquery_data_transfer_config` resource now supports `secret_access_key_wo` [write-only arguments](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)
* The `google_secret_version` resource now supports `secret_data_wo` [write-only arguments](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)

IMPROVEMENTS:
* sql: added `password_wo` and `password_wo_version` fields to `google_sql_user` resource ([#21616](https://github.com/hashicorp/terraform-provider-google/pull/21616))
* bigquerydatatransfer: added `secret_access_key_wo` and `secret_access_key_wo_version` fields to `google_bigquery_data_transfer_config` resource ([#21617](https://github.com/hashicorp/terraform-provider-google/pull/21617))
* secretmanager: added `secret_data_wo` and `secret_data_wo_version` fields to `google_secret_version` resource ([#21618](https://github.com/hashicorp/terraform-provider-google/pull/21618))

## 6.22.0 (Feb 24, 2025)

NOTES:
* provider: The Terraform Provider for Google Cloud's regular release date will move from Monday to Tuesday in early March. The 2025/03/10 release will be made on 2025/03/11.

DEPRECATIONS:
* datacatalog: deprecated `google_data_catalog_tag_template`. Use `google_dataplex_aspect_type` instead. For steps to transition your Data Catalog users, workloads, and content to Dataplex Catalog, see https://cloud.google.com/dataplex/docs/transition-to-dataplex-catalog. ([#9347](https://github.com/hashicorp/terraform-provider-google-beta/pull/9347))
* datacatalog: deprecated `google_data_catalog_entry_group`. Use `google_dataplex_entry_group` instead. For steps to transition your Data Catalog users, workloads, and content to Dataplex Catalog, see https://cloud.google.com/dataplex/docs/transition-to-dataplex-catalog. ([#9349](https://github.com/hashicorp/terraform-provider-google-beta/pull/9349))

FEATURES:
* **New Data Source:** `google_alloydb_cluster` ([#21496](https://github.com/hashicorp/terraform-provider-google/pull/21496))
* **New Data Source:** `google_project_ancestry` ([#21413](https://github.com/hashicorp/terraform-provider-google/pull/21413))
* **New Resource:** `google_gemini_data_sharing_with_google_setting_binding` ([#21479](https://github.com/hashicorp/terraform-provider-google/pull/21479))
* **New Resource:** `google_gemini_logging_setting_binding` ([#21429](https://github.com/hashicorp/terraform-provider-google/pull/21429))
* **New Resource:** `google_gemini_logging_setting` ([#21404](https://github.com/hashicorp/terraform-provider-google/pull/21404))
* **New Resource:** `google_spanner_instance_partition` ([#21475](https://github.com/hashicorp/terraform-provider-google/pull/21475))

IMPROVEMENTS:
* backupdr: promoted `google_backup_dr_management_server`, `google_backup_dr_backup_plan_association`, and `google_backup_dr_backup_plan` resources to GA
* compute: added `import_subnet_routes_with_public_ip` and `export_subnet_routes_with_public_ip` fields to `google_compute_network_peering_routes_config` resource ([#21405](https://github.com/hashicorp/terraform-provider-google/pull/21405))
* developerconnect: added `bitbucket_cloud_config` and `bitbucket_data_center_config` fields to `google_developer_connect_connection` resource ([#21433](https://github.com/hashicorp/terraform-provider-google/pull/21433))
* gemini: promoted `google_gemini_release_channel_setting` resource to GA ([#21481](https://github.com/hashicorp/terraform-provider-google/pull/21481))
* iam: added `extra_attributes_oauth2_client` field to `google_iam_workforce_pool_provider` resource ([#21430](https://github.com/hashicorp/terraform-provider-google/pull/21430))
* iambeta: promoted `google_iam_workload_identity_pool` and `google_iam_workload_identity_pool_provider` data sources to GA ([#21408](https://github.com/hashicorp/terraform-provider-google/pull/21408))
* redis: added `kms_key` field to `google_redis_cluster` resource ([#21428](https://github.com/hashicorp/terraform-provider-google/pull/21428))
* tpuv2: added `network_config` field to `google_tpu_v2_queued_resource` resource ([#21426](https://github.com/hashicorp/terraform-provider-google/pull/21426))

BUG FIXES:
* apigee: fixed error when deleting `google_apigee_organization` ([#21473](https://github.com/hashicorp/terraform-provider-google/pull/21473))
* bigtable: fixed a bug where sometimes updating an instance's cluster list could result in an error if there was an existing cluster with autoscaling enabled ([#21503](https://github.com/hashicorp/terraform-provider-google/pull/21503))
* chronicle: fixed bug setting `enabled` on creation in `google_chronicle_rule_deployment` ([#21460](https://github.com/hashicorp/terraform-provider-google/pull/21460))

## 6.21.0 (Feb 18, 2025)

NOTES:
* provider: The Terraform Provider for Google Cloud's regular release date will move from Monday to Tuesday in early March. The 2025/03/10 release will be made on 2025/03/11.
  
FEATURES:
* **New Data Source:** `google_alloydb_instance` ([#21383](https://github.com/hashicorp/terraform-provider-google/pull/21383))
* **New Resource:** `google_firebase_data_connect_service` ([#21368](https://github.com/hashicorp/terraform-provider-google/pull/21368))
* **New Resource:** `google_gemini_data_sharing_with_google_setting` ([#21393](https://github.com/hashicorp/terraform-provider-google/pull/21393))
* **New Resource:** `google_gemini_gemini_gcp_enablement_setting` ([#21357](https://github.com/hashicorp/terraform-provider-google/pull/21357))
* **New Resource:** `google_gemini_logging_setting_binding` ([#21354](https://github.com/hashicorp/terraform-provider-google/pull/21354))
* **New Resource:** `google_gemini_release_channel_setting` ([#21387](https://github.com/hashicorp/terraform-provider-google/pull/21387)
* **New Resource:** `google_gemini_release_channel_setting_binding` ([#21387](https://github.com/hashicorp/terraform-provider-google/pull/21387)
* **New Resource:** `google_netapp_volume_quota_rule` ([#21283](https://github.com/hashicorp/terraform-provider-google/pull/21283))

IMPROVEMENTS:
* accesscontextmanager: added `etag` to access context manager directional policy resources `google_access_context_manager_service_perimeter_dry_run_egress_policy`, `google_access_context_manager_service_perimeter_dry_run_ingress_policy`, `google_access_context_manager_service_perimeter_egress_policy` and `google_access_context_manager_service_perimeter_ingress_policy` to prevent overriding changes ([#21366](https://github.com/hashicorp/terraform-provider-google/pull/21366))
* accesscontextmanager: added `title` field to policy blocks under `google_access_context_manager_service_perimeter` and variants ([#21302](https://github.com/hashicorp/terraform-provider-google/pull/21302))
* artifactregistry: set pageSize to 1000 to speedup `google_artifact_registry_docker_image` data source queries ([#21360](https://github.com/hashicorp/terraform-provider-google/pull/21360))
* compute: added `labels` field to `google_compute_ha_vpn_gateway` resource ([#21385](https://github.com/hashicorp/terraform-provider-google/pull/21385))
* compute: added validation for disk names in `google_compute_disk` ([#21335](https://github.com/hashicorp/terraform-provider-google/pull/21335))
* container: added new fields `container_log_max_size`, `container_log_max_files`, `image_gc_low_threshold_percent`, `image_gc_high_threshold_percent`, `image_minimum_gc_age`, `image_maximum_gc_age`, and `allowed_unsafe_sysctls` to `node_kubelet_config` block in `google_container_cluster` resource. ([#21319](https://github.com/hashicorp/terraform-provider-google/pull/21319))
* monitoring: added `condition_sql` field to `google_monitoring_alert_policy` resource ([#21277](https://github.com/hashicorp/terraform-provider-google/pull/21277))
* networkservices: added `location` field to `google_network_services_mesh` resource ([#21337](https://github.com/hashicorp/terraform-provider-google/pull/21337))
* securitycenter: added `type`, `expiry_time` field to `google_scc_mute_config` resource ([#21318](https://github.com/hashicorp/terraform-provider-google/pull/21318))

BUG FIXES:
* chronicle: fixed creation issues when optional fields were missing for `google_chronicle_rule_deployment` resource ([#21389](https://github.com/hashicorp/terraform-provider-google/pull/21389))
* databasemigrationservice: fixed error details type on `google_database_migration_service_migration_job` ([#21279](https://github.com/hashicorp/terraform-provider-google/pull/21279))
* networkservices: fixed a bug with `google_network_services_authz_extension.wire_format` sending an invalid default value by removing the Terraform default and letting the API set the default. ([#21280](https://github.com/hashicorp/terraform-provider-google/pull/21280))

## 6.20.0 (Feb 10, 2025)

NOTES:
* provider: The Terraform Provider for Google Cloud's regular release date will move from Monday to Tuesday in early March. The 2025/03/10 release will be made on 2025/03/11.
* compute: `google_compute_firewall_policy` now uses MMv1 engine instead of DCL. ([#21235](https://github.com/hashicorp/terraform-provider-google/pull/21235))

FEATURES:
* **New Data Source:** `google_beyondcorp_application_iam_policy` ([#21199](https://github.com/hashicorp/terraform-provider-google/pull/21199))
* **New Data Source:** `google_parameter_manager_parameter_version_render` ([#21104](https://github.com/hashicorp/terraform-provider-google/pull/21104))
* **New Resource:** `google_beyondcorp_application` ([#21199](https://github.com/hashicorp/terraform-provider-google/pull/21199))
* **New Resource:** `google_beyondcorp_application_iam_binding` ([#21199](https://github.com/hashicorp/terraform-provider-google/pull/21199))
* **New Resource:** `google_beyondcorp_application_iam_member` ([#21199](https://github.com/hashicorp/terraform-provider-google/pull/21199))
* **New Resource:** `google_beyondcorp_application_iam_policy` ([#21199](https://github.com/hashicorp/terraform-provider-google/pull/21199))
* **New Resource:** `google_bigquery_analytics_hub_listing_subscription` ([#21189](https://github.com/hashicorp/terraform-provider-google/pull/21189))
* **New Resource:** `google_colab_notebook_execution` ([#21100](https://github.com/hashicorp/terraform-provider-google/pull/21100))
* **New Resource:** `google_colab_schedule` ([#21233](https://github.com/hashicorp/terraform-provider-google/pull/21233))

IMPROVEMENTS:
* accesscontextmanager: added `resource` to `sources` in `egress_from` under resources `google_access_context_manager_service_perimeter`, `google_access_context_manager_service_perimeters`, `google_access_context_manager_service_perimeter_egress_policy`, `google_access_context_manager_service_perimeter_dry_run_egress_policy` ([#21190](https://github.com/hashicorp/terraform-provider-google/pull/21190))
* cloudrunv2: added `base_image_uri` and `build_info` to `google_cloud_run_v2_service` ([#21236](https://github.com/hashicorp/terraform-provider-google/pull/21236))
* colab: added `auto_upgrade` field to `google_colab_runtime` ([#21214](https://github.com/hashicorp/terraform-provider-google/pull/21214))
* colab: added `software_config.post_startup_script_config` field to `google_colab_runtime_template` ([#21200](https://github.com/hashicorp/terraform-provider-google/pull/21200))
* colab: added `desired_state` field to `google_colab_runtime`, making it startable/stoppable. ([#21207](https://github.com/hashicorp/terraform-provider-google/pull/21207))
* compute: added `ip_collection` field to `google_compute_forwarding_rule ` resource ([#21188](https://github.com/hashicorp/terraform-provider-google/pull/21188))
* compute: added `mode` and `allocatable_prefix_length` fields to `google_compute_public_delegated_prefix` resource ([#21216](https://github.com/hashicorp/terraform-provider-google/pull/21216))
* compute: allow parallelization of `google_compute_per_instance_config` and `google_compute_region_per_instance_config` deletions by not locking on the parent resource, but including instance name. ([#21095](https://github.com/hashicorp/terraform-provider-google/pull/21095))
* container: added `auto_monitoring_config` field and subfields to the `google_container_cluster` resource ([#21229](https://github.com/hashicorp/terraform-provider-google/pull/21229))
* filestore: added `initial_replication` field for peer instance configuration and `effective_replication` output for replication configuration output to `google_filestore_instance` ([#21194](https://github.com/hashicorp/terraform-provider-google/pull/21194))
* memorystore: added `CLUSTER_DISABLED`  to `mode` field  in  `google_memorystore_instance` ([#21092](https://github.com/hashicorp/terraform-provider-google/pull/21092))
* networkservices: added `compression_mode` and `allowed_methods` fields to `google_network_services_edge_cache_service` resource ([#21195](https://github.com/hashicorp/terraform-provider-google/pull/21195))
* privateca: added `user_defined_access_urls` and subfields to `google_privateca_certificate_authority` resource to add support for custom CDP AIA URLs ([#21220](https://github.com/hashicorp/terraform-provider-google/pull/21220))
* workbench: added `enable_third_party_identity` field to `google_workbench_instance` resource ([#21265](https://github.com/hashicorp/terraform-provider-google/pull/21265))

BUG FIXES:
* appengine: added a mitigation for an upcoming default change to `standard_scheduler_settings.max_instances` for new `google_app_engine_standard_app_version` resources. If the field is not specified in configuration, diffs will now be ignored. ([#21257](https://github.com/hashicorp/terraform-provider-google/pull/21257))
* bigquery: added diff suppression for legacy values in `renewal_plan` field in `google_bigquery_capacity_commitment` resource ([#21103](https://github.com/hashicorp/terraform-provider-google/pull/21103))
* compute: fixed `google_compute_(region_)resize_request` requiring region/zone to be specified in all cases. They can now be pulled from the provider. ([#21264](https://github.com/hashicorp/terraform-provider-google/pull/21264))
* container: reverted locking behavior in `google_container_node_pool` that caused regression of operation apply time spike started in `v6.15` ([#21102](https://github.com/hashicorp/terraform-provider-google/pull/21102))
* gemini: fixed a bug where the `force_destroy` field in resource `gemini_code_repository_index` did not work properly ([#21212](https://github.com/hashicorp/terraform-provider-google/pull/21212))
* workbench: fixed a bug with `google_workbench_instance` metadata removal not working as expected ([#21204](https://github.com/hashicorp/terraform-provider-google/pull/21204))

## 6.18.1 (January 29, 2025)

BUG FIXES:
* container: fixed a diff caused by server-side set values for `node_config.resource_labels` ([#21082](https://github.com/hashicorp/terraform-provider-google/pull/21082))

## 6.18.0 (January 27, 2025)

FEATURES:
* **New Data Source:** `google_compute_instance_template_iam_policy` ([#20954](https://github.com/hashicorp/terraform-provider-google/pull/20954))
* **New Data Source:** `google_kms_key_handles` ([#20985](https://github.com/hashicorp/terraform-provider-google/pull/20985))
* **New Data Source:** `google_organizations` ([#20965](https://github.com/hashicorp/terraform-provider-google/pull/20965))
* **New Data Source:** `google_parameter_manager_parameter` ([#20953](https://github.com/hashicorp/terraform-provider-google/pull/20953))
* **New Data Source:** `google_parameter_manager_regional_parameters` ([#20958](https://github.com/hashicorp/terraform-provider-google/pull/20958))
* **New Resource:** `google_apihub_api_hub_instance` ([#20948](https://github.com/hashicorp/terraform-provider-google/pull/20948))
* **New Resource:** `google_chronicle_retrohunt` ([#20962](https://github.com/hashicorp/terraform-provider-google/pull/20962))
* **New Resource:** `google_colab_runtime` ([#20940](https://github.com/hashicorp/terraform-provider-google/pull/20940))
* **New Resource:** `google_colab_runtime_template_iam_binding` ([#20963](https://github.com/hashicorp/terraform-provider-google/pull/20963))
* **New Resource:** `google_colab_runtime_template_iam_member` ([#20963](https://github.com/hashicorp/terraform-provider-google/pull/20963))
* **New Resource:** `google_colab_runtime_template_iam_policy` ([#20963](https://github.com/hashicorp/terraform-provider-google/pull/20963))
* **New Resource:** `google_compute_instance_template_iam_binding` ([#20954](https://github.com/hashicorp/terraform-provider-google/pull/20954))
* **New Resource:** `google_compute_instance_template_iam_member` ([#20954](https://github.com/hashicorp/terraform-provider-google/pull/20954))
* **New Resource:** `google_compute_instance_template_iam_policy` ([#20954](https://github.com/hashicorp/terraform-provider-google/pull/20954))
* **New Resource:** `google_gemini_code_repository_index` (GA) ([#20941](https://github.com/hashicorp/terraform-provider-google/pull/20941))
* **New Resource:** `google_gemini_repository_group` (GA) ([#20941](https://github.com/hashicorp/terraform-provider-google/pull/20941))
* **New Resource:** `google_gemini_repository_group_iam_member` (GA) ([#20941](https://github.com/hashicorp/terraform-provider-google/pull/20941))
* **New Resource:** `google_gemini_repository_group_iam_binding` (GA) ([#20941](https://github.com/hashicorp/terraform-provider-google/pull/20941))
* **New Resource:** `google_gemini_repository_group_iam_policy` (GA) ([#20941](https://github.com/hashicorp/terraform-provider-google/pull/20941))
* **New Resource:** `google_parameter_manager_parameter_version` ([#20992](https://github.com/hashicorp/terraform-provider-google/pull/20992))
* **New Resource:** `google_redis_cluster_user_created_connections` ([#20977](https://github.com/hashicorp/terraform-provider-google/pull/20977))

IMPROVEMENTS:
* alloydb: added support for `skip_await_major_version_upgrade` field in `google_alloydb_cluster` resource, allowing for `major_version` to be updated ([#20923](https://github.com/hashicorp/terraform-provider-google/pull/20923))
* apigee: added `properties` field to `google_apigee_environment` resource ([#20932](https://github.com/hashicorp/terraform-provider-google/pull/20932))
* bug: added support for setting `custom_learned_route_priority` to 0 in 'google_compute_router_peer' by adding the `zero_custom_learned_route_priority` field ([#20952](https://github.com/hashicorp/terraform-provider-google/pull/20952))
* cloudrunv2: added `build_config` to `google_cloud_run_v2_service` ([#20979](https://github.com/hashicorp/terraform-provider-google/pull/20979))
* compute: added `pdp_scope` field to `google_compute_public_advertised_prefix` resource ([#20972](https://github.com/hashicorp/terraform-provider-google/pull/20972))
* compute: adding `labels` field to `google_compute_interconnect_attachment` ([#20971](https://github.com/hashicorp/terraform-provider-google/pull/20971))
* compute: fixed a issue where `custom_learned_route_priority` was accidentally set to 0 during updates in 'google_compute_router_peer' ([#20952](https://github.com/hashicorp/terraform-provider-google/pull/20952))
* filestore: added support for `tags` field to `google_filestore_instance` resource ([#20955](https://github.com/hashicorp/terraform-provider-google/pull/20955))
* networksecurity: added `custom_mirroring_profile` and `custom_intercept_profile` fields to `google_network_security_security_profile` and `google_network_security_security_profile_group`  resources ([#20990](https://github.com/hashicorp/terraform-provider-google/pull/20990))
* pubsub: added `enforce_in_transit` fields to `google_pubsub_topic` resource ([#20926](https://github.com/hashicorp/terraform-provider-google/pull/20926))
* pubsub: added `ingestion_data_source_settings.azure_event_hubs` field to `google_pubsub_topic` resource ([#20922](https://github.com/hashicorp/terraform-provider-google/pull/20922))
* redis: added `psc_service_attachments` field to `google_redis_cluster` resource, to enable use of the fine-grained resource `google_redis_cluster_user_created_connections` ([#20977](https://github.com/hashicorp/terraform-provider-google/pull/20977))

BUG FIXES:
* apigee: fixed `properties` field update on `google_apigee_environment` resource ([#20987](https://github.com/hashicorp/terraform-provider-google/pull/20987))
* artifactregistry: fixed perma-diff in `google_artifact_registry_repository` ([#20989](https://github.com/hashicorp/terraform-provider-google/pull/20989))
* compute: fixed failure when creating `google_compute_global_forwarding_rule` with labels targeting PSC endpoint ([#20986](https://github.com/hashicorp/terraform-provider-google/pull/20986))
* container: fixed `additive_vpc_scope_dns_domain` being ignored in Autopilot cluster definition ([#20937](https://github.com/hashicorp/terraform-provider-google/pull/20937))
* container: fixed propagation of `node_pool_defaults.node_config_defaults.insecure_kubelet_readonly_port_enabled` in node config. ([#20936](https://github.com/hashicorp/terraform-provider-google/pull/20936))
* iam: fixed missing result by adding pagination for data source `google_service_accounts`. ([#20966](https://github.com/hashicorp/terraform-provider-google/pull/20966))
* metastore: increased timeout on google_dataproc_metastore_service operations to 75m from 60m. This will expose server-returned reasons for operation failure instead of masking them with a Terraform timeout. ([#20981](https://github.com/hashicorp/terraform-provider-google/pull/20981))
* resourcemanager: added a slightly longer wait (two 10s checks bumped to 15s) for issues with billing associations in `google_project`. Default network deletion should succeed more often. ([#20982](https://github.com/hashicorp/terraform-provider-google/pull/20982))

## 6.17.0 (January 21, 2025)

FEATURES:
* **New Resource:** `google_apigee_environment_addons_config` ([#20851](https://github.com/hashicorp/terraform-provider-google/pull/20851))
* **New Resource:** `google_chronicle_reference_list` (beta) ([#20895](https://github.com/hashicorp/terraform-provider-google/pull/20895))
* **New Resource:** `google_chronicle_rule_deployment` ([#20888](https://github.com/hashicorp/terraform-provider-google/pull/20888))
* **New Resource:** `google_chronicle_rule` ([#20868](https://github.com/hashicorp/terraform-provider-google/pull/20868))
* **New Resource:** `google_colab_runtime_template` ([#20898](https://github.com/hashicorp/terraform-provider-google/pull/20898))
* **New Resource:** `google_edgenetwork_interconnect_attachment` ([#20856](https://github.com/hashicorp/terraform-provider-google/pull/20856))
* **New Resource:** `google_parameter_manager_parameter` ([#20886](https://github.com/hashicorp/terraform-provider-google/pull/20886))
* **New Resource:** `google_parameter_manager_regional_parameter_version` ([#20914](https://github.com/hashicorp/terraform-provider-google/pull/20914))
* **New Resource:** `google_parameter_manager_regional_parameter` ([#20858](https://github.com/hashicorp/terraform-provider-google/pull/20858))

IMPROVEMENTS:
* accesscontextmanager: added `etag` to `google_access_context_manager_service_perimeter_resource` to prevent overriding list of resources ([#20910](https://github.com/hashicorp/terraform-provider-google/pull/20910))
* compute: added `BPS_100G` enum value to `bandwidth` field of `google_compute_interconnect_attachment`. ([#20884](https://github.com/hashicorp/terraform-provider-google/pull/20884))
* compute: added support for `IPV6_ONLY` stack_type to `google_compute_subnetwork`, `google_compute_instance`, `google_compute_instance_template` and `google_compute_region_instance_template`. ([#20850](https://github.com/hashicorp/terraform-provider-google/pull/20850))
* compute: promoted `bgp_best_path_selection_mode `,`bgp_bps_always_compare_med` and `bgp_bps_inter_region_cost ` fields in `google_compute_network` from Beta to Ga ([#20865](https://github.com/hashicorp/terraform-provider-google/pull/20865))
* compute: promoted `next_hop_origin `,`next_hop_med ` and `next_hop_inter_region_cost ` output fields in `google_compute_route` form Beta to GA ([#20865](https://github.com/hashicorp/terraform-provider-google/pull/20865))
* discoveryengine: added `advanced_site_search_config` field to `google_discovery_engine_data_store` resource ([#20912](https://github.com/hashicorp/terraform-provider-google/pull/20912))
* gemini: added `force_destroy` field to resource `google_code_repository_index`, enabling deletion of the resource even when it has dependent RepositoryGroups ([#20881](https://github.com/hashicorp/terraform-provider-google/pull/20881))
* networkservices: added in-place update support for `ports` field on `google_network_services_gateway` resource ([#20908](https://github.com/hashicorp/terraform-provider-google/pull/20908))
* sql: `sql_source_representation_instance` now uses `string` representation of `databaseVersion` ([#20859](https://github.com/hashicorp/terraform-provider-google/pull/20859))
* sql: added `replication_cluster` field to `google_sql_database_instance` resource ([#20889](https://github.com/hashicorp/terraform-provider-google/pull/20889))
* sql: added support of switchover for MySQL and PostgreSQL in `google_sql_database_instance` resource ([#20889](https://github.com/hashicorp/terraform-provider-google/pull/20889))
* workbench: changed `container_image` field of `google_workbench_instance` resource to modifiable. ([#20894](https://github.com/hashicorp/terraform-provider-google/pull/20894))

BUG FIXES:
* apigee: fixed error 404 for `organization` update requests. ([#20854](https://github.com/hashicorp/terraform-provider-google/pull/20854))
* artifactregistry: fixed `artifact_registry_repository` not accepting durations with 'm', 'h' or 'd' ([#20902](https://github.com/hashicorp/terraform-provider-google/pull/20902))
* networkservices: fixed bug where `google_network_services_gateway` could not be updated in place ([#20908](https://github.com/hashicorp/terraform-provider-google/pull/20908))
* storagetransfer: fixed a permadiff with `transfer_spec.aws_s3_data_source.aws_access_key` in `google_storage_transfer_job` ([#20849](https://github.com/hashicorp/terraform-provider-google/pull/20849))

## 6.16.0 (January 13, 2025)

FEATURES:
* **New Resource:** `google_beyondcorp_security_gateway` ([#20844](https://github.com/hashicorp/terraform-provider-google/pull/20844))
* **New Resource:** `google_developer_connect_connection` ([#20823](https://github.com/hashicorp/terraform-provider-google/pull/20823))
* **New Resource:** `google_developer_connect_git_repository_link` ([#20823](https://github.com/hashicorp/terraform-provider-google/pull/20823))

IMPROVEMENTS:
* compute: promoted `standby_policy`, `target_suspended_size`, and `target_stopped_size` fields in `google_compute_region_instance_group_manager` and `google_compute_instance_group_manager` resource from beta to ga ([#20821](https://github.com/hashicorp/terraform-provider-google/pull/20821))
* dns: added `health_check` and `external_endpoints` fields to `google_dns_record_set` resource ([#20843](https://github.com/hashicorp/terraform-provider-google/pull/20843))
* sql: added `server_ca_pool` field to `google_sql_database_instance` resource ([#20834](https://github.com/hashicorp/terraform-provider-google/pull/20834))
* vmwareengine: allowed import of non-STANDARD private clouds in `google_vmwareengine_private_cloud` ([#20832](https://github.com/hashicorp/terraform-provider-google/pull/20832))

BUG FIXES:
* dataproc: fixed boolean fields in `shielded_instance_config` in the `google_dataproc_cluster` resource ([#20828](https://github.com/hashicorp/terraform-provider-google/pull/20828))
* gkeonprem: fixed permadiff on `vcenter` field in `google_gkeonprem_vmware_cluster` resource ([#20837](https://github.com/hashicorp/terraform-provider-google/pull/20837))
* networkservices: fixed `google_network_services_gateway` resource so that it correctly waits for the router to be deleted on `terraform destroy` ([#20817](https://github.com/hashicorp/terraform-provider-google/pull/20817))
* provider: fixed issue where `GOOGLE_CLOUD_QUOTA_PROJECT` env var would override explicit `billing_project` ([#20839](https://github.com/hashicorp/terraform-provider-google/pull/20839))

## 6.15.0 (January 6, 2025)

NOTES:
* compute: `google_compute_firewall_policy_association` now uses MMv1 engine instead of DCL. ([#20744](https://github.com/hashicorp/terraform-provider-google/pull/20744))

DEPRECATIONS:
* compute: deprecated `numeric_id` (string) field in `google_compute_network` resource. Use the new `network_id` (integer)  field instead ([#20698](https://github.com/hashicorp/terraform-provider-google/pull/20698))

FEATURES:
* **New Data Source:** `google_gke_hub_feature` ([#20721](https://github.com/hashicorp/terraform-provider-google/pull/20721))
* **New Resource:** `google_storage_folder` ([#20767](https://github.com/hashicorp/terraform-provider-google/pull/20767))

IMPROVEMENTS:
* artifactregistry: added `vulnerability_scanning_config` field to `google_artifact_registry_repository` resource ([#20726](https://github.com/hashicorp/terraform-provider-google/pull/20726))
* backupdr: promoted datasource `google_backup_dr_backup` to ga ([#20677](https://github.com/hashicorp/terraform-provider-google/pull/20677))
* backupdr: promoted datasource `google_backup_dr_data_source` to ga ([#20677](https://github.com/hashicorp/terraform-provider-google/pull/20677))
* bigquery: added `condition` field to `google_bigquery_dataset_access` resource ([#20707](https://github.com/hashicorp/terraform-provider-google/pull/20707))
* bigquery: added `condition` field to `google_bigquery_dataset` resource ([#20707](https://github.com/hashicorp/terraform-provider-google/pull/20707))
* composer: added `airflow_metadata_retention_config` field to `google_composer_environment` ([#20769](https://github.com/hashicorp/terraform-provider-google/pull/20769))
* compute: added back the validation for `target_service` field on the `google_compute_service_attachment` resource to validade a `ForwardingRule` or `Gateway` URL ([#20711](https://github.com/hashicorp/terraform-provider-google/pull/20711))
* compute: added `availability_domain` field to `google_compute_instance`, `google_compute_instance_template` and `google_compute_region_instance_template` resources ([#20694](https://github.com/hashicorp/terraform-provider-google/pull/20694))
* compute: added `network_id` (integer) field to `google_compute_network` resource and data source ([#20698](https://github.com/hashicorp/terraform-provider-google/pull/20698))
* compute: added `preset_topology` field to `google_network_connectivity_hub` resource ([#20720](https://github.com/hashicorp/terraform-provider-google/pull/20720))
* compute: added `subnetwork_id` field to `google_compute_subnetwork` data source ([#20666](https://github.com/hashicorp/terraform-provider-google/pull/20666))
* compute: made setting resource policies for `google_compute_instance` outside of terraform or using `google_compute_disk_resource_policy_attachment` no longer affect the `boot_disk.initialize_params.resource_policies` field ([#20764](https://github.com/hashicorp/terraform-provider-google/pull/20764))
* container: changed `google_container_cluster` to apply maintenance policy updates after upgrades during cluster update ([#20708](https://github.com/hashicorp/terraform-provider-google/pull/20708))
* container: made nodepool concurrent operations scale better for `google_container_cluster` and `google_container_node_pool` resources ([#20738](https://github.com/hashicorp/terraform-provider-google/pull/20738))
* datastream: added `gtid` and `binary_log_position` fields to `google_datastream_stream` resource ([#20777](https://github.com/hashicorp/terraform-provider-google/pull/20777))
* developerconnect: added support for setting up a `google_developer_connect_connection` resource without specifying the `authorizer_credentials` field ([#20756](https://github.com/hashicorp/terraform-provider-google/pull/20756))
* filestore: added `tags` field to `google_filestore_backup` to allow setting tags for backups at creation time ([#20718](https://github.com/hashicorp/terraform-provider-google/pull/20718))
* networkconnectivity: added `group` field to `google_network_connectivity_spoke` resource ([#20689](https://github.com/hashicorp/terraform-provider-google/pull/20689))
* networkmanagement: promoted `google_network_management_vpc_flow_logs_config` resource to ga ([#20701](https://github.com/hashicorp/terraform-provider-google/pull/20701))
* parallelstore: added `deployment_type` field to `google_parallelstore_instance` resource ([#20785](https://github.com/hashicorp/terraform-provider-google/pull/20785))
* storagetransfer: added `replication_spec` field to `google_storage_transfer_job` resource ([#20788](https://github.com/hashicorp/terraform-provider-google/pull/20788))
* workbench: made `gcs-data-bucket` metadata key modifiable in `google_workbench_instance` resource ([#20728](https://github.com/hashicorp/terraform-provider-google/pull/20728))

BUG FIXES:
* accesscontextmanager: fixed permadiff due to reordering on `google_access_context_manager_service_perimeter_dry_run_egress_policy` `egress_from.identities` ([#20794](https://github.com/hashicorp/terraform-provider-google/pull/20794))
* accesscontextmanager: fixed permadiff due to reordering on `google_access_context_manager_service_perimeter_dry_run_ingress_policy` `ingress_from.identities` ([#20794](https://github.com/hashicorp/terraform-provider-google/pull/20794))
* accesscontextmanager: fixed permadiff due to reordering on `google_access_context_manager_service_perimeter_egress_policy` `egress_from.identities` ([#20794](https://github.com/hashicorp/terraform-provider-google/pull/20794))
* accesscontextmanager: fixed permadiff due to reordering on `google_access_context_manager_service_perimeter_ingress_policy` `ingress_from.identities` ([#20794](https://github.com/hashicorp/terraform-provider-google/pull/20794))
* apigee: fixed 404 error when updating `google_apigee_environment` ([#20745](https://github.com/hashicorp/terraform-provider-google/pull/20745))
* bigquery: fixed DROP COLUMN error with bigquery flexible column names in `google_bigquery_table` ([#20797](https://github.com/hashicorp/terraform-provider-google/pull/20797))
* compute: allowed Service Attachment with Project Number to be used as `google_compute_forwarding_rule.target` ([#20790](https://github.com/hashicorp/terraform-provider-google/pull/20790))
* compute: fixed an issue where `terraform plan -refresh=false` with `google_compute_ha_vpn_gateway.gateway_ip_version` would plan a resource replacement if a full refresh had not been run yet. Terraform now assumes that the value is the default value, `IPV4`, until a refresh is completed. ([#20682](https://github.com/hashicorp/terraform-provider-google/pull/20682))
* compute: fixed panic when zonal resize request fails on `google_compute_resize_request` ([#20734](https://github.com/hashicorp/terraform-provider-google/pull/20734))
* compute: fixed perma-destroy for `psc_data` in `google_compute_region_network_endpoint_group` resource ([#20783](https://github.com/hashicorp/terraform-provider-google/pull/20783))
* compute: fixed `google_compute_instance_guest_attributes` to return an empty list when queried values don't exist instead of throwing an error ([#20760](https://github.com/hashicorp/terraform-provider-google/pull/20760))
* integrationconnectors: allowed `AUTH_TYPE_UNSPECIFIED` option in `google_integration_connectors_connection` resource to support non-standard auth types ([#20782](https://github.com/hashicorp/terraform-provider-google/pull/20782))
* logging: fixed bug in `google_logging_project_bucket_config` when providing `project` in the format of `<project-id-only>` ([#20709](https://github.com/hashicorp/terraform-provider-google/pull/20709))
* networkconnectivity: made `include_export_ranges` and `exclude_export_ranges` fields mutable in `google_network_connectivity_spoke` to avoid recreation of resources ([#20742](https://github.com/hashicorp/terraform-provider-google/pull/20742))
* sql: fixed permadiff when `settings.data_cache_config` is set to false for `google_sql_database_instance` resource ([#20656](https://github.com/hashicorp/terraform-provider-google/pull/20656))
* storage: made `resource_google_storage_bucket_object` generate diff for `md5hash`, `generation`, `crc32c` if content changes ([#20687](https://github.com/hashicorp/terraform-provider-google/pull/20687))
* vertexai: made `contents_delta_uri` an optional field in `google_vertex_ai_index` ([#20780](https://github.com/hashicorp/terraform-provider-google/pull/20780))
* workbench: fixed an issue where a server-added `metadata` tag of `"resource-url"` would not be ignored on `google_workbench_instance` ([#20717](https://github.com/hashicorp/terraform-provider-google/pull/20717))

## 6.14.1 (December 18, 2024)

BUG FIXES:
* compute: fixed an issue where `google_compute_firewall_policy_rule` was incorrectly removed from the Terraform state ([#20733](https://github.com/hashicorp/terraform-provider-google/pull/20733))

## 6.14.0 (December 16, 2024)

FEATURES:
* **New Resource:** `google_network_security_intercept_deployment_group` ([#20615](https://github.com/hashicorp/terraform-provider-google/pull/20615))
* **New Resource:** `google_network_security_intercept_deployment` ([#20634](https://github.com/hashicorp/terraform-provider-google/pull/20634))
* **New Resource:** `google_network_security_authz_policy` ([#20595](https://github.com/hashicorp/terraform-provider-google/pull/20595))
* **New Resource:** `google_network_services_authz_extension` ([#20595](https://github.com/hashicorp/terraform-provider-google/pull/20595))

IMPROVEMENTS:
* compute: `google_compute_instance` is no longer recreated when changing `boot_disk.auto_delete` ([#20580](https://github.com/hashicorp/terraform-provider-google/pull/20580))
* compute: added `CA_ENTERPRISE_ANNUAL` option for field `cloud_armor_tier` in `google_compute_project_cloud_armor_tier` resource ([#20596](https://github.com/hashicorp/terraform-provider-google/pull/20596))
* compute: added `network_tier` field to `google_compute_global_forwarding_rule` resource ([#20582](https://github.com/hashicorp/terraform-provider-google/pull/20582))
* compute: added `rule.rate_limit_options.enforce_on_key_configs` field to `google_compute_security_policy` resource ([#20597](https://github.com/hashicorp/terraform-provider-google/pull/20597))
* compute: made `metadata_startup_script` able to be updated via graceful switch in `google_compute_instance` ([#20655](https://github.com/hashicorp/terraform-provider-google/pull/20655))
* container: added field `enable_fqdn_network_policy` to resource `google_container_cluster` ([#20609](https://github.com/hashicorp/terraform-provider-google/pull/20609))
* identityplatform: marked `quota.0.sign_up_quota_config` subfields conditionally required in `google_identity_platform_config` to move errors from apply time up to plan time, and clarified the rule in documentation ([#20627](https://github.com/hashicorp/terraform-provider-google/pull/20627))
* networkconnectivity: added support for updating `linked_vpn_tunnels.include_import_ranges`, `linked_interconnect_attachments.include_import_ranges`, `linked_router_appliance_instances. instances` and `linked_router_appliance_instances.include_import_ranges` in `google_network_connectivity_spoke` ([#20650](https://github.com/hashicorp/terraform-provider-google/pull/20650))
* storage: added `hdfs_data_source` field to `google_storage_transfer_job` resource ([#20583](https://github.com/hashicorp/terraform-provider-google/pull/20583))
* tpuv2: added `network_configs` and `network_config.queue_count` fields to `google_tpu_v2_vm` resource ([#20621](https://github.com/hashicorp/terraform-provider-google/pull/20621))

BUG FIXES:
* accesscontextmanager: fixed an update bug in `google_access_context_manager_perimeter` by removing the broken output-only `etag` field in `google_access_context_manager_perimeter` and `google_access_context_manager_perimeters` ([#20691](https://github.com/hashicorp/terraform-provider-google/pull/20691))
* compute: fixed permadiff on the `recaptcha_options` field for `google_compute_security_policy` resource ([#20617](https://github.com/hashicorp/terraform-provider-google/pull/20617))
* compute: fixed issue where updating labels on `resource_google_compute_resource_policy` would fail because of a patch error with `guest_flush` ([#20632](https://github.com/hashicorp/terraform-provider-google/pull/20632))
* networkconnectivity: fixed `linked_router_appliance_instances.instances.virtual_machine` and `linked_router_appliance_instances.instances.ip_address` attributes in `google_network_connectivity_spoke` to be correctly marked as required. Otherwise the request to create the resource will fail. ([#20650](https://github.com/hashicorp/terraform-provider-google/pull/20650))
* privateca: fixed an issue which causes error when updating labels for activated sub-CA ([#20630](https://github.com/hashicorp/terraform-provider-google/pull/20630))
* sql: fixed permadiff when 'settings.data_cache_config' is set to false for 'google_sql_database_instance' resource ([#20656](https://github.com/hashicorp/terraform-provider-google/pull/20656))

## 6.13.0 (December 9, 2024)

NOTES:
* New [ephemeral resources](https://developer.hashicorp.com/terraform/language/v1.10.x/resources/ephemeral) `google_service_account_access_token`, `google_service_account_id_token`, `google_service_account_jwt`, `google_service_account_key` now support [ephemeral values](https://developer.hashicorp.com/terraform/language/v1.10.x/values/variables#exclude-values-from-state).
* iam3: promoted resources `google_iam_principal_access_boundary_policy`, `google_iam_organizations_policy_binding`, `google_iam_folders_policy_binding` and `google_iam_projects_policy_binding` to GA ([#20475](https://github.com/hashicorp/terraform-provider-google/pull/20475))
DEPRECATIONS:
* gkehub: deprecated `configmanagement.config_sync.metrics_gcp_service_account_email` in `google_gke_hub_feature_membership` resource ([#20561](https://github.com/hashicorp/terraform-provider-google/pull/20561))

FEATURES:
* **New Ephemeral Resource:** `google_service_account_access_token` ([#20542](https://github.com/hashicorp/terraform-provider-google/pull/20542))
* **New Ephemeral Resource:** `google_service_account_id_token` ([#20542](https://github.com/hashicorp/terraform-provider-google/pull/20542))
* **New Ephemeral Resource:** `google_service_account_jwt` ([#20542](https://github.com/hashicorp/terraform-provider-google/pull/20542))
* **New Ephemeral Resource:** `google_service_account_key` ([#20542](https://github.com/hashicorp/terraform-provider-google/pull/20542))
* **New Data Source:** `google_backup_dr_backup_vault` ([#20468](https://github.com/hashicorp/terraform-provider-google/pull/20468))
* **New Data Source:** `google_composer_user_workloads_config_map` (GA) ([#20478](https://github.com/hashicorp/terraform-provider-google/pull/20478))
* **New Data Source:** `google_composer_user_workloads_secret` (GA) ([#20478](https://github.com/hashicorp/terraform-provider-google/pull/20478))
* **New Resource:** `google_composer_user_workloads_config_map` (GA) ([#20478](https://github.com/hashicorp/terraform-provider-google/pull/20478))
* **New Resource:** `google_composer_user_workloads_secret` (GA) ([#20478](https://github.com/hashicorp/terraform-provider-google/pull/20478))
* **New Resource:** `google_gemini_code_repository_index` ([#20474](https://github.com/hashicorp/terraform-provider-google/pull/20474))
* **New Resource:** `google_network_security_mirroring_deployment` ([#20489](https://github.com/hashicorp/terraform-provider-google/pull/20489))
* **New Resource:** `google_network_security_mirroring_deployment_group` ([#20489](https://github.com/hashicorp/terraform-provider-google/pull/20489))
* **New Resource:** `google_network_security_mirroring_endpoint_group_association` ([#20489](https://github.com/hashicorp/terraform-provider-google/pull/20489))
* **New Resource:** `google_network_security_mirroring_endpoint_group` ([#20489](https://github.com/hashicorp/terraform-provider-google/pull/20489))

IMPROVEMENTS:
* accesscontextmanager: added `etag` to `google_access_context_manager_service_perimeter` and `google_access_context_manager_service_perimeters` ([#20455](https://github.com/hashicorp/terraform-provider-google/pull/20455))
* alloydb: increased default timeout on `google_alloydb_cluster` to 120m from 30m ([#20547](https://github.com/hashicorp/terraform-provider-google/pull/20547))
* bigtable: added `row_affinity` field to `google_bigtable_app_profile` resource ([#20435](https://github.com/hashicorp/terraform-provider-google/pull/20435))
* cloudbuild: added `private_service_connect` field to `google_cloudbuild_worker_pool` resource ([#20561](https://github.com/hashicorp/terraform-provider-google/pull/20561))
* clouddeploy: added `associated_entities` field to `google_clouddeploy_target` resource ([#20561](https://github.com/hashicorp/terraform-provider-google/pull/20561))
* clouddeploy: added `serial_pipeline.strategy.canary.runtime_config.kubernetes.gateway_service_mesh.route_destinations` field to `google_clouddeploy_delivery_pipeline` resource ([#20561](https://github.com/hashicorp/terraform-provider-google/pull/20561))
* composer: added multiple composer 3 related fields to `google_composer_environment` (GA) ([#20478](https://github.com/hashicorp/terraform-provider-google/pull/20478))
* compute: `google_compute_instance`, `google_compute_instance_template`, `google_compute_region_instance_template` now supports `advanced_machine_features.enable_uefi_networking` field ([#20531](https://github.com/hashicorp/terraform-provider-google/pull/20531))
* compute: added support for specifying storage pool with name or partial url ([#20502](https://github.com/hashicorp/terraform-provider-google/pull/20502))
* compute: added `numeric_id` to the `google_compute_network` data source ([#20548](https://github.com/hashicorp/terraform-provider-google/pull/20548))
* compute: added `threshold_configs` field to `google_compute_security_policy` resource ([#20545](https://github.com/hashicorp/terraform-provider-google/pull/20545))
* compute: added server generated id as `forwarding_rule_id` to `google_compute_global_forwarding_rule` ([#20404](https://github.com/hashicorp/terraform-provider-google/pull/20404))
* compute: added server generated id as `health_check_id` to `google_region_health_check` ([#20404](https://github.com/hashicorp/terraform-provider-google/pull/20404))
* compute: added server generated id as `instance_group_manager_id` to `google_instance_group_manager` ([#20404](https://github.com/hashicorp/terraform-provider-google/pull/20404))
* compute: added server generated id as `instance_group_manager_id` to `google_region_instance_group_manager` ([#20404](https://github.com/hashicorp/terraform-provider-google/pull/20404))
* compute: added server generated id as `network_endpoint_id` to `google_region_network_endpoint` ([#20404](https://github.com/hashicorp/terraform-provider-google/pull/20404))
* compute: added server generated id as `subnetwork_id` to `google_subnetwork` ([#20404](https://github.com/hashicorp/terraform-provider-google/pull/20404))
* compute: added the `psc_data` field to the `google_compute_region_network_endpoint_group` resource ([#20454](https://github.com/hashicorp/terraform-provider-google/pull/20454))
* container: added `enterprise_config` field to `google_container_cluster` resource ([#20534](https://github.com/hashicorp/terraform-provider-google/pull/20534))
* container: added `node_pool_autoconfig.linux_node_config.cgroup_mode` field to `google_container_cluster` resource ([#20460](https://github.com/hashicorp/terraform-provider-google/pull/20460))
* dataproc: added `autotuning_config` and `cohort` fields to `google_dataproc_batch` ([#20410](https://github.com/hashicorp/terraform-provider-google/pull/20410))
* dataproc: added `cluster_config.preemptible_worker_config.instance_flexibility_policy.provisioning_model_mix` field to `google_dataproc_cluster` resource ([#20396](https://github.com/hashicorp/terraform-provider-google/pull/20396))
* dataproc: added `confidential_instance_config` field to `google_dataproc_cluster` resource ([#20488](https://github.com/hashicorp/terraform-provider-google/pull/20488))
* discoveryengine: added `HEALTHCARE_FHIR` to `industry_vertical` field in `google_discovery_engine_search_engine` ([#20471](https://github.com/hashicorp/terraform-provider-google/pull/20471))
* gkehub: added `configmanagement.config_sync.stop_syncing` field to `google_gke_hub_feature_membership` resource ([#20561](https://github.com/hashicorp/terraform-provider-google/pull/20561))
* monitoring: added `disable_metric_validation` field to `google_monitoring_alert_policy` resource ([#20544](https://github.com/hashicorp/terraform-provider-google/pull/20544))
* oracledatabase: added `deletion_protection` field to `google_oracle_database_autonomous_database` ([#20484](https://github.com/hashicorp/terraform-provider-google/pull/20484))
* oracledatabase: added `deletion_protection` field to `google_oracle_database_cloud_exadata_infrastructure` ([#20485](https://github.com/hashicorp/terraform-provider-google/pull/20485))
* oracledatabase: added `deletion_protection` field to `google_oracle_database_cloud_vm_cluster ` ([#20392](https://github.com/hashicorp/terraform-provider-google/pull/20392))
* parallelstore: added `deployment_type` to `google_parallelstore_instance` ([#20457](https://github.com/hashicorp/terraform-provider-google/pull/20457))
* resourcemanager: made `google_service_account` `email` and `member` fields available during plan ([#20510](https://github.com/hashicorp/terraform-provider-google/pull/20510))

BUG FIXES:
* apigee: made `google_apigee_organization` wait for deletion operation to complete. ([#20504](https://github.com/hashicorp/terraform-provider-google/pull/20504))
* cloudfunctions: fixed issue when updating `vpc_connector_egress_settings` field for `google_cloudfunctions_function` resource. ([#20437](https://github.com/hashicorp/terraform-provider-google/pull/20437))
* dataproc: ensured oneOf condition is honored when expanding the job configuration for Hive, Pig, Spark-sql, and Presto in `google_dataproc_job`. ([#20453](https://github.com/hashicorp/terraform-provider-google/pull/20453))
* gkehub: fixed allowable value `INSTALLATION_UNSPECIFIED` in `template_library.installation` ([#20567](https://github.com/hashicorp/terraform-provider-google/pull/20567))
* sql: fixed edition downgrade failure for an `ENTERPRISE_PLUS` instance with data cache enabled. ([#20393](https://github.com/hashicorp/terraform-provider-google/pull/20393))


## 6.12.0 (November 18, 2024)

FEATURES:
* **New Data Source:** `google_access_context_manager_access_policy` ([#20295](https://github.com/hashicorp/terraform-provider-google/pull/20295))
* **New Resource:** `google_dataproc_gdc_spark_application` ([#20242](https://github.com/hashicorp/terraform-provider-google/pull/20242))
* **New Resource:** `google_managed_kafka_cluster` and `google_managed_kafka_topic` ([#20237](https://github.com/hashicorp/terraform-provider-google/pull/20237))

IMPROVEMENTS:
* artifactregistry: added `common_repository` field to `google_artifact_registry_repository` resource ([#20305](https://github.com/hashicorp/terraform-provider-google/pull/20305))
* cloudrunv2: added `urls` output field to `google_cloud_run_v2_service` resource ([#20313](https://github.com/hashicorp/terraform-provider-google/pull/20313))
* compute: added `IDPF` as a possible value for the `network_interface.nic_type` field in `google_compute_instance` resource ([#20250](https://github.com/hashicorp/terraform-provider-google/pull/20250))
* compute: added `IDPF` as a possible value for the `guest_os_features.type` field in `google_compute_image` resource ([#20250](https://github.com/hashicorp/terraform-provider-google/pull/20250))
* compute: added `replica_names` field to `sql_database_instance` resource ([#20202](https://github.com/hashicorp/terraform-provider-google/pull/20202))
* filestore: added `performance_config` field to `google_filestore_instance` ([#20218](https://github.com/hashicorp/terraform-provider-google/pull/20218))
* redis: added `persistence_config` to `google_redis_cluster`. ([#20212](https://github.com/hashicorp/terraform-provider-google/pull/20212))
* securesourcemanager: added `workforce_identity_federation_config` field to `google_secure_source_manager_instance` resource ([#20290](https://github.com/hashicorp/terraform-provider-google/pull/20290))
* spanner: added `default_backup_schedule_type` field to  `google_spanner_instance` ([#20213](https://github.com/hashicorp/terraform-provider-google/pull/20213))
* sql: added `psc_auto_connections` fields to `google_sql_database_instance` resource ([#20307](https://github.com/hashicorp/terraform-provider-google/pull/20307))

BUG FIXES:
* accesscontextmanager: fixed permadiff in perimeter `google_access_context_manager_service_perimeter_ingress_policy` and `google_access_context_manager_service_perimeter_egress_policy` resources when there are duplicate resources in the rules ([#20294](https://github.com/hashicorp/terraform-provider-google/pull/20294))
* * accesscontextmanager: fixed comparison of `identity_type` in `ingress_from` and `egress_from` when the `IDENTITY_TYPE_UNSPECIFIED` is set ([#20221](https://github.com/hashicorp/terraform-provider-google/pull/20221))
* compute: fixed permadiff on attempted `type` field updates in `google_computer_security_policy`, updating this field will now force recreation of the resource ([#20316](https://github.com/hashicorp/terraform-provider-google/pull/20316))
* identityplatform: fixed perma-diff originating from the `sign_in.anonymous.enabled` field in `google_identity_platform_config` ([#20244](https://github.com/hashicorp/terraform-provider-google/pull/20244))

## 6.11.2 (November 15, 2024)

BUG FIXES:
* vertexai: fixed issue with google_vertex_ai_endpoint where upgrading to 6.11.0 would delete all traffic splits that were set outside Terraform (which was previously a required step for all meaningful use of this resource). ([#20350](https://github.com/hashicorp/terraform-provider-google/pull/20350))

## 6.11.1 (November 12, 2024)

BUG FIXES:
* container: fixed diff on `google_container_cluster.user_managed_keys_config` field for resources that had not set it. ([#20314](https://github.com/hashicorp/terraform-provider-google/pull/20314))
* container: marked `google_container_cluster.user_managed_keys_config` as immutable because it can't be updated in place. ([#20314](https://github.com/hashicorp/terraform-provider-google/pull/20314))

## 6.11.0 (November 11, 2024)

NOTES:
* compute: migrated `google_compute_firewall_policy_rule` from DCL engine to MMv1 engine. ([#20160](https://github.com/hashicorp/terraform-provider-google/pull/20160))

BREAKING CHANGES:
* looker: made `oauth_config` a required field in `google_looker_instance`, as creating this resource without that field always triggers an API error ([#20196](https://github.com/hashicorp/terraform-provider-google/pull/20196))

FEATURES:
* **New Data Source:** `google_spanner_database` ([#20114](https://github.com/hashicorp/terraform-provider-google/pull/20114))
* **New Resource:** `google_apigee_api` ([#20113](https://github.com/hashicorp/terraform-provider-google/pull/20113))
* **New Resource:** `google_dataproc_gdc_application_environment` ([#20165](https://github.com/hashicorp/terraform-provider-google/pull/20165))
* **New Resource:** `google_dataproc_gdc_service_instance` ([#20147](https://github.com/hashicorp/terraform-provider-google/pull/20147))
* **New Resource:** `google_memorystore_instance` ([#20108](https://github.com/hashicorp/terraform-provider-google/pull/20108))

IMPROVEMENTS:
* apigee: added in-place update support for `google_apigee_env_references` ([#20182](https://github.com/hashicorp/terraform-provider-google/pull/20182))
* apigee: added in-place update support for `google_apigee_environment` resource ([#20189](https://github.com/hashicorp/terraform-provider-google/pull/20189))
* cloudrun: added `empty_dir` field to `google_cloud_run_service` ([#20185](https://github.com/hashicorp/terraform-provider-google/pull/20185))
* cloudrunv2: added `empty_dir` field to `google_cloud_run_v2_service` and `google_cloud_run_v2_job` ([#20185](https://github.com/hashicorp/terraform-provider-google/pull/20185))
* compute: added `disks` field to `google_compute_node_template` resource ([#20180](https://github.com/hashicorp/terraform-provider-google/pull/20180))
* compute: added `preconfigured_waf_config` field  to `google_compute_security_policy` resource ([#20183](https://github.com/hashicorp/terraform-provider-google/pull/20183))
* compute: added `replica_names` field to `sql_database_instance` resource ([#20202](https://github.com/hashicorp/terraform-provider-google/pull/20202))
* compute: added `instance_flexibility_policy` field to `google_compute_region_instance_group_manager` resource ([#20132](https://github.com/hashicorp/terraform-provider-google/pull/20132))
* compute: increased `google_compute_security_policy` timeouts from 20 minutes to 30 minutes ([#20145](https://github.com/hashicorp/terraform-provider-google/pull/20145))
* container: added `control_plane_endpoints_config` field to `google_container_cluster` resource. ([#20193](https://github.com/hashicorp/terraform-provider-google/pull/20193))
* container: added `parallelstore_csi_driver_config` field to `google_container_cluster` resource. ([#20163](https://github.com/hashicorp/terraform-provider-google/pull/20163))
* container: added `user_managed_keys_config` field to `google_container_cluster` resource. ([#20105](https://github.com/hashicorp/terraform-provider-google/pull/20105))
* firestore: allowed single field indexes to support `__name__ DESC` indexes in `google_firestore_index` resources ([#20124](https://github.com/hashicorp/terraform-provider-google/pull/20124))
* privateca: added support for `google_privateca_certificate_authority` with type = "SUBORDINATE" to be activated into "STAGED" state ([#20103](https://github.com/hashicorp/terraform-provider-google/pull/20103))
* spanner: added `default_backup_schedule_type` field to  `google_spanner_instance` ([#20213](https://github.com/hashicorp/terraform-provider-google/pull/20213))
* vertexai: added `traffic_split`, `private_service_connect_config`, `predict_request_response_logging_config`, `dedicated_endpoint_enabled`, and `dedicated_endpoint_dns` fields to `google_vertex_ai_endpoint` resource ([#20179](https://github.com/hashicorp/terraform-provider-google/pull/20179))
* workflows: added `deletion_protection` field to `google_workflows_workflow` resource ([#20106](https://github.com/hashicorp/terraform-provider-google/pull/20106))

BUG FIXES:
* compute: fixed a diff based on server-side reordering of `match.src_address_groups` and `match.dest_address_groups` in `google_compute_network_firewall_policy_rule` ([#20148](https://github.com/hashicorp/terraform-provider-google/pull/20148))
* compute: fixed permadiff on the `preconfigured_waf_config` field for `google_compute_security_policy` resource ([#20183](https://github.com/hashicorp/terraform-provider-google/pull/20183))
* container: fixed in-place updates for `node_config.containerd_config` in `google_container_cluster` and `google_container_node_pool` ([#20112](https://github.com/hashicorp/terraform-provider-google/pull/20112))

## 6.10.0 (November 4, 2024)

FEATURES:
* **New Data Source:** `google_compute_instance_guest_attributes` ([#20095](https://github.com/hashicorp/terraform-provider-google/pull/20095))
* **New Data Source:** `google_service_accounts` ([#20062](https://github.com/hashicorp/terraform-provider-google/pull/20062))
* **New Resource:** `google_iap_settings` ([#20085](https://github.com/hashicorp/terraform-provider-google/pull/20085))

IMPROVEMENTS:
* apphub: added `GLOBAL` enum value to `scope.type` field in `google_apphub_application` resource ([#20015](https://github.com/hashicorp/terraform-provider-google/pull/20015))
* assuredworkloads: added `workload_options` field to `google_assured_workloads_workload` resource ([#19985](https://github.com/hashicorp/terraform-provider-google/pull/19985))
* bigquery: added `external_catalog_dataset_options` fields to `google_bigquery_dataset` resource (beta) ([#20097](https://github.com/hashicorp/terraform-provider-google/pull/20097))
* bigquery: added descriptive validation errors for missing required fields in `google_bigquery_job` destination table configuration ([#20077](https://github.com/hashicorp/terraform-provider-google/pull/20077))
* compute: `desired_status` on google_compute_instance can now be set to `TERMINATED` or `SUSPENDED` on instance creation ([#20031](https://github.com/hashicorp/terraform-provider-google/pull/20031))
* compute: added `header_action` and `redirect_options` fields  to `google_compute_security_policy_rule` resource ([#20079](https://github.com/hashicorp/terraform-provider-google/pull/20079))
* compute: added `interface.ipv6-address` field in `google_compute_external_vpn_gateway` resource ([#20091](https://github.com/hashicorp/terraform-provider-google/pull/20091))
* compute: added `propagated_connection_limit` and `connected_endpoints.propagated_connection_count` fields to `google_compute_service_attachment` resource ([#20016](https://github.com/hashicorp/terraform-provider-google/pull/20016))
* compute: added plan-time validation to `name` on `google_compute_instance` ([#20036](https://github.com/hashicorp/terraform-provider-google/pull/20036))
* compute: added support for `advanced_machine_features.turbo_mode` to `google_compute_instance`, `google_compute_instance_template`, and `google_compute_region_instance_template` ([#20090](https://github.com/hashicorp/terraform-provider-google/pull/20090))
* container: added in-place update support for `labels`, `resource_manager_tags` and `workload_metadata_config` in `google_container_cluster.node_config` ([#20038](https://github.com/hashicorp/terraform-provider-google/pull/20038))
* filestore: added `protocol` property to resource `google_filestore_instance` ([#19982](https://github.com/hashicorp/terraform-provider-google/pull/19982))
* memorystore: added `mode` flag to `google_memorystore_instance` ([#19988](https://github.com/hashicorp/terraform-provider-google/pull/19988))
* netapp: added `zone` and `replica_zone` fields to `google_netapp_storage_pool` resource ([#19980](https://github.com/hashicorp/terraform-provider-google/pull/19980))
* netapp: added `zone` and `replica_zone` fields to `google_netapp_volume` resource ([#19980](https://github.com/hashicorp/terraform-provider-google/pull/19980))
* networksecurity: added `tls_inspection_policy` field to `google_network_security_gateway_security_policy` ([#19986](https://github.com/hashicorp/terraform-provider-google/pull/19986))
* resourcemanager: added `disabled` to `google_service_account` datasource ([#20034](https://github.com/hashicorp/terraform-provider-google/pull/20034))
* spanner: added `asymmetric_autoscaling_options` field to  `google_spanner_instance` ([#20014](https://github.com/hashicorp/terraform-provider-google/pull/20014))
* sql: removed the client-side default of `ENTERPRISE` for `edition` in `google_sql_database_instance` so that `edition` is determined by the API when unset. This will cause new instances to use `ENTERPRISE_PLUS` as the default for POSTGRES_16. ([#19977](https://github.com/hashicorp/terraform-provider-google/pull/19977))
* vmwareengine: added `autoscaling_settings` to `google_vmwareengine_private_cloud` resource ([#20057](https://github.com/hashicorp/terraform-provider-google/pull/20057))

BUG FIXES:
* accesscontextmanager: fixed permadiff for perimeter ingress / egress rule resources ([#20046](https://github.com/hashicorp/terraform-provider-google/pull/20046))
* compute: fixed an error in `google_compute_security_policy_rule` that prevented updating the default rule ([#20066](https://github.com/hashicorp/terraform-provider-google/pull/20066))
* container: fixed missing in-place updates for some `google_container_cluster.node_config` subfields ([#20038](https://github.com/hashicorp/terraform-provider-google/pull/20038))

## 6.9.0 (October 28, 2024)

DEPRECATIONS:
* containerattached: deprecated `security_posture_config` field in `google_container_attached_cluster` resource ([#19912](https://github.com/hashicorp/terraform-provider-google/pull/19912))

FEATURES:
* **New Data Source:** `google_oracle_database_autonomous_database` ([#19903](https://github.com/hashicorp/terraform-provider-google/pull/19903))
* **New Data Source:** `google_oracle_database_autonomous_databases` ([#19901](https://github.com/hashicorp/terraform-provider-google/pull/19901))
* **New Data Source:** `google_oracle_database_cloud_exadata_infrastructures` ([#19884](https://github.com/hashicorp/terraform-provider-google/pull/19884))
* **New Data Source:** `google_oracle_database_cloud_vm_clusters` ([#19900](https://github.com/hashicorp/terraform-provider-google/pull/19900))
* **New Resource:** `google_apigee_app_group` ([#19921](https://github.com/hashicorp/terraform-provider-google/pull/19921))
* **New Resource:** `google_apigee_developer` ([#19911](https://github.com/hashicorp/terraform-provider-google/pull/19911))
* **New Resource:** `google_network_connectivity_group` ([#19902](https://github.com/hashicorp/terraform-provider-google/pull/19902))

IMPROVEMENTS:
* compute: `google_compute_network_firewall_policy_association` now uses MMv1 engine instead of DCL. ([#19976](https://github.com/hashicorp/terraform-provider-google/pull/19976))
* compute: `google_compute_region_network_firewall_policy_association` now uses MMv1 engine instead of DCL. ([#19976](https://github.com/hashicorp/terraform-provider-google/pull/19976))
* compute: added `creation_timestamp` field to `google_compute_instance`, `google_compute_instance_template`, `google_compute_region_instance_template` ([#19906](https://github.com/hashicorp/terraform-provider-google/pull/19906))
* compute: added `key_revocation_action_type` to `google_compute_instance` and related resources ([#19952](https://github.com/hashicorp/terraform-provider-google/pull/19952))
* looker: added `deletion_policy` to `google_looker_instance` to allow force-destroying instances with nested resources by setting `deletion_policy = FORCE` ([#19924](https://github.com/hashicorp/terraform-provider-google/pull/19924))
* monitoring: added `alert_strategy.notification_prompts` field to `google_monitoring_alert_policy` ([#19928](https://github.com/hashicorp/terraform-provider-google/pull/19928))
* storage: added `hierarchical_namespace` to `google_storage_bucket` resource ([#19882](https://github.com/hashicorp/terraform-provider-google/pull/19882))
* sql: removed the client-side default of `ENTERPRISE` for `edition` in `google_sql_database_instance` so that `edition` is determined by the API when unset. This will cause new instances to use `ENTERPRISE_PLUS` as the default for POSTGRES_16. ([#19977](https://github.com/hashicorp/terraform-provider-google/pull/19977))
* vmwareengine: added `autoscaling_settings` to `google_vmwareengine_cluster` resource ([#19962](https://github.com/hashicorp/terraform-provider-google/pull/19962))
* workstations: added `max_usable_workstations` field to `google_workstations_workstation_config` resource. ([#19872](https://github.com/hashicorp/terraform-provider-google/pull/19872))

BUG FIXES:
* compute: fixed an issue where immutable `distribution_zones` was incorrectly sent to the API when updating `distribution_policy_target_shape` in `google_compute_region_instance_group_manager` resource ([#19949](https://github.com/hashicorp/terraform-provider-google/pull/19949))
* container: fixed a crash in `google_container_node_pool` caused by an occasional nil pointer ([#19922](https://github.com/hashicorp/terraform-provider-google/pull/19922))
* essentialcontacts: fixed `google_essential_contacts_contact` import to include required parent field. ([#19877](https://github.com/hashicorp/terraform-provider-google/pull/19877))
* sql: made `google_sql_database_instance.0.settings.0.data_cache_config` accept server-side changes when unset. When unset, no diffs will be created when instances change in `edition` and the feature is enabled or disabled as a result. ([#19972](https://github.com/hashicorp/terraform-provider-google/pull/19972))
* storage: removed retry on 404s during refresh for `google_storage_bucket`, preventing hanging when refreshing deleted buckets ([#19964](https://github.com/hashicorp/terraform-provider-google/pull/19964))

## 6.8.0 (October 21, 2024)

FEATURES:
* **New Data Source:** `google_oracle_database_cloud_exadata_infrastructure` ([#19856](https://github.com/hashicorp/terraform-provider-google/pull/19856))
* **New Data Source:** `google_oracle_database_cloud_vm_cluster` ([#19859](https://github.com/hashicorp/terraform-provider-google/pull/19859))
* **New Data Source:** `google_oracle_database_db_nodes` ([#19871](https://github.com/hashicorp/terraform-provider-google/pull/19871))
* **New Data Source:** `google_oracle_database_db_servers` ([#19823](https://github.com/hashicorp/terraform-provider-google/pull/19823))
* **New Resource:** `google_oracle_database_autonomous_database` ([#19860](https://github.com/hashicorp/terraform-provider-google/pull/19860))
* **New Resource:** `google_oracle_database_cloud_exadata_infrastructure` ([#19798](https://github.com/hashicorp/terraform-provider-google/pull/19798))
* **New Resource:** `google_oracle_database_cloud_vm_cluster` ([#19837](https://github.com/hashicorp/terraform-provider-google/pull/19837))
* **New Resource:** `google_transcoder_job_template` ([#19854](https://github.com/hashicorp/terraform-provider-google/pull/19854))
* **New Resource:** `google_transcoder_job` ([#19854](https://github.com/hashicorp/terraform-provider-google/pull/19854))

IMPROVEMENTS:
* cloudfunctions: increased the timeouts to 20 minutes for `google_cloudfunctions_function` resource ([#19799](https://github.com/hashicorp/terraform-provider-google/pull/19799))
* cloudrunv2: added `invoker_iam_disabled` field to `google_cloud_run_v2_service` ([#19833](https://github.com/hashicorp/terraform-provider-google/pull/19833))
* compute: made `google_compute_network_firewall_policy_rule` use MMv1 engine instead of DCL. ([#19862](https://github.com/hashicorp/terraform-provider-google/pull/19862))
* compute: made `google_compute_region_network_firewall_policy_rule` use MMv1 engine instead of DCL. ([#19862](https://github.com/hashicorp/terraform-provider-google/pull/19862))
* compute: added `ip_address_selection_policy` field to `google_compute_backend_service` and `google_compute_region_backend_service`. ([#19863](https://github.com/hashicorp/terraform-provider-google/pull/19863))
* compute: added `provisioned_throughput` field to `google_compute_instance_template` resource ([#19852](https://github.com/hashicorp/terraform-provider-google/pull/19852))
* compute: added `provisioned_throughput` field to `google_compute_region_instance_template` resource ([#19852](https://github.com/hashicorp/terraform-provider-google/pull/19852))
* container: added support for additional values `KCP_CONNECTION`, and `KCP_SSHD`in `google_container_cluster.logging_config` ([#19812](https://github.com/hashicorp/terraform-provider-google/pull/19812))
* dialogflowcx: added `advanced_settings.logging_settings` and `advanced_settings.speech_settings` to `google_dialogflow_cx_agent` and `google_dialogflow_cx_flow` ([#19801](https://github.com/hashicorp/terraform-provider-google/pull/19801))
* networkconnectivity: added `linked_producer_vpc_network` field to `google_network_connectivity_spoke` resource ([#19806](https://github.com/hashicorp/terraform-provider-google/pull/19806))
* secretmanager: added `is_secret_data_base64` field to `google_secret_manager_secret_version` and `google_secret_manager_secret_version_access` datasources ([#19831](https://github.com/hashicorp/terraform-provider-google/pull/19831))
* secretmanager: added `is_secret_data_base64` field to `google_secret_manager_regional_secret_version` and `google_secret_manager_regional_secret_version_access` datasources ([#19831](https://github.com/hashicorp/terraform-provider-google/pull/19831))
* spanner: added `kms_key_names` to `encryption_config` in `google_spanner_database` ([#19846](https://github.com/hashicorp/terraform-provider-google/pull/19846))
* workstations: added `max_usable_workstations` field to `google_workstations_workstation_config` resource ([#19872](https://github.com/hashicorp/terraform-provider-google/pull/19872))
* workstations: added field `allowed_ports` to `google_workstations_workstation_config` ([#19845](https://github.com/hashicorp/terraform-provider-google/pull/19845))

BUG FIXES:
* bigquery: fixed a regression that caused `google_bigquery_dataset_iam_*` resources to attempt to set deleted IAM members, thereby triggering an API error ([#19857](https://github.com/hashicorp/terraform-provider-google/pull/19857))
* compute: fixed an issue in `google_compute_backend_service` and `google_compute_region_backend_service` to allow sending `false` for `iap.enabled` ([#19795](https://github.com/hashicorp/terraform-provider-google/pull/19795))
* container: `node_config.linux_node_config`, `node_config.workload_metadata_config` and `node_config.kubelet_config` will now successfully send empty messages to the API when `terraform plan` indicates they are being removed, rather than null, which caused an error. The sole reliable case is `node_config.linux_node_config` when the block is removed, where there will still be a permadiff, but the update request that's triggered will no longer error and other changes displayed in the plan should go through. ([#19842](https://github.com/hashicorp/terraform-provider-google/pull/19842))

## 6.7.0 (October 14, 2024)

FEATURES:
* **New Resource:** `google_healthcare_pipeline_job` ([#19717](https://github.com/hashicorp/terraform-provider-google/pull/19717))
* **New Resource:** `google_secure_source_manager_branch_rule` ([#19773](https://github.com/hashicorp/terraform-provider-google/pull/19773))

IMPROVEMENTS:
* container: `google_container_cluster` will now accept server-specified values for `node_pool_auto_config.0.node_kubelet_config` when it is not defined in configuration and will not detect drift. Note that this means that removing the value from configuration will now preserve old settings instead of reverting the old settings. ([#19817](https://github.com/hashicorp/terraform-provider-google/pull/19817))
* discoveryengine: added `chat_engine_config.dialogflow_agent_to_link` field to `google_discovery_engine_chat_engine` resource ([#19723](https://github.com/hashicorp/terraform-provider-google/pull/19723))
* networkconnectivity: added field `migration` to resource `google_network_connectivity_internal_range` ([#19757](https://github.com/hashicorp/terraform-provider-google/pull/19757))
* networkservices: added `routing_mode` field to `google_network_services_gateway` resource ([#19764](https://github.com/hashicorp/terraform-provider-google/pull/19764))

BUG FIXES:
* bigtable: fixed an error where BigTable IAM resources could be created with conditions but the condition was not stored in state ([#19725](https://github.com/hashicorp/terraform-provider-google/pull/19725))
* container: fixed issue which caused to not being able to disable `enable_cilium_clusterwide_network_policy` field on `google_container_cluster`. ([#19736](https://github.com/hashicorp/terraform-provider-google/pull/19736))
* container: fixed a diff triggered by a new API-side default value for `node_config.0.kubelet_config.0.insecure_kubelet_readonly_port_enabled`. Terraform will now accept server-specified values for `node_config.0.kubelet_config` when it is not defined in configuration and will not detect drift. Note that this means that removing the value from configuration will now preserve old settings instead of reverting the old settings. ([#19817](https://github.com/hashicorp/terraform-provider-google/pull/19817))
* dataproc: fixed a bug in `google_dataproc_cluster` that prevented creation of clusters with `internal_ip_only` set to false ([#19782](https://github.com/hashicorp/terraform-provider-google/pull/19782))
* iam: addressed `google_service_account` creation issues caused by the eventual consistency of the GCP IAM API by ignoring 403 errors returned on polling the service account after creation. ([#19727](https://github.com/hashicorp/terraform-provider-google/pull/19727))
* logging: fixed the whitespace permadiff on `exclusions.filter` field in `google_logging_billing_account_sink`, `google_logging_folder_sink`, `google_logging_organization_sink` and `google_logging_project_sink` resources ([#19744](https://github.com/hashicorp/terraform-provider-google/pull/19744))
* pubsub: fixed permadiff with configuring an empty `retry_policy` in `google_pubsub_subscription`.  This will result in `minimum_backoff` and `maximum_backoff` using server-side defaults. To use "immedate retry", do not specify a `retry_policy` block at all. ([#19784](https://github.com/hashicorp/terraform-provider-google/pull/19784))
* secretmanager: fixed the issue of unpopulated fields `labels`, `annotations` and `version_destroy_ttl` in the terraform state for the `google_secret_manager_secrets` datasource ([#19748](https://github.com/hashicorp/terraform-provider-google/pull/19748))

## 6.6.0 (October 7, 2024)

FEATURES:
* **New Resource:** `google_dataproc_batch` ([#19686](https://github.com/hashicorp/terraform-provider-google/pull/19686))
* **New Resource:** `google_healthcare_pipeline_job` ([#19717](https://github.com/hashicorp/terraform-provider-google/pull/19717))
* **New Resource:** `google_site_verification_owner` ([#19641](https://github.com/hashicorp/terraform-provider-google/pull/19641))

IMPROVEMENTS:
* assuredworkloads: added `HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS` and `HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS_WITH_US_SUPPORT` enum values to `compliance_regime` in the `google_assuredworkload_workload` resource ([#19714](https://github.com/hashicorp/terraform-provider-google/pull/19714))
* compute: added `bgp_best_path_selection_mode `,`bgp_bps_always_compare_med` and `bgp_bps_inter_region_cost ` fields to `google_compute_network` resource ([#19708](https://github.com/hashicorp/terraform-provider-google/pull/19708))
* compute: added `next_hop_origin `,`next_hop_med ` and `next_hop_inter_region_cost ` output fields to `google_compute_route` resource ([#19708](https://github.com/hashicorp/terraform-provider-google/pull/19708))
* compute: added enum `STATEFUL_COOKIE_AFFINITY` and `strong_session_affinity_cookie` field to `google_compute_backend_service` and `google_compute_region_backend_service` resource ([#19665](https://github.com/hashicorp/terraform-provider-google/pull/19665))
* compute: moved `TDX` instance option for `confidential_instance_type` in `google_compute_instance` from Beta to GA ([#19706](https://github.com/hashicorp/terraform-provider-google/pull/19706))
* containeraws: added `kubelet_config` field group to the `google_container_aws_node_pool` resource ([#19714](https://github.com/hashicorp/terraform-provider-google/pull/19714))
* pubsub: added GCS ingestion settings and platform log settings to `google_pubsub_topic` resource ([#19669](https://github.com/hashicorp/terraform-provider-google/pull/19669))
* sourcerepo: added `create_ignore_already_exists` field to `google_sourcerepo_repository` resource ([#19716](https://github.com/hashicorp/terraform-provider-google/pull/19716))
* sql: added in-place update support for `settings.time_zone` in `google_sql_database_instance` resource ([#19654](https://github.com/hashicorp/terraform-provider-google/pull/19654))
* tags: increased maximum accepted input length for the `short_name` field in `google_tags_tag_key` and `google_tags_tag_value` resources ([#19712](https://github.com/hashicorp/terraform-provider-google/pull/19712))

BUG FIXES:
* bigquery: fixed `google_bigquery_dataset_iam_member` to be able to delete itself and overwrite the existing iam members for bigquery dataset keeping the authorized datasets as they are. ([#19682](https://github.com/hashicorp/terraform-provider-google/pull/19682))
* bigquery: fixed an error which could occur with service account field values containing non-lower-case characters in `google_bigquery_dataset_access` ([#19705](https://github.com/hashicorp/terraform-provider-google/pull/19705))
* compute: fixed an issue where the `boot_disk.initialize_params.resource_policies` field in `google_compute_instance` forced a resource recreation when used in combination with `google_compute_disk_resource_policy_attachment` ([#19692](https://github.com/hashicorp/terraform-provider-google/pull/19692))
* compute: fixed the issue that `labels` is not set when creating the resource `google_compute_interconnect` ([#19632](https://github.com/hashicorp/terraform-provider-google/pull/19632))
* tags:  removed `google_tags_location_tag_binding` resource from the Terraform state when its parent resource has been removed outside of Terraform ([#19693](https://github.com/hashicorp/terraform-provider-google/pull/19693))
* workbench: fixed a bug in the `google_workbench_instance` resource where the removal of `labels` was not functioning as expected. ([#19620](https://github.com/hashicorp/terraform-provider-google/pull/19620))

## 6.5.0 (September 30, 2024)
DEPRECATIONS:
* compute: deprecated `macsec.pre_shared_keys.fail_open` field in `google_compute_interconnect` resource. Use the new `macsec.fail_open` field instead ([#19572](https://github.com/hashicorp/terraform-provider-google/pull/19572))

FEATURES:
* **New Data Source:** `google_compute_region_instance_group_manager` ([#19589](https://github.com/hashicorp/terraform-provider-google/pull/19589))
* **New Data Source:** `google_privileged_access_manager_entitlement` ([#19580](https://github.com/hashicorp/terraform-provider-google/pull/19580))
* **New Data Source:** `google_secret_manager_regional_secret_version_access` ([#19538](https://github.com/hashicorp/terraform-provider-google/pull/19538))
* **New Data Source:** `google_secret_manager_regional_secret_version` ([#19514](https://github.com/hashicorp/terraform-provider-google/pull/19514))
* **New Data Source:** `google_secret_manager_regional_secrets` ([#19532](https://github.com/hashicorp/terraform-provider-google/pull/19532))
* **New Resource:** `google_compute_router_nat_address` ([#19550](https://github.com/hashicorp/terraform-provider-google/pull/19550))
* **New Resource:** `google_logging_log_scope` ([#19559](https://github.com/hashicorp/terraform-provider-google/pull/19559))

IMPROVEMENTS:
* apigee: added `activate` field to `google_apigee_nat_address` resource ([#19591](https://github.com/hashicorp/terraform-provider-google/pull/19591))
* bigquery: added `biglake_configuration` field to `google_bigquery_table` resource to support BigLake Managed Tables ([#19541](https://github.com/hashicorp/terraform-provider-google/pull/19541))
* cloudrunv2: promoted `scaling` field in `google_cloud_run_v2_service` resource to GA ([#19588](https://github.com/hashicorp/terraform-provider-google/pull/19588))
* composer: promoted `config.workloads_config.cloud_data_lineage_integration` field in `google_composer_environment` resource to GA ([#19612](https://github.com/hashicorp/terraform-provider-google/pull/19612))
* compute: added `existing_reservations` field to `google_compute_region_commitment` resource ([#19585](https://github.com/hashicorp/terraform-provider-google/pull/19585))
* compute: added `hostname` field to `google_compute_instance` data source ([#19607](https://github.com/hashicorp/terraform-provider-google/pull/19607))
* compute: added `initial_nat_ip` field to `google_compute_router_nat` resource ([#19550](https://github.com/hashicorp/terraform-provider-google/pull/19550))
* compute: added `macsec.fail_open` field to `google_compute_interconnect` resource ([#19572](https://github.com/hashicorp/terraform-provider-google/pull/19572))
* compute: added `SUSPENDED` as a possible value to `desired_state` field in `google_compute_instance` resource ([#19586](https://github.com/hashicorp/terraform-provider-google/pull/19586))
* compute: added import support for `projects/{{project}}/meta-data/{{key}}` format for `google_compute_project_metadata_item` resource ([#19613](https://github.com/hashicorp/terraform-provider-google/pull/19613))
* compute: marked `customer_name` and `location` fields as optional in `google_compute_interconnect` resource to support cross cloud interconnect ([#19619](https://github.com/hashicorp/terraform-provider-google/pull/19619))
* container: added `linux_node_config.hugepages_config` field to `google_container_node_pool` resource ([#19521](https://github.com/hashicorp/terraform-provider-google/pull/19521))
* container: promoted `gcfs_config` field in `google_container_cluster` resource to GA ([#19617](https://github.com/hashicorp/terraform-provider-google/pull/19617))
* looker: added `psc_enabled` and `psc_config` fields to `google_looker_instance` resource ([#19523](https://github.com/hashicorp/terraform-provider-google/pull/19523))
* networkconnectivity: added `include_import_ranges` field to `google_network_connectivity_spoke` resource for `linked_vpn_tunnels`, `linked_interconnect_attachments` and `linked_router_appliance_instances` ([#19530](https://github.com/hashicorp/terraform-provider-google/pull/19530))
* secretmanagerregional: added `version_aliases` field to `google_secret_manager_regional_secret` resource ([#19514](https://github.com/hashicorp/terraform-provider-google/pull/19514))
* workbench: increased create timeout to 20 minutes for `google_workbench_instance` resource ([#19551](https://github.com/hashicorp/terraform-provider-google/pull/19551))

BUG FIXES:
* bigquery: fixed in-place update of `google_bigquery_table` resource when `external_data_configuration.schema` field is set ([#19558](https://github.com/hashicorp/terraform-provider-google/pull/19558))
* bigquerydatapolicy: fixed permadiff on `policy_tag` field in `google_bigquery_datapolicy_data_policy` resource ([#19563](https://github.com/hashicorp/terraform-provider-google/pull/19563))
* composer: fixed `storage_config.bucket` field to support a bucket name with or without "gs://" prefix ([#19552](https://github.com/hashicorp/terraform-provider-google/pull/19552))
* container: added support for setting `addons_config.gcp_filestore_csi_driver_config` and `enable_autopilot` in the same `google_container_cluster` ([#19590](https://github.com/hashicorp/terraform-provider-google/pull/19590))
* container: fixed `node_config.kubelet_config` updates in `google_container_cluster` resource ([#19562](https://github.com/hashicorp/terraform-provider-google/pull/19562))
* container: fixed a bug where specifying `node_pool_defaults.node_config_defaults` with `enable_autopilot = true` would cause `google_container_cluster` resource creation failure ([#19543](https://github.com/hashicorp/terraform-provider-google/pull/19543))
* workbench: fixed a bug in the `google_workbench_instance` resource where the removal of `labels` was not functioning as expected ([#19620](https://github.com/hashicorp/terraform-provider-google/pull/19620))


## 6.4.0 (September 23, 2024)

DEPRECATIONS:
* securitycenterv2: deprecated `google_scc_v2_organization_scc_big_query_exports`. Use `google_scc_v2_organization_scc_big_query_export` instead. ([#19457](https://github.com/hashicorp/terraform-provider-google/pull/19457))

FEATURES:
* **New Data Source:** `google_secret_manager_regional_secret_version` ([#19514](https://github.com/hashicorp/terraform-provider-google/pull/19514))
* **New Data Source:** `google_secret_manager_regional_secret` ([#19491](https://github.com/hashicorp/terraform-provider-google/pull/19491))
* **New Resource:** `google_database_migration_service_migration_job` ([#19488](https://github.com/hashicorp/terraform-provider-google/pull/19488))
* **New Resource:** `google_discovery_engine_target_site` ([#19469](https://github.com/hashicorp/terraform-provider-google/pull/19469))
* **New Resource:** `google_healthcare_workspace` ([#19476](https://github.com/hashicorp/terraform-provider-google/pull/19476))
* **New Resource:** `google_scc_folder_scc_big_query_export` ([#19480](https://github.com/hashicorp/terraform-provider-google/pull/19480))
* **New Resource:** `google_scc_organization_scc_big_query_export` ([#19465](https://github.com/hashicorp/terraform-provider-google/pull/19465))
* **New Resource:** `google_scc_project_scc_big_query_export` ([#19466](https://github.com/hashicorp/terraform-provider-google/pull/19466))
* **New Resource:** `google_scc_v2_organization_scc_big_query_export` ([#19457](https://github.com/hashicorp/terraform-provider-google/pull/19457))
* **New Resource:** `google_secret_manager_regional_secret_version` ([#19504](https://github.com/hashicorp/terraform-provider-google/pull/19504))
* **New Resource:** `google_secret_manager_regional_secret` ([#19461](https://github.com/hashicorp/terraform-provider-google/pull/19461))
* **New Resource:** `google_site_verification_web_resource` ([#19477](https://github.com/hashicorp/terraform-provider-google/pull/19477))
* **New Resource:** `google_spanner_backup_schedule` ([#19449](https://github.com/hashicorp/terraform-provider-google/pull/19449))

IMPROVEMENTS:
* alloydb: added `enable_outbound_public_ip` field to `google_alloydb_instance` resource ([#19444](https://github.com/hashicorp/terraform-provider-google/pull/19444))
* apigee: added in-place update for `consumer_accept_list` field in `google_apigee_instance` resource ([#19442](https://github.com/hashicorp/terraform-provider-google/pull/19442))
* compute: added `interface` field to `google_compute_attached_disk` resource ([#19440](https://github.com/hashicorp/terraform-provider-google/pull/19440))
* compute: added in-place update in `google_compute_interconnect` resource, except for `remote_location` and `requested_features` fields ([#19508](https://github.com/hashicorp/terraform-provider-google/pull/19508))
* filestore: added `deletion_protection_enabled` and `deletion_protection_reason` fields to `google_filestore_instance` resource ([#19446](https://github.com/hashicorp/terraform-provider-google/pull/19446))
* looker: added `fips_enabled` field to `google_looker_instance` resource ([#19511](https://github.com/hashicorp/terraform-provider-google/pull/19511))
* metastore: added `deletion_protection` field to `google_dataproc_metastore_service` resource ([#19505](https://github.com/hashicorp/terraform-provider-google/pull/19505))
* netapp: added `allow_auto_tiering` field to `google_netapp_storage_pool` resource ([#19454](https://github.com/hashicorp/terraform-provider-google/pull/19454))
* netapp: added `tiering_policy` field to `google_netapp_volume` resource ([#19454](https://github.com/hashicorp/terraform-provider-google/pull/19454))
* secretmanagerregional: added `version_aliases` field to `google_secret_manager_regional_secret` resource ([#19514](https://github.com/hashicorp/terraform-provider-google/pull/19514))
* spanner: added `edition` field to `google_spanner_instance` resource ([#19449](https://github.com/hashicorp/terraform-provider-google/pull/19449))

BUG FIXES:
* compute: fixed a permadiff on `iap` field in `google_compute_backend` and `google_compute_region_backend` resources ([#19509](https://github.com/hashicorp/terraform-provider-google/pull/19509))
* container: fixed a bug where specifying `node_pool_defaults.node_config_defaults` with `enable_autopilot = true` will cause `google_container_cluster` resource creation failure ([#19543](https://github.com/hashicorp/terraform-provider-google/pull/19543))
* container: fixed a permadiff on `node_config.gcfs_config` field in `google_container_cluster` and `google_container_node_pool` resources ([#19512](https://github.com/hashicorp/terraform-provider-google/pull/19512))
* container: fixed the in-place update for `node_config.gcfs_config` field in `google_container_cluster` and `google_container_node_pool` resources ([#19512](https://github.com/hashicorp/terraform-provider-google/pull/19512))
* container: made `node_config.kubelet_config.cpu_manager_policy` field optional to fix its update in `google_container_cluster` resource ([#19464](https://github.com/hashicorp/terraform-provider-google/pull/19464))
* dns: fixed a permadiff on `dnssec_config` field in `google_dns_managed_zone` resource ([#19456](https://github.com/hashicorp/terraform-provider-google/pull/19456))
* pubsub: allowed `filter` field to contain line breaks in `google_pubsub_subscription` resource ([#19451](https://github.com/hashicorp/terraform-provider-google/pull/19451))

## 6.3.0 (September 16, 2024)

FEATURES:
* **New Data Source:** `google_bigquery_tables` ([#19402](https://github.com/hashicorp/terraform-provider-google/pull/19402))
* **New Resource:** `google_developer_connect_connection` ([#19431](https://github.com/hashicorp/terraform-provider-google/pull/19431))
* **New Resource:** `google_developer_connect_git_repository_link` ([#19431](https://github.com/hashicorp/terraform-provider-google/pull/19431))
* **New Resource:** `google_memorystore_instance` ([#19398](https://github.com/hashicorp/terraform-provider-google/pull/19398))

IMPROVEMENTS:
* compute: added `connected_endpoints.consumer_network` and `connected_endpoints.psc_connection_id` fields to `google_compute_service_attachment` resource ([#19426](https://github.com/hashicorp/terraform-provider-google/pull/19426))
* compute: added field `http_keep_alive_timeout_sec` to `google_region_compute_target_https_proxy` and `google_region_compute_target_http_proxy` resources ([#19432](https://github.com/hashicorp/terraform-provider-google/pull/19432))
* compute: added support for `boot_disk.initialize_params.resource_policies` in `google_compute_instance` and `google_instance_template` ([#19407](https://github.com/hashicorp/terraform-provider-google/pull/19407))
* container: added `storage_pools` to `node_config` in `google_container_cluster` and `google_container_node_pool` ([#19423](https://github.com/hashicorp/terraform-provider-google/pull/19423))
* containerattached: added `security_posture_config` field to `google_container_attached_cluster` resource ([#19411](https://github.com/hashicorp/terraform-provider-google/pull/19411))
* netapp: added `large_capacity` and `multiple_endpoints` to `google_netapp_volume` resource ([#19384](https://github.com/hashicorp/terraform-provider-google/pull/19384))
* resourcemanager: added `tags` field to `google_folder` to allow setting tags for folders at creation time ([#19380](https://github.com/hashicorp/terraform-provider-google/pull/19380))

BUG FIXES:
* compute: setting `network_ip` to "" will no longer cause diff and will be treated the same as `null` ([#19400](https://github.com/hashicorp/terraform-provider-google/pull/19400))
* dataproc: updated `google_dataproc_cluster` to protect against handling nil `kerberos_config` values ([#19401](https://github.com/hashicorp/terraform-provider-google/pull/19401))
* dns: added a mutex to `google_dns_record_set` to prevent conflicts when multiple resources attempt to operate on the same record set ([#19416](https://github.com/hashicorp/terraform-provider-google/pull/19416))
* managedkafka: added 5 second wait post `google_managed_kafka_topic` creation to fix eventual consistency errors ([#19429](https://github.com/hashicorp/terraform-provider-google/pull/19429))

## 6.2.0 (September 9, 2024)

FEATURES:
* **New Data Source:** `google_certificate_manager_certificates` ([#19361](https://github.com/hashicorp/terraform-provider-google/pull/19361))
* **New Resource:** `google_network_security_server_tls_policy` ([#19314](https://github.com/hashicorp/terraform-provider-google/pull/19314))
* **New Resource:** `google_scc_v2_folder_scc_big_query_export` ([#19327](https://github.com/hashicorp/terraform-provider-google/pull/19327))
* **New Resource:** `google_scc_v2_project_scc_big_query_export` ([#19311](https://github.com/hashicorp/terraform-provider-google/pull/19311))

IMPROVEMENTS:
* assuredworkload: added field `partner_service_billing_account` to `google_assured_workloads_workload` ([#19358](https://github.com/hashicorp/terraform-provider-google/pull/19358))
* bigtable: added support for `column_family.type` in `google_bigtable_table` ([#19302](https://github.com/hashicorp/terraform-provider-google/pull/19302))
* cloudrun: promoted support for nfs and csi volumes (for Cloud Storage FUSE) for `google_cloud_run_service` to GA ([#19359](https://github.com/hashicorp/terraform-provider-google/pull/19359))
* cloudrunv2: promoted support for nfs and gcs volumes for `google_cloud_run_v2_job` to GA ([#19359](https://github.com/hashicorp/terraform-provider-google/pull/19359))
* compute: added `boot_disk.interface` field to `google_compute_instance` resource ([#19319](https://github.com/hashicorp/terraform-provider-google/pull/19319))
* container: added `node_pool_auto_config.node_kublet_config.insecure_kubelet_readonly_port_enabled` field to `google_container_cluster`. ([#19320](https://github.com/hashicorp/terraform-provider-google/pull/19320))
* container: added `insecure_kubelet_readonly_port_enabled` to `node_pool.node_config.kubelet_config` and `node_config.kubelet_config` in `google_container_node_pool` resource. ([#19312](https://github.com/hashicorp/terraform-provider-google/pull/19312))
* container: added `insecure_kubelet_readonly_port_enabled` to `node_pool_defaults.node_config_defaults`, `node_pool.node_config.kubelet_config`, and `node_config.kubelet_config` in `google_container_cluster` resource. ([#19312](https://github.com/hashicorp/terraform-provider-google/pull/19312))
* container: added support for in-place updates for `google_compute_node_pool.node_config.gcfs_config` and `google_container_cluster.node_config.gcfs_cluster` and `google_container_cluster.node_pool.node_config.gcfs_cluster` ([#19365](https://github.com/hashicorp/terraform-provider-google/pull/19365))
* container: promoted the `additive_vpc_scope_dns_domain` field on the `google_container_cluster` resource to GA ([#19313](https://github.com/hashicorp/terraform-provider-google/pull/19313))
* iambeta: added `x509` field to `google_iam_workload_identity_pool_provider ` resource ([#19375](https://github.com/hashicorp/terraform-provider-google/pull/19375))
* networkconnectivity: added `include_export_ranges` to `google_network_connectivity_spoke` ([#19346](https://github.com/hashicorp/terraform-provider-google/pull/19346))
* pubsub: added `cloud_storage_config.max_messages` and `cloud_storage_config.avro_config.use_topic_schema` fields to `google_pubsub_subscription` resource ([#19338](https://github.com/hashicorp/terraform-provider-google/pull/19338))
* redis: added the `maintenance_policy` field to the `google_redis_cluster` resource ([#19341](https://github.com/hashicorp/terraform-provider-google/pull/19341))
* resourcemanager: added `tags` field to `google_project` to allow setting tags for projects at creation time ([#19351](https://github.com/hashicorp/terraform-provider-google/pull/19351))
* securitycenter: added support for empty `streaming_config.filter` values in `google_scc_notification_config` resources ([#19369](https://github.com/hashicorp/terraform-provider-google/pull/19369))

BUG FIXES:
* compute: fixed `google_compute_interconnect` to support correct `available_features` option of `IF_MACSEC` ([#19330](https://github.com/hashicorp/terraform-provider-google/pull/19330))
* compute: fixed a bug where `advertised_route_priority` was accidentally set to 0 during updates in `google_compute_router_peer` ([#19366](https://github.com/hashicorp/terraform-provider-google/pull/19366))
* compute: fixed a permadiff caused by setting `start_time` in an incorrect H:mm format in `google_compute_resource_policies` resources ([#19297](https://github.com/hashicorp/terraform-provider-google/pull/19297))
* compute: fixed `network_interface.subnetwork_project` validation to match with the project in `network_interface.subnetwork` field when `network_interface.subnetwork` has full self_link in `google_compute_instance` resource ([#19348](https://github.com/hashicorp/terraform-provider-google/pull/19348))
* container: removed unnecessary force replacement in node pool `gcfs_config` ([#19365](https://github.com/hashicorp/terraform-provider-google/pull/19365)
* kms: updated the `google_kms_autokey_config` resource's `folder` field to accept values that are either full resource names (`folders/{folder_id}`) or just the folder id (`{folder_id}` only) ([#19364](https://github.com/hashicorp/terraform-provider-google/pull/19364)))
* storage: added retry support for 429 errors in `google_storage_bucket` resource ([#19353](https://github.com/hashicorp/terraform-provider-google/pull/19353))


## 6.1.0 (September 4, 2024)

FEATURES:
* **New Data Source:** `google_kms_crypto_key_latest_version` ([#19249](https://github.com/hashicorp/terraform-provider-google/pull/19249))
* **New Data Source:** `google_kms_crypto_key_versions` ([#19241](https://github.com/hashicorp/terraform-provider-google/pull/19241))

IMPROVEMENTS:
* databasemigrationservice: added support in `google_database_migration_service_connection_profile` for creating DMS connection profiles that link to existing Cloud SQL instances/AlloyDB clusters. ([#19291](https://github.com/hashicorp/terraform-provider-google/pull/19291))
* alloydb: added `subscription_type` and `trial_metadata` field to `google_alloydb_cluster` resource ([#19262](https://github.com/hashicorp/terraform-provider-google/pull/19262))
* bigquery: added `encryption_configuration` field to `google_bigquery_data_transfer_config` resource ([#19267](https://github.com/hashicorp/terraform-provider-google/pull/19267))
* bigqueryanalyticshub: added `selected_resources`, and `restrict_direct_table_access` to `google_bigquery_analytics_hub_listing` resource ([#19244](https://github.com/hashicorp/terraform-provider-google/pull/19244))
* bigqueryanalyticshub: added `sharing_environment_config` to `google_bigquery_analytics_hub_data_exchange` resource ([#19244](https://github.com/hashicorp/terraform-provider-google/pull/19244))
* cloudtasks: added `http_target` field to `google_cloud_tasks_queue` resource ([#19253](https://github.com/hashicorp/terraform-provider-google/pull/19253))
* compute: added `accelerators` field to `google_compute_node_template` resource ([#19292](https://github.com/hashicorp/terraform-provider-google/pull/19292))
* compute: allowed disabling `server_tls_policy` during update in `google_compute_target_https_proxy` resources ([#19233](https://github.com/hashicorp/terraform-provider-google/pull/19233))
* container: added `secret_manager_config` field to `google_container_cluster` resource ([#19288](https://github.com/hashicorp/terraform-provider-google/pull/19288))
* datastream: added `transaction_logs` and `change_tables` to the `datastream_stream` resource ([#19248](https://github.com/hashicorp/terraform-provider-google/pull/19248))
* discoveryengine: added `chunking_config` and `layout_parsing_config` fields to `google_discovery_engine_data_store` resource ([#19274](https://github.com/hashicorp/terraform-provider-google/pull/19274))
* dlp: added `inspect_template_modified_cadence` field to `big_query_target` and `cloud_sql_target` in `google_data_loss_prevention_discovery_config` resource ([#19282](https://github.com/hashicorp/terraform-provider-google/pull/19282))
* dlp: added `tag_resources` field to `google_data_loss_prevention_discovery_config` resource ([#19282](https://github.com/hashicorp/terraform-provider-google/pull/19282))
* networksecurity: promoted `google_network_security_client_tls_policy` to GA ([#19293](https://github.com/hashicorp/terraform-provider-google/pull/19293))

BUG FIXES:
* bigquery: fixed an error which could occur with email field values containing non-lower-case characters in `google_bigquery_dataset_access` resource ([#19259](https://github.com/hashicorp/terraform-provider-google/pull/19259))
* bigqueryanalyticshub: made `bigquery_dataset` immutable in `google_bigquery_analytics_hub_listing` as it was not updatable in the API. Now modifying the field in Terraform will correctly recreate the resource rather than causing Terraform to report it would attempt an invalid update. ([#19244](https://github.com/hashicorp/terraform-provider-google/pull/19244))
* container: fixed update inconsistency in `google_container_cluster` resource ([#19247](https://github.com/hashicorp/terraform-provider-google/pull/19247))
* pubsub: fixed a validation bug that didn't allow empty filter definitions for `google_pubsub_subscription` resources ([#19284](https://github.com/hashicorp/terraform-provider-google/pull/19284))
* resourcemanager: fixed a bug where data.google_client_config failed silently when inadequate credentials were used to configure the provider ([#19286](https://github.com/hashicorp/terraform-provider-google/pull/19286))
* sql: fixed importing `google_sql_user` where `host` is an IPv4 CIDR ([#19243](https://github.com/hashicorp/terraform-provider-google/pull/19243))
* sql: fixed overwriting of `name` field for IAM Group user in `google_sql_user` resource ([#19234](https://github.com/hashicorp/terraform-provider-google/pull/19234))

## 6.0.1 (August 26, 2024)

BREAKING CHANGES:

* sql: removed `settings.ip_configuration.require_ssl` from `google_sql_database_instance` in favor of `settings.ip_configuration.ssl_mode`. This field was intended to be removed in 6.0.0. ([#19263](https://github.com/hashicorp/terraform-provider-google/pull/19263))

## 6.0.0 (August 26, 2024)

[Terraform Google Provider 6.0.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_6_upgrade)

BREAKING CHANGES:
* provider: changed provider labels to add the `goog-terraform-provisioned: true` label by default. ([#19190](https://github.com/hashicorp/terraform-provider-google/pull/19190))
* activedirectory: added `deletion_protection` field to `google_active_directory_domain` resource. This field defaults to `true`, preventing accidental deletions. To delete the resource, you must first set `deletion_protection = false` before destroying the resource. ([#18906](https://github.com/hashicorp/terraform-provider-google/pull/18906))
* alloydb: removed `network` in `google_alloy_db_cluster`. Use `network_config.network` instead. ([#19181](https://github.com/hashicorp/terraform-provider-google/pull/19181))
* bigquery: added client-side validation to prevent table view creation if schema contains required fields for `google_bigquery_table` resource ([#18767](https://github.com/hashicorp/terraform-provider-google/pull/18767))
* bigquery: removed `allow_resource_tags_on_deletion` from `google_bigquery_table`. Resource tags are now always allowed on table deletion. ([#19077](https://github.com/hashicorp/terraform-provider-google/pull/19077))
* bigqueryreservation: removed `multi_region_auxiliary` from `google_bigquery_reservation` ([#18922](https://github.com/hashicorp/terraform-provider-google/pull/18922))
* billing: revised the format of `id` for `google_billing_project_info` ([#18823](https://github.com/hashicorp/terraform-provider-google/pull/18823))
* cloudrunv2: added `deletion_protection` field to `google_cloudrunv2_service`.  This field defaults to `true`, preventing accidental deletions. To delete the resource, you must first set `deletion_protection = false` before destroying the resource.([#19019](https://github.com/hashicorp/terraform-provider-google/pull/19019))
* cloudrunv2: changed `liveness_probe` to no longer infer a default value from api on `google_cloud_run_v2_service`. Removing this field and applying the change will now  remove liveness probe from the Cloud Run service. ([#18764](https://github.com/hashicorp/terraform-provider-google/pull/18764))
* cloudrunv2: retyped `containers.env` to SET from ARRAY for `google_cloud_run_v2_service` and `google_cloud_run_v2_job`. ([#18855](https://github.com/hashicorp/terraform-provider-google/pull/18855))
* composer: `ip_allocation_policy = []` in `google_composer_environment` is no longer valid configuration. Removing the field from configuration should not produce a diff. ([#19207](https://github.com/hashicorp/terraform-provider-google/pull/19207))
* compute: added new required field `enabled` in `google_compute_backend_service` and `google_compute_region_backend_service` ([#18772](https://github.com/hashicorp/terraform-provider-google/pull/18772))
* compute: changed `certifcate_id` in `google_compute_managed_ssl_certificate` to correctly be output only. ([#19069](https://github.com/hashicorp/terraform-provider-google/pull/19069))
* compute: revised and in some cases removed default values  of `connection_draining_timeout_sec`, `balancing_mode` and `outlier_detection` in `google_compute_region_backend_service` and `google_compute_backend_service`. ([#18720](https://github.com/hashicorp/terraform-provider-google/pull/18720))
* compute: revised the format of `id`  for `compute_network_endpoints` ([#18844](https://github.com/hashicorp/terraform-provider-google/pull/18844))
* compute: `guest_accelerator = []` is no longer valid configuration in `google_compute_instance`. To explicitly set an empty list of objects, set guest_accelerator.count = 0. ([#19207](https://github.com/hashicorp/terraform-provider-google/pull/19207))
* compute: `google_compute_instance_from_template` and `google_compute_instance_from_machine_image` `network_interface.alias_ip_range, network_interface.access_config, attached_disk, guest_accelerator, service_account, scratch_disk` can no longer be set to an empty block `[]`. Removing the fields from configuration should not produce a diff. ([#19207](https://github.com/hashicorp/terraform-provider-google/pull/19207))
* compute: `secondary_ip_ranges = []` in `google_compute_subnetwork` is no longer valid configuration. To set an explicitly empty list, use `send_secondary_ip_range_if_empty` and completely remove `secondary_ip_range` from config.  ([#19207](https://github.com/hashicorp/terraform-provider-google/pull/19207))
* container: made `advanced_datapath_observability_config.enable_relay` required in `google_container_cluster` ([#19060](https://github.com/hashicorp/terraform-provider-google/pull/19060))
* container: removed deprecated field `advanced_datapath_observability_config.relay_mode` from `google_container_cluster` resource. Users are expected to use `enable_relay` field instead. ([#19060](https://github.com/hashicorp/terraform-provider-google/pull/19060))
* container: three label-related fields are now in `google_container_cluster` resource. `resource_labels` field is non-authoritative and only manages the labels defined by the users on the resource through Terraform. The new output-only `terraform_labels` field merges the labels defined by the users on the resource through Terraform and the default labels configured on the provider. The new output-only `effective_labels` field lists all of labels present on the resource in GCP, including the labels configured through Terraform, the system, and other clients. ([#19062](https://github.com/hashicorp/terraform-provider-google/pull/19062))
* container: made three fields `resource_labels`, `terraform_labels`, and `effective_labels` be present in `google_container_cluster` datasources. All three fields will have all of labels present on the resource in GCP including the labels configured through Terraform, the system, and other clients, equivalent to `effective_labels` on the resource. ([#19062](https://github.com/hashicorp/terraform-provider-google/pull/19062))
* container: `guest_accelerator = []` is no longer valid configuration in `google_container_cluster` and `google_container_node_pool`. To explicitly set an empty list of objects, set guest_accelerator.count = 0. ([#19207](https://github.com/hashicorp/terraform-provider-google/pull/19207))
* container: `guest_accelerator.gpu_driver_installation_config = []` and `guest_accelerator.gpu_sharing_config = []` are no longer valid configuration in `google_container_cluster` and `google_container_node_pool`. Removing the fields from configuration should not produce a diff. ([#19207](https://github.com/hashicorp/terraform-provider-google/pull/19207))
* datastore: removed `google_datastore_index` in favor of `google_firestore_index` ([#19160](https://github.com/hashicorp/terraform-provider-google/pull/19160))
* edgenetwork: three label-related fields are now in `google_edgenetwork_network ` and `google_edgenetwork_subnet` resources. `labels` field is non-authoritative and only manages the labels defined by the users on the resource through Terraform. The new output-only `terraform_labels` field merges the labels defined by the users on the resource through Terraform and the default labels configured on the provider. The new output-only `effective_labels` field lists all of labels present on the resource in GCP, including the labels configured through Terraform, the system, and other clients. ([#19062](https://github.com/hashicorp/terraform-provider-google/pull/19062))
* identityplatform: removed resource `google_identity_platform_project_default_config` in favor of `google_identity_platform_project_config` ([#18992](https://github.com/hashicorp/terraform-provider-google/pull/18992))
* pubsub: allowed `schema_settings` in `google_pubsub_topic` to be removed ([#18631](https://github.com/hashicorp/terraform-provider-google/pull/18631))
* integrations: removed `create_sample_workflows` and `provision_gmek` from `google_integrations_client` ([#19148](https://github.com/hashicorp/terraform-provider-google/pull/19148))
* redis: added a `deletion_protection_enabled` field to the `google_redis_cluster` resource.  This field defaults to `true`, preventing accidental deletions. To delete the resource, you must first set `deletion_protection_enabled = false` before destroying the resource. ([#19173](https://github.com/hashicorp/terraform-provider-google/pull/19173))
* resourcemanager: added `deletion_protection` field to `google_folder` to make deleting them require an explicit intent. Folder resources now cannot be destroyed unless `deletion_protection = false` is set for the resource. ([#19021](https://github.com/hashicorp/terraform-provider-google/pull/19021))
* resourcemanager: made `deletion_policy` in `google_project` 'PREVENT' by default. This makes deleting them require an explicit intent. `google_project` resources cannot be destroyed unless `deletion_policy` is set to 'ABANDON' or 'DELETE' for the resource. ([#19114](https://github.com/hashicorp/terraform-provider-google/pull/19114))
* sql: removed `settings.ip_configuration.require_ssl` in `google_sql_database_instance`. Please use `settings.ip_configuration.ssl_mode` instead. ([#18843](https://github.com/hashicorp/terraform-provider-google/pull/18843))
* storage: removed `no_age` field from  `lifecycle_rule.condition` in the `google_storage_bucket` resource ([#19048](https://github.com/hashicorp/terraform-provider-google/pull/19048))
* vpcaccess: removed default values for `min_throughput` and `min_instances` fields on `google_vpc_access_connector` and made them default to values returned from the API when not provided by users ([#18697](https://github.com/hashicorp/terraform-provider-google/pull/18697))
* vpcaccess: added a conflicting fields restriction between `min_throughput` and `min_instances` fields on `google_vpc_access_connector` ([#18697](https://github.com/hashicorp/terraform-provider-google/pull/18697))
* vpcaccess: added a conflicting fields restriction between `max_throughput` and `max_instances` fields on `google_vpc_access_connector` ([#18697](https://github.com/hashicorp/terraform-provider-google/pull/18697))
* workstation: defaulted `host.gce_instance.disable_ssh` to true for `google_workstations_workstation_config` ([#19101](https://github.com/hashicorp/terraform-provider-google/pull/19101))
IMPROVEMENTS:
* compute: added fields `reserved_internal_range` and `secondary_ip_ranges[].reserved_internal_range` to `google_compute_subnetwork` resource ([#19151](https://github.com/hashicorp/terraform-provider-google/pull/19151))
* compute: changed the behavior of `name_prefix` in multiple Compute resources to allow for a longer max length of 54 characters. See the upgrade guide and resource documentation for more details. ([#19152](https://github.com/hashicorp/terraform-provider-google/pull/19152))
BUG FIXES:
* compute: fixed an issue regarding sending `enabled` field by default for null `iap` message in `google_compute_backend_service` and `google_compute_region_backend_service` ([#18772](https://github.com/hashicorp/terraform-provider-google/pull/18772))
