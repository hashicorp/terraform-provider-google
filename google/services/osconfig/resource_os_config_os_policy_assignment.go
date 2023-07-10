// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package osconfig

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceOSConfigOSPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOSConfigOSPolicyAssignmentCreate,
		Read:   resourceOSConfigOSPolicyAssignmentRead,
		Update: resourceOSConfigOSPolicyAssignmentUpdate,
		Delete: resourceOSConfigOSPolicyAssignmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceOSConfigOSPolicyAssignmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_filter": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Filter to select VMs.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Target all VMs in the project. If true, no other criteria is permitted.`,
						},
						"exclusion_labels": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `List of label sets used for VM exclusion.
If the list has more than one label set, the VM is excluded if any of the label sets are applicable for the VM.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"labels": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: `Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"inclusion_labels": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `List of label sets used for VM inclusion.
If the list has more than one 'LabelSet', the VM is included if any of the label sets are applicable for the VM.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"labels": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: `Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"inventories": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `List of inventories to select VMs.
A VM is selected if its inventory data matches at least one of the following inventories.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"os_short_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The OS short name`,
									},
									"os_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of '7', specify the following value for this field '7.*' An empty string matches all OS versions.`,
									},
								},
							},
						},
					},
				},
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The location for the resource`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Resource name.`,
			},
			"os_policies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `List of OS policies to be applied to the VMs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The id of the OS policy with the following restrictions:
* Must contain only lowercase letters, numbers, and hyphens.
* Must start with a letter.
* Must be between 1-63 characters.
* Must end with a number or a letter.
* Must be unique within the assignment.`,
						},
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidateEnum([]string{"MODE_UNSPECIFIED", "VALIDATION", "ENFORCEMENT"}),
							Description:  `Policy mode Possible values: ["MODE_UNSPECIFIED", "VALIDATION", "ENFORCEMENT"]`,
						},
						"resource_groups": {
							Type:     schema.TypeList,
							Required: true,
							Description: `List of resource groups for the policy. For a particular VM, resource groups are evaluated in the order specified and the first resource group that is applicable is selected and the rest are ignored.
