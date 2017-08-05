package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"regexp"
)

func stringScopeHashcode(v interface{}) int {
	v = canonicalizeServiceScope(v.(string))
	return schema.HashString(v)
}

func resourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceCreate,
		Read:   resourceComputeInstanceRead,
		Update: resourceComputeInstanceUpdate,
		Delete: resourceComputeInstanceDelete,

		SchemaVersion: 2,
		MigrateState:  resourceComputeInstanceMigrateState,

		Schema: map[string]*schema.Schema{
			"boot_disk": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
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
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 {
												errors = append(errors, fmt.Errorf(
													"%q must be greater than 0", k))
											}
											return
										},
									},

									"type": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd"}, false),
									},

									"image": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},

						"source": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: []string{"boot_disk.initialize_params"},
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

			"disk": &schema.Schema{
				Type:       schema.TypeList,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Use boot_disk, scratch_disk, and attached_disk instead",
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

			// Preferred way of adding persistent disks to an instance.
			// Use this instead of `disk` when possible.
			"attached_disk": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true, // TODO(danawillow): Remove this, support attaching/detaching
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"device_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"disk_encryption_key_raw": &schema.Schema{
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},

						"disk_encryption_key_sha256": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"machine_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
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
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: linkDiffSuppress,
						},

						"subnetwork": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: linkDiffSuppress,
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
									},

									"assigned_nat_ip": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"network": &schema.Schema{
				Type:       schema.TypeList,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use network_interface",
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
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

			"service_account": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
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

			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"label_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
		},
	}
}

func getInstance(config *Config, d *schema.ResourceData) (*compute.Instance, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	instance, err := config.clientCompute.Instances.Get(
		project, d.Get("zone").(string), d.Id()).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Instance %s", d.Get("name").(string)))
	}

	return instance, nil
}

func resourceComputeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	log.Printf("[DEBUG] Loading zone: %s", d.Get("zone").(string))
	zone, err := config.clientCompute.Zones.Get(
		project, d.Get("zone").(string)).Do()
	if err != nil {
		return fmt.Errorf(
			"Error loading zone '%s': %s", d.Get("zone").(string), err)
	}

	// Get the machine type
	log.Printf("[DEBUG] Loading machine type: %s", d.Get("machine_type").(string))
	machineType, err := config.clientCompute.MachineTypes.Get(
		project, zone.Name, d.Get("machine_type").(string)).Do()
	if err != nil {
		return fmt.Errorf(
			"Error loading machine type: %s",
			err)
	}

	// Build up the list of disks

	disks := []*compute.AttachedDisk{}
	var hasBootDisk bool
	if _, hasBootDisk = d.GetOk("boot_disk"); hasBootDisk {
		bootDisk, err := expandBootDisk(d, config, zone, project)
		if err != nil {
			return err
		}
		disks = append(disks, bootDisk)
	}

	var hasScratchDisk bool
	if _, hasScratchDisk := d.GetOk("scratch_disk"); hasScratchDisk {
		scratchDisks, err := expandScratchDisks(d, config, zone)
		if err != nil {
			return err
		}
		disks = append(disks, scratchDisks...)
	}

	disksCount := d.Get("disk.#").(int)
	attachedDisksCount := d.Get("attached_disk.#").(int)

	if disksCount+attachedDisksCount == 0 && !hasBootDisk {
		return fmt.Errorf("At least one disk, attached_disk, or boot_disk must be set")
	}
	for i := 0; i < disksCount; i++ {
		prefix := fmt.Sprintf("disk.%d", i)

		// var sourceLink string

		// Build the disk
		var disk compute.AttachedDisk
		disk.Type = "PERSISTENT"
		disk.Mode = "READ_WRITE"
		disk.Boot = i == 0 && !hasBootDisk
		disk.AutoDelete = d.Get(prefix + ".auto_delete").(bool)

		if _, ok := d.GetOk(prefix + ".disk"); ok {
			if _, ok := d.GetOk(prefix + ".type"); ok {
				return fmt.Errorf(
					"Error: cannot define both disk and type.")
			}
		}

		hasSource := false
		// Load up the disk for this disk if specified
		if v, ok := d.GetOk(prefix + ".disk"); ok {
			diskName := v.(string)
			diskData, err := config.clientCompute.Disks.Get(
				project, zone.Name, diskName).Do()
			if err != nil {
				return fmt.Errorf(
					"Error loading disk '%s': %s",
					diskName, err)
			}

			disk.Source = diskData.SelfLink
			hasSource = true
		} else {
			// Create a new disk
			disk.InitializeParams = &compute.AttachedDiskInitializeParams{}
		}

		if v, ok := d.GetOk(prefix + ".scratch"); ok {
			if v.(bool) {
				if hasScratchDisk {
					return fmt.Errorf("Cannot set scratch disks using both `scratch_disk` and `disk` properties")
				}
				disk.Type = "SCRATCH"
			}
		}

		// Load up the image for this disk if specified
		if v, ok := d.GetOk(prefix + ".image"); ok && !hasSource {
			imageName := v.(string)

			imageUrl, err := resolveImage(config, imageName)
			if err != nil {
				return fmt.Errorf(
					"Error resolving image name '%s': %s",
					imageName, err)
			}

			disk.InitializeParams.SourceImage = imageUrl
		} else if ok && hasSource {
			return fmt.Errorf("Cannot specify disk image when referencing an existing disk")
		}

		if v, ok := d.GetOk(prefix + ".type"); ok && !hasSource {
			diskTypeName := v.(string)
			diskType, err := readDiskType(config, zone, diskTypeName)
			if err != nil {
				return fmt.Errorf(
					"Error loading disk type '%s': %s",
					diskTypeName, err)
			}

			disk.InitializeParams.DiskType = diskType.SelfLink
		} else if ok && hasSource {
			return fmt.Errorf("Cannot specify disk type when referencing an existing disk")
		}

		if v, ok := d.GetOk(prefix + ".size"); ok && !hasSource {
			diskSizeGb := v.(int)
			disk.InitializeParams.DiskSizeGb = int64(diskSizeGb)
		} else if ok && hasSource {
			return fmt.Errorf("Cannot specify disk size when referencing an existing disk")
		}

		if v, ok := d.GetOk(prefix + ".device_name"); ok {
			disk.DeviceName = v.(string)
		}

		if v, ok := d.GetOk(prefix + ".disk_encryption_key_raw"); ok {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{}
			disk.DiskEncryptionKey.RawKey = v.(string)
		}

		disks = append(disks, &disk)
	}

	for i := 0; i < attachedDisksCount; i++ {
		prefix := fmt.Sprintf("attached_disk.%d", i)
		disk := compute.AttachedDisk{
			Source:     d.Get(prefix + ".source").(string),
			AutoDelete: false, // Don't allow autodelete; let terraform handle disk deletion
		}

		disk.Boot = i == 0 && disksCount == 0 && !hasBootDisk

		if v, ok := d.GetOk(prefix + ".device_name"); ok {
			disk.DeviceName = v.(string)
		}

		if v, ok := d.GetOk(prefix + ".disk_encryption_key_raw"); ok {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey: v.(string),
			}
		}

		disks = append(disks, &disk)
	}

	networksCount := d.Get("network.#").(int)
	networkInterfacesCount := d.Get("network_interface.#").(int)

	if networksCount > 0 && networkInterfacesCount > 0 {
		return fmt.Errorf("Error: cannot define both networks and network_interfaces.")
	}
	if networksCount == 0 && networkInterfacesCount == 0 {
		return fmt.Errorf("Error: Must define at least one network_interface.")
	}

	var networkInterfaces []*compute.NetworkInterface

	if networksCount > 0 {
		// TODO: Delete this block when removing network { }
		// Build up the list of networkInterfaces
		networkInterfaces = make([]*compute.NetworkInterface, 0, networksCount)
		for i := 0; i < networksCount; i++ {
			prefix := fmt.Sprintf("network.%d", i)
			// Load up the name of this network
			networkName := d.Get(prefix + ".source").(string)
			network, err := config.clientCompute.Networks.Get(
				project, networkName).Do()
			if err != nil {
				return fmt.Errorf(
					"Error loading network '%s': %s",
					networkName, err)
			}

			// Build the networkInterface
			var iface compute.NetworkInterface
			iface.AccessConfigs = []*compute.AccessConfig{
				&compute.AccessConfig{
					Type:  "ONE_TO_ONE_NAT",
					NatIP: d.Get(prefix + ".address").(string),
				},
			}
			iface.Network = network.SelfLink

			networkInterfaces = append(networkInterfaces, &iface)
		}
	}

	if networkInterfacesCount > 0 {
		// Build up the list of networkInterfaces
		networkInterfaces = make([]*compute.NetworkInterface, 0, networkInterfacesCount)
		for i := 0; i < networkInterfacesCount; i++ {
			prefix := fmt.Sprintf("network_interface.%d", i)
			// Load up the name of this network_interface
			networkName := d.Get(prefix + ".network").(string)
			subnetworkName := d.Get(prefix + ".subnetwork").(string)
			address := d.Get(prefix + ".address").(string)
			var networkLink, subnetworkLink string

			if networkName != "" && subnetworkName != "" {
				return fmt.Errorf("Cannot specify both network and subnetwork values.")
			} else if networkName != "" {
				networkLink, err = getNetworkLink(d, config, prefix+".network")
				if err != nil {
					return fmt.Errorf(
						"Error referencing network '%s': %s",
						networkName, err)
				}
			} else {
				subnetworkLink, err = getSubnetworkLink(d, config, prefix+".subnetwork", prefix+".subnetwork_project", "zone")
				if err != nil {
					return err
				}
			}

			// Build the networkInterface
			var iface compute.NetworkInterface
			iface.Network = networkLink
			iface.Subnetwork = subnetworkLink
			iface.NetworkIP = address

			// Handle access_config structs
			accessConfigsCount := d.Get(prefix + ".access_config.#").(int)
			iface.AccessConfigs = make([]*compute.AccessConfig, accessConfigsCount)
			for j := 0; j < accessConfigsCount; j++ {
				acPrefix := fmt.Sprintf("%s.access_config.%d", prefix, j)
				iface.AccessConfigs[j] = &compute.AccessConfig{
					Type:  "ONE_TO_ONE_NAT",
					NatIP: d.Get(acPrefix + ".nat_ip").(string),
				}
			}

			networkInterfaces = append(networkInterfaces, &iface)
		}
	}

	serviceAccountsCount := d.Get("service_account.#").(int)
	serviceAccounts := make([]*compute.ServiceAccount, 0, serviceAccountsCount)
	for i := 0; i < serviceAccountsCount; i++ {
		prefix := fmt.Sprintf("service_account.%d", i)

		scopesSet := d.Get(prefix + ".scopes").(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, v := range scopesSet.List() {
			scopes[i] = canonicalizeServiceScope(v.(string))
		}

		email := "default"
		if v := d.Get(prefix + ".email"); v != nil {
			email = v.(string)
		}

		serviceAccount := &compute.ServiceAccount{
			Email:  email,
			Scopes: scopes,
		}

		serviceAccounts = append(serviceAccounts, serviceAccount)
	}

	prefix := "scheduling.0"
	scheduling := &compute.Scheduling{}

	if val, ok := d.GetOk(prefix + ".automatic_restart"); ok {
		scheduling.AutomaticRestart = googleapi.Bool(val.(bool))
	}

	if val, ok := d.GetOk(prefix + ".preemptible"); ok {
		scheduling.Preemptible = val.(bool)
	}

	if val, ok := d.GetOk(prefix + ".on_host_maintenance"); ok {
		scheduling.OnHostMaintenance = val.(string)
	}
	scheduling.ForceSendFields = []string{"AutomaticRestart", "Preemptible"}

	// Read create timeout
	var createTimeout int
	if v, ok := d.GetOk("create_timeout"); ok {
		createTimeout = v.(int)
	}

	metadata, err := resourceInstanceMetadata(d)
	if err != nil {
		return fmt.Errorf("Error creating metadata: %s", err)
	}

	// Create the instance information
	instance := compute.Instance{
		CanIpForward:      d.Get("can_ip_forward").(bool),
		Description:       d.Get("description").(string),
		Disks:             disks,
		MachineType:       machineType.SelfLink,
		Metadata:          metadata,
		Name:              d.Get("name").(string),
		NetworkInterfaces: networkInterfaces,
		Tags:              resourceInstanceTags(d),
		Labels:            resourceInstanceLabels(d),
		ServiceAccounts:   serviceAccounts,
		Scheduling:        scheduling,
	}

	log.Printf("[INFO] Requesting instance creation")
	op, err := config.clientCompute.Instances.Insert(
		project, zone.Name, &instance).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	// Store the ID now
	d.SetId(instance.Name)

	// Wait for the operation to complete
	waitErr := computeOperationWaitTime(config, op, project, "instance to create", createTimeout)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	instance, err := getInstance(config, d)
	if err != nil || instance == nil {
		return err
	}

	md := flattenMetadata(instance.Metadata)

	if _, scriptExists := d.GetOk("metadata_startup_script"); scriptExists {
		d.Set("metadata_startup_script", md["startup-script"])
		// Note that here we delete startup-script from our metadata list. This is to prevent storing the startup-script
		// as a value in the metadata since the config specifically tracks it under 'metadata_startup_script'
		delete(md, "startup-script")
	}

	existingMetadata := d.Get("metadata").(map[string]interface{})

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

	machineTypeResource := strings.Split(instance.MachineType, "/")
	machineType := machineTypeResource[len(machineTypeResource)-1]
	d.Set("machine_type", machineType)

	// Set the service accounts
	serviceAccounts := make([]map[string]interface{}, 0, 1)
	for _, serviceAccount := range instance.ServiceAccounts {
		scopes := make([]interface{}, len(serviceAccount.Scopes))
		for i, scope := range serviceAccount.Scopes {
			scopes[i] = scope
		}
		serviceAccounts = append(serviceAccounts, map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": schema.NewSet(stringScopeHashcode, scopes),
		})
	}
	d.Set("service_account", serviceAccounts)

	networksCount := d.Get("network.#").(int)
	networkInterfacesCount := d.Get("network_interface.#").(int)

	if networksCount > 0 && networkInterfacesCount > 0 {
		return fmt.Errorf("Error: cannot define both networks and network_interfaces.")
	}
	if networksCount == 0 && networkInterfacesCount == 0 {
		return fmt.Errorf("Error: Must define at least one network_interface.")
	}

	// Set the networks
	// Use the first external IP found for the default connection info.
	externalIP := ""
	internalIP := ""
	networks := make([]map[string]interface{}, 0, 1)
	if networksCount > 0 {
		// TODO: Remove this when realizing deprecation of .network
		for i, iface := range instance.NetworkInterfaces {
			var natIP string
			for _, config := range iface.AccessConfigs {
				if config.Type == "ONE_TO_ONE_NAT" {
					natIP = config.NatIP
					break
				}
			}

			if externalIP == "" && natIP != "" {
				externalIP = natIP
			}

			network := make(map[string]interface{})
			network["name"] = iface.Name
			network["external_address"] = natIP
			network["internal_address"] = iface.NetworkIP
			network["source"] = d.Get(fmt.Sprintf("network.%d.source", i))
			networks = append(networks, network)
		}
	}
	d.Set("network", networks)

	networkInterfaces := make([]map[string]interface{}, 0, 1)
	if networkInterfacesCount > 0 {
		for i, iface := range instance.NetworkInterfaces {
			// The first non-empty ip is left in natIP
			var natIP string
			accessConfigs := make(
				[]map[string]interface{}, 0, len(iface.AccessConfigs))
			for j, config := range iface.AccessConfigs {
				accessConfigs = append(accessConfigs, map[string]interface{}{
					"nat_ip":          d.Get(fmt.Sprintf("network_interface.%d.access_config.%d.nat_ip", i, j)),
					"assigned_nat_ip": config.NatIP,
				})

				if natIP == "" {
					natIP = config.NatIP
				}
			}

			if externalIP == "" {
				externalIP = natIP
			}

			if internalIP == "" {
				internalIP = iface.NetworkIP
			}

			networkInterfaces = append(networkInterfaces, map[string]interface{}{
				"name":               iface.Name,
				"address":            iface.NetworkIP,
				"network":            iface.Network,
				"subnetwork":         iface.Subnetwork,
				"subnetwork_project": getProjectFromSubnetworkLink(iface.Subnetwork),
				"access_config":      accessConfigs,
			})
		}
	}
	d.Set("network_interface", networkInterfaces)

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
	}

	if len(instance.Labels) > 0 {
		d.Set("labels", instance.Labels)
	}

	if instance.LabelFingerprint != "" {
		d.Set("label_fingerprint", instance.LabelFingerprint)
	}

	disksCount := d.Get("disk.#").(int)
	attachedDisksCount := d.Get("attached_disk.#").(int)
	scratchDisksCount := d.Get("scratch_disk.#").(int)

	if _, ok := d.GetOk("boot_disk"); ok {
		disksCount++
	}
	if expectedDisks := disksCount + attachedDisksCount + scratchDisksCount; len(instance.Disks) != expectedDisks {
		return fmt.Errorf("Expected %d disks, API returned %d", expectedDisks, len(instance.Disks))
	}

	attachedDiskSources := make(map[string]struct{}, attachedDisksCount)
	for i := 0; i < attachedDisksCount; i++ {
		attachedDiskSources[d.Get(fmt.Sprintf("attached_disk.%d.source", i)).(string)] = struct{}{}
	}

	dIndex := 0
	adIndex := 0
	sIndex := 0
	disks := make([]map[string]interface{}, 0, disksCount)
	attachedDisks := make([]map[string]interface{}, 0, attachedDisksCount)
	scratchDisks := make([]map[string]interface{}, 0, scratchDisksCount)
	for _, disk := range instance.Disks {
		if _, ok := d.GetOk("boot_disk"); ok && disk.Boot {
			// This disk is a boot disk and there is a boot disk set in the config, therefore
			// this is the boot disk set in the config.
			d.Set("boot_disk", flattenBootDisk(d, disk))
		} else if _, ok := d.GetOk("scratch_disk"); ok && disk.Type == "SCRATCH" {
			// This disk is a scratch disk and there are scratch disks set in the config, therefore
			// this is a scratch disk set in the config.
			scratchDisks = append(scratchDisks, flattenScratchDisk(disk))
			sIndex++
		} else if _, ok := attachedDiskSources[disk.Source]; !ok {
			di := map[string]interface{}{
				"disk":                    d.Get(fmt.Sprintf("disk.%d.disk", dIndex)),
				"image":                   d.Get(fmt.Sprintf("disk.%d.image", dIndex)),
				"type":                    d.Get(fmt.Sprintf("disk.%d.type", dIndex)),
				"scratch":                 d.Get(fmt.Sprintf("disk.%d.scratch", dIndex)),
				"auto_delete":             d.Get(fmt.Sprintf("disk.%d.auto_delete", dIndex)),
				"size":                    d.Get(fmt.Sprintf("disk.%d.size", dIndex)),
				"device_name":             d.Get(fmt.Sprintf("disk.%d.device_name", dIndex)),
				"disk_encryption_key_raw": d.Get(fmt.Sprintf("disk.%d.disk_encryption_key_raw", dIndex)),
			}
			if disk.DiskEncryptionKey != nil && disk.DiskEncryptionKey.Sha256 != "" {
				di["disk_encryption_key_sha256"] = disk.DiskEncryptionKey.Sha256
			}
			disks = append(disks, di)
			dIndex++
		} else {
			di := map[string]interface{}{
				"source":                  disk.Source,
				"device_name":             disk.DeviceName,
				"disk_encryption_key_raw": d.Get(fmt.Sprintf("attached_disk.%d.disk_encryption_key_raw", adIndex)),
			}
			if disk.DiskEncryptionKey != nil && disk.DiskEncryptionKey.Sha256 != "" {
				di["disk_encryption_key_sha256"] = disk.DiskEncryptionKey.Sha256
			}
			attachedDisks = append(attachedDisks, di)
			adIndex++
		}
	}
	d.Set("disk", disks)
	d.Set("attached_disk", attachedDisks)
	d.Set("scratch_disk", scratchDisks)

	scheduling, _ := flattenScheduling(instance.Scheduling)
	d.Set("scheduling", scheduling)

	d.Set("self_link", instance.SelfLink)
	d.SetId(instance.Name)

	return nil
}

func resourceComputeInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)

	instance, err := getInstance(config, d)
	if err != nil {
		return err
	}

	// Enable partial mode for the resource since it is possible
	d.Partial(true)

	// If the Metadata has changed, then update that.
	if d.HasChange("metadata") {
		o, n := d.GetChange("metadata")
		if script, scriptExists := d.GetOk("metadata_startup_script"); scriptExists {
			if _, ok := n.(map[string]interface{})["startup-script"]; ok {
				return fmt.Errorf("Only one of metadata.startup-script and metadata_startup_script may be defined")
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

			MetadataUpdate(o.(map[string]interface{}), n.(map[string]interface{}), md)

			if err != nil {
				return fmt.Errorf("Error updating metadata: %s", err)
			}
			op, err := config.clientCompute.Instances.SetMetadata(
				project, zone, d.Id(), md).Do()
			if err != nil {
				return fmt.Errorf("Error updating metadata: %s", err)
			}

			opErr := computeOperationWait(config, op, project, "metadata to update")
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
		op, err := config.clientCompute.Instances.SetTags(
			project, zone, d.Id(), tags).Do()
		if err != nil {
			return fmt.Errorf("Error updating tags: %s", err)
		}

		opErr := computeOperationWait(config, op, project, "tags to update")
		if opErr != nil {
			return opErr
		}

		d.SetPartial("tags")
	}

	if d.HasChange("labels") {
		labels := resourceInstanceLabels(d)
		labelFingerprint := d.Get("label_fingerprint").(string)
		req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: labelFingerprint}

		op, err := config.clientCompute.Instances.SetLabels(project, zone, d.Id(), &req).Do()
		if err != nil {
			return fmt.Errorf("Error updating labels: %s", err)
		}

		opErr := computeOperationWait(config, op, project, "labels to update")
		if opErr != nil {
			return opErr
		}

		d.SetPartial("labels")
	}

	if d.HasChange("scheduling") {
		prefix := "scheduling.0"
		scheduling := &compute.Scheduling{}

		if val, ok := d.GetOk(prefix + ".automatic_restart"); ok {
			scheduling.AutomaticRestart = googleapi.Bool(val.(bool))
		}
		if val, ok := d.GetOk(prefix + ".preemptible"); ok {
			scheduling.Preemptible = val.(bool)
		}
		if val, ok := d.GetOk(prefix + ".on_host_maintenance"); ok {
			scheduling.OnHostMaintenance = val.(string)
		}
		scheduling.ForceSendFields = []string{"AutomaticRestart", "Preemptible"}

		op, err := config.clientCompute.Instances.SetScheduling(project,
			zone, d.Id(), scheduling).Do()

		if err != nil {
			return fmt.Errorf("Error updating scheduling policy: %s", err)
		}

		opErr := computeOperationWait(config, op, project, "scheduling policy update")
		if opErr != nil {
			return opErr
		}

		d.SetPartial("scheduling")
	}

	networkInterfacesCount := d.Get("network_interface.#").(int)
	if networkInterfacesCount > 0 {
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
					opErr := computeOperationWait(config, op, project, "old access_config to delete")
					if opErr != nil {
						return opErr
					}
				}

				// Create new ones
				accessConfigsCount := d.Get(prefix + ".access_config.#").(int)
				for j := 0; j < accessConfigsCount; j++ {
					acPrefix := fmt.Sprintf("%s.access_config.%d", prefix, j)
					ac := &compute.AccessConfig{
						Type:  "ONE_TO_ONE_NAT",
						NatIP: d.Get(acPrefix + ".nat_ip").(string),
					}
					op, err := config.clientCompute.Instances.AddAccessConfig(
						project, zone, d.Id(), networkName, ac).Do()
					if err != nil {
						return fmt.Errorf("Error adding new access_config: %s", err)
					}
					opErr := computeOperationWait(config, op, project, "new access_config to add")
					if opErr != nil {
						return opErr
					}
				}
			}
		}
	}

	// We made it, disable partial mode
	d.Partial(false)

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)
	log.Printf("[INFO] Requesting instance deletion: %s", d.Id())
	op, err := config.clientCompute.Instances.Delete(project, zone, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	// Wait for the operation to complete
	opErr := computeOperationWait(config, op, project, "instance to delete")
	if opErr != nil {
		return opErr
	}

	d.SetId("")
	return nil
}

