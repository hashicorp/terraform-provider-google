// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataproc

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/dataproc/v1"
)

var (
	resolveDataprocImageVersion = regexp.MustCompile(`(?P<Major>[^\s.-]+)\.(?P<Minor>[^\s.-]+)(?:\.(?P<Subminor>[^\s.-]+))?(?:\-(?P<Distr>[^\s.-]+))?`)

	virtualClusterConfigKeys = []string{
		"virtual_cluster_config.0.staging_bucket",
		"virtual_cluster_config.0.auxiliary_services_config",
		"virtual_cluster_config.0.kubernetes_cluster_config",
	}

	auxiliaryServicesConfigKeys = []string{
		"virtual_cluster_config.0.auxiliary_services_config.0.metastore_config",
		"virtual_cluster_config.0.auxiliary_services_config.0.spark_history_server_config",
	}

	auxiliaryServicesMetastoreConfigKeys = []string{
		"virtual_cluster_config.0.auxiliary_services_config.0.metastore_config.0.dataproc_metastore_service",
	}

	auxiliaryServicesSparkHistoryServerConfigKeys = []string{
		"virtual_cluster_config.0.auxiliary_services_config.0.spark_history_server_config.0.dataproc_cluster",
	}

	kubernetesClusterConfigKeys = []string{
		"virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_namespace",
	}

	gkeClusterConfigKeys = []string{
		"virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.gke_cluster_target",
		"virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target",
	}

	gceClusterConfigKeys = []string{
		"cluster_config.0.gce_cluster_config.0.zone",
		"cluster_config.0.gce_cluster_config.0.network",
		"cluster_config.0.gce_cluster_config.0.subnetwork",
		"cluster_config.0.gce_cluster_config.0.tags",
		"cluster_config.0.gce_cluster_config.0.service_account",
		"cluster_config.0.gce_cluster_config.0.service_account_scopes",
		"cluster_config.0.gce_cluster_config.0.internal_ip_only",
		"cluster_config.0.gce_cluster_config.0.shielded_instance_config",
		"cluster_config.0.gce_cluster_config.0.metadata",
		"cluster_config.0.gce_cluster_config.0.reservation_affinity",
		"cluster_config.0.gce_cluster_config.0.node_group_affinity",
	}

	schieldedInstanceConfigKeys = []string{
		"cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_secure_boot",
		"cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_vtpm",
		"cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_integrity_monitoring",
	}

	reservationAffinityKeys = []string{
		"cluster_config.0.gce_cluster_config.0.reservation_affinity.0.consume_reservation_type",
		"cluster_config.0.gce_cluster_config.0.reservation_affinity.0.key",
		"cluster_config.0.gce_cluster_config.0.reservation_affinity.0.values",
	}

	preemptibleWorkerDiskConfigKeys = []string{
		"cluster_config.0.preemptible_worker_config.0.disk_config.0.num_local_ssds",
		"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_size_gb",
		"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_type",
	}

	clusterSoftwareConfigKeys = []string{
		"cluster_config.0.software_config.0.image_version",
		"cluster_config.0.software_config.0.override_properties",
		"cluster_config.0.software_config.0.optional_components",
	}

	dataprocMetricConfigKeys = []string{
		"cluster_config.0.dataproc_metric_config.0.metrics",
	}

	metricKeys = []string{
		"cluster_config.0.dataproc_metric_config.0.metrics.0.metric_source",
		"cluster_config.0.dataproc_metric_config.0.metrics.0.metric_overrides",
	}

	clusterConfigKeys = []string{
		"cluster_config.0.staging_bucket",
		"cluster_config.0.temp_bucket",
		"cluster_config.0.gce_cluster_config",
		"cluster_config.0.master_config",
		"cluster_config.0.worker_config",
		"cluster_config.0.preemptible_worker_config",
		"cluster_config.0.security_config",
		"cluster_config.0.software_config",
		"cluster_config.0.initialization_action",
		"cluster_config.0.encryption_config",
		"cluster_config.0.autoscaling_config",
		"cluster_config.0.metastore_config",
		"cluster_config.0.lifecycle_config",
		"cluster_config.0.endpoint_config",
		"cluster_config.0.dataproc_metric_config",
	}
)

const resourceDataprocGoogleProvidedLabelPrefix = "labels.goog-dataproc"

func resourceDataprocLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasPrefix(k, resourceDataprocGoogleProvidedLabelPrefix) && new == "" {
		return true
	}

	// Let diff be determined by labels (above)
	if strings.HasPrefix(k, "labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

const resourceDataprocGoogleProvidedDPGKEPrefix = "virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_software_config.0.properties.dpgke"
const resourceDataprocGoogleProvidedSparkPrefix = "virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_software_config.0.properties.spark"

func resourceDataprocPropertyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the properties provided by API
	if strings.HasPrefix(k, resourceDataprocGoogleProvidedDPGKEPrefix) && new == "" {
		return true
	}

	if strings.HasPrefix(k, resourceDataprocGoogleProvidedSparkPrefix) && new == "" {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

func ResourceDataprocCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocClusterCreate,
		Read:   resourceDataprocClusterRead,
		Update: resourceDataprocClusterUpdate,
		Delete: resourceDataprocClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the cluster, unique within the project and zone.`,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 55 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 55 characters", k))
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

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the cluster will exist. If it is not provided, the provider project is used.`,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				ForceNew:    true,
				Description: `The region in which the cluster and associated nodes will be created in. Defaults to global.`,
			},

			"graceful_decommission_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0s",
				Description: `The timeout duration which allows graceful decomissioning when you change the number of worker nodes directly through a terraform apply`,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// GCP automatically adds labels
				DiffSuppressFunc: resourceDataprocLabelDiffSuppress,
				Computed:         true,
				Description:      `The list of labels (key/value pairs) to be applied to instances in the cluster. GCP generates some itself including goog-dataproc-cluster-name which is the name of the cluster.`,
			},

			"virtual_cluster_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The virtual cluster config is used when creating a Dataproc cluster that does not directly control the underlying compute resources, for example, when creating a Dataproc-on-GKE cluster. Dataproc may set default values, and values may change when clusters are updated. Exactly one of config or virtualClusterConfig must be specified.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"staging_bucket": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: virtualClusterConfigKeys,
							ForceNew:     true,
							Description:  `A Cloud Storage bucket used to stage job dependencies, config files, and job driver console output. If you do not specify a staging bucket, Cloud Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's staging bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket.`,
						},

						"auxiliary_services_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							MaxItems:     1,
							AtLeastOneOf: virtualClusterConfigKeys,
							Description:  `Auxiliary services configuration for a Cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"metastore_config": {
										Type:         schema.TypeList,
										Optional:     true,
										MaxItems:     1,
										AtLeastOneOf: auxiliaryServicesConfigKeys,
										Description:  `The Hive Metastore configuration for this workload.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"dataproc_metastore_service": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													AtLeastOneOf: auxiliaryServicesMetastoreConfigKeys,
													Description:  `The Hive Metastore configuration for this workload.`,
												},
											},
										},
									},

									"spark_history_server_config": {
										Type:         schema.TypeList,
										Optional:     true,
										MaxItems:     1,
										AtLeastOneOf: auxiliaryServicesConfigKeys,
										Description:  `The Spark History Server configuration for the workload.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"dataproc_cluster": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													AtLeastOneOf: auxiliaryServicesSparkHistoryServerConfigKeys,
													Description:  `Resource name of an existing Dataproc Cluster to act as a Spark History Server for the workload.`,
												},
											},
										},
									},
								},
							},
						},

						"kubernetes_cluster_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							MaxItems:     1,
							AtLeastOneOf: virtualClusterConfigKeys,
							Description:  `The configuration for running the Dataproc cluster on Kubernetes.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"kubernetes_namespace": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										AtLeastOneOf: kubernetesClusterConfigKeys,
										Description:  `A namespace within the Kubernetes cluster to deploy into. If this namespace does not exist, it is created. If it exists, Dataproc verifies that another Dataproc VirtualCluster is not installed into it. If not specified, the name of the Dataproc Cluster is used.`,
									},

									"kubernetes_software_config": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: `The software configuration for this Dataproc cluster running on Kubernetes.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"component_version": {
													Type:        schema.TypeMap,
													Required:    true,
													ForceNew:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: `The components that should be installed in this Dataproc cluster. The key must be a string from the KubernetesComponent enumeration. The value is the version of the software to be installed.`,
												},

												"properties": {
													Type:             schema.TypeMap,
													Optional:         true,
													ForceNew:         true,
													DiffSuppressFunc: resourceDataprocPropertyDiffSuppress,
													Elem:             &schema.Schema{Type: schema.TypeString},
													Computed:         true,
													Description:      `The properties to set on daemon config files. Property keys are specified in prefix:property format, for example spark:spark.kubernetes.container.image.`,
												},
											},
										},
									},

									"gke_cluster_config": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: `The configuration for running the Dataproc cluster on GKE.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"gke_cluster_target": {
													Type:         schema.TypeString,
													ForceNew:     true,
													Optional:     true,
													AtLeastOneOf: gkeClusterConfigKeys,
													Description:  `A target GKE cluster to deploy to. It must be in the same project and region as the Dataproc cluster (the GKE cluster can be zonal or regional). Format: 'projects/{project}/locations/{location}/clusters/{cluster_id}'`,
												},

												"node_pool_target": {
													Type:         schema.TypeList,
													Optional:     true,
													AtLeastOneOf: gkeClusterConfigKeys,
													MinItems:     1,
													Description:  `GKE node pools where workloads will be scheduled. At least one node pool must be assigned the DEFAULT GkeNodePoolTarget.Role. If a GkeNodePoolTarget is not specified, Dataproc constructs a DEFAULT GkeNodePoolTarget.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{

															"node_pool": {
																Type:        schema.TypeString,
																ForceNew:    true,
																Required:    true,
																Description: `The target GKE node pool. Format: 'projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{nodePool}'`,
															},

															"roles": {
																Type:        schema.TypeSet,
																Elem:        &schema.Schema{Type: schema.TypeString},
																ForceNew:    true,
																Required:    true,
																Description: `The roles associated with the GKE node pool.`,
															},

															"node_pool_config": {
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																MaxItems:    1,
																Description: `Input only. The configuration for the GKE node pool.`,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{

																		"config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Computed:    true,
																			MaxItems:    1,
																			Description: `The node pool configuration.`,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{

																					"machine_type": {
																						Type:        schema.TypeString,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The name of a Compute Engine machine type.`,
																					},

																					"local_ssd_count": {
																						Type:        schema.TypeInt,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The minimum number of nodes in the node pool. Must be >= 0 and <= maxNodeCount.`,
																					},

																					"preemptible": {
																						Type:        schema.TypeBool,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `Whether the nodes are created as preemptible VM instances. Preemptible nodes cannot be used in a node pool with the CONTROLLER role or in the DEFAULT node pool if the CONTROLLER role is not assigned (the DEFAULT node pool will assume the CONTROLLER role).`,
																					},

																					"min_cpu_platform": {
																						Type:        schema.TypeString,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `Minimum CPU platform to be used by this instance. The instance may be scheduled on the specified or a newer CPU platform. Specify the friendly names of CPU platforms, such as "Intel Haswell" or "Intel Sandy Bridge".`,
																					},

																					"spot": {
																						Type:        schema.TypeBool,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `Spot flag for enabling Spot VM, which is a rebrand of the existing preemptible flag.`,
																					},
																				},
																			},
																		},

																		"locations": {
																			Type:        schema.TypeSet,
																			Elem:        &schema.Schema{Type: schema.TypeString},
																			ForceNew:    true,
																			Required:    true,
																			Description: `The list of Compute Engine zones where node pool nodes associated with a Dataproc on GKE virtual cluster will be located.`,
																		},

																		"autoscaling": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Computed:    true,
																			MaxItems:    1,
																			Description: `The autoscaler configuration for this node pool. The autoscaler is enabled only when a valid configuration is present.`,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{

																					"min_node_count": {
																						Type:        schema.TypeInt,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The minimum number of nodes in the node pool. Must be >= 0 and <= maxNodeCount.`,
																					},

																					"max_node_count": {
																						Type:        schema.TypeInt,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The maximum number of nodes in the node pool. Must be >= minNodeCount, and must be > 0.`,
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"cluster_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Allows you to configure various aspects of the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"staging_bucket": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							ForceNew:     true,
							Description:  `The Cloud Storage staging bucket used to stage files, such as Hadoop jars, between client machines and the cluster. Note: If you don't explicitly specify a staging_bucket then GCP will auto create / assign one for you. However, you are not guaranteed an auto generated bucket which is solely dedicated to your cluster; it may be shared with other clusters in the same region/zone also choosing to use the auto generation option.`,
						},
						// If the user does not specify a staging bucket, GCP will allocate one automatically.
						// The staging_bucket field provides a way for the user to supply their own
						// staging bucket. The bucket field is purely a computed field which details
						// the definitive bucket allocated and in use (either the user supplied one via
						// staging_bucket, or the GCP generated one)
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ` The name of the cloud storage bucket ultimately used to house the staging data for the cluster. If staging_bucket is specified, it will contain this value, otherwise it will be the auto generated name.`,
						},

						"temp_bucket": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: clusterConfigKeys,
							ForceNew:     true,
							Description:  `The Cloud Storage temp bucket used to store ephemeral cluster and jobs data, such as Spark and MapReduce history files. Note: If you don't explicitly specify a temp_bucket then GCP will auto create / assign one for you.`,
						},

						"gce_cluster_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `Common config settings for resources of Google Compute Engine cluster instances, applicable to all instances in the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"zone": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Description:  `The GCP zone where your data is stored and used (i.e. where the master and the worker nodes will be created in). If region is set to 'global' (default) then zone is mandatory, otherwise GCP is able to make use of Auto Zone Placement to determine this automatically for you. Note: This setting additionally determines and restricts which computing resources are available for use with other configs such as cluster_config.master_config.machine_type and cluster_config.worker_config.machine_type.`,
									},

									"network": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										AtLeastOneOf:     gceClusterConfigKeys,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.subnetwork"},
										DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
										Description:      `The name or self_link of the Google Compute Engine network to the cluster will be part of. Conflicts with subnetwork. If neither is specified, this defaults to the "default" network.`,
									},

									"subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										AtLeastOneOf:     gceClusterConfigKeys,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.network"},
										DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
										Description:      `The name or self_link of the Google Compute Engine subnetwork the cluster will be part of. Conflicts with network.`,
									},

									"tags": {
										Type:         schema.TypeSet,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Elem:         &schema.Schema{Type: schema.TypeString},
										Description:  `The list of instance tags applied to instances in the cluster. Tags are used to identify valid sources or targets for network firewalls.`,
									},

									"service_account": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Description:  `The service account to be used by the Node VMs. If not specified, the "default" service account is used.`,
									},

									"service_account_scopes": {
										Type:         schema.TypeSet,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Description:  `The set of Google API scopes to be made available on all of the node VMs under the service_account specified. These can be either FQDNs, or scope aliases.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											StateFunc: func(v interface{}) string {
												return tpgresource.CanonicalizeServiceScope(v.(string))
											},
										},
										Set: tpgresource.StringScopeHashcode,
									},

									"internal_ip_only": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Default:      false,
										Description:  `By default, clusters are not restricted to internal IP addresses, and will have ephemeral external IP addresses assigned to each instance. If set to true, all instances in the cluster will only have internal IP addresses. Note: Private Google Access (also known as privateIpGoogleAccess) must be enabled on the subnetwork that the cluster will be launched in.`,
									},

									"metadata": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ForceNew:     true,
										Description:  `A map of the Compute Engine metadata entries to add to all instances`,
									},

									"shielded_instance_config": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Computed:     true,
										MaxItems:     1,
										Description:  `Shielded Instance Config for clusters using Compute Engine Shielded VMs.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_secure_boot": {
													Type:         schema.TypeBool,
													Optional:     true,
													Default:      false,
													AtLeastOneOf: schieldedInstanceConfigKeys,
													ForceNew:     true,
													Description:  `Defines whether instances have Secure Boot enabled.`,
												},
												"enable_vtpm": {
													Type:         schema.TypeBool,
													Optional:     true,
													Default:      false,
													AtLeastOneOf: schieldedInstanceConfigKeys,
													ForceNew:     true,
													Description:  `Defines whether instances have the vTPM enabled.`,
												},
												"enable_integrity_monitoring": {
													Type:         schema.TypeBool,
													Optional:     true,
													Default:      false,
													AtLeastOneOf: schieldedInstanceConfigKeys,
													ForceNew:     true,
													Description:  `Defines whether instances have integrity monitoring enabled.`,
												},
											},
										},
									},

									"reservation_affinity": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Computed:     true,
										MaxItems:     1,
										Description:  `Reservation Affinity for consuming Zonal reservation.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"consume_reservation_type": {
													Type:         schema.TypeString,
													Optional:     true,
													AtLeastOneOf: reservationAffinityKeys,
													ForceNew:     true,
													ValidateFunc: validation.StringInSlice([]string{"NO_RESERVATION", "ANY_RESERVATION", "SPECIFIC_RESERVATION"}, false),
													Description:  `Type of reservation to consume.`,
												},
												"key": {
													Type:         schema.TypeString,
													Optional:     true,
													AtLeastOneOf: reservationAffinityKeys,
													ForceNew:     true,
													Description:  `Corresponds to the label key of reservation resource.`,
												},
												"values": {
													Type:         schema.TypeSet,
													Elem:         &schema.Schema{Type: schema.TypeString},
													Optional:     true,
													AtLeastOneOf: reservationAffinityKeys,
													ForceNew:     true,
													Description:  `Corresponds to the label values of reservation resource.`,
												},
											},
										},
									},

									"node_group_affinity": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Computed:     true,
										MaxItems:     1,
										Description:  `Node Group Affinity for sole-tenant clusters.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_group_uri": {
													Type:             schema.TypeString,
													ForceNew:         true,
													Required:         true,
													Description:      `The URI of a sole-tenant that the cluster will be created on.`,
													DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
												},
											},
										},
									},
								},
							},
						},

						"master_config": instanceConfigSchema("master_config"),
						"worker_config": instanceConfigSchema("worker_config"),
						// preemptible_worker_config has a slightly different config
						"preemptible_worker_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `The Google Compute Engine config settings for the additional (aka preemptible) instances in a cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_instances": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the number of preemptible nodes to create. Defaults to 0.`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
									},

									// API does not honour this if set ...
									// It always uses whatever is specified for the worker_config
									// "machine_type": { ... }
									// "min_cpu_platform": { ... }
									"preemptibility": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the preemptibility of the secondary nodes. Defaults to PREEMPTIBLE.`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice([]string{"PREEMPTIBILITY_UNSPECIFIED", "NON_PREEMPTIBLE", "PREEMPTIBLE", "SPOT"}, false),
										Default:      "PREEMPTIBLE",
									},

									"disk_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: `Disk Config`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
										MaxItems: 1,

										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_local_ssds": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													AtLeastOneOf: preemptibleWorkerDiskConfigKeys,
													ForceNew:     true,
													Description:  `The amount of local SSD disks that will be attached to each preemptible worker node. Defaults to 0.`,
												},

												"boot_disk_size_gb": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													AtLeastOneOf: preemptibleWorkerDiskConfigKeys,
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(10),
													Description:  `Size of the primary disk attached to each preemptible worker node, specified in GB. The smallest allowed disk size is 10GB. GCP will default to a predetermined computed value if not set (currently 500GB). Note: If SSDs are not attached, it also contains the HDFS data blocks and Hadoop working directories.`,
												},

												"boot_disk_type": {
													Type:         schema.TypeString,
													Optional:     true,
													AtLeastOneOf: preemptibleWorkerDiskConfigKeys,
													ForceNew:     true,
													Default:      "pd-standard",
													Description:  `The disk type of the primary disk attached to each preemptible worker node. Such as "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
												},
											},
										},
									},

									"instance_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `List of preemptible instance names which have been assigned to the cluster.`,
									},
								},
							},
						},

						"security_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Security related configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kerberos_config": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Kerberos related configuration",
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cross_realm_trust_admin_server": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The admin server (IP or hostname) for the remote trusted realm in a cross realm trust relationship.`,
												},
												"cross_realm_trust_kdc": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The KDC (IP or hostname) for the remote trusted realm in a cross realm trust relationship.`,
												},
												"cross_realm_trust_realm": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The remote realm the Dataproc on-cluster KDC will trust, should the user enable cross realm trust.`,
												},
												"cross_realm_trust_shared_password_uri": {
													Type:     schema.TypeString,
													Optional: true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the shared password between the on-cluster
Kerberos realm and the remote trusted realm, in a cross realm trust relationship.`,
												},
												"enable_kerberos": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: `Flag to indicate whether to Kerberize the cluster.`,
												},
												"kdc_db_key_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the master key of the KDC database.`,
												},
												"key_password_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the password to the user provided key. For the self-signed certificate, this password is generated by Dataproc.`,
												},
												"keystore_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of the keystore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.`,
												},
												"keystore_password_uri": {
													Type:     schema.TypeString,
													Optional: true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing
the password to the user provided keystore. For the self-signed certificate, this password is generated
by Dataproc`,
												},
												"kms_key_uri": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The uri of the KMS key used to encrypt various sensitive files.`,
												},
												"realm": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The name of the on-cluster Kerberos realm. If not specified, the uppercased domain of hostnames will be the realm.`,
												},
												"root_principal_password_uri": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The cloud Storage URI of a KMS encrypted file containing the root principal password.`,
												},
												"tgt_lifetime_hours": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: `The lifetime of the ticket granting ticket, in hours.`,
												},
												"truststore_password_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the password to the user provided truststore. For the self-signed certificate, this password is generated by Dataproc.`,
												},
												"truststore_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of the truststore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.`,
												},
											},
										},
									},
								},
							},
						},

						"software_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `The config settings for software inside the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_version": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										AtLeastOneOf:     clusterSoftwareConfigKeys,
										ForceNew:         true,
										DiffSuppressFunc: dataprocImageVersionDiffSuppress,
										Description:      `The Cloud Dataproc image version to use for the cluster - this controls the sets of software versions installed onto the nodes when you create clusters. If not specified, defaults to the latest version.`,
									},
									"override_properties": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: clusterSoftwareConfigKeys,
										ForceNew:     true,
										Elem:         &schema.Schema{Type: schema.TypeString},
										Description:  `A list of override and additional properties (key/value pairs) used to modify various aspects of the common configuration files used when creating a cluster.`,
									},

									"properties": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `A list of the properties used to set the daemon config files. This will include any values supplied by the user via cluster_config.software_config.override_properties`,
									},

									// We have two versions of the properties field here because by default
									// dataproc will set a number of default properties for you out of the
									// box. If you want to override one or more, if we only had one field,
									// you would need to add in all these values as well otherwise you would
									// get a diff. To make this easier, 'properties' simply contains the computed
									// values (including overrides) for all properties, whilst override_properties
									// is only for properties the user specifically wants to override. If nothing
									// is overridden, this will be empty.

									"optional_components": {
										Type:         schema.TypeSet,
										Optional:     true,
										AtLeastOneOf: clusterSoftwareConfigKeys,
										Description:  `The set of optional components to activate on the cluster.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"initialization_action": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							ForceNew:     true,
							Description:  `Commands to execute on each node after config is completed. You can specify multiple versions of these.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"script": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `The script to be executed during initialization of the cluster. The script must be a GCS file with a gs:// prefix.`,
									},

									"timeout_sec": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     300,
										ForceNew:    true,
										Description: `The maximum duration (in seconds) which script is allowed to take to execute its action. GCP will default to a predetermined computed value if not set (currently 300).`,
									},
								},
							},
						},
						"encryption_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							MaxItems:     1,
							Description:  `The Customer managed encryption keys settings for the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The Cloud KMS key name to use for PD disk encryption for all instances in the cluster.`,
									},
								},
							},
						},
						"autoscaling_config": {
							Type:             schema.TypeList,
							Optional:         true,
							AtLeastOneOf:     clusterConfigKeys,
							MaxItems:         1,
							Description:      `The autoscaling policy config associated with the cluster.`,
							DiffSuppressFunc: tpgresource.EmptyOrUnsetBlockDiffSuppress,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_uri": {
										Type:             schema.TypeString,
										Required:         true,
										Description:      `The autoscaling policy used by the cluster.`,
										DiffSuppressFunc: tpgresource.LocationDiffSuppress,
									},
								},
							},
						},
						"metastore_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							MaxItems:     1,
							Description:  `Specifies a Metastore configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataproc_metastore_service": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `Resource name of an existing Dataproc Metastore service.`,
									},
								},
							},
						},
						"lifecycle_config": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							AtLeastOneOf: clusterConfigKeys,
							Description:  `The settings for auto deletion cluster schedule.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"idle_delete_ttl": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The duration to keep the cluster alive while idling (no jobs running). After this TTL, the cluster will be deleted. Valid range: [10m, 14d].`,
										AtLeastOneOf: []string{
											"cluster_config.0.lifecycle_config.0.idle_delete_ttl",
											"cluster_config.0.lifecycle_config.0.auto_delete_time",
										},
									},
									"idle_start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Time when the cluster became idle (most recent job finished) and became eligible for deletion due to idleness.`,
									},
									// the API also has the auto_delete_ttl option in its request, however,
									// the value is not returned in the response, rather the auto_delete_time
									// after calculating ttl with the update time is returned, thus, for now
									// we will only allow auto_delete_time to updated.
									"auto_delete_time": {
										Type:             schema.TypeString,
										Optional:         true,
										Description:      `The time when cluster will be auto-deleted. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".`,
										DiffSuppressFunc: tpgresource.TimestampDiffSuppress(time.RFC3339Nano),
										AtLeastOneOf: []string{
											"cluster_config.0.lifecycle_config.0.idle_delete_ttl",
											"cluster_config.0.lifecycle_config.0.auto_delete_time",
										},
									},
								},
							},
						},
						"endpoint_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							MaxItems:     1,
							Description:  `The config settings for port access on the cluster. Structure defined below.`,
							AtLeastOneOf: clusterConfigKeys,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_http_port_access": {
										Type:        schema.TypeBool,
										Required:    true,
										ForceNew:    true,
										Description: `The flag to enable http access to specific ports on the cluster from external sources (aka Component Gateway). Defaults to false.`,
									},
									"http_ports": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `The map of port descriptions to URLs. Will only be populated if enable_http_port_access is true.`,
									},
								},
							},
						},

						"dataproc_metric_config": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Description:  `The config for Dataproc metrics.`,
							AtLeastOneOf: clusterConfigKeys,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metrics": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `Metrics sources to enable.`,
										Elem:        metricsSchema(),
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

