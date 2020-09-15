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
				Optional:    true,
				Elem:        cidrBlockConfig,
				Description: `External networks that can access the Kubernetes cluster master through HTTPS.`,
			},
		},
	}
	cidrBlockConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.CIDRNetwork(0, 32),
				Description:  `External network that can access Kubernetes master through HTTPS. Must be specified in CIDR notation.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Field for users to identify CIDR blocks.`,
			},
		},
	}

	ipAllocationCidrBlockFields = []string{"ip_allocation_policy.0.cluster_ipv4_cidr_block", "ip_allocation_policy.0.services_ipv4_cidr_block"}
	ipAllocationRangeFields     = []string{"ip_allocation_policy.0.cluster_secondary_range_name", "ip_allocation_policy.0.services_secondary_range_name"}

	addonsConfigKeys = []string{
		"addons_config.0.http_load_balancing",
		"addons_config.0.horizontal_pod_autoscaling",
		"addons_config.0.network_policy_config",
		"addons_config.0.cloudrun_config",
	}

	forceNewClusterNodeConfigFields = []string{}
)

// This uses the node pool nodeConfig schema but sets
// node-pool-only updatable fields to ForceNew
func clusterSchemaNodeConfig() *schema.Schema {
	nodeConfigSch := schemaNodeConfig()
	schemaMap := nodeConfigSch.Elem.(*schema.Resource).Schema
	for _, k := range forceNewClusterNodeConfigFields {
		if sch, ok := schemaMap[k]; ok {
			changeFieldSchemaToForceNew(sch)
		}
	}
	return nodeConfigSch
}

