## 5.1.0 (Unreleased)

## 5.0.0 (October 2, 2023)

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