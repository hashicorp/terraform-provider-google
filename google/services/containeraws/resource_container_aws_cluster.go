// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package containeraws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	containeraws "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceContainerAwsCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerAwsClusterCreate,
		Read:   resourceContainerAwsClusterRead,
		Update: resourceContainerAwsClusterUpdate,
		Delete: resourceContainerAwsClusterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceContainerAwsClusterImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"authorization": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configuration related to the cluster RBAC settings.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterAuthorizationSchema(),
			},

			"aws_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AWS region where the cluster runs. Each Google Cloud region supports a subset of nearby AWS regions. You can call to list all supported AWS regions within a given Google Cloud region.",
			},

			"control_plane": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Configuration related to the cluster control plane.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneSchema(),
			},

			"fleet": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Fleet configuration.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterFleetSchema(),
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of this resource.",
			},

			"networking": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Cluster-wide networking configuration.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterNetworkingSchema(),
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Annotations on the cluster. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Key can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. A human readable description of this cluster. Cannot be longer than 255 UTF-8 encoded bytes.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this cluster was created.",
			},

			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The endpoint of the cluster's API server.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"reconciling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. If set, there are currently changes in flight to the cluster.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The current state of the cluster. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. A globally unique identifier for the cluster.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this cluster was last updated.",
			},

			"workload_identity_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Workload Identity settings.",
				Elem:        ContainerAwsClusterWorkloadIdentityConfigSchema(),
			},
		},
	}
}

func ContainerAwsClusterAuthorizationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"admin_users": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Users to perform operations as a cluster admin. A managed ClusterRoleBinding will be created to grant the `cluster-admin` ClusterRole to the users. Up to ten admin users can be provided. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles",
				Elem:        ContainerAwsClusterAuthorizationAdminUsersSchema(),
			},
		},
	}
}

func ContainerAwsClusterAuthorizationAdminUsersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user, e.g. `my-gcp-id@gmail.com`.",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"aws_services_authentication": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Authentication configuration for management of AWS resources.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneAwsServicesAuthenticationSchema(),
			},

			"config_encryption": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The ARN of the AWS KMS key used to encrypt cluster configuration.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneConfigEncryptionSchema(),
			},

			"database_encryption": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The ARN of the AWS KMS key used to encrypt cluster secrets.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneDatabaseEncryptionSchema(),
			},

			"iam_instance_profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the AWS IAM instance pofile to assign to each control plane replica.",
			},

			"subnet_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The list of subnets where control plane replicas will run. A replica will be provisioned on each subnet and up to three values can be provided. Each subnet must be in a different AWS Availability Zone (AZ).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Kubernetes version to run on control plane replicas (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling .",
			},

			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. The AWS instance type. When unspecified, it defaults to `m5.large`.",
			},

			"main_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Configuration related to the main volume provisioned for each control plane replica. The main volume is in charge of storing all of the cluster's etcd state. Volumes will be provisioned in the availability zone associated with the corresponding subnet. When unspecified, it defaults to 8 GiB with the GP2 volume type.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneMainVolumeSchema(),
			},

			"proxy_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Proxy configuration for outbound HTTP(S) traffic.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneProxyConfigSchema(),
			},

			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Configuration related to the root volume provisioned for each control plane replica. Volumes will be provisioned in the availability zone associated with the corresponding subnet. When unspecified, it defaults to 32 GiB with the GP2 volume type.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneRootVolumeSchema(),
			},

			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. The IDs of additional security groups to add to control plane replicas. The Anthos Multi-Cloud API will automatically create and manage security groups with the minimum rules needed for a functioning cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"ssh_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. SSH configuration for how to access the underlying control plane machines.",
				MaxItems:    1,
				Elem:        ContainerAwsClusterControlPlaneSshConfigSchema(),
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. A set of AWS resource tags to propagate to all underlying managed AWS resources. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ContainerAwsClusterControlPlaneAwsServicesAuthenticationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Amazon Resource Name (ARN) of the role that the Anthos Multi-Cloud API will assume when managing AWS resources on your account.",
			},

			"role_session_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. An identifier for the assumed role session. When unspecified, it defaults to `multicloud-service-agent`.",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneConfigEncryptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kms_key_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of the AWS KMS key used to encrypt cluster configuration.",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneDatabaseEncryptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kms_key_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARN of the AWS KMS key used to encrypt cluster secrets.",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneMainVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.",
			},

			"kms_key_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.",
			},

			"size_gib": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
			},

			"throughput": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3.",
			},

			"volume_type": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareCaseInsensitive,
				Description:      "Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneProxyConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"secret_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.",
			},

			"secret_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version string of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneRootVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.",
			},

			"kms_key_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.",
			},

			"size_gib": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
			},

			"throughput": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3.",
			},

			"volume_type": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareCaseInsensitive,
				Description:      "Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3",
			},
		},
	}
}

func ContainerAwsClusterControlPlaneSshConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ec2_key_pair": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the EC2 key pair used to login into cluster machines.",
			},
		},
	}
}

func ContainerAwsClusterFleetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The number of the Fleet host project where this cluster will be registered.",
			},

			"membership": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the managed Hub Membership resource associated to this cluster. Membership names are formatted as projects/<project-number>/locations/global/membership/<cluster-id>.",
			},
		},
	}
}

func ContainerAwsClusterNetworkingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_address_cidr_blocks": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "All pods in the cluster are assigned an RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"service_address_cidr_blocks": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "All services in the cluster are assigned an RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC associated with the cluster. All component clusters (i.e. control plane and node pools) run on a single VPC. This field cannot be changed after creation.",
			},

			"per_node_pool_sg_rules_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable the per node pool subnet security group rules on the control plane security group. When set to true, you must also provide one or more security groups that ensure node pools are able to send requests to the control plane on TCP/443 and TCP/8132. Failure to do so may result in unavailable node pools.",
			},
		},
	}
}

func ContainerAwsClusterWorkloadIdentityConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identity_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the OIDC Identity Provider (IdP) associated to the Workload Identity Pool.",
			},

			"issuer_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OIDC issuer URL for this cluster.",
			},

			"workload_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Workload Identity Pool associated to the cluster.",
			},
		},
	}
}

func resourceContainerAwsClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containeraws.Cluster{
		Authorization: expandContainerAwsClusterAuthorization(d.Get("authorization")),
		AwsRegion:     dcl.String(d.Get("aws_region").(string)),
		ControlPlane:  expandContainerAwsClusterControlPlane(d.Get("control_plane")),
		Fleet:         expandContainerAwsClusterFleet(d.Get("fleet")),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Networking:    expandContainerAwsClusterNetworking(d.Get("networking")),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		Description:   dcl.String(d.Get("description").(string)),
		Project:       dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAwsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyCluster(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Cluster: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Cluster %q: %#v", d.Id(), res)

	return resourceContainerAwsClusterRead(d, meta)
}

func resourceContainerAwsClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containeraws.Cluster{
		Authorization: expandContainerAwsClusterAuthorization(d.Get("authorization")),
		AwsRegion:     dcl.String(d.Get("aws_region").(string)),
		ControlPlane:  expandContainerAwsClusterControlPlane(d.Get("control_plane")),
		Fleet:         expandContainerAwsClusterFleet(d.Get("fleet")),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Networking:    expandContainerAwsClusterNetworking(d.Get("networking")),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		Description:   dcl.String(d.Get("description").(string)),
		Project:       dcl.String(project),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAwsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetCluster(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ContainerAwsCluster %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("authorization", flattenContainerAwsClusterAuthorization(res.Authorization)); err != nil {
		return fmt.Errorf("error setting authorization in state: %s", err)
	}
	if err = d.Set("aws_region", res.AwsRegion); err != nil {
		return fmt.Errorf("error setting aws_region in state: %s", err)
	}
	if err = d.Set("control_plane", flattenContainerAwsClusterControlPlane(res.ControlPlane)); err != nil {
		return fmt.Errorf("error setting control_plane in state: %s", err)
	}
	if err = d.Set("fleet", flattenContainerAwsClusterFleet(res.Fleet)); err != nil {
		return fmt.Errorf("error setting fleet in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("networking", flattenContainerAwsClusterNetworking(res.Networking)); err != nil {
		return fmt.Errorf("error setting networking in state: %s", err)
	}
	if err = d.Set("annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("endpoint", res.Endpoint); err != nil {
		return fmt.Errorf("error setting endpoint in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("reconciling", res.Reconciling); err != nil {
		return fmt.Errorf("error setting reconciling in state: %s", err)
	}
	if err = d.Set("state", res.State); err != nil {
		return fmt.Errorf("error setting state in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}
	if err = d.Set("workload_identity_config", flattenContainerAwsClusterWorkloadIdentityConfig(res.WorkloadIdentityConfig)); err != nil {
		return fmt.Errorf("error setting workload_identity_config in state: %s", err)
	}

	return nil
}
func resourceContainerAwsClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containeraws.Cluster{
		Authorization: expandContainerAwsClusterAuthorization(d.Get("authorization")),
		AwsRegion:     dcl.String(d.Get("aws_region").(string)),
		ControlPlane:  expandContainerAwsClusterControlPlane(d.Get("control_plane")),
		Fleet:         expandContainerAwsClusterFleet(d.Get("fleet")),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Networking:    expandContainerAwsClusterNetworking(d.Get("networking")),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		Description:   dcl.String(d.Get("description").(string)),
		Project:       dcl.String(project),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAwsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyCluster(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Cluster: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Cluster %q: %#v", d.Id(), res)

	return resourceContainerAwsClusterRead(d, meta)
}

func resourceContainerAwsClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &containeraws.Cluster{
		Authorization: expandContainerAwsClusterAuthorization(d.Get("authorization")),
		AwsRegion:     dcl.String(d.Get("aws_region").(string)),
		ControlPlane:  expandContainerAwsClusterControlPlane(d.Get("control_plane")),
		Fleet:         expandContainerAwsClusterFleet(d.Get("fleet")),
		Location:      dcl.String(d.Get("location").(string)),
		Name:          dcl.String(d.Get("name").(string)),
		Networking:    expandContainerAwsClusterNetworking(d.Get("networking")),
		Annotations:   tpgresource.CheckStringMap(d.Get("annotations")),
		Description:   dcl.String(d.Get("description").(string)),
		Project:       dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Cluster %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLContainerAwsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteCluster(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Cluster: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Cluster %q", d.Id())
	return nil
}

func resourceContainerAwsClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/awsClusters/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/awsClusters/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandContainerAwsClusterAuthorization(o interface{}) *containeraws.ClusterAuthorization {
	if o == nil {
		return containeraws.EmptyClusterAuthorization
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterAuthorization
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterAuthorization{
		AdminUsers: expandContainerAwsClusterAuthorizationAdminUsersArray(obj["admin_users"]),
	}
}

func flattenContainerAwsClusterAuthorization(obj *containeraws.ClusterAuthorization) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"admin_users": flattenContainerAwsClusterAuthorizationAdminUsersArray(obj.AdminUsers),
	}

	return []interface{}{transformed}

}
func expandContainerAwsClusterAuthorizationAdminUsersArray(o interface{}) []containeraws.ClusterAuthorizationAdminUsers {
	if o == nil {
		return make([]containeraws.ClusterAuthorizationAdminUsers, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]containeraws.ClusterAuthorizationAdminUsers, 0)
	}

	items := make([]containeraws.ClusterAuthorizationAdminUsers, 0, len(objs))
	for _, item := range objs {
		i := expandContainerAwsClusterAuthorizationAdminUsers(item)
		items = append(items, *i)
	}

	return items
}

func expandContainerAwsClusterAuthorizationAdminUsers(o interface{}) *containeraws.ClusterAuthorizationAdminUsers {
	if o == nil {
		return containeraws.EmptyClusterAuthorizationAdminUsers
	}

	obj := o.(map[string]interface{})
	return &containeraws.ClusterAuthorizationAdminUsers{
		Username: dcl.String(obj["username"].(string)),
	}
}

func flattenContainerAwsClusterAuthorizationAdminUsersArray(objs []containeraws.ClusterAuthorizationAdminUsers) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenContainerAwsClusterAuthorizationAdminUsers(&item)
		items = append(items, i)
	}

	return items
}

func flattenContainerAwsClusterAuthorizationAdminUsers(obj *containeraws.ClusterAuthorizationAdminUsers) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"username": obj.Username,
	}

	return transformed

}

func expandContainerAwsClusterControlPlane(o interface{}) *containeraws.ClusterControlPlane {
	if o == nil {
		return containeraws.EmptyClusterControlPlane
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterControlPlane
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlane{
		AwsServicesAuthentication: expandContainerAwsClusterControlPlaneAwsServicesAuthentication(obj["aws_services_authentication"]),
		ConfigEncryption:          expandContainerAwsClusterControlPlaneConfigEncryption(obj["config_encryption"]),
		DatabaseEncryption:        expandContainerAwsClusterControlPlaneDatabaseEncryption(obj["database_encryption"]),
		IamInstanceProfile:        dcl.String(obj["iam_instance_profile"].(string)),
		SubnetIds:                 tpgdclresource.ExpandStringArray(obj["subnet_ids"]),
		Version:                   dcl.String(obj["version"].(string)),
		InstanceType:              dcl.StringOrNil(obj["instance_type"].(string)),
		MainVolume:                expandContainerAwsClusterControlPlaneMainVolume(obj["main_volume"]),
		ProxyConfig:               expandContainerAwsClusterControlPlaneProxyConfig(obj["proxy_config"]),
		RootVolume:                expandContainerAwsClusterControlPlaneRootVolume(obj["root_volume"]),
		SecurityGroupIds:          tpgdclresource.ExpandStringArray(obj["security_group_ids"]),
		SshConfig:                 expandContainerAwsClusterControlPlaneSshConfig(obj["ssh_config"]),
		Tags:                      tpgresource.CheckStringMap(obj["tags"]),
	}
}

func flattenContainerAwsClusterControlPlane(obj *containeraws.ClusterControlPlane) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"aws_services_authentication": flattenContainerAwsClusterControlPlaneAwsServicesAuthentication(obj.AwsServicesAuthentication),
		"config_encryption":           flattenContainerAwsClusterControlPlaneConfigEncryption(obj.ConfigEncryption),
		"database_encryption":         flattenContainerAwsClusterControlPlaneDatabaseEncryption(obj.DatabaseEncryption),
		"iam_instance_profile":        obj.IamInstanceProfile,
		"subnet_ids":                  obj.SubnetIds,
		"version":                     obj.Version,
		"instance_type":               obj.InstanceType,
		"main_volume":                 flattenContainerAwsClusterControlPlaneMainVolume(obj.MainVolume),
		"proxy_config":                flattenContainerAwsClusterControlPlaneProxyConfig(obj.ProxyConfig),
		"root_volume":                 flattenContainerAwsClusterControlPlaneRootVolume(obj.RootVolume),
		"security_group_ids":          obj.SecurityGroupIds,
		"ssh_config":                  flattenContainerAwsClusterControlPlaneSshConfig(obj.SshConfig),
		"tags":                        obj.Tags,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneAwsServicesAuthentication(o interface{}) *containeraws.ClusterControlPlaneAwsServicesAuthentication {
	if o == nil {
		return containeraws.EmptyClusterControlPlaneAwsServicesAuthentication
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterControlPlaneAwsServicesAuthentication
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneAwsServicesAuthentication{
		RoleArn:         dcl.String(obj["role_arn"].(string)),
		RoleSessionName: dcl.StringOrNil(obj["role_session_name"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneAwsServicesAuthentication(obj *containeraws.ClusterControlPlaneAwsServicesAuthentication) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"role_arn":          obj.RoleArn,
		"role_session_name": obj.RoleSessionName,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneConfigEncryption(o interface{}) *containeraws.ClusterControlPlaneConfigEncryption {
	if o == nil {
		return containeraws.EmptyClusterControlPlaneConfigEncryption
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterControlPlaneConfigEncryption
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneConfigEncryption{
		KmsKeyArn: dcl.String(obj["kms_key_arn"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneConfigEncryption(obj *containeraws.ClusterControlPlaneConfigEncryption) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"kms_key_arn": obj.KmsKeyArn,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneDatabaseEncryption(o interface{}) *containeraws.ClusterControlPlaneDatabaseEncryption {
	if o == nil {
		return containeraws.EmptyClusterControlPlaneDatabaseEncryption
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterControlPlaneDatabaseEncryption
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneDatabaseEncryption{
		KmsKeyArn: dcl.String(obj["kms_key_arn"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneDatabaseEncryption(obj *containeraws.ClusterControlPlaneDatabaseEncryption) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"kms_key_arn": obj.KmsKeyArn,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneMainVolume(o interface{}) *containeraws.ClusterControlPlaneMainVolume {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneMainVolume{
		Iops:       dcl.Int64OrNil(int64(obj["iops"].(int))),
		KmsKeyArn:  dcl.String(obj["kms_key_arn"].(string)),
		SizeGib:    dcl.Int64OrNil(int64(obj["size_gib"].(int))),
		Throughput: dcl.Int64OrNil(int64(obj["throughput"].(int))),
		VolumeType: containeraws.ClusterControlPlaneMainVolumeVolumeTypeEnumRef(obj["volume_type"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneMainVolume(obj *containeraws.ClusterControlPlaneMainVolume) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"iops":        obj.Iops,
		"kms_key_arn": obj.KmsKeyArn,
		"size_gib":    obj.SizeGib,
		"throughput":  obj.Throughput,
		"volume_type": obj.VolumeType,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneProxyConfig(o interface{}) *containeraws.ClusterControlPlaneProxyConfig {
	if o == nil {
		return containeraws.EmptyClusterControlPlaneProxyConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterControlPlaneProxyConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneProxyConfig{
		SecretArn:     dcl.String(obj["secret_arn"].(string)),
		SecretVersion: dcl.String(obj["secret_version"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneProxyConfig(obj *containeraws.ClusterControlPlaneProxyConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"secret_arn":     obj.SecretArn,
		"secret_version": obj.SecretVersion,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneRootVolume(o interface{}) *containeraws.ClusterControlPlaneRootVolume {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneRootVolume{
		Iops:       dcl.Int64OrNil(int64(obj["iops"].(int))),
		KmsKeyArn:  dcl.String(obj["kms_key_arn"].(string)),
		SizeGib:    dcl.Int64OrNil(int64(obj["size_gib"].(int))),
		Throughput: dcl.Int64OrNil(int64(obj["throughput"].(int))),
		VolumeType: containeraws.ClusterControlPlaneRootVolumeVolumeTypeEnumRef(obj["volume_type"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneRootVolume(obj *containeraws.ClusterControlPlaneRootVolume) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"iops":        obj.Iops,
		"kms_key_arn": obj.KmsKeyArn,
		"size_gib":    obj.SizeGib,
		"throughput":  obj.Throughput,
		"volume_type": obj.VolumeType,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterControlPlaneSshConfig(o interface{}) *containeraws.ClusterControlPlaneSshConfig {
	if o == nil {
		return containeraws.EmptyClusterControlPlaneSshConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterControlPlaneSshConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterControlPlaneSshConfig{
		Ec2KeyPair: dcl.String(obj["ec2_key_pair"].(string)),
	}
}

func flattenContainerAwsClusterControlPlaneSshConfig(obj *containeraws.ClusterControlPlaneSshConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"ec2_key_pair": obj.Ec2KeyPair,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterFleet(o interface{}) *containeraws.ClusterFleet {
	if o == nil {
		return containeraws.EmptyClusterFleet
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterFleet
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterFleet{
		Project: dcl.StringOrNil(obj["project"].(string)),
	}
}

func flattenContainerAwsClusterFleet(obj *containeraws.ClusterFleet) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"project":    obj.Project,
		"membership": obj.Membership,
	}

	return []interface{}{transformed}

}

func expandContainerAwsClusterNetworking(o interface{}) *containeraws.ClusterNetworking {
	if o == nil {
		return containeraws.EmptyClusterNetworking
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return containeraws.EmptyClusterNetworking
	}
	obj := objArr[0].(map[string]interface{})
	return &containeraws.ClusterNetworking{
		PodAddressCidrBlocks:       tpgdclresource.ExpandStringArray(obj["pod_address_cidr_blocks"]),
		ServiceAddressCidrBlocks:   tpgdclresource.ExpandStringArray(obj["service_address_cidr_blocks"]),
		VPCId:                      dcl.String(obj["vpc_id"].(string)),
		PerNodePoolSgRulesDisabled: dcl.Bool(obj["per_node_pool_sg_rules_disabled"].(bool)),
	}
}

func flattenContainerAwsClusterNetworking(obj *containeraws.ClusterNetworking) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"pod_address_cidr_blocks":         obj.PodAddressCidrBlocks,
		"service_address_cidr_blocks":     obj.ServiceAddressCidrBlocks,
		"vpc_id":                          obj.VPCId,
		"per_node_pool_sg_rules_disabled": obj.PerNodePoolSgRulesDisabled,
	}

	return []interface{}{transformed}

}

func flattenContainerAwsClusterWorkloadIdentityConfig(obj *containeraws.ClusterWorkloadIdentityConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"identity_provider": obj.IdentityProvider,
		"issuer_uri":        obj.IssuerUri,
		"workload_pool":     obj.WorkloadPool,
	}

	return []interface{}{transformed}

}
