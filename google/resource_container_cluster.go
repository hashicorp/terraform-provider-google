package google

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	containerBeta "google.golang.org/api/container/v1beta1"
)

var (
	instanceGroupManagerURL = regexp.MustCompile(fmt.Sprintf("projects/(%s)/zones/([a-z0-9-]*)/instanceGroupManagers/([^/]*)", ProjectRegex))

	networkConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr_blocks": {
				Type: schema.TypeSet,
				// Despite being the only entry in a nested block, this should be kept
				// Optional. Expressing the parent with no entries and omitting the
				// parent entirely are semantically different.
				Optional: true,
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

	ipAllocationCidrBlockFields = []string{"ip_allocation_policy.0.cluster_ipv4_cidr_block", "ip_allocation_policy.0.services_ipv4_cidr_block"}
	ipAllocationRangeFields     = []string{"ip_allocation_policy.0.cluster_secondary_range_name", "ip_allocation_policy.0.services_secondary_range_name"}

	addonsConfigKeys = []string{
		"addons_config.0.http_load_balancing",
		"addons_config.0.horizontal_pod_autoscaling",
		"addons_config.0.network_policy_config",
	}
)

func resourceContainerCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerClusterCreate,
		Read:   resourceContainerClusterRead,
		Update: resourceContainerClusterUpdate,
		Delete: resourceContainerClusterDelete,

		CustomizeDiff: customdiff.All(
			resourceNodeConfigEmptyGuestAccelerator,
			containerClusterPrivateClusterConfigCustomDiff,
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
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

			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use location instead",
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use location instead",
			},

			"node_locations": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"additional_zones": {
				Type:     schema.TypeSet,
				Optional: true,
				Removed:  "Use node_locations instead",
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
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"horizontal_pod_autoscaling": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"kubernetes_dashboard": {
							Type:     schema.TypeList,
							Optional: true,
							Removed:  "The Kubernetes Dashboard addon is removed for clusters on GKE.",
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
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"cluster_autoscaling": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// This field is Optional + Computed because we automatically set the
				// enabled value to false if the block is not returned in API responses.
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"resource_limits": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"minimum": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"maximum": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"auto_provisioning_defaults": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oauth_scopes": {
										Type:             schema.TypeList,
										Optional:         true,
										Elem:             &schema.Schema{Type: schema.TypeString},
										DiffSuppressFunc: containerClusterAddedScopesSuppress,
										ExactlyOneOf: []string{
											"cluster_autoscaling.0.auto_provisioning_defaults.0.oauth_scopes",
											"cluster_autoscaling.0.auto_provisioning_defaults.0.service_account",
										},
									},
									"service_account": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "default",
										ExactlyOneOf: []string{
											"cluster_autoscaling.0.auto_provisioning_defaults.0.oauth_scopes",
											"cluster_autoscaling.0.auto_provisioning_defaults.0.service_account",
										},
									},
								},
							},
						},
					},
				},
			},

			"cluster_ipv4_cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  orEmpty(validateRFC1918Network(8, 32)),
				ConflictsWith: []string{"ip_allocation_policy"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"enable_binary_authorization": {
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
				Computed: true,
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_kubernetes_alpha": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"enable_tpu": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
				Computed: true,
			},

			"enable_legacy_abac": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"authenticator_groups_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"initial_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"logging_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "logging.googleapis.com/kubernetes",
				ValidateFunc: validation.StringInSlice([]string{"logging.googleapis.com", "logging.googleapis.com/kubernetes", "none"}, false),
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
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"master_auth.0.password", "master_auth.0.username", "master_auth.0.client_certificate_config"},
							Sensitive:    true,
						},

						"username": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"master_auth.0.password", "master_auth.0.username", "master_auth.0.client_certificate_config"},
						},

						// Ideally, this would be Optional (and not Computed).
						// In past versions (incl. 2.X series) of the provider
						// though, being unset was considered identical to set
						// and the issue_client_certificate value being true.
						"client_certificate_config": {
							Type:         schema.TypeList,
							MaxItems:     1,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: []string{"master_auth.0.password", "master_auth.0.username", "master_auth.0.client_certificate_config"},
							ForceNew:     true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issue_client_certificate": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "monitoring.googleapis.com/kubernetes",
				ValidateFunc: validation.StringInSlice([]string{"monitoring.googleapis.com", "monitoring.googleapis.com/kubernetes", "none"}, false),
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
							Required: true,
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
				// Remove return nil from expand when this is removed for good.
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
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

			"services_ipv4_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_allocation_policy": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"cluster_ipv4_cidr"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// GKE creates/deletes secondary ranges in VPC
						"cluster_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: cidrOrSizeDiffSuppress,
						},

						"services_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: cidrOrSizeDiffSuppress,
						},

						// User manages secondary ranges manually
						"cluster_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
						},

						"services_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
						},

						"use_ip_aliases": {
							Type:     schema.TypeBool,
							Removed:  "This field is removed as of 3.0.0. If previously set to true, remove it from your config. If false, remove it.",
							Computed: true,
							Optional: true,
						},

						// GKE creates subnetwork automatically
						"create_subnetwork": {
							Type:     schema.TypeBool,
							Removed:  "This field is removed as of 3.0.0. Define an explicit google_compute_subnetwork and use subnetwork instead.",
							Computed: true,
							Optional: true,
						},

						"subnetwork_name": {
							Type:     schema.TypeString,
							Removed:  "This field is removed as of 3.0.0. Define an explicit google_compute_subnetwork and use subnetwork instead.",
							Computed: true,
							Optional: true,
						},

						"node_ipv4_cidr_block": {
							Type:     schema.TypeString,
							Removed:  "This field is removed as of 3.0.0. Define an explicit google_compute_subnetwork and use subnetwork instead.",
							Computed: true,
							Optional: true,
						},
					},
				},
			},

			"remove_default_node_pool": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"private_cluster_config": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_private_endpoint": {
							Type:             schema.TypeBool,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
						},
						"enable_private_nodes": {
							Type:             schema.TypeBool,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
						},
						"master_ipv4_cidr_block": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: orEmpty(validation.CIDRNetwork(28, 28)),
						},

						"private_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"resource_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"default_max_pods_per_node": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vertical_pod_autoscaling": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"enable_intranode_visibility": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
			},
		},
	}
}

