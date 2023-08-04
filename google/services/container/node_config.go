// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"google.golang.org/api/container/v1"
)

// Matches gke-default scope from https://cloud.google.com/sdk/gcloud/reference/container/clusters/create
var defaultOauthScopes = []string{
	"https://www.googleapis.com/auth/devstorage.read_only",
	"https://www.googleapis.com/auth/logging.write",
	"https://www.googleapis.com/auth/monitoring",
	"https://www.googleapis.com/auth/service.management.readonly",
	"https://www.googleapis.com/auth/servicecontrol",
	"https://www.googleapis.com/auth/trace.append",
}

func schemaLoggingVariant() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  `Type of logging agent that is used as the default value for node pools in the cluster. Valid values include DEFAULT and MAX_THROUGHPUT.`,
		Default:      "DEFAULT",
		ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "MAX_THROUGHPUT"}, false),
	}
}

func schemaGcfsConfig(forceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `GCFS configuration for this node.`,
		ForceNew:    forceNew,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:        schema.TypeBool,
					Required:    true,
					ForceNew:    forceNew,
					Description: `Whether or not GCFS is enabled`,
				},
			},
		},
	}
}

func schemaNodeConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The configuration of the nodepool`,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"disk_size_gb": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntAtLeast(10),
					Description:  `Size of the disk attached to each node, specified in GB. The smallest allowed disk size is 10GB.`,
				},

				"disk_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `Type of the disk attached to each node. Such as pd-standard, pd-balanced or pd-ssd`,
				},

				"guest_accelerator": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					ForceNew: true,
					// Legacy config mode allows removing GPU's from an existing resource
					// See https://www.terraform.io/docs/configuration/attr-as-blocks.html
					ConfigMode:  schema.SchemaConfigModeAttr,
					Description: `List of the type and count of accelerator cards attached to the instance.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"count": {
								Type:        schema.TypeInt,
								Required:    true,
								ForceNew:    true,
								Description: `The number of the accelerator cards exposed to an instance.`,
							},
							"type": {
								Type:             schema.TypeString,
								Required:         true,
								ForceNew:         true,
								DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
								Description:      `The accelerator type resource name.`,
							},
							"gpu_driver_installation_config": {
								Type:        schema.TypeList,
								MaxItems:    1,
								Optional:    true,
								ForceNew:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: `Configuration for auto installation of GPU driver.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"gpu_driver_version": {
											Type:         schema.TypeString,
											Required:     true,
											ForceNew:     true,
											Description:  `Mode for how the GPU driver is installed.`,
											ValidateFunc: validation.StringInSlice([]string{"GPU_DRIVER_VERSION_UNSPECIFIED", "INSTALLATION_DISABLED", "DEFAULT", "LATEST"}, false),
										},
									},
								},
							},
							"gpu_partition_size": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `Size of partitions to create on the GPU. Valid values are described in the NVIDIA mig user guide (https://docs.nvidia.com/datacenter/tesla/mig-user-guide/#partitioning)`,
							},
							"gpu_sharing_config": {
								Type:        schema.TypeList,
								MaxItems:    1,
								Optional:    true,
								ForceNew:    true,
								ConfigMode:  schema.SchemaConfigModeAttr,
								Description: `Configuration for GPU sharing.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"gpu_sharing_strategy": {
											Type:        schema.TypeString,
											Required:    true,
											ForceNew:    true,
											Description: `The type of GPU sharing strategy to enable on the GPU node. Possible values are described in the API package (https://pkg.go.dev/google.golang.org/api/container/v1#GPUSharingConfig)`,
										},
										"max_shared_clients_per_gpu": {
											Type:        schema.TypeInt,
											Required:    true,
											ForceNew:    true,
											Description: `The maximum number of containers that can share a GPU.`,
										},
									},
								},
							},
						},
					},
				},

				"image_type": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					DiffSuppressFunc: tpgresource.CaseDiffSuppress,
					Description:      `The image type to use for this node. Note that for a given image type, the latest version of it will be used.`,
				},

				"labels": {
					Type:     schema.TypeMap,
					Optional: true,
					// Computed=true because GKE Sandbox will automatically add labels to nodes that can/cannot run sandboxed pods.
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The map of Kubernetes labels (key/value pairs) to be applied to each node. These will added in addition to any default label(s) that Kubernetes may apply to the node.`,
				},

				"resource_labels": {
					Type:        schema.TypeMap,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The GCE resource labels (a map of key/value pairs) to be applied to the node pool.`,
				},

				"local_ssd_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `The number of local SSD disks to be attached to the node.`,
				},

				"logging_variant": schemaLoggingVariant(),

				"ephemeral_storage_local_ssd_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Parameters for the ephemeral storage filesystem. If unspecified, ephemeral storage is backed by the boot disk.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"local_ssd_count": {
								Type:         schema.TypeInt,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(0),
								Description:  `Number of local SSDs to use to back ephemeral storage. Uses NVMe interfaces. Each local SSD must be 375 or 3000 GB in size, and all local SSDs must share the same size.`,
							},
						},
					},
				},

				"local_nvme_ssd_block_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Parameters for raw-block local NVMe SSDs.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"local_ssd_count": {
								Type:         schema.TypeInt,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(0),
								Description:  `Number of raw-block local NVMe SSD disks to be attached to the node. Each local SSD is 375 GB in size.`,
							},
						},
					},
				},

				"gcfs_config": schemaGcfsConfig(true),

				"gvnic": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Enable or disable gvnic in the node pool.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:        schema.TypeBool,
								Required:    true,
								ForceNew:    true,
								Description: `Whether or not gvnic is enabled`,
							},
						},
					},
				},

				"machine_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The name of a Google Compute Engine machine type.`,
				},

				"metadata": {
					Type:        schema.TypeMap,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The metadata key/value pairs assigned to instances in the cluster.`,
				},

				"min_cpu_platform": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `Minimum CPU platform to be used by this instance. The instance may be scheduled on the specified or newer CPU platform.`,
				},

				"oauth_scopes": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The set of Google API scopes to be made available on all of the node VMs.`,
					Elem: &schema.Schema{
						Type: schema.TypeString,
						StateFunc: func(v interface{}) string {
							return tpgresource.CanonicalizeServiceScope(v.(string))
						},
					},
					DiffSuppressFunc: containerClusterAddedScopesSuppress,
					Set:              tpgresource.StringScopeHashcode,
				},

				"preemptible": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
					Description: `Whether the nodes are created as preemptible VM instances.`,
				},
				"reservation_affinity": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `The reservation affinity configuration for the node pool.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"consume_reservation_type": {
								Type:         schema.TypeString,
								Required:     true,
								ForceNew:     true,
								Description:  `Corresponds to the type of reservation consumption.`,
								ValidateFunc: validation.StringInSlice([]string{"UNSPECIFIED", "NO_RESERVATION", "ANY_RESERVATION", "SPECIFIC_RESERVATION"}, false),
							},
							"key": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `The label key of a reservation resource.`,
							},
							"values": {
								Type:        schema.TypeSet,
								Description: "The label values of the reservation resource.",
								ForceNew:    true,
								Optional:    true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"spot": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
					Description: `Whether the nodes are created as spot VM instances.`,
				},

				"service_account": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The Google Cloud Platform Service Account to be used by the node VMs.`,
				},

				"tags": {
					Type:        schema.TypeList,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The list of instance tags applied to all nodes.`,
				},

				"shielded_instance_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `Shielded Instance options.`,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enable_secure_boot": {
								Type:        schema.TypeBool,
								Optional:    true,
								ForceNew:    true,
								Default:     false,
								Description: `Defines whether the instance has Secure Boot enabled.`,
							},
							"enable_integrity_monitoring": {
								Type:        schema.TypeBool,
								Optional:    true,
								ForceNew:    true,
								Default:     true,
								Description: `Defines whether the instance has integrity monitoring enabled.`,
							},
						},
					},
				},

				"taint": {
					Type:     schema.TypeList,
					Optional: true,
					// Computed=true because GKE Sandbox will automatically add taints to nodes that can/cannot run sandboxed pods.
					Computed: true,
					ForceNew: true,
					// Legacy config mode allows explicitly defining an empty taint.
					// See https://www.terraform.io/docs/configuration/attr-as-blocks.html
					ConfigMode:  schema.SchemaConfigModeAttr,
					Description: `List of Kubernetes taints to be applied to each node.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
								Type:        schema.TypeString,
								Required:    true,
								ForceNew:    true,
								Description: `Key for taint.`,
							},
							"value": {
								Type:        schema.TypeString,
								Required:    true,
								ForceNew:    true,
								Description: `Value for taint.`,
							},
							"effect": {
								Type:         schema.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice([]string{"NO_SCHEDULE", "PREFER_NO_SCHEDULE", "NO_EXECUTE"}, false),
								Description:  `Effect for taint.`,
							},
						},
					},
				},

				"workload_metadata_config": {
					Computed:    true,
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `The workload metadata configuration for this node.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mode": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"MODE_UNSPECIFIED", "GCE_METADATA", "GKE_METADATA"}, false),
								Description:  `Mode is the configuration for how to expose metadata to workloads running on the node.`,
							},
						},
					},
				},

				"boot_disk_kms_key": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: `The Customer Managed Encryption Key used to encrypt the boot disk attached to each node in the node pool.`,
				},
				// Note that AtLeastOneOf can't be set because this schema is reused by
				// two different resources.
				"kubelet_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Node kubelet configs.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cpu_manager_policy": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"static", "none", ""}, false),
								Description:  `Control the CPU management policy on the node.`,
							},
							"cpu_cfs_quota": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: `Enable CPU CFS quota enforcement for containers that specify CPU limits.`,
							},
							"cpu_cfs_quota_period": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: `Set the CPU CFS quota period value 'cpu.cfs_period_us'.`,
							},
							"pod_pids_limit": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Controls the maximum number of processes allowed to run in a pod.`,
							},
						},
					},
				},

				"linux_node_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Parameters that can be configured on Linux nodes.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"sysctls": {
								Type:        schema.TypeMap,
								Required:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Description: `The Linux kernel parameters to be applied to the nodes and all pods running on the nodes.`,
							},
						},
					},
				},
				"node_group": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: `Setting this field will assign instances of this pool to run on the specified node group. This is useful for running workloads on sole tenant nodes.`,
				},

				"advanced_machine_features": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Specifies options for controlling advanced machine features.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"threads_per_core": {
								Type:        schema.TypeInt,
								Required:    true,
								ForceNew:    true,
								Description: `The number of threads per physical core. To disable simultaneous multithreading (SMT) set this to 1. If unset, the maximum number of threads supported per core by the underlying processor is assumed.`,
							},
						},
					},
				},
				"sole_tenant_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: `Node affinity options for sole tenant node pools.`,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"node_affinity": {
								Type:        schema.TypeSet,
								Required:    true,
								ForceNew:    true,
								Description: `.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"key": {
											Type:        schema.TypeString,
											Required:    true,
											ForceNew:    true,
											Description: `.`,
										},
										"operator": {
											Type:         schema.TypeString,
											Required:     true,
											ForceNew:     true,
											Description:  `.`,
											ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN"}, false),
										},
										"values": {
											Type:        schema.TypeList,
											Required:    true,
											ForceNew:    true,
											Description: `.`,
											Elem:        &schema.Schema{Type: schema.TypeString},
										},
									},
								},
							},
						},
					},
				},
				"host_maintenance_policy": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: `The maintenance policy for the hosts on which the GKE VMs run on.`,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"maintenance_interval": {
								Type:         schema.TypeString,
								Required:     true,
								ForceNew:     true,
								Description:  `.`,
								ValidateFunc: validation.StringInSlice([]string{"MAINTENANCE_INTERVAL_UNSPECIFIED", "AS_NEEDED", "PERIODIC"}, false),
							},
						},
					},
				},
			},
		},
	}
}

func expandNodeConfigDefaults(configured interface{}) *container.NodeConfigDefaults {
	configs := configured.([]interface{})
	if len(configs) == 0 || configs[0] == nil {
		return nil
	}
	config := configs[0].(map[string]interface{})

	nodeConfigDefaults := &container.NodeConfigDefaults{}
	if variant, ok := config["logging_variant"]; ok {
		nodeConfigDefaults.LoggingConfig = &container.NodePoolLoggingConfig{
			VariantConfig: &container.LoggingVariantConfig{
				Variant: variant.(string),
			},
		}
	}
	return nodeConfigDefaults
}

func expandNodeConfig(v interface{}) *container.NodeConfig {
	nodeConfigs := v.([]interface{})
	nc := &container.NodeConfig{
		// Defaults can't be set on a list/set in the schema, so set the default on create here.
		OauthScopes: defaultOauthScopes,
	}
	if len(nodeConfigs) == 0 {
		return nc
	}

	nodeConfig := nodeConfigs[0].(map[string]interface{})

	if v, ok := nodeConfig["machine_type"]; ok {
		nc.MachineType = v.(string)
	}

	if v, ok := nodeConfig["guest_accelerator"]; ok {
		accels := v.([]interface{})
		guestAccelerators := make([]*container.AcceleratorConfig, 0, len(accels))
		for _, raw := range accels {
			data := raw.(map[string]interface{})
			if data["count"].(int) == 0 {
				continue
			}
			guestAcceleratorConfig := &container.AcceleratorConfig{
				AcceleratorCount: int64(data["count"].(int)),
				AcceleratorType:  data["type"].(string),
				GpuPartitionSize: data["gpu_partition_size"].(string),
			}

			if v, ok := data["gpu_driver_installation_config"]; ok && len(v.([]interface{})) > 0 {
				gpuDriverInstallationConfig := data["gpu_driver_installation_config"].([]interface{})[0].(map[string]interface{})
				guestAcceleratorConfig.GpuDriverInstallationConfig = &container.GPUDriverInstallationConfig{
					GpuDriverVersion: gpuDriverInstallationConfig["gpu_driver_version"].(string),
				}
			}

			if v, ok := data["gpu_sharing_config"]; ok && len(v.([]interface{})) > 0 {
				gpuSharingConfig := data["gpu_sharing_config"].([]interface{})[0].(map[string]interface{})
				guestAcceleratorConfig.GpuSharingConfig = &container.GPUSharingConfig{
					GpuSharingStrategy:     gpuSharingConfig["gpu_sharing_strategy"].(string),
					MaxSharedClientsPerGpu: int64(gpuSharingConfig["max_shared_clients_per_gpu"].(int)),
				}
			}

			guestAccelerators = append(guestAccelerators, guestAcceleratorConfig)
		}
		nc.Accelerators = guestAccelerators
	}

	if v, ok := nodeConfig["disk_size_gb"]; ok {
		nc.DiskSizeGb = int64(v.(int))
	}

	if v, ok := nodeConfig["disk_type"]; ok {
		nc.DiskType = v.(string)
	}

	if v, ok := nodeConfig["local_ssd_count"]; ok {
		nc.LocalSsdCount = int64(v.(int))
	}

	if v, ok := nodeConfig["logging_variant"]; ok {
		nc.LoggingConfig = &container.NodePoolLoggingConfig{
			VariantConfig: &container.LoggingVariantConfig{
				Variant: v.(string),
			},
		}
	}

	if v, ok := nodeConfig["local_nvme_ssd_block_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.LocalNvmeSsdBlockConfig = &container.LocalNvmeSsdBlockConfig{
			LocalSsdCount: int64(conf["local_ssd_count"].(int)),
		}
	}

	if v, ok := nodeConfig["ephemeral_storage_local_ssd_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.EphemeralStorageLocalSsdConfig = &container.EphemeralStorageLocalSsdConfig{
			LocalSsdCount: int64(conf["local_ssd_count"].(int)),
		}
	}

	if v, ok := nodeConfig["gcfs_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.GcfsConfig = &container.GcfsConfig{
			Enabled: conf["enabled"].(bool),
		}
	}

	if v, ok := nodeConfig["gvnic"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.Gvnic = &container.VirtualNIC{
			Enabled: conf["enabled"].(bool),
		}
	}

	if v, ok := nodeConfig["reservation_affinity"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		valuesSet := conf["values"].(*schema.Set)
		values := make([]string, valuesSet.Len())
		for i, value := range valuesSet.List() {
			values[i] = value.(string)
		}

		nc.ReservationAffinity = &container.ReservationAffinity{
			ConsumeReservationType: conf["consume_reservation_type"].(string),
			Key:                    conf["key"].(string),
			Values:                 values,
		}
	}

	if scopes, ok := nodeConfig["oauth_scopes"]; ok {
		scopesSet := scopes.(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, scope := range scopesSet.List() {
			scopes[i] = tpgresource.CanonicalizeServiceScope(scope.(string))
		}

		nc.OauthScopes = scopes
	}

	if v, ok := nodeConfig["service_account"]; ok {
		nc.ServiceAccount = v.(string)
	}

	if v, ok := nodeConfig["metadata"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.Metadata = m
	}

	if v, ok := nodeConfig["image_type"]; ok {
		nc.ImageType = v.(string)
	}

	if v, ok := nodeConfig["labels"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.Labels = m
	}

	if v, ok := nodeConfig["resource_labels"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.ResourceLabels = m
	}

	if v, ok := nodeConfig["tags"]; ok {
		tagsList := v.([]interface{})
		tags := []string{}
		for _, v := range tagsList {
			if v != nil {
				tags = append(tags, v.(string))
			}
		}
		nc.Tags = tags
	}

	if v, ok := nodeConfig["shielded_instance_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.ShieldedInstanceConfig = &container.ShieldedInstanceConfig{
			EnableSecureBoot:          conf["enable_secure_boot"].(bool),
			EnableIntegrityMonitoring: conf["enable_integrity_monitoring"].(bool),
		}
	}

	// Preemptible Is Optional+Default, so it always has a value
	nc.Preemptible = nodeConfig["preemptible"].(bool)

	// Spot Is Optional+Default, so it always has a value
	nc.Spot = nodeConfig["spot"].(bool)

	if v, ok := nodeConfig["min_cpu_platform"]; ok {
		nc.MinCpuPlatform = v.(string)
	}

	if v, ok := nodeConfig["taint"]; ok && len(v.([]interface{})) > 0 {
		taints := v.([]interface{})
		nodeTaints := make([]*container.NodeTaint, 0, len(taints))
		for _, raw := range taints {
			data := raw.(map[string]interface{})
			taint := &container.NodeTaint{
				Key:    data["key"].(string),
				Value:  data["value"].(string),
				Effect: data["effect"].(string),
			}
			nodeTaints = append(nodeTaints, taint)
		}
		nc.Taints = nodeTaints
	}

	if v, ok := nodeConfig["workload_metadata_config"]; ok {
		nc.WorkloadMetadataConfig = expandWorkloadMetadataConfig(v)
	}

	if v, ok := nodeConfig["boot_disk_kms_key"]; ok {
		nc.BootDiskKmsKey = v.(string)
	}

	if v, ok := nodeConfig["kubelet_config"]; ok {
		nc.KubeletConfig = expandKubeletConfig(v)
	}

	if v, ok := nodeConfig["linux_node_config"]; ok {
		nc.LinuxNodeConfig = expandLinuxNodeConfig(v)
	}

	if v, ok := nodeConfig["node_group"]; ok {
		nc.NodeGroup = v.(string)
	}

	if v, ok := nodeConfig["advanced_machine_features"]; ok && len(v.([]interface{})) > 0 {
		advanced_machine_features := v.([]interface{})[0].(map[string]interface{})
		nc.AdvancedMachineFeatures = &container.AdvancedMachineFeatures{
			ThreadsPerCore: int64(advanced_machine_features["threads_per_core"].(int)),
		}
	}

	if v, ok := nodeConfig["sole_tenant_config"]; ok && len(v.([]interface{})) > 0 {
		nc.SoleTenantConfig = expandSoleTenantConfig(v)
	}

	return nc
}

func expandWorkloadMetadataConfig(v interface{}) *container.WorkloadMetadataConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	wmc := &container.WorkloadMetadataConfig{}

	cfg := ls[0].(map[string]interface{})

	if v, ok := cfg["mode"]; ok {
		wmc.Mode = v.(string)
	}

	return wmc
}

func expandKubeletConfig(v interface{}) *container.NodeKubeletConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	cfg := ls[0].(map[string]interface{})
	kConfig := &container.NodeKubeletConfig{}
	if cpuManagerPolicy, ok := cfg["cpu_manager_policy"]; ok {
		kConfig.CpuManagerPolicy = cpuManagerPolicy.(string)
	}
	if cpuCfsQuota, ok := cfg["cpu_cfs_quota"]; ok {
		kConfig.CpuCfsQuota = cpuCfsQuota.(bool)
		kConfig.ForceSendFields = append(kConfig.ForceSendFields, "CpuCfsQuota")
	}
	if cpuCfsQuotaPeriod, ok := cfg["cpu_cfs_quota_period"]; ok {
		kConfig.CpuCfsQuotaPeriod = cpuCfsQuotaPeriod.(string)
	}
	if podPidsLimit, ok := cfg["pod_pids_limit"]; ok {
		kConfig.PodPidsLimit = int64(podPidsLimit.(int))
	}
	return kConfig
}

func expandLinuxNodeConfig(v interface{}) *container.LinuxNodeConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	cfg := ls[0].(map[string]interface{})
	sysCfgRaw, ok := cfg["sysctls"]
	if !ok {
		return nil
	}
	m := make(map[string]string)
	for k, v := range sysCfgRaw.(map[string]interface{}) {
		m[k] = v.(string)
	}
	return &container.LinuxNodeConfig{
		Sysctls: m,
	}
}

func expandSoleTenantConfig(v interface{}) *container.SoleTenantConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	cfg := ls[0].(map[string]interface{})
	affinitiesRaw, ok := cfg["node_affinity"]
	if !ok {
		return nil
	}
	affinities := make([]*container.NodeAffinity, 0)
	for _, v := range affinitiesRaw.(*schema.Set).List() {
		na := v.(map[string]interface{})

		affinities = append(affinities, &container.NodeAffinity{
			Key:      na["key"].(string),
			Operator: na["operator"].(string),
			Values:   tpgresource.ConvertStringArr(na["values"].([]interface{})),
		})
	}
	return &container.SoleTenantConfig{
		NodeAffinities: affinities,
	}
}

func flattenNodeConfigDefaults(c *container.NodeConfigDefaults) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if c == nil {
		return result
	}

	result = append(result, map[string]interface{}{})

	result[0]["logging_variant"] = flattenLoggingVariant(c.LoggingConfig)

	return result
}

func flattenNodeConfig(c *container.NodeConfig) []map[string]interface{} {
	config := make([]map[string]interface{}, 0, 1)

	if c == nil {
		return config
	}

	config = append(config, map[string]interface{}{
		"machine_type":                       c.MachineType,
		"disk_size_gb":                       c.DiskSizeGb,
		"disk_type":                          c.DiskType,
		"guest_accelerator":                  flattenContainerGuestAccelerators(c.Accelerators),
		"local_ssd_count":                    c.LocalSsdCount,
		"logging_variant":                    flattenLoggingVariant(c.LoggingConfig),
		"local_nvme_ssd_block_config":        flattenLocalNvmeSsdBlockConfig(c.LocalNvmeSsdBlockConfig),
		"ephemeral_storage_local_ssd_config": flattenEphemeralStorageLocalSsdConfig(c.EphemeralStorageLocalSsdConfig),
		"gcfs_config":                        flattenGcfsConfig(c.GcfsConfig),
		"gvnic":                              flattenGvnic(c.Gvnic),
		"reservation_affinity":               flattenGKEReservationAffinity(c.ReservationAffinity),
		"service_account":                    c.ServiceAccount,
		"metadata":                           c.Metadata,
		"image_type":                         c.ImageType,
		"labels":                             c.Labels,
		"resource_labels":                    c.ResourceLabels,
		"tags":                               c.Tags,
		"preemptible":                        c.Preemptible,
		"spot":                               c.Spot,
		"min_cpu_platform":                   c.MinCpuPlatform,
		"shielded_instance_config":           flattenShieldedInstanceConfig(c.ShieldedInstanceConfig),
		"taint":                              flattenTaints(c.Taints),
		"workload_metadata_config":           flattenWorkloadMetadataConfig(c.WorkloadMetadataConfig),
		"boot_disk_kms_key":                  c.BootDiskKmsKey,
		"kubelet_config":                     flattenKubeletConfig(c.KubeletConfig),
		"linux_node_config":                  flattenLinuxNodeConfig(c.LinuxNodeConfig),
		"node_group":                         c.NodeGroup,
		"advanced_machine_features":          flattenAdvancedMachineFeaturesConfig(c.AdvancedMachineFeatures),
		"sole_tenant_config":                 flattenSoleTenantConfig(c.SoleTenantConfig),
	})

	if len(c.OauthScopes) > 0 {
		config[0]["oauth_scopes"] = schema.NewSet(tpgresource.StringScopeHashcode, tpgresource.ConvertStringArrToInterface(c.OauthScopes))
	}

	return config
}

func flattenAdvancedMachineFeaturesConfig(c *container.AdvancedMachineFeatures) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"threads_per_core": c.ThreadsPerCore,
		})
	}
	return result
}

func flattenContainerGuestAccelerators(c []*container.AcceleratorConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, accel := range c {
		accelerator := map[string]interface{}{
			"count":              accel.AcceleratorCount,
			"type":               accel.AcceleratorType,
			"gpu_partition_size": accel.GpuPartitionSize,
		}
		if accel.GpuDriverInstallationConfig != nil {
			accelerator["gpu_driver_installation_config"] = []map[string]interface{}{
				{
					"gpu_driver_version": accel.GpuDriverInstallationConfig.GpuDriverVersion,
				},
			}
		}
		if accel.GpuSharingConfig != nil {
			accelerator["gpu_sharing_config"] = []map[string]interface{}{
				{
					"gpu_sharing_strategy":       accel.GpuSharingConfig.GpuSharingStrategy,
					"max_shared_clients_per_gpu": accel.GpuSharingConfig.MaxSharedClientsPerGpu,
				},
			}
		}
		result = append(result, accelerator)
	}
	return result
}

func flattenShieldedInstanceConfig(c *container.ShieldedInstanceConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"enable_secure_boot":          c.EnableSecureBoot,
			"enable_integrity_monitoring": c.EnableIntegrityMonitoring,
		})
	}
	return result
}

func flattenLocalNvmeSsdBlockConfig(c *container.LocalNvmeSsdBlockConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"local_ssd_count": c.LocalSsdCount,
		})
	}
	return result
}

func flattenEphemeralStorageLocalSsdConfig(c *container.EphemeralStorageLocalSsdConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"local_ssd_count": c.LocalSsdCount,
		})
	}
	return result
}

func flattenLoggingVariant(c *container.NodePoolLoggingConfig) string {
	variant := "DEFAULT"
	if c != nil && c.VariantConfig != nil && c.VariantConfig.Variant != "" {
		variant = c.VariantConfig.Variant
	}
	return variant
}

func flattenGcfsConfig(c *container.GcfsConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"enabled": c.Enabled,
		})
	}
	return result
}

func flattenGvnic(c *container.VirtualNIC) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"enabled": c.Enabled,
		})
	}
	return result
}

func flattenGKEReservationAffinity(c *container.ReservationAffinity) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"consume_reservation_type": c.ConsumeReservationType,
			"key":                      c.Key,
			"values":                   c.Values,
		})
	}
	return result
}

func flattenTaints(c []*container.NodeTaint) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, taint := range c {
		result = append(result, map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		})
	}
	return result
}

func flattenWorkloadMetadataConfig(c *container.WorkloadMetadataConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"mode": c.Mode,
		})
	}
	return result
}

func flattenKubeletConfig(c *container.NodeKubeletConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"cpu_cfs_quota":        c.CpuCfsQuota,
			"cpu_cfs_quota_period": c.CpuCfsQuotaPeriod,
			"cpu_manager_policy":   c.CpuManagerPolicy,
			"pod_pids_limit":       c.PodPidsLimit,
		})
	}
	return result
}

func flattenLinuxNodeConfig(c *container.LinuxNodeConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"sysctls": c.Sysctls,
		})
	}
	return result
}

func flattenSoleTenantConfig(c *container.SoleTenantConfig) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c == nil {
		return result
	}
	affinities := []map[string]interface{}{}
	for _, affinity := range c.NodeAffinities {
		affinities = append(affinities, map[string]interface{}{
			"key":      affinity.Key,
			"operator": affinity.Operator,
			"values":   affinity.Values,
		})
	}
	return append(result, map[string]interface{}{
		"node_affinity": affinities,
	})
}
