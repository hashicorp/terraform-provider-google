## 3.18.0 (Unreleased)
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

## 2.20.3 (March 10, 2020)

NOTES:
* `2.20.3` is a backport release, and some changes will not appear in `3.X` series releases until `3.12.0`.
To upgrade to `3.X` you will need to perform a large jump in versions, and it is _strongly_ advised that you attempt to upgrade to `3.X` instead of using this release.
* `2.20.3` is primarily a preventative fix, in anticipation of a change in API response messages adding a default value.

BUG FIXES:
* compute: fixed error when reading `google_compute_instance_template` resources with `network_interface[*].name` set. ([#5812](https://github.com/terraform-providers/terraform-provider-google/pull/5812))

## 2.20.2 (February 03, 2020)

BUG FIXES:
* bigtable: fixed diff for DEVELOPMENT instances that are returned from the API with one node ([#5557](https://github.com/terraform-providers/terraform-provider-google/pull/5557))

## 2.20.1 (December 13, 2019)

**Note:** 2.20.1 is a backport release. The changes in it are unavailable in 3.0.0-beta.1 through 3.2.0.

BUG FIXES:
* iam: Fixed a bug that causes badRequest errors on IAM resources due to deleted serviceAccount principals ([#5142](https://github.com/terraform-providers/terraform-provider-google/pull/5142))

## 2.20.0 (November 13, 2019)

BREAKING CHANGES:
* compute: the `backend.group` field is now required for `google_compute_region_backend_service`. Configurations without this would not have worked, so this isn't considered an API break. ([#4772](https://github.com/terraform-providers/terraform-provider-google/pull/4772))

IMPROVEMENTS:
* bigtable: added import support to `google_bigtable_table` ([#4849](https://github.com/terraform-providers/terraform-provider-google/pull/4849))
* compute: `load_balancing_scheme` for `google_compute_forwarding_rule` now accepts `INTERNAL_MANAGED` as a value. ([#4772](https://github.com/terraform-providers/terraform-provider-google/pull/4772))
* compute: extended backend configuration options for `google_compute_region_backend_service` to include `backend.balancing_mode`, `backend.capacity_scaler`, `backend.max_connections`, `backend.max_connections_per_endpoint`, `backend.max_connections_per_instance`, `backend.max_rate`, `backend.max_rate_per_endpoint`, `backend.max_rate_per_instance`, and `backend.max_utilization` ([#4772](https://github.com/terraform-providers/terraform-provider-google/pull/4772))
* iam: changed the `id` for many IAM resources to the reference resource long name. Updated `instance_name` on `google_compute_instance_iam` and `subnetwork` on `google_compute_subnetwork` to their respective long names in state ([#4866](https://github.com/terraform-providers/terraform-provider-google/pull/4866))
* logging: added `display_name` field to `google_logging_metric` resource ([#4839](https://github.com/terraform-providers/terraform-provider-google/pull/4839))
* monitoring: Added `validate_ssl` to `google_monitoring_uptime_check_config` ([#4637](https://github.com/terraform-providers/terraform-provider-google/pull/4637))
* project: added batching functionality to `google_project_service` read calls, so fewer API requests are made ([#4854](https://github.com/terraform-providers/terraform-provider-google/pull/4854))
* storage: added `notification_id` field to `google_storage_notification` ([#4879](https://github.com/terraform-providers/terraform-provider-google/pull/4879))

BUG FIXES:
* compute: fixed issue where setting a 0 for `min_replicas` in `google_compute_autoscaler` and `google_compute_region_autoscaler` would set that field to its server-side default instead of 0. ([#4851](https://github.com/terraform-providers/terraform-provider-google/pull/4851))
* dns: fixed crash when `network` blocks are defined without `network_url`s ([#4840](https://github.com/terraform-providers/terraform-provider-google/pull/4840))
* google: used the correct update method for google_service_account.description ([#4870](https://github.com/terraform-providers/terraform-provider-google/pull/4870))
* logging: fixed issue where logging exclusion resources silently failed when being mutated in parallel ([#4814](https://github.com/terraform-providers/terraform-provider-google/pull/4814))

## 2.19.0 (November 05, 2019)

DEPRECATIONS:
* `compute`: deprecated `enable_flow_logs` on `google_compute_subnetwork`. The presence of the `log_config` block signals that flow logs are enabled for a subnetwork ([#4791](https://github.com/terraform-providers/terraform-provider-google/pull/4791))
* `compute`: deprecated `instance_template` for `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` . Use `version.instance_template` instead. ([#4763](https://github.com/terraform-providers/terraform-provider-google/pull/4763))
* `compute`: deprecated `update_strategy` for `google_compute_instance_group_manager` . Use `update_policy` instead. ([#4763](https://github.com/terraform-providers/terraform-provider-google/pull/4763))
* `container`: deprecated `google_container_cluster` `ip_allocation_policy.create_subnetwork`, `ip_allocation_policy.subnetwork_name`, `ip_allocation_policy.node_ipv4_cidr_block`. Define an explicit `google_compute_subnetwork` and use `subnetwork` instead. ([#4774](https://github.com/terraform-providers/terraform-provider-google/pull/4774))
* `container`: deprecated `google_container_cluster` `ip_allocation_policy.use_ip_aliases`. If it's set to true, remove it from your config. If false, remove `ip_allocation_policy` as a whole. ([#4774](https://github.com/terraform-providers/terraform-provider-google/pull/4774))
* `iam`: Deprecated `pgp_key` on `google_service_account_key` resource. See https://www.terraform.io/docs/extend/best-practices/sensitive-state.html for more information. ([#4810](https://github.com/terraform-providers/terraform-provider-google/pull/4810))

BREAKING CHANGES:
* `google_service_account_iam_*` resources now support IAM Conditions. If any conditions had been created out of band before this release, take extra care to ensure they are present in your Terraform config so the provider doesn't try to create new bindings with no conditions. Terraform will show a diff that it is adding the condition to the resource, which is safe to apply. ([#4541](https://github.com/terraform-providers/terraform-provider-google/pull/4541))

FEATURES:
* `compute`: added `google_compute_router` datasource ([#4614](https://github.com/terraform-providers/terraform-provider-google/pull/4614))

IMPROVEMENTS:
* `cloudbuild`: added ability to specify `name` for `cloud_build_trigger` to avoid name collisions when creating multiple triggers at once. ([#4709](https://github.com/terraform-providers/terraform-provider-google/pull/4709))
* `compute`: `log_config` is now available in GA for `google_compute_subnetwork` ([#4791](https://github.com/terraform-providers/terraform-provider-google/pull/4791))
* `compute`: added support for multiple versions of `instance_template` and granular control of the update policies for `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager`. ([#4763](https://github.com/terraform-providers/terraform-provider-google/pull/4763))
* `container`: added `maintenance_policy.recurring_window` support to `google_container_cluster`, significantly increasing expressive range. ([#4736](https://github.com/terraform-providers/terraform-provider-google/pull/4736))
* `container`: added `taint` field in GKE resources to the GA `google` provider ([#4743](https://github.com/terraform-providers/terraform-provider-google/pull/4743))
* `container`: fix a diff created in the cloud console when `MaintenanceExclusions` are added. ([#4764](https://github.com/terraform-providers/terraform-provider-google/pull/4764))
* `compute`: added `google_compute_instance` support for display device (Virtual Displays) ([#4775](https://github.com/terraform-providers/terraform-provider-google/pull/4775))
* `iam`: added support for IAM Conditions to the `google_service_account_iam_*` resources (beta provider only) ([#4541](https://github.com/terraform-providers/terraform-provider-google/pull/4541))
* `iam`: added `description` to `google_service_account`. ([#4734](https://github.com/terraform-providers/terraform-provider-google/pull/4734))

BUG FIXES:
* `appengine`: Resolved permadiff in `google_app_engine_domain_mapping.ssl_settings.certificate_id`. ([#4754](https://github.com/terraform-providers/terraform-provider-google/pull/4754))
* `storage`: Fixed error in `google_storage_bucket` where locked retention policies would cause a bucket to report failure on all updates (even though updates were applied correctly). ([#4761](https://github.com/terraform-providers/terraform-provider-google/pull/4761))

## 2.18.1 (October 25, 2019)

BUGS:
* `resourcemanager`: fixed deleting the default network in `google_project` ([#4748](https://github.com/terraform-providers/terraform-provider-google/pull/4748))

## 2.18.0 (October 23, 2019)

KNOWN ISSUES:
* `resourcemanager`: `google_project` `auto_create_network` is failing to delete networks when set to `false`. Use an earlier provider version to resolve.

DEPRECATIONS:
* `container`: The `kubernetes_dashboard` addon is deprecated for `google_container_cluster`. ([#4648](https://github.com/terraform-providers/terraform-provider-google/pull/4648))

FEATURES:
* **New Resource:** `google_app_engine_application_url_dispatch_rules` ([#4674](https://github.com/terraform-providers/terraform-provider-google/pull/4674))

IMPROVEMENTS:
* `all`: increased support for custom endpoints across the provider ([#4641](https://github.com/terraform-providers/terraform-provider-google/pull/4641))
* `appengine`: added the ability to delete the parent service of `google_app_engine_standard_app_version` ([#4596](https://github.com/terraform-providers/terraform-provider-google/pull/4596))
* `container`: Added `shielded_instance_config` attribute to `node_config` ([#4554](https://github.com/terraform-providers/terraform-provider-google/pull/4554))
* `dataflow`: added `ip_configuration` option to `job`. ([#4726](https://github.com/terraform-providers/terraform-provider-google/pull/4726))
* `pubsub`: Added field `oidc_token` to `google_pubsub_subscription` ([#4679](https://github.com/terraform-providers/terraform-provider-google/pull/4679))
* `sql`: added `location` field to `backup_configuration` block in `google_sql_database_instance` ([#4681](https://github.com/terraform-providers/terraform-provider-google/pull/4681))

BUGS:
* `all`: fixed the custom endpoint version used by older legacy REST clients ([#4695](https://github.com/terraform-providers/terraform-provider-google/pull/4695))
* `bigquery`: fix issue with `google_bigquery_data_transfer_config` `params` crashing on boolean values ([#4676](https://github.com/terraform-providers/terraform-provider-google/pull/4676))
* `cloudrun`: fixed the apiVersion sent in `google_cloud_run_domain_mapping` requests ([#4657](https://github.com/terraform-providers/terraform-provider-google/pull/4657))
* `compute`: added support for updating multiple fields at once to `google_compute_subnetwork` ([#4688](https://github.com/terraform-providers/terraform-provider-google/pull/4688))
* `compute`: fixed diffs in `google_compute_instance_group`'s `network` field when equivalent values were specified ([#4728](https://github.com/terraform-providers/terraform-provider-google/pull/4728))
* `compute`: fixed issues updating `google_compute_instance_group`'s `instances` field when config/state values didn't match ([#4728](https://github.com/terraform-providers/terraform-provider-google/pull/4728))
* `iam`: fixed bug where IAM binding wouldn't replace members if they were deleted outside of terraform. ([#4693](https://github.com/terraform-providers/terraform-provider-google/pull/4693))
* `pubsub`: Fixed permadiff due to interaction of organization policies and `google_pubsub_topic`. ([#4721](https://github.com/terraform-providers/terraform-provider-google/pull/4721))

## 2.17.0 (October 08, 2019)

NOTES:
* An [upgrade guide](https://www.terraform.io/docs/providers/google/version_3_upgrade.html) has been started for the upcoming 3.0.0 release. ([#4594](https://github.com/terraform-providers/terraform-provider-google/pull/4594))
* `google_project_services` users of provider versions prior to `2.17.0` should update, as past versions of the provider will not handle an upcoming rename of `bigquery-json.googleapis.com` to `bigquery.googleapis.com` well. See https://github.com/terraform-providers/terraform-provider-google/issues/4590 for details. ([#4616](https://github.com/terraform-providers/terraform-provider-google/pull/4616))

DEPRECATIONS:
* `google_project_services` ([#4587](https://github.com/terraform-providers/terraform-provider-google/pull/4587))

FEATURES:
* **New Resource:** `google_bigtable_gc_policy` ([#4578](https://github.com/terraform-providers/terraform-provider-google/pull/4578))
* **New Resource:** `google_binary_authorization_attestor_iam_policy` ([#4517](https://github.com/terraform-providers/terraform-provider-google/pull/4517))
* **New Resource:** `google_compute_region_ssl_certificate` ([#4537](https://github.com/terraform-providers/terraform-provider-google/pull/4537))
* **New Resource:** `google_compute_region_target_http_proxy` ([#4537](https://github.com/terraform-providers/terraform-provider-google/pull/4537))
* **New Resource:** `google_compute_region_target_https_proxy` ([#4537](https://github.com/terraform-providers/terraform-provider-google/pull/4537))
* **New Resource:** `google_iap_app_engine_service_iam_*` ([#4566](https://github.com/terraform-providers/terraform-provider-google/pull/4566))
* **New Resource:** `google_iap_app_engine_version_iam_*` ([#4566](https://github.com/terraform-providers/terraform-provider-google/pull/4566))
* **New Resource:** `google_storage_bucket_access_control` ([#4531](https://github.com/terraform-providers/terraform-provider-google/pull/4531))

IMPROVEMENTS:
* all: made `monitoring-read` scope available. ([#4569](https://github.com/terraform-providers/terraform-provider-google/pull/4569))
* bigquery: Added support for default customer-managed encryption keys (CMEK) for BigQuery datasets. ([#4312](https://github.com/terraform-providers/terraform-provider-google/pull/4312))
* bigtable: import support added to `google_bigtable_instance` ([#4598](https://github.com/terraform-providers/terraform-provider-google/pull/4598))
* container: moved `default_max_pods_per_node` to ga. ([#4621](https://github.com/terraform-providers/terraform-provider-google/pull/4621))
* containeranalysis: moved `google_containeranalysis_note` to ga ([#4517](https://github.com/terraform-providers/terraform-provider-google/pull/4517))
* projectservice: added mitigations for bigquery-json to bigquery rename in project service resources. ([#4616](https://github.com/terraform-providers/terraform-provider-google/pull/4616))

BUGS:
* cloudscheduler: Fixed permadiff for `app_engine_http_target.app_engine_routing` on `google_cloud_scheduler_job` ([#4444](https://github.com/terraform-providers/terraform-provider-google/pull/4444))
* compute: Added ability to set `quic_override` on `google_compute_https_target_proxy` to empty. ([#4588](https://github.com/terraform-providers/terraform-provider-google/pull/4588))
* compute: Fix bug where changes to `region_backend_service.backends.failover` was not detected. ([#4622](https://github.com/terraform-providers/terraform-provider-google/pull/4622))
* compute: fixed `google_compute_router_peer` to default if empty for `advertise_mode` ([#4503](https://github.com/terraform-providers/terraform-provider-google/pull/4503))
* compute: fixed perma-diff in `google_compute_router_nat` when referencing subnetwork via `name` ([#4549](https://github.com/terraform-providers/terraform-provider-google/pull/4549))
* container: fixed an overly-aggressive validation for `master_ipv4_cidr_block` in `google_container_cluster` ([#4577](https://github.com/terraform-providers/terraform-provider-google/pull/4577))

## 2.16.0 (September 24, 2019)

KNOWN ISSUES:
* Based on an upstream change, users of the `google_project_services` resource may have seen the `bigquery.googleapis.com` service added and the `bigquery-json.googleapis.com` service removed, causing a diff. This was later reverted, causing another diff. This issue is being tracked as https://github.com/terraform-providers/terraform-provider-google/issues/4590.

FEATURES:
* **New Resource**: `google_compute_region_url_map` is now available. To support this, the `protocol` for `google_compute_region_backend_service` can now be set to `HTTP`, `HTTPS`, `HTTP2`, and `SSL`. ([#4496](https://github.com/terraform-providers/terraform-provider-google/issues/4496))
* **New Resource**: Adds `google_runtimeconfig_config_iam_*` resources ([#4454](https://github.com/terraform-providers/terraform-provider-google/issues/4454))
* **New Resource**: Added `google_compute_resource_policy` and `google_compute_disk_resource_policy_attachment` to manage `google_compute_disk` resource policies as fine-grained resources ([#4409](https://github.com/terraform-providers/terraform-provider-google/issues/4409))

ENHANCEMENTS:
* composer: Add `python_version` and ability to set `image_version` in `google_composer_environment` in the GA provider ([#4465](https://github.com/terraform-providers/terraform-provider-google/issues/4465))
* compute: `google_compute_global_forwarding_rule` now supports `metadata_filters`. ([#4495](https://github.com/terraform-providers/terraform-provider-google/issues/4495))
* compute: `google_compute_backend_service` now supports `locality_lb_policy`, `outlier_detection`, `consistent_hash`, and `circuit_breakers`. ([#4412](https://github.com/terraform-providers/terraform-provider-google/issues/4412))
* compute: Add support for `guest_os_features` to resource `google_compute_image` ([#4483](https://github.com/terraform-providers/terraform-provider-google/issues/4483))
* compute: `google_compute_router_nat` now supports `drain_nat_ips` field ([#4480](https://github.com/terraform-providers/terraform-provider-google/issues/4480))
* container: `google_container_node_pool` now supports node_locations to specify specific node zones. ([#4478](https://github.com/terraform-providers/terraform-provider-google/issues/4478))
* googleapis: `google_netblock_ip_ranges` data source now has a `private-googleapis` field, for the IP addresses used for Private Google Access for services that do not support VPC Service Controls API access. ([#4367](https://github.com/terraform-providers/terraform-provider-google/issues/4367))
* project: `google_project_iam_*` Properly set the `project` field in state ([#4488](https://github.com/terraform-providers/terraform-provider-google/issues/4488))

BUG FIXES:
* cloudiot: Fixed error where `subfolder_matches` were not set in `google_cloudiot_registry` `event_notification_configs` ([#4527](https://github.com/terraform-providers/terraform-provider-google/issues/4527))

## 2.15.0 (September 17, 2019)

FEATURES:
* **New Resource**: `google_iap_web_iam_binding/_member/_policy` are now available for managing IAP web IAM permissions ([#4253](https://github.com/terraform-providers/terraform-provider-google/issues/4253))
* **New Resource**: `google_iap_web_backend_service_binding/_member/_policy` are now available for managing IAM permissions on IAP enabled backend services ([#4253](https://github.com/terraform-providers/terraform-provider-google/issues/4253))
* **New Resource**: `google_iap_web_type_compute_iam_binding/_member/_policy` are now available for managing IAM permissions on IAP enabled compute services ([#4253](https://github.com/terraform-providers/terraform-provider-google/issues/4253))
* **New Resource**: `google_iap_web_type_app_engine_iam_binding/_member/_policy` are now available for managing IAM permissions on IAP enabled App Engine applications ([#4253](https://github.com/terraform-providers/terraform-provider-google/issues/4253))
* **New Resource**: Add the new resource `google_app_engine_domain_mapping` ([#4310](https://github.com/terraform-providers/terraform-provider-google/issues/4310))
* **New Resource**: `google_cloudfunctions_function_iam_policy`, `google_cloudfunctions_function_iam_binding`, and `google_cloudfunctions_function_iam_member` have been added ([#4420](https://github.com/terraform-providers/terraform-provider-google/issues/4420))
* **New Resource**: `google_compute_reservation` allows you to reserve instance capacity in GCE. ([#4332](https://github.com/terraform-providers/terraform-provider-google/issues/4332))
* **New Resource**: `google_compute_region_health_check` is now available. This and `google_compute_health_check` now include additional support for HTTP2 health checks. ([#4270](https://github.com/terraform-providers/terraform-provider-google/issues/4270))

ENHANCEMENTS:
* compute: Add all options to `google_compute_router_peer` ([#4371](https://github.com/terraform-providers/terraform-provider-google/issues/4371))
* compute: add `tunnel_id` to `google_compute_vpn_tunnel` and `gateway_id` to `google_compute_vpn_gateway` ([#4373](https://github.com/terraform-providers/terraform-provider-google/issues/4373))
* compute: `google_compute_subnetwork` now includes the `purpose` and `role` fields. ([#4261](https://github.com/terraform-providers/terraform-provider-google/issues/4261))
* compute: add `purpose` field to `google_compute_address` ([#4400](https://github.com/terraform-providers/terraform-provider-google/issues/4400))
* compute: add `mode` option to `google_compute_instance.boot_disk` ([#4413](https://github.com/terraform-providers/terraform-provider-google/issues/4413))
* compute: `google_compute_firewall` does not show a diff if allowed or denied rules are specified with uppercase protocol values ([#4467](https://github.com/terraform-providers/terraform-provider-google/issues/4467))
* logging: added `metric_descriptor.unit` to `google_logging_metric` resource ([#4407](https://github.com/terraform-providers/terraform-provider-google/issues/4407))

BUG FIXES:
* all: More classes of generic HTTP errors are retried provider-wide.
* container: Fix error when `master_authorized_networks_config` is removed from the `google_container_cluster` configuration. ([#4446](https://github.com/terraform-providers/terraform-provider-google/issues/4446))
* iam: Make `google_service_account_` and `google_service_account_iam_*` validation less restrictive to allow for more default service accounts ([#4377](https://github.com/terraform-providers/terraform-provider-google/issues/4377))
* iam: set auditconfigs in state for google_\*\_iam_policy resources ([#4447](https://github.com/terraform-providers/terraform-provider-google/issues/4447))
* logging: `google_logging_metric` `explicit` bucket option can now be set ([#4358](https://github.com/terraform-providers/terraform-provider-google/issues/4358))
* pubsub: Add retry for Pubsub Topic creation when project is still initializing org policies ([#4352](https://github.com/terraform-providers/terraform-provider-google/issues/4352))
* servicenetworking: remove need for provider-level project to delete connection ([#4445](https://github.com/terraform-providers/terraform-provider-google/issues/4445))
* sql: Add more retries for operationInProgress 409 errors for `google_sql_database_instance` ([#4376](https://github.com/terraform-providers/terraform-provider-google/issues/4376))

MISC:
* The User-Agent header that Terraform sends has been updated to correctly report the version of Terraform being run, and has minorly changed the formatting on the Terraform string. ([#4374](https://github.com/terraform-providers/terraform-provider-google/issues/4374))


## 2.14.0 (August 28, 2019)

DEPRECATIONS:
* cloudiot: `resource_cloudiot_registry`'s `event_notification_config` field has been deprecated. ([#4282](https://github.com/terraform-providers/terraform-provider-google/issues/4282))

FEATURES:
* **New Resource**: `google_bigtable_app_profile` is now available. ([#4126](https://github.com/terraform-providers/terraform-provider-google/issues/4126))
* **New Resource**: `google_ml_engine_model` ([#4053](https://github.com/terraform-providers/terraform-provider-google/issues/4053))
* **New Resource**: `google_dataproc_autoscaling_policy` ([#2220](https://github.com/terraform-providers/terraform-provider-google/issues/2220))
* **New Data Source**: `google_kms_secret_ciphertext` ([#4204](https://github.com/terraform-providers/terraform-provider-google/issues/4204))

ENHANCEMENTS:
* bigquery: Add support for clustering/partitioning to bigquery_table ([#4223](https://github.com/terraform-providers/terraform-provider-google/issues/4223))
* bigtable: `num_nodes` can now be updated in `google_bigtable_instance` ([#4026](https://github.com/terraform-providers/terraform-provider-google/issues/4026))
* cloudiot: `resource_cloudiot_registry` now has fields plural `event_notification_configs` and `log_level`, and `event_notification_config` has been deprecated. ([#4282](https://github.com/terraform-providers/terraform-provider-google/issues/4282))
* cloud_run: New output-only fields have been added to google_cloud_run_service' status. ([#3799](https://github.com/terraform-providers/terraform-provider-google/issues/3799))
* compute: Adding bandwidth attribute to interconnect attachment. ([#4212](https://github.com/terraform-providers/terraform-provider-google/issues/4212))
* compute: `google_compute_region_instance_group_manager.update_policy` now supports `instance_redistribution_type` ([#4301](https://github.com/terraform-providers/terraform-provider-google/issues/4301))
* compute: adds admin_enabled to google_compute_interconnect_attachment ([#4300](https://github.com/terraform-providers/terraform-provider-google/issues/4300))
* compute: The compute routes includes next_hop_ilb attribute support in beta. ([#4311](https://github.com/terraform-providers/terraform-provider-google/issues/4311))
* scheduler: Add support for `oauth_token` and `oidc_token` on resource `google_cloud_scheduler_job` ([#4222](https://github.com/terraform-providers/terraform-provider-google/issues/4222))

BUG FIXES:
* containerregistry: Correctly handle domain-scoped projects ([#4129](https://github.com/terraform-providers/terraform-provider-google/issues/4129))
* iam: Fixed regression in 2.13.0 for permadiff on empty members in IAM policy bindings. ([#4347](https://github.com/terraform-providers/terraform-provider-google/issues/4347))
* project: `google_project_iam_custom_role` now sets the project properly on import. ([#4343](https://github.com/terraform-providers/terraform-provider-google/issues/4343))
* sql: Added back a missing import format for `google_sql_database`. ([#4279](https://github.com/terraform-providers/terraform-provider-google/issues/4279))

## 2.13.0 (August 15, 2019)

KNOWN ISSUES:
* `bigtable`: `google_bigtable_instance` may cause a panic on Terraform `0.11`. This was resolved in `2.17.0`.

FEATURES:
* **New Resource**: added the `google_vpc_access_connector` resource and the `vpc_connector` option on the `google_cloudfunctions_function` resource. ([#4189](https://github.com/terraform-providers/terraform-provider-google/issues/4189))
* **New Resource**: Add `google_scc_source` resource for managing Cloud Security Command Center sources in Terraform ([#4236](https://github.com/terraform-providers/terraform-provider-google/issues/4236))
* **New Data Source**: `google_compute_network_endpoint_group` ([#4173](https://github.com/terraform-providers/terraform-provider-google/issues/4173))

ENHANCEMENTS:
* bigquery: Added support for `google_bigquery_data_transfer_config` (which include scheduled queries). ([#4102](https://github.com/terraform-providers/terraform-provider-google/issues/4102))
* bigtable: `google_bigtable_instance` max number of `cluster` blocks is now 4 ([#4156](https://github.com/terraform-providers/terraform-provider-google/issues/4156))
* binary_authorization: Added `globalPolicyEvaluationMode` to `google_binary_authorization_policy`. ([#4124](https://github.com/terraform-providers/terraform-provider-google/issues/4124))
* cloudfunctions: Allow partial URIs in google_cloudfunctions_function event_trigger.resource ([#4201](https://github.com/terraform-providers/terraform-provider-google/issues/4201))
* compute: Enable update for `google_compute_router_nat`
* netblock: Extended `google_netblock_ip_ranges` to supportmultiple useful IP address ranges that have a special meaning on GCP. ([#4121](https://github.com/terraform-providers/terraform-provider-google/issues/4121))
* project: Wrapped API requests with retries for `google_project`, `google_folder`, and `google_*_organization_policy` ([#4098](https://github.com/terraform-providers/terraform-provider-google/issues/4098))
* project: IAM and service requests are now batched ([#4207](https://github.com/terraform-providers/terraform-provider-google/issues/4207))
* provider: allow provider's region to be specified as a self_link ([#4219](https://github.com/terraform-providers/terraform-provider-google/issues/4219))
* provider: Adds new provider-level field `user_project_override`, which allows billing, quota checks, and service enablement checks to occur against the project a resource is in instead of the project the credentials are from. ([#4202](https://github.com/terraform-providers/terraform-provider-google/issues/4202))
* pubsub: Pub/Sub topic geo restriction support. ([#4131](https://github.com/terraform-providers/terraform-provider-google/issues/4131))

BUG FIXES:
* binary_authorization: don't diff when attestation authority note public keys don't have an ID in the config ([#4246](https://github.com/terraform-providers/terraform-provider-google/issues/4246))
* compute: google_compute_instance's description field is now set in state ([#4136](https://github.com/terraform-providers/terraform-provider-google/issues/4136))
* project: ignore errors when deleting a default network that doesn't exist ([#4137](https://github.com/terraform-providers/terraform-provider-google/issues/4137))

## 2.12.0 (August 01, 2019)

FEATURES:
* **New Data Source**: google_kms_crypto_key_version - Provides access to KMS key version data with Google Cloud KMS. ([#4078](https://github.com/terraform-providers/terraform-provider-google/issues/4078))
* **New Resource**: `google_cloud_run_service` - Set up a cloud run service ([#3714](https://github.com/terraform-providers/terraform-provider-google/issues/3714))
* **New Resource**: `google_cloud_run_domain_mapping` - Allows custom domains to map to a cloud run service ([#3714](https://github.com/terraform-providers/terraform-provider-google/issues/3714))
* `google_binary_authorization_attestor` and `google_binary_authorization_policy` are available in the GA provider ([#3960](https://github.com/terraform-providers/terraform-provider-google/issues/3960))

ENHANCEMENTS:
* binary_authorization: Adds support for Cloud KMS PKIX keys to `binary_authorization_attestor`. ([#4078](https://github.com/terraform-providers/terraform-provider-google/issues/4078))
* composer: Add private IP config for `google_composer_environment` ([#3952](https://github.com/terraform-providers/terraform-provider-google/issues/3952))
* compute: add support for port_specification to resource `google_compute_health_check` ([#4001](https://github.com/terraform-providers/terraform-provider-google/issues/4001))
* compute: Fixed import formats for `google_compute_network_endpoint` and add location-only import formats ([#4037](https://github.com/terraform-providers/terraform-provider-google/issues/4037))
* compute: Support labelling for compute_instance boot_disks and compute_instance_template disks. ([#4117](https://github.com/terraform-providers/terraform-provider-google/issues/4117))
* container: validate that master_ipv4_cidr_block is set if enable_private_nodes is true ([#4038](https://github.com/terraform-providers/terraform-provider-google/issues/4038))
* dataflow: added support for user-defined `labels` on resource `google_dataflow_job` ([#4095](https://github.com/terraform-providers/terraform-provider-google/issues/4095))
* dataproc: add support for `optional_components` to resource `resource_dataproc_cluster` ([#4073](https://github.com/terraform-providers/terraform-provider-google/issues/4073))
* project: add checks to import to prevent importing by project number instead of id ([#4051](https://github.com/terraform-providers/terraform-provider-google/issues/4051))
* storage: add support for `retention_policy` to resource `google_storage_bucket` ([#4044](https://github.com/terraform-providers/terraform-provider-google/issues/4044))

BUG FIXES:
* access_context_manager: import format checking ([#4047](https://github.com/terraform-providers/terraform-provider-google/issues/4047))
dataproc: Suppress diff for `google_dataproc_cluster` `software_config.0.image_version` to prevent permadiff when server uses more specific versions of config value ([#4088](https://github.com/terraform-providers/terraform-provider-google/issues/4088))
* organization: Add auditConfigs to update masks for setting org and folder IAM policy (`google_organization_iam_policy`, `google_folder_iam_policy`) ([#4084](https://github.com/terraform-providers/terraform-provider-google/issues/4084))
* storage: `google_storage_bucket` Set website metadata during read ([#3977](https://github.com/terraform-providers/terraform-provider-google/issues/3977))

## 2.11.0 (July 16, 2019)

NOTES:
* container: We have changed the way container clusters handle cluster state, and they should now wait until the cluster is ready when creating, updating, or refreshing cluster state. This is meant to decrease the frequency of errors where Terraform is operating on a cluster that isn't ready to be operated on. If this change causes a problem, please open an issue with as much information as you can provide, especially [debug logs](https://www.terraform.io/docs/internals/debugging.html). See [[#3989](https://github.com/terraform-providers/terraform-provider-google/issues/3989)] for more info.

FEATURES:
* **New Resources**: `google_bigtable_instance_iam_binding`, `google_bigtable_instance_iam_member`, and `google_bigtable_instance_iam_policy` are now available. ([#3939](https://github.com/terraform-providers/terraform-provider-google/issues/3939))
* **New Resources**: Add support for source repo repository IAM resources `google_sourcerepo_repository_iam_*` ([#3961](https://github.com/terraform-providers/terraform-provider-google/issues/3961))

ENHANCEMENTS:
* bigquery: Added support for `external_data_configuration` to `google_bigquery_table`. ([#3602](https://github.com/terraform-providers/terraform-provider-google/issues/3602))
* compute: Avoid getting project if no diff found for `google_compute_instance_template` ([#4000](https://github.com/terraform-providers/terraform-provider-google/issues/4000))
* firestore: `google_firestore_index` `query_scope` can have `COLLECTION_GROUP` specified. ([#3972](https://github.com/terraform-providers/terraform-provider-google/issues/3972))

BUG FIXES:
* compute: Allow security policy to be removed from `google_backend_service` ([#3969](https://github.com/terraform-providers/terraform-provider-google/issues/3969))
* compute: Mark instance KMS self link field `kms_key_self_link` as computed ([#3802](https://github.com/terraform-providers/terraform-provider-google/issues/3802))
* container: Fix panic for nil nested objects when reading cluster maintenance window ([#4002](https://github.com/terraform-providers/terraform-provider-google/issues/4002))
* container: `google_container_cluster` keep clusters in state if they are created in an error state and don't get correctly cleaned up. ([#3995](https://github.com/terraform-providers/terraform-provider-google/issues/3995))
* container: `google_container_cluster` will now wait to act until the cluster can be operated on, respecting timeouts. ([#3989](https://github.com/terraform-providers/terraform-provider-google/issues/3989))
* container: `google_container_node_pool` Correctly set nodepool autoscaling in state when disabled in the API ([#3997](https://github.com/terraform-providers/terraform-provider-google/issues/3997))
* monitoring: Fix diff in `google_monitoring_uptime_check_config` on a deprecated field. ([#4019](https://github.com/terraform-providers/terraform-provider-google/issues/4019))
* servicenetworking: `google_service_networking_connection` correctly delete the connection when the resource is destroyed. ([#4003](https://github.com/terraform-providers/terraform-provider-google/issues/4003))
* spanner: Wait for spanner databases to create before returning. Don't wait for databases to delete before returning anymore. ([#3975](https://github.com/terraform-providers/terraform-provider-google/issues/3975))
* storage: Fixed an issue where `google_storage_transfer_job` `schedule_end_date` caused requests to fail if unset. ([#4005](https://github.com/terraform-providers/terraform-provider-google/issues/4005))
* storage: `google_storage_object_acl` Prevent panic when using interpolated object names. ([#3970](https://github.com/terraform-providers/terraform-provider-google/issues/3970))

## 2.10.0 (July 02, 2019)

DEPRECATIONS:
* monitoring: Deprecated non-existent fields `is_internal` and `internal_checkers` from `google_monitoring_uptime_check_config`. ([#3919](https://github.com/terraform-providers/terraform-provider-google/issues/3919))

FEATURES:
* **New Resource**: `google_compute_project_default_network_tier` ([#3907](https://github.com/terraform-providers/terraform-provider-google/issues/3907))

ENHANCEMENTS:
* compute: Added fields for managing network endpoint group backends in `google_compute_backend_service`, including `max_connections_per_endpoint` and `max_rate_per_endpoint` ([#3863](https://github.com/terraform-providers/terraform-provider-google/issues/3863))
* compute: Support custom timeouts in `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager` ([#3955](https://github.com/terraform-providers/terraform-provider-google/issues/3955))
* container: `logging.googleapis.com/kubernetes` and `monitoring.googleapis.com/kubernetes` are now GA for cluster logging/monitoring service
* folder: `google_folder` improve error message on delete ([#3902](https://github.com/terraform-providers/terraform-provider-google/issues/3902))
* iam: sort bindings in `google_*_iam_policy` resources to get simpler diffs ([#3855](https://github.com/terraform-providers/terraform-provider-google/issues/3855))
* kms: `google_kms_crypto_key` now supports labels. ([#3910](https://github.com/terraform-providers/terraform-provider-google/issues/3910))
* pubsub: `google_pubsub_topic` supports KMS keys with `kms_key_name`. ([#3925](https://github.com/terraform-providers/terraform-provider-google/issues/3925))

BUG FIXES:
* iam: the member field in iam_* resources is now case-insensitive ([#3900](https://github.com/terraform-providers/terraform-provider-google/issues/3900))
* servicenetworking: `google_service_networking_connection` fix update ([#3887](https://github.com/terraform-providers/terraform-provider-google/issues/3887))

## 2.9.1 (June 21, 2019)

BUG FIXES:
* kms: fix regression when reading existing `google_kms_crypto_key` resources ([#3893](https://github.com/terraform-providers/terraform-provider-google/issues/3893))
* storage: `google_storage_bucket` fix for crash that occurs when running plan on old buckets ([#3886](https://github.com/terraform-providers/terraform-provider-google/issues/3886))
* storage: `google_storage_bucket` allow updating `bucket_policy_only` to false ([#3886](https://github.com/terraform-providers/terraform-provider-google/issues/3886))

## 2.9.0 (June 19, 2019)

FEATURES:
* **Custom Endpoint Support**: The Google provider supports custom endpoints, allowing you to use GCP-like APIs such as emulators. See the [Provider Reference](https://www.terraform.io/docs/providers/google/provider_reference.html) for details. ([#3787](https://github.com/terraform-providers/terraform-provider-google/issues/3787))
* **New Resource** Network endpoint groups (`google_compute_network_endpoint_group`) and fine-grained resource endpoints (`google_compute_network_endpoint`) are now available. ([#3832](https://github.com/terraform-providers/terraform-provider-google/issues/3832))
* **New Resource** `google_service_networking_connection` is now available (previously beta-only)


ENHANCEMENTS:
* increased default timeouts for `google_compute_instance`, `google_container_cluster`, `google_dataproc_cluster`, and `google_sql_database_instance` ([#3872](https://github.com/terraform-providers/terraform-provider-google/issues/3872))
* compute: `google_compute_global_address` supports `prefix_length`, `purpose`, and `network` ([#3877](https://github.com/terraform-providers/terraform-provider-google/issues/3877))
* dns: `google_dns_record_set`: allow importing dns record sets in any project ([#3862](https://github.com/terraform-providers/terraform-provider-google/issues/3862))
* kms: `kms_crypto_key` supports `purpose` ([#3843](https://github.com/terraform-providers/terraform-provider-google/issues/3843))
* storage: `google_storage_bucket` now supports enabling `bucket_policy_only` access control. ([#1878](https://github.com/terraform-providers/terraform-provider-google/pull/1878)
* storage: IAM resources for storage buckets (`google_storage_bucket_iam_*`) now all support import ([#3830](https://github.com/terraform-providers/terraform-provider-google/issues/3830))
* pubsub: `google_pubsub_topic` Updates for labels are now supported ([#3828](https://github.com/terraform-providers/terraform-provider-google/issues/3828))


BUG FIXES:
* bigquery: `google_bigquery_dataset` Relax IAM role restrictions on BQ datasets ([#3451](https://github.com/terraform-providers/terraform-provider-google/issues/3451))
* compute: `google_project_iam` When importing resources `project` no longer needs to be set in the config post import ([#3777](https://github.com/terraform-providers/terraform-provider-google/issues/3777))
* compute: `google_compute_instance_template` Fixed issue so project can now be specified by interpolated varibles. ([#3798](https://github.com/terraform-providers/terraform-provider-google/issues/3798))
* compute: `google_compute_instance_template` Throw error when using incompatible disk fields instead of continual plan diff ([#3789](https://github.com/terraform-providers/terraform-provider-google/issues/3789))
* compute: `google_compute_instance_from_template` Make sure disk type is expanded to a URL ([#3717](https://github.com/terraform-providers/terraform-provider-google/issues/3717))
* compute: `google_compute_instance_template` Attempt to put disks in state in the same order they were specified ([#3717](https://github.com/terraform-providers/terraform-provider-google/issues/3717))
* container: `google_container_cluster` Stop guest_accelerator from having a permadiff for accelerators with `count=0` ([#3860](https://github.com/terraform-providers/terraform-provider-google/issues/3860))
* container: `google_container_cluster` and `google_node_pool` now retry correctly when polling for status of an operation. ([#3801](https://github.com/terraform-providers/terraform-provider-google/issues/3801))
* dns: `google_dns_record_set` overrides all existing record types on create, not just NS ([#3859](https://github.com/terraform-providers/terraform-provider-google/issues/3859))
* monitoring: `google_monitoring_notification_channel` Allow setting enabled to false ([#3874](https://github.com/terraform-providers/terraform-provider-google/issues/3874))
* pubsub: `google_pubsub_subscription` and `google_pubsub_topic` resources can be created inside VPC service controls. ([#3818](https://github.com/terraform-providers/terraform-provider-google/issues/3818))
* redis: `google_redis_instance` Fall back to region from `location_id` when region isn't specified ([#3846](https://github.com/terraform-providers/terraform-provider-google/issues/3846))
* sql: `google_sql_user` User's can now be updated to change their password ([#3785](https://github.com/terraform-providers/terraform-provider-google/issues/3785))
* sql: Providing an non-empty host for a Postgres `google_sql_user` now correctly actually registers that the user was created and gives a slightly more understandable error/diff, instead of returning a generic "provider error" ([#3857](https://github.com/terraform-providers/terraform-provider-google/issues/3857))

## 2.8.0 (June 04, 2019)


DEPRECATIONS:
* compute: The `auto_create_routes` field on `google_compute_network_peering` has been deprecated because it is not user configurable. ([#3394](https://github.com/terraform-providers/terraform-provider-google/issues/3394))

FEATURES:
* **New Datasource**: `google_compute_ssl_certificate`  ([#3683](https://github.com/terraform-providers/terraform-provider-google/pull/3683))
* **New Datasource**: `google_composer_image_versions` ([#3694](https://github.com/terraform-providers/terraform-provider-google/pull/3694))

ENHANCEMENTS:
* app_engine: Update allowed `app_engine_application` locations. ([#3674](https://github.com/terraform-providers/terraform-provider-google/pull/3674))
* composer: Make `google_composer_environment` image version updateable. ([#3681](https://github.com/terraform-providers/terraform-provider-google/pull/3681))
* compute: `google_compute_router_interface` now supports specifying an `interconnect_attachment`. ([#3715](https://github.com/terraform-providers/terraform-provider-google/pull/3715))
* compute: `google_compute_router_nat` now supports specifying a `log_config` block ([#3684](https://github.com/terraform-providers/terraform-provider-google/pull/3684))
* compute: `google_compute_router_nat` now supports more import formats. ([#3744](https://github.com/terraform-providers/terraform-provider-google/pull/3744))
* compute: `google_compute_network_peering` now supports importing/exporting custom routes ([#3699](https://github.com/terraform-providers/terraform-provider-google/pull/3699))
* compute: Add support for INTERNAL_SELF_MANAGED backend services. Changed Resources: `google_compute_backend_service`, `google_compute_global_forwarding_rule`. ([#3719](https://github.com/terraform-providers/terraform-provider-google/pull/3719))
* container: Expose the `services_ipv4_cidr` for `container_cluster`. ([#3776](https://github.com/terraform-providers/terraform-provider-google/pull/3776))
* dns: `google_dns_managed_zone` now supports DNSSec. ([#3677](https://github.com/terraform-providers/terraform-provider-google/pull/3677))
* dataflow: `google_dataflow_job` now supports setting machine type ([#1862](https://github.com/GoogleCloudPlatform/magic-modules/pull/1862))
* kms: `google_kms_key_ring` is now autogenerated using Magic Modules ([#3689](https://github.com/terraform-providers/terraform-provider-google/pull/3689))
* pubsub: `google_pubsub_subscription` supports setting an `expiration_policy` with no `ttl`. ([#3742](https://github.com/terraform-providers/terraform-provider-google/pull/3742))

BUG FIXES:
* compute: Allow setting firewall priority to 0. ([#3700](https://github.com/terraform-providers/terraform-provider-google/pull/3700))
* compute: Resolved an issue where `google_compute_region_backend_service` was unable to perform a state migration. ([#3731](https://github.com/terraform-providers/terraform-provider-google/pull/3731))
* compute: Allow empty metadata.startup-script on instances. ([#3732](https://github.com/terraform-providers/terraform-provider-google/pull/3732))
* compute: Fix expanding of routing config in `google_compute_network`. ([#3741](https://github.com/terraform-providers/terraform-provider-google/pull/3741))
* container: Allow going from no ip_allocation_policy to a blank-equivalent one. ([#3723](https://github.com/terraform-providers/terraform-provider-google/pull/3723))
* container: `google_container_cluster` will no longer diff unnecessarily on `issue_client_certificate`. ([#3751](https://github.com/terraform-providers/terraform-provider-google/pull/3751))
* container: `google_container_cluster` can enable client certificates on GKE `1.12+` series releases. ([#3751](https://github.com/terraform-providers/terraform-provider-google/pull/3751))
* container: `google_container_cluster` now retries the call to remove default node pools during cluster creation ([#3769](https://github.com/terraform-providers/terraform-provider-google/pull/3769))
* storage: Fix occasional crash when updating storage buckets ([#3686](https://github.com/terraform-providers/terraform-provider-google/pull/3686))

## 2.7.0 (May 21, 2019)

NOTE:
* Several resources were previously undocumented on the site or changelog; they should be added to both with this release. `google_compute_backend_bucket_signed_url_key` and `google_compute_backend_service_signed_url_key` were introduced in `2.4.0`.

BACKWARDS INCOMPATIBILITIES:
* cloudfunctions: `google_cloudfunctions_function.runtime` now has an explicit default value of `nodejs6`. Users who have a different value set in the API but the value undefined in their config will see a diff. ([#3605](https://github.com/terraform-providers/terraform-provider-google/issues/3605))

FEATURES:
* **New Resources**: `google_compute_instance_iam_binding`, `google_compute_instance_iam_member`, and `google_compute_instance_iam_policy` are now available. ([#3551](https://github.com/terraform-providers/terraform-provider-google/pull/3551))
* **New Resources**: IAM resources for Dataproc jobs and clusters (`google_dataproc_job_iam_policy`, `google_dataproc_job_iam_member`, `google_dataproc_job_iam_binding`, `google_dataproc_cluster_iam_policy`, `google_dataproc_cluster_iam_member`, `google_dataproc_cluster_iam_binding`) are now available. [#3632](https://github.com/terraform-providers/terraform-provider-google/pull/3632)

ENHANCEMENTS:
* provider: Add GCP zone to `google_client_config` datasource ([#3262](https://github.com/terraform-providers/terraform-provider-google/issues/3262))
* compute: `google_compute_backend_service` now supports `HTTP2` protocol (beta-only feature, use with GA provider at own risk)[#3631](https://github.com/terraform-providers/terraform-provider-google/pull/3631)
* compute: `interconnect_attachment` Make vlanTag8021q computed for using PARTNER attachments ([#3600](https://github.com/terraform-providers/terraform-provider-google/issues/3600))
* compute: Add support for creating instances with CMEK ([#3481](https://github.com/terraform-providers/terraform-provider-google/issues/3481))
* compute: Can now specify project when importing instance groups ([#2504](https://github.com/terraform-providers/terraform-provider-google/issues/2504))
* compute: `google_compute_organization_policies*` Allow all organization policies to be removed/unset from a constraint. ([#3611](https://github.com/terraform-providers/terraform-provider-google/issues/3611))
* compute: `google_compute_instance` now supports `shielded_instance_config` for verifiable integrity of your VM instances. ([#3531](https://github.com/terraform-providers/terraform-provider-google/issues/3531))
* compute: `google_compute_instance_template` now supports `shielded_instance_config` for verifiable integrity of your VM instances. ([#3531](https://github.com/terraform-providers/terraform-provider-google/issues/3531))
* container: use the cluster subnet to look up the node cidr block ([#3654](https://github.com/terraform-providers/terraform-provider-google/issues/3654))

BUG FIXES:
* cloudfunctions: `google_cloudfunctions_function.runtime` now has an explicit default value of `nodejs6`. ([#3605](https://github.com/terraform-providers/terraform-provider-google/issues/3605))
* compute: Fix panic in `compute_backend_service` hash function ([#3610](https://github.com/terraform-providers/terraform-provider-google/issues/3610))
* monitoring: updating `google_monitoring_alert_policy` is more likely to succeed ([#3587](https://github.com/terraform-providers/terraform-provider-google/issues/3587))
* kms: `google_kms_crypto_key` now (in addition to marking all crypto key versions for destruction) correctly disables auto-rotation for destroyed keys [[#3624](https://github.com/terraform-providers/terraform-provider-google/issues/3624)](https://github.com/terraform-providers/terraform-provider-google/pull/3624)
* iam: Increase IAM custom role length validation to match API. ([#3660](https://github.com/terraform-providers/terraform-provider-google/issues/3660))

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
* **New Resource**: `google_compute_backend_bucket_signed_url_key` is now available. [GH-3229]
* **New Resource**: `google_compute_backend_service_signed_url_key` is now available. [GH-3359]
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