// Setting a guest accelerator block to count=0 is the equivalent to omitting the block: it won't get
// sent to the API and it won't be stored in state. This diffFunc will try to compare the old + new state
// by only comparing the blocks with a positive count and ignoring those with count=0
//
// One quirk with this approach is that configs with mixed count=0 and count>0 accelerator blocks will
// show a confusing diff if one of there are config changes that result in a legitimate diff as the count=0
// blocks will not be in state.
//
// This could also be modelled by setting `guest_accelerator = []` in the config. However since the
// previous syntax requires that schema.SchemaConfigModeAttr is set on the field it is advisable that
// we have a work around for removing guest accelerators. Also Terraform 0.11 cannot use dynamic blocks
// so this isn't a solution for module authors who want to dynamically omit guest accelerators
// See https://github.com/terraform-providers/terraform-provider-google/issues/3786
func resourceNodeConfigEmptyGuestAccelerator(diff *schema.ResourceDiff, meta interface{}) error {
	old, new := diff.GetChange("node_config.0.guest_accelerator")
	oList := old.([]interface{})
	nList := new.([]interface{})

	if len(nList) == len(oList) || len(nList) == 0 {
		return nil
	}
	var hasAcceleratorWithEmptyCount bool
	// the list of blocks in a desired state with count=0 accelerator blocks in it
	// will be longer than the current state.
	// this index tracks the location of positive count accelerator blocks
	index := 0
	for i, item := range nList {
		accel := item.(map[string]interface{})
		if accel["count"].(int) == 0 {
			hasAcceleratorWithEmptyCount = true
			// Ignore any 'empty' accelerators because they aren't sent to the API
			continue
		}
		if index >= len(oList) {
			// Return early if there are more positive count accelerator blocks in the desired state
			// than the current state since a difference in 'legit' blocks is a valid diff.
			// This will prevent array index overruns
			return nil
		}
		if !reflect.DeepEqual(nList[i], oList[index]) {
			return nil
		}
		index += 1
	}

	if hasAcceleratorWithEmptyCount && index == len(oList) {
		// If the number of count>0 blocks match, there are count=0 blocks present and we
		// haven't already returned due to a legitimate diff
		err := diff.Clear("node_config.0.guest_accelerator")
		if err != nil {
			return err
		}
	}

	return nil
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
		Name:                           clusterName,
		InitialNodeCount:               int64(d.Get("initial_node_count").(int)),
		MaintenancePolicy:              expandMaintenancePolicy(d, meta),
		MasterAuthorizedNetworksConfig: expandMasterAuthorizedNetworksConfig(d.Get("master_authorized_networks_config")),
		InitialClusterVersion:          d.Get("min_master_version").(string),
		ClusterIpv4Cidr:                d.Get("cluster_ipv4_cidr").(string),
		Description:                    d.Get("description").(string),
		LegacyAbac: &containerBeta.LegacyAbac{
			Enabled:         d.Get("enable_legacy_abac").(bool),
			ForceSendFields: []string{"Enabled"},
		},
		LoggingService:          d.Get("logging_service").(string),
		MonitoringService:       d.Get("monitoring_service").(string),
		NetworkPolicy:           expandNetworkPolicy(d.Get("network_policy")),
		AddonsConfig:            expandClusterAddonsConfig(d.Get("addons_config")),
		EnableKubernetesAlpha:   d.Get("enable_kubernetes_alpha").(bool),
		IpAllocationPolicy:      expandIPAllocationPolicy(d.Get("ip_allocation_policy")),
		PodSecurityPolicyConfig: expandPodSecurityPolicyConfig(d.Get("pod_security_policy_config")),
		Autoscaling:             expandClusterAutoscaling(d.Get("cluster_autoscaling"), d),
		MasterAuth:              expandMasterAuth(d.Get("master_auth")),
		ResourceLabels:          expandStringMap(d, "resource_labels"),
	}

	if v, ok := d.GetOk("default_max_pods_per_node"); ok {
		cluster.DefaultMaxPodsConstraint = expandDefaultMaxPodsConstraint(v)
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

	if v, ok := d.GetOk("node_locations"); ok {
		locationsSet := v.(*schema.Set)
		if locationsSet.Contains(location) {
			return fmt.Errorf("when using a multi-zonal cluster, node_locations should not contain the original 'zone'")
		}

		// GKE requires a full list of node locations
		// but when using a multi-zonal cluster our schema only asks for the
		// additional zones, so append the cluster location if it's a zone
		if isZone(location) {
			locationsSet.Add(location)
		}
		cluster.Locations = convertStringSet(locationsSet)
	}

	if v, ok := d.GetOk("network"); ok {
		network, err := ParseNetworkFieldValue(v.(string), d, config)
		if err != nil {
			return err
		}
		cluster.Network = network.RelativeLink()
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		subnetwork, err := ParseSubnetworkFieldValue(v.(string), d, config)
		if err != nil {
			return err
		}
		cluster.Subnetwork = subnetwork.RelativeLink()
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

	if v, ok := d.GetOk("authenticator_groups_config"); ok {
		cluster.AuthenticatorGroupsConfig = expandAuthenticatorGroupsConfig(v)
	}

	if v, ok := d.GetOk("private_cluster_config"); ok {
		cluster.PrivateClusterConfig = expandPrivateClusterConfig(v)
	}

	if v, ok := d.GetOk("vertical_pod_autoscaling"); ok {
		cluster.VerticalPodAutoscaling = expandVerticalPodAutoscaling(v)
	}

	req := &containerBeta.CreateClusterRequest{
		Cluster: cluster,
	}

	mutexKV.Lock(containerClusterMutexKey(project, location, clusterName))
	defer mutexKV.Unlock(containerClusterMutexKey(project, location, clusterName))

	parent := fmt.Sprintf("projects/%s/locations/%s", project, location)
	var op *containerBeta.Operation
	err = retry(func() error {
		op, err = config.clientContainerBeta.Projects.Locations.Clusters.Create(parent, req).Do()
		return err
	})
	if err != nil {
		return err
	}

	d.SetId(containerClusterFullName(project, location, clusterName))

	// Wait until it's created
	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := containerOperationWait(config, op, project, location, "creating GKE cluster", timeoutInMinutes)
	if waitErr != nil {
		// Try a GET on the cluster so we can see the state in debug logs. This will help classify error states.
		_, getErr := config.clientContainerBeta.Projects.Locations.Clusters.Get(containerClusterFullName(project, location, clusterName)).Do()
		if getErr != nil {
			log.Printf("[WARN] Cluster %s was created in an error state and not found", clusterName)
			d.SetId("")
		}

		if deleteErr := cleanFailedContainerCluster(d, meta); deleteErr != nil {
			log.Printf("[WARN] Unable to clean up cluster from failed creation: %s", deleteErr)
			// Leave ID set as the cluster likely still exists and should not be removed from state yet.
		} else {
			log.Printf("[WARN] Verified failed creation of cluster %s was cleaned up", d.Id())
			d.SetId("")
		}
		// The resource didn't actually create
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been created", clusterName)

	if d.Get("remove_default_node_pool").(bool) {
		parent := fmt.Sprintf("%s/nodePools/%s", containerClusterFullName(project, location, clusterName), "default-pool")
		err = retry(func() error {
			op, err = config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(parent).Do()
			return err
		})
		if err != nil {
			return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
		}
		err = containerOperationWait(config, op, project, location, "removing default node pool", timeoutInMinutes)
		if err != nil {
			return errwrap.Wrapf("Error while waiting to delete default node pool: {{err}}", err)
		}
	}

	if err := resourceContainerClusterRead(d, meta); err != nil {
		return err
	}

	if err := waitForContainerClusterReady(config, project, location, clusterName, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}
	return nil
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

	clusterName := d.Get("name").(string)
	name := containerClusterFullName(project, location, clusterName)
	cluster, err := config.clientContainerBeta.Projects.Locations.Clusters.Get(name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Container Cluster %q", d.Get("name").(string)))
	}
	if cluster.Status == "ERROR" || cluster.Status == "DEGRADED" {
		return fmt.Errorf("Cluster %q has status %q with message %q", d.Get("name"), cluster.Status, cluster.StatusMessage)
	}

	d.Set("name", cluster.Name)
	if err := d.Set("network_policy", flattenNetworkPolicy(cluster.NetworkPolicy)); err != nil {
		return err
	}

	d.Set("location", cluster.Location)

	locations := schema.NewSet(schema.HashString, convertStringArrToInterface(cluster.Locations))
	locations.Remove(cluster.Zone) // Remove the original zone since we only store additional zones
	d.Set("node_locations", locations)

	d.Set("endpoint", cluster.Endpoint)
	if err := d.Set("maintenance_policy", flattenMaintenancePolicy(cluster.MaintenancePolicy)); err != nil {
		return err
	}
	if err := d.Set("master_auth", flattenMasterAuth(cluster.MasterAuth)); err != nil {
		return err
	}
	if err := d.Set("master_authorized_networks_config", flattenMasterAuthorizedNetworksConfig(cluster.MasterAuthorizedNetworksConfig)); err != nil {
		return err
	}
	d.Set("initial_node_count", cluster.InitialNodeCount)
	d.Set("master_version", cluster.CurrentMasterVersion)
	d.Set("node_version", cluster.CurrentNodeVersion)
	d.Set("cluster_ipv4_cidr", cluster.ClusterIpv4Cidr)
	d.Set("services_ipv4_cidr", cluster.ServicesIpv4Cidr)
	d.Set("description", cluster.Description)
	d.Set("enable_kubernetes_alpha", cluster.EnableKubernetesAlpha)
	d.Set("enable_legacy_abac", cluster.LegacyAbac.Enabled)
	d.Set("logging_service", cluster.LoggingService)
	d.Set("monitoring_service", cluster.MonitoringService)
	d.Set("network", cluster.NetworkConfig.Network)
	d.Set("subnetwork", cluster.NetworkConfig.Subnetwork)
	if err := d.Set("cluster_autoscaling", flattenClusterAutoscaling(cluster.Autoscaling)); err != nil {
		return err
	}
	if err := d.Set("authenticator_groups_config", flattenAuthenticatorGroupsConfig(cluster.AuthenticatorGroupsConfig)); err != nil {
		return err
	}
	if cluster.DefaultMaxPodsConstraint != nil {
		d.Set("default_max_pods_per_node", cluster.DefaultMaxPodsConstraint.MaxPodsPerNode)
	}
	if err := d.Set("node_config", flattenNodeConfig(cluster.NodeConfig)); err != nil {
		return err
	}
	d.Set("project", project)
	if err := d.Set("addons_config", flattenClusterAddonsConfig(cluster.AddonsConfig)); err != nil {
		return err
	}
	nps, err := flattenClusterNodePools(d, config, cluster.NodePools)
	if err != nil {
		return err
	}
	if err := d.Set("node_pool", nps); err != nil {
		return err
	}

	if err := d.Set("ip_allocation_policy", flattenIPAllocationPolicy(cluster, d, config)); err != nil {
		return err
	}

	if err := d.Set("private_cluster_config", flattenPrivateClusterConfig(cluster.PrivateClusterConfig)); err != nil {
		return err
	}

	igUrls, err := getInstanceGroupUrlsFromManagerUrls(config, cluster.InstanceGroupUrls)
	if err != nil {
		return err
	}
	if err := d.Set("instance_group_urls", igUrls); err != nil {
		return err
	}

	if err := d.Set("vertical_pod_autoscaling", flattenVerticalPodAutoscaling(cluster.VerticalPodAutoscaling)); err != nil {
		return err
	}

	d.Set("resource_labels", cluster.ResourceLabels)

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

	if err := waitForContainerClusterReady(config, project, location, clusterName, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

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
			return containerOperationWait(config, op, project, location, updateDescription, timeoutInMinutes)
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

	if d.HasChange("cluster_autoscaling") {
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredClusterAutoscaling: expandClusterAutoscaling(d.Get("cluster_autoscaling"), d),
			}}

		updateF := updateFunc(req, "updating GKE cluster autoscaling")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s's cluster-wide autoscaling has been updated", d.Id())

		d.SetPartial("cluster_autoscaling")
	}

	if d.HasChange("maintenance_policy") {
		req := &containerBeta.SetMaintenancePolicyRequest{
			MaintenancePolicy: expandMaintenancePolicy(d, meta),
		}

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetMaintenancePolicy(name, req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE cluster maintenance policy", timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s maintenance policy has been updated", d.Id())

		d.SetPartial("maintenance_policy")
	}

	if d.HasChange("node_locations") {
		azSetOldI, azSetNewI := d.GetChange("node_locations")
		azSetNew := azSetNewI.(*schema.Set)
		azSetOld := azSetOldI.(*schema.Set)
		if azSetNew.Contains(location) {
			return fmt.Errorf("for multi-zonal clusters, node_locations should not contain the primary 'zone'")
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

		updateF := updateFunc(req, "updating GKE cluster node locations")
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

			updateF := updateFunc(req, "updating GKE cluster node locations")
			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}
		}

		log.Printf("[INFO] GKE cluster %s node locations have been updated to %v", d.Id(), azSet.List())

		d.SetPartial("node_locations")
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
			err = containerOperationWait(config, op, project, location, "updating GKE legacy ABAC", timeoutInMinutes)
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

	if d.HasChange("monitoring_service") || d.HasChange("logging_service") {
		logging := d.Get("logging_service").(string)
		monitoring := d.Get("monitoring_service").(string)

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredMonitoringService: monitoring,
					DesiredLoggingService:    logging,
				},
			}
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE logging+monitoring service", timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s: logging service has been updated to %s, monitoring service has been updated to %s", d.Id(), logging, monitoring)
		d.SetPartial("logging_service")
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
			err = containerOperationWait(config, op, project, location, "updating GKE cluster network policy", timeoutInMinutes)
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

	// It's not super important that this come after updating the node pools, but it still seems like a better
	// idea than doing it before.
	if d.HasChange("node_version") {
		foundDefault := false
		if n, ok := d.GetOk("node_pool.#"); ok {
			for i := 0; i < n.(int); i++ {
				key := fmt.Sprintf("node_pool.%d.", i)
				if d.Get(key+"name").(string) == "default-pool" {
					desiredNodeVersion := d.Get("node_version").(string)
					req := &containerBeta.UpdateClusterRequest{
						Update: &containerBeta.ClusterUpdate{
							DesiredNodeVersion: desiredNodeVersion,
							DesiredNodePoolId:  "default-pool",
						},
					}
					updateF := updateFunc(req, "updating GKE default node pool node version")
					// Call update serially.
					if err := lockedCall(lockKey, updateF); err != nil {
						return err
					}
					log.Printf("[INFO] GKE cluster %s: default node pool has been updated to %s", d.Id(),
						desiredNodeVersion)
					foundDefault = true
				}
			}
		}

		if !foundDefault {
			return fmt.Errorf("node_version was updated but default-pool was not found. To update the version for a non-default pool, use the version attribute on that pool.")
		}

		d.SetPartial("node_version")
	}

	if d.HasChange("node_config") {
		if d.HasChange("node_config.0.image_type") {
			it := d.Get("node_config.0.image_type").(string)
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredImageType: it,
				},
			}

			updateF := func() error {
				name := containerClusterFullName(project, location, clusterName)
				op, err := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req).Do()
				if err != nil {
					return err
				}

				// Wait until it's updated
				return containerOperationWait(config, op, project, location, "updating GKE image type", timeoutInMinutes)
			}

			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}

			log.Printf("[INFO] GKE cluster %s: image type has been updated to %s", d.Id(), it)
		}
		d.SetPartial("node_config")
	}

	if d.HasChange("master_auth") {
		var req *containerBeta.SetMasterAuthRequest
		if ma, ok := d.GetOk("master_auth"); ok {
			req = &containerBeta.SetMasterAuthRequest{
				Action: "SET_USERNAME",
				Update: expandMasterAuth(ma),
			}
		} else {
			req = &containerBeta.SetMasterAuthRequest{
				Action: "SET_USERNAME",
				Update: &containerBeta.MasterAuth{
					Username: "admin",
				},
			}
		}

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetMasterAuth(name, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating master auth", timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s: master auth has been updated", d.Id())
		d.SetPartial("master_auth")
	}

	if d.HasChange("vertical_pod_autoscaling") {
		if ac, ok := d.GetOk("vertical_pod_autoscaling"); ok {
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredVerticalPodAutoscaling: expandVerticalPodAutoscaling(ac),
				},
			}

			updateF := updateFunc(req, "updating GKE cluster vertical pod autoscaling")
			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}

			log.Printf("[INFO] GKE cluster %s vertical pod autoscaling has been updated", d.Id())

			d.SetPartial("vertical_pod_autoscaling")
		}
	}

	if d.HasChange("resource_labels") {
		resourceLabels := d.Get("resource_labels").(map[string]interface{})
		req := &containerBeta.SetLabelsRequest{
			ResourceLabels: convertStringMap(resourceLabels),
		}
		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.SetResourceLabels(name, req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE resource labels", timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		d.SetPartial("resource_labels")
	}

	if d.HasChange("remove_default_node_pool") && d.Get("remove_default_node_pool").(bool) {
		name := fmt.Sprintf("%s/nodePools/%s", containerClusterFullName(project, location, clusterName), "default-pool")
		op, err := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(name).Do()
		if err != nil {
			if !isGoogleApiErrorWithCode(err, 404) {
				return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
			}
			log.Printf("[WARN] Container cluster %q default node pool already removed, no change", d.Id())
		} else {
			err = containerOperationWait(config, op, project, location, "removing default node pool", timeoutInMinutes)
			if err != nil {
				return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
			}
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

	location, err := getLocation(d, config)
	if err != nil {
		return err
	}

	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	if err := waitForContainerClusterReady(config, project, location, clusterName, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting GKE cluster %s", d.Get("name").(string))
	mutexKV.Lock(containerClusterMutexKey(project, location, clusterName))
	defer mutexKV.Unlock(containerClusterMutexKey(project, location, clusterName))

	var op *containerBeta.Operation
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
	waitErr := containerOperationWait(config, op, project, location, "deleting GKE cluster", timeoutInMinutes)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been deleted", d.Id())

	d.SetId("")

	return nil
}

// cleanFailedContainerCluster deletes clusters that failed but were
// created in an error state. Similar to resourceContainerClusterDelete
// but implemented in separate function as it doesn't try to lock already
// locked cluster state, does different error handling, and doesn't do retries.
func cleanFailedContainerCluster(d *schema.ResourceData, meta interface{}) error {
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
	fullName := containerClusterFullName(project, location, clusterName)

	log.Printf("[DEBUG] Cleaning up failed GKE cluster %s", d.Get("name").(string))
	op, err := config.clientContainerBeta.Projects.Locations.Clusters.Delete(fullName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Container Cluster %q", d.Get("name").(string)))
	}

	// Wait until it's deleted
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())
	waitErr := containerOperationWait(config, op, project, location, "deleting GKE cluster", timeoutInMinutes)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

func waitForContainerClusterReady(config *Config, project, location, clusterName string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		name := containerClusterFullName(project, location, clusterName)
		cluster, err := config.clientContainerBeta.Projects.Locations.Clusters.Get(name).Do()
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if cluster.Status == "PROVISIONING" || cluster.Status == "RECONCILING" || cluster.Status == "STOPPING" {
			return resource.RetryableError(fmt.Errorf("Cluster %q has status %q with message %q", clusterName, cluster.Status, cluster.StatusMessage))
		} else if cluster.Status == "RUNNING" {
			log.Printf("Cluster %q has status 'RUNNING'.", clusterName)
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("Cluster %q has terminal state %q with message %q.", clusterName, cluster.Status, cluster.StatusMessage))
		}
	})
}

