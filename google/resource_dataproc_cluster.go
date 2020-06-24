package google

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	dataproc "google.golang.org/api/dataproc/v1beta2"
)

var (
	resolveDataprocImageVersion = regexp.MustCompile(`(?P<Major>[^\s.-]+)\.(?P<Minor>[^\s.-]+)(?:\.(?P<Subminor>[^\s.-]+))?(?:\-(?P<Distr>[^\s.-]+))?`)

	gceClusterConfigKeys = []string{
		"cluster_config.0.gce_cluster_config.0.zone",
		"cluster_config.0.gce_cluster_config.0.network",
		"cluster_config.0.gce_cluster_config.0.subnetwork",
		"cluster_config.0.gce_cluster_config.0.tags",
		"cluster_config.0.gce_cluster_config.0.service_account",
		"cluster_config.0.gce_cluster_config.0.service_account_scopes",
		"cluster_config.0.gce_cluster_config.0.internal_ip_only",
		"cluster_config.0.gce_cluster_config.0.metadata",
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

	clusterConfigKeys = []string{
		"cluster_config.0.staging_bucket",
		"cluster_config.0.gce_cluster_config",
		"cluster_config.0.master_config",
		"cluster_config.0.worker_config",
		"cluster_config.0.preemptible_worker_config",
		"cluster_config.0.security_config",
		"cluster_config.0.software_config",
		"cluster_config.0.initialization_action",
		"cluster_config.0.encryption_config",
		"cluster_config.0.autoscaling_config",
	}
)

func resourceDataprocCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocClusterCreate,
		Read:   resourceDataprocClusterRead,
		Update: resourceDataprocClusterUpdate,
		Delete: resourceDataprocClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// GCP automatically adds two labels
				//    'goog-dataproc-cluster-uuid'
				//    'goog-dataproc-cluster-name'
				Computed:    true,
				Description: `The list of labels (key/value pairs) to be applied to instances in the cluster. GCP generates some itself including goog-dataproc-cluster-name which is the name of the cluster.`,
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
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The name or self_link of the Google Compute Engine network to the cluster will be part of. Conflicts with subnetwork. If neither is specified, this defaults to the "default" network.`,
									},

									"subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										AtLeastOneOf:     gceClusterConfigKeys,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.network"},
										DiffSuppressFunc: compareSelfLinkOrResourceName,
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
												return canonicalizeServiceScope(v.(string))
											},
										},
										Set: stringScopeHashcode,
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
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
									},

									// API does not honour this if set ...
									// It always uses whatever is specified for the worker_config
									// "machine_type": { ... }
									// "min_cpu_platform": { ... }
									"disk_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: `Disk Config`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
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
													ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd", ""}, false),
													Default:      "pd-standard",
													Description:  `The disk type of the primary disk attached to each preemptible worker node. One of "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
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
											ValidateFunc: validation.StringInSlice([]string{"COMPONENT_UNSPECIFIED", "ANACONDA", "DRUID", "HBASE", "HIVE_WEBHCAT",
												"JUPYTER", "KERBEROS", "PRESTO", "RANGER", "SOLR", "ZEPPELIN", "ZOOKEEPER"}, false),
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
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							MaxItems:     1,
							Description:  `The autoscaling policy config associated with the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_uri": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The autoscaling policy used by the cluster.`,
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

func instanceConfigSchema(parent string) *schema.Schema {
	var instanceConfigKeys = []string{
		"cluster_config.0." + parent + ".0.num_instances",
		"cluster_config.0." + parent + ".0.image_uri",
		"cluster_config.0." + parent + ".0.machine_type",
		"cluster_config.0." + parent + ".0.min_cpu_platform",
		"cluster_config.0." + parent + ".0.disk_config",
		"cluster_config.0." + parent + ".0.accelerators",
	}

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
								Description: `The disk type of the primary disk attached to each node. One of "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
								AtLeastOneOf: []string{
									"cluster_config.0." + parent + ".0.disk_config.0.num_local_ssds",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_size_gb",
									"cluster_config.0." + parent + ".0.disk_config.0.boot_disk_type",
								},
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd", ""}, false),
								Default:      "pd-standard",
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
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	cluster := &dataproc.Cluster{
		ClusterName: d.Get("name").(string),
		ProjectId:   project,
	}

	cluster.Config, err = expandClusterConfig(d, config)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("labels"); ok {
		cluster.Labels = expandLabels(d)
	}

	// Checking here caters for the case where the user does not specify cluster_config
	// at all, as well where it is simply missing from the gce_cluster_config
	if region == "global" && cluster.Config.GceClusterConfig.ZoneUri == "" {
		return errors.New("zone is mandatory when region is set to 'global'")
	}

	// Create the cluster
	op, err := config.clientDataprocBeta.Projects.Regions.Clusters.Create(
		project, region, cluster).Do()
	if err != nil {
		return fmt.Errorf("Error creating Dataproc cluster: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/clusters/%s", project, region, cluster.ClusterName))

	// Wait until it's created
	waitErr := dataprocClusterOperationWait(config, op, "creating Dataproc cluster", d.Timeout(schema.TimeoutCreate))
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

func expandClusterConfig(d *schema.ResourceData, config *Config) (*dataproc.ClusterConfig, error) {
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
		if conf.SecondaryWorkerConfig.NumInstances > 0 {
			conf.SecondaryWorkerConfig.IsPreemptible = true
		}
	}
	return conf, nil
}

func expandGceClusterConfig(d *schema.ResourceData, config *Config) (*dataproc.GceClusterConfig, error) {
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
		nf, err := ParseNetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", v, err)
		}

		conf.NetworkUri = nf.RelativeLink()
	}
	if v, ok := cfg["subnetwork"]; ok {
		snf, err := ParseSubnetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", v, err)
		}

		conf.SubnetworkUri = snf.RelativeLink()
	}
	if v, ok := cfg["tags"]; ok {
		conf.Tags = convertStringSet(v.(*schema.Set))
	}
	if v, ok := cfg["service_account"]; ok {
		conf.ServiceAccount = v.(string)
	}
	if scopes, ok := cfg["service_account_scopes"]; ok {
		scopesSet := scopes.(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, scope := range scopesSet.List() {
			scopes[i] = canonicalizeServiceScope(scope.(string))
		}
		conf.ServiceAccountScopes = scopes
	}
	if v, ok := cfg["internal_ip_only"]; ok {
		conf.InternalIpOnly = v.(bool)
	}
	if v, ok := cfg["metadata"]; ok {
		conf.Metadata = convertStringMap(v.(map[string]interface{}))
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
	return icg
}

func expandInstanceGroupConfig(cfg map[string]interface{}) *dataproc.InstanceGroupConfig {
	icg := &dataproc.InstanceGroupConfig{}

	if v, ok := cfg["num_instances"]; ok {
		icg.NumInstances = int64(v.(int))
	}
	if v, ok := cfg["machine_type"]; ok {
		icg.MachineTypeUri = GetResourceNameFromSelfLink(v.(string))
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
	config := meta.(*Config)

	project, err := getProject(d, config)
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

	if len(updMask) > 0 {
		patch := config.clientDataprocBeta.Projects.Regions.Clusters.Patch(
			project, region, clusterName, cluster)
		op, err := patch.UpdateMask(strings.Join(updMask, ",")).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := dataprocClusterOperationWait(config, op, "updating Dataproc cluster ", d.Timeout(schema.TimeoutUpdate))
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Dataproc cluster %s has been updated ", d.Id())
	}

	return resourceDataprocClusterRead(d, meta)
}

func resourceDataprocClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)

	cluster, err := config.clientDataprocBeta.Projects.Regions.Clusters.Get(
		project, region, clusterName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataproc Cluster %q", clusterName))
	}

	d.Set("name", cluster.ClusterName)
	d.Set("project", project)
	d.Set("region", region)
	d.Set("labels", cluster.Labels)

	cfg, err := flattenClusterConfig(d, cluster.Config)
	if err != nil {
		return err
	}

	err = d.Set("cluster_config", cfg)
	if err != nil {
		return err
	}
	return nil
}

func flattenClusterConfig(d *schema.ResourceData, cfg *dataproc.ClusterConfig) ([]map[string]interface{}, error) {

	data := map[string]interface{}{
		"staging_bucket": d.Get("cluster_config.0.staging_bucket").(string),

		"bucket":                    cfg.ConfigBucket,
		"gce_cluster_config":        flattenGceClusterConfig(d, cfg.GceClusterConfig),
		"security_config":           flattenSecurityConfig(d, cfg.SecurityConfig),
		"software_config":           flattenSoftwareConfig(d, cfg.SoftwareConfig),
		"master_config":             flattenInstanceGroupConfig(d, cfg.MasterConfig),
		"worker_config":             flattenInstanceGroupConfig(d, cfg.WorkerConfig),
		"preemptible_worker_config": flattenPreemptibleInstanceGroupConfig(d, cfg.SecondaryWorkerConfig),
		"encryption_config":         flattenEncryptionConfig(d, cfg.EncryptionConfig),
		"autoscaling_config":        flattenAutoscalingConfig(d, cfg.AutoscalingConfig),
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

func flattenAccelerators(accelerators []*dataproc.AcceleratorConfig) interface{} {
	acceleratorsTypeSet := schema.NewSet(schema.HashResource(acceleratorsSchema()), []interface{}{})
	for _, accelerator := range accelerators {
		data := map[string]interface{}{
			"accelerator_type":  GetResourceNameFromSelfLink(accelerator.AcceleratorTypeUri),
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

	gceConfig := map[string]interface{}{
		"tags":             schema.NewSet(schema.HashString, convertStringArrToInterface(gcc.Tags)),
		"service_account":  gcc.ServiceAccount,
		"zone":             GetResourceNameFromSelfLink(gcc.ZoneUri),
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
		gceConfig["service_account_scopes"] = schema.NewSet(stringScopeHashcode, convertStringArrToInterface(gcc.ServiceAccountScopes))
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
		data["machine_type"] = GetResourceNameFromSelfLink(icg.MachineTypeUri)
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
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)

	log.Printf("[DEBUG] Deleting Dataproc cluster %s", clusterName)
	op, err := config.clientDataprocBeta.Projects.Regions.Clusters.Delete(
		project, region, clusterName).Do()
	if err != nil {
		return err
	}

	// Wait until it's deleted
	waitErr := dataprocClusterOperationWait(config, op, "deleting Dataproc cluster", d.Timeout(schema.TimeoutDelete))
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
