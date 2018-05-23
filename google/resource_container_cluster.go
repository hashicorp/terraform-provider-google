package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	containerBeta "google.golang.org/api/container/v1beta1"
)

var (
	instanceGroupManagerURL = regexp.MustCompile(fmt.Sprintf("^https://www.googleapis.com/compute/v1/projects/(%s)/zones/([a-z0-9-]*)/instanceGroupManagers/([^/]*)", ProjectRegex))

	networkConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr_blocks": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 10,
				Elem:     cidrBlockConfig,
			},
		},
	}
	cidrBlockConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.CIDRNetwork(0, 32),
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
)

func resourceContainerCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerClusterCreate,
		Read:   resourceContainerClusterRead,
		Update: resourceContainerClusterUpdate,
		Delete: resourceContainerClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		MigrateState:  resourceContainerClusterMigrateState,

		Importer: &schema.ResourceImporter{
			State: resourceContainerClusterStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 40 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 40 characters", k))
					}
					if !regexp.MustCompile("^[a-z0-9-]+$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q can only contain lowercase letters, numbers and hyphens", k))
					}
					if !regexp.MustCompile("^[a-z]").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must start with a letter", k))
					}
					if !regexp.MustCompile("[a-z0-9]$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must end with a number or a letter", k))
					}
					return
				},
			},

			"region": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"zone"},
			},

			"zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"region"},
			},

			"additional_zones": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"addons_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_load_balancing": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"horizontal_pod_autoscaling": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"kubernetes_dashboard": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"network_policy_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"cluster_ipv4_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: orEmpty(validateRFC1918Network(8, 32)),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"enable_kubernetes_alpha": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"enable_legacy_abac": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"initial_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"logging_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"logging.googleapis.com", "none"}, false),
			},

			"maintenance_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_maintenance_window": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateFunc:     validateRFC3339Time,
										DiffSuppressFunc: rfc3339TimeDiffSuppress,
									},
									"duration": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"master_auth": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"client_certificate_config": {
							Type:             schema.TypeList,
							MaxItems:         1,
							Optional:         true,
							DiffSuppressFunc: masterAuthClientCertCfgSuppress,
							ForceNew:         true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issue_client_certificate": {
										Type:             schema.TypeBool,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: masterAuthClientCertCfgSuppress,
									},
								},
							},
						},

						"client_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"client_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"cluster_ca_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"master_authorized_networks_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     networkConfig,
			},

			"min_master_version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"monitoring_service": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "default",
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"network_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"provider": {
							Type:             schema.TypeString,
							Default:          "PROVIDER_UNSPECIFIED",
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROVIDER_UNSPECIFIED", "CALICO"}, false),
							DiffSuppressFunc: emptyOrDefaultStringSuppress("PROVIDER_UNSPECIFIED"),
						},
					},
				},
			},

			"node_config": schemaNodeConfig,

			"node_pool": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true, // TODO(danawillow): Add ability to add/remove nodePools
				Elem: &schema.Resource{
					Schema: schemaNodePool,
				},
			},

			"node_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pod_security_policy_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_group_urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"master_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_allocation_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_secondary_range_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"services_secondary_range_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"remove_default_node_pool": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"private_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"master_ipv4_cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.CIDRNetwork(28, 28),
			},
		},
	}
}

func resourceContainerClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}

	clusterName := d.Get("name").(string)

	cluster := &containerBeta.Cluster{
		Name:             clusterName,
		InitialNodeCount: int64(d.Get("initial_node_count").(int)),
	}

	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())

	if v, ok := d.GetOk("maintenance_policy"); ok {
		cluster.MaintenancePolicy = expandMaintenancePolicy(v)
	}

	if v, ok := d.GetOk("master_auth"); ok {
		masterAuths := v.([]interface{})
		masterAuth := masterAuths[0].(map[string]interface{})
		cluster.MasterAuth = &containerBeta.MasterAuth{
			Password: masterAuth["password"].(string),
			Username: masterAuth["username"].(string),
		}
		if certConfigV, ok := masterAuth["client_certificate_config"]; ok {
			certConfigs := certConfigV.([]interface{})
			if len(certConfigs) > 0 {
				certConfig := certConfigs[0].(map[string]interface{})
				cluster.MasterAuth.ClientCertificateConfig = &containerBeta.ClientCertificateConfig{
					IssueClientCertificate: certConfig["issue_client_certificate"].(bool),
				}
			}
		}
	}

	if v, ok := d.GetOk("master_authorized_networks_config"); ok {
		cluster.MasterAuthorizedNetworksConfig = expandMasterAuthorizedNetworksConfig(v)
	}

	if v, ok := d.GetOk("min_master_version"); ok {
		cluster.InitialClusterVersion = v.(string)
	}

	// Only allow setting node_version on create if it's set to the equivalent master version,
	// since `InitialClusterVersion` only accepts valid master-style versions.
	if v, ok := d.GetOk("node_version"); ok {
		// ignore -gke.X suffix for now. if it becomes a problem later, we can fix it.
		mv := strings.Split(cluster.InitialClusterVersion, "-")[0]
		nv := strings.Split(v.(string), "-")[0]
		if mv != nv {
			return fmt.Errorf("node_version and min_master_version must be set to equivalent values on create")
		}
	}

	if v, ok := d.GetOk("additional_zones"); ok {
		locationsSet := v.(*schema.Set)
		if locationsSet.Contains(location) {
			return fmt.Errorf("additional_zones should not contain the original 'zone'")
		}
		if isZone(location) {
			// GKE requires a full list of locations (including the original zone),
			// but our schema only asks for additional zones, so append the original.
			locationsSet.Add(location)
		}
		cluster.Locations = convertStringSet(locationsSet)
	}

	if v, ok := d.GetOk("cluster_ipv4_cidr"); ok {
		cluster.ClusterIpv4Cidr = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		cluster.Description = v.(string)
	}

	cluster.LegacyAbac = &containerBeta.LegacyAbac{
		Enabled:         d.Get("enable_legacy_abac").(bool),
		ForceSendFields: []string{"Enabled"},
	}

	if v, ok := d.GetOk("logging_service"); ok {
		cluster.LoggingService = v.(string)
	}

	if v, ok := d.GetOk("monitoring_service"); ok {
		cluster.MonitoringService = v.(string)
	}

	if v, ok := d.GetOk("network"); ok {
		network, err := ParseNetworkFieldValue(v.(string), d, config)
		if err != nil {
			return err
		}
		cluster.Network = network.RelativeLink()
	}

	if v, ok := d.GetOk("network_policy"); ok && len(v.([]interface{})) > 0 {
		cluster.NetworkPolicy = expandNetworkPolicy(v)
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		subnetwork, err := ParseSubnetworkFieldValue(v.(string), d, config)
		if err != nil {
			return err
		}
		cluster.Subnetwork = subnetwork.RelativeLink()
	}

	if v, ok := d.GetOk("addons_config"); ok {
		cluster.AddonsConfig = expandClusterAddonsConfig(v)
	}

	if v, ok := d.GetOk("enable_kubernetes_alpha"); ok {
		cluster.EnableKubernetesAlpha = v.(bool)
	}

	nodePoolsCount := d.Get("node_pool.#").(int)
	if nodePoolsCount > 0 {
		nodePools := make([]*containerBeta.NodePool, 0, nodePoolsCount)
		for i := 0; i < nodePoolsCount; i++ {
			prefix := fmt.Sprintf("node_pool.%d.", i)
			nodePool, err := expandNodePool(d, prefix)
			if err != nil {
				return err
			}
			nodePools = append(nodePools, nodePool)
		}
		cluster.NodePools = nodePools
	} else {
		// Node Configs have default values that are set in the expand function,
		// but can only be set if node pools are unspecified.
		cluster.NodeConfig = expandNodeConfig([]interface{}{})
	}

	if v, ok := d.GetOk("node_config"); ok {
		cluster.NodeConfig = expandNodeConfig(v)
	}

	if v, ok := d.GetOk("ip_allocation_policy"); ok {
		cluster.IpAllocationPolicy, err = expandIPAllocationPolicy(v)
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("pod_security_policy_config"); ok {
		cluster.PodSecurityPolicyConfig = expandPodSecurityPolicyConfig(v)
	}

	if v, ok := d.GetOk("master_ipv4_cidr_block"); ok {
		cluster.MasterIpv4CidrBlock = v.(string)
	}

	if v, ok := d.GetOk("private_cluster"); ok {
		if cluster.PrivateCluster = v.(bool); cluster.PrivateCluster {
			if cluster.MasterIpv4CidrBlock == "" {
				return fmt.Errorf("master_ipv4_cidr_block is mandatory when private_cluster=true")
			}
			if cluster.IpAllocationPolicy == nil {
				return fmt.Errorf("ip_allocation_policy is mandatory when private_cluster=true")
			}
		}
	}

	req := &containerBeta.CreateClusterRequest{
		Cluster: cluster,
	}

	mutexKV.Lock(containerClusterMutexKey(project, location, clusterName))
	defer mutexKV.Unlock(containerClusterMutexKey(project, location, clusterName))

	parent := fmt.Sprintf("projects/%s/locations/%s", project, location)
	op, err := config.clientContainerBeta.Projects.Locations.Clusters.Create(parent, req).Do()
	if err != nil {
		return err
	}

	d.SetId(clusterName)

	// Wait until it's created
	waitErr := containerSharedOperationWait(config, op, project, location, "creating GKE cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been created", clusterName)

	if d.Get("remove_default_node_pool").(bool) {
		parent := fmt.Sprintf("%s/nodePools/%s", containerClusterFullName(project, location, clusterName), "default-pool")
		op, err = config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(parent).Do()
		if err != nil {
			return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
		}
		err = containerSharedOperationWait(config, op, project, location, "removing default node pool", timeoutInMinutes, 3)
		if err != nil {
			return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
		}
	}

	return resourceContainerClusterRead(d, meta)
}

func resourceContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}

	cluster := &containerBeta.Cluster{}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		name := containerClusterFullName(project, location, d.Get("name").(string))
		cluster, err = config.clientContainerBeta.Projects.Locations.Clusters.Get(name).Do()
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if cluster.Status != "RUNNING" {
			return resource.RetryableError(fmt.Errorf("Cluster %q has status %q with message %q", d.Get("name"), cluster.Status, cluster.StatusMessage))
		}
		return nil
	})
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Container Cluster %q", d.Get("name").(string)))
	}

	d.Set("name", cluster.Name)

	d.Set("network_policy", flattenNetworkPolicy(cluster.NetworkPolicy))

	d.Set("zone", cluster.Zone)

	locations := schema.NewSet(schema.HashString, convertStringArrToInterface(cluster.Locations))
	locations.Remove(cluster.Zone) // Remove the original zone since we only store additional zones
	d.Set("additional_zones", locations)

	d.Set("endpoint", cluster.Endpoint)

	if cluster.MaintenancePolicy != nil {
		d.Set("maintenance_policy", flattenMaintenancePolicy(cluster.MaintenancePolicy))
	}

	masterAuth := []map[string]interface{}{
		{
			"username":               cluster.MasterAuth.Username,
			"password":               cluster.MasterAuth.Password,
			"client_certificate":     cluster.MasterAuth.ClientCertificate,
			"client_key":             cluster.MasterAuth.ClientKey,
			"cluster_ca_certificate": cluster.MasterAuth.ClusterCaCertificate,
		},
	}
	if len(cluster.MasterAuth.ClientCertificate) == 0 {
		masterAuth[0]["client_certificate_config"] = []map[string]interface{}{
			{"issue_client_certificate": false},
		}
	}
	d.Set("master_auth", masterAuth)

	if cluster.MasterAuthorizedNetworksConfig != nil {
		d.Set("master_authorized_networks_config", flattenMasterAuthorizedNetworksConfig(cluster.MasterAuthorizedNetworksConfig))
	}

	d.Set("initial_node_count", cluster.InitialNodeCount)
	d.Set("master_version", cluster.CurrentMasterVersion)
	d.Set("node_version", cluster.CurrentNodeVersion)
	d.Set("cluster_ipv4_cidr", cluster.ClusterIpv4Cidr)
	d.Set("description", cluster.Description)
	d.Set("enable_kubernetes_alpha", cluster.EnableKubernetesAlpha)
	d.Set("enable_legacy_abac", cluster.LegacyAbac.Enabled)
	d.Set("logging_service", cluster.LoggingService)
	d.Set("monitoring_service", cluster.MonitoringService)
	d.Set("network", cluster.NetworkConfig.Network)
	d.Set("subnetwork", cluster.NetworkConfig.Subnetwork)
	if err := d.Set("node_config", flattenNodeConfig(cluster.NodeConfig)); err != nil {
		return err
	}
	d.Set("project", project)
	if cluster.AddonsConfig != nil {
		d.Set("addons_config", flattenClusterAddonsConfig(cluster.AddonsConfig))
	}
	nps, err := flattenClusterNodePools(d, config, cluster.NodePools)
	if err != nil {
		return err
	}
	d.Set("node_pool", nps)

	if cluster.IpAllocationPolicy != nil {
		if err := d.Set("ip_allocation_policy", flattenIPAllocationPolicy(cluster.IpAllocationPolicy)); err != nil {
			return err
		}
	}

	if igUrls, err := getInstanceGroupUrlsFromManagerUrls(config, cluster.InstanceGroupUrls); err != nil {
		return err
	} else {
		d.Set("instance_group_urls", igUrls)
	}

	if cluster.PodSecurityPolicyConfig != nil {
		if err := d.Set("pod_security_policy_config", flattenPodSecurityPolicyConfig(cluster.PodSecurityPolicyConfig)); err != nil {
			return err
		}
	}

	d.Set("private_cluster", cluster.PrivateCluster)
	d.Set("master_ipv4_cidr_block", cluster.MasterIpv4CidrBlock)

	return nil
}

func resourceContainerClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}

	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())

	d.Partial(true)

	lockKey := containerClusterMutexKey(project, location, clusterName)

	updateFunc := func(req *containerBeta.UpdateClusterRequest, updateDescription string) func() error {
		return func() error {
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req).Do()
			if err != nil {
				return err
			}
			// Wait until it's updated
			return containerSharedOperationWait(config, op, project, location, updateDescription, timeoutInMinutes, 2)
		}
	}

	// The ClusterUpdate object that we use for most of these updates only allows updating one field at a time,
	// so we have to make separate calls for each field that we want to update. The order here is fairly arbitrary-
	// if the order of updating fields does matter, it is called out explicitly.
	if d.HasChange("master_authorized_networks_config") {
		c := d.Get("master_authorized_networks_config")
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredMasterAuthorizedNetworksConfig: expandMasterAuthorizedNetworksConfig(c),
			},
		}

		updateF := updateFunc(req, "updating GKE cluster master authorized networks")
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}
		log.Printf("[INFO] GKE cluster %s master authorized networks config has been updated", d.Id())

		d.SetPartial("master_authorized_networks_config")
	}

	// The master must be updated before the nodes
	if d.HasChange("min_master_version") {
		desiredMasterVersion := d.Get("min_master_version").(string)
		currentMasterVersion := d.Get("master_version").(string)
		des, err := version.NewVersion(desiredMasterVersion)
		if err != nil {
			return err
		}
		cur, err := version.NewVersion(currentMasterVersion)
		if err != nil {
			return err
		}

		// Only upgrade the master if the current version is lower than the desired version
		if cur.LessThan(des) {
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredMasterVersion: desiredMasterVersion,
				},
			}

			updateF := updateFunc(req, "updating GKE master version")
			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}
			log.Printf("[INFO] GKE cluster %s: master has been updated to %s", d.Id(), desiredMasterVersion)
		}
		d.SetPartial("min_master_version")
	}

	if d.HasChange("node_version") {
		desiredNodeVersion := d.Get("node_version").(string)
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredNodeVersion: desiredNodeVersion,
			},
		}

		updateF := updateFunc(req, "updating GKE node version")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}
		log.Printf("[INFO] GKE cluster %s: nodes have been updated to %s", d.Id(),
			desiredNodeVersion)

		d.SetPartial("node_version")
	}

	if d.HasChange("addons_config") {
		if ac, ok := d.GetOk("addons_config"); ok {
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredAddonsConfig: expandClusterAddonsConfig(ac),
				},
			}

			updateF := updateFunc(req, "updating GKE cluster addons")
			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}

			log.Printf("[INFO] GKE cluster %s addons have been updated", d.Id())

			d.SetPartial("addons_config")
		}
	}

	if d.HasChange("maintenance_policy") {
		var req *containerBeta.SetMaintenancePolicyRequest
		if mp, ok := d.GetOk("maintenance_policy"); ok {
			req = &containerBeta.SetMaintenancePolicyRequest{
				MaintenancePolicy: expandMaintenancePolicy(mp),
			}
		} else {
			req = &containerBeta.SetMaintenancePolicyRequest{
				NullFields: []string{"MaintenancePolicy"},
			}
		}

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetMaintenancePolicy(name, req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerSharedOperationWait(config, op, project, location, "updating GKE cluster maintenance policy", timeoutInMinutes, 2)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s maintenance policy has been updated", d.Id())

		d.SetPartial("maintenance_policy")
	}

	if d.HasChange("additional_zones") {
		azSetOldI, azSetNewI := d.GetChange("additional_zones")
		azSetNew := azSetNewI.(*schema.Set)
		azSetOld := azSetOldI.(*schema.Set)
		if azSetNew.Contains(location) {
			return fmt.Errorf("additional_zones should not contain the original 'zone'")
		}
		// Since we can't add & remove zones in the same request, first add all the
		// zones, then remove the ones we aren't using anymore.
		azSet := azSetOld.Union(azSetNew)

		if isZone(location) {
			azSet.Add(location)
		}

		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredLocations: convertStringSet(azSet),
			},
		}

		updateF := updateFunc(req, "updating GKE cluster locations")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		if isZone(location) {
			azSetNew.Add(location)
		}
		if !azSet.Equal(azSetNew) {
			req = &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredLocations: convertStringSet(azSetNew),
				},
			}

			updateF := updateFunc(req, "updating GKE cluster locations")
			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}
		}

		log.Printf("[INFO] GKE cluster %s locations have been updated to %v", d.Id(), azSet.List())

		d.SetPartial("additional_zones")
	}

	if d.HasChange("enable_legacy_abac") {
		enabled := d.Get("enable_legacy_abac").(bool)
		req := &containerBeta.SetLegacyAbacRequest{
			Enabled:         enabled,
			ForceSendFields: []string{"Enabled"},
		}

		updateF := func() error {
			log.Println("[DEBUG] updating enable_legacy_abac")
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetLegacyAbac(name, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			err = containerSharedOperationWait(config, op, project, location, "updating GKE legacy ABAC", timeoutInMinutes, 2)
			log.Println("[DEBUG] done updating enable_legacy_abac")
			return err
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s legacy ABAC has been updated to %v", d.Id(), enabled)

		d.SetPartial("enable_legacy_abac")
	}

	if d.HasChange("monitoring_service") {
		desiredMonitoringService := d.Get("monitoring_service").(string)

		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredMonitoringService: desiredMonitoringService,
			},
		}

		updateF := updateFunc(req, "updating GKE cluster monitoring service")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}
		log.Printf("[INFO] Monitoring service for GKE cluster %s has been updated to %s", d.Id(),
			desiredMonitoringService)

		d.SetPartial("monitoring_service")
	}

	if d.HasChange("network_policy") {
		np := d.Get("network_policy")
		req := &containerBeta.SetNetworkPolicyRequest{
			NetworkPolicy: expandNetworkPolicy(np),
		}

		updateF := func() error {
			log.Println("[DEBUG] updating network_policy")
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetNetworkPolicy(name, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			err = containerSharedOperationWait(config, op, project, location, "updating GKE cluster network policy", timeoutInMinutes, 2)
			log.Println("[DEBUG] done updating network_policy")
			return err
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Network policy for GKE cluster %s has been updated", d.Id())

		d.SetPartial("network_policy")

	}

	if n, ok := d.GetOk("node_pool.#"); ok {
		for i := 0; i < n.(int); i++ {
			nodePoolInfo, err := extractNodePoolInformationFromCluster(d, config, clusterName)
			if err != nil {
				return err
			}

			if err := nodePoolUpdate(d, meta, nodePoolInfo, fmt.Sprintf("node_pool.%d.", i), timeoutInMinutes); err != nil {
				return err
			}
		}
		d.SetPartial("node_pool")
	}

	if d.HasChange("logging_service") {
		logging := d.Get("logging_service").(string)

		req := &containerBeta.SetLoggingServiceRequest{
			LoggingService: logging,
		}
		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetLogging(name, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerSharedOperationWait(config, op, project, location, "updating GKE logging service", timeoutInMinutes, 2)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s: logging service has been updated to %s", d.Id(),
			logging)
		d.SetPartial("logging_service")
	}

	if d.HasChange("pod_security_policy_config") {
		c := d.Get("pod_security_policy_config")
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredPodSecurityPolicyConfig: expandPodSecurityPolicyConfig(c),
			},
		}

		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.Zones.Clusters.Update(project, location, clusterName, req).Do()
			if err != nil {
				return err
			}
			// Wait until it's updated
			return containerSharedOperationWait(config, op, project, location, "updating GKE cluster pod security policy config", timeoutInMinutes, 2)
		}
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}
		log.Printf("[INFO] GKE cluster %s pod security policy config has been updated", d.Id())

		d.SetPartial("pod_security_policy_config")
	}

	if d.HasChange("remove_default_node_pool") && d.Get("remove_default_node_pool").(bool) {
		name := fmt.Sprintf("%s/nodePools/%s", containerClusterFullName(project, location, clusterName), "default-pool")
		op, err := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(name).Do()
		if err != nil {
			return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
		}
		err = containerSharedOperationWait(config, op, project, location, "removing default node pool", timeoutInMinutes, 3)
		if err != nil {
			return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
		}
	}

	d.Partial(false)

	return resourceContainerClusterRead(d, meta)
}

func resourceContainerClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	var location string
	locations := []string{}
	if regionName, isRegionalCluster := d.GetOk("region"); !isRegionalCluster {
		location, err = getZone(d, config)
		if err != nil {
			return err
		}
		locations = append(locations, location)
	} else {
		location = regionName.(string)
	}

	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	log.Printf("[DEBUG] Deleting GKE cluster %s", d.Get("name").(string))
	mutexKV.Lock(containerClusterMutexKey(project, location, clusterName))
	defer mutexKV.Unlock(containerClusterMutexKey(project, location, clusterName))

	var op interface{}
	var count = 0
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		count++

		name := containerClusterFullName(project, location, clusterName)
		op, err = config.clientContainerBeta.Projects.Locations.Clusters.Delete(name).Do()

		if err != nil {
			log.Printf("[WARNING] Cluster is still not ready to delete, retrying %s", clusterName)
			return resource.RetryableError(err)
		}

		if count == 15 {
			return resource.NonRetryableError(fmt.Errorf("Error retrying to delete cluster %s", clusterName))
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting Cluster: %s", err)
	}

	// Wait until it's deleted
	waitErr := containerSharedOperationWait(config, op, project, location, "deleting GKE cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been deleted", d.Id())

	d.SetId("")

	return nil
}