// container engine's API returns the instance group manager's URL instead of the instance
// group's URL in its responses, while the field is named as if it should have been the group
// and not the manager. This shim should be supported for backwards compatibility reasons.
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
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	config := l[0].(map[string]interface{})
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

	if v, ok := config["network_policy_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.NetworkPolicyConfig = &containerBeta.NetworkPolicyConfig{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
	}

	return ac
}

func expandIPAllocationPolicy(configured interface{}) *containerBeta.IPAllocationPolicy {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &containerBeta.IPAllocationPolicy{
			UseIpAliases:    false,
			ForceSendFields: []string{"UseIpAliases"},
		}
	}

	config := l[0].(map[string]interface{})
	return &containerBeta.IPAllocationPolicy{
		UseIpAliases:          true,
		ClusterIpv4CidrBlock:  config["cluster_ipv4_cidr_block"].(string),
		ServicesIpv4CidrBlock: config["services_ipv4_cidr_block"].(string),

		ClusterSecondaryRangeName:  config["cluster_secondary_range_name"].(string),
		ServicesSecondaryRangeName: config["services_secondary_range_name"].(string),
		ForceSendFields:            []string{"UseIpAliases"},
	}
}

func expandMaintenancePolicy(d *schema.ResourceData, meta interface{}) *containerBeta.MaintenancePolicy {
	config := meta.(*Config)
	// We have to perform a full Get() as part of this, to get the fingerprint.  We can't do this
	// at any other time, because the fingerprint update might happen between plan and apply.
	// We can omit error checks, since to have gotten this far, a project is definitely configured.
	project, _ := getProject(d, config)
	location, _ := getLocation(d, config)
	clusterName := d.Get("name").(string)
	name := containerClusterFullName(project, location, clusterName)
	cluster, _ := config.clientContainerBeta.Projects.Locations.Clusters.Get(name).Do()
	resourceVersion := ""
	// If the cluster doesn't exist or if there is a read error of any kind, we will pass in an empty
	// resourceVersion.  If there happens to be a change to maintenance policy, we will fail at that
	// point.  This is a compromise between code cleanliness and a slightly worse user experience in
	// an unlikely error case - we choose code cleanliness.
	if cluster != nil && cluster.MaintenancePolicy != nil {
		resourceVersion = cluster.MaintenancePolicy.ResourceVersion
	}
	exclusions := make(map[string]containerBeta.TimeWindow)
	if cluster != nil && cluster.MaintenancePolicy != nil && cluster.MaintenancePolicy.Window != nil {
		exclusions = cluster.MaintenancePolicy.Window.MaintenanceExclusions
	}

	configured := d.Get("maintenance_policy")
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &containerBeta.MaintenancePolicy{
			ResourceVersion: resourceVersion,
			Window: &containerBeta.MaintenanceWindow{
				MaintenanceExclusions: exclusions,
			},
		}
	}
	maintenancePolicy := l[0].(map[string]interface{})

	if dailyMaintenanceWindow, ok := maintenancePolicy["daily_maintenance_window"]; ok && len(dailyMaintenanceWindow.([]interface{})) > 0 {
		dmw := dailyMaintenanceWindow.([]interface{})[0].(map[string]interface{})
		startTime := dmw["start_time"].(string)
		return &containerBeta.MaintenancePolicy{
			Window: &containerBeta.MaintenanceWindow{
				MaintenanceExclusions: exclusions,
				DailyMaintenanceWindow: &containerBeta.DailyMaintenanceWindow{
					StartTime: startTime,
				},
			},
			ResourceVersion: resourceVersion,
		}
	}

	return nil
}

