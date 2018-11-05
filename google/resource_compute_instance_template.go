package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/googleapi"
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
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validateGCPName,
			},

			"name_prefix": &schema.Schema{
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

			"disk": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"boot": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"device_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"disk_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"disk_size_gb": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"disk_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"source_image": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"interface": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"source": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"type": &schema.Schema{
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

			"machine_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"automatic_restart": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Removed:  "Use 'scheduling.automatic_restart' instead.",
			},

			"can_ip_forward": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instance_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"metadata_startup_script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"metadata_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"address": &schema.Schema{
							Type:       schema.TypeString,
							Computed:   true, // Computed because it is set if network_ip is set.
							Optional:   true,
							ForceNew:   true,
							Deprecated: "Please use network_ip",
						},

						"network_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true, // Computed because it is set if address is set.
							Optional: true,
							ForceNew: true,
						},

						"subnetwork": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"subnetwork_project": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"access_config": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
									"network_tier": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PREMIUM", "STANDARD"}, false),
									},
									"assigned_nat_ip": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
										Removed:  "Use network_interface.access_config.nat_ip instead.",
									},
								},
							},
						},

						"alias_ip_range": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_cidr_range": &schema.Schema{
										Type:             schema.TypeString,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: ipCidrRangeDiffSuppress,
									},
									"subnetwork_range_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"on_host_maintenance": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Use 'scheduling.on_host_maintenance' instead.",
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"scheduling": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preemptible": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},

						"automatic_restart": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"on_host_maintenance": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_account": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"scopes": &schema.Schema{
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

			"guest_accelerator": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"type": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: linkDiffSuppress,
						},
					},
				},
			},

			"min_cpu_platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"labels": &schema.Schema{
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
	project, err := getProjectFromDiff(diff, config)
	if err != nil {
		return err
	}
	numDisks := diff.Get("disk.#").(int)
	for i := 0; i < numDisks; i++ {
		key := fmt.Sprintf("disk.%d.source_image", i)
		if diff.HasChange(key) {
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

	instanceProperties := &computeBeta.InstanceProperties{
		CanIpForward:   d.Get("can_ip_forward").(bool),
		Description:    d.Get("instance_description").(string),
		MachineType:    d.Get("machine_type").(string),
		MinCpuPlatform: d.Get("min_cpu_platform").(string),
	}

	disks, err := buildDisks(d, config)
	if err != nil {
		return err
	}
	instanceProperties.Disks = disks

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return err
	}
	instanceProperties.Metadata = metadata
	networks, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return err
	}
	instanceProperties.NetworkInterfaces = networks

	instanceProperties.Scheduling = &computeBeta.Scheduling{}
	instanceProperties.Scheduling.OnHostMaintenance = "MIGRATE"

	forceSendFieldsScheduling := make([]string, 0, 3)
	var hasSendMaintenance bool
	hasSendMaintenance = false
	if v, ok := d.GetOk("scheduling"); ok {
		_schedulings := v.([]interface{})
		if len(_schedulings) > 1 {
			return fmt.Errorf("Error, at most one `scheduling` block can be defined")
		}
		_scheduling := _schedulings[0].(map[string]interface{})

		// "automatic_restart" has a default value and is always safe to dereference
		automaticRestart := _scheduling["automatic_restart"].(bool)
		instanceProperties.Scheduling.AutomaticRestart = googleapi.Bool(automaticRestart)
		forceSendFieldsScheduling = append(forceSendFieldsScheduling, "AutomaticRestart")

		if vp, okp := _scheduling["on_host_maintenance"]; okp {
			instanceProperties.Scheduling.OnHostMaintenance = vp.(string)
			forceSendFieldsScheduling = append(forceSendFieldsScheduling, "OnHostMaintenance")
			hasSendMaintenance = true
		}

		if vp, okp := _scheduling["preemptible"]; okp {
			instanceProperties.Scheduling.Preemptible = vp.(bool)
			forceSendFieldsScheduling = append(forceSendFieldsScheduling, "Preemptible")
			if vp.(bool) && !hasSendMaintenance {
				instanceProperties.Scheduling.OnHostMaintenance = "TERMINATE"
				forceSendFieldsScheduling = append(forceSendFieldsScheduling, "OnHostMaintenance")
			}
		}
	}
	instanceProperties.Scheduling.ForceSendFields = forceSendFieldsScheduling

	instanceProperties.ServiceAccounts = expandServiceAccounts(d.Get("service_account").([]interface{}))

	instanceProperties.GuestAccelerators = expandInstanceTemplateGuestAccelerators(d, config)

	instanceProperties.Tags = resourceInstanceTags(d)
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