// container engine's API currently mistakenly returns the instance group manager's
// URL instead of the instance group's URL in its responses. This shim detects that
// error, and corrects it, by fetching the instance group manager URL and retrieving
// the instance group manager, then using that to look up the instance group URL, which
// is then substituted.
//
// This should be removed when the API response is fixed.
func getInstanceGroupUrlsFromManagerUrls(config *Config, igmUrls []string) ([]string, error) {
	instanceGroupURLs := make([]string, 0, len(igmUrls))
	for _, u := range igmUrls {
		if !instanceGroupManagerURL.MatchString(u) {
			instanceGroupURLs = append(instanceGroupURLs, u)
			continue
		}
		matches := instanceGroupManagerURL.FindStringSubmatch(u)
		instanceGroupManager, err := config.clientCompute.InstanceGroupManagers.Get(matches[1], matches[2], matches[3]).Do()
		if err != nil {
			return nil, fmt.Errorf("Error reading instance group manager returned as an instance group URL: %s", err)
		}
		instanceGroupURLs = append(instanceGroupURLs, instanceGroupManager.InstanceGroup)
	}
	return instanceGroupURLs, nil
}

func expandClusterAddonsConfig(configured interface{}) *containerBeta.AddonsConfig {
	config := configured.([]interface{})[0].(map[string]interface{})
	ac := &containerBeta.AddonsConfig{}

	if v, ok := config["http_load_balancing"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.HttpLoadBalancing = &containerBeta.HttpLoadBalancing{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["horizontal_pod_autoscaling"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.HorizontalPodAutoscaling = &containerBeta.HorizontalPodAutoscaling{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["kubernetes_dashboard"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.KubernetesDashboard = &containerBeta.KubernetesDashboard{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["network_policy_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.NetworkPolicyConfig = &containerBeta.NetworkPolicyConfig{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	return ac
}

func expandIPAllocationPolicy(configured interface{}) (*containerBeta.IPAllocationPolicy, error) {
	ap := &containerBeta.IPAllocationPolicy{}
	l := configured.([]interface{})
	if len(l) > 0 {
		if config, ok := l[0].(map[string]interface{}); ok {
			ap.UseIpAliases = true
			if v, ok := config["cluster_secondary_range_name"]; ok {
				ap.ClusterSecondaryRangeName = v.(string)
			}

			if v, ok := config["services_secondary_range_name"]; ok {
				ap.ServicesSecondaryRangeName = v.(string)
			}
		} else {
			return nil, fmt.Errorf("clusters using IP aliases must specify secondary ranges")
		}
	}

	return ap, nil
}

func expandMaintenancePolicy(configured interface{}) *containerBeta.MaintenancePolicy {
	result := &containerBeta.MaintenancePolicy{}
	if len(configured.([]interface{})) > 0 {
		maintenancePolicy := configured.([]interface{})[0].(map[string]interface{})
		dailyMaintenanceWindow := maintenancePolicy["daily_maintenance_window"].([]interface{})[0].(map[string]interface{})
		startTime := dailyMaintenanceWindow["start_time"].(string)
		result.Window = &containerBeta.MaintenanceWindow{
			DailyMaintenanceWindow: &containerBeta.DailyMaintenanceWindow{
				StartTime: startTime,
			},
		}
	}
	return result
}

func expandMasterAuthorizedNetworksConfig(configured interface{}) *containerBeta.MasterAuthorizedNetworksConfig {
	result := &containerBeta.MasterAuthorizedNetworksConfig{}
	if len(configured.([]interface{})) > 0 {
		result.Enabled = true
		config := configured.([]interface{})[0].(map[string]interface{})
		if _, ok := config["cidr_blocks"]; ok {
			cidrBlocks := config["cidr_blocks"].(*schema.Set).List()
			result.CidrBlocks = make([]*containerBeta.CidrBlock, 0)
			for _, v := range cidrBlocks {
				cidrBlock := v.(map[string]interface{})
				result.CidrBlocks = append(result.CidrBlocks, &containerBeta.CidrBlock{
					CidrBlock:   cidrBlock["cidr_block"].(string),
					DisplayName: cidrBlock["display_name"].(string),
				})
			}
		}
	}
	return result
}

func expandNetworkPolicy(configured interface{}) *containerBeta.NetworkPolicy {
	result := &containerBeta.NetworkPolicy{}
	if configured != nil && len(configured.([]interface{})) > 0 {
		config := configured.([]interface{})[0].(map[string]interface{})
		if enabled, ok := config["enabled"]; ok && enabled.(bool) {
			result.Enabled = true
			if provider, ok := config["provider"]; ok {
				result.Provider = provider.(string)
			}
		}
	}
	return result
}

func expandPodSecurityPolicyConfig(configured interface{}) *containerBeta.PodSecurityPolicyConfig {
	result := &containerBeta.PodSecurityPolicyConfig{}
	if len(configured.([]interface{})) > 0 {
		config := configured.([]interface{})[0].(map[string]interface{})
		result.Enabled = config["enabled"].(bool)
		result.ForceSendFields = []string{"Enabled"}
	}
	return result
}

func flattenNetworkPolicy(c *containerBeta.NetworkPolicy) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"enabled":  c.Enabled,
			"provider": c.Provider,
		})
	} else {
		// Explicitly set the network policy to the default.
		result = append(result, map[string]interface{}{
			"enabled":  false,
			"provider": "PROVIDER_UNSPECIFIED",
		})
	}
	return result
}

func flattenClusterAddonsConfig(c *containerBeta.AddonsConfig) []map[string]interface{} {
	result := make(map[string]interface{})
	if c.HorizontalPodAutoscaling != nil {
		result["horizontal_pod_autoscaling"] = []map[string]interface{}{
			{
				"disabled": c.HorizontalPodAutoscaling.Disabled,
			},
		}
	}
	if c.HttpLoadBalancing != nil {
		result["http_load_balancing"] = []map[string]interface{}{
			{
				"disabled": c.HttpLoadBalancing.Disabled,
			},
		}
	}
	if c.KubernetesDashboard != nil {
		result["kubernetes_dashboard"] = []map[string]interface{}{
			{
				"disabled": c.KubernetesDashboard.Disabled,
			},
		}
	}
	if c.NetworkPolicyConfig != nil {
		result["network_policy_config"] = []map[string]interface{}{
			{
				"disabled": c.NetworkPolicyConfig.Disabled,
			},
		}
	}

	return []map[string]interface{}{result}
}

func flattenClusterNodePools(d *schema.ResourceData, config *Config, c []*containerBeta.NodePool) ([]map[string]interface{}, error) {
	nodePools := make([]map[string]interface{}, 0, len(c))

	for i, np := range c {
		nodePool, err := flattenNodePool(d, config, np, fmt.Sprintf("node_pool.%d.", i))
		if err != nil {
			return nil, err
		}
		nodePools = append(nodePools, nodePool)
	}

	return nodePools, nil
}

func flattenIPAllocationPolicy(c *containerBeta.IPAllocationPolicy) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"cluster_secondary_range_name":  c.ClusterSecondaryRangeName,
			"services_secondary_range_name": c.ServicesSecondaryRangeName,
		},
	}
}

