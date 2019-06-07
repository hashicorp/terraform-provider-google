package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

func resourceComputeInstanceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceTemplateCreate,
		Read:   resourceComputeInstanceTemplateRead,
		Delete: resourceComputeInstanceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		CustomizeDiff: resourceComputeInstanceTemplateSourceImageCustomizeDiff,
		MigrateState:  resourceComputeInstanceTemplateMigrateState,

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
			},

			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"boot": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"device_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"disk_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"disk_size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"source_image": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"interface": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"disk_encryption_key": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_self_link": {
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkRelativePaths,
									},
								},
							},
						},
					},
				},
			},

			"machine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"automatic_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Removed:  "Use 'scheduling.automatic_restart' instead.",
			},

			"can_ip_forward": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instance_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"metadata_startup_script": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"metadata_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_interface": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"network_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"subnetwork_project": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"access_config": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
									"network_tier": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PREMIUM", "STANDARD"}, false),
									},
									"assigned_nat_ip": {
										Type:     schema.TypeString,
										Computed: true,
										Removed:  "Use network_interface.access_config.nat_ip instead.",
									},
								},
							},
						},

						"alias_ip_range": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_cidr_range": {
										Type:             schema.TypeString,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: ipCidrRangeDiffSuppress,
									},
									"subnetwork_range_name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},

						"address": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							Removed:  "Please use network_ip",
						},
					},
				},
			},

			"on_host_maintenance": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Use 'scheduling.on_host_maintenance' instead.",
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"scheduling": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preemptible": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},

						"automatic_restart": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"on_host_maintenance": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"node_affinities": {
							Type:             schema.TypeSet,
							Optional:         true,
							ForceNew:         true,
							Elem:             instanceSchedulingNodeAffinitiesElemSchema(),
							DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
						},
					},
				},
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_account": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"scopes": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				// Since this block is used by the API based on which
				// image being used, the field needs to be marked as Computed.
				Computed:         true,
				DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_secure_boot": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},

						"enable_vtpm": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"enable_integrity_monitoring": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},
					},
				},
			},

			"guest_accelerator": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: linkDiffSuppress,
						},
					},
				},
			},

			"min_cpu_platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceComputeInstanceTemplateSourceImageCustomizeDiff(diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*Config)

	numDisks := diff.Get("disk.#").(int)
	for i := 0; i < numDisks; i++ {
		key := fmt.Sprintf("disk.%d.source_image", i)
		if diff.HasChange(key) {
			// project must be retrieved once we know there is a diff to resolve, otherwise it will
			// attempt to retrieve project during `plan` before all calculated fields are ready
			// see https://github.com/terraform-providers/terraform-provider-google/issues/2878
			project, err := getProjectFromDiff(diff, config)
			if err != nil {
				return err
			}
			old, new := diff.GetChange(key)
			if old == "" || new == "" {
				// no sense in resolving empty strings
				err = diff.ForceNew(key)
				if err != nil {
					return err
				}
				continue
			}
			oldResolved, err := resolveImage(config, project, old.(string))
			if err != nil {
				return err
			}
			oldResolved, err = resolvedImageSelfLink(project, oldResolved)
			if err != nil {
				return err
			}
			newResolved, err := resolveImage(config, project, new.(string))
			if err != nil {
				return err
			}
			newResolved, err = resolvedImageSelfLink(project, newResolved)
			if err != nil {
				return err
			}
			if oldResolved != newResolved {
				err = diff.ForceNew(key)
				if err != nil {
					return err
				}
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
		CanIpForward:      d.Get("can_ip_forward").(bool),
		Description:       d.Get("instance_description").(string),
		GuestAccelerators: expandInstanceTemplateGuestAccelerators(d, config),
		MachineType:       d.Get("machine_type").(string),
		MinCpuPlatform:    d.Get("min_cpu_platform").(string),
		Disks:             disks,
		Metadata:          metadata,
		NetworkInterfaces: networks,
		Scheduling:        scheduling,
		ServiceAccounts:   expandServiceAccounts(d.Get("service_account").([]interface{})),
		Tags:              resourceInstanceTags(d),
		ShieldedVmConfig:  expandShieldedVmConfigs(d),
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
	d.SetId(instanceTemplate.Name)

	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating Instance Template")
	if err != nil {
		return err
	}

	return resourceComputeInstanceTemplateRead(d, meta)
}

func flattenDisks(disks []*computeBeta.AttachedDisk, d *schema.ResourceData, defaultProject string) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(disks))
	for _, disk := range disks {
		diskMap := make(map[string]interface{})
		if disk.InitializeParams != nil {
			if disk.InitializeParams.SourceImage != "" {
				selfLink, err := resolvedImageSelfLink(defaultProject, disk.InitializeParams.SourceImage)
				if err != nil {
					return nil, errwrap.Wrapf("Error expanding source image input to self_link: {{err}}", err)
				}
				path, err := getRelativePath(selfLink)
				if err != nil {
					return nil, errwrap.Wrapf("Error getting relative path for source image: {{err}}", err)
				}
				diskMap["source_image"] = path
			} else {
				diskMap["source_image"] = ""
			}
			diskMap["disk_type"] = disk.InitializeParams.DiskType
			diskMap["disk_name"] = disk.InitializeParams.DiskName
			diskMap["disk_size_gb"] = disk.InitializeParams.DiskSizeGb
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
		result = append(result, diskMap)
	}
	return result, nil
}

func resourceComputeInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceTemplate, err := config.clientComputeBeta.InstanceTemplates.Get(project, d.Id()).Do()
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
		d.Set("tags_fingerprint", "")
	}
	if instanceTemplate.Properties.Labels != nil {
		d.Set("labels", instanceTemplate.Properties.Labels)
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
		if err = d.Set("shielded_instance_config", flattenShieldedVmConfig(instanceTemplate.Properties.ShieldedVmConfig)); err != nil {
			return fmt.Errorf("Error setting shielded_instance_config: %s", err)
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

	op, err := config.clientCompute.InstanceTemplates.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting instance template: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Deleting Instance Template")
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