func resourceInstanceMetadata(d *schema.ResourceData) (*compute.Metadata, error) {
	m := &compute.Metadata{}
	mdMap := d.Get("metadata").(map[string]interface{})
	if v, ok := d.GetOk("metadata_startup_script"); ok && v.(string) != "" {
		mdMap["startup-script"] = v
	}
	if len(mdMap) > 0 {
		m.Items = make([]*compute.MetadataItems, 0, len(mdMap))
		for key, val := range mdMap {
			v := val.(string)
			m.Items = append(m.Items, &compute.MetadataItems{
				Key:   key,
				Value: &v,
			})
		}

		// Set the fingerprint. If the metadata has never been set before
		// then this will just be blank.
		m.Fingerprint = d.Get("metadata_fingerprint").(string)
	}

	return m, nil
}

func resourceInstanceLabels(d *schema.ResourceData) map[string]string {
	labels := map[string]string{}
	if v, ok := d.GetOk("labels"); ok {
		labelMap := v.(map[string]interface{})
		for k, v := range labelMap {
			labels[k] = v.(string)
		}
	}
	return labels
}

func resourceInstanceTags(d *schema.ResourceData) *compute.Tags {
	// Calculate the tags
	var tags *compute.Tags
	if v := d.Get("tags"); v != nil {
		vs := v.(*schema.Set)
		tags = new(compute.Tags)
		tags.Items = make([]string, vs.Len())
		for i, v := range vs.List() {
			tags.Items[i] = v.(string)
		}

		tags.Fingerprint = d.Get("tags_fingerprint").(string)
	}

	return tags
}

