package google

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mitchellh/hashstructure"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func resourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceCreate,
		Read:   resourceComputeInstanceRead,
		Update: resourceComputeInstanceUpdate,
		Delete: resourceComputeInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeInstanceImportState,
		},

		SchemaVersion: 6,
		MigrateState:  resourceComputeInstanceMigrateState,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},

		// A compute instance is more or less a superset of a compute instance
		// template. Please attempt to maintain consistency with the
		// resource_compute_instance_template schema when updating this one.
		Schema: map[string]*schema.Schema{
			"boot_disk": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"device_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"disk_encryption_key_raw": &schema.Schema{
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"disk_encryption_key_sha256": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"initialize_params": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": &schema.Schema{
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},

									"type": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd"}, false),
									},

									"image": &schema.Schema{
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										ForceNew:         true,
										DiffSuppressFunc: diskImageDiffSuppress,
									},
								},
							},
						},

						"source": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    []string{"boot_disk.initialize_params"},
							DiffSuppressFunc: linkDiffSuppress,
						},
					},
				},
			},

			"machine_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_interface": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"subnetwork": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"subnetwork_project": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"address": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							ForceNew:   true,
							Computed:   true,
							Deprecated: "Please use network_ip",
						},

						"network_ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},

						"access_config": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"network_tier": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PREMIUM", "STANDARD"}, false),
									},

									// It's unclear why this field exists, as
									// nat_ip can be both optional and computed.
									// Consider deprecating it.
									"assigned_nat_ip": &schema.Schema{
										Type:       schema.TypeString,
										Computed:   true,
										Deprecated: "Use network_interface.access_config.nat_ip instead.",
									},

									"public_ptr_domain_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"alias_ip_range": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_cidr_range": &schema.Schema{
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: ipCidrRangeDiffSuppress,
									},
									"subnetwork_range_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"allow_stopping_for_update": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"attached_disk": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},

						"device_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"mode": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "READ_WRITE",
							ValidateFunc: validation.StringInSlice([]string{"READ_WRITE", "READ_ONLY"}, false),
						},

						"disk_encryption_key_raw": &schema.Schema{
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"disk_encryption_key_sha256": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"can_ip_forward": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"create_timeout": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Default:    4,
				Deprecated: "Use timeouts block instead.",
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"deletion_protection": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"disk": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Removed:  "Use boot_disk, scratch_disk, and attached_disk instead",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// TODO(mitchellh): one of image or disk is required

						"disk": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"image": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"scratch": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},

						"auto_delete": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},

						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"device_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"disk_encryption_key_raw": &schema.Schema{
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"disk_encryption_key_sha256": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"guest_accelerator": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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

			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"metadata_startup_script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"min_cpu_platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"network": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Removed:  "Please use network_interface",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"address": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"internal_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"external_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"scheduling": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_host_maintenance": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"automatic_restart": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"preemptible": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
					},
				},
			},

			"scratch_disk": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SCSI",
							ValidateFunc: validation.StringInSlice([]string{"SCSI", "NVME"}, false),
						},
					},
				},
			},

			"service_account": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"scopes": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
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

			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"cpu_platform": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"label_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.If(
				func(d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange("guest_accelerator")
				},
				suppressEmptyGuestAcceleratorDiff,
			),
		),
	}
}

func getInstance(config *Config, d *schema.ResourceData) (*computeBeta.Instance, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}
	zone, err := getZone(d, config)
	if err != nil {
		return nil, err
	}
	instance, err := config.clientComputeBeta.Instances.Get(project, zone, d.Id()).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Instance %s", d.Get("name").(string)))
	}
	return instance, nil
}

func getDisk(diskUri string, d *schema.ResourceData, config *Config) (*compute.Disk, error) {
	source, err := ParseDiskFieldValue(diskUri, d, config)
	if err != nil {
		return nil, err
	}

	disk, err := config.clientCompute.Disks.Get(source.Project, source.Zone, source.Name).Do()
	if err != nil {
		return nil, err
	}

	return disk, err
}

