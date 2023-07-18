## 4.75.0 (Unreleased)

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
