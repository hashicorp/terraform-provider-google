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
* **New Resource**: `google_compute_backend_bucket_signed_url_key` is now available. ([#3229](https://github.com/terraform-providers/terraform-provider-google/issues/3229))
* **New Resource**: `google_compute_backend_service_signed_url_key` is now available. ([#3359](https://github.com/terraform-providers/terraform-provider-google/issues/3359))
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