// We need to pull metrics' schema out so we can use it to make a set hash func
func metricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric_source": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MONITORING_AGENT_DEFAULTS", "HDFS", "SPARK", "YARN", "SPARK_HISTORY_SERVER", "HIVESERVER2"}, false),
				Description:  `A source for the collection of Dataproc OSS metrics (see [available OSS metrics] (https://cloud.google.com//dataproc/docs/guides/monitoring#available_oss_metrics)).`,
			},
			"metric_overrides": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: `Specify one or more [available OSS metrics] (https://cloud.google.com/dataproc/docs/guides/monitoring#available_oss_metrics) to collect.`,
			},
		},
	}
}

func instanceConfigSchema(parent string) *schema.Schema {
	var instanceConfigKeys = []string{
		"cluster_config.0." + parent + ".0.num_instances",
		"cluster_config.0." + parent + ".0.image_uri",
		"cluster_config.0." + parent + ".0.machine_type",
		"cluster_config.0." + parent + ".0.min_cpu_platform",
		"cluster_config.0." + parent + ".0.disk_config",
		"cluster_config.0." + parent + ".0.accelerators",
	}

	masterConfig := strings.Contains(parent, "master")

	return &schema.Schema{
		Type:         schema.TypeList,
		Optional:     true,
		Computed:     true,
		AtLeastOneOf: clusterConfigKeys,
		MaxItems:     1,
		Description:  `The Google Compute Engine config settings for the master/worker instances in a cluster.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"num_instances": {
					Type:         schema.TypeInt,
					Optional:     true,
					ForceNew:     masterConfig,
					Computed:     true,
					Description:  `Specifies the number of master/worker nodes to create. If not specified, GCP will default to a predetermined computed value.`,
					AtLeastOneOf: instanceConfigKeys,
				},

				"image_uri": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					AtLeastOneOf: instanceConfigKeys,
					ForceNew:     true,
					Description:  `The URI for the image to use for this master/worker`,
				},

				"machine_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					AtLeastOneOf: instanceConfigKeys,
					ForceNew:     true,
					Description:  `The name of a Google Compute Engine machine type to create for the master/worker`,
				},

				"min_cpu_platform": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					AtLeastOneOf: instanceConfigKeys,
					ForceNew:     true,
					Description:  `The name of a minimum generation of CPU family for the master/worker. If not specified, GCP will default to a predetermined computed value for each zone.`,
				},
				"disk_config": {
					Type:         schema.TypeList,
					Optional:     true,
					Computed:     true,
					AtLeastOneOf: instanceConfigKeys,
					MaxItems:     1,
					Description:  `Disk Config`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"num_local_ssds": {
								Type:        schema.TypeInt,
								Optional:    true,
								Computed:    true,
								Description: `The amount of local SSD disks that will be attached to each master cluster node. Defaults to 0.`,
								AtLeastOneOf: []string{
									"cluster_config.0." + parent + ".0.disk_config.0.num_local_ssds",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_size_gb",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_type",
								},
								ForceNew: true,
							},

							"boot_disk_size_gb": {
								Type:        schema.TypeInt,
								Optional:    true,
								Computed:    true,
								Description: `Size of the primary disk attached to each node, specified in GB. The primary disk contains the boot volume and system libraries, and the smallest allowed disk size is 10GB. GCP will default to a predetermined computed value if not set (currently 500GB). Note: If SSDs are not attached, it also contains the HDFS data blocks and Hadoop working directories.`,
								AtLeastOneOf: []string{
									"cluster_config.0." + parent + ".0.disk_config.0.num_local_ssds",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_size_gb",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_type",
								},
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(10),
							},

							"boot_disk_type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: `The disk type of the primary disk attached to each node. Such as "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
								AtLeastOneOf: []string{
									"cluster_config.0." + parent + ".0.disk_config.0.num_local_ssds",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_size_gb",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_type",
								},
								ForceNew: true,
								Default:  "pd-standard",
							},
						},
					},
				},

				// Note: preemptible workers don't support accelerators
				"accelerators": {
					Type:         schema.TypeSet,
					Optional:     true,
					AtLeastOneOf: instanceConfigKeys,
					ForceNew:     true,
					Elem:         acceleratorsSchema(),
					Description:  `The Compute Engine accelerator (GPU) configuration for these instances. Can be specified multiple times.`,
				},

				"instance_names": {
					Type:        schema.TypeList,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `List of master/worker instance names which have been assigned to the cluster.`,
				},
			},
		},
	}
}

