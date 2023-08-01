// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/hashstructure"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/compute/v1"
)

var (
	bootDiskKeys = []string{
		"boot_disk.0.auto_delete",
		"boot_disk.0.device_name",
		"boot_disk.0.disk_encryption_key_raw",
		"boot_disk.0.kms_key_self_link",
		"boot_disk.0.initialize_params",
		"boot_disk.0.mode",
		"boot_disk.0.source",
	}

	initializeParamsKeys = []string{
		"boot_disk.0.initialize_params.0.size",
		"boot_disk.0.initialize_params.0.type",
		"boot_disk.0.initialize_params.0.image",
		"boot_disk.0.initialize_params.0.labels",
		"boot_disk.0.initialize_params.0.resource_manager_tags",
	}

	schedulingKeys = []string{
		"scheduling.0.on_host_maintenance",
		"scheduling.0.automatic_restart",
		"scheduling.0.preemptible",
		"scheduling.0.node_affinities",
		"scheduling.0.min_node_cpus",
		"scheduling.0.provisioning_model",
		"scheduling.0.instance_termination_action",
		"scheduling.0.local_ssd_recovery_timeout",
	}

	shieldedInstanceConfigKeys = []string{
		"shielded_instance_config.0.enable_secure_boot",
		"shielded_instance_config.0.enable_vtpm",
		"shielded_instance_config.0.enable_integrity_monitoring",
	}
)

// network_interface.[d].network_ip can only change when subnet/network
// is also changing. Validate that if network_ip is changing this scenario
// holds up to par.
func forceNewIfNetworkIPNotUpdatable(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return forceNewIfNetworkIPNotUpdatableFunc(d)
}