func expandClusterAutoscaling(configured interface{}, d *schema.ResourceData) *containerBeta.ClusterAutoscaling {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &containerBeta.ClusterAutoscaling{
			EnableNodeAutoprovisioning: false,
			ForceSendFields:            []string{"EnableNodeAutoprovisioning"},
		}
	}

	config := l[0].(map[string]interface{})

	// Conditionally provide an empty list to preserve a legacy 2.X behaviour
	// when `enabled` is false and resource_limits is unset, allowing users to
	// explicitly disable the feature. resource_limits don't work when node
	// auto-provisioning is disabled at time of writing. This may change API-side
	// in the future though, as the feature is intended to apply to both node
	// auto-provisioning and node autoscaling.
	var resourceLimits []*containerBeta.ResourceLimit
	if limits, ok := config["resource_limits"]; ok {
		resourceLimits = make([]*containerBeta.ResourceLimit, 0)
		if lmts, ok := limits.([]interface{}); ok {
			for _, v := range lmts {
				limit := v.(map[string]interface{})
				resourceLimits = append(resourceLimits,
					&containerBeta.ResourceLimit{
						ResourceType: limit["resource_type"].(string),
						// Here we're relying on *not* setting ForceSendFields for 0-values.
						Minimum: int64(limit["minimum"].(int)),
						Maximum: int64(limit["maximum"].(int)),
					})
			}
		}
	}
	return &containerBeta.ClusterAutoscaling{
		EnableNodeAutoprovisioning:       config["enabled"].(bool),
		ResourceLimits:                   resourceLimits,
		AutoprovisioningNodePoolDefaults: expandAutoProvisioningDefaults(config["auto_provisioning_defaults"], d),
	}
}