// We need to pull accelerators' schema out so we can use it to make a set hash func
func acceleratorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerator_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The short name of the accelerator type to expose to this instance. For example, nvidia-tesla-k80.`,
			},

			"accelerator_count": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The number of the accelerator cards of this type exposed to this instance. Often restricted to one of 1, 2, 4, or 8.`,
			},
		},
	}
}

func resourceDataprocClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	cluster := &dataproc.Cluster{
		ClusterName: d.Get("name").(string),
		ProjectId:   project,
	}

	if _, ok := d.GetOk("virtual_cluster_config"); ok {
		cluster.VirtualClusterConfig, err = expandVirtualClusterConfig(d, config)
	} else {
		cluster.Config, err = expandClusterConfig(d, config)
	}

	if err != nil {
		return err
	}

	if _, ok := d.GetOk("labels"); ok {
		cluster.Labels = tpgresource.ExpandLabels(d)
	}

	// Checking here caters for the case where the user does not specify cluster_config
	// at all, as well where it is simply missing from the gce_cluster_config
	if region == "global" && cluster.Config.GceClusterConfig.ZoneUri == "" {
		return errors.New("zone is mandatory when region is set to 'global'")
	}

	// Create the cluster
	op, err := config.NewDataprocClient(userAgent).Projects.Regions.Clusters.Create(
		project, region, cluster).Do()
	if err != nil {
		return fmt.Errorf("Error creating Dataproc cluster: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/clusters/%s", project, region, cluster.ClusterName))

	// Wait until it's created
	waitErr := DataprocClusterOperationWait(config, op, "creating Dataproc cluster", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		// The resource didn't actually create
		// Note that we do not remove the ID here - this resource tends to leave
		// partially created clusters behind, so we'll let the next Read remove
		// it.
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been created", cluster.ClusterName)
	return resourceDataprocClusterRead(d, meta)
}

func expandVirtualClusterConfig(d *schema.ResourceData, config *transport_tpg.Config) (*dataproc.VirtualClusterConfig, error) {
	conf := &dataproc.VirtualClusterConfig{}

	if v, ok := d.GetOk("virtual_cluster_config"); ok {
		confs := v.([]interface{})
		if (len(confs)) == 0 {
			return conf, nil
		}
	}

	if v, ok := d.GetOk("virtual_cluster_config.0.staging_bucket"); ok {
		conf.StagingBucket = v.(string)
	}

	if cfg, ok := configOptions(d, "virtual_cluster_config.0.auxiliary_services_config"); ok {
		conf.AuxiliaryServicesConfig = expandAuxiliaryServicesConfig(d, cfg)
	}

	if cfg, ok := configOptions(d, "virtual_cluster_config.0.kubernetes_cluster_config"); ok {
		conf.KubernetesClusterConfig = expandKubernetesClusterConfig(d, cfg)
	}
	return conf, nil
}

func expandAuxiliaryServicesConfig(d *schema.ResourceData, cfg map[string]interface{}) *dataproc.AuxiliaryServicesConfig {
	conf := &dataproc.AuxiliaryServicesConfig{}
	if mcfg, ok := configOptions(d, "virtual_cluster_config.0.auxiliary_services_config.0.metastore_config"); ok {
		conf.MetastoreConfig = expandVCMetastoreConfig(mcfg)
	}

	if shscfg, ok := configOptions(d, "virtual_cluster_config.0.auxiliary_services_config.0.spark_history_server_config"); ok {
		conf.SparkHistoryServerConfig = expandSparkHistoryServerConfig(shscfg)
	}
	return conf
}

func expandVCMetastoreConfig(cfg map[string]interface{}) *dataproc.MetastoreConfig {
	conf := &dataproc.MetastoreConfig{}
	if v, ok := cfg["dataproc_metastore_service"]; ok {
		conf.DataprocMetastoreService = v.(string)
	}
	return conf
}

func expandSparkHistoryServerConfig(cfg map[string]interface{}) *dataproc.SparkHistoryServerConfig {
	conf := &dataproc.SparkHistoryServerConfig{}
	if v, ok := cfg["dataproc_cluster"]; ok {
		conf.DataprocCluster = v.(string)
	}
	return conf
}

func expandKubernetesClusterConfig(d *schema.ResourceData, cfg map[string]interface{}) *dataproc.KubernetesClusterConfig {
	conf := &dataproc.KubernetesClusterConfig{}

	if v, ok := cfg["kubernetes_namespace"]; ok {
		conf.KubernetesNamespace = v.(string)
	}

	if kscfg, ok := d.GetOk("virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_software_config"); ok {
		conf.KubernetesSoftwareConfig = expandKubernetesSoftwareConfig(kscfg.([]interface{})[0].(map[string]interface{}))
	}

	if gkeccfg, ok := d.GetOk("virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config"); ok {
		conf.GkeClusterConfig = expandGkeClusterConfig(d, gkeccfg.([]interface{})[0].(map[string]interface{}))
	}
	return conf
}

func expandKubernetesSoftwareConfig(cfg map[string]interface{}) *dataproc.KubernetesSoftwareConfig {
	conf := &dataproc.KubernetesSoftwareConfig{}
	if compSet, ok := cfg["component_version"]; ok {
		components := map[string]string{}

		for k, val := range compSet.(map[string]interface{}) {
			components[k] = val.(string)
		}

		conf.ComponentVersion = components
	}

	if propSet, ok := cfg["properties"]; ok {
		properties := map[string]string{}

		for k, val := range propSet.(map[string]interface{}) {
			properties[k] = val.(string)
		}

		conf.Properties = properties
	}
	return conf
}

func expandGkeClusterConfig(d *schema.ResourceData, cfg map[string]interface{}) *dataproc.GkeClusterConfig {
	conf := &dataproc.GkeClusterConfig{}

	if clusterAddress, ok := cfg["gke_cluster_target"]; ok {
		conf.GkeClusterTarget = clusterAddress.(string)

		if v, ok := d.GetOk("virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target"); ok {
			conf.NodePoolTarget = expandGkeNodePoolTarget(d, v, clusterAddress.(string))
		}
	}
	return conf
}

func expandGkeNodePoolTarget(d *schema.ResourceData, v interface{}, clusterAddress string) []*dataproc.GkeNodePoolTarget {
	nodePools := v.([]interface{})

	nodePoolList := []*dataproc.GkeNodePoolTarget{}
	for i, v1 := range nodePools {
		data := v1.(map[string]interface{})
		nodePool := dataproc.GkeNodePoolTarget{
			NodePool: clusterAddress + "/nodePools/" + data["node_pool"].(string),
			Roles:    tpgresource.ConvertStringSet(data["roles"].(*schema.Set)),
		}

		if v, ok := d.GetOk(fmt.Sprintf("virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target.%d.node_pool_config", i)); ok {
			nodePool.NodePoolConfig = expandGkeNodePoolConfig(v.([]interface{})[0].(map[string]interface{}))
		}

		nodePoolList = append(nodePoolList, &nodePool)
	}

	return nodePoolList
}

