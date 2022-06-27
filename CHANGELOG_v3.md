## 3.90.1 (November 02, 2021)

DEPRECATIONS:

* container: fixed an overly-broad deprecation on `master_auth`, constraining it to `master_auth.username` and `master_auth.password`

## 3.90.0 (October 26, 2021)

DEPRECATIONS:
* container: deprecated `workload_identity_config.0.identity_namespace` and it will be removed in a future major release as it has been deprecated in the API. Use `workload_identity_config.0.workload_pool` instead. Switching your configuration from one value to the other will trigger a diff at plan time, and a spurious update. ([#10327](https://github.com/hashicorp/terraform-provider-google/pull/10327))
* container: deprecated the following `google_container_cluster` fields: `instance_group_urls` and `master_auth` ([#10356](https://github.com/hashicorp/terraform-provider-google/pull/10356))

IMPROVEMENTS:
* container: added `node_config.0.guest_accelerator.0.gpu_partition_size` field to google_container_node_pool ([#10339](https://github.com/hashicorp/terraform-provider-google/pull/10339))
* container: added `workload_identity_config.0.workload_pool` to `google_container_cluster` ([#10327](https://github.com/hashicorp/terraform-provider-google/pull/10327))
* container_cluster: Updated `monitoring_config` to accept `WORKLOAD` ([#10321](https://github.com/hashicorp/terraform-provider-google/pull/10321))
* provider: Added links to nested types documentation for manually generated pages ([#10333](https://github.com/hashicorp/terraform-provider-google/pull/10333))

BUG FIXES:
* cloudrun: fixed a permadiff on the field `template.spec.containers.ports.name` of the `google_cloud_run_service` resource ([#10340](https://github.com/hashicorp/terraform-provider-google/pull/10340))
* composer: removed `config.node_config.zone` requirement on `google_composer_environment` ([#10353](https://github.com/hashicorp/terraform-provider-google/pull/10353))
* compute: fixed permadiff for `failover_policy` on `google_compute_region_backend_service` ([#10316](https://github.com/hashicorp/terraform-provider-google/pull/10316))
* compute: fixed to make `description` updatable without recreation on `google_compute_instance_group_manager` ([#10329](https://github.com/hashicorp/terraform-provider-google/pull/10329))
* container: fixed a permadiff on `google_container_node_pool.workload_metadata_config.mode` ([#10313](https://github.com/hashicorp/terraform-provider-google/pull/10313))
* iam: fixed request batching bug where failed requests would show unnecessary backslash escaping to the user. ([#10303](https://github.com/hashicorp/terraform-provider-google/pull/10303))
* securitycenter: fixed bug where `google_scc_notification_config.streaming_config.filter` was not updating. ([#10315](https://github.com/hashicorp/terraform-provider-google/pull/10315))

## 3.89.0 (October 18, 2021)

DEPRECATIONS:
* compute: deprecated the `enable_display` field in `google_compute_instance_template` in the `google` (GA) provider. It will only be available in the `google-beta` provider in a future release, as the underlying feature is in beta. ([#10281](https://github.com/hashicorp/terraform-provider-google/pull/10281))

BUG FIXES:
* compute: fixed bug where `google_compute_router_peer` could not set an advertised route priority of 0, causing permadiff. ([#10292](https://github.com/hashicorp/terraform-provider-google/pull/10292))
* container: fixed a crash on `monitoring_config` of `google_container_cluster` ([#10290](https://github.com/hashicorp/terraform-provider-google/pull/10290))
* iam: fixed request batching bug where failed requests would show unnecessary backslash escaping to the user. ([#10303](https://github.com/hashicorp/terraform-provider-google/pull/10303))
* storage: fixed a bug to better handle eventual consistency among `google_storage_bucket` resources. ([#10287](https://github.com/hashicorp/terraform-provider-google/pull/10287))

## 3.88.0 (October 11, 2021)

NOTES:
* reorganized documentation to group all Compute Engine and Monitoring (Stackdriver) resources together. ([#10205](https://github.com/hashicorp/terraform-provider-google/pull/10205))

DEPRECATIONS:
* container: deprecated `workload_metadata_configuration.node_metadata` in favor of `workload_metadata_configuration.mode` in `google_container_cluster` ([#10238](https://github.com/hashicorp/terraform-provider-google/pull/10238))
* dataproc: deprecated the `google_dataproc_workflow_template.version` field, as it wasn't actually useful. The field is used during updates, but updates aren't currently possible with the resource. ([#10183](https://github.com/hashicorp/terraform-provider-google/pull/10183))
* runtimeconfig: deprecated the Runtime Configurator service in the `google` (GA) provider including `google_runtimeconfig_config`, `google_runtimeconfig_variable`, `google_runtimeconfig_config_iam_policy`, `google_runtimeconfig_config_iam_binding`, `google_runtimeconfig_config_iam_member`, `data.google_runtimeconfig_config`. They will only be available in the `google-beta` provider in a future release, as the underlying service is in beta. ([#10232](https://github.com/hashicorp/terraform-provider-google/pull/10232))
BREAKING CHANGES:
* gke_hub: made the `config_membership` field in `google_gke_hub_feature` required, disallowing invalid configurations ([#10199](https://github.com/hashicorp/terraform-provider-google/pull/10199))
* gke_hub: made the `configmanagement`, `feature`, `location`, `membership` fields in `google_gke_hub_feature_membership` required, disallowing invalid configurations ([#10199](https://github.com/hashicorp/terraform-provider-google/pull/10199))

FEATURES:
* **New Data Source:** `google_service_networking_peered_dns_domain` ([#10229](https://github.com/hashicorp/terraform-provider-google/pull/10229))
* **New Data Source:** `google_sourcerepo_repository` ([#10203](https://github.com/hashicorp/terraform-provider-google/pull/10203))
* **New Data Source:** `google_storage_bucket` ([#10190](https://github.com/hashicorp/terraform-provider-google/pull/10190))
* **New Resource:** `google_pubsub_lite_reservation` ([#10263](https://github.com/hashicorp/terraform-provider-google/pull/10263))
* **New Resource:** `google_service_networking_peered_dns_domain` ([#10229](https://github.com/hashicorp/terraform-provider-google/pull/10229))

IMPROVEMENTS:
* composer: added support for composer v2 fields `workloads_config` and `cloud_composer_network_ipv4_cidr_block` to `composer_environment` ([10269](https://github.com/hashicorp/terraform-provider-google/pull/10269))
* compute: added external IPv6 support on `google_compute_subnetwork` and `google_compute_instance.network_interfaces` ([#10189](https://github.com/hashicorp/terraform-provider-google/pull/10189))
* container: added support for `workload_metadata_configuration.mode` in `google_container_cluster` ([#10238](https://github.com/hashicorp/terraform-provider-google/pull/10238))
* eventarc: added support for `uid` output field, `cloud_function` destination to `google_eventarc_trigger` ([#10199](https://github.com/hashicorp/terraform-provider-google/pull/10199))
* gke_hub: added support for `gcp_service_account_email` when configuring Git sync in `google_gke_hub_feature_membership` ([#10199](https://github.com/hashicorp/terraform-provider-google/pull/10199))
* gke_hub: added support for `resource_state`, `state` outputs to `google_gke_hub_feature` ([#10199](https://github.com/hashicorp/terraform-provider-google/pull/10199))
* pubsub:  Added support for references to `google_pubsub_lite_reservation` to `google_pubsub_lite_topic`. ([#10263](https://github.com/hashicorp/terraform-provider-google/pull/10263))

BUG FIXES:
* monitoring: fixed typo in `google_monitoring_uptime_check_config` where `NOT_MATCHES_REGEX` could not be specified. ([#10249](https://github.com/hashicorp/terraform-provider-google/pull/10249))

## 3.87.0 (October 04, 2021)

DEPRECATIONS:
* dataproc: deprecated the `google_dataproc_workflow_template.version` field, as it wasn't actually useful. The field is used during updates, but updates aren't currently possible with the resource. ([#10183](https://github.com/hashicorp/terraform-provider-google/pull/10183))

FEATURES:
* **New Resource:** `google_org_policy_policy` ([#10111](https://github.com/hashicorp/terraform-provider-google/pull/10111))

IMPROVEMENTS:
* cloudbuild: added field `service_account` to `google_cloudbuild_trigger` ([#10159](https://github.com/hashicorp/terraform-provider-google/pull/10159))
* composer: added field `scheduler_count` to `google_composer_environment` ([#10158](https://github.com/hashicorp/terraform-provider-google/pull/10158))
* compute: Disabled recreation of GCE instances when updating `resource_policies` property ([#10173](https://github.com/hashicorp/terraform-provider-google/pull/10173))
* container: added support for `logging_config` and `monitoring_config` to `google_container_cluster` ([#10125](https://github.com/hashicorp/terraform-provider-google/pull/10125))
* kms: added support for `import_only` to `google_kms_crypto_key` ([#10157](https://github.com/hashicorp/terraform-provider-google/pull/10157))
* networkservices: boosted the default timeout for `google_network_services_edge_cache_origin` from 30m to 60m ([#10182](https://github.com/hashicorp/terraform-provider-google/pull/10182))

BUG FIXES:
* container: fixed an issue where a node pool created with error (eg. GKE_STOCKOUT) would not be captured in state ([#10137](https://github.com/hashicorp/terraform-provider-google/pull/10137))
* filestore: Allowed updating `reserved_ip_range` on `google_filestore_instance` via recreation of the instance ([#10146](https://github.com/hashicorp/terraform-provider-google/pull/10146))
* serviceusage: enabled the service api to retry on failed operation calls in anticipation of transient errors that occur when first enabling the service. ([#10171](https://github.com/hashicorp/terraform-provider-google/pull/10171))

## 3.86.0 (September 27, 2021)

IMPROVEMENTS:
* healthcare: promoted `google_healthcare_hl7_v2_store.parseConfig.version` to GA ([#10099](https://github.com/hashicorp/terraform-provider-google/pull/10099))

BUG FIXES:
* dns: fixed an issue in `google_dns_record_set` where `rrdatas` could not be updated ([#10089](https://github.com/hashicorp/terraform-provider-google/pull/10089))
* dns: fixed an issue in `google_dns_record_set` where creating the resource would result in an 409 error ([#10089](https://github.com/hashicorp/terraform-provider-google/pull/10089))
* platform: fixed a bug in wrongly writing to state when creation failed on `google_organization_policy` ([#10082](https://github.com/hashicorp/terraform-provider-google/pull/10082))

## 3.85.0 (September 20, 2021)
IMPROVEMENTS:
* bigtable: enabled support for `user_project_override` in `google_bigtable_instance` and `google_bigtable_table` ([#10060](https://github.com/hashicorp/terraform-provider-google/pull/10060))
* compute: added `iap` fields to `google_compute_region_backend_service` ([#10038](https://github.com/hashicorp/terraform-provider-google/pull/10038))
* compute: allowed passing an IP address to the `nextHopIlb` field of `google_compute_route` resource ([#10048](https://github.com/hashicorp/terraform-provider-google/pull/10048))
* iam: added `disabled` field to `google_service_account` resource ([#10033](https://github.com/hashicorp/terraform-provider-google/pull/10033))
* provider: added links to nested types documentation within a resource ([#10063](https://github.com/hashicorp/terraform-provider-google/pull/10063))
* storage: added field `path` to `google_storage_transfer_job` ([#10047](https://github.com/hashicorp/terraform-provider-google/pull/10047))

BUG FIXES:
* appengine: fixed bug where `deployment.container.image` would update to an old version even if in `ignore_changes` ([#10058](https://github.com/hashicorp/terraform-provider-google/pull/10058))
* bigquery: fixed a bug where `destination_encryption_config.kms_key_name` stored the version rather than the key name. ([#10068](https://github.com/hashicorp/terraform-provider-google/pull/10068))
* redis: extended the default timeouts on `google_redis_instance` ([#10037](https://github.com/hashicorp/terraform-provider-google/pull/10037))
* serviceusage: fixed an issue in `google_project_service` where users could not reenable services that were disabled outside of Terraform. ([#10045](https://github.com/hashicorp/terraform-provider-google/pull/10045))

## 3.84.0 (September 13, 2021)
FEATURES:
* **New Data Source:** `google_secret_manager_secret` ([#9983](https://github.com/hashicorp/terraform-provider-google/pull/9983))

IMPROVEMENTS:
* compute: added update support to `google_compute_service_attachment` ([#9982](https://github.com/hashicorp/terraform-provider-google/pull/9982))

BUG FIXES:
* container: fixed a bug in failing to remove `maintenance_exclusion` on `google_container_cluster` ([#10025](https://github.com/hashicorp/terraform-provider-google/pull/10025))
* compute: fixed an issue in `google_compute_router_nat` where removing `log_config` resulted in a perma-diff ([#9950](https://github.com/hashicorp/terraform-provider-google/pull/9950))
* compute: fixed `advanced_machine_features` error messages in `google_compute_instance` ([#10023](https://github.com/hashicorp/terraform-provider-google/pull/10023))
* eventarc: fixed bug where resources deleted outside of Terraform would cause errors ([#9997](https://github.com/hashicorp/terraform-provider-google/pull/9997))
* functions: fixed an error message on `google_cloudfunctions_function` ([#10011](https://github.com/hashicorp/terraform-provider-google/pull/10011))
* logging: fixed the data type for `bucket_options.linear_buckets.width` on `google_logging_metric` ([#9985](https://github.com/hashicorp/terraform-provider-google/pull/9985))
* osconfig: fixed import on `google_os_config_guest_policies` ([#10019](https://github.com/hashicorp/terraform-provider-google/pull/10019))
* storage: fixed an undetected change on `days_since_noncurrent_time` of `google_storage_bucket` ([#10024](https://github.com/hashicorp/terraform-provider-google/pull/10024))


## 3.83.0 (September 09, 2021)
FEATURES:
* **New Resource:** `google_privateca_certificate_template` ([#9905](https://github.com/hashicorp/terraform-provider-google/pull/9905))

IMPROVEMENTS:
* privateca: added `certificate_template` to `google_privateca_certificate`. ([#9915](https://github.com/hashicorp/terraform-provider-google/pull/9915))
* compute: allowed setting `ip_address` field of `google_compute_router_peer` ([#9913](https://github.com/hashicorp/terraform-provider-google/pull/9913))
* compute: promoted `google_compute_service_attachment` to ga ([#9914](https://github.com/hashicorp/terraform-provider-google/pull/9914))
* compute: promoted `role` and `purpose` fields in `google_compute_subnetwork` to ga ([#9914](https://github.com/hashicorp/terraform-provider-google/pull/9914))
* kms: added support for `destroy_scheduled_duration` to `google_kms_crypto_key` ([#9911](https://github.com/hashicorp/terraform-provider-google/pull/9911))

BUG FIXES:
* endpoints: fixed a timezone discrepancy in `config_id` on `google_endpoints_service` ([#9912](https://github.com/hashicorp/terraform-provider-google/pull/9912))
* cloudbuild: marked `google_cloudbuild_trigger` as requiring one of branch_name/tag_name/commit_sha  within build.source.repo_source ([#9952](https://github.com/hashicorp/terraform-provider-google/pull/9952))
* compute: fixed a crash on `enable` field of `google_compute_router_peer` ([#9940](https://github.com/hashicorp/terraform-provider-google/pull/9940))
* compute: fixed a permanent diff for `next_hop_instance_zone` on `google_compute_route` when `next_hop_instance` was set to a self link ([#9931](https://github.com/hashicorp/terraform-provider-google/pull/9931))
* compute: fixed an issue in `google_compute_router_nat` where removing `log_config` resulted in a perma-diff ([#9950](https://github.com/hashicorp/terraform-provider-google/pull/9950))
* privateca: fixed a permadiff bug for `publishing_options` on `google_privateca_ca_pool` when both attributes set false ([#9926](https://github.com/hashicorp/terraform-provider-google/pull/9926))
* spanner: fixed instance updates to processing units ([#9933](https://github.com/hashicorp/terraform-provider-google/pull/9933))
* storage: added support for timeouts on `google_storage_bucket_object` ([#9937](https://github.com/hashicorp/terraform-provider-google/pull/9937))

## 3.82.0 (August 30, 2021)
FEATURES:
* **New Resource:** `google_privateca_certificate_template` ([#9905](https://github.com/hashicorp/terraform-provider-google/pull/9905))
* **New Resource:** `google_compute_firewall_policy` ([#9887](https://github.com/hashicorp/terraform-provider-google/pull/9887))
* **New Resource:** `google_compute_firewall_policy_association` ([#9887](https://github.com/hashicorp/terraform-provider-google/pull/9887))
* **New Resource:** `google_compute_firewall_policy_rule` ([#9887](https://github.com/hashicorp/terraform-provider-google/pull/9887))

IMPROVEMENTS:
* sql: added field `collation` to `google_sql_database_instance` ([#9888](https://github.com/hashicorp/terraform-provider-google/pull/9888))

BUG FIXES:
* apigateway: fixed import functionality for all `apigateway` resources ([#9871](https://github.com/hashicorp/terraform-provider-google/pull/9871))
* dns: fixed not-exists error message on data source `google_dns_managed_zone` ([#9898](https://github.com/hashicorp/terraform-provider-google/pull/9898))
* healthcare: fixed bug where changes to `google_healthcare_hl7_v2_store.parser_config` subfields would error with "...parser_config.version field is immutable..." ([#9900](https://github.com/hashicorp/terraform-provider-google/pull/9900))
* os_config: fixed imports for `google_os_config_guest_policies` ([#9872](https://github.com/hashicorp/terraform-provider-google/pull/9872))
* pubsub: added polling to `google_pubsub_schema` to deal with eventually consistent deletes ([#9863](https://github.com/hashicorp/terraform-provider-google/pull/9863))
* secretmanager: fixed an issue where `replication` fields would not update in `google_secret_manager_secret` ([#9894](https://github.com/hashicorp/terraform-provider-google/pull/9894))
* service_usage: fixed imports on `google_service_usage_consumer_quota_override` ([#9876](https://github.com/hashicorp/terraform-provider-google/pull/9876))
* sql: fixed a permadiff bug for `type` when BUILT_IN on `google_sql_user` ([#9864](https://github.com/hashicorp/terraform-provider-google/pull/9864))
* sql: fixed bug in `google_sql_user` with CLOUD_IAM_USERs on POSTGRES. ([#9859](https://github.com/hashicorp/terraform-provider-google/pull/9859))

## 3.81.0 (August 23, 2021)

IMPROVEMENTS:
* compute: Added `enable` attribute to `google_compute_router_peer` ([#9776](https://github.com/hashicorp/terraform-provider-google/pull/9776))
* compute: added support for `L3_DEFAULT` as `ip_protocol` for `google_compute_forwarding_rule` and `UNSPECIFIED` as `protocol` for `google_compute_region_backend_service` to support network load balancers that forward all protocols and ports. ([#9799](https://github.com/hashicorp/terraform-provider-google/pull/9799))
* compute: added support for `security_settings` to `google_compute_backend_service` ([#9797](https://github.com/hashicorp/terraform-provider-google/pull/9797))
* essentialcontacts: promoted `google_essential_contacts_contact` to GA ([#9822](https://github.com/hashicorp/terraform-provider-google/pull/9822))
* gkehub: added `google_gke_hub_membership` support for both `//container.googleapis.com/${google_container_cluster.my-cluster.id}` and `google_container_cluster.my-cluster.id` in `endpoint.0.gke_cluster.0.resource_link` ([#9765](https://github.com/hashicorp/terraform-provider-google/pull/9765))
* provider: Added provider support for `request_reason` ([#9794](https://github.com/hashicorp/terraform-provider-google/pull/9794))
* provider: added support for `billing_project` across all resources. If `user_project_override` is set to `true` and a `billing_project` is set, the `X-Goog-User-Project` header will be sent for all resources. ([#9852](https://github.com/hashicorp/terraform-provider-google/pull/9852))

BUG FIXES:
* assuredworkloads: fixed resource deletion so `google_assured_workloads_workload` can delete what it creates ([#9835](https://github.com/hashicorp/terraform-provider-google/pull/9835))
* bigquery: fixed the permadiff bug on `location` of the `google_bigquery_dataset` ([#9810](https://github.com/hashicorp/terraform-provider-google/pull/9810))
* composer: fixed environment version regexp to explicitly require . (dot) instead of any character after 'preview' (example: composer-2.0.0-preview.0-airflow-2.1.1) ([#9804](https://github.com/hashicorp/terraform-provider-google/pull/9804))
* compute: changed `wait_for_instances` in `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` to no longer block plan / refresh, waiting on managed instance statuses during apply instead ([#9832](https://github.com/hashicorp/terraform-provider-google/pull/9832))
* compute: fixed a bug where `negative_caching_policy` cannot be set always revalidate on `google_compute_backend_service` ([#9821](https://github.com/hashicorp/terraform-provider-google/pull/9821))
* compute: fixed instances where compute resource calls would have their urls appended with a redundant `/projects` after the host ([#9834](https://github.com/hashicorp/terraform-provider-google/pull/9834))
* firestore: removed diff for server generated field `__name__` on `google_firestore_index` ([#9820](https://github.com/hashicorp/terraform-provider-google/pull/9820))
* privateca: fixed the creation of subordinate `google_privateca_certificate_authority` with `max_issuer_path_length = 0`. ([#9856](https://github.com/hashicorp/terraform-provider-google/pull/9856))
* privateca: Fixed null for `ignore_active_certificates_on_deletion` on the imported `google_privateca_certificate_authority` ([#9781](https://github.com/hashicorp/terraform-provider-google/pull/9781))

## 3.80.0 (August 16, 2021)

FEATURES:
* **New Resource:** `google_dialogflow_cx_environment` ([#9738](https://github.com/hashicorp/terraform-provider-google/pull/9738))

IMPROVEMENTS:
* gkehub: added support for both `//container.googleapis.com/${google_container_cluster.my-cluster.id}` and `google_container_cluster.my-cluster.id` references in `google_gke_hub_membership.endpoint.0.gke_cluster.0.resource_link` ([#9765](https://github.com/hashicorp/terraform-provider-google/pull/9765))
* kms: added `name` field to `google_kms_crypto_key_version` datasource ([#9762](https://github.com/hashicorp/terraform-provider-google/pull/9762))

BUG FIXES:
* apigee: fixed update behavior on `google_apigee_envgroup` ([#9740](https://github.com/hashicorp/terraform-provider-google/pull/9740))
* privateca: fixed a failure to create `google_privateca_certificate_authority` of type `SUBORDINATE` due to an invalid attempt to activate it on creation. ([#9761](https://github.com/hashicorp/terraform-provider-google/pull/9761))

## 3.79.0 (August 09, 2021)

NOTES:
* spanner: The `num_nodes` field on `google_spanner_instance` will have its default removed in a future major release, and either `num_nodes` or `processing_units` will be required. ([#9716](https://github.com/hashicorp/terraform-provider-google/pull/9716))

FEATURES:
* **New Resource:** `google_dialogflow_cx_entity_type` ([#9717](https://github.com/hashicorp/terraform-provider-google/pull/9717))
* **New Resource:** `google_dialogflow_cx_page` ([#9683](https://github.com/hashicorp/terraform-provider-google/pull/9683))

IMPROVEMENTS:
* spanner: added `processing_units` to `google_spanner_instance` ([#9716](https://github.com/hashicorp/terraform-provider-google/pull/9716))
* storage: added support for `customer_encryption` on `resource_storage_bucket_object` ([#9704](https://github.com/hashicorp/terraform-provider-google/pull/9704))


## 3.78.0 (August 02, 2021)
FEATURES:
* **New Resource:** `google_gke_hub_membership` ([#9616](https://github.com/hashicorp/terraform-provider-google/pull/9616))

IMPROVEMENTS:
* servicenetworking: added support for `user_project_override` and `billing_project ` to `google_service_networking_connection` ([#9668](https://github.com/hashicorp/terraform-provider-google/pull/9668))

BUG FIXES:
* storagetransfer: Fixed a crash on `azure_blob_storage_data_source` for `google_storage_transfer_job` ([#9644](https://github.com/hashicorp/terraform-provider-google/pull/9644))
* sql: fixed bug that wouldn't insert the `google_sql_user` in state for iam users. ([#9625](https://github.com/hashicorp/terraform-provider-google/pull/9625))
* storage: fixed a crash when `azure_credentials` was defined in `google_storage_transfer_job` ([#9671](https://github.com/hashicorp/terraform-provider-google/pull/9671))

## 3.77.0 (July 26, 2021)

FEATURES:
* **New Resource:** `google_scc_notification_config` ([#9578](https://github.com/hashicorp/terraform-provider-google/pull/9578))

IMPROVEMENTS:
* compute: fixed a permadiff bug in `log_config` field of `google_compute_region_backend_service` ([#9568](https://github.com/hashicorp/terraform-provider-google/pull/9568))
* dlp: added `crypto_replace_ffx_fpe_config` and `crypto_replace_ffx_fpe_config` as primitive transformation types to `google_data_loss_prevention_deidentify_template` ([#9572](https://github.com/hashicorp/terraform-provider-google/pull/9572))

BUG FIXES:
* bigquerydatatransfer: fixed a bug where `destination_dataset_id` was required, it is now optional. ([#9605](https://github.com/hashicorp/terraform-provider-google/pull/9605))
* billing: Fixed ordering of `budget_filter. projects` on `google_billing_budget` ([#9598](https://github.com/hashicorp/terraform-provider-google/pull/9598))
* compute: removed default value of `0.8` from `google_backend_service.backend.max_utilization` and it will now default from API. All `max_connections_xxx` and `max_rate_xxx` will also default from API as these are all conditional on balancing mode. ([#9587](https://github.com/hashicorp/terraform-provider-google/pull/9587))
* sql: fixed bug where the provider would retry on an error if the database instance name couldn't be reused. ([#9591](https://github.com/hashicorp/terraform-provider-google/pull/9591))

## 3.76.0 (July 19, 2021)

FEATURES:
* **New Resource:** `google_dialogflow_cx_flow` ([#9551](https://github.com/hashicorp/terraform-provider-google/pull/9551))
* **New Resource:** `google_dialogflow_cx_intent` ([#9537](https://github.com/hashicorp/terraform-provider-google/pull/9537))
* **New Resource:** `google_dialogflow_cx_version` ([#9554](https://github.com/hashicorp/terraform-provider-google/pull/9554))
* **New Resource:** `google_network_services_edge_cache_keyset` ([#9540](https://github.com/hashicorp/terraform-provider-google/pull/9540))
* **New Resource:** `google_network_services_edge_cache_origin` ([#9540](https://github.com/hashicorp/terraform-provider-google/pull/9540))
* **New Resource:** `google_network_services_edge_cache_service` ([#9540](https://github.com/hashicorp/terraform-provider-google/pull/9540))

IMPROVEMENTS:
* apigee: Added SLASH_22 support for `peering_cidr_range` on `google_apigee_instance` ([#9558](https://github.com/hashicorp/terraform-provider-google/pull/9558))
* cloudbuild: Added `pubsub_config` and `webhook_config` parameter to `google_cloudbuild_trigger`. ([#9541](https://github.com/hashicorp/terraform-provider-google/pull/9541))

BUG FIXES:
* pubsub: fixed pubsublite update issues ([#9544](https://github.com/hashicorp/terraform-provider-google/pull/9544))

## 3.75.0 (July 12, 2021)

FEATURES:
* **New Resource:** google_privateca_ca_pool ([#9480](https://github.com/hashicorp/terraform-provider-google/pull/9480))
* **New Resource:** google_privateca_certificate ([#9480](https://github.com/hashicorp/terraform-provider-google/pull/9480))
* **New Resource:** google_privateca_certificate_authority ([#9480](https://github.com/hashicorp/terraform-provider-google/pull/9480))

IMPROVEMENTS:
* bigquery: added `kms_key_version` as an output on `bigquery_table.encryption_configuration` and the `destination_encryption_configuration` blocks of `bigquery_job.query`, `bigquery_job.load`, and `bigquery_copy`. ([#9500](https://github.com/hashicorp/terraform-provider-google/pull/9500))
* compute: added `advanced_machine_features` to `google_compute_instance` ([#9470](https://github.com/hashicorp/terraform-provider-google/pull/9470))
* compute: promoted all `cdn_policy` sub fields in `google_compute_backend_service`, `google_compute_region_backend_service` and `google_compute_backend_bucket` to GA ([#9432](https://github.com/hashicorp/terraform-provider-google/pull/9432))
* dlp: Added `replace_with_info_type_config` to `dlp_deidentify_template`. ([#9446](https://github.com/hashicorp/terraform-provider-google/pull/9446))
* storage: added `temporary_hold` and `event_based_hold` attributes to `google_storage_bucket_object` ([#9487](https://github.com/hashicorp/terraform-provider-google/pull/9487))

BUG FIXES:
* bigquery: Fixed permadiff due to lowercase mode/type in `google_bigquery_table.schema` ([#9499](https://github.com/hashicorp/terraform-provider-google/pull/9499))
* billing: made `all_updates_rule.*` fields updatable on `google_billing_budget` ([#9473](https://github.com/hashicorp/terraform-provider-google/pull/9473))
* billing: made `amount.specified_amount.units` updatable on `google_billing_budget` ([#9465](https://github.com/hashicorp/terraform-provider-google/pull/9465))
* compute: fixed perma-diff in `google_compute_instance` ([#9460](https://github.com/hashicorp/terraform-provider-google/pull/9460))
* storage: fixed handling of object paths that contain slashes for `google_storage_object_access_control` ([#9502](https://github.com/hashicorp/terraform-provider-google/pull/9502))

## 3.74.0 (June 28, 2021)
FEATURES:
* **New Resource:** `google_app_engine_service_network_settings` ([#9414](https://github.com/hashicorp/terraform-provider-google/pull/9414))
* **New Resource:** `google_vertex_ai_dataset` ([#9411](https://github.com/hashicorp/terraform-provider-google/pull/9411))
* **New Resource:** `google_cloudbuild_worker_pool` ([#9417](https://github.com/hashicorp/terraform-provider-google/pull/9417))

IMPROVEMENTS:
* bigtable: added `cluster.kms_key_name` field to `google_bigtable_instance` ([#9393](https://github.com/hashicorp/terraform-provider-google/pull/9393))
* compute: promoted all `cdn_policy` sub fields in `google_compute_backend_service`, `google_compute_region_backend_service` and `google_compute_backend_bucket` to GA ([#9432](https://github.com/hashicorp/terraform-provider-google/pull/9432))
* secretmanager: added `ttl`, `expire_time`, `topics` and `rotation` fields to `google_secret_manager_secret` ([#9398](https://github.com/hashicorp/terraform-provider-google/pull/9398))

BUG FIXES:
* container: allowed setting `node_config.service_account` at the same time as `enable_autopilot = true` for `google_container_cluster` ([#9399](https://github.com/hashicorp/terraform-provider-google/pull/9399))
* container: fixed issue where creating a node pool with a name that already exists would import that resource. `google_container_node_pool` ([#9424](https://github.com/hashicorp/terraform-provider-google/pull/9424))
* dataproc: fixed crash when creating `google_dataproc_workflow_template` with `secondary_worker_config` empty except for `num_instances = 0` ([#9381](https://github.com/hashicorp/terraform-provider-google/pull/9381))
* filestore: fixed an issue in `google_filestore_instance` where creating two instances simultaneously resulted in an error. ([#9396](https://github.com/hashicorp/terraform-provider-google/pull/9396))
* sql: added support for `binary_logging` on replica instances for `googe_sql_database_instance` ([#9428](https://github.com/hashicorp/terraform-provider-google/pull/9428))

## 3.73.0 (June 21, 2021)
FEATURES:
* **New Resource:** `google_dialogflow_cx_agent` ([#9338](https://github.com/hashicorp/terraform-provider-google/pull/9338))

IMPROVEMENTS:
* provider: added support for [mtls authentication](https://google.aip.dev/auth/4114) ([#9382](https://github.com/hashicorp/terraform-provider-google/pull/9382))
* compute: added `advanced_machine_features` fields to `google_compute_instance_template` ([#9363](https://github.com/hashicorp/terraform-provider-google/pull/9363))
* compute: promoted `custom_response_headers` to GA for `google_compute_backend_service` and `google_compute_backend_bucket` ([#9374](https://github.com/hashicorp/terraform-provider-google/pull/9374))
* redis: allowed `redis_version` to be upgraded on `google_redis_instance` ([#9378](https://github.com/hashicorp/terraform-provider-google/pull/9378))
* redis: promoted fields `transit_encryption_mode` and `server_ca_certs` to GA on `google_redis_instance` ([#9378](https://github.com/hashicorp/terraform-provider-google/pull/9378))

BUG FIXES:
* apigee: added SLASH_23 support for `peering_cidr_range` on `google_apigee_instance` ([#9343](https://github.com/hashicorp/terraform-provider-google/pull/9343))
* cloudrun: fixed a bug where plan would should a diff on `google_cloud_run_service` if the order of the `template.spec.containers.env` list was re-ordered outside of terraform. ([#9340](https://github.com/hashicorp/terraform-provider-google/pull/9340))
* container: added `user_project_override` support to the ContainerOperationWaiter used by `google_container_cluster` ([#9379](https://github.com/hashicorp/terraform-provider-google/pull/9379))

## 3.72.0 (June 14, 2021)
IMPROVEMENTS:
* compute: added support for IPsec-encrypted Interconnect in the form of new fields on `google_compute_router`, `google_compute_ha_vpn_gateway`, `google_compute_interconnect_attachment` and `google_compute_address`([#9288](https://github.com/hashicorp/terraform-provider-google/pull/9288))
* container: Allowed specifying a cluster id field for `google_container_node_pool.cluster` to ensure that a node pool is recreated if the associated cluster is recreated. ([#9309](https://github.com/hashicorp/terraform-provider-google/pull/9309))
* storagetransfer: added support for `azure_blob_storage_data_source` to `google_storage_transfer_job` ([#9311](https://github.com/hashicorp/terraform-provider-google/pull/9311))

BUG FIXES:
* bigquery: Fixed `google_bigquery_table.schema` handling of policyTags ([#9302](https://github.com/hashicorp/terraform-provider-google/pull/9302))
* bigtable: fixed bug that would error if creating multiple bigtable gc policies at the same time ([#9305](https://github.com/hashicorp/terraform-provider-google/pull/9305))
* compute: fixed bug where `encryption` showed a perma-diff on resources created prior to the feature being released. ([#9303](https://github.com/hashicorp/terraform-provider-google/pull/9303))

## 3.71.0 (June 07, 2021)
FEATURES:
* **New Resource:** `google_dialogflow_fulfillment` ([#9253](https://github.com/hashicorp/terraform-provider-google/pull/9253))

IMPROVEMENTS:
* compute: added `reservation_affinity` to `google_compute_instance` and `google_compute_instance_template` ([#9256](https://github.com/hashicorp/terraform-provider-google/pull/9256))
* compute: added support for `wait_for_instances_status` on `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#9231](https://github.com/hashicorp/terraform-provider-google/pull/9231))
* compute: added support for output-only `status` field on `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#9231](https://github.com/hashicorp/terraform-provider-google/pull/9231))
* compute: promoted `log_config` field of `google_compute_health_check` and `google_compute_region_health_check` to GA ([#9274](https://github.com/hashicorp/terraform-provider-google/pull/9274))
* compute: set the default value for log_config.enable on `google_compute_region_health_check` to avoid permanent diff on plan/apply. ([#9274](https://github.com/hashicorp/terraform-provider-google/pull/9274))

BUG FIXES:
* composer: fixed a check that did not allow for preview versions in `google_composer_environment` ([#9255](https://github.com/hashicorp/terraform-provider-google/pull/9255))
* storage: fixed error when `matches_storage_class` is set empty on `google_storage_bucket` ([#9221](https://github.com/hashicorp/terraform-provider-google/pull/9221))
* vpcaccess: fixed permadiff when `max_throughput` is not set on `google_vpc_access_connector` ([#9282](https://github.com/hashicorp/terraform-provider-google/pull/9282))

## 3.70.0 (June 01, 2021)

IMPROVEMENTS:
* compute: added `provisioned_iops` to `google_compute_disk` ([#9193](https://github.com/hashicorp/terraform-provider-google/pull/9193))
* compute: promoted `distribution_policy_target_shape` field in `google_compute_region_instance_group_manager` to GA. ([#9186](https://github.com/hashicorp/terraform-provider-google/pull/9186))
* sql: added field `disk_autoresize_limit` to `sql_database_instance` ([#9203](https://github.com/hashicorp/terraform-provider-google/pull/9203))

BUG FIXES:
* cloudrun: fixed a bug where resources would return successfully due to responses based on a previous version of the resource ([#9213](https://github.com/hashicorp/terraform-provider-google/pull/9213))
* storage: fixed error when `matches_storage_class` is set empty on `google_storage_bucket` ([#9221](https://github.com/hashicorp/terraform-provider-google/pull/9221))

## 3.69.0 (May 24, 2021)

IMPROVEMENTS:
* compute: added "description" field to "google_compute_resource_policy" resource ([#9176](https://github.com/hashicorp/terraform-provider-google/pull/9176))
* compute: added "instance_schedule_policy" field to "google_compute_resource_policy" resource ([#9176](https://github.com/hashicorp/terraform-provider-google/pull/9176))
* compute: promoted field `autoscaling_policy.scaling_schedules` on `google_compute_autoscaler` and `google_compute_region_autoscaler` to ga ([#9165](https://github.com/hashicorp/terraform-provider-google/pull/9165))
* compute: promoted `autoscaling_policy.cpu_utilization.predictive_method` on `google_compute_autoscaler` and `google_compute_region_autoscaler` to ga. ([#9156](https://github.com/hashicorp/terraform-provider-google/pull/9156))

BUG FIXES:
* cloudidentity: fixed recreation on the `initial_group_config` of `google_cloud_identity_group` ([#9143](https://github.com/hashicorp/terraform-provider-google/pull/9143))
* compute: added mutex in `google_compute_metadata_item` to reduce retries + quota errors ([#9168](https://github.com/hashicorp/terraform-provider-google/pull/9168))
* container: fixed bug where `enable_shielded_nodes` could not be false on resource `google_container_cluster` ([#9131](https://github.com/hashicorp/terraform-provider-google/pull/9131))

## 3.68.0 (May 18, 2021)
FEATURES:
* **New Resource:** `google_pubsub_schema` ([#9116](https://github.com/hashicorp/terraform-provider-google/pull/9116))

IMPROVEMENTS:
* compute: added `initial_size`  in resource `google_compute_node_group` to account for scenarios where size may change under the hood ([#9078](https://github.com/hashicorp/terraform-provider-google/pull/9078))
* compute: added support for setting `kms_key_name` on `google_compute_machine_image` ([#9107](https://github.com/hashicorp/terraform-provider-google/pull/9107))
* dataflow: enabled updates for `google_dataflow_flex_template_job` ([#9123](https://github.com/hashicorp/terraform-provider-google/pull/9123))

BUG FIXES:
* compute: fixed bug where, when an organization security policy association was removed outside of terraform, the next plan/apply would fail. ([#9095](https://github.com/hashicorp/terraform-provider-google/pull/9095))
* container: added validation to check that both `node_version` and `remove_default_node_pool` cannot be set on `google_container_cluster` ([#9100](https://github.com/hashicorp/terraform-provider-google/pull/9100))
* dns: suppressed spurious diffs due to case changes in DS records ([#9099](https://github.com/hashicorp/terraform-provider-google/pull/9099))

## 3.67.0 (May 10, 2021)
FEATURES:
* **New Resource:** google_memcache_instance ([#8982](https://github.com/hashicorp/terraform-provider-google/pull/8982))

NOTES:
* all: changed default HTTP request timeout from 30 seconds to 120 seconds ([#8966](https://github.com/hashicorp/terraform-provider-google/pull/8966))
DEPRECATIONS:
* compute: deprecated `distribution_policy_target_shape` in `google_compute_region_instance_group_manager` Use the `google-beta` provider to continue using this field ([#8970](https://github.com/hashicorp/terraform-provider-google/pull/8970))
* compute: deprecated `min_ready_sec` in `google_compute_region_instance_group_manager` & `google_compute_instance_group_manager` Use the `google-beta` provider to continue using this field ([#8970](https://github.com/hashicorp/terraform-provider-google/pull/8970))
* container: deprecated `pod_security_policy_config` field on resource `google_container_cluster`. Use the `google-beta` provider to continue using this field ([#8970](https://github.com/hashicorp/terraform-provider-google/pull/8970))

BREAKING CHANGES:
* bigquery: updating `dataset_id` or `project_id` in `google_bigquery_dataset` will now recreate the resource ([#8973](https://github.com/hashicorp/terraform-provider-google/pull/8973))

IMPROVEMENTS:
* accesscontextmanager: added support for `require_verified_chrome_os` in basic access levels. ([#9071](https://github.com/hashicorp/terraform-provider-google/pull/9071))
* billingbudget: added support for import of `google_billing_budget` ([#8990](https://github.com/hashicorp/terraform-provider-google/pull/8990))
* cloud_identity: added support for `initial_group_config` to the google_cloud_identity_group resource ([#9035](https://github.com/hashicorp/terraform-provider-google/pull/9035))
* cloudrun: added support to bind secrets from Secret Manager to environment variables or files to `google_cloud_run_service` ([#9073](https://github.com/hashicorp/terraform-provider-google/pull/9073))
* compute: added `initial_size` to account for scenarios where size may change under the hood in resource `google_compute_node_group` ([#9078](https://github.com/hashicorp/terraform-provider-google/pull/9078))
* healthcare: added support for `stream_configs` in `google_healthcare_dicom_store` ([#8986](https://github.com/hashicorp/terraform-provider-google/pull/8986))
* secretmanager: added support for setting a CMEK on `google_secret_manager_secret` ([#9046](https://github.com/hashicorp/terraform-provider-google/pull/9046))
* spanner: added `force_destroy` to `google_spanner_instance` to delete instances that have backups enabled. ([#9076](https://github.com/hashicorp/terraform-provider-google/pull/9076))
* spanner: added support for setting a CMEK on `google_spanner_database` ([#8966](https://github.com/hashicorp/terraform-provider-google/pull/8966))
* workflows: marked `source_contents` and `service_account` as updatable on `google_workflows_workflow` ([#9018](https://github.com/hashicorp/terraform-provider-google/pull/9018))

BUG FIXES:
* bigquery: fixed `dataset_id` to force new resource if name is changed. ([#8973](https://github.com/hashicorp/terraform-provider-google/pull/8973))
* cloudrun: fixed permadiff on `google_cloud_run_domain_mapping.metadata.labels` ([#8971](https://github.com/hashicorp/terraform-provider-google/pull/8971))
* composer: changed `google_composer_environment.master_ipv4_cidr_block` to draw default from the API ([#9017](https://github.com/hashicorp/terraform-provider-google/pull/9017))
* container: fixed container node pool not removed from the state when received 404 error on delete call for the resource `google_container_node_pool` ([#9034](https://github.com/hashicorp/terraform-provider-google/pull/9034))
* dns: fixed empty `rrdatas` list on `google_dns_record_set` for AAAA records ([#9029](https://github.com/hashicorp/terraform-provider-google/pull/9029))
* kms: fixed indirectly force replacement via `skip_initial_version_creation` on `google_kms_crypto_key` ([#8988](https://github.com/hashicorp/terraform-provider-google/pull/8988))
* logging: fixed `metric_descriptor.labels` can't be updated on 'google_logging_metric' ([#9057](https://github.com/hashicorp/terraform-provider-google/pull/9057))
* pubsub: fixed diff for `minimum_backoff` & `maximum_backoff` on `google_pubsub_subscription` ([#9048](https://github.com/hashicorp/terraform-provider-google/pull/9048))
* resourcemanager: fixed broken handling of IAM conditions for `google_organization_iam_member`, `google_organization_iam_binding`, and `google_organization_iam_policy` ([#9047](https://github.com/hashicorp/terraform-provider-google/pull/9047))
* serviceusage: added `google_project_service.service` validation to reject invalid service domains that don't contain a period ([#8987](https://github.com/hashicorp/terraform-provider-google/pull/8987))
* storage: fixed bug where `role_entity` user wouldn't update if the role changed. ([#9008](https://github.com/hashicorp/terraform-provider-google/pull/9008))

## 3.66.1 (April 29, 2021)
BUG FIXES:
* compute: fixed bug where terraform would crash if updating from no `service_account.scopes` to more. ([#9032](https://github.com/hashicorp/terraform-provider-google/pull/9032))

## 3.66.0 (April 28, 2021)
NOTES:
* all: changed default HTTP request timeout from 30 seconds to 120 seconds ([#8966](https://github.com/hashicorp/terraform-provider-google/pull/8966))

BREAKING CHANGES:
* datacatalog: updating `parent` in `google_data_catalog_tag` will now recreate the resource ([#8964](https://github.com/hashicorp/terraform-provider-google/pull/8964))

FEATURES:
* **New Data Source:** `google_compute_ha_vpn_gateway` ([#8952](https://github.com/hashicorp/terraform-provider-google/pull/8952))
* **New Resource:** `google_dataproc_workflow_template` ([#8962](https://github.com/hashicorp/terraform-provider-google/pull/8962))

IMPROVEMENTS:
* bigquery: Added BigTable source format in BigQuery table ([#8923](https://github.com/hashicorp/terraform-provider-google/pull/8923))
* cloudfunctions: removed bounds on the supported memory range in `google_cloudfunctions_function.available_memory_mb` ([#8946](https://github.com/hashicorp/terraform-provider-google/pull/8946))
* compute: marked scheduling.0.node_affinities as updatable in `google_compute_instance` ([#8927](https://github.com/hashicorp/terraform-provider-google/pull/8927))
* dataproc: added `shielded_instance_config` fields to `google_dataproc_cluster` ([#8910](https://github.com/hashicorp/terraform-provider-google/pull/8910))
* spanner: added support for setting a CMEK on `google_spanner_database` ([#8966](https://github.com/hashicorp/terraform-provider-google/pull/8966))

BUG FIXES:
* compute: fixed error when creating empty `scopes` on `google_compute_instance` ([#8953](https://github.com/hashicorp/terraform-provider-google/pull/8953))
* container: fixed a bug that allowed specifying `node_config` on `google_container_cluster` when autopilot is used ([#8905](https://github.com/hashicorp/terraform-provider-google/pull/8905))
* datacatalog: fixed an issue where `parent` in `google_data_catalog_tag` attempted to update the resource when change instead of recreating it ([#8964](https://github.com/hashicorp/terraform-provider-google/pull/8964))
* datacatalog: set default false for `force_delete` on `google_data_catalog_tag_template` ([#8922](https://github.com/hashicorp/terraform-provider-google/pull/8922))
* dns: added missing record types to `google_dns_record_set` resource ([#8919](https://github.com/hashicorp/terraform-provider-google/pull/8919))
* sql: set `clone.point_in_time` optional for `google_sql_database_instance` ([#8965](https://github.com/hashicorp/terraform-provider-google/pull/8965))

## 3.65.0 (April 20, 2021)

FEATURES:
* **New Resource:** google_eventarc_trigger ([#8895](https://github.com/hashicorp/terraform-provider-google/pull/8895))

IMPROVEMENTS:
* compute: added the ability to specify `google_compute_forwarding_rule.ip_address` by a reference in addition to raw IP address ([#8877](https://github.com/hashicorp/terraform-provider-google/pull/8877))
* compute: enabled fields `advertiseMode`, `advertisedGroups`, `peerAsn`, and `peerIpAddress` to be updatable on resource `google_compute_router_peer` ([#8862](https://github.com/hashicorp/terraform-provider-google/pull/8862))

BUG FIXES:
* cloud_identity: fixed google_cloud_identity_group_membership import/update ([#8867](https://github.com/hashicorp/terraform-provider-google/pull/8867))
* compute: fixed an issue in `google_compute_instance` where `min_node_cpus` could not be set ([#8865](https://github.com/hashicorp/terraform-provider-google/pull/8865))
* compute: removed minimum for `scopes` field on `google_compute_instance` resource ([#8893](https://github.com/hashicorp/terraform-provider-google/pull/8893))
* iam: fixed issue with principle and principleSet members not retaining their casing ([#8860](https://github.com/hashicorp/terraform-provider-google/pull/8860))
* workflows: fixed a bug in `google_workflows_workflow` that could cause inconsistent final plan errors when using the `name` field in other resources ([#8869](https://github.com/hashicorp/terraform-provider-google/pull/8869))

## 3.64.0 (April 12, 2021)

FEATURES:
* **New Resource:** `google_tags_tag_key_iam_binding` ([#8844](https://github.com/hashicorp/terraform-provider-google/pull/8844))
* **New Resource:** `google_tags_tag_key_iam_member` ([#8844](https://github.com/hashicorp/terraform-provider-google/pull/8844))
* **New Resource:** `google_tags_tag_key_iam_policy` ([#8844](https://github.com/hashicorp/terraform-provider-google/pull/8844))
* **New Resource:** `google_tags_tag_value_iam_binding` ([#8844](https://github.com/hashicorp/terraform-provider-google/pull/8844))
* **New Resource:** `google_tags_tag_value_iam_member` ([#8844](https://github.com/hashicorp/terraform-provider-google/pull/8844))
* **New Resource:** `google_tags_tag_value_iam_policy` ([#8844](https://github.com/hashicorp/terraform-provider-google/pull/8844))
* **New Resource:** `google_apigee_envgroup_attachment` ([#8853](https://github.com/hashicorp/terraform-provider-google/pull/8853))
* **New Resource:** `google_tags_tag_binding` ([#8854](https://github.com/hashicorp/terraform-provider-google/pull/8854))
* **New Resource:** `google_tags_tag_key` ([#8854](https://github.com/hashicorp/terraform-provider-google/pull/8854))
* **New Resource:** `google_tags_tag_value` ([#8854](https://github.com/hashicorp/terraform-provider-google/pull/8854))

IMPROVEMENTS:
* bigquery: added `require_partition_filter` field to `google_bigquery_table` when provisioning `hive_partitioning_options` ([#8775](https://github.com/hashicorp/terraform-provider-google/pull/8775))
* compute: added field `maintenance_window.start_time` to `google_compute_node_group` ([#8847](https://github.com/hashicorp/terraform-provider-google/pull/8847))
* compute: added gVNIC support for `google_compute_instance_template` ([#8842](https://github.com/hashicorp/terraform-provider-google/pull/8842))
* datacatalog: added `description` field to `google_data_catalog_tag_template ` resource ([#8851](https://github.com/hashicorp/terraform-provider-google/pull/8851))
* iam: added support for third party identities via the principle and principleSet IAM members ([#8860](https://github.com/hashicorp/terraform-provider-google/pull/8860))
* tags: promoted `google_tags_tag_key` to GA ([#8854](https://github.com/hashicorp/terraform-provider-google/pull/8854))
* tags: promoted `google_tags_tag_value` to GA ([#8854](https://github.com/hashicorp/terraform-provider-google/pull/8854))

BUG FIXES:
* compute: reverted datatype change for `mtu` in `google_compute_interconnect_attachment` as it was incompatible with existing state representation ([#8829](https://github.com/hashicorp/terraform-provider-google/pull/8829))
* iam: fixed issue with principle and principleSet members not retaining their casing ([#8860](https://github.com/hashicorp/terraform-provider-google/pull/8860))
* storage: fixed intermittent `Provider produced inconsistent result after apply` error when creating `google_storage_hmac_key` ([#8817](https://github.com/hashicorp/terraform-provider-google/pull/8817))

## 3.63.0 (April 5, 2021)

FEATURES:
* **New Data Source:** `google_monitoring_istio_canonical_service` ([#8789](https://github.com/hashicorp/terraform-provider-google/pull/8789))
* **New Resource:** `google_apigee_instance_attachment` ([#8795](https://github.com/hashicorp/terraform-provider-google/pull/8795))

IMPROVEMENTS:
* added support for Apple silicon chip (updated to go 1.16) ([#8693](https://github.com/hashicorp/terraform-provider-google/pull/8693))
* container: 
  * added support for GKE Autopilot in `google_container_cluster`([#8805](https://github.com/hashicorp/terraform-provider-google/pull/8805))
  * promoted `networking_mode` to GA in `google_container_cluster` ([#8805](https://github.com/hashicorp/terraform-provider-google/pull/8805))
  * added `private_ipv6_google_access` field to `google_container_cluster` ([#8798](https://github.com/hashicorp/terraform-provider-google/pull/8798))
* sql: changed the default timeout of `google_sql_database_instance` to 30m from 20m ([#8802](https://github.com/hashicorp/terraform-provider-google/pull/8802))

BUG FIXES:
* bigquery: fixed issue where you couldn't extend an existing `schema` with additional columns in `google_bigquery_table` ([#8803](https://github.com/hashicorp/terraform-provider-google/pull/8803))
* cloudidentity: modified `google_cloud_identity_groups` and `google_cloud_identity_group_memberships ` to respect the `user_project_override` and `billing_project` configurations and send the appropriate headers to establish a quota project ([#8762](https://github.com/hashicorp/terraform-provider-google/pull/8762))
* compute: added minimum for `scopes` field to `google_compute_instance` resource ([#8801](https://github.com/hashicorp/terraform-provider-google/pull/8801))
* notebooks: fixed permadiff on labels for `google_notebook_instance` ([#8799](https://github.com/hashicorp/terraform-provider-google/pull/8799))
* secretmanager: set required on `secret_data` in `google_secret_manager_secret_version` ([#8797](https://github.com/hashicorp/terraform-provider-google/pull/8797))


## 3.62.0 (March 29, 2021)

FEATURES:
* **New Data Source:** `google_compute_health_check` ([#8725](https://github.com/hashicorp/terraform-provider-google/pull/8725))
* **New Data Source:** `google_kms_secret_asymmetric` ([#8745](https://github.com/hashicorp/terraform-provider-google/pull/8745))
* **New Resource:** `google_data_catalog_tag_template_iam_*` ([#8730](https://github.com/hashicorp/terraform-provider-google/pull/8730))

IMPROVEMENTS:
* accesscontextmanager: added support for ingress and egress policies to `google_access_context_manager_service_perimeter` ([#8723](https://github.com/hashicorp/terraform-provider-google/pull/8723))
* compute: added `proxy_bind` to `google_compute_target_tcp_proxy`, `google_compute_target_http_proxy` and `google_compute_target_https_proxy` ([#8706](https://github.com/hashicorp/terraform-provider-google/pull/8706))

BUG FIXES:
* compute: fixed an issue where exceeding the operation rate limit would fail without retrying ([#8746](https://github.com/hashicorp/terraform-provider-google/pull/8746))
* compute: corrected underlying type to integer for field `mtu` in `google_compute_interconnect_attachment` ([#8744](https://github.com/hashicorp/terraform-provider-google/pull/8744))


## 3.61.0 (March 23, 2021)

IMPROVEMENTS:
* provider: The provider now supports [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation). The federated json credentials must be loaded through the `GOOGLE_APPLICATION_CREDENTIALS` environment variable. ([#8671](https://github.com/hashicorp/terraform-provider-google/issues/8671))
* compute: added `proxy_bind` to `google_compute_target_tcp_proxy`, `google_compute_target_http_proxy` and `google_compute_target_https_proxy` ([#8706](https://github.com/hashicorp/terraform-provider-google/pull/8706))
* compute: changed `google_compute_subnetwork` to accept more values in the `purpose` field ([#8647](https://github.com/hashicorp/terraform-provider-google/pull/8647))
* compute: promoted field compute_instance.scheduling.min_node_cpus and related fields to ga ([#8697](https://github.com/hashicorp/terraform-provider-google/pull/8697))
* dataflow: added `enable_streaming_engine` argument to `google_dataflow_job` ([#8670](https://github.com/hashicorp/terraform-provider-google/pull/8670))
* healthcare: promoted `google_healthcare_consent_store*` to GA support ([#8681](https://github.com/hashicorp/terraform-provider-google/pull/8681))

BUG FIXES:
* container: Fixed updates on `export_custom_routes` and `import_custom_routes` in `google_compute_network_peering` ([#8650](https://github.com/hashicorp/terraform-provider-google/pull/8650))

## 3.60.0 (March 15, 2021)

NOTES: From this release onwards [`google_compute_shared_vpc_service_project`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_shared_vpc_service_project) will not recognise the Shared VPC Admin role when assigned at the folder level in the GA provider, as that functionality is not enabled in the GA API. If you have folder-level IAM configured for Shared VPC permissions, use the `google-beta` provider instead.

FEATURES:
* **New Resource:** google_apigee_envgroup ([#8641](https://github.com/hashicorp/terraform-provider-google/pull/8641))
* **New Resource:** google_apigee_environment ([#8596](https://github.com/hashicorp/terraform-provider-google/pull/8596))

IMPROVEMENTS:
* cloudrun: suppressed metadata.labels["cloud.googleapis.com/location"] value in `google_cloud_run_service` ([#8574](https://github.com/hashicorp/terraform-provider-google/pull/8574))
* compute: added `mtu` field to `google_compute_interconnect_attachment` ([#8575](https://github.com/hashicorp/terraform-provider-google/pull/8575))
* compute: added support for `nic_type` to `google_compute_instance` (GA only) ([#8562](https://github.com/hashicorp/terraform-provider-google/pull/8562))
* datafusion: added support for the `DEVELOPER` instance type to `google_data_fusion_instance`  ([#8590](https://github.com/hashicorp/terraform-provider-google/pull/8590))
* monitoring: added windows based availability sli to the resource `google_monitoring_slo` ([#8588](https://github.com/hashicorp/terraform-provider-google/pull/8588))
* sql: added `settings.0.backup_configuration.transaction_log_retention_days` and `settings.0.backup_configuration.transaction_log_retention_days` fields to `google_sql_database_instance` ([#8582](https://github.com/hashicorp/terraform-provider-google/pull/8582))
* storage: added `kms_key_name` to `google_storage_bucket_object` resource ([#8615](https://github.com/hashicorp/terraform-provider-google/pull/8615))

BUG FIXES:
* bigquery: fixed materialized view to be recreated when query changes ([#8628](https://github.com/hashicorp/terraform-provider-google/pull/8628))
* bigtable: fixed bug where gc_policy would attempt to recreate the resource when switching from deprecated attribute but maintaining the same underlying value ([#8639](https://github.com/hashicorp/terraform-provider-google/pull/8639))
* bigtable: required resource recreation if any fields change on `resource_bigtable_gc_policy` ([#8552](https://github.com/hashicorp/terraform-provider-google/pull/8552))
* binaryauthorization: fixed permadiff in `google_binary_authorization_attestor` ([#8636](https://github.com/hashicorp/terraform-provider-google/pull/8636))
* cloudfunction: added retry logic for `google_cloudfunctions_function` updates ([#8554](https://github.com/hashicorp/terraform-provider-google/pull/8554))
* cloudidentity: fixed a bug where `google_cloud_identity_group` would periodically fail with a 403 ([#8585](https://github.com/hashicorp/terraform-provider-google/pull/8585))
* compute: fixed a perma-diff for `nat_ips` that were specified as short forms in `google_compute_router_nat` ([#8576](https://github.com/hashicorp/terraform-provider-google/pull/8576))
* compute: fixed perma-diff for cos-family disk images ([#8602](https://github.com/hashicorp/terraform-provider-google/pull/8602))
* compute: Fixed service account scope alias to be updated. ([#8604](https://github.com/hashicorp/terraform-provider-google/pull/8604))
* container: fixed container cluster not removed from the state when received 404 error on delete call for the resource `google_container_cluster` ([#8594](https://github.com/hashicorp/terraform-provider-google/pull/8594))
* container: Fixed failure in deleting `maintenance_exclusion` for `google_container_cluster` ([#8589](https://github.com/hashicorp/terraform-provider-google/pull/8589))
* container: fixed an issue where release channel UNSPECIFIED could not be set ([#8595](https://github.com/hashicorp/terraform-provider-google/pull/8595))
* essentialcontacts: made `language_tag` required for `google_essential_contacts_contact` ([#8557](https://github.com/hashicorp/terraform-provider-google/pull/8557))

## 3.59.0 (March 08, 2021)
FEATURES:
* **New Resource:** `google_workflows_workflow` ([#8549](https://github.com/hashicorp/terraform-provider-google/pull/8549))
* **New Resource:** google_apigee_instance ([#8546](https://github.com/hashicorp/terraform-provider-google/pull/8546))

IMPROVEMENTS:
* compute: Added graceful termination to `google_container_node_pool` create calls so that partially created node pools will resume the original operation if the Terraform process is killed mid create. ([#8492](https://github.com/hashicorp/terraform-provider-google/pull/8492))
* compute: Promoted gVNIC support for `google_compute_instance` resource to GA ([#8506](https://github.com/hashicorp/terraform-provider-google/pull/8506))
* compute: added autoscaling_policy.cpu_utilization.predictive_method field to `google_compute_autoscaler` and `google_compute_region_autoscaler` ([#8547](https://github.com/hashicorp/terraform-provider-google/pull/8547))
* redis : marked `auth_string` on the `resource_redis_instance` resource as sensitive ([#8513](https://github.com/hashicorp/terraform-provider-google/pull/8513))

BUG FIXES:
* apigee: fixed IDs when importing `google_apigee_organization` resource ([#8488](https://github.com/hashicorp/terraform-provider-google/pull/8488))
* artifactregistry: fixed issue where updating `google_artifact_registry_repository` always failed ([#8491](https://github.com/hashicorp/terraform-provider-google/pull/8491))
* compute : fixed a bug where `guest_flush` could not be set to false for the resource `google_compute_resource_policy` ([#8517](https://github.com/hashicorp/terraform-provider-google/pull/8517))
* compute: fixed a panic on empty `target_size` in `google_compute_region_instance_group_manager` ([#8528](https://github.com/hashicorp/terraform-provider-google/pull/8528))
* redis: fixed invalid value error on `auth_string` in `google_redis_instance` ([#8493](https://github.com/hashicorp/terraform-provider-google/pull/8493))

## 3.58.0 (February 23, 2021)

NOTES:
* `google_bigquery_table` resources now cannot be destroyed unless `deletion_protection = false` is set in state for the resource. ([#8453](https://github.com/hashicorp/terraform-provider-google/pull/8453))

FEATURES:
* **New Data Source:** `google_iap_client` ([#8450](https://github.com/hashicorp/terraform-provider-google/pull/8450))

IMPROVEMENTS:
* bigquery: added `deletion_protection` field to `google_bigquery_table` to make deleting them require an explicit intent. ([#8453](https://github.com/hashicorp/terraform-provider-google/pull/8453))
* cloudrun: updated retry logic to attempt to retry 409 errors from the Cloud Run API, which may be returned intermittently on create. ([#8440](https://github.com/hashicorp/terraform-provider-google/pull/8440))
* compute: removed max items limit from `google_compute_target_ssl_proxy`. The API currently allows upto 15 Certificates. ([#8478](https://github.com/hashicorp/terraform-provider-google/pull/8478))
* compute: added support for Private Services Connect for Google APIs in `google_compute_global_address` and `google_compute_global_forwarding_rule`([#8458](https://github.com/hashicorp/terraform-provider-google/pull/8458))
* iam: added a retry condition that retries editing `iam_binding` and `iam_member` resources on policies that have frequently deleted service accounts ([#8476](https://github.com/hashicorp/terraform-provider-google/pull/8476))
* sql: added `insights_config` block to `google_sql_database_instance` resource ([#8434](https://github.com/hashicorp/terraform-provider-google/pull/8434))

BUG FIXES:
* compute: fixed an issue where the provider could return an error on a successful delete operation ([#8463](https://github.com/hashicorp/terraform-provider-google/pull/8463))
* dataproc : fixed an issue where `max_failure_per_hour` was not set correctly for `google_dataproc_job` ([#8441](https://github.com/hashicorp/terraform-provider-google/pull/8441))
* dlp : modified `google_data_loss_prevention_stored_info_type` `regex.group_indexes` field to trigger resource recreation on update ([#8439](https://github.com/hashicorp/terraform-provider-google/pull/8439))
* sql: fixed diffs based on case for `charset` in `google_sql_database` ([#8462](https://github.com/hashicorp/terraform-provider-google/pull/8462))

## 3.57.0 (February 16, 2021)

DEPRECATIONS:
* compute: deprecated `source_disk_url` field in `google_compute_snapshot`. ([#8410](https://github.com/hashicorp/terraform-provider-google/pull/8410))
* kms: deprecated `self_link` field in `google_kms_keyring` and `google_kms_cryptokey` resource as it is identical value to `id` field. ([#8410](https://github.com/hashicorp/terraform-provider-google/pull/8410))
* pubsub: deprecated `path` field in `google_pubsub_subscription` resource as it is identical value to `id` field. ([#8410](https://github.com/hashicorp/terraform-provider-google/pull/8410))

FEATURES:
* **New Resource:** `google_essential_contacts_contact` ([#8426](https://github.com/hashicorp/terraform-provider-google/pull/8426))

IMPROVEMENTS:
* bigquery: added `status` field to `google_bigquery_job` ([#8377](https://github.com/hashicorp/terraform-provider-google/pull/8377))
* compute: added `disk.resource_policies` field to resource `google_compute_instance_template` ([#8393](https://github.com/hashicorp/terraform-provider-google/pull/8393))
* pubsub: marked `kms_key_name` field in `google_pubsub_topic` as updatable ([#8424](https://github.com/hashicorp/terraform-provider-google/pull/8424))

BUG FIXES:
* appengine: added retry for P4SA propagation delay ([#8409](https://github.com/hashicorp/terraform-provider-google/pull/8409))
* compute: fixed overly-aggressive detection of changes to google_compute_security_policy rules ([#8417](https://github.com/hashicorp/terraform-provider-google/pull/8417))

## 3.56.0 (February 8, 2021)

FEATURES:
* **New Resource:** `google_privateca_certificate` ([#8371](https://github.com/hashicorp/terraform-provider-google/pull/8371))

IMPROVEMENTS:
* all: added plan time validations for fields that expect base64 values. ([#8304](https://github.com/hashicorp/terraform-provider-google/pull/8304))
* sql: added support for point-in-time-recovery to `google_sql_database_instance` ([#8367](https://github.com/hashicorp/terraform-provider-google/pull/8367))
* monitoring : added `availability` sli metric support for the resource `google_monitoring_slo` ([#8315](https://github.com/hashicorp/terraform-provider-google/pull/8315))

BUG FIXES:
* bigquery: fixed bug where you could not reorder columns on `schema` for resource `google_bigquery_table` ([#8321](https://github.com/hashicorp/terraform-provider-google/pull/8321))
* cloudrun: suppressed `run.googleapis.com/ingress-status` annotation in `google_cloud_run_service` ([#8361](https://github.com/hashicorp/terraform-provider-google/pull/8361))
* serviceaccount: loosened restrictions on `account_id` for datasource `google_service_account` ([#8344](https://github.com/hashicorp/terraform-provider-google/pull/8344))

## 3.55.0 (February 1, 2021)

BREAKING CHANGES:
* Reverted `* bigquery: made incompatible changes to the `google_bigquery_table.schema` field to cause the resource to be recreated ([#8232](https://github.com/hashicorp/terraform-provider-google/pull/8232))` due to unintended interactions with a bug introduced in an earlier version of the resource.

FEATURES:
* **New Data Source:** `google_runtimeconfig_config` ([#8268](https://github.com/hashicorp/terraform-provider-google/pull/8268))

IMPROVEMENTS:
* compute: added `distribution_policy_target_shape` field to `google_compute_region_instance_group_manager` resource ([#8277](https://github.com/hashicorp/terraform-provider-google/pull/8277))
* container: promoted `master_global_access_config`, `tpu_ipv4_cidr_block`, `default_snat_status` and `datapath_provider` fields of `google_container_cluster` to GA. ([#8303](https://github.com/hashicorp/terraform-provider-google/pull/8303))
* dataproc: Added field `temp_bucket` to `google_dataproc_cluster` cluster config. ([#8131](https://github.com/hashicorp/terraform-provider-google/pull/8131))
* notebooks: added `tags`, `service_account_scopes`,`shielded_instance_config` to `google_notebooks_instance` ([#8289](https://github.com/hashicorp/terraform-provider-google/pull/8289))
* provider: added plan time validations for fields that expect base64 values. ([#8304](https://github.com/hashicorp/terraform-provider-google/pull/8304))

BUG FIXES:
* bigquery: fixed permadiff on expiration_ms for `google_bigquery_table` ([#8298](https://github.com/hashicorp/terraform-provider-google/pull/8298))
* billing: fixed perma-diff on currency_code in `google_billing_budget` ([#8266](https://github.com/hashicorp/terraform-provider-google/pull/8266))
 * compute: changed private_ipv6_google_access in `google_compute_subnetwork` to correctly send a fingerprint ([#8290](https://github.com/hashicorp/terraform-provider-google/pull/8290))
* healthcare: add retry logic on healthcare dataset not initialized error ([#8256](https://github.com/hashicorp/terraform-provider-google/pull/8256))

## 3.54.0 (January 25, 2021)

KNOWN ISSUES: New `google_bigquery_table` behaviour introduced in this version had unintended consequences, and may incorrectly flag tables for recreation. We expect to revert this for `3.55.0`.

FEATURES:
* **New Data Source:** `google_cloud_run_locations` ([#8192](https://github.com/hashicorp/terraform-provider-google/pull/8192))
* **New Resource:** `google_privateca_certificate_authority` ([#8233](https://github.com/hashicorp/terraform-provider-google/pull/8233))
* **New Resource:** `google_privateca_certificate_authority_iam_binding` ([#8249](https://github.com/hashicorp/terraform-provider-google/pull/8249))
* **New Resource:** `google_privateca_certificate_authority_iam_member` ([#8249](https://github.com/hashicorp/terraform-provider-google/pull/8249))
* **New Resource:** `google_privateca_certificate_authority_iam_policy` ([#8249](https://github.com/hashicorp/terraform-provider-google/pull/8249))

IMPROVEMENTS:
* bigquery: made incompatible changes to the `google_bigquery_table.schema` field to cause the resource to be recreated ([#8232](https://github.com/hashicorp/terraform-provider-google/pull/8232))
* bigtable: fixed an issue where the `google_bigtable_instance` resource was not inferring the zone from the provider. ([#8222](https://github.com/hashicorp/terraform-provider-google/pull/8222))
* cloudscheduler: fixed unnecessary recreate for `google_cloud_scheduler_job` ([#8248](https://github.com/hashicorp/terraform-provider-google/pull/8248))
* compute: added `scaling_schedules` fields to `google_compute_autoscaler` and `google_compute_region_autoscaler` (beta) ([#8245](https://github.com/hashicorp/terraform-provider-google/pull/8245))
* compute: fixed an issue where `google_compute_region_per_instance_config`, `google_compute_per_instance_config`, `google_compute_region_instance_group_manager` resources were not inferring the region/zone from the provider. ([#8224](https://github.com/hashicorp/terraform-provider-google/pull/8224))
* memcache: fixed an issue where `google_memcached_instance` resource was not inferring the region from the provider. ([#8188](https://github.com/hashicorp/terraform-provider-google/pull/8188))
* tpu: fixed an issue where `google_tpu_node` resource was not inferring the zone from the provider. ([#8188](https://github.com/hashicorp/terraform-provider-google/pull/8188))
* vpcaccess: fixed an issue where `google_vpc_access_connector` resource was not inferring the region from the provider. ([#8188](https://github.com/hashicorp/terraform-provider-google/pull/8188))

BUG FIXES:
* bigquery: fixed an issue in `bigquery_dataset_iam_member` where deleted members were not handled correctly ([#8231](https://github.com/hashicorp/terraform-provider-google/pull/8231))
* compute: fixed a perma-diff on `google_compute_health_check` when `log_config.enable` is set to false ([#8209](https://github.com/hashicorp/terraform-provider-google/pull/8209))
* notebooks: fixed permadiff on noRemoveDataDisk for `google_notebooks_instance` ([#8246](https://github.com/hashicorp/terraform-provider-google/pull/8246))
* resourcemanager: fixed an inconsistent result when IAM conditions are specified with `google_folder_iam_*` ([#8235](https://github.com/hashicorp/terraform-provider-google/pull/8235))
* healthcare: added retry logic on healthcare dataset not initialized error ([#8256](https://github.com/hashicorp/terraform-provider-google/pull/8256))

## 3.53.0 (January 19, 2021)

FEATURES:
* **New Data Source:** `google_compute_instance_template` ([#8137](https://github.com/hashicorp/terraform-provider-google/pull/8137))
* **New Resource:** `google_apigee_organization` ([#8178](https://github.com/hashicorp/terraform-provider-google/pull/8178))

IMPROVEMENTS:
* accesscontextmanager: added support for `google_access_context_manager_gcp_user_access_binding` ([#8168](https://github.com/hashicorp/terraform-provider-google/pull/8168))
* cloudbuild: promoted `github` fields in `google_cloud_build_trigger` to GA ([#8167](https://github.com/hashicorp/terraform-provider-google/pull/8167))
* memcached: fixed an issue where `google_memcached_instance` resource was not inferring the region from the provider. ([More info](https://github.com/hashicorp/terraform-provider-google/issues/8027))
* serviceaccount: added a `keepers` field to `google_service_account_key` that recreates the field when it is modified ([#8097](https://github.com/hashicorp/terraform-provider-google/pull/8097))
* sql: added restore from backup support to `google_sql_database_instance` ([#8138](https://github.com/hashicorp/terraform-provider-google/pull/8138))
* sql: added support for MYSQL_8_0 on resource `google_sql_source_representation_instance` ([#8135](https://github.com/hashicorp/terraform-provider-google/pull/8135))
* tpu: fixed an issue where `google_tpu_node` resource was not inferring the zone from the provider. ([More info](https://github.com/hashicorp/terraform-provider-google/issues/8027))
* vpcaccess: fixed an issue where `google_vpc_access_connector` resource was not inferring the region from the provider. ([More info](https://github.com/hashicorp/terraform-provider-google/issues/8027))

BUG FIXES:
* bigquery: enhanced diff suppress to ignore certain api divergences on resource `table` ([#8134](https://github.com/hashicorp/terraform-provider-google/pull/8134))
* container: fixed crash due to nil exclusions object when updating an existent cluster with maintenance_policy but without exclusions ([#8126](https://github.com/hashicorp/terraform-provider-google/pull/8126))
* project: fixed a bug in `google_project_access_approval_settings` where the default `project` was used rather than `project_id` ([#8169](https://github.com/hashicorp/terraform-provider-google/pull/8169))

## 3.52.0 (January 11, 2021)

BREAKING CHANGES:
* billing: removed import support for `google_billing_budget` as it never functioned correctly ([#8023](https://github.com/hashicorp/terraform-provider-google/pull/8023))

FEATURES:
* **New Data Source:** `google_sql_backup_run` ([#8100](https://github.com/hashicorp/terraform-provider-google/pull/8100))
* **New Data Source:** `google_storage_bucket_object_content` ([#8016](https://github.com/hashicorp/terraform-provider-google/pull/8016))
* **New Resource:** `google_billing_subaccount` ([#8022](https://github.com/hashicorp/terraform-provider-google/pull/8022))
* **New Resource:** `google_pubsub_lite_subscription` ([#8011](https://github.com/hashicorp/terraform-provider-google/pull/8011))
* **New Resource:** `google_pubsub_lite_topic` ([#8011](https://github.com/hashicorp/terraform-provider-google/pull/8011))

IMPROVEMENTS:
* bigquery: promoted bigquery reservation to GA. ([#8079](https://github.com/hashicorp/terraform-provider-google/pull/8079))
* bigtable: added support for specifying `duration` for `bigtable_gc_policy` to allow durations shorter than a day ([#7879](https://github.com/hashicorp/terraform-provider-google/pull/7879))
* bigtable: added support for specifying `duration` for `bigtable_gc_policy` to allow durations shorter than a day ([#8081](https://github.com/hashicorp/terraform-provider-google/pull/8081))
* billing: promoted `google_billing_budget` to GA ([#8023](https://github.com/hashicorp/terraform-provider-google/pull/8023))
* compute: Added support for Google Virtual Network Interface (gVNIC) for `google_compute_image` ([#8007](https://github.com/hashicorp/terraform-provider-google/pull/8007))
* compute: added SHARED_LOADBALANCER_VIP as a valid option for `google_compute_address.purpose` ([#7987](https://github.com/hashicorp/terraform-provider-google/pull/7987))
* compute: added field `multiwriter` to resource `disk` (beta) ([#8098](https://github.com/hashicorp/terraform-provider-google/pull/8098))
* compute: added support for `enable_independent_endpoint_mapping` to `google_compute_router_nat` resource ([#8049](https://github.com/hashicorp/terraform-provider-google/pull/8049))
* compute: added support for `filter.direction` to `google_compute_packet_mirroring` ([#8102](https://github.com/hashicorp/terraform-provider-google/pull/8102))
* compute: promoted `confidential_instance_config` field in `google_compute_instance` and `google_compute_instance_template` to GA ([#8089](https://github.com/hashicorp/terraform-provider-google/pull/8089))
* compute: promoted `google_compute_forwarding_rule` `is_mirroring_collector` to GA ([#8102](https://github.com/hashicorp/terraform-provider-google/pull/8102))
* compute: promoted `google_compute_packet_mirroring` to GA ([#8102](https://github.com/hashicorp/terraform-provider-google/pull/8102))
* dataflow: Added optional `kms_key_name` field for `google_dataflow_job` ([#8116](https://github.com/hashicorp/terraform-provider-google/pull/8116))
* dataflow: added documentation about using `parameters` for custom service account and other pipeline options to `google_dataflow_flex_template_job` ([#7999](https://github.com/hashicorp/terraform-provider-google/pull/7999))
* redis: added `auth_string` output to `google_redis_instance` when `auth_enabled` is `true` ([#8090](https://github.com/hashicorp/terraform-provider-google/pull/8090))
* redis: promoted `google_redis_instance.auth_enabled` to GA ([#8090](https://github.com/hashicorp/terraform-provider-google/pull/8090))
* sql: added support for setting the `type` field on `google_sql_user` to support IAM authentication ([#8017](https://github.com/hashicorp/terraform-provider-google/pull/8017))
* sql: added support for setting the `type` field on `google_sql_user` to support IAM authentication ([#8047](https://github.com/hashicorp/terraform-provider-google/pull/8047))

BUG FIXES:
* compute: removed requirement for `google_compute_region_url_map` default_service, as it should be a choice of default_service or default_url_redirect ([#2810](https://github.com/hashicorp/terraform-provider-google-beta/pull/2810))
* cloud_tasks: fixed permadiff on retry_config.max_retry_duration for `google_cloud_tasks_queue` when the 0s is supplied ([#8078](https://github.com/hashicorp/terraform-provider-google/pull/8078))
* cloudfunctions: fixed a bug where `google_cloudfunctions_function` would sometimes fail to update after being imported from gcloud ([#8010](https://github.com/hashicorp/terraform-provider-google/pull/8010))
* cloudrun: fixed a permanent diff on `google_cloud_run_domain_mapping` `spec.force_override` field ([#8026](https://github.com/hashicorp/terraform-provider-google/pull/8026))
* container: added plan time validation to ensure `enable_private_nodes` is true if `master_ipv4_cidr_block` is set on resource `cluster` ([#8066](https://github.com/hashicorp/terraform-provider-google/pull/8066))
* dataproc: updated jobs to no longer wait for job completion during create ([#8064](https://github.com/hashicorp/terraform-provider-google/pull/8064))
* filestore: updated retry logic to fail fast on quota error which cannot succeed on retry. ([#8080](https://github.com/hashicorp/terraform-provider-google/pull/8080))
* logging: fixed updating on disabled in `google_logging_project_sink` ([#8093](https://github.com/hashicorp/terraform-provider-google/pull/8093))
* scheduler: Fixed syntax error in the Cloud Scheduler HTTP target example. ([#8004](https://github.com/hashicorp/terraform-provider-google/pull/8004))
* sql: fixed a bug in `google_sql_database_instance` that caused a permadiff on `settings.replication_type` ([#8006](https://github.com/hashicorp/terraform-provider-google/pull/8006))
* storage: updated IAM resources to refresh etag sooner on an IAM conflict error, which will make applications of multiple IAM resources much faster. ([#8080](https://github.com/hashicorp/terraform-provider-google/pull/8080))

## 3.51.1 (January 07, 2021)

BUG FIXES:
* all: fixed a bug that would occur in various resources due to comparison of large integers ([#8103](https://github.com/hashicorp/terraform-provider-google/pull/8103))

## 3.51.0 (December 14, 2020)

FEATURES:
* **New Resource:** `google_firestore_document` ([#7932](https://github.com/hashicorp/terraform-provider-google/pull/7932))
* **New Resource:** `google_notebooks_instance` ([#7933](https://github.com/hashicorp/terraform-provider-google/pull/7933))
* **New Resource:** `google_notebooks_environment` ([#7933](https://github.com/hashicorp/terraform-provider-google/pull/7933))

IMPROVEMENTS:
* compute: added CDN features to `google_compute_region_backend_service`. ([#7941](https://github.com/hashicorp/terraform-provider-google/pull/7941))
* compute: added Flexible Cache Control features to `google_compute_backend_service`. ([#7941](https://github.com/hashicorp/terraform-provider-google/pull/7941))
* compute: added `replacement_method` field to `update_policy` block of `google_compute_instance_group_manager` ([#7918](https://github.com/hashicorp/terraform-provider-google/pull/7918))
* compute: added `replacement_method` field to `update_policy` block of `google_compute_region_instance_group_manager` ([#7918](https://github.com/hashicorp/terraform-provider-google/pull/7918))
* logging: added plan time validation for `unique_writer_identity` on `google_logging_project_sink` ([#7974](https://github.com/hashicorp/terraform-provider-google/pull/7974))
* storage: added more lifecycle conditions to `google_storage_bucket` resource ([#7937](https://github.com/hashicorp/terraform-provider-google/pull/7937))

BUG FIXES:
* all: bump default request timeout to avoid conflicts if creating a resource takes longer than expected ([#7976](https://github.com/hashicorp/terraform-provider-google/pull/7976))
* compute: removed `custom_response_headers` from GA `google_compute_backend_service` since it only works in the Beta version ([#7943](https://github.com/hashicorp/terraform-provider-google/pull/7943))
* project: fixed a bug where `google_project_default_service_accounts` would delete all IAM bindings on a project when run with `action = "DEPRIVILEGE"` ([#7984](https://github.com/hashicorp/terraform-provider-google/pull/7984))
* spanner: fixed an issue in `google_spanner_database` where multi-statement updates were not formatted correctly ([#7970](https://github.com/hashicorp/terraform-provider-google/pull/7970))
* sql: fixed a bug in `google_sql_database_instance` that caused a permadiff on `settings.replication_type` ([#8006](https://github.com/hashicorp/terraform-provider-google/pull/8006))

## 3.50.0 (December 7, 2020)

FEATURES:
* **New Data Source:** `google_composer_environment` ([#7902](https://github.com/hashicorp/terraform-provider-google/pull/7902))
* **New Data Source:** `google_monitoring_cluster_istio_service` ([#7847](https://github.com/hashicorp/terraform-provider-google/pull/7847))
* **New Data Source:** `google_monitoring_mesh_istio_service` ([#7847](https://github.com/hashicorp/terraform-provider-google/pull/7847))

IMPROVEMENTS:
* compute: added `replacement_method` field to `update_policy` block of `google_compute_instance_group_manager` ([#7918](https://github.com/hashicorp/terraform-provider-google/pull/7918))
* compute: added `replacement_method` field to `update_policy` block of `google_compute_region_instance_group_manager` ([#7918](https://github.com/hashicorp/terraform-provider-google/pull/7918))
* compute: added more fields to cdn_policy block of `google_compute_backend_bucket` ([#7888](https://github.com/hashicorp/terraform-provider-google/pull/7888))
* compute: promoted `google_compute_managed_ssl_certificate` to GA ([#7914](https://github.com/hashicorp/terraform-provider-google/pull/7914))
* compute: promoted `google_compute_resource_policy` to GA ([#7917](https://github.com/hashicorp/terraform-provider-google/pull/7917))
* compute: updated `google_compute_url_map`'s fields referring to backend services to be able to refer to backend buckets. ([#7916](https://github.com/hashicorp/terraform-provider-google/pull/7916))
* container: added cluster state check before proceeding on the node pool activities ([#7887](https://github.com/hashicorp/terraform-provider-google/pull/7887))
* google: added support for more import formats to google_project_iam_custom_role ([#7862](https://github.com/hashicorp/terraform-provider-google/pull/7862))
* project: added new restore_policy `REVERT_AND_IGNORE_FAILURE` to `google_project_default_service_accounts` ([#7906](https://github.com/hashicorp/terraform-provider-google/pull/7906))

BUG FIXES:
* bigqueryconnection: fixed failure to import a resource if it has a non-default project or location. ([#7903](https://github.com/hashicorp/terraform-provider-google/pull/7903))
* iam: fixed iam conflict handling so that optimistic-locking retries will succeed more often. ([#7915](https://github.com/hashicorp/terraform-provider-google/pull/7915))
* storage: fixed an issue in `google_storage_bucket` where `cors` could not be removed ([#7858](https://github.com/hashicorp/terraform-provider-google/pull/7858))

## 3.49.0 (November 23, 2020)

FEATURES:
* **New Resource:** google_healthcare_consent_store ([#7803](https://github.com/hashicorp/terraform-provider-google/pull/7803))
* **New Resource:** google_healthcare_consent_store_iam_binding ([#7803](https://github.com/hashicorp/terraform-provider-google/pull/7803))
* **New Resource:** google_healthcare_consent_store_iam_member ([#7803](https://github.com/hashicorp/terraform-provider-google/pull/7803))
* **New Resource:** google_healthcare_consent_store_iam_policy ([#7803](https://github.com/hashicorp/terraform-provider-google/pull/7803))

IMPROVEMENTS:
* bigquery: added `ORC` as a valid option to `source_format` field of  `google_bigquery_table` resource ([#7804](https://github.com/hashicorp/terraform-provider-google/pull/7804))
* cloud_identity: promoted `google_cloud_identity_group_membership` to GA ([#7786](https://github.com/hashicorp/terraform-provider-google/pull/7786))
* cloud_identity: promoted `google_cloud_identity_group` to GA ([#7786](https://github.com/hashicorp/terraform-provider-google/pull/7786))
* cloud_identity: promoted data source `google_cloud_identity_group_memberships` to GA ([#7786](https://github.com/hashicorp/terraform-provider-google/pull/7786))
* cloud_identity: promoted data source `google_cloud_identity_groups` to GA ([#7786](https://github.com/hashicorp/terraform-provider-google/pull/7786))
* compute: added `custom_response_headers` field to `google_compute_backend_service` resource ([#7824](https://github.com/hashicorp/terraform-provider-google/pull/7824))
* container: added maintenance_exclusions_window to `google_container_cluster` ([#7830](https://github.com/hashicorp/terraform-provider-google/pull/7830))
* logging: added description and disabled to logging sinks ([#7809](https://github.com/hashicorp/terraform-provider-google/pull/7809))
* runtimeconfig: marked value and text fields in `google_runtimeconfig_variable` resource as sensitive ([#7808](https://github.com/hashicorp/terraform-provider-google/pull/7808))
* sql: added `deletion_policy` field to `google_sql_user` to enable abandoning users rather than deleting them ([#7820](https://github.com/hashicorp/terraform-provider-google/pull/7820))

BUG FIXES:
* bigtable: added ignore_warnings flag to create call for `google_bigtable_app_profile` ([#7806](https://github.com/hashicorp/terraform-provider-google/pull/7806))

## 3.48.0 (November 16, 2020)

FEATURES:
* **New Data Source:** `google_iam_workload_identity_pool_provider` ([#7733](https://github.com/hashicorp/terraform-provider-google/pull/7733))

IMPROVEMENTS:
* apigateway: added api_config_id_prefix field to `google_api_gateway_api_config` resoure ([#7753](https://github.com/hashicorp/terraform-provider-google/pull/7753))
* cloudfunctions: fixed a bug with `google_cloudfunction_function` that blocked updates when Organization Policies are enabled. ([#7723](https://github.com/hashicorp/terraform-provider-google/pull/7723))
* compute: added `autoscaling_policy.0.scale_in_control` fields to `google_compute_autoscaler` ([#7773](https://github.com/hashicorp/terraform-provider-google/pull/7773))
* compute: added `autoscaling_policy.0.scale_in_control` fields to `google_compute_region_autoscaler` ([#7773](https://github.com/hashicorp/terraform-provider-google/pull/7773))
* compute: added update support for `google_compute_interconnect_attachment` `bandwidth` field ([#7762](https://github.com/hashicorp/terraform-provider-google/pull/7762))
* dataproc: added "FLINK", "DOCKER", "HBASE" as valid options for field `cluster_config.0.software_config.0.optional_components` of `google_dataproc_cluster` resource ([#7726](https://github.com/hashicorp/terraform-provider-google/pull/7726))

BUG FIXES:
* cloudrun: added diff suppress function for `google_cloud_run_domain_mapping` `metadata.annotations` to ignore API-set fields ([#7764](https://github.com/hashicorp/terraform-provider-google/pull/7764))
* spanner: marked `google_spanner_instance.config` as ForceNew as is not updatable ([#7763](https://github.com/hashicorp/terraform-provider-google/pull/7763))

## 3.47.0 (November 09, 2020)

FEATURES:
* **New Data Source:** `google_iam_workload_identity_pool` ([#7704](https://github.com/hashicorp/terraform-provider-google/pull/7704))
* **New Resource:** `google_iam_workload_identity_pool_provider` ([#7712](https://github.com/hashicorp/terraform-provider-google/pull/7712))
* **New Resource:** `google_project_default_service_accounts` ([#7709](https://github.com/hashicorp/terraform-provider-google/pull/7709))

IMPROVEMENTS:
* cloudfunctions: fixed a bug with `google_cloudfunction_function` that blocked updates when Organization Policies are enabled. ([#7723](https://github.com/hashicorp/terraform-provider-google/pull/7723))
* functions: added 4096 as a valid value for available_memory_mb field of `google_cloudfunction_function` ([#7707](https://github.com/hashicorp/terraform-provider-google/pull/7707))
* cloudrun: patched `google_cloud_run_service` to suppress Google generated annotations ([#7721](https://github.com/hashicorp/terraform-provider-google/pull/7721))

BUG FIXES:
* dataflow: removed required validation for zone for `google_data_flow_job` when region is given in the config ([#7703](https://github.com/hashicorp/terraform-provider-google/pull/7703))
* monitoring: Fixed type of `google_monitoring_slo`'s `range` values - some `range` values are doubles, others are integers. ([#7676](https://github.com/hashicorp/terraform-provider-google/pull/7676))
* pubsub: Fixed permadiff on push_config.attributes. ([#7714](https://github.com/hashicorp/terraform-provider-google/pull/7714))
* storage: fixed an issue in `google_storage_bucket` where `lifecycle_rules` were always included in update requests ([#7727](https://github.com/hashicorp/terraform-provider-google/pull/7727))

## 3.46.0 (November 02, 2020)
NOTES:
* compute: updated `google_compute_machine_image` resource to complete once the Image is ready. ([#7629](https://github.com/hashicorp/terraform-provider-google/pull/7629))

FEATURES:
* **New Resource:** `google_api_gateway_api_config_iam_binding` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api_config_iam_member` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api_config_iam_policy` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api_config` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api_iam_binding` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api_iam_member` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api_iam_policy` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_api` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_gateway_iam_binding` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_gateway_iam_member` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_gateway_iam_policy` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_api_gateway_gateway` ([#7626](https://github.com/hashicorp/terraform-provider-google/pull/7626))
* **New Resource:** `google_compute_instance_from_machine_image` ([#7629](https://github.com/hashicorp/terraform-provider-google/pull/7629))
* **New Resource:** `google_compute_machine_image_iam_binding` ([#7629](https://github.com/hashicorp/terraform-provider-google/pull/7629))
* **New Resource:** `google_compute_machine_image_iam_member` ([#7629](https://github.com/hashicorp/terraform-provider-google/pull/7629))
* **New Resource:** `google_compute_machine_image_iam_policy` ([#7629](https://github.com/hashicorp/terraform-provider-google/pull/7629))
* **New Resource:** `google_iap_tunnel_iam_binding` ([#7635](https://github.com/hashicorp/terraform-provider-google/pull/7635))
* **New Resource:** `google_iap_tunnel_iam_member` ([#7635](https://github.com/hashicorp/terraform-provider-google/pull/7635))
* **New Resource:** `google_iap_tunnel_iam_policy` ([#7635](https://github.com/hashicorp/terraform-provider-google/pull/7635))
* **New Resource:** compute: promoted `google_compute_region_network_endpoint_group` to GA ([#7618](https://github.com/hashicorp/terraform-provider-google/pull/7618))

IMPROVEMENTS:
* asset: added conditions to Cloud Asset Feeds ([#7632](https://github.com/hashicorp/terraform-provider-google/pull/7632))
* bigquery: added `email_preferences ` field to `google_bigquery_data_transfer_config` resource ([#7665](https://github.com/hashicorp/terraform-provider-google/pull/7665))
* bigquery: added `schedule_options` field to `google_bigquery_data_transfer_config` resource ([#7633](https://github.com/hashicorp/terraform-provider-google/pull/7633))
* compute: added `private_ipv6_google_access` field to `google_compute_subnetwork` ([#7651](https://github.com/hashicorp/terraform-provider-google/pull/7651))
* compute: added storage_locations & cmek fields to `google_compute_machine_image` resource ([#7629](https://github.com/hashicorp/terraform-provider-google/pull/7629))
* compute: added support for non-destructive updates to `export_custom_routes` and `import_custom_routes` for `google_compute_network_peering` ([#7619](https://github.com/hashicorp/terraform-provider-google/pull/7619))
* compute: relax `load_balancing_scheme` validation of `google_compute_region_backend_service` to support external network load-balancers ([#7592](https://github.com/hashicorp/terraform-provider-google/pull/7592))
* datacatalog: Add taxonomy and policy_tag to `google_data_catalog` ([#7588](https://github.com/hashicorp/terraform-provider-google/pull/7588))
* dlp: added `custom_info_types` to `google_dlp_inspect_template` ([#7650](https://github.com/hashicorp/terraform-provider-google/pull/7650))
* functions: added `build_environment_variables` field to `google_cloudfunction_function` ([#7596](https://github.com/hashicorp/terraform-provider-google/pull/7596))
* kms: added `skip_initial_version_creation` to `google_kms_crypto_key` ([#7647](https://github.com/hashicorp/terraform-provider-google/pull/7647))
* monitoring: Added Monitoring Query Language based alerting for `google_monitoring_alert_policy` ([#7664](https://github.com/hashicorp/terraform-provider-google/pull/7664))

BUG FIXES:
* compute: fixed an issue where `google_compute_health_check` `port` values caused a diff when `port_specification` was unset or set to `""` ([#7623](https://github.com/hashicorp/terraform-provider-google/pull/7623))
* monitoring: added more retries for potential failed monitoring operations ([#7631](https://github.com/hashicorp/terraform-provider-google/pull/7631))
* osconfig: fixed an issue where the `rollout.disruption_budget.percentage` field in `google_os_config_patch_deployment` did not correspond to a field in the API ([#7641](https://github.com/hashicorp/terraform-provider-google/pull/7641))
* sql: fixed a case in `google_sql_database_instance` where we inadvertently required the `projects.get` permission for a service networking precheck introduced in `v3.44.0` ([#7622](https://github.com/hashicorp/terraform-provider-google/pull/7622))

## 3.45.0 (October 28, 2020)

BREAKING CHANGES:
* pubsub: changing the value of `google_pubsub_subscription.enable_message_ordering` will now recreate the resource. Previously, an error was returned. ([#7584](https://github.com/hashicorp/terraform-provider-google/pull/7584))
* spanner: `google_spanner_database` resources now cannot be destroyed unless `deletion_protection = false` is set in state for the resource. ([#7557](https://github.com/hashicorp/terraform-provider-google/pull/7557))

NOTES:
* compute: added a warning to `google_compute_vpn_gateway` ([#7547](https://github.com/hashicorp/terraform-provider-google/pull/7547))

FEATURES:
* **New Data Source:** `google_spanner_instance` ([#7537](https://github.com/hashicorp/terraform-provider-google/pull/7537))
* **New Resource:** `access_context_manager_access_level_condition` ([#7524](https://github.com/hashicorp/terraform-provider-google/pull/7524))
* **New Resource:** `google_bigquery_routine` ([#7579](https://github.com/hashicorp/terraform-provider-google/pull/7579))

IMPROVEMENTS:
* billing_budget: added `disable_default_iam_recipients ` field to `google_billing_budget` to allow disable sending email notifications to default recipients. ([#7544](https://github.com/hashicorp/terraform-provider-google/pull/7544))
* compute: added `interface` attribute to `google_compute_disk` ([#7554](https://github.com/hashicorp/terraform-provider-google/pull/7554))
* compute: added support for updating `network_interface.[d].network_ip` on `google_compute_instance` when changing network or subnetwork ([#7515](https://github.com/hashicorp/terraform-provider-google/pull/7515))
* compute: added `mtu` field to `google_compute_network` resource ([#7567](https://github.com/hashicorp/terraform-provider-google/pull/7567))
* compute: promoted HA VPN fields in `google_compute_vpn_tunnel` to GA ([#7547](https://github.com/hashicorp/terraform-provider-google/pull/7547))
* compute: promoted `google_compute_external_vpn_gateway` to GA ([#7547](https://github.com/hashicorp/terraform-provider-google/pull/7547))
* compute: promoted `google_compute_ha_vpn_gateway` to GA ([#7547](https://github.com/hashicorp/terraform-provider-google/pull/7547))
* provider: added support for service account impersonation. ([#7542](https://github.com/hashicorp/terraform-provider-google/pull/7542))
* spanner: added `deletion_protection` field to `google_spanner_database` to make deleting them require an explicit intent. ([#7557](https://github.com/hashicorp/terraform-provider-google/pull/7557))

BUG FIXES:
* all: fixed misleading "empty non-retryable error" message that was appearing in debug logs ([#7569](https://github.com/hashicorp/terraform-provider-google/pull/7569))
* compute: fixed incorrect import format for `google_compute_global_network_endpoint` ([#7523](https://github.com/hashicorp/terraform-provider-google/pull/7523))
* compute: fixed issue where `google_compute_[region_]backend_service.backend.max_utilization` could not be updated ([#7575](https://github.com/hashicorp/terraform-provider-google/pull/7575))
* iap: fixed an eventual consistency bug causing creates for `google_iap_brand` to fail ([#7520](https://github.com/hashicorp/terraform-provider-google/pull/7520))
* provider: fixed an issue where the request headers would grow proportionally to the number of resources in a given `terraform apply` ([#7576](https://github.com/hashicorp/terraform-provider-google/pull/7576))
* serviceusage: fixed bug where concurrent activations/deactivations of project services would fail, now they retry ([#7519](https://github.com/hashicorp/terraform-provider-google/pull/7519))

## 3.44.0 (October 19, 2020)

BREAKING CHANGE:
* Added `deletion_protection` to `google_sql_database_instance`, which defaults to true. SQL instances can no longer be destroyed without setting `deletion_protection = false`. ([#7499](https://github.com/hashicorp/terraform-provider-google/pull/7499))

FEATURES:
* **New Data Source:** `google_app_engine_default_service_account` ([#7472](https://github.com/hashicorp/terraform-provider-google/pull/7472))
* **New Data Source:** `google_pubsub_topic` ([#7448](https://github.com/hashicorp/terraform-provider-google/pull/7448))

IMPROVEMENTS:
* bigquery: added ability for `google_bigquery_dataset_access` to retry quota errors since quota refreshes quickly. ([#7507](https://github.com/hashicorp/terraform-provider-google/pull/7507))
* bigquery: added `MONTH` and `YEAR` as allowed values in `google_bigquery_table.time_partitioning.type` ([#7461](https://github.com/hashicorp/terraform-provider-google/pull/7461))
* cloud_tasks: added `stackdriver_logging_config` field to `cloud_tasks_queue` resource ([#7487](https://github.com/hashicorp/terraform-provider-google/pull/7487))
* compute: added support for updating `network_interface.[d].network_ip` on `google_compute_instance` when changing network or subnetwork ([#7515](https://github.com/hashicorp/terraform-provider-google/pull/7515))
* compute: added `maintenance_policy` field to `google_compute_node_group` ([#7510](https://github.com/hashicorp/terraform-provider-google/pull/7510))
* compute: added filter field to google_compute_image datasource ([#7488](https://github.com/hashicorp/terraform-provider-google/pull/7488))
* compute: promoted `autoscaling_policy` field in `google_compute_node_group` to GA ([#7510](https://github.com/hashicorp/terraform-provider-google/pull/7510))
* dataproc: Added `graceful_decomissioning_timeout` field to `dataproc_cluster` resource ([#7485](https://github.com/hashicorp/terraform-provider-google/pull/7485))
* iam: fixed `google_service_account_id_token` datasource to work with User ADCs and Impersonated Credentials ([#7457](https://github.com/hashicorp/terraform-provider-google/pull/7457))
* logging: added bucket creation based on custom-id given for the resource `google_logging_project_bucket_config` ([#7492](https://github.com/hashicorp/terraform-provider-google/pull/7492))
* logging: Added support for exclusions options for `google_logging_project_sink` ([#7335](https://github.com/hashicorp/terraform-provider-google/pull/7335))
* oslogin: added ability to set a `project` on `google_os_login_ssh_public_key` ([#7505](https://github.com/hashicorp/terraform-provider-google/pull/7505))
* resourcemanager: added a precheck that the serviceusage API is enabled to `google_project` when `auto_create_network` is false, as configuring the GCE API is required in that circumstance ([#7447](https://github.com/hashicorp/terraform-provider-google/pull/7447))
* sql: added a check to `google_sql_database_instance` to catch failures early by seeing if Service Networking Connections already exists for the private network of the instance. ([#7499](https://github.com/hashicorp/terraform-provider-google/pull/7499))

BUG FIXES:
* accessapproval: fixed issue where, due to a recent API change, `google_*_access_approval.enrolled_services.cloud_product` entries specified as a URL would result in a permadiff ([#7468](https://github.com/hashicorp/terraform-provider-google/pull/7468))
* compute: fixed ability to clear `description` field on `google_compute_health_check` and `google_compute_region_health_check` ([#7500](https://github.com/hashicorp/terraform-provider-google/pull/7500))
* monitoring: fixed bug where deleting a `google_monitoring_dashboard` would give an "unsupported protocol scheme" error ([#7453](https://github.com/hashicorp/terraform-provider-google/pull/7453))

## 3.43.0 (October 12, 2020)

FEATURES:
* **New Data Source:** `google_pubsub_topic` ([#7426](https://github.com/hashicorp/terraform-provider-google/pull/7426))
* **New Data Source:** `google_compute_global_forwarding_rule` ([#7434](https://github.com/hashicorp/terraform-provider-google/pull/7434))
* **New Data Source:** `google_cloud_run_service` ([#7388](https://github.com/hashicorp/terraform-provider-google/pull/7388))
* **New Resource:** `google_bigtable_table_iam_member` ([#7410](https://github.com/hashicorp/terraform-provider-google/pull/7410))
* **New Resource:** `google_bigtable_table_iam_binding` ([#7410](https://github.com/hashicorp/terraform-provider-google/pull/7410))
* **New Resource:** `google_bigtable_table_iam_policy` ([#7410](https://github.com/hashicorp/terraform-provider-google/pull/7410))

IMPROVEMENTS:
* appengine: added ability to manage pre-firestore appengine applications. ([#7408](https://github.com/hashicorp/terraform-provider-google/pull/7408))
* bigquery: added support for `google_bigquery_table` `materialized_view` field ([#7080](https://github.com/hashicorp/terraform-provider-google/pull/7080))
* compute: Marked `google_compute_per_instance_config` as GA ([#7429](https://github.com/hashicorp/terraform-provider-google/pull/7429))
* compute: Marked `google_compute_region_per_instance_config` as GA ([#7429](https://github.com/hashicorp/terraform-provider-google/pull/7429))
* compute: Marked `stateful_disk` as GA in `google_compute_instance_group_manager` ([#7429](https://github.com/hashicorp/terraform-provider-google/pull/7429))
* compute: Marked `stateful_disk` as GA in `google_compute_region_instance_group_manager` ([#7429](https://github.com/hashicorp/terraform-provider-google/pull/7429))
* compute: added additional fields to the `google_compute_forwarding_rule` datasource. ([#7437](https://github.com/hashicorp/terraform-provider-google/pull/7437))
* dns: added `forwarding_path` field to `google_dns_policy` resource ([#7416](https://github.com/hashicorp/terraform-provider-google/pull/7416))
* netblock: changed `google_netblock_ip_ranges` to read from cloud.json file rather than DNS record ([#7157](https://github.com/hashicorp/terraform-provider-google/pull/7157))

BUG FIXES:
* accessapproval: fixed issue where, due to a recent API change, `google_*_access_approval.enrolled_services.cloud_product` entries specified as a URL would result in a permadiff
* bigquery: fixed an issue in `google_bigquery_job` where non-US locations could not be read ([#7418](https://github.com/hashicorp/terraform-provider-google/pull/7418))
* cloudrun: fixed an issue in `google_cloud_run_domain_mapping` where labels provided by Google would cause a diff ([#7407](https://github.com/hashicorp/terraform-provider-google/pull/7407))
* compute: Fixed an issue where `google_compute_region_backend_service` required `healthChecks` for a serverless network endpoint group. ([#7433](https://github.com/hashicorp/terraform-provider-google/pull/7433))
* container: fixed `node_config.image_type` perma-diff when specified in lower case. ([#7412](https://github.com/hashicorp/terraform-provider-google/pull/7412))
* datacatalog: fixed an error in `google_data_catalog_tag` when trying to set boolean field to `false` ([#7409](https://github.com/hashicorp/terraform-provider-google/pull/7409))
* monitoring: fixed bug where deleting a `google_monitoring_dashboard` would give an "unsupported protocol scheme" error

## 3.42.0 (October 05, 2020)

FEATURES:
* **New Resource:** google_data_loss_prevention_deidentify_template ([#7378](https://github.com/hashicorp/terraform-provider-google/pull/7378))

IMPROVEMENTS:
* compute: added support for updating `network_interface.[d].network` and `network_interface.[d].subnetwork` properties on `google_compute_instance`. ([#7358](https://github.com/hashicorp/terraform-provider-google/pull/7358))
* healthcare: added field `parser_config.version` to `google_healthcare_hl7_v2_store` ([#7357](https://github.com/hashicorp/terraform-provider-google/pull/7357))

BUG FIXES:
* bigquery: fixed an issue where `google_bigquery_table` would crash while reading an empty schema ([#7359](https://github.com/hashicorp/terraform-provider-google/pull/7359))
* compute: fixed an issue where `google_compute_instance_template` would throw an error for unspecified `disk_size_gb` values while upgrading the provider. ([#7355](https://github.com/hashicorp/terraform-provider-google/pull/7355))
* resourcemanager: fixed an issue in retrieving `google_active_folder` data source when the display name included whitespace ([#7395](https://github.com/hashicorp/terraform-provider-google/pull/7395))

## 3.41.0 (September 28, 2020)

IMPROVEMENTS:
* compute: added `SEV_CAPABLE` option to `guest_os_features` in `google_compute_image` resource. ([#7313](https://github.com/hashicorp/terraform-provider-google/pull/7313))
* tpu: added `use_service_networking` to `google_tpu_node` which enables Shared VPC Support. ([#7294](https://github.com/hashicorp/terraform-provider-google/pull/7294))

## 3.40.0 (September 21, 2020)

DEPRECATIONS:
* bigtable: Deprecated `instance_type` for `google_bigtable_instance` - it is now recommended to leave field unspecified. ([#7253](https://github.com/hashicorp/terraform-provider-google/pull/7253))

FEATURES:
* **New Data Source:** `google_compute_region_ssl_certificate` ([#7252](https://github.com/hashicorp/terraform-provider-google/pull/7252))
* **New Resource:** `google_compute_target_grpc_proxy` ([#7277](https://github.com/hashicorp/terraform-provider-google/pull/7277))

IMPROVEMENTS:
* cloudfunctions: added the ALLOW_INTERNAL_AND_GCLB option to `ingress_settings` of `google_cloudfunctions_function` resource. ([#7287](https://github.com/hashicorp/terraform-provider-google/pull/7287))
* cloudlbuild: added `options` and `artifacts` properties to `google_cloudbuild_trigger` ([#7280](https://github.com/hashicorp/terraform-provider-google/pull/7280))
* compute: added GRPC as a valid value for `google_compute_backend_service.protocol` (and regional equivalent) ([#7254](https://github.com/hashicorp/terraform-provider-google/pull/7254))
* compute: added support for configuring Internal load balancer for Cloud Run for Anthos ([#7268](https://github.com/hashicorp/terraform-provider-google/pull/7268))
* compute: added 'all' option for `google_compute_firewall` ([#7225](https://github.com/hashicorp/terraform-provider-google/pull/7225))
* dataflow : added `transformnameMapping` to `google_dataflow_job` ([#7259](https://github.com/hashicorp/terraform-provider-google/pull/7259))
* dns: added `force_destroy` option to `google_dns_managed_zone` to delete records created outside of Terraform ([#7289](https://github.com/hashicorp/terraform-provider-google/pull/7289))
* serviceusage: added ability to pass `google.project.id` to `google_project_service.project` ([#7255](https://github.com/hashicorp/terraform-provider-google/pull/7255))
* spanner: added schema update/update ddl support for `google_spanner_database` ([#7279](https://github.com/hashicorp/terraform-provider-google/pull/7279))

BUG FIXES:
* bigtable: fixed the update behaviour of the `single_cluster_routing` sub-fields in `google_bigtable_app_profile` ([#7266](https://github.com/hashicorp/terraform-provider-google/pull/7266))
* dataproc: fixed issues where updating `google_dataproc_cluster.cluster_config.autoscaling_policy` would do nothing, and where there was no way to remove a policy. ([#7269](https://github.com/hashicorp/terraform-provider-google/pull/7269))
* osconfig: fixed a potential crash in `google_os_config_patch_deployment` due to an unchecked nil value in `recurring_schedule` ([#7265](https://github.com/hashicorp/terraform-provider-google/pull/7265))
* serviceusage: fixed intermittent failure when a service is already being modified - added retries [#7230](https://github.com/hashicorp/terraform-provider-google/pull/7230))
* serviceusage: fixed an issue where `bigquery.googleapis.com` was getting enabled as the `bigquery-json.googleapis.com` alias instead, incorrectly. This had no user impact yet, but the alias may go away in the future. ([#7230](https://github.com/hashicorp/terraform-provider-google/pull/7230))

## 3.39.0 (September 15, 2020)

IMPROVEMENTS:
* compute: added `storage_locations` field to `google_compute_snapshot` ([#7201](https://github.com/hashicorp/terraform-provider-google/pull/7201))
* compute: added `kms_key_service_account`, `kms_key_self_link ` fields to `snapshot_encryption_key` field in `google_compute_snapshot` ([#7201](https://github.com/hashicorp/terraform-provider-google/pull/7201))
* compute: added `source_disk_encryption_key.kms_key_service_account` field to `google_compute_snapshot` ([#7201](https://github.com/hashicorp/terraform-provider-google/pull/7201))
* container: added `self_link` to `google_container_cluster` ([#7191](https://github.com/hashicorp/terraform-provider-google/pull/7191))
* container: marked `workload_metadata_config` as GA in `google_container_node_pool` ([#7192](https://github.com/hashicorp/terraform-provider-google/pull/7192))

BUG FIXES:
* bigquery: fixed a bug when a BigQuery table schema didn't have `name` in the schema. Previously it would panic; now it logs an error. ([#7215](https://github.com/hashicorp/terraform-provider-google/pull/7215))
* bigquery: fixed bug where updating `clustering` would force a new resource rather than update. ([#7195](https://github.com/hashicorp/terraform-provider-google/pull/7195))
* bigquerydatatransfer: fixed `params.secret_access_key` perma-diff for AWS S3 data transfer config types by adding a `sensitive_params` block with the `secret_access_key` attribute. ([#7174](https://github.com/hashicorp/terraform-provider-google/pull/7174))
* compute: fixed bug where `delete_default_routes_on_create=true` was not actually deleting the default routes on create. ([#7199](https://github.com/hashicorp/terraform-provider-google/pull/7199))

## 3.38.0 (September 08, 2020)

DEPRECATIONS:
* storage: deprecated `bucket_policy_only` field in `google_storage_bucket` in favour of `uniform_bucket_level_access` ([#7143](https://github.com/hashicorp/terraform-provider-google/pull/7143))

FEATURES:
* **New Resource:** google_compute_disk_iam_binding ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** google_compute_disk_iam_member ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** google_compute_disk_iam_policy ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** google_compute_region_disk_iam_binding ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** google_compute_region_disk_iam_member ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** google_compute_region_disk_iam_policy ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** google_data_loss_prevention_inspect_template ([#7123](https://github.com/hashicorp/terraform-provider-google/pull/7123))
* **New Resource:** google_data_loss_prevention_job_trigger ([#7123](https://github.com/hashicorp/terraform-provider-google/pull/7123))
* **New Resource:** google_data_loss_prevention_stored_info_type ([#7145](https://github.com/hashicorp/terraform-provider-google/pull/7145))

IMPROVEMENTS:
* compute: Added graceful termination to `google_compute_instance_group_manager` create calls so that partially created instance group managers will resume the original operation if the Terraform process is killed mid create. ([#7153](https://github.com/hashicorp/terraform-provider-google/pull/7153))
* container: added project override support to `google_container_cluster` and `google_container_nodepool` ([#7114](https://github.com/hashicorp/terraform-provider-google/pull/7114))
* osconfig: added rollout field to `google_os_config_patch_deployment` ([#7172](https://github.com/hashicorp/terraform-provider-google/pull/7172))
* provider: added a new field `billing_project` to the provider that's associated as a billing/quota project with most requests when `user_project_override` is true ([#7113](https://github.com/hashicorp/terraform-provider-google/pull/7113))
* resourcemanager: added additional fields to `google_projects` datasource ([#7139](https://github.com/hashicorp/terraform-provider-google/pull/7139))
* serviceusage: added project override support to `google_project_service` ([#7114](https://github.com/hashicorp/terraform-provider-google/pull/7114))

BUG FIXES:
* bigquerydatatransfer: fixed `params.secret_access_key` perma-diff for AWS S3 data transfer config types by adding a `sensitive_params` block with the `secret_access_key` attribute. ([#7174](https://github.com/hashicorp/terraform-provider-google/pull/7174))
* compute: Fixed bug with `google_netblock_ip_ranges` data source failing to read from the correct URL ([#7156](https://github.com/hashicorp/terraform-provider-google/pull/7156))
* compute: fixed updating `google_compute_instance.shielded_instance_config` by adding it to the `allow_stopping_for_update` list ([#7132](https://github.com/hashicorp/terraform-provider-google/pull/7132))

## 3.37.0 (August 31, 2020)
NOTES:
* Drop recommendation to use -provider= on import in documentation ([#7100](https://github.com/hashicorp/terraform-provider-google/pull/7100))

FEATURES:
* **New Resource:** `google_compute_image_iam_binding` ([#7070](https://github.com/hashicorp/terraform-provider-google/pull/7070))
* **New Resource:** `google_compute_image_iam_member` ([#7070](https://github.com/hashicorp/terraform-provider-google/pull/7070))
* **New Resource:** `google_compute_image_iam_policy` ([#7070](https://github.com/hashicorp/terraform-provider-google/pull/7070))
* **New Resource:** `google_compute_disk_iam_binding` ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** `google_compute_disk_iam_member` ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** `google_compute_disk_iam_policy` ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** `google_compute_region_disk_iam_binding` ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** `google_compute_region_disk_iam_member` ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))
* **New Resource:** `google_compute_region_disk_iam_policy` ([#7110](https://github.com/hashicorp/terraform-provider-google/pull/7110))

IMPROVEMENTS:
* appengine: added `vpc_access_connector` field to `google_app_engine_standard_app_version` resource ([#7062](https://github.com/hashicorp/terraform-provider-google/pull/7062))
* bigquery: added `notification_pubsub_topic` field to `google_bigquery_data_transfer_config` resource ([#7076](https://github.com/hashicorp/terraform-provider-google/pull/7076))
* compute: Added custom metadata fields and filter expressions to `google_compute_subnetwork` flow log configuration ([#7099](https://github.com/hashicorp/terraform-provider-google/pull/7099))
* compute: Added support to `google_compute_backend_service` for setting a serverless regional network endpoint group as `backend.group` ([#7066](https://github.com/hashicorp/terraform-provider-google/pull/7066))
* compute: added support for pd-balanced disk type for `google_compute_instance` ([#7108](https://github.com/hashicorp/terraform-provider-google/pull/7108))
* container: added support for pd-balanced disk type for `google_container_node_pool` ([#7108](https://github.com/hashicorp/terraform-provider-google/pull/7108))
* pubsub: added `retry_policy` to `google_pubsub_subscription` resource ([#7077](https://github.com/hashicorp/terraform-provider-google/pull/7077))

BUG FIXES:
* compute: fixed an issue where `google_compute_url_map` `path_matcher.default_route_action` would conflict with `default_url_redirect` ([#7063](https://github.com/hashicorp/terraform-provider-google/pull/7063))
* kms: updated `data_source_secret_manager_secret_version` to have consistent id value ([#7098](https://github.com/hashicorp/terraform-provider-google/pull/7098))

## 3.36.0 (August 24, 2020)

FEATURES:
* **New Resource:** `google_active_directory_domain_trust` ([#7056](https://github.com/hashicorp/terraform-provider-google/pull/7056))
* **New Resource:** `google_access_context_manager_service_perimeters` ([#7027](https://github.com/hashicorp/terraform-provider-google/pull/7027))
* **New Resource:** `google_access_context_manager_access_levels` ([#7027](https://github.com/hashicorp/terraform-provider-google/pull/7027))
* **New Resource:** `google_folder_access_approval_settings` ([#7010](https://github.com/hashicorp/terraform-provider-google/pull/7010))
* **New Resource:** `google_organization_access_approval_settings` ([#7010](https://github.com/hashicorp/terraform-provider-google/pull/7010))
* **New Resource:** `google_project_access_approval_settings` ([#7010](https://github.com/hashicorp/terraform-provider-google/pull/7010))
* **New Resource:** `google_bigquery_table_iam_policy` ([#7041](https://github.com/hashicorp/terraform-provider-google/pull/7041))
* **New Resource:** `google_bigquery_table_iam_binding` ([#7041](https://github.com/hashicorp/terraform-provider-google/pull/7041))
* **New Resource:** `google_bigquery_table_iam_member` ([#7041](https://github.com/hashicorp/terraform-provider-google/pull/7041))

IMPROVEMENTS:
* compute: added grpc_health_check block to compute_health_check ([#7038](https://github.com/hashicorp/terraform-provider-google/pull/7038))
* compute: added grpc_health_check block to compute_region_health_check ([#7038](https://github.com/hashicorp/terraform-provider-google/pull/7038))
* pubsub: added `enable_message_ordering` support to `google_pubsub_subscription` ([#7039](https://github.com/hashicorp/terraform-provider-google/pull/7039))
* sql: added project field to `google_sql_database_instance` datasource. ([#7007](https://github.com/hashicorp/terraform-provider-google/pull/7007))
* storage: added `ARCHIVE` as an accepted class for `google_storage_bucket` and `google_storage_bucket_object` ([#7030](https://github.com/hashicorp/terraform-provider-google/pull/7030))

BUG FIXES:
* all: updated base urls for compute, dns, storage, and bigquery APIs to their recommended endpoints ([#7045](https://github.com/hashicorp/terraform-provider-google/pull/7045))
* bigquery: fixed a bug where `dataset_access.iam_member` would produce inconsistent results after apply. ([#7047](https://github.com/hashicorp/terraform-provider-google/pull/7047))
* bigquery: fixed an issue with `use_legacy_sql` not being set to `false`. ([#7012](https://github.com/hashicorp/terraform-provider-google/pull/7012))
* dns: fixed an issue where `google_dns_managed_zone` would not remove `private_visibility_config` on updates ([#7022](https://github.com/hashicorp/terraform-provider-google/pull/7022))
* sql: fixed an issue where `google_sql_database_instance` would throw an error when removing `private_network`. Removing `private_network` now recreates the resource. ([#7054](https://github.com/hashicorp/terraform-provider-google/pull/7054))

## 3.35.0 (August 17, 2020)
NOTES:
* all: Updated lists of enums to display the enum options in the documentation pages. ([#6946](https://github.com/hashicorp/terraform-provider-google/pull/6946))

FEATURES:
* **New Resource:** `google_compute_region_network_endpoint_group` (supports serverless NEGs) ([#6960](https://github.com/hashicorp/terraform-provider-google/pull/6960))
* **New Resource:** `google_game_services_game_server_cluster` ([#6983](https://github.com/hashicorp/terraform-provider-google/pull/6983))
* **New Resource:** `google_game_services_game_server_config` ([#6983](https://github.com/hashicorp/terraform-provider-google/pull/6983))
* **New Resource:** `google_game_services_game_server_deployment_rollout` ([#6983](https://github.com/hashicorp/terraform-provider-google/pull/6983))
* **New Resource:** `google_game_services_game_server_deployment` ([#6983](https://github.com/hashicorp/terraform-provider-google/pull/6983))
* **New Resource:** `google_game_services_realm` ([#6983](https://github.com/hashicorp/terraform-provider-google/pull/6983))

IMPROVEMENTS:
* appengine: converted `google_app_engine_standard_app_version`'s `inbound_services` to an enum array, which enhances docs and provides some client-side validation. ([#6956](https://github.com/hashicorp/terraform-provider-google/pull/6956))
* cloudbuild: added tags, source, queue_ttl, logs_bucket, substitutions, and secrets to `google_cloudbuild_trigger` ([#6942](https://github.com/hashicorp/terraform-provider-google/pull/6942))
* cloudfunctions: Updated the `google_cloudfunctions_function` datasource to include new fields available in the API. ([#6935](https://github.com/hashicorp/terraform-provider-google/pull/6935))
* compute: added `source_image` and `source_snapshot` to `google_compute_image` ([#6980](https://github.com/hashicorp/terraform-provider-google/pull/6980))
* compute: added confidential_instance_config block to google_compute_instance ([#7000](https://github.com/hashicorp/terraform-provider-google/pull/7000))
* compute: added confidential_instance_config block to google_compute_instance_template ([#7000](https://github.com/hashicorp/terraform-provider-google/pull/7000))
* container: added `release_channel_default_version` field to `data.google_container_engine_versions` (GA) ([#6963](https://github.com/hashicorp/terraform-provider-google/pull/6963))
* container: added `release_channel` to `google_container-cluster` (GA) ([#6955](https://github.com/hashicorp/terraform-provider-google/pull/6955))
* iam: Added `public_key_type` field to `google_service_account_key ` ([#6999](https://github.com/hashicorp/terraform-provider-google/pull/6999))
* pubsub: added `filter` field to `google_pubsub_subscription` resource ([#6997](https://github.com/hashicorp/terraform-provider-google/pull/6997))
* resource-manager: updated documentation for `folder_iam_*` and `organization_iam_*` resources. ([#6991](https://github.com/hashicorp/terraform-provider-google/pull/6991))
* sql: added support for point_in_time_recovery for `google_sql_database_instance` ([#6944](https://github.com/hashicorp/terraform-provider-google/pull/6944))

BUG FIXES:
* appengine: Set `iap` to computed in `google_app_engine_application` ([#6951](https://github.com/hashicorp/terraform-provider-google/pull/6951))
* artifactrepository: Fixed import failure of `google_artifact_registry_repository`. ([#6957](https://github.com/hashicorp/terraform-provider-google/pull/6957))
* compute: fixed shielded instance config, which had been failing to apply due to a field rename on the GCP side. ([#6943](https://github.com/hashicorp/terraform-provider-google/pull/6943))
* monitoring: fixed validation rules for `google_monitoring_slo` `windows_based_sli.metric_sum_in_range.max` field ([#6974](https://github.com/hashicorp/terraform-provider-google/pull/6974))
* osconfig: fixed `google_os_config_patch_deployment` `windows_update.classifications` field to work correctly, accepting multiple values. ([#6946](https://github.com/hashicorp/terraform-provider-google/pull/6946))

## 3.34.0 (August 11, 2020)
NOTES:
* redis: explicitly noted in `google_redis_instance` documentation that `"REDIS_5_0"` is supported ([#6917](https://github.com/terraform-providers/terraform-provider-google/pull/6917))
* all: fix markdown formatting while showing enum values in documentation ([#6924](https://github.com/terraform-providers/terraform-provider-google/pull/6924)).

IMPROVEMENTS:
* bigtable: added support for labels in `google_bigtable_instance` ([#6921](https://github.com/terraform-providers/terraform-provider-google/pull/6921))
* cloudfunctions: updated `google_cloudfunctions_function` datasource to include new fields. ([#6935](https://github.com/terraform-providers/terraform-provider-google/pull/6935))
* redis: added `persistence_iam_identity` output field to `google_redis_instance` ([#6917](https://github.com/terraform-providers/terraform-provider-google/pull/6917))
* storage: added google_storage_bucket_object.media_link. ([#6897](https://github.com/terraform-providers/terraform-provider-google/pull/6897))

BUG FIXES:
* all: fixed crash due to nil context when loading credentials ([#6903](https://github.com/terraform-providers/terraform-provider-google/pull/6903))
* compute: fixed issue where the `project` field in `data.google_compute_network_endpoint_group` was returning an error when specified ([#6918](https://github.com/terraform-providers/terraform-provider-google/pull/6918))
* sourcerepo: fixed perma-diff in `google_sourcerepo_repository` ([#6886](https://github.com/terraform-providers/terraform-provider-google/pull/6886))

## 3.33.0 (August 04, 2020)

DEPRECATIONS:
* compute: deprecated `enable_logging` on `google_compute_firewall`, define `log_config.metadata` to enable logging instead. ([#6871](https://github.com/terraform-providers/terraform-provider-google/pull/6871))

FEATURES:
* **New Resource:** `google_active_directory_domain` ([#6866](https://github.com/terraform-providers/terraform-provider-google/pull/6866))

IMPROVEMENTS:
* cloudrun: added `ports` field to `google_cloud_run_service` `templates.spec.containers` ([#6873](https://github.com/terraform-providers/terraform-provider-google/pull/6873))
* compute: added `log_config.metadata` to `google_compute_firewall`, defining this will enable logging. ([#6871](https://github.com/terraform-providers/terraform-provider-google/pull/6871))

BUG FIXES:
* container: Fixed a crash in `google_container_cluster` when `""` was specified for `resource_usage_export_config.bigquery_destination.dataset_id`. ([#6839](https://github.com/terraform-providers/terraform-provider-google/pull/6839))
* endpoints: Fixed a crash when `google_endpoints_service` is used on a machine without timezone data ([#6849](https://github.com/terraform-providers/terraform-provider-google/pull/6849))
* resourcemanager: bumped `google_project` timeout defaults to 10 minutes (from 4) ([#6859](https://github.com/terraform-providers/terraform-provider-google/pull/6859))

## 3.32.0 (July 27, 2020)
FEATURES:
* **New Data Source:** `google_sql_database_instance`  #2841 ([#6797](https://github.com/terraform-providers/terraform-provider-google/pull/6797))
* **New Resource:** `google_cloud_asset_folder_feed` ([#6821](https://github.com/terraform-providers/terraform-provider-google/pull/6821))
* **New Resource:** `google_cloud_asset_organization_feed` ([#6821](https://github.com/terraform-providers/terraform-provider-google/pull/6821))
* **New Resource:** `google_cloud_asset_project_feed` ([#6821](https://github.com/terraform-providers/terraform-provider-google/pull/6821))
* **New Resource:** `google_monitoring_metric_descriptor` ([#6829](https://github.com/terraform-providers/terraform-provider-google/pull/6829))

IMPROVEMENTS:
* filestore: Added support for filestore high scale tier. ([#6828](https://github.com/terraform-providers/terraform-provider-google/pull/6828))
* resourcemanager: Added `folder_id` as computed attribute to `google_folder` resource and datasource. ([#6823](https://github.com/terraform-providers/terraform-provider-google/pull/6823))
* compute: Added support to `google_compute_backend_service` for setting a network endpoint group as `backend.group`. ([#6853](https://github.com/terraform-providers/terraform-provider-google/pull/6853))

BUG FIXES:
* container: Fixed a crash in `google_container_cluster` when `""` was specified for `resource_usage_export_config.bigquery_destination.dataset_id`. ([#6839](https://github.com/terraform-providers/terraform-provider-google/pull/6839))
* bigquery: Fixed bug where a permadiff would show up when adding a column to the middle of a `bigquery_table.schema` ([#6803](https://github.com/terraform-providers/terraform-provider-google/pull/6803))

## 3.31.0 (July 20, 2020)

FEATURES:
* **New Data Source:** `google_service_account_id_token` ([#6791](https://github.com/terraform-providers/terraform-provider-google/pull/6791))
* **New Resource:** `google_cloudiot_device` ([#6785](https://github.com/terraform-providers/terraform-provider-google/pull/6785))

IMPROVEMENTS:
* bigquery: added support for BigQuery custom schemas for external data using CSV / NDJSON ([#6772](https://github.com/terraform-providers/terraform-provider-google/pull/6772))

## 3.30.0 (July 13, 2020)
FEATURES:
* **New Resource:** `google_os_config_patch_deployment` ([#6741](https://github.com/terraform-providers/terraform-provider-google/pull/6741))

IMPROVEMENTS:
* iam: made the `condition` block GA for all IAM resource and datasource types. ([#6748](https://github.com/terraform-providers/terraform-provider-google/pull/6748))

BUG FIXES:
* container: added the ability to update `database_encryption` without recreating the cluster. ([#6757](https://github.com/terraform-providers/terraform-provider-google/pull/6757))
* endpoints: fixed `google_endpoints_service` to allow dependent resources to plan based on the `config_id` value. ([#6722](https://github.com/terraform-providers/terraform-provider-google/pull/6722))
* runtimeconfig: fixed `Requested entity was not found.` error when config was deleted outside of terraform. ([#6753](https://github.com/terraform-providers/terraform-provider-google/pull/6753))

## 3.29.0 (July 06, 2020)
NOTES:
* added the `https://www.googleapis.com/auth/cloud-identity` scope to the provider by default ([#6681](https://github.com/terraform-providers/terraform-provider-google/pull/6681))
* `google_app_engine_*_version`'s `service` field is required; previously it would have passed validation but failed on apply if it were absent. ([#6720](https://github.com/terraform-providers/terraform-provider-google/pull/6720))

FEATURES:
* **New Resource:** `google_kms_key_ring_import_job` ([#6682](https://github.com/terraform-providers/terraform-provider-google/pull/6682))
* **New Resource:** `google_folder_iam_audit_config` ([#6708](https://github.com/terraform-providers/terraform-provider-google/pull/6708))

IMPROVEMENTS:
* bigquery: Added `"HOUR"` option for `google_bigquery_table` time partitioning (`type`) ([#6702](https://github.com/terraform-providers/terraform-provider-google/pull/6702))
* bigquery: Added support for BigQuery hourly time partitioning  ([#6675](https://github.com/terraform-providers/terraform-provider-google/pull/6675))
* compute: Added `mode` to `google_compute_region_autoscaler` `autoscaling_policy` ([#6685](https://github.com/terraform-providers/terraform-provider-google/pull/6685))
* container: Promoted `google_container_cluster` `database_encryption` to GA. ([#6701](https://github.com/terraform-providers/terraform-provider-google/pull/6701))
* endpoints: `google_endpoints_service` now allows dependent resources to plan based on the `config_id` value. ([#6722](https://github.com/terraform-providers/terraform-provider-google/pull/6722))
* monitoring: added `request_method`, `content_type`, and `body` fields within the `http_check` object to `google_monitoring_uptime_check_config` resource ([#6700](https://github.com/terraform-providers/terraform-provider-google/pull/6700))

BUG FIXES:
* compute: fixed an issue in `compute_url_map` where `path_matcher` sub-fields would conflict with `default_service` ([#6721](https://github.com/terraform-providers/terraform-provider-google/pull/6721))

## 3.28.0 (June 29, 2020)

FEATURES:
* **New Data Source:** `google_redis_instance` ([#6649](https://github.com/terraform-providers/terraform-provider-google/pull/6649))
* **New Resource:** `google_notebook_environment` ([#6639](https://github.com/terraform-providers/terraform-provider-google/pull/6639))
* **New Resource:** `google_notebook_instance` ([#6639](https://github.com/terraform-providers/terraform-provider-google/pull/6639))

IMPROVEMENTS:
* appengine: Enabled provisioning Firestore on a new project by adding the option to specify `database_type` in `google_app_engine_application` ([#6629](https://github.com/terraform-providers/terraform-provider-google/pull/6629))
* compute: Added `mode` to `google_compute_autoscaler` `autoscaling_policy` ([#6664](https://github.com/terraform-providers/terraform-provider-google/pull/6664))
* dns: enabled google_dns_policy to accept network id ([#6624](https://github.com/terraform-providers/terraform-provider-google/pull/6624))

BUG FIXES:
* appengine: Added polling to `google_app_engine_firewall_rule` to prevent issues with eventually consistent creation ([#6633](https://github.com/terraform-providers/terraform-provider-google/pull/6633))
* compute: Allowed updating `google_compute_network_peering_routes_config ` `import_custom_routes` and  `export_custom_routes` to false ([#6625](https://github.com/terraform-providers/terraform-provider-google/pull/6625))
* netblock: fixed the google netblock ranges returned by the `google_netblock_ip_ranges` by targeting json on gstatic domain instead of reading SPF dns records (solution provided by network team) ([#6650](https://github.com/terraform-providers/terraform-provider-google/pull/6650))

## 3.27.0 (June 23, 2020)

IMPROVEMENTS:
* accesscontextmanager: Added `custom` config to `google_access_context_manager_access_level` ([#6611](https://github.com/terraform-providers/terraform-provider-google/pull/6611))
* cloudbuild: Added `invert_regex` flag in Github PullRequestFilter and PushFilter in triggerTemplate ([#6594](https://github.com/terraform-providers/terraform-provider-google/pull/6594))
* cloudrun: Added `template.spec.timeout_seconds` to `google_cloud_run_service` ([#6575](https://github.com/terraform-providers/terraform-provider-google/pull/6575))
* compute: Added `export_subnet_routes_with_public_ip` and `import_subnet_routes_with_public_ip` to `google_compute_network_peering` ([#6586](https://github.com/terraform-providers/terraform-provider-google/pull/6586))
* compute: Added support for `google_compute_instance_group` `instances` to accept instance id field as well as self_link ([#6569](https://github.com/terraform-providers/terraform-provider-google/pull/6569))
* dns: Added support for `google_dns_policy` network to accept `google_compute_network.id` ([#6624](https://github.com/terraform-providers/terraform-provider-google/pull/6624))
* redis: Added validation for name attribute in `redis_instance` ([#6581](https://github.com/terraform-providers/terraform-provider-google/pull/6581))
* sql: Promoted `google_sql_database_instance` `root_password` (MS SQL) to GA ([#6601](https://github.com/terraform-providers/terraform-provider-google/pull/6601))

BUG FIXES:
* bigquery: Fixed `range_partitioning.range.start` so that the value `0` is sent in `google_bigquery_table` ([#6562](https://github.com/terraform-providers/terraform-provider-google/pull/6562))
* container: Fixed a regression in `google_container_cluster` where the location was not inferred when using a `subnetwork` shortname value like `name` ([#6568](https://github.com/terraform-providers/terraform-provider-google/pull/6568))
* datastore: Added retries to `google_datastore_index` requests when under contention. ([#6563](https://github.com/terraform-providers/terraform-provider-google/pull/6563))
* kms: Fixed the `id` value in the `google_kms_crypto_key_version` datasource to include a `/v1` part following `//cloudkms.googleapis.com/`, making it useful for interpolation into Binary Authorization. ([#6576](https://github.com/terraform-providers/terraform-provider-google/pull/6576))

## 3.26.0 (June 15, 2020)

FEATURES:
* **New Resource:** `google_data_catalog_tag` ([#6550](https://github.com/terraform-providers/terraform-provider-google/pull/6550))
* **New Resource:** `google_bigquery_dataset_iam_binding` ([#6553](https://github.com/terraform-providers/terraform-provider-google/pull/6553))
* **New Resource:** `google_bigquery_dataset_iam_member` ([#6553](https://github.com/terraform-providers/terraform-provider-google/pull/6553))
* **New Resource:** `google_bigquery_dataset_iam_policy` ([#6553](https://github.com/terraform-providers/terraform-provider-google/pull/6553))
* **New Resource:** `google_memcache_instance` ([#6540](https://github.com/terraform-providers/terraform-provider-google/pull/6540))
* **New Resource:** `google_network_management_connectivity_test` ([#6529](https://github.com/terraform-providers/terraform-provider-google/pull/6529))

IMPROVEMENTS:
* compute: added `default_route_action` to `compute_url_map` and `compute_url_map.path_matchers` ([#6547](https://github.com/terraform-providers/terraform-provider-google/pull/6547))
* dialogflow: Changed `google_dialogflow_agent.time_zone` to be updatable ([#6519](https://github.com/terraform-providers/terraform-provider-google/pull/6519))
* dns: enabled google_dns_managed_zone to accept network id for two attributes ([#6533](https://github.com/terraform-providers/terraform-provider-google/pull/6533))
* healthcare: Added support for `streaming_configs` to `google_healthcare_fhir_store` ([#6551](https://github.com/terraform-providers/terraform-provider-google/pull/6551))
* monitoring: added `matcher` attribute to `content_matchers` block for `google_monitoring_uptime_check_config` ([#6558](https://github.com/terraform-providers/terraform-provider-google/pull/6558))

BUG FIXES:
* compute: fixed issue where trying to update the region of `google_compute_subnetwork` would fail instead of destroying/recreating the subnetwork ([#6522](https://github.com/terraform-providers/terraform-provider-google/pull/6522))
* dataflow: added retries in `google_dataflow_job` for common retryable API errors when waiting for job to update ([#6552](https://github.com/terraform-providers/terraform-provider-google/pull/6552))
* dataflow: changed the update logic for `google_dataflow_job` to wait for the replacement job to start successfully before modifying the resource ID to point to the replacement job ([#6534](https://github.com/terraform-providers/terraform-provider-google/pull/6534))

## 3.25.0 (June 08, 2020)

FEATURES:
* **New Resource:** `google_data_catalog_tag_template` ([#6485](https://github.com/terraform-providers/terraform-provider-google/pull/6485))
* **New Resource:** `google_container_analysis_occurence` ([#6474](https://github.com/terraform-providers/terraform-provider-google/pull/6474))

IMPROVEMENTS:
* appengine: added `inbound_services` to `StandardAppVersion` resource ([#6514](https://github.com/terraform-providers/terraform-provider-google/pull/6514))
* bigquery: Promoted `google_bigquery_table` `range_partitioning` to GA ([#6488](https://github.com/terraform-providers/terraform-provider-google/pull/6488))
* bigquery: Added support for `google_bigquery_table` `hive_partitioning_options` ([#6488](https://github.com/terraform-providers/terraform-provider-google/pull/6488))
* container: Promoted `google_container_cluster.workload_identity_config` to GA. ([#6490](https://github.com/terraform-providers/terraform-provider-google/pull/6490))
* container_analysis: Added top-level generic note fields to `google_container_analysis_note` ([#6474](https://github.com/terraform-providers/terraform-provider-google/pull/6474))

BUG FIXES:
* bigquery: Fixed an issue where `google_bigquery_job` would return "was present, but now absent" error after job creation ([#6489](https://github.com/terraform-providers/terraform-provider-google/pull/6489))
* container: Changed retry logic for `google_container_node_pool` deletion to use timeouts and retry errors more specifically when cluster is updating. ([#6335](https://github.com/terraform-providers/terraform-provider-google/pull/6335))
* dataflow: fixed an issue where `google_dataflow_job` would try to update `max_workers` ([#6468](https://github.com/terraform-providers/terraform-provider-google/pull/6468))
* dataflow: fixed an issue where updating `on_delete` in `google_dataflow_job` would cause the job to be replaced ([#6468](https://github.com/terraform-providers/terraform-provider-google/pull/6468))
* compute: fixed issue where removing all target pools from `google_compute_instance_group_manager` or `google_compute_region_instance_group_manager` had no effect ([#6492](https://github.com/terraform-providers/terraform-provider-google/pull/6492))
* functions: Added retry to `google_cloudfunctions_function` creation when API returns error while pulling source from GCS ([#6476](https://github.com/terraform-providers/terraform-provider-google/pull/6476))
* provider: Removed credentials from output error when provider cannot parse given credentials ([#6473](https://github.com/terraform-providers/terraform-provider-google/pull/6473))

## 3.24.0 (June 01, 2020)

FEATURES:
* **New Data Source:** `google_secret_manager_secret_version` ([#6432](https://github.com/terraform-providers/terraform-provider-google/pull/6432))
* **New Resources:** `google_data_catalog_entry_group_iam_*` ([#6438](https://github.com/terraform-providers/terraform-provider-google/pull/6438))
* **New Resource:** `google_data_catalog_entry_group` ([#6438](https://github.com/terraform-providers/terraform-provider-google/pull/6438))
* **New Resource:** `google_data_catalog_entry` ([#6444](https://github.com/terraform-providers/terraform-provider-google/pull/6444))
* **New Resource:** `google_dns_policy` is now GA ([#6439](https://github.com/terraform-providers/terraform-provider-google/pull/6439))
* **New Resource:** `google_secret_manager_secret` ([#6432](https://github.com/terraform-providers/terraform-provider-google/pull/6432))
* **New Resources:** `google_secret_manager_secret_iam_*` ([#6432](https://github.com/terraform-providers/terraform-provider-google/pull/6432))
* **New Resource:** `google_secret_manager_secret_version` ([#6432](https://github.com/terraform-providers/terraform-provider-google/pull/6432))

IMPROVEMENTS:
* appengine: added `handlers` to `google_flexible_app_version` ([#6449](https://github.com/terraform-providers/terraform-provider-google/pull/6449))
* bigquery: suppressed diffs between fully qualified URLs and relative paths that reference the same table or dataset in `google_bigquery_job` ([#6451](https://github.com/terraform-providers/terraform-provider-google/pull/6451))
* dns: Promoted the following `google_dns_managed_zone ` fields to GA: `forwarding_config`, `peering_config` ([#6439](https://github.com/terraform-providers/terraform-provider-google/pull/6439))

BUG FIXES:
* appengine: added ability to fully sync `StandardAppVersion` resources ([#6435](https://github.com/terraform-providers/terraform-provider-google/pull/6435))
* bigquery: Fixed an issue with `google_bigquery_dataset_access` failing for primitive role `roles/bigquery.dataViewer` ([#6431](https://github.com/terraform-providers/terraform-provider-google/pull/6431))
* dataflow: fixed an issue where `google_dataflow_job` would try to update `max_workers` ([#6468](https://github.com/terraform-providers/terraform-provider-google/pull/6468))
* dataflow: fixed an issue where updating `on_delete` in `google_dataflow_job` would cause the job to be replaced ([#6468](https://github.com/terraform-providers/terraform-provider-google/pull/6468))
* os_login: Fixed `google_os_login_ssh_public_key` `key` field attempting to update in-place ([#6433](https://github.com/terraform-providers/terraform-provider-google/pull/6433))

## 3.23.0 (May 26, 2020)

BREAKING CHANGES:
* The base url for the `monitoring` endpoint no longer includes the API version (previously "v3/"). If you use a `monitoring_custom_endpoint`, remove the trailing "v3/". ([#6424](https://github.com/terraform-providers/terraform-provider-google/pull/6424))

FEATURES:
* **New Data Source:** `google_iam_testable_permissions` ([#6382](https://github.com/terraform-providers/terraform-provider-google/pull/6382))
* **New Resource:** `google_monitoring_dashboard` ([#6424](https://github.com/terraform-providers/terraform-provider-google/pull/6424))

IMPROVEMENTS:

* bigquery: Added ability for various `table_id` fields (and one `dataset_id` field) in `google_bigquery_job` to specify a relative path instead of just the table id ([#6404](https://github.com/terraform-providers/terraform-provider-google/pull/6404))
* composer: Added support for `google_composer_environment` `config.private_environment_config.cloud_sql_ipv4_cidr_block` ([#6392](https://github.com/terraform-providers/terraform-provider-google/pull/6392))
* composer: Added support for `google_composer_environment` `config.private_environment_config.web_server_ipv4_cidr_block` ([#6392](https://github.com/terraform-providers/terraform-provider-google/pull/6392))
* container: Added update support for `node_config.workload_metadata_config` to `google_container_node_pool` ([#6430](https://github.com/terraform-providers/terraform-provider-google/pull/6430))
* container: Added the ability to unspecify `google_container_cluster`'s `min_master_version` field ([#6373](https://github.com/terraform-providers/terraform-provider-google/pull/6373))
* monitoring: Added window-based SLI to `google_monitoring_slo` ([#6381](https://github.com/terraform-providers/terraform-provider-google/pull/6381))


BUG FIXES:
* compute: Fixed an issue where `google_compute_route` creation failed while VPC peering was in progress. ([#6410](https://github.com/terraform-providers/terraform-provider-google/pull/6410))
* Fixed an issue where data source `google_organization` would ignore exact domain matches if multiple domains were found ([#6420](https://github.com/terraform-providers/terraform-provider-google/pull/6420))
* compute: Fixed `google_compute_interconnect_attachment` `edge_availability_domain` diff when the field is unspecified ([#6419](https://github.com/terraform-providers/terraform-provider-google/pull/6419))
* compute: fixed error where plan would error if `google_compute_region_disk_resource_policy_attachment` had been deleted outside of terraform. ([#6367](https://github.com/terraform-providers/terraform-provider-google/pull/6367))
* compute: raise limit on number of `src_ip_ranges` values in `google_compute_security_policy` to supported 10 ([#6394](https://github.com/terraform-providers/terraform-provider-google/pull/6394))
* iam: Fixed an issue where `google_service_account` shows an error after creating the resource ([#6391](https://github.com/terraform-providers/terraform-provider-google/pull/6391))

## 3.22.0 (May 18, 2020)
BREAKING CHANGE:
* `google_bigtable_instance` resources now cannot be destroyed unless `deletion_protection = false` is set in state for the resource. ([#6357](https://github.com/terraform-providers/terraform-provider-google/pull/6357))

FEATURES:
* **New Resource:** `google_dialogflow_entity_type` ([#6339](https://github.com/terraform-providers/terraform-provider-google/pull/6339))

IMPROVEMENTS:
* bigtable: added `deletion_protection` field to `google_bigtable_instance` to make deleting them require an explicit intent. ([#6357](https://github.com/terraform-providers/terraform-provider-google/pull/6357))
* compute: Added `google_compute_region_backend_service` `port_name` parameter ([#6327](https://github.com/terraform-providers/terraform-provider-google/pull/6327))
* dataproc: Updated `google_dataproc_cluster.software_config.optional_components` to include new options. ([#6330](https://github.com/terraform-providers/terraform-provider-google/pull/6330))
* monitoring: Added `request_based` SLI support to `google_monitoring_slo` ([#6353](https://github.com/terraform-providers/terraform-provider-google/pull/6353))
* storage: added `google_storage_bucket` bucket name to the error message when the bucket can't be deleted because it's not empty ([#6355](https://github.com/terraform-providers/terraform-provider-google/pull/6355))

BUG FIXES:
* bigquery: Fixed error where `google_bigquery_dataset_access` resources could not be found post-creation if role was set to a predefined IAM role with an equivalent primitive role (e.g. `roles/bigquery.dataOwner` and `OWNER`) ([#6307](https://github.com/terraform-providers/terraform-provider-google/pull/6307))
* compute: Fixed permadiff in `google_compute_instance_template`'s `network_tier`. ([#6344](https://github.com/terraform-providers/terraform-provider-google/pull/6344))
* compute: Removed permadiff or errors on update for `google_compute_backend_service` and `google_compute_region_backend_service` when `consistent_hash` values were previously set on  backend service but are not supported by updated value of `locality_lb_policy` ([#6316](https://github.com/terraform-providers/terraform-provider-google/pull/6316))
* sql: Fixed occasional failure to delete `google_sql_database_instance` and `google_sql_user`. ([#6318](https://github.com/terraform-providers/terraform-provider-google/pull/6318))

## 3.21.0 (May 11, 2020)

FEATURES:
* **New Resource:** `google_compute_region_target_http_proxy` is now GA ([#6245](https://github.com/terraform-providers/terraform-provider-google/pull/6245))
* **New Resource:** `google_compute_region_target_https_proxy` is now GA ([#6245](https://github.com/terraform-providers/terraform-provider-google/pull/6245))
* **New Resource:** `google_compute_region_url_map` is now GA ([#6245](https://github.com/terraform-providers/terraform-provider-google/pull/6245))
* **New Resource:** `google_logging_billing_account_bucket_config` ([#6227](https://github.com/terraform-providers/terraform-provider-google/pull/6227))
* **New Resource:** `google_logging_folder_bucket_config` ([#6227](https://github.com/terraform-providers/terraform-provider-google/pull/6227))
* **New Resource:** `google_logging_organization_bucket_config` ([#6227](https://github.com/terraform-providers/terraform-provider-google/pull/6227))
* **New Resource:** `google_logging_project_bucket_config` ([#6227](https://github.com/terraform-providers/terraform-provider-google/pull/6227))

IMPROVEMENTS:
* all: added configurable timeouts to several resources that did not previously have them ([#6226](https://github.com/terraform-providers/terraform-provider-google/pull/6226))
* bigquery: added `service_account_name` field to `google_bigquery_data_transfer_config` resource ([#6221](https://github.com/terraform-providers/terraform-provider-google/pull/6221))
* cloudfunctions: Added validation to label keys for `google_cloudfunctions_function` as API errors aren't useful. ([#6228](https://github.com/terraform-providers/terraform-provider-google/pull/6228))
* compute: Promoted the following `google_compute_backend_service` fields to GA: `circuit_breakers`, `consistent_hash`, `custom_request_headers`, `locality_lb_policy`, `outlier_detection` ([#6245](https://github.com/terraform-providers/terraform-provider-google/pull/6245))
* compute: Promoted the following `google_compute_region_backend_service` fields to GA: `affinity_cookie_ttl_sec`,`circuit_breakers`, `consistent_hash`, `failover_policy`, `locality_lb_policy`, `outlier_detection`, `log_config`, `failover` ([#6245](https://github.com/terraform-providers/terraform-provider-google/pull/6245))
* container: Promoted `google_container_cluster.addons_config.cloudrun_config` from beta to GA. ([#6304](https://github.com/terraform-providers/terraform-provider-google/pull/6304))
* container: Promoted `google_container_cluster.enable_shielded_nodes` from beta to GA. ([#6303](https://github.com/terraform-providers/terraform-provider-google/pull/6303))
* container: Promoted `node_locations` to `google_container_node_pool` and `google_container_cluster.node_pool` from beta to GA ([#6253](https://github.com/terraform-providers/terraform-provider-google/pull/6253))
* dataflow: Added drift detection for `google_dataflow_job` `template_gcs_path` and `temp_gcs_location` fields ([#6257](https://github.com/terraform-providers/terraform-provider-google/pull/6257))
* dataflow: Added support for update-by-replacement to `google_dataflow_job` ([#6257](https://github.com/terraform-providers/terraform-provider-google/pull/6257))
* dataflow: Added support for providing additional experiments to Dataflow job ([#6196](https://github.com/terraform-providers/terraform-provider-google/pull/6196))
* storage: Added retries for `google_storage_bucket_iam_*` on 412 (precondition not met) errors for eventually consistent bucket creation. ([#6235](https://github.com/terraform-providers/terraform-provider-google/pull/6235))

BUG FIXES:
* all: fixed bug where timeouts specified in units other than minutes were getting incorrectly rounded. Also fixed several instances of timeout values being used from the wrong method. ([#6218](https://github.com/terraform-providers/terraform-provider-google/pull/6218))
* accesscontextmanager: Fixed setting `require_screen_lock` to true for `google_access_context_manager_access_level` ([#6234](https://github.com/terraform-providers/terraform-provider-google/pull/6234))
* appengine: Changed `google_app_engine_application` to respect updates in `iap` ([#6216](https://github.com/terraform-providers/terraform-provider-google/pull/6216))
* bigquery: Fixed error where `google_bigquery_dataset_access` resources could not be found post-creation if role was set to a predefined IAM role with an equivalent primative role (e.g. `roles/bigquery.dataOwner` and `OWNER`) ([#6307](https://github.com/terraform-providers/terraform-provider-google/pull/6307))
* bigquery: Fixed the `google_sheets_options` at least one of logic. ([#6280](https://github.com/terraform-providers/terraform-provider-google/pull/6280))
* cloudscheduler: Fixed permadiff for `google_cloud_scheduler_job.retry_config.*` block when API provides default values ([#6278](https://github.com/terraform-providers/terraform-provider-google/pull/6278))
* compute: Added lock to prevent `google_compute_route` from changing while peering operations are happening on its network ([#6243](https://github.com/terraform-providers/terraform-provider-google/pull/6243))
* compute: fixed issue where the default value for the attribute `advertise_mode` on `google_compte_router_peer` was not populated on import ([#6265](https://github.com/terraform-providers/terraform-provider-google/pull/6265))
* container: Fix occasional error with `container_node_pool` partially-successful creations not being recorded if an error occurs on the GCP side. ([#6305](https://github.com/terraform-providers/terraform-provider-google/pull/6305))
* container: fixed issue where terraform would error if a gke instance group was deleted out-of-band ([#6242](https://github.com/terraform-providers/terraform-provider-google/pull/6242))
* storage: Fixed setting/reading `google_storage_bucket_object`  metadata on API object ([#6271](https://github.com/terraform-providers/terraform-provider-google/pull/6271))
* storage: Marked the credentials field in `google_storage_object_signed_url` as sensitive so it doesn't expose private credentials. ([#6272](https://github.com/terraform-providers/terraform-provider-google/pull/6272))

## 3.20.0 (May 04, 2020)

FEATURES:
* **New Resource:** `google_healthcare_dataset_iam_binding` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_dataset_iam_member` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_dataset_iam_policy` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_dataset` is now GA ([#6164](https://github.com/terraform-providers/terraform-provider-google/pull/6164))
* **New Resource:** `google_healthcare_dicom_store_iam_binding` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_dicom_store_iam_member` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_dicom_store_iam_policy` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_dicom_store` is now GA ([#6164](https://github.com/terraform-providers/terraform-provider-google/pull/6164))
* **New Resource:** `google_healthcare_fhir_store_iam_binding` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_fhir_store_iam_member` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_fhir_store_iam_policy` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_fhir_store` is now GA ([#6164](https://github.com/terraform-providers/terraform-provider-google/pull/6164))
* **New Resource:** `google_healthcare_hl7_v2_store_iam_binding` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_hl7_v2_store_iam_member` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_hl7_v2_store_iam_policy` is now GA ([#6193](https://github.com/terraform-providers/terraform-provider-google/pull/6193))
* **New Resource:** `google_healthcare_hl7_v2_store` is now GA ([#6164](https://github.com/terraform-providers/terraform-provider-google/pull/6164))


IMPROVEMENTS:
* appengine: Added `automatic_scaling`, `basic_scaling`, and `manual_scaling` to `google_app_engine_standard_app_version` ([#6183](https://github.com/terraform-providers/terraform-provider-google/pull/6183))
* bigquery: added `service_account_name` field to `google_bigquery_data_transfer_config` resource ([#6221](https://github.com/terraform-providers/terraform-provider-google/pull/6221))
* bigtable: added ability to add/remove column families in `google_bigtable_table` ([#6188](https://github.com/terraform-providers/terraform-provider-google/pull/6188))
* cloudfunctions: Added validation to label keys for `google_cloudfunctions_function` as API errors aren't useful. ([#6228](https://github.com/terraform-providers/terraform-provider-google/pull/6228))
* compute: Added support for default URL redirects to `google_compute_url_map` and `google_compute_region_url_map` ([#6203](https://github.com/terraform-providers/terraform-provider-google/pull/6203))
* dataflow: Add support for `additional_experiments` to `google_dataflow_job` ([#6196](https://github.com/terraform-providers/terraform-provider-google/pull/6196))

BUG FIXES:
* accesscontextmanager: Fixed setting `require_screen_lock` to true for `google_access_context_manager_access_level` ([#6234](https://github.com/terraform-providers/terraform-provider-google/pull/6234))
* appengine: Changed `google_app_engine_application` to respect updates in `iap` ([#6216](https://github.com/terraform-providers/terraform-provider-google/pull/6216))
* sql: Allowed `binary_log_enabled` to be disabled. ([#6163](https://github.com/terraform-providers/terraform-provider-google/pull/6163))
* storage: Added retries for `google_storage_bucket_iam_*` on 412 (precondition not met) errors for eventually consistent bucket creation. ([#6235](https://github.com/terraform-providers/terraform-provider-google/pull/6235))

## 3.19.0 (April 27, 2020)

FEATURES:
* **New Resource:** `google_bigquery_job` ([#6134](https://github.com/terraform-providers/terraform-provider-google/pull/6134))
* **New Resource:** `google_monitoring_slo` ([#6118](https://github.com/terraform-providers/terraform-provider-google/pull/6118))
* **New Resource:** `google_service_directory_endpoint` ([#6150](https://github.com/terraform-providers/terraform-provider-google/pull/6150))
* **New Resource:** `google_service_directory_namespace` ([#6150](https://github.com/terraform-providers/terraform-provider-google/pull/6150))
* **New Resource:** `google_service_directory_service` ([#6150](https://github.com/terraform-providers/terraform-provider-google/pull/6150))

IMPROVEMENTS:
* bigtable: Reduced the minimum number of nodes for the `bigtable_instace` resource from 3 to 1. ([#6159](https://github.com/terraform-providers/terraform-provider-google/pull/6159))
* container: Added support for `google_container_cluster` Compute Engine persistent disk CSI driver ([#6160](https://github.com/terraform-providers/terraform-provider-google/pull/6160))
* compute: Added support for `google_compute_instance` `resource_policies` field ([#6125](https://github.com/terraform-providers/terraform-provider-google/pull/6125))
* compute: Added support for `google_compute_resource_policy` group placement policies ([#6125](https://github.com/terraform-providers/terraform-provider-google/pull/6125))

BUG FIXES:
* dataproc: Fixed diff when `google_dataproc_cluster` `preemptible_worker_config.0.num_instances` is sized to 0 and other `preemptible_worker_config` subfields are set ([#6123](https://github.com/terraform-providers/terraform-provider-google/pull/6123))
* resourcemanager: Added a wait to `google_project` so that projects are more likely to be ready before the resource finishes creation ([#6161](https://github.com/terraform-providers/terraform-provider-google/pull/6161))
* sql: Allowed `binary_log_enabled` to be disabled. ([#6163](https://github.com/terraform-providers/terraform-provider-google/pull/6163))
* sql: Fixed behaviour in `google_sql_database` when the parent instance is deleted, removing it from state ([#6162](https://github.com/terraform-providers/terraform-provider-google/pull/6162))

## 3.18.0 (April 20, 2020)

FEATURES:
* **New Data Source:** `google_monitoring_app_engine_service` ([#6078](https://github.com/terraform-providers/terraform-provider-google/pull/6078))
* **New Resource:** `google_monitoring_custom_service` ([#6078](https://github.com/terraform-providers/terraform-provider-google/pull/6078))
* **New Resource:** `google_compute_global_network_endpoint` ([#6095](https://github.com/terraform-providers/terraform-provider-google/pull/6095))
* **New Resource:** `google_compute_global_network_endpoint_group` ([#6095](https://github.com/terraform-providers/terraform-provider-google/pull/6095))
* **New Resource:** `google_monitoring_slo` ([#6118](https://github.com/terraform-providers/terraform-provider-google/pull/6118))

IMPROVEMENTS:
* appengine: Added `iap.enabled` field to `google_app_engine_application` resource ([#6076](https://github.com/terraform-providers/terraform-provider-google/pull/6076))
* iam: Added `name` field to `google_organization_iam_custom_role` ([#6111](https://github.com/terraform-providers/terraform-provider-google/pull/6111))
* iam: Added `name` field to `google_project_iam_custom_role` ([#6111](https://github.com/terraform-providers/terraform-provider-google/pull/6111))

BUG FIXES:
* container: Fixed importing/reading `google_container_node_pool` resources in non-RUNNING states ([#6115](https://github.com/terraform-providers/terraform-provider-google/pull/6115))
* monitoring: Made `display_name` optional on `google_monitoring_notification_channel ` ([#6090](https://github.com/terraform-providers/terraform-provider-google/pull/6090))

## 3.17.0 (April 13, 2020)

FEATURES:
* **New Resource:** `google_bigquery_dataset_access` ([#6035](https://github.com/terraform-providers/terraform-provider-google/pull/6035))
* **New Resource:** `google_dialogflow_intent` ([#6061](https://github.com/terraform-providers/terraform-provider-google/pull/6061))
* **New Resource:** `google_os_login_ssh_public_key` ([#6026](https://github.com/terraform-providers/terraform-provider-google/pull/6026))

IMPROVEMENTS:
* accesscontextmanager: added `spec` and `use_explicit_dry_run_spec` to `google_access_context_manager_service_perimeter` to test perimeter configurations in dry-run mode. ([#6071](https://github.com/terraform-providers/terraform-provider-google/pull/6071))
* compute: Added update support for `google_compute_interconnect_attachment` `admin_enabled` ([#6046](https://github.com/terraform-providers/terraform-provider-google/pull/6046))
* compute: Added field `log_config` to `google_compute_health_check` and `google_compute_region_health_check` to enable health check logging. ([#6058](https://github.com/terraform-providers/terraform-provider-google/pull/6058))
* compute: Added more import formats for `google_compute_instance` ([#6023](https://github.com/terraform-providers/terraform-provider-google/pull/6023))
* sourcerepo: allowed `google_sourcerepo_repo` `pubsub_configs.topic` to accept short topic names in addition to full references. ([#6069](https://github.com/terraform-providers/terraform-provider-google/pull/6069))

BUG FIXES:
* compute: Fixed diff on default value for `google_compute_interconnect_attachment` `admin_enabled` ([#6046](https://github.com/terraform-providers/terraform-provider-google/pull/6046))
* compute: Fixed perma-diff on `google_compute_interconnect_attachment` `candidate_subnets` ([#6046](https://github.com/terraform-providers/terraform-provider-google/pull/6046))
* compute: fixed bug where `google_compute_instance_from_template` instance defaults were overriding `scheduling` ([#6070](https://github.com/terraform-providers/terraform-provider-google/pull/6070))
* iap: `project` can now be unset in `iap_web_iam_member` and will read from the default `project` ([#6060](https://github.com/terraform-providers/terraform-provider-google/pull/6060))
* serviceusage: fixed issue where `google_project_services` attempted to read a project before enabling the API that allows that read ([#6062](https://github.com/terraform-providers/terraform-provider-google/pull/6062))
* sql: fixed error that occurred on `google_sql_database_instance` when `settings.ip_configuration` was set but `ipv4_enabled` was not set to true and `private_network` was not configured, by defaulting `ipv4_enabled` to true. ([#6041](https://github.com/terraform-providers/terraform-provider-google/pull/6041))
* storage: fixed an issue where `google_storage_bucket_iam_member` showed a diff for bucket self links ([#6019](https://github.com/terraform-providers/terraform-provider-google/pull/6019))
* storage: fixed bug where deleting a `google_storage_bucket` that contained non-deletable objects would retry indefinitely ([#6044](https://github.com/terraform-providers/terraform-provider-google/pull/6044))

## 3.16.0 (April 06, 2020)
FEATURES:
* **New Data Source:** `google_monitoring_uptime_check_ips` ([#6009](https://github.com/terraform-providers/terraform-provider-google/pull/6009))

IMPROVEMENTS:
* cloudfunctions: Added `ingress_settings` field to `google_cloudfunctions_function` ([#5981](https://github.com/terraform-providers/terraform-provider-google/pull/5981))
* cloudfunctions: added support for `vpc_connector_egress_settings` to `google_cloudfunctions_function` ([#5984](https://github.com/terraform-providers/terraform-provider-google/pull/5984))
* accesscontextmanager: added `status.vpc_accessible_services` to `google_access_context_manager_service_perimeter` to control which services are available from the perimeter's VPC networks to the restricted Google APIs IP address range. ([#6006](https://github.com/terraform-providers/terraform-provider-google/pull/6006))
* cloudrun: added ability to autogenerate revision name ([#5987](https://github.com/terraform-providers/terraform-provider-google/pull/5987))
* compute: added ability to resize `google_compute_reservation` ([#5999](https://github.com/terraform-providers/terraform-provider-google/pull/5999))
* container: added `resource_usage_export_config` to `google_container_cluster`, previously only available in `google-beta` ([#5990](https://github.com/terraform-providers/terraform-provider-google/pull/5990))
* dns: added ability to update `google_dns_managed_zone.dnssec_config` ([#6011](https://github.com/terraform-providers/terraform-provider-google/pull/6011))
* pubsub: Added `dead_letter_policy` support to `google_pubsub_subscription` ([#6010](https://github.com/terraform-providers/terraform-provider-google/pull/6010))

BUG FIXES:
* compute: Fixed an issue where `port` could not be removed from health checks ([#5997](https://github.com/terraform-providers/terraform-provider-google/pull/5997))
* storage: fixed an issue where `google_storage_bucket_iam_member` showed a diff for bucket self links ([#6019](https://github.com/terraform-providers/terraform-provider-google/pull/6019))

## 3.15.0 (March 30, 2020)

FEATURES:
* **New Resource:** `google_compute_instance_group_named_port` ([#5932](https://github.com/terraform-providers/terraform-provider-google/pull/5932))
* **New Resource:** `google_service_usage_consumer_quota_override` ([#5966](https://github.com/terraform-providers/terraform-provider-google/pull/5966))
* **New Resource:** `google_iap_brand` ([#5881](https://github.com/terraform-providers/terraform-provider-google/pull/5881))
* **New Resource:** `google_iap_client` ([#5881](https://github.com/terraform-providers/terraform-provider-google/pull/5881))
* **New Resource:** `google_appengine_flexible_app_version` ([#5882](https://github.com/terraform-providers/terraform-provider-google/pull/5882))

IMPROVEMENTS:
* accesscontextmanager: Added `regions` field to `google_access_context_manager_access_level` ([#5961](https://github.com/terraform-providers/terraform-provider-google/pull/5961))
* compute: added field `network` to `google_compute_region_backend_service`, which allows internal load balancers to target the non-primary interface of an instance. ([#5957](https://github.com/terraform-providers/terraform-provider-google/pull/5957))
* compute: added support for IAM conditions in `google_compute_subnet_iam_*` IAM resources ([#5954](https://github.com/terraform-providers/terraform-provider-google/pull/5954))
* container: Added field `maintenance_policy.recurring_window` to  `google_container_cluster` ([#5962](https://github.com/terraform-providers/terraform-provider-google/pull/5962))
* kms: Added new field `additional_authenticated_data` for Cloud KMS data source `google_kms_secret` ([#5968](https://github.com/terraform-providers/terraform-provider-google/pull/5968))
* kms: Added new field `additional_authenticated_data` for Cloud KMS resource `google_kms_secret_ciphertext` ([#5968](https://github.com/terraform-providers/terraform-provider-google/pull/5968))

BUG FIXES:
* kms: Fixed an issue in `google_kms_crypto_key_version` where `public_key` would return empty after apply ([#5956](https://github.com/terraform-providers/terraform-provider-google/pull/5956))
* logging: Fixed import issue with `google_logging_metric` in a non-default project. ([#5944](https://github.com/terraform-providers/terraform-provider-google/pull/5944))
* provider: Fixed an error with resources failing to upload large files (e.g. with `google_storage_bucket_object`) during retried requests ([#5977](https://github.com/terraform-providers/terraform-provider-google/pull/5977))

## 3.14.0 (March 23, 2020)

FEATURES:
* **New Data Source:** `google_compute_instance_serial_port` ([#5911](https://github.com/terraform-providers/terraform-provider-google/pull/5911))
* **New Resource:** `google_compute_region_ssl_certificate` ([#5913](https://github.com/terraform-providers/terraform-provider-google/pull/5913))

IMPROVEMENTS:
* compute: Added new attribute reference `current_status` to the `google_compute_instance` resource ([#5903](https://github.com/terraform-providers/terraform-provider-google/pull/5903))
* compute: Added `allow_global_access` to `google_compute_forwarding_rule` resource. ([#5912](https://github.com/terraform-providers/terraform-provider-google/pull/5912))
* container: Added `dns_cache_config` field to `google_container_cluster` resource ([#5887](https://github.com/terraform-providers/terraform-provider-google/pull/5887))
* container: Added `upgrade_settings` to `google_container_node_pool` resource ([#5910](https://github.com/terraform-providers/terraform-provider-google/pull/5910))
* provider: Added provider-wide request retries for common temporary GCP error codes and network errors ([#5902](https://github.com/terraform-providers/terraform-provider-google/pull/5902))
* redis: Added `connect_mode` field to `google_redis_instance` resource ([#5888](https://github.com/terraform-providers/terraform-provider-google/pull/5888))

## 3.13.0 (March 16, 2020)

BREAKING CHANGES:
* dialogflow: Changed `google_dialogflow_agent.time_zone` to ForceNew. Updating this field will require recreation. This is due to a change in API behavior. ([#5831](https://github.com/terraform-providers/terraform-provider-google/pull/5831))

FEATURES:
* **New Resource:** `google_compute_region_disk_resource_policy_attachment` ([#5849](https://github.com/terraform-providers/terraform-provider-google/pull/5849))
* **New Resource:** `google_sql_source_representation_instance` ([#5839](https://github.com/terraform-providers/terraform-provider-google/pull/5839))

IMPROVEMENTS:
* bigtable: Added support for full-name/id `instance_name` value in `google_bigtable_table` and `google_bigtable_gc_policy` ([#5837](https://github.com/terraform-providers/terraform-provider-google/pull/5837))
* compute: Added `autoscaling_policy` to `google_compute_node_group` ([#5864](https://github.com/terraform-providers/terraform-provider-google/pull/5864))
* compute: Added support for full-name/id `network_endpoint_group` value in `google_network_endpoint` ([#5838](https://github.com/terraform-providers/terraform-provider-google/pull/5838))
* compute: Added support for `google_compute_router_nat` `drain_nat_ips` (previously beta-only). ([#5821](https://github.com/terraform-providers/terraform-provider-google/pull/5821))
* dialogflow: Changed `google_dialogflow_agent` to not read `tier` status ([#5835](https://github.com/terraform-providers/terraform-provider-google/pull/5835))
* monitoring: Added `sensitive_labels` to `google_monitoring_notification_channel` so that labels like `password` and `auth_token` can be managed separately from the other labels and marked as sensitive. ([#5873](https://github.com/terraform-providers/terraform-provider-google/pull/5873))

BUG FIXES:
* all: fixed issue where nested objects were getting sent as null values to GCP on create instead of being omitted from requests ([#5825](https://github.com/terraform-providers/terraform-provider-google/pull/5825))
* cloudfunctions: fixed `vpc_connector` to be updated properly in `google_cloudfunctions_function` ([#5829](https://github.com/terraform-providers/terraform-provider-google/pull/5829))
* compute: fixed `google_compute_security_policy` from allowing two rules with the same priority. ([#5834](https://github.com/terraform-providers/terraform-provider-google/pull/5834))
* compute: fixed bug where `google_compute_instance.scheduling.node_affinities.operator` would incorrectly accept `NOT` rather than `NOT_IN`. ([#5841](https://github.com/terraform-providers/terraform-provider-google/pull/5841))
* container: Fixed issue where `google_container_node_pool` resources created in the 2.X series were failing to update after 3.11. ([#5877](https://github.com/terraform-providers/terraform-provider-google/pull/5877))
## 3.12.0 (March 09, 2020)
IMPROVEMENTS:
* serviceusage: `google_project_service` no longer attempts to enable a service that is already enabled. ([#5810](https://github.com/terraform-providers/terraform-provider-google/pull/5810))
* bigtable: Added support for full-name/id `instance` value in `google_bigtable_app_profile` ([#5780](https://github.com/terraform-providers/terraform-provider-google/pull/5780))
* compute: Added `google_compute_router_nat` `drain_nat_ips` field (formerly beta). ([#5821](https://github.com/terraform-providers/terraform-provider-google/pull/5821))
* pubsub: Added polling to ensure correct resource state for negative-cached PubSub resources ([#5813](https://github.com/terraform-providers/terraform-provider-google/pull/5813))

BUG FIXES:
* compute: Fixed a scenario where `google_compute_instance_template` would cause a crash. ([#5808](https://github.com/terraform-providers/terraform-provider-google/pull/5808))
* container: Fixed panic when upgrading `google_container_cluster` with autoscaling block ([#5782](https://github.com/terraform-providers/terraform-provider-google/pull/5782))
* storage: Added check for bucket retention policy list being empty. ([#5793](https://github.com/terraform-providers/terraform-provider-google/pull/5793))
* storage: Added locking for operations involving `google_storage_*_access_control` resources to prevent errors from ACLs being added at the same time. ([#5791](https://github.com/terraform-providers/terraform-provider-google/pull/5791))

## 3.11.0 (March 02, 2020)

FEATURES:
* **New Data Source:** `google_compute_backend_bucket` ([#5720](https://github.com/terraform-providers/terraform-provider-google/pull/5720))
* **New Resource:** `google_app_engine_service_split_traffic` ([#5729](https://github.com/terraform-providers/terraform-provider-google/pull/5729))
* **New Resource:** `google_compute_packet_mirroring` ([#5755](https://github.com/terraform-providers/terraform-provider-google/pull/5755))
* **New Resource:** `google_vpc_access_connector` (GA provider) ([#5752](https://github.com/terraform-providers/terraform-provider-google/pull/5752))

IMPROVEMENTS:
* bigquery: Landed support for range-based partitioning in `google_bigquery_table` ([#5723](https://github.com/terraform-providers/terraform-provider-google/pull/5723))
* compute: added check on `google_compute_router` for non-empty advertised_groups or advertised_ip_ranges values when advertise_mode is DEFAULT in the bgp block. ([#5718](https://github.com/terraform-providers/terraform-provider-google/pull/5718))
* compute: added the ability to manage the status of `google_compute_instance` resources with the `desired_status` field ([#4797](https://github.com/terraform-providers/terraform-provider-google/pull/4797))
* iam: `google_project_iam_member` and `google_project_iam_binding`'s `project` field can be specified with an optional `projects/` prefix ([#5722](https://github.com/terraform-providers/terraform-provider-google/pull/5722))
* storage: added `metadata` to `google_storage_bucket_object`. ([#5721](https://github.com/terraform-providers/terraform-provider-google/pull/5721))

BUG FIXES:
* compute: Updated `google_project` to check for valid permissions on the parent billing account before creating and tainting the resource. ([#5719](https://github.com/terraform-providers/terraform-provider-google/pull/5719))
* container: Fixed panic when upgrading `google_container_cluster` with `autoscaling` block ([#5782](https://github.com/terraform-providers/terraform-provider-google/pull/5782))

## 3.10.0 (February 25, 2020)

BREAKING CHANGES:
* container: Fully removed `use_ip_aliases` and `create_subnetwork` fields to fix misleading diff for removed fields ([#5666](https://github.com/terraform-providers/terraform-provider-google/pull/5666))

FEATURES:
* **New Data Source:** `google_dns_keys` ([#5703](https://github.com/terraform-providers/terraform-provider-google/pull/5703))
* **New Resource:** `google_storage_hmac_key` ([#5679](https://github.com/terraform-providers/terraform-provider-google/pull/5679))
* **New Resource:** `google_datastore_index` ([#5655](https://github.com/terraform-providers/terraform-provider-google/pull/5655))
* **New Resource:** `google_endpoints_service_iam_binding` ([#5668](https://github.com/terraform-providers/terraform-provider-google/pull/5668))
* **New Resource:** `google_endpoints_service_iam_member` ([#5668](https://github.com/terraform-providers/terraform-provider-google/pull/5668))
* **New Resource:** `google_endpoints_service_iam_policy` ([#5668](https://github.com/terraform-providers/terraform-provider-google/pull/5668))

IMPROVEMENTS:
* container: Allowed import/update/deletion of `google_container_cluster` in error states. ([#5663](https://github.com/terraform-providers/terraform-provider-google/pull/5663))
* container: Changed `google_container_node_pool` so node pools created in an error state will be marked as tainted on creation. ([#5662](https://github.com/terraform-providers/terraform-provider-google/pull/5662))
* container: Allowed import/update/deletion of `google_container_node_pool` in error states and updated resource to wait for a stable state after any changes. ([#5662](https://github.com/terraform-providers/terraform-provider-google/pull/5662))
* container: added label_fingerprint to `google_container_cluster` ([#5647](https://github.com/terraform-providers/terraform-provider-google/pull/5647))
* container: Enabled configuring autoscaling profile in `google_container_cluster` (https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-autoscaler#autoscaling_profiles) ([#5659](https://github.com/terraform-providers/terraform-provider-google/pull/5659))
* dataflow: added `job_id` attribute ([#5644](https://github.com/terraform-providers/terraform-provider-google/pull/5644))
* dataflow: added computed `type` field to `google_dataflow_job`. ([#5709](https://github.com/terraform-providers/terraform-provider-google/pull/5709))
* provider: Added retries for common network errors we've encountered. ([#5675](https://github.com/terraform-providers/terraform-provider-google/pull/5675))

## 3.9.0 (February 18, 2020)

FEATURES:
* **New Resource:** `google_container_registry` ([#5593](https://github.com/terraform-providers/terraform-provider-google/pull/5593))

IMPROVEMENTS:
* all: improve error handling of 404s. ([#5601](https://github.com/terraform-providers/terraform-provider-google/pull/5601))
* bigtable: added update support for `display_name` and `instance_type` ([#5648](https://github.com/terraform-providers/terraform-provider-google/pull/5648))
* container: `google_container_cluster` will wait for a stable state after updates. ([#5616](https://github.com/terraform-providers/terraform-provider-google/pull/5616))
* container: added `boot_disk_kms_key` to `node_config` block. ([#5615](https://github.com/terraform-providers/terraform-provider-google/pull/5615))
* dataflow: added `job_id` field to `google_dataflow_job` ([#5653](https://github.com/terraform-providers/terraform-provider-google/pull/5653))
* dialogflow: improve error handling by increasing retry count ([#5603](https://github.com/terraform-providers/terraform-provider-google/pull/5603))
* resourcemanager: fixed retry behavior for updates in `google_project`, added retries for billing metadata requests ([#5578](https://github.com/terraform-providers/terraform-provider-google/pull/5578))
* sql: add `encryption_key_name` to `google_sql_database_instance` ([#5591](https://github.com/terraform-providers/terraform-provider-google/pull/5591))

BUG FIXES:
* cloudrun: fixed permadiff caused by new API default values on `annotations` and `limits` ([#5600](https://github.com/terraform-providers/terraform-provider-google/pull/5600))
* compute: Fixed bug where `google_project` would fail to create if the `auto_create_network` was false and the `compute-skipDefaultNetworkCreation` organization policies was enforced. ([#5601](https://github.com/terraform-providers/terraform-provider-google/pull/5601))
* container: Removed restriction on `auto_provisioning_defaults` to allow both `oauth_scopes` and `service_account` to be set ([#5642](https://github.com/terraform-providers/terraform-provider-google/pull/5642))
* firestore: fixed import of `google_firestore_index` when database or collection were non-default. ([#5626](https://github.com/terraform-providers/terraform-provider-google/pull/5626))
* iam: Fixed an erroneous error during import of IAM resources when a provider default project/zone/region is not defined. ([#5613](https://github.com/terraform-providers/terraform-provider-google/pull/5613))
* kms: Fixed issue where `google_kms_crypto_key_version` datasource would throw an Invalid Index error on plan ([#5619](https://github.com/terraform-providers/terraform-provider-google/pull/5619))

## 3.8.0 (February 10, 2020)

NOTES:
* provider: added documentation for the `id` field for many resources, including format ([#5543](https://github.com/terraform-providers/terraform-provider-google/pull/5543))
BREAKING CHANGES:
* compute: Added conditional requirement of `google_compute_**region**_backend_service` `backend.capacity_scaler` to no longer accept the API default if not INTERNAL. Non-INTERNAL backend services must now specify `capacity_scaler` explicitly and have a total capacity greater than 0. In addition, API default of 1.0 must now be explicitly set and will be treated as nil or zero if not set in config. ([#5561](https://github.com/terraform-providers/terraform-provider-google/pull/5561))

FEATURES:
* **New Data Source:** `secret_manager_secret_version` ([#5562](https://github.com/terraform-providers/terraform-provider-google/pull/5562))
* **New Resource:** `google_access_context_manager_service_perimeter_resource` ([#5574](https://github.com/terraform-providers/terraform-provider-google/pull/5574))
* **New Resource:** `secret_manager_secret_version` ([#5562](https://github.com/terraform-providers/terraform-provider-google/pull/5562))
* **New Resource:** `secret_manager_secret` ([#5562](https://github.com/terraform-providers/terraform-provider-google/pull/5562))
* **New Resource:** `google_dialogflow_agent` ([#5559](https://github.com/terraform-providers/terraform-provider-google/pull/5559))

IMPROVEMENTS:
* appengine: added support for `google_app_engine_application.iap` ([#5556](https://github.com/terraform-providers/terraform-provider-google/pull/5556))
* compute: `google_compute_security_policy` `rule.match.expr` field is now GA ([#5532](https://github.com/terraform-providers/terraform-provider-google/pull/5532))
* compute: added additional validation to `google_cloud_router`'s `bgp.asn` field. ([#5547](https://github.com/terraform-providers/terraform-provider-google/pull/5547))

BUG FIXES:
* bigtable: fixed diff for DEVELOPMENT instances that are returned from the API with one node ([#5557](https://github.com/terraform-providers/terraform-provider-google/pull/5557))
* compute: Fixed `backend.capacity_scaler` to actually set zero (0.0) value. ([#5561](https://github.com/terraform-providers/terraform-provider-google/pull/5561))
* compute: Fixed `google_compute_**region**_backend_service` so it no longer has a permadiff if `backend.capacity_scaler` is unset in config by requiring capacity scaler. ([#5561](https://github.com/terraform-providers/terraform-provider-google/pull/5561))
* compute: updated `google_compute_project_metadata_item` to fail on create if its key is already present in the project metadata. ([#5576](https://github.com/terraform-providers/terraform-provider-google/pull/5576))
* logging: updated `bigquery_options` so the default value from the api will be set in state. ([#5534](https://github.com/terraform-providers/terraform-provider-google/pull/5534))
* sql: undeprecated `settings.ip_configuration.authorized_networks.expiration_time` ([#5531](https://github.com/terraform-providers/terraform-provider-google/pull/5531))

## 3.7.0 (February 03, 2020)

BREAKING CHANGES:
* iam: starts reading/writing IAM policies at version 3 in the GA provider. If you have an IAM resource defined in your config that has a condition on it created outside of Terraform, you should start using the beta provider and defining the condition in your config to avoid unexpected behavior. ([#5469](https://github.com/terraform-providers/terraform-provider-google/pull/5469))

IMPROVEMENTS:
* dns: `google_dns_managed_zone` added support for Non-RFC1918 fields for reverse lookup and fowarding paths. ([#5493](https://github.com/terraform-providers/terraform-provider-google/pull/5493))
* monitoring: Added `labels` and `user_labels` filters to data source `google_monitoring_notification_channel` ([#5470](https://github.com/terraform-providers/terraform-provider-google/pull/5470))

BUG FIXES:
* bigtable: fixed diff for DEVELOPMENT instances that are returned from the API with one node ([#5557](https://github.com/terraform-providers/terraform-provider-google/pull/5557))
* compute: `google_compute_instance_template` added plan time check for any disks marked `boot` outside of the first disk ([#5491](https://github.com/terraform-providers/terraform-provider-google/pull/5491))
* container: Fixed perma-diff in `google_container_cluster`'s `cluster_autoscaling.auto_provisioning_defaults`. ([#5486](https://github.com/terraform-providers/terraform-provider-google/pull/5486))
* iam: fixed issue where users of the GA provider who used IAM conditions outside of Terraform were getting an error ([#5469](https://github.com/terraform-providers/terraform-provider-google/pull/5469))
* logging: updated `bigquery_options` so the default value from the api will be set in state. ([#5534](https://github.com/terraform-providers/terraform-provider-google/pull/5534))
* storage: Stopped `project-owner` showing up in the diff for `google_storage_bucket_acl` ([#5479](https://github.com/terraform-providers/terraform-provider-google/pull/5479))

## 3.6.0 (January 29, 2020)

KNOWN ISSUES:

* bigtable: due to API changes, bigtable DEVELOPMENT instances may show a diff on `num_nodes`. There will be a fix in the 3.7.0 release of the provider. No known workarounds exist at the moment, but will be tracked in https://github.com/terraform-providers/terraform-provider-google/issues/5492.

FEATURES:
* **New Data Source:** google_monitoring_notification_channel ([#5405](https://github.com/terraform-providers/terraform-provider-google/pull/5405))
* **New Resource:** Added `google_iap_tunnel_instance_iam_*` IAM resources for IAP Tunnel Instances ([#5429](https://github.com/terraform-providers/terraform-provider-google/pull/5429))
* **New Resource:** google_compute_network_peering_routes_config ([#5426](https://github.com/terraform-providers/terraform-provider-google/pull/5426))

IMPROVEMENTS:
* compute: added waiting logic to `google_compute_interconnect_attachment` to avoid modifications when the attachment is UNPROVISIONED ([#5459](https://github.com/terraform-providers/terraform-provider-google/pull/5459))
* compute: made the `google_compute_network_peering` routes fields available in GA ([#5419](https://github.com/terraform-providers/terraform-provider-google/pull/5419))
* container: Promoted `enable_binary_authorization` from beta into ga. ([#5456](https://github.com/terraform-providers/terraform-provider-google/pull/5456))
* scheduler: Added `attempt_deadline` to `google_cloud_scheduler_job`. ([#5399](https://github.com/terraform-providers/terraform-provider-google/pull/5399))
* storage: added `default_event_based_hold` to `google_storage_bucket` ([#5373](https://github.com/terraform-providers/terraform-provider-google/pull/5373))

BUG FIXES:
* compute: Fixed `google_compute_instance_from_template` with existing boot disks ([#5430](https://github.com/terraform-providers/terraform-provider-google/pull/5430))
* compute: Fixed a bug in `google_compute_instance` when attempting to update a field that requires stopping and starting an instance with an encrypted disk ([#5436](https://github.com/terraform-providers/terraform-provider-google/pull/5436))

## 3.5.0 (January 22, 2020)

DEPRECATIONS:
* kms: deprecated `data.google_kms_secret_ciphertext` as there was no way to make it idempotent. Instead, use the `google_kms_secret_ciphertext` resource. ([#5314](https://github.com/terraform-providers/terraform-provider-google/pull/5314))
* sql: deprecated first generation-only fields on `google_sql_database_instance` ([#5376](https://github.com/terraform-providers/terraform-provider-google/pull/5376))

FEATURES:
* **New Resource:** `google_kms_secret_ciphertext` ([#5314](https://github.com/terraform-providers/terraform-provider-google/pull/5314))

IMPROVEMENTS:
* bigtable: added the ability to add/remove clusters from `google_bigtable_instance` ([#5318](https://github.com/terraform-providers/terraform-provider-google/pull/5318))
* compute: added support for other resource types (like a Proxy) as a `target` to `google_compute_forwarding_rule`. ([#5383](https://github.com/terraform-providers/terraform-provider-google/pull/5383))
* dataproc: added `lifecycle_config` to `google_dataproc_cluster.cluster_config` ([#5323](https://github.com/terraform-providers/terraform-provider-google/pull/5323))
* iam: updated to allow for empty bindings in `data_source_google_iam_policy` data source ([#4525](https://github.com/terraform-providers/terraform-provider-google/pull/4525))
* provider: added retries for batched requests so failed batches will retry each single request separately. ([#5355](https://github.com/terraform-providers/terraform-provider-google/pull/5355))
* resourcemanager: restricted the length of the `description` field of `google_service_account`. It is now limited to 256 characters. ([#5409](https://github.com/terraform-providers/terraform-provider-google/pull/5409))

BUG FIXES:
* bigtable: Fixed error on reading non-existent `google_bigtable_gc_policy`,  `google_bigtable_instance`,  `google_bigtable_table` ([#5331](https://github.com/terraform-providers/terraform-provider-google/pull/5331))
* cloudfunctions: Fixed validation of `google_cloudfunctions_function` name to allow for 63 characters. ([#5400](https://github.com/terraform-providers/terraform-provider-google/pull/5400))
* cloudtasks: Changed `max_dispatches_per_second` to a double instead of an integer. ([#5393](https://github.com/terraform-providers/terraform-provider-google/pull/5393))
* compute: Added validation for `compute_resource_policy` to no longer allow invalid `start_time` values that weren't hourly. ([#5342](https://github.com/terraform-providers/terraform-provider-google/pull/5342))
* compute: Fixed errors from concurrent creation/deletion of overlapping `google_compute_network_peering` resources. ([#5338](https://github.com/terraform-providers/terraform-provider-google/pull/5338))
* compute: Stopped panic when using `usage_export_bucket` and the setting had been disabled manually. ([#5349](https://github.com/terraform-providers/terraform-provider-google/pull/5349))
* compute: fixed `google_compute_router_nat` timeout fields causing a diff when using a long-lived resource ([#5353](https://github.com/terraform-providers/terraform-provider-google/pull/5353))
* compute: fixed `google_compute_target_https_proxy.quic_override` causing a diff when using a long-lived resource ([#5351](https://github.com/terraform-providers/terraform-provider-google/pull/5351))
* identityplatform: fixed `google_identity_platform_default_supported_idp_config` to correctly allow configuration of both `idp_id` and `client_id` separately ([#5398](https://github.com/terraform-providers/terraform-provider-google/pull/5398))
* monitoring: Stopped `labels` from causing a perma diff on `AlertPolicy` ([#5367](https://github.com/terraform-providers/terraform-provider-google/pull/5367))

## 3.4.0 (January 07, 2020)

DEPRECATIONS:
* kms: deprecated `data.google_kms_secret_ciphertext` as there was no way to make it idempotent. Instead, use the `google_kms_secret_ciphertext` resource. ([#5314](https://github.com/terraform-providers/terraform-provider-google/pull/5314))

BREAKING CHANGES:
* cloudrun: Changed `google_cloud_run_domain_mapping` to correctly match Cloud Run API expected format for `spec.route_name`, {serviceName}, instead of invalid projects/{project}/global/services/{serviceName} ([#5264](https://github.com/terraform-providers/terraform-provider-google/pull/5264))
* compute: Added back ConflictsWith restrictions for ExactlyOneOf restrictions that were removed in v3.3.0 for `google_compute_firewall`, `google_compute_health_check`, and `google_compute_region_health_check`. This effectively changes an API-side failure that was only accessible in v3.3.0 to a plan-time one. ([#5220](https://github.com/terraform-providers/terraform-provider-google/pull/5220))
* logging: Changed `google_logging_metric.metric_descriptors.labels` from a list to a set ([#5258](https://github.com/terraform-providers/terraform-provider-google/pull/5258))
* resourcemanager: Added back ConflictsWith restrictions for ExactlyOneOf restrictions that were removed in v3.3.0 for `google_organization_policy`, `google_folder_organization_policy`, and `google_project_organization_policy`. This effectively changes an API-side failure that was only accessible in v3.3.0 to a plan-time one. ([#5220](https://github.com/terraform-providers/terraform-provider-google/pull/5220))

FEATURES:
* **New Data Source:** google_sql_ca_certs ([#5306](https://github.com/terraform-providers/terraform-provider-google/pull/5306))
* **New Resource:** `google_identity_platform_default_supported_idp_config` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_identity_platform_inbound_saml_config` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_identity_platform_oauth_idp_config` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_identity_platform_tenant_default_supported_idp_config` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_identity_platform_tenant_inbound_saml_config` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_identity_platform_tenant_oauth_idp_config` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_identity_platform_tenant` ([#5199](https://github.com/terraform-providers/terraform-provider-google/pull/5199))
* **New Resource:** `google_kms_crypto_key_iam_policy` ([#5247](https://github.com/terraform-providers/terraform-provider-google/pull/5247))
* **New Resource:** `google_kms_secret_ciphertext` ([#5314](https://github.com/terraform-providers/terraform-provider-google/pull/5314))

IMPROVEMENTS:
* composer: Increased default timeouts for `google_composer_environment` ([#5223](https://github.com/terraform-providers/terraform-provider-google/pull/5223))
* compute: Added graceful termination to `container_cluster` create calls so that partially created clusters will resume the original operation if the Terraform process is killed mid create. ([#5217](https://github.com/terraform-providers/terraform-provider-google/pull/5217))
* compute: Fixed `google_compute_disk_resource_policy_attachment` parsing of region from zone to allow for provider-level zone and make error message more accurate` ([#5257](https://github.com/terraform-providers/terraform-provider-google/pull/5257))
* provider: Reduced default `send_after` controlling the time interval after which a batched request sends. ([#5268](https://github.com/terraform-providers/terraform-provider-google/pull/5268))

BUG FIXES:
* all: fixed issue where many fields that were removed in 3.0.0 would show a diff when they were removed from config ([#5313](https://github.com/terraform-providers/terraform-provider-google/pull/5313))
* bigquery: fixed `bigquery_table.encryption_configuration` to correctly recreate the table when modified ([#5321](https://github.com/terraform-providers/terraform-provider-google/pull/5321))
* cloudrun:  Changed `google_cloud_run_domain_mapping` to correctly match Cloud Run API expected format for `spec.route_name`, {serviceName}, instead of invalid projects/{project}/global/services/{serviceName} ([#5264](https://github.com/terraform-providers/terraform-provider-google/pull/5264))
* cloudrun: Changed `cloud_run_domain_mapping` to poll for success or failure and throw an appropriate error when ready status returns as false. ([#5267](https://github.com/terraform-providers/terraform-provider-google/pull/5267))
* cloudrun: Fixed `google_cloudrun_service` to allow update instead of force-recreation for changes in `spec` `env` and `command` fields ([#5269](https://github.com/terraform-providers/terraform-provider-google/pull/5269))
* cloudrun: Removed unsupported update for `google_cloud_run_domain_mapping` to allow force-recreation. ([#5253](https://github.com/terraform-providers/terraform-provider-google/pull/5253))
* cloudrun: Stopped returning an error when a `cloud_run_domain_mapping` was waiting on DNS verification. ([#5315](https://github.com/terraform-providers/terraform-provider-google/pull/5315))
* compute: Fixed `google_compute_backend_service` to allow updating `cdn_policy.cache_key_policy.*` fields to false or empty. ([#5276](https://github.com/terraform-providers/terraform-provider-google/pull/5276))
* compute: Fixed behaviour where `google_compute_subnetwork` did not record a value for `name` when `self_link` was specified. ([#5288](https://github.com/terraform-providers/terraform-provider-google/pull/5288))
* container: fixed issue where an empty variable in `tags` would cause a crash ([#5226](https://github.com/terraform-providers/terraform-provider-google/pull/5226))
* endpoints: Added operation wait for `google_endpoints_service` to fix 403 "Service not found" errors during initial creation ([#5259](https://github.com/terraform-providers/terraform-provider-google/pull/5259))
* logging: Made `google_logging_metric.metric_descriptors.labels` a set to prevent diff from ordering ([#5258](https://github.com/terraform-providers/terraform-provider-google/pull/5258))
* resourcemanager: added retries for `data.google_organization` ([#5246](https://github.com/terraform-providers/terraform-provider-google/pull/5246))

## 3.3.0 (December 17, 2019)

FEATURES:
* **New Resource:** `google_compute_region_health_check` is now available in GA ([#5149](https://github.com/terraform-providers/terraform-provider-google/pull/5149))
* **New Resource:** `google_deployment_manager_deployment` ([#5139](https://github.com/terraform-providers/terraform-provider-google/pull/5139))

IMPROVEMENTS:
* bigquery: added `PARQUET` as an option in `google_bigquery_table.external_data_configuration.source_format` ([#5170](https://github.com/terraform-providers/terraform-provider-google/pull/5170))
* compute: Added support for `next_hop_ilb` to `google_compute_route` ([#5162](https://github.com/terraform-providers/terraform-provider-google/pull/5162))
* dataproc: added support for `security_config` to `google_dataproc_cluster` ([#5129](https://github.com/terraform-providers/terraform-provider-google/pull/5129))
* storage: updated `id` and `bucket` fields for `google_storage_bucket_iam_*` resources to use `b/{bucket_name}` ([#5099](https://github.com/terraform-providers/terraform-provider-google/pull/5099))

BUG FIXES:
* compute: Fixed an issue where interpolated values caused plan-time errors in `google_compute_router_interface`. ([#5178](https://github.com/terraform-providers/terraform-provider-google/pull/5178))
* compute: relaxed ExactlyOneOf restrictions on `google_compute_firewall`, `google_compute_health_check`, and `google_compute_region_health_check` to enable the use of dynamic blocks with those resources. ([#5194](https://github.com/terraform-providers/terraform-provider-google/pull/5194))
* iam: Fixed a bug that causes badRequest errors on IAM resources due to deleted serviceAccount principals ([#5142](https://github.com/terraform-providers/terraform-provider-google/pull/5142))
* resourcemanager: relaxed ExactlyOneOf restrictions on `google_organization_policy `, `google_folder_organization_policy `, and `google_project_organization_policy ` to enable the use of dynamic blocks with those resources. ([#5194](https://github.com/terraform-providers/terraform-provider-google/pull/5194))
* sourcerepo: Fixed a bug preventing repository IAM resources from referencing repositories with the `/` character in their name ([#5195](https://github.com/terraform-providers/terraform-provider-google/pull/5195))
* sql: fixed bug where terraform would keep retrying to create new `google_sql_database_instance` with the name of a previously deleted instance ([#5141](https://github.com/terraform-providers/terraform-provider-google/pull/5141))

## 3.2.0 (December 11, 2019)

DEPRECATIONS:
* compute: deprecated `fingerprint` field in `google_compute_subnetwork`. Its value is now always `""`. ([#5105](https://github.com/terraform-providers/terraform-provider-google/pull/5105))

FEATURES:
* **New Data Source:** `data_source_google_bigquery_default_service_account` ([#5081](https://github.com/terraform-providers/terraform-provider-google/pull/5081))
* **New Resource:** cloudrun: Added support for `google_cloud_run_service` IAM resources: `google_cloud_run_service_iam_policy`, `google_cloud_run_service_iam_binding`, `google_cloud_run_service_iam_member` ([#5051](https://github.com/terraform-providers/terraform-provider-google/pull/5051))

IMPROVEMENTS:
* all: Added `synchronous_timeout` to provider block to allow setting higher per-operation-poll timeouts. ([#5013](https://github.com/terraform-providers/terraform-provider-google/pull/5013))
* bigquery: Added KMS support to `google_bigquery_table` ([#5081](https://github.com/terraform-providers/terraform-provider-google/pull/5081))
* cloudresourcemanager: Added `org_id` field to `google_organization` datasource to expose the raw organization id ([#5115](https://github.com/terraform-providers/terraform-provider-google/pull/5115))
* cloudrun: Stopped requiring the root `metadata` block for `google_cloud_run_service`. ([#5094](https://github.com/terraform-providers/terraform-provider-google/pull/5094))
* compute: added support for `expr` to `google_compute_security_policy.rule.match` ([#5070](https://github.com/terraform-providers/terraform-provider-google/pull/5070))
* compute: added support for `path_rules` to `google_compute_region_url_map` ([#5122](https://github.com/terraform-providers/terraform-provider-google/pull/5122))
* compute: added support for `path_rules` to `google_compute_url_map` ([#5106](https://github.com/terraform-providers/terraform-provider-google/pull/5106))
* compute: added support for `route_rules` to `google_compute_region_url_map` ([#5130](https://github.com/terraform-providers/terraform-provider-google/pull/5130))
* compute: added support for header actions and route rules to `google_compute_url_map` ([#4992](https://github.com/terraform-providers/terraform-provider-google/pull/4992))
* dns: Added `visibility` field to `google_dns_managed_zone` data source ([#5063](https://github.com/terraform-providers/terraform-provider-google/pull/5063))
* sourcerepo: added support for `pubsub_configs` to `google_sourcerepo_repository` ([#5050](https://github.com/terraform-providers/terraform-provider-google/pull/5050))

BUG FIXES:
* dns: fixed 503s caused by high numbers of `dns_record_set`s. ([#5093](https://github.com/terraform-providers/terraform-provider-google/pull/5093))
* logging: updated `exponential_buckets.growth_factor` from integer to double. ([#5111](https://github.com/terraform-providers/terraform-provider-google/pull/5111))
* storage: fixed bug where users without storage.objects.list permissions couldn't delete empty buckets ([#5006](https://github.com/terraform-providers/terraform-provider-google/pull/5006))

## 3.1.0 (December 05, 2019)

BREAKING CHANGES:
* compute: field `peer_ip_address` in `google_compute_router_peer` is now required, to match the API behavior. ([#4923](https://github.com/terraform-providers/terraform-provider-google/pull/4923))

FEATURES:
* **New Resource:** `google_billing_budget` ([#5005](https://github.com/terraform-providers/terraform-provider-google/pull/5005))
* **New Resource:** `google_cloud_tasks_queue` ([#4880](https://github.com/terraform-providers/terraform-provider-google/pull/4880))
* **New Resource:** `google_organization_iam_audit_config` ([#4977](https://github.com/terraform-providers/terraform-provider-google/pull/4977))

IMPROVEMENTS:
* accesscontextmanager: added support for `requireAdminApproval` and `requireCorpOwned` in `google_access_context_manager_access_level`'s `devicePolicy`. ([#4931](https://github.com/terraform-providers/terraform-provider-google/pull/4931))
* all: added retries for timeouts while fetching operations ([#4605](https://github.com/terraform-providers/terraform-provider-google/pull/4605))
* cloudbuild: Added build timeout to `google_cloudbuild_trigger` ([#4938](https://github.com/terraform-providers/terraform-provider-google/pull/4938))
* cloudresourcemanager: added support for importing `google_folder` in the form of the bare folder id, rather than requiring `folders/{bare_id}` ([#4981](https://github.com/terraform-providers/terraform-provider-google/pull/4981))
* compute: Updated default timeouts on `google_compute_project_metadata_item`. ([#4995](https://github.com/terraform-providers/terraform-provider-google/pull/4995))
* compute: `google_compute_disk` `disk_encryption_key.raw_key` is now sensitive ([#5009](https://github.com/terraform-providers/terraform-provider-google/pull/5009))
* compute: `google_compute_firewall` `enable_logging` is now GA ([#4999](https://github.com/terraform-providers/terraform-provider-google/pull/4999))
* compute: `google_compute_network_peering` resource can now be imported ([#4998](https://github.com/terraform-providers/terraform-provider-google/pull/4998))
* compute: computed attribute `management_type` in `google_compute_router_peer` is now available. ([#4923](https://github.com/terraform-providers/terraform-provider-google/pull/4923))
* container: `authenticator_groups_config` in `google_container_cluster` is now GA ([#4969](https://github.com/terraform-providers/terraform-provider-google/pull/4969))
* container: `google_container_cluster.vertical_pod_autoscaling` is now GA ([#5033](https://github.com/terraform-providers/terraform-provider-google/pull/5033))
* container: added `auto_provisioning_defaults` to `google_container_cluster.cluster_autoscaling` ([#4991](https://github.com/terraform-providers/terraform-provider-google/pull/4991))
* container: added `upgrade_settings` support  to `google_container_node_pool` ([#4926](https://github.com/terraform-providers/terraform-provider-google/pull/4926))
* container: increased timeouts on `google_container_cluster` and `google_container_node_pool` ([#4902](https://github.com/terraform-providers/terraform-provider-google/pull/4902))
* dataproc: `google_dataproc_autoscaling_policy` is now GA. `google_dataproc_cluster.autoscaling_config` is also available in GA ([#4966](https://github.com/terraform-providers/terraform-provider-google/pull/4966))
* dataproc: `google_dataproc_cluster` `min_cpu_platform` on both `worker_config` and `master_config` is now GA ([#4968](https://github.com/terraform-providers/terraform-provider-google/pull/4968))
* kms: enabled use of `user_project_override` for the `kms_crypto_key` resource ([#4967](https://github.com/terraform-providers/terraform-provider-google/pull/4967))
* kms: enabled use of `user_project_override` for the `kms_secret_ciphertext` data source ([#4985](https://github.com/terraform-providers/terraform-provider-google/pull/4985))
* sql: added `root_password` field to `google_sql_database_instance` resource ([#4983](https://github.com/terraform-providers/terraform-provider-google/pull/4983))

BUG FIXES:
* bigquery: fixed an issue where bigquery table id formats from the `2.X` series caused an error at plan time ([#5012](https://github.com/terraform-providers/terraform-provider-google/pull/5012))
* cloudbuild: Fixed incorrect dependency between `trigger_template` and `github` in `google_cloud_build_trigger`. ([#4946](https://github.com/terraform-providers/terraform-provider-google/pull/4946))
* cloudfunctions: Fixed inability to set `google_cloud_functions_function` update timeout. ([#5011](https://github.com/terraform-providers/terraform-provider-google/pull/5011))
* cloudrun: Wait for the cloudrun resource to reach a ready state before returning success. ([#4945](https://github.com/terraform-providers/terraform-provider-google/pull/4945))
* compute: `self_link` in several datasources will now error on invalid values instead of crashing ([#4887](https://github.com/terraform-providers/terraform-provider-google/pull/4887))
* compute: field `advertised_ip_ranges` in `google_compute_router_peer` can now be updated without recreating the resource. ([#4923](https://github.com/terraform-providers/terraform-provider-google/pull/4923))
* compute: marked `min_cpu_platform` on `google_compute_instance` as computed so if it is not specified it will not cause diffs ([#4980](https://github.com/terraform-providers/terraform-provider-google/pull/4980))
* dns: Fixed issue causing `google_dns_record_set` deletion to fail when the managed zone ceased to exist before the deletion event. ([#5010](https://github.com/terraform-providers/terraform-provider-google/pull/5010))
* iam: disallowed `deleted:` principals in IAM resources ([#4958](https://github.com/terraform-providers/terraform-provider-google/pull/4958))
* sql: added retries to `google_sql_user` create and update to reduce flakiness ([#4860](https://github.com/terraform-providers/terraform-provider-google/pull/4860))

## 3.0.0 (December 04, 2019)

NOTES:

These are the changes between 3.0.0-beta.1 and the 3.0.0 final release. For changes since 2.20.0, see also the 3.0.0-beta.1 changelog entry below.

**Please see [the 3.0.0 upgrade guide](https://www.terraform.io/docs/providers/google/guides/version_3_upgrade.html) for upgrade guidance.**

BREAKING CHANGES:
* cloudrun: updated `cloud_run_service` to v1. Significant updates have been made to the resource including a breaking schema change. ([#4972](https://github.com/terraform-providers/terraform-provider-google/issues/4972))

BUG FIXES:
* compute: fixed a bug in `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` that created an artificial diff when removing a now-removed field from a config ([#4929](https://github.com/terraform-providers/terraform-provider-google/issues/4929))
* dns: Fixed bug causing `google_dns_managed_zone` datasource to always return a 404 ([#4940](https://github.com/terraform-providers/terraform-provider-google/issues/4940))
* service_networking: fixed "An unknown error occurred" bug when creating multiple google_service_networking_connection resources in parallel ([#4646](https://github.com/terraform-providers/terraform-provider-google/issues/4646))

## 3.0.0-beta.1 (November 15, 2019)

BREAKING CHANGES:

* access_context_manager: Made `os_type` required on block `google_access_context_manager_access_level.basic.conditions.device_policy.os_constraints`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* all: changed any id values that could not be interpolated as self_links into values that could [MM#2461](https://github.com/GoogleCloudPlatform/magic-modules/pull/2461)
* app_engine: Made `ssl_management_type` required on `google_app_engine_domain_mapping.ssl_settings` [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* app_engine: Made `shell` required on `google_app_engine_standard_app_version.entrypoint`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* app_engine: Made `source_url` required on `google_app_engine_standard_app_version.deployment.files` and `google_app_engine_standard_app_version.deployment.zip`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* app_engine: Made `split_health_checks ` required on `google_app_engine_application.feature_settings` [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* app_engine: Made `script_path` required on `google_app_engine_standard_app_version.handlers.script`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* bigtable: Made `cluster_id` required on `google_bigtable_app_profile.single_cluster_routing`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* bigquery: Made at least one of `range` or `skip_leading_rows` required on `google_bigquery_table.external_data_configuration.google_sheets_options`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* bigquery: Made `role` required on `google_bigquery_dataset.access`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* bigtable: Made exactly one of `single_cluster_routing` or `multi_cluster_routing_use_any` required on `google_bigtable_app_profile`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* binary_authorization: Made `name_pattern` required on `google_binary_authorization_policy.admission_whitelist_patterns`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* binary_authorization: Made `evaluation_mode` and `enforcement_mode` required on `google_binary_authorization_policy.cluster_admission_rules`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* cloudbuild: made Cloud Build Trigger's trigger template required to match API requirements. [MM#2352](https://github.com/GoogleCloudPlatform/magic-modules/pull/2352)
* cloudbuild: Made `branch` required on `google_cloudbuild_trigger.github`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* cloudbuild: Made `steps` required on `google_cloudbuild_trigger.build`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* cloudbuild: Made `name` required on `google_cloudbuild_trigger.build.steps`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* cloudbuild: Made `name` and `path` required on `google_cloudbuild_trigger.build.steps.volumes`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* cloudbuild: Made exactly one of `filename` or `build` required on `google_cloudbuild_trigger`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* cloudfunctions: deprecated `nodejs6` as option for `runtime` in `function` and made it required. [MM#2499](https://github.com/GoogleCloudPlatform/magic-modules/pull/2499)
* cloudscheduler: Made exactly one of `pubsub_target`, `http_target` or `app_engine_http_target` required on `google_cloudscheduler_job`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* cloudiot: removed `event_notification_config` (singular) from `google_cloudiot_registry`. Use plural `event_notification_configs` instead. [MM#2390](https://github.com/GoogleCloudPlatform/magic-modules/pull/2390)
* cloudiot: Made `public_key_certificate` required on `google_cloudiot_registry. credentials `. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* cloudscheduler: Made `service_account_email` required on `google_cloudscheduler_job.http_target.oauth_token` and `google_cloudscheduler_job.http_target.oidc_token`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* composer: Made at least one of `airflow_config_overrides`, `pypi_packages`, `env_variables, `image_version`, or `python_version` required on `google_composer_environment.config.software_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* composer: Made `use_ip_aliases` required on `google_composer_environment.config.node_config.ip_allocation_policy`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* composer: Made `enable_private_endpoint` required on `google_composer_environment.config.private_environment_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* composer: Made at least one of `enable_private_endpoint` or `master_ipv4_cidr_block` required on `google_composer_environment.config.private_environment_config` [MM#2682](https://github.com/GoogleCloudPlatform/magic-modules/pull/2682)
* composer: Made at least one of `node_count`, `node_config`, `software_config` or `private_environment_config` required on `google_composer_environment.config` [MM#2682](https://github.com/GoogleCloudPlatform/magic-modules/pull/2682)
* compute: `google_compute_backend_service`'s `backend` field field now requires the `group` subfield to be set. [MM#2373](https://github.com/GoogleCloudPlatform/magic-modules/pull/2373)
* compute: permanently removed `ip_version` field from `google_compute_forwarding_rule` [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* compute: permanently removed `ipv4_range` field from `google_compute_network`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* compute: permanently removed `auto_create_routes` field from `google_compute_network_peering`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* compute: permanently removed `update_strategy` field from `google_compute_region_instance_group_manager`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* compute: added check to only allow `google_compute_instance_template`s with 375gb scratch disks [MM#2495](https://github.com/GoogleCloudPlatform/magic-modules/pull/2495)
* compute: made `google_compute_instance_template` fail at plan time when scratch disks do not have `disk_type` `"local-ssd"`. [MM#2282](https://github.com/GoogleCloudPlatform/magic-modules/pull/2282)
* compute: removed `enable_flow_logs` field from `google_compute_subnetwork`. This is now controlled by the presence of the `log_config` block [MM#2597](https://github.com/GoogleCloudPlatform/magic-modules/pull/2597)
* compute: Made `raw_key` required on `google_compute_snapshot.snapshot_encryption_key`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made at least one of `auto_delete`, `device_name`, `disk_encryption_key_raw`, `kms_key_self_link`, `initialize_params`, `mode` or `source` required on `google_compute_instance.boot_disk`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made at least one of `size`, `type`, `image`, or `labels` required on `google_compute_instance.boot_disk.initialize_params`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made at least one of `enable_secure_boot`, `enable_vtpm`, or `enable_integrity_monitoring` required on `google_compute_instance.shielded_instance_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made at least one of `on_host_maintenance`, `automatic_restart`, `preemptible`, or `node_affinities` required on `google_compute_instance.scheduling`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made `interface` required on `google_compute_instance.scratch_disk`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made at least one of `enable_secure_boot`, `enable_vtpm`, or `enable_integrity_monitoring` required on `google_compute_instance_template.shielded_instance_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made at least one of `on_host_maintenance`, `automatic_restart`, `preemptible`, or `node_affinities` are now required on `google_compute_instance_template.scheduling`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made `kms_key_self_link` required on `google_compute_instance_template.disk.disk_encryption_key`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made `range` required on `google_compute_router_peer. advertised_ip_ranges`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Made `channel` required on `google_container_cluster.release_channel`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* compute: Removed `instance_template` for `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager`. Use `version.instance_template` instead. [MM#2595](https://github.com/GoogleCloudPlatform/magic-modules/pull/2595)
* compute: removed `update_strategy` for `google_compute_instance_group_manager`. Use `update_policy` instead. [MM#2595](https://github.com/GoogleCloudPlatform/magic-modules/pull/2595)
* compute: stopped allowing selfLink or path style references as IP addresses for `google_compute_forwarding_rule` or `google_compute_global_forwarding_rule` [MM#2620](https://github.com/GoogleCloudPlatform/magic-modules/pull/2620)
* compute: Made exactly one of `http_health_check`, `https_health_check`, `http2_health_check`, `tcp_health_check` or `ssl_health_check` required on `google_compute_health_check`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* compute: Made exactly one of `http_health_check`, `https_health_check`, `http2_health_check`, `tcp_health_check` or `ssl_health_check` required on `google_compute_region_health_check`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* container: permanently removed `zone` and `region` fields from data source `google_container_engine_versions`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* container: permanently removed `zone`, `region` and `additional_zones` fields from `google_container_cluster`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* container: permanently removed `zone` and `region` fields from `google_container_node_pool`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* container: set `google_container_cluster`'s `logging_service` and `monitoring_service` defaults to enable GKE Stackdriver Monitoring. [MM#2471](https://github.com/GoogleCloudPlatform/magic-modules/pull/2471)
* container: removed `kubernetes_dashboard` from `google_container_cluster.addons_config` [MM#2551](https://github.com/GoogleCloudPlatform/magic-modules/pull/2551)
* container: removed automatic suppression of GPU taints in GKE `taint` [MM#2537](https://github.com/GoogleCloudPlatform/magic-modules/pull/2537)
* container: Made `disabled` required on `google_container_cluster.addons_config.http_load_balancing`, `google_container_cluster.addons_config.horizontal_pod_autoscaling`, `google_container_cluster.addons_config.network_policy_config`, `google_container_cluster.addons_config.cloudrun_config`, and `google_container_cluster.addons_config.istio_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: Made at least one of `http_load_balancing`, `horizontal_pod_autoscaling` , `network_policy_config`, `cloudrun_config`, or `istio_config` required on `google_container_cluster.addons_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: Made `enabled` required on `google_container_cluster.network_policy`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: Made `enable_private_endpoint` required on `google_container_cluster.private_cluster_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: Made `enabled` required on `google_container_cluster.vertical_pod_autoscaling`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: Made `cidr_blocks` required on `google_container_cluster.master_authorized_networks_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: Made at least one of `username`, `password` or `client_certificate_config` required on `google_container_cluster.master_auth`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* container: removed `google_container_cluster` `ip_allocation_policy.use_ip_aliases`. If it's set to true, remove it from your config. If false, remove `ip_allocation_policy` as a whole. [MM#2615](https://github.com/GoogleCloudPlatform/magic-modules/pull/2615)
* container: removed `google_container_cluster` `ip_allocation_policy.create_subnetwork`, `ip_allocation_policy.subnetwork_name`, `ip_allocation_policy.node_ipv4_cidr_block`. Define an explicit `google_compute_subnetwork` and use `subnetwork` instead. [MM#2615](https://github.com/GoogleCloudPlatform/magic-modules/pull/2615)
* dataproc: Made at least one of `staging_bucket`, `gce_cluster_config`, `master_config`, `worker_config`, `preemptible_worker_config`, `software_config`, `initialization_action` or `encryption_config` required on `google_dataproc_cluster.cluster_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `zone`, `network`, `subnetwork`, `tags`, `service_account`, `service_account_scopes`, `internal_ip_only` or `metadata` required on `google_dataproc_cluster.cluster_config.gce_cluster_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `num_instances`, `image_uri`, `machine_type`, `min_cpu_platform`, `disk_config`, or `accelerators` required on `google_dataproc_cluster.cluster_config.master_config` and `google_dataproc_cluster.cluster_config.worker_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `num_local_ssds`, `boot_disk_size_gb` or `boot_disk_type` required on `google_dataproc_cluster.cluster_config.preemptible_worker_config.disk_config`, `google_dataproc_cluster.cluster_config.master_config.disk_config` and `google_dataproc_cluster.cluster_config.worker_config.disk_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `num_instances` or `disk_config` required on `google_dataproc_cluster.cluster_config.preemptible_worker_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `image_version`, `override_properties` or `optional_components` is now required on `google_dataproc_cluster.cluster_config.software_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made `policy_uri` required on `google_dataproc_cluster.cluster_config.autoscaling_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made `max_failures_per_hour` required on `google_dataproc_job.scheduling`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made `driver_log_levels` required on `google_dataproc_job.pyspark_config.logging_config`, `google_dataproc_job.spark_config.logging_config`, `google_dataproc_job.hadoop_config.logging_config`, `google_dataproc_job.hive_config.logging_config`, `google_dataproc_job.pig_config.logging_config`, `google_dataproc_job.sparksql_config.logging_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `main_class` or `main_jar_file_uri` required on `google_dataproc_job.spark_config` and `google_dataproc_job.hadoop_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dataproc: Made at least one of `query_file_uri` or `query_list` required on `google_dataproc_job.hive_config`, `google_dataproc_job.pig_config`, and `google_dataproc_job.sparksql_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dns: Made `networks` required on `google_dns_managed_zone.private_visibility_config`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* dns: Made `network_url` required on `google_dns_managed_zone.private_visibility_config.networks`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* iam: made `iam_audit_config` resources overwrite existing audit config on create. Previous implementations merged config with existing audit configs on create. [MM#2438](https://github.com/GoogleCloudPlatform/magic-modules/pull/2438)
* iam: Made exactly one of `list_policy`, `boolean_policy`, or `restore_policy` required on `google_organization_policy`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* iam: Made exactly one of `all` or `values` required on `google_organization_policy.list_policy.allow` and `google_organization_policy.list_policy.deny`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* iam: `google_project_iam_policy` can handle the `project` field in either of the following forms: `project-id` or `projects/project-id` [MM#2700](https://github.com/GoogleCloudPlatform/magic-modules/pull/2700)
* iam: Made exactly one of `allow` or `deny` required on `google_organization_policy.list_policy` [MM#2682](https://github.com/GoogleCloudPlatform/magic-modules/pull/2682)
* iam: removed the deprecated `pgp_key`, `private_key_encrypted` and `private_key_fingerprint` from `google_service_account_key` [MM#2680](https://github.com/GoogleCloudPlatform/magic-modules/pull/2680)
* monitoring: permanently removed `is_internal` and `internal_checkers` fields from `google_monitoring_uptime_check_config`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* monitoring: permanently removed `labels` field from `google_monitoring_alert_policy`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* monitoring: Made `content` required on `google_monitoring_uptime_check_config.content_matchers`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* monitoring: Made exactly one of `http_check` or `tcp_check` is now required on `google_monitoring_uptime_check_config`. [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* monitoring: Made at least one of `auth_info`, `port`, `headers`, `path`, `use_ssl`, or `mask_headers` is now required on `google_monitoring_uptime_check_config.http_check` [MM#2665](https://github.com/GoogleCloudPlatform/magic-modules/pull/2665)
* provider: added the `https://www.googleapis.com/auth/userinfo.email` scope to the provider by default [MM#2473](https://github.com/GoogleCloudPlatform/magic-modules/pull/2473)
* pubsub: removed ability to set a full path for `google_pubsub_subscription.name` (e.g. `projects/my-project/subscriptions/my-subscription`). `name` now must be the shortname (e.g. `my-subscription`) [MM#2561](https://github.com/GoogleCloudPlatform/magic-modules/pull/2561)
* resourcemanager: converted `google_folder_organization_policy` and `google_organization_policy` import format to use slashes instead of colons. [MM#2638](https://github.com/GoogleCloudPlatform/magic-modules/pull/2638)
* serviceusage: removed `google_project_services` [MM#2403](https://github.com/GoogleCloudPlatform/magic-modules/pull/2403)
* serviceusage: stopped accepting `bigquery-json.googleapis.com` in `google_project_service`. Specify `biquery.googleapis.com` instead. [MM#2626](https://github.com/GoogleCloudPlatform/magic-modules/pull/2626)
* sql: Made `name` and `value` required on `google_sql_database_instance.settings.database_flags`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* sql: Made at least one of `binary_log_enabled`, `enabled`, `start_time`, and `location` required on `google_sql_database_instance.settings.backup_configuration`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* sql: Made at least one of `authorized_networks`, `ipv4_enabled`, `require_ssl`, and `private_network` required on `google_sql_database_instance.settings.ip_configuration`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* sql: Made at least one of `day`, `hour`, and `update_track` required on `google_sql_database_instance.settings.maintenance_window`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* sql: Made at least one of `cert`, `common_name`, `create_time`, `expiration_time`, or `sha1_fingerprint` required on `google_sql_database_instance.settings.server_ca_cert`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* sql: Made at least one of `ca_certificate`, `client_certificate`, `client_key`, `connect_retry_interval`, `dump_file_path`, `failover_target`, `master_heartbeat_period`, `password`, `ssl_cipher`, `username`, and `verify_server_certificate` required on `google_sql_database_instance.settings.replica_configuration`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* sql: Made `value` required on `google_sql_database_instance.settings.ip_configuration.authorized_networks`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* storage: permanently removed `is_live` flag from `google_storage_bucket`. [MM#2436](https://github.com/GoogleCloudPlatform/magic-modules/pull/2436)
* storage: Made at least one of `main_page_suffix` or `not_found_page` required on `google_storage_bucket.website`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* storage: Made at least one of `min_time_elapsed_since_last_modification`, `max_time_elapsed_since_last_modification`, `include_prefixes`, or `exclude_prefixes` required on `google_storage_transfer_job.transfer_spec.object_conditions`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* storage: Made at least one of `overwrite_objects_already_existing_in_sink`, `delete_objects_unique_in_sink`, and `delete_objects_from_source_after_transfer` required on `google_storage_transfer_job.transfer_spec.transfer_options`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
* storage: Made at least one of `gcs_data_source`, `aws_s3_data_source`, or `http_data_source` required on `google_storage_transfer_job.transfer_options`. [MM#2608](https://github.com/GoogleCloudPlatform/magic-modules/pull/2608)
