## 4.48.0 (Unreleased)
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