func forceNewIfNetworkIPNotUpdatableFunc(d tpgresource.TerraformResourceDiff) error {
	oldCount, newCount := d.GetChange("network_interface.#")
	if oldCount.(int) != newCount.(int) {
		return nil
	}

	for i := 0; i < newCount.(int); i++ {
		prefix := fmt.Sprintf("network_interface.%d", i)
		networkKey := prefix + ".network"
		subnetworkKey := prefix + ".subnetwork"
		subnetworkProjectKey := prefix + ".subnetwork_project"
		networkIPKey := prefix + ".network_ip"
		if d.HasChange(networkIPKey) {
			if !d.HasChange(networkKey) && !d.HasChange(subnetworkKey) && !d.HasChange(subnetworkProjectKey) {
				if err := d.ForceNew(networkIPKey); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// User may specify AUTOMATIC using any case; the API will accept it and return an empty string.
func ComputeInstanceMinCpuPlatformEmptyOrAutomaticDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	old = strings.ToLower(old)
	new = strings.ToLower(new)
	defaultVal := "automatic"
	return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
}

func ResourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceCreate,
		Read:   resourceComputeInstanceRead,
		Update: resourceComputeInstanceUpdate,
		Delete: resourceComputeInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeInstanceImportState,
		},

		SchemaVersion: 6,
		MigrateState:  ResourceComputeInstanceMigrateState,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		// A compute instance is more or less a superset of a compute instance
		// template. Please attempt to maintain consistency with the
		// resource_compute_instance_template schema when updating this one.
		Schema: map[string]*schema.Schema{
			"boot_disk": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The boot disk for the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							Default:      true,
							ForceNew:     true,
							Description:  `Whether the disk will be auto-deleted when the instance is deleted.`,
						},

						"device_name": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							Computed:     true,
							ForceNew:     true,
							Description:  `Name with which attached disk will be accessible under /dev/disk/by-id/`,
						},

						"disk_encryption_key_raw": {
							Type:          schema.TypeString,
							Optional:      true,
							AtLeastOneOf:  bootDiskKeys,
							ForceNew:      true,
							ConflictsWith: []string{"boot_disk.0.kms_key_self_link"},
							Sensitive:     true,
							Description:   `A 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to encrypt this disk. Only one of kms_key_self_link and disk_encryption_key_raw may be set.`,
						},

						"disk_encryption_key_sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.`,
						},

						"kms_key_self_link": {
							Type:             schema.TypeString,
							Optional:         true,
							AtLeastOneOf:     bootDiskKeys,
							ForceNew:         true,
							ConflictsWith:    []string{"boot_disk.0.disk_encryption_key_raw"},
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Computed:         true,
							Description:      `The self_link of the encryption key that is stored in Google Cloud KMS to encrypt this disk. Only one of kms_key_self_link and disk_encryption_key_raw may be set.`,
						},

						"initialize_params": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							Computed:     true,
							ForceNew:     true,
							MaxItems:     1,
							Description:  `Parameters with which a disk was created alongside the instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:         schema.TypeInt,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntAtLeast(1),
										Description:  `The size of the image in gigabytes.`,
									},

									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `The Google Compute Engine disk type. Such as pd-standard, pd-ssd or pd-balanced.`,
									},

									"image": {
										Type:             schema.TypeString,
										Optional:         true,
										AtLeastOneOf:     initializeParamsKeys,
										Computed:         true,
										ForceNew:         true,
										DiffSuppressFunc: DiskImageDiffSuppress,
										Description:      `The image from which this disk was initialised.`,
									},

									"labels": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `A set of key/value label pairs assigned to the disk.`,
									},

									"resource_manager_tags": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										ForceNew:     true,
										Description:  `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
									},
								},
							},
						},

						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							ForceNew:     true,
							Default:      "READ_WRITE",
							ValidateFunc: validation.StringInSlice([]string{"READ_WRITE", "READ_ONLY"}, false),
							Description:  `Read/write mode for the disk. One of "READ_ONLY" or "READ_WRITE".`,
						},

						"source": {
							Type:             schema.TypeString,
							Optional:         true,
							AtLeastOneOf:     bootDiskKeys,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    []string{"boot_disk.initialize_params"},
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the disk attached to this instance.`,
						},
					},
				},
			},

			"machine_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The machine type to create.`,
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the instance. One of name or self_link must be provided.`,
			},

			"network_interface": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `The networks attached to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the network attached to this interface.`,
						},

						"subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the subnetwork attached to this interface.`,
						},

						"subnetwork_project": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The project in which the subnetwork belongs.`,
						},

						"network_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The private IP address assigned to the instance.`,
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the interface`,
						},
						"nic_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"GVNIC", "VIRTIO_NET"}, false),
							Description:  `The type of vNIC to be used on this interface. Possible values:GVNIC, VIRTIO_NET`,
						},
						"access_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Access configurations, i.e. IPs via which this instance can be accessed via the Internet.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The IP address that is be 1:1 mapped to the instance's network ip.`,
									},

									"network_tier": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The networking tier used for configuring this instance. One of PREMIUM or STANDARD.`,
									},

									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The DNS domain name for the public PTR record.`,
									},
								},
							},
						},

						"alias_ip_range": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `An array of alias IP ranges for this network interface.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_cidr_range": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: tpgresource.IpCidrRangeDiffSuppress,
										Description:      `The IP CIDR range represented by this alias IP range.`,
									},
									"subnetwork_range_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The subnetwork secondary range name specifying the secondary range from which to allocate the IP CIDR range for this alias IP range.`,
									},
								},
							},
						},

						"stack_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"IPV4_ONLY", "IPV4_IPV6", ""}, false),
							Description:  `The stack type for this network interface to identify whether the IPv6 feature is enabled or not. If not specified, IPV4_ONLY will be used.`,
						},

						"ipv6_access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `One of EXTERNAL, INTERNAL to indicate whether the IP can be accessed from the Internet. This field is always inherited from its subnetwork.`,
						},

						"ipv6_access_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `An array of IPv6 access configurations for this interface. Currently, only one IPv6 access config, DIRECT_IPV6, is supported. If there is no ipv6AccessConfig specified, then this instance will have no external IPv6 Internet access.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_tier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The service-level to be provided for IPv6 traffic when the subnet has an external subnet. Only PREMIUM tier is valid for IPv6`,
									},
									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The domain name to be used when creating DNSv6 records for the external IPv6 ranges.`,
									},
									"external_ipv6": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										ForceNew:         true,
										DiffSuppressFunc: ipv6RepresentationDiffSuppress,
										Description:      `The first IPv6 address of the external IPv6 range associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig. To use a static external IP address, it must be unused and in the same region as the instance's zone. If not specified, Google Cloud will automatically assign an external IPv6 address from the instance's subnetwork.`,
									},
									"external_ipv6_prefix_length": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: `The prefix length of the external IPv6 range.`,
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: `The name of this access configuration. In ipv6AccessConfigs, the recommended name is External IPv6.`,
									},
								},
							},
						},

						"queue_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `The networking queue count that's specified by users for the network interface. Both Rx and Tx queues will be set to this number. It will be empty if not specified.`,
						},
					},
				},
			},
			"network_performance_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Configures network performance settings for the instance. If not specified, the instance will be created with its default network performance configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_egress_bandwidth_tier": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"TIER_1", "DEFAULT"}, false),
							Description:  `The egress bandwidth tier to enable. Possible values:TIER_1, DEFAULT`,
						},
					},
				},
			},
			"allow_stopping_for_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `If true, allows Terraform to stop the instance to update its properties. If you try to update a property that requires stopping the instance without setting this field, the update will fail.`,
			},

			"attached_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `List of disks attached to the instance`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the disk attached to this instance.`,
						},

						"device_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Name with which the attached disk is accessible under /dev/disk/by-id/`,
						},

						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "READ_WRITE",
							ValidateFunc: validation.StringInSlice([]string{"READ_WRITE", "READ_ONLY"}, false),
							Description:  `Read/write mode for the disk. One of "READ_ONLY" or "READ_WRITE".`,
						},

						"disk_encryption_key_raw": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `A 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to encrypt this disk. Only one of kms_key_self_link and disk_encryption_key_raw may be set.`,
						},

						"kms_key_self_link": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Computed:         true,
							Description:      `The self_link of the encryption key that is stored in Google Cloud KMS to encrypt this disk. Only one of kms_key_self_link and disk_encryption_key_raw may be set.`,
						},

						"disk_encryption_key_sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.`,
						},
					},
				},
			},

			"can_ip_forward": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether sending and receiving of packets with non-matching source or destination IPs is allowed.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A brief description of the resource.`,
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether deletion protection is enabled on this instance.`,
			},

			"enable_display": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the instance has virtual displays enabled.`,
			},

			"guest_accelerator": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Description: `List of the type and count of accelerator cards attached to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: `The number of the guest accelerator cards exposed to this instance.`,
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The accelerator type resource exposed to this instance. E.g. nvidia-tesla-k80.`,
						},
					},
				},
			},

			"params": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Stores additional params passed with the request, but not persisted as part of resource payload.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_manager_tags": {
							Type:     schema.TypeMap,
							Optional: true,
							// This field is intentionally not updatable. The API overrides all existing tags on the field when updated.  See go/gce-tags-terraform-support for details.
							ForceNew:    true,
							Description: `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
						},
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A set of key/value label pairs assigned to the instance.`,
			},

			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Metadata key/value pairs made available within the instance.`,
			},

			"metadata_startup_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Metadata startup scripts made available within the instance.`,
			},

			"min_cpu_platform": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      `The minimum CPU platform specified for the VM instance.`,
				DiffSuppressFunc: ComputeInstanceMinCpuPlatformEmptyOrAutomaticDiffSuppress,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If self_link is provided, this value is ignored. If neither self_link nor project are provided, the provider project is used.`,
			},

			"scheduling": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `The scheduling strategy being used by the instance.`,
				Elem: &schema.Resource{
					// !!! IMPORTANT !!!
					// We have a custom diff function for the scheduling block due to issues with Terraform's
					// diff on schema.Set. If changes are made to this block, they must be reflected in that
					// method. See schedulingHasChangeWithoutReboot in compute_instance_helpers.go
					Schema: map[string]*schema.Schema{
						"on_host_maintenance": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Describes maintenance behavior for the instance. One of MIGRATE or TERMINATE,`,
						},

						"automatic_restart": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
							Default:      true,
							Description:  `Specifies if the instance should be restarted if it was terminated by Compute Engine (not a user).`,
						},

						"preemptible": {
							Type:         schema.TypeBool,
							Optional:     true,
							Default:      false,
							AtLeastOneOf: schedulingKeys,
							ForceNew:     true,
							Description:  `Whether the instance is preemptible.`,
						},

						"node_affinities": {
							Type:             schema.TypeSet,
							Optional:         true,
							AtLeastOneOf:     schedulingKeys,
							Elem:             instanceSchedulingNodeAffinitiesElemSchema(),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress(""),
							Description:      `Specifies node affinities or anti-affinities to determine which sole-tenant nodes your instances and managed instance groups will use as host systems.`,
						},

						"min_node_cpus": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
						},

						"provisioning_model": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Whether the instance is spot. If this is set as SPOT.`,
						},

						"instance_termination_action": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Specifies the action GCE should take when SPOT VM is preempted.`,
						},
						"local_ssd_recovery_timeout": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `Specifies the maximum amount of time a Local Ssd Vm should wait while
  recovery of the Local Ssd state is attempted. Its value should be in
  between 0 and 168 hours with hour granularity and the default value being 1
  hour.`,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"seconds": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
										Description: `Span of time at a resolution of a second.
Must be from 0 to 315,576,000,000 inclusive.`,
									},
									"nanos": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Description: `Span of time that's a fraction of a second at nanosecond
resolution. Durations less than one second are represented
with a 0 seconds field and a positive nanos field. Must
be from 0 to 999,999,999 inclusive.`,
									},
								},
							},
						},
					},
				},
			},

			"scratch_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The scratch disks attached to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SCSI", "NVME"}, false),
							Description:  `The disk interface used for attaching this disk. One of SCSI or NVME.`,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(375),
							Default:      375,
							Description:  `The size of the disk in gigabytes. One of 375 or 3000.`,
						},
					},
				},
			},

			"service_account": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: serviceAccountDiffSuppress,
				Description:      `The service account to attach to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The service account e-mail address.`,
						},

						"scopes": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: `A list of service scopes.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(v interface{}) string {
									return tpgresource.CanonicalizeServiceScope(v.(string))
								},
							},
							Set: tpgresource.StringScopeHashcode,
						},
					},
				},
			},

			"shielded_instance_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// Since this block is used by the API based on which
				// image being used, the field needs to be marked as Computed.
				Computed:         true,
				DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress(""),
				Description:      `The shielded vm config being used by the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_secure_boot": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceConfigKeys,
							Default:      false,
							Description:  `Whether secure boot is enabled for the instance.`,
						},

						"enable_vtpm": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceConfigKeys,
							Default:      true,
							Description:  `Whether the instance uses vTPM.`,
						},

						"enable_integrity_monitoring": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceConfigKeys,
							Default:      true,
							Description:  `Whether integrity monitoring is enabled for the instance.`,
						},
					},
				},
			},
			"advanced_machine_features": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `Controls for advanced machine-related behavior features.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_nested_virtualization": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: []string{"advanced_machine_features.0.enable_nested_virtualization", "advanced_machine_features.0.threads_per_core"},
							Description:  `Whether to enable nested virtualization or not.`,
						},
						"threads_per_core": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: []string{"advanced_machine_features.0.enable_nested_virtualization", "advanced_machine_features.0.threads_per_core"},
							Description:  `The number of threads per physical core. To disable simultaneous multithreading (SMT) set this to 1. If unset, the maximum number of threads supported per core by the underlying processor is assumed.`,
						},
						"visible_core_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: []string{"advanced_machine_features.0.enable_nested_virtualization", "advanced_machine_features.0.threads_per_core", "advanced_machine_features.0.visible_core_count"},
							Description:  `The number of physical cores to expose to an instance. Multiply by the number of threads per core to compute the total number of virtual CPUs to expose to the instance. If unset, the number of cores is inferred from the instance\'s nominal CPU count and the underlying platform\'s SMT width.`,
						},
					},
				},
			},
			"confidential_instance_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The Confidential VM config being used by the instance.  on_host_maintenance has to be set to TERMINATE or this will fail to create.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_confidential_compute": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Defines whether the instance should have confidential compute enabled.`,
						},
					},
				},
			},
			"desired_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RUNNING", "TERMINATED"}, false),
				Description:  `Desired status of the instance. Either "RUNNING" or "TERMINATED".`,
			},
			"current_status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `
					Current status of the instance.
					This could be one of the following values: PROVISIONING, STAGING, RUNNING, STOPPING, SUSPENDING, SUSPENDED, REPAIRING, and TERMINATED.
					For more information about the status of the instance, see [Instance life cycle](https://cloud.google.com/compute/docs/instances/instance-life-cycle).`,
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: `The list of tags attached to the instance.`,
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The zone of the instance. If self_link is provided, this value is ignored. If neither self_link nor zone are provided, the provider zone is used.`,
			},

			"cpu_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU platform used by this instance.`,
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The server-assigned unique identifier of this instance.`,
			},

			"label_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique fingerprint of the labels.`,
			},

			"metadata_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique fingerprint of the metadata.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},

			"tags_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique fingerprint of the tags.`,
			},

			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A custom hostname for the instance. Must be a fully qualified DNS name and RFC-1035-valid. Valid format is a series of labels 1-63 characters long matching the regular expression [a-z]([-a-z0-9]*[a-z0-9]), concatenated with periods. The entire hostname must not exceed 253 characters. Changing this forces a new resource to be created.`,
			},

			"resource_policies": {
				Type:             schema.TypeList,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
				Optional:         true,
				MaxItems:         1,
				Description:      `A list of self_links of resource policies to attach to the instance. Currently a max of 1 resource policy is supported.`,
			},

			"reservation_affinity": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the reservations that this instance can consume from.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"ANY_RESERVATION", "SPECIFIC_RESERVATION", "NO_RESERVATION"}, false),
							Description:  `The type of reservation from which this instance can consume resources.`,
						},

						"specific_reservation": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the label selector for the reservation to use.`,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `Corresponds to the label key of a reservation resource. To target a SPECIFIC_RESERVATION by name, specify compute.googleapis.com/reservation-name as the key and specify the name of your reservation as the only value.`,
									},
									"values": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Required:    true,
										ForceNew:    true,
										Description: `Corresponds to the label values of a reservation resource.`,
									},
								},
							},
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.If(
				func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange("guest_accelerator")
				},
				suppressEmptyGuestAcceleratorDiff,
			),
			desiredStatusDiff,
			forceNewIfNetworkIPNotUpdatable,
		),
		UseJSONNumber: true,
	}
}