func expandAutoProvisioningDefaults(configured interface{}, d *schema.ResourceData) *containerBeta.AutoprovisioningNodePoolDefaults {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &containerBeta.AutoprovisioningNodePoolDefaults{}
	}
	config := l[0].(map[string]interface{})

	return &containerBeta.AutoprovisioningNodePoolDefaults{
		OauthScopes:    convertStringArr(config["oauth_scopes"].([]interface{})),
		ServiceAccount: config["service_account"].(string),
	}
}

func expandAuthenticatorGroupsConfig(configured interface{}) *containerBeta.AuthenticatorGroupsConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	result := &containerBeta.AuthenticatorGroupsConfig{}
	config := l[0].(map[string]interface{})
	if securityGroup, ok := config["security_group"]; ok {
		result.Enabled = true
		result.SecurityGroup = securityGroup.(string)
	}
	return result
}

func expandMasterAuth(configured interface{}) *containerBeta.MasterAuth {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	masterAuth := l[0].(map[string]interface{})
	result := &containerBeta.MasterAuth{
		Username: masterAuth["username"].(string),
		Password: masterAuth["password"].(string),
	}

	if v, ok := masterAuth["client_certificate_config"]; ok {
		if len(v.([]interface{})) > 0 {
			clientCertificateConfig := masterAuth["client_certificate_config"].([]interface{})[0].(map[string]interface{})

			result.ClientCertificateConfig = &containerBeta.ClientCertificateConfig{
				IssueClientCertificate: clientCertificateConfig["issue_client_certificate"].(bool),
			}
		}
	}

	return result
}

