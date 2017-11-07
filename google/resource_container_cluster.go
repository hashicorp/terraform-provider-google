package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/container/v1"
)

var (
	instanceGroupManagerURL = regexp.MustCompile("^https://www.googleapis.com/compute/v1/projects/([a-z][a-z0-9-]{5}(?:[-a-z0-9]{0,23}[a-z0-9])?)/zones/([a-z0-9-]*)/instanceGroupManagers/([^/]*)")
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

			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
					},
				},
			},

			"cluster_ipv4_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateRFC1918Network(8, 32),
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
				Default:  true,
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
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_maintenance_window": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validateRFC3339Time,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_blocks": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 10,
							Elem: &schema.Resource{
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
							},
						},
					},
				},
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
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "default",
				ForceNew:  true,
				StateFunc: StoreResourceName,
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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"subnetwork": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
		},
	}
}

func resourceContainerClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zoneName := d.Get("zone").(string)
	clusterName := d.Get("name").(string)

	cluster := &container.Cluster{
		Name:             clusterName,
		InitialNodeCount: int64(d.Get("initial_node_count").(int)),
	}

	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())

	if v, ok := d.GetOk("maintenance_policy"); ok {
		maintenancePolicy := v.([]interface{})[0].(map[string]interface{})
		dailyMaintenanceWindow := maintenancePolicy["daily_maintenance_window"].([]interface{})[0].(map[string]interface{})
		startTime := dailyMaintenanceWindow["start_time"].(string)
		cluster.MaintenancePolicy = &container.MaintenancePolicy{
			Window: &container.MaintenanceWindow{
				DailyMaintenanceWindow: &container.DailyMaintenanceWindow{
					StartTime: startTime,
				},
			},
		}
	}

	if v, ok := d.GetOk("master_auth"); ok {
		masterAuths := v.([]interface{})
		masterAuth := masterAuths[0].(map[string]interface{})
		cluster.MasterAuth = &container.MasterAuth{
			Password: masterAuth["password"].(string),
			Username: masterAuth["username"].(string),
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
		locationsList := v.(*schema.Set).List()
		locations := []string{}
		for _, v := range locationsList {
			location := v.(string)
			locations = append(locations, location)
			if location == zoneName {
				return fmt.Errorf("additional_zones should not contain the original 'zone'.")
			}
		}
		locations = append(locations, zoneName)
		cluster.Locations = locations
	}

	if v, ok := d.GetOk("cluster_ipv4_cidr"); ok {
		cluster.ClusterIpv4Cidr = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		cluster.Description = v.(string)
	}

	cluster.LegacyAbac = &container.LegacyAbac{
		Enabled:         d.Get("enable_legacy_abac").(bool),
		ForceSendFields: []string{"Enabled"},
	}

	if v, ok := d.GetOk("logging_service"); ok {
		cluster.LoggingService = v.(string)
	}

	if v, ok := d.GetOk("monitoring_service"); ok {
		cluster.MonitoringService = v.(string)
	}

	if _, ok := d.GetOk("network"); ok {
		network, err := getNetworkName(d, "network")
		if err != nil {
			return err
		}
		cluster.Network = network
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		cluster.Subnetwork = v.(string)
	}

	if v, ok := d.GetOk("addons_config"); ok {
		cluster.AddonsConfig = expandClusterAddonsConfig(v)
	}

	if v, ok := d.GetOk("node_config"); ok {
		cluster.NodeConfig = expandNodeConfig(v)
	}

	if v, ok := d.GetOk("enable_kubernetes_alpha"); ok {
		cluster.EnableKubernetesAlpha = v.(bool)
	}

	nodePoolsCount := d.Get("node_pool.#").(int)
	if nodePoolsCount > 0 {
		nodePools := make([]*container.NodePool, 0, nodePoolsCount)
		for i := 0; i < nodePoolsCount; i++ {
			prefix := fmt.Sprintf("node_pool.%d.", i)
			nodePool, err := expandNodePool(d, prefix)
			if err != nil {
				return err
			}
			nodePools = append(nodePools, nodePool)
		}
		cluster.NodePools = nodePools
	}

	req := &container.CreateClusterRequest{
		Cluster: cluster,
	}

	op, err := config.clientContainer.Projects.Zones.Clusters.Create(
		project, zoneName, req).Do()
	if err != nil {
		return err
	}

	// Wait until it's created
	waitErr := containerOperationWait(config, op, project, zoneName, "creating GKE cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been created", clusterName)

	d.SetId(clusterName)

	return resourceContainerClusterRead(d, meta)
}

func resourceContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zoneName := d.Get("zone").(string)

	var cluster *container.Cluster
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		cluster, err = config.clientContainer.Projects.Zones.Clusters.Get(
			project, zoneName, d.Get("name").(string)).Do()
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
	d.Set("zone", cluster.Zone)

	locations := []string{}
	if len(cluster.Locations) > 1 {
		for _, location := range cluster.Locations {
			if location != cluster.Zone {
				locations = append(locations, location)
			}
		}
	}
	d.Set("additional_zones", locations)

	d.Set("endpoint", cluster.Endpoint)

	if cluster.MaintenancePolicy != nil && cluster.MaintenancePolicy.Window != nil && cluster.MaintenancePolicy.Window.DailyMaintenanceWindow != nil {
		maintenancePolicy := []map[string]interface{}{
			{
				"daily_maintenance_window": []map[string]interface{}{
					{
						"start_time": cluster.MaintenancePolicy.Window.DailyMaintenanceWindow.StartTime,
						"duration":   cluster.MaintenancePolicy.Window.DailyMaintenanceWindow.Duration,
					},
				},
			},
		}
		d.Set("maintenance_policy", maintenancePolicy)
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
	d.Set("network", cluster.Network)
	d.Set("subnetwork", cluster.Subnetwork)
	d.Set("node_config", flattenNodeConfig(cluster.NodeConfig))
	if cluster.AddonsConfig != nil {
		d.Set("addons_config", flattenClusterAddonsConfig(cluster.AddonsConfig))
	}
	nps, err := flattenClusterNodePools(d, config, cluster.NodePools)
	if err != nil {
		return err
	}
	d.Set("node_pool", nps)

	if igUrls, err := getInstanceGroupUrlsFromManagerUrls(config, cluster.InstanceGroupUrls); err != nil {
		return err
	} else {
		d.Set("instance_group_urls", igUrls)
	}

	return nil
}

func resourceContainerClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zoneName := d.Get("zone").(string)
	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())

	d.Partial(true)

	if d.HasChange("master_authorized_networks_config") {
		c := d.Get("master_authorized_networks_config")
		req := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredMasterAuthorizedNetworksConfig: expandMasterAuthorizedNetworksConfig(c),
			},
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.Update(
			project, zoneName, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE cluster master authorized networks", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
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
			req := &container.UpdateClusterRequest{
				Update: &container.ClusterUpdate{
					DesiredMasterVersion: desiredMasterVersion,
				},
			}
			op, err := config.clientContainer.Projects.Zones.Clusters.Update(
				project, zoneName, clusterName, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE master version", timeoutInMinutes, 2)
			if waitErr != nil {
				return waitErr
			}

			log.Printf("[INFO] GKE cluster %s: master has been updated to %s", d.Id(),
				desiredMasterVersion)
		}

		d.SetPartial("min_master_version")
	}

	if d.HasChange("node_version") {
		desiredNodeVersion := d.Get("node_version").(string)

		req := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredNodeVersion: desiredNodeVersion,
			},
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.Update(
			project, zoneName, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE node version", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] GKE cluster %s: nodes have been updated to %s", d.Id(),
			desiredNodeVersion)

		d.SetPartial("node_version")
	}

	if d.HasChange("addons_config") {
		if ac, ok := d.GetOk("addons_config"); ok {
			req := &container.UpdateClusterRequest{
				Update: &container.ClusterUpdate{
					DesiredAddonsConfig: expandClusterAddonsConfig(ac),
				},
			}
			op, err := config.clientContainer.Projects.Zones.Clusters.Update(
				project, zoneName, clusterName, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE cluster addons", timeoutInMinutes, 2)
			if waitErr != nil {
				return waitErr
			}

			log.Printf("[INFO] GKE cluster %s addons have been updated", d.Id())

			d.SetPartial("addons_config")
		}
	}

	if d.HasChange("additional_zones") {
		azSet := d.Get("additional_zones").(*schema.Set)
		if azSet.Contains(zoneName) {
			return fmt.Errorf("additional_zones should not contain the original 'zone'.")
		}
		azs := convertStringArr(azSet.List())
		locations := append(azs, zoneName)
		req := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredLocations: locations,
			},
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.Update(
			project, zoneName, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE cluster locations", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] GKE cluster %s locations have been updated to %v", d.Id(),
			locations)

		d.SetPartial("additional_zones")
	}

	if d.HasChange("enable_legacy_abac") {
		enabled := d.Get("enable_legacy_abac").(bool)
		req := &container.SetLegacyAbacRequest{
			Enabled:         enabled,
			ForceSendFields: []string{"Enabled"},
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.LegacyAbac(project, zoneName, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE legacy ABAC", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] GKE cluster %s legacy ABAC has been updated to %v", d.Id(), enabled)

		d.SetPartial("enable_legacy_abac")
	}

	if d.HasChange("monitoring_service") {
		desiredMonitoringService := d.Get("monitoring_service").(string)

		req := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredMonitoringService: desiredMonitoringService,
			},
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.Update(
			project, zoneName, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE cluster monitoring service", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}
		log.Printf("[INFO] Monitoring service for GKE cluster %s has been updated to %s", d.Id(),
			desiredMonitoringService)

		d.SetPartial("monitoring_service")
	}

	if n, ok := d.GetOk("node_pool.#"); ok {
		for i := 0; i < n.(int); i++ {
			if err := nodePoolUpdate(d, meta, clusterName, fmt.Sprintf("node_pool.%d.", i), timeoutInMinutes); err != nil {
				return err
			}
		}
		d.SetPartial("node_pool")
	}

	if d.HasChange("logging_service") {
		logging := d.Get("logging_service").(string)

		req := &container.SetLoggingServiceRequest{
			LoggingService: logging,
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.Logging(
			project, zoneName, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zoneName, "updating GKE logging service", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] GKE cluster %s: logging service has been updated to %s", d.Id(),
			logging)
		d.SetPartial("logging_service")
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

	zoneName := d.Get("zone").(string)
	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	log.Printf("[DEBUG] Deleting GKE cluster %s", d.Get("name").(string))
	op, err := config.clientContainer.Projects.Zones.Clusters.Delete(
		project, zoneName, clusterName).Do()
	if err != nil {
		return err
	}

	// Wait until it's deleted
	waitErr := containerOperationWait(config, op, project, zoneName, "deleting GKE cluster", timeoutInMinutes, 3)
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

func expandClusterAddonsConfig(configured interface{}) *container.AddonsConfig {
	config := configured.([]interface{})[0].(map[string]interface{})
	ac := &container.AddonsConfig{}

	if v, ok := config["http_load_balancing"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.HttpLoadBalancing = &container.HttpLoadBalancing{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["horizontal_pod_autoscaling"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.HorizontalPodAutoscaling = &container.HorizontalPodAutoscaling{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	if v, ok := config["kubernetes_dashboard"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.KubernetesDashboard = &container.KubernetesDashboard{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}
	return ac
}

func expandMasterAuthorizedNetworksConfig(configured interface{}) *container.MasterAuthorizedNetworksConfig {
	result := &container.MasterAuthorizedNetworksConfig{}
	if len(configured.([]interface{})) > 0 {
		result.Enabled = true
		config := configured.([]interface{})[0].(map[string]interface{})
		if _, ok := config["cidr_blocks"]; ok {
			cidrBlocks := config["cidr_blocks"].(*schema.Set).List()
			result.CidrBlocks = make([]*container.CidrBlock, 0)
			for _, v := range cidrBlocks {
				cidrBlock := v.(map[string]interface{})
				result.CidrBlocks = append(result.CidrBlocks, &container.CidrBlock{
					CidrBlock:   cidrBlock["cidr_block"].(string),
					DisplayName: cidrBlock["display_name"].(string),
				})
			}
		}
	}
	return result
}

func flattenClusterAddonsConfig(c *container.AddonsConfig) []map[string]interface{} {
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
	return []map[string]interface{}{result}
}

func flattenClusterNodePools(d *schema.ResourceData, config *Config, c []*container.NodePool) ([]map[string]interface{}, error) {
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

func flattenMasterAuthorizedNetworksConfig(c *container.MasterAuthorizedNetworksConfig) []map[string]interface{} {
	result := make(map[string]interface{})
	if c.Enabled && len(c.CidrBlocks) > 0 {
		cidrBlocks := make([]map[string]interface{}, 0, len(c.CidrBlocks))
		for _, v := range c.CidrBlocks {
			cidrBlocks = append(cidrBlocks, map[string]interface{}{
				"cidr_block":   v.CidrBlock,
				"display_name": v.DisplayName,
			})
		}
		result["cidr_blocks"] = cidrBlocks
	}
	return []map[string]interface{}{result}
}

func resourceContainerClusterStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid container cluster specifier. Expecting {zone}/{name}")
	}

	d.Set("zone", parts[0])
	d.Set("name", parts[1])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