func expandBootDisk(d *schema.ResourceData, config *Config, zone *compute.Zone, project string) (*compute.AttachedDisk, error) {
	disk := &compute.AttachedDisk{
		AutoDelete: d.Get("boot_disk.0.auto_delete").(bool),
		Boot:       true,
	}

	if v, ok := d.GetOk("boot_disk.0.device_name"); ok {
		disk.DeviceName = v.(string)
	}

	if v, ok := d.GetOk("boot_disk.0.disk_encryption_key_raw"); ok {
		disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
			RawKey: v.(string),
		}
	}

	if v, ok := d.GetOk("boot_disk.0.source"); ok {
		diskName := v.(string)
		diskData, err := config.clientCompute.Disks.Get(
			project, zone.Name, diskName).Do()
		if err != nil {
			return nil, fmt.Errorf("Error loading disk '%s': %s", diskName, err)
		}
		disk.Source = diskData.SelfLink
	}

	if _, ok := d.GetOk("boot_disk.0.initialize_params"); ok {
		disk.InitializeParams = &compute.AttachedDiskInitializeParams{}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.size"); ok {
			disk.InitializeParams.DiskSizeGb = int64(v.(int))
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.type"); ok {
			diskTypeName := v.(string)
			diskType, err := readDiskType(config, zone, diskTypeName)
			if err != nil {
				return nil, fmt.Errorf("Error loading disk type '%s': %s", diskTypeName, err)
			}
			disk.InitializeParams.DiskType = diskType.SelfLink
		}

		if v, ok := d.GetOk("boot_disk.0.initialize_params.0.image"); ok {
			imageName := v.(string)
			imageUrl, err := resolveImage(config, imageName)
			if err != nil {
				return nil, fmt.Errorf("Error resolving image name '%s': %s", imageName, err)
			}

			disk.InitializeParams.SourceImage = imageUrl
		}
	}

	return disk, nil
}