func expandMasterAuthorizedNetworksConfig(configured interface{}) *containerBeta.MasterAuthorizedNetworksConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return &containerBeta.MasterAuthorizedNetworksConfig{
			Enabled: false,
		}
	}
	result := &containerBeta.MasterAuthorizedNetworksConfig{
		Enabled: true,
	}
	if config, ok := l[0].(map[string]interface{}); ok {
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
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	result := &containerBeta.NetworkPolicy{}
	config := l[0].(map[string]interface{})
	if enabled, ok := config["enabled"]; ok && enabled.(bool) {
		result.Enabled = true
		if provider, ok := config["provider"]; ok {
			result.Provider = provider.(string)
		}
	}
	return result
}

func expandPrivateClusterConfig(configured interface{}) *containerBeta.PrivateClusterConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &containerBeta.PrivateClusterConfig{
		EnablePrivateEndpoint: config["enable_private_endpoint"].(bool),
		EnablePrivateNodes:    config["enable_private_nodes"].(bool),
		MasterIpv4CidrBlock:   config["master_ipv4_cidr_block"].(string),
		ForceSendFields:       []string{"EnablePrivateEndpoint", "EnablePrivateNodes", "MasterIpv4CidrBlock"},
	}
}

func expandVerticalPodAutoscaling(configured interface{}) *containerBeta.VerticalPodAutoscaling {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &containerBeta.VerticalPodAutoscaling{
		Enabled: config["enabled"].(bool),
	}
}

