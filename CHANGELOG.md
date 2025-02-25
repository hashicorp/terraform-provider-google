## 6.23.0 (Unreleased)

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

## 6.18.1 (January 29, 2025)

BUG FIXES:
* container: fixed a diff caused by server-side set values for `node_config.resource_labels` ([#21082](https://github.com/hashicorp/terraform-provider-google/pull/21082))

## 5.45.1 (January 29, 2025)

NOTES:
* 5.45.1 is a backport release, responding to a new GKE label being applied that can cause unwanted diffs in node pools. The changes in this release will be available in 6.18.1 and users upgrading to 6.X should upgrade to that version or higher.

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

## 5.45.0 (November 11, 2024)

NOTES:
* 5.45.0 is a backport release, responding to a new Spanner feature that may result in creation of unwanted backups for users. The changes in this release will be available in 6.11.0 and users upgrading to 6.X should upgrade to that version or higher.

IMPROVEMENTS:
* spanner: added `default_backup_schedule_type` field to  `google_spanner_instance` ([#20213](https://github.com/hashicorp/terraform-provider-google/pull/20213))

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

## 5.44.2 (October 14, 2024)

Notes:
* 5.44.2 is a backport release, responding to a GKE rollout that created permadiffs for many users. The changes in this release will be available in 6.7.0 and users upgrading to 6.X should upgrade to that version or higher.

IMPROVEMENTS:
* container: `google_container_cluster` will now accept server-specified values for `node_pool_auto_config.0.node_kubelet_config` when it is not defined in configuration and will not detect drift. Note that this means that removing the value from configuration will now preserve old settings instead of reverting the old settings. ([#19817](https://github.com/hashicorp/terraform-provider-google/pull/19817))

BUG FIXES:
* container: fixed a diff triggered by a new API-side default value for `node_config.0.kubelet_config.0.insecure_kubelet_readonly_port_enabled`. Terraform will now accept server-specified values for `node_config.0.kubelet_config` when it is not defined in configuration and will not detect drift. Note that this means that removing the value from configuration will now preserve old settings instead of reverting the old settings. ([#19817](https://github.com/hashicorp/terraform-provider-google/pull/19817))

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

## 4.85.0 (June 12, 2024)

NOTES:
* The `4.85.0` release backports configuration for the retention period for Cloud Storage soft delete (https://cloud.google.com/resources/storage/soft-delete-announce) so that customers who have not yet upgraded to `5.22.0`+ are able to configure the retention period of objects in their buckets. By upgrading to this version and configuring or otherwise interacting with the `google_storage_bucket.soft_delete_policy` values, you will need to upgrade directly to `5.22.0`+ from `4.85.0` when upgrading to `5.X` in the future.

IMPROVEMENTS:
* storage: added `soft_delete_policy` to `google_storage_bucket` resource ([#17624](https://github.com/hashicorp/terraform-provider-google/pull/17624))

## 4.84.0 (September 26, 2023)

DEPRECATIONS:
* alloydb: deprecated `network` field in favor of `network_config` on `google_alloydb_cluster`. ([#15881](https://github.com/hashicorp/terraform-provider-google/pull/15881))
* identityplayform: deprecated `google_identity_platform_project_default_config` resource. Use `google_identity_platform_config` resource instead ([#15876](https://github.com/hashicorp/terraform-provider-google/pull/15876))

FEATURES:
* **New Data Source:** `google_certificate_manager_certificate_map` ([#15906](https://github.com/hashicorp/terraform-provider-google/pull/15906))
* **New Resource:** `google_artifact_registry_vpcsc_config` ([#15840](https://github.com/hashicorp/terraform-provider-google/pull/15840))
* **New Resource:** `google_dialogflow_cx_security_settings` ([#15886](https://github.com/hashicorp/terraform-provider-google/pull/15886))
* **New Resource:** `google_gke_backup_restore_plan` ([#15858](https://github.com/hashicorp/terraform-provider-google/pull/15858))
* **New Resource:** `google_edgenetwork_network` ([#15891](https://github.com/hashicorp/terraform-provider-google/pull/15891))
* **New Resource:** `google_edgenetwork_subnet` ([#15891](https://github.com/hashicorp/terraform-provider-google/pull/15891))

IMPROVEMENTS:
* alloydb: added `network_config` field to support named IP ranges on `google_alloydb_cluster`. ([#15881](https://github.com/hashicorp/terraform-provider-google/pull/15881))
* cloudrunv2: added fields `network_interfaces` to resource `google_cloud_run_v2_job` to support Direct VPC egress. ([#15870](https://github.com/hashicorp/terraform-provider-google/pull/15870))
* cloudrunv2: added fields `network_interfaces` to resource `google_cloud_run_v2_service` to support Direct VPC egress. ([#15870](https://github.com/hashicorp/terraform-provider-google/pull/15870))
* compute: updated the `autoscaling_policy.mode` to accept `ONLY_SCALE_OUT` on `google_compute_autoscaler` ([#15890](https://github.com/hashicorp/terraform-provider-google/pull/15890))
* compute: added `server_tls_policy` argument to `google_compute_target_https_proxy` resource ([#15845](https://github.com/hashicorp/terraform-provider-google/pull/15845))
* compute: added `member` attribute to `google_compute_default_service_account` datasource ([#15897](https://github.com/hashicorp/terraform-provider-google/pull/15897))
* compute: added output field `internal_ipv6_prefix` to `google_compute_subnetwork` resource ([#15892](https://github.com/hashicorp/terraform-provider-google/pull/15892))
* container: added `node_config.fast_socket` field to `google_container_node_pool` ([#15872](https://github.com/hashicorp/terraform-provider-google/pull/15872))
* container: promoted `node_pool_auto_config` field in `google_container_cluster` from beta provider to GA provider. ([#15884](https://github.com/hashicorp/terraform-provider-google/pull/15884))
* container: promoted field `placement_policy.tpu_topology` in resource `google_container_node_pool` to GA ([#15869](https://github.com/hashicorp/terraform-provider-google/pull/15869))
* containeraws: added support for `auto_repair` in `google_container_aws_node_pool` ([#15862](https://github.com/hashicorp/terraform-provider-google/pull/15862))
* containerazure: added support for `auto_repair` in `google_container_azure_node_pool` ([#15862](https://github.com/hashicorp/terraform-provider-google/pull/15862))
* filestore: added support for the `"ZONAL"` value to `tier` in `google_filestore_instance` ([#15889](https://github.com/hashicorp/terraform-provider-google/pull/15889))
* firestore: added `delete_protection_state` field to `google_firestore_database` resource. ([#15878](https://github.com/hashicorp/terraform-provider-google/pull/15878))
* identityplatform: added `sign-in` field to `google_identity_platform_config` resource ([#15876](https://github.com/hashicorp/terraform-provider-google/pull/15876))
* networkconnectivity: added support for `linked_vpc_network` in `google_network_connectivity_spoke` ([#15862](https://github.com/hashicorp/terraform-provider-google/pull/15862))
* networkservices: increased default timeout for `google_network_services_edge_cache_origin` to 120m from 60m ([#15855](https://github.com/hashicorp/terraform-provider-google/pull/15855))
* networkservices: increased default timeout for `google_network_services_edge_cache_service` to 60m from 30m ([#15861](https://github.com/hashicorp/terraform-provider-google/pull/15861))
* secretmanager: added `is_secret_data_base64` field to `google_secret_manager_secret_version` resource ([#15853](https://github.com/hashicorp/terraform-provider-google/pull/15853))

BUG FIXES:
* bigquery: updated documentation for `google_bigquery_table.time_partitioning.expiration_ms` ([#15873](https://github.com/hashicorp/terraform-provider-google/pull/15873))
* bigtable: added a read timeout to `google_bigtable_instance` ([#15856](https://github.com/hashicorp/terraform-provider-google/pull/15856))
* bigtable: improved regional reliability when instance overlaps a downed region in the resource `google_bigtable_instance` ([#15900](https://github.com/hashicorp/terraform-provider-google/pull/15900))
* eventarc: resolved permadiff on `google_eventarc_trigger.event_data_content_type` by defaulting to the value returned by the API if not set in the configuration. ([#15862](https://github.com/hashicorp/terraform-provider-google/pull/15862))
* identityplatform: fixed a potential perma-diff for `sign_in` in `google_identity_platform_config` resource ([#15907](https://github.com/hashicorp/terraform-provider-google/pull/15907))
* monitoring: fixed scaling issues when deploying terraform changes with many `google_monitoring_monitored_project` ([#15828](https://github.com/hashicorp/terraform-provider-google/pull/15828))
* monitoring: fixed validation of `service_id` on `google_monitoring_custom_service` and `slo_id` on `google_monitoring_slo` ([#15841](https://github.com/hashicorp/terraform-provider-google/pull/15841))
* osconfig: fixed no more than one setting is allowed under `patch_config.windows_update` on `google_os_config_patch_deployment` ([#15904](https://github.com/hashicorp/terraform-provider-google/pull/15904))
* provider: addressed a bug where configuring the provider with unknown values did not behave as expected ([#15898](https://github.com/hashicorp/terraform-provider-google/pull/15898))
* provider: fixed the provider so it resumes ignoring empty strings set in the `provider` block ([#15844](https://github.com/hashicorp/terraform-provider-google/pull/15844))
* secretmanager: replaced the panic block with an error in import function of `google_secret_manager_secret_version` resource ([#15880](https://github.com/hashicorp/terraform-provider-google/pull/15880))
* secretmanager: fixed an issue in `google_secretmanager_secret` where replacing `replication.automatic` with `replication.auto` would destroy and recreate the resource ([#15922](https://github.com/hashicorp/terraform-provider-google/pull/15922))

## 4.83.0 (September 18, 2023)

DEPRECATIONS:
* secretmanager: deprecated `automatic` field on `google_secret_manager_secret`. Use `auto` instead. ([#15793](https://github.com/hashicorp/terraform-provider-google/pull/15793))

FEATURES:
* **New Resource:** `google_biglake_table` ([#15736](https://github.com/hashicorp/terraform-provider-google/pull/15736))
* **New Resource:** `google_data_pipeline_pipeline` ([#15785](https://github.com/hashicorp/terraform-provider-google/pull/15785))
* **New Resource:** `google_dialogflow_cx_test_case` ([#15814](https://github.com/hashicorp/terraform-provider-google/pull/15814))
* **New Resource:** `google_storage_insights_report_config` ([#15819](https://github.com/hashicorp/terraform-provider-google/pull/15819))
* **New Resource:** `google_apigee_target_server` ([#15751](https://github.com/hashicorp/terraform-provider-google/pull/15751))

IMPROVEMENTS:
* gkehub: added `labels` fields to `google_gke_hub_membership_binding` resource ([#15753](https://github.com/hashicorp/terraform-provider-google/pull/15753))
* bigquery: added `allow_non_incremental_definition` to `google_bigquery_table` resource ([#15813](https://github.com/hashicorp/terraform-provider-google/pull/15813))
* bigquery: added `table_constraints` field to `google_bigquery_table` resource ([#15815](https://github.com/hashicorp/terraform-provider-google/pull/15815))
* compute: added internal IPV6 support for `google_compute_address` and `google_compute_instance` resources ([#15780](https://github.com/hashicorp/terraform-provider-google/pull/15780))
* containerattached: added `binary_authorization` field to `google_container_attached_cluster` resource ([#15822](https://github.com/hashicorp/terraform-provider-google/pull/15822))
* containeraws: added update support for `config.instance_type` in `container_aws_node_pool` ([#15862](https://github.com/hashicorp/terraform-provider-google/pull/15862))
* firestore: added `point_in_time_recovery_enablement` field to `google_firestore_database` resource ([#15795](https://github.com/hashicorp/terraform-provider-google/pull/15795))
* firestore: added `update_time` and `uid` fields to `google_firestore_database` resource ([#15823](https://github.com/hashicorp/terraform-provider-google/pull/15823))
* gkehub2: added `labels`, `namespace_labels` fields to `google_gke_hub_namespace` resource ([#15732](https://github.com/hashicorp/terraform-provider-google/pull/15732))
* gkehub: added `labels` fields to `google_gke_hub_scope` resource ([#15801](https://github.com/hashicorp/terraform-provider-google/pull/15801))
* gkeonprem: added `upgrade_policy` and `binary_authorization` fields in `google_gkeonprem_bare_metal_cluster` resource (beta) ([#15765](https://github.com/hashicorp/terraform-provider-google/pull/15765))
* gkeonprem: added `upgrade_policy` field in `google_gkeonprem_vmware_cluster` resource (beta) ([#15765](https://github.com/hashicorp/terraform-provider-google/pull/15765))
* secretmanager: added `auto` field to `google_secret_manager_secret` resource ([#15793](https://github.com/hashicorp/terraform-provider-google/pull/15793))
* secretmanager: added `deletion_policy` field to `google_secret_manager_secret_version` resource ([#15818](https://github.com/hashicorp/terraform-provider-google/pull/15818))
* storage: supported in-place update for `autoclass` field in `google_storage_bucket` resource ([#15782](https://github.com/hashicorp/terraform-provider-google/pull/15782))
* vertexai: added `public_endpoint_enabled` to `google_vertex_ai_index_endpoint` ([#15741](https://github.com/hashicorp/terraform-provider-google/pull/15741))

BUG FIXES:
* bigquerydatatransfer: fixed a bug when importing `location` of `google_bigquery_data_transfer_config` ([#15734](https://github.com/hashicorp/terraform-provider-google/pull/15734))
* container: fixed concurrent ops' quota-error to be retriable in `google_container_node_pool ` ([#15820](https://github.com/hashicorp/terraform-provider-google/pull/15820))
* eventarc: resolved permadiff on `event_content_type` in `eventarc_trigger`, the field will now default to a value returned by the API when not set in configuration ([#15862](https://github.com/hashicorp/terraform-provider-google/pull/15862))
* pipeline: fixed issue where certain `google_dataflow_job` instances would crash the provider ([#15821](https://github.com/hashicorp/terraform-provider-google/pull/15821))
* provider: fixed a bug where `user_project_override` would not be not used correctly when provisioning resources implemented using the plugin framework. Currently there are no resources implemented this way, so no-one should have been impacted. ([#15776](https://github.com/hashicorp/terraform-provider-google/pull/15776))
* pubsub: fixed issue where setting `no_wrapper.write_metadata` to false wasn't passed to the API for `google_pubsub_subscription` ([#15758](https://github.com/hashicorp/terraform-provider-google/pull/15758))
* serviceaccount: added retries for reads after `google_service_account` creation if 403 Forbidden is returned. ([#15760](https://github.com/hashicorp/terraform-provider-google/pull/15760))
* storage: fixed the failure in building a plan when a `content` value is expected on `google_storage_bucket_object_content` ([#15735](https://github.com/hashicorp/terraform-provider-google/pull/15735))

## 4.82.0 (September 11, 2023)

IMPROVEMENTS:
* compute: added in-place update support for field `enable_proxy_protocol` in `google_compute_service_attachment` resource ([#15716](https://github.com/hashicorp/terraform-provider-google/pull/15716))
* compute: added in-place update support for field `reconcile_connections` in `google_compute_service_attachment` resource ([#15706](https://github.com/hashicorp/terraform-provider-google/pull/15706))
* compute: added in-place update support for field `allowPscGlobalAccess` in `google_compute_forwarding_rule` resource ([#15691](https://github.com/hashicorp/terraform-provider-google/pull/15691))
* compute: promoted `google_compute_region_instance_template` to GA ([#15710](https://github.com/hashicorp/terraform-provider-google/pull/15710))
* container: added additional options for field `monitoring_config.enable_components` in `google_container_cluster` resource ([#15727](https://github.com/hashicorp/terraform-provider-google/pull/15727))
* gkehub: added `labels` field to `google_gke_hub_scope_rbac_role_binding` resource ([#15729](https://github.com/hashicorp/terraform-provider-google/pull/15729))
* logging: added in-place update support for field `unique_writer_identity` in `google_logging_project_sink` resource ([#15721](https://github.com/hashicorp/terraform-provider-google/pull/15721))
* networkconnectivity: added `psc_connections.error.details` field to `google_network_connectivity_service_connection_policy` resource ([#15726](https://github.com/hashicorp/terraform-provider-google/pull/15726))
* secretmanager: added in-place update support for field `replication.user_managed.replicas.customer_managed_encryption` in `google_secret_manager_secret` resource ([#15685](https://github.com/hashicorp/terraform-provider-google/pull/15685))

BUG FIXES:
* bigquery: made `params.destination_table_name_template` and `params.data_path` immutable as updating these fields if value of `data_source_id` is `amazon_s3` in `google_bigquery_data_transfer_config` resource ([#15723](https://github.com/hashicorp/terraform-provider-google/pull/15723))
* dns: fixed hash function for `network_url` in `google_dns_managed_zone` and `google_dns_policy` resources to make sure that the private DNS zone or DNS policy can be attatched to all of the networks in different projects, even though the network name is the same across of those projects. ([#15728](https://github.com/hashicorp/terraform-provider-google/pull/15728))

## 4.81.0 (September 05, 2023)

FEATURES:
* **New Resource:** `google_biglake_catalog` ([#15634](https://github.com/hashicorp/terraform-provider-google/pull/15634))
* **New Resource:** `google_redis_cluster` ([#15645](https://github.com/hashicorp/terraform-provider-google/pull/15645))
* **New Resource:** `google_biglake_database` ([#15651](https://github.com/hashicorp/terraform-provider-google/pull/15651))
* **New Resource:** `google_compute_network_attachment` ([#15648](https://github.com/hashicorp/terraform-provider-google/pull/15648))
* **New Resource:** `google_gke_hub_feature_membership` ([#15604](https://github.com/hashicorp/terraform-provider-google/pull/15604))
* **New Resource:** `google_gke_hub_membership_binding` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))
* **New Resource:** `google_gke_hub_namespace` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))
* **New Resource:** `google_gke_hub_scope` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))
* **New Resource:** `google_gke_hub_scope_iam_member` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))
* **New Resource:** `google_gke_hub_scope_iam_policy` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))
* **New Resource:** `google_gke_hub_membership_binding` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))
* **New Resource:** `google_gke_hub_scope_rbac_role_binding` ([#15670](https://github.com/hashicorp/terraform-provider-google/pull/15670))

IMPROVEMENTS:
* compute: made the field `distribution_policy_target_shape` of `google_compute_region_instance_group_manager` not cause recreation of the resource. ([#15641](https://github.com/hashicorp/terraform-provider-google/pull/15641))
* compute: promoted the `ssl_policy` field on the `google_compute_region_target_https_proxy` resource to GA. ([#15608](https://github.com/hashicorp/terraform-provider-google/pull/15608))
* container: added `enable_fqdn_network_policy` field to `google_container_cluster` ([#15642](https://github.com/hashicorp/terraform-provider-google/pull/15642))
* container: added `node_config.confidential_compute` field to `google_container_node_pool` resource ([#15662](https://github.com/hashicorp/terraform-provider-google/pull/15662))
* datastream: made `password` in `google_datastream_connection_profile` not cause recreation of the resource. ([#15610](https://github.com/hashicorp/terraform-provider-google/pull/15610))
* dialogflowcx: added `response_type`, `channel`, `payload`, `conversation_success`, `output_audio_text`, `live_agent_handoff`, `play_audo`, `telephony_transfer_call`, `reprompt_event_handlers`, `set_parameter_actions`, and `conditional_cases` fields to `google_dialogflow_cx_page` resource ([#15668](https://github.com/hashicorp/terraform-provider-google/pull/15668))
* dialogflowcx: added `response_type`, `channel`, `payload`, `conversation_success`, `output_audio_text`, `live_agent_handoff`, `play_audo`, `telephony_transfer_call`, `set_parameter_actions`, and `conditional_cases` fields to `google_dialogflow_cx_flow` resource ([#15668](https://github.com/hashicorp/terraform-provider-google/pull/15668))
* iam: added `web_sso_config.additional_scopes` field to `google_iam_workforce_pool_provider` resource under ([#15616](https://github.com/hashicorp/terraform-provider-google/pull/15616))
* monitoring: added `synthetic_monitor` to `google_monitoring_uptime_check_config` resource ([#15623](https://github.com/hashicorp/terraform-provider-google/pull/15623))
* provider: improved error message when resource creation fails to to invalid API response ([#15629](https://github.com/hashicorp/terraform-provider-google/pull/15629))

BUG FIXES:
* cloudrunv2: changed `template.volumes.secret.items.mode` field in `google_cloud_run_v2_job` resource to a non-required field. ([#15638](https://github.com/hashicorp/terraform-provider-google/pull/15638))
* cloudrunv2: changed `template.volumes.secret.items.mode` field in `google_cloud_run_v2_service` resource to a non-required field. ([#15638](https://github.com/hashicorp/terraform-provider-google/pull/15638))
* filestore: fixed a bug causing permadiff on `reserved_ip_range` field in `google_filestore_instance` ([#15614](https://github.com/hashicorp/terraform-provider-google/pull/15614))
* identityplatform: fixed a permadiff on `authorized_domains` in `google_identity_platform_config` resource ([#15607](https://github.com/hashicorp/terraform-provider-google/pull/15607))


## 4.80.0 (August 28, 2023)

DEPRECATIONS:
* dataplex: deprecated the following `google_dataplex_datascan` fields: `dataProfileResult` and `dataQualityResult` ([#15528](https://github.com/hashicorp/terraform-provider-google/pull/15528))
* firebase: deprecated `google_firebase_project_location` in favor of `google_firebase_storage_bucket` and `google_firestore_database` ([#15526](https://github.com/hashicorp/terraform-provider-google/pull/15526))

FEATURES:
* **New Data Source:** `google_sql_database_instance_latest_recovery_time` ([#15551](https://github.com/hashicorp/terraform-provider-google/pull/15551))
* **New Resource:** `google_certificate_manager_trust_config` ([#15562](https://github.com/hashicorp/terraform-provider-google/pull/15562))
* **New Resource:** `google_compute_region_security_policy_rule` ([#15523](https://github.com/hashicorp/terraform-provider-google/pull/15523))
* **New Resource:** `google_iam_deny_policy` ([#15571](https://github.com/hashicorp/terraform-provider-google/pull/15571))
* **New Resource:** `google_bigquery_bi_reservation` ([#15527](https://github.com/hashicorp/terraform-provider-google/pull/15527))
* **New Resource:** `google_gke_hub_feature_membership` ([#15604](https://github.com/hashicorp/terraform-provider-google/pull/15604))

IMPROVEMENTS:
* alloydb: added `restore_backup_source` and `restore_continuous_backup_source` fields to support restore feature in `google_alloydb_cluster` resource. ([#15580](https://github.com/hashicorp/terraform-provider-google/pull/15580))
* artifactregistry: added `cleanup_policies` and `cleanup_policy_dry_run` fields to resource `google_artifact_registry_repository` ([#15561](https://github.com/hashicorp/terraform-provider-google/pull/15561))
* clouddeploy: added `multi_target` to in `google_clouddelploy_target` ([#15564](https://github.com/hashicorp/terraform-provider-google/pull/15564))
* compute: added `security_policy` field to `google_compute_target_instance` resource (beta) ([#15566](https://github.com/hashicorp/terraform-provider-google/pull/15566))
* compute: added support for `security_policy` field to `google_compute_target_pool` ([#15569](https://github.com/hashicorp/terraform-provider-google/pull/15569))
* compute: added support for `user_defined_fields` to `google_compute_region_security_policy` ([#15523](https://github.com/hashicorp/terraform-provider-google/pull/15523))
* compute: added support for specifying regional disks for `google_compute_instance` `boot_disk.source` ([#15597](https://github.com/hashicorp/terraform-provider-google/pull/15597))
* container: added `additional_pod_ranges_config` field to `google_container_cluster` resource ([#15600](https://github.com/hashicorp/terraform-provider-google/pull/15600))
* containeraws: made `config.labels` updatable in `google_container_aws_node_pool` ([#15564](https://github.com/hashicorp/terraform-provider-google/pull/15564))
* dataplex: added fields `data_profile_spec.post_scan_actions`, `data_profile_spec.include_fields` and `data_profile_spec.exclude_fields` ([#15545](https://github.com/hashicorp/terraform-provider-google/pull/15545))
* dns: added support for removing the networks block from the configuration in the resource `google_dns_response_policy` ([#15557](https://github.com/hashicorp/terraform-provider-google/pull/15557))
* firebase: added `api_key_id` field to `google_firebase_web_app`, `google_firebase_android_app`, and `google_firebase_apple_app`. ([#15577](https://github.com/hashicorp/terraform-provider-google/pull/15577))
* sql: added `psc_config` , `psc_service_attachment_link`, and `dns_name` fields to `google_sql_database_instance` ([#15563](https://github.com/hashicorp/terraform-provider-google/pull/15563))
* workstations: added `enable_nested_virtualization` field to `google_workstations_workstation_config` resource ([#15567](https://github.com/hashicorp/terraform-provider-google/pull/15567))

BUG FIXES:
* bigquery: added support to unset policy tags in table schema ([#15547](https://github.com/hashicorp/terraform-provider-google/pull/15547))
* bigtable: fixed permadiff in `google_bigtable_gc_policy.gc_rules` when `max_age` is specified using increments larger than hours ([#15595](https://github.com/hashicorp/terraform-provider-google/pull/15595))
* bigtable: fixed permadiff in `google_bigtable_gc_policy.gc_rules` when `mode` is specified ([#15595](https://github.com/hashicorp/terraform-provider-google/pull/15595))
* container: updated `resource_container_cluster` to ignore `dns_config` diff when `enable_autopilot = true` ([#15549](https://github.com/hashicorp/terraform-provider-google/pull/15549))
* containerazure: added diff suppression for case changes of enum values in `google_container_azure_cluster` ([#15536](https://github.com/hashicorp/terraform-provider-google/pull/15536))


## 4.79.0 (August 21, 2023)
FEATURES:
* **New Resource:** `google_backup_dr_management_server` ([#15479](https://github.com/hashicorp/terraform-provider-google/pull/15479))
* **New Resource:** `google_compute_region_security_policy_rule` ([#15523](https://github.com/hashicorp/terraform-provider-google/pull/15523))

IMPROVEMENTS:
* cloudbuild: added `git_file_source.bitbucket_server_config` and `source_to_build.bitbucket_server_config` fields to `google_cloudbuild_trigger` resource ([#15475](https://github.com/hashicorp/terraform-provider-google/pull/15475))
* cloudrunv2: added the following output only fields to `google_cloud_run_v2_job` and `google_cloud_run_v2_service` resources: `create_time`, `update_time`, `delete_time`, `expire_time`, `creator` and `last_modifier` ([#15502](https://github.com/hashicorp/terraform-provider-google/pull/15502))
* composer: added `config.private_environment_config.connection_type` field to `google_composer_environment` resource ([#15460](https://github.com/hashicorp/terraform-provider-google/pull/15460))
* compute: added `disk.provisioned_iops` field to `google_compute_instance_template` and `google_compute_region_instance_template` resources ([#15506](https://github.com/hashicorp/terraform-provider-google/pull/15506))
* compute: added `user_defined_fields` field to `google_compute_region_security_policy` resource ([#15523](https://github.com/hashicorp/terraform-provider-google/pull/15523))
* databasemigrationservice: added `edition` field to `google_database_migration_service_connection_profile` resource ([#15510](https://github.com/hashicorp/terraform-provider-google/pull/15510))
* dns: allowed `globalL7ilb` value for the `routing_policy.load_balancer_type` field in `google_dns_record_set` resource ([#15521](https://github.com/hashicorp/terraform-provider-google/pull/15521))
* healthcare: added `default_search_handling_strict` field to `google_healthcare_fhir_store` resource ([#15514](https://github.com/hashicorp/terraform-provider-google/pull/15514))
* metastore: added `scaling_config` field to `google_dataproc_metastore_service` resource ([#15476](https://github.com/hashicorp/terraform-provider-google/pull/15476))
* secretmanager: added `version_aliases` field to `google_secret_manager_secret` resource ([#15483](https://github.com/hashicorp/terraform-provider-google/pull/15483))

BUG FIXES:
* alloydb: fixed a permadiff on `google_alloydb_cluster` when `backup_window`, `enabled` or `location` fields are unset ([#15444](https://github.com/hashicorp/terraform-provider-google/pull/15444))
* containeraws: fixed permadiffs on `google_container_aws_cluster` and `google_container_aws_node_pool` resources ([#15491](https://github.com/hashicorp/terraform-provider-google/pull/15491))
* dataplex: fixed a bug when importing `google_dataplex_datascan` after running a job ([#15468](https://github.com/hashicorp/terraform-provider-google/pull/15468))
* dns: changed `private_visibility_config.networks` from `required` to requiring at least one of `private_visibility_config.networks` or `private_visibility_config.gke_clusters` in `google_dns_managed_zone` resource ([#15443](https://github.com/hashicorp/terraform-provider-google/pull/15443))


## 4.78.0 (August 15, 2023)

FEATURES:
* **New Resource:** `google_billing_project_info` ([#15400](https://github.com/hashicorp/terraform-provider-google/pull/15400))
* **New Resource:** `google_network_connectivity_service_connection_policy` ([#15381](https://github.com/hashicorp/terraform-provider-google/pull/15381))

IMPROVEMENTS:
* alloydb: added `continuous_backup_config` and `continuous_backup_info` fields to `cluster` resource ([#15370](https://github.com/hashicorp/terraform-provider-google/pull/15370))
* bigquery: added `external_data_configuration.file_set_spec_type` to `google_bigquery_table` ([#15402](https://github.com/hashicorp/terraform-provider-google/pull/15402))
* bigquery: added `max_staleness` to `google_bigquery_table` ([#15395](https://github.com/hashicorp/terraform-provider-google/pull/15395))
* billingbudget: added `resource_ancestors` field to `google_billing_budget` resource ([#15393](https://github.com/hashicorp/terraform-provider-google/pull/15393))
* cloudfunctions2: added support for GCF Gen2 CMEK ([#15385](https://github.com/hashicorp/terraform-provider-google/pull/15385))
* cloudidentity: added field `type` to `google_cloud_identity_group_memberships` ([#15398](https://github.com/hashicorp/terraform-provider-google/pull/15398))
* compute: added `subnetwork` field to the resource `google_compute_global_forwarding_rule` ([#15424](https://github.com/hashicorp/terraform-provider-google/pull/15424))
* compute: added support for `INTERNAL_MANAGED` to the field `load_balancing_scheme` in the resource `google_compute_backend_service` ([#15424](https://github.com/hashicorp/terraform-provider-google/pull/15424))
* compute: added support for `INTERNAL_MANAGED` to the field `load_balancing_scheme` in the resource `google_compute_global_forwarding_rule` ([#15424](https://github.com/hashicorp/terraform-provider-google/pull/15424))
* compute: added support for `ip_version` to `google_compute_forwarding_rule` ([#15388](https://github.com/hashicorp/terraform-provider-google/pull/15388))
* container: marked `master_ipv4_cidr_block` as not required when `private_endpoint_subnetwork` is provided for `google_container_cluster` ([#15422](https://github.com/hashicorp/terraform-provider-google/pull/15422))
* container: added support for `advanced_datapath_observability_config` to `google_container_cluster` ([#15425](https://github.com/hashicorp/terraform-provider-google/pull/15425))
* eventarc: added field `event_data_content_type` to `google_eventarc_trigger` ([#15433](https://github.com/hashicorp/terraform-provider-google/pull/15433))
* healthcare: added `send_previous_resource_on_delete` field to `notification_configs` of `google_healthcare_fhir_store` ([#15380](https://github.com/hashicorp/terraform-provider-google/pull/15380))
* pubsub: added `cloud_storage_config` field to `google_pubsub_subscription` resource ([#15420](https://github.com/hashicorp/terraform-provider-google/pull/15420))
* secretmanager: added `annotations` field to `google_secret_manager_secret` resource ([#15392](https://github.com/hashicorp/terraform-provider-google/pull/15392))

BUG FIXES:
* certificatemanager: added recreation behavior to the `google_certificate_manager_certificate` resource when its location changes ([#15432](https://github.com/hashicorp/terraform-provider-google/pull/15432))
* cloudfunctions2: fixed creation failure state inconsistency in `google_cloudfunctions2_function` ([#15418](https://github.com/hashicorp/terraform-provider-google/pull/15418))
* monitoring: updated `evaluation_interval` on `condition_prometheus_query_language` to be optional ([#15429](https://github.com/hashicorp/terraform-provider-google/pull/15429))

## 4.77.0 (August 7, 2023)

NOTES:
* vpcaccess: reverted the ability to update the number of instances for resource `google_vpc_access_connector` ([#15313](https://github.com/hashicorp/terraform-provider-google/pull/15313))

FEATURES:
* **New Resource:** `google_document_ai_warehouse_document_schema` ([#15326](https://github.com/hashicorp/terraform-provider-google/pull/15326))
* **New Resource:** `google_document_ai_warehouse_location` ([#15326](https://github.com/hashicorp/terraform-provider-google/pull/15326))

IMPROVEMENTS:
* alloydb: added `continuous_backup_config` and `continuous_backup_info` fields to `cluster` resource ([#15370](https://github.com/hashicorp/terraform-provider-google/pull/15370))
* cloudbuild: removed the validation function for the values of `machine_type` field on the `google_cloudbuild_trigger` resource ([#15357](https://github.com/hashicorp/terraform-provider-google/pull/15357))
* compute: add future_limit in quota exceeded error details for compute resources. ([#15346](https://github.com/hashicorp/terraform-provider-google/pull/15346))
* compute: added `ipv6_endpoint_type` and `ip_version` to `google_compute_address` ([#15358](https://github.com/hashicorp/terraform-provider-google/pull/15358))
* compute: added `local_ssd_recovery_timeout` field to `google_compute_instance` resource ([#15366](https://github.com/hashicorp/terraform-provider-google/pull/15366))
* compute: added `local_ssd_recovery_timeout` field to `google_compute_instance_template` resource ([#15366](https://github.com/hashicorp/terraform-provider-google/pull/15366))
* compute: added `network_interface.ipv6_access_config.external_ipv6_prefix_length` to `google_compute_instance` ([#15358](https://github.com/hashicorp/terraform-provider-google/pull/15358))
* compute: added `network_interface.ipv6_access_config.name` to `google_compute_instance` ([#15358](https://github.com/hashicorp/terraform-provider-google/pull/15358))
* compute: added a new type `GLOBAL_MANAGED_PROXY` for the field `purpose` in the resource `google_compute_subnetwork` ([#15345](https://github.com/hashicorp/terraform-provider-google/pull/15345))
* compute: added field `instance_lifecycle_policy` to `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#15322](https://github.com/hashicorp/terraform-provider-google/pull/15322))
* compute: added protocol type: UNSPECIFIED in `google_compute_backend_service` as per [release note](https://cloud.google.com/load-balancing/docs/release-notes#July_24_2023)
 ([#15328](https://github.com/hashicorp/terraform-provider-google/pull/15328))
* compute: made `network_interface.ipv6_access_config.external_ipv6` configurable in `google_compute_instance` ([#15358](https://github.com/hashicorp/terraform-provider-google/pull/15358))
* container: added `enable_k8s_beta_apis.enabled_apis` field to `google_container_cluster` ([#15320](https://github.com/hashicorp/terraform-provider-google/pull/15320))
* container: added `node_config.host_maintenance_policy` field to `google_container_cluster` and `google_container_node_pool` ([#15347](https://github.com/hashicorp/terraform-provider-google/pull/15347))
* container: added `placement_policy.policy_name` field to `google_container_node_pool` resource ([#15367](https://github.com/hashicorp/terraform-provider-google/pull/15367))
* container: allowed `enabled_private_endpoint` to be settable on creation for PSC-based clusters ([#15361](https://github.com/hashicorp/terraform-provider-google/pull/15361))
* container: unsuppressed `private_cluster_config` when `master_global_access_config` is set in `google_container_cluster` ([#15369](https://github.com/hashicorp/terraform-provider-google/pull/15369))
* gkeonprem: added taint on failed resource creation for `google_gkeonprem_bare_metal_admin_cluster` ([#15362](https://github.com/hashicorp/terraform-provider-google/pull/15362))
* gkeonprem: increased timeout for resources `google_gkeonprem_bare_metal_cluster` and `google_gkeonprem_bare_metal_admin_cluster` ([#15362](https://github.com/hashicorp/terraform-provider-google/pull/15362))
* identityplayform: added support for `blocking_functions` `quota` and `authorized_domains` in `google_identity_platform_config` ([#15325](https://github.com/hashicorp/terraform-provider-google/pull/15325))
* monitoring: added update support for `period` in `google_monitoring_uptime_check_config` ([#15315](https://github.com/hashicorp/terraform-provider-google/pull/15315))
* pubsub: added `no_wrapper` field to `google_pubsub_subscription` resource ([#15334](https://github.com/hashicorp/terraform-provider-google/pull/15334))

BUG FIXES:
* bigquery: fixed a bug in update support for several fields in `google_bigquery_data_transfer_config` ([#15359](https://github.com/hashicorp/terraform-provider-google/pull/15359))
* cloudfunctions2: fixed an issue where `google_cloudfunctions2_function.build_config.source.storage_source.generation` created a diff when not set in config ([#15364](https://github.com/hashicorp/terraform-provider-google/pull/15364))
* monitoring: fixed an issue in `google_monitoring_monitored_project` where project numbers were not accepted for `name` ([#15305](https://github.com/hashicorp/terraform-provider-google/pull/15305))
* vpcaccess: reverted new behaviour introduced by resource `google_vpc_access_connector` in `4.75.0`. `min_throughput` and `max_throughput` fields lost their default value, and customers could not make deployment due to that change. ([#15313](https://github.com/hashicorp/terraform-provider-google/pull/15313))

## 4.76.0 (July 31, 2023)

FEATURES:
* **New Resource:** `google_compute_region_ssl_policy` ([#15299](https://github.com/hashicorp/terraform-provider-google/pull/15299))
* **New Resource:** `google_dataplex_task` ([#15226](https://github.com/hashicorp/terraform-provider-google/pull/15226))
* **New Resource:** `google_iap_web_region_backend_service_iam_binding` ([#15285](https://github.com/hashicorp/terraform-provider-google/pull/15285))
* **New Resource:** `google_iap_web_region_backend_service_iam_member` ([#15285](https://github.com/hashicorp/terraform-provider-google/pull/15285))
* **New Resource:** `google_iap_web_region_backend_service_iam_policy` ([#15285](https://github.com/hashicorp/terraform-provider-google/pull/15285))

IMPROVEMENTS:
* cloudrun: added `status.traffic` output fields to `google_cloud_run_service` resource ([#15284](https://github.com/hashicorp/terraform-provider-google/pull/15284))
* cloudrunv2: added field `custom_audiences` to resource `google_cloud_run_v2_service ` ([#15268](https://github.com/hashicorp/terraform-provider-google/pull/15268))
* composer: added support for updating `resilience_mode` in `google_composer_environment` ([#15238](https://github.com/hashicorp/terraform-provider-google/pull/15238))
* compute: added `reconcile_connections` for `google_compute_service_attachment`. ([#15288](https://github.com/hashicorp/terraform-provider-google/pull/15288))
* container : added `gcs_fuse_csi_driver_config` field to `addons_config` in `google_container_cluster` resource. ([#15290](https://github.com/hashicorp/terraform-provider-google/pull/15290))
* container: added `allow_net_admin` field to `google_container_cluster` resource ([#15275](https://github.com/hashicorp/terraform-provider-google/pull/15275))
* container: allowed user to set up to 20 maintenance exclusions for `google_container_cluster` resource ([#15291](https://github.com/hashicorp/terraform-provider-google/pull/15291))
* healthcare: added `last_updated_partition_config` field to `google_healthcare_fhir_store` resource ([#15271](https://github.com/hashicorp/terraform-provider-google/pull/15271))
* monitoring: added `condition_prometheus_query_language` field to `google_monitoring_alert_policy` resource ([#15301](https://github.com/hashicorp/terraform-provider-google/pull/15301))
* networkservices: made `scope` field optional in `google_network_services_gateway` resource ([#15273](https://github.com/hashicorp/terraform-provider-google/pull/15273))
* spanner: added `enable_drop_protection` to `google_spanner_database` resource([#15283](https://github.com/hashicorp/terraform-provider-google/pull/15283))

BUG FIXES:
* alloydb: fixed permadiffs when setting 0 as start time (midnight) for `automated_backup_policy` in `google_alloydb_cluster` resource ([#15219](https://github.com/hashicorp/terraform-provider-google/pull/15219))
* artifactregistry: fixed reading back maven_config state in `google_artifact_registry_repository` ([#15269](https://github.com/hashicorp/terraform-provider-google/pull/15269))
* cloudtasks: suppressed time-unit permadiffs on `google_cloud_tasks_queue` min and max backoff settings ([#15237](https://github.com/hashicorp/terraform-provider-google/pull/15237))
* cloudrun: fixed the bug where default system labels set in `service.spec.template.metadata.labels` were treated as a diff. ([#15302](https://github.com/hashicorp/terraform-provider-google/pull/15302))
* compute: fixed wrongly triggered recreation on changes of `enforce_on_key_configs` on `google_compute_security_policy` ([#15248](https://github.com/hashicorp/terraform-provider-google/pull/15248))
* monitoring: fixed an issue in `google_monitoring_monitored_project` where project numbers were not accepted for `name` ([#15305](https://github.com/hashicorp/terraform-provider-google/pull/15305))

## 4.75.1 (July 27, 2023)

BUG FIXES:

* vpcaccess: reverted new behaviour introduced by resource `google_vpc_access_connector` in `4.75.0`. `min_throughput` and `max_throughput` fields lost their default value, and customers could not make deployment due to that change.

* vpcaccess: reverted the ability to update the number of instances for resource `google_vpc_access_connector`

## 4.75.0 (July 24, 2023)

FEATURES:
* **New Resource:** `google_dns_response_policy_rule`([#15146](https://github.com/hashicorp/terraform-provider-google/pull/15146))
* **New Resource:** `google_dns_response_policy`([#15146](https://github.com/hashicorp/terraform-provider-google/pull/15146))
* **New Resource:** `google_looker_instance` ([#15188](https://github.com/hashicorp/terraform-provider-google/pull/15188))

IMPROVEMENTS:
* apigee: added `disable_vpc_peering` field to `google_apigee_organization` resource ([#15186](https://github.com/hashicorp/terraform-provider-google/pull/15186))
* bigquery: added `external_data_configuration.json_options` and `external_data_configuration.parquet_options` fields to `google_bigquery_table` ([#15197](https://github.com/hashicorp/terraform-provider-google/pull/15197))
* bigtable: added `change_stream_retention` field to `google_bigtable_table.table` resource ([#15152](https://github.com/hashicorp/terraform-provider-google/pull/15152))
* compute: added `most_recent` argument to `google_compute_image` datasource ([#15187](https://github.com/hashicorp/terraform-provider-google/pull/15187))
* compute: added field `enable_confidential_compute` for `google_compute_disk` resource ([#15180](https://github.com/hashicorp/terraform-provider-google/pull/15180))
* container: added `gpu_driver_installation_config.gpu_driver_version` field to `google_container_node_pool` ([#15182](https://github.com/hashicorp/terraform-provider-google/pull/15182))
* gkebackup: added `state` and `state_reason` output-only fields to `google_gkebackup_backupplan` resource ([#15201](https://github.com/hashicorp/terraform-provider-google/pull/15201))
* healthcare: added `complex_data_type_reference_parsing ` field to `google_healthcare_fhir_store` resource ([#15159](https://github.com/hashicorp/terraform-provider-google/pull/15159))
* networkservices: increased max_size to 20 for both `included_query_parameters` and `excluded_query_parameters` on `google_network_services_edge_cache_service` ([#15168](https://github.com/hashicorp/terraform-provider-google/pull/15168))
* vpcaccess: added support for updates to `google_vpc_access_connector` resource ([#15176](https://github.com/hashicorp/terraform-provider-google/pull/15176))

BUG FIXES:
* alloydb: fixed `google_alloydb_cluster` handling of automated backup policy midnight start time ([#15219](https://github.com/hashicorp/terraform-provider-google/pull/15219))
* compute: fixed logic when unsetting `google_compute_instance.min_cpu_platform` and switching to a `machine_type` that does not support `min_cpu_platform` at the same time ([#15217](https://github.com/hashicorp/terraform-provider-google/pull/15217))
* tags: fixed race condition when modifying `google_tags_location_tag_binding` ([#15189](https://github.com/hashicorp/terraform-provider-google/pull/15189))


## 4.74.0 (July 18, 2023)

FEATURES:
* **New Resource:** `google_cloudbuildv2_connection` ([#15098](https://github.com/hashicorp/terraform-provider-google/pull/15098))
* **New Resource:** `google_cloudbuildv2_repository` ([#15098](https://github.com/hashicorp/terraform-provider-google/pull/15098))
* **New Resource:** `google_gkeonprem_bare_metal_admin_cluster` ([#15099](https://github.com/hashicorp/terraform-provider-google/pull/15099))
* **New Resource:** `google_network_security_address_group` ([#15111](https://github.com/hashicorp/terraform-provider-google/pull/15111))
* **New Resource:** `google_network_security_gateway_security_policy_rule` ([#15112](https://github.com/hashicorp/terraform-provider-google/pull/15112))
* **New Resource:** `google_network_security_gateway_security_policy` ([#15112](https://github.com/hashicorp/terraform-provider-google/pull/15112))
* **New Resource:** `google_network_security_url_lists` ([#15112](https://github.com/hashicorp/terraform-provider-google/pull/15112))
* **New Resource:** `google_network_services_gateway` ([#15112](https://github.com/hashicorp/terraform-provider-google/pull/15112))

IMPROVEMENTS:
* bigquery: added `storage_billing_model` argument to `google_bigquery_dataset` ([#15115](https://github.com/hashicorp/terraform-provider-google/pull/15115))
* bigquery: added `external_data_configuration.metadata_cache_mode` and `external_data_configuration.object_metadata` to `google_bigquery_table` ([#15096](https://github.com/hashicorp/terraform-provider-google/pull/15096))
* bigquery: made `external_data_configuration.source_fomat` optional in `google_bigquery_table` ([#15096](https://github.com/hashicorp/terraform-provider-google/pull/15096))
* certificatemanager: added `issuance_config` field to `google_certificate_manager_certificate` resource ([#15101](https://github.com/hashicorp/terraform-provider-google/pull/15101))
* cloudbuild: added `repository_event_config` field to `google_cloudbuild_trigger` resource ([#15098](https://github.com/hashicorp/terraform-provider-google/pull/15098))
* compute: added field `http_keep_alive_timeout_sec` to resource `google_compute_target_http_proxy` ([#15109](https://github.com/hashicorp/terraform-provider-google/pull/15109))
* compute: added field `http_keep_alive_timeout_sec` to resource `google_compute_target_https_proxy` ([#15109](https://github.com/hashicorp/terraform-provider-google/pull/15109))
* compute: added support for updating labels in `google_compute_external_vpn_gateway` ([#15134](https://github.com/hashicorp/terraform-provider-google/pull/15134))
* container: made `monitoring_config.enable_components` optional on `google_container_cluster` ([#15131](https://github.com/hashicorp/terraform-provider-google/pull/15131))
* container: added field `tpu_topology` under `placement_policy` in resource `google_container_node_pool` ([#15130](https://github.com/hashicorp/terraform-provider-google/pull/15130))
* gkehub: promoted the `google_gke_hub_feature` resource's `fleetobservability` block to GA. ([#15105](https://github.com/hashicorp/terraform-provider-google/pull/15105))
* iamworkforcepool: added `oidc.client_secret` field to `google_iam_workforce_pool_provider` and new enum values `CODE` and `MERGE_ID_TOKEN_OVER_USER_INFO_CLAIMS` to `oidc.web_sso_config.response_type` and `oidc.web_sso_config.assertion_claims_behavior` respectively ([#15069](https://github.com/hashicorp/terraform-provider-google/pull/15069))
* sql: added `settings.data_cache_config` to `sql_database_instance` resource. ([#15127](https://github.com/hashicorp/terraform-provider-google/pull/15127))
* sql: added `settings.edition` field to `sql_database_instance` resource. ([#15127](https://github.com/hashicorp/terraform-provider-google/pull/15127))
* vertexai: supported `shard_size` in `google_vertex_ai_index` ([#15133](https://github.com/hashicorp/terraform-provider-google/pull/15133))

BUG FIXES:
* compute: made `google_compute_router_peer.peer_ip_address` optional ([#15095](https://github.com/hashicorp/terraform-provider-google/pull/15095))
* redis: fixed issue with `google_redis_instance` populating output-only field `maintenance_schedule`. ([#15063](https://github.com/hashicorp/terraform-provider-google/pull/15063))
* orgpolicy: fixed forcing recreation on imported state for `google_org_policy_policy` ([#15132](https://github.com/hashicorp/terraform-provider-google/pull/15132))
* osconfig: fixed validation of file resource `state` fields in `google_os_config_os_policy_assignment` ([#15107](https://github.com/hashicorp/terraform-provider-google/pull/15107))

## 4.73.2 (July 17, 2023)

BUG FIXES:
* monitoring: fixed an issue which occurred when `name` field of `google_monitoring_monitored_project` was long-form

## 4.73.1 (July 13, 2023)

BUG FIXES:
* monitoring: fixed an issue causing `google_monitoring_monitored_project` to appear to be deleted

## 4.73.0 (July 10, 2023)

FEATURES:
* **New Resource:** `google_firebase_extensions_instance` ([#15013](https://github.com/hashicorp/terraform-provider-google/pull/15013))

IMPROVEMENTS:
* compute: added the `no_automate_dns_zone` field to `google_compute_forwarding_rule`. ([#15028](https://github.com/hashicorp/terraform-provider-google/pull/15028))
* compute: promoted `google_compute_disk_async_replication` resource to GA. ([#15029](https://github.com/hashicorp/terraform-provider-google/pull/15029))
* compute: promoted `async_primary_disk` field in `google_compute_disk` resource to GA. ([#15029](https://github.com/hashicorp/terraform-provider-google/pull/15029))
* compute: promoted `async_primary_disk` field in `google_compute_region_disk` resource to GA. ([#15029](https://github.com/hashicorp/terraform-provider-google/pull/15029))
* compute: promoted `disk_consistency_group_policy` field in `google_compute_resource_policy` resource to GA. ([#15029](https://github.com/hashicorp/terraform-provider-google/pull/15029))
* resourcemanager: fixed handling of `google_service_account_id_token` when authenticated with GCE metadata credentials ([#15003](https://github.com/hashicorp/terraform-provider-google/pull/15003))

BUG FIXES:
* networkservices: increased default timeout for `google_network_services_edge_cache_keyset` to 90m ([#15024](https://github.com/hashicorp/terraform-provider-google/pull/15024))

## 4.72.1 (July 6, 2023)

BUG FIXES:
* compute: fixed an issue in `google_compute_instance_template` where initialize params stopped the `disk.disk_size_gb` field being used ([#15054](https://github.com/hashicorp/terraform-provider-google/pull/15054))

## 4.72.0 (July 3, 2023)

FEATURES:
* **New Resource:** `google_public_ca_external_account_key` ([#14983](https://github.com/hashicorp/terraform-provider-google/pull/14983))

IMPROVEMENTS:
* compute: added `provisioned_throughput` field to `google_compute_disk` used by `hyperdisk-throughput` pd type ([#14985](https://github.com/hashicorp/terraform-provider-google/pull/14985))
* container: added field `security_posture_config` to resource `google_container_cluster` ([#14999](https://github.com/hashicorp/terraform-provider-google/pull/14999))
* logging: added support for `locked` to `google_logging_project_bucket_config` ([#14977](https://github.com/hashicorp/terraform-provider-google/pull/14977))

BUG FIXES:
* bigquery: fixed an issue where api default value for `edition` field of `google_bigquery_reservation` was not handled ([#14961](https://github.com/hashicorp/terraform-provider-google/pull/14961))
* cloudfunction2: fixed permadiffs of some fields of `service_config` in `google_cloudfunctions2_function` resource ([#14975](https://github.com/hashicorp/terraform-provider-google/pull/14975))
* compute: fixed an issue with setting project field to long form in `google_compute_forwarding_rule` and `google_compute_global_forwarding_rule` ([#14996](https://github.com/hashicorp/terraform-provider-google/pull/14996))
* gkehub: fixed an issue with setting project field to long form in `google_gke_hub_feature` ([#14996](https://github.com/hashicorp/terraform-provider-google/pull/14996))

## 4.71.0 (June 27, 2023)

FEATURES:
* **New Resource:** `google_gke_hub_feature_iam_*` ([#14912](https://github.com/hashicorp/terraform-provider-google/pull/14912))
* **New Resource:** `google_gke_hub_feature` ([#14912](https://github.com/hashicorp/terraform-provider-google/pull/14912))
* **New Resource:** `google_vmwareengine_cluster` ([#14917](https://github.com/hashicorp/terraform-provider-google/pull/14917))
* **New Resource:** `google_vmwareengine_private_cloud` ([#14917](https://github.com/hashicorp/terraform-provider-google/pull/14917))

IMPROVEMENTS:
* apigee: added output-only field `apigee_project_id` to resource `google_apigee_organization` ([#14911](https://github.com/hashicorp/terraform-provider-google/pull/14911))
* bigtable: increased default timeout for instance operations to 1 hour in resoure `google_bigtable_instance` ([#14909](https://github.com/hashicorp/terraform-provider-google/pull/14909))
* cloudrunv2: added fields `annotations` and `template.annotations` to resource `google_cloud_run_v2_job` ([#14948](https://github.com/hashicorp/terraform-provider-google/pull/14948))
* composer: added field `resilience_mode` to resource `google_composer_environment` ([#14939](https://github.com/hashicorp/terraform-provider-google/pull/14939))
* compute: added support for `params.resource_manager_tags` and `boot_disk.initialize_params.resource_manager_tags` to resource `google_compute_instance` ([#14924](https://github.com/hashicorp/terraform-provider-google/pull/14924))
* bigquerydatatransfer: made field `service_account_name` mutable in resource `google_bigquery_data_transfer_config` ([#14907](https://github.com/hashicorp/terraform-provider-google/pull/14907))
* iambeta: added field `jwks_json` to resource `google_iam_workload_identity_pool_provider` ([#14938](https://github.com/hashicorp/terraform-provider-google/pull/14938))

BUG FIXES:
* bigtable: validated that `cluster_id` values are unique within resource `google_bigtable_instance` ([#14908](https://github.com/hashicorp/terraform-provider-google/pull/14908))
* storage: fixed a bug that caused a permadiff when the `autoclass.enabled` field was explicitly set to false in resource `google_storage_bucket` ([#14902](https://github.com/hashicorp/terraform-provider-google/pull/14902))

## 4.70.0 (June 20, 2023)

FEATURES:
* **New Resource:** `google_compute_network_endpoints` ([#14869](https://github.com/hashicorp/terraform-provider-google/pull/14869))
* **New Resource:** `vertex_ai_index_endpoint` ([#14842](https://github.com/hashicorp/terraform-provider-google/pull/14842))

IMPROVEMENTS:
* bigtable: added 20 minutes timeout support to `google_bigtable_gc_policy` ([#14861](https://github.com/hashicorp/terraform-provider-google/pull/14861))
* cloudfunctions2: added `url` output field to `google_cloudfunctions2_function` ([#14851](https://github.com/hashicorp/terraform-provider-google/pull/14851))
* compute: added field `network_attachment` to `google_compute_instance_template` ([#14874](https://github.com/hashicorp/terraform-provider-google/pull/14874))
* compute: surfaced additional information about quota exceeded errors for compute resources. ([#14879](https://github.com/hashicorp/terraform-provider-google/pull/14879))
* compute: added `path_template_match` and `path_template_rewrite` to `google_compute_url_map`. ([#14873](https://github.com/hashicorp/terraform-provider-google/pull/14873))
* compute: added ability to update Hyperdisk PD IOPS without recreation to `google_compute_disk` ([#14844](https://github.com/hashicorp/terraform-provider-google/pull/14844))
* container: added `sole_tenant_config` to `node_config` in `google_container_node_pool` and `google_container_cluster` ([#14897](https://github.com/hashicorp/terraform-provider-google/pull/14897))
* dataform: added field `workspace_compilation_overrides` to resource `google_dataform_repository` (beta) ([#14839](https://github.com/hashicorp/terraform-provider-google/pull/14839))
* dlp: added `crypto_hash_config` to `google_data_loss_prevention_deidentify_template` ([#14870](https://github.com/hashicorp/terraform-provider-google/pull/14870))
* dlp: added `trigger_id` field to `google_data_loss_prevention_job_trigger` ([#14892](https://github.com/hashicorp/terraform-provider-google/pull/14892))
* dlp: added missing file types `POWERPOINT` and `EXCEL` in `inspect_job.storage_config.cloud_storage_options.file_types` enum to `google_data_loss_prevention_job_trigger` resource ([#14856](https://github.com/hashicorp/terraform-provider-google/pull/14856))
* dlp: added multiple `sensitivity_score` field to `google_data_loss_prevention_deidentify_template` resource ([#14880](https://github.com/hashicorp/terraform-provider-google/pull/14880))
* dlp: added multiple `sensitivity_score` field to `google_data_loss_prevention_inspect_template` resource ([#14871](https://github.com/hashicorp/terraform-provider-google/pull/14871))
* dlp: added multiple `sensitivity_score` field to `google_data_loss_prevention_job_trigger` resource ([#14881](https://github.com/hashicorp/terraform-provider-google/pull/14881))
* dlp: changed `inspect_template_name` field from required to optional in `google_data_loss_prevention_job_trigger` resource ([#14845](https://github.com/hashicorp/terraform-provider-google/pull/14845))
* pubsub: allowed `definition` field of `google_pubsub_schema` updatable. (https://cloud.google.com/pubsub/docs/schemas#commit-schema-revision) ([#14857](https://github.com/hashicorp/terraform-provider-google/pull/14857))
* sql: added `POSTGRES_15` to version docs for `database_version` field to `google_sql_database_instance` ([#14891](https://github.com/hashicorp/terraform-provider-google/pull/14891))
* vpcaccess: added `connected_projects` field to resource `google_vpc_access_connector`. ([#14835](https://github.com/hashicorp/terraform-provider-google/pull/14835))

BUG FIXES:
* provider: fixed an issue on multiple resources where non-retryable quota errors were considered retryable ([#14850](https://github.com/hashicorp/terraform-provider-google/pull/14850))
* vertexai: made `google_vertex_ai_featurestore_entitytype_feature` always use region corresponding to parent's region ([#14843](https://github.com/hashicorp/terraform-provider-google/pull/14843))

## 4.69.1 (June 12, 2023)

NOTE:
* Added a new user guide to the provider documentation ([#14886](https://github.com/hashicorp/terraform-provider-google/pull/14886))

## 4.69.0 (June 12, 2023)

FEATURES:
* **New Data Source:** `google_vmwareengine_network` ([#14821](https://github.com/hashicorp/terraform-provider-google/pull/14821))
* **New Resource:** `google_access_context_manager_service_perimeter_egress_policy` ([#14817](https://github.com/hashicorp/terraform-provider-google/pull/14817))
* **New Resource:** `google_access_context_manager_service_perimeter_ingress_policy` ([#14817](https://github.com/hashicorp/terraform-provider-google/pull/14817))
* **New Resource:** `google_certificate_manager_certificate_issuance_config` ([#14798](https://github.com/hashicorp/terraform-provider-google/pull/14798))
* **New Resource:** `google_dataplex_datascan` ([#14798](https://github.com/hashicorp/terraform-provider-google/pull/14798))
* **New Resource:** `google_dataplex_datascan_iam_*` ([#14828](https://github.com/hashicorp/terraform-provider-google/pull/14828))
* **New Resource:** `google_vmwareengine_network` ([#14821](https://github.com/hashicorp/terraform-provider-google/pull/14821))

IMPROVEMENTS:
* billing: added `lookup_projects` to `google_billing_account` datasource that skips reading the list of associated projects ([#14815](https://github.com/hashicorp/terraform-provider-google/pull/14815))
* dlp: added `info_type_transformations` block in the `record_transformations` field to `google_data_loss_prevention_deidentify_template` resource. ([#14827](https://github.com/hashicorp/terraform-provider-google/pull/14827))
* dlp: added `redact_config`, `fixed_size_bucketing_config`, `bucketing_config`, `time_part_config` and `date_shift_config`  fields to `google_data_loss_prevention_deidentify_template` resource ([#14797](https://github.com/hashicorp/terraform-provider-google/pull/14797))
* dlp: added `stored_info_type_id` field to `google_data_loss_prevention_stored_info_type` resource ([#14791](https://github.com/hashicorp/terraform-provider-google/pull/14791))
* dlp: added `template_id` field to `google_data_loss_prevention_deidentify_template` and `google_data_loss_prevention_inspect_template` ([#14823](https://github.com/hashicorp/terraform-provider-google/pull/14823))
* dlp: changed `actions` field from required to optional in `google_data_loss_prevention_job_trigger` resource ([#14803](https://github.com/hashicorp/terraform-provider-google/pull/14803))
* kms: removed validation for `purpose` in `google_kms_crypto_key` to allow newly added values for the field ([#14799](https://github.com/hashicorp/terraform-provider-google/pull/14799))
* pubsub: allowed `schema_settings` of `google_pubsub_topic` to change without deleting and recreating the resource ([#14819](https://github.com/hashicorp/terraform-provider-google/pull/14819))

BUG FIXES:
* tags: fixed providing `projects/<project_id` to `parent` causing recreation on `google_tags_tag_key` ([#14809](https://github.com/hashicorp/terraform-provider-google/pull/14809))

## 4.68.0 (June 5, 2023)

FEATURES:
* **New Resource:** `google_container_analysis_note_iam_*` ([#14706](https://github.com/hashicorp/terraform-provider-google/pull/14706))

IMPROVEMENTS:
* compute: promoted `allow_psc_global_access` field in `google_compute_forwarding_rule` to GA ([#14754](https://github.com/hashicorp/terraform-provider-google/pull/14754))
* dlp: added `included_fields` and `excluded_fields` fields to `google_data_loss_prevention_job_trigger` ([#14736](https://github.com/hashicorp/terraform-provider-google/pull/14736))
* dns: added `regionalL7ilb` enum support to the `routing_policy.load_balancer_type` field in `google_dns_record_set` ([#14710](https://github.com/hashicorp/terraform-provider-google/pull/14710))

BUG FIXES:
* accesscontextmanager: fixed incorrect validations for `spec` and `status` in `google_access_context_manager_service_perimeter` ([#14705](https://github.com/hashicorp/terraform-provider-google/pull/14705))
* alloydb: increased timeouts for `google_alloydb_instance` from 20m to 40m ([#14713](https://github.com/hashicorp/terraform-provider-google/pull/14713))
* apigee: fixed bug where updating `config_bundle` in `google_apigee_sharedflow` that's attached to `google_apigee_sharedflow_deployment` causes an error ([#14725](https://github.com/hashicorp/terraform-provider-google/pull/14725))
* compute: increased timeout for `compute_security_policy` from 4m to 8m ([#14712](https://github.com/hashicorp/terraform-provider-google/pull/14712))
* dataproc: fixed crash when reading `google_dataproc_cluster.virtual_cluster_config` ([#14744](https://github.com/hashicorp/terraform-provider-google/pull/14744))

## 4.67.0 (May 30, 2023)

FEATURES:
* **New Data Source:** `google_*_iam_policy` ([#14662](https://github.com/hashicorp/terraform-provider-google/pull/14662))
* **New Data Source:** `google_vertex_ai_index` ([#14640](https://github.com/hashicorp/terraform-provider-google/pull/14640))

IMPROVEMENTS:
* cloudrun: added `template.spec.containers.name` field to `google_cloud_run_service` ([#14647](https://github.com/hashicorp/terraform-provider-google/pull/14647))
* compute: added `network_performance_config` field to `google_compute_instance` and `google_compute_instance_template` ([#14678](https://github.com/hashicorp/terraform-provider-google/pull/14678))
* compute: added `guest_os_features` and `licenses` fields to `google_compute_disk` and `google_compute_region_disk` ([#14660](https://github.com/hashicorp/terraform-provider-google/pull/14660))
* datastream: added `mysql_source_config.max_concurrent_backfill_tasks` field to `google_datastream_stream` ([#14639](https://github.com/hashicorp/terraform-provider-google/pull/14639))
* firebase: added additional import formats for `google_firebase_webapp` ([#14638](https://github.com/hashicorp/terraform-provider-google/pull/14638))
* notebooks: added update support for `google_notebooks_instance.metadata` field ([#14650](https://github.com/hashicorp/terraform-provider-google/pull/14650))
* privateca: added `encoding_format` field to `google_privateca_ca_pool` ([#14663](https://github.com/hashicorp/terraform-provider-google/pull/14663))

BUG FIXES:
* apigee: increased `google_apigee_organization` timeout defaults to 45m from 20m ([#14643](https://github.com/hashicorp/terraform-provider-google/pull/14643))
* cloudresourcemanager: added retries to handle internal error: type: "googleapis.com" subject: "160009" ([#14727](https://github.com/hashicorp/terraform-provider-google/pull/14727))
* cloudrun: fixed a permadiff for `metadata.annotation` in `google_cloud_run_service` ([#14642](https://github.com/hashicorp/terraform-provider-google/pull/14642))
* container: fixed a crash scenario in `google_container_node_pool` ([#14693](https://github.com/hashicorp/terraform-provider-google/pull/14693))
* gkeonprem: changed `hostname` (under `ip_block`) from required to optional for `google_gkeonprem_vmware_cluster` ([#14690](https://github.com/hashicorp/terraform-provider-google/pull/14690))
* serviceusage: added retries to handle internal error: type: "googleapis.com" subject: "160009" when activating services ([#14727](https://github.com/hashicorp/terraform-provider-google/pull/14727))

## 4.66.0 (May 22, 2023)
NOTE:
* Upgraded to Go 1.19.9 ([#14561](https://github.com/hashicorp/terraform-provider-google/pull/14561))

FEATURES:
* **New Resource:** `google_network_security_server_tls_policy` ([#14557](https://github.com/hashicorp/terraform-provider-google/pull/14557))

IMPROVEMENTS:
* bigquery: added `ICEBERG` as an enum for `external_data_configuration.source_format` field in `google_bigquery_table` ([#14562](https://github.com/hashicorp/terraform-provider-google/pull/14562))
* cloudfunctions: added `status` attribute to the `google_cloudfunctions_function` resource and data source ([#14574](https://github.com/hashicorp/terraform-provider-google/pull/14574))
* compute: added `storage_location` field in `google_compute_image` resource ([#14619](https://github.com/hashicorp/terraform-provider-google/pull/14619))
* compute: added support for additional machine types in `google_compute_region_commitment` ([#14593](https://github.com/hashicorp/terraform-provider-google/pull/14593))
* monitoring: added `forecast_options` field to `google_monitoring_alert_policy` resource ([#14616](https://github.com/hashicorp/terraform-provider-google/pull/14616))
* monitoring: added `notification_channel_strategy` field to `google_monitoring_alert_policy` resource ([#14563](https://github.com/hashicorp/terraform-provider-google/pull/14563))
* sql: added `advanced_machine_features` field in `google_sql_database_instance` ([#14604](https://github.com/hashicorp/terraform-provider-google/pull/14604))
* storagetransfer: added field `path` to `transfer_spec.aws_s3_data_source` in `google_storage_transfer_job` ([#14610](https://github.com/hashicorp/terraform-provider-google/pull/14610))

BUG FIXES:
* artifactregistry: fixed new repositories ignoring the provider region if location is unset in `google_artifact_registry_repository`. ([#14596](https://github.com/hashicorp/terraform-provider-google/pull/14596))
* compute: fixed permadiff on `log_config.sample_rate` of `google_compute_backend_service` ([#14590](https://github.com/hashicorp/terraform-provider-google/pull/14590))
* container: fixed permadiff on `gateway_api_config.channel` of `google_container_cluster` ([#14576](https://github.com/hashicorp/terraform-provider-google/pull/14576))
* dataflow: fixed inconsistent final plan when labels are added to `google_dataflow_job` ([#14594](https://github.com/hashicorp/terraform-provider-google/pull/14594))
* provider: fixed an issue where mtls transports were not used consistently(initial implementation in v4.65.0, reverted in v4.65.1) ([#14621](https://github.com/hashicorp/terraform-provider-google/pull/14621))
* storage: fixed inconsistent final plan when labels are added to `google_storage_bucket` ([#14594](https://github.com/hashicorp/terraform-provider-google/pull/14594))

## 4.65.2 (May 16, 2023)

BUG FIXES:
* provider: fixed an issue where `google_client_config` datasource return `null` for all attributes when region or zone is unset in provider config

## 4.65.1 (May 15, 2023)

BUG FIXES:
* provider: fixed an issue where `google_client_config` datasource return `null` for `access_token`

## 4.65.0 (May 15, 2023)

FEATURES:
* **New Data Source:** `google_datastream_static_ips` ([#14487](https://github.com/hashicorp/terraform-provider-google/pull/14487))
* **New Resource:** `google_compute_disk_async_replication` ([#14489](https://github.com/hashicorp/terraform-provider-google/pull/14489))
* **New Resource:** `google_firestore_field` ([#14512](https://github.com/hashicorp/terraform-provider-google/pull/14512))

IMPROVEMENTS:
* bigquery: added general field `load.parquet_options` to `google_bigquery_job` ([#14497](https://github.com/hashicorp/terraform-provider-google/pull/14497))
* cloudbuild: added `allow_failure` and `allow_exit_codes` to `build.step` in `google_cloudbuild_trigger` resource ([#14498](https://github.com/hashicorp/terraform-provider-google/pull/14498))
* compute: added enumeration values `SEV_SNP_CAPABLE`, `SUSPEND_RESUME_COMPATIBLE`, `TDX_CAPABLE` for the `guest_os_features` of `google_compute_image` ([#14518](https://github.com/hashicorp/terraform-provider-google/pull/14518))
* compute: added support for `stack_type` to `google_compute_network_peering` ([#14509](https://github.com/hashicorp/terraform-provider-google/pull/14509))
* dlp: added `publish_to_stackdriver` field to `google_data_loss_prevention_job_trigger` resource ([#14539](https://github.com/hashicorp/terraform-provider-google/pull/14539))

BUG FIXES:
* certificatemanager: fixed an issue where `self_managed.pem_certificate` and `self_managed.pem_certificate` can't be updated on `google_certificate_manager_certificate` ([#14521](https://github.com/hashicorp/terraform-provider-google/pull/14521))
* compute: fixed crash on `terraform destroy -refresh=false` for instance group managers with `wait_for_instances = "true"` if the instance group manager was not found ([#14543](https://github.com/hashicorp/terraform-provider-google/pull/14543))
* container: fixed node auto-provisioning not working when `auto_provisioning_defaults.management` is not provided on `google_container_cluster` ([#14519](https://github.com/hashicorp/terraform-provider-google/pull/14519))
* provider: fixed an issue where mtls transports were not used consistently ([#14550](https://github.com/hashicorp/terraform-provider-google/pull/14550))

## 4.64.0 (May 8, 2023)

FEATURES:
* **New Data Source:** `google_alloydb_locations` ([#14355](https://github.com/hashicorp/terraform-provider-google/pull/14355))
* **New Data Source:** `google_sql_tiers` ([#14420](https://github.com/hashicorp/terraform-provider-google/pull/14420))
* **New Resource:** `google_database_migration_service_connection_profile` ([#14383](https://github.com/hashicorp/terraform-provider-google/pull/14383))

IMPROVEMENTS:
* alloydb: added `encryption_config` and `encryption_info` fields in `google_alloydb_cluster`, to allow CMEK encryption of the cluster's data. ([#14426](https://github.com/hashicorp/terraform-provider-google/pull/14426))
* alloydb: added support for CMEK in `google_alloydb_backup` resource ([#14421](https://github.com/hashicorp/terraform-provider-google/pull/14421))
* alloydb: added the `encryption_config` field inside the `automated_backup_policy` block in`google_alloydb_cluster`, to allow CMEK encryption of automated backups. ([#14426](https://github.com/hashicorp/terraform-provider-google/pull/14426))
* certificatemanager: added `location` field to `certificatemanager` certificate resource ([#14432](https://github.com/hashicorp/terraform-provider-google/pull/14432))
* cloudrun: promoted `startup_probe` and `liveness_probe` in resource `google_cloud_run_service` to GA. ([#14363](https://github.com/hashicorp/terraform-provider-google/pull/14363))
* cloudrunv2: added field `port` to `http_get` to resource `google_cloud_run_v2_service` ([#14358](https://github.com/hashicorp/terraform-provider-google/pull/14358))
* cloudrunv2: added field `startupCpuBoost` to resource `service` ([#14372](https://github.com/hashicorp/terraform-provider-google/pull/14372))
* cloudrunv2: added support for `session_affinity` to `google_cloud_run_v2_service` ([#14367](https://github.com/hashicorp/terraform-provider-google/pull/14367))
* compute: added `dest_fqdns`, `dest_region_codes`, `dest_threat_intelligences`, `src_fqdns`, `src_region_codes`, and `src_threat_intelligences` to `google_compute_firewall_policy_rule` resource. ([#14378](https://github.com/hashicorp/terraform-provider-google/pull/14378))
* compute: added `source_ip_ranges` and `base_forwarding_rule` to `google_compute_forwarding_rule` resource ([#14378](https://github.com/hashicorp/terraform-provider-google/pull/14378))
* compute: added `bypass_cache_on_request_headers` to `cdn_policy` in `google_compute_backend_service` resource ([#14446](https://github.com/hashicorp/terraform-provider-google/pull/14446))
* compute: added `dest_address_groups` and `src_address_groups` fields to `google_compute_firewall_policy_rule` and `google_compute_network_firewall_policy_rule` ([#14396](https://github.com/hashicorp/terraform-provider-google/pull/14396))
* compute: added new field `async_primary_disk` to `google_compute_disk` and `google_compute_region_disk` ([#14431](https://github.com/hashicorp/terraform-provider-google/pull/14431))
* compute: added new field `disk_consistency_group_policy` to `google_compute_resource_policy` ([#14431](https://github.com/hashicorp/terraform-provider-google/pull/14431))
* compute: added support for IPv6 prefix exchange in `google_compute_router_peer` ([#14397](https://github.com/hashicorp/terraform-provider-google/pull/14397))
* compute: made `network_firewall_policy_enforcement_order` field mutable in `google_compute_network`. ([#14364](https://github.com/hashicorp/terraform-provider-google/pull/14364))
* dlp: added `exclude_by_hotword` exclusion rule to `google_data_loss_prevention_inspect_template` resource ([#14433](https://github.com/hashicorp/terraform-provider-google/pull/14433))
* dlp: added `image_transformations` field to `google_data_loss_prevention_deidentify_template` resource ([#14434](https://github.com/hashicorp/terraform-provider-google/pull/14434))
* dlp: added `inspectConfig` field to `google_data_loss_prevention_job_trigger` resource ([#14401](https://github.com/hashicorp/terraform-provider-google/pull/14401))
* dlp: added `replace_dictionary_config` field to `info_type_transformations` in `google_data_loss_prevention_deidentify_template` resource ([#14434](https://github.com/hashicorp/terraform-provider-google/pull/14434))
* dlp: added `surrogate_type` custom type to `google_data_loss_prevention_inspect_template` resource ([#14433](https://github.com/hashicorp/terraform-provider-google/pull/14433))
* dlp: added `version` field for multiple `info_type` blocks to `google_data_loss_prevention_inspect_template` resource ([#14433](https://github.com/hashicorp/terraform-provider-google/pull/14433))
* gkehub: moved `google_gke_hub_feature` from beta to ga ([#14396](https://github.com/hashicorp/terraform-provider-google/pull/14396))
* sql: Added support for Postgres in `google_sql_source_representation_instance` ([#14436](https://github.com/hashicorp/terraform-provider-google/pull/14436))
* vertexai: added `region` field to `google_vertex_ai_endpoint` ([#14362](https://github.com/hashicorp/terraform-provider-google/pull/14362))
* workflows: added `crypto_key_name` field to `google_workflows_workflow` resource ([#14357](https://github.com/hashicorp/terraform-provider-google/pull/14357))

BUG FIXES:
* accesscontextmanager: fixed test for `google_access_context_manager_ingress_policy` ([#14361](https://github.com/hashicorp/terraform-provider-google/pull/14361))
* cloudplatform: added validation for `role_id` on `google_organization_iam_custom_role` ([#14454](https://github.com/hashicorp/terraform-provider-google/pull/14454))
* compute: fixed an import bug for `google_compute_router_interface` that happened when project was not set in the provider configuration or via environment variable ([#14356](https://github.com/hashicorp/terraform-provider-google/pull/14356))
* dns: fixed bug in `google_dns_keys` data source where list attributes could not be used at plan-time ([#14418](https://github.com/hashicorp/terraform-provider-google/pull/14418))
* firebase: specified required argument `bundle_id` in `google_firebase_apple_app` ([#14469](https://github.com/hashicorp/terraform-provider-google/pull/14469))

## 4.63.1 (April 26, 2023)

BUG FIXES:
* bigtable: fixed plan failure because of an unused zone being unavailable

## 4.63.0 (April 24, 2023)

NOTES:
* alloydb: changed `location` from `optional` to `required` for `google_alloydb_cluster` and `google_alloydb_backup` resources. `location` had previously been marked as optional, but operations failed if it was omitted, and there was no way for `location` to be inherited from the provider configuration or from an environment variable. This means there was no way to have a working configuration without `location` specified. ([#14330](https://github.com/hashicorp/terraform-provider-google/pull/14330), [#14334](https://github.com/hashicorp/terraform-provider-google/pull/14334))

FEATURES:
* **New Resource:** `google_access_context_manager_ingress_policy` ([#14302](https://github.com/hashicorp/terraform-provider-google/pull/14302))
* **New Resource:** `google_compute_public_advertised_prefix` ([#14303](https://github.com/hashicorp/terraform-provider-google/pull/14303))
* **New Resource:** `google_compute_public_delegated_prefix` ([#14303](https://github.com/hashicorp/terraform-provider-google/pull/14303))
* **New Resource:** `google_compute_region_commitment` ([#14301](https://github.com/hashicorp/terraform-provider-google/pull/14301))
* **New Resource:** `google_network_services_http_route` ([#14294](https://github.com/hashicorp/terraform-provider-google/pull/14294))

IMPROVEMENTS:
* dlp: added `inspect_job.actions.job_notification_emails` and `inspect_job.actions.deidentify`  fields to `google_data_loss_prevention_job_trigger` resource ([#14309](https://github.com/hashicorp/terraform-provider-google/pull/14309))
* dlp: added `triggers.manual` and `inspect_job.storage_config.hybrid_options` to `google_data_loss_prevention_job_trigger` ([#14326](https://github.com/hashicorp/terraform-provider-google/pull/14326))
* iam: added `oidc.web_sso_config` field to `google_iam_workforce_pool_provider` ([#14327](https://github.com/hashicorp/terraform-provider-google/pull/14327))

BUG FIXES:
* alloydb: changed `weekly_schedule` (under `automated_backup_policy`) from required to optional for `google_alloydb_cluster` ([#14335](https://github.com/hashicorp/terraform-provider-google/pull/14335))
* compute: fixed an issue with TTLs being sent when `USE_ORIGIN_HEADERS` is set in `google_compute_backend_bucket` ([#14323](https://github.com/hashicorp/terraform-provider-google/pull/14323))
* networkservices: increased default timeouts for `google_network_services_edge_cache_keyset` to 60m (from 30m) ([#14314](https://github.com/hashicorp/terraform-provider-google/pull/14314))
* sql: fixed an issue that prevented setting `enable_private_path_for_google_cloud_services` to `false` in `google_sql_database_instance` ([#14316](https://github.com/hashicorp/terraform-provider-google/pull/14316))

## 4.62.1 (April 19, 2023)

BUG FIXES:
* compute: fixed a diff that occurred when `stack_type` was unset on `google_compute_ha_vpn_gateway` ([#14311](https://github.com/hashicorp/terraform-provider-google/pull/14311))

## 4.62.0 (April 17, 2023)

FEATURES:
* **New Data Source:** `google_compute_region_instance_template` ([#14280](https://github.com/hashicorp/terraform-provider-google/pull/14280))
* **New Resource:** `google_compute_region_instance_template` ([#14280](https://github.com/hashicorp/terraform-provider-google/pull/14280))
* **New Resource:** `google_logging_linked_dataset` ([#14261](https://github.com/hashicorp/terraform-provider-google/pull/14261))

IMPROVEMENTS:
* cloudasset: added `OS_INVENTORY` value to `content_type` for `google_cloud_asset_*_feed` ([#14277](https://github.com/hashicorp/terraform-provider-google/pull/14277))
* clouddeploy: added canary deployment fields for resource `google_clouddeploy_delivery_pipeline` ([#14249](https://github.com/hashicorp/terraform-provider-google/pull/14249))
* compute: supported region instance template in`source_instance_template` field of `google_compute_instance_from_template` resource ([#14280](https://github.com/hashicorp/terraform-provider-google/pull/14280))
* container: added `pod_cidr_overprovision_config` field to `google_container_cluster` and  `google_container_node_pool` resources. ([#14281](https://github.com/hashicorp/terraform-provider-google/pull/14281))
* orgpolicy: accepted variable cases for booleans such as true, True, and TRUE in `google_org_policy_policy` ([#14240](https://github.com/hashicorp/terraform-provider-google/pull/14240))

BUG FIXES:
* cloudidentity: fixed immutability issue on `initialGroupConfig` field for resource `google_cloud_identity_group` ([#14257](https://github.com/hashicorp/terraform-provider-google/pull/14257))
* provider: fixed an error resulting from leaving `batching.send_after` unspecified and `batching` specified ([#14263](https://github.com/hashicorp/terraform-provider-google/pull/14263))
* provider: fixed bug where `credentials` field could not be set as an empty string ([#14279](https://github.com/hashicorp/terraform-provider-google/pull/14279))
* vertex: increased the default timeout for `google_vertex_ai_index` to 180m ([#14248](https://github.com/hashicorp/terraform-provider-google/pull/14248))

## 4.61.0 (April 10, 2023)

BREAKING CHANGES:
* cloudrunv2: set a default value of 3 for `max_retries` in `google_cloud_run_v2_job`. This should match the API's existing default, but may show a diff at plan time in limited circumstances as drift is now detected ([#14223](https://github.com/hashicorp/terraform-provider-google/pull/14223))

FEATURES:
* **New Data Source:** `google_firebase_android_app_config` ([#14202](https://github.com/hashicorp/terraform-provider-google/pull/14202))
* **New Resource:** `google_apigee_keystores_aliases_pkcs12` ([#14168](https://github.com/hashicorp/terraform-provider-google/pull/14168))
* **New Resource:** `google_apigee_keystores_aliases_self_signed_cert` ([#14140](https://github.com/hashicorp/terraform-provider-google/pull/14140))
* **New Resource:** `google_network_security_url_lists` ([#14232](https://github.com/hashicorp/terraform-provider-google/pull/14232))
* **New Resource:** `google_network_services_mesh` ([#14139](https://github.com/hashicorp/terraform-provider-google/pull/14139))

IMPROVEMENTS:
* alloydb: added update support for `initial_user` and `automated_backup_policy.weekly_schedule` to `google_alloydb_cluster` ([#14187](https://github.com/hashicorp/terraform-provider-google/pull/14187))
* artifactregistry: added support for tag immutability ([#14206](https://github.com/hashicorp/terraform-provider-google/pull/14206))
* artifactregistry: promoted `mode`, `virtual_repository_config`, and `remote_repository_config` to GA ([#14204](https://github.com/hashicorp/terraform-provider-google/pull/14204))
* bigqueryreservation: added `edition` and `autoscale` to `google_bigquery_reservation` and `edition` to `bigquery_capacity_commitment` ([#14148](https://github.com/hashicorp/terraform-provider-google/pull/14148))
* compute: added support for `SEV_LIVE_MIGRATABLE` to `guest_os_features.type` in `google_compute_image` ([#14200](https://github.com/hashicorp/terraform-provider-google/pull/14200))
* compute: added support for `stack_type` to `google_compute_ha_vpn_gateway` ([#14141](https://github.com/hashicorp/terraform-provider-google/pull/14141))
* container: added support for `ephemeral_storage_local_ssd_config` to `google_container_cluster.node_config`, `google_container_cluster.node_pools.node_config`, `google_container_node_pool.node_config` ([#14150](https://github.com/hashicorp/terraform-provider-google/pull/14150))
* dlp: Changed `dictionary`, `regex`, `regex.group_indexes` and `large_custom_dictionary` fields in `google_data_loss_prevention_stored_info_type` to be update-in-place ([#14207](https://github.com/hashicorp/terraform-provider-google/pull/14207))
* logging: added support for `disabled` to `google_logging_metric` ([#14198](https://github.com/hashicorp/terraform-provider-google/pull/14198))
* networkservices: increased the max count for `route_rule` to 200 on `google_network_services_edge_cache_service` ([#14224](https://github.com/hashicorp/terraform-provider-google/pull/14224))
* storagetransfer: added support for 'last_modified_since' and 'last_modified_before' fields to 'google_storage_transfer_job' resource ([#14147](https://github.com/hashicorp/terraform-provider-google/pull/14147))

BUG FIXES:
* bigquery: fixed the import logic in `google_bigquery_capacity_commitment` ([#14226](https://github.com/hashicorp/terraform-provider-google/pull/14226))
* cloudrunv2: fixed the bug where setting `max_retries` to 0 in `google_cloud_run_v2_job` was not respected. ([#14223](https://github.com/hashicorp/terraform-provider-google/pull/14223))
* container: fixed a bug creating a diff adding a `stack_type` when GKE omitted `stackType` in API responses from older GKE clusters ([#14208](https://github.com/hashicorp/terraform-provider-google/pull/14208))
* dataproc: fixed validation of `optional_components` ([#14167](https://github.com/hashicorp/terraform-provider-google/pull/14167))
* provider: fixed an issue where the `USER_PROJECT_OVERRIDE` environment variable was not being read ([#14238](https://github.com/hashicorp/terraform-provider-google/pull/14238))
* provider: fixed an issue where the provider crashed when "batching" was set in `4.60.0`/`4.60.1` ([#14235](https://github.com/hashicorp/terraform-provider-google/pull/14235))

## 4.60.2 (April 6, 2023)

BUG FIXES:
* provider: fixed an issue where the provider crashed when "batching" was set in `4.60.0`/`4.60.1`
* provider: fixed an issue where the `USER_PROJECT_OVERRIDE` environment variable was not being read

## 4.60.1 (April 5, 2023)

BUG FIXES:
* container: fixed a bug creating a diff adding a `stack_type` when GKE omitted `stackType` in API responses from older GKE clusters

## 4.60.0 (April 4, 2023)

FEATURES:
* **New Resource:** `google_apigee_keystores_aliases_key_cert_file` ([#14130](https://github.com/hashicorp/terraform-provider-google/pull/14130))

IMPROVEMENTS:
* compute: added `address_type`, `network`, `network_tier`, `prefix_length`, `purpose`, `subnetwork` and `users` field for `google_compute_address` and `google_compute_global_address` datasource ([#14078](https://github.com/hashicorp/terraform-provider-google/pull/14078))
* compute: added `network_firewall_policy_enforcement_order` field to `google_compute_network` resource ([#14111](https://github.com/hashicorp/terraform-provider-google/pull/14111))
* compute: added output-only attribute `self_link_unique` for `google_compute_instance_template` to point to the unique id of the resource instead of its name ([#14128](https://github.com/hashicorp/terraform-provider-google/pull/14128))
* container: added `stack_type` field to `google_container_cluster` resource ([#14079](https://github.com/hashicorp/terraform-provider-google/pull/14079))
* container: added `advanced_machine_features` field to `google_container_cluster` resource ([#14106](https://github.com/hashicorp/terraform-provider-google/pull/14106))
* networkservice: updated the max number of `host_rule` on `google_network_services_edge_cache_service` ([#14112](https://github.com/hashicorp/terraform-provider-google/pull/14112))
* sql: added support of single-database-recovery for SQL Server PITR with `database_names` attribute to `google_sql_instance` ([#14088](https://github.com/hashicorp/terraform-provider-google/pull/14088))

BUG FIXES:
* cloudrun: fixed race condition when polling for status during an update of a `google_cloud_run_service` ([#14087](https://github.com/hashicorp/terraform-provider-google/pull/14087))
* cloudsql: fixed the error in any subsequent apply on `google_sql_user` after its `google_sql_database_instance` is deleted ([#14098](https://github.com/hashicorp/terraform-provider-google/pull/14098))
* datacatalog: fixed `google_data_catalog_tag` only allowing 10 tags by increasing the page size to 1000 ([#14077](https://github.com/hashicorp/terraform-provider-google/pull/14077))
* firebase: fixed `google_firebase_project` to succeed on apply when the project already has firebase enabled ([#14121](https://github.com/hashicorp/terraform-provider-google/pull/14121))


## 4.59.0 (March 28, 2023)

FEATURES:
* **New Resource:** `google_dataplex_asset_iam_*` ([#14046](https://github.com/hashicorp/terraform-provider-google/pull/14046))
* **New Resource:** `google_dataplex_lake_iam_*` ([#14046](https://github.com/hashicorp/terraform-provider-google/pull/14046))
* **New Resource:** `google_dataplex_zone_iam_*` ([#14046](https://github.com/hashicorp/terraform-provider-google/pull/14046))
* **New Resource:** `google_network_services_gateway` ([#14057](https://github.com/hashicorp/terraform-provider-google/pull/14057))

IMPROVEMENTS:
* auth: added support for oauth2 token exchange over mTLS ([#14032](https://github.com/hashicorp/terraform-provider-google/pull/14032))
* bigquery: added `is_case_insensitive` and `default_collation` fields to `google_bigquery_dataset` resource ([#14031](https://github.com/hashicorp/terraform-provider-google/pull/14031))
* bigquerydatapolicy: promoted `google_bigquery_datapolicy_data_policy` to GA ([#13991](https://github.com/hashicorp/terraform-provider-google/pull/13991))
* compute: added `scratch_disk.size` field on `google_compute_instance` ([#14061](https://github.com/hashicorp/terraform-provider-google/pull/14061))
* compute: added 3000 as allowable value for `disk_size_gb` for SCRATCH disks in `google_compute_instance_template` ([#14061](https://github.com/hashicorp/terraform-provider-google/pull/14061))
* compute: added `WEIGHED_MAGLEV` to `locality_lb_policy` enum for backend service resources ([#14055](https://github.com/hashicorp/terraform-provider-google/pull/14055))
* container: added `local_nvme_ssd_block` to `node_config` block in the `google_container_node_pool` ([#14008](https://github.com/hashicorp/terraform-provider-google/pull/14008))
* logging: added `enable_analytics` field to `google_logging_project_bucket_config` ([#14043](https://github.com/hashicorp/terraform-provider-google/pull/14043))
* networkservices: updated max allowed items to 25 for `expose_headers`, `allow_headers`, `request_header_to_remove`, `request_header_to_add`, `response_header_to_add` and `response_header_to_remove` of `google_network_services_edge_cache_service` ([#14041](https://github.com/hashicorp/terraform-provider-google/pull/14041))
* networkservices: updated max allowed items to 25 for `request_headers_to_add` of `google_network_services_edge_cache_origin` ([#14041](https://github.com/hashicorp/terraform-provider-google/pull/14041))

BUG FIXES:
* certificatemanager: fixed `managed.dns_authorizations` not being included during import of `google_certificate_manager_certificate` ([#13992](https://github.com/hashicorp/terraform-provider-google/pull/13992))
* certificatemanager: fixed a bug where modifying non-updatable fields `hostname` and `matcher` in `google_certificate_manager_certificate_map_entry` would fail with API errors; now updating them will recreate the resource ([#13994](https://github.com/hashicorp/terraform-provider-google/pull/13994))
* compute: fixed bug where `enforce_on_key_name` could not be unset on `google_compute_security_policy` ([#13993](https://github.com/hashicorp/terraform-provider-google/pull/13993))
* datastream: fixed bug where field `dataset_id` could not utilize the id from bigquery directly ([#14003](https://github.com/hashicorp/terraform-provider-google/pull/14003))
* workstations: fixed permadiff on `service_account` of `google_workstations_workstation_config` ([#13989](https://github.com/hashicorp/terraform-provider-google/pull/13989))

## 4.58.0 (March 21, 2023)

FEATURES:
* **New Resource:** `google_apigee_sharedflow` ([#13938](https://github.com/hashicorp/terraform-provider-google/pull/13938))
* **New Resource:** `google_apigee_sharedflow_deployment` ([#13938](https://github.com/hashicorp/terraform-provider-google/pull/13938))
* **New Resource:** `google_apigee_flowhook` ([#13938](https://github.com/hashicorp/terraform-provider-google/pull/13938))

IMPROVEMENTS:
* datafusion: added support for `accelerators` field to `google_datafusion_instance` resource. ([#13946](https://github.com/hashicorp/terraform-provider-google/pull/13946))
* privateca: added support for X.509 name constraints to `google_privateca_pool`, `google_privateca_certificate`, and `google_privateca_certificate_authority` ([#13969](https://github.com/hashicorp/terraform-provider-google/pull/13969))

BUG FIXES:
* alloydb: fixed permadiff on `automated_backup_policy.weekly_schedule` of `google_alloydb_cluster` ([#13948](https://github.com/hashicorp/terraform-provider-google/pull/13948))
* bigquery: fixed a permadiff when `friendly_name` is removed from `google_bigquery_dataset` ([#13973](https://github.com/hashicorp/terraform-provider-google/pull/13973))
* redis: fixed a bug causing diff detection on `reserved_ip_range` in `google_redis_instance` ([#13958](https://github.com/hashicorp/terraform-provider-google/pull/13958))

## 4.57.0 (March 13, 2023)

FEATURES:
* **New Resource:** `google_access_context_manager_authorized_orgs_desc` ([#13925](https://github.com/hashicorp/terraform-provider-google/pull/13925))
* **New Resource:** `google_bigquery_capacity_commitment` ([#13902](https://github.com/hashicorp/terraform-provider-google/pull/13902))
* **New Resource:** `google_workstations_workstation` ([#13885](https://github.com/hashicorp/terraform-provider-google/pull/13885))
* **New Resource:** `google_apigee_env_keystore` ([#13876](https://github.com/hashicorp/terraform-provider-google/pull/13876))
* **New Resource:** `google_apigee_env_references` ([#13876](https://github.com/hashicorp/terraform-provider-google/pull/13876))
* **New Resource:** `google_firestore_database` ([#13874](https://github.com/hashicorp/terraform-provider-google/pull/13874))

BUG FIXES:
* cloudidentity: fixed an issue on `google_cloud_identity_group` `initial_group_config` field when importing ([#13875](https://github.com/hashicorp/terraform-provider-google/pull/13875))
* compute: fixed the error of invalid value for field `failover_policy` when UDP is selected on `google_compute_region_backend_service` ([#13897](https://github.com/hashicorp/terraform-provider-google/pull/13897))
* firebase: allowed specifying a `project` field on datasources for `google_firebase_android_app`, `google_firebase_web_app`, and `google_firebase_apple_app`. ([#13927](https://github.com/hashicorp/terraform-provider-google/pull/13927))
* tags: fixed a bug preventing use of `google_tags_location_tag_binding` with zonal parent resources ([#13880](https://github.com/hashicorp/terraform-provider-google/pull/13880))

## 4.56.0 (March 6, 2023)

FEATURES:
* **New Resource:** google_data_catalog_policy_tag ([#13818](https://github.com/hashicorp/terraform-provider-google/pull/13848))
* **New Resource:** google_data_catalog_taxonomy ([#13818](https://github.com/hashicorp/terraform-provider-google/pull/13848))
* **New Resource:** google_scc_mute_config ([#13818](https://github.com/hashicorp/terraform-provider-google/pull/13818))
* **New Resource:** google_workstations_workstation_config ([#13832](https://github.com/hashicorp/terraform-provider-google/pull/13832))

IMPROVEMENTS:
* cloudbuild: added `peered_network_ip_range` field to `google_cloudbuild_worker_pool` resource ([#13854](https://github.com/hashicorp/terraform-provider-google/pull/13854))
* cloudrun: added `template.0.containers0.liveness_probe.grpc`, `template.0.containers0.startup_probe.grpc` fields to `google_cloud_run_v2_service` resource ([#13855](https://github.com/hashicorp/terraform-provider-google/pull/13855))
* compute: added `max_distance` field to `resource-policy` resource ([#13853](https://github.com/hashicorp/terraform-provider-google/pull/13853))
* compute: added field `deletion_policy` to resource `google_compute_shared_vpc_service_project` ([#13822](https://github.com/hashicorp/terraform-provider-google/pull/13822))
* containerazure: added `azure_services_authentication` to `google_container_azure_cluster` ([#13854](https://github.com/hashicorp/terraform-provider-google/pull/13854))
* networkservices: increased maximum `allow_origins` from 5 to 25 on `network_services_edge_cache_service` ([#13808](https://github.com/hashicorp/terraform-provider-google/pull/13808))
* storagetransfer: added general field `sink_agent_pool_name` and `source_agent_pool_name` to `google_storage_transfer_job` ([#13865](https://github.com/hashicorp/terraform-provider-google/pull/13865))

BUG FIXES:
* cloudfunctions: fixed no diff found on `event_trigger.resource` of `google_cloudfunctions_function` ([#13862](https://github.com/hashicorp/terraform-provider-google/pull/13862))
* dataproc: fixed an issue where `master_config.num_instances` would not force recreation when changed in `google_dataproc_cluster` ([#13837](https://github.com/hashicorp/terraform-provider-google/pull/13837))
* spanner: fixed the error when updating `deletion_protection` on `google_spanner_database` ([#13821](https://github.com/hashicorp/terraform-provider-google/pull/13821))
* spanner: fixed the error when updating `force_destroy` on `google_spanner_instance` ([#13821](https://github.com/hashicorp/terraform-provider-google/pull/13821))

## 4.55.0 (February 27, 2023)

FEATURES:
* **New Resource:** `google_cloudbuild_bitbucket_server_config` ([#13767](https://github.com/hashicorp/terraform-provider-google/pull/13767))
* **New Resource:** `google_firebase_hosting_release` ([#13793](https://github.com/hashicorp/terraform-provider-google/pull/13793))
* **New Resource:** `google_firebase_hosting_version` ([#13793](https://github.com/hashicorp/terraform-provider-google/pull/13793))

IMPROVEMENTS:
* container: added support for `node_config.kubelet_config.pod_pids_limit` on `google_container_node_pool` ([#13762](https://github.com/hashicorp/terraform-provider-google/pull/13762))
* storage: changed the default create timeout of `google_storage_bucket` to 10m from 4m ([#13774](https://github.com/hashicorp/terraform-provider-google/pull/13774))

BUG FIXES:
* container: fixed a crash when leaving `placement_policy` blank on `google_container_node_pool` ([#13797](https://github.com/hashicorp/terraform-provider-google/pull/13797))

## 4.54.0 (February 22, 2023)

FEATURES:
* **New Data Source:** `google_firebase_hosting_channel` ([#13686](https://github.com/hashicorp/terraform-provider-google/pull/13686))
* **New Data Source:** `google_logging_sink` ([#13742](https://github.com/hashicorp/terraform-provider-google/pull/13742))
* **New Data Source:** `google_sql_databases` ([#13738](https://github.com/hashicorp/terraform-provider-google/pull/13738))

IMPROVEMENTS:
* cloudbuild: added `bitbucket_server_trigger_config` field to `google_cloudbuild_trigger` resource ([#13728](https://github.com/hashicorp/terraform-provider-google/pull/13728))
* cloudbuild: added `github.enterprise_config_resource_name` field to `google_cloudbuild_trigger` resource ([#13739](https://github.com/hashicorp/terraform-provider-google/pull/13739))
* compute: added field `rsa_encrypted_key` to `google_compute_disk` resource ([#13685](https://github.com/hashicorp/terraform-provider-google/pull/13685))
* sql: added replica promotion support to `google_sql_database_instance`. This change will allow users to promote read replica as stand alone primary instance. ([#13682](https://github.com/hashicorp/terraform-provider-google/pull/13682))

BUG FIXES:
* bigquery: fixed permadiff on `max_time_travel_hours` of `google_bigquery_dataset` ([#13691](https://github.com/hashicorp/terraform-provider-google/pull/13691))
* compute: added possibility to remove `stateful_disk` in `compute_instance_group_manager` and `compute_region_instance_group_manager`. ([#13737](https://github.com/hashicorp/terraform-provider-google/pull/13737))
* sql: fixed an issue with updating the `settings.activation_policy` field in `google_sql_database_instance`([#13736](https://github.com/hashicorp/terraform-provider-google/pull/13736))

## 4.53.1 (February 14, 2023)

BUG FIXES:
* provider: fixed crash when trying to configure the provider with invalid credentials

## 4.53.0 (February 13, 2023)

FEATURES:
* **New Resource:** `google_apigee_addons_config` ([#13654](https://github.com/hashicorp/terraform-provider-google/pull/13654))
* **New Resource:** `google_alloydb_backup` ([#13639](https://github.com/hashicorp/terraform-provider-google/pull/13639))
* **New Resource:** `google_alloydb_cluster` ([#13639](https://github.com/hashicorp/terraform-provider-google/pull/13639))
* **New Resource:** `google_alloydb_instance` ([#13639](https://github.com/hashicorp/terraform-provider-google/pull/13639))
* **New Resource:** `google_compute_region_target_tcp_proxy` ([#13640](https://github.com/hashicorp/terraform-provider-google/pull/13640))
* **New Resource:** `google_firestore_database` ([#13675](https://github.com/hashicorp/terraform-provider-google/pull/13675))
* **New Resource:** `google_workstations_workstation_cluster` ([#13619](https://github.com/hashicorp/terraform-provider-google/pull/13619))

IMPROVEMENTS:
* compute: added `resource_policies` field to `google_compute_instance_template` ([#13677](https://github.com/hashicorp/terraform-provider-google/pull/13677))
* compute: added the `labels` field to the `google_compute_external_vpn_gateway` resource ([#13642](https://github.com/hashicorp/terraform-provider-google/pull/13642))
* datastream: added `postgresql_source_config` & `oracle_source_config` in `google_datastream_stream` ([#13646](https://github.com/hashicorp/terraform-provider-google/pull/13646))
* datastream: added support for creating `google_datastream_stream` with `desired_state=RUNNING` ([#13646](https://github.com/hashicorp/terraform-provider-google/pull/13646))
* datastream: exposed validation errors during `google_datastream_stream` creation ([#13646](https://github.com/hashicorp/terraform-provider-google/pull/13646))
* firebase: marked `deletion_policy` as updatable without recreation on `google_firebase_android_app` and `google_firebase_apple_app` ([#13643](https://github.com/hashicorp/terraform-provider-google/pull/13643))
* sql: added `enable_private_path_for_google_cloud_services` field to `google_sql_database_instance` resource ([#13668](https://github.com/hashicorp/terraform-provider-google/pull/13668))
* vertex_ai: added the field `description` to `google_vertex_ai_featurestore_entitytype` ([#13641](https://github.com/hashicorp/terraform-provider-google/pull/13641))

BUG FIXES:
* composer: fixed an issue with cleaning up environments created in an error state ([#13644](https://github.com/hashicorp/terraform-provider-google/pull/13644))
* compute: fixed wrong maximum limit description for possible VPC MTUs ([#13674](https://github.com/hashicorp/terraform-provider-google/pull/13674))
* datafusion: fixed `version` can't be updated on `google_data_fusion_instance` ([#13658](https://github.com/hashicorp/terraform-provider-google/pull/13658))

## 4.52.0 (February 6, 2023)

FEATURES:
* **New Data Source:** `google_secret_manager_secret_version_access` ([#13605](https://github.com/hashicorp/terraform-provider-google/pull/13605))
* **New Resource:** `google_workstations_workstation_cluster` ([#13619](https://github.com/hashicorp/terraform-provider-google/pull/13619))

IMPROVEMENTS:
* bigquery: added support for federated Azure identities to BigQuery Omni connections. ([#13614](https://github.com/hashicorp/terraform-provider-google/pull/13614))
* bigquery: added `cloud_spanner.use_serverless_analytics` field ([#13588](https://github.com/hashicorp/terraform-provider-google/pull/13588))
* bigquery: added `cloud_sql.service_account_id` and `azure.identity` output fields ([#13588](https://github.com/hashicorp/terraform-provider-google/pull/13588))
* compute: added `locality_lb_policies` field to `google_compute_backend_service` ([#13604](https://github.com/hashicorp/terraform-provider-google/pull/13604))
* sql: updated the `settings.deletion_protection_enabled` property documentation. ([#13581](https://github.com/hashicorp/terraform-provider-google/pull/13581))
* sql: made `root_password` field updatable in `google_sql_database_instance` ([#13574](https://github.com/hashicorp/terraform-provider-google/pull/13574))

BUG FIXES:
* cloudfunctions: updated max_instances field to take API's result as default value ([#13575](https://github.com/hashicorp/terraform-provider-google/pull/13575))
* container: fixed an issue with resuming failed cluster creation ([#13580](https://github.com/hashicorp/terraform-provider-google/pull/13580))
* gke: fixed the error of Invalid address to set on `config_connector_config` of the data source `google_container_cluster` ([#13566](https://github.com/hashicorp/terraform-provider-google/pull/13566))
* secretmanager: fixed incorrect required_with for topics in `google_secret_managed_secret` ([#13612](https://github.com/hashicorp/terraform-provider-google/pull/13612))

## 4.51.0 (January 30, 2023)

DEPRECATIONS:
* cloudrunv2: deprecated `liveness_probe.tcp_socket` field from `google_cloud_run_v2_service` resource as it is not supported by the API and it will be removed in a future major release ([#13563](https://github.com/hashicorp/terraform-provider-google/pull/13563))
* cloudrunv2: deprecated `startup_probe` and `liveness_probe` fields from `google_cloud_run_v2_job` resource as they are not supported by the API and they will be removed in a future major release ([#13531](https://github.com/hashicorp/terraform-provider-google/pull/13531))

FEATURES:
* **New Resource:** `google_iam_access_boundary_policy` ([#13565](https://github.com/hashicorp/terraform-provider-google/pull/13565))
* **New Resource:** `google_tags_location_tag_bindings` ([#13524](https://github.com/hashicorp/terraform-provider-google/pull/13524))

IMPROVEMENTS:
* cloudbuild: added `github_enterprise_config` fields to `google_cloudbuild_trigger` resource. ([#13518](https://github.com/hashicorp/terraform-provider-google/pull/13518))
* cloudrunV2: added `annotations` to `google_cloud_run_v2_service` resource ([#13509](https://github.com/hashicorp/terraform-provider-google/pull/13509))
* compute:  added `tcp_time_wait_timeout_sec` field to `google_compute_router_nat` resource ([#13554](https://github.com/hashicorp/terraform-provider-google/pull/13554))
* compute: added `share_settings` field to the `google_compute_node_group` resource. ([#13522](https://github.com/hashicorp/terraform-provider-google/pull/13522))
* containerattached: added `deletion_policy` field to `google_container_attached_cluster` resource. ([#13551](https://github.com/hashicorp/terraform-provider-google/pull/13551))
* datastream: added `customer_managed_encryption_key` and `destination_config.bigquery_destination_config.source_hierarchy_datasets.dataset_template.kms_key_name` fields to `google_datastream_stream` resource ([#13549](https://github.com/hashicorp/terraform-provider-google/pull/13549))
* dlp: added `publish_findings_to_cloud_data_catalog` and `publish_summary_to_cscc` to `google_data_loss_prevention_job_trigger` resource ([#13562](https://github.com/hashicorp/terraform-provider-google/pull/13562))
* sql: added point_in_time_recovery_enabled for SQLServer in `google_sql_database_instance` ([#13555](https://github.com/hashicorp/terraform-provider-google/pull/13555))
* spanner: added support for IAM conditions with `google_spanner_database_iam_member` and `google_spanner_instance_iam_member` ([#13556](https://github.com/hashicorp/terraform-provider-google/pull/13556))
* sql: added additional fields to `google_sql_source_representation_instance` ([#13523](https://github.com/hashicorp/terraform-provider-google/pull/13523))

BUG FIXES:
* bigquery: fixed bug where valid iam member values for bigquery were prevented from actuation by validation ([#13520](https://github.com/hashicorp/terraform-provider-google/pull/13520))
* bigquery: fixed permadiff on `external_data_configuration.connection_id` of `google_bigquery_table` ([#13560](https://github.com/hashicorp/terraform-provider-google/pull/13560))
* gke: fixed the error of Invalid address to set on `config_connector_config` of the data source `google_container_cluster` ([#13566](https://github.com/hashicorp/terraform-provider-google/pull/13566))
* google_project: fixes misleading examples that could cause `firebase:enabled` label to be accidentally removed. ([#13552](https://github.com/hashicorp/terraform-provider-google/pull/13552))

## 4.50.0 (January 23, 2023)

FEATURES:
* **New Data Source:** `google_compute_network_peering` ([#13476](https://github.com/hashicorp/terraform-provider-google/pull/13476))
* **New Data Source:** `google_compute_router_nat` ([#13475](https://github.com/hashicorp/terraform-provider-google/pull/13475))
* **New Resource:** `google_cloud_run_v2_job_iam_binding` ([#13492](https://github.com/hashicorp/terraform-provider-google/pull/13492))
* **New Resource:** `google_cloud_run_v2_job_iam_member` ([#13492](https://github.com/hashicorp/terraform-provider-google/pull/13492))
* **New Resource:** `google_cloud_run_v2_job_iam_policy` ([#13492](https://github.com/hashicorp/terraform-provider-google/pull/13492))
* **New Resource:** `google_cloud_run_v2_service_iam_binding` ([#13492](https://github.com/hashicorp/terraform-provider-google/pull/13492))
* **New Resource:** `google_cloud_run_v2_service_iam_member` ([#13492](https://github.com/hashicorp/terraform-provider-google/pull/13492))
* **New Resource:** `google_cloud_run_v2_service_iam_policy` ([#13492](https://github.com/hashicorp/terraform-provider-google/pull/13492))
* **New Resource:** `google_gke_backup_backup_plan_iam_binding` ([#13508](https://github.com/hashicorp/terraform-provider-google/pull/13508))
* **New Resource:** `google_gke_backup_backup_plan_iam_member` ([#13508](https://github.com/hashicorp/terraform-provider-google/pull/13508))
* **New Resource:** `google_gke_backup_backup_plan_iam_policy` ([#13508](https://github.com/hashicorp/terraform-provider-google/pull/13508))

IMPROVEMENTS:
* bigquery_table - added `reference_file_schema_uri` ([#13493](https://github.com/hashicorp/terraform-provider-google/pull/13493))
* billingbudget: made fields `credit_types` and `subaccounts` updatable for `google_billing_budget` ([#13466](https://github.com/hashicorp/terraform-provider-google/pull/13466))
* cloudrunV2: added `annotations` to `CloudRunV2_service` resource ([#13509](https://github.com/hashicorp/terraform-provider-google/pull/13509))
* composer: added `recovery_config` in `google_composer_environment` resource ([#13504](https://github.com/hashicorp/terraform-provider-google/pull/13504))
* compute: added support for 'edge_security_policy' field to 'google_compute_backend_service' resource. ([#13494](https://github.com/hashicorp/terraform-provider-google/pull/13494))
* compute: added `max_run_duration` field to `google_compute_instance` and `google_compute_instance_template` resource (beta) ([#13489](https://github.com/hashicorp/terraform-provider-google/pull/13489))
* dataproc: added support for `dataproc_metric_config` to resource `google_dataproc_cluster` ([#13480](https://github.com/hashicorp/terraform-provider-google/pull/13480))
* dlp: added all subfields under `deidentify_template.record_transformations.field_transformations.primitive_transformation` to `google_data_loss_prevention_deidentify_template` ([#13498](https://github.com/hashicorp/terraform-provider-google/pull/13498))
* sql: changed the default create timeout of `google_sql_database_instance` to 40m from 30m ([#13481](https://github.com/hashicorp/terraform-provider-google/pull/13481))

BUG FIXES:
* certificatemanager: removed incorrect indication that the `self_managed` field in `google_certificate_manager_certificate` was treated as sensitive, and marked `self_managed.pem_private_key` as sensitive ([#13505](https://github.com/hashicorp/terraform-provider-google/pull/13505))
* cloudplatform: fixed the error with header `X-Goog-User-Project` on `google_client_openid_userinfo` ([#13474](https://github.com/hashicorp/terraform-provider-google/pull/13474))
* cloudsql: fixed `disk_type` can't be updated on `google_sql_database_instance` ([#13483](https://github.com/hashicorp/terraform-provider-google/pull/13483))
* vertexai: fixed updating value_type in google_vertex_ai_featurestore_entitytype_feature ([#13491](https://github.com/hashicorp/terraform-provider-google/pull/13491))

## 4.49.0 (January 17, 2023)

FEATURES:
* **New Data Source:** `google_project_service` ([#13434](https://github.com/hashicorp/terraform-provider-google/pull/13434))
* **New Data Source:** `google_sql_database_instances` ([#13433](https://github.com/hashicorp/terraform-provider-google/pull/13433))
* **New Data Source:** `google_container_attached_install_manifest` ([#13443](https://github.com/hashicorp/terraform-provider-google/pull/13443))
* **New Data Source:** `google_container_attached_install_manifest` ([#13455](https://github.com/hashicorp/terraform-provider-google/pull/13455))
* **New Data Source:** `google_container_attached_versions` ([#13443](https://github.com/hashicorp/terraform-provider-google/pull/13443))
* **New Resource:** `google_datastream_stream` ([#13385](https://github.com/hashicorp/terraform-provider-google/pull/13385))

IMPROVEMENTS:
* android_app: added general fields `sha1_hashes`, `sha256_hashes` and `etag` to `google_firebase_android_app`. ([#13444](https://github.com/hashicorp/terraform-provider-google/pull/13444))
* cloudids: added `threat_exception` field to `google_cloud_ids_endpoint` resource ([#13442](https://github.com/hashicorp/terraform-provider-google/pull/13442))
* compute: added deletion for `statefulIps` fields in `instance_group_manager` and `region_instance_group_manager`. ([#13428](https://github.com/hashicorp/terraform-provider-google/pull/13428))
* compute: added field `expire_time` to resource `google_compute_region_ssl_certificate` ([#13392](https://github.com/hashicorp/terraform-provider-google/pull/13392))
* compute: added field `expire_time` to resource `google_compute_ssl_certificate` ([#13392](https://github.com/hashicorp/terraform-provider-google/pull/13392))
* container: added `release_channel_latest_version` in `google_container_engine_versions` datasource ([#13384](https://github.com/hashicorp/terraform-provider-google/pull/13384))
* container: added `google_container_aws_node_pool` `autoscaling_metrics_collection` field ([#13462](https://github.com/hashicorp/terraform-provider-google/pull/13462))
* container: added update support for `google_container_aws_node_pool` `tags` field ([#13462](https://github.com/hashicorp/terraform-provider-google/pull/13462))
* container: added `config_connector_config` addon field to `google_container_cluster` ([#13380](https://github.com/hashicorp/terraform-provider-google/pull/13380))
* container: added `kubelet_config` field to `google_container_node_pool` ([#13423](https://github.com/hashicorp/terraform-provider-google/pull/13423))
* dataproc: added support for `node_group_affinity.` in `google_dataproc_cluster` ([#13400](https://github.com/hashicorp/terraform-provider-google/pull/13400))
* dataproc: added support for `reservation_affinity` in `google_dataproc_cluster` ([#13393](https://github.com/hashicorp/terraform-provider-google/pull/13393))
* dlp: added field `identifying_fields` to `big_query_options` for creating DLP jobs. ([#13463](https://github.com/hashicorp/terraform-provider-google/pull/13463))
* metastore: added `telemetry_config` field to `google_dataproc_metastore_service` ([#13432](https://github.com/hashicorp/terraform-provider-google/pull/13432))
* sql: added the ability to set `point_in_time_recovery_enabled` flag for `google_sql_database_instance` `SQLSERVER` instances ([#13454](https://github.com/hashicorp/terraform-provider-google/pull/13454))
* sql: added `instance_type` field to `google_sql_database_instance` resource ([#13406](https://github.com/hashicorp/terraform-provider-google/pull/13406))
* vertexai: added `scaling` field in `google_vertex_ai_featurestore` ([#13458](https://github.com/hashicorp/terraform-provider-google/pull/13458))

BUG FIXES:
* android_app: modified the `package_name` field suffix to always start with a letter in `google_firebase_android_app`. ([#13444](https://github.com/hashicorp/terraform-provider-google/pull/13444))
* bigqueryconnection: fixed a bug where `aws.access_role.iam_role_id` cannot be updated on `google_bigquery_connection` ([#13460](https://github.com/hashicorp/terraform-provider-google/pull/13460))
* cloudplatform: fixed a bug where `google_folder` deletion would fail to handle async operations ([#13377](https://github.com/hashicorp/terraform-provider-google/pull/13377))
* container: fixed a bug preventing updates to `master_global_access_config` in `google_container_cluster` ([#13383](https://github.com/hashicorp/terraform-provider-google/pull/13383))
* spanner: fixed crash when `google_spanner_database.ddl` item was nil ([#13441](https://github.com/hashicorp/terraform-provider-google/pull/13441))

## 4.48.0 (January 9, 2023)

FEATURES:
* **New Data Source:** `google_beyondcorp_app_connection` ([#13336](https://github.com/hashicorp/terraform-provider-google/pull/13336))
* **New Data Source:** `google_beyondcorp_app_connector` ([#13305](https://github.com/hashicorp/terraform-provider-google/pull/13305))
* **New Data Source:** `google_beyondcorp_app_gateway` ([#13305](https://github.com/hashicorp/terraform-provider-google/pull/13305))
* **New Data Source:** `google_cloudbuild_trigger` ([#13329](https://github.com/hashicorp/terraform-provider-google/pull/13329))
* **New Data Source:** `google_compute_instance_group_manager` ([#13297](https://github.com/hashicorp/terraform-provider-google/pull/13297))
* **New Data Source:** `google_firebase_apple_app` ([#13239](https://github.com/hashicorp/terraform-provider-google/pull/13239))
* **New Data Source:** `google_pubsub_subscription` ([#13296](https://github.com/hashicorp/terraform-provider-google/pull/13296))
* **New Data Source:** `google_sql_database` ([#13376](https://github.com/hashicorp/terraform-provider-google/pull/13376))
* **New Resource:** `google_apigee_sync_authorization` ([#13324](https://github.com/hashicorp/terraform-provider-google/pull/13324))
* **New Resource:** `google_beyondcorp_app_connection` ([#13318](https://github.com/hashicorp/terraform-provider-google/pull/13318))
* **New Resource:** `google_container_attached_cluster` ([#13374](https://github.com/hashicorp/terraform-provider-google/pull/13374))
* **New Resource:** `google_dns_managed_zone_iam_*` ([#13304](https://github.com/hashicorp/terraform-provider-google/pull/13304))
* **New Resource:** `google_gke_backup_backup_plan` ([#13359](https://github.com/hashicorp/terraform-provider-google/pull/13359))
* **New Resource:** `google_iam_workforce_pool_provider` ([#13299](https://github.com/hashicorp/terraform-provider-google/pull/13299))
* **New Resource:** `google_iam_workforce_pool` ([#13299](https://github.com/hashicorp/terraform-provider-google/pull/13299))

IMPROVEMENTS:
* cloudfunctions2: added `available_cpu` and `max_instance_request_concurrency` to support concurrency in `google_cloudfunctions2_function` resource ([#13315](https://github.com/hashicorp/terraform-provider-google/pull/13315))
* compute: added support for local IP ranges in `google_compute_firewall` ([#13240](https://github.com/hashicorp/terraform-provider-google/pull/13240))
* compute: added `router_appliance_instance` field to `google_compute_router_bgp_peer` ([#13373](https://github.com/hashicorp/terraform-provider-google/pull/13373))
* compute: added support for `generated_id` field in `google_compute_backend_service` to get the value of `id` defined by the server ([#13242](https://github.com/hashicorp/terraform-provider-google/pull/13242))
* compute: added support for `image_encryption_key` to `google_compute_image` ([#13253](https://github.com/hashicorp/terraform-provider-google/pull/13253))
* compute: added support for `source_snapshot`, `source_snapshot_encyption_key`, and `source_image_encryption_key` to `google_compute_instance_template` ([#13253](https://github.com/hashicorp/terraform-provider-google/pull/13253))
* container: promoted `google_container_node_pool.placement_policy` to GA ([#13372](https://github.com/hashicorp/terraform-provider-google/pull/13372))
* container: added `gateway_api_config` block to `google_container_cluster` resource for supporting the gke gateway api controller ([#13233](https://github.com/hashicorp/terraform-provider-google/pull/13233))
* container: supported in-place update for `labels` in `google_container_node_pool` ([#13284](https://github.com/hashicorp/terraform-provider-google/pull/13284))
* dataproc: added support for `SPOT` option for `preemptibility` in `google_dataproc_cluster` ([#13335](https://github.com/hashicorp/terraform-provider-google/pull/13335))
* dlp: added field `deidentify_config.record_transformations.field_transformations` to `google_data_loss_prevention_deidentify_template` ([#13282](https://github.com/hashicorp/terraform-provider-google/pull/13282))
* dlp: added field `deidentify_config.record_transformations.record_suppressions` to `google_data_loss_prevention_deidentify_template` ([#13300](https://github.com/hashicorp/terraform-provider-google/pull/13300))
* dlp: added `version` field to `google_data_loss_prevention_inspect_template` resource ([#13366](https://github.com/hashicorp/terraform-provider-google/pull/13366))
* osconfig: added support for `skip_await_rollout` in `google_os_config_os_policy_assignment` ([#13340](https://github.com/hashicorp/terraform-provider-google/pull/13340))
* sql: added [new deletion protection](https://cloud.google.com/sql/docs/mysql/deletion-protection) feature `deletion_protection_enabled` in `google_sql_database_instance` to guard against deletion from all surfaces ([#13249](https://github.com/hashicorp/terraform-provider-google/pull/13249))
* sql: made `settings.sql_server_audit_config.bucket` field in `google_sql_database_instance` to be optional. ([#13252](https://github.com/hashicorp/terraform-provider-google/pull/13252))
* storagetransfer: supported in-place update for `schedule` in `google_storage_transfer_job` ([#13262](https://github.com/hashicorp/terraform-provider-google/pull/13262))

BUG FIXES:
* bigquery: fixed a permadiff on `labels` of `google_bigquery_dataset` when it is referenced in `google_dataplex_asset` ([#13333](https://github.com/hashicorp/terraform-provider-google/pull/13333))
* compute: fixed a permadiff on `private_ip_google_access` of `google_compute_subnetwork` ([#13244](https://github.com/hashicorp/terraform-provider-google/pull/13244))
* compute: fixed an issue where `enable_dynamic_port_allocation` was not able to set to `false` in `google_compute_router_nat` ([#13243](https://github.com/hashicorp/terraform-provider-google/pull/13243))
* container: fixed a permadiff on `location_policy` of `google_container_cluster` and `google_container_node_pool` ([#13283](https://github.com/hashicorp/terraform-provider-google/pull/13283))
* identityplatform: fixed issues with `google_identity_platform_config` creation ([#13301](https://github.com/hashicorp/terraform-provider-google/pull/13301))
* resourcemanager: fixed the `google_project` datasource silently returning empty results when the project was not found or not in the ACTIVE state. Now, an error will be surfaced instead. ([#13358](https://github.com/hashicorp/terraform-provider-google/pull/13358))
* sql: fixed `sql_database_instance` leaking root users ([#13258](https://github.com/hashicorp/terraform-provider-google/pull/13258))

## 4.47.0 (December 21, 2022)

NOTES:
* sql: fixed an issue where `google_sql_database` was abandoned by default as of version `4.45.0`. Users who have upgraded to `4.45.0` or `4.46.0` will see a diff when running their next `terraform apply` after upgrading this version, indicating the `deletion_policy` field's value has changed from `"ABANDON"` to `"DELETE"`. This will create a no-op call against the API, but can otherwise be safely applied. ([#13226](https://github.com/hashicorp/terraform-provider-google/pull/13226))

FEATURES:
* **New Resource:** `google_alloydb_backup` ([#13202](https://github.com/hashicorp/terraform-provider-google/pull/13202))
* **New Resource:** `google_filestore_backup` ([#13209](https://github.com/hashicorp/terraform-provider-google/pull/13209))

IMPROVEMENTS:
* bigtable: added `deletion_protection` field to `google_bigtable_table` ([#13232](https://github.com/hashicorp/terraform-provider-google/pull/13232))
* compute: made `google_compute_subnetwork.ipv6_access_type` field updatable in-place ([#13211](https://github.com/hashicorp/terraform-provider-google/pull/13211))
* container: added `auto_provisioning_defaults.cluster_autoscaling.upgrade_settings` in `google_container_cluster` ([#13199](https://github.com/hashicorp/terraform-provider-google/pull/13199))
* container: added `gateway_api_config` block to `google_container_cluster` resource for supporting the gke gateway api controller ([#13233](https://github.com/hashicorp/terraform-provider-google/pull/13233))
* container: promoted `gke_backup_agent_config` in `google_container_cluster` to GA ([#13223](https://github.com/hashicorp/terraform-provider-google/pull/13223))
* container: promoted `min_cpu_platform` in `google_container_cluster` to GA ([#13199](https://github.com/hashicorp/terraform-provider-google/pull/13199))
* datacatalog: added update support for `fields` in `google_data_catalog_tag_template` ([#13216](https://github.com/hashicorp/terraform-provider-google/pull/13216))
* iam: Added plan-time validation for IAM members ([#13203](https://github.com/hashicorp/terraform-provider-google/pull/13203))
* logging: added `bucket_name` field to `google_logging_metric` ([#13210](https://github.com/hashicorp/terraform-provider-google/pull/13210))
* logging: made `metric_descriptor` field optional for `google_logging_metric` ([#13225](https://github.com/hashicorp/terraform-provider-google/pull/13225))

BUG FIXES:
* composer: fixed a crash when updating `ip_allocation_policy` of `google_composer_environment` ([#13188](https://github.com/hashicorp/terraform-provider-google/pull/13188))
* sql: fixed an issue where `google_sql_database` was abandoned by default as of version `4.45.0`. Users who have upgraded to `4.45.0` or `4.46.0` will see a diff when running their next `terraform apply` after upgrading this version, indicating the `deletion_policy` field's value has changed from `"ABANDON"` to `"DELETE"`. This will create a no-op call against the API, but can otherwise be safely applied. ([#13226](https://github.com/hashicorp/terraform-provider-google/pull/13226))

## 4.46.0 (December 12, 2022)

FEATURES:
* **New Data Source:** `google_firebase_android_app` ([#13186](https://github.com/hashicorp/terraform-provider-google/pull/13186))
* **New Resource:** `google_cloud_run_v2_job` ([#13154](https://github.com/hashicorp/terraform-provider-google/pull/13154))
* **New Resource:** `google_cloud_run_v2_service` ([#13166](https://github.com/hashicorp/terraform-provider-google/pull/13166))
* **New Resource:** `google_gke_backup_backup_plan` (beta) ([#13176](https://github.com/hashicorp/terraform-provider-google/pull/13176))
* **New Resource:** google_firebase_storage_bucket ([#13183](https://github.com/hashicorp/terraform-provider-google/pull/13183))

IMPROVEMENTS:
* network_services: added `origin_override_action` and `origin_redirect` to `google_network_services_edge_cache_origin` ([#13153](https://github.com/hashicorp/terraform-provider-google/pull/13153))
* bigquerydatatransfer: recreate `google_bigquery_data_transfer_config` for Cloud Storage transfers when immutable params `data_path_template` and `destination_table_name_template` are changed ([#13137](https://github.com/hashicorp/terraform-provider-google/pull/13137))
* compute: Added fields to resource `google_compute_security_policy` to support Cloud Armor bot management ([#13159](https://github.com/hashicorp/terraform-provider-google/pull/13159))
* container: Added support for concurrent node pool mutations on a cluster. Previously, node pool mutations were restricted to run synchronously clientside. NOTE: While this feature is supported in Terraform from this release onwards, only a limited number of GCP projects will support this behavior initially. The provider will automatically process mutations concurrently as the feature rolls out generally. ([#13173](https://github.com/hashicorp/terraform-provider-google/pull/13173))
* container: promoted `managed_prometheus` field in `google_container_cluster` to GA ([#13150](https://github.com/hashicorp/terraform-provider-google/pull/13150))
* metastore: added general field `network_config` to `google_dataproc_metastore_service` ([#13184](https://github.com/hashicorp/terraform-provider-google/pull/13184))
* storage: added support for `autoclass` in `google_storage_bucket` resource ([#13185](https://github.com/hashicorp/terraform-provider-google/pull/13185))

BUG FIXES:
* alloydb: made `machine_config.cpu_count` updatable on `google_alloydb_instance` ([#13144](https://github.com/hashicorp/terraform-provider-google/pull/13144))
* composer: fixed a crash when updating `ip_allocation_policy` of `google_composer_environment` ([#13188](https://github.com/hashicorp/terraform-provider-google/pull/13188))
* container: fixed GKE permadiff/thrashing when `update_settings. max_surge` or `update_settings. max_unavailable` values are updating on `google_container_node_pool` ([#13171](https://github.com/hashicorp/terraform-provider-google/pull/13171))
* datastream: fixed `google_datastream_private_connection` ignoring failures during creation ([#13160](https://github.com/hashicorp/terraform-provider-google/pull/13160))
* kms: fixed issues with deleting crypto key versions in states other than ENABLED ([#13167](https://github.com/hashicorp/terraform-provider-google/pull/13167))

## 4.45.0 (December 5, 2022)

FEATURES:
* **New Data Source:** `google_logging_project_cmek_settings` ([#13078](https://github.com/hashicorp/terraform-provider-google/pull/13078))
* **New Resource:** `google_vertex_ai_tensorboard` ([#13065](https://github.com/hashicorp/terraform-provider-google/pull/13065))
* **New Resource:** `google_data_fusion_instance_iam_binding` ([#13134](https://github.com/hashicorp/terraform-provider-google/pull/13134))
* **New Resource:** `google_data_fusion_instance_iam_member` ([#13134](https://github.com/hashicorp/terraform-provider-google/pull/13134))
* **New Resource:** `google_data_fusion_instance_iam_policy` ([#13134](https://github.com/hashicorp/terraform-provider-google/pull/13134))
* **New Resource:** `google_eventarc_google_channel_config` ([#13080](https://github.com/hashicorp/terraform-provider-google/pull/13080))
* **New Resource:** `google_vertex_ai_index` ([#13132](https://github.com/hashicorp/terraform-provider-google/pull/13132))

IMPROVEMENTS:
* bigquerydatatransfer: forced recreation on `google_bigquery_data_transfer_config` for Cloud Storage transfers when immutable params `data_path_template` and `destination_table_name_template` are changed ([#13137](https://github.com/hashicorp/terraform-provider-google/pull/13137))
* bigtable: added support for abandoning GC policy ([#13066](https://github.com/hashicorp/terraform-provider-google/pull/13066))
* cloudsql: added `connector_enforcement` field to `google_sql_database_instance` resource ([#13059](https://github.com/hashicorp/terraform-provider-google/pull/13059))
* compute: added `default_route_action.cors_policy` field to `google_compute_region_url_map` resource ([#13063](https://github.com/hashicorp/terraform-provider-google/pull/13063))
* compute: added `default_route_action.fault_injection_policy` field to `google_compute_region_url_map` resource ([#13063](https://github.com/hashicorp/terraform-provider-google/pull/13063))
* compute: added `default_route_action.timeout` field to `google_compute_region_url_map` resource ([#13063](https://github.com/hashicorp/terraform-provider-google/pull/13063))
* compute: added `default_route_action.url_rewrite` field to `google_compute_region_url_map` resource ([#13063](https://github.com/hashicorp/terraform-provider-google/pull/13063))
* compute: added `include_http_headers` field to the `cdn_policy` field of `google_compute_backend_service` resource ([#13093](https://github.com/hashicorp/terraform-provider-google/pull/13093))
* compute: added field `list_managed_instances_results` to `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#13079](https://github.com/hashicorp/terraform-provider-google/pull/13079))
* compute: added subnetwork and private_ip_address arguments to resource_compute_router_interface ([#13105](https://github.com/hashicorp/terraform-provider-google/pull/13105))
* container: added `resource_labels` field to `node_config` resource ([#13104](https://github.com/hashicorp/terraform-provider-google/pull/13104))
* container: added field `enable_private_nodes` in `network_config` to `google_container_node_pool` ([#13128](https://github.com/hashicorp/terraform-provider-google/pull/13128))
* container: added field `gcp_public_cidrs_access_enabled` and `private_endpoint_subnetwork` to `google_container_cluster` ([#13128](https://github.com/hashicorp/terraform-provider-google/pull/13128))
* container: added update support for `enable_private_endpoint` and `enable_private_nodes` in `google_container_cluster` ([#13128](https://github.com/hashicorp/terraform-provider-google/pull/13128))
* container: promoted `network_config` in `google_container_node_pool` to GA. ([#13128](https://github.com/hashicorp/terraform-provider-google/pull/13128))
* datafusion: added `api_endpoint` and `p4_service_account ` attributes to `google_data_fusion_instance` ([#13134](https://github.com/hashicorp/terraform-provider-google/pull/13134))
* datafusion: added `zone`, `display_name`, `crypto_key_config`, `event_publish_config`, and `enable_rbac` args to `google_data_fusion_instance` ([#13134](https://github.com/hashicorp/terraform-provider-google/pull/13134))
* logging: added `cmek_settings` field to `google_logging_project_bucket_config` resource ([#13078](https://github.com/hashicorp/terraform-provider-google/pull/13078))
* sql: added 'deny_maintenance_period' field for 'google_sql_database_instance' within which 'end_date', 'start_date' and 'time' fields are present. ([#13106](https://github.com/hashicorp/terraform-provider-google/pull/13106))
* sql: added field `deletion_policy` to resource `google_sql_database` ([#13107](https://github.com/hashicorp/terraform-provider-google/pull/13107))

BUG FIXES:
* compute: fixed a crash with `google_compute_instance_template` on a newly released field when `advanced_machine_features` was set ([#13108](https://github.com/hashicorp/terraform-provider-google/pull/13108))
* compute: fixed a failure in updating `most_disruptive_allowed_action` on `google_compute_per_instance_config` and `google_compute_region_per_instance_config` ([#13067](https://github.com/hashicorp/terraform-provider-google/pull/13067))
* compute: fixed the error when `metadata` and `machine_type` are updated while `metadata_startup_script` was already provided on `google_compute_instance` ([#13077](https://github.com/hashicorp/terraform-provider-google/pull/13077))
* container: fixed the inability to update `authenticator_groups_config` on `google_container_cluster` ([#13111](https://github.com/hashicorp/terraform-provider-google/pull/13111))
* container: fixed the data source `google_container_cluster` to return an error if it does not exist ([#13070](https://github.com/hashicorp/terraform-provider-google/pull/13070))
* sql: fixed `googe_sql_database_instance` to include `backup_configuration` in initial create request ([#13092](https://github.com/hashicorp/terraform-provider-google/pull/13092))
* storage: fixed permdiff when `website`, `website.main_page_suffix`, `website.not_found_page` are removed on `google_storage_bucket` ([#13069](https://github.com/hashicorp/terraform-provider-google/pull/13069))
## 4.44.1 (November 22, 2022)

BUG FIXES:
* compute: fixed a crash with `google_compute_instance_template` on a newly released field when `advanced_machine_features` was set ([#13108](https://github.com/hashicorp/terraform-provider-google/pull/13108))

## 4.44.0 (November 21, 2022)

FEATURES:
* **New Resource:** `google_alloydb_instance` ([#12981](https://github.com/hashicorp/terraform-provider-google/pull/12981))
* **New Resource:** `google_beyondcorp_app_connector` ([#13011](https://github.com/hashicorp/terraform-provider-google/pull/13011))
* **New Resource:** `google_beyondcorp_app_gateway` ([#13011](https://github.com/hashicorp/terraform-provider-google/pull/13011))
* **New Resource:** `google_compute_network_firewall_policy_association` ([#13013](https://github.com/hashicorp/terraform-provider-google/pull/13013))
* **New Resource:** `google_compute_network_firewall_policy_rule` ([#13031](https://github.com/hashicorp/terraform-provider-google/pull/13031))
* **New Resource:** `google_compute_network_firewall_policy` ([#12969](https://github.com/hashicorp/terraform-provider-google/pull/12969))
* **New Resource:** `google_compute_region_network_firewall_policy_association` ([#13013](https://github.com/hashicorp/terraform-provider-google/pull/13013))
* **New Resource:** `google_compute_region_network_firewall_policy_rule` ([#13031](https://github.com/hashicorp/terraform-provider-google/pull/13031))
* **New Resource:** `google_compute_region_network_firewall_policy` ([#12969](https://github.com/hashicorp/terraform-provider-google/pull/12969))
* **New Resource:** `google_eventarc_channel` ([#13021](https://github.com/hashicorp/terraform-provider-google/pull/13021))
* **New Resource:** `google_firebase_apple_app` ([#13047](https://github.com/hashicorp/terraform-provider-google/pull/13047))
* **New Resource:** `google_firebase_hosting_channel` ([#13053](https://github.com/hashicorp/terraform-provider-google/pull/13053))
* **New Resource:** `google_firebase_hosting_site` ([#12960](https://github.com/hashicorp/terraform-provider-google/pull/12960))
* **New Resource:** `google_kms_crypto_key_versions` ([#12926](https://github.com/hashicorp/terraform-provider-google/pull/12926))
* **New Resource:** `google_storage_transfer_agent_pool` ([#12945](https://github.com/hashicorp/terraform-provider-google/pull/12945))
* **New Resource:** `google_identity_platform_project_default_config` ([#12977](https://github.com/hashicorp/terraform-provider-google/pull/12977))

IMPROVEMENTS:
* bigquery: supported authorized routines on resource `bigquery_dataset` and `bigquery_dataset_access` ([#12979](https://github.com/hashicorp/terraform-provider-google/pull/12979))
* cloudidentity: made security label settable by making labels updatable in `google_cloud_identity_groups` ([#12943](https://github.com/hashicorp/terraform-provider-google/pull/12943))
* cloudsql: added `connector_enforcement` field to `google_sql_database_instance` resource ([#13059](https://github.com/hashicorp/terraform-provider-google/pull/13059))
* compute: added optional `redundant_interface` argument to `google_compute_router_interface` resource ([#13032](https://github.com/hashicorp/terraform-provider-google/pull/13032))
* compute: added `default_route_action.request_mirror_policy` field to `google_compute_region_url_map` resource ([#13030](https://github.com/hashicorp/terraform-provider-google/pull/13030))
* compute: added `default_route_action.retry_policy` field to `google_compute_region_url_map` resource ([#13030](https://github.com/hashicorp/terraform-provider-google/pull/13030))
* compute: added `default_route_action.weighted_backend_services` field to `google_compute_region_url_map` resource ([#13030](https://github.com/hashicorp/terraform-provider-google/pull/13030))
* compute: modified machine_type field in compute instance resource to accept short name. ([#12965](https://github.com/hashicorp/terraform-provider-google/pull/12965))
* compute: added `visible_core_count` field to `google_compute_instance` ([#13043](https://github.com/hashicorp/terraform-provider-google/pull/13043))
* container: added `enable_l4_ilb_subsetting` to GA `google_container_cluster` ([#12988](https://github.com/hashicorp/terraform-provider-google/pull/12988))
* container: added `node_config.logging_variant` to `google_container_node_pool`. ([#13049](https://github.com/hashicorp/terraform-provider-google/pull/13049))
* container: added `node_pool_defaults.node_config_defaults.logging_variant`, `node_pool.node_config.logging_variant`, and `node_config.logging_variant` to `google_container_cluster`. ([#13049](https://github.com/hashicorp/terraform-provider-google/pull/13049))
* container: added support for Shielded Instance configuration for node auto-provisioning to `google_container_cluster` ([#12930](https://github.com/hashicorp/terraform-provider-google/pull/12930))
* container: added management attribute to the `google_container_cluster` resource ([#12987](https://github.com/hashicorp/terraform-provider-google/pull/12987))
* container: added field `blue_green_settings` to `google_container_node_pool` ([#12984](https://github.com/hashicorp/terraform-provider-google/pull/12984))
* container: added field `strategy` to `google_container_node_pool` ([#12984](https://github.com/hashicorp/terraform-provider-google/pull/12984))
* container: added support for additional values `APISERVER`, `CONTROLLER_MANAGER`, and `SCHEDULER` in `google_container_cluster.monitoring_config` ([#12978](https://github.com/hashicorp/terraform-provider-google/pull/12978))
* datafusion: added `enable_rbac` field to `google_data_fusion_instance` resource ([#12992](https://github.com/hashicorp/terraform-provider-google/pull/12992))
* dlp: added fields `rows_limit`, `rows_limit_percent`, and `sample_method` to `big_query_options` in `google_data_loss_prevention_job_trigger` ([#12980](https://github.com/hashicorp/terraform-provider-google/pull/12980))
* dlp: added pubsub action to `google_data_loss_prevention_job_trigger` ([#12929](https://github.com/hashicorp/terraform-provider-google/pull/12929))
* dns: added `gke_clusters` field to `google_dns_managed_zone` resource ([#13048](https://github.com/hashicorp/terraform-provider-google/pull/13048))
* dns: added `gke_clusters` field to `google_dns_response_policy` resource ([#13048](https://github.com/hashicorp/terraform-provider-google/pull/13048))
* eventarc: added field `channel` to `google_eventarc_trigger` ([#13021](https://github.com/hashicorp/terraform-provider-google/pull/13021))
* gkehub: added `mesh` field and `management` subfield to resource `feature_membership` ([#13012](https://github.com/hashicorp/terraform-provider-google/pull/13012))
* networkservices: added `aws_v4_authentication ` field to `google_network_services_edge_cache_origin ` to support S3-compatible Origins ([#13020](https://github.com/hashicorp/terraform-provider-google/pull/13020))
* networkservices: added `signed_token_options` and `add_signatures` field to `google_network_services_edge_cache_service` and `validation_shared_keys` to `google_network_services_edge_cache_keyset` to support dual-token authentication ([#13041](https://github.com/hashicorp/terraform-provider-google/pull/13041))
* sql: added `query_plan_per_minute` field to `insights_config` in `google_sql_database_instance` resource ([#12951](https://github.com/hashicorp/terraform-provider-google/pull/12951))
* vertexai: added fields to `vertex_ai_featurestore_entitytype` to support feature value monitoring ([#12983](https://github.com/hashicorp/terraform-provider-google/pull/12983))

BUG FIXES:
* apigee: fixed permadiff on `consumer_accept_list` for `google_apigee_instance` ([#13037](https://github.com/hashicorp/terraform-provider-google/pull/13037))
* appengine: fixed permadiff on `serviceaccount` for 'google_app_engine_flexible_app_version' ([#12982](https://github.com/hashicorp/terraform-provider-google/pull/12982))
* bigtable: updated ForceNew logic for `kms_key_name` ([#13018](https://github.com/hashicorp/terraform-provider-google/pull/13018))
* bigtable: updated the error handling logic to remove the resource on resource not found error only ([#12953](https://github.com/hashicorp/terraform-provider-google/pull/12953))
* billingbudget: fixed a bug where `budget_filter.credit_types_treatment` in `google_billing_budget` resource was not updating. ([#12947](https://github.com/hashicorp/terraform-provider-google/pull/12947))
* cloudbuild: fixed a failure when BITBUCKET is provided for `repo_type` on `google_cloudbuild_trigger` ([#13027](https://github.com/hashicorp/terraform-provider-google/pull/13027))
* cloudids: fixed `endpoint_forwarding_rule` and `endpoint_ip` attributes for `google_cloud_ids_endpoint` ([#12957](https://github.com/hashicorp/terraform-provider-google/pull/12957))
* compute: fixed perma-diff on `google_compute_disk` for new amd64 images ([#12961](https://github.com/hashicorp/terraform-provider-google/pull/12961))
* compute: made `target_https_proxy` possible to set `ssl_certificates` and `certificate_map` in `google_compute_target_https_proxy` at the same time ([#12950](https://github.com/hashicorp/terraform-provider-google/pull/12950))
* container: fixed a bug where `cluster_autoscaling.auto_provisioning_defaults.service_account` can not be set when `enable_autopilot = true` for `google_container_cluster` ([#13024](https://github.com/hashicorp/terraform-provider-google/pull/13024))
* dialogflowcx: fixed a deployment issue for `google_dialogflow_cx_version` and `google_dialogflow_cx_environment` when they are deployed to a non-global location ([#13014](https://github.com/hashicorp/terraform-provider-google/pull/13014))
* dns: fixed apply failure when `description` is set to empty string on `google_dns_managed_zone` ([#12948](https://github.com/hashicorp/terraform-provider-google/pull/12948))
* provider: fixed a crash during provider authentication for certain environments ([#13056](https://github.com/hashicorp/terraform-provider-google/pull/13056))
* storage: fixed a crash when `log_bucket` is updated with empty body on `google_storage_bucket` ([#13058](https://github.com/hashicorp/terraform-provider-google/pull/13058))
* vertexai: made google_vertex_ai_featurestore_entitytype always use regional endpoint corresponding to parent's region ([#12959](https://github.com/hashicorp/terraform-provider-google/pull/12959))

## 4.43.0 (November 7, 2022)

FEATURES:
* **New Resource:** `google_kms_crypto_key_version` ([#12926](https://github.com/hashicorp/terraform-provider-google/pull/12926))

## 4.42.1 (November 2, 2022)

BUG FIXES:
* storage: fixed a crash in `google_storage_bucket` when upgrading provider to version `4.42.0` with `lifecycle_rule.condition.age` unset ([#12922](https://github.com/hashicorp/terraform-provider-google/pull/12922))

## 4.42.0 (October 31, 2022)

FEATURES:
* **New Data Source:** `google_compute_addresses` ([#12829](https://github.com/hashicorp/terraform-provider-google/pull/12829))
* **New Data Source:** `google_compute_region_network_endpoint_group` ([#12849](https://github.com/hashicorp/terraform-provider-google/pull/12849))
* **New Resource:** `google_alloydb_cluster` ([#12772](https://github.com/hashicorp/terraform-provider-google/pull/12772))
* **New Resource:** `google_bigquery_analytics_hub_data_exchange_iam` ([#12845](https://github.com/hashicorp/terraform-provider-google/pull/12845))
* **New Resource:** `google_bigquery_analytics_hub_data_exchange` ([#12845](https://github.com/hashicorp/terraform-provider-google/pull/12845))
* **New Resource:** `google_bigquery_analytics_hub_listing_iam` ([#12845](https://github.com/hashicorp/terraform-provider-google/pull/12845))
* **New Resource:** `google_bigquery_analytics_hub_listing` ([#12845](https://github.com/hashicorp/terraform-provider-google/pull/12845))
* **New Resource:** `google_iam_workforce_pool` ([#12863](https://github.com/hashicorp/terraform-provider-google/pull/12863))
* **New Resource:** `google_monitoring_generic_service` ([#12796](https://github.com/hashicorp/terraform-provider-google/pull/12796))
* **New Resource:** `google_scc_source_iam_binding` ([#12840](https://github.com/hashicorp/terraform-provider-google/pull/12840))
* **New Resource:** `google_scc_source_iam_member` ([#12840](https://github.com/hashicorp/terraform-provider-google/pull/12840))
* **New Resource:** `google_scc_source_iam_policy` ([#12840](https://github.com/hashicorp/terraform-provider-google/pull/12840))
* **New Resource:** `google_vertex_ai_endpoint` ([#12858](https://github.com/hashicorp/terraform-provider-google/pull/12858))
* **New Resource:** `google_vertex_ai_featurestore_entitytype_feature` ([#12797](https://github.com/hashicorp/terraform-provider-google/pull/12797))
* **New Resource:** `google_vertex_ai_featurestore_entitytype` ([#12797](https://github.com/hashicorp/terraform-provider-google/pull/12797))
* **New Resource:** `google_vertex_ai_featurestore` ([#12797](https://github.com/hashicorp/terraform-provider-google/pull/12797))

IMPROVEMENTS:
* appengine: added `member` field to `google_app_engine_default_service_account` datasource ([#12768](https://github.com/hashicorp/terraform-provider-google/pull/12768))
* bigquery: added `max_time_travel_hours` field in `google_bigquery_dataset` resource ([#12830](https://github.com/hashicorp/terraform-provider-google/pull/12830))
* bigquery: added `member` field to `google_bigquery_default_service_account` datasource ([#12768](https://github.com/hashicorp/terraform-provider-google/pull/12768))
* cloudbuild: added `script` field to `google_cloudbuild_trigger` resource ([#12841](https://github.com/hashicorp/terraform-provider-google/pull/12841))
* cloudplatform: validated `project_id` for `google_project` data-source ([#12846](https://github.com/hashicorp/terraform-provider-google/pull/12846))
* compute: added `source_disk` field to `google_compute_disk` and `google_compute_region_disk` resource ([#12779](https://github.com/hashicorp/terraform-provider-google/pull/12779))
* compute: added general field `rules` to `google_compute_router_nat` ([#12815](https://github.com/hashicorp/terraform-provider-google/pull/12815))
* container: added support for in-place update of `node_config.0.tags` for `google_container_node_pool` resource ([#12773](https://github.com/hashicorp/terraform-provider-google/pull/12773))
* container: added support for the Disk type and size configuration on the GKE Node Auto-provisioning ([#12786](https://github.com/hashicorp/terraform-provider-google/pull/12786))
* container: promote `enable_cost_allocation` field in `google_container_cluster` to GA ([#12866](https://github.com/hashicorp/terraform-provider-google/pull/12866))
* datastream: added `private_connectivity` field to `google_datastream_connection_profile` ([#12844](https://github.com/hashicorp/terraform-provider-google/pull/12844))
* dns: added `enable_geo_fencing` to `routing_policy` block of `google_dns_record_set` resource ([#12859](https://github.com/hashicorp/terraform-provider-google/pull/12859))
* dns: added `health_checked_targets` to `wrr` and `geo` blocks of `google_dns_record_set` resource ([#12859](https://github.com/hashicorp/terraform-provider-google/pull/12859))
* dns: added `primary_backup` to `routing_policy` block of `google_dns_record_set` resource ([#12859](https://github.com/hashicorp/terraform-provider-google/pull/12859))
* firebase: added deletion support and new field `deletion_policy` for `google_firebase_web_app` ([#12812](https://github.com/hashicorp/terraform-provider-google/pull/12812))
* privateca: added a new field `skip_grace_period` to skip the grace period when deleting a CertificateAuthority. ([#12784](https://github.com/hashicorp/terraform-provider-google/pull/12784))
* serviceaccount: added `member` field to `google_service_account` resource and datasource ([#12768](https://github.com/hashicorp/terraform-provider-google/pull/12768))
* sql: added `time_zone` field in `google_sql_database_instance` ([#12760](https://github.com/hashicorp/terraform-provider-google/pull/12760))
* storage: added `member` field to `google_storage_project_service_account` and `google_storage_transfer_project_service_account` datasource ([#12768](https://github.com/hashicorp/terraform-provider-google/pull/12768))
* storage: promoted `public_access_prevention` field on `google_storage_bucket` resource to GA ([#12766](https://github.com/hashicorp/terraform-provider-google/pull/12766))
* vpcaccess: promoted `machine_type`, `min_instances`, `max_instances`, and `subnet` in `google_vpc_access_connector` to GA ([#12838](https://github.com/hashicorp/terraform-provider-google/pull/12838))

BUG FIXES:
* compute: made `vm_count` in `google_compute_resource_policy` optional ([#12807](https://github.com/hashicorp/terraform-provider-google/pull/12807))
* container: fixed inability to update `datapath_provider` on `google_container_cluster` by making field changes trigger resource recreation ([#12887](https://github.com/hashicorp/terraform-provider-google/pull/12887))
* pubsub: ensured topics are recreated when their schemas change. ([#12806](https://github.com/hashicorp/terraform-provider-google/pull/12806))
* redis: updated `persistence_config.rdb_snapshot_period` to optional in the `google_redis_instance` resource. ([#12872](https://github.com/hashicorp/terraform-provider-google/pull/12872))

## 4.41.0 (October 17, 2022)

KNOWN ISSUES:
* container: This release introduced a new field, `node_config.0.guest_accelerator.0.gpu_sharing_config`, to an https://www.terraform.io/language/attr-as-blocks field (`node_config.0.guest_accelerator`). As detailed on the linked page, this may cause issues for modules and/or formats other than HCL.

BREAKING CHANGES:
* sql: updated `google_sql_user.sql_server_user_details` to be read only. Any configuration attempting to set this field is invalid and will cause the provider to fail during plan time. ([#12742](https://github.com/hashicorp/terraform-provider-google/pull/12742))

FEATURES:
* **New Resource:**  `google_cloud_ids_endpoint` ([#12744](https://github.com/hashicorp/terraform-provider-google/pull/12744))

IMPROVEMENTS:
* appengine: added support for `service_account` field to `google_app_engine_standard_app_version` resource ([#12732](https://github.com/hashicorp/terraform-provider-google/pull/12732))
* bigquery: added `avro_options` field to `google_bigquery_table` resource ([#12750](https://github.com/hashicorp/terraform-provider-google/pull/12750))
* compute: added `node_config.0.guest_accelerator.0.gpu_sharing_config` field to `google_container_node_pool` resource ([#12733](https://github.com/hashicorp/terraform-provider-google/pull/12733))
* datafusion: added `crypto_key_config` field to `google_data_fusion_instance` resource ([#12737](https://github.com/hashicorp/terraform-provider-google/pull/12737))
* filestore: removed constraint that forced multiple `google_filestore_instance` creations to occur serially ([#12753](https://github.com/hashicorp/terraform-provider-google/pull/12753))

BUG FIXES:
* kms: fixed apply failure when `google_kms_crypto_key` is removed after its versions were destroyed earlier ([#12752](https://github.com/hashicorp/terraform-provider-google/pull/12752))
* monitoring: fixed a bug causing a perma-diff in `google_monitoring_alert_policy` when `cross_series_reducer` was set to "REDUCE_NONE" ([#12741](https://github.com/hashicorp/terraform-provider-google/pull/12741))


## 4.40.0 (October 10, 2022)

FEATURES:
* **New Data Source:** `google_cloudfunctions2_function` ([#12673](https://github.com/hashicorp/terraform-provider-google/pull/12673))
* **New Data Source:** `google_compute_snapshot` ([#12671](https://github.com/hashicorp/terraform-provider-google/pull/12671))
* **New Resource:** `google_compute_region_target_tcp_proxy` ([#12715](https://github.com/hashicorp/terraform-provider-google/pull/12715))
* **New Resource:** `google_identity_platform_config` ([#12665](https://github.com/hashicorp/terraform-provider-google/pull/12665))
* **New Resource:** `google_bigquery_datapolicy_data_policy` ([#12725](https://github.com/hashicorp/terraform-provider-google/pull/12725))
* **New Resource:** `google_bigquery_datapolicy_data_policy_iam_binding` ([#12725](https://github.com/hashicorp/terraform-provider-google/pull/12725))
* **New Resource:** `google_bigquery_datapolicy_data_policy_iam_member` ([#12725](https://github.com/hashicorp/terraform-provider-google/pull/12725))
* **New Resource:** `google_bigquery_datapolicy_data_policy_iam_policy` ([#12725](https://github.com/hashicorp/terraform-provider-google/pull/12725))
* **New Resource:** `google_org_policy_custom_constraint` ([#12691](https://github.com/hashicorp/terraform-provider-google/pull/12691))

IMPROVEMENTS:
* bigqueryreservation: added `concurrency` and `multiRegionAuxiliary` to `google_bigquery_reservation` ([#12687](https://github.com/hashicorp/terraform-provider-google/pull/12687))
* bigtable: added additional retry GC policy operations with a longer poll interval to avoid quota issues ([#12717](https://github.com/hashicorp/terraform-provider-google/pull/12717))
* bigtable: improved error messaging ([#12707](https://github.com/hashicorp/terraform-provider-google/pull/12707))
* compute: added support for `compression_mode` field in `google_compute_backend_bucket` and `google_compute_backend_service` ([#12674](https://github.com/hashicorp/terraform-provider-google/pull/12674))
* datastream: added field `bigquery_profile` to `google_datastream_connection_profile` ([#12693](https://github.com/hashicorp/terraform-provider-google/pull/12693))
* dns: added field `cloud_logging_config` to `google_dns_managed_zone` ([#12675](https://github.com/hashicorp/terraform-provider-google/pull/12675))
* metastore: added support `BIGQUERY` as a value in `metastore_type` for `google_dataproc_metastore_service` ([#12724](https://github.com/hashicorp/terraform-provider-google/pull/12724))
* storage: added `custom_placement_config` field to `google_storage_bucket` resource to support custom dual-region GCS buckets ([#12723](https://github.com/hashicorp/terraform-provider-google/pull/12723))
* sql: added  `password_policy` field to `google_sql_user` resource ([#12668](https://github.com/hashicorp/terraform-provider-google/pull/12668))

BUG FIXES:
* storage: fixed a bug where user specified labels get overwritten by Dataplex auto generated labels ([#12694](https://github.com/hashicorp/terraform-provider-google/pull/12694))
* storagetransfer: fixed a bug in `google_storagetransfer_job` refreshes when `transfer_schedule` was empty ([#12704](https://github.com/hashicorp/terraform-provider-google/pull/12704))

## 4.39.0 (October 3, 2022)

FEATURES:
* **New Data Source:** `google_artifact_registry_repository` ([#12637](https://github.com/hashicorp/terraform-provider-google/pull/12637))
* **New Resource:** `google_identity_platform_config` ([#12665](https://github.com/hashicorp/terraform-provider-google/pull/12665))

IMPROVEMENTS:
* certificatemanager: added public/private PEM fields `pem_certificate` / `pem_private_key` and deprecated `certificate_pem` / `private_key_pem` ([#12664](https://github.com/hashicorp/terraform-provider-google/pull/12664))
* clouddeploy: added `serial_pipeline.stages.strategy` field to `google_clouddeploy_delivery_pipeline` ([#12619](https://github.com/hashicorp/terraform-provider-google/pull/12619))
* container: added `notification_config.pubsub.filter` field to `google_container_cluster` ([#12643](https://github.com/hashicorp/terraform-provider-google/pull/12643))
* eventarc: added `channels` and `conditions` fields to `google_eventarc_trigger` ([#12619](https://github.com/hashicorp/terraform-provider-google/pull/12619))
* healthcare: added `notification_configs ` field to `google_healthcare_fhir_store` resource ([#12646](https://github.com/hashicorp/terraform-provider-google/pull/12646))
* iap: added ability to import `google_iap_brand` using ID using {{project}}/{{brand_id}} format ([#12633](https://github.com/hashicorp/terraform-provider-google/pull/12633))
* secretmanager: added output field 'version' to resource 'secret_manager_secret_version' ([#12658](https://github.com/hashicorp/terraform-provider-google/pull/12658))
* sql: added `maintenance_version` and `available_maintenance_versions` fields to `google_sql_database_instance` resource ([#12659](https://github.com/hashicorp/terraform-provider-google/pull/12659))
* storagetransfer: added `notification_config` field to `google_storage_transfer_job` resource ([#12625](https://github.com/hashicorp/terraform-provider-google/pull/12625))
* tags: added `purpose` and `purpose_data` properties to `google_tags_tag_key` ([#12649](https://github.com/hashicorp/terraform-provider-google/pull/12649))

BUG FIXES:
* bigquery: fixed a bug where `allow_quoted_newlines` and `allow_jagged_rows` could not be set to false on `google_bigquery_table` ([#12627](https://github.com/hashicorp/terraform-provider-google/pull/12627))
* cloudfunction: fixed inability to update `docker_repository` and `kms_key_name` on `google_cloudfunctions_function` ([#12662](https://github.com/hashicorp/terraform-provider-google/pull/12662))
* compute: fixed inability to manage Cloud Armor `adaptive_protection_config` on `google_compute_security_policy` ([#12661](https://github.com/hashicorp/terraform-provider-google/pull/12661))
* iam: fixed diffs between `policy_data` from `google_iam_policy` data source and policy data in API responses ([#12652](https://github.com/hashicorp/terraform-provider-google/pull/12652))
* iam: fixed permadiff resulting from empty fields being sent in requests to set conditional IAM policies ([#12653](https://github.com/hashicorp/terraform-provider-google/pull/12653))
* secretmanager: fixed a bug where `google_secret_manager_secret_version` that was destroyed outside of Terraform would not be recreated on apply ([#12644](https://github.com/hashicorp/terraform-provider-google/pull/12644))
* storagetransfer: fixed a crash in `google_storagetransfer_job` when `transfer_schedule` is empty ([#12704](https://github.com/hashicorp/terraform-provider-google/pull/12704))

## 4.38.0 (September 26, 2022)

FEATURES:
* **New Data Source:** `google_vpc_access_connector` ([#12580](https://github.com/hashicorp/terraform-provider-google/pull/12580))
* **New Resource:** `google_datastream_private_connection` ([#12574](https://github.com/hashicorp/terraform-provider-google/pull/12574))

IMPROVEMENTS:
* appengine: Added `egress_setting` for field `vpc_access_connector` to `google_app_engine_standard_app_version` ([#12606](https://github.com/hashicorp/terraform-provider-google/pull/12606))
* bigquery: added `json_extension` field to the `load` block of `google_bigquery_job` resource ([#12597](https://github.com/hashicorp/terraform-provider-google/pull/12597))
* cloudfunctions: Added `build_worker_pool` to `google_cloudfunctions_function` ([#12591](https://github.com/hashicorp/terraform-provider-google/pull/12591))
* compute: added `json_custom_config` field to `google_compute_security_policy` resource ([#12611](https://github.com/hashicorp/terraform-provider-google/pull/12611))
* redis: Added support `persistence_config` field to `google_redis_instance` resource. ([#12569](https://github.com/hashicorp/terraform-provider-google/pull/12569))
* storage: added support for `overwriteWhen` field to `transfer_options` in `google_storage_transfer_job` resource ([#12573](https://github.com/hashicorp/terraform-provider-google/pull/12573))

BUG FIXES:
* bigtable: added drift detection on `gc_rules` for `google_bigtable_gc_policy` ([#12568](https://github.com/hashicorp/terraform-provider-google/pull/12568))
* compute: fixed the inability to update `most_disruptive_allowed_action` for both `google_compute_per_instance_config` and `google_compute_region_per_instance_config` ([#12566](https://github.com/hashicorp/terraform-provider-google/pull/12566))
* container: fixed allow passing empty list to `monitoring_config` and `logging_config` in `google_container_cluster` ([#12605](https://github.com/hashicorp/terraform-provider-google/pull/12605))
* sql: fixed a bug causing a perma-diff on `disk_type` due to API values being downcased ([#12567](https://github.com/hashicorp/terraform-provider-google/pull/12567))
* storage: fixed the inability to set 0 for `lifecycle_rule.condition.age` on `google_storage_bucket` ([#12593](https://github.com/hashicorp/terraform-provider-google/pull/12593))

## 4.37.0 (September 19, 2022)

FEATURES:
* **New Resource:** `google_apigee_nat_address` ([#12536](https://github.com/hashicorp/terraform-provider-google/pull/12536))
* **New Resource:** `google_dialogflow_cx_webhook` ([#12498](https://github.com/hashicorp/terraform-provider-google/pull/12498))
* **New Resource:** `google_filestore_snapshot` ([#12490](https://github.com/hashicorp/terraform-provider-google/pull/12490))

IMPROVEMENTS:
* apigee: added read-only field `connection_state` to `google_apigee_endpoint_attachment` ([#12500](https://github.com/hashicorp/terraform-provider-google/pull/12500))
* bigtable: added support for `autoscaling_config.storage_target` to `google_bigtable_instance` ([#12510](https://github.com/hashicorp/terraform-provider-google/pull/12510))
* cloudbuild: added support for `BITBUCKET` option to `git_source.repo_type` in `google_cloudbuild_trigger` ([#12542](https://github.com/hashicorp/terraform-provider-google/pull/12542))
* dns: added in validation for trailing dot at end of DNS record name ([#12521](https://github.com/hashicorp/terraform-provider-google/pull/12521))
* project: added validation for field `project_id` in `google_project` datasource. ([#12553](https://github.com/hashicorp/terraform-provider-google/pull/12553))
* serviceaccount: added `expires_in` attribute for generating `exp` claim  to `google_service_account_jwt` datasource ([#12539](https://github.com/hashicorp/terraform-provider-google/pull/12539))

BUG FIXES:
* notebooks: fixed perma-diff in `google_notebooks_instance` ([#12493](https://github.com/hashicorp/terraform-provider-google/pull/12493))
* privateca: fixed an issue that blocked subordinate CA data sources when `state` was not `AWAITING_USER_ACTIVATION` ([#12511](https://github.com/hashicorp/terraform-provider-google/pull/12511))
* storage: fixed permdiff on the field `versioning` of `google_storage_bucket` ([#12495](https://github.com/hashicorp/terraform-provider-google/pull/12495))

## 4.36.0 (September 12, 2022)

FEATURES:
* **New Resource:** `google_datastream_connection_profile` ([#12475](https://github.com/hashicorp/terraform-provider-google/pull/12475))

IMPROVEMENTS:
* appengine: added field `service_account` to `google_app_engine_flexible_app_version` ([#12463](https://github.com/hashicorp/terraform-provider-google/pull/12463))
* bigtable: increased timeout in `google_bigtable_table` creation. ([#12468](https://github.com/hashicorp/terraform-provider-google/pull/12468))
* cloudbuild: added `location` field to `google_cloudbuild_trigger` resource ([#12450](https://github.com/hashicorp/terraform-provider-google/pull/12450))
* compute: added `certificate_map` to `compute_target_ssl_proxy` resource ([#12467](https://github.com/hashicorp/terraform-provider-google/pull/12467))
* compute: added field `chain_name` to `google_compute_resource_policy.snapshot_properties` ([#12481](https://github.com/hashicorp/terraform-provider-google/pull/12481))
* compute: added field `chain_name` to resource `google_compute_snapshot` ([#12481](https://github.com/hashicorp/terraform-provider-google/pull/12481))
* container: added `autoscaling.total_min_node_count`, `autoscaling.total_max_node_count`, and `autoscaling.location_policy` to `google_container_cluster.node_pool` ([#12453](https://github.com/hashicorp/terraform-provider-google/pull/12453))
* container: added field `node_pool_defaults` to `resource_container_cluster`. ([#12452](https://github.com/hashicorp/terraform-provider-google/pull/12452))
* dataproc: added option `shielded_instance_config` to resource `google_dataproc_workflow_template`. ([#12451](https://github.com/hashicorp/terraform-provider-google/pull/12451))
* metastore: extended default timeouts for `google_dataproc_metastore_service` from 40m to 60m ([#12462](https://github.com/hashicorp/terraform-provider-google/pull/12462))
* pubsub: made `google_pubsub_subscription.enable_exactly_once_delivery` mutable so that it updates subscription without recreation. ([#12438](https://github.com/hashicorp/terraform-provider-google/pull/12438))

## 4.35.0 (September 6, 2022)

IMPROVEMENTS:
* apigee: added support for `nodeConfig` in `google_apigee_environment` ([#12394](https://github.com/hashicorp/terraform-provider-google/pull/12394))
* apigee: added a `properties` field to `google_apigee_organization` ([#12433](https://github.com/hashicorp/terraform-provider-google/pull/12433))
* cloudfunctions2: added `secret_environment_variables` and `secret_volumes` to `google_cloudfunctions2_function` ([#12417](https://github.com/hashicorp/terraform-provider-google/pull/12417))
* compute: added support for param `visible_core_count` in `google_compute_instance` and `google_compute_instance_template` under `advanced_machine_features` ([#12404](https://github.com/hashicorp/terraform-provider-google/pull/12404))
* compute: added support documentation links to error messages for certain Compute Operation errors. ([#12418](https://github.com/hashicorp/terraform-provider-google/pull/12418))
* container: added `service_external_ips_config` support to `cluster_container` resource. ([#12415](https://github.com/hashicorp/terraform-provider-google/pull/12415))
* container: added `enable_cost_allocation` to `google_container_cluster` ([#12416](https://github.com/hashicorp/terraform-provider-google/pull/12416))
* dns: added `behavior` field to `google_dns_response_policy_rule` resource ([#12407](https://github.com/hashicorp/terraform-provider-google/pull/12407))
* monitoring: added `force_delete` field to `google_monitoring_notification_channel` resource ([#12414](https://github.com/hashicorp/terraform-provider-google/pull/12414))

BUG FIXES:
* compute: fixed the `id` format of the data source `google_compute_instance` ([#12405](https://github.com/hashicorp/terraform-provider-google/pull/12405))

## 4.34.0 (August 29, 2022)
NOTES:
* updated Bigtable go client version from 1.13 to 1.16. ([#12349](https://github.com/hashicorp/terraform-provider-google/pull/12349))

IMPROVEMENTS:
* apigee: added support for specifying retention when deleting `google_apigee_organization` ([#12336](https://github.com/hashicorp/terraform-provider-google/pull/12336))
* appengine: added `app_engine_apis` field to `google_app_engine_standard_app_version` resource ([#12339](https://github.com/hashicorp/terraform-provider-google/pull/12339))
* cloudfunction2: promoted to `google_cloudfunctions2_function` ga ([#12322](https://github.com/hashicorp/terraform-provider-google/pull/12322))
* compute: improved error messaging for compute errors ([#12333](https://github.com/hashicorp/terraform-provider-google/pull/12333))
* container: added general field `reservation_affinity` to `google_container_node_pool` ([#12375](https://github.com/hashicorp/terraform-provider-google/pull/12375))
* container: added field `auto_provisioning_network_tags` to `google_container_cluster` (beta) ([#12347](https://github.com/hashicorp/terraform-provider-google/pull/12347))
* sql: added support for major version upgrade to `google_sql_database_instance ` resource ([#12338](https://github.com/hashicorp/terraform-provider-google/pull/12338))

BUG FIXES:
* bigtable: fixed comparing column family name when reading a GC policy. ([#12381](https://github.com/hashicorp/terraform-provider-google/pull/12381))
* bigtable: passed `isTopeLevel` in getGCPolicyFromJSON() instead of hardcoding it to true. ([#12351](https://github.com/hashicorp/terraform-provider-google/pull/12351))
* composer: corrected the description of `image_version` field. ([#12329](https://github.com/hashicorp/terraform-provider-google/pull/12329))

## 4.33.0 (August 22, 2022)

FEATURES:
* **New Resource:** `google_cloudfunctions2_function` ([#12322](https://github.com/hashicorp/terraform-provider-google/pull/12322))

IMPROVEMENTS:
* container: added update support for `authenticator_groups_config` in `google_container_cluster` ([#12310](https://github.com/hashicorp/terraform-provider-google/pull/12310))
* dataflow: added ability to import `google_dataflow_job` ([#12316](https://github.com/hashicorp/terraform-provider-google/pull/12316))
* dns: added `managed_zone_id` attribute to `google_dns_managed_zone` data source ([#12312](https://github.com/hashicorp/terraform-provider-google/pull/12312))
* monitoring: added `accepted_response_status_codes` to `monitoring_uptime_check` ([#12313](https://github.com/hashicorp/terraform-provider-google/pull/12313))
* sql: added `password_validation_policy` field to `google_cloud_sql` resource ([#12320](https://github.com/hashicorp/terraform-provider-google/pull/12320))

BUG FIXES:
* bigquery: removed force replacement for `display_name` on `google_bigquery_data_transfer_config` ([#12311](https://github.com/hashicorp/terraform-provider-google/pull/12311))
* compute: fixed permadiff for `instance_termination_action` in `google_compute_instance_template` ([#12309](https://github.com/hashicorp/terraform-provider-google/pull/12309))

## 4.32.0 (August 15, 2022)

NOTES:
* Updated to Golang 1.18 ([#12246](https://github.com/hashicorp/terraform-provider-google/pull/12246))

FEATURES:
* **New Resource:** `google_dataplex_asset` ([#12210](https://github.com/hashicorp/terraform-provider-google/pull/12210))
* **New Resource:** `google_gke_hub_membership_iam_binding` ([#12280](https://github.com/hashicorp/terraform-provider-google/pull/12280))
* **New Resource:** `google_gke_hub_membership_iam_member` ([#12280](https://github.com/hashicorp/terraform-provider-google/pull/12280))
* **New Resource:** `google_gke_hub_membership_iam_policy` ([#12280](https://github.com/hashicorp/terraform-provider-google/pull/12280))

IMPROVEMENTS:
* certificatemanager: added `state`, `authorization_attempt_info` and `provisioning_issue` output fields to `google_certificate_manager_certificate` ([#12224](https://github.com/hashicorp/terraform-provider-google/pull/12224))
* compute: added `certificate_map` to `compute_target_https_proxy` resource ([#12227](https://github.com/hashicorp/terraform-provider-google/pull/12227))
* compute: added validation for name field on `google_compute_network` ([#12271](https://github.com/hashicorp/terraform-provider-google/pull/12271))
* compute: made `port` optional in `google_compute_network_endpoint` to allow network endpoints to be associated with `GCE_VM_IP` network endpoint groups ([#12267](https://github.com/hashicorp/terraform-provider-google/pull/12267))
* container: added support for additional values `APISERVER`, `CONTROLLER_MANAGER`, and `SCHEDULER` in `google_container_cluster.monitoring_config` ([#12247](https://github.com/hashicorp/terraform-provider-google/pull/12247))
* gkehub: added `monitoring` and `mutation_enabled` fields to resource `feature_membership` ([#12265](https://github.com/hashicorp/terraform-provider-google/pull/12265))
* gkehub: added better support for import for `google_gke_hub_membership` ([#12207](https://github.com/hashicorp/terraform-provider-google/pull/12207))
* pubsub: added `bigquery_config` to `google_pubsub_subscription` ([#12216](https://github.com/hashicorp/terraform-provider-google/pull/12216))
* scheduler: added `paused` field to `google_cloud_scheduler_job` ([#12190](https://github.com/hashicorp/terraform-provider-google/pull/12190))
* scheduler: added `state` output field to `google_cloud_scheduler_job` ([#12190](https://github.com/hashicorp/terraform-provider-google/pull/12190))

BUG FIXES:
* apigee: fixed an issue where `google_apigee_instance` creation would fail due to multiple concurrent instances ([#12289](https://github.com/hashicorp/terraform-provider-google/pull/12289))
* billingbudget: fixed a bug where `google_billing_budget.budget_filter.services` was not updating. ([#12270](https://github.com/hashicorp/terraform-provider-google/pull/12270))
* compute: fixed perma-diff on `google_compute_disk` for new arm64 images ([#12184](https://github.com/hashicorp/terraform-provider-google/pull/12184))
* dataflow: fixed bug where permadiff would show on `google_dataflow_job.additional_experiments` ([#12268](https://github.com/hashicorp/terraform-provider-google/pull/12268))
* storage: fixed a bug in `google_storage_bucket` where `name` was incorrectly validated. ([#12248](https://github.com/hashicorp/terraform-provider-google/pull/12248))

## 4.31.0 (Aug 1, 2022)

FEATURES:
* **New Resource:** `google_dataplex_zone` ([#12146](https://github.com/hashicorp/terraform-provider-google/pull/12146))

IMPROVEMENTS:
* bucket: added support for `matches_prefix` and `matches_suffix` in `condition` of a `lifecycle_rule` in  `google_storage_bucket` ([#12175](https://github.com/hashicorp/terraform-provider-google/pull/12175))
* compute: added `network` and `subnetwork` fields to `google_compute_region_network_endpoint_group` for PSC. ([#12176](https://github.com/hashicorp/terraform-provider-google/pull/12176))
* container: added field `boot_disk_kms_key` to `auto_provisioning_defaults` in `google_container_cluster` ([#12173](https://github.com/hashicorp/terraform-provider-google/pull/12173))
* notebooks: added `bootDiskType` support for `PD_EXTREME` in `google_notebooks_instance` ([#12181](https://github.com/hashicorp/terraform-provider-google/pull/12181))
* notebooks: added `softwareConfig.upgradeable`, `softwareConfig.postStartupScriptBehavior`, `softwareConfig.kernels` in `google_notebooks_runtime` ([#12181](https://github.com/hashicorp/terraform-provider-google/pull/12181))
* notebooks: promoted `nicType` and `reservationAffinity` in `google_notebooks_instance` to GA ([#12181](https://github.com/hashicorp/terraform-provider-google/pull/12181))
* storage: added name validation for `google_storage_bucket` ([#12183](https://github.com/hashicorp/terraform-provider-google/pull/12183))

BUG FIXES:
* Cloud IAM: fixed incorrect basePath for `IAMBetaBasePathKey` on `google_iam_workload_identity_pool` (ga) ([#12145](https://github.com/hashicorp/terraform-provider-google/pull/12145))
* compute: fixed perma-diff on `google_compute_disk` for new arm64 images ([#12184](https://github.com/hashicorp/terraform-provider-google/pull/12184))
* dns: fixed a bug where `google_dns_record_set` would create an inconsistent plan when using interpolated values in `rrdatas` ([#12157](https://github.com/hashicorp/terraform-provider-google/pull/12157))
* kms: fixed setting of resource id post-import for `google_kms_crypto_key` ([#12164](https://github.com/hashicorp/terraform-provider-google/pull/12164))
* provider: fixed a bug where user-agent was showing "dev" rather than the provider version ([#12137](https://github.com/hashicorp/terraform-provider-google/pull/12137))

## 4.30.0 (July 25, 2022)

FEATURES:
* **New Data Source:** `google_service_account_jwt` ([#12107](https://github.com/hashicorp/terraform-provider-google/pull/12107))
* **New Resource:** `google_certificate_map_entry` ([#12127](https://github.com/hashicorp/terraform-provider-google/pull/12127))
* **New Resource:** `google_certificate_map` ([#12127](https://github.com/hashicorp/terraform-provider-google/pull/12127))

IMPROVEMENTS:
* billingbudget: made `thresholdRules` optional in `google_billing_budget` ([#12087](https://github.com/hashicorp/terraform-provider-google/pull/12087))
* compute: added `instance_termination_action` field to `google_compute_instance_template` resource to support Spot VM termination action ([#12105](https://github.com/hashicorp/terraform-provider-google/pull/12105))
* compute: added `instance_termination_action` field to `google_compute_instance` resource to support Spot VM termination action ([#12105](https://github.com/hashicorp/terraform-provider-google/pull/12105))
* compute: added `request_coalescing` and `bypass_cache_on_request_headers` fields to `compute_backend_bucket` ([#12098](https://github.com/hashicorp/terraform-provider-google/pull/12098))
* compute: added support for `esp` protocol in `google_compute_packet_mirroring.filters.ip_protocols` ([#12118](https://github.com/hashicorp/terraform-provider-google/pull/12118))
* compute: promoted `rules.rate_limit_options`,  `rules.redirect_options`,  `adaptive_protection_config` in `compute_security_policy` to GA ([#12085](https://github.com/hashicorp/terraform-provider-google/pull/12085))
* dataproc: promoted `lifecycle_config` and `endpoint_config` in `google_dataproc_cluster` to GA ([#12129](https://github.com/hashicorp/terraform-provider-google/pull/12129))
* monitoring: added `evaluation_missing_data` field to `google_monitoring_alert_policy` ([#12128](https://github.com/hashicorp/terraform-provider-google/pull/12128))
* notebooks: added `reserved_ip_range` field to `google_notebooks_runtime` ([#12113](https://github.com/hashicorp/terraform-provider-google/pull/12113))

BUG FIXES:
* bigtable: fixed an incorrect diff when adding two or more clusters ([#12109](https://github.com/hashicorp/terraform-provider-google/pull/12109))
* compute: allowed properly updating `adaptive_protection_config` in `compute_security_policy` ([#12085](https://github.com/hashicorp/terraform-provider-google/pull/12085))
* notebooks: fixed a bug where `google_notebooks_runtime` can't be updated ([#12113](https://github.com/hashicorp/terraform-provider-google/pull/12113))
* sql: fixed an issue in `google_sql_database_instance` where updates would fail because of the `collation` field ([#12131](https://github.com/hashicorp/terraform-provider-google/pull/12131))

## 4.29.0 (July 18, 2022)

FEATURES:
* **New Resource:** `google_artifact_registry_repository_iam_binding` ([#12063](https://github.com/hashicorp/terraform-provider-google/pull/12063))
* **New Resource:** `google_artifact_registry_repository_iam_member` ([#12063](https://github.com/hashicorp/terraform-provider-google/pull/12063))
* **New Resource:** `google_artifact_registry_repository_iam_policy` ([#12063](https://github.com/hashicorp/terraform-provider-google/pull/12063))
* **New Resource:** `google_artifact_registry_repository` ([#12063](https://github.com/hashicorp/terraform-provider-google/pull/12063))
* **New Resource:** `google_iam_workload_identity_pool_provider` ([#12065](https://github.com/hashicorp/terraform-provider-google/pull/12065))
* **New Resource:** `google_iam_workload_identity_pool` ([#12065](https://github.com/hashicorp/terraform-provider-google/pull/12065))
* **New Resource:** `google_cloudiot_registry_iam_binding` ([#12036](https://github.com/hashicorp/terraform-provider-google/pull/12036))
* **New Resource:** `google_cloudiot_registry_iam_member` ([#12036](https://github.com/hashicorp/terraform-provider-google/pull/12036))
* **New Resource:** `google_cloudiot_registry_iam_policy` ([#12036](https://github.com/hashicorp/terraform-provider-google/pull/12036))
* **New Resource:** `google_compute_snapshot_iam_binding` ([#12028](https://github.com/hashicorp/terraform-provider-google/pull/12028))
* **New Resource:** `google_compute_snapshot_iam_member` ([#12028](https://github.com/hashicorp/terraform-provider-google/pull/12028))
* **New Resource:** `google_compute_snapshot_iam_policy` ([#12028](https://github.com/hashicorp/terraform-provider-google/pull/12028))
* **New Resource:** `google_dataproc_metastore_service` ([#12026](https://github.com/hashicorp/terraform-provider-google/pull/12026))

IMPROVEMENTS:
* container: added `binauthz_evaluation_mode` field to `resource_container_cluster`. ([#12035](https://github.com/hashicorp/terraform-provider-google/pull/12035))
* dataproc: added Support for Dataproc on GKE in `google_dataproc_cluster` ([#12076](https://github.com/hashicorp/terraform-provider-google/pull/12076))
* dataproc: added `metastore_config` in `google_dataproc_cluster` ([#12040](https://github.com/hashicorp/terraform-provider-google/pull/12040))
* metastore: add `databaseType`, `releaseChannel`, and `hiveMetastoreConfig.endpointProtocol` arguments ([#12026](https://github.com/hashicorp/terraform-provider-google/pull/12026))
* sql: added attribute "encryption_key_name" to `google_sql_database_instance` resource. ([#12039](https://github.com/hashicorp/terraform-provider-google/pull/12039))

BUG FIXES:
* bigquery: fixed case-sensitive for `user_by_email` and `group_by_email` on `google_bigquery_dataset_access` ([#12029](https://github.com/hashicorp/terraform-provider-google/pull/12029))
* clouddeploy: fixed permadiff on `execution_configs` in `google_clouddeploy_target` resource ([#12033](https://github.com/hashicorp/terraform-provider-google/pull/12033))
* cloudscheduler: fixed a diff on the last slash of uri on `google_cloud_scheduler_job` ([#12027](https://github.com/hashicorp/terraform-provider-google/pull/12027))
* compute: fixed force recreation on `provisioned_iops` of `google_compute_disk` ([#12058](https://github.com/hashicorp/terraform-provider-google/pull/12058))
* compute: fixed missing `network_interface.0.ipv6_access_config.0.external_ipv6` output on `google_compute_instance` ([#12072](https://github.com/hashicorp/terraform-provider-google/pull/12072))
* documentai: fixed a bug where eu region could not be utilized for documentai resources ([#12074](https://github.com/hashicorp/terraform-provider-google/pull/12074))
* gkehub: fixed a bug where `issuer` can't be updated on `google_gke_hub_membership` ([#12073](https://github.com/hashicorp/terraform-provider-google/pull/12073))

## 4.28.0 (July 11, 2022)

FEATURES:
* **New Resource:** google_bigquery_connection_iam_binding ([#12004](https://github.com/hashicorp/terraform-provider-google/pull/12004))
* **New Resource:** google_bigquery_connection_iam_member ([#12004](https://github.com/hashicorp/terraform-provider-google/pull/12004))
* **New Resource:** google_bigquery_connection_iam_policy ([#12004](https://github.com/hashicorp/terraform-provider-google/pull/12004))
* **New Resource:** google_cloud_tasks_queue_iam_binding ([#11987](https://github.com/hashicorp/terraform-provider-google/pull/11987))
* **New Resource:** google_cloud_tasks_queue_iam_member ([#11987](https://github.com/hashicorp/terraform-provider-google/pull/11987))
* **New Resource:** google_cloud_tasks_queue_iam_policy ([#11987](https://github.com/hashicorp/terraform-provider-google/pull/11987))
* **New Resource:** google_dataproc_autoscaling_policy_iam_binding ([#12008](https://github.com/hashicorp/terraform-provider-google/pull/12008))
* **New Resource:** google_dataproc_autoscaling_policy_iam_member ([#12008](https://github.com/hashicorp/terraform-provider-google/pull/12008))
* **New Resource:** google_dataproc_autoscaling_policy_iam_policy ([#12008](https://github.com/hashicorp/terraform-provider-google/pull/12008))
* **New Resource:** monitoring: Promoted 'monitoredproject' to GA ([#11974](https://github.com/hashicorp/terraform-provider-google/pull/11974))

IMPROVEMENTS:
* bigquery: fixed a permadiff in `google_bigquery_job.query. destination_table` ([#11936](https://github.com/hashicorp/terraform-provider-google/pull/11936))
* billing: added `calendar_period` and `custom_period` fields to `google_billing_budget` ([#11993](https://github.com/hashicorp/terraform-provider-google/pull/11993))
* cloudsql: added attribute `project` to data source `google_sql_backup_run` ([#11938](https://github.com/hashicorp/terraform-provider-google/pull/11938))
* composer: added CMEK, PUPI and IP_masq_agent support for Composer 2 in `google_composer_environment` resource ([#11994](https://github.com/hashicorp/terraform-provider-google/pull/11994))
* compute: added `max_ports_per_vm` field to `google_compute_router_nat` resource ([#11933](https://github.com/hashicorp/terraform-provider-google/pull/11933))
* compute: added `GCE_VM_IP` support to `google_compute_network_endpoint_group` resource. ([#11997](https://github.com/hashicorp/terraform-provider-google/pull/11997))
* compute: promoted `disk_encryption_key.kms_key_name` on `google_compute_region_disk` ([#11976](https://github.com/hashicorp/terraform-provider-google/pull/11976))
* container: promoted `gce_persistent_disk_csi_driver_config` addon in `google_container_cluster` resource to GA ([#11999](https://github.com/hashicorp/terraform-provider-google/pull/11999))
* container: promoted `notification_config` and `dns_cache_config` on `google_container_cluster` ([#11944](https://github.com/hashicorp/terraform-provider-google/pull/11944))
* privateca: added support to subordinate CA activation ([#11980](https://github.com/hashicorp/terraform-provider-google/pull/11980))
* redis: added CMEK key field `customer_managed_key` in `google_redis_instance ` ([#11998](https://github.com/hashicorp/terraform-provider-google/pull/11998))
* spanner: added field `version_retention_period` to `google_spanner_database` resource ([#11982](https://github.com/hashicorp/terraform-provider-google/pull/11982))
* sql: added `settings.location_preference.secondary_zone` field in `google_sql_database_instance` ([#11996](https://github.com/hashicorp/terraform-provider-google/pull/11996))
* sql: added `sql_server_audit_config` field in `google_sql_database_instance` ([#11941](https://github.com/hashicorp/terraform-provider-google/pull/11941))

BUG FIXES:
* composer: fixed a problem with updating Cloud Composer's `scheduler_count` field (https://github.com/hashicorp/terraform-provider-google/issues/11940) ([#11951](https://github.com/hashicorp/terraform-provider-google/pull/11951))
* composer: fixed permadiff on `private_environment_config.cloud_composer_connection_subnetwork` ([#11954](https://github.com/hashicorp/terraform-provider-google/pull/11954))
* container: fixed an issue where `node_config.min_cpu_platform` could cause a perma-diff in `google_container_cluster` ([#11986](https://github.com/hashicorp/terraform-provider-google/pull/11986))
* filestore: fixed a case where `google_filestore_instance.networks.network` would incorrectly see a diff between state and config when the network `id` format was used ([#11995](https://github.com/hashicorp/terraform-provider-google/pull/11995))

## 4.27.0 (June 27, 2022)

IMPROVEMENTS:
* clouddeploy: added `suspend` field to `google_clouddeploy_delivery_pipeline` resource ([#11914](https://github.com/hashicorp/terraform-provider-google/pull/11914))
* compute: added maxPortsPerVm field to `google_compute_router_nat` resource ([#11933](https://github.com/hashicorp/terraform-provider-google/pull/11933))
* compute: added `psc_connection_id` and `psc_connection_status` output fields to `google_compute_forwarding_rule` and `google_compute_global_forwarding_rule` resources ([#11892](https://github.com/hashicorp/terraform-provider-google/pull/11892))
* containeraws: made `config.instance_type` field updatable in `google_container_aws_node_pool` ([#11892](https://github.com/hashicorp/terraform-provider-google/pull/11892))

BUG FIXES:
* compute: fixed default handling for `enable_dynamic_port_allocation ` to be managed by the api ([#11887](https://github.com/hashicorp/terraform-provider-google/pull/11887))
* vertexai: Fixed a bug where terraform crashes when `force_destroy` is set in `google_vertex_ai_featurestore` resource ([#11928](https://github.com/hashicorp/terraform-provider-google/pull/11928))

## 4.26.0 (June 21, 2022)

FEATURES:
* **New Resource:** `google_cloudfunctions2_function_iam_binding` ([#11853](https://github.com/hashicorp/terraform-provider-google/pull/11853))
* **New Resource:** `google_cloudfunctions2_function_iam_member` ([#11853](https://github.com/hashicorp/terraform-provider-google/pull/11853))
* **New Resource:** `google_cloudfunctions2_function_iam_policy` ([#11853](https://github.com/hashicorp/terraform-provider-google/pull/11853))
* **New Resource:** `google_documentai_processor` ([#11879](https://github.com/hashicorp/terraform-provider-google/pull/11879))
* **New Resource:** `google_documentai_processor_default_version` ([#11879](https://github.com/hashicorp/terraform-provider-google/pull/11879))

IMPROVEMENTS:
* accesscontextmanager: Added `external_resources` to `egress_to` in `google_access_context_manager_service_perimeter` and `google_access_context_manager_service_perimeters` resource ([#11857](https://github.com/hashicorp/terraform-provider-google/pull/11857))
* cloudbuild: Added `include_build_logs` to `google_cloudbuild_trigger` ([#11866](https://github.com/hashicorp/terraform-provider-google/pull/11866))
* composer: Promoted `config.privately_used_public_ips` and `config.ip_masq_agent` in `google_composer_environment` resource to GA. ([#11849](https://github.com/hashicorp/terraform-provider-google/pull/11849))

BUG FIXES:
* dns: fixed a bug where `google_dns_record_set` resource can not be changed from default routing to Geo routing policy. ([#11872](https://github.com/hashicorp/terraform-provider-google/pull/11872))

## 4.25.0 (June 15, 2022)

IMPROVEMENTS:
* bigquery: added `connection_id` to `external_data_configuration` for `google_bigquery_table` ([#11836](https://github.com/hashicorp/terraform-provider-google/pull/11836))
* composer: promoted `config.master_authorized_networks_config` in `google_composer_environment` resource to GA. ([#11810](https://github.com/hashicorp/terraform-provider-google/pull/11810))
* compute: added `advanced_options_config` to `google_compute_security_policy` ([#11809](https://github.com/hashicorp/terraform-provider-google/pull/11809))
* compute: added `cache_key_policy` field to `google_compute_backend_bucket` resource ([#11791](https://github.com/hashicorp/terraform-provider-google/pull/11791))
* compute: added `include_named_cookies` to `cdn_policy` on `compute_backend_service` resource ([#11818](https://github.com/hashicorp/terraform-provider-google/pull/11818))
* compute: added internal IPv6 support on `google_compute_network` and `google_compute_subnetwork` ([#11842](https://github.com/hashicorp/terraform-provider-google/pull/11842))
* container: added `spot` field to `node_config` sub-resource ([#11796](https://github.com/hashicorp/terraform-provider-google/pull/11796))
* monitoring: added support for JSONPath content matchers to `google_monitoring_uptime_check_config` resource ([#11829](https://github.com/hashicorp/terraform-provider-google/pull/11829))
* monitoring: added support for `user_labels` in `google_monitoring_slo` resource ([#11833](https://github.com/hashicorp/terraform-provider-google/pull/11833)
* sql: added `sql_server_user_details` field to `google_sql_user` resource ([#11834](https://github.com/hashicorp/terraform-provider-google/pull/11834))

BUG FIXES:
* certificatemanager: fixed bug where `DEFAULT` scope would permadiff and force replace the certificate. ([#11811](https://github.com/hashicorp/terraform-provider-google/pull/11811))
* dns: fixed perma-diff for updated labels in `google_dns_managed_zone` ([#11846](https://github.com/hashicorp/terraform-provider-google/pull/11846))
* storagetransfer: fixed perm diff on transfer_options for `google_storage_transfer_job` ([#11812](https://github.com/hashicorp/terraform-provider-google/pull/11812))

## 4.24.0 (June 6, 2022)

IMPROVEMENTS:
* compute: added `cache_key_policy` field to `google_compute_backend_bucket` resource ([#11791](https://github.com/hashicorp/terraform-provider-google/pull/11791))

## 4.23.0 (June 1, 2022)

FEATURES:
* **New Data Source:** `google_tags_tag_key` ([#11753](https://github.com/hashicorp/terraform-provider-google/pull/11753))
* **New Data Source:** `google_tags_tag_value` ([#11753](https://github.com/hashicorp/terraform-provider-google/pull/11753))
* **New Resource:** `google_dataplex_lake` ([#11769](https://github.com/hashicorp/terraform-provider-google/pull/11769))

IMPROVEMENTS:
* bigqueryconnection: updated connection types to support v1 ga ([#11728](https://github.com/hashicorp/terraform-provider-google/pull/11728))
* cloudfunctions: added docker registry support for Cloud Functions ([#11729](https://github.com/hashicorp/terraform-provider-google/pull/11729))
* memcache: added `maintenance_policy` and `maintenance_schedule` to `google_memcache_instance` ([#11759](https://github.com/hashicorp/terraform-provider-google/pull/11759))

BUG FIXES:
* binaryauthorization: fixed permadiff in `google_binary_authorization_attestor` ([#11731](https://github.com/hashicorp/terraform-provider-google/pull/11731))
* service: added re-polling for service account after creation, 404s sometimes due to [eventual consistency](https://cloud.google.com/iam/docs/overview#consistency) ([#11749](https://github.com/hashicorp/terraform-provider-google/pull/11749))

## 4.22.0 (May 24, 2022)

FEATURES:
* **New Resource:** `google_bigquery_connection` ([#11701](https://github.com/hashicorp/terraform-provider-google/pull/11701))
* **New Resource:** `google_certificate_manager_certificate` ([#11685](https://github.com/hashicorp/terraform-provider-google/pull/11685))
* **New Resource:** `google_certificate_manager_dns_authorization` ([#11685](https://github.com/hashicorp/terraform-provider-google/pull/11685))
* **New Resource:** `google_clouddeploy_delivery_pipeline` ([#11658](https://github.com/hashicorp/terraform-provider-google/pull/11658))
* **New Resource:** `google_clouddeploy_target` ([#11658](https://github.com/hashicorp/terraform-provider-google/pull/11658))

IMPROVEMENTS:
* bigquery: Added connection of type cloud_resource for `google_bigquery_connection` ([#11701](https://github.com/hashicorp/terraform-provider-google/pull/11701))
* cloudfunctions: added `https_trigger_security_level` to `google_cloudfunctions_function` ([#11672](https://github.com/hashicorp/terraform-provider-google/pull/11672))
* cloudrun: added `traffic.tag` and `traffic.url` fields to `google_cloud_run_service` ([#11641](https://github.com/hashicorp/terraform-provider-google/pull/11641))
* compute: Added `enable_dynamic_port_allocation` to `google_compute_router_nat` ([#11707](https://github.com/hashicorp/terraform-provider-google/pull/11707))
* compute: added field `update_policy.most_disruptive_allowed_action` to `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#11640](https://github.com/hashicorp/terraform-provider-google/pull/11640))
* compute: added support for NEG type `PRIVATE_SERVICE_CONNECT` in `NetworkEndpointGroup` ([#11687](https://github.com/hashicorp/terraform-provider-google/pull/11687))
* compute: added support for `domain_names` attribute in `google_compute_service_attachment` ([#11702](https://github.com/hashicorp/terraform-provider-google/pull/11702))
* compute: added value `REFRESH` to field update_policy.minimal_action` in `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#11640](https://github.com/hashicorp/terraform-provider-google/pull/11640))
* container: added field `exclusion_options` to `google_container_cluster` ([#11662](https://github.com/hashicorp/terraform-provider-google/pull/11662))
* monitoring: Added `checker_type` field to `google_monitoring_uptime_check_config` resource ([#11686](https://github.com/hashicorp/terraform-provider-google/pull/11686))
* privateca: add a new field `desired_state` to manage CertificateAuthority state. ([#11638](https://github.com/hashicorp/terraform-provider-google/pull/11638))
* sql: added `active_directory_config` field in `google_sql_database_instance` ([#11678](https://github.com/hashicorp/terraform-provider-google/pull/11678))
* sql: removed requirement that Cloud SQL Insight is only allowed for Postgres in `google_sql_database_instance` ([#11699](https://github.com/hashicorp/terraform-provider-google/pull/11699))

BUG FIXES:
* compute: fixed extra diffs generated on `google_security_policy` `rules` when modifying a rule ([#11656](https://github.com/hashicorp/terraform-provider-google/pull/11656))
* container: fixed Autopilot cluster couldn't omit master ipv4 cidr in `google_container_cluster` ([#11639](https://github.com/hashicorp/terraform-provider-google/pull/11639))
* resourcemanager: fixed a bug in wrongly writing to state when creation failed on `google_project_organization_policy` ([#11676](https://github.com/hashicorp/terraform-provider-google/pull/11676))
* storage: not specifying `content` or `source` for `google_storage_bucket_object` now fails at plan-time instead of apply-time. ([#11663](https://github.com/hashicorp/terraform-provider-google/pull/11663))



## 4.21.0 (May 16, 2022)

IMPROVEMENTS:
* cloudfunctions: added CMEK support for Cloud Functions ([#11627](https://github.com/hashicorp/terraform-provider-google/pull/11627))
* compute: added `service_directory_registrations` to `google_compute_forwarding_rule` resource ([#11635](https://github.com/hashicorp/terraform-provider-google/pull/11635))
* compute: removed validation checking against a fixed set of persistent disk types ([#11630](https://github.com/hashicorp/terraform-provider-google/pull/11630))
* container: removed validation checking against a fixed set of persistent disk types ([#11630](https://github.com/hashicorp/terraform-provider-google/pull/11630))
* containeraws: added `proxy_config` to `google_container_aws_node_pool` resource ([#11635](https://github.com/hashicorp/terraform-provider-google/pull/11635))
* containerazure: added `proxy_config` to `google_container_azure_node_pool` resource ([#11635](https://github.com/hashicorp/terraform-provider-google/pull/11635))
* dataproc: removed validation checking against a fixed set of persistent disk types ([#11630](https://github.com/hashicorp/terraform-provider-google/pull/11630))
* dns: added `routing_policy` to `google_dns_record_set` resource ([#11610](https://github.com/hashicorp/terraform-provider-google/pull/11610))

BUG FIXES:
* compute: fixed a crash in `google_compute_instance` when the instance is deleted outside of Terraform ([#11602](https://github.com/hashicorp/terraform-provider-google/pull/11602))
* provider: removed printing credentials to the console if malformed JSON is given ([#11614](https://github.com/hashicorp/terraform-provider-google/pull/11614))

## 4.20.0 (May 2, 2022)

NOTES:
* `google_privateca_certificate_authority` resources now cannot be destroyed unless `deletion_protection = false` is set in state for the resource. ([#11551](https://github.com/hashicorp/terraform-provider-google/pull/11551))

FEATURES:
* **New Data Source:** `google_compute_disk` ([#11584](https://github.com/hashicorp/terraform-provider-google/pull/11584))

IMPROVEMENTS:
* apigee: added `consumer_accept_list` and `service_attachment` to `google_apigee_instance`. ([#11595](https://github.com/hashicorp/terraform-provider-google/pull/11595))
* compute: added `provisioning_model` field to `google_compute_instance_template` and `google_compute_instance` resources to support Spot VM ([#11552](https://github.com/hashicorp/terraform-provider-google/pull/11552))
* privateca: added `deletion_protection` for `google_privateca_certificate_authority`. ([#11551](https://github.com/hashicorp/terraform-provider-google/pull/11551))
* privateca: added new output fields on `google_privateca_certificate` including `issuer_certificate_authority`, `pem_certificate_chain` and `certificate_description.x509_description` ([#11553](https://github.com/hashicorp/terraform-provider-google/pull/11553))
* redis: added multi read replica field `read_replicas_mode` and `secondary_ip_range` in `google_redis_instance` ([#11592](https://github.com/hashicorp/terraform-provider-google/pull/11592))

BUG FIXES:
* compute: fixed a crash when `compute.instance` is not found ([#11602](https://github.com/hashicorp/terraform-provider-google/pull/11602))
* provider: removed printing credentials to the console if malformed JSON is given ([#11599](https://github.com/hashicorp/terraform-provider-google/pull/11599))
* sql: fixed bug where `encryption_key_name` was not being propagated to the API. ([#11601](https://github.com/hashicorp/terraform-provider-google/pull/11601))

## 4.19.0 (April 25, 2022)

IMPROVEMENTS:
* cloudbuild: made `CLOUD_LOGGING_ONLY` available as a cloud build logging option. ([#11511](https://github.com/hashicorp/terraform-provider-google/pull/11511))
* compute: added `redirect_options` field for `google_compute_security_policy` rules ([#11492](https://github.com/hashicorp/terraform-provider-google/pull/11492))
* compute: added `FIXED_STANDARD` and `STANDARD` as valid values to the field `network_interface.0.access_configs.0.network_tier` of  `google_compute_instance_template` resource ([#11536](https://github.com/hashicorp/terraform-provider-google/pull/11536))
* compute: added `FIXED_STANDARD` and `STANDARD` as valid values to the field `network_interface.0.access_configs.0.network_tier` of  `google_compute_instance` resource ([#11536](https://github.com/hashicorp/terraform-provider-google/pull/11536))
* filestore: added `kms_key_name` field to `google_filestore_instance` resource to support CMEK ([#11493](https://github.com/hashicorp/terraform-provider-google/pull/11493))
* filestore: promoted enterprise features to GA ([#11493](https://github.com/hashicorp/terraform-provider-google/pull/11493))
* logging: made `google_logging_*_bucket_config` deletable ([#11538](https://github.com/hashicorp/terraform-provider-google/pull/11538))
* notebooks: updated `container_images` on `google_notebooks_runtime` to default to the value returned by the API if not set ([#11491](https://github.com/hashicorp/terraform-provider-google/pull/11491))
* provider: modified request retry logic to retry all per-minute quota limits returned with a 403 error code. Previously, only read requests were retried. This will generally affect Google Compute Engine resources. ([#11508](https://github.com/hashicorp/terraform-provider-google/pull/11508))

BUG FIXES:
* bigquery: fixed a bug where `encryption_configuration.kms_key_name` stored the version rather than the key name. ([#11496](https://github.com/hashicorp/terraform-provider-google/pull/11496))
* compute: fixed url_mask required mis-annotation in `google_compute_region_network_endpoint_group`, making it optional ([#11517](https://github.com/hashicorp/terraform-provider-google/pull/11517))
* spanner: fixed escaping of database names with Postgres dialect in `google_spanner_database` ([#11518](https://github.com/hashicorp/terraform-provider-google/pull/11518))

## 4.18.0 (April 18, 2022)

FEATURES:
* **New Resource:** `google_privateca_certificate_template_iam_binding` ([#11464](https://github.com/hashicorp/terraform-provider-google/pull/11464))
* **New Resource:** `google_privateca_certificate_template_iam_member` ([#11464](https://github.com/hashicorp/terraform-provider-google/pull/11464))
* **New Resource:** `google_privateca_certificate_template_iam_policy` ([#11464](https://github.com/hashicorp/terraform-provider-google/pull/11464))

IMPROVEMENTS:
* bigtable: added `gc_rules` to `google_bigtable_gc_policy` resource. ([#11481](https://github.com/hashicorp/terraform-provider-google/pull/11481))
* dialogflow: added support for location based dialogflow resources ([#11470](https://github.com/hashicorp/terraform-provider-google/pull/11470))
* metastore: added support for encryption_config during service creation. ([#11468](https://github.com/hashicorp/terraform-provider-google/pull/11468))
* privateca: added support for update on CertificateAuthority and Certificate ([#11476](https://github.com/hashicorp/terraform-provider-google/pull/11476))

BUG FIXES:
* apigee: updated mutex on google_apigee_instance_attachment to lock on org_id. ([#11467](https://github.com/hashicorp/terraform-provider-google/pull/11467))
* vpcaccess: fixed an issue where `google_vpc_access_connector` would be repeatedly recreated when `network` was not specified ([#11469](https://github.com/hashicorp/terraform-provider-google/pull/11469))

## 4.17.0 (April 11, 2022)

FEATURES:
* **New Data Source:** `google_access_approval_folder_service_account` ([#11407](https://github.com/hashicorp/terraform-provider-google/pull/11407))
* **New Data Source:** `google_access_approval_organization_service_account` ([#11407](https://github.com/hashicorp/terraform-provider-google/pull/11407))
* **New Data Source:** `google_access_approval_project_service_account` ([#11407](https://github.com/hashicorp/terraform-provider-google/pull/11407))
* **New Resource:** `google_access_context_manager_access_policy_iam_binding` ([#11409](https://github.com/hashicorp/terraform-provider-google/pull/11409))
* **New Resource:** `google_access_context_manager_access_policy_iam_member` ([#11409](https://github.com/hashicorp/terraform-provider-google/pull/11409))
* **New Resource:** `google_access_context_manager_access_policy_iam_policy` ([#11409](https://github.com/hashicorp/terraform-provider-google/pull/11409))
* **New Resource:** `google_endpoints_service_consumers_iam_binding` ([#11372](https://github.com/hashicorp/terraform-provider-google/pull/11372))
* **New Resource:** `google_endpoints_service_consumers_iam_member` ([#11372](https://github.com/hashicorp/terraform-provider-google/pull/11372))
* **New Resource:** `google_endpoints_service_consumers_iam_policy` ([#11372](https://github.com/hashicorp/terraform-provider-google/pull/11372))
* **New Resource:** `google_iam_deny_policy` ([#11446](https://github.com/hashicorp/terraform-provider-google/pull/11446))

IMPROVEMENTS:
* access approval: added `active_key_version`, `ancestor_has_active_key_version`, and `invalid_key_version` fields to `google_folder_access_approval_settings`, `google_organization_access_approval_settings`, and `google_project_access_approval_settings` resources ([#11407](https://github.com/hashicorp/terraform-provider-google/pull/11407))
* access context manager: added support for scoped policies in `google_access_context_manager_access_policy` ([#11409](https://github.com/hashicorp/terraform-provider-google/pull/11409))
* apigee: added `deployment_type` and `api_proxy_type` to `google_apigee_environment` ([#11405](https://github.com/hashicorp/terraform-provider-google/pull/11405))
* bigtable: updated the examples to show users can create all 3 different flavors of AppProfile ([#11394](https://github.com/hashicorp/terraform-provider-google/pull/11394))
* cloudbuild: added `approval_config` to `google_cloudbuild_trigger` ([#11375](https://github.com/hashicorp/terraform-provider-google/pull/11375))
* composer: added support for `airflow-1` and `airflow-2` aliases in image version argument ([#11422](https://github.com/hashicorp/terraform-provider-google/pull/11422))
* dataflow: added `skip_wait_on_job_termination` attribute to `google_dataflow_job` and `google_dataflow_flex_template_job` resources (issue #10559) ([#11452](https://github.com/hashicorp/terraform-provider-google/pull/11452))
* dataproc: added `presto_config` to `dataproc_job` ([#11393](https://github.com/hashicorp/terraform-provider-google/pull/11393))
* healthcare: added support V3 parser version for Healthcare HL7 stores. ([#11430](https://github.com/hashicorp/terraform-provider-google/pull/11430))
* healthcare: added support for `ANALYTICS_V2 `and `LOSSLESS` BigQueryDestination schema types to `google_healthcare_fhir_store` ([#11426](https://github.com/hashicorp/terraform-provider-google/pull/11426))
* os-config: added field `migInstancesAllowed` to resource `os_config_patch_deployment` ([#11447](https://github.com/hashicorp/terraform-provider-google/pull/11447))
* privateca: added support for IAM conditions to CaPool ([#11392](https://github.com/hashicorp/terraform-provider-google/pull/11392))
* pubsub: added `enable_exactly_once_delivery` to `google_pubsub_subscription` ([#11384](https://github.com/hashicorp/terraform-provider-google/pull/11384))
* spanner: added support for setting database_dialect on `google_spanner_database` ([#11363](https://github.com/hashicorp/terraform-provider-google/pull/11363))

BUG FIXES:
* redis: fixed an issue where older redis instances had a dangerous diff on the field `read_replicas_mode`, adding a default of `READ_REPLICAS_DISABLED`. Now, if the field is not set in config, the value of the field will keep the old value from state. ([#11420](https://github.com/hashicorp/terraform-provider-google/pull/11420))
* tags: fixed issue where tags could not be applied sequentially to the same parent in `google_tags_tag_binding` ([#11442](https://github.com/hashicorp/terraform-provider-google/pull/11442))

## 4.16.0 (April 4, 2022)
NOTE: We're marked a change in this release as a `BREAKING CHANGE` to indicate that the change may cause undesirable behavior for users in some circumstances. This is done to increase visibility on the change, which otherwise would have been marked under the `BUG FIXES` category, and it is not believed to be a change that breaks the backwards compatibility of the provider requiring a major version change.

BREAKING CHANGES:
* composer: made the `google_composer_environment.config.software_config.image_version` field immutable as updating this field is only available in beta. ([#11309](https://github.com/hashicorp/terraform-provider-google/pull/11309))

FEATURES:
* **New Resource:** `google_firebaserules_release` ([#11297](https://github.com/hashicorp/terraform-provider-google/pull/11297))
* **New Resource:** `google_firebaserules_ruleset` ([#11297](https://github.com/hashicorp/terraform-provider-google/pull/11297))

IMPROVEMENTS:
* apigee: added field `billing_type`([#11285](https://github.com/hashicorp/terraform-provider-google/pull/11285))
* bigtable: added support for `autoscaling_config` to `google_bigtable_instance` ([#11344](https://github.com/hashicorp/terraform-provider-google/pull/11344))
* composer: Added support for `composer-1` and `composer-2` aliases in image version argument ([#11296](https://github.com/hashicorp/terraform-provider-google/pull/11296))
* compute: added support for attaching a `edge_security_policy` to `google_compute_backend_bucket` ([#11350](https://github.com/hashicorp/terraform-provider-google/pull/11350))
* compute: added support for field `type` to `google_compute_security_policy` ([#11350](https://github.com/hashicorp/terraform-provider-google/pull/11350))
* eventarc: added gke and workflows destination for eventarc trigger resource. ([#11347](https://github.com/hashicorp/terraform-provider-google/pull/11347))
* networkservices: added `included_cookie_names` to cache key policy configuration ([#11333](https://github.com/hashicorp/terraform-provider-google/pull/11333))
* redis: added read replica field `replicaCount `, `nodes`,  `readEndpoint`, `readEndpointPort`, `readReplicasMode` in `google_redis_instance` ([#11330](https://github.com/hashicorp/terraform-provider-google/pull/11330))
* spanner: added support for setting database_dialect on `google_spanner_database` ([#11363](https://github.com/hashicorp/terraform-provider-google/pull/11363))
* storagetransfer: added `repeat_interval` field to `google_storage_transfer_job` resource ([#11328](https://github.com/hashicorp/terraform-provider-google/pull/11328))

BUG FIXES:
* apikeys: fixed a bug where `google_apikeys_key.key_string` was not being set. ([#11308](https://github.com/hashicorp/terraform-provider-google/pull/11308))
* container: fixed a bug where `google_container_cluster.authenticator_groups_config` could not be set in tandem with `enable_autopilot` ([#11310](https://github.com/hashicorp/terraform-provider-google/pull/11310))
* iam: fixed an issue where special identifiers `allAuthenticatedUsers` and `allUsers` were flattened to lower case in IAM members. ([#11359](https://github.com/hashicorp/terraform-provider-google/pull/11359))
* logging: fixed bug where `google_logging_project_bucket_config` would erroneously write to state after it errored out and wasn't actually created. ([#11314](https://github.com/hashicorp/terraform-provider-google/pull/11314))
* monitoring: fixed a permadiff when `google_monitoring_uptime_check_config.http_check.path` does not begin with "/" ([#11301](https://github.com/hashicorp/terraform-provider-google/pull/11301))
* osconfig: fixed a bug where `recurring_schedule.time_of_day` can not be set to 12am exact time in `google_os_config_patch_deployment` resource ([#11293](https://github.com/hashicorp/terraform-provider-google/pull/11293))
* storage: fixed a bug where `google_storage_bucket` data source would retry for 20 min when bucket was not found. ([#11295](https://github.com/hashicorp/terraform-provider-google/pull/11295))
* storage: fixed bug where `google_storage_transfer_job` that was deleted outside of Terraform would not be recreated on apply. ([#11307](https://github.com/hashicorp/terraform-provider-google/pull/11307))

## 4.15.0 (March 21, 2022)

FEATURES:
* **New Resource:** google_logging_log_view ([#11282](https://github.com/hashicorp/terraform-provider-google/pull/11282))

IMPROVEMENTS:
* apigee: added `billing_type` attribute to `google_apigee_organization` resource. ([#11285](https://github.com/hashicorp/terraform-provider-google/pull/11285))
* networkservices: added `disable_http2` property to `google_network_services_edge_cache_service` resource ([#11258](https://github.com/hashicorp/terraform-provider-google/pull/11258))
* networkservices: updated `google_network_services_edge_cache_origin` resource to read and write the `timeout` property, including a new `read_timeout` field. ([#11277](https://github.com/hashicorp/terraform-provider-google/pull/11277))
* networkservices: updated `google_network_services_edge_cache_origin` to retry_conditions to include `FORBIDDEN` ([#11277](https://github.com/hashicorp/terraform-provider-google/pull/11277))

BUG FIXES:
* dataproc: fixed a crash when `logging_config` only contains `nil` entry  in `google_dataproc_workflow_template` ([#11280](https://github.com/hashicorp/terraform-provider-google/pull/11280))
* sql: fixed crash when one of `settings.database_flags` is nil. ([#11279](https://github.com/hashicorp/terraform-provider-google/pull/11279))

## 4.14.0 (March 14, 2022)

FEATURES:
* **New Resource:** `google_bigqueryreservation_assignment` ([#11215](https://github.com/hashicorp/terraform-provider-google/pull/11215))
* **New Resource:** `google_apikeys_key` ([#11249](https://github.com/hashicorp/terraform-provider-google/pull/11249))

IMPROVEMENTS:
* artifactregistry: added maven config for `google_artifact_registry_repository` ([#11246](https://github.com/hashicorp/terraform-provider-google/pull/11246))
* cloudbuild: added support for manual builds, git source for webhook/pubsub triggered builds and filter field ([#11219](https://github.com/hashicorp/terraform-provider-google/pull/11219))
* composer: added support for Private Service Connect by adding `cloud_composer_connection_subnetwork` field in `google_composer_environment` ([#11223](https://github.com/hashicorp/terraform-provider-google/pull/11223))
* container: added support for gvnic to `google_container_node_pool` ([#11240](https://github.com/hashicorp/terraform-provider-google/pull/11240))
* dataproc: added `preemptibility` field to the `preemptible_worker_config` of `google_dataproc_cluster` ([#11230](https://github.com/hashicorp/terraform-provider-google/pull/11230))
* serviceusage: supported `force` behavior for deleting consumer quota override ([#11205](https://github.com/hashicorp/terraform-provider-google/pull/11205))

BUG FIXES:
* dataproc: fixed a crash when `logging_config` only contains `nil` entry  in `google_dataproc_job` ([#11232](https://github.com/hashicorp/terraform-provider-google/pull/11232))

## 4.13.0 (March 7, 2022)

FEATURES:
* **New Resource:** `google_apigee_endpoint_attachment` ([#11157](https://github.com/hashicorp/terraform-provider-google/pull/11157))
* **New Datasource:** `google_dns_record_set` ([#11180](https://github.com/hashicorp/terraform-provider-google/pull/11180))
* **New Datasource:** `google_privateca_certificate_authority` ([#11182](https://github.com/hashicorp/terraform-provider-google/pull/11182))

IMPROVEMENTS:
* composer: added support for Cloud Composer maintenance window in GA ([#11170](https://github.com/hashicorp/terraform-provider-google/pull/11170))
* compute: added support for `keepalive_interval` to `google_compute_router.bgp` ([#11188](https://github.com/hashicorp/terraform-provider-google/pull/11188))
* compute: added update support for `google_compute_reservation.share_settings` ([#11202](https://github.com/hashicorp/terraform-provider-google/pull/11202))
* storagetransfer: added attribute `subject_id` to data source `google_storage_transfer_project_service_account` ([#11156](https://github.com/hashicorp/terraform-provider-google/pull/11156))

BUG FIXES:
* composer: allow region to be undefined in configuration for `google_composer_environment` ([#11178](https://github.com/hashicorp/terraform-provider-google/pull/11178))
* container: fixed a bug where `vertical_pod_autoscaling` would cause autopilot clusters to recreate ([#11167](https://github.com/hashicorp/terraform-provider-google/pull/11167))

## 4.12.0 (February 28, 2022)

NOTE:
* updated to go 1.16.14 ([#11132](https://github.com/hashicorp/terraform-provider-google/pull/11132))

IMPROVEMENTS:
* bigquery: added support for authorized datasets to `google_bigquery_dataset.access` and `google_bigquery_dataset_access` ([#11091](https://github.com/hashicorp/terraform-provider-google/pull/11091))
* bigtable: added `multi_cluster_routing_cluster_ids` fields to `google_bigtable_app_profile` ([#11097](https://github.com/hashicorp/terraform-provider-google/pull/11097))
* compute: updated `instance` attribute for `google_compute_network_endpoint` to be optional, as Hybrid connectivity NEGs use network endpoints with just IP and Port. ([#11147](https://github.com/hashicorp/terraform-provider-google/pull/11147))
* compute: added `NON_GCP_PRIVATE_IP_PORT` value for `network_endpoint_type` in the `google_compute_network_endpoint_group` resource ([#11147](https://github.com/hashicorp/terraform-provider-google/pull/11147))
* datafusion: promoted `google_datafusion_instance` to GA ([#11087](https://github.com/hashicorp/terraform-provider-google/pull/11087))
* provider: added retries for `ReadRequest` errors incorrectly coded as `403` errors, particularly in Google Compute Engine ([#11129](https://github.com/hashicorp/terraform-provider-google/pull/11129))

BUG FIXES:
* apigee: fixed a bug where multiple `google_apigee_instance` could not be used on the same `google_apigee_organization` ([#11121](https://github.com/hashicorp/terraform-provider-google/pull/11121))
* compute: corrected an issue in `google_compute_security_policy` where only alpha values for certain enums were accepted ([#11095](https://github.com/hashicorp/terraform-provider-google/pull/11095))

## 4.11.0 (February 16, 2022)

IMPROVEMENTS:
* cloudfunctions: Added SecretManager integration support to `google_cloudfunctions_function`. ([#11062](https://github.com/hashicorp/terraform-provider-google/pull/11062))
* dataproc: increased the default timeout for `google_dataproc_cluster` from 20m to 45m ([#11026](https://github.com/hashicorp/terraform-provider-google/pull/11026))
* sql: added field `clone.allocated_ip_range` to support address range picker for clone in resource `google_sql_database_instance` ([#11058](https://github.com/hashicorp/terraform-provider-google/pull/11058))
* storagetransfer: added support for POSIX data source and data sink to `google_storage_transfer_job` via `transfer_spec.posix_data_source` and `transfer_spec.posix_data_sink` fields ([#11039](https://github.com/hashicorp/terraform-provider-google/pull/11039))

BUG FIXES:
* cloudrun: updated `containers.ports.container_port` to be optional instead of required on `google_cloud_run_service` ([#11040](https://github.com/hashicorp/terraform-provider-google/pull/11040))
* compute: marked `project` field optional in `google_compute_instance_template` data source ([#11041](https://github.com/hashicorp/terraform-provider-google/pull/11041))

## 4.10.0 (February 7, 2022)

FEATURES:
* **New Resource:** `google_backend_service_iam_*` ([#11010](https://github.com/hashicorp/terraform-provider-google/pull/11010))

IMPROVEMENTS:
* compute: added `EXTERNAL_MANAGED` as option for `load_balancing_scheme` in `google_compute_global_forwarding_rule` resource ([#10985](https://github.com/hashicorp/terraform-provider-google/pull/10985))
* compute: promoted `EXTERNAL_MANAGED` value for `load_balancing_scheme` in `google_compute_backend_service ` and `google_compute_global_forwarding_rule` to GA ([#11018](https://github.com/hashicorp/terraform-provider-google/pull/11018))
* container: added support for image type configuration on the GKE Node Auto-provisioning ([#11015](https://github.com/hashicorp/terraform-provider-google/pull/11015))
* container: added support for GCPFilestoreCSIDriver addon to `google_container_cluster` resource. ([#10998](https://github.com/hashicorp/terraform-provider-google/pull/10998))
* dataproc: increased the default timeout for `google_dataproc_cluster` from 20m to 45m ([#11026](https://github.com/hashicorp/terraform-provider-google/pull/11026))
* redis: added `maintenance_policy` and `maintenance_schedule` to `google_redis_instance` ([#10978](https://github.com/hashicorp/terraform-provider-google/pull/10978))
* vpcaccess: updated field `network` in `google_vpc_access_connector` to accept `self_link` or `name` ([#10988](https://github.com/hashicorp/terraform-provider-google/pull/10988))

BUG FIXES:
* storage: Fixed bug where the provider crashes when `Object.owner` is missing when using `google_storage_object_acl` ([#11006](https://github.com/hashicorp/terraform-provider-google/pull/11006))

## 4.9.0 (January 31, 2022)

BREAKING CHANGES:
* cloudrun: changed the `location` of `google_cloud_run_service` so that modifying the `location` field will recreate the resource rather than causing Terraform to report it would attempt an invalid update ([#10948](https://github.com/hashicorp/terraform-provider-google/pull/10948))

IMPROVEMENTS:
* provider: changed the default timeout for many resources to 20 minutes, the current Terraform default, where it was less than 20 minutes previously ([#10954](https://github.com/hashicorp/terraform-provider-google/pull/10954))
* redis: added `maintenance_policy` and `maintenance_schedule` to `google_redis_instance` ([#10978](https://github.com/hashicorp/terraform-provider-google/pull/10978))
* storage: added field `transfer_spec.aws_s3_data_source.role_arn` to `google_storage_transfer_job` ([#10950](https://github.com/hashicorp/terraform-provider-google/pull/10950))

BUG FIXES:
* cloudrun: fixed a bug where changing the non-updatable `location` of a `google_cloud_run_service` would not force resource recreation ([#10948](https://github.com/hashicorp/terraform-provider-google/pull/10948))
* compute: fixed a bug where `google_compute_firewall` would incorrectly find `source_ranges` to be empty during validation ([#10976](https://github.com/hashicorp/terraform-provider-google/pull/10976))
* notebooks: fixed permadiff in `google_notebooks_runtime.software_config` ([#10947](https://github.com/hashicorp/terraform-provider-google/pull/10947))

## 4.8.0 (January 24, 2022)

BREAKING CHANGES:
* dlp: renamed the `characters_to_ignore.character_to_skip` field to `characters_to_ignore.characters_to_skip` in `google_data_loss_prevention_deidentify_template`. Any affected configurations will have been failing with an error at apply time already. ([#10910](https://github.com/hashicorp/terraform-provider-google/pull/10910))

FEATURES:
* **New Resource:** `google_network_connectivity_spoke` ([#10921](https://github.com/hashicorp/terraform-provider-google/pull/10921))

IMPROVEMENTS:
* apigee: added `ip_range` field to `google_apigee_instance` ([#10928](https://github.com/hashicorp/terraform-provider-google/pull/10928))
* cloudrun: added support for `default_mode` and `mode` settings for created files within `secrets` in `google_cloud_run_service` ([#10911](https://github.com/hashicorp/terraform-provider-google/pull/10911))
* compute: Added `share_settings` in `google_compute_reservation` ([#10899](https://github.com/hashicorp/terraform-provider-google/pull/10899))
* container: promoted `dns_config` field of `google_container_cluster` to GA ([#10892](https://github.com/hashicorp/terraform-provider-google/pull/10892))

BUG FIXES:
* all: Fixed operation polling to support custom endpoints. ([#10913](https://github.com/hashicorp/terraform-provider-google/pull/10913))
* cloudrun: Fixed permadiff in `google_cloud_run_service`'s `template.spec.service_account_name`. ([#10940](https://github.com/hashicorp/terraform-provider-google/pull/10940))
* dlp: Fixed typo in name of `characters_to_ignore.characters_to_skip` field for `google_data_loss_prevention_deidentify_template` ([#10910](https://github.com/hashicorp/terraform-provider-google/pull/10910))
* storagetransfer: fixed bug where `schedule` was required, but really it is optional. ([#10942](https://github.com/hashicorp/terraform-provider-google/pull/10942))

## 4.7.0 (January 19, 2022)

IMPROVEMENTS:
* compute: added `EXTERNAL_MANAGED` as option for `load_balancing_scheme` in `google_compute_backend_service` resource ([#10889](https://github.com/hashicorp/terraform-provider-google/pull/10889))
* container: promoted `dns_config` field of `google_container_cluster` to GA ([#10892](https://github.com/hashicorp/terraform-provider-google/pull/10892))
* monitoring: added `conditionMatchedLog` and `alertStrategy` fields to `google_monitoring_alert_policy` resource ([#10865](https://github.com/hashicorp/terraform-provider-google/pull/10865))

## 4.6.0 (January 10, 2022)

BREAKING CHANGES:
* pubsub: changed `google_pubsub_schema` so that modifiying fields will recreate the resource rather than causing Terraform to report it would attempt an invalid update ([#10768](https://github.com/hashicorp/terraform-provider-google/pull/10768))

FEATURES:
* **New Resource:** `google_apigee_nat_address` ([#10789](https://github.com/hashicorp/terraform-provider-google/pull/10789))
* **New Resource:** `google_network_connectivity_hub` ([#10812](https://github.com/hashicorp/terraform-provider-google/pull/10812))

IMPROVEMENTS:
* bigquery: added ability to create a table with both a schema and view simultaneously to `google_bigquery_table` ([#10819](https://github.com/hashicorp/terraform-provider-google/pull/10819))
* cloud_composer: Added GA support for following fields:  `web_server_network_access_control`, `database_config`, `web_server_config`, `encryption_config`. ([#10827](https://github.com/hashicorp/terraform-provider-google/pull/10827))
* cloud_composer: Added support for Cloud Composer master authorized networks flag ([#10780](https://github.com/hashicorp/terraform-provider-google/pull/10780))
* cloud_composer: Added support for Cloud Composer v2 in GA. ([#10795](https://github.com/hashicorp/terraform-provider-google/pull/10795))
* container: promoted `node_config.0.boot_disk_kms_key` of `google_container_node_pool` to GA ([#10829](https://github.com/hashicorp/terraform-provider-google/pull/10829))
* osconfig: Added daily os config patch deployments ([#10807](https://github.com/hashicorp/terraform-provider-google/pull/10807))
* storage: added configurable read timeout to `google_storage_bucket` ([#10781](https://github.com/hashicorp/terraform-provider-google/pull/10781))

BUG FIXES:
* billingbudget: fixed a bug where `google_billing_budget.budget_filter.labels` was not updating. ([#10767](https://github.com/hashicorp/terraform-provider-google/pull/10767))
* compute: fixed scenario where `region_instance_group_manager` would not start update if `wait_for_instances` was set and initial status was not `STABLE` ([#10818](https://github.com/hashicorp/terraform-provider-google/pull/10818))
* healthcare: Added back `self_link` functionality which was accidentally removed in `4.0.0` release. ([#10808](https://github.com/hashicorp/terraform-provider-google/pull/10808))
* pubsub: fixed update failure when attempting to change non-updatable resource `google_pubsub_schema` ([#10768](https://github.com/hashicorp/terraform-provider-google/pull/10768))
* storage: fixed a bug where `google_storage_bucket.lifecycle_rule.condition.days_since_custom_time` was not updating. ([#10778](https://github.com/hashicorp/terraform-provider-google/pull/10778))
* vpcaccess: Added back `self_link` functionality which was accidentally removed in `4.0.0` release. ([#10808](https://github.com/hashicorp/terraform-provider-google/pull/10808))

## 4.5.0 (December 20, 2021)

FEATURES:
* **New Data Source:** google_container_aws_versions ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))
* **New Data Source:** google_container_azure_versions ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))
* **New Resource:** google_container_aws_cluster ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))
* **New Resource:** google_container_aws_node_pool ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))
* **New Resource:** google_container_azure_client ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))
* **New Resource:** google_container_azure_cluster ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))
* **New Resource:** google_container_azure_node_pool ([#10754](https://github.com/hashicorp/terraform-provider-google/pull/10754))

IMPROVEMENTS:
* bigquery: added the `return_table_type` field to `google_bigquery_routine` ([#10743](https://github.com/hashicorp/terraform-provider-google/pull/10743))
* cloudbuild: added support for `available_secrets` to `google_cloudbuild_trigger` ([#10714](https://github.com/hashicorp/terraform-provider-google/pull/10714))
* cloudfunctions: added support for `min_instances` to `google_cloudfunctions_function` ([#10712](https://github.com/hashicorp/terraform-provider-google/pull/10712))
* composer: added support for Private Service Connect by adding field `cloud_composer_connection_subnetwork` in `google_composer_environment` ([#10724](https://github.com/hashicorp/terraform-provider-google/pull/10724))
* compute: fixed bug where `google_compute_instance`'s `can_ip_forward` could not be updated without recreating or restarting the instance. ([#10741](https://github.com/hashicorp/terraform-provider-google/pull/10741))
* compute: added field `public_access_prevention` to resource `bucket` (beta) ([#10740](https://github.com/hashicorp/terraform-provider-google/pull/10740))
* compute: added support for regional external HTTP(S) load balancer ([#10738](https://github.com/hashicorp/terraform-provider-google/pull/10738))
* privateca: added support for setting default values for basic constraints for `google_privateca_certificate`, `google_privateca_certificate_authority`, and `google_privateca_ca_pool` via the `non_ca` and `zero_max_issuer_path_length` fields ([#10702](https://github.com/hashicorp/terraform-provider-google/pull/10702))
* provider: enabled gRPC requests and response logging ([#10721](https://github.com/hashicorp/terraform-provider-google/pull/10721))

BUG FIXES:
* assuredworkloads: fixed a bug preventing `google_assured_workloads_workload` from being created in any region other than us-central1 ([#10749](https://github.com/hashicorp/terraform-provider-google/pull/10749))

## 4.4.0 (December 13, 2021)

DEPRECATIONS:
* filestore: deprecated `zone` on `google_filestore_instance` in favor of `location` to allow for regional instances ([#10662](https://github.com/hashicorp/terraform-provider-google/pull/10662))

FEATURES:
* **New Resource:** `google_os_config_os_policy_assignment` ([#10676](https://github.com/hashicorp/terraform-provider-google/pull/10676))
* **New Resource:** `google_recaptcha_enterprise_key` ([#10672](https://github.com/hashicorp/terraform-provider-google/pull/10672))
* **New Resource:** `google_spanner_instance_iam_policy` ([#10695](https://github.com/hashicorp/terraform-provider-google/pull/10695))
* **New Resource:** `google_spanner_instance_iam_binding` ([#10695](https://github.com/hashicorp/terraform-provider-google/pull/10695))
* **New Resource:** `google_spanner_instance_iam_member` ([#10695](https://github.com/hashicorp/terraform-provider-google/pull/10695))

IMPROVEMENTS:
* filestore: added support for `ENTERPRISE` value on `google_filestore_instance` `tier` ([#10662](https://github.com/hashicorp/terraform-provider-google/pull/10662))
* privateca: added support for setting default values for basic constraints for `google_privateca_certificate`, `google_privateca_certificate_authority`, and `google_privateca_ca_pool` via the `non_ca` and `zero_max_issuer_path_length` fields ([#10702](https://github.com/hashicorp/terraform-provider-google/pull/10702))
* sql: added field `allocated_ip_range` to resource `google_sql_database_instance` ([#10687](https://github.com/hashicorp/terraform-provider-google/pull/10687))

BUG FIXES:
* compute: fixed incorrectly failing validation for `INTERNAL_MANAGED` `google_compute_region_backend_service`. ([#10664](https://github.com/hashicorp/terraform-provider-google/pull/10664))
* compute: fixed scenario where `instance_group_manager` would not start update if `wait_for_instances` was set and initial status was not `STABLE` ([#10680](https://github.com/hashicorp/terraform-provider-google/pull/10680))
* container: fixed the `ROUTES` value for the `networking_mode` field in `google_container_cluster`. A recent API change unintentionally changed the default to a `VPC_NATIVE` cluster, and removed the ability to create a `ROUTES`-based one. Provider versions prior to this one will default to `VPC_NATIVE` due to this change, and are unable to create `ROUTES` clusters. ([#10686](https://github.com/hashicorp/terraform-provider-google/pull/10686))

## 4.3.0 (December 7, 2021)

FEATURES:
* **New Data Source:** `google_compute_router_status` ([#10573](https://github.com/hashicorp/terraform-provider-google/pull/10573))
* **New Data Source:** `google_folders` ([#10658](https://github.com/hashicorp/terraform-provider-google/pull/10658))
* **New Resource:** `google_notebooks_runtime` ([#10627](https://github.com/hashicorp/terraform-provider-google/pull/10627))
* **New Resource:** `google_vertex_ai_metadata_store` ([#10657](https://github.com/hashicorp/terraform-provider-google/pull/10657))
* **New Resource:** `google_cloudbuild_worker_pool` ([#10617](https://github.com/hashicorp/terraform-provider-google/pull/10617))

IMPROVEMENTS:
* apigee: Added IAM support for `google_apigee_environment`. ([#10608](https://github.com/hashicorp/terraform-provider-google/pull/10608))
* apigee: Added supported values for 'peeringCidrRange' in `google_apigee_instance`. ([#10636](https://github.com/hashicorp/terraform-provider-google/pull/10636))
* cloudbuild: added display_name and annotations to google_cloudbuild_worker_pool for compatibility with new GA. ([#10617](https://github.com/hashicorp/terraform-provider-google/pull/10617))
* container: added `node_group` to `node_config` for container clusters and node pools to support sole tenancy ([#10646](https://github.com/hashicorp/terraform-provider-google/pull/10646))
* redis: Added Multi read replica field `replicaCount `, `nodes`,  `readEndpoint`, `readEndpointPort`, `readReplicasMode` in `google_redis_instance ` ([#10607](https://github.com/hashicorp/terraform-provider-google/pull/10607))

BUG FIXES:
* essentialcontacts: marked updating `email` in `google_essential_contacts_contact` as requiring recreation ([#10592](https://github.com/hashicorp/terraform-provider-google/pull/10592))
* privateca: fixed crlAccessUrls in `CertificateAuthority ` ([#10577](https://github.com/hashicorp/terraform-provider-google/pull/10577))

## 4.2.0 (December 2, 2021)

FEATURES:
* **New Data Source:** `google_compute_router_status` ([#10573](https://github.com/hashicorp/terraform-provider-google/pull/10573))

IMPROVEMENTS:
* compute: added support for `queue_count` to `google_compute_instance.network_interface` and `google_compute_instance_template.network_interface` ([#10571](https://github.com/hashicorp/terraform-provider-google/pull/10571))

BUG FIXES:
* all: fixed an issue where some documentation for new resources was not showing up in the GA provider if it was beta-only. ([#10545](https://github.com/hashicorp/terraform-provider-google/pull/10545))
* bigquery: fixed update failure when attempting to change non-updatable fields in `google_bigquery_routine`. ([#10546](https://github.com/hashicorp/terraform-provider-google/pull/10546))
* compute: fixed a bug when `cache_mode` is set to FORCE_CACHE_ALL on `google_compute_backend_bucket` ([#10572](https://github.com/hashicorp/terraform-provider-google/pull/10572))
* compute: fixed a perma-diff on `google_compute_region_health_check` when `log_config.enable` is set to false ([#10553](https://github.com/hashicorp/terraform-provider-google/pull/10553))
* servicedirectory: added support for vpc network configuration in `google_service_directory_endpoint`. ([#10569](https://github.com/hashicorp/terraform-provider-google/pull/10569))

## 4.1.0 (November 15, 2021)

IMPROVEMENTS:
* cloudrun: Added support for secrets to GA provider. ([#10519](https://github.com/hashicorp/terraform-provider-google/pull/10519))
* compute: Added `bfd` to `google_compute_router_peer` ([#10487](https://github.com/hashicorp/terraform-provider-google/pull/10487))
* container: added `gcfs_config` to `node_config` of `google_container_node_pool` resource ([#10499](https://github.com/hashicorp/terraform-provider-google/pull/10499))
* container: promoted `confidential_nodes` field in `google_container_cluster` to GA ([#10531](https://github.com/hashicorp/terraform-provider-google/pull/10531))
* provider: added retries for the `resourceNotReady` error returned when attempting to add resources to a recently-modified subnetwork ([#10498](https://github.com/hashicorp/terraform-provider-google/pull/10498))
* pubsub: added `message_retention_duration` field to `google_pubsub_topic` ([#10501](https://github.com/hashicorp/terraform-provider-google/pull/10501))

BUG FIXES:
* apigee: fixed a bug where multiple `google_apigee_instance_attachment` could not be used on the same `google_apigee_instance` ([#10520](https://github.com/hashicorp/terraform-provider-google/pull/10520))
* bigquery: fixed a bug following import where schema is empty on `google_bigquery_table` ([#10521](https://github.com/hashicorp/terraform-provider-google/pull/10521))
* billingbudget: fixed unable to provide `labels` on `google_billing_budget` ([#10490](https://github.com/hashicorp/terraform-provider-google/pull/10490))
* compute: allowed `source_disk` to accept full image path on `google_compute_snapshot` ([#10516](https://github.com/hashicorp/terraform-provider-google/pull/10516))
* compute: fixed a bug in `google_compute_firewall` that would cause changes in `source_ranges` to not correctly be applied ([#10515](https://github.com/hashicorp/terraform-provider-google/pull/10515))
* logging: fixed a bug with updating `description` on `google_logging_project_sink`, `google_logging_folder_sink` and `google_logging_organization_sink` ([#10493](https://github.com/hashicorp/terraform-provider-google/pull/10493))

## 4.0.0 (November 02, 2021)

NOTES:
* compute: Google Compute Engine resources will now call the endpoint appropriate to the provider version rather than the beta endpoint by default ([#10429](https://github.com/hashicorp/terraform-provider-google/pull/10429))
* container: Google Kubernetes Engine resources will now call the endpoint appropriate to the provider version rather than the beta endpoint by default ([#10430](https://github.com/hashicorp/terraform-provider-google/pull/10430))

BREAKING CHANGES:
* appengine: marked `google_app_engine_standard_app_version` `entrypoint` as required ([#10425](https://github.com/hashicorp/terraform-provider-google/pull/10425))
* compute: removed the ability to specify the `trace-append` or `trace-ro` as scopes in `google_compute_instance`, use `trace` instead ([#10377](https://github.com/hashicorp/terraform-provider-google/pull/10377))
* compute: changed `advanced_machine_features` on `google_compute_instance_template` to track changes when the block is undefined in a user's config ([#10427](https://github.com/hashicorp/terraform-provider-google/pull/10427))
* compute: changed `source_ranges` in `google_compute_firewall_rule` to track changes when it is not set in a config file ([#10439](https://github.com/hashicorp/terraform-provider-google/pull/10439))
* compute: changed the import / drift detection behaviours for `metadata_startup_script`, `metadata.startup-script` in `google_compute_instance`. Now, `metadata.startup-script` will be set by default, and `metadata_startup_script` will only be set if present. ([#10392](https://github.com/hashicorp/terraform-provider-google/pull/10392))
* compute: removed `source_disk_link` field from `google_compute_snapshot` ([#10424](https://github.com/hashicorp/terraform-provider-google/pull/10424))
* compute: removed the `enable_display` field from `google_compute_instance_template` ([#10410](https://github.com/hashicorp/terraform-provider-google/pull/10410))
* compute: removed the `update_policy.min_ready_sec` field from `google_compute_instance_group_manager`, `google_compute_region_instance_group_manager` ([#10410](https://github.com/hashicorp/terraform-provider-google/pull/10410))
* container: `instance_group_urls` has been removed in favor of `node_pool.managed_instance_group_urls` ([#10442](https://github.com/hashicorp/terraform-provider-google/pull/10442))
* container: changed default for `enable_shielded_nodes` to true for `google_container_cluster` ([#10403](https://github.com/hashicorp/terraform-provider-google/pull/10403))
* container: changed `master_auth.client_certificate_config` to required ([#10441](https://github.com/hashicorp/terraform-provider-google/pull/10441))
* container: removed `master_auth.username` and `master_auth.password` from `google_container_cluster` ([#10441](https://github.com/hashicorp/terraform-provider-google/pull/10441))
* container: removed `workload_metadata_configuration.node_metadata` in favor of `workload_metadata_configuration.mode` in `google_container_cluster` ([#10400](https://github.com/hashicorp/terraform-provider-google/pull/10400))
* container: removed the `pod_security_policy_config` field from `google_container_cluster` ([#10410](https://github.com/hashicorp/terraform-provider-google/pull/10410))
* container: removed the `workload_identity_config.0.identity_namespace` field from `google_container_cluster`, use `workload_identity_config.0.workload_pool` instead ([#10410](https://github.com/hashicorp/terraform-provider-google/pull/10410))
* project: removed ability to specify `bigquery-json.googleapis.com`, the provider will no longer convert it as the upstream API migration is finished. Use `bigquery.googleapis.com` instead. ([#10370](https://github.com/hashicorp/terraform-provider-google/pull/10370))
* provider: changed `credentials`, `access_token` precedence so that `credentials` values in configuration take precedence over `access_token` values assigned through environment variables ([#10393](https://github.com/hashicorp/terraform-provider-google/pull/10393))
* provider: removed redundant default scopes. The provider's default scopes when authenticating with credentials are now exclusively "https://www.googleapis.com/auth/cloud-platform" and "https://www.googleapis.com/auth/userinfo.email". ([#10374](https://github.com/hashicorp/terraform-provider-google/pull/10374))
* pubsub: removed `path` field from `google_pubsub_subscription` ([#10424](https://github.com/hashicorp/terraform-provider-google/pull/10424))
* resourcemanager: made `google_project` remove `org_id` and `folder_id` from state when they are removed from config ([#10373](https://github.com/hashicorp/terraform-provider-google/pull/10373))
* resourcemanager: added conflict between `org_id`, `folder_id` at plan time in `google_project` ([#10373](https://github.com/hashicorp/terraform-provider-google/pull/10373))
* resourcemanager: changed the `project` field to `Required` in all `google_project_iam_*` resources ([#10394](https://github.com/hashicorp/terraform-provider-google/pull/10394))
* runtimeconfig: removed the Runtime Configurator service from the `google` (GA) provider including `google_runtimeconfig_config`, `google_runtimeconfig_variable`, `google_runtimeconfig_config_iam_policy`, `google_runtimeconfig_config_iam_binding`, `google_runtimeconfig_config_iam_member`, `data.google_runtimeconfig_config`. They are only available in the `google-beta` provider, as the underlying service is in beta. ([#10410](https://github.com/hashicorp/terraform-provider-google/pull/10410))
* sql: added drift detection to the following `google_sql_database_instance` fields: `activation_policy` (defaults `ALWAYS`), `availability_type` (defaults `ZONAL`), `disk_type` (defaults `PD_SSD`), `encryption_key_name` ([#10412](https://github.com/hashicorp/terraform-provider-google/pull/10412))
* sql: changed the `database_version` field to `Required` in `google_sql_database_instance` resource ([#10398](https://github.com/hashicorp/terraform-provider-google/pull/10398))
* sql: removed the following `google_sql_database_instance` fields: `authorized_gae_applications`, `crash_safe_replication`, `replication_type` ([#10412](https://github.com/hashicorp/terraform-provider-google/pull/10412))
* storage: removed `bucket_policy_only` from `google_storage_bucket` ([#10397](https://github.com/hashicorp/terraform-provider-google/pull/10397))
* storage: changed the `location` field to required in `google_storage_bucket` ([#10399](https://github.com/hashicorp/terraform-provider-google/pull/10399))

VALIDATION CHANGES:
* bigquery: at least one of `statement_timeout_ms`, `statement_byte_budget`, or `key_result_statement` is required on `google_bigquery_job.query.script_options.` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* bigquery: exactly one of `query`, `load`, `copy` or `extract` is required on `google_bigquery_job` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* bigquery: exactly one of `source_table` or `source_model` is required on `google_bigquery_job.extract` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* cloudbuild: exactly one of `branch_name`, `commit_sha` or `tag_name` is required on `google_cloudbuild_trigger.build.source.repo_source` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `fixed_delay` or `percentage` is required on `google_compute_url_map.default_route_action.fault_injection_policy.delay` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `fixed` or `percent` is required on `google_compute_autoscaler.autoscaling_policy.scale_down_control.max_scaled_down_replicas` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `fixed` or `percent` is required on `google_compute_autoscaler.autoscaling_policy.scale_in_control.max_scaled_in_replicas` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `fixed` or `percent` is required on `google_compute_region_autoscaler.autoscaling_policy.scale_down_control.max_scaled_down_replicas` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `fixed` or `percent` is required on `google_compute_region_autoscaler.autoscaling_policy.scale_in_control.max_scaled_in_replicas` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `max_scaled_down_replicas` or `time_window_sec` is required on `google_compute_autoscaler.autoscaling_policy.scale_down_control` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `max_scaled_down_replicas` or `time_window_sec` is required on `google_compute_region_autoscaler.autoscaling_policy.scale_down_control` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `max_scaled_in_replicas` or `time_window_sec` is required on `google_compute_autoscaler.autoscaling_policy.scale_in_control.0.` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: at least one of `max_scaled_in_replicas` or `time_window_sec` is required on `google_compute_region_autoscaler.autoscaling_policy.scale_in_control.0.` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* compute: required one of `source_tags`, `source_ranges` or `source_service_accounts` on INGRESS `google_compute_firewall` resources ([#10369](https://github.com/hashicorp/terraform-provider-google/pull/10369))
* dlp: at least one of `start_time` or `end_time` is required on `google_data_loss_prevention_trigger.inspect_job.storage_config.timespan_config` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* dlp: exactly one of `url` or `regex_file_set` is required on `google_data_loss_prevention_trigger.inspect_job.storage_config.cloud_storage_options.file_set` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* kms: removed `self_link` field from `google_kms_crypto_key` and `google_kms_key_ring` ([#10424](https://github.com/hashicorp/terraform-provider-google/pull/10424))
* osconfig: at least one of `linux_exec_step_config` or `windows_exec_step_config` is required on `google_os_config_patch_deployment.patch_config.post_step` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* osconfig: at least one of `linux_exec_step_config` or `windows_exec_step_config` is required on `google_os_config_patch_deployment.patch_config.pre_step` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* osconfig: at least one of `reboot_config`, `apt`, `yum`, `goo` `zypper`, `windows_update`, `pre_step` or `pre_step` is required on `google_os_config_patch_deployment.patch_config` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* osconfig: at least one of `security`, `minimal`, `excludes` or `exclusive_packages` is required on `google_os_config_patch_deployment.patch_config.yum` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* osconfig: at least one of `type`, `excludes` or `exclusive_packages` is required on `google_os_config_patch_deployment.patch_config.apt` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* osconfig: at least one of `with_optional`, `with_update`, `categories`, `severities`, `excludes` or `exclusive_patches` is required on `google_os_config_patch_deployment.patch_config.zypper` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* osconfig: exactly one of `classifications`, `excludes` or `exclusive_patches` is required on `google_os_config_patch_deployment.inspect_job.patch_config.windows_update` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))
* spanner: at least one of `num_nodes` or `processing_units` is required on `google_spanner_instance` ([#10371](https://github.com/hashicorp/terraform-provider-google/pull/10371))

IMPROVEMENTS:
* compute: added `encrypted_interconnect_router` to `google_compute_router` ([#10454](https://github.com/hashicorp/terraform-provider-google/pull/10454))
* container: added `managed_instance_group_urls` to `google_container_node_pool` to replace `instance_group_urls` on `google_container_cluster` ([#10467](https://github.com/hashicorp/terraform-provider-google/pull/10467))
* kms: added support for EKM to `google_kms_crypto_key.protection_level` ([#10391](https://github.com/hashicorp/terraform-provider-google/pull/10391))
* project: added support for `billing_project` on `google_project_service` ([#10395](https://github.com/hashicorp/terraform-provider-google/pull/10395))
* spanner: increased the default timeout on `google_spanner_instance` operations from 4 minutes to 20 minutes, significantly reducing the likelihood that resources will time out ([#10437](https://github.com/hashicorp/terraform-provider-google/pull/10437))

BUG FIXES:
* bigquery: fixed a bug of cannot add required fields to an existing schema on `google_bigquery_table` ([#10421](https://github.com/hashicorp/terraform-provider-google/pull/10421))
* compute: fixed a bug in updating multiple `ttl` fields on `google_compute_backend_bucket` ([#10375](https://github.com/hashicorp/terraform-provider-google/pull/10375))
* compute: fixed a permadiff on `subnetwork` when it is optional on `google_compute_network_endpoint_group` ([#10420](https://github.com/hashicorp/terraform-provider-google/pull/10420))
* compute: fixed perma-diff bug on `log_config.enable` of both `google_compute_backend_service` and `google_compute_region_backend_service` ([#10378](https://github.com/hashicorp/terraform-provider-google/pull/10378))
* compute: fixed the `google_compute_instance_group_manager.update_policy.0.min_ready_sec` field so that updating it to `0` works ([#10457](https://github.com/hashicorp/terraform-provider-google/pull/10457))
* compute: fixed the `google_compute_region_instance_group_manager.update_policy.0.min_ready_sec` field so that updating it to `0` works ([#10457](https://github.com/hashicorp/terraform-provider-google/pull/10457))
* spanner: fixed the schema for `data.google_spanner_instance` so that non-configurable fields are considered outputs ([#10450](https://github.com/hashicorp/terraform-provider-google/pull/10450))