func expandGkeNodePoolConfig(cfg map[string]interface{}) *dataproc.GkeNodePoolConfig {
	conf := &dataproc.GkeNodePoolConfig{}

	if nodecfg, ok := cfg["config"]; ok {
		conf.Config = expandGkeNodeConfig(nodecfg.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := cfg["locations"]; ok {
		conf.Locations = tpgresource.ConvertStringSet(v.(*schema.Set))
	}

	if autoscalingcfg, ok := cfg["autoscaling"]; ok {
		conf.Autoscaling = expandGkeNodePoolAutoscalingConfig(autoscalingcfg.([]interface{})[0].(map[string]interface{}))
	}
	return conf
}

func expandGkeNodeConfig(cfg map[string]interface{}) *dataproc.GkeNodeConfig {
	conf := &dataproc.GkeNodeConfig{}

	if v, ok := cfg["local_ssd_count"]; ok {
		conf.LocalSsdCount = int64(v.(int))
	}

	if v, ok := cfg["machine_type"]; ok {
		conf.MachineType = v.(string)
	}

	if v, ok := cfg["preemptible"]; ok {
		conf.Preemptible = v.(bool)
	}

	if v, ok := cfg["min_cpu_platform"]; ok {
		conf.MinCpuPlatform = v.(string)
	}

	if v, ok := cfg["spot"]; ok {
		conf.Spot = v.(bool)
	}
	return conf
}

func expandGkeNodePoolAutoscalingConfig(cfg map[string]interface{}) *dataproc.GkeNodePoolAutoscalingConfig {
	conf := &dataproc.GkeNodePoolAutoscalingConfig{}

	if v, ok := cfg["min_node_count"]; ok {
		conf.MinNodeCount = int64(v.(int))
	}

	if v, ok := cfg["max_node_count"]; ok {
		conf.MaxNodeCount = int64(v.(int))
	}
	return conf
}

func expandClusterConfig(d *schema.ResourceData, config *transport_tpg.Config) (*dataproc.ClusterConfig, error) {
	conf := &dataproc.ClusterConfig{
		// SDK requires GceClusterConfig to be specified,
		// even if no explicit values specified
		GceClusterConfig: &dataproc.GceClusterConfig{},
	}

	if v, ok := d.GetOk("cluster_config"); ok {
		confs := v.([]interface{})
		if (len(confs)) == 0 {
			return conf, nil
		}
	}

	if v, ok := d.GetOk("cluster_config.0.staging_bucket"); ok {
		conf.ConfigBucket = v.(string)
	}

	if v, ok := d.GetOk("cluster_config.0.temp_bucket"); ok {
		conf.TempBucket = v.(string)
	}

	c, err := expandGceClusterConfig(d, config)
	if err != nil {
		return nil, err
	}
	conf.GceClusterConfig = c

	if cfg, ok := configOptions(d, "cluster_config.0.security_config"); ok {
		conf.SecurityConfig = expandSecurityConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.software_config"); ok {
		conf.SoftwareConfig = expandSoftwareConfig(cfg)
	}

	if v, ok := d.GetOk("cluster_config.0.initialization_action"); ok {
		conf.InitializationActions = expandInitializationActions(v)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.encryption_config"); ok {
		conf.EncryptionConfig = expandEncryptionConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.autoscaling_config"); ok {
		conf.AutoscalingConfig = expandAutoscalingConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.metastore_config"); ok {
		conf.MetastoreConfig = expandMetastoreConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.lifecycle_config"); ok {
		conf.LifecycleConfig = expandLifecycleConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.endpoint_config"); ok {
		conf.EndpointConfig = expandEndpointConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.dataproc_metric_config"); ok {
		conf.DataprocMetricConfig = expandDataprocMetricConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.master_config"); ok {
		log.Println("[INFO] got master_config")
		conf.MasterConfig = expandInstanceGroupConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.worker_config"); ok {
		log.Println("[INFO] got worker config")
		conf.WorkerConfig = expandInstanceGroupConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.preemptible_worker_config"); ok {
		log.Println("[INFO] got preemptible worker config")
		conf.SecondaryWorkerConfig = expandPreemptibleInstanceGroupConfig(cfg)
	}
	return conf, nil
}

func expandGceClusterConfig(d *schema.ResourceData, config *transport_tpg.Config) (*dataproc.GceClusterConfig, error) {
	conf := &dataproc.GceClusterConfig{}

	v, ok := d.GetOk("cluster_config.0.gce_cluster_config")
	if !ok {
		return conf, nil
	}
	cfg := v.([]interface{})[0].(map[string]interface{})

	if v, ok := cfg["zone"]; ok {
		conf.ZoneUri = v.(string)
	}
	if v, ok := cfg["network"]; ok {
		nf, err := tpgresource.ParseNetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", v, err)
		}

		conf.NetworkUri = nf.RelativeLink()
	}
	if v, ok := cfg["subnetwork"]; ok {
		snf, err := tpgresource.ParseSubnetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", v, err)
		}

		conf.SubnetworkUri = snf.RelativeLink()
	}
	if v, ok := cfg["tags"]; ok {
		conf.Tags = tpgresource.ConvertStringSet(v.(*schema.Set))
	}
	if v, ok := cfg["service_account"]; ok {
		conf.ServiceAccount = v.(string)
	}
	if scopes, ok := cfg["service_account_scopes"]; ok {
		scopesSet := scopes.(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, scope := range scopesSet.List() {
			scopes[i] = tpgresource.CanonicalizeServiceScope(scope.(string))
		}
		conf.ServiceAccountScopes = scopes
	}
	if v, ok := cfg["internal_ip_only"]; ok {
		conf.InternalIpOnly = v.(bool)
	}
	if v, ok := cfg["metadata"]; ok {
		conf.Metadata = tpgresource.ConvertStringMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("cluster_config.0.gce_cluster_config.0.shielded_instance_config"); ok {
		cfgSic := v.([]interface{})[0].(map[string]interface{})
		conf.ShieldedInstanceConfig = &dataproc.ShieldedInstanceConfig{}
		if v, ok := cfgSic["enable_integrity_monitoring"]; ok {
			conf.ShieldedInstanceConfig.EnableIntegrityMonitoring = v.(bool)
		}
		if v, ok := cfgSic["enable_secure_boot"]; ok {
			conf.ShieldedInstanceConfig.EnableSecureBoot = v.(bool)
		}
		if v, ok := cfgSic["enable_vtpm"]; ok {
			conf.ShieldedInstanceConfig.EnableVtpm = v.(bool)
		}
	}
	if v, ok := d.GetOk("cluster_config.0.gce_cluster_config.0.reservation_affinity"); ok {
		cfgRa := v.([]interface{})[0].(map[string]interface{})
		conf.ReservationAffinity = &dataproc.ReservationAffinity{}
		if v, ok := cfgRa["consume_reservation_type"]; ok {
			conf.ReservationAffinity.ConsumeReservationType = v.(string)
		}
		if v, ok := cfgRa["key"]; ok {
			conf.ReservationAffinity.Key = v.(string)
		}
		if v, ok := cfgRa["values"]; ok {
			conf.ReservationAffinity.Values = tpgresource.ConvertStringSet(v.(*schema.Set))
		}
	}
	if v, ok := d.GetOk("cluster_config.0.gce_cluster_config.0.node_group_affinity"); ok {
		cfgNga := v.([]interface{})[0].(map[string]interface{})
		conf.NodeGroupAffinity = &dataproc.NodeGroupAffinity{}
		if v, ok := cfgNga["node_group_uri"]; ok {
			conf.NodeGroupAffinity.NodeGroupUri = v.(string)
		}
	}
	return conf, nil
}

func expandSecurityConfig(cfg map[string]interface{}) *dataproc.SecurityConfig {
	conf := &dataproc.SecurityConfig{}
	if kfg, ok := cfg["kerberos_config"]; ok {
		conf.KerberosConfig = expandKerberosConfig(kfg.([]interface{})[0].(map[string]interface{}))
	}
	return conf
}

func expandKerberosConfig(cfg map[string]interface{}) *dataproc.KerberosConfig {
	conf := &dataproc.KerberosConfig{}
	if v, ok := cfg["enable_kerberos"]; ok {
		conf.EnableKerberos = v.(bool)
	}
	if v, ok := cfg["root_principal_password_uri"]; ok {
		conf.RootPrincipalPasswordUri = v.(string)
	}
	if v, ok := cfg["kms_key_uri"]; ok {
		conf.KmsKeyUri = v.(string)
	}
	if v, ok := cfg["keystore_uri"]; ok {
		conf.KeystoreUri = v.(string)
	}
	if v, ok := cfg["truststore_uri"]; ok {
		conf.TruststoreUri = v.(string)
	}
	if v, ok := cfg["keystore_password_uri"]; ok {
		conf.KeystorePasswordUri = v.(string)
	}
	if v, ok := cfg["key_password_uri"]; ok {
		conf.KeyPasswordUri = v.(string)
	}
	if v, ok := cfg["truststore_password_uri"]; ok {
		conf.TruststorePasswordUri = v.(string)
	}
	if v, ok := cfg["cross_realm_trust_realm"]; ok {
		conf.CrossRealmTrustRealm = v.(string)
	}
	if v, ok := cfg["cross_realm_trust_kdc"]; ok {
		conf.CrossRealmTrustKdc = v.(string)
	}
	if v, ok := cfg["cross_realm_trust_admin_server"]; ok {
		conf.CrossRealmTrustAdminServer = v.(string)
	}
	if v, ok := cfg["cross_realm_trust_shared_password_uri"]; ok {
		conf.CrossRealmTrustSharedPasswordUri = v.(string)
	}
	if v, ok := cfg["kdc_db_key_uri"]; ok {
		conf.KdcDbKeyUri = v.(string)
	}
	if v, ok := cfg["tgt_lifetime_hours"]; ok {
		conf.TgtLifetimeHours = int64(v.(int))
	}
	if v, ok := cfg["realm"]; ok {
		conf.Realm = v.(string)
	}

	return conf
}

func expandSoftwareConfig(cfg map[string]interface{}) *dataproc.SoftwareConfig {
	conf := &dataproc.SoftwareConfig{}
	if v, ok := cfg["override_properties"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		conf.Properties = m
	}
	if v, ok := cfg["image_version"]; ok {
		conf.ImageVersion = v.(string)
	}
	if components, ok := cfg["optional_components"]; ok {
		compSet := components.(*schema.Set)
		components := make([]string, compSet.Len())
		for i, component := range compSet.List() {
			components[i] = component.(string)
		}
		conf.OptionalComponents = components
	}
	return conf
}

func expandEncryptionConfig(cfg map[string]interface{}) *dataproc.EncryptionConfig {
	conf := &dataproc.EncryptionConfig{}
	if v, ok := cfg["kms_key_name"]; ok {
		conf.GcePdKmsKeyName = v.(string)
	}
	return conf
}

func expandAutoscalingConfig(cfg map[string]interface{}) *dataproc.AutoscalingConfig {
	conf := &dataproc.AutoscalingConfig{}
	if v, ok := cfg["policy_uri"]; ok {
		conf.PolicyUri = v.(string)
	}
	return conf
}

func expandLifecycleConfig(cfg map[string]interface{}) *dataproc.LifecycleConfig {
	conf := &dataproc.LifecycleConfig{}
	if v, ok := cfg["idle_delete_ttl"]; ok {
		conf.IdleDeleteTtl = v.(string)
	}
	if v, ok := cfg["auto_delete_time"]; ok {
		conf.AutoDeleteTime = v.(string)
	}
	return conf
}

func expandEndpointConfig(cfg map[string]interface{}) *dataproc.EndpointConfig {
	conf := &dataproc.EndpointConfig{}
	if v, ok := cfg["enable_http_port_access"]; ok {
		conf.EnableHttpPortAccess = v.(bool)
	}
	return conf
}

func expandDataprocMetricConfig(cfg map[string]interface{}) *dataproc.DataprocMetricConfig {
	conf := &dataproc.DataprocMetricConfig{}
	metricsConfigs := cfg["metrics"].([]interface{})
	metricsSet := make([]*dataproc.Metric, 0, len(metricsConfigs))

	for _, raw := range metricsConfigs {
		data := raw.(map[string]interface{})
		metric := dataproc.Metric{
			MetricSource:    data["metric_source"].(string),
			MetricOverrides: tpgresource.ConvertStringSet(data["metric_overrides"].(*schema.Set)),
		}
		metricsSet = append(metricsSet, &metric)
	}
	conf.Metrics = metricsSet
	return conf
}

func expandMetastoreConfig(cfg map[string]interface{}) *dataproc.MetastoreConfig {
	conf := &dataproc.MetastoreConfig{}
	if v, ok := cfg["dataproc_metastore_service"]; ok {
		conf.DataprocMetastoreService = v.(string)
	}
	return conf
}

func expandInitializationActions(v interface{}) []*dataproc.NodeInitializationAction {
	actionList := v.([]interface{})

	actions := []*dataproc.NodeInitializationAction{}
	for _, v1 := range actionList {
		actionItem := v1.(map[string]interface{})
		action := &dataproc.NodeInitializationAction{
			ExecutableFile: actionItem["script"].(string),
		}
		if x, ok := actionItem["timeout_sec"]; ok {
			action.ExecutionTimeout = strconv.Itoa(x.(int)) + "s"
		}
		actions = append(actions, action)
	}

	return actions
}

func expandPreemptibleInstanceGroupConfig(cfg map[string]interface{}) *dataproc.InstanceGroupConfig {
	icg := &dataproc.InstanceGroupConfig{}

	if v, ok := cfg["num_instances"]; ok {
		icg.NumInstances = int64(v.(int))
	}
	if dc, ok := cfg["disk_config"]; ok {
		d := dc.([]interface{})
		if len(d) > 0 {
			dcfg := d[0].(map[string]interface{})
			icg.DiskConfig = &dataproc.DiskConfig{}

			if v, ok := dcfg["boot_disk_size_gb"]; ok {
				icg.DiskConfig.BootDiskSizeGb = int64(v.(int))
			}
			if v, ok := dcfg["num_local_ssds"]; ok {
				icg.DiskConfig.NumLocalSsds = int64(v.(int))
			}
			if v, ok := dcfg["boot_disk_type"]; ok {
				icg.DiskConfig.BootDiskType = v.(string)
			}
		}
	}
	if p, ok := cfg["preemptibility"]; ok {
		icg.Preemptibility = p.(string)
	}
	return icg
}

func expandInstanceGroupConfig(cfg map[string]interface{}) *dataproc.InstanceGroupConfig {
	icg := &dataproc.InstanceGroupConfig{}

	if v, ok := cfg["num_instances"]; ok {
		icg.NumInstances = int64(v.(int))
	}
	if v, ok := cfg["machine_type"]; ok {
		icg.MachineTypeUri = tpgresource.GetResourceNameFromSelfLink(v.(string))
	}
	if v, ok := cfg["min_cpu_platform"]; ok {
		icg.MinCpuPlatform = v.(string)
	}
	if v, ok := cfg["image_uri"]; ok {
		icg.ImageUri = v.(string)
	}

	if dc, ok := cfg["disk_config"]; ok {
		d := dc.([]interface{})
		if len(d) > 0 {
			dcfg := d[0].(map[string]interface{})
			icg.DiskConfig = &dataproc.DiskConfig{}

			if v, ok := dcfg["boot_disk_size_gb"]; ok {
				icg.DiskConfig.BootDiskSizeGb = int64(v.(int))
			}
			if v, ok := dcfg["num_local_ssds"]; ok {
				icg.DiskConfig.NumLocalSsds = int64(v.(int))
			}
			if v, ok := dcfg["boot_disk_type"]; ok {
				icg.DiskConfig.BootDiskType = v.(string)
			}
		}
	}

	icg.Accelerators = expandAccelerators(cfg["accelerators"].(*schema.Set).List())
	return icg
}

func expandAccelerators(configured []interface{}) []*dataproc.AcceleratorConfig {
	accelerators := make([]*dataproc.AcceleratorConfig, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		accelerator := dataproc.AcceleratorConfig{
			AcceleratorTypeUri: data["accelerator_type"].(string),
			AcceleratorCount:   int64(data["accelerator_count"].(int)),
		}

		accelerators = append(accelerators, &accelerator)
	}

	return accelerators
}

func resourceDataprocClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)

	cluster := &dataproc.Cluster{
		ClusterName: clusterName,
		ProjectId:   project,
		Config:      &dataproc.ClusterConfig{},
	}

	updMask := []string{}

	if d.HasChange("labels") {
		v := d.Get("labels")
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cluster.Labels = m

		updMask = append(updMask, "labels")
	}

	if d.HasChange("cluster_config.0.worker_config.0.num_instances") {
		desiredNumWorks := d.Get("cluster_config.0.worker_config.0.num_instances").(int)
		cluster.Config.WorkerConfig = &dataproc.InstanceGroupConfig{
			NumInstances: int64(desiredNumWorks),
		}

		updMask = append(updMask, "config.worker_config.num_instances")
	}

	if d.HasChange("cluster_config.0.preemptible_worker_config.0.num_instances") {
		desiredNumWorks := d.Get("cluster_config.0.preemptible_worker_config.0.num_instances").(int)
		cluster.Config.SecondaryWorkerConfig = &dataproc.InstanceGroupConfig{
			NumInstances: int64(desiredNumWorks),
		}

		updMask = append(updMask, "config.secondary_worker_config.num_instances")
	}

	if d.HasChange("cluster_config.0.autoscaling_config") {
		desiredPolicy := d.Get("cluster_config.0.autoscaling_config.0.policy_uri").(string)
		cluster.Config.AutoscalingConfig = &dataproc.AutoscalingConfig{
			PolicyUri: desiredPolicy,
		}

		updMask = append(updMask, "config.autoscaling_config.policy_uri")
	}

	if d.HasChange("cluster_config.0.lifecycle_config.0.idle_delete_ttl") {
		idleDeleteTtl := d.Get("cluster_config.0.lifecycle_config.0.idle_delete_ttl").(string)
		cluster.Config.LifecycleConfig = &dataproc.LifecycleConfig{
			IdleDeleteTtl: idleDeleteTtl,
		}

		updMask = append(updMask, "config.lifecycle_config.idle_delete_ttl")
	}

	if d.HasChange("cluster_config.0.lifecycle_config.0.auto_delete_time") {
		desiredDeleteTime := d.Get("cluster_config.0.lifecycle_config.0.auto_delete_time").(string)
		if cluster.Config.LifecycleConfig != nil {
			cluster.Config.LifecycleConfig.AutoDeleteTime = desiredDeleteTime
		} else {
			cluster.Config.LifecycleConfig = &dataproc.LifecycleConfig{
				AutoDeleteTime: desiredDeleteTime,
			}
		}

		updMask = append(updMask, "config.lifecycle_config.auto_delete_time")
	}

	if len(updMask) > 0 {
		gracefulDecommissionTimeout := d.Get("graceful_decommission_timeout").(string)

		patch := config.NewDataprocClient(userAgent).Projects.Regions.Clusters.Patch(
			project, region, clusterName, cluster)
		patch.GracefulDecommissionTimeout(gracefulDecommissionTimeout)
		patch.UpdateMask(strings.Join(updMask, ","))
		op, err := patch.Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := DataprocClusterOperationWait(config, op, "updating Dataproc cluster ", userAgent, d.Timeout(schema.TimeoutUpdate))
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Dataproc cluster %s has been updated ", d.Id())
	}

	return resourceDataprocClusterRead(d, meta)
}

func resourceDataprocClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)

	cluster, err := config.NewDataprocClient(userAgent).Projects.Regions.Clusters.Get(
		project, region, clusterName).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Dataproc Cluster %q", clusterName))
	}

	if err := d.Set("name", cluster.ClusterName); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("labels", cluster.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}

	var cfg []map[string]interface{}
	cfg, err = flattenClusterConfig(d, cluster.Config)

	if err != nil {
		return err
	}

	err = d.Set("cluster_config", cfg)
	virtualCfg, err := flattenVirtualClusterConfig(d, cluster.VirtualClusterConfig)

	if err != nil {
		return err
	}

	err = d.Set("virtual_cluster_config", virtualCfg)

	if err != nil {
		return err
	}

	return nil
}

func flattenVirtualClusterConfig(d *schema.ResourceData, cfg *dataproc.VirtualClusterConfig) ([]map[string]interface{}, error) {
	if cfg == nil {
		return []map[string]interface{}{}, nil
	}

	data := map[string]interface{}{
		"staging_bucket":            d.Get("virtual_cluster_config.0.staging_bucket"),
		"auxiliary_services_config": flattenAuxiliaryServicesConfig(d, cfg.AuxiliaryServicesConfig),
		"kubernetes_cluster_config": flattenKubernetesClusterConfig(d, cfg.KubernetesClusterConfig),
	}

	return []map[string]interface{}{data}, nil
}