func expandPodSecurityPolicyConfig(configured interface{}) *containerBeta.PodSecurityPolicyConfig {
	// Removing lists is hard - the element count (#) will have a diff from nil -> computed
	// If we set this to empty on Read, it will be stable.
	return nil
}

func expandDefaultMaxPodsConstraint(v interface{}) *containerBeta.MaxPodsConstraint {
	if v == nil {
		return nil
	}

	return &containerBeta.MaxPodsConstraint{
		MaxPodsPerNode: int64(v.(int)),
	}
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
	if c == nil {
		return nil
	}
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

func flattenAuthenticatorGroupsConfig(c *containerBeta.AuthenticatorGroupsConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"security_group": c.SecurityGroup,
		},
	}
}

func flattenPrivateClusterConfig(c *containerBeta.PrivateClusterConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enable_private_endpoint": c.EnablePrivateEndpoint,
			"enable_private_nodes":    c.EnablePrivateNodes,
			"master_ipv4_cidr_block":  c.MasterIpv4CidrBlock,
			"private_endpoint":        c.PrivateEndpoint,
			"public_endpoint":         c.PublicEndpoint,
		},
	}
}

func flattenVerticalPodAutoscaling(c *containerBeta.VerticalPodAutoscaling) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"enabled": c.Enabled,
		},
	}
}

func flattenIPAllocationPolicy(c *containerBeta.Cluster, d *schema.ResourceData, config *Config) []map[string]interface{} {
	// If IP aliasing isn't enabled, none of the values in this block can be set.
	if c == nil || c.IpAllocationPolicy == nil || !c.IpAllocationPolicy.UseIpAliases {
		return nil
	}

	p := c.IpAllocationPolicy
	return []map[string]interface{}{
		{
			"cluster_ipv4_cidr_block":       p.ClusterIpv4CidrBlock,
			"services_ipv4_cidr_block":      p.ServicesIpv4CidrBlock,
			"cluster_secondary_range_name":  p.ClusterSecondaryRangeName,
			"services_secondary_range_name": p.ServicesSecondaryRangeName,
		},
	}
}

