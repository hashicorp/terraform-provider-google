//
package google

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

	computeBeta "google.golang.org/api/compute/v0.beta"
)

var (
	schedulingInstTemplateKeys = []string{
		"scheduling.0.on_host_maintenance",
		"scheduling.0.automatic_restart",
		"scheduling.0.preemptible",
		"scheduling.0.node_affinities",
	}

	shieldedInstanceTemplateConfigKeys = []string{
		"shielded_instance_config.0.enable_secure_boot",
		"shielded_instance_config.0.enable_vtpm",
		"shielded_instance_config.0.enable_integrity_monitoring",
	}
)

func resourceComputeInstanceTemplate() *schema.Resource {
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
				ValidateFunc:  validateGCPName,
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
							Description: `The size of the image in gigabytes. If not specified, it will inherit the size of its base image. For SCRATCH disks, the size must be exactly 375GB.`,
						},

						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: `The GCE disk type. Can be either "pd-ssd", "local-ssd", "pd-balanced" or "pd-standard".`,
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

						"source_image": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `The image from which to initialize this disk. This can be one of: the image's self_link, projects/{project}/global/images/{image}, projects/{project}/global/images/family/{family}, global/images/{image}, global/images/family/{family}, family/{family}, {project}/{family}, {project}/{image}, {family}, or {image}. ~> Note: Either source or source_image is required when creating a new instance except for when creating a local SSD.`,
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
							Description: `The type of GCE disk, can be either "SCRATCH" or "PERSISTENT".`,
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
										DiffSuppressFunc: compareSelfLinkRelativePaths,
										Description:      `The self link of the encryption key that is stored in Google Cloud KMS.`,
									},
								},
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

			"enable_display": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Enable Virtual Displays on this instance. Note: allow_stopping_for_update must be set to true in order to update this field.`,
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
							DiffSuppressFunc: compareSelfLinkOrResourceName,
							Description:      `The name or self_link of the network to attach this interface to. Use network attribute for Legacy or Auto subnetted networks and subnetwork for custom subnetted networks.`,
						},

						"subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
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
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ForceNew:     true,
										Description:  `The networking tier used for configuring this instance template. This field can take the following values: PREMIUM or STANDARD. If this field is not specified, it is assumed to be PREMIUM.`,
										ValidateFunc: validation.StringInSlice([]string{"PREMIUM", "STANDARD"}, false),
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
										DiffSuppressFunc: ipCidrRangeDiffSuppress,
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
							DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
							Description:      `Specifies node affinities or anti-affinities to determine which sole-tenant nodes your instances and managed instance groups will use as host systems.`,
						},
					},
				},
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
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
									return canonicalizeServiceScope(v.(string))
								},
							},
							Set: stringScopeHashcode,
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
				DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
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
							DiffSuppressFunc: compareSelfLinkOrResourceName,
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
		},
	}
}

func resourceComputeInstanceTemplateSourceImageCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*Config)

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
			project, err := getProjectFromDiff(diff, config)
			if err != nil {
				return err
			}
			oldResolved, err := resolveImage(config, project, old.(string))
			if err != nil {
				return err
			}
			oldResolved, err = resolveImageRefToRelativeURI(project, oldResolved)
			if err != nil {
				return err
			}
			newResolved, err := resolveImage(config, project, new.(string))
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

func resourceComputeInstanceTemplateScratchDiskCustomizeDiffFunc(diff TerraformResourceDiff) error {
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
		if typee == "SCRATCH" && diskSize != 375 {
			return fmt.Errorf("SCRATCH disks must be exactly 375GB, disk %d is %d", i, diskSize)
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

func buildDisks(d *schema.ResourceData, config *Config) ([]*computeBeta.AttachedDisk, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	disksCount := d.Get("disk.#").(int)

	disks := make([]*computeBeta.AttachedDisk, 0, disksCount)
	for i := 0; i < disksCount; i++ {
		prefix := fmt.Sprintf("disk.%d", i)

		// Build the disk
		var disk computeBeta.AttachedDisk
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
			disk.DiskEncryptionKey = &computeBeta.CustomerEncryptionKey{}
			if v, ok := d.GetOk(prefix + ".disk_encryption_key.0.kms_key_self_link"); ok {
				disk.DiskEncryptionKey.KmsKeyName = v.(string)
			}
		}

		if v, ok := d.GetOk(prefix + ".source"); ok {
			disk.Source = v.(string)
			conflicts := []string{"disk_size_gb", "disk_name", "disk_type", "source_image", "labels"}
			for _, conflict := range conflicts {
				if _, ok := d.GetOk(prefix + "." + conflict); ok {
					return nil, fmt.Errorf("Cannot use `source` with any of the fields in %s", conflicts)
				}
			}
		} else {
			disk.InitializeParams = &computeBeta.AttachedDiskInitializeParams{}

			if v, ok := d.GetOk(prefix + ".disk_name"); ok {
				disk.InitializeParams.DiskName = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".disk_size_gb"); ok {
				disk.InitializeParams.DiskSizeGb = int64(v.(int))
			}
			disk.InitializeParams.DiskType = "pd-standard"
			if v, ok := d.GetOk(prefix + ".disk_type"); ok {
				disk.InitializeParams.DiskType = v.(string)
			}

			if v, ok := d.GetOk(prefix + ".source_image"); ok {
				imageName := v.(string)
				imageUrl, err := resolveImage(config, project, imageName)
				if err != nil {
					return nil, fmt.Errorf(
						"Error resolving image name '%s': %s",
						imageName, err)
				}
				disk.InitializeParams.SourceImage = imageUrl
			}

			disk.InitializeParams.Labels = expandStringMap(d, prefix+".labels")
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
func expandInstanceTemplateGuestAccelerators(d TerraformResourceData, config *Config) []*computeBeta.AcceleratorConfig {
	configs, ok := d.GetOk("guest_accelerator")
	if !ok {
		return nil
	}
	accels := configs.([]interface{})
	guestAccelerators := make([]*computeBeta.AcceleratorConfig, 0, len(accels))
	for _, raw := range accels {
		data := raw.(map[string]interface{})
		if data["count"].(int) == 0 {
			continue
		}
		guestAccelerators = append(guestAccelerators, &computeBeta.AcceleratorConfig{
			AcceleratorCount: int64(data["count"].(int)),
			// We can't use ParseAcceleratorFieldValue here because an instance
			// template does not have a zone we can use.
			AcceleratorType: data["type"].(string),
		})
	}

	return guestAccelerators
}

func resourceComputeInstanceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
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

	instanceProperties := &computeBeta.InstanceProperties{
		CanIpForward:           d.Get("can_ip_forward").(bool),
		Description:            d.Get("instance_description").(string),
		GuestAccelerators:      expandInstanceTemplateGuestAccelerators(d, config),
		MachineType:            d.Get("machine_type").(string),
		MinCpuPlatform:         d.Get("min_cpu_platform").(string),
		Disks:                  disks,
		Metadata:               metadata,
		NetworkInterfaces:      networks,
		Scheduling:             scheduling,
		ServiceAccounts:        expandServiceAccounts(d.Get("service_account").([]interface{})),
		Tags:                   resourceInstanceTags(d),
		ShieldedInstanceConfig: expandShieldedVmConfigs(d),
		DisplayDevice:          expandDisplayDevice(d),
	}

	if _, ok := d.GetOk("labels"); ok {
		instanceProperties.Labels = expandLabels(d)
	}

	var itName string
	if v, ok := d.GetOk("name"); ok {
		itName = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		itName = resource.PrefixedUniqueId(v.(string))
	} else {
		itName = resource.UniqueId()
	}
	instanceTemplate := &computeBeta.InstanceTemplate{
		Description: d.Get("description").(string),
		Properties:  instanceProperties,
		Name:        itName,
	}

	op, err := config.clientComputeBeta.InstanceTemplates.Insert(project, instanceTemplate).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance template: %s", err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("projects/%s/global/instanceTemplates/%s", project, instanceTemplate.Name))

	err = computeOperationWaitTime(config, op, project, "Creating Instance Template", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeInstanceTemplateRead(d, meta)
}

type diskCharacteristics struct {
	mode        string
	diskType    string
	diskSizeGb  string
	autoDelete  bool
	sourceImage string
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
	return dc
}

func flattenDisk(disk *computeBeta.AttachedDisk, defaultProject string) (map[string]interface{}, error) {
	diskMap := make(map[string]interface{})
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
		diskMap["disk_name"] = disk.InitializeParams.DiskName
		diskMap["disk_size_gb"] = disk.InitializeParams.DiskSizeGb
		diskMap["labels"] = disk.InitializeParams.Labels
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
	diskMap["source"] = ConvertSelfLinkToV1(disk.Source)
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

func flattenDisks(disks []*computeBeta.AttachedDisk, d *schema.ResourceData, defaultProject string) ([]map[string]interface{}, error) {
	apiDisks := make([]map[string]interface{}, len(disks))

	for i, disk := range disks {
		d, err := flattenDisk(disk, defaultProject)
		if err != nil {
			return nil, err
		}
		apiDisks[i] = d
	}

	return reorderDisks(d.Get("disk").([]interface{}), apiDisks), nil
}

func resourceComputeInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	splits := strings.Split(d.Id(), "/")
	instanceTemplate, err := config.clientComputeBeta.InstanceTemplates.Get(project, splits[len(splits)-1]).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Instance Template %q", d.Get("name").(string)))
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
			return fmt.Errorf("Error reading tags_fingerprint: %s", err)
		}
	}
	if instanceTemplate.Properties.Labels != nil {
		if err := d.Set("labels", instanceTemplate.Properties.Labels); err != nil {
			return fmt.Errorf("Error reading labels: %s", err)
		}
	}
	if err = d.Set("self_link", instanceTemplate.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
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
	if instanceTemplate.Properties.ShieldedVmConfig != nil {
		if err = d.Set("shielded_instance_config", flattenShieldedVmConfig(instanceTemplate.Properties.ShieldedInstanceConfig)); err != nil {
			return fmt.Errorf("Error setting shielded_instance_config: %s", err)
		}
	}

	if instanceTemplate.Properties.DisplayDevice != nil {
		if err = d.Set("enable_display", flattenEnableDisplay(instanceTemplate.Properties.DisplayDevice)); err != nil {
			return fmt.Errorf("Error setting enable_display: %s", err)
		}
	}
	return nil
}

func resourceComputeInstanceTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	splits := strings.Split(d.Id(), "/")
	op, err := config.clientCompute.InstanceTemplates.Delete(
		project, splits[len(splits)-1]).Do()
	if err != nil {
		return fmt.Errorf("Error deleting instance template: %s", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting Instance Template", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// This wraps the general compute instance helper expandScheduling.
// Default value of OnHostMaintenance depends on the value of Preemptible,
// so we can't set a default in schema
func expandResourceComputeInstanceTemplateScheduling(d *schema.ResourceData, meta interface{}) (*computeBeta.Scheduling, error) {
	v, ok := d.GetOk("scheduling")
	if !ok || v == nil {
		// We can't set defaults for lists (e.g. scheduling)
		return &computeBeta.Scheduling{
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
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/global/instanceTemplates/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/global/instanceTemplates/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
