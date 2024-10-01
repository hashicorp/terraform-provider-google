// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/compute/v1"
)

func ResourceComputeRegionInstanceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionInstanceTemplateCreate,
		Read:   resourceComputeRegionInstanceTemplateRead,
		Update: resourceComputeRegionInstanceTemplateUpdate,
		Delete: resourceComputeRegionInstanceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeRegionInstanceTemplateImportState,
		},
		SchemaVersion: 1,
		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.DefaultProviderRegion,
			resourceComputeInstanceTemplateSourceImageCustomizeDiff,
			resourceComputeInstanceTemplateScratchDiskCustomizeDiff,
			resourceComputeInstanceTemplateBootDiskCustomizeDiff,
			tpgresource.SetLabelsDiff,
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		// A compute region instance template is more or less a subset of a compute
		// instance. Please attempt to maintain consistency with the
		// resource_compute_instance schema when updating this one.
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which the instance template is located. If it is not provided, the provider region is used.`,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  verify.ValidateGCEName,
				Description:   `The name of the instance template. If you leave this blank, Terraform will auto-generate a unique name.`,
			},

			"name_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Creates a unique name beginning with the specified prefix. Conflicts with name.`,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					// https://cloud.google.com/compute/docs/reference/latest/instanceTemplates#resource
					// uuid is 9 characters, limit the prefix to 54.
					value := v.(string)
					if len(value) > 54 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 54 characters, name is limited to 63", k))
					}
					return
				},
			},

			"disk": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `Disks to attach to instances created from this template. This can be specified multiple times for multiple disks.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: `Whether or not the disk should be auto-deleted. This defaults to true.`,
						},

						"boot": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `Indicates that this is a boot disk.`,
						},

						"device_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `A unique device name that is reflected into the /dev/ tree of a Linux operating system running within the instance. If not specified, the server chooses a default device name to apply to this disk.`,
						},

						"disk_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Name of the disk. When not provided, this defaults to the name of the instance.`,
						},

						"disk_size_gb": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `The size of the image in gigabytes. If not specified, it will inherit the size of its base image. For SCRATCH disks, the size must be one of 375 or 3000 GB, with a default of 375 GB.`,
						},

						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `The Google Compute Engine disk type. Such as "pd-ssd", "local-ssd", "pd-balanced" or "pd-standard".`,
						},

						"labels": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: `A set of key/value label pairs to assign to disks,`,
						},

						"provisioned_iops": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `Indicates how many IOPS to provision for the disk. This sets the number of I/O operations per second that the disk can handle. Values must be between 10,000 and 120,000. For more details, see the [Extreme persistent disk documentation](https://cloud.google.com/compute/docs/disks/extreme-persistent-disk).`,
						},

						"resource_manager_tags": {
							Type:        schema.TypeMap,
							Optional:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
						},

						"source_image": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `The image from which to initialize this disk. This can be one of: the image's self_link, projects/{project}/global/images/{image}, projects/{project}/global/images/family/{family}, global/images/{image}, global/images/family/{family}, family/{family}, {project}/{family}, {project}/{image}, {family}, or {image}. ~> Note: Either source or source_image is required when creating a new instance except for when creating a local SSD.`,
						},
						"source_image_encryption_key": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Description: `The customer-supplied encryption key of the source
image. Required if the source image is protected by a
customer-supplied encryption key.

Instance templates do not store customer-supplied
encryption keys, so you cannot create disks for
instances in a managed instance group if the source
images are encrypted with your own keys.`,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_service_account": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: `The service account being used for the encryption
request for the given KMS key. If absent, the Compute
Engine default service account is used.`,
									},
									"kms_key_self_link": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										Description: `The self link of the encryption key that is stored in
Google Cloud KMS.`,
									},
								},
							},
						},
						"source_snapshot": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: `The source snapshot to create this disk. When creating
a new instance, one of initializeParams.sourceSnapshot,
initializeParams.sourceImage, or disks.source is
required except for local SSD.`,
						},
						"source_snapshot_encryption_key": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: `The customer-supplied encryption key of the source snapshot.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_service_account": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: `The service account being used for the encryption
request for the given KMS key. If absent, the Compute
Engine default service account is used.`,
									},
									"kms_key_self_link": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										Description: `The self link of the encryption key that is stored in