func flattenMaintenancePolicy(mp *containerBeta.MaintenancePolicy) []map[string]interface{} {
	if mp == nil || mp.Window == nil {
		return nil
	}
	if mp.Window.DailyMaintenanceWindow != nil {
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

	return nil
}

func flattenMasterAuth(ma *containerBeta.MasterAuth) []map[string]interface{} {
	if ma == nil {
		return nil
	}
	masterAuth := []map[string]interface{}{
		{
			"username":               ma.Username,
			"password":               ma.Password,
			"client_certificate":     ma.ClientCertificate,
			"client_key":             ma.ClientKey,
			"cluster_ca_certificate": ma.ClusterCaCertificate,
		},
	}

	// No version of the GKE API returns the client_certificate_config value.
	// Instead, we need to infer whether or not it was set based on the
	// client cert being returned from the API or not.
	// Previous versions of the provider didn't record anything in state when
	// a client cert was enabled, only setting the block when it was false.
	masterAuth[0]["client_certificate_config"] = []map[string]interface{}{
		{
			"issue_client_certificate": len(ma.ClientCertificate) != 0,
		},
	}

	return masterAuth
}

func flattenClusterAutoscaling(a *containerBeta.ClusterAutoscaling) []map[string]interface{} {
	r := make(map[string]interface{})
	if a == nil || !a.EnableNodeAutoprovisioning {
		r["enabled"] = false
	} else {
		resourceLimits := make([]interface{}, 0, len(a.ResourceLimits))
		for _, rl := range a.ResourceLimits {
			resourceLimits = append(resourceLimits, map[string]interface{}{
				"resource_type": rl.ResourceType,
				"minimum":       rl.Minimum,
				"maximum":       rl.Maximum,
			})
		}
		r["resource_limits"] = resourceLimits
		r["enabled"] = true
		r["auto_provisioning_defaults"] = flattenAutoProvisioningDefaults(a.AutoprovisioningNodePoolDefaults)
	}
	return []map[string]interface{}{r}
}

func flattenAutoProvisioningDefaults(a *containerBeta.AutoprovisioningNodePoolDefaults) []map[string]interface{} {
	r := make(map[string]interface{})
	r["oauth_scopes"] = a.OauthScopes
	r["service_account"] = a.ServiceAccount

	return []map[string]interface{}{r}
}

func flattenMasterAuthorizedNetworksConfig(c *containerBeta.MasterAuthorizedNetworksConfig) []map[string]interface{} {
	if c == nil || !c.Enabled {
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

func resourceContainerClusterStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/clusters/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)", "(?P<location>[^/]+)/(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return nil, err
	}

	clusterName := d.Get("name").(string)

	d.Set("location", location)
	if err := waitForContainerClusterReady(config, project, location, clusterName, d.Timeout(schema.TimeoutCreate)); err != nil {
		return nil, err
	}

	d.SetId(containerClusterFullName(project, location, clusterName))

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

func cidrOrSizeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// If the user specified a size and the API returned a full cidr block, suppress.
	return strings.HasPrefix(new, "/") && strings.HasSuffix(old, new)
}

// Suppress unremovable default scope values from GCP.
// If the default service account would not otherwise have it, the `monitoring.write` scope
// is added to a GKE cluster's scopes regardless of what the user provided.
// monitoring.write is inherited from monitoring (rw) and cloud-platform, so it won't always
// be present.
// Enabling Stackdriver features through logging_service and monitoring_service may enable
// monitoring or logging.write. We've chosen not to suppress in those cases because they're
// removable by disabling those features.
func containerClusterAddedScopesSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("cluster_autoscaling.0.auto_provisioning_defaults.0.oauth_scopes")
	if o == nil || n == nil {
		return false
	}

	addedScopes := []string{
		"https://www.googleapis.com/auth/monitoring.write",
	}

	// combine what the default scopes are with what was passed
	m := golangSetFromStringSlice(append(addedScopes, convertStringArr(n.([]interface{}))...))
	combined := stringSliceFromGolangSet(m)

	// compare if the combined new scopes and default scopes differ from the old scopes
	if len(combined) != len(convertStringArr(o.([]interface{}))) {
		return false
	}

	for _, i := range combined {
		if stringInSlice(convertStringArr(o.([]interface{})), i) {
			continue
		}

		return false
	}

	return true
}

// We want to suppress diffs for empty/disabled private cluster config.
func containerClusterPrivateClusterConfigSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("private_cluster_config.0.enable_private_endpoint")
	suppressEndpoint := !o.(bool) && !n.(bool)

	o, n = d.GetChange("private_cluster_config.0.enable_private_nodes")
	suppressNodes := !o.(bool) && !n.(bool)

	if k == "private_cluster_config.0.enable_private_endpoint" {
		return suppressEndpoint
	} else if k == "private_cluster_config.0.enable_private_nodes" {
		return suppressNodes
	} else if k == "private_cluster_config.#" {
		return suppressEndpoint && suppressNodes
	}
	return false
}

func containerClusterPrivateClusterConfigCustomDiff(d *schema.ResourceDiff, meta interface{}) error {
	pcc, ok := d.GetOk("private_cluster_config")
	if !ok {
		return nil
	}
	pccList := pcc.([]interface{})
	if len(pccList) == 0 {
		return nil
	}
	config := pccList[0].(map[string]interface{})
	if config["enable_private_nodes"].(bool) {
		block := config["master_ipv4_cidr_block"]

		// We can only apply this validation if we know the final value of the field, and we may
		// not know the final value if users feed the value into their config in unintuitive ways.
		// https://github.com/terraform-providers/terraform-provider-google/issues/4186
		blockValueKnown := d.NewValueKnown("private_cluster_config.0.master_ipv4_cidr_block")

		if blockValueKnown && (block == nil || block == "") {
			return fmt.Errorf("master_ipv4_cidr_block must be set if enable_private_nodes == true")
		}
	}
	return nil
}
