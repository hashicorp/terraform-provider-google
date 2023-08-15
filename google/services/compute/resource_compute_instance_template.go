// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/compute/v1"
)

var (
	schedulingInstTemplateKeys = []string{
		"scheduling.0.on_host_maintenance",
		"scheduling.0.automatic_restart",
		"scheduling.0.preemptible",
		"scheduling.0.node_affinities",
		"scheduling.0.min_node_cpus",
		"scheduling.0.provisioning_model",
		"scheduling.0.instance_termination_action",
		"scheduling.0.local_ssd_recovery_timeout",
	}

	shieldedInstanceTemplateConfigKeys = []string{
		"shielded_instance_config.0.enable_secure_boot",
		"shielded_instance_config.0.enable_vtpm",
		"shielded_instance_config.0.enable_integrity_monitoring",
	}
)

var DEFAULT_SCRATCH_DISK_SIZE_GB = 375
var VALID_SCRATCH_DISK_SIZES_GB [2]int = [2]int{375, 3000}

func ResourceComputeInstanceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceTemplateCreate,
		Read:   resourceComputeInstanceTemplateRead,
		Delete: resourceComputeInstanceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeInstanceTemplateImportState,
		},
		SchemaVersion: 1,
		CustomizeDiff: customdiff.All(
			resourceComputeInstanceTemplateSourceImageCustomizeDiff,
			resourceComputeInstanceTemplateScratchDiskCustomizeDiff,
			resourceComputeInstanceTemplateBootDiskCustomizeDiff,
		),
		MigrateState: resourceComputeInstanceTemplateMigrateState,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		// A compute instance template is more or less a subset of a compute
		// instance. Please attempt to maintain consistency with the
		// resource_compute_instance schema when updating this one.
		Schema: map[string]*schema.Schema{
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
					// uuid is 26 characters, limit the prefix to 37.
					value := v.(string)
					if len(value) > 37 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 37 characters, name is limited to 63", k))
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
										DiffSuppressFunc: tpgresource.IpCidrRangeDiffSuppress,
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
									// Possibly configurable- this was added so we don't break if it's inadvertently set
									// (assuming the same ass access config)
									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The domain name to be used when creating DNSv6 records for the external IPv6 ranges.`,
									},
									"external_ipv6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The first IPv6 address of the external IPv6 range associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig. The field is output only, an IPv6 address from a subnetwork associated with the instance will be allocated dynamically.`,
									},
									"external_ipv6_prefix_length": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The prefix length of the external IPv6 range.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of this access configuration.`,
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

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `An instance template is a global resource that is not bound to a zone or a region. However, you can still specify some regional resources in an instance template, which restricts the template to the region where that resource resides. For example, a custom subnetwork resource is tied to a specific region. Defaults to the region of the Provider if no value is given.`,
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
						"local_ssd_recovery_timeout": {
							Type:     schema.TypeList,
							Optional: true,
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
				Description: `The URI of the created resource.`,
			},

			"self_link_unique": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A special URI of the created resource that uniquely identifies this instance template.`,
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
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: `Defines whether the instance should have confidential compute enabled.`,
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
				Description: `The unique fingerprint of the tags.`,
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: `A set of key/value label pairs to assign to instances created from this template,`,
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

func resourceComputeInstanceTemplateSourceImageCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	numDisks := diff.Get("disk.#").(int)
	for i := 0; i < numDisks; i++ {
		key := fmt.Sprintf("disk.%d.source_image", i)
		if diff.HasChange(key) {
			var err error
			old, new := diff.GetChange(key)
			if old == "" || new == "" {
				continue
			}
			// project must be retrieved once we know there is a diff to resolve, otherwise it will
			// attempt to retrieve project during `plan` before all calculated fields are ready
			// see https://github.com/hashicorp/terraform-provider-google/issues/2878
			project, err := tpgresource.GetProjectFromDiff(diff, config)
			if err != nil {
				return err
			}
			oldResolved, err := ResolveImage(config, project, old.(string), config.UserAgent)
			if err != nil {
				return err
			}
			oldResolved, err = resolveImageRefToRelativeURI(project, oldResolved)
			if err != nil {
				return err
			}
			newResolved, err := ResolveImage(config, project, new.(string), config.UserAgent)
			if err != nil {
				return err
			}
			newResolved, err = resolveImageRefToRelativeURI(project, newResolved)
			if err != nil {
				return err
			}
			if oldResolved != newResolved {
				continue
			}
			err = diff.Clear(key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceComputeInstanceTemplateScratchDiskCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return resourceComputeInstanceTemplateScratchDiskCustomizeDiffFunc(diff)
}

func resourceComputeInstanceTemplateScratchDiskCustomizeDiffFunc(diff tpgresource.TerraformResourceDiff) error {
	numDisks := diff.Get("disk.#").(int)
	for i := 0; i < numDisks; i++ {
		// misspelled on purpose, type is a special symbol
		typee := diff.Get(fmt.Sprintf("disk.%d.type", i)).(string)
		diskType := diff.Get(fmt.Sprintf("disk.%d.disk_type", i)).(string)
		if typee == "SCRATCH" && diskType != "local-ssd" {
			return fmt.Errorf("SCRATCH disks must have a disk_type of local-ssd. disk %d has disk_type %s", i, diskType)
		}

		if diskType == "local-ssd" && typee != "SCRATCH" {
			return fmt.Errorf("disks with a disk_type of local-ssd must be SCRATCH disks. disk %d is a %s disk", i, typee)
		}

		diskSize := diff.Get(fmt.Sprintf("disk.%d.disk_size_gb", i)).(int)
		if typee == "SCRATCH" && !(diskSize == 375 || diskSize == 3000) { // see VALID_SCRATCH_DISK_SIZES_GB
			return fmt.Errorf("SCRATCH disks must be one of %v GB, disk %d is %d", VALID_SCRATCH_DISK_SIZES_GB, i, diskSize)
		}

		interfacee := diff.Get(fmt.Sprintf("disk.%d.interface", i)).(string)
		if typee == "SCRATCH" && diskSize == 3000 && interfacee != "NVME" {
			return fmt.Errorf("SCRATCH disks with a size of 3000 GB must have an interface of NVME. disk %d has interface %s", i, interfacee)
		}
	}

	return nil
}

func resourceComputeInstanceTemplateBootDiskCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	numDisks := diff.Get("disk.#").(int)
	// No disk except the first can be the boot disk
	for i := 1; i < numDisks; i++ {
		key := fmt.Sprintf("disk.%d.boot", i)
		if v, ok := diff.GetOk(key); ok {
			if v.(bool) {
				return fmt.Errorf("Only the first disk specified in instance_template can be the boot disk. %s was true", key)
			}
		}
	}
	return nil
}

func buildDisks(d *schema.ResourceData, config *transport_tpg.Config) ([]*compute.AttachedDisk, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	disksCount := d.Get("disk.#").(int)

	disks := make([]*compute.AttachedDisk, 0, disksCount)
	for i := 0; i < disksCount; i++ {
		prefix := fmt.Sprintf("disk.%d", i)

		// Build the disk
		var disk compute.AttachedDisk
		disk.Type = "PERSISTENT"
		disk.Mode = "READ_WRITE"
		disk.Interface = "SCSI"
		disk.Boot = i == 0
		disk.AutoDelete = d.Get(prefix + ".auto_delete").(bool)

		if v, ok := d.GetOk(prefix + ".boot"); ok {
			disk.Boot = v.(bool)
		}

		if v, ok := d.GetOk(prefix + ".device_name"); ok {
			disk.DeviceName = v.(string)
		}

		if _, ok := d.GetOk(prefix + ".disk_encryption_key"); ok {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{}
			if v, ok := d.GetOk(prefix + ".disk_encryption_key.0.kms_key_self_link"); ok {
				disk.DiskEncryptionKey.KmsKeyName = v.(string)
			}
		}
		// Assign disk.DiskSizeGb and disk.InitializeParams.DiskSizeGb the same value
		if v, ok := d.GetOk(prefix + ".disk_size_gb"); ok {
			disk.DiskSizeGb = int64(v.(int))
		}
		if v, ok := d.GetOk(prefix + ".source"); ok {
			disk.Source = v.(string)
			conflicts := []string{"disk_size_gb", "disk_name", "disk_type", "provisioned_iops", "source_image", "source_snapshot", "labels"}
			for _, conflict := range conflicts {
				if _, ok := d.GetOk(prefix + "." + conflict); ok {
					return nil, fmt.Errorf("Cannot use `source` with any of the fields in %s", conflicts)
				}
			}
		} else {
			disk.InitializeParams = &compute.AttachedDiskInitializeParams{}

			if v, ok := d.GetOk(prefix + ".disk_name"); ok {
				disk.InitializeParams.DiskName = v.(string)
			}
			// Assign disk.DiskSizeGb and disk.InitializeParams.DiskSizeGb the same value
			if v, ok := d.GetOk(prefix + ".disk_size_gb"); ok {
				disk.InitializeParams.DiskSizeGb = int64(v.(int))
			}
			disk.InitializeParams.DiskType = "pd-standard"
			if v, ok := d.GetOk(prefix + ".disk_type"); ok {
				disk.InitializeParams.DiskType = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".provisioned_iops"); ok {
				disk.InitializeParams.ProvisionedIops = int64(v.(int))
			}

			disk.InitializeParams.Labels = tpgresource.ExpandStringMap(d, prefix+".labels")

			if v, ok := d.GetOk(prefix + ".source_image"); ok {
				imageName := v.(string)
				imageUrl, err := ResolveImage(config, project, imageName, userAgent)
				if err != nil {
					return nil, fmt.Errorf(
						"Error resolving image name '%s': %s",
						imageName, err)
				}
				disk.InitializeParams.SourceImage = imageUrl
			}

			if _, ok := d.GetOk(prefix + ".source_image_encryption_key"); ok {
				disk.InitializeParams.SourceImageEncryptionKey = &compute.CustomerEncryptionKey{}
				if v, ok := d.GetOk(prefix + ".source_image_encryption_key.0.kms_key_self_link"); ok {
					disk.InitializeParams.SourceImageEncryptionKey.KmsKeyName = v.(string)
				}
				if v, ok := d.GetOk(prefix + ".source_image_encryption_key.0.kms_key_service_account"); ok {
					disk.InitializeParams.SourceImageEncryptionKey.KmsKeyServiceAccount = v.(string)
				}
			}

			if v, ok := d.GetOk(prefix + ".source_snapshot"); ok {
				disk.InitializeParams.SourceSnapshot = v.(string)
			}

			if _, ok := d.GetOk(prefix + ".source_snapshot_encryption_key"); ok {
				disk.InitializeParams.SourceSnapshotEncryptionKey = &compute.CustomerEncryptionKey{}
				if v, ok := d.GetOk(prefix + ".source_snapshot_encryption_key.0.kms_key_self_link"); ok {
					disk.InitializeParams.SourceSnapshotEncryptionKey.KmsKeyName = v.(string)
				}
				if v, ok := d.GetOk(prefix + ".source_snapshot_encryption_key.0.kms_key_service_account"); ok {
					disk.InitializeParams.SourceSnapshotEncryptionKey.KmsKeyServiceAccount = v.(string)
				}
			}

			if _, ok := d.GetOk(prefix + ".resource_policies"); ok {
				// instance template only supports a resource name here (not uri)
				disk.InitializeParams.ResourcePolicies = expandInstanceTemplateResourcePolicies(d, prefix+".resource_policies")
			}
		}

		if v, ok := d.GetOk(prefix + ".interface"); ok {
			disk.Interface = v.(string)
		}

		if v, ok := d.GetOk(prefix + ".mode"); ok {
			disk.Mode = v.(string)
		}

		if v, ok := d.GetOk(prefix + ".type"); ok {
			disk.Type = v.(string)
		}

		disks = append(disks, &disk)
	}

	return disks, nil
}

// We don't share this code with compute instances because instances want a
// partial URL, but instance templates want the bare accelerator name (despite
// the docs saying otherwise).
//
// Using a partial URL on an instance template results in:
// Invalid value for field 'resource.properties.guestAccelerators[0].acceleratorType':
// 'zones/us-east1-b/acceleratorTypes/nvidia-tesla-k80'.
// Accelerator type 'zones/us-east1-b/acceleratorTypes/nvidia-tesla-k80'
// must be a valid resource name (not an url).
func expandInstanceTemplateGuestAccelerators(d tpgresource.TerraformResourceData, config *transport_tpg.Config) []*compute.AcceleratorConfig {
	configs, ok := d.GetOk("guest_accelerator")
	if !ok {
		return nil
	}
	accels := configs.([]interface{})
	guestAccelerators := make([]*compute.AcceleratorConfig, 0, len(accels))
	for _, raw := range accels {
		data := raw.(map[string]interface{})
		if data["count"].(int) == 0 {
			continue
		}
		guestAccelerators = append(guestAccelerators, &compute.AcceleratorConfig{
			AcceleratorCount: int64(data["count"].(int)),
			// We can't use ParseAcceleratorFieldValue here because an instance
			// template does not have a zone we can use.
			AcceleratorType: data["type"].(string),
		})
	}

	return guestAccelerators
}

func expandInstanceTemplateResourcePolicies(d tpgresource.TerraformResourceData, dataKey string) []string {
	return tpgresource.ConvertAndMapStringArr(d.Get(dataKey).([]interface{}), tpgresource.GetResourceNameFromSelfLink)
}

func resourceComputeInstanceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
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

	if _, ok := d.GetOk("labels"); ok {
		instanceProperties.Labels = tpgresource.ExpandLabels(d)
	}

	var itName string
	if v, ok := d.GetOk("name"); ok {
		itName = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		itName = resource.PrefixedUniqueId(v.(string))
	} else {
		itName = resource.UniqueId()
	}
	instanceTemplate := &compute.InstanceTemplate{
		Description: d.Get("description").(string),
		Properties:  instanceProperties,
		Name:        itName,
	}

	op, err := config.NewComputeClient(userAgent).InstanceTemplates.Insert(project, instanceTemplate).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance template: %s", err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("projects/%s/global/instanceTemplates/%s", project, instanceTemplate.Name))
	// And also the unique ID
	d.Set("self_link_unique", fmt.Sprintf("%v?uniqueId=%v", d.Id(), op.TargetId))

	err = ComputeOperationWaitTime(config, op, project, "Creating Instance Template", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeInstanceTemplateRead(d, meta)
}

type diskCharacteristics struct {
	mode            string
	diskType        string
	diskSizeGb      string
	autoDelete      bool
	sourceImage     string
	provisionedIops string
}

func diskCharacteristicsFromMap(m map[string]interface{}) diskCharacteristics {
	dc := diskCharacteristics{}
	if v := m["mode"]; v == nil || v.(string) == "" {
		// mode has an apply-time default of READ_WRITE
		dc.mode = "READ_WRITE"
	} else {
		dc.mode = v.(string)
	}

	if v := m["disk_type"]; v != nil {
		dc.diskType = v.(string)
	}

	if v := m["disk_size_gb"]; v != nil {
		// Terraform and GCP return ints as different types (int vs int64), so just
		// use strings to compare for simplicity.
		dc.diskSizeGb = fmt.Sprintf("%v", v)
	}

	if v := m["auto_delete"]; v != nil {
		dc.autoDelete = v.(bool)
	}

	if v := m["source_image"]; v != nil {
		dc.sourceImage = v.(string)
	}

	if v := m["provisioned_iops"]; v != nil {
		// Terraform and GCP return ints as different types (int vs int64), so just
		// use strings to compare for simplicity.
		dc.provisionedIops = fmt.Sprintf("%v", v)
	}

	return dc
}

func flattenDisk(disk *compute.AttachedDisk, configDisk map[string]any, defaultProject string) (map[string]interface{}, error) {
	diskMap := make(map[string]interface{})

	// These values are not returned by the API, so we copy them from the config.
	diskMap["source_image_encryption_key"] = configDisk["source_image_encryption_key"]
	diskMap["source_snapshot"] = configDisk["source_snapshot"]
	diskMap["source_snapshot_encryption_key"] = configDisk["source_snapshot_encryption_key"]

	if disk.InitializeParams != nil {
		if disk.InitializeParams.SourceImage != "" {
			path, err := resolveImageRefToRelativeURI(defaultProject, disk.InitializeParams.SourceImage)
			if err != nil {
				return nil, errwrap.Wrapf("Error expanding source image input to relative URI: {{err}}", err)
			}
			diskMap["source_image"] = path
		} else {
			diskMap["source_image"] = ""
		}
		diskMap["disk_type"] = disk.InitializeParams.DiskType
		diskMap["provisioned_iops"] = disk.InitializeParams.ProvisionedIops
		diskMap["disk_name"] = disk.InitializeParams.DiskName
		diskMap["labels"] = disk.InitializeParams.Labels
		// The API does not return a disk size value for scratch disks. They are largely only one size,
		// so we can assume that size here. Prefer disk.DiskSizeGb over the deprecated
		// disk.InitializeParams.DiskSizeGb.
		if disk.DiskSizeGb == 0 && disk.InitializeParams.DiskSizeGb == 0 && disk.Type == "SCRATCH" {
			diskMap["disk_size_gb"] = DEFAULT_SCRATCH_DISK_SIZE_GB
		} else if disk.DiskSizeGb != 0 {
			diskMap["disk_size_gb"] = disk.DiskSizeGb
		} else {
			diskMap["disk_size_gb"] = disk.InitializeParams.DiskSizeGb
		}

		diskMap["resource_policies"] = disk.InitializeParams.ResourcePolicies
	}

	if disk.DiskEncryptionKey != nil {
		encryption := make([]map[string]interface{}, 1)
		encryption[0] = make(map[string]interface{})
		encryption[0]["kms_key_self_link"] = disk.DiskEncryptionKey.KmsKeyName
		diskMap["disk_encryption_key"] = encryption
	}

	diskMap["auto_delete"] = disk.AutoDelete
	diskMap["boot"] = disk.Boot
	diskMap["device_name"] = disk.DeviceName
	diskMap["interface"] = disk.Interface
	diskMap["source"] = tpgresource.ConvertSelfLinkToV1(disk.Source)
	diskMap["mode"] = disk.Mode
	diskMap["type"] = disk.Type

	return diskMap, nil
}

func reorderDisks(configDisks []interface{}, apiDisks []map[string]interface{}) []map[string]interface{} {
	if len(apiDisks) != len(configDisks) {
		// There are different numbers of disks in state and returned from the API, so it's not
		// worth trying to reorder them since it'll be a diff anyway.
		return apiDisks
	}

	result := make([]map[string]interface{}, len(apiDisks))

	/*
		Disks aren't necessarily returned from the API in the same order they were sent, so gather
		information about the ones in state that we can use to map it back. We can't do this by
		just looping over all of the disks, because you could end up matching things in the wrong
		order. For example, if the config disks contain the following disks:
		disk 1: auto delete = false, size = 10
		disk 2: auto delete = false, size = 10, device name = "disk 2"
		disk 3: type = scratch
		And the disks returned from the API are:
		disk a: auto delete = false, size = 10, device name = "disk 2"
		disk b: auto delete = false, size = 10, device name = "disk 1"
		disk c: type = scratch
		Then disk a will match disk 1, disk b won't match any disk, and c will match 3, making the
		final order a, c, b, which is wrong. To get disk a to match disk 2, we have to go in order
		of fields most specifically able to identify a disk to least.
	*/
	disksByDeviceName := map[string]int{}
	scratchDisksByInterface := map[string][]int{}
	attachedDisksBySource := map[string]int{}
	attachedDisksByDiskName := map[string]int{}
	attachedDisksByCharacteristics := []int{}

	for i, d := range configDisks {
		if i == 0 {
			// boot disk
			continue
		}
		disk := d.(map[string]interface{})
		if v := disk["device_name"]; v.(string) != "" {
			disksByDeviceName[v.(string)] = i
		} else if v := disk["type"]; v.(string) == "SCRATCH" {
			iface := disk["interface"].(string)
			if iface == "" {
				// apply-time default
				iface = "SCSI"
			}
			scratchDisksByInterface[iface] = append(scratchDisksByInterface[iface], i)
		} else if v := disk["source"]; v.(string) != "" {
			attachedDisksBySource[v.(string)] = i
		} else if v := disk["disk_name"]; v.(string) != "" {
			attachedDisksByDiskName[v.(string)] = i
		} else {
			attachedDisksByCharacteristics = append(attachedDisksByCharacteristics, i)
		}
	}

	// Align the disks, going from the most specific criteria to the least.
	for _, apiDisk := range apiDisks {
		// 1. This resource only works if the boot disk is the first one (which should be fixed
		//	  separately), so put the boot disk first.
		if apiDisk["boot"].(bool) {
			result[0] = apiDisk

			// 2. All disks have a unique device name
		} else if i, ok := disksByDeviceName[apiDisk["device_name"].(string)]; ok {
			result[i] = apiDisk

			// 3. Scratch disks are all the same except device name and interface, so match them by
			//    interface.
		} else if apiDisk["type"].(string) == "SCRATCH" {
			iface := apiDisk["interface"].(string)
			indexes := scratchDisksByInterface[iface]
			if len(indexes) > 0 {
				result[indexes[0]] = apiDisk
				scratchDisksByInterface[iface] = indexes[1:]
			} else {
				result = append(result, apiDisk)
			}

			// 4. Each attached disk will have a different source, so match by that.
		} else if i, ok := attachedDisksBySource[apiDisk["source"].(string)]; ok {
			result[i] = apiDisk

			// 5. If a disk was created for this resource via initializeParams, it will have a
			//    unique name.
		} else if v, ok := apiDisk["disk_name"]; ok && attachedDisksByDiskName[v.(string)] != 0 {
			result[attachedDisksByDiskName[v.(string)]] = apiDisk

			// 6. If no unique keys exist on this disk, then use a combination of its remaining
			//    characteristics to see whether it matches exactly.
		} else {
			found := false
			for arrayIndex, i := range attachedDisksByCharacteristics {
				configDisk := configDisks[i].(map[string]interface{})
				stateDc := diskCharacteristicsFromMap(configDisk)
				readDc := diskCharacteristicsFromMap(apiDisk)
				if reflect.DeepEqual(stateDc, readDc) {
					result[i] = apiDisk
					attachedDisksByCharacteristics = append(attachedDisksByCharacteristics[:arrayIndex], attachedDisksByCharacteristics[arrayIndex+1:]...)
					found = true
					break
				}
			}
			if !found {
				result = append(result, apiDisk)
			}
		}
	}

	// Remove nils from map in case there were disks that could not be matched
	ds := []map[string]interface{}{}
	for _, d := range result {
		if d != nil {
			ds = append(ds, d)
		}
	}
	return ds
}

func flattenDisks(disks []*compute.AttachedDisk, d *schema.ResourceData, defaultProject string) ([]map[string]interface{}, error) {
	apiDisks := make([]map[string]interface{}, len(disks))

	for i, disk := range disks {
		configDisk := d.Get(fmt.Sprintf("disk.%d", i)).(map[string]any)
		apiDisk, err := flattenDisk(disk, configDisk, defaultProject)
		if err != nil {
			return nil, err
		}
		apiDisks[i] = apiDisk
	}

	return reorderDisks(d.Get("disk").([]interface{}), apiDisks), nil
}

func resourceComputeInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	idStr := d.Id()
	if v, ok := d.GetOk("self_link_unique"); ok && v != "" {
		idStr = ConvertToUniqueIdWhenPresent(v.(string))
	}

	splits := strings.Split(idStr, "/")
	instanceTemplate, err := config.NewComputeClient(userAgent).InstanceTemplates.Get(project, splits[len(splits)-1]).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Instance Template %q", d.Get("name").(string)))
	}
	// Set the metadata fingerprint if there is one.
	if instanceTemplate.Properties.Metadata != nil {
		if err = d.Set("metadata_fingerprint", instanceTemplate.Properties.Metadata.Fingerprint); err != nil {
			return fmt.Errorf("Error setting metadata_fingerprint: %s", err)
		}

		md := instanceTemplate.Properties.Metadata

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
	if instanceTemplate.Properties.Tags != nil {
		if err = d.Set("tags_fingerprint", instanceTemplate.Properties.Tags.Fingerprint); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
	} else {
		if err := d.Set("tags_fingerprint", ""); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
	}
	if instanceTemplate.Properties.Labels != nil {
		if err := d.Set("labels", instanceTemplate.Properties.Labels); err != nil {
			return fmt.Errorf("Error setting labels: %s", err)
		}
	}
	if err = d.Set("self_link", instanceTemplate.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err = d.Set("self_link_unique", fmt.Sprintf("%v?uniqueId=%v", instanceTemplate.SelfLink, instanceTemplate.Id)); err != nil {
		return fmt.Errorf("Error setting self_link_unique: %s", err)
	}
	if err = d.Set("name", instanceTemplate.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if instanceTemplate.Properties.Disks != nil {
		disks, err := flattenDisks(instanceTemplate.Properties.Disks, d, project)
		if err != nil {
			return fmt.Errorf("error flattening disks: %s", err)
		}
		if err = d.Set("disk", disks); err != nil {
			return fmt.Errorf("Error setting disk: %s", err)
		}
	}
	if err = d.Set("description", instanceTemplate.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("machine_type", instanceTemplate.Properties.MachineType); err != nil {
		return fmt.Errorf("Error setting machine_type: %s", err)
	}
	if err = d.Set("min_cpu_platform", instanceTemplate.Properties.MinCpuPlatform); err != nil {
		return fmt.Errorf("Error setting min_cpu_platform: %s", err)
	}

	if err = d.Set("can_ip_forward", instanceTemplate.Properties.CanIpForward); err != nil {
		return fmt.Errorf("Error setting can_ip_forward: %s", err)
	}

	if err = d.Set("instance_description", instanceTemplate.Properties.Description); err != nil {
		return fmt.Errorf("Error setting instance_description: %s", err)
	}
	if err = d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("network_performance_config", flattenNetworkPerformanceConfig(instanceTemplate.Properties.NetworkPerformanceConfig)); err != nil {
		return err
	}
	if instanceTemplate.Properties.NetworkInterfaces != nil {
		networkInterfaces, region, _, _, err := flattenNetworkInterfaces(d, config, instanceTemplate.Properties.NetworkInterfaces)
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
	if instanceTemplate.Properties.Scheduling != nil {
		scheduling := flattenScheduling(instanceTemplate.Properties.Scheduling)
		if err = d.Set("scheduling", scheduling); err != nil {
			return fmt.Errorf("Error setting scheduling: %s", err)
		}
	}
	if instanceTemplate.Properties.Tags != nil {
		if err = d.Set("tags", instanceTemplate.Properties.Tags.Items); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	} else {
		if err = d.Set("tags", nil); err != nil {
			return fmt.Errorf("Error setting empty tags: %s", err)
		}
	}
	if instanceTemplate.Properties.ServiceAccounts != nil {
		if err = d.Set("service_account", flattenServiceAccounts(instanceTemplate.Properties.ServiceAccounts)); err != nil {
			return fmt.Errorf("Error setting service_account: %s", err)
		}
	}
	if instanceTemplate.Properties.GuestAccelerators != nil {
		if err = d.Set("guest_accelerator", flattenGuestAccelerators(instanceTemplate.Properties.GuestAccelerators)); err != nil {
			return fmt.Errorf("Error setting guest_accelerator: %s", err)
		}
	}
	if instanceTemplate.Properties.ShieldedInstanceConfig != nil {
		if err = d.Set("shielded_instance_config", flattenShieldedVmConfig(instanceTemplate.Properties.ShieldedInstanceConfig)); err != nil {
			return fmt.Errorf("Error setting shielded_instance_config: %s", err)
		}
	}

	if instanceTemplate.Properties.ConfidentialInstanceConfig != nil {
		if err = d.Set("confidential_instance_config", flattenConfidentialInstanceConfig(instanceTemplate.Properties.ConfidentialInstanceConfig)); err != nil {
			return fmt.Errorf("Error setting confidential_instance_config: %s", err)
		}
	}
	if instanceTemplate.Properties.AdvancedMachineFeatures != nil {
		if err = d.Set("advanced_machine_features", flattenAdvancedMachineFeatures(instanceTemplate.Properties.AdvancedMachineFeatures)); err != nil {
			return fmt.Errorf("Error setting advanced_machine_features: %s", err)
		}
	}

	if instanceTemplate.Properties.ResourcePolicies != nil {
		if err = d.Set("resource_policies", instanceTemplate.Properties.ResourcePolicies); err != nil {
			return fmt.Errorf("Error setting resource_policies: %s", err)
		}
	}

	if reservationAffinity := instanceTemplate.Properties.ReservationAffinity; reservationAffinity != nil {
		if err = d.Set("reservation_affinity", flattenReservationAffinity(reservationAffinity)); err != nil {
			return fmt.Errorf("Error setting reservation_affinity: %s", err)
		}
	}

	return nil
}

func resourceComputeInstanceTemplateDelete(d *schema.ResourceData, meta interface{}) error {
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
	op, err := config.NewComputeClient(userAgent).InstanceTemplates.Delete(
		project, splits[len(splits)-1]).Do()
	if err != nil {
		return fmt.Errorf("Error deleting instance template: %s", err)
	}

	err = ComputeOperationWaitTime(config, op, project, "Deleting Instance Template", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// This wraps the general compute instance helper expandScheduling.
// Default value of OnHostMaintenance depends on the value of Preemptible,
// so we can't set a default in schema
func expandResourceComputeInstanceTemplateScheduling(d *schema.ResourceData, meta interface{}) (*compute.Scheduling, error) {
	v, ok := d.GetOk("scheduling")
	if !ok || v == nil {
		// We can't set defaults for lists (e.g. scheduling)
		return &compute.Scheduling{
			OnHostMaintenance: "MIGRATE",
		}, nil
	}

	expanded, err := expandScheduling(v)
	if err != nil {
		return nil, err
	}

	// Make sure we have an appropriate value for OnHostMaintenance if Preemptible
	if expanded.Preemptible && expanded.OnHostMaintenance == "" {
		expanded.OnHostMaintenance = "TERMINATE"
	}
	return expanded, nil
}

func resourceComputeInstanceTemplateImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{"projects/(?P<project>[^/]+)/global/instanceTemplates/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/global/instanceTemplates/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