func rfc5545RecurrenceDiffSuppress(k, o, n string, d *schema.ResourceData) bool {
	// This diff gets applied in the cloud console if you specify
	// "FREQ=DAILY" in your config and add a maintenance exclusion.
	if o == "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU" && n == "FREQ=DAILY" {
		return true
	}
	// Writing a full diff suppress for identical recurrences would be
	// complex and error-prone - it's not a big problem if a user
	// changes the recurrence and it's textually difference but semantically
	// identical.
	return false
}

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
			Read:   schema.DefaultTimeout(40 * time.Minute),
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the cluster, unique within the project and location.`,
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

			"operation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The location (region or zone) in which the cluster master will be created, as well as the default node location. If you specify a zone (such as us-central1-a), the cluster will be a zonal cluster with a single cluster master. If you specify a region (such as us-west1), the cluster will be a regional cluster with multiple masters spread across zones in the region, and with default node locations in those zones as well.`,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Removed:     "Use location instead",
				Computed:    true,
				Description: `The region in which the cluster master will be created. Zone and region have been removed in favor of location.`,
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Removed:     "Use location instead",
				Computed:    true,
				Description: `The zone in which the cluster master will be created. Zone and region have been removed in favor of location.`,
			},

			"node_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of zones in which the cluster's nodes are located. Nodes must be in the region of their regional cluster or in the same region as their cluster's zone for zonal clusters. If this is specified for a zonal cluster, omit the cluster's zone.`,
			},

			"additional_zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				Removed:     "Use node_locations instead",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Additional_zones has been removed in favor of node_locations.`,
			},

			"addons_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The configuration for addons supported by GKE.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_load_balancing": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the HTTP (L7) load balancing controller addon, which makes it easy to set up HTTP load balancers for services in a cluster. It is enabled by default; set disabled = true to disable.`,
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
							Description:  `The status of the Horizontal Pod Autoscaling addon, which increases or decreases the number of replica pods a replication controller has based on the resource usage of the existing pods. It ensures that a Heapster pod is running in the cluster, which is also used by the Cloud Monitoring service. It is enabled by default; set disabled = true to disable.`,
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
							Type:        schema.TypeList,
							Optional:    true,
							Removed:     "The Kubernetes Dashboard addon is removed for clusters on GKE.",
							Computed:    true,
							MaxItems:    1,
							Description: `The status of  Kubernetes Dashboard addon.`,
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
							Description:  `Whether we should enable the network policy addon for the master. This must be enabled in order to enable network policy for the nodes. To enable this, you must also define a network_policy block, otherwise nothing will happen. It can only be disabled if the nodes already do not have network policies enabled. Defaults to disabled; set disabled = false to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"cloudrun_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the CloudRun addon. It is disabled by default. Set disabled = false to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"load_balancer_type": {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"LOAD_BALANCER_TYPE_INTERNAL"}, false),
										Optional:     true,
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
				Optional:    true,
				Computed:    true,
				Description: `Per-cluster configuration of Node Auto-Provisioning with Cluster Autoscaler to automatically adjust the size of the cluster and create/delete node pools based on the current needs of the cluster's workload. See the guide to using Node Auto-Provisioning for more details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether node auto-provisioning is enabled. Resource limits for cpu and memory must be defined to enable node auto-provisioning.`,
						},
						"resource_limits": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Global constraints for machine resources in the cluster. Configuring the cpu and memory types is required if node auto-provisioning is enabled. These limits will apply to node pool autoscaling in addition to node auto-provisioning.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The type of the resource. For example, cpu and memory. See the guide to using Node Auto-Provisioning for a list of types.`,
									},
									"minimum": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Minimum amount of the resource in the cluster.`,
									},
									"maximum": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Maximum amount of the resource in the cluster.`,
									},
								},
							},
						},
						"auto_provisioning_defaults": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `Contains defaults for a node pool created by NAP.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oauth_scopes": {
										Type:             schema.TypeList,
										Optional:         true,
										Computed:         true,
										Elem:             &schema.Schema{Type: schema.TypeString},
										DiffSuppressFunc: containerClusterAddedScopesSuppress,
										Description:      `Scopes that are used by NAP when creating node pools.`,
									},
									"service_account": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "default",
										Description: `The Google Cloud Platform Service Account to be used by the node VMs.`,
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
				Description:   `The IP address range of the Kubernetes pods in this cluster in CIDR notation (e.g. 10.96.0.0/14). Leave blank to have one automatically chosen or specify a /14 block in 10.0.0.0/8. This field will only work for routes-based clusters, where ip_allocation_policy is not defined.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: ` Description of the cluster.`,
			},

			"enable_binary_authorization": {
				Default:     false,
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enable Binary Authorization for this cluster. If enabled, all container images will be validated by Google Binary Authorization.`,
			},

			"enable_kubernetes_alpha": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: `Whether to enable Kubernetes Alpha features for this cluster. Note that when this option is enabled, the cluster cannot be upgraded and will be automatically deleted after 30 days.`,
			},

			"enable_tpu": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to enable Cloud TPU resources in this cluster.`,
				Removed:     "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
				Computed:    true,
			},

			"enable_legacy_abac": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether the ABAC authorizer is enabled for this cluster. When enabled, identities in the system, including service accounts, nodes, and controllers, will have statically granted permissions beyond those provided by the RBAC configuration or IAM. Defaults to false.`,
			},

			"enable_shielded_nodes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Enable Shielded Nodes features on all nodes in this cluster. Defaults to false.`,
			},

			"authenticator_groups_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `Configuration for the Google Groups for GKE feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The name of the RBAC security group for use with Google security groups in Kubernetes RBAC. Group name must be in format gke-security-groups@yourdomain.com.`,
						},
					},
				},
			},

			"initial_node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `The number of nodes to create in this cluster's default node pool. In regional or multi-zonal clusters, this is the number of nodes per zone. Must be set if node_pool is not set. If you're using google_container_node_pool objects with no default node pool, you'll need to set this to a value of at least 1, alongside setting remove_default_node_pool to true.`,
			},

			"logging_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"logging.googleapis.com", "logging.googleapis.com/kubernetes", "none"}, false),
				Description:  `The logging service that the cluster should write logs to. Available options include logging.googleapis.com(Legacy Stackdriver), logging.googleapis.com/kubernetes(Stackdriver Kubernetes Engine Logging), and none. Defaults to logging.googleapis.com/kubernetes.`,
			},

			"maintenance_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The maintenance policy to use for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_maintenance_window": {
							Type:     schema.TypeList,
							Optional: true,
							ExactlyOneOf: []string{
								"maintenance_policy.0.daily_maintenance_window",
								"maintenance_policy.0.recurring_window",
							},
							MaxItems:    1,
							Description: `Time window specified for daily maintenance operations. Specify start_time in RFC3339 format "HH:MM”, where HH : [00-23] and MM : [00-59] GMT.`,
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
						"recurring_window": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ExactlyOneOf: []string{
								"maintenance_policy.0.daily_maintenance_window",
								"maintenance_policy.0.recurring_window",
							},
							Description: `Time window for recurring maintenance operations.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateRFC3339Date,
									},
									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateRFC3339Date,
									},
									"recurrence": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: rfc5545RecurrenceDiffSuppress,
									},
								},
							},
						},
					},
				},
			},

			"master_auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `The authentication information for accessing the Kubernetes master. Some values in this block are only returned by the API if your service account has permission to get credentials for your GKE cluster. If you see an unexpected diff removing a username/password or unsetting your client cert, ensure you have the container.clusters.getCredentials permission.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"master_auth.0.password", "master_auth.0.username", "master_auth.0.client_certificate_config"},
							Sensitive:    true,
							Description:  `The password to use for HTTP basic authentication when accessing the Kubernetes master endpoint.`,
						},

						"username": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"master_auth.0.password", "master_auth.0.username", "master_auth.0.client_certificate_config"},
							Description:  `The username to use for HTTP basic authentication when accessing the Kubernetes master endpoint. If not present basic auth will be disabled.`,
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
							Description:  `Whether client certificate authorization is enabled for this cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issue_client_certificate": {
										Type:        schema.TypeBool,
										Required:    true,
										ForceNew:    true,
										Description: `Whether client certificate authorization is enabled for this cluster.`,
									},
								},
							},
						},

						"client_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Base64 encoded public certificate used by clients to authenticate to the cluster endpoint.`,
						},

						"client_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `Base64 encoded private key used by clients to authenticate to the cluster endpoint.`,
						},

						"cluster_ca_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Base64 encoded public certificate that is the root of trust for the cluster.`,
						},
					},
				},
			},

			"master_authorized_networks_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        networkConfig,
				Description: `The desired configuration options for master authorized networks. Omit the nested cidr_blocks attribute to disallow external access (except the cluster node IPs, which GKE automatically whitelists).`,
			},

			"min_master_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The minimum version of the master. GKE will auto-update the master to new versions, so this does not guarantee the current master version--use the read-only master_version field to obtain that. If unset, the cluster's version will be set by GKE to the version of the most recent official release (which is not necessarily the latest version).`,
			},

			"monitoring_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"monitoring.googleapis.com", "monitoring.googleapis.com/kubernetes", "none"}, false),
				Description:  `The monitoring service that the cluster should write metrics to. Automatically send metrics from pods in the cluster to the Google Cloud Monitoring API. VM metrics will be collected by Google Compute Engine regardless of this setting Available options include monitoring.googleapis.com(Legacy Stackdriver), monitoring.googleapis.com/kubernetes(Stackdriver Kubernetes Engine Monitoring), and none. Defaults to monitoring.googleapis.com/kubernetes.`,
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "default",
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name or self_link of the Google Compute Engine network to which the cluster is connected. For Shared VPC, set this to the self link of the shared network.`,
			},

			"network_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration options for the NetworkPolicy feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether network policy is enabled on the cluster.`,
						},
						"provider": {
							Type:             schema.TypeString,
							Default:          "PROVIDER_UNSPECIFIED",
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROVIDER_UNSPECIFIED", "CALICO"}, false),
							DiffSuppressFunc: emptyOrDefaultStringSuppress("PROVIDER_UNSPECIFIED"),
							Description:      `The selected network policy provider. Defaults to PROVIDER_UNSPECIFIED.`,
						},
					},
				},
			},

			"node_config": clusterSchemaNodeConfig(),

			"node_pool": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true, // TODO(danawillow): Add ability to add/remove nodePools
				Elem: &schema.Resource{
					Schema: schemaNodePool,
				},
				Description: `List of node pools associated with this cluster. See google_container_node_pool for schema. Warning: node pools defined inside a cluster can't be changed (or added/removed) after cluster creation without deleting and recreating the entire cluster. Unless you absolutely need the ability to say "these are the only node pools associated with this cluster", use the google_container_node_pool resource instead of this property.`,
			},

			"node_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The Kubernetes version on the nodes. Must either be unset or set to the same value as min_master_version on create. Defaults to the default version set by GKE which is not necessarily the latest version. This only affects nodes in the default node pool. While a fuzzy version can be specified, it's recommended that you specify explicit versions as Terraform will see spurious diffs when fuzzy versions are used. See the google_container_engine_versions data source's version_prefix field to approximate fuzzy versions in a Terraform-compatible way. To update nodes in other node pools, use the version attribute on the node pool.`,
			},

			"pod_security_policy_config": {
				// Remove return nil from expand when this is removed for good.
				Removed:     "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
				Computed:    true,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Configuration for the PodSecurityPolicy feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enable the PodSecurityPolicy controller for this cluster. If enabled, pods must be valid under a PodSecurityPolicy to be created.`,
						},
					},
				},
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name or self_link of the Google Compute Engine subnetwork in which the cluster's instances are launched.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Server-defined URL for the resource.`,
			},

			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of this cluster's Kubernetes master.`,
			},

			"instance_group_urls": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `List of instance group URLs which have been assigned to the cluster.`,
			},

			"master_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current version of the master in the cluster. This may be different than the min_master_version set in the config if the master has been updated by GKE.`,
			},

			"services_ipv4_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address range of the Kubernetes services in this cluster, in CIDR notation (e.g. 1.2.3.4/29). Service addresses are typically put in the last /16 from the container CIDR.`,
			},

			"ip_allocation_policy": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"cluster_ipv4_cidr"},
				Description:   `Configuration of cluster IP allocation for VPC-native clusters. Adding this block enables IP aliasing, making the cluster VPC-native instead of routes-based.`,
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
							Description:      `The IP address range for the cluster pod IPs. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.`,
						},

						"services_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: cidrOrSizeDiffSuppress,
							Description:      `The IP address range of the services IPs in this cluster. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.`,
						},

						// User manages secondary ranges manually
						"cluster_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
							Description:   `The name of the existing secondary range in the cluster's subnetwork to use for pod IP addresses. Alternatively, cluster_ipv4_cidr_block can be used to automatically create a GKE-managed one.`,
						},

						"services_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
							Description:   `The name of the existing secondary range in the cluster's subnetwork to use for service ClusterIPs. Alternatively, services_ipv4_cidr_block can be used to automatically create a GKE-managed one.`,
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
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `If true, deletes the default node pool upon cluster creation. If you're using google_container_node_pool resources with no default node pool, this should be set to true, alongside setting initial_node_count to at least 1.`,
			},

			"private_cluster_config": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
				Description:      `Configuration for private clusters, clusters with private nodes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_private_endpoint": {
							Type:             schema.TypeBool,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `Enables the private cluster feature, creating a private endpoint on the cluster. In a private cluster, nodes only have RFC 1918 private addresses and communicate with the master's private endpoint via private networking.`,
						},
						"enable_private_nodes": {
							Type:             schema.TypeBool,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `When true, the cluster's private endpoint is used as the cluster endpoint and access through the public endpoint is disabled. When false, either endpoint can be used. This field only applies to private clusters, when enable_private_nodes is true.`,
						},
						"master_ipv4_cidr_block": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: orEmpty(validation.CIDRNetwork(28, 28)),
							Description:  `The IP range in CIDR notation to use for the hosted master network. This range will be used for assigning private IP addresses to the cluster master(s) and the ILB VIP. This range must not overlap with any other ranges in use within the cluster's network, and it must be a /28 subnet. See Private Cluster Limitations for more details. This field only applies to private clusters, when enable_private_nodes is true.`,
						},
						"peering_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the peering between this cluster and the Google owned VPC.`,
						},
						"private_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The internal IP address of this cluster's master endpoint.`,
						},
						"public_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The external IP address of this cluster's master endpoint.`,
						},
					},
				},
			},

			"resource_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The GCE resource labels (a map of key/value pairs) to be applied to the cluster.`,
			},

			"label_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the set of labels for this cluster.`,
			},

			"default_max_pods_per_node": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The default maximum number of pods per node in this cluster. This doesn't work on "routes-based" clusters, clusters that don't have IP Aliasing enabled.`,
			},

			"vertical_pod_autoscaling": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enables vertical pod autoscaling.`,
						},
					},
				},
			},
			"workload_identity_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"database_encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Application-layer Secrets Encryption settings. The object format is {state = string, key_name = string}. Valid values of state are: "ENCRYPTED"; "DECRYPTED". key_name is the name of a CloudKMS key.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ENCRYPTED", "DECRYPTED"}, false),
							Description:  `ENCRYPTED or DECRYPTED.`,
						},
						"key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The key to use to encrypt/decrypt secrets.`,
						},
					},
				},
			},

			"release_channel": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration options for the Release channel feature, which provide more control over automatic upgrades of your GKE clusters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringInSlice([]string{"UNSPECIFIED", "RAPID", "REGULAR", "STABLE"}, false),
							DiffSuppressFunc: emptyOrDefaultStringSuppress("UNSPECIFIED"),
							Description:      `The selected release channel.`,
						},
					},
				},
			},

			"enable_intranode_visibility": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether Intra-node visibility is enabled for this cluster. This makes same node pod to pod traffic visible for VPC network.`,
				Removed:     "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/guides/provider_versions.html for more details.",
				Computed:    true,
			},

			"resource_usage_export_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `Configuration for the ResourceUsageExportConfig feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_network_egress_metering": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: `Whether to enable network egress metering for this cluster. If enabled, a daemonset will be created in the cluster to meter network egress traffic.`,
						},
						"enable_resource_consumption_metering": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Whether to enable resource consumption metering on this cluster. When enabled, a table will be created in the resource export BigQuery dataset to store resource consumption data. The resulting table can be joined with the resource usage table or with BigQuery billing export. Defaults to true.`,
						},
						"bigquery_destination": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: `Parameters for using BigQuery as the destination of resource usage export.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The ID of a BigQuery Dataset.`,
									},
								},
							},
						},
					},
				},
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
// See https://github.com/hashicorp/terraform-provider-google/issues/3786
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

	ipAllocationBlock, err := expandIPAllocationPolicy(d.Get("ip_allocation_policy"))
	if err != nil {
		return err
	}

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
		IpAllocationPolicy:      ipAllocationBlock,
		PodSecurityPolicyConfig: expandPodSecurityPolicyConfig(d.Get("pod_security_policy_config")),
		Autoscaling:             expandClusterAutoscaling(d.Get("cluster_autoscaling"), d),
		BinaryAuthorization: &containerBeta.BinaryAuthorization{
			Enabled:         d.Get("enable_binary_authorization").(bool),
			ForceSendFields: []string{"Enabled"},
		},
		ShieldedNodes: &containerBeta.ShieldedNodes{
			Enabled:         d.Get("enable_shielded_nodes").(bool),
			ForceSendFields: []string{"Enabled"},
		},
		ReleaseChannel: expandReleaseChannel(d.Get("release_channel")),
		MasterAuth:     expandMasterAuth(d.Get("master_auth")),
		ResourceLabels: expandStringMap(d, "resource_labels"),
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
		subnetwork, err := parseRegionalFieldValue("subnetworks", v.(string), "project", "location", "location", d, config, true) // variant of ParseSubnetworkFieldValue
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

	if v, ok := d.GetOk("database_encryption"); ok {
		cluster.DatabaseEncryption = expandDatabaseEncryption(v)
	}

	if v, ok := d.GetOk("workload_identity_config"); ok {
		cluster.WorkloadIdentityConfig = expandWorkloadIdentityConfig(v)
	}

	if v, ok := d.GetOk("resource_usage_export_config"); ok {
		cluster.ResourceUsageExportConfig = expandResourceUsageExportConfig(v)
	}

	req := &containerBeta.CreateClusterRequest{
		Cluster: cluster,
	}

	mutexKV.Lock(containerClusterMutexKey(project, location, clusterName))
	defer mutexKV.Unlock(containerClusterMutexKey(project, location, clusterName))

	parent := fmt.Sprintf("projects/%s/locations/%s", project, location)
	var op *containerBeta.Operation
	err = retry(func() error {
		clusterCreateCall := config.clientContainerBeta.Projects.Locations.Clusters.Create(parent, req)
		if config.UserProjectOverride {
			clusterCreateCall.Header().Add("X-Goog-User-Project", project)
		}
		op, err = clusterCreateCall.Do()
		return err
	})
	if err != nil {
		return err
	}

	d.SetId(containerClusterFullName(project, location, clusterName))

	// Wait until it's created
	waitErr := containerOperationWait(config, op, project, location, "creating GKE cluster", d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		// Check if the create operation failed because Terraform was prematurely terminated. If it was we can persist the
		// operation id to state so that a subsequent refresh of this resource will wait until the operation has terminated
		// before attempting to Read the state of the cluster. This allows a graceful resumption of a Create that was killed
		// by the upstream Terraform process exiting early such as a sigterm.
		select {
		case <-config.context.Done():
			log.Printf("[DEBUG] Persisting %s so this operation can be resumed \n", op.Name)
			d.Set("operation", op.Name)
			return nil
		default:
			// leaving default case to ensure this is non blocking
		}
		// Try a GET on the cluster so we can see the state in debug logs. This will help classify error states.
		clusterGetCall := config.clientContainerBeta.Projects.Locations.Clusters.Get(containerClusterFullName(project, location, clusterName))
		if config.UserProjectOverride {
			clusterGetCall.Header().Add("X-Goog-User-Project", project)
		}
		_, getErr := clusterGetCall.Do()
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
			clusterNodePoolDeleteCall := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(parent)
			if config.UserProjectOverride {
				clusterNodePoolDeleteCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err = clusterNodePoolDeleteCall.Do()
			return err
		})
		if err != nil {
			return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
		}
		err = containerOperationWait(config, op, project, location, "removing default node pool", d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return errwrap.Wrapf("Error while waiting to delete default node pool: {{err}}", err)
		}
	}

	if err := resourceContainerClusterRead(d, meta); err != nil {
		return err
	}

	state, err := containerClusterAwaitRestingState(config, project, location, clusterName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	if containerClusterRestingStates[state] == ErrorState {
		return fmt.Errorf("Cluster %s was created in the error state %q", clusterName, state)
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

	operation := d.Get("operation").(string)
	if operation != "" {
		log.Printf("[DEBUG] in progress operation detected at %v, attempting to resume", operation)
		op := &containerBeta.Operation{
			Name: operation,
		}
		d.Set("operation", "")
		waitErr := containerOperationWait(config, op, project, location, "resuming GKE cluster", d.Timeout(schema.TimeoutRead))
		if waitErr != nil {
			return waitErr
		}
	}

	clusterName := d.Get("name").(string)
	name := containerClusterFullName(project, location, clusterName)
	clusterGetCall := config.clientContainerBeta.Projects.Locations.Clusters.Get(name)
	if config.UserProjectOverride {
		clusterGetCall.Header().Add("X-Goog-User-Project", project)
	}

	cluster, err := clusterGetCall.Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Container Cluster %q", d.Get("name").(string)))
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
	d.Set("self_link", cluster.SelfLink)
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
	d.Set("enable_binary_authorization", cluster.BinaryAuthorization != nil && cluster.BinaryAuthorization.Enabled)
	if cluster.ShieldedNodes != nil {
		d.Set("enable_shielded_nodes", cluster.ShieldedNodes.Enabled)
	}
	if err := d.Set("release_channel", flattenReleaseChannel(cluster.ReleaseChannel)); err != nil {
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

	if err := d.Set("workload_identity_config", flattenWorkloadIdentityConfig(cluster.WorkloadIdentityConfig)); err != nil {
		return err
	}

	if err := d.Set("database_encryption", flattenDatabaseEncryption(cluster.DatabaseEncryption)); err != nil {
		return err
	}

	d.Set("resource_labels", cluster.ResourceLabels)
	d.Set("label_fingerprint", cluster.LabelFingerprint)

	if err := d.Set("resource_usage_export_config", flattenResourceUsageExportConfig(cluster.ResourceUsageExportConfig)); err != nil {
		return err
	}

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

	if _, err := containerClusterAwaitRestingState(config, project, location, clusterName, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	d.Partial(true)

	lockKey := containerClusterMutexKey(project, location, clusterName)

	updateFunc := func(req *containerBeta.UpdateClusterRequest, updateDescription string) func() error {
		return func() error {
			name := containerClusterFullName(project, location, clusterName)
			clusterUpdateCall := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req)
			if config.UserProjectOverride {
				clusterUpdateCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterUpdateCall.Do()
			if err != nil {
				return err
			}
			// Wait until it's updated
			return containerOperationWait(config, op, project, location, updateDescription, d.Timeout(schema.TimeoutUpdate))
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

	if d.HasChange("enable_binary_authorization") {
		enabled := d.Get("enable_binary_authorization").(bool)
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredBinaryAuthorization: &containerBeta.BinaryAuthorization{
					Enabled:         enabled,
					ForceSendFields: []string{"Enabled"},
				},
			},
		}

		updateF := updateFunc(req, "updating GKE binary authorization")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s's binary authorization has been updated to %v", d.Id(), enabled)

		d.SetPartial("enable_binary_authorization")
	}

	if d.HasChange("enable_shielded_nodes") {
		enabled := d.Get("enable_shielded_nodes").(bool)
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredShieldedNodes: &containerBeta.ShieldedNodes{
					Enabled:         enabled,
					ForceSendFields: []string{"Enabled"},
				},
			},
		}

		updateF := updateFunc(req, "updating GKE shielded nodes")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s's shielded nodes has been updated to %v", d.Id(), enabled)

		d.SetPartial("enable_shielded_nodes")
	}

	if d.HasChange("release_channel") {
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredReleaseChannel: expandReleaseChannel(d.Get("release_channel")),
			},
		}
		updateF := func() error {
			log.Println("[DEBUG] updating release_channel")
			name := containerClusterFullName(project, location, clusterName)
			clusterUpdateCall := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req)
			if config.UserProjectOverride {
				clusterUpdateCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterUpdateCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			err = containerOperationWait(config, op, project, location, "updating Release Channel", d.Timeout(schema.TimeoutUpdate))
			log.Println("[DEBUG] done updating release_channel")
			return err
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s Release Channel has been updated to %#v", d.Id(), req.Update.DesiredReleaseChannel)

		d.SetPartial("release_channel")
	}

	if d.HasChange("maintenance_policy") {
		req := &containerBeta.SetMaintenancePolicyRequest{
			MaintenancePolicy: expandMaintenancePolicy(d, meta),
		}

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			clusterSetMaintenancePolicyCall := config.clientContainerBeta.Projects.Locations.Clusters.SetMaintenancePolicy(name, req)
			if config.UserProjectOverride {
				clusterSetMaintenancePolicyCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterSetMaintenancePolicyCall.Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE cluster maintenance policy", d.Timeout(schema.TimeoutUpdate))
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
			clusterSetLegacyAbacCall := config.clientContainerBeta.Projects.Locations.Clusters.SetLegacyAbac(name, req)
			if config.UserProjectOverride {
				clusterSetLegacyAbacCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterSetLegacyAbacCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			err = containerOperationWait(config, op, project, location, "updating GKE legacy ABAC", d.Timeout(schema.TimeoutUpdate))
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
			clusterUpdateCall := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req)
			if config.UserProjectOverride {
				clusterUpdateCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterUpdateCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE logging+monitoring service", d.Timeout(schema.TimeoutUpdate))
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
			clusterSetNetworkPolicyCall := config.clientContainerBeta.Projects.Locations.Clusters.SetNetworkPolicy(name, req)
			if config.UserProjectOverride {
				clusterSetNetworkPolicyCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterSetNetworkPolicyCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			err = containerOperationWait(config, op, project, location, "updating GKE cluster network policy", d.Timeout(schema.TimeoutUpdate))
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

			if err := nodePoolUpdate(d, meta, nodePoolInfo, fmt.Sprintf("node_pool.%d.", i), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
		d.SetPartial("node_pool")
	}

	// The master must be updated before the nodes
	// If set to "", skip this step- any master version satisfies that minimum.
	if ver := d.Get("min_master_version").(string); d.HasChange("min_master_version") && ver != "" {
		des, err := version.NewVersion(ver)
		if err != nil {
			return err
		}

		currentMasterVersion := d.Get("master_version").(string)
		cur, err := version.NewVersion(currentMasterVersion)
		if err != nil {
			return err
		}

		// Only upgrade the master if the current version is lower than the desired version
		if cur.LessThan(des) {
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredMasterVersion: ver,
				},
			}

			updateF := updateFunc(req, "updating GKE master version")
			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}
			log.Printf("[INFO] GKE cluster %s: master has been updated to %s", d.Id(), ver)
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
				clusterUpdateCall := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req)
				if config.UserProjectOverride {
					clusterUpdateCall.Header().Add("X-Goog-User-Project", project)
				}
				op, err := clusterUpdateCall.Do()
				if err != nil {
					return err
				}

				// Wait until it's updated
				return containerOperationWait(config, op, project, location, "updating GKE image type", d.Timeout(schema.TimeoutUpdate))
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
			clusterSetMasterAuthCall := config.clientContainerBeta.Projects.Locations.Clusters.SetMasterAuth(name, req)
			if config.UserProjectOverride {
				clusterSetMasterAuthCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterSetMasterAuthCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating master auth", d.Timeout(schema.TimeoutUpdate))
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

	if d.HasChange("database_encryption") {
		c := d.Get("database_encryption")
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredDatabaseEncryption: expandDatabaseEncryption(c),
			},
		}

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			clusterUpdateCall := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req)
			if config.UserProjectOverride {
				clusterUpdateCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterUpdateCall.Do()
			if err != nil {
				return err
			}
			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE cluster database encryption config", d.Timeout(schema.TimeoutUpdate))
		}
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}
		log.Printf("[INFO] GKE cluster %s database encryption config has been updated", d.Id())

		d.SetPartial("database_encryption")
	}

	if d.HasChange("workload_identity_config") {
		// Because GKE uses a non-RESTful update function, when removing the
		// feature you need to specify a fairly full request body or it fails:
		// "update": {"desiredWorkloadIdentityConfig": {"identityNamespace": ""}}
		req := &containerBeta.UpdateClusterRequest{}
		if v, ok := d.GetOk("workload_identity_config"); !ok {
			req.Update = &containerBeta.ClusterUpdate{
				DesiredWorkloadIdentityConfig: &containerBeta.WorkloadIdentityConfig{
					IdentityNamespace: "",
					ForceSendFields:   []string{"IdentityNamespace"},
				},
			}
		} else {
			req.Update = &containerBeta.ClusterUpdate{
				DesiredWorkloadIdentityConfig: expandWorkloadIdentityConfig(v),
			}
		}

		updateF := updateFunc(req, "updating GKE cluster workload identity config")
		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE cluster %s workload identity config has been updated", d.Id())

		d.SetPartial("workload_identity_config")
	}

	if d.HasChange("resource_labels") {
		resourceLabels := d.Get("resource_labels").(map[string]interface{})
		labelFingerprint := d.Get("label_fingerprint").(string)
		req := &containerBeta.SetLabelsRequest{
			ResourceLabels:   convertStringMap(resourceLabels),
			LabelFingerprint: labelFingerprint,
		}
		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			clusterSetResourceLabelsCall := config.clientContainerBeta.Projects.Locations.Clusters.SetResourceLabels(name, req)
			if config.UserProjectOverride {
				clusterSetResourceLabelsCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterSetResourceLabelsCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE resource labels", d.Timeout(schema.TimeoutUpdate))
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		d.SetPartial("resource_labels")
	}

	if d.HasChange("remove_default_node_pool") && d.Get("remove_default_node_pool").(bool) {
		name := fmt.Sprintf("%s/nodePools/%s", containerClusterFullName(project, location, clusterName), "default-pool")
		clusterNodePoolDeleteCall := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(name)
		if config.UserProjectOverride {
			clusterNodePoolDeleteCall.Header().Add("X-Goog-User-Project", project)
		}
		op, err := clusterNodePoolDeleteCall.Do()
		if err != nil {
			if !isGoogleApiErrorWithCode(err, 404) {
				return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
			}
			log.Printf("[WARN] Container cluster %q default node pool already removed, no change", d.Id())
		} else {
			err = containerOperationWait(config, op, project, location, "removing default node pool", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return errwrap.Wrapf("Error deleting default node pool: {{err}}", err)
			}
		}
	}

	if d.HasChange("resource_usage_export_config") {
		c := d.Get("resource_usage_export_config")
		req := &containerBeta.UpdateClusterRequest{
			Update: &containerBeta.ClusterUpdate{
				DesiredResourceUsageExportConfig: expandResourceUsageExportConfig(c),
			},
		}

		updateF := func() error {
			name := containerClusterFullName(project, location, clusterName)
			clusterUpdateCall := config.clientContainerBeta.Projects.Locations.Clusters.Update(name, req)
			if config.UserProjectOverride {
				clusterUpdateCall.Header().Add("X-Goog-User-Project", project)
			}
			op, err := clusterUpdateCall.Do()
			if err != nil {
				return err
			}
			// Wait until it's updated
			return containerOperationWait(config, op, project, location, "updating GKE cluster resource usage export config", d.Timeout(schema.TimeoutUpdate))
		}
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}
		log.Printf("[INFO] GKE cluster %s resource usage export config has been updated", d.Id())

		d.SetPartial("resource_usage_export_config")
	}

	d.Partial(false)

	if _, err := containerClusterAwaitRestingState(config, project, location, clusterName, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

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

	if _, err := containerClusterAwaitRestingState(config, project, location, clusterName, d.Timeout(schema.TimeoutDelete)); err != nil {
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
		clusterDeleteCall := config.clientContainerBeta.Projects.Locations.Clusters.Delete(name)
		if config.UserProjectOverride {
			clusterDeleteCall.Header().Add("X-Goog-User-Project", project)
		}
		op, err = clusterDeleteCall.Do()

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
	waitErr := containerOperationWait(config, op, project, location, "deleting GKE cluster", d.Timeout(schema.TimeoutDelete))
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
	clusterDeleteCall := config.clientContainerBeta.Projects.Locations.Clusters.Delete(fullName)
	if config.UserProjectOverride {
		clusterDeleteCall.Header().Add("X-Goog-User-Project", project)
	}
	op, err := clusterDeleteCall.Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Container Cluster %q", d.Get("name").(string)))
	}

	// Wait until it's deleted
	waitErr := containerOperationWait(config, op, project, location, "deleting GKE cluster", d.Timeout(schema.TimeoutDelete))
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] GKE cluster %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

var containerClusterRestingStates = RestingStates{
	"RUNNING":  ReadyState,
	"DEGRADED": ErrorState,
	"ERROR":    ErrorState,
}

// returns a state with no error if the state is a resting state, and the last state with an error otherwise
func containerClusterAwaitRestingState(config *Config, project, location, clusterName string, timeout time.Duration) (state string, err error) {
	err = resource.Retry(timeout, func() *resource.RetryError {
		name := containerClusterFullName(project, location, clusterName)
		clusterGetCall := config.clientContainerBeta.Projects.Locations.Clusters.Get(name)
		if config.UserProjectOverride {
			clusterGetCall.Header().Add("X-Goog-User-Project", project)
		}
		cluster, gErr := clusterGetCall.Do()
		if gErr != nil {
			return resource.NonRetryableError(gErr)
		}

		state = cluster.Status

		switch stateType := containerClusterRestingStates[cluster.Status]; stateType {
		case ReadyState:
			log.Printf("[DEBUG] Cluster %q has status %q with message %q.", clusterName, state, cluster.StatusMessage)
			return nil
		case ErrorState:
			log.Printf("[DEBUG] Cluster %q has error state %q with message %q.", clusterName, state, cluster.StatusMessage)
			return nil
		default:
			return resource.RetryableError(fmt.Errorf("Cluster %q has state %q with message %q", clusterName, state, cluster.StatusMessage))
		}
	})

	return state, err
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
		if isGoogleApiErrorWithCode(err, 404) {
			// The IGM URL is stale; don't include it
			continue
		}
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

	if v, ok := config["cloudrun_config"]; ok && len(v.([]interface{})) > 0 {
		addon := v.([]interface{})[0].(map[string]interface{})
		ac.CloudRunConfig = &containerBeta.CloudRunConfig{
			Disabled:        addon["disabled"].(bool),
			ForceSendFields: []string{"Disabled"},
		}
		if addon["load_balancer_type"] != "" {
			ac.CloudRunConfig.LoadBalancerType = addon["load_balancer_type"].(string)
		}
	}

	return ac
}

func expandIPAllocationPolicy(configured interface{}) (*containerBeta.IPAllocationPolicy, error) {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &containerBeta.IPAllocationPolicy{
			UseIpAliases:    false,
			ForceSendFields: []string{"UseIpAliases"},
		}, nil
	}

	config := l[0].(map[string]interface{})
	return &containerBeta.IPAllocationPolicy{
		UseIpAliases:          true,
		ClusterIpv4CidrBlock:  config["cluster_ipv4_cidr_block"].(string),
		ServicesIpv4CidrBlock: config["services_ipv4_cidr_block"].(string),

		ClusterSecondaryRangeName:  config["cluster_secondary_range_name"].(string),
		ServicesSecondaryRangeName: config["services_secondary_range_name"].(string),
		ForceSendFields:            []string{"UseIpAliases"},
	}, nil
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
	clusterGetCall := config.clientContainerBeta.Projects.Locations.Clusters.Get(name)
	if config.UserProjectOverride {
		clusterGetCall.Header().Add("X-Goog-User-Project", project)
	}
	cluster, _ := clusterGetCall.Do()
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
	if recurringWindow, ok := maintenancePolicy["recurring_window"]; ok && len(recurringWindow.([]interface{})) > 0 {
		rw := recurringWindow.([]interface{})[0].(map[string]interface{})
		return &containerBeta.MaintenancePolicy{
			Window: &containerBeta.MaintenanceWindow{
				MaintenanceExclusions: exclusions,
				RecurringWindow: &containerBeta.RecurringTimeWindow{
					Window: &containerBeta.TimeWindow{
						StartTime: rw["start_time"].(string),
						EndTime:   rw["end_time"].(string),
					},
					Recurrence: rw["recurrence"].(string),
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

	npd := &containerBeta.AutoprovisioningNodePoolDefaults{
		OauthScopes:    convertStringArr(config["oauth_scopes"].([]interface{})),
		ServiceAccount: config["service_account"].(string),
	}

	return npd
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

func expandDatabaseEncryption(configured interface{}) *containerBeta.DatabaseEncryption {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &containerBeta.DatabaseEncryption{
		State:   config["state"].(string),
		KeyName: config["key_name"].(string),
	}
}

func expandReleaseChannel(configured interface{}) *containerBeta.ReleaseChannel {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &containerBeta.ReleaseChannel{
		Channel: config["channel"].(string),
	}
}

func expandWorkloadIdentityConfig(configured interface{}) *containerBeta.WorkloadIdentityConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &containerBeta.WorkloadIdentityConfig{
		IdentityNamespace: config["identity_namespace"].(string),
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
func expandResourceUsageExportConfig(configured interface{}) *containerBeta.ResourceUsageExportConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return &containerBeta.ResourceUsageExportConfig{}
	}

	resourceUsageConfig := l[0].(map[string]interface{})

	result := &containerBeta.ResourceUsageExportConfig{
		EnableNetworkEgressMetering: resourceUsageConfig["enable_network_egress_metering"].(bool),
		ConsumptionMeteringConfig: &containerBeta.ConsumptionMeteringConfig{
			Enabled:         resourceUsageConfig["enable_resource_consumption_metering"].(bool),
			ForceSendFields: []string{"Enabled"},
		},
		ForceSendFields: []string{"EnableNetworkEgressMetering"},
	}
	if _, ok := resourceUsageConfig["bigquery_destination"]; ok {
		destinationArr := resourceUsageConfig["bigquery_destination"].([]interface{})
		if len(destinationArr) > 0 && destinationArr[0] != nil {
			bigqueryDestination := destinationArr[0].(map[string]interface{})
			if _, ok := bigqueryDestination["dataset_id"]; ok {
				result.BigqueryDestination = &containerBeta.BigQueryDestination{
					DatasetId: bigqueryDestination["dataset_id"].(string),
				}
			}
		}
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

	if c.CloudRunConfig != nil {
		cloudRunConfig := map[string]interface{}{
			"disabled": c.CloudRunConfig.Disabled,
		}
		if c.CloudRunConfig.LoadBalancerType == "LOAD_BALANCER_TYPE_INTERNAL" {
			// Currently we only allow setting load_balancer_type to LOAD_BALANCER_TYPE_INTERNAL
			cloudRunConfig["load_balancer_type"] = "LOAD_BALANCER_TYPE_INTERNAL"
		}
		result["cloudrun_config"] = []map[string]interface{}{cloudRunConfig}
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
			"peering_name":            c.PeeringName,
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

func flattenReleaseChannel(c *containerBeta.ReleaseChannel) []map[string]interface{} {
	result := []map[string]interface{}{}
	if c != nil {
		result = append(result, map[string]interface{}{
			"channel": c.Channel,
		})
	} else {
		// Explicitly set the release channel to the default.
		result = append(result, map[string]interface{}{
			"channel": "UNSPECIFIED",
		})
	}
	return result
}

func flattenWorkloadIdentityConfig(c *containerBeta.WorkloadIdentityConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"identity_namespace": c.IdentityNamespace,
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
	if mp.Window.RecurringWindow != nil {
		return []map[string]interface{}{
			{
				"recurring_window": []map[string]interface{}{
					{
						"start_time": mp.Window.RecurringWindow.Window.StartTime,
						"end_time":   mp.Window.RecurringWindow.Window.EndTime,
						"recurrence": mp.Window.RecurringWindow.Recurrence,
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
	if a == nil {
		r["enabled"] = false
		return []map[string]interface{}{r}
	}

	if a.EnableNodeAutoprovisioning {
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
	} else {
		r["enabled"] = false
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

func flattenResourceUsageExportConfig(c *containerBeta.ResourceUsageExportConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}

	enableResourceConsumptionMetering := false
	if c.ConsumptionMeteringConfig != nil && c.ConsumptionMeteringConfig.Enabled == true {
		enableResourceConsumptionMetering = true
	}

	return []map[string]interface{}{
		{
			"enable_network_egress_metering":       c.EnableNetworkEgressMetering,
			"enable_resource_consumption_metering": enableResourceConsumptionMetering,
			"bigquery_destination": []map[string]interface{}{
				{"dataset_id": c.BigqueryDestination.DatasetId},
			},
		},
	}
}

func flattenDatabaseEncryption(c *containerBeta.DatabaseEncryption) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"state":    c.State,
			"key_name": c.KeyName,
		},
	}
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
	if _, err := containerClusterAwaitRestingState(config, project, location, clusterName, d.Timeout(schema.TimeoutCreate)); err != nil {
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
		// https://github.com/hashicorp/terraform-provider-google/issues/4186
		blockValueKnown := d.NewValueKnown("private_cluster_config.0.master_ipv4_cidr_block")

		if blockValueKnown && (block == nil || block == "") {
			return fmt.Errorf("master_ipv4_cidr_block must be set if enable_private_nodes == true")
		}
	}
	return nil
}
