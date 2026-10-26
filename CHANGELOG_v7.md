## 7.46.2 (October 23, 2026)

NOTES:
* Updated to Go 1.26.10 for security fixes to the crypto/x509 package described under https://github.com/golang/go/issues?q=milestone%3AGo1.26.10+label%3ACherryPickApproved

## 7.46.1 (September 4, 2026)

BUG FIXES:
* compute: fix permadiff regression when iap is omitted from `google_compute_backend_service` ([#29156](https://github.com/hashicorp/terraform-provider-google/pull/29156))

## 7.46.0 (August 25, 2026)

DEPRECATIONS:
* beyondcorp: deprecated `google_beyondcorp_app_connection`, `google_beyondcorp_app_connector`, and `google_beyondcorp_app_gateway` resources and data sources. Use `google_beyondcorp_security_gateway` and `google_beyondcorp_security_gateway_application` instead. ([#28976](https://github.com/hashicorp/terraform-provider-google/pull/28976))

FEATURES:
* **New Data Source:** `google_memorystore_acl_policy` ([#28984](https://github.com/hashicorp/terraform-provider-google/pull/28984))
* **New Data Source:** `google_redis_cluster_acl_policy` ([#28985](https://github.com/hashicorp/terraform-provider-google/pull/28985))
* **New List Resource:** `google_migration_center_assets_export_job` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_migration_center_discovery_client` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_migration_center_group` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_migration_center_import_job` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_migration_center_preference_set` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_migration_center_report_config` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_migration_center_source` ([#28904](https://github.com/hashicorp/terraform-provider-google/pull/28904))
* **New List Resource:** `google_network_services_authz_extension` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_consumer_association` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_domain_activation` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_domain_group` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_domain` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_group_consumer_activation` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_group_producer_activation` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_group_range_activation` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_group_range` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New List Resource:** `google_network_services_multicast_producer_association` ([#28978](https://github.com/hashicorp/terraform-provider-google/pull/28978))
* **New Resource:** `google_memorystore_acl_policy` ([#28984](https://github.com/hashicorp/terraform-provider-google/pull/28984))
* **New Resource:** `google_redis_cluster_acl_policy` ([#28985](https://github.com/hashicorp/terraform-provider-google/pull/28985))

IMPROVEMENTS:
* biglake: added `serde_info` field to `google_biglake_table` resource ([#28990](https://github.com/hashicorp/terraform-provider-google/pull/28990))
* ces: added `connector_toolset` and `timeout` fields to `google_ces_toolset` resource ([#28903](https://github.com/hashicorp/terraform-provider-google/pull/28903))
* compute: added write-only arguments for IAP `oauth2_client_id` and `oauth2_client_secret` to `google_compute_backend_service` resource ([#28809](https://github.com/hashicorp/terraform-provider-google/pull/28809))
* discoveryengine: made `google_discovery_engine_search_engine` `search_engine_config.required_subscription_tier` updatable ([#28991](https://github.com/hashicorp/terraform-provider-google/pull/28991))
* securesourcemanager: added `PULL_REQUEST_COMMENT` enum to `events` field in `google_secure_source_manager_hook` ([#28907](https://github.com/hashicorp/terraform-provider-google/pull/28907))
* sql: added `replication_lag_max_seconds` to `google_sql_database_instance` ([#28813](https://github.com/hashicorp/terraform-provider-google/pull/28813))

BUG FIXES:
* compute: fixed truncation of results at 500 images in `google_compute_images` data source ([#28909](https://github.com/hashicorp/terraform-provider-google/pull/28909))
* container: fixed a permadiff on `enable_private_endpoint` and `master_global_access_config.enabled` in `google_container_cluster` when `control_plane_endpoints_config.ip_endpoints_config.enabled` is set to `false` ([#28981](https://github.com/hashicorp/terraform-provider-google/pull/28981))
* sql: fixed `google_sql_user` returning `Missing Resource Identity After Read` when the parent Cloud SQL instance is stopped ([#28977](https://github.com/hashicorp/terraform-provider-google/pull/28977))


## 7.45.0 (August 18, 2026)

FEATURES:
* **New Data Source:** `google_iam_workload_identity_pool_openid_config` ([#28790](https://github.com/hashicorp/terraform-provider-google/pull/28790))
* **New Resource:** `google_agentic_applications_analyst_agent_persona` ([#28777](https://github.com/hashicorp/terraform-provider-google/pull/28777))
* **New Resource:** `google_firestore_change_stream` ([#28800](https://github.com/hashicorp/terraform-provider-google/pull/28800))

IMPROVEMENTS:
* accesscontextmanager: added `dry_run_access_levels` and `principal` fields to `google_access_context_manager_gcp_user_access_binding` resource ([#28767](https://github.com/hashicorp/terraform-provider-google/pull/28767))
* accesscontextmanager: updated `group_key` to be optional and conflict with `principal` on `google_access_context_manager_gcp_user_access_binding` resource ([#28767](https://github.com/hashicorp/terraform-provider-google/pull/28767))
* bigqueryreservation: added `labels` field to `google_bigquery_reservation` resource ([#28711](https://github.com/hashicorp/terraform-provider-google/pull/28711))
* certificatemanager: enabled write-only support for `pem_private_key` on `google_certificate_manager_certificate` resource ([#28794](https://github.com/hashicorp/terraform-provider-google/pull/28794))
* ces: added `snippets_config` field to `data_store_tool.modality_configs` and `service_directory_config` field to `python_function` in `google_ces_tool` ([#28759](https://github.com/hashicorp/terraform-provider-google/pull/28759))
* chronicle: added `schedule_customizations` field to `google_chronicle_rule_deployment` resource ([#28715](https://github.com/hashicorp/terraform-provider-google/pull/28715))
* cloudrunv2: added `templates.sandboxes` field to `google_cloud_run_v2_service` resource ([#28749](https://github.com/hashicorp/terraform-provider-google/pull/28749))
* colab: added `custom_environment_spec.shielded_instance_config` and `workbench_runtime.vm_image` fields to `google_colab_notebook_execution` resource ([#28774](https://github.com/hashicorp/terraform-provider-google/pull/28774))
* compute: added `custom_error_response_policy` and `default_error_response_policy` fields to `google_compute_url_map` resource ([#18511](https://github.com/hashicorp/terraform-provider-google/pull/18511))
* compute: promoted `max_run_duration` and `on_instance_stop_action` fields on `google_compute_instance`, `google_compute_instance_template`, and `google_compute_instance_from_machine_image` resources ([#18623](https://github.com/hashicorp/terraform-provider-google/pull/18623))
* dataplex: added `sql_assertion` field to `google_dataplex_datascan` resource ([#18559](https://github.com/hashicorp/terraform-provider-google/pull/18559))
* netapp: added `zone` and `replica_zone` fields to `google_netapp_storage_pool` resource ([#18609](https://github.com/hashicorp/terraform-provider-google/pull/18609))
* securityscanner: added `static_ip_scan` field to `google_security_scanner_scan_config` resource ([#28786](https://github.com/hashicorp/terraform-provider-google/pull/28786))
* vertexai: added `spec.container_spec.port` field to `google_vertex_ai_reasoning_engine` resource ([#28762](https://github.com/hashicorp/terraform-provider-google/pull/28762))
* workbench: added `compute_instance_id` field to `google_workbench_instance` resource ([#28773](https://github.com/hashicorp/terraform-provider-google/pull/28773))

BUG FIXES:
* alloydb: fixed an issue where `updateMask` URL parameter was dropped during cluster updates (e.g. `database_version` upgrade) due to variable shadowing ([#28801](https://github.com/hashicorp/terraform-provider-google/pull/28801))
* appengine: fixed permadiff in `google_app_engine_standard_app_version` ([#28745](https://github.com/hashicorp/terraform-provider-google/pull/28745))
* bigquery: fixed an issue where updating `google_bigquery_dataset` overwrote fine-grained IAM permissions ([#28775](https://github.com/hashicorp/terraform-provider-google/pull/28775))
* cloudsecuritycompliance: fixed state drift on `supported_enforcement_modes` in `google_cloud_security_compliance_framework` resource ([#28676](https://github.com/hashicorp/terraform-provider-google/pull/28676))
* colab: fixed drift detection on `direct_notebook_source.content` field in `google_colab_notebook_execution` resource ([#28774](https://github.com/hashicorp/terraform-provider-google/pull/28774))
* compute: fixed a panic in `google_compute_shared_vpc_service_project` during `terraform plan`/`refresh` when the shared VPC link had been removed outside of Terraform ([#28750](https://github.com/hashicorp/terraform-provider-google/pull/28750))
* compute: fixed permadiff on `adaptive_protection_config.layer_7_ddos_defense_config.enable` in `google_compute_security_policy` when field is not set in config ([#28751](https://github.com/hashicorp/terraform-provider-google/pull/28751))
* compute: fixed permadiffs for `google_compute_disk` on Fedora CoreOS images ([#28712](https://github.com/hashicorp/terraform-provider-google/pull/28712))
* networksecurity: fixed `google_network_security_gateway_security_policy` to force replacement when `name` or `location` is modified ([#28776](https://github.com/hashicorp/terraform-provider-google/pull/28776))

## 7.44.0 (August 11, 2026)

FEATURES:
* **New List Resource:** `google_bigquery_dataset_iam_member` ([#28578](https://github.com/hashicorp/terraform-provider-google/pull/28578))
* **New List Resource:** `google_bigquery_table` ([#28576](https://github.com/hashicorp/terraform-provider-google/pull/28576))
* **New Resource:** `google_chronicle_custom_list` ([#28618](https://github.com/hashicorp/terraform-provider-google/pull/28618))
* **New Resource:** `google_chronicle_soar_network` ([#28649](https://github.com/hashicorp/terraform-provider-google/pull/28649))
* **New Resource:** `google_dataform_repository` ([#28642](https://github.com/hashicorp/terraform-provider-google/pull/28642))
* **New Resource:** `google_dataform_repository_iam_binding` ([#28642](https://github.com/hashicorp/terraform-provider-google/pull/28642))
* **New Resource:** `google_dataform_repository_iam_member` ([#28642](https://github.com/hashicorp/terraform-provider-google/pull/28642))
* **New Resource:** `google_dataform_repository_iam_policy` ([#28642](https://github.com/hashicorp/terraform-provider-google/pull/28642))
* **New Resource:** `google_iam_folder_access_policy` ([#28664](https://github.com/hashicorp/terraform-provider-google/pull/28664))
* **New Resource:** `google_iam_organization_access_policy` ([#28664](https://github.com/hashicorp/terraform-provider-google/pull/28664))
* **New Resource:** `google_iam_project_access_policy` ([#28664](https://github.com/hashicorp/terraform-provider-google/pull/28664))
* **New Resource:** `google_vertex_ai_evaluation_metric` ([#28665](https://github.com/hashicorp/terraform-provider-google/pull/28665))

IMPROVEMENTS:
* bigquery: added resource identity support to `google_bigquery_table` ([#28574](https://github.com/hashicorp/terraform-provider-google/pull/28574))
* compute: promoted `host_error_timeout_seconds` field in `google_compute_instance`, `google_compute_instance_template` and `google_compute_region_instance_template` to GA ([#28671](https://github.com/hashicorp/terraform-provider-google/pull/28671))
* container: added `high_scale_checkpointing_config` block to `addons_config` in `google_container_cluster` ([#28669](https://github.com/hashicorp/terraform-provider-google/pull/28669))
* dataproc: added `attached_disk_config` to `disk_config` to support attached disks in the `google_dataproc_cluster` resource ([#28590](https://github.com/hashicorp/terraform-provider-google/pull/28590))
* memorystore: documented `TOKEN_AUTH` as a `google-beta`-only value for `authorization_mode` field on `google_memorystore_instance` resource ([#28584](https://github.com/hashicorp/terraform-provider-google/pull/28584))
* networkservices: added `allow_global_access` field to `google_network_services_gateway` ([#28589](https://github.com/hashicorp/terraform-provider-google/pull/28589))
* oracledatabase: added `identity_connector` to `google_oracle_database_exadb_vm_cluster` for CMEK support ([#28615](https://github.com/hashicorp/terraform-provider-google/pull/28615))
* securesourcemanager: added `service_account` and `scan_config` fields to `google_secure_source_manager_repository` ([#28662](https://github.com/hashicorp/terraform-provider-google/pull/28662))
* sql: added `password_secret_version` and `user` fields into `google_sql_provision_script` resource ([#28619](https://github.com/hashicorp/terraform-provider-google/pull/28619))
* sql: removed `ForceNew` config from `disk_type` field in `google_sql_database_instance`, allowing it to change in future without requiring the database instance to be destroyed and recreated ([#28616](https://github.com/hashicorp/terraform-provider-google/pull/28616))

BUG FIXES:
* backupdr: fixed issue where `google_backup_dr_restore_workload` dropped `resource_manager_tags` during restore requests, causing tags shown in plan to not be applied to restored resources ([#28586](https://github.com/hashicorp/terraform-provider-google/pull/28586))
* biglakeiceberg: fixed a permadiff in `google_biglake_iceberg_table` by suppressing diffs on `location` when the API-returned path contains a suffix folder ([#28577](https://github.com/hashicorp/terraform-provider-google/pull/28577))
* biglakeiceberg: fixed permadiff on `write.parquet.compression-codec` in `google_biglake_iceberg_table` ([#28650](https://github.com/hashicorp/terraform-provider-google/pull/28650))
* bigquery: fixed a provider panic when reading incomplete IAM conditions on `google_bigquery_dataset_iam_member` ([#28617](https://github.com/hashicorp/terraform-provider-google/pull/28617))
* cloudrunv2: fixed permadiff by setting `template.scaling.max_instance_count` to computed in `google_cloud_run_v2_service` ([#28644](https://github.com/hashicorp/terraform-provider-google/pull/28644))
* cloudrunv2: fixed permadiff in `template.containers.resources.limits` block in `google_cloud_run_v2_service` resource ([#28670](https://github.com/hashicorp/terraform-provider-google/pull/28670))
* cloudsecuritycompliance: fixed state drift on `supported_enforcement_modes` in `google_cloud_security_compliance_framework` resource ([#28676](https://github.com/hashicorp/terraform-provider-google/pull/28676))
* container: fixed a permadiff where `ignore_node_count_changes` was not persisted in `google_container_cluster.node_pool` ([#28573](https://github.com/hashicorp/terraform-provider-google/pull/28573))
* container: fixed an issue where `google_container_node_pool` and `google_container_cluster` failed to invalidate their Instance Group Manager cache when a resize occurred or when `ignore_node_count_changes` was active ([#28585](https://github.com/hashicorp/terraform-provider-google/pull/28585))
* networkconnectivity: made `network` field in `google_network_connectivity_transport` optional ([#28639](https://github.com/hashicorp/terraform-provider-google/pull/28639))
* recaptchaenterprise: fixed updates to `google_recaptcha_enterprise_key` so existing `challenge_settings.action_settings` entries are preserved when the action map changes ([#28673](https://github.com/hashicorp/terraform-provider-google/pull/28673))
* servicenetworking: fixed `google_service_networking_connection` ignoring the configured `delete` timeout, which used the `create` timeout instead ([#28641](https://github.com/hashicorp/terraform-provider-google/pull/28641))

## 7.43.0 (August 4, 2026)

NOTES:
* docs(workflows): added warning to `source_contents` field on `google_workflows_workflow` noting that it will become required in version 8.0.0 ([#28524](https://github.com/hashicorp/terraform-provider-google/pull/28524))
* firestore: clarified which resource to use for which kind of index in Standard and Enterprise editions ([#28529](https://github.com/hashicorp/terraform-provider-google/pull/28529))

FEATURES:
* **New List Resource:** `google_compute_instance` ([#28548](https://github.com/hashicorp/terraform-provider-google/pull/28548))
* **New List Resource:** `google_discovery_engine_assistant` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_chat_engine` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_cmek_config` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_control` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_data_store` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_license_config` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_recommendation_engine` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_schema` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_search_engine` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_serving_config` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_sitemap` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_target_site` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_discovery_engine_user_store` ([#28525](https://github.com/hashicorp/terraform-provider-google/pull/28525))
* **New List Resource:** `google_project_iam_custom_role` ([#28526](https://github.com/hashicorp/terraform-provider-google/pull/28526))
* **New List Resource:** `google_pubsub_subscription_iam_member` ([#28549](https://github.com/hashicorp/terraform-provider-google/pull/28549))
* **New List Resource:** `google_service_account_iam_member` ([#28553](https://github.com/hashicorp/terraform-provider-google/pull/28553))
* **New Resource:** `google_cloud_support_support_event_subscription` ([#28517](https://github.com/hashicorp/terraform-provider-google/pull/28517))
* **New Resource:** `google_compute_region_network_policy_traffic_classification_rule` ([#28504](https://github.com/hashicorp/terraform-provider-google/pull/28504))
* **New Resource:** `google_netapp_trial` ([#28503](https://github.com/hashicorp/terraform-provider-google/pull/28503))
* **New Resource:** `google_network_connectivity_gateway_advertised_route` GA promotion ([#28471](https://github.com/hashicorp/terraform-provider-google/pull/28471))

IMPROVEMENTS:
* apigee: added `service_account` field to `google_apigee_api_deployment` resource ([#28552](https://github.com/hashicorp/terraform-provider-google/pull/28552))
* apihub: added `source_project_id` field to `google_apihub_plugin_instance` resource ([#28530](https://github.com/hashicorp/terraform-provider-google/pull/28530))
* artifactregistry: added `no_cache` field to `google_artifact_registry_repository.remote_repository_config` ([#28506](https://github.com/hashicorp/terraform-provider-google/pull/28506))
* bigquery: added `data_governance_tags_info` field to `google_bigquery_table` resource ([#28532](https://github.com/hashicorp/terraform-provider-google/pull/28532))
* bigqueryanalyticshub: added `proposer` field to `google_bigquery_analytics_hub_query_template` ([#28518](https://github.com/hashicorp/terraform-provider-google/pull/28518))
* bigqueryanalyticshub: promoted `google_bigquery_analytics_hub_query_template` to GA ([#28518](https://github.com/hashicorp/terraform-provider-google/pull/28518))
* bigquerydatapolicyv2: added `data_governance_tag` field to `google_bigquery_datapolicyv2_data_policy` resource ([#28463](https://github.com/hashicorp/terraform-provider-google/pull/28463))
* bigqueryreservation: added `principal` field to `google_bigquery_reservation_assignment` resource ([#28508](https://github.com/hashicorp/terraform-provider-google/pull/28508))
* cloudrunv: added `sandbox_launcher` field to the containers of  `google_cloud_run_service` resource ([#28521](https://github.com/hashicorp/terraform-provider-google/pull/28521))
* cloudrunv2: added `sandbox_launcher` field to the containers of  `google_cloud_run_v2_service` resource ([#28521](https://github.com/hashicorp/terraform-provider-google/pull/28521))
* compute: promoted `ncc_gateway` field in `google_compute_router` to GA ([#28471](https://github.com/hashicorp/terraform-provider-google/pull/28471))
* discoveryengine: added `acl_enabled` field to `google_discovery_engine_data_store` resource ([#28465](https://github.com/hashicorp/terraform-provider-google/pull/28465))
* networkconnectivity: promoted `gateway` field in `google_network_connectivity_spoke` to GA ([#28471](https://github.com/hashicorp/terraform-provider-google/pull/28471))
* privateca: made `config.subject_config.subject.organization` optional in `google_privateca_certificate` ([#28464](https://github.com/hashicorp/terraform-provider-google/pull/28464))
* vertexai: added `traffic_config` field and live traffic split update support to `google_vertex_ai_reasoning_engine` ([#28473](https://github.com/hashicorp/terraform-provider-google/pull/28473))

BUG FIXES:
* apigee: fixed `continue_on_error` argument being dropped in `google_apigee_flowhook`, aligning provider behavior with the Apigee API ([#28554](https://github.com/hashicorp/terraform-provider-google/pull/28554))
* compute: fixed panic when setting scheduling attributes in `google_compute_region_instance_template` (ga) ([#28467](https://github.com/hashicorp/terraform-provider-google/pull/28467))
* gkehub2: made field `spec.workloadidentity.scopeTenancyPool` in resource `google_gke_hub_feature` not required. ([#28472](https://github.com/hashicorp/terraform-provider-google/pull/28472))

## 7.42.0 (July 28, 2026)

NOTES:
* compute: migrated `google_compute_region_instance_template` resource to use direct HTTP rather than a client library ([#28431](https://github.com/hashicorp/terraform-provider-google/pull/28431))

DEPRECATIONS:
* vertexai: deprecated `google_vertex_ai_schedule`, an accidentally-added duplicate resource; use `google_colab_schedule` instead. ([#28406](https://github.com/hashicorp/terraform-provider-google/pull/28406))

FEATURES:
* **New Data Source:** `google_cloud_quotas_quota_adjuster_settings` ([#28383](https://github.com/hashicorp/terraform-provider-google/pull/28383))
* **New List Resource:** `google_service_account_key` ([#28430](https://github.com/hashicorp/terraform-provider-google/pull/28430))
* **New Resource:** `google_agent_identity_auth_provider` ([#28447](https://github.com/hashicorp/terraform-provider-google/pull/28447))
* **New Resource:** `google_apihub_runtime_project_attachment` ([#28449](https://github.com/hashicorp/terraform-provider-google/pull/28449))
* **New Resource:** `google_chronicle_big_query_export` ([#28403](https://github.com/hashicorp/terraform-provider-google/pull/28403))
* **New Resource:** `google_compute_global_vm_extension_policy` ([#28445](https://github.com/hashicorp/terraform-provider-google/pull/28445))
* **New Resource:** `google_compute_rollout_plan` ([#28445](https://github.com/hashicorp/terraform-provider-google/pull/28445))
* **New Resource:** `google_vector_search_data_object` ([#28434](https://github.com/hashicorp/terraform-provider-google/pull/28434))
* **New Resource:** `google_vertex_ai_persistent_resource` ([#28435](https://github.com/hashicorp/terraform-provider-google/pull/28435))

IMPROVEMENTS:
* bigquery: added `table_type` field to `google_bigquery_routine` resource ([#28446](https://github.com/hashicorp/terraform-provider-google/pull/28446))
* cloudrunv2: added `start_execution_token` and `run_execution_token` fields to `google_cloud_run_v2_job`resource ([#28384](https://github.com/hashicorp/terraform-provider-google/pull/28384))
* colab: added `catch_up`, `create_pipeline_job_request`, `create_time`, `last_pause_time`, `last_resume_time`, `last_scheduled_run_response`, `max_concurrent_active_run_count`, `next_run_time`, `started_run_count`, and `update_time` fields, and sub-fields under `create_notebook_execution_job_request.notebook_execution_job` (`create_time`, `custom_environment_spec`, `encryption_spec`, `job_state`, `kernel_name`, `labels`, `name`, `schedule_resource_name`, `workbench_runtime`) and under `create_notebook_execution_job_request` (`notebook_execution_job_id`, `parent`) to `google_colab_schedule` resource ([#28406](https://github.com/hashicorp/terraform-provider-google/pull/28406))
* compute: added `effective_location` field to `google_compute_interconnect` resource ([#28416](https://github.com/hashicorp/terraform-provider-google/pull/28416))
* compute: added `request_headers` and `response_headers` fields to `log_config` on `google_compute_backend_service` and `google_compute_region_backend_service` resources ([#28421](https://github.com/hashicorp/terraform-provider-google/pull/28421))
* compute: added identity support to `google_compute_instance`, allowing resource import using an `identity` block ([#28433](https://github.com/hashicorp/terraform-provider-google/pull/28433))
* compute: changed `location` field to mutable for `google_compute_interconnect` resource ([#28416](https://github.com/hashicorp/terraform-provider-google/pull/28416))
* container: added `addons_config.node_readiness_config` field to `google_container_cluster` resource ([#28417](https://github.com/hashicorp/terraform-provider-google/pull/28417))
* container: added `rollback_safe_upgrade`, `desired_emulated_version`, and `emulated_version` fields to `google_container_cluster` resource ([#28442](https://github.com/hashicorp/terraform-provider-google/pull/28442))
* container: increased default timeout to 2 hours for `google_container_node_pool`resource ([#28382](https://github.com/hashicorp/terraform-provider-google/pull/28382))
* dataproc: added `confidential_instance_type` field to `google_dataproc_cluster` resource ([#28371](https://github.com/hashicorp/terraform-provider-google/pull/28371))
* gkehub: added `min_control_plane_version`, `min_node_version`, `target_control_plane_version`, `target_node_version`, and `operational_state` fields to `google_gke_hub_rollout_sequence` resource ([#28429](https://github.com/hashicorp/terraform-provider-google/pull/28429))
* hypercomputecluster: increased default timeouts for `google_hypercomputecluster_cluster` to 120 minutes ([#28448](https://github.com/hashicorp/terraform-provider-google/pull/28448))
* modelarmor: added field `template_metadata.filter_version_selector` to `google_model_armor_template` resource ([#28402](https://github.com/hashicorp/terraform-provider-google/pull/28402))
* sql: added identity support to `google_sql_user` for `terraform query` support ([#28428](https://github.com/hashicorp/terraform-provider-google/pull/28428))

BUG FIXES:
* bigtable: fixed an issue where `bigtable_custom_endpoint` and `universe_domain` were ignored when creating Bigtable resources ([#28404](https://github.com/hashicorp/terraform-provider-google/pull/28404))
* compute: fixed an issue where diffs in `google_compute_security_policy` were not detected ([#28420](https://github.com/hashicorp/terraform-provider-google/pull/28420))
* gkehub: fixed `rollout_creation_scope` and `upgrade_types` fields in `google_gke_hub_rollout_sequence` resource ([#28429](https://github.com/hashicorp/terraform-provider-google/pull/28429))
* osconfig: added client-side validation to ensure `resource_hierarchy_selector` and `location_selector` are not set at the same time in `google_os_config_v2_policy_orchestrator`, `google_os_config_v2_policy_orchestrator_for_folder`, and `google_os_config_v2_policy_orchestrator_for_organization` ([#28407](https://github.com/hashicorp/terraform-provider-google/pull/28407))
* secretmanager: fixed an issue where `google_secret_manager_secret_version` would fail at apply time if neither `secret_data` nor `secret_data_wo` was set ([#28419](https://github.com/hashicorp/terraform-provider-google/pull/28419))
* sql: fixed issue where updates to `settings.ip_configuration.psc_config.allowed_consumer_projects` in `google_sql_database_instance` were silently ignored on in-place updates ([#28444](https://github.com/hashicorp/terraform-provider-google/pull/28444))
* vertexai: fixed `google_vertex_ai_endpoint_with_model_garden_deployment` destroying and recreating the endpoint when `min_replica_count`, `max_replica_count`, `required_replica_count`, or `autoscaling_metric_specs` changed ([#28401](https://github.com/hashicorp/terraform-provider-google/pull/28401))

## 7.41.0 (July 17, 2026)

FEATURES:
* **New Resource:** `google_chronicle_environment_group` ([#28338](https://github.com/hashicorp/terraform-provider-google/pull/28338))
* **New Resource:** `google_compute_router_named_set` ([#28326](https://github.com/hashicorp/terraform-provider-google/pull/28326))
* **New List Resource:** `google_compute_backend_bucket_signed_url_key` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_backend_service_signed_url_key` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_network_firewall_policy` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_network_firewall_policy_association` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_network_firewall_policy_packet_mirroring_rule` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_preview_feature` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_public_advertised_prefix` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_region_backend_bucket` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))
* **New List Resource:** `google_compute_region_network_firewall_policy` ([#28357](https://github.com/hashicorp/terraform-provider-google/pull/28357))

IMPROVEMENTS:
* accesscontextmanager: added `allowed_service_patterns` and `service_patterns_enforcement_scopes` fields to `google_access_context_manager_service_perimeter` to support VPC Service Controls for non-GCP APIs. ([#28349](https://github.com/hashicorp/terraform-provider-google/pull/28349))
* accesscontextmanager: added `pscEndpoint` to `sources` in `ingress_from` and `egress_from` under resources `google_access_context_manager_service_perimeter` and variants ([#28307](https://github.com/hashicorp/terraform-provider-google/pull/28307))
* backupdr: added `backup_blocked_by_vault_access_restriction` to `data.google_backup_dr_data_source` resource ([#28361](https://github.com/hashicorp/terraform-provider-google/pull/28361))
* backupdr: added `force_update_access_restriction` to `google_backup_dr_backup_vault` resource ([#28361](https://github.com/hashicorp/terraform-provider-google/pull/28361))
* backupdr: added update support for `access_restriction` to `google_backup_dr_backup_vault` resource ([#28361](https://github.com/hashicorp/terraform-provider-google/pull/28361))
* certificatemanager: added in-place update support for the `self_managed` certificate data (`pem_certificate` / `pem_private_key`) on `google_certificate_manager_certificate`; changing the certificate data is now applied via update instead of forcing recreation ([#28337](https://github.com/hashicorp/terraform-provider-google/pull/28337))
* cloudrunv2: added `tags` field to `google_cloud_run_v2_service` and `google_cloud_run_v2_job` resources to allow setting tags for services and jobs at creation time. ([#28328](https://github.com/hashicorp/terraform-provider-google/pull/28328))
* cloudsql: added `max_custom_on_demand_retention_days` to create backup_plan example for sqladmin ([#28343](https://github.com/hashicorp/terraform-provider-google/pull/28343))
* compute: added 3500GB and 7000GB SSD partition size to `google_compute_instance_template` resource ([#28352](https://github.com/hashicorp/terraform-provider-google/pull/28352))
* compute: added `FLEX_START` and `RESERVATION_BOUND` support to `google_compute_instance`, `google_compute_instance_template`, and `google_compute_region_instance_template` resources ([#28365](https://github.com/hashicorp/terraform-provider-google/pull/28365))
* container: added the support for updating `node_image_config` and `image_type` fields at the same time ([#28283](https://github.com/hashicorp/terraform-provider-google/pull/28283))
* dataproc: added `confidential_instance_type` to `google_dataproc_cluster` resource ([#28371](https://github.com/hashicorp/terraform-provider-google/pull/28371))
* dataproc: added `instance_selection.disk_config` field to `google_dataproc_cluster` resource ([#28339](https://github.com/hashicorp/terraform-provider-google/pull/28339))
* discoveryengine: added `enable_llm_layout_parsing` and `enable_get_processed_document` fields to `google_discovery_engine_data_store` resource ([#28284](https://github.com/hashicorp/terraform-provider-google/pull/28284))
* sql: added `instance_auto_dns_status` and `write_endpoint_auto_dns_status` output fields to `psc_auto_connections` block in `google_sql_database_instance` resource ([#28331](https://github.com/hashicorp/terraform-provider-google/pull/28331))
* sql: added `include_replicas_for_major_version_upgrade` field to `google_sql_database_instance` resource ([#28345](https://github.com/hashicorp/terraform-provider-google/pull/28345))
* sql: added `switch_transaction_logs_to_cloud_storage_enabled` field to `google_sql_database_instance` resource ([#28318](https://github.com/hashicorp/terraform-provider-google/pull/28318))
* vertexai: promoted `google_vertex_ai_semantic_governance_policy_engine` resource to GA ([#28347](https://github.com/hashicorp/terraform-provider-google/pull/28347))
* workbench: added `enable_deletion_protection` field to `google_workbench_instance` resource ([#28355](https://github.com/hashicorp/terraform-provider-google/pull/28355))
* workbench: added `resource_policies` field to `google_workbench_instance` resource ([#28354](https://github.com/hashicorp/terraform-provider-google/pull/28354))
* workbench: added support for `min_cpu_platform` in `google_workbench_instance` resource ([#28369](https://github.com/hashicorp/terraform-provider-google/pull/28369))
* workstations: added `instance_metadata` field to `google_workstations_workstation_config` resource ([#28342](https://github.com/hashicorp/terraform-provider-google/pull/28342))

BUG FIXES:
* compute: fixed bug where a permadiff on `google_compute_reservation_region_commitment.existing_reservations` would persist after upgrading ([#28353](https://github.com/hashicorp/terraform-provider-google/pull/28353))
* compute: fixed permadiff on `keepalive_interval` for `google_compute_router` `bgp` block when set to default value ([#28285](https://github.com/hashicorp/terraform-provider-google/pull/28285))
* resourcemanager: fixed validation of `target_service_account` and `delegates` in the `google_service_account_access_token`, `google_service_account_id_token`, and `google_service_account_jwt` data sources, and of `name` in the `google_service_account_key` data source, to reject identifiers that contain path separators ([#28308](https://github.com/hashicorp/terraform-provider-google/pull/28308))
* securityposture: fixed a bug where the `enforce` field in `google_securityposture_posture` was always set, causing failures for list constraints. ([#28359](https://github.com/hashicorp/terraform-provider-google/pull/28359))
* vmwareengine: added correct `update_mask` value to `google_vmwareengine_private_cloud` updates ([#28360](https://github.com/hashicorp/terraform-provider-google/pull/28360))

## 7.40.0 (July 14, 2026)
DEPRECATIONS:
* storage: the `admit-on-second-miss` value for `google_storage_anywhere_cache.admission_policy` is deprecated and will be removed in a future major release. The backend will ignore this attribute and treat it as `admit-on-first-miss`. ([#28210](https://github.com/hashicorp/terraform-provider-google/pull/28210))

NOTES:
* compute: migrated `google_compute_instance` code related to advanced machine features to use direct HTTP rather than a client library ([#28160](https://github.com/hashicorp/terraform-provider-google/pull/28160))

FEATURES:
* **New Data Source:** `google_data_catalog_taxonomy` ([#28237](https://github.com/hashicorp/terraform-provider-google/pull/28237))
* **New Data Source:** `google_oracle_database_exascale_db_storage_vault` ([#28260](https://github.com/hashicorp/terraform-provider-google/pull/28260))
* **New List Resource:** 'google_project' ([#28041](https://github.com/hashicorp/terraform-provider-google/pull/28041))
* **New List Resource:** `google_compute_instant_snapshot` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_region_instant_snapshot` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_region_target_http_proxy` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_region_target_tcp_proxy` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_region_url_map` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_rollout_plan` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_target_grpc_proxy` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_target_http_proxy` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_target_ssl_proxy` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_compute_target_tcp_proxy` ([#28256](https://github.com/hashicorp/terraform-provider-google/pull/28256))
* **New List Resource:** `google_dns_managed_zone` ([#28257](https://github.com/hashicorp/terraform-provider-google/pull/28257))
* **New List Resource:** `google_oracle_database_exascale_db_storage_vaults` ([#28260](https://github.com/hashicorp/terraform-provider-google/pull/28260))
* **New Resource:** `google_chronicle_findings_refinement_deployment` ([#28240](https://github.com/hashicorp/terraform-provider-google/pull/28240))
* **New Resource:** `google_chronicle_soar_domain` ([#28214](https://github.com/hashicorp/terraform-provider-google/pull/28214))
* **New Resource:** `google_iap_agent_registry_agent_iam_binding` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_agent_iam_member` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_agent_iam_policy` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_endpoint_iam_binding` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_endpoint_iam_member` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_endpoint_iam_policy` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_mcp_server_iam_binding` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_mcp_server_iam_member` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_iap_agent_registry_mcp_server_iam_policy` ([#28231](https://github.com/hashicorp/terraform-provider-google/pull/28231))
* **New Resource:** `google_tags_tag_binding_collection` ([#28180](https://github.com/hashicorp/terraform-provider-google/pull/28180))
* **New Resource:** `google_vector_search_index` ([#28238](https://github.com/hashicorp/terraform-provider-google/pull/28238))
* **New Resource:** `google_chronicle_environment` ([#28206](https://github.com/hashicorp/terraform-provider-google/pull/28206))
* **New Resource:** `google_chronicle_data_export` ([#28239](https://github.com/hashicorp/terraform-provider-google/pull/28239))

IMPROVEMENTS:
* agentregistry: added `name` field to `google_agent_registry_binding` resource ([#28207](https://github.com/hashicorp/terraform-provider-google/pull/28207))
* agentregistry: added `name` field to `google_agent_registry_service` resource ([#28207](https://github.com/hashicorp/terraform-provider-google/pull/28207))
* appengine: added `app_engine_bundled_services` field to `google_app_engine_standard_app_version` resource ([#28213](https://github.com/hashicorp/terraform-provider-google/pull/28213))
* biglakeiceberg: add support for `CATALOG_TYPE_FEDERATED` with `federated_catalog_options` to `google_biglake_iceberg_catalog` ([#28241](https://github.com/hashicorp/terraform-provider-google/pull/28241))
* compute: added `target_type` and `target_forwarding_rules` to `google_compute_region_network_firewall_policy_with_rules` resource ([#28061](https://github.com/hashicorp/terraform-provider-google/pull/28061))
* compute: added `workload_identity_config` fields to `google_compute_instance` and `google_compute_instance_template` resources ([#28266](https://github.com/hashicorp/terraform-provider-google/pull/28266))
* compute: added `instance_lifecycle_policy.on_repair.allow_changing_zone` field to `google_compute_instance_group_manager` and `google_compute_instance_region_group_manager`. ([#28174](https://github.com/hashicorp/terraform-provider-google/pull/28174))
* container: added `ANY_RESERVATION_THEN_FAIL` option to `consume_reservation_type` field in `google_container_cluster` and `google_container_node_pool` resources ([#28060](https://github.com/hashicorp/terraform-provider-google/pull/28060))
* container: added `custom_node_init` configuration block to `node_config` (supporting Cloud Storage and Secret Manager) for both `google_container_cluster` and `google_container_node_pool` resources. ([#28262](https://github.com/hashicorp/terraform-provider-google/pull/28262))
* container: added `maintenance_policy` field to `google_container_node_pool` resource ([#28217](https://github.com/hashicorp/terraform-provider-google/pull/28217))
* container: added `recurring_maintenance_window` field to `google_container_cluster` resource ([#28227](https://github.com/hashicorp/terraform-provider-google/pull/28227))
* dataplex: added `catalog_publishing_enabled` field to `google_dataplex_datascan` resource ([#28232](https://github.com/hashicorp/terraform-provider-google/pull/28232))
* dlp: added `inspect_config.min_likelihood_per_info_type` to `google_data_loss_prevention_inspect_template` ([#28236](https://github.com/hashicorp/terraform-provider-google/pull/28236))
* firestore: added `skip_wait` field to `google_firestore_field` resource, skipping the wait for index creation ([#28222](https://github.com/hashicorp/terraform-provider-google/pull/28222))
* oracledatabase: added support for configuring Exascale-based VM clusters on top of dedicated storage vaults via the `exascale_db_storage_vault` parameter and `storage_management_type` to determine if VM Cluster is ASM or EXASCALE in `google_oracle_database_cloud_vm_cluster` ([#28197](https://github.com/hashicorp/terraform-provider-google/pull/28197))
* oracledatabase: added `exadata_infrastructure` field to `google_oracle_database_exascale_db_storage_vault` resource ([#28177](https://github.com/hashicorp/terraform-provider-google/pull/28177))
* oracledatabase: added `exascale_db_storage_vault` and `storage_management_type` fields to `google_oracle_database_cloud_vm_clusters` data source ([#28260](https://github.com/hashicorp/terraform-provider-google/pull/28260))
* sql: added `enforce_new_sql_network_architecture` field to `google_sql_database_instance` resource ([#28233](https://github.com/hashicorp/terraform-provider-google/pull/28233))
* sql: added `psc_auto_connection_policy_enabled` field and output-only `service_connection_policy` and `service_connection_policy_creation_result` fields to `google_sql_database_instance` resource ([#28225](https://github.com/hashicorp/terraform-provider-google/pull/28225))
* storagetransfer: added `private_network_service` to resource `google_storage_transfer_job` ([#28178](https://github.com/hashicorp/terraform-provider-google/pull/28178))

BUG FIXES:
* bigtable: fixed a bug where `row_affinity` updates did not persist on `google_bigtable_app_profile` ([#28215](https://github.com/hashicorp/terraform-provider-google/pull/28215))
* bug: fixed labels diff in `google_project` ([#28229](https://github.com/hashicorp/terraform-provider-google/pull/28229))
* chronicle: suppressed a permadiff on `google_chronicle_rule.text` caused by the Chronicle API appending a trailing newline to every stored rule body ([#28216](https://github.com/hashicorp/terraform-provider-google/pull/28216))
* cloudscheudler: added retries for "409: sync mutate calls cannot be queued" error for `google_cloud_scheduler_job` ([#28164](https://github.com/hashicorp/terraform-provider-google/pull/28164))
* compute: fixed an issue where `preview = false` updates for `google_compute_organization_security_policy_rule` were omitted from API requests. ([#28223](https://github.com/hashicorp/terraform-provider-google/pull/28223))
* compute: fixed bug where it wasn't possible to disable `enable_proxy_protocol` on `google_compute_service_attachment` resource ([#28264](https://github.com/hashicorp/terraform-provider-google/pull/28264))
* datastream: fixed a bug in update functionality in `google_datastream_connection_profile` `mongodb_profile.additional_options` ([#28254](https://github.com/hashicorp/terraform-provider-google/pull/28254))
* eventarc: fixed a type mismatch when an `google_eventarc_trigger` resource returns non-empty `conditions`. ([#28226](https://github.com/hashicorp/terraform-provider-google/pull/28226))
* filestore: aligned `google_filestore_instance` resource timeouts with the Filestore service instance operations TTLs ([#28208](https://github.com/hashicorp/terraform-provider-google/pull/28208))
* gemini: fixed truncated timeouts in `google_gemini_code_tools_setting`, `google_gemini_data_sharing_with_google_setting_binding`, `google_gemini_gemini_gcp_enablement_setting_binding`, and `google_gemini_release_channel_setting_binding` ([#28220](https://github.com/hashicorp/terraform-provider-google/pull/28220))
* hypercomputecluster: fixed 20-minute timeout limit during `google_hypercomputecluster_cluster` resource creation ([#28182](https://github.com/hashicorp/terraform-provider-google/pull/28182))
* logging: fixed an issue where errors on update would not be propagated in `google_logging_project_bucket_config` ([#28055](https://github.com/hashicorp/terraform-provider-google/pull/28055))
* observability: fixed unintentionally long timeouts in `google_observability_folder_settings`, `google_observability_organization_settings`, and `google_observability_project_settings` ([#28220](https://github.com/hashicorp/terraform-provider-google/pull/28220))
* oracledatabase: fixed early client-side timeouts and aligned default schema timeouts with backend async polling limits on `google_oracle_database_exadb_vm_cluster`, `google_oracle_database_odb_network`, `google_oracle_database_odb_subnet`, `google_oracle_database_goldengate_connection`, `google_oracle_database_goldengate_deployment`, and `google_oracle_database_goldengate_connection_assignment`. ([#28228](https://github.com/hashicorp/terraform-provider-google/pull/28228))
* oracledatabase: fixed truncated timeouts in `google_oracle_database_exadb_vm_cluster`, `google_oracle_database_goldengate_connection`, and `google_oracle_database_odb_subnet` ([#28220](https://github.com/hashicorp/terraform-provider-google/pull/28220))
* privilegedaccessmanager: fixed a permadiff on `google_privileged_access_manager_entitlement` for entitlements created without an approval workflow ([#28224](https://github.com/hashicorp/terraform-provider-google/pull/28224))
* provider: fixed validation of `external_credentials.identity_token` to reject malformed JWTs containing empty segments ([#28258](https://github.com/hashicorp/terraform-provider-google/pull/28258))

## 7.39.0 (June 30, 2026)
NOTES:
* compute: migrated `google_compute_instance_template` resource partially to use direct HTTP rather than a client library ([#28010](https://github.com/hashicorp/terraform-provider-google/pull/28010))
* compute: migrated `google_compute_network_peering` resource to use direct HTTP rather than a client library ([#28021](https://github.com/hashicorp/terraform-provider-google/pull/28021))
* compute: migrated metadata handling to use direct HTTP rather than a client library ([#27968](https://github.com/hashicorp/terraform-provider-google/pull/27968))

FEATURES:
* **New Data Source:** `google_agent_registry_agent` ([#28028](https://github.com/hashicorp/terraform-provider-google/pull/28028))
* **New Data Source:** `google_agent_registry_endpoint` ([#28028](https://github.com/hashicorp/terraform-provider-google/pull/28028))
* **New Data Source:** `google_agent_registry_mcp_server` ([#28028](https://github.com/hashicorp/terraform-provider-google/pull/28028))
* **New Data Source:** `google_compute_instance_groups` ([#27981](https://github.com/hashicorp/terraform-provider-google/pull/27981))
* **New Data Source:** `google_storage_control_folder_intelligence_findings_summary` ([#28019](https://github.com/hashicorp/terraform-provider-google/pull/28019))
* **New Data Source:** `google_storage_control_organization_intelligence_findings_summary` ([#28019](https://github.com/hashicorp/terraform-provider-google/pull/28019))
* **New Data Source:** `google_storage_control_project_intelligence_findings_summary` ([#28019](https://github.com/hashicorp/terraform-provider-google/pull/28019))
* **New Resource:** `google_agent_registry_binding` ([#28028](https://github.com/hashicorp/terraform-provider-google/pull/28028))
* **New Resource:** `google_agent_registry_service` ([#28028](https://github.com/hashicorp/terraform-provider-google/pull/28028))
* **New Resource:** `google_artifact_registry_project_config` ([#28009](https://github.com/hashicorp/terraform-provider-google/pull/28009))
* **New Resource:** `google_chronicle_findings_refinement` ([#28035](https://github.com/hashicorp/terraform-provider-google/pull/28035))
* **New Resource:** `google_compute_bulk_per_instance_config` ([#28031](https://github.com/hashicorp/terraform-provider-google/pull/28031))
* **New Resource:** `google_compute_firewall_policy_iam_binding` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_firewall_policy_iam_member` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_firewall_policy_iam_policy` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_network_firewall_policy_iam_binding` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_network_firewall_policy_iam_member` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_network_firewall_policy_iam_policy` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_region_network_firewall_policy_iam_binding` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_region_network_firewall_policy_iam_member` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_region_network_firewall_policy_iam_policy` ([#27978](https://github.com/hashicorp/terraform-provider-google/pull/27978))
* **New Resource:** `google_compute_region_resize_request` ([#27984](https://github.com/hashicorp/terraform-provider-google/pull/27984))
* **New Resource:** `google_compute_zone_vm_extension_policy` ([#28034](https://github.com/hashicorp/terraform-provider-google/pull/28034))
* **New Resource:** `google_gke_hub_rollout_sequence` ([#28007](https://github.com/hashicorp/terraform-provider-google/pull/28007))
* **New Resource:** `google_iap_agent_registry_iam_binding` ([#28032](https://github.com/hashicorp/terraform-provider-google/pull/28032))
* **New Resource:** `google_iap_agent_registry_iam_member` ([#28032](https://github.com/hashicorp/terraform-provider-google/pull/28032))
* **New Resource:** `google_iap_agent_registry_iam_policy` ([#28032](https://github.com/hashicorp/terraform-provider-google/pull/28032))
* **New Resource:** `google_iap_location_web_iam_binding` ([#28032](https://github.com/hashicorp/terraform-provider-google/pull/28032))
* **New Resource:** `google_iap_location_web_iam_member` ([#28032](https://github.com/hashicorp/terraform-provider-google/pull/28032))
* **New Resource:** `google_iap_location_web_iam_policy` ([#28032](https://github.com/hashicorp/terraform-provider-google/pull/28032))
* **New Resource:** `google_oracle_database_cloud_exadata_infrastructure_exascale_config` ([#28033](https://github.com/hashicorp/terraform-provider-google/pull/28033))
* **New List Resource:** `google_bigquery_dataset` ([#28005](https://github.com/hashicorp/terraform-provider-google/pull/28005))
* **New List Resource:** `google_compute_cross_site_network` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_external_vpn_gateway` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_global_network_endpoint_group` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_ha_vpn_gateway` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_interconnect_attachment_group` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_interconnect_group` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_public_delegated_prefix` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_region_commitment` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_region_network_endpoint_group` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_vpn_gateway` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_compute_wire_group` ([#28018](https://github.com/hashicorp/terraform-provider-google/pull/28018))
* **New List Resource:** `google_folder_iam_member` ([#27993](https://github.com/hashicorp/terraform-provider-google/pull/27993))
* **New List Resource:** `google_kms_crypto_key_version` ([#28006](https://github.com/hashicorp/terraform-provider-google/pull/28006))
* **New List Resource:** `google_project` ([#28041](https://github.com/hashicorp/terraform-provider-google/pull/28041))
* **New List Resource:** `google_project_service` ([#27989](https://github.com/hashicorp/terraform-provider-google/pull/27989))

IMPROVEMENTS:
* bigquery: added `external_runtime_options.container_request_concurrency` field to `google_bigquery_routine` resource ([#28029](https://github.com/hashicorp/terraform-provider-google/pull/28029))
* compute: added `instance_lifecycle_policy.on_failed_health_check` field in resources `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` (ga) ([#27992](https://github.com/hashicorp/terraform-provider-google/pull/27992))
* container: added new fields `shutdown_grace_period_seconds` and `shutdown_grace_period_critical_pods_seconds` to `node_kubelet_config` block. ([#28015](https://github.com/hashicorp/terraform-provider-google/pull/28015))
* container: promoted `agent_sandbox_config` addon field under `addons_config` in `google_container_cluster` to GA ([#28017](https://github.com/hashicorp/terraform-provider-google/pull/28017))
* dataplex: added `icon` field to `google_dataplex_data_product` resource ([#27986](https://github.com/hashicorp/terraform-provider-google/pull/27986))
* dataplex: added `name` field to `google_dataplex_data_product_data_asset` resource ([#28020](https://github.com/hashicorp/terraform-provider-google/pull/28020))
* dlp: added `allow_limited_availability_info_types` to `google_data_loss_prevention_inspect_template` ([#28024](https://github.com/hashicorp/terraform-provider-google/pull/28024))
* networkservices: added `forward_attributes` field to `google_network_services_lb_edge_extension`, `google_network_services_lb_route_extension`, and `google_network_services_lb_traffic_extension` resources ([#28012](https://github.com/hashicorp/terraform-provider-google/pull/28012))

BUG FIXES:
* compute: fixed a panic in `google_compute_project_metadata` and `google_compute_project_metadata_item` when project common instance metadata items contain null/empty values. ([#28008](https://github.com/hashicorp/terraform-provider-google/pull/28008))
* compute: fixed a validation error on `google_compute_instance` (`Provisioned IOPS cannot be specified with disk type pd-balanced`) that occurred during updates on instances with Hyperdisk Balanced boot disks. ([#27975](https://github.com/hashicorp/terraform-provider-google/pull/27975))
* dataproc: fixed a bug where changing `policy_id` on `google_dataproc_autoscaling_policy` planned an in-place update and failed; it now correctly forces resource replacement (destroy and recreate). ([#28036](https://github.com/hashicorp/terraform-provider-google/pull/28036))
* firestore: added retries on 409 errors in `google_firestore_user_creds` resource ([#27972](https://github.com/hashicorp/terraform-provider-google/pull/27972))
* iam: fixed ephemeral `google_service_account_key` producing a 404 due to duplicate `/keys` in the URL when `fetch_key = true` and `name` is provided ([#27980](https://github.com/hashicorp/terraform-provider-google/pull/27980))

## 7.38.0 (June 23, 2026)
NOTES:
* bigquery: migrated `google_bigquery_table` resource to use direct HTTP rather than a client library ([#27909](https://github.com/hashicorp/terraform-provider-google/pull/27909))
* compute: migrated `resource_compute_instance_template_test.go.tmpl` resource to use direct HTTP rather than a client library ([#27859](https://github.com/hashicorp/terraform-provider-google/pull/27859))
* compute: migrated `google_compute_resource_compute_instance`  to use direct HTTP rather than a client library ([#27925](https://github.com/hashicorp/terraform-provider-google/pull/27925))
* compute: migrated parts of `google_compute_instance` and shared instance functions to use direct HTTP ([#27815](https://github.com/hashicorp/terraform-provider-google/pull/27815))

FEATURES:
* **New Data Source:** `google_storage_control_project_intelligence_finding_revision` ([#27912](https://github.com/hashicorp/terraform-provider-google/pull/27912))
* **New Data Source:** `google_storage_control_project_intelligence_finding_revisions` ([#27912](https://github.com/hashicorp/terraform-provider-google/pull/27912))
* **New Resource:** `google_biglake_hive_catalog` ([#27892](https://github.com/hashicorp/terraform-provider-google/pull/27892))
* **New Resource:** `google_chronicle_feed` ([#27860](https://github.com/hashicorp/terraform-provider-google/pull/27860))
* **New Resource:** `google_chronicle_parser_extension` ([#27906](https://github.com/hashicorp/terraform-provider-google/pull/27906))
* **New Resource:** `google_dataplex_metadata_feed` ([#27934](https://github.com/hashicorp/terraform-provider-google/pull/27934))
* **New Resource:** `google_network_services_agent_gateway` ([#27803](https://github.com/hashicorp/terraform-provider-google/pull/27803))
* **New Resource:** `google_vertex_ai_schedule` ([#27895](https://github.com/hashicorp/terraform-provider-google/pull/27895))
* **New Resource:** `google_vertex_ai_tensorboard_run` ([#27913](https://github.com/hashicorp/terraform-provider-google/pull/27913))
* **New List Resource:** `google_compute_address` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_cross_site_network` ([#27864](https://github.com/hashicorp/terraform-provider-google/pull/27864))
* **New List Resource:** `google_compute_https_health_check` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_node_template` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_packet_mirroring` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_region_autoscaler` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_region_composite_health_check` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_region_health_aggregation_policy` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_compute_region_health_source` ([#27917](https://github.com/hashicorp/terraform-provider-google/pull/27917))
* **New List Resource:** `google_project_iam_member` ([#27905](https://github.com/hashicorp/terraform-provider-google/pull/27905))
* **New List Resource:** `google_pubsub_topic` ([#27914](https://github.com/hashicorp/terraform-provider-google/pull/27914))
* **New List Resource:** `google_secret_manager_secret` ([#27910](https://github.com/hashicorp/terraform-provider-google/pull/27910))

IMPROVEMENTS:
* apigee: added `consumer_key` and `consumer_secret` fields to `google_apigee_developer_app` to allow specifying a static credential ([#27820](https://github.com/hashicorp/terraform-provider-google/pull/27820))
* artifactregistry: added update support for `upstream_credentials` to `google_artifact_registry_repository` ([#27819](https://github.com/hashicorp/terraform-provider-google/pull/27819))
* biglakeiceberg: added `CATALOG_TYPE_BIGLAKE` enum to `catalog_type` field and added `restricted_locations_config.restricted_locations` field in `google_biglake_iceberg_catalog ` resource ([#27930](https://github.com/hashicorp/terraform-provider-google/pull/27930))
* biglakeiceberg: added `sort_order` field to `google_biglake_iceberg_table` resource ([#27865](https://github.com/hashicorp/terraform-provider-google/pull/27865))
* ces: added `timeout` and `tool_fake_config` fields to `google_ces_tool` and `google_ces_toolset` resource ([#27907](https://github.com/hashicorp/terraform-provider-google/pull/27907))
* compute: added `params.resource_manager_tags` field to `google_compute_snapshot` resource ([#27869](https://github.com/hashicorp/terraform-provider-google/pull/27869))
* compute: made `network_endpoints.ip_address` optional in `google_compute_network_endpoints` resource to support attaching endpoints to a network endpoint group of type `GCE_VM_IP_DEDICATED_BACKEND` ([#27870](https://github.com/hashicorp/terraform-provider-google/pull/27870))
* container: added `dataplane_optimization_mode` in `google_container_cluster` ([#27861](https://github.com/hashicorp/terraform-provider-google/pull/27861))
* container: added `ignore_node_count_changes` field to `google_container_cluster` and `google_container_node_pool` resources. When set to true, the provider ignores drift via external node count changes and skips related IGM API queries, resolving long plan times on clusters with a large number of instance groups. ([#27896](https://github.com/hashicorp/terraform-provider-google/pull/27896))
* container: added `skip_node_pool_refresh` field to `google_container_cluster` resource. When set to true, the `google_container_cluster` skips refreshing and setting `node_pools` from the API, resolving long plan times on clusters with a large number of node pools. Note that this results in `node_pools` being set to an empty list in state ([#27896](https://github.com/hashicorp/terraform-provider-google/pull/27896))
* container: added `taint_config` block to `google_container_cluster` and `google_container_node_pool` ([#27884](https://github.com/hashicorp/terraform-provider-google/pull/27884))
* container: improved GKE node pool read performance by caching instance group metadata longer ([#27896](https://github.com/hashicorp/terraform-provider-google/pull/27896))
* datastream: added `additional_options` field to `google_datastream_connection_profile` resource ([#27915](https://github.com/hashicorp/terraform-provider-google/pull/27915))
* iamworkforcepool: write-only support for `oidc.client_secret` in `google_iam_workforce_pool_provider` ([#27867](https://github.com/hashicorp/terraform-provider-google/pull/27867))
* kms: added resource identity support for `google_kms_crypto_key_version` resource ([#27883](https://github.com/hashicorp/terraform-provider-google/pull/27883))
* networkservices: added `dns_peering_config` field to `google_network_services_agent_gateway` resource ([#27813](https://github.com/hashicorp/terraform-provider-google/pull/27813))
* sql: added `mode`, `dns_servers`, `admin_credential_secret_name`, and `organizational_unit` fields to `active_directory_config` block in `google_sql_database_instance` resource for SQL Server instances ([#27862](https://github.com/hashicorp/terraform-provider-google/pull/27862))
* storage: added `lifecycle_rule.condition.size_above_bytes` and `lifecycle_rule.condition.size_below_bytes` fields to `google_storage_bucket` resource ([#27857](https://github.com/hashicorp/terraform-provider-google/pull/27857))

BUG FIXES:
* apigee: `google_apigee_developer_app` now updates `api_products` and `scopes` on the existing credential instead of creating a new credential (consumer key) on update ([#27929](https://github.com/hashicorp/terraform-provider-google/pull/27929))
* biglake: allow `location` to be set on `google_biglake_iceberg_namespace` ([#27814](https://github.com/hashicorp/terraform-provider-google/pull/27814))
* biglake: fixed creation failure of `google_biglake_iceberg_table` resource when the referenced `google_biglake_iceberg_catalog` has `credential_mode` set to `CREDENTIAL_MODE_VENDED_CREDENTIALS` due to a missing `X-Iceberg-Access-Delegation` header ([#27903](https://github.com/hashicorp/terraform-provider-google/pull/27903))
* compute: fixed broken import of `share_settings` on `google_compute_reservation` ([#27916](https://github.com/hashicorp/terraform-provider-google/pull/27916))
* datastream: fixed a positional diff when adding objects to the `salesforce_source_config.include_objects` field in `google_datastream_stream` resource ([#27926](https://github.com/hashicorp/terraform-provider-google/pull/27926))
* iamworkforcepool: marked sensitive and ignore_read as true for `security_token` in `google_iam_workforce_pool_provider_scim_token` resource ([#27812](https://github.com/hashicorp/terraform-provider-google/pull/27812))
* networkconnectivity: fixed `google_network_connectivity_regional_endpoint` being recreated on every apply when `address` is set to a resource URI ([#27923](https://github.com/hashicorp/terraform-provider-google/pull/27923))
* networkservices: fixed `name` field expansion for `google_network_services_agent_gateway` resources so that short names are automatically expanded to full resource names, preventing API validation errors on create and update. ([#27902](https://github.com/hashicorp/terraform-provider-google/pull/27902))

## 7.37.0 (June 16, 2026)

NOTES:
* compute: migrated `EnableDisplay` fields in `google_compute_instance` resources to use direct HTTP rather than a client library ([#27778](https://github.com/hashicorp/terraform-provider-google/pull/27778))
* compute: migrated `desired_status` block and `startInstanceOperation` in `resource_compute_instance.go.tmpl` to use direct HTTP rather than a client library ([#27755](https://github.com/hashicorp/terraform-provider-google/pull/27755))
* compute: migrated `getInstance`, `getDisk`, `Delete` and the `setMetadata` update block in `resource_compute_instance.go.tmpl` to use direct HTTP rather than a client library ([#27716](https://github.com/hashicorp/terraform-provider-google/pull/27716))
* compute: migrated part of `google_compute_instance` to use direct HTTP rather than a client library ([#27788](https://github.com/hashicorp/terraform-provider-google/pull/27788))

DEPRECATIONS:
* cloudrunv2: deprecated `http_get.http_headers.port` field in container startup probe and liveness probe in `google_cloud_run_v2_worker_pool` resource because it is not supported in Cloud Run API. This field will be removed in a future major release. ([#27800](https://github.com/hashicorp/terraform-provider-google/pull/27800))
* cloudsecuritycompliance: deprecated the `organization` field on `google_cloud_security_compliance_cloud_control`, `google_cloud_security_compliance_framework`, and `google_cloud_security_compliance_framework_deployment`. Use `parent` instead ([#27769](https://github.com/hashicorp/terraform-provider-google/pull/27769))
* networkservices: deprecated `protocols` on `google_network_services_agent_gateway` ([#27802](https://github.com/hashicorp/terraform-provider-google/pull/27802))

FEATURES:
* **New Data Source:** `google_oracle_database_goldengate_deployment_versions` ([#27771](https://github.com/hashicorp/terraform-provider-google/pull/27771))
* **New Data Source:** `google_storage_control_project_intelligence_finding` ([#27764](https://github.com/hashicorp/terraform-provider-google/pull/27764))
* **New Data Source:** `google_storage_control_project_intelligence_findings` ([#27764](https://github.com/hashicorp/terraform-provider-google/pull/27764))
* **New Resource:** `google_chronicle_parser` ([#27801](https://github.com/hashicorp/terraform-provider-google/pull/27801))
* **New Resource:** `google_migration_center_import_data_file` ([#27721](https://github.com/hashicorp/terraform-provider-google/pull/27721))
* **New Resource:** `google_network_services_agent_gateway` ([#27803](https://github.com/hashicorp/terraform-provider-google/pull/27803))
* **New Resource:** `google_vertex_ai_tensorboard_experiment` ([#27796](https://github.com/hashicorp/terraform-provider-google/pull/27796))
* **New List Resource:** `google_bigquery_dataset_access` ([#27758](https://github.com/hashicorp/terraform-provider-google/pull/27758))
* **New List Resource:** `google_cloud_scheduler_job` ([#27758](https://github.com/hashicorp/terraform-provider-google/pull/27758))
* **New List Resource:** `google_dns_record_set` ([#27792](https://github.com/hashicorp/terraform-provider-google/pull/27792))
* **New List Resource:** `google_monitoring_alert_policy` ([#27758](https://github.com/hashicorp/terraform-provider-google/pull/27758))
* **New List Resource:** `google_pubsub_subscription` ([#27758](https://github.com/hashicorp/terraform-provider-google/pull/27758))

IMPROVEMENTS:
* apigee: added new resource `google_apigee_environment_debugmask` for managing Apigee environment debug masks ([#27719](https://github.com/hashicorp/terraform-provider-google/pull/27719))
* backupdr: added support for `use_project_service_account` flag in `google_backup_dr_restore_workload` disk and compute restores ([#27797](https://github.com/hashicorp/terraform-provider-google/pull/27797))
* cloudrunv2: added `http_get.http_headers.name` field to container startup probe and liveness probe in `google_cloud_run_v2_worker_pool` resource ([#27800](https://github.com/hashicorp/terraform-provider-google/pull/27800))
* cloudrunv2: added `template.client` and `template.client_version` fields to `google_cloud_run_v2_worker_pool` resource ([#27757](https://github.com/hashicorp/terraform-provider-google/pull/27757))
* cloudsecuritycompliance: added support for project parent to `google_cloud_security_compliance_cloud_control`, `google_cloud_security_compliance_framework`, and `google_cloud_security_compliance_framework_deployment` via the new `parent` field. The `organization` field has been deprecated. ([#27769](https://github.com/hashicorp/terraform-provider-google/pull/27769))
* compute: added `params.resource_manager_tags` field to `google_compute_reservation` resource ([#27770](https://github.com/hashicorp/terraform-provider-google/pull/27770))
* compute: added data sources for `google_compute_target_http_proxy`, `google_compute_target_https_proxy`, `google_compute_region_target_http_proxy`, and `google_compute_region_target_https_proxy` ([#27767](https://github.com/hashicorp/terraform-provider-google/pull/27767))
* container: added `addons_config.slurm_operator_config` field to `google_container_cluster` resource ([#27765](https://github.com/hashicorp/terraform-provider-google/pull/27765))
* container: added `node_image_config` field to `google_container_node_pool` and `google_container_cluster` resources ([#27794](https://github.com/hashicorp/terraform-provider-google/pull/27794))
* databasemigrationservice: added `state` and `stop_on_warnings` fields to `google_database_migration_service_migration_job` resource ([#27731](https://github.com/hashicorp/terraform-provider-google/pull/27731))
* dns: added resource identity support for `google_dns_record_set` resource ([#27792](https://github.com/hashicorp/terraform-provider-google/pull/27792))
* networksecurity: added `network_rules` field on `google_network_security_authz_policy` resource ([#27821](https://github.com/hashicorp/terraform-provider-google/pull/27821))
* pubsub: added `first_revision_id` and `last_revision_id` fields to `google_pubsub_topic` resource ([#27718](https://github.com/hashicorp/terraform-provider-google/pull/27718))
* sql: added `settings.ip_configuration.psc_config.psc_auto_dns_enabled` and `settings.ip_configuration.psc_config.psc_write_endpoint_dns_enabled` fields to `google_sql_database_instance` resource ([#27776](https://github.com/hashicorp/terraform-provider-google/pull/27776))

BUG FIXES:
* apigee: fixed `google_apigee_api` not detecting local bundle changes due to a missing default on `detect_md5hash`, and fixed the test sweeper's list URL ([#27791](https://github.com/hashicorp/terraform-provider-google/pull/27791))
* apigee: fixed `google_apigee_security_action` update failure by enabling PATCH-based updates now that the Apigee Security Actions API supports mutations ([#27768](https://github.com/hashicorp/terraform-provider-google/pull/27768))
* apigee: fixed a perma-diff for `api_products` and `scopes` fields in `google_apigee_developer_app` resource when updating them with multiple items ([#27789](https://github.com/hashicorp/terraform-provider-google/pull/27789))
* apigee: fixed an issue where the resource would attempt recreation if the `key_expires_in` field was set in `google_apigee_developer_app` resource ([#27779](https://github.com/hashicorp/terraform-provider-google/pull/27779))
* ces: fixed persistent diff in `google_ces_guardrail` when `llm_prompt_security` is configured with `default_settings` ([#27766](https://github.com/hashicorp/terraform-provider-google/pull/27766))
* cloudrun: fixed a permadiff for the `run.googleapis.com/gpu-zonal-redundancy-disabled` annotation in `google_cloud_run_service` ([#27787](https://github.com/hashicorp/terraform-provider-google/pull/27787))
* cloudrunv2: fixed bug where only one `http_get.http_headers` block could be specified in container startup probe and liveness probe in `google_cloud_run_v2_worker_pool` resource ([#27800](https://github.com/hashicorp/terraform-provider-google/pull/27800))
* compute: fixed an issue in `google_compute_subnetwork` where `secondary_ip_range` entries linked to an `internal_range` could not be removed and adding new ranges would sometimes fail due to positional shifts ([#27175](https://github.com/hashicorp/terraform-provider-google/issues/27175)) ([#27720](https://github.com/hashicorp/terraform-provider-google/pull/27720))
* compute: fixed diff when using `existing_reservations` field in `google_region_commitment` ([#27775](https://github.com/hashicorp/terraform-provider-google/pull/27775))
* compute: fixed rules in `google_compute_security_policy` being unnecessarily recreated due to TypeSet hash instability ([#27754](https://github.com/hashicorp/terraform-provider-google/pull/27754))
* sql: fixed inconsistent result after apply error when adding `users` of type `CLOUD_IAM_GROUP` with capitalized domain names for MySQL ([#27784](https://github.com/hashicorp/terraform-provider-google/pull/27784))
* storage: fixed OOM issue for `google_storage_bucket` `force_destroy` by limiting the number of outstanding tasks to 2000 ([#27777](https://github.com/hashicorp/terraform-provider-google/pull/27777))

## 7.36.0 (June 9, 2026)

FEATURES:
* **New Data Source:** `google_apigee_instance` ([#27683](https://github.com/hashicorp/terraform-provider-google/pull/27683))
* **New Data Source:** `google_oracle_database_goldengate_deployment_types` ([#27634](https://github.com/hashicorp/terraform-provider-google/pull/27634))
* **New Resource:** `google_apigee_datastore` ([#27607](https://github.com/hashicorp/terraform-provider-google/pull/27607))
* **New Resource:** `google_discovery_engine_search_engine_iam_binding` ([#27703](https://github.com/hashicorp/terraform-provider-google/pull/27703))
* **New Resource:** `google_license_manager_configuration` ([#27707](https://github.com/hashicorp/terraform-provider-google/pull/27707))
* **New Resource:** `google_migration_center_import_job` ([#27599](https://github.com/hashicorp/terraform-provider-google/pull/27599))
* **New List Resource:** `google_compute_disk` ([#27608](https://github.com/hashicorp/terraform-provider-google/pull/27608))
* **New List Resource:** `google_compute_image` ([#27608](https://github.com/hashicorp/terraform-provider-google/pull/27608))
* **New List Resource:** `google_compute_snapshot` ([#27608](https://github.com/hashicorp/terraform-provider-google/pull/27608))
* **New List Resource:** `google_storage_hmac_key` ([#27637](https://github.com/hashicorp/terraform-provider-google/pull/27637))

IMPROVEMENTS:
* accesscontextmanager: added in-place update for `egress_from` and `egress_to` fields in `google_access_context_manager_service_perimeter_egress_policy` resource ([#27690](https://github.com/hashicorp/terraform-provider-google/pull/27690))
* accesscontextmanager: added in-place update for `egress_from` and `egress_to` fields in `google_access_context_manager_service_perimeter_ingress_policy` resource ([#27690](https://github.com/hashicorp/terraform-provider-google/pull/27690))
* bigquery: added IAM support (`google_bigquery_routine_iam_policy`, `google_bigquery_routine_iam_binding`, `google_bigquery_routine_iam_member`) for `google_bigquery_routine` resource ([#27704](https://github.com/hashicorp/terraform-provider-google/pull/27704))
* bigtable: added `automated_backup_policy.locations` field in `google_bigtable_table` resource ([#27646](https://github.com/hashicorp/terraform-provider-google/pull/27646))
* ces: added `agent_tool`, `file_search_tool`, and `widget_tool` fields to the `google_ces_tool` resource ([#27681](https://github.com/hashicorp/terraform-provider-google/pull/27681))
* ces: added `google_search_tool.prompt_config` and `data_store_tool.data_store_source` fields to the `google_ces_tool` resource ([#27681](https://github.com/hashicorp/terraform-provider-google/pull/27681))
* ces: exposed `remote_agent_tool`, `connector_tool`, and `mcp_tool` as read-only (output-only) attributes in `google_ces_tool` ([#27681](https://github.com/hashicorp/terraform-provider-google/pull/27681))
* container: added `node_creation_config` field to `google_container_cluster` resource ([#27702](https://github.com/hashicorp/terraform-provider-google/pull/27702))
* container: added `node_drain_config.pdb_timeout_duration`  and `node_drain_config.grace_termination_duration` fields to `google_container_node_pool` and `google_container_cluster` resources ([#27694](https://github.com/hashicorp/terraform-provider-google/pull/27694))
* data_catalog: added `RICHTEXT` to allowed values of `primitive_type` on `google_data_catalog_tag_template` fields. ([#27672](https://github.com/hashicorp/terraform-provider-google/pull/27672))
* dataplex: added IAM support for `google_dataplex_data_product` resource (`iam_policy`, `iam_binding`, `iam_member`) ([#27652](https://github.com/hashicorp/terraform-provider-google/pull/27652))
* dataplex: added `access_approval_config` field to `google_dataplex_data_product` resource ([#27652](https://github.com/hashicorp/terraform-provider-google/pull/27652))
* hypercomputecluster: marked `network_resources` field as required in `google_hypercomputecluster_cluster` resource to align with API validation ([#27655](https://github.com/hashicorp/terraform-provider-google/pull/27655))
* networksecurity: `google_network_security_ull_mirroring_engine`, `google_network_security_ull_mirroring_collector`, and `google_network_security_ull_mirroring_collector_rule` resources promoted to GA ([#27710](https://github.com/hashicorp/terraform-provider-google/pull/27710))
* securesourcemanager: added `psc_allowed_projects` field to `google_secure_source_manager_instance` resource ([#27695](https://github.com/hashicorp/terraform-provider-google/pull/27695))
* workbench: added `NVIDIA_RTX6000` to the supported `gce_setup.accelerator_configs.type` values on `google_workbench_instance` resource([#27709](https://github.com/hashicorp/terraform-provider-google/pull/27709))

BUG FIXES:
* apigee: send zero values for `ip_header_index` in `google_apigee_environment` resource ([#27670](https://github.com/hashicorp/terraform-provider-google/pull/27670))
* backupdr: fixed an issue where `google_backup_dr_restore_workload` did not use the correct API JSON names for networking/reservation fields ([#27680](https://github.com/hashicorp/terraform-provider-google/pull/27680))
* compute: fixed an issue where updating `connection_limit` in the `consumer_accept_lists` block of `google_compute_service_attachment` would not trigger a resource update. ([#27688](https://github.com/hashicorp/terraform-provider-google/pull/27688))
* compute: fixed regional backend reference in  `google_compute_regional_url_map` resource ([#27705](https://github.com/hashicorp/terraform-provider-google/pull/27705))
* dlp: fixed error when reading `google_data_loss_prevention_discovery_config` caused by nested error details ([#27669](https://github.com/hashicorp/terraform-provider-google/pull/27669))
* sql: fixed permadiff on `connection_pool_config` when `connection_pooling_enabled` is set to `false` ([#27711](https://github.com/hashicorp/terraform-provider-google/pull/27711))
* tags: fixed `google_tags_location_tag_binding` failing with `Operation location does not match service location 'global'` during creation ([#27668](https://github.com/hashicorp/terraform-provider-google/pull/27668))
* vertexai: fixed `terraform import` of `google_vertex_ai_index_endpoint_deployed_index` failing with "Cannot determine region" when provider-level `region`/`zone` is unset ([#27692](https://github.com/hashicorp/terraform-provider-google/pull/27692))

## 7.35.0 (June 2, 2026)

FEATURES:
* **New Data Source:** `google_oracle_database_goldengate_connection_types` ([#27567](https://github.com/hashicorp/terraform-provider-google/pull/27567))
* **New Resource:** `google_chronicle_findings_refinement` ([#27591](https://github.com/hashicorp/terraform-provider-google/pull/27591))
* **New Resource:** `google_dataplex_data_product` ([#27588](https://github.com/hashicorp/terraform-provider-google/pull/27588))
* **New Resource:** `google_dataplex_data_product_data_asset` ([#27588](https://github.com/hashicorp/terraform-provider-google/pull/27588))
* **New Resource:** `google_migration_center_discovery_client` ([#27572](https://github.com/hashicorp/terraform-provider-google/pull/27572))
* **New Resource:** `google_migration_center_report` ([#27548](https://github.com/hashicorp/terraform-provider-google/pull/27548))
* **New Resource:** `google_oracle_database_goldengate_connection_assignment` ([#27566](https://github.com/hashicorp/terraform-provider-google/pull/27566))
* **New Resource:** `google_oracle_database_goldengate_connection` ([#27587](https://github.com/hashicorp/terraform-provider-google/pull/27587))
* **New Resource:** `google_oracle_database_goldengate_deployment` ([#27575](https://github.com/hashicorp/terraform-provider-google/pull/27575))
* **New List Resource:** `google_compute_firewall` ([#27549](https://github.com/hashicorp/terraform-provider-google/pull/27549))
* **New List Resource:** `google_compute_global_address` ([#27549](https://github.com/hashicorp/terraform-provider-google/pull/27549))
* **New List Resource:** `google_compute_subnetwork` ([#27549](https://github.com/hashicorp/terraform-provider-google/pull/27549))
* **New List Resource:** `google_sql_database` ([#27552](https://github.com/hashicorp/terraform-provider-google/pull/27552))

IMPROVEMENTS:
* compute: added `target_type` and `target_forwarding_rules` fields to `google_compute_network_firewall_policy_rule` resource ([#27538](https://github.com/hashicorp/terraform-provider-google/pull/27538))
* container: added `crash_loop_back_off.max_container_restart_period` field to `google_container_node_pool` and `google_container_cluster` resources ([#27574](https://github.com/hashicorp/terraform-provider-google/pull/27574))
* container: added additional value `KCP_VPA` for `logging_config.enable_components` field to `google_container_cluster` resource ([#27546](https://github.com/hashicorp/terraform-provider-google/pull/27546))
* dataplex: added `service_account` support to `google_dataplex_data_product` access group principals ([#27588](https://github.com/hashicorp/terraform-provider-google/pull/27588))
* firestore: added `ttl_config.expiration_offset` field to `google_firestore_field` resource ([#27589](https://github.com/hashicorp/terraform-provider-google/pull/27589))
* netapp: added `ontap_source` field to `google_netapp_backup` resource ([#27584](https://github.com/hashicorp/terraform-provider-google/pull/27584))
* networkmanagement: added `gke_pod` and `network_type` fields to `google_network_management_connectivity_test` resource ([#27585](https://github.com/hashicorp/terraform-provider-google/pull/27585))

BUG FIXES:
* resourcemanager: fixed a bug where ephemeral `google_service_account_key` failed on deletion if the parent service account had already been deleted ([#27541](https://github.com/hashicorp/terraform-provider-google/pull/27541))
* storage: fixed missing identity error when updating values in `google_storage_bucket` ([#27605](https://github.com/hashicorp/terraform-provider-google/pull/27605))

## 7.34.0 (May 27, 2026)

NOTES:
* compute: migrated `google_compute_region_instance_template` to use direct HTTP rather than a client library ([#27471](https://github.com/hashicorp/terraform-provider-google/pull/27471))
* compute: migrated `google_compute_instance_group_manager` resource to use direct HTTP rather than a client library ([#27441](https://github.com/hashicorp/terraform-provider-google/pull/27441))

FEATURES:
* **New Data Source:** `google_compute_service_attachment` ([#27526](https://github.com/hashicorp/terraform-provider-google/pull/27526))
* **New Data Source:** `google_oracle_database_goldengate_deployment_environments` ([#27499](https://github.com/hashicorp/terraform-provider-google/pull/27499))
* **New Resource:** `google_config_deployment` ([#27438](https://github.com/hashicorp/terraform-provider-google/pull/27438))
* **New Resource:** `google_dialogflow_sip_trunk` ([#27468](https://github.com/hashicorp/terraform-provider-google/pull/27468))
* **New Resource:** `google_migration_center_assets_export_job` ([#27466](https://github.com/hashicorp/terraform-provider-google/pull/27466))
* **New Resource:** `google_migration_center_report_config` ([#27395](https://github.com/hashicorp/terraform-provider-google/pull/27395))
* **New Resource:** `google_migration_center_settings` ([#27465](https://github.com/hashicorp/terraform-provider-google/pull/27465))
* **New Resource:** `google_migration_center_source` ([#27496](https://github.com/hashicorp/terraform-provider-google/pull/27496))

IMPROVEMENTS:
* bigtable: added `edition` field to `google_bigtable_instance` resource ([#27507](https://github.com/hashicorp/terraform-provider-google/pull/27507))
* ces: added `fail_open` field to `llm_prompt_security` block in `google_ces_guardrail` resource ([#27497](https://github.com/hashicorp/terraform-provider-google/pull/27497))
* ces: added read-only `fail_open` field to `llm_prompt_security` block in `google_ces_app_version` resource ([#27497](https://github.com/hashicorp/terraform-provider-google/pull/27497))
* compute: added `ip_version` and `ip_collection` fields to `secondary_ip_range` field in `google_compute_subnetwork` resource ([#27432](https://github.com/hashicorp/terraform-provider-google/pull/27432))
* compute: added `post_quantum_key_exchange` field to `google_compute_ssl_policy` and `google_compute_region_ssl_policy` resources ([#27479](https://github.com/hashicorp/terraform-provider-google/pull/27479))
* compute: added support in the `google_compute_network` datasource for looking up a network by `self_link` in addition to `name` ([#27509](https://github.com/hashicorp/terraform-provider-google/pull/27509))
* container: added `agent_sandbox_config` field to `google_container_cluster` resource ([#27482](https://github.com/hashicorp/terraform-provider-google/pull/27482))
* container: added `node_config.gpudirect_strategy` and `node_pool.node_config.gpudirect_strategy` to `cluster` resource, added `node_config.gpudirect_strategy` to `node_pool` resource ([#27495](https://github.com/hashicorp/terraform-provider-google/pull/27495))
* dataflow: Added  `create_ignore_already_exists`  field to  `google_dataflow_flex_template_job` resource to handle 409 conflicts ([#27476](https://github.com/hashicorp/terraform-provider-google/pull/27476))
* datafusion: added `maintenance_policy` field to `google_data_fusion_instance` resource ([#27470](https://github.com/hashicorp/terraform-provider-google/pull/27470))
* iam: add resource identity support for `iam_member` resources ([#27383](https://github.com/hashicorp/terraform-provider-google/pull/27383))
* networkconnectivity: `google_network_connectivity_transport` resource promoted to GA ([#27440](https://github.com/hashicorp/terraform-provider-google/pull/27440))
* oracledatabase: added `identity_connector` to `google_oracle_database_cloud_vm_cluster` for CMEK support ([#27435](https://github.com/hashicorp/terraform-provider-google/pull/27435))
* project: added Resource Identity support to `google_project_iam_binding` ([#27502](https://github.com/hashicorp/terraform-provider-google/pull/27502))
* project: added Resource Identity support to `google_project_iam_policy` ([#27503](https://github.com/hashicorp/terraform-provider-google/pull/27503))
* sql: promoted Hyperdisk fields, `data_disk_provisioned_iops` and `data_disk_provisioned_throughput`  to GA ([#27437](https://github.com/hashicorp/terraform-provider-google/pull/27437))

BUG FIXES:
* bigtable: fixed an issue where `bigtable_custom_endpoint` and `universe_domain` were ignored when creating Bigtable resources. ([#27515](https://github.com/hashicorp/terraform-provider-google/pull/27515))
* compute: fixed an issue in `google_compute_subnetwork` where `secondary_ip_range` entries linked to an `internal_range` could not be removed and adding new ranges would sometimes fail due to positional shifts ([#27175](https://github.com/hashicorp/terraform-provider-google/issues/27175)) ([#27512](https://github.com/hashicorp/terraform-provider-google/pull/27512))
* compute: marked encryption keys as immutable and sensitive across compute and backupdr resources ([#27508](https://github.com/hashicorp/terraform-provider-google/pull/27508))
* dialogflow: corrected `AUDIOENCODING_SPEEX_WITH_HEADER_BYTE` enum value to `AUDIO_ENCODING_SPEEX_WITH_HEADER_BYTE` for `audio_encoding` field in `google_dialogflow_conversation_profile` resource ([#27459](https://github.com/hashicorp/terraform-provider-google/pull/27459))
* resourcemanager: resolved a one-time diff for `deletion_policy` that would occur on existing and imported `google_project_service` resources following upgrading to v7.32.0 ([#27484](https://github.com/hashicorp/terraform-provider-google/pull/27484))

## 7.33.0 (May 19, 2026)

NOTES:
* compute: migrated `google_compute_target_pool` resource to use direct HTTP rather than a client library ([#12212](https://github.com/hashicorp/terraform-provider-google-beta/pull/12212))
* compute: migrated `google_compute_instance_group_manager` resource  to use direct HTTP rather than a client library ([#12206](https://github.com/hashicorp/terraform-provider-google-beta/pull/12206))
* compute: migrated `google_compute_project_default_network_tier` resource to use direct HTTP rather than a client library ([#12201](https://github.com/hashicorp/terraform-provider-google-beta/pull/12201))
* compute: migrated `google_compute_router_status` data source to use direct HTTP rather than a client library ([#12174](https://github.com/hashicorp/terraform-provider-google-beta/pull/12174))
* compute: migrated `google_compute_instance_group_manager` resource to use direct HTTP rather than a client library ([#12216](https://github.com/hashicorp/terraform-provider-google-beta/pull/12216))
* compute: partially migrated `google_compute_instance` resource to use direct HTTP rather then a client library ([#12205](https://github.com/hashicorp/terraform-provider-google-beta/pull/12205))

FEATURES:
* **New Data Source:** `google_logging_log_view` ([#12226](https://github.com/hashicorp/terraform-provider-google-beta/pull/12226))
* **New Resource:** `google_apigee_data_collector` ([#12190](https://github.com/hashicorp/terraform-provider-google-beta/pull/12190))
* **New Resource:** `google_chronicle_native_dashboard` (ga) ([#12188](https://github.com/hashicorp/terraform-provider-google-beta/pull/12188))
* **New Resource:** `google_contact_center_insights_encryption_spec` ([#12225](https://github.com/hashicorp/terraform-provider-google-beta/pull/12225))

IMPROVEMENTS:
* backupdr: added `guest_flush` field to `google_backup_dr_backup_plan` resource and `google_backup_dr_backup` data source. ([#12229](https://github.com/hashicorp/terraform-provider-google-beta/pull/12229))
* backupdr: added `guest_flush` field to `google_backup_dr_backup_plan` resource and `google_backup_dr_backup` data source. ([#12230](https://github.com/hashicorp/terraform-provider-google-beta/pull/12230))
* ces: added `security_settings` field to `google_ces_deployment` resource ([#12227](https://github.com/hashicorp/terraform-provider-google-beta/pull/12227))
* ces: added `tool_execution_mode` field to `google_ces_app` resource ([#12221](https://github.com/hashicorp/terraform-provider-google-beta/pull/12221))
* compute: added `stabilization_period` field to `google_compute_autoscaler` and `google_compute_region_autoscaler` resources ([#12232](https://github.com/hashicorp/terraform-provider-google-beta/pull/12232))
* compute: added support for "ARP_BROADCAST_PRIMARY_RANGE" values to the `resolve_subnet_mask` field in `google_compute_subnetwork` resource ([#12176](https://github.com/hashicorp/terraform-provider-google-beta/pull/12176))
* compute: added support for "GCE_VM_IP_DEDICATED_BACKEND" to the  `network_endpoint_type` field in `google_compute_network_endpoint_group` resource ([#12176](https://github.com/hashicorp/terraform-provider-google-beta/pull/12176))
* compute: migrated `data_source_google_compute_regions` to use direct HTTP rather than a client library ([#12202](https://github.com/hashicorp/terraform-provider-google-beta/pull/12202))
* container: added `pod_snapshot_config` field to `google_container_cluster` resource (GA) ([#12196](https://github.com/hashicorp/terraform-provider-google-beta/pull/12196))
* container: added `secret_sync_config` field to `google_container_cluster` resource (ga) ([#12215](https://github.com/hashicorp/terraform-provider-google-beta/pull/12215))
* databasemigrationservice: added `database` and `private_connectivity` fields to `google_database_migration_service_connection_profile` resource ([#12203](https://github.com/hashicorp/terraform-provider-google-beta/pull/12203))
* databasemigrationservice: added `postgres_homogeneous_config` field to `google_database_migration_service_migration_job` resource ([#12203](https://github.com/hashicorp/terraform-provider-google-beta/pull/12203))
* databasemigrationservice: added `psc_interface_config` field to `google_database_migration_service_private_connection` resource ([#12184](https://github.com/hashicorp/terraform-provider-google-beta/pull/12184))
* hypercomputecluster: added `terminal_storage_class` and `per_unit_storage_throughput` fields to the `google_hypercomputecluster_cluster` resource ([#12234](https://github.com/hashicorp/terraform-provider-google-beta/pull/12234))
* netapp: added `ontap_source` field to `google_netapp_backup` resource (beta) ([#12231](https://github.com/hashicorp/terraform-provider-google-beta/pull/12231))
* provider: support for a `deletion_policy` field has been added to almost all resources in the provider. Details on its usage can be found within individual resource documentation if supported. ([#12183](https://github.com/hashicorp/terraform-provider-google-beta/pull/12183))
* storagebatchoperations: added `description` field to `google_storage_batch_operations_job` resource ([#12207](https://github.com/hashicorp/terraform-provider-google-beta/pull/12207))
* workstations: added `workstation_authorization_url` and `workstation_launch_url` fields to the `google_workstations_workstation_cluster ` resource. ([#12185](https://github.com/hashicorp/terraform-provider-google-beta/pull/12185))

BUG FIXES:
* apigee: fixed forced replacement when importing `google_apigee_sharedflow_deployment` resource, where `service_account` read as null ([#12228](https://github.com/hashicorp/terraform-provider-google-beta/pull/12228))
* bigqueryconnection: fixed an issue where `configuration.authentication.username_password.password.secret_type` is not populated and a diff on `configuration.authentication.username_password.username` after import in `google_bigquery_connection` resource ([#12179](https://github.com/hashicorp/terraform-provider-google-beta/pull/12179))
* bigqueryreservation: Fixed `google_bigquery_reservation_assignment` returning a confusing 404 error when `reservation` is a bare name and `location` is not set ([#12210](https://github.com/hashicorp/terraform-provider-google-beta/pull/12210))
* ces: updated supported values for `channel_type`, `modality`, and `theme` in `google_ces_deployment` ([#12227](https://github.com/hashicorp/terraform-provider-google-beta/pull/12227))
* compute: updated `google_compute_forwarding_rule` resource to properly prompt for resource recreation when updating the `target` field between different "serviceAttachments", rather than having an in-place update blocked by an API error. ([#12214](https://github.com/hashicorp/terraform-provider-google-beta/pull/12214))
* modelarmor: fixed permadiff and `REQUEST_FIELD_MISSING` error when `template_metadata` is omitted from `google_model_armor_template` ([#12222](https://github.com/hashicorp/terraform-provider-google-beta/pull/12222))
* networkconnectivity: fixed an issue where `google_network_connectivity_destination` was not recognizing the `name` field as mapping to an API value ([#12224](https://github.com/hashicorp/terraform-provider-google-beta/pull/12224))
* networkconnectivity: fixed an issue where `google_network_connectivity_multicloud_data_transfer_config` was not recognizing the `name` field as mapping to an API value ([#12224](https://github.com/hashicorp/terraform-provider-google-beta/pull/12224))
* resourcemanager: added verification polling to `google_service_account` updates to ensure the resource is consistent before succeeding ([#12217](https://github.com/hashicorp/terraform-provider-google-beta/pull/12217))

## 7.32.0 (May 12, 2026)

NOTES:
* compute: migrated `google_compute_instance_from_machine` resource  to use direct HTTP rather than a client library ([#27260](https://github.com/hashicorp/terraform-provider-google/pull/27260))
* compute: migrated `google_compute_instance_group_manager` resource to use direct HTTP rather than a client library ([#27259](https://github.com/hashicorp/terraform-provider-google/pull/27259))
* compute: migrated `google_compute_zones` data source to use direct HTTP rather than a client library ([#27261](https://github.com/hashicorp/terraform-provider-google/pull/27261))
* compute: migrated `google_compute_project_metadata_item` resource to use direct HTTP rather than a client library ([#27200](https://github.com/hashicorp/terraform-provider-google/pull/27200))

FEATURES:
* **New Data Source:** `google_compute_region_instant_snapshot_iam_policy` ([#27281](https://github.com/hashicorp/terraform-provider-google/pull/27281))
* **New Resource:** `google_chronicle_dashboard_chart` ([#27275](https://github.com/hashicorp/terraform-provider-google/pull/27275))
* **New Resource:** `google_compute_region_instant_snapshot_iam_binding` ([#27281](https://github.com/hashicorp/terraform-provider-google/pull/27281))
* **New Resource:** `google_compute_region_instant_snapshot_iam_member` ([#27281](https://github.com/hashicorp/terraform-provider-google/pull/27281))
* **New Resource:** `google_compute_region_instant_snapshot_iam_policy` ([#27281](https://github.com/hashicorp/terraform-provider-google/pull/27281))
* **New Resource:** `google_compute_region_instant_snapshot` ([#27281](https://github.com/hashicorp/terraform-provider-google/pull/27281))

IMPROVEMENTS:
* compute: added `IDPF` value to `nic_type` in `resource_compute_instance_template` ([#27244](https://github.com/hashicorp/terraform-provider-google/pull/27244))
* compute: added `IDPF` value to `nic_type` in `resource_compute_instance` ([#27244](https://github.com/hashicorp/terraform-provider-google/pull/27244))
* compute: added `IDPF` value to `nic_type` in `resource_compute_region_instance_template` ([#27244](https://github.com/hashicorp/terraform-provider-google/pull/27244))
* compute: added `address_id` field to `google_compute_address` resource ([#27216](https://github.com/hashicorp/terraform-provider-google/pull/27216))
* compute: added `advanced_options_config` field on `google_compute_organization_security_policy` resource ([#27255](https://github.com/hashicorp/terraform-provider-google/pull/27255))
* compute: added `connection_tracking_policy` field to `google_compute_region_backend_service` resource ([#27217](https://github.com/hashicorp/terraform-provider-google/pull/27217))
* compute: added `image`, `source_image_encryption_key`, and `source_image_id` fields to `google_compute_region_disk` resource. This field is currently behind an allowlist. ([#27243](https://github.com/hashicorp/terraform-provider-google/pull/27243))
* compute: added `replica_zones` field to `google_compute_instance` resource ([#27258](https://github.com/hashicorp/terraform-provider-google/pull/27258))
* compute: added `request_body` field on `google_compute_security_policy_rule` resource ([#27252](https://github.com/hashicorp/terraform-provider-google/pull/27252))
* compute: added update support for `ip_collection` field to `google_compute_subnetwork` resource ([#27265](https://github.com/hashicorp/terraform-provider-google/pull/27265))
* discoveryengine: added `config_id` attribute to `google_discovery_engine_widget_config` ([#27278](https://github.com/hashicorp/terraform-provider-google/pull/27278))
* networksecurity: added support for project `parent` values to `google_network_security_firewall_endpoint` ([#27222](https://github.com/hashicorp/terraform-provider-google/pull/27222))
* recaptchaenterprise: added `POLICY_BASED_CHALLENGE` value to `integration_type` field and added new `challenge_settings` field to `google_recaptcha_enterprise_key` ([#27221](https://github.com/hashicorp/terraform-provider-google/pull/27221))
* redis: added new node types supported in `google_redis_cluster`. ([#27242](https://github.com/hashicorp/terraform-provider-google/pull/27242))
* resourcemanager: add `private_key` and `private_key_type` fields to ephemeral `google_service_account_key` resource ([#27279](https://github.com/hashicorp/terraform-provider-google/pull/27279))
* storage: added `ingest_on_write` field for `google_storage_anywhere_cache` resource ([#27271](https://github.com/hashicorp/terraform-provider-google/pull/27271))
* workstations: added `gce_hd` field to `google_workstations_workstation_config` resource ([#27201](https://github.com/hashicorp/terraform-provider-google/pull/27201))

BUG FIXES:
* cloudfunctions2: fixed bug where `all_traffic_on_latest_revision = false` was ignored in `google_cloudfunctions2_function` ([#27256](https://github.com/hashicorp/terraform-provider-google/pull/27256))
* compute: fixed permadiff when removing `preconfigured_waf_config` from a `google_compute_security_policy` rule ([#27276](https://github.com/hashicorp/terraform-provider-google/pull/27276))

## 7.31.0 (May 5, 2026)

NOTES:
* compute: migrated `google_compute_instance.network_interface` field to use direct HTTP rather than a client library ([#27104](https://github.com/hashicorp/terraform-provider-google/pull/27104))
* compute: migrated `google_compute_image` datasource to use direct HTTP rather then a client library ([#27179](https://github.com/hashicorp/terraform-provider-google/pull/27179))
* compute: migrated `partner_metadata` field on `google_compute_instance`, `google_compute_instance_template`, and `google_compute_region_instance_template` to use direct HTTP rather than a client library ([#27131](https://github.com/hashicorp/terraform-provider-google/pull/27131))
* compute: migrated `google_compute_node_types` data source to use direct HTTP rather than a client library ([#27184](https://github.com/hashicorp/terraform-provider-google/pull/27184))
* compute: migrated `google_compute_region_instance_group` data source to use direct HTTP rather than a client library ([#27178](https://github.com/hashicorp/terraform-provider-google/pull/27178))
* compute: migrated `google_compute_subnetwork` data source to use direct HTTP rather than a client library ([#27167](https://github.com/hashicorp/terraform-provider-google/pull/27167))
* compute: migrated `google_compute_vpn_gateway` data source to use direct HTTP rather than a client library ([#27168](https://github.com/hashicorp/terraform-provider-google/pull/27168))

FEATURES:
* **New Data Source:** `google_artifact_registry_file` ([#27183](https://github.com/hashicorp/terraform-provider-google/pull/27183))
* **New Resource:** `google_ces_app_root_agent_association` ([#27123](https://github.com/hashicorp/terraform-provider-google/pull/27123))
* **New Resource:** `google_contact_center_insights_qa_question` ([#27169](https://github.com/hashicorp/terraform-provider-google/pull/27169))
* **New Resource:** `google_contact_center_insights_qa_scorecard_revision` ([#27169](https://github.com/hashicorp/terraform-provider-google/pull/27169))
* **New Resource:** `google_contact_center_insights_qa_scorecard` ([#27169](https://github.com/hashicorp/terraform-provider-google/pull/27169))
* **New Resource:** `google_firebase_app_check_resource_policy` ([#27185](https://github.com/hashicorp/terraform-provider-google/pull/27185))

IMPROVEMENTS:
* clouddeploy: added `default_pool` and `private_pool` fields to `google_clouddeploy_target` resource ([#27187](https://github.com/hashicorp/terraform-provider-google/pull/27187))
* clouddeploy: added `tasks` and `analysis` fields to `google_clouddeploy_delivery_pipeline` resource ([#27187](https://github.com/hashicorp/terraform-provider-google/pull/27187))
* compute: added `params.resource_manager_tags` field to `google_compute_image` ([#27107](https://github.com/hashicorp/terraform-provider-google/pull/27107))
* compute: added `params.resource_manager_tags` field to `google_compute_region_commitment` resource ([#27181](https://github.com/hashicorp/terraform-provider-google/pull/27181))
* compute: added `resource_policies.workload_policy` to `google_compute_region_instance_group_manager` resource ([#27170](https://github.com/hashicorp/terraform-provider-google/pull/27170))
* compute: marked csek disk encryption key fields as sensitive in compute resources ([#27193](https://github.com/hashicorp/terraform-provider-google/pull/27193))
* container: added `node_pool.network_config.accelerator_network_profile` to `google_container_cluster` resource and `network_config.accelerator_network_profile` to `google_container_node_pool` resource ([#27171](https://github.com/hashicorp/terraform-provider-google/pull/27171))
* databasemigrationservice: added `objects_config` field to `google_database_migration_service_migration_job` resource ([#27180](https://github.com/hashicorp/terraform-provider-google/pull/27180))
* dataplex: added `attributes`, `template_reference`, `enable_catalog_basedRules`, and `filter` fields to `google_dataplex_datascan` resource ([#27130](https://github.com/hashicorp/terraform-provider-google/pull/27130))
* firestore: added `search_config` field to `google_firestore_index` resource ([#27108](https://github.com/hashicorp/terraform-provider-google/pull/27108))
* oracle_database: added `pluggable_database_id`, `pluggable_database_name` fields to `google_oracle_database_db_system` resource ([#27127](https://github.com/hashicorp/terraform-provider-google/pull/27127))

BUG FIXES:
* provider: fixed a bad `timeouts` diff across a number of resources that had resource identity support added in `7.29.0` ([#27189](https://github.com/hashicorp/terraform-provider-google/pull/27189))
* assuredworkloads: made assuredworkloads resources use GA endpoint instead of beta ([#27122](https://github.com/hashicorp/terraform-provider-google/pull/27122))
* bigquery: fixed `ignore_auto_generated_schema` evaluation for `google_bigquery_table` external tables which caused spurious replacement ([#27188](https://github.com/hashicorp/terraform-provider-google/pull/27188))
* cloudscheduler: fixed perpetual diff on `google_cloud_scheduler_job.http_target.headers` when `oidc_token` or `oauth_token` is set ([#27173](https://github.com/hashicorp/terraform-provider-google/pull/27173))
* servicenetworking: fixed a permadiff issue of `reserved_peering_ranges` in `google_service_networking_connection` ([#27132](https://github.com/hashicorp/terraform-provider-google/pull/27132))
* storage: fix inconsistent plan issue for `google_storage_notification.custom_attributes` field ([#27129](https://github.com/hashicorp/terraform-provider-google/pull/27129))

## 7.30.0 (Apr 28, 2026)

BREAKING CHANGES:
* apigee: fixed `google_apigee_env_keystore` to require the `name` field which is mandatory in the Apigee API ([#27006](https://github.com/hashicorp/terraform-provider-google/pull/27006))

FEATURES:
* **New Data Source:** `google_data_lineage_config` ([#27098](https://github.com/hashicorp/terraform-provider-google/pull/27098))
* **New Resource:** `google_artifact_registry_rule` ([#27049](https://github.com/hashicorp/terraform-provider-google/pull/27049))
* **New Resource:** `google_data_lineage_config` ([#27098](https://github.com/hashicorp/terraform-provider-google/pull/27098))
* **New Resource:** `google_document_ai_schema` ([#27102](https://github.com/hashicorp/terraform-provider-google/pull/27102))
* **New Resource:** `google_firebase_remote_config_remote_config` ([#27050](https://github.com/hashicorp/terraform-provider-google/pull/27050))

IMPROVEMENTS:
* provider: added support for `prefer_global_endpoints` and `prefer_regional_endpoints` to the provider configuration. Support for regional endpoints will be rolled out on a per-product level ([#27014](https://github.com/hashicorp/terraform-provider-google/pull/27014))
* artifactregistry: added support for regionalized endpoints ([#27014](https://github.com/hashicorp/terraform-provider-google/pull/27014))
* assuredworkloads: added `SPAIN_DATA_BOUNDARY_BY_TELEFONICA` value to `partner` field on `google_assured_workloads_workload` resource ([#27027](https://github.com/hashicorp/terraform-provider-google/pull/27027))
* bigqueryconnection: added `configuration` block to `google_bigquery_connection` resource to support AlloyDB and other connector types via the BigQuery Connector framework ([#27029](https://github.com/hashicorp/terraform-provider-google/pull/27029))
* bigtable: added support for `tags` to `google_bigtable_instance` ([#27060](https://github.com/hashicorp/terraform-provider-google/pull/27060))
* cloudrunv2: added `DISK` fields to `google_cloud_run_v2_job` resource ([#27052](https://github.com/hashicorp/terraform-provider-google/pull/27052))
* cloudrunv2: added `DISK` fields to `google_cloud_run_v2_worker_pool` resource ([#27048](https://github.com/hashicorp/terraform-provider-google/pull/27048))
* compute: add `params.resourceManagerTags` field to the `google_compute_storage_pool` ([#27051](https://github.com/hashicorp/terraform-provider-google/pull/27051))
* compute: added `cache_policy` field to `google_compute_url_map` ([#27011](https://github.com/hashicorp/terraform-provider-google/pull/27011))
* compute: added `params.resource_manager_tags` field to `google_compute_instant_snapshot` resource ([#27087](https://github.com/hashicorp/terraform-provider-google/pull/27087))
* compute: added `resource_manager_tags` field to `google_compute_machine_image` resource ([#27075](https://github.com/hashicorp/terraform-provider-google/pull/27075))
* container: added `node_config.linux_node_config.accurate_time_config` field to `google_container_node_pool` resource ([#27064](https://github.com/hashicorp/terraform-provider-google/pull/27064))
* container: added `node_pool.node_config.linux_node_config.accurate_time_config` and `node_config.linux_node_config.accurate_time_config` fields to `google_container_cluster` resource ([#27064](https://github.com/hashicorp/terraform-provider-google/pull/27064))
* container: added `node_pool.node_config.linux_node_config.swap_config` field to `google_container_node_pool` resource ([#26982](https://github.com/hashicorp/terraform-provider-google/pull/26982))
* container: increased default timeout for `google_container_cluster` to 90 minutes (from 40/60 depending on operation) and `google_container_node_pool` to 60 minutes (from 30) ([#27101](https://github.com/hashicorp/terraform-provider-google/pull/27101))
* discoveryengine: added `destionation_configs.destionations.port` and `destionation_configs.params` fields to `google_discovery_engine_data_connector ` resource ([#27058](https://github.com/hashicorp/terraform-provider-google/pull/27058))
* dns: added support for IAM conditions to `google_dns_managed_zone` resource ([#27010](https://github.com/hashicorp/terraform-provider-google/pull/27010))
* datastream: added `deletion_policy` field to control whether child routes are force-deleted to `google_datastream_private_connection` ([#27033](https://github.com/hashicorp/terraform-provider-google/pull/27033))
* networkconnectivity: added support for IAM conditions to `google_network_connectivity_hub` resource ([#27005](https://github.com/hashicorp/terraform-provider-google/pull/27005))
* networksecurity: added `parent` field to `google_network_security_address_groups` data source ([#27082](https://github.com/hashicorp/terraform-provider-google/pull/27082))
* workbench: added support for new disk types and accelerators to `google_workbench_instance` ([#27061](https://github.com/hashicorp/terraform-provider-google/pull/27061))

BUG FIXES:
* alloydb: fixed `google_alloydb_cluster` so that `maintenance_update_policy.maintenance_windows.start_time.hours` can be set to `0` (midnight) ([#26981](https://github.com/hashicorp/terraform-provider-google/pull/26981))
* ces: fixed type mismatch in `google_ces_app` variable default value ([#27084](https://github.com/hashicorp/terraform-provider-google/pull/27084))
* compute: fixed an issue where an erroneous error could occur for having an unset `zone` field in `google_compute_instance_template` ([#27076](https://github.com/hashicorp/terraform-provider-google/pull/27076))
* compute: fixed permadiff for `iap.oauth2_client_id` in `google_compute_backend_service` and `google_compute_region_backend_service` when the API returns a single space ([#26975](https://github.com/hashicorp/terraform-provider-google/pull/26975))
* container: fixed a permadiff in `google_container_cluster` where `database_encryption.state` returning `ALL_OBJECTS_ENCRYPTION_ENABLED` instead of the configured `ENCRYPTED` caused unintended reapplies ([#27040](https://github.com/hashicorp/terraform-provider-google/pull/27040))
* dataplex: fixed acceptance test failure for one time scans ([#27095](https://github.com/hashicorp/terraform-provider-google/pull/27095))
* dialogflowcx: fixed a perma-diff in `google_dialogflow_cx_test_case` when `session_parameters` was omitted from the configuration ([#26985](https://github.com/hashicorp/terraform-provider-google/pull/26985))
* hypercomputecluster: fixed a permadiff in `google_hypercomputecluster_cluster` when `count`, `static_node_count`, or `max_dynamic_node_count` were explicitly set to `0`. ([#27073](https://github.com/hashicorp/terraform-provider-google/pull/27073))
* identityplatform: fixed a premadiff on `multi_tenant` in `google_identity_platform_config` resource. Removing the value from config will now preserve the existing settings instead of removing them. ([#26986](https://github.com/hashicorp/terraform-provider-google/pull/26986))
* memorystore: fixed an issue preventing updating multiple properties at once for `google_redis_cluster` ([#27077](https://github.com/hashicorp/terraform-provider-google/pull/27077))

NOTES:
* compute: Migrate `resource_compute_instance_group.go.tmpl` resource to use direct HTTP rather then a client library ([#27080](https://github.com/hashicorp/terraform-provider-google/pull/27080))
* compute: migrated `compute-operation` resource to use direct HTTP rather then a client library ([#27053](https://github.com/hashicorp/terraform-provider-google/pull/27053))
* compute: migrated `compute_backend_bucket_security_policy` resource to use direct HTTP rather than a client library ([#27012](https://github.com/hashicorp/terraform-provider-google/pull/27012))
* compute: migrated `compute_instance_network_interface_helpers` resource to use direct HTTP rather than a client library ([#27104](https://github.com/hashicorp/terraform-provider-google/pull/27104))
* compute: migrated `data_source_google_compute_addresses.go.tmpl` data source to use direct HTTP rather then a client library ([#27016](https://github.com/hashicorp/terraform-provider-google/pull/27016))
* compute: migrated `data_source_google_compute_machine_types` datasource to use direct HTTP rather than a client library ([#27017](https://github.com/hashicorp/terraform-provider-google/pull/27017))
* compute: migrated `google_disk_test` to use direct HTTP rather than a client library ([#27079](https://github.com/hashicorp/terraform-provider-google/pull/27079))
* compute: migrated `resource_compute_disk_async_replication` resource to use direct HTTP rather then a client library ([#27028](https://github.com/hashicorp/terraform-provider-google/pull/27028))
* compute: migrated `resource_compute_http_health_check_test.go.tmpl` resource to use direct HTTP rather then a client library ([#27057](https://github.com/hashicorp/terraform-provider-google/pull/27057))

## 7.29.0 (Apr 21, 2026)

NOTES:
* provider: List resources are now supported in both google and google-beta providers with the introduction of `google_service_account` list resource - more info can be found [here](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query) ([#26938](https://github.com/hashicorp/terraform-provider-google/pull/26938))

FEATURES:
* **New Data Source:** `google_firebase_admin_sdk_config` ([#26901](https://github.com/hashicorp/terraform-provider-google/pull/26901))
* **New Resource:** `google_chronicle_datatable_row` ([#26960](https://github.com/hashicorp/terraform-provider-google/pull/26960))
* **New Resource:** `google_chronicle_datatable` ([#26895](https://github.com/hashicorp/terraform-provider-google/pull/26895))
* **New Resource:** `google_dataform_folder` ([#26881](https://github.com/hashicorp/terraform-provider-google/pull/26881))
* **New Resource:** `google_dataform_team_folder` ([#26881](https://github.com/hashicorp/terraform-provider-google/pull/26881))
* **New Resource:** `google_firebase_storage_default_bucket` ([#26965](https://github.com/hashicorp/terraform-provider-google/pull/26965))

IMPROVEMENTS:
* alloydb: added `track_client_address` field to `google_alloydb_instance` resource ([#26964](https://github.com/hashicorp/terraform-provider-google/pull/26964))
* clouddeploy: added `tasks` field to `google_clouddeploy_custom_target_type` resource ([#26941](https://github.com/hashicorp/terraform-provider-google/pull/26941))
* compute: added `header_action` and `redirect_options` fields to `google_compute_organization_security_policy_rule` resource ([#26942](https://github.com/hashicorp/terraform-provider-google/pull/26942))
* dataplex: added `execution_identity` field to `google_dataplex_datascan` resource ([#26924](https://github.com/hashicorp/terraform-provider-google/pull/26924))
* dataproc: added `cluster_config.engine` field to `google_dataproc_cluster` resource ([#26962](https://github.com/hashicorp/terraform-provider-google/pull/26962))
* iambeta: added `trust_default_shared_ca` field to `google_iam_workload_identity_pool` resource ([#26974](https://github.com/hashicorp/terraform-provider-google/pull/26974))
* netapp: added `large_capacity_config` field to `google_netapp_volume` resource([#26927](https://github.com/hashicorp/terraform-provider-google/pull/26927))
* netapp: added `kms_config`, `encryption_state` and `backups_crypto_key_version` fields to `google_netapp_backup_vault` resource ([#26939](https://github.com/hashicorp/terraform-provider-google/pull/26939))
* resourcemanager: add resource-identity support to `google_service_account` resource ([#26938](https://github.com/hashicorp/terraform-provider-google/pull/26938))
* sql: added `entraid_config` field to `google_sql_database_instance` resource ([#26921](https://github.com/hashicorp/terraform-provider-google/pull/26921))
* vectorsearch: added `encryption_spec` field to `google_vector_search_collection` resource ([#26972](https://github.com/hashicorp/terraform-provider-google/pull/26972))

BUG FIXES:
* apigee: fixed ignoring `is_enabled = false` on create and update in `google_apigee_target_server` resource ([#26878](https://github.com/hashicorp/terraform-provider-google/pull/26878))
* bigquery: fixed inability to set `default_collation` to empty string in `google_bigquery_dataset` ([#26925](https://github.com/hashicorp/terraform-provider-google/pull/26925))
* ces: fixed a diff on `logging_settings` when unspecified in `google_ces_app`. Removing the value from config will now preserve the existing settings instead of removing them. ([#26899](https://github.com/hashicorp/terraform-provider-google/pull/26899))
* compute: fixed a permadiff on `iap.oauth2_client_id` in `google_compute_backend_service` and `google_compute_region_backend_service` when the API returns a single space ([#26975](https://github.com/hashicorp/terraform-provider-google/pull/26975))
* container: fixed a bug in `google_container_cluster` where setting multiple fields in `dns_endpoint_config` failed to apply all changes ([#26968](https://github.com/hashicorp/terraform-provider-google/pull/26968))
* workstations: fixed a permadiff on `persistent_directories.gce_pd.reclaim_policy` in `google_workstations_workstation_config` resource ([#26971](https://github.com/hashicorp/terraform-provider-google/pull/26971))

## 7.28.0 (Apr 14, 2026)

NOTES:
* compute: migrated `data_source_google_compute_instance_template` datasource to use direct HTTP rather then a client library ([#26831](https://github.com/hashicorp/terraform-provider-google/pull/26831))
* compute: migrated `google_compute_instance_guest_attributes` datasource to use direct HTTP rather then a client library ([#26826](https://github.com/hashicorp/terraform-provider-google/pull/26826))
* provider: added provider-wide `Identity()` schema support, allowing imports with MMv1 resources to occur using the identity block instead of id field ([#26783](https://github.com/hashicorp/terraform-provider-google/pull/26783))

FEATURES:
* **New Data Source:** `google_vertex_ai_reasoning_engine_query` ([#26787](https://github.com/hashicorp/terraform-provider-google/pull/26787))
* **New Resource:** `google_apigee_space` ([#26857](https://github.com/hashicorp/terraform-provider-google/pull/26857))
* **New Resource:** `google_vertex_ai_reasoning_engine_iam_binding` ([#26785](https://github.com/hashicorp/terraform-provider-google/pull/26785))
* **New Resource:** `google_vertex_ai_reasoning_engine_iam_member` ([#26785](https://github.com/hashicorp/terraform-provider-google/pull/26785))
* **New Resource:** `google_vertex_ai_reasoning_engine_iam_policy` ([#26785](https://github.com/hashicorp/terraform-provider-google/pull/26785))
* **New Resource:** `google_workload_identity_service_agent` ([#26780](https://github.com/hashicorp/terraform-provider-google/pull/26780))

IMPROVEMENTS:
* bigqueryanalyticshub: added `replica_locations` and `effective_replicas` fields to `google_bigquery_analytics_hub_listing` resource ([#26843](https://github.com/hashicorp/terraform-provider-google/pull/26843))
* bigqueryanalyticshub: added `replica_locations` field to `google_bigquery_analytics_hub_listing_subscription` resource ([#26843](https://github.com/hashicorp/terraform-provider-google/pull/26843))
* composer: increased `google_composer_environment` default delete timeout to 120m from 30m ([#26851](https://github.com/hashicorp/terraform-provider-google/pull/26851))
* compute: added `target_size_policy` field to `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` resources ([#26849](https://github.com/hashicorp/terraform-provider-google/pull/26849))
* compute: increased `google_compute_security_policy` default timeout to 60m from 30m ([#26850](https://github.com/hashicorp/terraform-provider-google/pull/26850))
* compute: supported simultaneous updates for Hyperdisk IOPS and throughput in `google_compute_disk` and `google_compute_region_disk` resources ([#26815](https://github.com/hashicorp/terraform-provider-google/pull/26815))
* container: added `autopilot_cluster_policy_config` field to `google_container_cluster` resource ([#26822](https://github.com/hashicorp/terraform-provider-google/pull/26822))
* container: added `disable_multi_nic` field to `lustre_csi_driver_config` in `google_container_cluster` resource ([#26759](https://github.com/hashicorp/terraform-provider-google/pull/26759))
* developerconnect: added `custom_oauth_config`, `etag`, and `proxy_config` fields to `google_developer_connect_account_connector` resource ([#26751](https://github.com/hashicorp/terraform-provider-google/pull/26751))
* netapp: added `scale_type` field to `google_netapp_storage_pool` resource ([#26821](https://github.com/hashicorp/terraform-provider-google/pull/26821))
* netapp: added `mode` field to `google_netapp_storage_pool` resource ([#26778](https://github.com/hashicorp/terraform-provider-google/pull/26778))
* networkservices: added `all_ports` field to `google_network_services_gateway` resource ([#26808](https://github.com/hashicorp/terraform-provider-google/pull/26808))
* sql: added `SQLSERVER_2025` value to `database_version` field in `database_instance` resource ([#26845](https://github.com/hashicorp/terraform-provider-google/pull/26845))
* vertexai: add `labels` field to `google_vertex_ai_reasoning_engine` resource ([#26825](https://github.com/hashicorp/terraform-provider-google/pull/26825))
* vertexai: added `spec.source_code_spec.image_spec` field to `google_vertex_ai_reasoning_engine` resource ([#26790](https://github.com/hashicorp/terraform-provider-google/pull/26790))
* vertexai: added `container_spec` field to `google_vertex_ai_reasoning_engine` resource ([#26813](https://github.com/hashicorp/terraform-provider-google/pull/26813))
* vertexai: added `spec.identity_type` and `spec.effective_identity` fields to `google_vertex_ai_reasoning_engine` resource ([#26788](https://github.com/hashicorp/terraform-provider-google/pull/26788))

BUG FIXES:
* apigee: fixed a crash in `google_apigee_environment_addons_config` resource when analytics are not configured  ([#26810](https://github.com/hashicorp/terraform-provider-google/pull/26810))
* apigee: fixed overly restrictive validation of `name` field in `google_apigee_api_product` that rejected uppercase letters, aligning provider behavior with the Apigee API ([#26756](https://github.com/hashicorp/terraform-provider-google/pull/26756))
* bigquery: fixed crash when `hive_partitioning_options` is defined with all null values in `google_bigquery_table` resource ([#26846](https://github.com/hashicorp/terraform-provider-google/pull/26846))
* firebaseailogic: fixed permadiff on `traffic_filter` field in `google_firebase_ai_logic_config` resource ([#26749](https://github.com/hashicorp/terraform-provider-google/pull/26749))
* networksecurity: fixed permadiff on `policy_profile` field in `google_network_security_authz_policy` resource ([#26865](https://github.com/hashicorp/terraform-provider-google/pull/26865))
* vertexai: added 10-second wait before reading the updated resource in `google_vertex_ai_reasoning_engine`, preventing stale values getting written to state ([#26852](https://github.com/hashicorp/terraform-provider-google/pull/26852))

## 7.27.0 (Apr 07, 2026)

BREAKING CHANGES:
* lustre: marked `maintenance_policy.weekly_maintenance_windows` field required in `google_lustre_instance` resource. Configuring `maintenance_policy` without `weekly_maintenance_windows` will cause an API error. ([#26741](https://github.com/hashicorp/terraform-provider-google/pull/26741))

FEATURES:
* **New Data Source:** `google_discovery_engine_data_store` ([#26651](https://github.com/hashicorp/terraform-provider-google/pull/26651))
* **New Data Source:** `google_discovery_engine_data_stores` ([#26651](https://github.com/hashicorp/terraform-provider-google/pull/26651))
* **New Data Source:** `google_dns_record_sets` ([#26736](https://github.com/hashicorp/terraform-provider-google/pull/26736))
* **New Resource:** `google_chronicle_dashboard_chart` ([#26707](https://github.com/hashicorp/terraform-provider-google/pull/26707))
* **New Resource:** `google_chronicle_feed` ([#26742](https://github.com/hashicorp/terraform-provider-google/pull/26742))
* **New Resource:** `google_network_connectivity_transport` ([#26626](https://github.com/hashicorp/terraform-provider-google/pull/26626))
* **New Resource:** `google_iam_workload_identity_pool_managed_identity` ([#26732](https://github.com/hashicorp/terraform-provider-google/pull/26732))
* **New Resource:** `google_iam_workload_identity_pool_namespace` ([#26647](https://github.com/hashicorp/terraform-provider-google/pull/26647))

IMPROVEMENTS:
* compute: added `SEV_LIVE_MIGRATABLE_V2` to `guest_os_features` enum for `google_compute_region_disk` resource ([#26735](https://github.com/hashicorp/terraform-provider-google/pull/26735))
* compute: added `SNP_SVSM_CAPABLE` to `guest_os_features` enum for `google_compute_image` and `google_compute_region_disk` resources ([#26735](https://github.com/hashicorp/terraform-provider-google/pull/26735))
* compute: added `excluded_folders` and `excluded_projects` fields to `google_compute_organization_security_policy_association` resource ([#26694](https://github.com/hashicorp/terraform-provider-google/pull/26694))
* compute: supported in-place update for `secondary_ip_range` field in `google_compute_subnetwork` resource ([#26689](https://github.com/hashicorp/terraform-provider-google/pull/26689))
* container: added `autopilot_privileged_admission` field to `google_container_cluster` resource for Customer-Driven Allowlisting ([#26668](https://github.com/hashicorp/terraform-provider-google/pull/26668))
* dataplex: added `aspects` field to `google_dataplex_entry_link` resource ([#26664](https://github.com/hashicorp/terraform-provider-google/pull/26664))
* dataplex: supported in-place update for `aspects` field in `google_dataplex_entry_link` resource ([#26702](https://github.com/hashicorp/terraform-provider-google/pull/26702))
* dataproc: added `boot_disk_provisioned_iops` and `boot_disk_provisioned_throughput` fields to `cluster_config.worker_config.disk_config` in `google_dataproc_cluster` resource ([#26691](https://github.com/hashicorp/terraform-provider-google/pull/26691))
* dataproc: added value `AUTO` to `runtime_config.autotuning_config.scenarios` field in `google_dataproc_batch` resource ([#26646](https://github.com/hashicorp/terraform-provider-google/pull/26646))
* iambeta: added `attestation_rules` field to `google_iam_workload_identity_pool` resource ([#26706](https://github.com/hashicorp/terraform-provider-google/pull/26706))
* lustre: added `dynamic_tier_options` field to `google_lustre_instance` resource ([#26741](https://github.com/hashicorp/terraform-provider-google/pull/26741))
* migrationcenter: added `virtual_machine_preferences.compute_engine_preferences.persistent_disk_type` field to `google_migration_center_preference_set` resource ([#26693](https://github.com/hashicorp/terraform-provider-google/pull/26693))
* networkconnectivity: added `exclude_import_ranges`, `include_export_ranges`, `exclude_export_ranges` fields to `google_network_connectivity_spoke` resource ([#26730](https://github.com/hashicorp/terraform-provider-google/pull/26730))
* pubsub: added `ai_inference` field to `google_pubsub_topic` and `google_pubsub_subscription` resources ([#26738](https://github.com/hashicorp/terraform-provider-google/pull/26738))
* sql: added `clone_context.source_project` field to `google_sql_database_instance` resource to support cross project clone ([#26652](https://github.com/hashicorp/terraform-provider-google/pull/26652))

BUG FIXES:
* compute: fixed a permadiff on the `adaptive_protection_config` field in `google_compute_security_policy` resource ([#26692](https://github.com/hashicorp/terraform-provider-google/pull/26692))
* compute: fixed panic when setting `google_compute_project_metadata` on a project with no existing metadata ([#26630](https://github.com/hashicorp/terraform-provider-google/pull/26630))
* biglakeiceberg: changed the `primary-location` parameter to `primary_location` in the create URL of google_biglake_iceberg_catalog resource ([#26695](https://github.com/hashicorp/terraform-provider-google/pull/26695))
* securityposture: always sent value of `enforce` in `policies.constraint.org_policy_constraint.policy_rules` to the api in `google_securityposture_posture` resource ([#26645](https://github.com/hashicorp/terraform-provider-google/pull/26645))
* vertexai: fixed missing Private Service Connect service attachment for `service_attachment` field in `google_vertex_ai_endpoint_with_model_garden_deployment` resource ([#26690](https://github.com/hashicorp/terraform-provider-google/pull/26690))
* workstations: fixed update of `private_cluster_config.allowed_projects` in `google_workstations_workstation_cluster` resource ([#26705](https://github.com/hashicorp/terraform-provider-google/pull/26705))

## 7.26.0 (Mar 31, 2026)

BREAKING CHANGES:
* compute: Removed `google_compute_region_backend_bucket` from the `google` (GA) provider. It is currently beta-only, and calls to the nonexistent GA API always returned a 404. Until released in `google`, use `google-beta` instead. ([#26597](https://github.com/hashicorp/terraform-provider-google/pull/26597))

FEATURES:
* **New Data Source:** `google_network_security_address_groups` ([#26562](https://github.com/hashicorp/terraform-provider-google/pull/26562))
* **New Data Source:** `google_iam_workload_identity_pool_iam_policy` ([#26598](https://github.com/hashicorp/terraform-provider-google/pull/26598))
* **New Resource:** `google_bigqueryreservation_reservation_group` ([#26560](https://github.com/hashicorp/terraform-provider-google/pull/26560))
* **New Resource:** `google_compute_region_composite_health_check` ([#26591](https://github.com/hashicorp/terraform-provider-google/pull/26591))
* **New Resource:** `google_compute_region_health_aggregation_policy` ([#26591](https://github.com/hashicorp/terraform-provider-google/pull/26591))
* **New Resource:** `google_compute_region_health_source` ([#26591](https://github.com/hashicorp/terraform-provider-google/pull/26591))
* **New Resource:** `google_contact_center_insights_assessment_rule` ([#26530](https://github.com/hashicorp/terraform-provider-google/pull/26530))
* **New Resource:** `google_iam_workload_identity_pool_iam_*` ([#26598](https://github.com/hashicorp/terraform-provider-google/pull/26598))
* **New Resource:** `google_workstations_workstation` ([#26561](https://github.com/hashicorp/terraform-provider-google/pull/26561))
* **New Resource:** `google_workstations_workstation_iam_*` ([#26561](https://github.com/hashicorp/terraform-provider-google/pull/26561))
* **New Resource:** `google_workstations_workstation_cluster` ([#26561](https://github.com/hashicorp/terraform-provider-google/pull/26561))
* **New Resource:** `google_workstations_workstation_config` ([#26561](https://github.com/hashicorp/terraform-provider-google/pull/26561))
* **New Resource:** `google_workstations_workstation_config_iam_*` ([#26561](https://github.com/hashicorp/terraform-provider-google/pull/26561))

IMPROVEMENTS:
* bigqueryreservation: added `reservation_group` field to `google_bigquery_reservation` resource ([#26560](https://github.com/hashicorp/terraform-provider-google/pull/26560))
* ces: added `remote_dialogflow_agent.respect_response_interruption_settings` field to `google_ces_agent` resource ([#26578](https://github.com/hashicorp/terraform-provider-google/pull/26578))
* clusterdirector: made `boot_disk.size_gb` and `boot_disk.type` editable within nodesets and login nodes in `google_hypercomputecluster_cluster` ([#26615](https://github.com/hashicorp/terraform-provider-google/pull/26615))
* colab: added `colab_image` field to `google_colab_runtime_template` resource ([#26582](https://github.com/hashicorp/terraform-provider-google/pull/26582))
* colab: made `google_colab_runtime_template` resource updatable ([#26582](https://github.com/hashicorp/terraform-provider-google/pull/26582))
* compute: added `hyperdisk-balanced` as an option for `disk_type` field in `google_container_cluster` resource ([#26581](https://github.com/hashicorp/terraform-provider-google/pull/26581))
* compute: made `backend_service` field optional for `google_compute_target_tcp_proxy` resource ([#26519](https://github.com/hashicorp/terraform-provider-google/pull/26519))
* compute: promoted `resolve_subnet_field` field in `google_compute_subnetwork` resource to GA ([#26570](https://github.com/hashicorp/terraform-provider-google/pull/26570))
* iambeta: promoted `mode`, `inline_certificate_issuance_config`, and `inline_trust_config` fields in `google_iam_workload_identity_pool` resource to GA ([#26598](https://github.com/hashicorp/terraform-provider-google/pull/26598))
* spanner: added autoscaling config for instance partition and missing asymmetric autoscaling override fields to `google_spanner_instance` resource ([#26577](https://github.com/hashicorp/terraform-provider-google/pull/26577))
* sql: added `server_certificate_rotation_mode` field to `google_sql_database_instance` resource ([#26572](https://github.com/hashicorp/terraform-provider-google/pull/26572))
* storage: added `google_managed_encryption_enforcement_config`, `customer_managed_encryption_enforcement_config` and `customer_supplied_encryption_enforcement_config` to `google_storage_bucket` resource ([#26529](https://github.com/hashicorp/terraform-provider-google/pull/26529))

BUG FIXES:
* alloydb: fixed an issue where `password_wo` and `password_wo_version` fields were not functioning properly during update requests in `google_alloydb_user` resource ([#26571](https://github.com/hashicorp/terraform-provider-google/pull/26571))
* biglake: fixed erroneous diff for the `properties` field in the `google_biglake_iceberg_table` and `google_biglake_iceberg_namespace` resources ([#26595](https://github.com/hashicorp/terraform-provider-google/pull/26595))
* cloudfunctionsv2: fixed validation to only allow one of `direct_vpc_network_interface` or `vpc_connector` on `google_cloudfunctions2_function` resource ([#26567](https://github.com/hashicorp/terraform-provider-google/pull/26567))
* cloudrunv2: fixed validation to only allow one of `network_interfaces` or `connector` on `google_cloud_run_v2_service` and `google_cloud_run_v2_job` resources ([#26567](https://github.com/hashicorp/terraform-provider-google/pull/26567))
* compute: fixed `google_compute_region_backend_bucket` being present in the `google` (GA) provider. It is currently beta-only, and calls to the nonexistent GA API always returned a 404. ([#26597](https://github.com/hashicorp/terraform-provider-google/pull/26597))
* compute: fixed invalid update mask used for `rate_limit_options` field in `google_compute_region_security_policy_rule` resource ([#26527](https://github.com/hashicorp/terraform-provider-google/pull/26527))
* compute: fixed invalid update mask used for `rate_limit_options` field in `google_compute_security_policy` and `google_compute_security_policy_rule` resources ([#26526](https://github.com/hashicorp/terraform-provider-google/pull/26526))
* iambeta: fixed a perma-diff on `mode` field for `google_iam_workload_identity_pool` resource ([#26601](https://github.com/hashicorp/terraform-provider-google/pull/26601))
* provider: fixed an issue when custom endpoints use `http://` ([#26600](https://github.com/hashicorp/terraform-provider-google/pull/26600))
* vertexai: fixed operation calls in `google_vertex_ai_` resources not respecting `universe_domain` and `vertex_custom_endpoint` ([#26556](https://github.com/hashicorp/terraform-provider-google/pull/26556))


## 7.25.0 (Mar 24, 2026)

FEATURES:
* **New Data Source:** `google_compute_network_endpoint_groups` ([#26515](https://github.com/hashicorp/terraform-provider-google/pull/26515))
* **New Resource:** `google_dialogflow_environment` ([#26489](https://github.com/hashicorp/terraform-provider-google/pull/26489))
* **New Resource:** `google_kms_project_autokey_config` ([#26501](https://github.com/hashicorp/terraform-provider-google/pull/26501))

IMPROVEMENTS:
* backupdr: added `disk_backup_plan_properties` field to `google_backup_dr_backup_plan` resource ([#26497](https://github.com/hashicorp/terraform-provider-google/pull/26497))
* backupdr: made `backup_rules` optional in `google_backup_dr_backup_plan` resource  ([#26494](https://github.com/hashicorp/terraform-provider-google/pull/26494))
* blockchainnodeengine: added `ethereum_details.validator_config.beacon_fee_recipient` field to `google_blockchain_node_engine_blockchain_nodes` resource ([#26499](https://github.com/hashicorp/terraform-provider-google/pull/26499))
* ces: added `custom_headers` field to MCP toolset in CES `google_ces_toolset` resource ([#26473](https://github.com/hashicorp/terraform-provider-google/pull/26473))
* compute: added `expr` field  to `google_compute_organization_security_policy_rule` resource ([#26506](https://github.com/hashicorp/terraform-provider-google/pull/26506))
* compute: added `location` field to `google_network_services_tls_route` resource ([#26514](https://github.com/hashicorp/terraform-provider-google/pull/26514))
* compute: added `target_proxies` field to `google_network_services_tls_route` resource ([#26516](https://github.com/hashicorp/terraform-provider-google/pull/26516))
* compute: made `backend_service` field optional for resource `google_compute_target_tcp_proxy` ([#26519](https://github.com/hashicorp/terraform-provider-google/pull/26519))
* compute: made `backend_service` field optional for resource `google_compute_region_target_tcp_proxy` ([#26493](https://github.com/hashicorp/terraform-provider-google/pull/26493))
* iamworkforcepool: added `detailed_audit_logging` field to `google_iam_workforce_pool_provider` resource ([#26500](https://github.com/hashicorp/terraform-provider-google/pull/26500))
* kms: added `key_project_resolution_mode` field to `google_kms_autokey_config` resource ([#26501](https://github.com/hashicorp/terraform-provider-google/pull/26501))
* lustre: added `maintenance_policy` field to `google_lustre_instance` resource ([#26512](https://github.com/hashicorp/terraform-provider-google/pull/26512))
* sql: added `point_in_time_restore_context.region` field to `google_sql_database_instance` resource ([#26510](https://github.com/hashicorp/terraform-provider-google/pull/26510))
* vertexai: added `deletion_policy` field to `resource_vertex_ai_reasoning_engine` resource ([#26518](https://github.com/hashicorp/terraform-provider-google/pull/26518))

BUG FIXES:
* vertexai: fixed permadiff on `spec` field in `google_vertex_ai_reasoning_engine` resource ([#26470](https://github.com/hashicorp/terraform-provider-google/pull/26470))

## 7.24.0 (Mar 17, 2026)

DEPRECATIONS:
* iamworkforcepool: deprecated `extended_attributes_oauth2_client` on `google_iam_workforce_pool_provider`. Use `scim_usage` instead. ([#26388](https://github.com/hashicorp/terraform-provider-google/pull/26388))

FEATURES:
* **New Resource:** `google_biglake_iceberg_table` ([#26394](https://github.com/hashicorp/terraform-provider-google/pull/26394))
* **New Resource:** `google_contact_center_insights_auto_labeling_rule` ([#26426](https://github.com/hashicorp/terraform-provider-google/pull/26426))
* **New Resource:** `google_observability_trace_scope` ([#26428](https://github.com/hashicorp/terraform-provider-google/pull/26428))
* **New Resource:** `google_sql_provision_script` ([#26432](https://github.com/hashicorp/terraform-provider-google/pull/26432))

IMPROVEMENTS:
* ces: added Service Account OAuth `scopes` fields to `google_ces_toolset` resource ([#26368](https://github.com/hashicorp/terraform-provider-google/pull/26368))
* cloudrunv2: added `DISK` fields to `google_cloud_run_v2_service` resource ([#26418](https://github.com/hashicorp/terraform-provider-google/pull/26418))
* cloudsql: added `max_custom_on_demand_retention_days` field to `sqladmin` resource ([#26407](https://github.com/hashicorp/terraform-provider-google/pull/26407))
* compute: added `ForwardProxy` field in `google_compute_region_backend_service` resource ([#26449](https://github.com/hashicorp/terraform-provider-google/pull/26449))
* compute: added `accelerator_topology_mode` field to `google_compute_resource_policy` resource ([#26383](https://github.com/hashicorp/terraform-provider-google/pull/26383))
* compute: added `target_type` and `target_forwarding_rules` on `google_compute_region_network_firewall_policy_rule` resource ([#26369](https://github.com/hashicorp/terraform-provider-google/pull/26369))
* compute: promoted the `endpoint_url` field in  `google_compute_service_attachment` to GA ([#26434](https://github.com/hashicorp/terraform-provider-google/pull/26434))
* container: marked `subnetwork` as settable in `google_container_node_pool` ([#26416](https://github.com/hashicorp/terraform-provider-google/pull/26416))
* container: added `disruption_budget` field to `google_container_cluster` resource ([#26425](https://github.com/hashicorp/terraform-provider-google/pull/26425))
* discoveryengine: added `search_engine_config.required_subscription_tier ` field to `google_discovery_engine_search_engine` resource ([#26398](https://github.com/hashicorp/terraform-provider-google/pull/26398))
* discoveryengine: marked `content_config` as optional field in `google_discovery_engine_data_store` ([#26398](https://github.com/hashicorp/terraform-provider-google/pull/26398))
* memorystore: added `server_ca_mode` and `server_ca_pool` fields to `google_memorystore_instance` resource ([#26437](https://github.com/hashicorp/terraform-provider-google/pull/26437))
* networkservices: relaxed `authority` validation in `google_network_services_authz_extension` for different target types ([#26386](https://github.com/hashicorp/terraform-provider-google/pull/26386))
* redis: added `server_ca_mode` and `server_ca_pool` fields to `google_redis_cluster` resource ([#26437](https://github.com/hashicorp/terraform-provider-google/pull/26437))
* sql: added `clone_context.source_project` field to `google_sql_database_instance` resource to support cross project clone (beta) ([#26384](https://github.com/hashicorp/terraform-provider-google/pull/26384))
* transport: added automatic retry for GCE 403 errors with reason `CONCURRENT_OPERATIONS_QUOTA_EXCEEDED` ([#26417](https://github.com/hashicorp/terraform-provider-google/pull/26417))

BUG FIXES:
* compute: fixed perpetual diff for `oauth2_client_id` in `iap` block of `google_compute_backend_service` and `google_compute_region_backend_service` when disabling IAP ([#26385](https://github.com/hashicorp/terraform-provider-google/pull/26385))
* datastream: fixed an issue in `google_datastream_stream` where `source_config.mysql_source_config.binary_log_position` would show a diff when unset ([#26435](https://github.com/hashicorp/terraform-provider-google/pull/26435))
* workbench: marked `install-nvidia-driver` metadata key as settable for `google_workbench_instance` ([#26402](https://github.com/hashicorp/terraform-provider-google/pull/26402))

## 7.23.0 (Mar 10, 2026)

DEPRECATIONS:
* notebooks: `google_notebooks_environment` is deprecated and will be removed in a future major release. Use `google_workbench_instance` instead ([#26288](https://github.com/hashicorp/terraform-provider-google/pull/26288))
* provider: `google_*_iam_*` resources and datasources will now show deprecation messages when their parent resource has been deprecated ([#26288](https://github.com/hashicorp/terraform-provider-google/pull/26288))

FEATURES:
* **New Data Source:** `google_oracle_database_odb_network` ([#26290](https://github.com/hashicorp/terraform-provider-google/pull/26290))
* **New Data Source:** `google_oracle_database_odb_subnet` ([#26290](https://github.com/hashicorp/terraform-provider-google/pull/26290))
* **New Resource:** `google_vector_search_collection` ([#26353](https://github.com/hashicorp/terraform-provider-google/pull/26353))

IMPROVEMENTS:
* alloydb: added `dataplex_config` field to `google_alloydb_cluster` resource ([#26304](https://github.com/hashicorp/terraform-provider-google/pull/26304))
* biglake: added `primary_location` to `google_biglake_iceberg_catalog` resource ([#26307](https://github.com/hashicorp/terraform-provider-google/pull/26307))
* compute: added `params` field to `google_compute_external_vpn_gateway` resource ([#26348](https://github.com/hashicorp/terraform-provider-google/pull/26348))
* compute: added `params` field to `google_compute_ha_vpn_gateway` resource ([#26348](https://github.com/hashicorp/terraform-provider-google/pull/26348))
* compute: added `params` field to `google_compute_vpn_gateway` resource ([#26348](https://github.com/hashicorp/terraform-provider-google/pull/26348))
* compute: added `params` field to `google_compute_vpn_tunnel` resource ([#26348](https://github.com/hashicorp/terraform-provider-google/pull/26348))
* compute: added `storage_pool` support to `google_compute_instance_template` and `google_compute_region_instance_template` disks ([#26347](https://github.com/hashicorp/terraform-provider-google/pull/26347))
* container: added `control_plane_disk_encryption_key_versions` field to `user_managed_keys_config` in `google_container_cluster` resource ([#26289](https://github.com/hashicorp/terraform-provider-google/pull/26289))
* dataproc: added `cluster_type` to `google_dataproc_cluster` resource ([#26350](https://github.com/hashicorp/terraform-provider-google/pull/26350))
* dlp: added `actions.publish_to_scc`, `actions.publish_to_chronicle`, `actions.export_data.sample_findings_table` and `targets.big_query_target.filter.table_reference.project_id` fields to `google_data_loss_prevention_discovery_config` resource ([#26281](https://github.com/hashicorp/terraform-provider-google/pull/26281))
* gkebackup: added `protected_namespace_count` field to `google_gke_backup_backup_plan` resource ([#26283](https://github.com/hashicorp/terraform-provider-google/pull/26283))
* netapp: added `mode` field to `google_netapp_storage_pool` resource ([#26319](https://github.com/hashicorp/terraform-provider-google/pull/26319))
* osconfig: added `patch_config.skip_unpatchable_vms` field to `google_os_config_patch_deployment` resource ([#26282](https://github.com/hashicorp/terraform-provider-google/pull/26282))
* pubsub: added `text_config` field to `google_pubsub_subscription` resource ([#26329](https://github.com/hashicorp/terraform-provider-google/pull/26329))

BUG FIXES:
* tags: fixed iam read-after-write consistency issue with conditions in `google_tags_tag_key_iam_member` resource ([#26330](https://github.com/hashicorp/terraform-provider-google/pull/26330))

## 7.22.0 (Mar 3, 2026)

DEPRECATIONS:
* dataplex: deprecated `google_dataplex_data_asset`. Use `google_dataplex_data_product_data_asset` instead. ([#26256](https://github.com/hashicorp/terraform-provider-google/pull/26256))

FEATURES:
* **New Resource:** `google_compute_organization_security_policy_rule` ([#26202](https://github.com/hashicorp/terraform-provider-google/pull/26202))
* **New Resource:** `google_hypercomputecluster_cluster` ([#26180](https://github.com/hashicorp/terraform-provider-google/pull/26180))

IMPROVEMENTS:
* compute: `initialize_params.size` is now updatable in-place in the `google_compute_instance` resource ([#26195](https://github.com/hashicorp/terraform-provider-google/pull/26195))
* compute: added `dest_network_context`, `src_network_context` and `src_networks` fields to `google_compute_firewall_policy_rule` resource ([#26227](https://github.com/hashicorp/terraform-provider-google/pull/26227))
* compute: added `dest_network_context`, `src_network_context` and `src_networks` fields to `google_compute_network_firewall_policy_rule` resource ([#26227](https://github.com/hashicorp/terraform-provider-google/pull/26227))
* compute: added `dest_network_context`, `src_network_context` and `src_networks` fields to `google_compute_region_network_firewall_policy_rule` resource ([#26227](https://github.com/hashicorp/terraform-provider-google/pull/26227))
* container: promoted `sandbox_config` field in `google_container_cluster` and `google_container_node_pool` resources to GA ([#26247](https://github.com/hashicorp/terraform-provider-google/pull/26247))
* developerconnect: added `http_config` field to `google_developer_connect_connection` resource ([#26232](https://github.com/hashicorp/terraform-provider-google/pull/26232))
* filestore: added `source_backupdr_backup` field to `google_filestore_instance` resource ([#26238](https://github.com/hashicorp/terraform-provider-google/pull/26238))
* gkehub2: added field `spec.workloadidentity` to resource `google_gke_hub_feature` ([#26259](https://github.com/hashicorp/terraform-provider-google/pull/26259))
* iam: added AZURE_AD_GROUPS_DISPLAY_NAME enum value to `extra_attributes_oauth2_client.attribute-type` field in `google_iam_workforce_pool_provider` resource ([#26226](https://github.com/hashicorp/terraform-provider-google/pull/26226))
* kms: added a KMS AutokeyConfig-specific 10s post-create/post-update ([#26236](https://github.com/hashicorp/terraform-provider-google/pull/26236))
* networksecurity: added `url_filtering_profile` field  to `google_network_security_security_profile_group` resource ([#26266](https://github.com/hashicorp/terraform-provider-google/pull/26266))
* networksecurity: added `url_filtering_profile` field  to `google_network_security_security_profile` resource ([#26266](https://github.com/hashicorp/terraform-provider-google/pull/26266))
* networkservices: added support for use of multiple `ports` for `google_network_services_gateway` resources of type `SECURE_WEB_GATEWAY` ([#26265](https://github.com/hashicorp/terraform-provider-google/pull/26265))
* sql: added `auto_upgrade_enabled` field to `google_sql_database_instance` resource. ([#26205](https://github.com/hashicorp/terraform-provider-google/pull/26205))
* sql: added `data_api_access` field to `google_sql_database_instance` resource ([#26217](https://github.com/hashicorp/terraform-provider-google/pull/26217))
* sql: added `enhanced_query_insights_enabled` field to `google_sql_database_instance` resource ([#26244](https://github.com/hashicorp/terraform-provider-google/pull/26244))

BUG FIXES:
* datastream: fixed permadiff where `google_datastream_connection_profile.salesforce_profile.oauth2_client_credentials.client_id` is not read properly from the API ([#26201](https://github.com/hashicorp/terraform-provider-google/pull/26201))
* servicenetworking: added retry when creating `google_service_networking_connection` if it looks like the service account permissions haven't yet propagated ([#26220](https://github.com/hashicorp/terraform-provider-google/pull/26220))

## 7.21.0 (Feb 24, 2026)

FEATURES:
* **New Data Source:** `google_vmwareengine_announcements` ([#26145](https://github.com/hashicorp/terraform-provider-google/pull/26145))
* **New Data Source:** `google_vmwareengine_upgrades` ([#26174](https://github.com/hashicorp/terraform-provider-google/pull/26174))
* **New Resource:** `google_compute_region_backend_bucket` ([#26144](https://github.com/hashicorp/terraform-provider-google/pull/26144))
* **New Resource:** `google_hypercomputecluster_cluster` ([#26180](https://github.com/hashicorp/terraform-provider-google/pull/26180))
* **New Resource:** `google_network_services_agent_gateway` (beta) ([#26140](https://github.com/hashicorp/terraform-provider-google/pull/26140))

IMPROVEMENTS:
* beyondcorp: added `logging` field to `google_beyondcorp_security_gateway` resource ([#26159](https://github.com/hashicorp/terraform-provider-google/pull/26159))
* cloudfunctions2: added `direct_vpc_network_interface` and `direct_vpc_egress` fields to `google_cloudfunctions2_function` resource. Users who directly enabled DirectVPC on the underlying Cloud Run service will see a diff as a result of this update. ([#26142](https://github.com/hashicorp/terraform-provider-google/pull/26142))
* cloudrunv2: added the `iap_enabled` field to `google_cloud_run_v2_service` resource ([#26161](https://github.com/hashicorp/terraform-provider-google/pull/26161))
* dataproc: added `wait_for_completion` to `google_dataproc_job` resource ([#26177](https://github.com/hashicorp/terraform-provider-google/pull/26177))
* discoveryengine: added `disable_analytics` field to `google_discovery_engine_search_engine` resource ([#26171](https://github.com/hashicorp/terraform-provider-google/pull/26171))
* dlp: added `targets.cloud_storage_target.filter.collection.include_tags` block to `google_data_loss_prevention_discovery_config` resource ([#26178](https://github.com/hashicorp/terraform-provider-google/pull/26178))
* iap: added `client_id`, `client_secret`, and `client_secret_sha256` fields to `google_iap_settings` resource ([#26170](https://github.com/hashicorp/terraform-provider-google/pull/26170))
* networksecurity: added `mirroring_deployment_groups` and `mirroring_endpoint_group_type` fields to `google_network_security_security_profile` resource ([#26137](https://github.com/hashicorp/terraform-provider-google/pull/26137))

BUG FIXES:
* cloudrun: fixed perma-diff on `http_target.uri_override.query_override` in `google_cloud_tasks_queue` ([#26172](https://github.com/hashicorp/terraform-provider-google/pull/26172))
* storage: fixed a bug in `google_storage_bucket` where `force_destroy = true` would fail to delete buckets with large number of objects due to missing pagination ([#26164](https://github.com/hashicorp/terraform-provider-google/pull/26164))

## 7.20.0 (Feb 17, 2026)

FEATURES:
* **New Data Source:** `google_access_context_manager_supported_service` ([#26092](https://github.com/hashicorp/terraform-provider-google/pull/26092))
* **New Data Source:** `google_access_context_manager_supported_services` ([#26092](https://github.com/hashicorp/terraform-provider-google/pull/26092))
* **New Data Source:** `google_backup_dr_data_sources` ([#26080](https://github.com/hashicorp/terraform-provider-google/pull/26080))
* **New Data Source:** `google_kms_secret_asymmetric` ([#26096](https://github.com/hashicorp/terraform-provider-google/pull/26096))
* **New Data Source:** `google_storage_bucket_object_contents` ([#26054](https://github.com/hashicorp/terraform-provider-google/pull/26054))
* **New Resource:** `google_biglake_iceberg_namespace` ([#26076](https://github.com/hashicorp/terraform-provider-google/pull/26076))
* **New Resource:** `google_compute_rollout_plan` ([#26093](https://github.com/hashicorp/terraform-provider-google/pull/26093))
* **New Resource:** `google_oracle_database_exadb_vm_cluster` ([#26021](https://github.com/hashicorp/terraform-provider-google/pull/26021))
* **New Resource:** `google_vector_search_collection` ([#26098](https://github.com/hashicorp/terraform-provider-google/pull/26098))

IMPROVEMENTS:
* alloydb: added write-only support for `initial_user.password_wo` to `google_alloydb_cluster` ([#26074](https://github.com/hashicorp/terraform-provider-google/pull/26074))
* ces: added `mcp_toolset` field to `google_ces_toolset` resource ([#26025](https://github.com/hashicorp/terraform-provider-google/pull/26025))
* compute: added `allow_subnet_cidr_routes_overlap` field to `google_compute_subnetwork` resource ([#26019](https://github.com/hashicorp/terraform-provider-google/pull/26019))
* compute: added write-only support for `private_key` to `google_compute_region_ssl_certificate` resource ([#26072](https://github.com/hashicorp/terraform-provider-google/pull/26072))
* compute: added write-only support for `private_key` to `google_compute_ssl_certificate` resource ([#26072](https://github.com/hashicorp/terraform-provider-google/pull/26072))
* compute: added `enable` field to `google_compute_packet_mirroring` resource ([#26064](https://github.com/hashicorp/terraform-provider-google/pull/26064))
* compute: added `params` field to `google_compute_external_vpn_gateway` resource ([#26089](https://github.com/hashicorp/terraform-provider-google/pull/26089))
* compute: added `params` field to `google_compute_ha_vpn_gateway` resource ([#26089](https://github.com/hashicorp/terraform-provider-google/pull/26089))
* compute: added `params` field to `google_compute_interconnect_attachment` resource ([#26042](https://github.com/hashicorp/terraform-provider-google/pull/26042))
* compute: added `params` field to `google_compute_vpn_gateway` resource ([#26089](https://github.com/hashicorp/terraform-provider-google/pull/26089))
* compute: added `params` field to `google_compute_vpn_tunnel` resource ([#26089](https://github.com/hashicorp/terraform-provider-google/pull/26089))
* compute: added `slice_controller_config` field to `google_container_cluster` resource ([#26023](https://github.com/hashicorp/terraform-provider-google/pull/26023))
* container: added `additional_ip_ranges_config.status` to `google_container_cluster` resource ([#26061](https://github.com/hashicorp/terraform-provider-google/pull/26061))
* dataproc: added `instance_flexibility_policy` to `master_config` and `worker_config` in `google_dataproc_cluster` resource ([#26058](https://github.com/hashicorp/terraform-provider-google/pull/26058))
* developerconnect: added `target_projects` field to `google_developer_connect_insights_config` resource ([#26073](https://github.com/hashicorp/terraform-provider-google/pull/26073))
* filestore: added `replica_action` to `google_filestore_instance` resource ([#26082](https://github.com/hashicorp/terraform-provider-google/pull/26082))
* networksecurity: added `policy_profile`, `http_rules.0.to.0.operations.0.mcp` to `google_network_security_authz_policy` resource ([#26090](https://github.com/hashicorp/terraform-provider-google/pull/26090))
* networkservices: added `ull_multicast_domain` field to `google_network_services_multicast_domain` resource ([#26071](https://github.com/hashicorp/terraform-provider-google/pull/26071))
* networkservices: relaxed `load_balancing_scheme` validation to support non-Backend Service targets in `google_network_services_authz_extension` ([#26090](https://github.com/hashicorp/terraform-provider-google/pull/26090))
* spanner: added support for `user_project_override` in `google_spanner_database_iam` and `google_spanner_instance_iam` resources ([#26052](https://github.com/hashicorp/terraform-provider-google/pull/26052))
* vmwareengine: added `datastore_mount_config` field to `google_vmwareengine_cluster` resource ([#26083](https://github.com/hashicorp/terraform-provider-google/pull/26083))

BUG FIXES:
* bigquery: fixed permadiff with the `collation` field in `google_bigquery_table.schema` when it inherits the value from `google_bigquery_dataset.default_collation` ([#26065](https://github.com/hashicorp/terraform-provider-google/pull/26065))
* bigqueryanalyticshub: fixed update failure for `replica_locations` in `google_bigquery_analytics_hub_listing` ([#26046](https://github.com/hashicorp/terraform-provider-google/pull/26046))
* iam: fixed an issue where iam resources not retry on error 409 concurrent policy changes ([#26095](https://github.com/hashicorp/terraform-provider-google/pull/26095))
* publicca: fixed `mac_key` fields not being properly set in `google_public_ca_external_account_key` ([#26099](https://github.com/hashicorp/terraform-provider-google/pull/26099))

## 7.19.0 (Feb 10, 2026)

DEPRECATIONS:
* backupdr: `google_backupdr_restore_workload.name` is deprecated and will be removed in a future major release. The backup is identified by the parameters (location, backup_vault_id, data_source_id, backup_id). ([#25986](https://github.com/hashicorp/terraform-provider-google/pull/25986))
* publicca: `google_public_ca_external_account_key.b64url_mac_key` is deprecated and will be removed in a future major release. Use `mac_key` instead. ([#25964](https://github.com/hashicorp/terraform-provider-google/pull/25964))

FEATURES:
* **New Resource:** `google_network_security_mirroring_endpoint` ([#25988](https://github.com/hashicorp/terraform-provider-google/pull/25988))
* **New Resource:** `google_network_security_mirroring_endpoint_group` ([#25988](https://github.com/hashicorp/terraform-provider-google/pull/25988))
* **New Resource:** `google_backup_dr_restore_workload` ([#26013](https://github.com/hashicorp/terraform-provider-google/pull/26013))

IMPROVEMENTS:
* compute: added `network_pass_through_lb_traffic_policy` field to `google_compute_region_backend_service` resource ([#25994](https://github.com/hashicorp/terraform-provider-google/pull/25994))
* compute: added `RDMA_FALCON_POLICY` and `ULL_POLICY` values to `policy_type` field in `google_compute_region_network_firewall_policy`, `google_compute_region_network_firewall_policy_with_rules` ([#25985](https://github.com/hashicorp/terraform-provider-google/pull/25985))
* compute: added support for `network_interface.network_attachment` to `google_compute_instance_template` ([#25995](https://github.com/hashicorp/terraform-provider-google/pull/25995))
* compute: added support for `network_interface.network_attachment` to `google_compute_region_instance_template` ([#25995](https://github.com/hashicorp/terraform-provider-google/pull/25995))
* compute: added support for `network_interface.vlan` to `google_compute_instance_template`, enabling dynamic NIC ([#25995](https://github.com/hashicorp/terraform-provider-google/pull/25995))
* compute: added support for `network_interface.vlan` to `google_compute_instance`, enabling dynamic NIC. Creating and deleting from an existing instance is not yet supported. ([#25995](https://github.com/hashicorp/terraform-provider-google/pull/25995))
* compute: added support for `network_interface.vlan` to `google_compute_region_instance_template`, enabling dynamic NIC ([#25995](https://github.com/hashicorp/terraform-provider-google/pull/25995))
* discoveryengine: added `knowledge_graph_config` field to `google_discovery_engine_search_engine` resource ([#25980](https://github.com/hashicorp/terraform-provider-google/pull/25980))
* firestore: added `firestore_data_access_mode`, `mongodb_compatible_data_acess_mode`, and `realtime_updates_mode` fields to the `google_firestore_database` resource ([#26000](https://github.com/hashicorp/terraform-provider-google/pull/26000))
* firestore: added `deletion_policy` virtual field to `google_firestore_index` resource ([#25984](https://github.com/hashicorp/terraform-provider-google/pull/25984))
* monitoring: added write-only variants (`auth_token_wo` + `auth_token_wo_version`, `password_wo` + `password_wo_version`, `service_key_wo` + `service_key_wo_version`) for `google_monitoring_notification_channel.sensitive_labels` ([#25983](https://github.com/hashicorp/terraform-provider-google/pull/25983))
* networkconnectivity: added support for update operation on `google_network_connectivity_gateway_advertised_route` resource ([#25945](https://github.com/hashicorp/terraform-provider-google/pull/25945))
* provider: added a configurable `poll_interval` field to the provider for rare cases where it is being used in latency-sensitive situations. This can be set to a custom duration to change operation polling intervals. The default is unchanged, at `10s`. ([#26008](https://github.com/hashicorp/terraform-provider-google/pull/26008))
* publicca: added `mac_key` to `google_public_ca_external_account_key` ([#25964](https://github.com/hashicorp/terraform-provider-google/pull/25964))
* run: added `readiness_probe` field to `google_cloud_run_v2_service` resource ([#26003](https://github.com/hashicorp/terraform-provider-google/pull/26003))
* vertexai: added support for `developer_connect_source` to `spec.source_code_spec` in `google_vertex_ai_reasoning_engine` ([#26011](https://github.com/hashicorp/terraform-provider-google/pull/26011))

BUG FIXES:
* compute: fixed issue where it wasn't possible to set both `ssl_certificates` and `certificate_map` in `google_compute_target_ssl_proxy` ([#26012](https://github.com/hashicorp/terraform-provider-google/pull/26012))
* container: fixed an issue when toggling `default_compute_class_enabled` in `google_container_cluster` with Autopilot enabled ([#25966](https://github.com/hashicorp/terraform-provider-google/pull/25966))
* firebaseailogic: fixed bug in `google_firebase_ai_logic_config.generative_language_config.api_key_wo` where the value set wouldn't be sent to the API. ([#25983](https://github.com/hashicorp/terraform-provider-google/pull/25983))
* publicca: fixed `b64url_mac_key` sometimes being empty in `google_public_ca_external_account_key` ([#25964](https://github.com/hashicorp/terraform-provider-google/pull/25964))

## 7.18.0 (Feb 3, 2026)

BREAKING CHANGES:
* alloydb: removed the incorrect top-level field `last_successful_backup_consistency_time` from `google_backup_dr_backup_plan_association`. No value has been present in this output-only field. ([#25928](https://github.com/hashicorp/terraform-provider-google/pull/25928))

FEATURES:
* **New Resource:** `google_dataplex_data_asset` ([#25922](https://github.com/hashicorp/terraform-provider-google/pull/25922))
* **New Resource:** `google_logging_saved_query` ([#25921](https://github.com/hashicorp/terraform-provider-google/pull/25921))

IMPROVEMENTS:
* alloydb: added `restore_backupdr_backup_source`, `restore_backupdr_pitr_source`, and `backupdr_backup_source` to `google_alloydb_cluster` ([#25928](https://github.com/hashicorp/terraform-provider-google/pull/25928))
* alloydb: added `rules_config_info.last_successful_backup_consistency_time` to `google_backup_dr_backup_plan_association` ([#25928](https://github.com/hashicorp/terraform-provider-google/pull/25928))
* compute: updated `target_service` field to support update-in-place in `google_compute_service_attachment` resource ([#25924](https://github.com/hashicorp/terraform-provider-google/pull/25924))
* datafusion: added `patch_revision` field to `google_data_fusion_instance` resource ([#25923](https://github.com/hashicorp/terraform-provider-google/pull/25923))
* firestore: added `skip_wait` field to `google_firestore_index` resource, skipping the wait for index creation ([#25934](https://github.com/hashicorp/terraform-provider-google/pull/25934))
* gkeonprem: added `skip_validations` field to `google_gkeonprem_vmware_cluster` resource ([#25917](https://github.com/hashicorp/terraform-provider-google/pull/25917))
* sql: added `database_role` field and `iam_email` field to `google_sql_user` resource to support managing Cloud SQL users with database roles. ([#25926](https://github.com/hashicorp/terraform-provider-google/pull/25926))

BUG FIXES:
* cloudbuild: fixed `google_cloudbuild_trigger` to allow creation without source configuration for manual triggers ([#25925](https://github.com/hashicorp/terraform-provider-google/pull/25925))
* cloudrunv2: fix permadiff on `scaling.scaling_mode` in `google_cloud_run_v2_worker_pool` ([#25927](https://github.com/hashicorp/terraform-provider-google/pull/25927))
* compute: resolved issues where `show_nat_ips` and `nat_ips` in `google_compute_service_attachment` were causing test failures due to an underlying API problem. These fields are now temporarily non-functional and will be ignored. ([#25908](https://github.com/hashicorp/terraform-provider-google/pull/25908))
* container: fixed a bug in `google_container_node_pool` that prevented creation when `blue_green_settings` was specified ([#25916](https://github.com/hashicorp/terraform-provider-google/pull/25916))
* container: fixed perma-diff in `google_container_cluster` when setting `resource_limits` with disabled node autoprovisioning ([#25929](https://github.com/hashicorp/terraform-provider-google/pull/25929))

## 7.17.0 (January 27, 2026)

BREAKING CHANGES:
* networkconnectivity: changed `services` in `google_network_connectivity_multicloud_data_transfer_config` from TypeList to TypeSet. The order of or value of interpolations referencing the field may change. ([#25767](https://github.com/hashicorp/terraform-provider-google/pull/25767))

FEATURES:
* **New Resource:** `google_dataplex_data_product` ([#25844](https://github.com/hashicorp/terraform-provider-google/pull/25844))
* **New Resource:** `google_dialogflow_cx_tool_version` ([#25809](https://github.com/hashicorp/terraform-provider-google/pull/25809))
* **New Resource:** `google_firebase_ai_logic_config` ([#25846](https://github.com/hashicorp/terraform-provider-google/pull/25846))
* **New Resource:** `google_firebase_ai_logic_prompt_template` ([#25862](https://github.com/hashicorp/terraform-provider-google/pull/25862))
* **New Resource:** `google_firebase_ai_logic_prompt_template_lock` ([#25877](https://github.com/hashicorp/terraform-provider-google/pull/25877))
* **New Resource:** `google_saas_runtime_unit_operation` ([#25760](https://github.com/hashicorp/terraform-provider-google/pull/25760))
* **New Resource:** `google_vmwareengine_datastore` ([#25845](https://github.com/hashicorp/terraform-provider-google/pull/25845))
* **New Data Source:** `google_vmwareengine_datastore` ([#25845](https://github.com/hashicorp/terraform-provider-google/pull/25845))

IMPROVEMENTS:
* backupdr: added support for restore compute instance and disk ([#25723](https://github.com/hashicorp/terraform-provider-google/pull/25723))
* bigquery: added `source_column_match` field to `csv_options` in `google_bigquery_table` resource ([#25868](https://github.com/hashicorp/terraform-provider-google/pull/25868))
* compute: added `FIPS_202205` enum to `PROFILE` field in `SSL_POLICY` and `REGION_SSL_POLICY` resources, and added `TLS_1_3` enum to `MIN_TLS_VERSION` field in `SSL_POLICY` and `REGION_SSL_POLICY` resources. ([#25777](https://github.com/hashicorp/terraform-provider-google/pull/25777))
* compute: added `attachments` field to `google_compute_interconnect_attachment_group.logicalStructure.regions.metros.facilities.zones` and deprecated `attachment` field ([#25842](https://github.com/hashicorp/terraform-provider-google/pull/25842))
* compute: added `enable_enhanced_ipv4_allocation` field to `google_compute_public_delegated_prefix` resource ([#25732](https://github.com/hashicorp/terraform-provider-google/pull/25732))
* compute: added `ip_collection` field to `google_compute_address` resource ([#25732](https://github.com/hashicorp/terraform-provider-google/pull/25732))
* compute: added `source_instant_snapshot` field to `google_compute_snapshot` resource ([#25780](https://github.com/hashicorp/terraform-provider-google/pull/25780))
* compute: added support for "IF_L2_FORWARDING" as a value for the `availableFeatures` field of the `google_compute_interconnect` resource ([#25751](https://github.com/hashicorp/terraform-provider-google/pull/25751))
* compute: added support for "IF_L2_FORWARDING" as a value for the `requestedFeatures` field  of the `google_compute_interconnect` resource ([#25751](https://github.com/hashicorp/terraform-provider-google/pull/25751))
* compute: added support for "L2_DEDICATED" as a value for the `type` field of the `google_compute_interconnect_attachment` resource. ([#25751](https://github.com/hashicorp/terraform-provider-google/pull/25751))
* compute: added support for `igmp_query` field in `google_compute_instance`, `google_compute_instance_template`, and related instance resources. ([#25752](https://github.com/hashicorp/terraform-provider-google/pull/25752))
* compute: added support for the `l2Forwarding` field to `google_compute_interconnect_attachment` ([#25751](https://github.com/hashicorp/terraform-provider-google/pull/25751))
* compute: promoted `request_body_inspection_size` to GA in `google_compute_security_policy` resource (ga) ([#25775](https://github.com/hashicorp/terraform-provider-google/pull/25775))
* container: added `accelerator_network_config` field to `node_pool` resource ([#25856](https://github.com/hashicorp/terraform-provider-google/pull/25856))
* container: added `managed_opentelemetry_config` to `google_container_cluster` resource ([#25861](https://github.com/hashicorp/terraform-provider-google/pull/25861))
* container: added `node_drain_config` field to `google_container_node_pool` resources ([#25791](https://github.com/hashicorp/terraform-provider-google/pull/25791))
* container: improved `google_container_cluster` reconciliation time by caching node pools and instance group managers after a list call instead of getting each one seperately. ([#25784](https://github.com/hashicorp/terraform-provider-google/pull/25784))
* datastream: added `backfill_all.spanner_excluded_objects` and  `source_config.spanner_source_config` fields to `google_datastream_stream` ([#25804](https://github.com/hashicorp/terraform-provider-google/pull/25804))
* datastream: added `spanner_profile` field to `google_datastream_connection_profile` ([#25804](https://github.com/hashicorp/terraform-provider-google/pull/25804))
* dialogflowcx: added `serviceAccountAuthConfig ` field to `google_dialogflow_cx_webhook` resource ([#25781](https://github.com/hashicorp/terraform-provider-google/pull/25781))
* oracledatabase: added `peerAutonomousDatabases`, `disasterRecoverySupportedLocations`, `sourceConfig` fields to Autonomous database resource. ([#25859](https://github.com/hashicorp/terraform-provider-google/pull/25859))
* tags: added `allowed_values_regex` field to `google_tags_tag_key` resource ([#25869](https://github.com/hashicorp/terraform-provider-google/pull/25869))
* tags: added support for dynamic tag keys in `google_tags_tag_binding` and `google_tags_location_tag_binding`  resources ([#25874](https://github.com/hashicorp/terraform-provider-google/pull/25874))
* vertex_ai: added `deployment_spec.psc_interface_config` to `google_vertex_ai_reasoning_engine` ([#25765](https://github.com/hashicorp/terraform-provider-google/pull/25765))

BUG FIXES:
* bigquery: fixed permadiff with the `collation` field in `google_bigquery_table.schema` ([#25762](https://github.com/hashicorp/terraform-provider-google/pull/25762))
* cloudasset: fixed bug in `google_cloud_asset_folder_feed` where `folder_id` was always empty ([#25798](https://github.com/hashicorp/terraform-provider-google/pull/25798))
* cloudbuild: fixed permadiff on `google_cloudbuild_trigger.pubsub_config.service_account_email` ([#25792](https://github.com/hashicorp/terraform-provider-google/pull/25792))
* compute: fix crash when specifying an empty `instance_flexibility_policy` block on the `google_compute_region_instance_group_manager` resource ([#25731](https://github.com/hashicorp/terraform-provider-google/pull/25731))
* compute: fixed a permadiff that could occur when using mixed short and long form IPv6 addresses in the `source_ranges` field of `google_compute_firewall` ([#25867](https://github.com/hashicorp/terraform-provider-google/pull/25867))
* iambeta: fixed a permadiff that could occur in the `jwks_json` field for `google_iam_workload_identity_pool_provider` resource ([#25847](https://github.com/hashicorp/terraform-provider-google/pull/25847))
* netapp: fixed export_policy update bug with squash_mode in netapp volume ([#25776](https://github.com/hashicorp/terraform-provider-google/pull/25776))
* networkconnectivity: fixed a diff on `services` in `google_network_connectivity_multicloud_data_transfer_config` reordering elements ([#25767](https://github.com/hashicorp/terraform-provider-google/pull/25767))
* sql: fixed an issue where transient server errors caused false failures for SQL operations that eventually completed successfully ([#25735](https://github.com/hashicorp/terraform-provider-google/pull/25735))
* workbench: made `enable-jupyterlab4` metadata key settable for `google_workbench_instance` ([#25769](https://github.com/hashicorp/terraform-provider-google/pull/25769))

## 7.16.0 (January 13, 2026)

DEPRECATIONS:
* cloudrunv2: deprecated `custom_audience` field in the `google_cloud_run_v2_worker_pool` resource, as this field is not applicable to the WorkerPools resource ([#25688](https://github.com/hashicorp/terraform-provider-google/pull/25688))

FEATURES:
* **New Data Source:** `google_compute_routers` ([#25715](https://github.com/hashicorp/terraform-provider-google/pull/25715))
* **New Resource:** `google_backup_dr_restore_workload` ([#25723](https://github.com/hashicorp/terraform-provider-google/pull/25723))

IMPROVEMENTS:
* backupdr: added `max_custom_on_demand_retention_days` field to `google_backup_dr_backup_plan` resource ([#25704](https://github.com/hashicorp/terraform-provider-google/pull/25704))
* bigquery: added support for merge and update operations for dataPolicies in `schema` field in `google_bigquery_table` resource when `ignore_schema_changes` is defined ([#25721](https://github.com/hashicorp/terraform-provider-google/pull/25721))
* bigtable: added `etag` field to `google_bigtable_schema_bundle` resource ([#25687](https://github.com/hashicorp/terraform-provider-google/pull/25687))
* compute: added `BPS_400G` enum value to `bandwidth` field in `google_compute_interconnect_attachment` resource ([#25714](https://github.com/hashicorp/terraform-provider-google/pull/25714))
* container: added `registry_hosts` field to `containerd_config` in `google_container_cluster` and `google_container_node_pool` resources ([#25705](https://github.com/hashicorp/terraform-provider-google/pull/25705))
* dataplex: added `one_time` field to `google_dataplex_datascan` resource ([#25695](https://github.com/hashicorp/terraform-provider-google/pull/25695))
* datastream: added `postgresql_profile.ssl_config` to `google_datastream_connection_profile` resource ([#25671](https://github.com/hashicorp/terraform-provider-google/pull/25671))
* networkservices: added `EXT_AUTHZ_GRPC` enum value to `wire_format` field in `google_network_services_authz_extension` resource ([#25706](https://github.com/hashicorp/terraform-provider-google/pull/25706))
* networkservices: added `disable_placement_policy` field to `google_network_services_multicast_domain_activation` resource ([#25720](https://github.com/hashicorp/terraform-provider-google/pull/25720))
* networkservices: added `metadata`, `supported_events`, `request_body_send_mode`, and `observability_mode` fields to `google_network_services_lb_route_extension` resource ([#25702](https://github.com/hashicorp/terraform-provider-google/pull/25702))
* securitycenterv2: added support for supplying `location` values other than "GLOBAL" to the `google_scc_v2_project_notification_config` resource ([#25698](https://github.com/hashicorp/terraform-provider-google/pull/25698))
* storageinsights: added `activity_data_retention_period_days` field to `google_storage_insights_dataset_config` resource ([#25703](https://github.com/hashicorp/terraform-provider-google/pull/25703))
* workbench: added support to set post-startup script metadata keys with managed EUC in `google_workbench_instance` resource ([#25719](https://github.com/hashicorp/terraform-provider-google/pull/25719))

## 7.15.0 (January 6, 2026)

NOTES:
* lustre: increased delete and update operation timeouts from 20 minutes to 60 minutes for `google_lustre_instance` resource ([#25662](https://github.com/hashicorp/terraform-provider-google/pull/25662))

BREAKING CHANGES:
* compute: changed `cipher_suite` fields in the `google_compute_vpn_tunnel` resource to track order ([#25657](https://github.com/hashicorp/terraform-provider-google/pull/25657))

FEATURES:
* **New Resource:** `google_apigee_security_feedback` ([#25589](https://github.com/hashicorp/terraform-provider-google/pull/25589))
* **New Resource:** `google_apphub_boundary` ([#25640](https://github.com/hashicorp/terraform-provider-google/pull/25640))
* **New Resource:** `google_biglake_iceberg_catalog_iam_binding` ([#25638](https://github.com/hashicorp/terraform-provider-google/pull/25638))
* **New Resource:** `google_biglake_iceberg_catalog_iam_member` ([#25638](https://github.com/hashicorp/terraform-provider-google/pull/25638))
* **New Resource:** `google_biglake_iceberg_catalog_iam_policy` ([#25638](https://github.com/hashicorp/terraform-provider-google/pull/25638))
* **New Resource:** `google_biglake_iceberg_catalog` ([#25528](https://github.com/hashicorp/terraform-provider-google/pull/25528))
* **New Resource:** `google_compute_organization_security_policy_association` ([#25643](https://github.com/hashicorp/terraform-provider-google/pull/25643))
* **New Resource:** `google_network_connectivity_destination` ([#25663](https://github.com/hashicorp/terraform-provider-google/pull/25663))
* **New Resource:** `google_network_connectivity_multicloud_data_transfer_config` ([#25609](https://github.com/hashicorp/terraform-provider-google/pull/25609))
* **New Resource:** `google_network_security_dns_threat_detector` ([#25634](https://github.com/hashicorp/terraform-provider-google/pull/25634))

IMPROVEMENTS:
* backupdr: added ignore_read to `encryption_config` field in `google_backup_dr_backup_vault` resource ([#25685](https://github.com/hashicorp/terraform-provider-google/pull/25685))
* biglakeiceberg: made `google_biglake_iceberg_catalog` use the resource project as the quota project when `user_project_override` is `true` ([#25638](https://github.com/hashicorp/terraform-provider-google/pull/25638))
* composer: added new enum `ENVIRONMENT_SIZE_EXTRA_LARGE` to `environment_size` field to `google_composer_environment` resource ([#25531](https://github.com/hashicorp/terraform-provider-google/pull/25531))
* compute: added `candidate_cloud_router_ip_address`, `candidate_customer_router_ip_address`, `candidate_cloud_router_ipv6_address`, and `candidate_customer_router_ipv6_address` fields to `google_compute_interconnect_attachment` resource ([#25581](https://github.com/hashicorp/terraform-provider-google/pull/25581))
* compute: added `prefix_length` field to `google_compute_addresses` data source ([#25654](https://github.com/hashicorp/terraform-provider-google/pull/25654))
* compute: added `client_destination_port` and `instance` fields to `google_compute_region_network_endpoints` resource ([#25621](https://github.com/hashicorp/terraform-provider-google/pull/25621))
* datastream: added support for the `rule_sets` field in the `google_datastream_stream` resource, allowing configuration of customization rules, such as BigQuery destinations partitioning and clustering. ([#25529](https://github.com/hashicorp/terraform-provider-google/pull/25529))
* iamworkforcepool: added `hard_delete` support in `google_iam_workforce_pool_provider_scim_tenant` resource ([#25656](https://github.com/hashicorp/terraform-provider-google/pull/25656))
* looker: added `periodic_export_config` field to `google_looker_instance` resource ([#25610](https://github.com/hashicorp/terraform-provider-google/pull/25610))
* lustre: added `access_rules_options` field to `google_lustre_instance` resource to support root squashing and IP-based access control configuration ([#25617](https://github.com/hashicorp/terraform-provider-google/pull/25617))
* managedkafka: replaced `disk_size_gb` with `disk_size_gib` in `broker_capacity_config` within the `google_managed_kafka_cluster` resource ([#25613](https://github.com/hashicorp/terraform-provider-google/pull/25613))
* networkservices: added `state` field to `google_network_services_multicast_domain` resource ([#25532](https://github.com/hashicorp/terraform-provider-google/pull/25532))
* redis: added `labels` to `google_redis_cluster` ([#25639](https://github.com/hashicorp/terraform-provider-google/pull/25639))
* sql: marked `replication_cluster.psa_write_endpoint` field as Computed in `google_sql_database_instance` resource ([#25573](https://github.com/hashicorp/terraform-provider-google/pull/25573))
* sql: set `replication_cluster` when update `google_sql_database_instance` resource if there is a disaster recovery(DR) replica set or there is a PSA write endpoint ([#25573](https://github.com/hashicorp/terraform-provider-google/pull/25573))
* storage: updated datasource `google_storage_object_signed_url.signed_url` to use virtual style hosted url ([#25568](https://github.com/hashicorp/terraform-provider-google/pull/25568))
* vertexai: added `bigtable`, `zone`, `encryption_spec`, and `bigtable_options` fields to `google_vertex_ai_feature_online_store` resource ([#25601](https://github.com/hashicorp/terraform-provider-google/pull/25601))
* vertexai: added `psc_automation_configs` to resource `google_vertex_ai_index_endpoint` ([#25570](https://github.com/hashicorp/terraform-provider-google/pull/25570))

BUG FIXES:
* provider: fixed an issue where error type 409 and 412 were not being correctly retried. This commonly shows up in IAM resources, but can appear in other resources as well ([#25596](https://github.com/hashicorp/terraform-provider-google/pull/25596))
* alloydb: fixed an issue where boolean fields were ignored when set to `false` for `google_alloydb_cluster` and `google_alloydb_instance` ([#25561](https://github.com/hashicorp/terraform-provider-google/pull/25561))
* cloudrunv2: fixed a permadiff when default values of the `scaling` block were explicitly declared on the `google_cloud_run_v2_service` resource ([#25569](https://github.com/hashicorp/terraform-provider-google/pull/25569))
* compute: fixed a crash in `google_compute_disk`/`google_compute_region_disk` when deleting a disk attached to an instance that had any scratch disks attached ([#25641](https://github.com/hashicorp/terraform-provider-google/pull/25641))
* compute: fixed issue where `endpoints.interconnects.vlan_tags` wouldn't be read correctly from the API  in `google_compute_wire_group` resource ([#25602](https://github.com/hashicorp/terraform-provider-google/pull/25602))
* compute: fixed update logic that causes empty instance being sent for `google_compute_network_endpoints` ([#25621](https://github.com/hashicorp/terraform-provider-google/pull/25621))
* datacatalog: fixed issue where `fields.display_name` wouldn't be read correctly from the API in `google_data_catalog_tag` resource ([#25602](https://github.com/hashicorp/terraform-provider-google/pull/25602))
* discoveryengine: marked `cmek_config_id` field in `google_discovery_engine_cmek_config` resource as required ([#25527](https://github.com/hashicorp/terraform-provider-google/pull/25527))
* securitygateway: allowed empty field for `service_discovery` in `google_beyondcorp_security_gateway` ([#25653](https://github.com/hashicorp/terraform-provider-google/pull/25653))
* securitygateway: allowed empty fields for `user_info`, `group_info` and `device_info` in `google_beyondcorp_security_gateway` ([#25653](https://github.com/hashicorp/terraform-provider-google/pull/25653))
* servicedirectory: fixed an issue where `google_service_directory_endpoint` or `google_service_directory_service` without `metadata` specified would have other fields removed ([#25588](https://github.com/hashicorp/terraform-provider-google/pull/25588))
* storage: fixed the behavior in `google_storage_bucket` resource when `force_destroy` is set to `true`. Previously, failing to list anywhere caches would prevent destroying objects on the bucket. Now, both objects and caches are processed independently. ([#25655](https://github.com/hashicorp/terraform-provider-google/pull/25655))

## 7.14.0 (December 16, 2025)

DEPRECATIONS:
* managedkafka: added deprecation warning for `google_managed_kafka_connect_cluster` `additional_subnets` field ([#25487](https://github.com/hashicorp/terraform-provider-google/pull/25487))

FEATURES:
* **New Data Source:** `google_artifact_registry_versions` ([#25512](https://github.com/hashicorp/terraform-provider-google/pull/25512))
* **New Data Source:** `google_cloud_identity_policies` ([#25513](https://github.com/hashicorp/terraform-provider-google/pull/25513))
* **New Data Source:** `google_compute_region_security_policy` ([#25488](https://github.com/hashicorp/terraform-provider-google/pull/25488))
* **New Data Source:** `google_compute_storage_pool` ([#25485](https://github.com/hashicorp/terraform-provider-google/pull/25485))
* **New Resource:** `google_compute_cross_site_network` ([#25479](https://github.com/hashicorp/terraform-provider-google/pull/25479))
* **New Resource:** `google_compute_wire_group` ([#25479](https://github.com/hashicorp/terraform-provider-google/pull/25479))
* **New Resource:** `google_network_services_multicast_group_consumer_activation` ([#25515](https://github.com/hashicorp/terraform-provider-google/pull/25515))
* **New Resource:** `google_network_services_multicast_group_producer_activation` ([#25472](https://github.com/hashicorp/terraform-provider-google/pull/25472))

IMPROVEMENTS:
* alloydb: added `connection_pool_config`, `connection_pool_config.enabled` and `connection_pool_config.flags` in `google_alloydb_instance` resource ([#25484](https://github.com/hashicorp/terraform-provider-google/pull/25484))
* colab: added `software_config.post_startup_script_config` field to `google_colab_runtime_template` ([#25509](https://github.com/hashicorp/terraform-provider-google/pull/25509))
* compute: added new field `instance_flexibility_policy.instance_selection.min_cpu_platform` & `instance_flexibility_policy.instance_selection.disks` to `google_compute_region_instance_group_manager` ([#25444](https://github.com/hashicorp/terraform-provider-google/pull/25444))
* dataplex: removed the need for import in `google_dataplex_entry` when using first party source systems ([#25507](https://github.com/hashicorp/terraform-provider-google/pull/25507))
* dataproc: added `auto_stop_time` and  `idle_stop_ttl` to `google_dataproc_cluster` resource ([#25456](https://github.com/hashicorp/terraform-provider-google/pull/25456))
* eventarc: added `retry_policy` field to `google_eventarc_trigger` resource ([#25467](https://github.com/hashicorp/terraform-provider-google/pull/25467))
* networksecurity: enabled in-place update for `custom_mirroring_profile.mirroring_deployment_groups` on `google_network_security_security_profile` ([#25508](https://github.com/hashicorp/terraform-provider-google/pull/25508))
* spanner: added `autoscaling_config.autoscaling_targets.total_cpu_utilization_percent` field to  `google_spanner_instance` resource ([#25495](https://github.com/hashicorp/terraform-provider-google/pull/25495))
* sql: added changes to ignore changes in backup configuration's fields like `enabled`, `binary_log_enabled`, `start_time`, `point_in_time_recovery_enabled`, `transaction_log_retention_days` and `backup_retention_settings.retained_backups` in `google_sql_database_instance` if the instance is managed by Google Cloud Backup and Disaster (DR) Recovery Service. ([#25516](https://github.com/hashicorp/terraform-provider-google/pull/25516))

BUG FIXES:
* compute: fixed `google_compute_network` in-place update to enable `enable_ula_internal_ipv6`. ([#25468](https://github.com/hashicorp/terraform-provider-google/pull/25468))
* iam: fixed error 409 concurrency policy changes by correctly detecting the error type. ([#25473](https://github.com/hashicorp/terraform-provider-google/pull/25473))
* sql: fixed an issue where the computed `psc_service_attachment_link` attribute was not being exported properly in `google_sql_database_instance` resource and datasources ([#25510](https://github.com/hashicorp/terraform-provider-google/pull/25510))

## 7.13.0 (December 9, 2025)

NOTES:
* alloydb: reverted requiring `initial_user.password` as required on create for new `google_alloydb_cluster` resources, instead `initial_user.password` or `initial_user.user` must be set if `initial_user` is specified for `google_alloydb_cluster` resources ([#25366](https://github.com/hashicorp/terraform-provider-google/pull/25366))
* privateca: modified `encryption_spec` field from `google_privateca_ca_pool` resource to be mutable and allow cmek key rotation ([#25267](https://github.com/hashicorp/terraform-provider-google/pull/25267))

DEPRECATIONS:
* cloudquotas: deprecated `effective_container` and `effective_enablement` fields in the `google_cloud_quotas_quota_adjuster_settings` resource ([#25443](https://github.com/hashicorp/terraform-provider-google/pull/25443))
* dlp: deprecated `publish_findings_to_cloud_data_catalog` field in `google_data_loss_prevention_job_trigger` resource. Use `publish_findings_to_dataplex_catalog` field instead. ([#25250](https://github.com/hashicorp/terraform-provider-google/pull/25250))
* networkservices: removed `google_service_binding` resource due to service binding support being disabled ([#25367](https://github.com/hashicorp/terraform-provider-google/pull/25367))

FEATURES:
* **New Resource:** `google_ces_app_version` ([#25297](https://github.com/hashicorp/terraform-provider-google/pull/25297))
* **New Resource:** `google_compute_organization_security_policy` ([#25322](https://github.com/hashicorp/terraform-provider-google/pull/25322))
* **New Resource:** `google_dialogflow_generator` ([#25340](https://github.com/hashicorp/terraform-provider-google/pull/25340))
* **New Resource:** `google_dialogflow_version` ([#25179](https://github.com/hashicorp/terraform-provider-google/pull/25179))
* **New Resource:** `google_discovery_engine_widget_config` ([#25378](https://github.com/hashicorp/terraform-provider-google/pull/25378))
* **New Resource:** `google_iam_workforce_pool_provider_scim_token` ([#25270](https://github.com/hashicorp/terraform-provider-google/pull/25270))
* **New Resource:** `google_network_services_lb_edge_extension` ([#25299](https://github.com/hashicorp/terraform-provider-google/pull/25299))
* **New Resource:** `google_network_services_multicast_consumer_association` ([#25321](https://github.com/hashicorp/terraform-provider-google/pull/25321))
* **New Resource:** `google_network_services_multicast_group_range_activation` ([#25386](https://github.com/hashicorp/terraform-provider-google/pull/25386))
* **New Resource:** `google_network_services_multicast_group_range` ([#25353](https://github.com/hashicorp/terraform-provider-google/pull/25353))
* **New Resource:** `google_network_services_multicast_producer_association` ([#25291](https://github.com/hashicorp/terraform-provider-google/pull/25291))

IMPROVEMENTS:
* alloydb: added `password_wo` and `password_wo_version` fields to `google_alloydb_user` resource ([#25266](https://github.com/hashicorp/terraform-provider-google/pull/25266))
* apphub: added `identity` field to `google_apphub_service` and `google_apphub_workload` resources ([#25363](https://github.com/hashicorp/terraform-provider-google/pull/25363))
* backupdr: added `encryption_config` field to `google_backup_dr_backup_vault` resource ([#25221](https://github.com/hashicorp/terraform-provider-google/pull/25221))
* ces: added `client_function.parameters.max_items`, `client_function.parameters.min_items`, `client_function.parameters.maximum`, `client_function.parameters.minimum`, `client_function.parameters.title`, `client_function.response.max_items`, `client_function.response.min_items`, `client_function.response.maximum`, `client_function.response.minimum`, and `client_function.response.title` fields to `google_ces_tool` resource ([#25309](https://github.com/hashicorp/terraform-provider-google/pull/25309))
* ces: added `entry_agent` field to `google_ces_example` resource ([#25182](https://github.com/hashicorp/terraform-provider-google/pull/25182))
* ces: added `google_search_tool.context_urls`, `google_search_tool.preferred_domains`, and `open_api_tool.api_authentication.bearer_token_config` fields to `google_ces_tool` resource ([#25309](https://github.com/hashicorp/terraform-provider-google/pull/25309))
* ces: added `message.chunk.tool_response` and `message.chunk.tool_call` fields to `google_ces_example` resource ([#25182](https://github.com/hashicorp/terraform-provider-google/pull/25182))
* ces: added `pinned` and `variable_declarations.schema.title` fields to `google_ces_app` resource ([#25233](https://github.com/hashicorp/terraform-provider-google/pull/25233))
* cloudsecuritycompliance: added `cloud_control_details.parameters.parameter_value.oneof_value` fields to `google_cloud_security_compliance_framework_deployment` resource ([#25382](https://github.com/hashicorp/terraform-provider-google/pull/25382))
* cloudsecuritycompliance: added `cloud_control_details.parameters.parameter_value.oneof_value` fields to `google_cloud_security_compliance_framework` resource ([#25382](https://github.com/hashicorp/terraform-provider-google/pull/25382))
* cloudsecuritycompliance: added `parameter_spec.default_value.oneof_value` and `validation.allowed_values.values.oneof_value` fields to `google_cloud_security_compliance_cloud_control ` resource ([#25441](https://github.com/hashicorp/terraform-provider-google/pull/25441))
* cloudsecuritycompliance: added `sub_parameters` field to `google_cloud_security_compliance_cloud_control ` resource ([#25441](https://github.com/hashicorp/terraform-provider-google/pull/25441))
* colab: added `custom_environment_spec` field to `google_colab_notebook_execution` resource ([#25379](https://github.com/hashicorp/terraform-provider-google/pull/25379))
* compute: added `network_pass_through_lb_traffic_policy` field to `google_compute_region_backend_service` resource. ([#25223](https://github.com/hashicorp/terraform-provider-google/pull/25223))
* compute: added `params` field to `google_compute_interconnect` resource ([#25350](https://github.com/hashicorp/terraform-provider-google/pull/25350))
* compute: added `show_nat_ips` and `nat_ips` fields to `google_compute_service_attachment` ([#25296](https://github.com/hashicorp/terraform-provider-google/pull/25296))
* compute: added `snapshot_type` field to `google_compute_snapshot` resource ([#25348](https://github.com/hashicorp/terraform-provider-google/pull/25348))
* compute: added new field `instance_flexibility_policy.instance_selection.min_cpu_platform` & `instance_flexibility_policy.instance_selection.disks` to `google_compute_region_instance_group_manager` ([#25444](https://github.com/hashicorp/terraform-provider-google/pull/25444))
* container: added `autoscaled_rollout_policy` field to `google_container_node_pool` resource (beta) ([#25362](https://github.com/hashicorp/terraform-provider-google/pull/25362))
* container: added `node_kernel_module_loading.policy` field to `google_container_node_pool` and `google_container_cluster` resources ([#25383](https://github.com/hashicorp/terraform-provider-google/pull/25383))
* filestore: added support for updating `directory_services` fields in place in `google_filestore_instance` ([#25315](https://github.com/hashicorp/terraform-provider-google/pull/25315))
* iamworkforcepool: added `claim_mapping`, `purge_time`, and `service_agent` fields to `google_iam_workforce_pool_provider_scim_tenant` resource ([#25270](https://github.com/hashicorp/terraform-provider-google/pull/25270))
* looker: added `controlled_egress_enabled` and `controlled_egress_config` fields to `google_looker_instance` resource ([#25214](https://github.com/hashicorp/terraform-provider-google/pull/25214))
* lustre: added `kms_key` field to `google_lustre_instance` resource ([#25261](https://github.com/hashicorp/terraform-provider-google/pull/25261))
* modelarmor: added `google_mcp_server_floor_setting` field to `google_model_armor_floorsetting ` resource ([#25313](https://github.com/hashicorp/terraform-provider-google/pull/25313))
* monitoring: fixes an issue with `google_monitoring_alert_policy` where it ignores the resource project during Import ([#25287](https://github.com/hashicorp/terraform-provider-google/pull/25287))
* netapp: added public docs link for `google_netapp_host_group` resource ([#25368](https://github.com/hashicorp/terraform-provider-google/pull/25368))
* netapp: added 'nfsv4' to custom update export_policy object in `google_netapp_volume` resource ([#25442](https://github.com/hashicorp/terraform-provider-google/pull/25442))
* oracledatabase: added `properties.cpu_core_count`, `properties.secret_id`, and `properties.vault_id` fields to `google_oracle_database_autonomous` resource ([#25264](https://github.com/hashicorp/terraform-provider-google/pull/25264))
* oracledatabase: added `properties.time_zone.version` field to `google_oracle_database_cloud_vm_cluster` resource ([#25264](https://github.com/hashicorp/terraform-provider-google/pull/25264))
* servicedirectory: promoted `google_service_directory_namespace`, `google_service_directory_service`, and `google_service_directory_endpoint` to GA ([#25177](https://github.com/hashicorp/terraform-provider-google/pull/25177))
* servicedirectory: replaced `metadata` KeyValuePair with `annotations` KeyValueAnnotations in  `google_service_directory_service`, and `google_service_directory_endpoint` resources ([#25177](https://github.com/hashicorp/terraform-provider-google/pull/25177))
* sql: added write-only argument for `root_password` in `google_sql_database_instance` resource ([#25252](https://github.com/hashicorp/terraform-provider-google/pull/25252))
* storage: added `contexts` for resource `google_storage_bucket_object` ([#25346](https://github.com/hashicorp/terraform-provider-google/pull/25346))
* vertex_ai: added `resourceLimits`, `minInstances`, `maxInstances`, `containerConcurrency` and `sourceCodeSpec` fields to `google_vertex_ai_reasoning_engine` resource ([#25349](https://github.com/hashicorp/terraform-provider-google/pull/25349))

BUG FIXES:
* bigquery: fixed the permadiff when email field values contain non-lower-case characters in `access` in `google_bigquery_dataset` ([#25317](https://github.com/hashicorp/terraform-provider-google/pull/25317))
* bigquery: fixed the permadiff when table schema is unchanged for a `google_bigquery_table` with row access policies ([#25256](https://github.com/hashicorp/terraform-provider-google/pull/25256))
* cloudrunv2: fixed permadiff if `scaling` field is unset on resource `google_cloud_run_v2_service` ([#25310](https://github.com/hashicorp/terraform-provider-google/pull/25310))
* compute: fixed an issue where the `bgp_always_compare_med` field could not be unset in in `google_compute_network`. It can now be unset by configuring the new field `delete_bgp_always_compare_med` to a value of `true`. ([#25288](https://github.com/hashicorp/terraform-provider-google/pull/25288))
* compute: fixed crashes when no `network_endpoints` block specified in `google_compute_network_endpoints` resource or no network endpoints exist ([#25220](https://github.com/hashicorp/terraform-provider-google/pull/25220))
* compute: fixed the `terms` field in `google_compute_router_route_policy` to be updatable without forcing resource recreation ([#25289](https://github.com/hashicorp/terraform-provider-google/pull/25289))
* container: fixed a perpetual diff in `google_container_cluster` resource when `enable_l4_ilb_subsetting` is enabled by the GKE control plane and not explicitly set in the configuration ([#25323](https://github.com/hashicorp/terraform-provider-google/pull/25323))
* dialogflowcx: fixed update_mask in `google_dialogflow_cx_playbook` where a granular update mask is required. ([#25254](https://github.com/hashicorp/terraform-provider-google/pull/25254))
* discoveryengine: fixed a permadiff on `advanced_site_search_config` in `google_discovery_engine_data_store` resource ([#25387](https://github.com/hashicorp/terraform-provider-google/pull/25387))
* iamworkforcepool: fixed bug in `google_iam_workforce_pool_provider_scim_token` where `base_uri` wasn't set correctly from the API ([#25270](https://github.com/hashicorp/terraform-provider-google/pull/25270))
* logging: fixed an issue with `google_logging_*_sink.include_children` fields not being updatable to true ([#25247](https://github.com/hashicorp/terraform-provider-google/pull/25247))
* memorystore: fixed an issue where a permadiff on `desired_auto_created_endpoints` caused the `google_memorystore_instance` resource to recreated. ([#25278](https://github.com/hashicorp/terraform-provider-google/pull/25278))
* spanner: prevented recreation when `kms_key_name` and `kms_key_names` are same for `google_spanner_database` ([#25215](https://github.com/hashicorp/terraform-provider-google/pull/25215))

## 7.12.0 (November 18, 2025)

DEPRECATIONS:
* backupdr: deprecated `required_type` in `google_backup_dr_backup_plan_associations` and `google_backup_dr_data_source_references`. Both resources no longer have functionality, and will be removed in the next major release. ([#25107](https://github.com/hashicorp/terraform-provider-google/pull/25107))

FEATURES:
* **New Resource:** `google_ces_agent` ([#25106](https://github.com/hashicorp/terraform-provider-google/pull/25106))
* **New Resource:** `google_ces_guardrail` ([#25112](https://github.com/hashicorp/terraform-provider-google/pull/25112))
* **New Resource:** `google_ces_tool` ([#25113](https://github.com/hashicorp/terraform-provider-google/pull/25113))
* **New Resource:** `google_cloud_security_compliance_cloud_control` ([#25137](https://github.com/hashicorp/terraform-provider-google/pull/25137))
* **New Resource:** `google_cloud_security_compliance_framework_deployment` ([#25138](https://github.com/hashicorp/terraform-provider-google/pull/25138))
* **New Resource:** `google_cloud_security_compliance_framework` ([#25111](https://github.com/hashicorp/terraform-provider-google/pull/25111))
* **New Resource:** `google_discovery_engine_serving_config` ([#25105](https://github.com/hashicorp/terraform-provider-google/pull/25105))
* **New Resource:** `google_oracle_database_exascale_db_storage_vault` ([#25129](https://github.com/hashicorp/terraform-provider-google/pull/25129))

IMPROVEMENTS:
* apphub: added `functional_type`, `registration_type`, and `extended_metadata` fields to `google_apphub_service` and `google_apphub_workload` resources ([#25145](https://github.com/hashicorp/terraform-provider-google/pull/25145))
* ces: added `bearer_token_config` field to `google_ces_toolset` resource ([#25119](https://github.com/hashicorp/terraform-provider-google/pull/25119))
* ces: added `client_certificate_settings` field to `google_ces_app` resource ([#25117](https://github.com/hashicorp/terraform-provider-google/pull/25117))
* compute: added `block_names` field to `google_compute_reservation` resource ([#25121](https://github.com/hashicorp/terraform-provider-google/pull/25121))
* compute: added `sub_block_names` field to `google_compute_reservation_block` data source ([#25121](https://github.com/hashicorp/terraform-provider-google/pull/25121))
* compute: added `tls_settings` field to `google_compute_regional_backend_service` resource ([#25068](https://github.com/hashicorp/terraform-provider-google/pull/25068))
* container: added `end_time_behavior` field to `google_container_cluster` resource ([#25120](https://github.com/hashicorp/terraform-provider-google/pull/25120))
* container: added `writable_cgroups` field to `node_config.defaults.containerd_config` in `google_container_cluster` resource ([#25140](https://github.com/hashicorp/terraform-provider-google/pull/25140))
* dataplex: added `catalog_publishing_enabled` field to `data_profile_spec` in `google_dataplex_datascan` resource ([#25143](https://github.com/hashicorp/terraform-provider-google/pull/25143))
* dns: added `forwarding_config.target_name_servers.ipv6_address` argument to `google_dns_managed_zone` resource ([#25131](https://github.com/hashicorp/terraform-provider-google/pull/25131))
* gkeonprem: added `advanced_networking`, `multiple_network_interfaces_config` and `bgp_lb_config` fields to `google_gkeonprem_bare_metal_cluster` resource ([#25136](https://github.com/hashicorp/terraform-provider-google/pull/25136))
* managedkafka: added `broker_capacity_config` field to `google_managed_kafka_cluster` resource ([#25074](https://github.com/hashicorp/terraform-provider-google/pull/25074))
* networksecurity: added `endpoint_settings.jumbo_frames_enabled` field to `google_network_security_firewall_endpoint` resource ([#25073](https://github.com/hashicorp/terraform-provider-google/pull/25073))
* run: added `readiness_probe` field to `cloud_run_service` resource ([#25114](https://github.com/hashicorp/terraform-provider-google/pull/25114))

BUG FIXES:
* backupdr: updated `google_backup_dr_backup_plan_associations` and `google_backup_dr_data_source_references` to use LIST APIs, and require the correct List permissions ([#25107](https://github.com/hashicorp/terraform-provider-google/pull/25107))
* provider: an issue preventing X.509 certificates from being used for authentication when supplied as Application Default Credentials as been resolved ([#25144](https://github.com/hashicorp/terraform-provider-google/pull/25144))

## 7.11.0 (November 11, 2025)

DEPRECATIONS:
* `pubsublite`: `google_pubsub_lite_reservation` will be turned down effective March 18, 2026. Use `google_pubsub_reservation` instead. ([#25058](https://github.com/hashicorp/terraform-provider-google/pull/25058))
* `pubsublite`: `google_pubsub_lite_subscription` will be turned down effective March 18, 2026. Use `google_pubsub_subscription` instead. ([#25058](https://github.com/hashicorp/terraform-provider-google/pull/25058))
* `pubsublite`: `google_pubsub_lite_topic` will be turned down effective March 18, 2026. Use `google_pubsub_topic` instead. ([#25058](https://github.com/hashicorp/terraform-provider-google/pull/25058))

BREAKING CHANGES:
* netapp: made `google_netapp_volume.export_policy.rules.squash_mode` not preserve values returned by the API. Without this change, unsetting `squash_mode` in the provider can cause an API error. ([#25059](https://github.com/hashicorp/terraform-provider-google/pull/25059))

FEATURES:
* **New Data Source:** `google_artifact_registry_python_packages` ([#25053](https://github.com/hashicorp/terraform-provider-google/pull/25053))
* **New Data Source:** `google_cloud_identity_policy` ([#24946](https://github.com/hashicorp/terraform-provider-google/pull/24946))
* **New Data Source:** `google_compute_reservation_block` ([#25034](https://github.com/hashicorp/terraform-provider-google/pull/25034))
* **New Data Source:** `google_compute_reservation_sub_block` ([#25034](https://github.com/hashicorp/terraform-provider-google/pull/25034))
* **New Resource:** `google_ces_deployment` ([#24945](https://github.com/hashicorp/terraform-provider-google/pull/24945))
* **New Resource:** `google_ces_example` ([#25056](https://github.com/hashicorp/terraform-provider-google/pull/25056))
* **New Resource:** `google_discovery_engine_user_store` ([#25054](https://github.com/hashicorp/terraform-provider-google/pull/25054))

IMPROVEMENTS:
* bigquery: added `external_data_configuration.decimal_target_types` to `google_bigquery_table` ([#24936](https://github.com/hashicorp/terraform-provider-google/pull/24936))
* compute: added `internal_ipv6_prefix` field to the `google_compute_subnetwork` resource ([#25037](https://github.com/hashicorp/terraform-provider-google/pull/25037))
* compute: added `ipv6_access_type` field and `INTERNAL_IPV6_SUBNETWORK_CREATION` as a supported value for the `mode` field in `google_compute_public_delegated_prefix` resource ([#24940](https://github.com/hashicorp/terraform-provider-google/pull/24940))
* compute: added `ipv6_access_type` field to `google_compute_public_advertised_prefix` resource ([#24911](https://github.com/hashicorp/terraform-provider-google/pull/24911))
* dataplex: added `data_documentation_spec` field to `google_dataplex_datascan` resource to support the `DATA_DOCUMENTATION` scan type ([#25044](https://github.com/hashicorp/terraform-provider-google/pull/25044))
* dataproc: added `resource_manager_tags` to `google_dataproc_cluster` resource ([#25057](https://github.com/hashicorp/terraform-provider-google/pull/25057))
* lustre: added `placement_policy` field to `google_lustre_instance` resource ([#25042](https://github.com/hashicorp/terraform-provider-google/pull/25042))
* netapp: added `cache_parameters` field to `google_netapp_volume` resource ([#24909](https://github.com/hashicorp/terraform-provider-google/pull/24909))
* secretmanager: added project and short name support for `secret` on `google_secret_manager_secret_version` ([#25045](https://github.com/hashicorp/terraform-provider-google/pull/25045))
* secretmanager: added project and short name support for `secret` on ephemeral `google_secret_manager_secret_version` ([#25045](https://github.com/hashicorp/terraform-provider-google/pull/25045))

BUG FIXES:
* alloydb: fixed issue with creation when `initial_user.password` was set to a computed value in `google_alloydb_cluster` ([#25036](https://github.com/hashicorp/terraform-provider-google/pull/25036))
* bigquery: fixed extraneous diffs in `google_bigquery_table.external_data_configuration.schema` ([#24936](https://github.com/hashicorp/terraform-provider-google/pull/24936))
* compute: fixed a breaking change in `google_compute_instance` introduced in 7.9.0 where a destroy-diff is prompted for instances with preset GPUs ([#25021](https://github.com/hashicorp/terraform-provider-google/pull/25021))
* container: added `KUBE_DNS` as an accepted value for `cluster_dns` field on `google_container_cluster` ([#24953](https://github.com/hashicorp/terraform-provider-google/pull/24953))
* netapp: fixed bug where unsetting `export_policy.rules.squash_mode` on `google_netapp_volume` can cause an API error ([#25059](https://github.com/hashicorp/terraform-provider-google/pull/25059))
* pubsub: fixed bug where `google_pubsub_subscription` could only be updated if `bigquery_config` was modified ([#24952](https://github.com/hashicorp/terraform-provider-google/pull/24952))
* sql: fixed bug where `final_backup_description` in `google_sql_database_instance` resource wasn't set on the final backup on delete ([#25055](https://github.com/hashicorp/terraform-provider-google/pull/25055))
* storage: fixed bug where certain changes to `google_storage_bucket_acl.role_entity` were ignored ([#24949](https://github.com/hashicorp/terraform-provider-google/pull/24949))
* workstations: fixed bug in `google_workstations_workstation` where setting `source_workstation` caused a permadiff that forced recreation ([#24941](https://github.com/hashicorp/terraform-provider-google/pull/24941))
* vmwareengine: made deletion of `google_vmwareengine_private_cloud` wait until the deletion completes ([#25040](https://github.com/hashicorp/terraform-provider-google/pull/25040))

## 7.10.0 (November 4, 2025)

BREAKING CHANGES:
* alloydb: marked `initial_user.password` as required on create of new `google_alloydb_cluster` resources. This change aligns the provider with existing API constraints to surface errors earlier. ([#25022](https://github.com/hashicorp/terraform-provider-google/pull/25022))

FEATURES:
* **New Resource:** `google_ces_app` ([#24861](https://github.com/hashicorp/terraform-provider-google/pull/24861))
* **New Resource:** `google_ces_toolset` ([#24885](https://github.com/hashicorp/terraform-provider-google/pull/24885))
* **New Resource:** `google_discovery_engine_control` ([#24883](https://github.com/hashicorp/terraform-provider-google/pull/24883))
* **New Resource:** `google_netapp_host_group` ([#24876](https://github.com/hashicorp/terraform-provider-google/pull/24876))
* **New Resource:** `google_network_management_organization_vpc_flow_logs_config` ([#24896](https://github.com/hashicorp/terraform-provider-google/pull/24896))
* **New Resource:** `google_network_services_multicast_domain` ([#24864](https://github.com/hashicorp/terraform-provider-google/pull/24864))
* **New Resource:** `google_privileged_access_manager_settings` ([#24878](https://github.com/hashicorp/terraform-provider-google/pull/24878))
* **New Ephemeral Resource:** `google_client_config` ([#24900](https://github.com/hashicorp/terraform-provider-google/pull/24900))

IMPROVEMENTS:
* cloudfunctions2: added `direct_vpc_network_interface` and `direct_vpc_egress` field to `google_cloudfunctions2_function` resource ([#24895](https://github.com/hashicorp/terraform-provider-google/pull/24895))
* cloudrunv2: added `template.container.depends_on` field to `google_cloud_run_v2_worker_pool` resource ([#24893](https://github.com/hashicorp/terraform-provider-google/pull/24893))
* compute: added `grpc_tls_health_check` field to `google_compute_healthcheck` resource ([#24872](https://github.com/hashicorp/terraform-provider-google/pull/24872))
* container: added `network_tier_config` to `google_container_cluster` resource. ([#24877](https://github.com/hashicorp/terraform-provider-google/pull/24877))
* eventarc: added `labels` field to `google_eventarc_channel` resource ([#24854](https://github.com/hashicorp/terraform-provider-google/pull/24854))
* netapp: added `block_devices` field and `ISCSI` protocol support to `goolge_netapp_volume` resource, and increased timeouts on its operations ([#24898](https://github.com/hashicorp/terraform-provider-google/pull/24898))
* netapp: added `type` field to `google_netapp_storage_pool` resource ([#24867](https://github.com/hashicorp/terraform-provider-google/pull/24867))
* vertexai: added `psc_automation_configs` field to `google_vertex_ai_endpoint` resource ([#24870](https://github.com/hashicorp/terraform-provider-google/pull/24870))
* vertexai: added `sync_config.continuous` field to `google_vertex_ai_feature_online_store_featureview` ([#24881](https://github.com/hashicorp/terraform-provider-google/pull/24881))

BUG FIXES:
* accesscontextmanager: fixed issue where `google_access_context_manager_service_perimeter_[dry_run_][egress|ingress]_policy` caused the provider to crash when a provided identity casing was invalid. ([#24886](https://github.com/hashicorp/terraform-provider-google/pull/24886))
* apigee: fixed issue where `credentials` block was not populated in the Terraform state in `google_apigee_developer_app` resource ([#24880](https://github.com/hashicorp/terraform-provider-google/pull/24880))
* compute: fixed `google_compute_network_firewall_policy_rule` staying disabled after apply with `disabled = false` ([#24879](https://github.com/hashicorp/terraform-provider-google/pull/24879))
* compute: fixed a breaking change in `google_compute_instance` introduced in 7.9.0 where a destroy-diff is prompted for instances with preset GPUs ([#25020](https://github.com/hashicorp/terraform-provider-google/pull/25020)
* compute: resolve permadiff for `display_name` in new deployments of `google_compute_organization_security_policy` ([#24882](https://github.com/hashicorp/terraform-provider-google/pull/24882))
* storage: fixed a conversion error in `google_storage_bucket` state migration. This bug impacted Pulumi users. ([#24853](https://github.com/hashicorp/terraform-provider-google/pull/24853))

## 7.9.0 (October 28, 2025)

BREAKING CHANGES:
* beyondcorp: made the `ports` field in `endpoint_matchers` required in response to a change in the API surface. ([#24770](https://github.com/hashicorp/terraform-provider-google/pull/24770))

FEATURES:
* **New Resource:** `google_firestore_user_creds` ([#24794](https://github.com/hashicorp/terraform-provider-google/pull/24794))
* **New Resource:** `google_network_security_dns_threat_detector` ([#24744](https://github.com/hashicorp/terraform-provider-google/pull/24744))

IMPROVEMENTS:
* appengine: added `ssl_policy` to `application` on `google_app_engine_application` resource ([#24786](https://github.com/hashicorp/terraform-provider-google/pull/24786))
* bigquery: added support for IAM conditions in `google_bigquery_dataset_iam_*` ([#24778](https://github.com/hashicorp/terraform-provider-google/pull/24778))
* compute: promoted `policy_type` to GA in `google_compute_network_firewall_policy`, `google_compute_network_firewall_policy_with_rules`, `google_compute_region_network_firewall_policy`, `google_compute_region_network_firewall_policy_with_rules`. ([#24769](https://github.com/hashicorp/terraform-provider-google/pull/24769))
* container: added `dns_endpoint_confg.enable_k8s_tokens_via_dns` and  `dns_endpoint_config.enable_k8s_certs_via_dns` fields to `google_container_cluster` resource ([#24774](https://github.com/hashicorp/terraform-provider-google/pull/24774))
* container: added `fleet.membership_type` field to `google_container_cluster` resource ([#24759](https://github.com/hashicorp/terraform-provider-google/pull/24759))
* dataplex: added `data_classification` field to `google_dataplex_aspect_type` resource ([#24807](https://github.com/hashicorp/terraform-provider-google/pull/24807))
* iamworkforcepool: added `scim_usage` field to `workforce_pool_provider` resource ([#24787](https://github.com/hashicorp/terraform-provider-google/pull/24787))
* memorystore: added `available_maintenance_versions` field to `google_memorystore_instance` resource ([#24745](https://github.com/hashicorp/terraform-provider-google/pull/24745))
* memorystore: added `maintenance_version` field to `google_memorystore_instance` resource ([#24740](https://github.com/hashicorp/terraform-provider-google/pull/24740))
* redis: added `available_maintenance_versions` field to `google_redis_cluster` resource ([#24745](https://github.com/hashicorp/terraform-provider-google/pull/24745))
* redis: added `maintenance_version` field to `google_redis_cluster` resource ([#24740](https://github.com/hashicorp/terraform-provider-google/pull/24740))
* storagetransfer: added `transfer_manifest` field to `google_storage_transfer_job` resource ([#24768](https://github.com/hashicorp/terraform-provider-google/pull/24768))

BUG FIXES:
* bigquery: added validation for `target_types` in `google_bigquery_dataset_access` ([#24810](https://github.com/hashicorp/terraform-provider-google/pull/24810))
* cloudquotas: resolved permadiff for `preferred_value` in `google_cloud_quotas_quota_preference` ([#24776](https://github.com/hashicorp/terraform-provider-google/pull/24776))
* compute: fixed scenario where `google_compute_instance` would not be staged for recreation if `guest_accelerator.count` was updated to 0 from non-zero value ([#24762](https://github.com/hashicorp/terraform-provider-google/pull/24762))
* sql: fixed an issue where `dataDiskSize` was unintentionally null instead of set to the current value in API requests, triggering unrelated errors ([#24790](https://github.com/hashicorp/terraform-provider-google/pull/24790))

## 7.8.0 (October 21, 2025)

FEATURES:
* **New Data Source:** `google_artifact_registry_packages` ([#24696](https://github.com/hashicorp/terraform-provider-google/pull/24696))
* **New Data Source:** `google_network_management_connectivity_tests` ([#24635](https://github.com/hashicorp/terraform-provider-google/pull/24635))
* **New Resource:** `google_apigee_environment_api_revision_deployment` ([#24657](https://github.com/hashicorp/terraform-provider-google/pull/24657))
* **New Resource:** `google_dataplex_entry_link` ([#24737](https://github.com/hashicorp/terraform-provider-google/pull/24737))
* **New Resource:** `google_discovery_engine_assistant` ([#24724](https://github.com/hashicorp/terraform-provider-google/pull/24724))
* **New Resource:** `google_oracle_database_db_system` ([#24733](https://github.com/hashicorp/terraform-provider-google/pull/24733))
* **New Resource:** `google_saas_runtime_unit` ([#24692](https://github.com/hashicorp/terraform-provider-google/pull/24692))

IMPROVEMENTS:
* compute: added `IN_FLIGHT` to `balancing_mode` on `google_compute_backend_service` resource ([#24710](https://github.com/hashicorp/terraform-provider-google/pull/24710))
* compute: added new field `instance_lifecycle_policy.on_repair.allow_changing_zone` to `google_compute_region_instance_group_manager` & `google_compute_instance_group_manager` ([#24706](https://github.com/hashicorp/terraform-provider-google/pull/24706))
* compute: promoted `security_policy` in `compute_region_backend_service` resource to GA ([#24693](https://github.com/hashicorp/terraform-provider-google/pull/24693))
* compute: promoted the `google_compute_preview_feature` resource to GA. ([#24725](https://github.com/hashicorp/terraform-provider-google/pull/24725))
* compute: the `activation_status` attribute within the `google_compute_preview_feature` resource now uses the `ACTIVATION_STATE_UNSPECIFIED` value instead of `DISABLED`. Support for `DISABLED` will be added in a future release. ([#24725](https://github.com/hashicorp/terraform-provider-google/pull/24725))
* datastream: added `backfill_all.mongodb_excluded_objects` and  `source_config.mongodb_source_config` fields to `google_datastream_stream` ([#24727](https://github.com/hashicorp/terraform-provider-google/pull/24727))
* datastream: added `mongodb_profile` field to `google_datastream_connection_profile` ([#24727](https://github.com/hashicorp/terraform-provider-google/pull/24727))
* discoveryengine: added `connector_modes`, `sync_mode`, `incremental_refresh_interval`, `auto_run_disabled`, and `incremental_sync_disabled` fields to `google_discovery_engine_data_connector` resource ([#24658](https://github.com/hashicorp/terraform-provider-google/pull/24658))
* discoveryengine: added `kms_key_name` field to `google_discovery_engine_search_engine` resource ([#24658](https://github.com/hashicorp/terraform-provider-google/pull/24658))
* discoveryengine: added in-place update support for `entities.params` and `entities.key_property_mappings` in `google_discovery_engine_data_connector` ([#24739](https://github.com/hashicorp/terraform-provider-google/pull/24739))
* dlp: added `publish_findings_to_dataplex_catalog` field to `google_data_loss_prevention_job_trigger ` ([#24722](https://github.com/hashicorp/terraform-provider-google/pull/24722))
* iambeta: allowed GKE workload identity pool pattern in `workload_identity_pool_id` field of `google_iam_workload_identity_pool` resource. ([#24656](https://github.com/hashicorp/terraform-provider-google/pull/24656))
* memorystore: added `maintenance_version` field to `google_memorystore_instance` resource ([#24740](https://github.com/hashicorp/terraform-provider-google/pull/24740))
* memorystore: added `available_maintenance_versions` field to `google_memorystore_instance` resource ([#24745](https://github.com/hashicorp/terraform-provider-google/pull/24745))
* networkconnectivity: added `HYBRID_INSPECTION` enum value to `preset_topology` field in `google_network_connectivity_hub` resource ([#24738](https://github.com/hashicorp/terraform-provider-google/pull/24738))
* networkservices: added `isolationConfig` on `google_network_services_service_lb_policies` resource ([#24652](https://github.com/hashicorp/terraform-provider-google/pull/24652))
* redis: added `deletion_protection` field to `redis_instance` to make deleting them require an explicit intent. `redis_instance` resources now cannot be destroyed unless `deletion_protection = false` is set for the resource. ([#24654](https://github.com/hashicorp/terraform-provider-google/pull/24654))
* redis: added `maintenance_version` field to `google_redis_cluster` resource ([#24740](https://github.com/hashicorp/terraform-provider-google/pull/24740))
* redis: added `available_maintenance_versions` field to `google_redis_cluster` resource ([#24745](https://github.com/hashicorp/terraform-provider-google/pull/24745))
* saas_runtime: added `default_release` field to `google_saas_runtime_unit_kind` resource ([#24726](https://github.com/hashicorp/terraform-provider-google/pull/24726))
* sql: added `read_pool_auto_scale_config` support to `sql_database_instance` resource ([#24723](https://github.com/hashicorp/terraform-provider-google/pull/24723))

BUG FIXES:
* bigquery: fixed the issue where `google_bigquery_table` detected an incorrect `schema` diff on tables with row access policies when the schema was unchanged. ([#24711](https://github.com/hashicorp/terraform-provider-google/pull/24711))
* compute: allow `requested_link_count` to be updated in-place in `google_compute_interconnect` resource ([#24705](https://github.com/hashicorp/terraform-provider-google/pull/24705))

## 7.7.0 (October 14, 2025)
BREAKING CHANGES:
* discoveryengine: changed type of `google_discovery_engine_data_connector.entities.params`. Previously, it was a map of string keys to string values; now, it must be a [JSON-encoded](https://developer.hashicorp.com/terraform/language/functions/jsonencode) string containing an object. This change is being made in a minor release because the field wasn't usable as intended – specifically, all current valid uses require mapping strings to _lists_ of strings. ([#24658](https://github.com/hashicorp/terraform-provider-google/pull/24658))

FEATURES:
* **New Data Source:** `google_network_management_connectivity_tests` ([#24635](https://github.com/hashicorp/terraform-provider-google/pull/24635))
* **New Resource:** `google_apigee_developer_app` ([#24625](https://github.com/hashicorp/terraform-provider-google/pull/24625))
* **New Resource:** `google_discovery_engine_license_config` ([#24619](https://github.com/hashicorp/terraform-provider-google/pull/24619))
* **New Resource:** `google_iam_workforce_pool_provider_scim_tenant` ([#24587](https://github.com/hashicorp/terraform-provider-google/pull/24587))
* **New Resource:** `google_kms_project_kaj_policy_config` ([#24622](https://github.com/hashicorp/terraform-provider-google/pull/24622))
* **New Resource:** `google_saas_runtime_tenant` ([#24608](https://github.com/hashicorp/terraform-provider-google/pull/24608))

IMPROVEMENTS:
* apigee: updated the `scopes` argument in `google_apigee_api_product` resource to be order-insensitive. ([#24625](https://github.com/hashicorp/terraform-provider-google/pull/24625))
* beyondcorp: added `proxy_protocol_config` and `service_discovery` fields to `google_beyondcorp_security_gateway` resource ([#24609](https://github.com/hashicorp/terraform-provider-google/pull/24609))
* cloudrunv2: added `default_uri_disabled` field to `google_cloud_run_v2_service` resource. (GA promotion) ([#24602](https://github.com/hashicorp/terraform-provider-google/pull/24602))
* cloudrunv2: added `health_check_disabled` field to `google_cloud_run_v2_service` resource. ([#24602](https://github.com/hashicorp/terraform-provider-google/pull/24602))
* compute: added `params` field to `google_compute_router` resource (GA) ([#24611](https://github.com/hashicorp/terraform-provider-google/pull/24611))
* discoveryengine: added `connector_modes`, `sync_mode`, `incremental_refresh_interval`, `auto_run_disabled`, and `incremental_sync_disabled` fields to `google_discovery_engine_data_connector` resource ([#24658](https://github.com/hashicorp/terraform-provider-google/pull/24658))
* discoveryengine: added `kms_key_name` field to `google_discovery_engine_search_engine` resource ([#24658](https://github.com/hashicorp/terraform-provider-google/pull/24658))
* dlp: added `publish_to_dataplex_catalog` field to `discovery_config` resource ([#24621](https://github.com/hashicorp/terraform-provider-google/pull/24621))
* gkeonprem: made it possible to set the `on_prem_version` field on `google_gkeonprem_vmware_node_pool` (previously output-only) ([#24614](https://github.com/hashicorp/terraform-provider-google/pull/24614))
* memcache: added `deletion_protection` field to `memcache_instance` to make deleting them require an explicit intent. `memcache_instance` resources now cannot be destroyed unless `deletion_protection = false` is set for the resource. ([#24613](https://github.com/hashicorp/terraform-provider-google/pull/24613))
* metastore: added `tags` field to `google_dataproc_metastore_service` and 'google_dataproc_metastore_federation' resources to allow setting tags for services and federation at creation time ([#24633](https://github.com/hashicorp/terraform-provider-google/pull/24633))
* networksecurity: added `URL_FILTERING` option to enum field `type` for `google_network_security_security_profile` resource ([#24583](https://github.com/hashicorp/terraform-provider-google/pull/24583))
* networksecurity: added `url_filtering_profile` field  to `google_network_security_security_profile_group` resource (beta) ([#24583](https://github.com/hashicorp/terraform-provider-google/pull/24583))
* networksecurity: added `url_filtering_profile` field  to `google_network_security_security_profile` resource (beta) ([#24583](https://github.com/hashicorp/terraform-provider-google/pull/24583))
* sql: added `source_instance_deletion_time` field to `google_sql_database_instance_latest_recovery_time` data source ([#24576](https://github.com/hashicorp/terraform-provider-google/pull/24576))
* sql: added `source_instance_deletion_time` field to `google_sql_database_instance` resource ([#24576](https://github.com/hashicorp/terraform-provider-google/pull/24576))

BUG FIXES:
* bigqueryanalyticshub: fixed `google_bigquery_analytics_hub_listing_subscription` import ([#24634](https://github.com/hashicorp/terraform-provider-google/pull/24634))
* discoveryengine: fixed bug where it wasn't possible to specify values for `knowledgeBaseSysId` or `catalogSysId` in `google_discovery_engine_data_connector.entities.params`. ([#24658](https://github.com/hashicorp/terraform-provider-google/pull/24658))

## 7.6.0 (October 7, 2025)
DEPRECATIONS:
* networksecurity: deprecated `ignore_case`, `exact`, `prefix`, `suffix` and `contains` fields in `http_rules.from.not_sources.principals` and `http_rules.from.sources.principals` blocks in `google_network_security_authz_policy` resource. Use the equivalent fields in `http_rules.from.not_sources.principals.principal` or `http_rules.from.sources.principals.principal` instead. ([#24543](https://github.com/hashicorp/terraform-provider-google/pull/24543))

BREAKING CHANGES:
* container: `node_config` blocks that had set `kubelet_config` without explicitly setting `cpu_cfs_quota` implicitly set `cfu_cfs_quota` to `false` when unset. From this version onwards, an unset `cpu_cfs_quota` will instead match the API default of true `true`. Resources that are recreated will receive the new value; old resources are unaffected, and may change values by explicitly setting the intended one. ([#24569](https://github.com/hashicorp/terraform-provider-google/pull/24569))
* storageinsights: removed `activity_data_retention_period_days` field from `google_storage_insights_dataset_config` resource due to a delayed launch. It will be readded when the feature launches. ([#24570](https://github.com/hashicorp/terraform-provider-google/pull/24570))

FEATURES:
* **New Resource:** `google_kms_folder_kaj_policy_config` ([#24513](https://github.com/hashicorp/terraform-provider-google/pull/24513))
* **New Resource:** `google_vertex_ai_cache_config` ([#24541](https://github.com/hashicorp/terraform-provider-google/pull/24541))
* **New Resource:** `google_vertex_ai_reasoning_engine` ([#24512](https://github.com/hashicorp/terraform-provider-google/pull/24512))

IMPROVEMENTS:
* backupdr: added `data_source` and `rules_config_info` fields to `google_backup_dr_backup_plan_associations` datasource ([#24517](https://github.com/hashicorp/terraform-provider-google/pull/24517))
* beyondcorp: added `external`, `proxy_protocol`, and `schema` fields to `google_beyondcorp_security_gateway_application` resource ([#24542](https://github.com/hashicorp/terraform-provider-google/pull/24542))
* beyondcorp: changed `endpoint_matchers` field to not be required anymore in the `google_beyondcorp_security_gateway_application` resource ([#24542](https://github.com/hashicorp/terraform-provider-google/pull/24542))
* cloudrunv2: added `default_uri_disabled` field to `google_cloud_run_v2_service` resource ([#24556](https://github.com/hashicorp/terraform-provider-google/pull/24556))
* compute: added `shared_secret_wo` and `shared_secret_wo_version` fields to `google_compute_vpn_tunnel` resource, enabling write-only management of the shared secret. ([#24491](https://github.com/hashicorp/terraform-provider-google/pull/24491))
* dlp: added `SENSITIVITY_UNKNOWN` as possible enum value for `actions.tag_resources.tag_conditions.sensitivity_score.score` in `google_data_loss_prevention_discovery_config` resource ([#24564](https://github.com/hashicorp/terraform-provider-google/pull/24564))
* dlp: added `actions.save_findings.output_config.storage_path` field to `google_data_loss_prevention_job_trigger` resource ([#24558](https://github.com/hashicorp/terraform-provider-google/pull/24558))
* filestore: added `file_shares.nfs_export_options.network` and `networks.psc_config.endpoint_project` fields to `google_filestore_instance` resource ([#24567](https://github.com/hashicorp/terraform-provider-google/pull/24567))
* lustre: increased creation timeout from 20min to 40min for `google_lustre_instance` resource ([#24559](https://github.com/hashicorp/terraform-provider-google/pull/24559))
* netapp: added `hybrid_replication_user_commands` field  with subfield `commands` to `google_netapp_volume_replication` resource ([#24554](https://github.com/hashicorp/terraform-provider-google/pull/24554))
* netapp: added `replication_schedule`, `hybrid_replication_type`, `large_volume_constituent_count` fields to `hybrid_replication_parameters` field in `google_netapp_volume` resource ([#24554](https://github.com/hashicorp/terraform-provider-google/pull/24554))
* networksecurity: added `ip_blocks` field to `google_network_security_authz_policy` resource ([#24543](https://github.com/hashicorp/terraform-provider-google/pull/24543))
* secretmanager: added ephemeral support for `google_secret_manager_secret_version` resource ([#24566](https://github.com/hashicorp/terraform-provider-google/pull/24566))
* sql: added `source_instance_deletion_time` field to `google_sql_database_instance_latest_recovery_time` data source ([#24576](https://github.com/hashicorp/terraform-provider-google/pull/24576))
* sql: added `source_instance_deletion_time` field to `google_sql_database_instance` resource ([#24576](https://github.com/hashicorp/terraform-provider-google/pull/24576))
* storagetransfer: added `user_project_override` and  `billing_project` fields to `google_storage_transfer_job` resource ([#24504](https://github.com/hashicorp/terraform-provider-google/pull/24504))

BUG FIXES:
* container: fixed the default for `node_config.kubelet_config.cpu_cfs_quota` on `google_container_cluster`, `google_container_node_pool`, `google_container_cluster.node_pool` to align with the API. Terraform will now send a `true` value when the field is unset on creation, and preserve any previously set value when unset. Explicitly set values will work as defined in configuration. ([#24569](https://github.com/hashicorp/terraform-provider-google/pull/24569))


## 7.5.0 (September 30, 2025)

BREAKING CHANGES:
* netapp: changed `peer_ip_addresses` field type from String to Array in `google_netapp_volume` resource, as it was unusable otherwise ([#24428](https://github.com/hashicorp/terraform-provider-google/pull/24428))

FEATURES:
* **New Data Source:** `google_artifact_registry_maven_artifacts` ([#24487](https://github.com/hashicorp/terraform-provider-google/pull/24487))
* **New Data Source:** `google_artifact_registry_npm_packages` ([#24486](https://github.com/hashicorp/terraform-provider-google/pull/24486))
* **New Resource:** `google_apigee_api_deployment` ([#24469](https://github.com/hashicorp/terraform-provider-google/pull/24469))
* **New Resource:** `google_discovery_engine_data_connector` ([#24472](https://github.com/hashicorp/terraform-provider-google/pull/24472))
* **New Resource:** `google_managed_kafka_connect_cluster` ([#24443](https://github.com/hashicorp/terraform-provider-google/pull/24443))
* **New Resource:** `google_managed_kafka_connector` ([#24443](https://github.com/hashicorp/terraform-provider-google/pull/24443))
* **New Resource:** `google_kms_organization_kaj_policy_config` ([#24471](https://github.com/hashicorp/terraform-provider-google/pull/24471))
* **New Resource:** `google_saas_runtime_rollout_kind` ([#24447](https://github.com/hashicorp/terraform-provider-google/pull/24447))

IMPROVEMENTS:
* cloudrunv2: added `mount_options` in gcsfuse volumes for `google_cloud_run_v2_service`, `google_cloud_run_v2_job`, and `google_cloud_run_v2_workerpool` resources. ([#24413](https://github.com/hashicorp/terraform-provider-google/pull/24413))
* cloudrunv2: added `startup_probe` and `liveness_probe` to `google_cloud_run_v2_worker_pool` resource ([#24418](https://github.com/hashicorp/terraform-provider-google/pull/24418))
* compute: added `bandwidth_allocation` field to `google_compute_wire_group` resource ([#24460](https://github.com/hashicorp/terraform-provider-google/pull/24460))
* compute: added `shared_secret_wo` and `shared_secret_wo_version` fields for `google_compute_vpn_tunnel` resource, enabling write-only management of the shared secret. ([#24491](https://github.com/hashicorp/terraform-provider-google/pull/24491))
* dialogflow:  added `new_recognition_result_notification_config` field to `google_dialogflow_conversation_profile ` resource ([#24468](https://github.com/hashicorp/terraform-provider-google/pull/24468))
* discoveryengine: added `features` field to `google_discovery_engine_search_engine` resource ([#24445](https://github.com/hashicorp/terraform-provider-google/pull/24445))
* dlp: added `other_cloud_target` and `other_cloud_starting_location` to `google_data_loss_prevention_discovery_config` ([#24463](https://github.com/hashicorp/terraform-provider-google/pull/24463))
* gkebackup: added `backup_config.selected_namespace_labels` field to `google_gke_backup_backup_plan` resource ([#24427](https://github.com/hashicorp/terraform-provider-google/pull/24427))
* looker: added `gemini_enabled` field to `google_looker_instance` resource ([#24461](https://github.com/hashicorp/terraform-provider-google/pull/24461))
* netapp: added `hot_tier_bypass_mode_enabled` and `hot_tier_size_used_gib` fields to `google_netapp_volume` ([#24454](https://github.com/hashicorp/terraform-provider-google/pull/24454))
* netapp: added `hot_tier_size_gib`, `enable_hot_tier_auto_resize`, `cold_tier_size_used_gib` and `hot_tier_size_used_gib` fields to `google_netapp_storage_pool` ([#24454](https://github.com/hashicorp/terraform-provider-google/pull/24454))
* oracledatabase: added `gcp_oracle_zone` field to `google_oracle_database_odb_network` resource ([#24456](https://github.com/hashicorp/terraform-provider-google/pull/24456))
* privilegedaccessmanager: added `approval_workflow.steps.id` field to `google_privileged_access_manager_entitlement` resource ([#24419](https://github.com/hashicorp/terraform-provider-google/pull/24419))
* pubsub: added support for `tags` field to `google_pubsub_topic` and `google_pubsub_subscription` resources ([#24442](https://github.com/hashicorp/terraform-provider-google/pull/24442))
* sql: added `point_in_time_restore_context` field to `google_sql_database_instance` ([#24489](https://github.com/hashicorp/terraform-provider-google/pull/24489))
* storage: added `force_destroy` field to `google_storage_insights_report_config` resource ([#24462](https://github.com/hashicorp/terraform-provider-google/pull/24462))
* storageinsights: added `activity_data_retention_period_days` field to `google_storage_insights_dataset_config` resource ([#24459](https://github.com/hashicorp/terraform-provider-google/pull/24459))
* vertexai: added `endpoint_config.private_service_connect_config` block to `google_vertex_ai_endpoint_with_model_garden_deployment` resource ([#24425](https://github.com/hashicorp/terraform-provider-google/pull/24425))
* vertexai: added `encryption_spec.kms_key_name` field to `google_vertex_ai_index_endpoint ` resource ([#24490](https://github.com/hashicorp/terraform-provider-google/pull/24490))
* vertexai: added `encryption_spec.kms_key_name` field to `google_vertex_ai_index` resource ([#24441](https://github.com/hashicorp/terraform-provider-google/pull/24441))

BUG FIXES:
* apihub: fixed a permadiff on `config_template` in `google_apihub_plugin` resource ([#24429](https://github.com/hashicorp/terraform-provider-google/pull/24429))
* storage: fixed a panic caused by empty `cors` blocks `google_storage_bucket` resource ([#24476](https://github.com/hashicorp/terraform-provider-google/pull/24476))

## 7.4.0 (September 23, 2025)

DEPRECATIONS:
* compute: deprecated the option to deploy a container during VM creation using the container startup agent in `google_compute_instance`. Use alternative services to run containers on your VMs. Learn more at https://cloud.google.com/compute/docs/containers/migrate-containers. ([#24375](https://github.com/hashicorp/terraform-provider-google/pull/24375))

FEATURES:
* **New Data Source:** `google_artifact_registry_maven_artifact` ([#24358](https://github.com/hashicorp/terraform-provider-google/pull/24358))
* **New Data Source:** `google_compute_interconnect_location` ([#24377](https://github.com/hashicorp/terraform-provider-google/pull/24377))
* **New Resource:** `google_network_services_wasm_plugin` ([#24406](https://github.com/hashicorp/terraform-provider-google/pull/24406))
* **New Resource:** `google_resource_manager_capability` ([#24404](https://github.com/hashicorp/terraform-provider-google/pull/24404))

IMPROVEMENTS:
* cloudrunv2: added `mount_options` in gcsfuse volumes for `google_cloud_run_v2_service`, `google_cloud_run_v2_job`, and `google_cloud_run_v2_workerpool` resources. ([#24413](https://github.com/hashicorp/terraform-provider-google/pull/24413))
* compute: added `cipher_suite` field to `google_compute_vpn_tunnel` resource. ([#24378](https://github.com/hashicorp/terraform-provider-google/pull/24378))
* container: added `auto_ipam_config` to `google_container_cluster` resource. ([#24396](https://github.com/hashicorp/terraform-provider-google/pull/24396))
* storage: added support for `timeouts`  to  `google_storage_bucket_iam_binding`, `google_storage_bucket_iam_member`, `google_storage_bucket_iam_policy` resources ([#24376](https://github.com/hashicorp/terraform-provider-google/pull/24376))

BUG FIXES:
* bigtable: fixed `node_scaling_factor` forcing new instance on `google_bigtable_instance` when adding new cluster ([#24410](https://github.com/hashicorp/terraform-provider-google/pull/24410))
* cloudscheduler: fixed a type assertion panic in `google_cloud_scheduler_job` when processing HTTP headers with nil or unexpected data types ([#24360](https://github.com/hashicorp/terraform-provider-google/pull/24360))
* compute: fixed the `Network field cannot be modified` issue in `google_compute_region_backend_service`. Now updating the `network` field will force the resource to be recreated. ([#24398](https://github.com/hashicorp/terraform-provider-google/pull/24398))
* netapp: fixed incorrect default value handling in `google_netapp_volume` for `export_policy.rules` attributes `has_root_access` and `squash_mode`. When not specified, these fields will now take on the API default value with no diff. ([#24395](https://github.com/hashicorp/terraform-provider-google/pull/24395))
* netapp: updated `google_netapp_storage_pool` to source the default value for the `qos_type` field from the API. If not specified in the configuration, `qos_type` will now default to the value provided by the NetApp Volumes API. ([#24394](https://github.com/hashicorp/terraform-provider-google/pull/24394))
* sql: fixed the permadiffs on `disk_size ` when `disk_autoresize` is enabled in `google_sql_database_instance` ([#24399](https://github.com/hashicorp/terraform-provider-google/pull/24399))
* workbench: added retry for `unable to queue the operation` 409 errors in `google_workbench_instance` resource. ([#24392](https://github.com/hashicorp/terraform-provider-google/pull/24392))

## 7.3.0 (September 16, 2025)

FEATURES:
* **New Data Source:** `google_backup_dr_data_source_reference` ([#24346](https://github.com/hashicorp/terraform-provider-google/pull/24346))
* **New Resource:** `google_bigquery_datapolicyv2_data_policy` ([#24313](https://github.com/hashicorp/terraform-provider-google/pull/24313))
* **New Resource:** `google_saas_runtime_release` ([#24289](https://github.com/hashicorp/terraform-provider-google/pull/24289))
* **New Resource:** `google_secure_source_manager_hook` ([#24345](https://github.com/hashicorp/terraform-provider-google/pull/24345))

IMPROVEMENTS:
* cloudrun: added `sub_path` field to `google_cloud_run_service` resource. ([#24341](https://github.com/hashicorp/terraform-provider-google/pull/24341))
* cloudrunv2: added `sub_path` field to `google_cloud_run_v2_service` `google_cloud_run_v2_job` and `google_cloud_run_v2_worker_pool` resource. ([#24341](https://github.com/hashicorp/terraform-provider-google/pull/24341))
* compute: added `labels` and `label_fingerprint` fields to `google_compute_security_policy` resource ([#24322](https://github.com/hashicorp/terraform-provider-google/pull/24322))
* compute: `labels` under `initialize_params` are now updatable on `google_compute_instance` ([#24349](https://github.com/hashicorp/terraform-provider-google/pull/24349))
* container: added new fields `memory_manager` and `topology_manager` to `node_kubelet_config` block ([#24277](https://github.com/hashicorp/terraform-provider-google/pull/24277))
* datastream: added `destination_config.bigquery_destination_config.source_hierarchy_datasets.project_id` field to `google_datastream_stream` resource ([#24340](https://github.com/hashicorp/terraform-provider-google/pull/24340))
* discoveryengine: added `app_type` field to `google_discovery_engine_search_engine` resource ([#24320](https://github.com/hashicorp/terraform-provider-google/pull/24320))
* gkeonprem: added `proxy` field to `google_gkeonprem_vmware_admin_cluster` resource ([#24338](https://github.com/hashicorp/terraform-provider-google/pull/24338))
* healthcare: added `validation_config` to `google_healthcare_fhir_store` resource ([#24336](https://github.com/hashicorp/terraform-provider-google/pull/24336))
* iamworkforcepool: added `extended_attributes` field to `workforce_pool_provider` resource ([#24308](https://github.com/hashicorp/terraform-provider-google/pull/24308))
* netapp: added `export_policy.rules.squash_mode` field to `google_netapp_volume` resource. ([#24350](https://github.com/hashicorp/terraform-provider-google/pull/24350))
* privateca: added `encryption_spec` field to `google_privateca_ca_pool` resource ([#24328](https://github.com/hashicorp/terraform-provider-google/pull/24328))
* run: added `connector`  to `vpc_access` on `google_cloud_run_v2_worker_pool` resource ([#24337](https://github.com/hashicorp/terraform-provider-google/pull/24337))
* tags: added the `DATA_GOVERNANCE` value to `google_tags_tag_key.purpose` ([#24307](https://github.com/hashicorp/terraform-provider-google/pull/24307))

BUG FIXES:
* bigquery: updated the schema change detection for `google_bigquery_table` to take into account presence of row access policy ([#24284](https://github.com/hashicorp/terraform-provider-google/pull/24284))
* compute: fixed `allow_global_access` to correctly be immutable for `google_compute_forwarding_rule` resources with load balancing scheme of INTERNAL_MANAGED ([#24312](https://github.com/hashicorp/terraform-provider-google/pull/24312))
* compute: fixed a crash in `google_compute_security_policy` due to a changed API response for empty `match.0.expr_options` blocks ([#24353](https://github.com/hashicorp/terraform-provider-google/pull/24353))
* dialogflow: added support for non-global endpoints for `google_dialogflow_conversation_profile` ([#24351](https://github.com/hashicorp/terraform-provider-google/pull/24351))
* publicca: use `RawURLEncoding` instead of `URLEncoding` for unpadded base64 encoding ([#24283](https://github.com/hashicorp/terraform-provider-google/pull/24283))
* secretmanager: fixed a panic in `google_secret_manager_secret_version` in a `secret_manager` ([#24326](https://github.com/hashicorp/terraform-provider-google/pull/24326))
* workbench: fixed issue that resource creation with computed `labels` field fails in `google_workbench_instance` resource ([#24311](https://github.com/hashicorp/terraform-provider-google/pull/24311))
* workbench: made `report-notebook-metrics` metadata key settable for `google_workbench_instance` ([#24310](https://github.com/hashicorp/terraform-provider-google/pull/24310))

## 7.2.0 (September 9, 2025)

FEATURES:
* **New Data Source:** `google_artifact_registry_python_package` ([#24267](https://github.com/hashicorp/terraform-provider-google/pull/24267))
* **New Data Source:** `google_backup_dr_data_source_references` ([#24268](https://github.com/hashicorp/terraform-provider-google/pull/24268))
* **New Resource:** `google_discovery_engine_acl_config` ([#24276](https://github.com/hashicorp/terraform-provider-google/pull/24276))
* **New Resource:** `google_saas_runtime_unit_kind` ([#24236](https://github.com/hashicorp/terraform-provider-google/pull/24236))

IMPROVEMENTS:
* chronicle: made the `scope_info` field in `google_chronicle_reference_list` configurable ([#24250](https://github.com/hashicorp/terraform-provider-google/pull/24250))
* compute: added `header_action` to `path_matcher` and `default_service` level on `google_compute_region_url_map` resource ([#24253](https://github.com/hashicorp/terraform-provider-google/pull/24253))
* container: added `secret_manager_config.rotation_config` field to `google_container_cluster` resource ([#24244](https://github.com/hashicorp/terraform-provider-google/pull/24244))
* container: added new fields `memory_manager` and `topology_manager` to `google_container_cluster.node_config.kubelet_config` and `google_container_node_pool.node_config.kubelet_config`  ([#24277](https://github.com/hashicorp/terraform-provider-google/pull/24277))
* sql: added `final_backup_description` and `final_backup_config` fields to `google_sql_database_instance` resource ([#24273](https://github.com/hashicorp/terraform-provider-google/pull/24273))
* storage: added `aws_s3_compatible_data_source` to `google_storage_transfer_job` resource ([#24241](https://github.com/hashicorp/terraform-provider-google/pull/24241))

BUG FIXES:
* provider: fixed an issue with `universe_domain` where the provider tried to connect to "googleapis.com" for user email logging when `universe_domain` was set ([#24238](https://github.com/hashicorp/terraform-provider-google/pull/24238))
* container: fixed a faulty diff for arrays on `user_managed_keys_config` that caused faulty cluster updates to be triggered in `google_container_cluster` ([#24256](https://github.com/hashicorp/terraform-provider-google/pull/24256))
* osconfig: fixed a permadiff in `google_osconfig_patch_deployment` where `patch_config.yum.minimal` doesn't send `false` for empty values ([#24247](https://github.com/hashicorp/terraform-provider-google/pull/24247))

## 7.1.1 (September 4, 2025)

BUG FIXES:
* bigtable: fixed an error encountered when applying `google_bigtable_table_iam_*` resources after upgrading to 7.x and replacing `instance` with `instance_name` ([#24255](https://github.com/hashicorp/terraform-provider-google/pull/24255))

## 7.1.0 (September 2, 2025)

DEPRECATIONS:
* container: deprecated `enterprise_config` field in `google_container_cluster` resource. GKE Enterprise features are now available without an Enterprise tier. ([#24210](https://github.com/hashicorp/terraform-provider-google/pull/24210))
* storage: removed deprecated status for field to `detect_md5hash`  in  `google_storage_bucket_object` resource ([#24147](https://github.com/hashicorp/terraform-provider-google/pull/24147))

FEATURES:
* **New Data Source:** `google_iap_web_forwarding_rule_service_iam_policy` ([#24178](https://github.com/hashicorp/terraform-provider-google/pull/24178))
* **New Resource:** `google_iap_web_forwarding_rule_service_iam_binding` ([#24178](https://github.com/hashicorp/terraform-provider-google/pull/24178))
* **New Resource:** `google_iap_web_forwarding_rule_service_iam_member` ([#24178](https://github.com/hashicorp/terraform-provider-google/pull/24178))
* **New Resource:** `google_iap_web_forwarding_rule_service_iam_policy` ([#24178](https://github.com/hashicorp/terraform-provider-google/pull/24178))

IMPROVEMENTS:
* artifactregistry: added `registry_uri` as attribute to `google_artifact_registry_repository` ([#24164](https://github.com/hashicorp/terraform-provider-google/pull/24164))
* backupdr: added 'supported_resource_types' field to `google_backup_dr_backup_plan` resource ([#24189](https://github.com/hashicorp/terraform-provider-google/pull/24189))
* backupdr: added `create_time` field to `google_backup_dr_backup` data source ([#24183](https://github.com/hashicorp/terraform-provider-google/pull/24183))
* cloudbuild: added `worker_config.enable_nested_virtualization` field to `google_cloudbuild_worker_pool` resource ([#24176](https://github.com/hashicorp/terraform-provider-google/pull/24176))
* cloudrunv2: added support for `multi_region_settings` field to `google_cloud_run_v2_service` resource ([#24149](https://github.com/hashicorp/terraform-provider-google/pull/24149))
* compute: add `params.resource_manager_tags` field to the `google_compute_region_backend_service` ([#24191](https://github.com/hashicorp/terraform-provider-google/pull/24191))
* compute: added `public_delegated_sub_prefixs` field to resource `google_compute_public_delegated_prefix` ([#24202](https://github.com/hashicorp/terraform-provider-google/pull/24202))
* compute: added `update_strategy`  field to `google_compute_network_peering ` resource ([#24180](https://github.com/hashicorp/terraform-provider-google/pull/24180))
* firestore: added `unique` field to `google_firestore_index` resource ([#24163](https://github.com/hashicorp/terraform-provider-google/pull/24163))
* netapp: added  `qos_type` and `available_throughput_mibps` fields to `google_netapp_storage_pool` resource ([#24161](https://github.com/hashicorp/terraform-provider-google/pull/24161))
* netapp: added  `throughput_mibps` field to `google_netapp_volume` resource ([#24161](https://github.com/hashicorp/terraform-provider-google/pull/24161))
* networkservices: allowed `EXPLICIT_ROUTING_MODE` for `routing_mode` on `google_network_services_gateway` resource ([#24151](https://github.com/hashicorp/terraform-provider-google/pull/24151))
* sql: added `consumer_network_status`, `ip_address`, and `status` fields to `psc_auto_connections` field on `google_sql_database_instance` resource ([#24201](https://github.com/hashicorp/terraform-provider-google/pull/24201))
* storagetransfer: added `service_account` field to `google_storage_transfer_job` resource ([#24193](https://github.com/hashicorp/terraform-provider-google/pull/24193))
* storagetransfer: added `transfer_spec.aws_s3_data_source.credentials_secret` to `google_storage_transfer_job` resource ([#24152](https://github.com/hashicorp/terraform-provider-google/pull/24152))

BUG FIXES:
* compute: fixed certain spurious diffs for `google_compute_region_backend_service.backend.group` ([#24157](https://github.com/hashicorp/terraform-provider-google/pull/24157))
* compute: fixed permadiff on `google_compute_region_network_endpoint_group` when no `network` is specified ([#24182](https://github.com/hashicorp/terraform-provider-google/pull/24182))
* memorystore: fixed permadiffs that cause destroy+recreate on new `google_memorystore_instance` when `desired_psc_auto_connections` is set ([#24212](https://github.com/hashicorp/terraform-provider-google/pull/24212))
* netapp: fixed a permadiff on `total_iops` in `google_netapp_storage_pool` resource ([#24207](https://github.com/hashicorp/terraform-provider-google/pull/24207))
* oracledatabase: fixed permadiffs on `google_oracle_database_autonomous_database` resource for the `odb_network` and `odb_subnet` fields ([#24184](https://github.com/hashicorp/terraform-provider-google/pull/24184))

## 7.0.1 (August 27, 2025)

BUG FIXES:
* storage: fixed a conversion crash in `google_storage_bucket` state migration [#24186](https://github.com/hashicorp/terraform-provider-google/pull/24186)

## 7.0.0 (August 26, 2025)

[Terraform Google Provider 7.0.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/version_7_upgrade)

BREAKING RESOURCE REMOVALS:
* beyondcorp: removed `google_beyondcorp_application`, its associated IAM resources `google_beyondcorp_application_iam_binding`, `google_beyondcorp_application_iam_member`, and `google_beyondcorp_application_iam_policy`, and the `google_beyondcorp_application_iam_policy` datasource. Use `google_beyondcorp_security_gateway_application` instead. [#23999](https://github.com/hashicorp/terraform-provider-google/pull/23999)
* notebooks: removed `google_notebooks_location` [#23607](https://github.com/hashicorp/terraform-provider-google/pull/23607)
* tpu: removed `google_tpu_node`. Use `google_tpu_v2_vm` instead. [#23964](https://github.com/hashicorp/terraform-provider-google/pull/23964)

BREAKING FIELD REMOVALS:
* cloudrunv2: removed `template.containers.depends_on` within `resource google_cloud_run_v2_worker_pool` [#23815](https://github.com/hashicorp/terraform-provider-google/pull/23815)
* colab: removed `post_startup_script_config` field from from `google_colab_runtime_template` resource [#24026](https://github.com/hashicorp/terraform-provider-google/pull/24026)
* compute: removed field `enable_flow_logs` from `google_compute_subnetwork` [#23704](https://github.com/hashicorp/terraform-provider-google/pull/23704)
* gkehub: removed `configmanagement.binauthz` field in `google_gke_hub_feature_membership` [#24076](https://github.com/hashicorp/terraform-provider-google/pull/24076)
* gkehub: removed `description` field in `google_gke_hub_membership` [#23587](https://github.com/hashicorp/terraform-provider-google/pull/23587)
* memorystore: removed `allow_fewer_zones_deployment` field from `google_memorystore_instance` resource because it isn't user-configurable [#24079](https://github.com/hashicorp/terraform-provider-google/pull/24079)
* redis: removed `allow_fewer_zones_deployment` field from `google_redis_cluster` resource because it isn't user-configurable [#24079](https://github.com/hashicorp/terraform-provider-google/pull/24079)
* resourcemanager: removed non-functional `project` field from `google_service_account_key` datasource [#24000](https://github.com/hashicorp/terraform-provider-google/pull/24000)
* vertexai: removed `enable_secure_private_service_connect` in `google_vertex_ai_endpoint` [#23843](https://github.com/hashicorp/terraform-provider-google/pull/23843)

BREAKING INCREASED VALIDATION:
* cloudfunctions2: made `event_type` a required field for `event_trigger` in `google_cloudfunctions2_function` [#23918](https://github.com/hashicorp/terraform-provider-google/pull/23918)
* networkservices: made `load_balancing_scheme` required in `google_network_services_lb_traffic_extension` [#23748](https://github.com/hashicorp/terraform-provider-google/pull/23748)
* sql: made `password_wo_version` required when `password_wo` is set in `google_sql_user` [#24083](https://github.com/hashicorp/terraform-provider-google/pull/24083)
* storage: added validation requiring the `topic` field to be in the form "projects/<project>/topics/<topic>" in `google_storage_notification` [#24135](https://github.com/hashicorp/terraform-provider-google/pull/24135)
* storagetransfer: added path validation for GCS path source and sink in `google_storage_transfer_job` [#23493](https://github.com/hashicorp/terraform-provider-google/pull/23493)
* vertexai: made `metadata`, and `metadata.config` required in `google_vertex_ai_index`. Resource creation would fail without these attributes already, so no change is necessary to existing configurations. [#23971](https://github.com/hashicorp/terraform-provider-google/pull/23971)

OTHER BREAKING CHANGES:
* alloydb: added `deletion_protection` field with a default value of `true` to `google_alloydb_cluster` resource [#24024](https://github.com/hashicorp/terraform-provider-google/pull/24024)
* apigee: changed `certs_info` field in `google_apigee_keystores_aliases_key_cert_file` to be output-only [#24135](https://github.com/hashicorp/terraform-provider-google/pull/24135)
* apigee: migrated `google_apigee_keystores_aliases_key_cert_file` to the plugin framework [#24135](https://github.com/hashicorp/terraform-provider-google/pull/24135)
* artifactregistry: removed the default values for `public_repository` fields in `google_artifact_registry_repository`. If your state is reliant on them, they will now need to be manually included in your configuration. [#23970](https://github.com/hashicorp/terraform-provider-google/pull/23970)
* bigquery: removed the default value of `view.use_legacy_sql` in `google_bigquery_table` [#24065](https://github.com/hashicorp/terraform-provider-google/pull/24065)
* bigtable: renamed instance to `instance_name` for bigtable_table_iam objects [#23399](https://github.com/hashicorp/terraform-provider-google/pull/23399)
* billing: made `budget_filter.credit types` and `budget_filter.subaccounts` no longer optional+computed, only optional, in `google_billing_budget` resource [#24078](https://github.com/hashicorp/terraform-provider-google/pull/24078)
* cloudfunctions2: changed `service_config.service` field in `google_cloudfunctions2_function` resource to be output-only [#23790](https://github.com/hashicorp/terraform-provider-google/pull/23790)
* compute: `subnetworks` and `instances` fields in `google_compute_packet_mirroring` have been converted from arrays to sets [#24021](https://github.com/hashicorp/terraform-provider-google/pull/24021)
* compute: `advertised_ip_ranges` field group in `google_compute_router` has been converted from a list to a set [#24030](https://github.com/hashicorp/terraform-provider-google/pull/24030)
* compute: `disk.type`, `disk.mode` and `disk.interface` no longer use provider configured default values and instead will be set by the API in `google_compute_instance_template` and `google_compute_region_instance_template` resources [#24055](https://github.com/hashicorp/terraform-provider-google/pull/24055)
* provider: fixed many import functions throughout the provider that erroneously matched a subset of the provided input, leading to unclear error messages when using `terraform input` with invalid resource IDs. [#24010](https://github.com/hashicorp/terraform-provider-google/pull/24010)
* resourcemanager: changed `disable_on_destroy` default value to `false` in `google_project_service` [#23951](https://github.com/hashicorp/terraform-provider-google/pull/23951)
* securesourcemanager: changed `deletion_policy` default value from `DELETE` to `PREVENT` [#23963](https://github.com/hashicorp/terraform-provider-google/pull/23963)
* storage: `retention_period` field in `google_storage_bucket` has been converted from `int` to `string` data type [#23535](https://github.com/hashicorp/terraform-provider-google/pull/23535)
* storage: migrated `google_storage_notification` to the plugin framework [#24135](https://github.com/hashicorp/terraform-provider-google/pull/24135)

FEATURES:
* **New Data Source:** `google_artifact_registry_npm_package` ([#24072](https://github.com/hashicorp/terraform-provider-google/pull/24072))
* **New Data Source:** `google_certificate_manager_dns_authorization` ([#24009](https://github.com/hashicorp/terraform-provider-google/pull/24009))
* **New Resource:** `google_iap_web_region_forwarding_rule_service_iam_binding` ([#24041](https://github.com/hashicorp/terraform-provider-google/pull/24041))
* **New Resource:** `google_iap_web_region_forwarding_rule_service_iam_member` ([#24041](https://github.com/hashicorp/terraform-provider-google/pull/24041))
* **New Resource:** `google_iap_web_region_forwarding_rule_service_iam_policy` ([#24041](https://github.com/hashicorp/terraform-provider-google/pull/24041))
* **New Resource:** `google_saas_runtime_saas` ([#24028](https://github.com/hashicorp/terraform-provider-google/pull/24028))

IMPROVEMENTS:
* cloudbuild: added `developer_connect_event_config` field to `google_cloudbuild_trigger` resource ([#24043](https://github.com/hashicorp/terraform-provider-google/pull/24043))
* cloudtasks: added `desired_state` field to `google_cloud_tasks_queue ` resource ([#24053](https://github.com/hashicorp/terraform-provider-google/pull/24053))
* cloudrunv2: added `max_instance_count` field to `google_cloud_run_v2_service` resource. ([#24031](https://github.com/hashicorp/terraform-provider-google/pull/24031))
* compute: added `params.resourceManagerTags` field to the `google_compute_backend_service` ([#24062](https://github.com/hashicorp/terraform-provider-google/pull/24062))
* compute: added `params.resource_manager_tags` field to `google_compute_backend_bucket` ([#24068](https://github.com/hashicorp/terraform-provider-google/pull/24068))
* compute: added `short_name` field to `google_compute_organization_security_policy` resource ([#24059](https://github.com/hashicorp/terraform-provider-google/pull/24059))
* container: added `cluster_autoscaling.default_compute_class_enabled` field to `google_container_cluster` resource ([#24023](https://github.com/hashicorp/terraform-provider-google/pull/24023))
* dialogflowcx: added `enableMultiLanguageTraining`, `locked`, `answerFeedbackSettings`, `personalizationSettings`, `clientCertificateSettings`, `startPlaybook`, `satisfiesPzs`, and `satisfiesPzi` to `google_dialogflow_cx_agent` resource. ([#24007](https://github.com/hashicorp/terraform-provider-google/pull/24007))
* lustre: increased `google_lustre_instance` resource create timeout to 120m from 20m ([#24056](https://github.com/hashicorp/terraform-provider-google/pull/24056))
* oracledatabase: enabled default_from_api flag for ODB Network related fields in `google_oracle_database_cloud_vm_cluster` resource ([#24045](https://github.com/hashicorp/terraform-provider-google/pull/24045))
* sql: added feature to restore `google_sql_database_instance` using `backupdr_backup` ([#24066](https://github.com/hashicorp/terraform-provider-google/pull/24066))
* ssm: made `ca_pool` argument optional for private instances that use Google-managed trusted certificates.` to `secure_source_manager` resource ([#24039](https://github.com/hashicorp/terraform-provider-google/pull/24039))

BUG FIXES:
* container: fixed issue where a failed creation on `google_container_node_pool` would result in an unrecoverable tainted state ([#24077](https://github.com/hashicorp/terraform-provider-google/pull/24077))
* gkeonprem: set `default_from_api` in image field in `google_vmware_node_pool` ([#24022](https://github.com/hashicorp/terraform-provider-google/pull/24022))
* workbench: made `install-monitoring-agent` metadata key settable for `google_workbench_instance` ([#24080](https://github.com/hashicorp/terraform-provider-google/pull/24080))