Google Cloud KMS.`,
									},
								},
							},
						},

						"interface": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `Specifies the disk interface to use for attaching this disk.`,
						},

						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `The mode in which to attach this disk, either READ_WRITE or READ_ONLY. If you are attaching or creating a boot disk, this must read-write mode.`,
						},

						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The name (not self_link) of the disk (such as those managed by google_compute_disk) to attach. ~> Note: Either source or source_image is required when creating a new instance except for when creating a local SSD.`,
						},

						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `The type of Google Compute Engine disk, can be either "SCRATCH" or "PERSISTENT".`,
						},

						"disk_encryption_key": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: `Encrypts or decrypts a disk using a customer-supplied encryption key.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_self_link": {
										Type:             schema.TypeString,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
										Description:      `The self link of the encryption key that is stored in Google Cloud KMS.`,
									},
								},
							},
						},

						"resource_policies": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: `A list (short name or id) of resource policies to attach to this disk. Currently a max of 1 resource policy is supported.`,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: tpgresource.CompareResourceNames,
							},
						},
					},
				},
			},

			"machine_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The machine type to create. To create a machine with a custom type (such as extended memory), format the value like custom-VCPUS-MEM_IN_MB like custom-6-20480 for 6 vCPU and 20GB of RAM.`,
			},

			"can_ip_forward": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: `Whether to allow sending and receiving of packets with non-matching source or destination IPs. This defaults to false.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A brief description of this resource.`,
			},

			"instance_description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A description of the instance.`,
			},

			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `Metadata key/value pairs to make available from within instances created from this template.`,
			},

			"metadata_startup_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `An alternative to using the startup-script metadata key, mostly to match the compute_instance resource. This replaces the startup-script metadata key on the created instance and thus the two mechanisms are not allowed to be used simultaneously.`,
			},

			"metadata_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The unique fingerprint of the metadata.`,
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

			"network_interface": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `Networks to attach to instances created from this template. This can be specified multiple times for multiple networks.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the network to attach this interface to. Use network attribute for Legacy or Auto subnetted networks and subnetwork for custom subnetted networks.`,
						},

						"subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name of the subnetwork to attach this interface to. The subnetwork must exist in the same region this instance will be created in. Either network or subnetwork must be provided.`,
						},

						"subnetwork_project": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `The ID of the project in which the subnetwork belongs. If it is not provided, the provider project is used.`,
						},

						"network_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The private IP address to assign to the instance. If empty, the address will be automatically assigned.`,
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    true,
							Description: `The name of the network_interface.`,
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
							ForceNew:    true,
							Description: `Access configurations, i.e. IPs via which this instance can be accessed via the Internet. Omit to ensure that the instance is not accessible from the Internet (this means that ssh provisioners will not work unless you are running Terraform can send traffic to the instance's network (e.g. via tunnel or because it is running on another cloud instance on that network). This block can be repeated multiple times.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Computed:    true,
										Description: `The IP address that will be 1:1 mapped to the instance's network ip. If not given, one will be generated.`,
									},
									"network_tier": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: `The networking tier used for configuring this instance template. This field can take the following values: PREMIUM, STANDARD, FIXED_STANDARD. If this field is not specified, it is assumed to be PREMIUM.`,
									},
									// Possibly configurable- this was added so we don't break if it's inadvertently set
									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: `The DNS domain name for the public PTR record.The DNS domain name for the public PTR record.`,
									},
								},
							},
						},

						"alias_ip_range": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: `An array of alias IP ranges for this network interface. Can only be specified for network interfaces on subnet-mode networks.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_cidr_range": {
										Type:             schema.TypeString,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: IpCidrRangeDiffSuppress,
										Description:      `The IP CIDR range represented by this alias IP range. This IP CIDR range must belong to the specified subnetwork and cannot contain IP addresses reserved by system or used by other network interfaces. At the time of writing only a netmask (e.g. /24) may be supplied, with a CIDR format resulting in an API error.`,
									},
									"subnetwork_range_name": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: `The subnetwork secondary range name specifying the secondary range from which to allocate the IP CIDR range for this alias IP range. If left unspecified, the primary range of the subnetwork will be used.`,
									},
								},
							},
						},

						"stack_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"IPV4_ONLY", "IPV4_IPV6", ""}, false),
							Description:  `The stack type for this network interface to identify whether the IPv6 feature is enabled or not. If not specified, IPV4_ONLY will be used.`,
						},

						"ipv6_access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    true,
							Description: `One of EXTERNAL, INTERNAL to indicate whether the IP can be accessed from the Internet. This field is always inherited from its subnetwork.`,
						},

						"ipv6_access_config": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: `An array of IPv6 access configurations for this interface. Currently, only one IPv6 access config, DIRECT_IPV6, is supported. If there is no ipv6AccessConfig specified, then this instance will have no external IPv6 Internet access.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_tier": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `The service-level to be provided for IPv6 traffic when the subnet has an external subnet. Only PREMIUM tier is valid for IPv6`,
									},
									// Possibly configurable- this was added so we don't break if it's inadvertently set
									// (assuming the same ass access config)
									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: `The domain name to be used when creating DNSv6 records for the external IPv6 ranges.`,
									},
									"external_ipv6": {
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: `The first IPv6 address of the external IPv6 range associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig. The field is output only, an IPv6 address from a subnetwork associated with the instance will be allocated dynamically.`,
									},
									"external_ipv6_prefix_length": {
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: `The prefix length of the external IPv6 range.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: `The name of this access configuration.`,
									},
								},
							},
						},
						"internal_ipv6_prefix_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `The prefix length of the primary internal IPv6 range.`,
						},
						"ipv6_address": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: ipv6RepresentationDiffSuppress,
							Description:      `An IPv6 internal network address for this network interface. If not specified, Google Cloud will automatically assign an internal IPv6 address from the instance's subnetwork.`,
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

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"resource_manager_tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Description: `A map of resource manager tags.
				Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
			},

			"scheduling": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The scheduling strategy to use.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preemptible": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: schedulingInstTemplateKeys,
							Default:      false,
							ForceNew:     true,
							Description:  `Allows instance to be preempted. This defaults to false.`,
						},

						"automatic_restart": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: schedulingInstTemplateKeys,
							Default:      true,
							ForceNew:     true,
							Description:  `Specifies whether the instance should be automatically restarted if it is terminated by Compute Engine (not terminated by a user). This defaults to true.`,
						},

						"on_host_maintenance": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: schedulingInstTemplateKeys,
							ForceNew:     true,
							Description:  `Defines the maintenance behavior for this instance.`,
						},

						"node_affinities": {
							Type:             schema.TypeSet,
							Optional:         true,
							AtLeastOneOf:     schedulingInstTemplateKeys,
							ForceNew:         true,
							Elem:             instanceSchedulingNodeAffinitiesElemSchema(),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress(""),
							Description:      `Specifies node affinities or anti-affinities to determine which sole-tenant nodes your instances and managed instance groups will use as host systems.`,
						},
						"min_node_cpus": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: schedulingInstTemplateKeys,
							Description:  `Minimum number of cpus for the instance.`,
						},
						"provisioning_model": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							AtLeastOneOf: schedulingInstTemplateKeys,
							Description:  `Whether the instance is spot. If this is set as SPOT.`,
						},
						"instance_termination_action": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: schedulingInstTemplateKeys,
							Description:  `Specifies the action GCE should take when SPOT VM is preempted.`,
						},
						"max_run_duration": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The timeout for new network connections to hosts.`,
							MaxItems:    1,
							ForceNew:    true,
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
						"on_instance_stop_action": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							ForceNew:    true,
							Description: `Defines the behaviour for instances with the instance_termination_action.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"discard_local_ssd": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `If true, the contents of any attached Local SSD disks will be discarded.`,
										Default:     false,
										ForceNew:    true,
									},
								},
							},
						},
						"local_ssd_recovery_timeout": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Description: `Specifies the maximum amount of time a Local Ssd Vm should wait while
  recovery of the Local Ssd state is attempted. Its value should be in
  between 0 and 168 hours with hour granularity and the default value being 1
  hour.`,

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

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The URI of the created resource.`,
			},

			"service_account": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Service account to attach to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `The service account e-mail address. If not given, the default Google Compute Engine service account is used.`,
						},

						"scopes": {
							Type:        schema.TypeSet,
							Required:    true,
							ForceNew:    true,
							Description: `A list of service scopes. Both OAuth2 URLs and gcloud short names are supported. To allow full access to all Cloud APIs, use the cloud-platform scope.`,
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Enable Shielded VM on this instance. Shielded VM provides verifiable integrity to prevent against malware and rootkits. Defaults to disabled. Note: shielded_instance_config can only be used with boot images with shielded vm support.`,
				// Since this block is used by the API based on which
				// image being used, the field needs to be marked as Computed.
				Computed:         true,
				DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress(""),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_secure_boot": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceTemplateConfigKeys,
							Default:      false,
							ForceNew:     true,
							Description:  `Verify the digital signature of all boot components, and halt the boot process if signature verification fails. Defaults to false.`,
						},

						"enable_vtpm": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceTemplateConfigKeys,
							Default:      true,
							ForceNew:     true,
							Description:  `Use a virtualized trusted platform module, which is a specialized computer chip you can use to encrypt objects like keys and certificates. Defaults to true.`,
						},

						"enable_integrity_monitoring": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceTemplateConfigKeys,
							Default:      true,
							ForceNew:     true,
							Description:  `Compare the most recent boot measurements to the integrity policy baseline and return a pair of pass/fail results depending on whether they match or not. Defaults to true.`,
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
				Description: `The Confidential VM config being used by the instance. on_host_maintenance has to be set to TERMINATE or this will fail to create.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_confidential_compute": {
							Type:         schema.TypeBool,
							Optional:     true,
							ForceNew:     true,
							Description:  `Defines whether the instance should have confidential compute enabled. Field will be deprecated in a future release.`,
							AtLeastOneOf: []string{"confidential_instance_config.0.enable_confidential_compute", "confidential_instance_config.0.confidential_instance_type"},
						},
						"confidential_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: `
								The confidential computing technology the instance uses.
								SEV is an AMD feature. TDX is an Intel feature. One of the following
								values is required: SEV, SEV_SNP, TDX. If SEV_SNP, min_cpu_platform =
								"AMD Milan" is currently required.`,
							AtLeastOneOf: []string{"confidential_instance_config.0.enable_confidential_compute", "confidential_instance_config.0.confidential_instance_type"},
						},
					},
				},
			},
			"advanced_machine_features": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Controls for advanced machine-related behavior features.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_nested_virtualization": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: `Whether to enable nested virtualization or not.`,
						},
						"threads_per_core": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    false,
							ForceNew:    true,
							Description: `The number of threads per physical core. To disable simultaneous multithreading (SMT) set this to 1. If unset, the maximum number of threads supported per core by the underlying processor is assumed.`,
						},
						"visible_core_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `The number of physical cores to expose to an instance. Multiply by the number of threads per core to compute the total number of virtual CPUs to expose to the instance. If unset, the number of cores is inferred from the instance\'s nominal CPU count and the underlying platform\'s SMT width.`,
						},
					},
				},
			},
			"guest_accelerator": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
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
							Description:      `The accelerator type resource to expose to this instance. E.g. nvidia-tesla-k80.`,
						},
					},
				},
			},

			"min_cpu_platform": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies a minimum CPU platform. Applicable values are the friendly names of CPU platforms, such as Intel Haswell or Intel Skylake.`,
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: `Tags to attach to the instance.`,
			},

			"tags_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The unique fingerprint of the tags.`,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Description: `A set of key/value label pairs to assign to instances created from this template,

				**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
				Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Set:         schema.HashString,
				Description: `The combination of labels configured directly on the resource and default labels configured on the provider.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				ForceNew:    true,
				Set:         schema.HashString,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"resource_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `A list of self_links of resource policies to attach to the instance. Currently a max of 1 resource policy is supported.`,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: tpgresource.CompareResourceNames,
				},
			},

			"reservation_affinity": {
				Type:        schema.TypeList,
				MaxItems:    1,
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
		UseJSONNumber: true,
	}
}

func resourceComputeRegionInstanceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	disks, err := buildDisks(d, config)
	if err != nil {
		return err
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return err
	}

	networks, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return err
	}

	scheduling, err := expandResourceComputeInstanceTemplateScheduling(d, config)
	if err != nil {
		return err
	}
	networkPerformanceConfig, err := expandNetworkPerformanceConfig(d, config)
	if err != nil {
		return nil
	}
	reservationAffinity, err := expandReservationAffinity(d)
	if err != nil {
		return err
	}
	resourcePolicies := expandInstanceTemplateResourcePolicies(d, "resource_policies")

	instanceProperties := &compute.InstanceProperties{
		CanIpForward:               d.Get("can_ip_forward").(bool),
		Description:                d.Get("instance_description").(string),
		GuestAccelerators:          expandInstanceTemplateGuestAccelerators(d, config),
		MachineType:                d.Get("machine_type").(string),
		MinCpuPlatform:             d.Get("min_cpu_platform").(string),
		Disks:                      disks,
		Metadata:                   metadata,
		NetworkInterfaces:          networks,
		NetworkPerformanceConfig:   networkPerformanceConfig,
		Scheduling:                 scheduling,
		ServiceAccounts:            expandServiceAccounts(d.Get("service_account").([]interface{})),
		Tags:                       resourceInstanceTags(d),
		ConfidentialInstanceConfig: expandConfidentialInstanceConfig(d),
		ShieldedInstanceConfig:     expandShieldedVmConfigs(d),
		AdvancedMachineFeatures:    expandAdvancedMachineFeatures(d),
		ResourcePolicies:           resourcePolicies,
		ReservationAffinity:        reservationAffinity,
	}

	if _, ok := d.GetOk("effective_labels"); ok {
		instanceProperties.Labels = tpgresource.ExpandEffectiveLabels(d)
	}

	if _, ok := d.GetOk("resource_manager_tags"); ok {
		instanceProperties.ResourceManagerTags = tpgresource.ExpandStringMap(d, "resource_manager_tags")
	}

	var itName string
	if v, ok := d.GetOk("name"); ok {
		itName = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		prefix := v.(string)
		if len(prefix) > 37 {
			itName = tpgresource.ReducedPrefixedUniqueId(prefix)
		} else {
			itName = id.PrefixedUniqueId(prefix)
		}
	} else {
		itName = id.UniqueId()
	}

	instanceTemplate := make(map[string]interface{})
	instanceTemplate["description"] = d.Get("description").(string)
	instanceTemplate["properties"] = instanceProperties
	instanceTemplate["name"] = itName

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceTemplates")
	if err != nil {
		return err
	}

	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      instanceTemplate,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating RegionInstanceTemplate: %s", err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceTemplates/%s", project, region, instanceTemplate["name"]))

	err = ComputeOperationWaitTime(config, op, project, "Creating Region Instance Template", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeRegionInstanceTemplateRead(d, meta)
}

func resourceComputeRegionInstanceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	// Only the field "labels" and "terraform_labels" is mutable
	return resourceComputeRegionInstanceTemplateRead(d, meta)
}

func resourceComputeRegionInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	splits := strings.Split(d.Id(), "/")
	name := splits[len(splits)-1]

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceTemplates/"+name)
	if err != nil {
		return err
	}

	instanceTemplate, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ComputeRegionInstanceTemplate %q", d.Id()))
	}

	instancePropertiesMap := instanceTemplate["properties"]

	instancePropertiesObj, err := json.Marshal(instancePropertiesMap)
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceProperties := compute.InstanceProperties{}

	if err := json.Unmarshal(instancePropertiesObj, &instanceProperties); err != nil {
		fmt.Println(err)
		return err
	}

	// Set the metadata fingerprint if there is one.
	if instanceProperties.Metadata != nil {
		if err = d.Set("metadata_fingerprint", instanceProperties.Metadata.Fingerprint); err != nil {
			return fmt.Errorf("Error setting metadata_fingerprint: %s", err)
		}

		md := instanceProperties.Metadata

		_md := flattenMetadataBeta(md)

		if script, scriptExists := d.GetOk("metadata_startup_script"); scriptExists {
			if err = d.Set("metadata_startup_script", script); err != nil {
				return fmt.Errorf("Error setting metadata_startup_script: %s", err)
			}

			delete(_md, "startup-script")
		}

		if err = d.Set("metadata", _md); err != nil {
			return fmt.Errorf("Error setting metadata: %s", err)
		}
	}

	// Set the tags fingerprint if there is one.
	if instanceProperties.Tags != nil {
		if err = d.Set("tags_fingerprint", instanceProperties.Tags.Fingerprint); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
	} else {
		if err := d.Set("tags_fingerprint", ""); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
	}
	if instanceProperties.Labels != nil {
		if err := tpgresource.SetLabels(instanceProperties.Labels, d, "labels"); err != nil {
			return fmt.Errorf("Error setting labels: %s", err)
		}
	}
	if err := tpgresource.SetLabels(instanceProperties.Labels, d, "terraform_labels"); err != nil {
		return fmt.Errorf("Error setting terraform_labels: %s", err)
	}
	if err := d.Set("effective_labels", instanceProperties.Labels); err != nil {
		return fmt.Errorf("Error setting effective_labels: %s", err)
	}
	if err = d.Set("self_link", instanceTemplate["selfLink"]); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err = d.Set("name", instanceTemplate["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if instanceProperties.Disks != nil {
		disks, err := flattenDisks(instanceProperties.Disks, d, project)
		if err != nil {
			return fmt.Errorf("error flattening disks: %s", err)
		}
		if err = d.Set("disk", disks); err != nil {
			return fmt.Errorf("Error setting disk: %s", err)
		}
	}
	if err = d.Set("description", instanceTemplate["description"]); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("machine_type", instanceProperties.MachineType); err != nil {
		return fmt.Errorf("Error setting machine_type: %s", err)
	}
	if err = d.Set("min_cpu_platform", instanceProperties.MinCpuPlatform); err != nil {
		return fmt.Errorf("Error setting min_cpu_platform: %s", err)
	}

	if err = d.Set("can_ip_forward", instanceProperties.CanIpForward); err != nil {
		return fmt.Errorf("Error setting can_ip_forward: %s", err)
	}

	if err = d.Set("instance_description", instanceProperties.Description); err != nil {
		return fmt.Errorf("Error setting instance_description: %s", err)
	}
	if err = d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("network_performance_config", flattenNetworkPerformanceConfig(instanceProperties.NetworkPerformanceConfig)); err != nil {
		return err
	}
	if instanceProperties.NetworkInterfaces != nil {
		networkInterfaces, region, _, _, err := flattenNetworkInterfaces(d, config, instanceProperties.NetworkInterfaces)
		if err != nil {
			return err
		}
		if err = d.Set("network_interface", networkInterfaces); err != nil {
			return fmt.Errorf("Error setting network_interface: %s", err)
		}
		// region is where to look up the subnetwork if there is one attached to the instance template
		if region != "" {
			if err = d.Set("region", region); err != nil {
				return fmt.Errorf("Error setting region: %s", err)
			}
		}
	}
	if instanceProperties.Scheduling != nil {
		scheduling := flattenScheduling(instanceProperties.Scheduling)
		if err = d.Set("scheduling", scheduling); err != nil {
			return fmt.Errorf("Error setting scheduling: %s", err)
		}
	}
	if instanceProperties.Tags != nil {
		if err = d.Set("tags", instanceProperties.Tags.Items); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	} else {
		if err = d.Set("tags", nil); err != nil {
			return fmt.Errorf("Error setting empty tags: %s", err)
		}
	}
	if instanceProperties.ServiceAccounts != nil {
		if err = d.Set("service_account", flattenServiceAccounts(instanceProperties.ServiceAccounts)); err != nil {
			return fmt.Errorf("Error setting service_account: %s", err)
		}
	}
	if instanceProperties.GuestAccelerators != nil {
		if err = d.Set("guest_accelerator", flattenGuestAccelerators(instanceProperties.GuestAccelerators)); err != nil {
			return fmt.Errorf("Error setting guest_accelerator: %s", err)
		}
	}
	if instanceProperties.ShieldedInstanceConfig != nil {
		if err = d.Set("shielded_instance_config", flattenShieldedVmConfig(instanceProperties.ShieldedInstanceConfig)); err != nil {
			return fmt.Errorf("Error setting shielded_instance_config: %s", err)
		}
	}

	if instanceProperties.ConfidentialInstanceConfig != nil {
		if err = d.Set("confidential_instance_config", flattenConfidentialInstanceConfig(instanceProperties.ConfidentialInstanceConfig)); err != nil {
			return fmt.Errorf("Error setting confidential_instance_config: %s", err)
		}
	}
	if instanceProperties.AdvancedMachineFeatures != nil {
		if err = d.Set("advanced_machine_features", flattenAdvancedMachineFeatures(instanceProperties.AdvancedMachineFeatures)); err != nil {
			return fmt.Errorf("Error setting advanced_machine_features: %s", err)
		}
	}

	if instanceProperties.ResourcePolicies != nil {
		if err = d.Set("resource_policies", instanceProperties.ResourcePolicies); err != nil {
			return fmt.Errorf("Error setting resource_policies: %s", err)
		}
	}

	if reservationAffinity := instanceProperties.ReservationAffinity; reservationAffinity != nil {
		if err = d.Set("reservation_affinity", flattenReservationAffinity(reservationAffinity)); err != nil {
			return fmt.Errorf("Error setting reservation_affinity: %s", err)
		}
	}

	return nil
}

func resourceComputeRegionInstanceTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/instanceTemplates/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "RegionInstanceTemplate")
	}

	err = ComputeOperationWaitTime(config, op, project, "Deleting Region Instance Template", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceComputeRegionInstanceTemplateImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/instanceTemplates/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/instanceTemplates/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