func expandComputeInstance(project string, zone *compute.Zone, d *schema.ResourceData, config *Config) (*computeBeta.Instance, error) {
	// Get the machine type
	var machineTypeUrl string
	if mt, ok := d.GetOk("machine_type"); ok {
		log.Printf("[DEBUG] Loading machine type: %s", mt.(string))
		machineType, err := config.clientCompute.MachineTypes.Get(
			project, zone.Name, mt.(string)).Do()
		if err != nil {
			return nil, fmt.Errorf(
				"Error loading machine type: %s",
				err)
		}
		machineTypeUrl = machineType.SelfLink
	}

	// Build up the list of disks

	disks := []*computeBeta.AttachedDisk{}
	if _, hasBootDisk := d.GetOk("boot_disk"); hasBootDisk {
		bootDisk, err := expandBootDisk(d, config, zone, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, bootDisk)
	}

	if _, hasScratchDisk := d.GetOk("scratch_disk"); hasScratchDisk {
		scratchDisks, err := expandScratchDisks(d, config, zone, project)
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

	prefix := "scheduling.0"
	scheduling := &computeBeta.Scheduling{
		AutomaticRestart:  googleapi.Bool(d.Get(prefix + ".automatic_restart").(bool)),
		Preemptible:       d.Get(prefix + ".preemptible").(bool),
		OnHostMaintenance: d.Get(prefix + ".on_host_maintenance").(string),
		ForceSendFields:   []string{"AutomaticRestart", "Preemptible"},
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating metadata: %s", err)
	}

	networkInterfaces, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network interfaces: %s", err)
	}

	accels, err := expandInstanceGuestAccelerators(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating guest accelerators: %s", err)
	}

	// Create the instance information
	return &computeBeta.Instance{
		CanIpForward:       d.Get("can_ip_forward").(bool),
		Description:        d.Get("description").(string),
		Disks:              disks,
		MachineType:        machineTypeUrl,
		Metadata:           metadata,
		Name:               d.Get("name").(string),
		NetworkInterfaces:  networkInterfaces,
		Tags:               resourceInstanceTags(d),
		Labels:             expandLabels(d),
		ServiceAccounts:    expandServiceAccounts(d.Get("service_account").([]interface{})),
		GuestAccelerators:  accels,
		MinCpuPlatform:     d.Get("min_cpu_platform").(string),
		Scheduling:         scheduling,
		DeletionProtection: d.Get("deletion_protection").(bool),
		ForceSendFields:    []string{"CanIpForward", "DeletionProtection"},
	}, nil
}

func resourceComputeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	z, err := getZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Loading zone: %s", z)
	zone, err := config.clientCompute.Zones.Get(
		project, z).Do()
	if err != nil {
		return fmt.Errorf("Error loading zone '%s': %s", z, err)
	}

	instance, err := expandComputeInstance(project, zone, d, config)
	if err != nil {
		return err
	}

	// Read create timeout
	// Until "create_timeout" is removed, use that timeout if set.
	createTimeout := int(d.Timeout(schema.TimeoutCreate).Minutes())
	if v, ok := d.GetOk("create_timeout"); ok && v != 4 {
		createTimeout = v.(int)
	}

	log.Printf("[INFO] Requesting instance creation")
	op, err := config.clientComputeBeta.Instances.Insert(project, zone.Name, instance).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	// Store the ID now
	d.SetId(instance.Name)

	// Wait for the operation to complete
	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project, createTimeout, "instance to create")
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance, err := getInstance(config, d)
	if err != nil || instance == nil {
		return err
	}

	md := flattenMetadataBeta(instance.Metadata)
	existingMetadata := d.Get("metadata").(map[string]interface{})

	// If the existing config specifies "metadata.startup-script" instead of "metadata_startup_script",
	// we shouldn't move the remote metadata.startup-script to metadata_startup_script.  Otherwise,
	// we should.
	if ss, ok := existingMetadata["startup-script"]; !ok || ss == "" {
		d.Set("metadata_startup_script", md["startup-script"])
		// Note that here we delete startup-script from our metadata list. This is to prevent storing the startup-script
		// as a value in the metadata since the config specifically tracks it under 'metadata_startup_script'
		delete(md, "startup-script")
	} else if _, ok := d.GetOk("metadata_startup_script"); ok {
		delete(md, "startup-script")
	}

	// Delete any keys not explicitly set in our config file
	for k := range md {
		if _, ok := existingMetadata[k]; !ok {
			delete(md, k)
		}
	}

	if err = d.Set("metadata", md); err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	d.Set("can_ip_forward", instance.CanIpForward)
	d.Set("machine_type", GetResourceNameFromSelfLink(instance.MachineType))

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

	// Set the metadata fingerprint if there is one.
	if instance.Metadata != nil {
		d.Set("metadata_fingerprint", instance.Metadata.Fingerprint)
	}

	// Set the tags fingerprint if there is one.
	if instance.Tags != nil {
		d.Set("tags_fingerprint", instance.Tags.Fingerprint)
		d.Set("tags", convertStringArrToInterface(instance.Tags.Items))
	}

	if err := d.Set("labels", instance.Labels); err != nil {
		return err
	}

	if instance.LabelFingerprint != "" {
		d.Set("label_fingerprint", instance.LabelFingerprint)
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
			source, err := ParseRegionDiskFieldValue(disk["source"].(string), d, config)
			if err != nil {
				return err
			}
			sourceLink = source.RelativeLink()
		} else {
			source, err := ParseDiskFieldValue(disk["source"].(string), d, config)
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
			d.Set("boot_disk", flattenBootDisk(d, disk, config))
		} else if disk.Type == "SCRATCH" {
			scratchDisks = append(scratchDisks, flattenScratchDisk(disk))
		} else {
			var sourceLink string
			if strings.Contains(disk.Source, "regions/") {
				source, err := ParseRegionDiskFieldValue(disk.Source, d, config)
				if err != nil {
					return err
				}
				sourceLink = source.RelativeLink()
			} else {
				source, err := ParseDiskFieldValue(disk.Source, d, config)
				if err != nil {
					return err
				}
				sourceLink = source.RelativeLink()
			}
			adIndex, inConfig := attachedDiskSources[sourceLink]
			di := map[string]interface{}{
				"source":      ConvertSelfLinkToV1(disk.Source),
				"device_name": disk.DeviceName,
				"mode":        disk.Mode,
			}
			if key := disk.DiskEncryptionKey; key != nil {
				if inConfig {
					di["disk_encryption_key_raw"] = d.Get(fmt.Sprintf("attached_disk.%d.disk_encryption_key_raw", adIndex))
				}
				di["disk_encryption_key_sha256"] = key.Sha256
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
	// Remove nils from map in case there were disks in the config that were not present on read;
	// i.e. a disk was detached out of band
	ads := []map[string]interface{}{}
	for _, d := range attachedDisks {
		if d != nil {
			ads = append(ads, d)
		}
	}

	d.Set("service_account", flattenServiceAccounts(instance.ServiceAccounts))
	d.Set("attached_disk", ads)
	d.Set("scratch_disk", scratchDisks)
	d.Set("scheduling", flattenScheduling(instance.Scheduling))
	d.Set("guest_accelerator", flattenGuestAccelerators(instance.GuestAccelerators))
	d.Set("cpu_platform", instance.CpuPlatform)
	d.Set("min_cpu_platform", instance.MinCpuPlatform)
	d.Set("deletion_protection", instance.DeletionProtection)
	d.Set("self_link", ConvertSelfLinkToV1(instance.SelfLink))
	d.Set("instance_id", fmt.Sprintf("%d", instance.Id))
	d.Set("project", project)
	d.Set("zone", GetResourceNameFromSelfLink(instance.Zone))
	d.Set("name", instance.Name)
	d.SetId(instance.Name)

	return nil
}

func resourceComputeInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	// Use beta api directly in order to read network_interface.fingerprint without having to put it in the schema.
	// Change back to getInstance(config, d) once updating alias ips is GA.
	instance, err := config.clientComputeBeta.Instances.Get(project, zone, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Instance %s", d.Get("name").(string)))
	}

	// Enable partial mode for the resource since it is possible
	d.Partial(true)

	// If the Metadata has changed, then update that.
	if d.HasChange("metadata") {
		o, n := d.GetChange("metadata")
		if script, scriptExists := d.GetOk("metadata_startup_script"); scriptExists && script != "" {
			if _, ok := n.(map[string]interface{})["startup-script"]; ok {
				return fmt.Errorf("Only one of metadata.startup-script and metadata_startup_script may be defined")
			}

			if err = d.Set("metadata", n); err != nil {
				return err
			}
			n.(map[string]interface{})["startup-script"] = script
		}

		updateMD := func() error {
			// Reload the instance in the case of a fingerprint mismatch
			instance, err = getInstance(config, d)
			if err != nil {
				return err
			}

			md := instance.Metadata

			BetaMetadataUpdate(o.(map[string]interface{}), n.(map[string]interface{}), md)

			if err != nil {
				return fmt.Errorf("Error updating metadata: %s", err)
			}

			mdV1 := &compute.Metadata{}
			err = Convert(md, mdV1)
			if err != nil {
				return err
			}

			op, err := config.clientCompute.Instances.SetMetadata(
				project, zone, d.Id(), mdV1).Do()
			if err != nil {
				return fmt.Errorf("Error updating metadata: %s", err)
			}

			opErr := computeOperationWaitTime(config.clientCompute, op, project, "metadata to update", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
			if opErr != nil {
				return opErr
			}

			d.SetPartial("metadata")

			return nil
		}

		MetadataRetryWrapper(updateMD)

	}

	if d.HasChange("tags") {
		tags := resourceInstanceTags(d)
		tagsV1 := &compute.Tags{}
		if err := Convert(tags, tagsV1); err != nil {
			return err
		}
		op, err := config.clientCompute.Instances.SetTags(
			project, zone, d.Id(), tagsV1).Do()
		if err != nil {
			return fmt.Errorf("Error updating tags: %s", err)
		}

		opErr := computeOperationWaitTime(config.clientCompute, op, project, "tags to update", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if opErr != nil {
			return opErr
		}

		d.SetPartial("tags")
	}

	if d.HasChange("labels") {
		labels := expandLabels(d)
		labelFingerprint := d.Get("label_fingerprint").(string)
		req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: labelFingerprint}

		op, err := config.clientCompute.Instances.SetLabels(project, zone, d.Id(), &req).Do()
		if err != nil {
			return fmt.Errorf("Error updating labels: %s", err)
		}

		opErr := computeOperationWaitTime(config.clientCompute, op, project, "labels to update", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if opErr != nil {
			return opErr
		}

		d.SetPartial("labels")
	}

	if d.HasChange("scheduling") {
		prefix := "scheduling.0"
		scheduling := &compute.Scheduling{
			AutomaticRestart:  googleapi.Bool(d.Get(prefix + ".automatic_restart").(bool)),
			Preemptible:       d.Get(prefix + ".preemptible").(bool),
			OnHostMaintenance: d.Get(prefix + ".on_host_maintenance").(string),
			ForceSendFields:   []string{"AutomaticRestart", "Preemptible"},
		}

		op, err := config.clientCompute.Instances.SetScheduling(project,
			zone, d.Id(), scheduling).Do()

		if err != nil {
			return fmt.Errorf("Error updating scheduling policy: %s", err)
		}

		opErr := computeOperationWaitTime(config.clientCompute, op, project, "scheduling policy update", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if opErr != nil {
			return opErr
		}

		d.SetPartial("scheduling")
	}

	networkInterfacesCount := d.Get("network_interface.#").(int)
	// Sanity check
	if networkInterfacesCount != len(instance.NetworkInterfaces) {
		return fmt.Errorf("Instance had unexpected number of network interfaces: %d", len(instance.NetworkInterfaces))
	}
	for i := 0; i < networkInterfacesCount; i++ {
		prefix := fmt.Sprintf("network_interface.%d", i)
		instNetworkInterface := instance.NetworkInterfaces[i]
		networkName := d.Get(prefix + ".name").(string)

		// TODO: This sanity check is broken by #929, disabled for now (by forcing the equality)
		networkName = instNetworkInterface.Name
		// Sanity check
		if networkName != instNetworkInterface.Name {
			return fmt.Errorf("Instance networkInterface had unexpected name: %s", instNetworkInterface.Name)
		}

		if d.HasChange(prefix + ".access_config") {

			// TODO: This code deletes then recreates accessConfigs.  This is bad because it may
			// leave the machine inaccessible from either ip if the creation part fails (network
			// timeout etc).  However right now there is a GCE limit of 1 accessConfig so it is
			// the only way to do it.  In future this should be revised to only change what is
			// necessary, and also add before removing.

			// Delete any accessConfig that currently exists in instNetworkInterface
			for _, ac := range instNetworkInterface.AccessConfigs {
				op, err := config.clientCompute.Instances.DeleteAccessConfig(
					project, zone, d.Id(), ac.Name, networkName).Do()
				if err != nil {
					return fmt.Errorf("Error deleting old access_config: %s", err)
				}
				opErr := computeOperationWaitTime(config.clientCompute, op, project, "old access_config to delete", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
				if opErr != nil {
					return opErr
				}
			}

			// Create new ones
			accessConfigsCount := d.Get(prefix + ".access_config.#").(int)
			for j := 0; j < accessConfigsCount; j++ {
				acPrefix := fmt.Sprintf("%s.access_config.%d", prefix, j)
				ac := &computeBeta.AccessConfig{
					Type:        "ONE_TO_ONE_NAT",
					NatIP:       d.Get(acPrefix + ".nat_ip").(string),
					NetworkTier: d.Get(acPrefix + ".network_tier").(string),
				}
				if ptr, ok := d.GetOk(acPrefix + ".public_ptr_domain_name"); ok && ptr != "" {
					ac.SetPublicPtr = true
					ac.PublicPtrDomainName = ptr.(string)
				}

				op, err := config.clientComputeBeta.Instances.AddAccessConfig(
					project, zone, d.Id(), networkName, ac).Do()
				if err != nil {
					return fmt.Errorf("Error adding new access_config: %s", err)
				}
				opErr := computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "new access_config to add")
				if opErr != nil {
					return opErr
				}
			}
		}

		if d.HasChange(prefix + ".alias_ip_range") {
			rereadFingerprint := false

			// Alias IP ranges cannot be updated; they must be removed and then added.
			if len(instNetworkInterface.AliasIpRanges) > 0 {
				ni := &computeBeta.NetworkInterface{
					Fingerprint:     instNetworkInterface.Fingerprint,
					ForceSendFields: []string{"AliasIpRanges"},
				}
				op, err := config.clientComputeBeta.Instances.UpdateNetworkInterface(project, zone, d.Id(), networkName, ni).Do()
				if err != nil {
					return errwrap.Wrapf("Error removing alias_ip_range: {{err}}", err)
				}
				opErr := computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "updaing alias ip ranges")
				if opErr != nil {
					return opErr
				}
				rereadFingerprint = true
			}

			ranges := d.Get(prefix + ".alias_ip_range").([]interface{})
			if len(ranges) > 0 {
				if rereadFingerprint {
					instance, err = config.clientComputeBeta.Instances.Get(project, zone, d.Id()).Do()
					if err != nil {
						return err
					}
					instNetworkInterface = instance.NetworkInterfaces[i]
				}
				ni := &computeBeta.NetworkInterface{
					AliasIpRanges: expandAliasIpRanges(ranges),
					Fingerprint:   instNetworkInterface.Fingerprint,
				}
				op, err := config.clientComputeBeta.Instances.UpdateNetworkInterface(project, zone, d.Id(), networkName, ni).Do()
				if err != nil {
					return errwrap.Wrapf("Error adding alias_ip_range: {{err}}", err)
				}
				opErr := computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "updaing alias ip ranges")
				if opErr != nil {
					return opErr
				}
			}
		}
		d.SetPartial("network_interface")
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
				err = Convert(computeDisk, computeDiskV1)
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
				op, err := config.clientCompute.Instances.DetachDisk(project, zone, instance.Name, deviceName).Do()
				if err != nil {
					return errwrap.Wrapf("Error detaching disk: %s", err)
				}

				opErr := computeOperationWaitTime(config.clientCompute, op, project, "detaching disk", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
				if opErr != nil {
					return opErr
				}
				log.Printf("[DEBUG] Successfully detached disk %s", deviceName)
			}
		}

		// Attach the new disks
		for _, disk := range attach {
			op, err := config.clientCompute.Instances.AttachDisk(project, zone, instance.Name, disk).Do()
			if err != nil {
				return errwrap.Wrapf("Error attaching disk : {{err}}", err)
			}

			opErr := computeOperationWaitTime(config.clientCompute, op, project, "attaching disk", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
			if opErr != nil {
				return opErr
			}
			log.Printf("[DEBUG] Successfully attached disk %s", disk.Source)
		}

		d.SetPartial("attached_disk")
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

		op, err := config.clientCompute.Instances.SetDeletionProtection(project, zone, d.Id()).DeletionProtection(nDeletionProtection).Do()
		if err != nil {
			return fmt.Errorf("Error updating deletion protection flag: %s", err)
		}

		opErr := computeOperationWaitTime(config.clientCompute, op, project, "deletion protection to update", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if opErr != nil {
			return opErr
		}

		d.SetPartial("deletion_protection")
	}

	// Attributes which can only be changed if the instance is stopped
	if scopesChange || d.HasChange("service_account.0.email") || d.HasChange("machine_type") || d.HasChange("min_cpu_platform") {
		if !d.Get("allow_stopping_for_update").(bool) {
			return fmt.Errorf("Changing the machine_type, min_cpu_platform, or service_account on an instance requires stopping it. " +
				"To acknowledge this, please set allow_stopping_for_update = true in your config.")
		}
		op, err := config.clientCompute.Instances.Stop(project, zone, instance.Name).Do()
		if err != nil {
			return errwrap.Wrapf("Error stopping instance: {{err}}", err)
		}

		opErr := computeOperationWaitTime(config.clientCompute, op, project, "stopping instance", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if opErr != nil {
			return opErr
		}

		if d.HasChange("machine_type") {
			mt, err := ParseMachineTypesFieldValue(d.Get("machine_type").(string), d, config)
			if err != nil {
				return err
			}
			req := &compute.InstancesSetMachineTypeRequest{
				MachineType: mt.RelativeLink(),
			}
			op, err = config.clientCompute.Instances.SetMachineType(project, zone, instance.Name, req).Do()
			if err != nil {
				return err
			}
			opErr := computeOperationWaitTime(config.clientCompute, op, project, "updating machinetype", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
			if opErr != nil {
				return opErr
			}
			d.SetPartial("machine_type")
		}

		if d.HasChange("min_cpu_platform") {
			minCpuPlatform, ok := d.GetOk("min_cpu_platform")
			// Even though you don't have to set minCpuPlatform on create, you do have to set it to an
			// actual value on update. "Automatic" is the default. This will be read back from the API as empty,
			// so we don't need to worry about diffs.
			if !ok {
				minCpuPlatform = "Automatic"
			}
			req := &compute.InstancesSetMinCpuPlatformRequest{
				MinCpuPlatform: minCpuPlatform.(string),
			}
			op, err = config.clientCompute.Instances.SetMinCpuPlatform(project, zone, instance.Name, req).Do()
			if err != nil {
				return err
			}
			opErr := computeOperationWaitTime(config.clientCompute, op, project, "updating min cpu platform", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
			if opErr != nil {
				return opErr
			}
			d.SetPartial("min_cpu_platform")
		}

		if d.HasChange("service_account.0.email") || scopesChange {
			sa := d.Get("service_account").([]interface{})
			req := &compute.InstancesSetServiceAccountRequest{ForceSendFields: []string{"email"}}
			if len(sa) > 0 && sa[0] != nil {
				saMap := sa[0].(map[string]interface{})
				req.Email = saMap["email"].(string)
				req.Scopes = canonicalizeServiceScopes(convertStringSet(saMap["scopes"].(*schema.Set)))
			}
			op, err = config.clientCompute.Instances.SetServiceAccount(project, zone, instance.Name, req).Do()
			if err != nil {
				return err
			}
			opErr := computeOperationWaitTime(config.clientCompute, op, project, "updating service account", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
			if opErr != nil {
				return opErr
			}
			d.SetPartial("service_account")
		}

		op, err = config.clientCompute.Instances.Start(project, zone, instance.Name).Do()
		if err != nil {
			return errwrap.Wrapf("Error starting instance: {{err}}", err)
		}

		opErr = computeOperationWaitTime(config.clientCompute, op, project, "starting instance", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if opErr != nil {
			return opErr
		}
	}

	// We made it, disable partial mode
	d.Partial(false)

	return resourceComputeInstanceRead(d, meta)
}

func expandAttachedDisk(diskConfig map[string]interface{}, d *schema.ResourceData, meta interface{}) (*computeBeta.AttachedDisk, error) {
	config := meta.(*Config)

	s := diskConfig["source"].(string)
	var sourceLink string
	if strings.Contains(s, "regions/") {
		source, err := ParseRegionDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	} else {
		source, err := ParseDiskFieldValue(s, d, config)
		if err != nil {
			return nil, err
		}
		sourceLink = source.RelativeLink()
	}

	disk := &computeBeta.AttachedDisk{
		Source: sourceLink,
	}

	if v, ok := diskConfig["mode"]; ok {
		disk.Mode = v.(string)
	}

	if v, ok := diskConfig["device_name"]; ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := diskConfig["disk_encryption_key_raw"]; ok {
		disk.DiskEncryptionKey = &computeBeta.CustomerEncryptionKey{
			RawKey: v.(string),
		}
	}
	return disk, nil
}

// See comment on expandInstanceTemplateGuestAccelerators regarding why this
// code is duplicated.
func expandInstanceGuestAccelerators(d TerraformResourceData, config *Config) ([]*computeBeta.AcceleratorConfig, error) {
	configs, ok := d.GetOk("guest_accelerator")
	if !ok {
		return nil, nil
	}
	accels := configs.([]interface{})
	guestAccelerators := make([]*computeBeta.AcceleratorConfig, 0, len(accels))
	for _, raw := range accels {
		data := raw.(map[string]interface{})
		if data["count"].(int) == 0 {
			continue
		}
		at, err := ParseAcceleratorFieldValue(data["type"].(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot parse accelerator type: %v", err)
		}
		guestAccelerators = append(guestAccelerators, &computeBeta.AcceleratorConfig{
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
// 		guest_accelerators {
//      	count = "${var.enable_gpu ? var.gpu_count : 0}"
//          ...
// 		}
// After reconciling the desired and actual state, we would otherwise see a
// perpetual resembling:
// 		[] != [{"count":0, "type": "nvidia-tesla-k80"}]
func suppressEmptyGuestAcceleratorDiff(d *schema.ResourceDiff, meta interface{}) error {
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

func resourceComputeInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Requesting instance deletion: %s", d.Id())

	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("Cannot delete instance %s: instance Deletion Protection is enabled. Set deletion_protection to false for this resource and run \"terraform apply\" before attempting to delete it.", d.Id())
	} else {
		op, err := config.clientCompute.Instances.Delete(project, zone, d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error deleting instance: %s", err)
		}

		// Wait for the operation to complete
		opErr := computeOperationWaitTime(config.clientCompute, op, project, "instance to delete", int(d.Timeout(schema.TimeoutDelete).Minutes()))
		if opErr != nil {
			return opErr
		}

		d.SetId("")
		return nil
	}
}

func resourceComputeInstanceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {project}/{zone}/{instance_name}", d.Id())
	}

	d.Set("project", parts[0])
	d.Set("zone", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func expandBootDisk(d *schema.ResourceData, config *Config, zone *compute.Zone, project string) (*computeBeta.AttachedDisk, error) {
	disk := &computeBeta.AttachedDisk{
		AutoDelete: d.Get("boot_disk.0.auto_delete").(bool),
		Boot:       true,
	}

	if v, ok := d.GetOk("boot_disk.0.device_name"); ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_raw"); ok {
		disk.DiskEncryptionKey = &computeBeta.CustomerEncryptionKey{
			RawKey: v.(string),
		}
	}

	if v, ok := d.GetOk("boot_disk.0.source"); ok {
		source, err := ParseDiskFieldValue(v.(string), d, config)
		if err != nil {
			return nil, err
		}
		disk.Source = source.RelativeLink()
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params"); ok {
		disk.InitializeParams = &computeBeta.AttachedDiskInitializeParams{}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.size"); ok {
			disk.InitializeParams.DiskSizeGb = int64(v.(int))
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.type"); ok {
			diskTypeName := v.(string)
			diskType, err := readDiskType(config, zone, project, diskTypeName)
			if err != nil {
				return nil, fmt.Errorf("Error loading disk type '%s': %s", diskTypeName, err)
			}
			disk.InitializeParams.DiskType = diskType.SelfLink
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.image"); ok {
			imageName := v.(string)
			imageUrl, err := resolveImage(config, project, imageName)
			if err != nil {
				return nil, fmt.Errorf("Error resolving image name '%s': %s", imageName, err)
			}

			disk.InitializeParams.SourceImage = imageUrl
		}
	}

	return disk, nil
}

func flattenBootDisk(d *schema.ResourceData, disk *computeBeta.AttachedDisk, config *Config) []map[string]interface{} {
	result := map[string]interface{}{
		"auto_delete": disk.AutoDelete,
		"device_name": disk.DeviceName,
		"source":      ConvertSelfLinkToV1(disk.Source),
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
			"type": GetResourceNameFromSelfLink(diskDetails.Type),
			// If the config specifies a family name that doesn't match the image name, then
			// the diff won't be properly suppressed. See DiffSuppressFunc for this field.
			"image": diskDetails.SourceImage,
			"size":  diskDetails.SizeGb,
		}}
	}

	if disk.DiskEncryptionKey != nil {
		result["disk_encryption_key_sha256"] = disk.DiskEncryptionKey.Sha256
	}

	return []map[string]interface{}{result}
}

func expandScratchDisks(d *schema.ResourceData, config *Config, zone *compute.Zone, project string) ([]*computeBeta.AttachedDisk, error) {
	diskType, err := readDiskType(config, zone, project, "local-ssd")
	if err != nil {
		return nil, fmt.Errorf("Error loading disk type 'local-ssd': %s", err)
	}

	n := d.Get("scratch_disk.#").(int)
	scratchDisks := make([]*computeBeta.AttachedDisk, 0, n)
	for i := 0; i < n; i++ {
		scratchDisks = append(scratchDisks, &computeBeta.AttachedDisk{
			AutoDelete: true,
			Type:       "SCRATCH",
			Interface:  d.Get(fmt.Sprintf("scratch_disk.%d.interface", i)).(string),
			InitializeParams: &computeBeta.AttachedDiskInitializeParams{
				DiskType: diskType.SelfLink,
			},
		})
	}

	return scratchDisks, nil
}

func flattenScratchDisk(disk *computeBeta.AttachedDisk) map[string]interface{} {
	result := map[string]interface{}{
		"interface": disk.Interface,
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