func flattenMaintenancePolicy(mp *containerBeta.MaintenancePolicy) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"daily_maintenance_window": []map[string]interface{}{
				{
					"start_time": mp.Window.DailyMaintenanceWindow.StartTime,
					"duration":   mp.Window.DailyMaintenanceWindow.Duration,
				},
			},
		},
	}
}

func flattenMasterAuthorizedNetworksConfig(c *containerBeta.MasterAuthorizedNetworksConfig) []map[string]interface{} {
	if len(c.CidrBlocks) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	if c.Enabled {
		cidrBlocks := make([]interface{}, 0, len(c.CidrBlocks))
		for _, v := range c.CidrBlocks {
			cidrBlocks = append(cidrBlocks, map[string]interface{}{
				"cidr_block":   v.CidrBlock,
				"display_name": v.DisplayName,
			})
		}
		result["cidr_blocks"] = schema.NewSet(schema.HashResource(cidrBlockConfig), cidrBlocks)
	}
	return []map[string]interface{}{result}
}

func flattenPodSecurityPolicyConfig(c *containerBeta.PodSecurityPolicyConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"enabled": c.Enabled,
		},
	}
}

func resourceContainerClusterStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	switch len(parts) {
	case 2:
		if loc := parts[0]; isZone(loc) {
			d.Set("zone", loc)
		} else {
			d.Set("region", loc)
		}
		d.Set("name", parts[1])
	case 3:
		d.Set("project", parts[0])
		if loc := parts[1]; isZone(loc) {
			d.Set("zone", loc)
		} else {
			d.Set("region", loc)
		}
		d.Set("name", parts[2])
	default:
		return nil, fmt.Errorf("Invalid container cluster specifier. Expecting {zone}/{name} or {project}/{zone}/{name}")
	}

	d.SetId(parts[len(parts)-1])
	return []*schema.ResourceData{d}, nil
}