func getInstance(config *transport_tpg.Config, d *schema.ResourceData) (*compute.Instance, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}
	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return nil, err
	}
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}
	instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, d.Get("name").(string)).Do()
	if err != nil {
		return nil, transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Instance %s", d.Get("name").(string)))
	}
	return instance, nil
}

func getDisk(diskUri string, d *schema.ResourceData, config *transport_tpg.Config) (*compute.Disk, error) {
	source, err := tpgresource.ParseDiskFieldValue(diskUri, d, config)
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	disk, err := config.NewComputeClient(userAgent).Disks.Get(source.Project, source.Zone, source.Name).Do()
	if err != nil {
		return nil, err
	}

	return disk, err
}

func expandComputeInstance(project string, d *schema.ResourceData, config *transport_tpg.Config) (*compute.Instance, error) {
	// Get the machine type
	var machineTypeUrl string
	if mt, ok := d.GetOk("machine_type"); ok {
		machineType, err := tpgresource.ParseMachineTypesFieldValue(mt.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf(
				"Error loading machine type: %s",
				err)
		}
		machineTypeUrl = machineType.RelativeLink()
	}

	// Build up the list of disks

	disks := []*compute.AttachedDisk{}
	if _, hasBootDisk := d.GetOk("boot_disk"); hasBootDisk {
		bootDisk, err := expandBootDisk(d, config, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, bootDisk)
	}

	if _, hasScratchDisk := d.GetOk("scratch_disk"); hasScratchDisk {
		scratchDisks, err := expandScratchDisks(d, config, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, scratchDisks...)
	}

	attachedDisksCount := d.Get("attached_disk.#").(int)

	for i := 0; i < attachedDisksCount; i++ {
		diskConfig := d.Get(fmt.Sprintf("attached_disk.%d", i)).(map[string]interface{})
		disk, err := expandAttachedDisk(diskConfig, d, config)
		if err != nil {
			return nil, err
		}

		disks = append(disks, disk)
	}

	scheduling, err := expandScheduling(d.Get("scheduling"))
	if err != nil {
		return nil, fmt.Errorf("Error creating scheduling: %s", err)
	}

	params, err := expandParams(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating params: %s", err)
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating metadata: %s", err)
	}

	networkInterfaces, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network interfaces: %s", err)
	}
	networkPerformanceConfig, err := expandNetworkPerformanceConfig(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network performance config: %s", err)
	}
	accels, err := expandInstanceGuestAccelerators(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating guest accelerators: %s", err)
	}

	reservationAffinity, err := expandReservationAffinity(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating reservation affinity: %s", err)
	}

	// Create the instance information
	return &compute.Instance{
		CanIpForward:               d.Get("can_ip_forward").(bool),
		Description:                d.Get("description").(string),
		Disks:                      disks,
		MachineType:                machineTypeUrl,
		Metadata:                   metadata,
		Name:                       d.Get("name").(string),
		NetworkInterfaces:          networkInterfaces,
		NetworkPerformanceConfig:   networkPerformanceConfig,
		Tags:                       resourceInstanceTags(d),
		Params:                     params,
		Labels:                     tpgresource.ExpandLabels(d),
		ServiceAccounts:            expandServiceAccounts(d.Get("service_account").([]interface{})),
		GuestAccelerators:          accels,
		MinCpuPlatform:             d.Get("min_cpu_platform").(string),
		Scheduling:                 scheduling,
		DeletionProtection:         d.Get("deletion_protection").(bool),
		Hostname:                   d.Get("hostname").(string),
		ForceSendFields:            []string{"CanIpForward", "DeletionProtection"},
		ConfidentialInstanceConfig: expandConfidentialInstanceConfig(d),
		AdvancedMachineFeatures:    expandAdvancedMachineFeatures(d),
		ShieldedInstanceConfig:     expandShieldedVmConfigs(d),
		DisplayDevice:              expandDisplayDevice(d),
		ResourcePolicies:           tpgresource.ConvertStringArr(d.Get("resource_policies").([]interface{})),
		ReservationAffinity:        reservationAffinity,
	}, nil
}

var computeInstanceStatus = []string{
	"PROVISIONING",
	"REPAIRING",
	"RUNNING",
	"STAGING",
	"STOPPED",
	"STOPPING",
	"SUSPENDED",
	"SUSPENDING",
	"TERMINATED",
}

// return all possible Compute instances status except the one passed as parameter
func getAllStatusBut(status string) []string {
	for i, s := range computeInstanceStatus {
		if status == s {
			return append(computeInstanceStatus[:i], computeInstanceStatus[i+1:]...)
		}
	}
	return computeInstanceStatus
}

func waitUntilInstanceHasDesiredStatus(config *transport_tpg.Config, d *schema.ResourceData) error {
	desiredStatus := d.Get("desired_status").(string)

	if desiredStatus != "" {
		stateRefreshFunc := func() (interface{}, string, error) {
			instance, err := getInstance(config, d)
			if err != nil || instance == nil {
				log.Printf("Error on InstanceStateRefresh: %s", err)
				return nil, "", err
			}
			return instance.Id, instance.Status, nil
		}
		stateChangeConf := resource.StateChangeConf{
			Delay:      5 * time.Second,
			Pending:    getAllStatusBut(desiredStatus),
			Refresh:    stateRefreshFunc,
			Target:     []string{desiredStatus},
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 2 * time.Second,
		}
		_, err := stateChangeConf.WaitForState()

		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance to reach desired status %s: %s", desiredStatus, err)
		}
	}

	return nil
}

func resourceComputeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	z, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Loading zone: %s", z)
	zone, err := config.NewComputeClient(userAgent).Zones.Get(
		project, z).Do()
	if err != nil {
		return fmt.Errorf("Error loading zone '%s': %s", z, err)
	}

	instance, err := expandComputeInstance(project, d, config)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Requesting instance creation")
	op, err := config.NewComputeClient(userAgent).Instances.Insert(project, zone.Name, instance).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, z, instance.Name))

	// Wait for the operation to complete
	waitErr := ComputeOperationWaitTime(config, op, project, "instance to create", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	err = waitUntilInstanceHasDesiredStatus(config, d)
	if err != nil {
		return fmt.Errorf("Error waiting for status: %s", err)
	}

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instance, err := getInstance(config, d)
	if err != nil || instance == nil {
		return err
	}

	md := flattenMetadataBeta(instance.Metadata)

	// If the existing state contains "metadata_startup_script" instead of "metadata.startup-script",
	// we should move the remote metadata.startup-script to metadata_startup_script to avoid
	// specifying it in two places.
	if _, ok := d.GetOk("metadata_startup_script"); ok {
		if err := d.Set("metadata_startup_script", md["startup-script"]); err != nil {
			return fmt.Errorf("Error setting metadata_startup_script: %s", err)
		}

		delete(md, "startup-script")
	}

	if err = d.Set("metadata", md); err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	if err := d.Set("metadata_fingerprint", instance.Metadata.Fingerprint); err != nil {
		return fmt.Errorf("Error setting metadata_fingerprint: %s", err)
	}
	if err := d.Set("can_ip_forward", instance.CanIpForward); err != nil {
		return fmt.Errorf("Error setting can_ip_forward: %s", err)
	}
	if err := d.Set("machine_type", tpgresource.GetResourceNameFromSelfLink(instance.MachineType)); err != nil {
		return fmt.Errorf("Error setting machine_type: %s", err)
	}
	if err := d.Set("network_performance_config", flattenNetworkPerformanceConfig(instance.NetworkPerformanceConfig)); err != nil {
		return err
	}
	// Set the networks
	// Use the first external IP found for the default connection info.
	networkInterfaces, _, internalIP, externalIP, err := flattenNetworkInterfaces(d, config, instance.NetworkInterfaces)
	if err != nil {
		return err
	}
	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	// Fall back on internal ip if there is no external ip.  This makes sense in the situation where
	// terraform is being used on a cloud instance and can therefore access the instances it creates
	// via their internal ips.
	sshIP := externalIP
	if sshIP == "" {
		sshIP = internalIP
	}

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": sshIP,
	})

	// Set the tags fingerprint if there is one.
	if instance.Tags != nil {
		if err := d.Set("tags_fingerprint", instance.Tags.Fingerprint); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
		if err := d.Set("tags", tpgresource.ConvertStringArrToInterface(instance.Tags.Items)); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}

	if err := d.Set("labels", instance.Labels); err != nil {
		return err
	}

	if instance.LabelFingerprint != "" {
		if err := d.Set("label_fingerprint", instance.LabelFingerprint); err != nil {
			return fmt.Errorf("Error setting label_fingerprint: %s", err)
		}
	}

	attachedDiskSources := make(map[string]int)
	for i, v := range d.Get("attached_disk").([]interface{}) {
		if v == nil {
			// There was previously a bug in this code that, when triggered,
			// would cause some nil values to end up in the list of attached disks.
			// Check for this case to make sure we don't try to parse the nil disk.
			continue
		}
		disk := v.(map[string]interface{})
		s := disk["source"].(string)
		var sourceLink string
		if strings.Contains(s, "regions/") {
			source, err := tpgresource.ParseRegionDiskFieldValue(disk["source"].(string), d, config)
			if err != nil {
				return err
			}
			sourceLink = source.RelativeLink()
		} else {
			source, err := tpgresource.ParseDiskFieldValue(disk["source"].(string), d, config)
			if err != nil {
				return err
			}
			sourceLink = source.RelativeLink()
		}
		attachedDiskSources[sourceLink] = i
	}

	attachedDisks := make([]map[string]interface{}, d.Get("attached_disk.#").(int))
	scratchDisks := []map[string]interface{}{}
	for _, disk := range instance.Disks {
		if disk.Boot {
			if err := d.Set("boot_disk", flattenBootDisk(d, disk, config)); err != nil {
				return fmt.Errorf("Error setting boot_disk: %s", err)
			}
		} else if disk.Type == "SCRATCH" {
			scratchDisks = append(scratchDisks, flattenScratchDisk(disk))
		} else {
			var sourceLink string
			if strings.Contains(disk.Source, "regions/") {
				source, err := tpgresource.ParseRegionDiskFieldValue(disk.Source, d, config)
				if err != nil {
					return err
				}
				sourceLink = source.RelativeLink()
			} else {
				source, err := tpgresource.ParseDiskFieldValue(disk.Source, d, config)
				if err != nil {
					return err
				}
				sourceLink = source.RelativeLink()
			}
			adIndex, inConfig := attachedDiskSources[sourceLink]
			di := map[string]interface{}{
				"source":      tpgresource.ConvertSelfLinkToV1(disk.Source),
				"device_name": disk.DeviceName,
				"mode":        disk.Mode,
			}
			if key := disk.DiskEncryptionKey; key != nil {
				if inConfig {
					rawKey := d.Get(fmt.Sprintf("attached_disk.%d.disk_encryption_key_raw", adIndex))
					if rawKey != "" {
						di["disk_encryption_key_raw"] = rawKey
					}
				}
				if key.KmsKeyName != "" {
					// The response for crypto keys often includes the version of the key which needs to be removed
					// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1
					di["kms_key_self_link"] = strings.Split(disk.DiskEncryptionKey.KmsKeyName, "/cryptoKeyVersions")[0]
				}
				if key.Sha256 != "" {
					di["disk_encryption_key_sha256"] = key.Sha256
				}
			}
			// We want the disks to remain in the order we set in the config, so if a disk
			// is present in the config, make sure it's at the correct index. Otherwise, append it.
			if inConfig {
				attachedDisks[adIndex] = di
			} else {
				attachedDisks = append(attachedDisks, di)
			}
		}
	}

	if err := d.Set("resource_policies", instance.ResourcePolicies); err != nil {
		return fmt.Errorf("Error setting resource_policies: %s", err)
	}

	// Remove nils from map in case there were disks in the config that were not present on read;
	// i.e. a disk was detached out of band
	ads := []map[string]interface{}{}
	for _, d := range attachedDisks {
		if d != nil {
			ads = append(ads, d)
		}
	}

	zone := tpgresource.GetResourceNameFromSelfLink(instance.Zone)

	if err := d.Set("service_account", flattenServiceAccounts(instance.ServiceAccounts)); err != nil {
		return fmt.Errorf("Error setting service_account: %s", err)
	}
	if err := d.Set("attached_disk", ads); err != nil {
		return fmt.Errorf("Error setting attached_disk: %s", err)
	}
	if err := d.Set("scratch_disk", scratchDisks); err != nil {
		return fmt.Errorf("Error setting scratch_disk: %s", err)
	}
	if err := d.Set("scheduling", flattenScheduling(instance.Scheduling)); err != nil {
		return fmt.Errorf("Error setting scheduling: %s", err)
	}
	if err := d.Set("guest_accelerator", flattenGuestAccelerators(instance.GuestAccelerators)); err != nil {
		return fmt.Errorf("Error setting guest_accelerator: %s", err)
	}
	if err := d.Set("shielded_instance_config", flattenShieldedVmConfig(instance.ShieldedInstanceConfig)); err != nil {
		return fmt.Errorf("Error setting shielded_instance_config: %s", err)
	}
	if err := d.Set("enable_display", flattenEnableDisplay(instance.DisplayDevice)); err != nil {
		return fmt.Errorf("Error setting enable_display: %s", err)
	}
	if err := d.Set("cpu_platform", instance.CpuPlatform); err != nil {
		return fmt.Errorf("Error setting cpu_platform: %s", err)
	}
	if err := d.Set("min_cpu_platform", instance.MinCpuPlatform); err != nil {
		return fmt.Errorf("Error setting min_cpu_platform: %s", err)
	}
	if err := d.Set("deletion_protection", instance.DeletionProtection); err != nil {
		return fmt.Errorf("Error setting deletion_protection: %s", err)
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(instance.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("instance_id", fmt.Sprintf("%d", instance.Id)); err != nil {
		return fmt.Errorf("Error setting instance_id: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("name", instance.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", instance.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("hostname", instance.Hostname); err != nil {
		return fmt.Errorf("Error setting hostname: %s", err)
	}
	if err := d.Set("current_status", instance.Status); err != nil {
		return fmt.Errorf("Error setting current_status: %s", err)
	}
	if err := d.Set("confidential_instance_config", flattenConfidentialInstanceConfig(instance.ConfidentialInstanceConfig)); err != nil {
		return fmt.Errorf("Error setting confidential_instance_config: %s", err)
	}
	if err := d.Set("advanced_machine_features", flattenAdvancedMachineFeatures(instance.AdvancedMachineFeatures)); err != nil {
		return fmt.Errorf("Error setting advanced_machine_features: %s", err)
	}
	if d.Get("desired_status") != "" {
		if err := d.Set("desired_status", instance.Status); err != nil {
			return fmt.Errorf("Error setting desired_status: %s", err)
		}
	}
	if err := d.Set("reservation_affinity", flattenReservationAffinity(instance.ReservationAffinity)); err != nil {
		return fmt.Errorf("Error setting reservation_affinity: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, instance.Name))

	return nil
}

func resourceComputeInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}

	// Use beta api directly in order to read network_interface.fingerprint without having to put it in the schema.
	// Change back to getInstance(config, d) once updating alias ips is GA.
	instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, d.Get("name").(string)).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Instance %s", d.Get("name").(string)))
	}

	// Enable partial mode for the resource since it is possible
	d.Partial(true)

	if d.HasChange("metadata") {
		metadata, err := resourceInstanceMetadata(d)
		if err != nil {
			return fmt.Errorf("Error parsing metadata: %s", err)
		}

		metadataV1 := &compute.Metadata{}
		if err := tpgresource.Convert(metadata, metadataV1); err != nil {
			return err
		}

		// We're retrying for an error 412 where the metadata fingerprint is out of date
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				// retrieve up-to-date metadata from the API in case several updates hit simultaneously. instances
				// sometimes but not always share metadata fingerprints.
				instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, instance.Name).Do()
				if err != nil {
					return fmt.Errorf("Error retrieving metadata: %s", err)
				}

				metadataV1.Fingerprint = instance.Metadata.Fingerprint

				op, err := config.NewComputeClient(userAgent).Instances.SetMetadata(project, zone, instance.Name, metadataV1).Do()
				if err != nil {
					return fmt.Errorf("Error updating metadata: %s", err)
				}

				opErr := ComputeOperationWaitTime(config, op, project, "metadata to update", userAgent, d.Timeout(schema.TimeoutUpdate))
				if opErr != nil {
					return opErr
				}

				return nil
			},
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		tags := resourceInstanceTags(d)
		tagsV1 := &compute.Tags{}
		if err := tpgresource.Convert(tags, tagsV1); err != nil {
			return err
		}
		op, err := config.NewComputeClient(userAgent).Instances.SetTags(
			project, zone, d.Get("name").(string), tagsV1).Do()
		if err != nil {
			return fmt.Errorf("Error updating tags: %s", err)
		}

		opErr := ComputeOperationWaitTime(config, op, project, "tags to update", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}

	if d.HasChange("labels") {
		labels := tpgresource.ExpandLabels(d)
		labelFingerprint := d.Get("label_fingerprint").(string)
		req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: labelFingerprint}

		op, err := config.NewComputeClient(userAgent).Instances.SetLabels(project, zone, instance.Name, &req).Do()
		if err != nil {
			return fmt.Errorf("Error updating labels: %s", err)
		}

		opErr := ComputeOperationWaitTime(config, op, project, "labels to update", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}

	if d.HasChange("resource_policies") {
		if len(instance.ResourcePolicies) > 0 {
			req := compute.InstancesRemoveResourcePoliciesRequest{ResourcePolicies: instance.ResourcePolicies}

			op, err := config.NewComputeClient(userAgent).Instances.RemoveResourcePolicies(project, zone, instance.Name, &req).Do()
			if err != nil {
				return fmt.Errorf("Error removing existing resource policies: %s", err)
			}

			opErr := ComputeOperationWaitTime(config, op, project, "resource policies to remove", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		resourcePolicies := tpgresource.ConvertStringArr(d.Get("resource_policies").([]interface{}))
		if len(resourcePolicies) > 0 {
			req := compute.InstancesAddResourcePoliciesRequest{ResourcePolicies: resourcePolicies}

			op, err := config.NewComputeClient(userAgent).Instances.AddResourcePolicies(project, zone, instance.Name, &req).Do()
			if err != nil {
				return fmt.Errorf("Error adding resource policies: %s", err)
			}

			opErr := ComputeOperationWaitTime(config, op, project, "resource policies to add", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}
	}

	bootRequiredSchedulingChange := schedulingHasChangeRequiringReboot(d)
	bootNotRequiredSchedulingChange := schedulingHasChangeWithoutReboot(d)
	if bootNotRequiredSchedulingChange {
		scheduling, err := expandScheduling(d.Get("scheduling"))
		if err != nil {
			return fmt.Errorf("Error creating request data to update scheduling: %s", err)
		}

		op, err := config.NewComputeClient(userAgent).Instances.SetScheduling(
			project, zone, instance.Name, scheduling).Do()
		if err != nil {
			return fmt.Errorf("Error updating scheduling policy: %s", err)
		}

		opErr := ComputeOperationWaitTime(
			config, op, project, "scheduling policy update", userAgent,
			d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}

	networkInterfaces, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return fmt.Errorf("Error getting network interface from config: %s", err)
	}

	// Sanity check
	if len(networkInterfaces) != len(instance.NetworkInterfaces) {
		return fmt.Errorf("Instance had unexpected number of network interfaces: %d", len(instance.NetworkInterfaces))
	}

	var updatesToNIWhileStopped []func(inst *compute.Instance) error
	for i := 0; i < len(networkInterfaces); i++ {
		prefix := fmt.Sprintf("network_interface.%d", i)
		networkInterface := networkInterfaces[i]
		instNetworkInterface := instance.NetworkInterfaces[i]

		networkName := d.Get(prefix + ".name").(string)
		subnetwork := networkInterface.Subnetwork
		updateDuringStop := d.HasChange(prefix+".subnetwork") || d.HasChange(prefix+".network") || d.HasChange(prefix+".subnetwork_project")

		// Sanity check
		if networkName != instNetworkInterface.Name {
			return fmt.Errorf("Instance networkInterface had unexpected name: %s", instNetworkInterface.Name)
		}

		// On creation the network is inferred if only subnetwork is given.
		// Unforunately for us there is no way to determine if the user is
		// explicitly assigning network or we are reusing the one that was inferred
		// from state. So here we check if subnetwork changed and network did not.
		// In this scenario we assume network was inferred and attempt to figure out
		// the new corresponding network.

		if d.HasChange(prefix + ".subnetwork") {
			if !d.HasChange(prefix + ".network") {
				subnetProjectField := prefix + ".subnetwork_project"
				sf, err := tpgresource.ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
				if err != nil {
					return fmt.Errorf("Cannot determine self_link for subnetwork %q: %s", subnetwork, err)
				}
				resp, err := config.NewComputeClient(userAgent).Subnetworks.Get(sf.Project, sf.Region, sf.Name).Do()
				if err != nil {
					return errwrap.Wrapf("Error getting subnetwork value: {{err}}", err)
				}
				nf, err := tpgresource.ParseNetworkFieldValue(resp.Network, d, config)
				if err != nil {
					return fmt.Errorf("Cannot determine self_link for network %q: %s", resp.Network, err)
				}
				networkInterface.Network = nf.RelativeLink()
			}
		}

		if !updateDuringStop && d.HasChange(prefix+".access_config") {
			// TODO: This code deletes then recreates accessConfigs.  This is bad because it may
			// leave the machine inaccessible from either ip if the creation part fails (network
			// timeout etc).  However right now there is a GCE limit of 1 accessConfig so it is
			// the only way to do it.  In future this should be revised to only change what is
			// necessary, and also add before removing.

			// Delete current access configs
			err := computeInstanceDeleteAccessConfigs(d, config, instNetworkInterface, project, zone, userAgent, instance.Name)
			if err != nil {
				return err
			}

			// Create new ones
			err = computeInstanceAddAccessConfigs(d, config, instNetworkInterface, networkInterface.AccessConfigs, project, zone, userAgent, instance.Name)
			if err != nil {
				return err
			}

			// re-read fingerprint
			instance, err = config.NewComputeClient(userAgent).Instances.Get(project, zone, instance.Name).Do()
			if err != nil {
				return err
			}
			instNetworkInterface = instance.NetworkInterfaces[i]
		}

		if !updateDuringStop && d.HasChange(prefix+".alias_ip_range") {
			// Alias IP ranges cannot be updated; they must be removed and then added
			// unless you are changing subnetwork/network
			if len(instNetworkInterface.AliasIpRanges) > 0 {
				ni := &compute.NetworkInterface{
					Fingerprint:     instNetworkInterface.Fingerprint,
					ForceSendFields: []string{"AliasIpRanges"},
				}
				op, err := config.NewComputeClient(userAgent).Instances.UpdateNetworkInterface(project, zone, instance.Name, networkName, ni).Do()
				if err != nil {
					return errwrap.Wrapf("Error removing alias_ip_range: {{err}}", err)
				}
				opErr := ComputeOperationWaitTime(config, op, project, "updating alias ip ranges", userAgent, d.Timeout(schema.TimeoutUpdate))
				if opErr != nil {
					return opErr
				}
				// re-read fingerprint
				instance, err = config.NewComputeClient(userAgent).Instances.Get(project, zone, instance.Name).Do()
				if err != nil {
					return err
				}
				instNetworkInterface = instance.NetworkInterfaces[i]
			}

			networkInterfacePatchObj := &compute.NetworkInterface{
				AliasIpRanges: networkInterface.AliasIpRanges,
				Fingerprint:   instNetworkInterface.Fingerprint,
			}
			updateCall := config.NewComputeClient(userAgent).Instances.UpdateNetworkInterface(project, zone, instance.Name, networkName, networkInterfacePatchObj).Do
			op, err := updateCall()
			if err != nil {
				return errwrap.Wrapf("Error updating network interface: {{err}}", err)
			}
			opErr := ComputeOperationWaitTime(config, op, project, "network interface to update", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if updateDuringStop {
			// Lets be explicit about what we are changing in the patch call
			networkInterfacePatchObj := &compute.NetworkInterface{
				Network:       networkInterface.Network,
				Subnetwork:    networkInterface.Subnetwork,
				AliasIpRanges: networkInterface.AliasIpRanges,
			}

			// network_ip can be inferred if not declared. Let's only patch if it's being changed by user
			// otherwise this could fail if the network ip is not compatible with the new Subnetwork/Network.
			if d.HasChange(prefix + ".network_ip") {
				networkInterfacePatchObj.NetworkIP = networkInterface.NetworkIP
			}

			// Access config can run into some issues since we can't tell the difference between
			// the users declared intent (config within their hcl file) and what we have inferred from the
			// server (terraform state). Access configs contain an ip subproperty that can be incompatible
			// with the subnetwork/network we are transitioning to. Due to this we only change access
			// configs if we notice the configuration (user intent) changes.
			accessConfigsHaveChanged := d.HasChange(prefix + ".access_config")

			updateCall := computeInstanceCreateUpdateWhileStoppedCall(d, config, networkInterfacePatchObj, networkInterface.AccessConfigs, accessConfigsHaveChanged, i, project, zone, userAgent, instance.Name)
			updatesToNIWhileStopped = append(updatesToNIWhileStopped, updateCall)
		}
	}

	if d.HasChange("attached_disk") {
		o, n := d.GetChange("attached_disk")

		// Keep track of disks currently in the instance. Because the google_compute_disk resource
		// can detach disks, it's possible that there are fewer disks currently attached than there
		// were at the time we ran terraform plan.
		currDisks := map[string]struct{}{}
		for _, disk := range instance.Disks {
			if !disk.Boot && disk.Type != "SCRATCH" {
				currDisks[disk.DeviceName] = struct{}{}
			}
		}

		// Keep track of disks currently in state.
		// Since changing any field within the disk needs to detach+reattach it,
		// keep track of the hash of the full disk.
		oDisks := map[uint64]string{}
		for _, disk := range o.([]interface{}) {
			diskConfig := disk.(map[string]interface{})
			computeDisk, err := expandAttachedDisk(diskConfig, d, config)
			if err != nil {
				return err
			}
			hash, err := hashstructure.Hash(*computeDisk, nil)
			if err != nil {
				return err
			}
			if _, ok := currDisks[computeDisk.DeviceName]; ok {
				oDisks[hash] = computeDisk.DeviceName
			}
		}

		// Keep track of new config's disks.
		// Since changing any field within the disk needs to detach+reattach it,
		// keep track of the hash of the full disk.
		// If a disk with a certain hash is only in the new config, it should be attached.
		nDisks := map[uint64]struct{}{}
		var attach []*compute.AttachedDisk
		for _, disk := range n.([]interface{}) {
			diskConfig := disk.(map[string]interface{})
			computeDisk, err := expandAttachedDisk(diskConfig, d, config)
			if err != nil {
				return err
			}
			hash, err := hashstructure.Hash(*computeDisk, nil)
			if err != nil {
				return err
			}
			nDisks[hash] = struct{}{}

			if _, ok := oDisks[hash]; !ok {
				computeDiskV1 := &compute.AttachedDisk{}
				err = tpgresource.Convert(computeDisk, computeDiskV1)
				if err != nil {
					return err
				}
				attach = append(attach, computeDiskV1)
			}
		}

		// If a source is only in the old config, it should be detached.
		// Detach the old disks.
		for hash, deviceName := range oDisks {
			if _, ok := nDisks[hash]; !ok {
				op, err := config.NewComputeClient(userAgent).Instances.DetachDisk(project, zone, instance.Name, deviceName).Do()
				if err != nil {
					return errwrap.Wrapf("Error detaching disk: %s", err)
				}

				opErr := ComputeOperationWaitTime(config, op, project, "detaching disk", userAgent, d.Timeout(schema.TimeoutUpdate))
				if opErr != nil {
					return opErr
				}
				log.Printf("[DEBUG] Successfully detached disk %s", deviceName)
			}
		}

		// Attach the new disks
		for _, disk := range attach {
			op, err := config.NewComputeClient(userAgent).Instances.AttachDisk(project, zone, instance.Name, disk).Do()
			if err != nil {
				return errwrap.Wrapf("Error attaching disk : {{err}}", err)
			}

			opErr := ComputeOperationWaitTime(config, op, project, "attaching disk", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
			log.Printf("[DEBUG] Successfully attached disk %s", disk.Source)
		}
	}

	// d.HasChange("service_account") is oversensitive: see https://github.com/hashicorp/terraform/issues/17411
	// Until that's fixed, manually check whether there is a change.
	o, n := d.GetChange("service_account")
	oList := o.([]interface{})
	nList := n.([]interface{})
	scopesChange := false
	if len(oList) != len(nList) {
		scopesChange = true
	} else if len(oList) == 1 {
		// service_account has MaxItems: 1
		// scopes is a required field and so will always be set
		oScopes := oList[0].(map[string]interface{})["scopes"].(*schema.Set)
		nScopes := nList[0].(map[string]interface{})["scopes"].(*schema.Set)
		scopesChange = !oScopes.Equal(nScopes)
	}

	if d.HasChange("deletion_protection") {
		nDeletionProtection := d.Get("deletion_protection").(bool)

		op, err := config.NewComputeClient(userAgent).Instances.SetDeletionProtection(project, zone, d.Get("name").(string)).DeletionProtection(nDeletionProtection).Do()
		if err != nil {
			return fmt.Errorf("Error updating deletion protection flag: %s", err)
		}

		opErr := ComputeOperationWaitTime(config, op, project, "deletion protection to update", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}

	if d.HasChange("can_ip_forward") {
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, instance.Name).Do()
				if err != nil {
					return fmt.Errorf("Error retrieving instance: %s", err)
				}

				instance.CanIpForward = d.Get("can_ip_forward").(bool)

				op, err := config.NewComputeClient(userAgent).Instances.Update(project, zone, instance.Name, instance).Do()
				if err != nil {
					return fmt.Errorf("Error updating instance: %s", err)
				}

				opErr := ComputeOperationWaitTime(config, op, project, "can_ip_forward, updating", userAgent, d.Timeout(schema.TimeoutUpdate))
				if opErr != nil {
					return opErr
				}

				return nil
			},
		})

		if err != nil {
			return err
		}
	}

	needToStopInstanceBeforeUpdating := scopesChange || d.HasChange("service_account.0.email") || d.HasChange("machine_type") || d.HasChange("min_cpu_platform") || d.HasChange("enable_display") || d.HasChange("shielded_instance_config") || len(updatesToNIWhileStopped) > 0 || bootRequiredSchedulingChange || d.HasChange("advanced_machine_features")

	if d.HasChange("desired_status") && !needToStopInstanceBeforeUpdating {
		desiredStatus := d.Get("desired_status").(string)

		if desiredStatus != "" {
			var op *compute.Operation

			if desiredStatus == "RUNNING" {
				op, err = startInstanceOperation(d, config)
				if err != nil {
					return errwrap.Wrapf("Error starting instance: {{err}}", err)
				}
			} else if desiredStatus == "TERMINATED" {
				op, err = config.NewComputeClient(userAgent).Instances.Stop(project, zone, instance.Name).Do()
				if err != nil {
					return err
				}
			}
			opErr := ComputeOperationWaitTime(
				config, op, project, "updating status", userAgent,
				d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}
	}

	// Attributes which can only be changed if the instance is stopped
	if needToStopInstanceBeforeUpdating {
		statusBeforeUpdate := instance.Status
		desiredStatus := d.Get("desired_status").(string)

		if statusBeforeUpdate == "RUNNING" && desiredStatus != "TERMINATED" && !d.Get("allow_stopping_for_update").(bool) {
			return fmt.Errorf("Changing the machine_type, min_cpu_platform, service_account, enable_display, shielded_instance_config, scheduling.node_affinities " +
				"or network_interface.[#d].(network/subnetwork/subnetwork_project) or advanced_machine_features on a started instance requires stopping it. " +
				"To acknowledge this, please set allow_stopping_for_update = true in your config. " +
				"You can also stop it by setting desired_status = \"TERMINATED\", but the instance will not be restarted after the update.")
		}

		if statusBeforeUpdate != "TERMINATED" {
			op, err := config.NewComputeClient(userAgent).Instances.Stop(project, zone, instance.Name).Do()
			if err != nil {
				return errwrap.Wrapf("Error stopping instance: {{err}}", err)
			}

			opErr := ComputeOperationWaitTime(config, op, project, "stopping instance", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if d.HasChange("min_cpu_platform") {
			minCpuPlatform := d.Get("min_cpu_platform")
			req := &compute.InstancesSetMinCpuPlatformRequest{
				MinCpuPlatform: minCpuPlatform.(string),
			}
			op, err := config.NewComputeClient(userAgent).Instances.SetMinCpuPlatform(project, zone, instance.Name, req).Do()
			if err != nil {
				return err
			}
			opErr := ComputeOperationWaitTime(config, op, project, "updating min cpu platform", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if d.HasChange("machine_type") {
			mt, err := tpgresource.ParseMachineTypesFieldValue(d.Get("machine_type").(string), d, config)
			if err != nil {
				return err
			}
			req := &compute.InstancesSetMachineTypeRequest{
				MachineType: mt.RelativeLink(),
			}
			op, err := config.NewComputeClient(userAgent).Instances.SetMachineType(project, zone, instance.Name, req).Do()
			if err != nil {
				return err
			}
			opErr := ComputeOperationWaitTime(config, op, project, "updating machinetype", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if d.HasChange("service_account.0.email") || scopesChange {
			sa := d.Get("service_account").([]interface{})
			req := &compute.InstancesSetServiceAccountRequest{ForceSendFields: []string{"email"}}
			if len(sa) > 0 && sa[0] != nil {
				saMap := sa[0].(map[string]interface{})
				req.Email = saMap["email"].(string)
				req.Scopes = tpgresource.CanonicalizeServiceScopes(tpgresource.ConvertStringSet(saMap["scopes"].(*schema.Set)))
			}
			op, err := config.NewComputeClient(userAgent).Instances.SetServiceAccount(project, zone, instance.Name, req).Do()
			if err != nil {
				return err
			}
			opErr := ComputeOperationWaitTime(config, op, project, "updating service account", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if d.HasChange("enable_display") {
			req := &compute.DisplayDevice{
				EnableDisplay:   d.Get("enable_display").(bool),
				ForceSendFields: []string{"EnableDisplay"},
			}
			op, err := config.NewComputeClient(userAgent).Instances.UpdateDisplayDevice(project, zone, instance.Name, req).Do()
			if err != nil {
				return fmt.Errorf("Error updating display device: %s", err)
			}
			opErr := ComputeOperationWaitTime(config, op, project, "updating display device", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if d.HasChange("shielded_instance_config") {
			shieldedVmConfig := expandShieldedVmConfigs(d)

			op, err := config.NewComputeClient(userAgent).Instances.UpdateShieldedInstanceConfig(project, zone, instance.Name, shieldedVmConfig).Do()
			if err != nil {
				return fmt.Errorf("Error updating shielded vm config: %s", err)
			}

			opErr := ComputeOperationWaitTime(config, op, project,
				"shielded vm config update", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if bootRequiredSchedulingChange {
			scheduling, err := expandScheduling(d.Get("scheduling"))
			if err != nil {
				return fmt.Errorf("Error creating request data to update scheduling: %s", err)
			}

			op, err := config.NewComputeClient(userAgent).Instances.SetScheduling(
				project, zone, instance.Name, scheduling).Do()
			if err != nil {
				return fmt.Errorf("Error updating scheduling policy: %s", err)
			}

			opErr := ComputeOperationWaitTime(
				config, op, project, "scheduling policy update", userAgent,
				d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}

		if d.HasChange("advanced_machine_features") {
			err = transport_tpg.Retry(transport_tpg.RetryOptions{
				RetryFunc: func() error {
					// retrieve up-to-date instance from the API in case several updates hit simultaneously. instances
					// sometimes but not always share metadata fingerprints.
					instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, instance.Name).Do()
					if err != nil {
						return fmt.Errorf("Error retrieving instance: %s", err)
					}

					instance.AdvancedMachineFeatures = expandAdvancedMachineFeatures(d)

					op, err := config.NewComputeClient(userAgent).Instances.Update(project, zone, instance.Name, instance).Do()
					if err != nil {
						return fmt.Errorf("Error updating instance: %s", err)
					}

					opErr := ComputeOperationWaitTime(config, op, project, "advanced_machine_features to update", userAgent, d.Timeout(schema.TimeoutUpdate))
					if opErr != nil {
						return opErr
					}

					return nil
				},
			})

			if err != nil {
				return err
			}
		}

		// If the instance stops it can invalidate the fingerprint for network interface.
		// refresh the instance to get a new fingerprint
		if len(updatesToNIWhileStopped) > 0 {
			instance, err = config.NewComputeClient(userAgent).Instances.Get(project, zone, instance.Name).Do()
			if err != nil {
				return err
			}
		}
		for _, patch := range updatesToNIWhileStopped {
			err := patch(instance)
			if err != nil {
				return err
			}
		}

		if (statusBeforeUpdate == "RUNNING" && desiredStatus != "TERMINATED") ||
			(statusBeforeUpdate == "TERMINATED" && desiredStatus == "RUNNING") {
			op, err := startInstanceOperation(d, config)
			if err != nil {
				return errwrap.Wrapf("Error starting instance: {{err}}", err)
			}

			opErr := ComputeOperationWaitTime(config, op, project,
				"starting instance", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}
	}

	// We made it, disable partial mode
	d.Partial(false)

	return resourceComputeInstanceRead(d, meta)
}

func startInstanceOperation(d *schema.ResourceData, config *transport_tpg.Config) (*compute.Operation, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	// Use beta api directly in order to read network_interface.fingerprint without having to put it in the schema.
	// Change back to getInstance(config, d) once updating alias ips is GA.
	instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, d.Get("name").(string)).Do()
	if err != nil {
		return nil, transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Instance %s", instance.Name))
	}

	// Retrieve instance from config to pull encryption keys if necessary
	instanceFromConfig, err := expandComputeInstance(project, d, config)
	if err != nil {
		return nil, err
	}

	var encrypted []*compute.CustomerEncryptionKeyProtectedDisk
	for _, disk := range instanceFromConfig.Disks {
		if disk.DiskEncryptionKey != nil {
			key := compute.CustomerEncryptionKey{RawKey: disk.DiskEncryptionKey.RawKey, KmsKeyName: disk.DiskEncryptionKey.KmsKeyName}
			eDisk := compute.CustomerEncryptionKeyProtectedDisk{Source: disk.Source, DiskEncryptionKey: &key}
			encrypted = append(encrypted, &eDisk)
		}
	}

	var op *compute.Operation

	if len(encrypted) > 0 {
		request := compute.InstancesStartWithEncryptionKeyRequest{Disks: encrypted}
		op, err = config.NewComputeClient(userAgent).Instances.StartWithEncryptionKey(project, zone, instance.Name, &request).Do()
	} else {
		op, err = config.NewComputeClient(userAgent).Instances.Start(project, zone, instance.Name).Do()
	}

	return op, err
}

func expandAttachedDisk(diskConfig map[string]interface{}, d *schema.ResourceData, meta interface{}) (*compute.AttachedDisk, error) {
	config := meta.(*transport_tpg.Config)

	s := diskConfig["source"].(string)
	var sourceLink string
	if strings.Contains(s, "regions/") {
		source, err := tpgresource.ParseRegionDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	} else {
		source, err := tpgresource.ParseDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	}

	disk := &compute.AttachedDisk{
		Source: sourceLink,
	}

	if v, ok := diskConfig["mode"]; ok {
		disk.Mode = v.(string)
	}

	if v, ok := diskConfig["device_name"]; ok {
		disk.DeviceName = v.(string)
	}

	keyValue, keyOk := diskConfig["disk_encryption_key_raw"]
	if keyOk {
		if keyValue != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey: keyValue.(string),
			}
		}
	}

	kmsValue, kmsOk := diskConfig["kms_key_self_link"]
	if kmsOk {
		if keyOk && keyValue != "" && kmsValue != "" {
			return nil, errors.New("Only one of kms_key_self_link and disk_encryption_key_raw can be set")
		}
		if kmsValue != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				KmsKeyName: kmsValue.(string),
			}
		}
	}
	return disk, nil
}

// See comment on expandInstanceTemplateGuestAccelerators regarding why this
// code is duplicated.
func expandInstanceGuestAccelerators(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]*compute.AcceleratorConfig, error) {
	configs, ok := d.GetOk("guest_accelerator")
	if !ok {
		return nil, nil
	}
	accels := configs.([]interface{})
	guestAccelerators := make([]*compute.AcceleratorConfig, 0, len(accels))
	for _, raw := range accels {
		data := raw.(map[string]interface{})
		if data["count"].(int) == 0 {
			continue
		}
		at, err := tpgresource.ParseAcceleratorFieldValue(data["type"].(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot parse accelerator type: %v", err)
		}
		guestAccelerators = append(guestAccelerators, &compute.AcceleratorConfig{
			AcceleratorCount: int64(data["count"].(int)),
			AcceleratorType:  at.RelativeLink(),
		})
	}

	return guestAccelerators, nil
}

// suppressEmptyGuestAcceleratorDiff is used to work around perpetual diff
// issues when a count of `0` guest accelerators is desired. This may occur when
// guest_accelerator support is controlled via a module variable. E.g.:
//
//	guest_accelerators {
// 		count = "${var.enable_gpu ? var.gpu_count : 0}"
// 		...
//	}

// After reconciling the desired and actual state, we would otherwise see a
// perpetual diff resembling:
//
//	[] != [{"count":0, "type": "nvidia-tesla-k80"}]
func suppressEmptyGuestAcceleratorDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	oldi, newi := d.GetChange("guest_accelerator")

	old, ok := oldi.([]interface{})
	if !ok {
		return fmt.Errorf("Expected old guest accelerator diff to be a slice")
	}

	new, ok := newi.([]interface{})
	if !ok {
		return fmt.Errorf("Expected new guest accelerator diff to be a slice")
	}

	if len(old) != 0 && len(new) != 1 {
		return nil
	}

	firstAccel, ok := new[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to type assert guest accelerator")
	}

	if firstAccel["count"].(int) == 0 {
		if err := d.Clear("guest_accelerator"); err != nil {
			return err
		}
	}

	return nil
}

// return an error if the desired_status field is set to a value other than RUNNING on Create.
func desiredStatusDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	// when creating an instance, name is not set
	oldName, _ := diff.GetChange("name")

	if oldName == nil || oldName == "" {
		_, newDesiredStatus := diff.GetChange("desired_status")

		if newDesiredStatus == nil || newDesiredStatus == "" {
			return nil
		} else if newDesiredStatus != "RUNNING" {
			return fmt.Errorf("When creating an instance, desired_status can only accept RUNNING value")
		}
		return nil
	}

	return nil
}

func resourceComputeInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Requesting instance deletion: %s", d.Get("name").(string))

	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("Cannot delete instance %s: instance Deletion Protection is enabled. Set deletion_protection to false for this resource and run \"terraform apply\" before attempting to delete it.", d.Get("name").(string))
	} else {
		op, err := config.NewComputeClient(userAgent).Instances.Delete(project, zone, d.Get("name").(string)).Do()
		if err != nil {
			return fmt.Errorf("Error deleting instance: %s", err)
		}

		// Wait for the operation to complete
		opErr := ComputeOperationWaitTime(config, op, project, "instance to delete", userAgent, d.Timeout(schema.TimeoutDelete))
		if opErr != nil {
			// Refresh operation to check status
			op, _ = config.NewComputeClient(userAgent).ZoneOperations.Get(project, zone, strconv.FormatUint(op.Id, 10)).Do()
			// Do not return an error if the operation actually completed
			if op == nil || op.Status != "DONE" {
				return opErr
			}
		}

		d.SetId("")
		return nil
	}
}

func resourceComputeInstanceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/zones/{{zone}}/instances/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandParams(d *schema.ResourceData) (*compute.InstanceParams, error) {
	params := &compute.InstanceParams{}

	if _, ok := d.GetOk("params.0.resource_manager_tags"); ok {
		params.ResourceManagerTags = tpgresource.ExpandStringMap(d, "params.0.resource_manager_tags")
	}

	return params, nil
}

func expandBootDisk(d *schema.ResourceData, config *transport_tpg.Config, project string) (*compute.AttachedDisk, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	disk := &compute.AttachedDisk{
		AutoDelete: d.Get("boot_disk.0.auto_delete").(bool),
		Boot:       true,
	}

	if v, ok := d.GetOk("boot_disk.0.device_name"); ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_raw"); ok {
		if v != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey: v.(string),
			}
		}
	}

	if v, ok := d.GetOk("boot_disk.0.kms_key_self_link"); ok {
		if v != "" {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				KmsKeyName: v.(string),
			}
		}
	}

	if v, ok := d.GetOk("boot_disk.0.source"); ok {
		source, err := tpgresource.ParseDiskFieldValue(v.(string), d, config)
		if err != nil {
			return nil, err
		}
		disk.Source = source.RelativeLink()
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params"); ok {
		disk.InitializeParams = &compute.AttachedDiskInitializeParams{}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.size"); ok {
			disk.InitializeParams.DiskSizeGb = int64(v.(int))
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.type"); ok {
			diskTypeName := v.(string)
			diskType, err := readDiskType(config, d, diskTypeName)
			if err != nil {
				return nil, fmt.Errorf("Error loading disk type '%s': %s", diskTypeName, err)
			}
			disk.InitializeParams.DiskType = diskType.RelativeLink()
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.image"); ok {
			imageName := v.(string)
			imageUrl, err := ResolveImage(config, project, imageName, userAgent)
			if err != nil {
				return nil, fmt.Errorf("Error resolving image name '%s': %s", imageName, err)
			}

			disk.InitializeParams.SourceImage = imageUrl
		}

		if _, ok := d.GetOk("boot_disk.0.initialize_params.0.labels"); ok {
			disk.InitializeParams.Labels = tpgresource.ExpandStringMap(d, "boot_disk.0.initialize_params.0.labels")
		}

		if _, ok := d.GetOk("boot_disk.0.initialize_params.0.resource_manager_tags"); ok {
			disk.InitializeParams.ResourceManagerTags = tpgresource.ExpandStringMap(d, "boot_disk.0.initialize_params.0.resource_manager_tags")
		}
	}

	if v, ok := d.GetOk("boot_disk.0.mode"); ok {
		disk.Mode = v.(string)
	}

	return disk, nil
}

func flattenBootDisk(d *schema.ResourceData, disk *compute.AttachedDisk, config *transport_tpg.Config) []map[string]interface{} {
	result := map[string]interface{}{
		"auto_delete": disk.AutoDelete,
		"device_name": disk.DeviceName,
		"mode":        disk.Mode,
		"source":      tpgresource.ConvertSelfLinkToV1(disk.Source),
		// disk_encryption_key_raw is not returned from the API, so copy it from what the user
		// originally specified to avoid diffs.
		"disk_encryption_key_raw": d.Get("boot_disk.0.disk_encryption_key_raw"),
	}

	diskDetails, err := getDisk(disk.Source, d, config)
	if err != nil {
		log.Printf("[WARN] Cannot retrieve boot disk details: %s", err)

		if _, ok := d.GetOk("boot_disk.0.initialize_params.#"); ok {
			// If we can't read the disk details due to permission for instance,
			// copy the initialize_params from what the user originally specified to avoid diffs.
			m := d.Get("boot_disk.0.initialize_params")
			result["initialize_params"] = m
		}
	} else {
		result["initialize_params"] = []map[string]interface{}{{
			"type": tpgresource.GetResourceNameFromSelfLink(diskDetails.Type),
			// If the config specifies a family name that doesn't match the image name, then
			// the diff won't be properly suppressed. See DiffSuppressFunc for this field.
			"image":                 diskDetails.SourceImage,
			"size":                  diskDetails.SizeGb,
			"labels":                diskDetails.Labels,
			"resource_manager_tags": d.Get("boot_disk.0.initialize_params.0.resource_manager_tags"),
		}}
	}

	if disk.DiskEncryptionKey != nil {
		if disk.DiskEncryptionKey.Sha256 != "" {
			result["disk_encryption_key_sha256"] = disk.DiskEncryptionKey.Sha256
		}
		if disk.DiskEncryptionKey.KmsKeyName != "" {
			// The response for crypto keys often includes the version of the key which needs to be removed
			// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1
			result["kms_key_self_link"] = strings.Split(disk.DiskEncryptionKey.KmsKeyName, "/cryptoKeyVersions")[0]
		}
	}

	return []map[string]interface{}{result}
}

func expandScratchDisks(d *schema.ResourceData, config *transport_tpg.Config, project string) ([]*compute.AttachedDisk, error) {
	diskType, err := readDiskType(config, d, "local-ssd")
	if err != nil {
		return nil, fmt.Errorf("Error loading disk type 'local-ssd': %s", err)
	}

	n := d.Get("scratch_disk.#").(int)
	scratchDisks := make([]*compute.AttachedDisk, 0, n)
	for i := 0; i < n; i++ {
		scratchDisks = append(scratchDisks, &compute.AttachedDisk{
			AutoDelete: true,
			Type:       "SCRATCH",
			Interface:  d.Get(fmt.Sprintf("scratch_disk.%d.interface", i)).(string),
			DiskSizeGb: int64(d.Get(fmt.Sprintf("scratch_disk.%d.size", i)).(int)),
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskType: diskType.RelativeLink(),
			},
		})
	}

	return scratchDisks, nil
}

func flattenScratchDisk(disk *compute.AttachedDisk) map[string]interface{} {
	result := map[string]interface{}{
		"interface": disk.Interface,
		"size":      disk.DiskSizeGb,
	}
	return result
}

func hash256(raw string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(decoded)
	return base64.StdEncoding.EncodeToString(h[:]), nil
}

func serviceAccountDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k != "service_account.#" {
		return false
	}

	o, n := d.GetChange("service_account")
	var l []interface{}
	if old == "0" && new == "1" {
		l = n.([]interface{})
	} else if new == "0" && old == "1" {
		l = o.([]interface{})
	} else {
		// we don't have one set and one unset, so don't suppress the diff
		return false
	}

	// suppress changes between { } and {scopes:[]}
	if l[0] != nil {
		contents := l[0].(map[string]interface{})
		if scopes, ok := contents["scopes"]; ok {
			a := scopes.(*schema.Set).List()
			if a != nil && len(a) > 0 {
				return false
			}
		}
	}
	return true
}