If none of the resource groups are applicable for a VM, the VM is considered to be non-compliant w.r.t this policy. This behavior can be toggled by the flag 'allow_no_resource_group_match'`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resources": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `List of resources configured for this resource group. The resources are executed in the exact order specified here.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Required: true,
													Description: `The id of the resource with the following restrictions:
* Must contain only lowercase letters, numbers, and hyphens.
* Must start with a letter.
* Must be between 1-63 characters.
* Must end with a number or a letter.
* Must be unique within the OS policy.`,
												},
												"exec": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: `Exec resource`,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"validate": {
																Type:        schema.TypeList,
																Required:    true,
																Description: `What to run to validate this resource is in the desired state. An exit code of 100 indicates "in desired state", and exit code of 101 indicates "not in desired state". Any other exit code indicates a failure running validate.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"interpreter": {
																			Type:         schema.TypeString,
																			Required:     true,
																			ValidateFunc: verify.ValidateEnum([]string{"INTERPRETER_UNSPECIFIED", "NONE", "SHELL", "POWERSHELL"}),
																			Description:  `The script interpreter to use. Possible values: ["INTERPRETER_UNSPECIFIED", "NONE", "SHELL", "POWERSHELL"]`,
																		},
																		"args": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `Optional arguments to pass to the source during execution.`,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"file": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `A remote or local file.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"allow_insecure": {
																						Type:     schema.TypeBool,
																						Optional: true,
																						Description: `Defaults to false. When false, files are subject to validations based on the file type:
Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.`,
																					},
																					"gcs": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A Cloud Storage object.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Bucket of the Cloud Storage object.`,
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Name of the Cloud Storage object.`,
																								},
																								"generation": {
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: `Generation number of the Cloud Storage object.`,
																								},
																							},
																						},
																					},
																					"local_path": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: `A local path within the VM to use.`,
																					},
																					"remote": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A generic remote file.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"uri": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `URI from which to fetch the object. It should contain both the protocol and path following the format '{protocol}://{location}'.`,
																								},
																								"sha256_checksum": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: `SHA256 checksum of the remote file.`,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"output_file_path": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.`,
																		},
																		"script": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `An inline script. The size of the script is limited to 1024 characters.`,
																		},
																	},
																},
															},
															"enforce": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `What to run to bring this resource into the desired state. An exit code of 100 indicates "success", any other exit code indicates a failure running enforce.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"interpreter": {
																			Type:         schema.TypeString,
																			Required:     true,
																			ValidateFunc: verify.ValidateEnum([]string{"INTERPRETER_UNSPECIFIED", "NONE", "SHELL", "POWERSHELL"}),
																			Description:  `The script interpreter to use. Possible values: ["INTERPRETER_UNSPECIFIED", "NONE", "SHELL", "POWERSHELL"]`,
																		},
																		"args": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `Optional arguments to pass to the source during execution.`,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"file": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `A remote or local file.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"allow_insecure": {
																						Type:        schema.TypeBool,
																						Optional:    true,
																						Description: `Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.`,
																					},
																					"gcs": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A Cloud Storage object.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Bucket of the Cloud Storage object.`,
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Name of the Cloud Storage object.`,
																								},
																								"generation": {
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: `Generation number of the Cloud Storage object.`,
																								},
																							},
																						},
																					},
																					"local_path": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: `A local path within the VM to use.`,
																					},
																					"remote": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A generic remote file.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"uri": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `URI from which to fetch the object. It should contain both the protocol and path following the format '{protocol}://{location}'.`,
																								},
																								"sha256_checksum": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: `SHA256 checksum of the remote file.`,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"output_file_path": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.`,
																		},
																		"script": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `An inline script. The size of the script is limited to 1024 characters.`,
																		},
																	},
																},
															},
														},
													},
												},
												"file": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: `File resource`,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Required:    true,
																Description: `The absolute path of the file within the VM.`,
															},
															"state": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: verify.ValidateEnum([]string{"DESIRED_STATE_UNSPECIFIED", "PRESENT", "ABSENT", "CONTENTS_MATCH"}),
																Description:  `Desired state of the file. Possible values: ["DESIRED_STATE_UNSPECIFIED", "PRESENT", "ABSENT", "CONTENTS_MATCH"]`,
															},
															"content": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: `A a file with this content. The size of the content is limited to 1024 characters.`,
															},
															"file": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A remote or local source.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"allow_insecure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: `Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.`,
																		},
																		"gcs": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `A Cloud Storage object.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: `Bucket of the Cloud Storage object.`,
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: `Name of the Cloud Storage object.`,
																					},
																					"generation": {
																						Type:        schema.TypeInt,
																						Optional:    true,
																						Description: `Generation number of the Cloud Storage object.`,
																					},
																				},
																			},
																		},
																		"local_path": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `A local path within the VM to use.`,
																		},
																		"remote": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `A generic remote file.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"uri": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: `URI from which to fetch the object. It should contain both the protocol and path following the format '{protocol}://{location}'.`,
																					},
																					"sha256_checksum": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: `SHA256 checksum of the remote file.`,
																					},
																				},
																			},
																		},
																	},
																},
															},
															"permissions": {
																Type:     schema.TypeString,
																Computed: true,
																Description: `Consists of three octal digits which represent, in order, the permissions of the owner, group, and other users for the file (similarly to the numeric mode used in the linux chmod utility). Each digit represents a three bit number with the 4 bit corresponding to the read permissions, the 2 bit corresponds to the write bit, and the one bit corresponds to the execute permission. Default behavior is 755.
Below are some examples of permissions and their associated values: read, write, and execute: 7 read and execute: 5 read and write: 6 read only: 4`,
															},
														},
													},
												},
												"pkg": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: `Package resource`,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"desired_state": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: verify.ValidateEnum([]string{"DESIRED_STATE_UNSPECIFIED", "INSTALLED", "REMOVED"}),
																Description:  `The desired state the agent should maintain for this package. Possible values: ["DESIRED_STATE_UNSPECIFIED", "INSTALLED", "REMOVED"]`,
															},
															"apt": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A package managed by Apt.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `Package name.`,
																		},
																	},
																},
															},
															"deb": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A deb package file.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"source": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: `A deb package.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"allow_insecure": {
																						Type:     schema.TypeBool,
																						Optional: true,
																						Description: `Defaults to false. When false, files are subject to validations based on the file type:
Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.`,
																					},
																					"gcs": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A Cloud Storage object.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Bucket of the Cloud Storage object.`,
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Name of the Cloud Storage object.`,
																								},
																								"generation": {
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: `Generation number of the Cloud Storage object.`,
																								},
																							},
																						},
																					},
																					"local_path": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: `A local path within the VM to use.`,
																					},
																					"remote": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A generic remote file.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"uri": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `URI from which to fetch the object. It should contain both the protocol and path following the format '{protocol}://{location}'.`,
																								},
																								"sha256_checksum": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: `SHA256 checksum of the remote file.`,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"pull_deps": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: `Whether dependencies should also be installed. - install when false: 'dpkg -i package' - install when true: 'apt-get update && apt-get -y install package.deb'`,
																		},
																	},
																},
															},
															"googet": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A package managed by GooGet.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `Package name.`,
																		},
																	},
																},
															},
															"msi": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `An MSI package.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"source": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: `The MSI package.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"allow_insecure": {
																						Type:     schema.TypeBool,
																						Optional: true,
																						Description: `Defaults to false. When false, files are subject to validations based on the file type:
Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.`,
																					},
																					"gcs": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A Cloud Storage object.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Bucket of the Cloud Storage object.`,
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Name of the Cloud Storage object.`,
																								},
																								"generation": {
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: `Generation number of the Cloud Storage object.`,
																								},
																							},
																						},
																					},
																					"local_path": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: `A local path within the VM to use.`,
																					},
																					"remote": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A generic remote file.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"uri": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `URI from which to fetch the object. It should contain both the protocol and path following the format '{protocol}://{location}'.`,
																								},
																								"sha256_checksum": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: `SHA256 checksum of the remote file.`,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"properties": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `Additional properties to use during installation. This should be in the format of Property=Setting. Appended to the defaults of 'ACTION=INSTALL REBOOT=ReallySuppress'.`,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"rpm": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `An rpm package file.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"source": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: `An rpm package.`,
																			MaxItems:    1,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"allow_insecure": {
																						Type:     schema.TypeBool,
																						Optional: true,
																						Description: `Defaults to false. When false, files are subject to validations based on the file type:
Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.`,
																					},
																					"gcs": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A Cloud Storage object.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Bucket of the Cloud Storage object.`,
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `Name of the Cloud Storage object.`,
																								},
																								"generation": {
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: `Generation number of the Cloud Storage object.`,
																								},
																							},
																						},
																					},
																					"local_path": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: `A local path within the VM to use.`,
																					},
																					"remote": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: `A generic remote file.`,
																						MaxItems:    1,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"uri": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: `URI from which to fetch the object. It should contain both the protocol and path following the format '{protocol}://{location}'.`,
																								},
																								"sha256_checksum": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: `SHA256 checksum of the remote file.`,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"pull_deps": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: `Whether dependencies should also be installed. - install when false: 'rpm --upgrade --replacepkgs package.rpm' - install when true: 'yum -y install package.rpm' or 'zypper -y install package.rpm'`,
																		},
																	},
																},
															},
															"yum": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A package managed by YUM.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `Package name.`,
																		},
																	},
																},
															},
															"zypper": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A package managed by Zypper.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `Package name.`,
																		},
																	},
																},
															},
														},
													},
												},
												"repository": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: `Package repository resource`,
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"apt": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `An Apt Repository.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"archive_type": {
																			Type:         schema.TypeString,
																			Required:     true,
																			ValidateFunc: verify.ValidateEnum([]string{"ARCHIVE_TYPE_UNSPECIFIED", "DEB", "DEB_SRC"}),
																			Description:  `Type of archive files in this repository. Possible values: ["ARCHIVE_TYPE_UNSPECIFIED", "DEB", "DEB_SRC"]`,
																		},
																		"components": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: `List of components for this repository. Must contain at least one item.`,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"distribution": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `Distribution of this repository.`,
																		},
																		"uri": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `URI for this repository.`,
																		},
																		"gpg_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `URI of the key file for this repository. The agent maintains a keyring at '/etc/apt/trusted.gpg.d/osconfig_agent_managed.gpg'.`,
																		},
																	},
																},
															},
															"goo": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A Goo Repository.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `The name of the repository.`,
																		},
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `The url of the repository.`,
																		},
																	},
																},
															},
															"yum": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A Yum Repository.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"base_url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `The location of the repository directory.`,
																		},
																		"id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `A one word, unique name for this repository. This is the 'repo id' in the yum config file and also the 'display_name' if 'display_name' is omitted. This id is also used as the unique identifier when checking for resource conflicts.`,
																		},
																		"display_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `The display name of the repository.`,
																		},
																		"gpg_keys": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `URIs of GPG keys.`,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"zypper": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `A Zypper Repository.`,
																MaxItems:    1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"base_url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `The location of the repository directory.`,
																		},
																		"id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `A one word, unique name for this repository. This is the 'repo id' in the zypper config file and also the 'display_name' if 'display_name' is omitted. This id is also used as the unique identifier when checking for GuestPolicy conflicts.`,
																		},
																		"display_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `The display name of the repository.`,
																		},
																		"gpg_keys": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: `URIs of GPG keys.`,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
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
									"inventory_filters": {
										Type:     schema.TypeList,
										Optional: true,
										Description: `List of inventory filters for the resource group.
The resources in this resource group are applied to the target VM if it satisfies at least one of the following inventory filters.
For example, to apply this resource group to VMs running either 'RHEL' or 'CentOS' operating systems, specify 2 items for the list with following values: inventory_filters[0].os_short_name='rhel' and inventory_filters[1].os_short_name='centos'
If the list is empty, this resource group will be applied to the target VM unconditionally.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"os_short_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The OS short name`,
												},
												"os_version": {
													Type:     schema.TypeString,
													Optional: true,
													Description: `The OS version
Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of '7', specify the following value for this field '7.*'
An empty string matches all OS versions.`,
												},
											},
										},
									},
								},
							},
						},
						"allow_no_resource_group_match": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `This flag determines the OS policy compliance status when none of the resource groups within the policy are applicable for a VM. Set this value to 'true' if the policy needs to be reported as compliant even if the policy has nothing to validate or enforce.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Policy description. Length of the description is limited to 1024 characters.`,
						},
					},
				},
			},
			"rollout": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Rollout to deploy the OS policy assignment. A rollout is triggered in the following situations: 1) OSPolicyAssignment is created. 2) OSPolicyAssignment is updated and the update contains changes to one of the following fields: - instance_filter - os_policies 3) OSPolicyAssignment is deleted.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disruption_budget": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `The maximum number (or percentage) of VMs per zone to disrupt at any given moment.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Specifies a fixed value.`,
									},
									"percent": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Specifies the relative value defined as a percentage, which will be multiplied by a reference value.`,
									},
								},
							},
						},
						"min_wait_duration": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      `This determines the minimum duration of time to wait after the configuration changes are applied through the current rollout. A VM continues to count towards the 'disruption_budget' at least until this duration of time has passed after configuration changes are applied.`,
							DiffSuppressFunc: compareDuration,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `OS policy assignment description. Length of the description is limited to 1024 characters.`,
			},
			"baseline": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Output only. Indicates that this revision has been successfully rolled out in this zone and new VMs will be assigned OS policies from this revision.
For a given OS policy assignment, there is only one revision with a value of 'true' for this field.`,
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Output only. Indicates that this revision deletes the OS policy assignment.`,
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The etag for this OS policy assignment. If this is provided on update, it must match the server's etag.`,
			},
			"reconciling": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Output only. Indicates that reconciliation is in progress for the revision. This value is 'true' when the 'rollout_state' is one of:
* IN_PROGRESS
* CANCELLING`,
			},
			"revision_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The timestamp that the revision was created.`,
			},
			"revision_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The assignment revision ID A new revision is committed whenever a rollout is triggered for a OS policy assignment`,
			},
			"rollout_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. OS policy assignment rollout state`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Server generated unique id for the OS policy assignment resource.`,
			},
			"skip_await_rollout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Set to true to skip awaiting rollout during resource creation and update.`,
			},
			"project": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},
		},
		UseJSONNumber: true,
	}
}

func resourceOSConfigOSPolicyAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandOSConfigOSPolicyAssignmentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	osPoliciesProp, err := expandOSConfigOSPolicyAssignmentOsPolicies(d.Get("os_policies"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("os_policies"); !tpgresource.IsEmptyValue(reflect.ValueOf(osPoliciesProp)) && (ok || !reflect.DeepEqual(v, osPoliciesProp)) {
		obj["osPolicies"] = osPoliciesProp
	}
	instanceFilterProp, err := expandOSConfigOSPolicyAssignmentInstanceFilter(d.Get("instance_filter"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("instance_filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(instanceFilterProp)) && (ok || !reflect.DeepEqual(v, instanceFilterProp)) {
		obj["instanceFilter"] = instanceFilterProp
	}
	rolloutProp, err := expandOSConfigOSPolicyAssignmentRollout(d.Get("rollout"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rollout"); !tpgresource.IsEmptyValue(reflect.ValueOf(rolloutProp)) && (ok || !reflect.DeepEqual(v, rolloutProp)) {
		obj["rollout"] = rolloutProp
	}

	log.Printf("[DEBUG] Creating new OSPolicyAssignment: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OSPolicyAssignment: %s", err)
	}
	// Shorten long form project id to short form.
	billingProject = tpgresource.GetResourceNameFromSelfLink(project)

	url, err := tpgresource.ReplaceVars(d, config, "{{OSConfigBasePath}}projects/{{project}}/locations/{{location}}/osPolicyAssignments?osPolicyAssignmentId={{name}}")
	if err != nil {
		return err
	}
	// Always use GA endpoints for this resource.
	url = strings.ReplaceAll(url, "https://osconfig.googleapis.com/v1beta", "https://osconfig.googleapis.com/v1")
	// Remove redundant projects/ from url.
	url = strings.ReplaceAll(url, "projects/projects/", "projects/")

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating OSPolicyAssignment: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	// Remove redundant projects/ from id.
	id = strings.ReplaceAll(id, "projects/projects/", "projects/")
	d.SetId(id)

	if skipAwaitRollout := d.Get("skip_await_rollout").(bool); !skipAwaitRollout {
		// Use the resource in the operation response to populate
		// identity fields and d.Id() before read
		var opRes map[string]interface{}
		err = OSConfigOperationWaitTimeWithResponse(
			config, res, &opRes, project, "Creating OSPolicyAssignment", userAgent,
			d.Timeout(schema.TimeoutCreate))
		if err != nil {
			// The resource didn't actually create
			d.SetId("")

			return fmt.Errorf("Error waiting to create OSPolicyAssignment: %s", err)
		}

		if err := d.Set("name", flattenOSConfigOSPolicyAssignmentName(opRes["name"], d, config)); err != nil {
			return err
		}

		// This may have caused the ID to update - update it if so.
		id, err = tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
		if err != nil {
			return fmt.Errorf("Error constructing id: %s", err)
		}
		// Remove redundant projects/ from id.
		id = strings.ReplaceAll(id, "projects/projects/", "projects/")
		d.SetId(id)
	}

	log.Printf("[DEBUG] Finished creating OSPolicyAssignment %q: %#v", d.Id(), res)

	return resourceOSConfigOSPolicyAssignmentRead(d, meta)
}

func resourceOSConfigOSPolicyAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OSPolicyAssignment: %s", err)
	}
	// Shorten long form project id to short form
	billingProject = tpgresource.GetResourceNameFromSelfLink(project)

	url, err := tpgresource.ReplaceVars(d, config, "{{OSConfigBasePath}}projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return err
	}
	// Always use GA endpoints for this resource.
	url = strings.ReplaceAll(url, "https://osconfig.googleapis.com/v1beta", "https://osconfig.googleapis.com/v1")
	// Remove redundant projects/ from url.
	url = strings.ReplaceAll(url, "projects/projects/", "projects/")

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("OSConfigOSPolicyAssignment %q", d.Id()))
	}

	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOkExists("skip_await_rollout"); !ok {
		if err := d.Set("skip_await_rollout", false); err != nil {
			return fmt.Errorf("Error setting skip_await_rollout: %s", err)
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}

	if err := d.Set("name", flattenOSConfigOSPolicyAssignmentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("description", flattenOSConfigOSPolicyAssignmentDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("os_policies", flattenOSConfigOSPolicyAssignmentOsPolicies(res["osPolicies"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("instance_filter", flattenOSConfigOSPolicyAssignmentInstanceFilter(res["instanceFilter"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("rollout", flattenOSConfigOSPolicyAssignmentRollout(res["rollout"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("revision_id", flattenOSConfigOSPolicyAssignmentRevisionId(res["revisionId"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("revision_create_time", flattenOSConfigOSPolicyAssignmentRevisionCreateTime(res["revisionCreateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("etag", flattenOSConfigOSPolicyAssignmentEtag(res["etag"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("rollout_state", flattenOSConfigOSPolicyAssignmentRolloutState(res["rolloutState"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("baseline", flattenOSConfigOSPolicyAssignmentBaseline(res["baseline"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("deleted", flattenOSConfigOSPolicyAssignmentDeleted(res["deleted"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("reconciling", flattenOSConfigOSPolicyAssignmentReconciling(res["reconciling"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}
	if err := d.Set("uid", flattenOSConfigOSPolicyAssignmentUid(res["uid"], d, config)); err != nil {
		return fmt.Errorf("Error reading OSPolicyAssignment: %s", err)
	}

	return nil
}

func resourceOSConfigOSPolicyAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OSPolicyAssignment: %s", err)
	}
	// Shorten long form project id to short form
	billingProject = tpgresource.GetResourceNameFromSelfLink(project)

	obj := make(map[string]interface{})
	descriptionProp, err := expandOSConfigOSPolicyAssignmentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	osPoliciesProp, err := expandOSConfigOSPolicyAssignmentOsPolicies(d.Get("os_policies"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("os_policies"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, osPoliciesProp)) {
		obj["osPolicies"] = osPoliciesProp
	}
	instanceFilterProp, err := expandOSConfigOSPolicyAssignmentInstanceFilter(d.Get("instance_filter"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("instance_filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, instanceFilterProp)) {
		obj["instanceFilter"] = instanceFilterProp
	}
	rolloutProp, err := expandOSConfigOSPolicyAssignmentRollout(d.Get("rollout"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rollout"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, rolloutProp)) {
		obj["rollout"] = rolloutProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{OSConfigBasePath}}projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return err
	}
	// Always use GA endpoints for this resource.
	url = strings.ReplaceAll(url, "https://osconfig.googleapis.com/v1beta", "https://osconfig.googleapis.com/v1")
	// Remove redundant projects/ from url.
	url = strings.ReplaceAll(url, "projects/projects/", "projects/")

	log.Printf("[DEBUG] Updating OSPolicyAssignment %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("os_policies") {
		updateMask = append(updateMask, "osPolicies")
	}

	if d.HasChange("instance_filter") {
		updateMask = append(updateMask, "instanceFilter")
	}

	if d.HasChange("rollout") {
		updateMask = append(updateMask, "rollout")
	}
	// updateMask is a URL parameter but not present in the schema, so tpgresource.ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
	})

	if err != nil {
		return fmt.Errorf("Error updating OSPolicyAssignment %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating OSPolicyAssignment %q: %#v", d.Id(), res)
	}

	if skipAwaitRollout := d.Get("skip_await_rollout").(bool); !skipAwaitRollout {
		err = OSConfigOperationWaitTime(
			config, res, project, "Updating OSPolicyAssignment", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceOSConfigOSPolicyAssignmentRead(d, meta)
}

func resourceOSConfigOSPolicyAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OSPolicyAssignment: %s", err)
	}
	// Shorten long form project id to short form
	billingProject = tpgresource.GetResourceNameFromSelfLink(project)

	url, err := tpgresource.ReplaceVars(d, config, "{{OSConfigBasePath}}projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return err
	}
	// Always use GA endpoints for this resource.
	url = strings.ReplaceAll(url, "https://osconfig.googleapis.com/v1beta", "https://osconfig.googleapis.com/v1")
	// Remove redundant projects/ from url.
	url = strings.ReplaceAll(url, "projects/projects/", "projects/")

	log.Printf("[DEBUG] Deleting OSPolicyAssignment %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "OSPolicyAssignment")
	}

	if skipAwaitRollout := d.Get("skip_await_rollout").(bool); !skipAwaitRollout {
		err = OSConfigOperationWaitTime(
			config, res, project, "Deleting OSPolicyAssignment", userAgent,
			d.Timeout(schema.TimeoutDelete))

		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Finished deleting OSPolicyAssignment %q: %#v", d.Id(), res)
	return nil
}

func resourceOSConfigOSPolicyAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/osPolicyAssignments/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	// Remove redundant projects/ from id.
	id = strings.ReplaceAll(id, "projects/projects/", "projects/")
	d.SetId(id)

	// Explicitly set virtual fields to default values on import
	if err := d.Set("skip_await_rollout", false); err != nil {
		return nil, fmt.Errorf("Error setting skip_await_rollout: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func compareDuration(_, old, new string, _ *schema.ResourceData) bool {
	oldDuration, err := time.ParseDuration(old)
	if err != nil {
		return false
	}
	newDuration, err := time.ParseDuration(new)
	if err != nil {
		return false
	}
	return oldDuration == newDuration
}

func flattenOSConfigOSPolicyAssignmentName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenOSConfigOSPolicyAssignmentDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPolicies(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"id":                            flattenOSConfigOSPolicyAssignmentOsPoliciesId(original["id"], d, config),
			"description":                   flattenOSConfigOSPolicyAssignmentOsPoliciesDescription(original["description"], d, config),
			"mode":                          flattenOSConfigOSPolicyAssignmentOsPoliciesMode(original["mode"], d, config),
			"resource_groups":               flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroups(original["resourceGroups"], d, config),
			"allow_no_resource_group_match": flattenOSConfigOSPolicyAssignmentOsPoliciesAllowNoResourceGroupMatch(original["allowNoResourceGroupMatch"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroups(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"inventory_filters": flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFilters(original["inventoryFilters"], d, config),
			"resources":         flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResources(original["resources"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFilters(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"os_short_name": flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsShortName(original["osShortName"], d, config),
			"os_version":    flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsVersion(original["osVersion"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsShortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsVersion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResources(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"id":         flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesId(original["id"], d, config),
			"pkg":        flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkg(original["pkg"], d, config),
			"repository": flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepository(original["repository"], d, config),
			"exec":       flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExec(original["exec"], d, config),
			"file":       flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFile(original["file"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkg(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["desired_state"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDesiredState(original["desiredState"], d, config)
	transformed["apt"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgApt(original["apt"], d, config)
	transformed["deb"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDeb(original["deb"], d, config)
	transformed["yum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYum(original["yum"], d, config)
	transformed["zypper"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypper(original["zypper"], d, config)
	transformed["rpm"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpm(original["rpm"], d, config)
	transformed["googet"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGooget(original["googet"], d, config)
	transformed["msi"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsi(original["msi"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDesiredState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgApt(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgAptName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgAptName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDeb(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["source"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSource(original["source"], d, config)
	transformed["pull_deps"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebPullDeps(original["pullDeps"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSource(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["remote"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemote(original["remote"], d, config)
	transformed["gcs"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcs(original["gcs"], d, config)
	transformed["local_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceLocalPath(original["localPath"], d, config)
	transformed["allow_insecure"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceAllowInsecure(original["allowInsecure"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemote(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteUri(original["uri"], d, config)
	transformed["sha256_checksum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteSha256Checksum(original["sha256Checksum"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteSha256Checksum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsBucket(original["bucket"], d, config)
	transformed["object"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsObject(original["object"], d, config)
	transformed["generation"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsGeneration(original["generation"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsObject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsGeneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceLocalPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceAllowInsecure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebPullDeps(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYumName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYumName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypper(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypperName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypperName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpm(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["source"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSource(original["source"], d, config)
	transformed["pull_deps"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmPullDeps(original["pullDeps"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSource(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["remote"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemote(original["remote"], d, config)
	transformed["gcs"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcs(original["gcs"], d, config)
	transformed["local_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceLocalPath(original["localPath"], d, config)
	transformed["allow_insecure"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceAllowInsecure(original["allowInsecure"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemote(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteUri(original["uri"], d, config)
	transformed["sha256_checksum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSha256Checksum(original["sha256Checksum"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSha256Checksum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsBucket(original["bucket"], d, config)
	transformed["object"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsObject(original["object"], d, config)
	transformed["generation"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsGeneration(original["generation"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsObject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsGeneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceLocalPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceAllowInsecure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmPullDeps(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGooget(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGoogetName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGoogetName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsi(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["source"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSource(original["source"], d, config)
	transformed["properties"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiProperties(original["properties"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSource(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["remote"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemote(original["remote"], d, config)
	transformed["gcs"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcs(original["gcs"], d, config)
	transformed["local_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceLocalPath(original["localPath"], d, config)
	transformed["allow_insecure"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceAllowInsecure(original["allowInsecure"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemote(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteUri(original["uri"], d, config)
	transformed["sha256_checksum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSha256Checksum(original["sha256Checksum"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSha256Checksum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsBucket(original["bucket"], d, config)
	transformed["object"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsObject(original["object"], d, config)
	transformed["generation"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsGeneration(original["generation"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsObject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsGeneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceLocalPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceAllowInsecure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiProperties(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepository(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["apt"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryApt(original["apt"], d, config)
	transformed["yum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYum(original["yum"], d, config)
	transformed["zypper"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypper(original["zypper"], d, config)
	transformed["goo"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGoo(original["goo"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryApt(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["archive_type"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptArchiveType(original["archiveType"], d, config)
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptUri(original["uri"], d, config)
	transformed["distribution"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptDistribution(original["distribution"], d, config)
	transformed["components"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptComponents(original["components"], d, config)
	transformed["gpg_key"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptGpgKey(original["gpgKey"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptArchiveType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptDistribution(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptComponents(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptGpgKey(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["id"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumId(original["id"], d, config)
	transformed["display_name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumDisplayName(original["displayName"], d, config)
	transformed["base_url"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumBaseUrl(original["baseUrl"], d, config)
	transformed["gpg_keys"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumGpgKeys(original["gpgKeys"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumBaseUrl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumGpgKeys(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypper(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["id"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperId(original["id"], d, config)
	transformed["display_name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperDisplayName(original["displayName"], d, config)
	transformed["base_url"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperBaseUrl(original["baseUrl"], d, config)
	transformed["gpg_keys"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperGpgKeys(original["gpgKeys"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperBaseUrl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperGpgKeys(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGoo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooName(original["name"], d, config)
	transformed["url"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooUrl(original["url"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooUrl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["validate"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidate(original["validate"], d, config)
	transformed["enforce"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforce(original["enforce"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["file"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFile(original["file"], d, config)
	transformed["script"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateScript(original["script"], d, config)
	transformed["args"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateArgs(original["args"], d, config)
	transformed["interpreter"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateInterpreter(original["interpreter"], d, config)
	transformed["output_file_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateOutputFilePath(original["outputFilePath"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFile(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["remote"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemote(original["remote"], d, config)
	transformed["gcs"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcs(original["gcs"], d, config)
	transformed["local_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileLocalPath(original["localPath"], d, config)
	transformed["allow_insecure"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileAllowInsecure(original["allowInsecure"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemote(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteUri(original["uri"], d, config)
	transformed["sha256_checksum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteSha256Checksum(original["sha256Checksum"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteSha256Checksum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsBucket(original["bucket"], d, config)
	transformed["object"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsObject(original["object"], d, config)
	transformed["generation"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsGeneration(original["generation"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsObject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsGeneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileLocalPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileAllowInsecure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateScript(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateArgs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateInterpreter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateOutputFilePath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforce(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["file"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFile(original["file"], d, config)
	transformed["script"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceScript(original["script"], d, config)
	transformed["args"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceArgs(original["args"], d, config)
	transformed["interpreter"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceInterpreter(original["interpreter"], d, config)
	transformed["output_file_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceOutputFilePath(original["outputFilePath"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFile(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["remote"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemote(original["remote"], d, config)
	transformed["gcs"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcs(original["gcs"], d, config)
	transformed["local_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileLocalPath(original["localPath"], d, config)
	transformed["allow_insecure"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileAllowInsecure(original["allowInsecure"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemote(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteUri(original["uri"], d, config)
	transformed["sha256_checksum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteSha256Checksum(original["sha256Checksum"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteSha256Checksum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsBucket(original["bucket"], d, config)
	transformed["object"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsObject(original["object"], d, config)
	transformed["generation"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsGeneration(original["generation"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsObject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsGeneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileLocalPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileAllowInsecure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceScript(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceArgs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceInterpreter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceOutputFilePath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFile(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["file"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFile(original["file"], d, config)
	transformed["content"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileContent(original["content"], d, config)
	transformed["path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePath(original["path"], d, config)
	transformed["state"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileState(original["state"], d, config)
	transformed["permissions"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePermissions(original["permissions"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFile(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["remote"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemote(original["remote"], d, config)
	transformed["gcs"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcs(original["gcs"], d, config)
	transformed["local_path"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileLocalPath(original["localPath"], d, config)
	transformed["allow_insecure"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileAllowInsecure(original["allowInsecure"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemote(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteUri(original["uri"], d, config)
	transformed["sha256_checksum"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteSha256Checksum(original["sha256Checksum"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteSha256Checksum(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcs(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsBucket(original["bucket"], d, config)
	transformed["object"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsObject(original["object"], d, config)
	transformed["generation"] =
		flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsGeneration(original["generation"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsBucket(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsObject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsGeneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileLocalPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileAllowInsecure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileContent(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePermissions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentOsPoliciesAllowNoResourceGroupMatch(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentInstanceFilter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["all"] =
		flattenOSConfigOSPolicyAssignmentInstanceFilterAll(original["all"], d, config)
	transformed["inclusion_labels"] =
		flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(original["inclusionLabels"], d, config)
	transformed["exclusion_labels"] =
		flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(original["exclusionLabels"], d, config)
	transformed["inventories"] =
		flattenOSConfigOSPolicyAssignmentInstanceFilterInventories(original["inventories"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentInstanceFilterAll(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"labels": flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsLabels(original["labels"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"labels": flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsLabels(original["labels"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInventories(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"os_short_name": flattenOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsShortName(original["osShortName"], d, config),
			"os_version":    flattenOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsVersion(original["osVersion"], d, config),
		})
	}
	return transformed
}
func flattenOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsShortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsVersion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentRollout(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["disruption_budget"] =
		flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudget(original["disruptionBudget"], d, config)
	transformed["min_wait_duration"] =
		flattenOSConfigOSPolicyAssignmentRolloutMinWaitDuration(original["minWaitDuration"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudget(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["fixed"] =
		flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudgetFixed(original["fixed"], d, config)
	transformed["percent"] =
		flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudgetPercent(original["percent"], d, config)
	return []interface{}{transformed}
}
func flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudgetFixed(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudgetPercent(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenOSConfigOSPolicyAssignmentRolloutMinWaitDuration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentRevisionId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentRevisionCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentEtag(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentRolloutState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentBaseline(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentDeleted(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentReconciling(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenOSConfigOSPolicyAssignmentUid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandOSConfigOSPolicyAssignmentName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPolicies(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedId, err := expandOSConfigOSPolicyAssignmentOsPoliciesId(original["id"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["id"] = transformedId
		}

		transformedDescription, err := expandOSConfigOSPolicyAssignmentOsPoliciesDescription(original["description"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["description"] = transformedDescription
		}

		transformedMode, err := expandOSConfigOSPolicyAssignmentOsPoliciesMode(original["mode"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMode); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["mode"] = transformedMode
		}

		transformedResourceGroups, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroups(original["resource_groups"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedResourceGroups); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["resourceGroups"] = transformedResourceGroups
		}

		transformedAllowNoResourceGroupMatch, err := expandOSConfigOSPolicyAssignmentOsPoliciesAllowNoResourceGroupMatch(original["allow_no_resource_group_match"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAllowNoResourceGroupMatch); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["allowNoResourceGroupMatch"] = transformedAllowNoResourceGroupMatch
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesMode(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroups(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedInventoryFilters, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFilters(original["inventory_filters"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedInventoryFilters); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["inventoryFilters"] = transformedInventoryFilters
		}

		transformedResources, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResources(original["resources"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedResources); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["resources"] = transformedResources
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFilters(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedOsShortName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsShortName(original["os_short_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedOsShortName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["osShortName"] = transformedOsShortName
		}

		transformedOsVersion, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsVersion(original["os_version"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedOsVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["osVersion"] = transformedOsVersion
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsShortName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsInventoryFiltersOsVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResources(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedId, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesId(original["id"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["id"] = transformedId
		}

		transformedPkg, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkg(original["pkg"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPkg); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["pkg"] = transformedPkg
		}

		transformedRepository, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepository(original["repository"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRepository); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["repository"] = transformedRepository
		}

		transformedExec, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExec(original["exec"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedExec); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["exec"] = transformedExec
		}

		transformedFile, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFile(original["file"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedFile); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["file"] = transformedFile
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkg(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDesiredState, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDesiredState(original["desired_state"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDesiredState); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["desiredState"] = transformedDesiredState
	}

	transformedApt, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgApt(original["apt"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedApt); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["apt"] = transformedApt
	}

	transformedDeb, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDeb(original["deb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDeb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["deb"] = transformedDeb
	}

	transformedYum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYum(original["yum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedYum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["yum"] = transformedYum
	}

	transformedZypper, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypper(original["zypper"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedZypper); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["zypper"] = transformedZypper
	}

	transformedRpm, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpm(original["rpm"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRpm); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["rpm"] = transformedRpm
	}

	transformedGooget, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGooget(original["googet"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGooget); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["googet"] = transformedGooget
	}

	transformedMsi, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsi(original["msi"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMsi); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["msi"] = transformedMsi
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDesiredState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgApt(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgAptName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgAptName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDeb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSource, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSource(original["source"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSource); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["source"] = transformedSource
	}

	transformedPullDeps, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebPullDeps(original["pull_deps"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPullDeps); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["pullDeps"] = transformedPullDeps
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRemote, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemote(original["remote"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRemote); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["remote"] = transformedRemote
	}

	transformedGcs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcs(original["gcs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gcs"] = transformedGcs
	}

	transformedLocalPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceLocalPath(original["local_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localPath"] = transformedLocalPath
	}

	transformedAllowInsecure, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceAllowInsecure(original["allow_insecure"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowInsecure); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowInsecure"] = transformedAllowInsecure
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemote(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedSha256Checksum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteSha256Checksum(original["sha256_checksum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256Checksum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha256Checksum"] = transformedSha256Checksum
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceRemoteSha256Checksum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedObject, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsObject(original["object"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedObject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["object"] = transformedObject
	}

	transformedGeneration, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsGeneration(original["generation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGeneration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["generation"] = transformedGeneration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceGcsGeneration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceLocalPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebSourceAllowInsecure(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgDebPullDeps(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYumName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgYumName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypper(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypperName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgZypperName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpm(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSource, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSource(original["source"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSource); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["source"] = transformedSource
	}

	transformedPullDeps, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmPullDeps(original["pull_deps"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPullDeps); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["pullDeps"] = transformedPullDeps
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRemote, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemote(original["remote"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRemote); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["remote"] = transformedRemote
	}

	transformedGcs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcs(original["gcs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gcs"] = transformedGcs
	}

	transformedLocalPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceLocalPath(original["local_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localPath"] = transformedLocalPath
	}

	transformedAllowInsecure, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceAllowInsecure(original["allow_insecure"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowInsecure); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowInsecure"] = transformedAllowInsecure
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemote(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedSha256Checksum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSha256Checksum(original["sha256_checksum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256Checksum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha256Checksum"] = transformedSha256Checksum
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSha256Checksum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedObject, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsObject(original["object"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedObject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["object"] = transformedObject
	}

	transformedGeneration, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsGeneration(original["generation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGeneration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["generation"] = transformedGeneration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceGcsGeneration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceLocalPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmSourceAllowInsecure(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgRpmPullDeps(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGooget(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGoogetName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgGoogetName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsi(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSource, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSource(original["source"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSource); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["source"] = transformedSource
	}

	transformedProperties, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiProperties(original["properties"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProperties); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["properties"] = transformedProperties
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRemote, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemote(original["remote"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRemote); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["remote"] = transformedRemote
	}

	transformedGcs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcs(original["gcs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gcs"] = transformedGcs
	}

	transformedLocalPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceLocalPath(original["local_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localPath"] = transformedLocalPath
	}

	transformedAllowInsecure, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceAllowInsecure(original["allow_insecure"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowInsecure); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowInsecure"] = transformedAllowInsecure
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemote(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedSha256Checksum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSha256Checksum(original["sha256_checksum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256Checksum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha256Checksum"] = transformedSha256Checksum
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSha256Checksum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedObject, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsObject(original["object"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedObject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["object"] = transformedObject
	}

	transformedGeneration, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsGeneration(original["generation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGeneration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["generation"] = transformedGeneration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceGcsGeneration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceLocalPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiSourceAllowInsecure(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesPkgMsiProperties(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepository(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedApt, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryApt(original["apt"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedApt); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["apt"] = transformedApt
	}

	transformedYum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYum(original["yum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedYum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["yum"] = transformedYum
	}

	transformedZypper, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypper(original["zypper"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedZypper); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["zypper"] = transformedZypper
	}

	transformedGoo, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGoo(original["goo"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGoo); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["goo"] = transformedGoo
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryApt(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedArchiveType, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptArchiveType(original["archive_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedArchiveType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["archiveType"] = transformedArchiveType
	}

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedDistribution, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptDistribution(original["distribution"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDistribution); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["distribution"] = transformedDistribution
	}

	transformedComponents, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptComponents(original["components"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedComponents); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["components"] = transformedComponents
	}

	transformedGpgKey, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptGpgKey(original["gpg_key"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGpgKey); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gpgKey"] = transformedGpgKey
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptArchiveType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptDistribution(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptComponents(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryAptGpgKey(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedId, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumId(original["id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["id"] = transformedId
	}

	transformedDisplayName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumDisplayName(original["display_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDisplayName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["displayName"] = transformedDisplayName
	}

	transformedBaseUrl, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumBaseUrl(original["base_url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBaseUrl); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["baseUrl"] = transformedBaseUrl
	}

	transformedGpgKeys, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumGpgKeys(original["gpg_keys"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGpgKeys); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gpgKeys"] = transformedGpgKeys
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumBaseUrl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryYumGpgKeys(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypper(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedId, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperId(original["id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["id"] = transformedId
	}

	transformedDisplayName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperDisplayName(original["display_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDisplayName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["displayName"] = transformedDisplayName
	}

	transformedBaseUrl, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperBaseUrl(original["base_url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBaseUrl); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["baseUrl"] = transformedBaseUrl
	}

	transformedGpgKeys, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperGpgKeys(original["gpg_keys"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGpgKeys); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gpgKeys"] = transformedGpgKeys
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperBaseUrl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryZypperGpgKeys(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGoo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedName, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	transformedUrl, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooUrl(original["url"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUrl); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["url"] = transformedUrl
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesRepositoryGooUrl(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedValidate, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidate(original["validate"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedValidate); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["validate"] = transformedValidate
	}

	transformedEnforce, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforce(original["enforce"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnforce); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enforce"] = transformedEnforce
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidate(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFile, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFile(original["file"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFile); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["file"] = transformedFile
	}

	transformedScript, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateScript(original["script"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedScript); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["script"] = transformedScript
	}

	transformedArgs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateArgs(original["args"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedArgs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["args"] = transformedArgs
	}

	transformedInterpreter, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateInterpreter(original["interpreter"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInterpreter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["interpreter"] = transformedInterpreter
	}

	transformedOutputFilePath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateOutputFilePath(original["output_file_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOutputFilePath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["outputFilePath"] = transformedOutputFilePath
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRemote, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemote(original["remote"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRemote); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["remote"] = transformedRemote
	}

	transformedGcs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcs(original["gcs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gcs"] = transformedGcs
	}

	transformedLocalPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileLocalPath(original["local_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localPath"] = transformedLocalPath
	}

	transformedAllowInsecure, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileAllowInsecure(original["allow_insecure"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowInsecure); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowInsecure"] = transformedAllowInsecure
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemote(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedSha256Checksum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteSha256Checksum(original["sha256_checksum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256Checksum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha256Checksum"] = transformedSha256Checksum
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileRemoteSha256Checksum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedObject, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsObject(original["object"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedObject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["object"] = transformedObject
	}

	transformedGeneration, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsGeneration(original["generation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGeneration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["generation"] = transformedGeneration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileGcsGeneration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileLocalPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateFileAllowInsecure(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateScript(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateArgs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateInterpreter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecValidateOutputFilePath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforce(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFile, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFile(original["file"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFile); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["file"] = transformedFile
	}

	transformedScript, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceScript(original["script"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedScript); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["script"] = transformedScript
	}

	transformedArgs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceArgs(original["args"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedArgs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["args"] = transformedArgs
	}

	transformedInterpreter, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceInterpreter(original["interpreter"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInterpreter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["interpreter"] = transformedInterpreter
	}

	transformedOutputFilePath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceOutputFilePath(original["output_file_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOutputFilePath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["outputFilePath"] = transformedOutputFilePath
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRemote, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemote(original["remote"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRemote); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["remote"] = transformedRemote
	}

	transformedGcs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcs(original["gcs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gcs"] = transformedGcs
	}

	transformedLocalPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileLocalPath(original["local_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localPath"] = transformedLocalPath
	}

	transformedAllowInsecure, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileAllowInsecure(original["allow_insecure"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowInsecure); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowInsecure"] = transformedAllowInsecure
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemote(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedSha256Checksum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteSha256Checksum(original["sha256_checksum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256Checksum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha256Checksum"] = transformedSha256Checksum
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileRemoteSha256Checksum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedObject, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsObject(original["object"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedObject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["object"] = transformedObject
	}

	transformedGeneration, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsGeneration(original["generation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGeneration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["generation"] = transformedGeneration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileGcsGeneration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileLocalPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceFileAllowInsecure(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceScript(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceArgs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceInterpreter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesExecEnforceOutputFilePath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFile, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFile(original["file"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFile); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["file"] = transformedFile
	}

	transformedContent, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileContent(original["content"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedContent); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["content"] = transformedContent
	}

	transformedPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePath(original["path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["path"] = transformedPath
	}

	transformedState, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileState(original["state"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedState); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["state"] = transformedState
	}

	transformedPermissions, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePermissions(original["permissions"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPermissions); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["permissions"] = transformedPermissions
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFile(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRemote, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemote(original["remote"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRemote); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["remote"] = transformedRemote
	}

	transformedGcs, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcs(original["gcs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["gcs"] = transformedGcs
	}

	transformedLocalPath, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileLocalPath(original["local_path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocalPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["localPath"] = transformedLocalPath
	}

	transformedAllowInsecure, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileAllowInsecure(original["allow_insecure"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowInsecure); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["allowInsecure"] = transformedAllowInsecure
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemote(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUri, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteUri(original["uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUri); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["uri"] = transformedUri
	}

	transformedSha256Checksum, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteSha256Checksum(original["sha256_checksum"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha256Checksum); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sha256Checksum"] = transformedSha256Checksum
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteUri(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileRemoteSha256Checksum(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucket, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsBucket(original["bucket"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBucket); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["bucket"] = transformedBucket
	}

	transformedObject, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsObject(original["object"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedObject); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["object"] = transformedObject
	}

	transformedGeneration, err := expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsGeneration(original["generation"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGeneration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["generation"] = transformedGeneration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsObject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileGcsGeneration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileLocalPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileFileAllowInsecure(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileContent(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFileState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesResourceGroupsResourcesFilePermissions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentOsPoliciesAllowNoResourceGroupMatch(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAll, err := expandOSConfigOSPolicyAssignmentInstanceFilterAll(original["all"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAll); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["all"] = transformedAll
	}

	transformedInclusionLabels, err := expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(original["inclusion_labels"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInclusionLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["inclusionLabels"] = transformedInclusionLabels
	}

	transformedExclusionLabels, err := expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(original["exclusion_labels"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExclusionLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["exclusionLabels"] = transformedExclusionLabels
	}

	transformedInventories, err := expandOSConfigOSPolicyAssignmentInstanceFilterInventories(original["inventories"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedInventories); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["inventories"] = transformedInventories
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterAll(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedLabels, err := expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsLabels(original["labels"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["labels"] = transformedLabels
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedLabels, err := expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsLabels(original["labels"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedLabels); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["labels"] = transformedLabels
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInventories(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedOsShortName, err := expandOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsShortName(original["os_short_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedOsShortName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["osShortName"] = transformedOsShortName
		}

		transformedOsVersion, err := expandOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsVersion(original["os_version"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedOsVersion); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["osVersion"] = transformedOsVersion
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsShortName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInventoriesOsVersion(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentRollout(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDisruptionBudget, err := expandOSConfigOSPolicyAssignmentRolloutDisruptionBudget(original["disruption_budget"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDisruptionBudget); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["disruptionBudget"] = transformedDisruptionBudget
	}

	transformedMinWaitDuration, err := expandOSConfigOSPolicyAssignmentRolloutMinWaitDuration(original["min_wait_duration"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMinWaitDuration); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["minWaitDuration"] = transformedMinWaitDuration
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentRolloutDisruptionBudget(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFixed, err := expandOSConfigOSPolicyAssignmentRolloutDisruptionBudgetFixed(original["fixed"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFixed); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["fixed"] = transformedFixed
	}

	transformedPercent, err := expandOSConfigOSPolicyAssignmentRolloutDisruptionBudgetPercent(original["percent"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPercent); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["percent"] = transformedPercent
	}

	return transformed, nil
}

func expandOSConfigOSPolicyAssignmentRolloutDisruptionBudgetFixed(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentRolloutDisruptionBudgetPercent(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandOSConfigOSPolicyAssignmentRolloutMinWaitDuration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