func containerClusterMutexKey(project, location, clusterName string) string {
	return fmt.Sprintf("google-container-cluster/%s/%s/%s", project, location, clusterName)
}

func containerClusterFullName(project, location, cluster string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, cluster)
}

func extractNodePoolInformationFromCluster(d *schema.ResourceData, config *Config, clusterName string) (*NodePoolInformation, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return nil, err
	}

	return &NodePoolInformation{
		project:  project,
		location: location,
		cluster:  d.Get("name").(string),
	}, nil
}

// We want to suppress diffs for empty or default client certificate configs, i.e:
// 	[{ "issue_client_certificate": true}] --> []
//  [] -> [{ "issue_client_certificate": true}]
func masterAuthClientCertCfgSuppress(k, old, new string, r *schema.ResourceData) bool {
	var clientConfig map[string]interface{}
	if v, ok := r.GetOk("master_auth"); ok {
		masterAuths := v.([]interface{})
		masterAuth := masterAuths[0].(map[string]interface{})
		cfgs := masterAuth["client_certificate_config"].([]interface{})
		if len(cfgs) > 0 {
			clientConfig = cfgs[0].(map[string]interface{})
		}
	}

	if strings.HasSuffix(k, "client_certificate_config.#") && old == "0" && new == "1" {
		// nil --> { "issue_client_certificate": true }
		if issueCert, ok := clientConfig["issue_client_certificate"]; ok {
			return issueCert.(bool)
		}
	}

	return strings.HasSuffix(k, ".issue_client_certificate") && old == "" && new == "true"
}