func flattenBootDisk(d *schema.ResourceData, disk *compute.AttachedDisk) []map[string]interface{} {
	sourceUrl := strings.Split(disk.Source, "/")
	result := map[string]interface{}{
		"auto_delete": disk.AutoDelete,
		"device_name": disk.DeviceName,
		"source":      sourceUrl[len(sourceUrl)-1],
		// disk_encryption_key_raw is not returned from the API, so copy it from what the user
		// originally specified to avoid diffs.
		"disk_encryption_key_raw": d.Get("boot_disk.0.disk_encryption_key_raw"),
	}
	if disk.DiskEncryptionKey != nil {
		result["disk_encryption_key_sha256"] = disk.DiskEncryptionKey.Sha256
	}
	if _, ok := d.GetOk("boot_disk.0.initialize_params.#"); ok {
		// initialize_params is not returned from the API, so copy it from what the user
		// originally specified to avoid diffs.
		m := d.Get("boot_disk.0.initialize_params")
		result["initialize_params"] = m
	}

	return []map[string]interface{}{result}
}

func expandScratchDisks(d *schema.ResourceData, config *Config, zone *compute.Zone) ([]*compute.AttachedDisk, error) {
	diskType, err := readDiskType(config, zone, "local-ssd")
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
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskType: diskType.SelfLink,
			},
		})
	}

	return scratchDisks, nil
}

func flattenScratchDisk(disk *compute.AttachedDisk) map[string]interface{} {
	result := map[string]interface{}{
		"interface": disk.Interface,
	}
	return result
}

func getProjectFromSubnetworkLink(subnetwork string) string {
	r := regexp.MustCompile(SubnetworkLinkRegex)
	if !r.MatchString(subnetwork) {
		return ""
	}

	return r.FindStringSubmatch(subnetwork)[1]
}