func flattenAuxiliaryServicesConfig(d *schema.ResourceData, cfg *dataproc.AuxiliaryServicesConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	data := map[string]interface{}{
		"metastore_config":            flattenVCMetastoreConfig(d, cfg.MetastoreConfig),
		"spark_history_server_config": flattenSparkHistoryServerConfig(d, cfg.SparkHistoryServerConfig),
	}

	return []map[string]interface{}{data}
}

func flattenVCMetastoreConfig(d *schema.ResourceData, cfg *dataproc.MetastoreConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	data := map[string]interface{}{
		"dataproc_metastore_service": cfg.DataprocMetastoreService,
	}

	return []map[string]interface{}{data}
}

func flattenSparkHistoryServerConfig(d *schema.ResourceData, cfg *dataproc.SparkHistoryServerConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	data := map[string]interface{}{
		"dataproc_cluster": cfg.DataprocCluster,
	}

	return []map[string]interface{}{data}
}

func flattenKubernetesClusterConfig(d *schema.ResourceData, cfg *dataproc.KubernetesClusterConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	data := map[string]interface{}{
		"gke_cluster_config":         flattenGkeClusterConfig(d, cfg.GkeClusterConfig),
		"kubernetes_namespace":       cfg.KubernetesNamespace,
		"kubernetes_software_config": flattenKubernetesSoftwareConfig(d, cfg.KubernetesSoftwareConfig),
	}

	return []map[string]interface{}{data}
}

func flattenGkeClusterConfig(d *schema.ResourceData, cfg *dataproc.GkeClusterConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	data := map[string]interface{}{
		"gke_cluster_target": cfg.GkeClusterTarget,
	}

	if len(cfg.NodePoolTarget) > 0 {
		val, err := flattenNodePoolTargetConfig(d, cfg.NodePoolTarget, cfg.GkeClusterTarget)
		if err != nil {
			return nil
		}

		data["node_pool_target"] = val
	}

	return []map[string]interface{}{data}
}

func flattenNodePoolTargetConfig(d *schema.ResourceData, nia []*dataproc.GkeNodePoolTarget, clusterAddress string) ([]map[string]interface{}, error) {
	nodePools := []map[string]interface{}{}
	for i, v := range nia {
		nodePoolAddress := strings.Split(v.NodePool, "/")
		nodePool := map[string]interface{}{
			"node_pool": nodePoolAddress[len(nodePoolAddress)-1],
			"roles":     v.Roles,
		}

		if npc, ok := d.GetOk(fmt.Sprintf("virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target.%d.node_pool_config", i)); ok {
			// We can't read the node pool config details due to `Input only`-field,
			// copy the initialize_params from what the user originally specified to avoid diffs.
			nodePool["node_pool_config"] = npc
		}

		nodePools = append(nodePools, nodePool)
	}

	return nodePools, nil
}

func flattenKubernetesSoftwareConfig(d *schema.ResourceData, cfg *dataproc.KubernetesSoftwareConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	data := map[string]interface{}{
		"component_version": cfg.ComponentVersion,
		"properties":        cfg.Properties,
	}

	return []map[string]interface{}{data}
}

func flattenClusterConfig(d *schema.ResourceData, cfg *dataproc.ClusterConfig) ([]map[string]interface{}, error) {
	if cfg == nil {
		return []map[string]interface{}{}, nil
	}

	data := map[string]interface{}{
		"staging_bucket": d.Get("cluster_config.0.staging_bucket").(string),

		"bucket":                    cfg.ConfigBucket,
		"temp_bucket":               cfg.TempBucket,
		"gce_cluster_config":        flattenGceClusterConfig(d, cfg.GceClusterConfig),
		"master_config":             flattenInstanceGroupConfig(d, cfg.MasterConfig),
		"worker_config":             flattenInstanceGroupConfig(d, cfg.WorkerConfig),
		"software_config":           flattenSoftwareConfig(d, cfg.SoftwareConfig),
		"encryption_config":         flattenEncryptionConfig(d, cfg.EncryptionConfig),
		"autoscaling_config":        flattenAutoscalingConfig(d, cfg.AutoscalingConfig),
		"security_config":           flattenSecurityConfig(d, cfg.SecurityConfig),
		"preemptible_worker_config": flattenPreemptibleInstanceGroupConfig(d, cfg.SecondaryWorkerConfig),
		"metastore_config":          flattenMetastoreConfig(d, cfg.MetastoreConfig),
		"lifecycle_config":          flattenLifecycleConfig(d, cfg.LifecycleConfig),
		"endpoint_config":           flattenEndpointConfig(d, cfg.EndpointConfig),
		"dataproc_metric_config":    flattenDataprocMetricConfig(d, cfg.DataprocMetricConfig),
	}

	if len(cfg.InitializationActions) > 0 {
		val, err := flattenInitializationActions(cfg.InitializationActions)
		if err != nil {
			return nil, err
		}
		data["initialization_action"] = val
	}
	return []map[string]interface{}{data}, nil
}

func flattenSecurityConfig(d *schema.ResourceData, sc *dataproc.SecurityConfig) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	data := map[string]interface{}{
		"kerberos_config": flattenKerberosConfig(d, sc.KerberosConfig),
	}

	return []map[string]interface{}{data}
}

func flattenKerberosConfig(d *schema.ResourceData, kfg *dataproc.KerberosConfig) []map[string]interface{} {
	data := map[string]interface{}{
		"enable_kerberos":                       kfg.EnableKerberos,
		"root_principal_password_uri":           kfg.RootPrincipalPasswordUri,
		"kms_key_uri":                           kfg.KmsKeyUri,
		"keystore_uri":                          kfg.KeystoreUri,
		"truststore_uri":                        kfg.TruststoreUri,
		"keystore_password_uri":                 kfg.KeystorePasswordUri,
		"key_password_uri":                      kfg.KeyPasswordUri,
		"truststore_password_uri":               kfg.TruststorePasswordUri,
		"cross_realm_trust_realm":               kfg.CrossRealmTrustRealm,
		"cross_realm_trust_kdc":                 kfg.CrossRealmTrustKdc,
		"cross_realm_trust_admin_server":        kfg.CrossRealmTrustAdminServer,
		"cross_realm_trust_shared_password_uri": kfg.CrossRealmTrustSharedPasswordUri,
		"kdc_db_key_uri":                        kfg.KdcDbKeyUri,
		"tgt_lifetime_hours":                    kfg.TgtLifetimeHours,
		"realm":                                 kfg.Realm,
	}

	return []map[string]interface{}{data}
}

func flattenSoftwareConfig(d *schema.ResourceData, sc *dataproc.SoftwareConfig) []map[string]interface{} {
	data := map[string]interface{}{
		"image_version":       sc.ImageVersion,
		"optional_components": sc.OptionalComponents,
		"properties":          sc.Properties,
		"override_properties": d.Get("cluster_config.0.software_config.0.override_properties").(map[string]interface{}),
	}

	return []map[string]interface{}{data}
}

func flattenEncryptionConfig(d *schema.ResourceData, ec *dataproc.EncryptionConfig) []map[string]interface{} {
	if ec == nil {
		return nil
	}

	data := map[string]interface{}{
		"kms_key_name": ec.GcePdKmsKeyName,
	}

	return []map[string]interface{}{data}
}

func flattenAutoscalingConfig(d *schema.ResourceData, ec *dataproc.AutoscalingConfig) []map[string]interface{} {
	if ec == nil {
		return nil
	}

	data := map[string]interface{}{
		"policy_uri": ec.PolicyUri,
	}

	return []map[string]interface{}{data}
}

func flattenLifecycleConfig(d *schema.ResourceData, lc *dataproc.LifecycleConfig) []map[string]interface{} {
	if lc == nil {
		return nil
	}

	data := map[string]interface{}{
		"idle_delete_ttl":  lc.IdleDeleteTtl,
		"auto_delete_time": lc.AutoDeleteTime,
	}

	return []map[string]interface{}{data}
}

func flattenEndpointConfig(d *schema.ResourceData, ec *dataproc.EndpointConfig) []map[string]interface{} {
	if ec == nil {
		return nil
	}

	data := map[string]interface{}{
		"enable_http_port_access": ec.EnableHttpPortAccess,
		"http_ports":              ec.HttpPorts,
	}

	return []map[string]interface{}{data}
}

func flattenDataprocMetricConfig(d *schema.ResourceData, dmc *dataproc.DataprocMetricConfig) []map[string]interface{} {
	if dmc == nil {
		return nil
	}

	metrics := map[string]interface{}{}
	metricsTypeList := schema.NewSet(schema.HashResource(metricsSchema()), []interface{}{}).List()
	for _, metric := range dmc.Metrics {
		data := map[string]interface{}{
			"metric_source":    metric.MetricSource,
			"metric_overrides": metric.MetricOverrides,
		}

		metricsTypeList = append(metricsTypeList, &data)
	}
	metrics["metrics"] = metricsTypeList

	return []map[string]interface{}{metrics}
}

func flattenMetastoreConfig(d *schema.ResourceData, ec *dataproc.MetastoreConfig) []map[string]interface{} {
	if ec == nil {
		return nil
	}

	data := map[string]interface{}{
		"dataproc_metastore_service": ec.DataprocMetastoreService,
	}

	return []map[string]interface{}{data}
}

func flattenAccelerators(accelerators []*dataproc.AcceleratorConfig) interface{} {
	acceleratorsTypeSet := schema.NewSet(schema.HashResource(acceleratorsSchema()), []interface{}{})
	for _, accelerator := range accelerators {
		data := map[string]interface{}{
			"accelerator_type":  tpgresource.GetResourceNameFromSelfLink(accelerator.AcceleratorTypeUri),
			"accelerator_count": int(accelerator.AcceleratorCount),
		}

		acceleratorsTypeSet.Add(data)
	}

	return acceleratorsTypeSet
}

func flattenInitializationActions(nia []*dataproc.NodeInitializationAction) ([]map[string]interface{}, error) {

	actions := []map[string]interface{}{}
	for _, v := range nia {
		action := map[string]interface{}{
			"script": v.ExecutableFile,
		}
		if len(v.ExecutionTimeout) > 0 {
			tsec, err := extractInitTimeout(v.ExecutionTimeout)
			if err != nil {
				return nil, err
			}
			action["timeout_sec"] = tsec
		}

		actions = append(actions, action)
	}
	return actions, nil

}

func flattenGceClusterConfig(d *schema.ResourceData, gcc *dataproc.GceClusterConfig) []map[string]interface{} {
	if gcc == nil {
		return []map[string]interface{}{}
	}

	gceConfig := map[string]interface{}{
		"tags":             schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(gcc.Tags)),
		"service_account":  gcc.ServiceAccount,
		"zone":             tpgresource.GetResourceNameFromSelfLink(gcc.ZoneUri),
		"internal_ip_only": gcc.InternalIpOnly,
		"metadata":         gcc.Metadata,
	}

	if gcc.NetworkUri != "" {
		gceConfig["network"] = gcc.NetworkUri
	}
	if gcc.SubnetworkUri != "" {
		gceConfig["subnetwork"] = gcc.SubnetworkUri
	}
	if len(gcc.ServiceAccountScopes) > 0 {
		gceConfig["service_account_scopes"] = schema.NewSet(tpgresource.StringScopeHashcode, tpgresource.ConvertStringArrToInterface(gcc.ServiceAccountScopes))
	}
	if gcc.ShieldedInstanceConfig != nil {
		gceConfig["shielded_instance_config"] = []map[string]interface{}{
			{
				"enable_integrity_monitoring": gcc.ShieldedInstanceConfig.EnableIntegrityMonitoring,
				"enable_secure_boot":          gcc.ShieldedInstanceConfig.EnableSecureBoot,
				"enable_vtpm":                 gcc.ShieldedInstanceConfig.EnableVtpm,
			},
		}
	}
	if gcc.ReservationAffinity != nil {
		gceConfig["reservation_affinity"] = []map[string]interface{}{
			{
				"consume_reservation_type": gcc.ReservationAffinity.ConsumeReservationType,
				"key":                      gcc.ReservationAffinity.Key,
				"values":                   gcc.ReservationAffinity.Values,
			},
		}
	}
	if gcc.NodeGroupAffinity != nil {
		gceConfig["node_group_affinity"] = []map[string]interface{}{
			{
				"node_group_uri": gcc.NodeGroupAffinity.NodeGroupUri,
			},
		}
	}

	return []map[string]interface{}{gceConfig}
}

func flattenPreemptibleInstanceGroupConfig(d *schema.ResourceData, icg *dataproc.InstanceGroupConfig) []map[string]interface{} {
	// if num_instances is 0, icg will always be returned nil. This means the
	// server has discarded diskconfig etc. However, the only way to remove the
	// preemptible group is to set the size to 0, because it's O+C. Many users
	// won't remove the rest of the config (eg disk config). Therefore, we need to
	// preserve the other set fields by using the old state to stop users from
	// getting a diff.
	if icg == nil {
		icgSchema := d.Get("cluster_config.0.preemptible_worker_config")
		log.Printf("[DEBUG] state of preemptible is %#v", icgSchema)
		if v, ok := icgSchema.([]interface{}); ok && len(v) > 0 {
			if m, ok := v[0].(map[string]interface{}); ok {
				return []map[string]interface{}{m}
			}
		}
	}

	disk := map[string]interface{}{}
	data := map[string]interface{}{}

	if icg != nil {
		data["num_instances"] = icg.NumInstances
		data["instance_names"] = icg.InstanceNames
		data["preemptibility"] = icg.Preemptibility
		if icg.DiskConfig != nil {
			disk["boot_disk_size_gb"] = icg.DiskConfig.BootDiskSizeGb
			disk["num_local_ssds"] = icg.DiskConfig.NumLocalSsds
			disk["boot_disk_type"] = icg.DiskConfig.BootDiskType
		}
	}

	data["disk_config"] = []map[string]interface{}{disk}
	return []map[string]interface{}{data}
}

func flattenInstanceGroupConfig(d *schema.ResourceData, icg *dataproc.InstanceGroupConfig) []map[string]interface{} {
	disk := map[string]interface{}{}
	data := map[string]interface{}{}

	if icg != nil {
		data["num_instances"] = icg.NumInstances
		data["machine_type"] = tpgresource.GetResourceNameFromSelfLink(icg.MachineTypeUri)
		data["min_cpu_platform"] = icg.MinCpuPlatform
		data["image_uri"] = icg.ImageUri
		data["instance_names"] = icg.InstanceNames
		if icg.DiskConfig != nil {
			disk["boot_disk_size_gb"] = icg.DiskConfig.BootDiskSizeGb
			disk["num_local_ssds"] = icg.DiskConfig.NumLocalSsds
			disk["boot_disk_type"] = icg.DiskConfig.BootDiskType
		}

		data["accelerators"] = flattenAccelerators(icg.Accelerators)
	}

	data["disk_config"] = []map[string]interface{}{disk}
	return []map[string]interface{}{data}
}

func extractInitTimeout(t string) (int, error) {
	d, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}
	return int(d.Seconds()), nil
}

func resourceDataprocClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)

	log.Printf("[DEBUG] Deleting Dataproc cluster %s", clusterName)
	op, err := config.NewDataprocClient(userAgent).Projects.Regions.Clusters.Delete(
		project, region, clusterName).Do()
	if err != nil {
		return err
	}

	// Wait until it's deleted
	waitErr := DataprocClusterOperationWait(config, op, "deleting Dataproc cluster", userAgent, d.Timeout(schema.TimeoutDelete))
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been deleted", d.Id())
	d.SetId("")

	return nil
}

func configOptions(d *schema.ResourceData, option string) (map[string]interface{}, bool) {
	if v, ok := d.GetOk(option); ok {
		clist := v.([]interface{})
		if len(clist) == 0 {
			return nil, false
		}

		if clist[0] != nil {
			return clist[0].(map[string]interface{}), true
		}
	}
	return nil, false
}

func dataprocImageVersionDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldV, err := parseDataprocImageVersion(old)
	if err != nil {
		return false
	}
	newV, err := parseDataprocImageVersion(new)
	if err != nil {
		return false
	}

	if newV.major != oldV.major {
		return false
	}
	if newV.minor != oldV.minor {
		return false
	}
	// Only compare subminor version if set in config version.
	if newV.subminor != "" && newV.subminor != oldV.subminor {
		return false
	}
	// Only compare os if it is set in config version.
	if newV.osName != "" && newV.osName != oldV.osName {
		return false
	}
	return true
}

type dataprocImageVersion struct {
	major    string
	minor    string
	subminor string
	osName   string
}

func parseDataprocImageVersion(version string) (*dataprocImageVersion, error) {
	matches := resolveDataprocImageVersion.FindStringSubmatch(version)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid image version %q", version)
	}

	return &dataprocImageVersion{
		major:    matches[1],
		minor:    matches[2],
		subminor: matches[3],
		osName:   matches[4],
	}, nil
}
