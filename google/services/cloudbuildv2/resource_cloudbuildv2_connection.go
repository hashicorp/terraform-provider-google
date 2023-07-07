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

package cloudbuildv2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	cloudbuildv2 "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuildv2"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceCloudbuildv2Connection() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudbuildv2ConnectionCreate,
		Read:   resourceCloudbuildv2ConnectionRead,
		Update: resourceCloudbuildv2ConnectionUpdate,
		Delete: resourceCloudbuildv2ConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudbuildv2ConnectionImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
				Description: "Immutable. The resource name of the connection, in the format `projects/{project}/locations/{location}/connections/{connection_id}`.",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Allows clients to store small amounts of arbitrary data.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If disabled is set to true, functionality is disabled for this connection. Repository based API methods and webhooks processing for repositories in this connection will be disabled.",
			},

			"github_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Configuration for connections to github.com.",
				MaxItems:      1,
				Elem:          Cloudbuildv2ConnectionGithubConfigSchema(),
				ConflictsWith: []string{"github_enterprise_config", "gitlab_config"},
			},

			"github_enterprise_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Configuration for connections to an instance of GitHub Enterprise.",
				MaxItems:      1,
				Elem:          Cloudbuildv2ConnectionGithubEnterpriseConfigSchema(),
				ConflictsWith: []string{"github_config", "gitlab_config"},
			},

			"gitlab_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Configuration for connections to gitlab.com or an instance of GitLab Enterprise.",
				MaxItems:      1,
				Elem:          Cloudbuildv2ConnectionGitlabConfigSchema(),
				ConflictsWith: []string{"github_config", "github_enterprise_config"},
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
				Description: "Output only. Server assigned timestamp for when the connection was created.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"installation_state": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Installation state of the Connection.",
				Elem:        Cloudbuildv2ConnectionInstallationStateSchema(),
			},

			"reconciling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Set to true when the connection is being set up or updated in the background.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Server assigned timestamp for when the connection was updated.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGithubConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app_installation_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "GitHub App installation id.",
			},

			"authorizer_credential": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "OAuth credential of the account that authorized the Cloud Build GitHub App. It is recommended to use a robot account instead of a human user account. The OAuth token must be tied to the Cloud Build GitHub App.",
				MaxItems:    1,
				Elem:        Cloudbuildv2ConnectionGithubConfigAuthorizerCredentialSchema(),
			},
		},
	}
}

func Cloudbuildv2ConnectionGithubConfigAuthorizerCredentialSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"oauth_token_secret_version": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "A SecretManager resource containing the OAuth token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.",
			},

			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The username associated to this token.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGithubEnterpriseConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The URI of the GitHub Enterprise host this connection is for.",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Id of the GitHub App created from the manifest.",
			},

			"app_installation_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the installation of the GitHub App.",
			},

			"app_slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL-friendly name of the GitHub App.",
			},

			"private_key_secret_version": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "SecretManager resource containing the private key of the GitHub App, formatted as `projects/*/secrets/*/versions/*`.",
			},

			"service_directory_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration for using Service Directory to privately connect to a GitHub Enterprise server. This should only be set if the GitHub Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitHub Enterprise server will be made over the public internet.",
				MaxItems:    1,
				Elem:        Cloudbuildv2ConnectionGithubEnterpriseConfigServiceDirectoryConfigSchema(),
			},

			"ssl_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL certificate to use for requests to GitHub Enterprise.",
			},

			"webhook_secret_secret_version": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "SecretManager resource containing the webhook secret of the GitHub App, formatted as `projects/*/secrets/*/versions/*`.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGithubEnterpriseConfigServiceDirectoryConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. The Service Directory service name. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGitlabConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authorizer_credential": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. A GitLab personal access token with the `api` scope access.",
				MaxItems:    1,
				Elem:        Cloudbuildv2ConnectionGitlabConfigAuthorizerCredentialSchema(),
			},

			"read_authorizer_credential": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. A GitLab personal access token with the minimum `read_api` scope access.",
				MaxItems:    1,
				Elem:        Cloudbuildv2ConnectionGitlabConfigReadAuthorizerCredentialSchema(),
			},

			"webhook_secret_secret_version": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. Immutable. SecretManager resource containing the webhook secret of a GitLab Enterprise project, formatted as `projects/*/secrets/*/versions/*`.",
			},

			"host_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The URI of the GitLab Enterprise host this connection is for. If not specified, the default value is https://gitlab.com.",
			},

			"service_directory_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration for using Service Directory to privately connect to a GitLab Enterprise server. This should only be set if the GitLab Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitLab Enterprise server will be made over the public internet.",
				MaxItems:    1,
				Elem:        Cloudbuildv2ConnectionGitlabConfigServiceDirectoryConfigSchema(),
			},

			"ssl_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL certificate to use for requests to GitLab Enterprise.",
			},

			"server_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Version of the GitLab Enterprise server running on the `host_uri`.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGitlabConfigAuthorizerCredentialSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_token_secret_version": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. A SecretManager resource containing the user token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.",
			},

			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The username associated to this token.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGitlabConfigReadAuthorizerCredentialSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_token_secret_version": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. A SecretManager resource containing the user token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.",
			},

			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The username associated to this token.",
			},
		},
	}
}

func Cloudbuildv2ConnectionGitlabConfigServiceDirectoryConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. The Service Directory service name. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.",
			},
		},
	}
}

func Cloudbuildv2ConnectionInstallationStateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Link to follow for next action. Empty string if the installation is already complete.",
			},

			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Message of what the user should do next to continue the installation. Empty string if the installation is already complete.",
			},

			"stage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Current step of the installation process. Possible values: STAGE_UNSPECIFIED, PENDING_CREATE_APP, PENDING_USER_OAUTH, PENDING_INSTALL_APP, COMPLETE",
			},
		},
	}
}

func resourceCloudbuildv2ConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuildv2.Connection{
		Location:               dcl.String(d.Get("location").(string)),
		Name:                   dcl.String(d.Get("name").(string)),
		Annotations:            tpgresource.CheckStringMap(d.Get("annotations")),
		Disabled:               dcl.Bool(d.Get("disabled").(bool)),
		GithubConfig:           expandCloudbuildv2ConnectionGithubConfig(d.Get("github_config")),
		GithubEnterpriseConfig: expandCloudbuildv2ConnectionGithubEnterpriseConfig(d.Get("github_enterprise_config")),
		GitlabConfig:           expandCloudbuildv2ConnectionGitlabConfig(d.Get("gitlab_config")),
		Project:                dcl.String(project),
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
	client := transport_tpg.NewDCLCloudbuildv2Client(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyConnection(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Connection: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Connection %q: %#v", d.Id(), res)

	return resourceCloudbuildv2ConnectionRead(d, meta)
}

func resourceCloudbuildv2ConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuildv2.Connection{
		Location:               dcl.String(d.Get("location").(string)),
		Name:                   dcl.String(d.Get("name").(string)),
		Annotations:            tpgresource.CheckStringMap(d.Get("annotations")),
		Disabled:               dcl.Bool(d.Get("disabled").(bool)),
		GithubConfig:           expandCloudbuildv2ConnectionGithubConfig(d.Get("github_config")),
		GithubEnterpriseConfig: expandCloudbuildv2ConnectionGithubEnterpriseConfig(d.Get("github_enterprise_config")),
		GitlabConfig:           expandCloudbuildv2ConnectionGitlabConfig(d.Get("gitlab_config")),
		Project:                dcl.String(project),
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
	client := transport_tpg.NewDCLCloudbuildv2Client(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetConnection(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("Cloudbuildv2Connection %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("disabled", res.Disabled); err != nil {
		return fmt.Errorf("error setting disabled in state: %s", err)
	}
	if err = d.Set("github_config", flattenCloudbuildv2ConnectionGithubConfig(res.GithubConfig)); err != nil {
		return fmt.Errorf("error setting github_config in state: %s", err)
	}
	if err = d.Set("github_enterprise_config", flattenCloudbuildv2ConnectionGithubEnterpriseConfig(res.GithubEnterpriseConfig)); err != nil {
		return fmt.Errorf("error setting github_enterprise_config in state: %s", err)
	}
	if err = d.Set("gitlab_config", flattenCloudbuildv2ConnectionGitlabConfig(res.GitlabConfig)); err != nil {
		return fmt.Errorf("error setting gitlab_config in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("installation_state", flattenCloudbuildv2ConnectionInstallationState(res.InstallationState)); err != nil {
		return fmt.Errorf("error setting installation_state in state: %s", err)
	}
	if err = d.Set("reconciling", res.Reconciling); err != nil {
		return fmt.Errorf("error setting reconciling in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceCloudbuildv2ConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuildv2.Connection{
		Location:               dcl.String(d.Get("location").(string)),
		Name:                   dcl.String(d.Get("name").(string)),
		Annotations:            tpgresource.CheckStringMap(d.Get("annotations")),
		Disabled:               dcl.Bool(d.Get("disabled").(bool)),
		GithubConfig:           expandCloudbuildv2ConnectionGithubConfig(d.Get("github_config")),
		GithubEnterpriseConfig: expandCloudbuildv2ConnectionGithubEnterpriseConfig(d.Get("github_enterprise_config")),
		GitlabConfig:           expandCloudbuildv2ConnectionGitlabConfig(d.Get("gitlab_config")),
		Project:                dcl.String(project),
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
	client := transport_tpg.NewDCLCloudbuildv2Client(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyConnection(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Connection: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Connection %q: %#v", d.Id(), res)

	return resourceCloudbuildv2ConnectionRead(d, meta)
}

func resourceCloudbuildv2ConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &cloudbuildv2.Connection{
		Location:               dcl.String(d.Get("location").(string)),
		Name:                   dcl.String(d.Get("name").(string)),
		Annotations:            tpgresource.CheckStringMap(d.Get("annotations")),
		Disabled:               dcl.Bool(d.Get("disabled").(bool)),
		GithubConfig:           expandCloudbuildv2ConnectionGithubConfig(d.Get("github_config")),
		GithubEnterpriseConfig: expandCloudbuildv2ConnectionGithubEnterpriseConfig(d.Get("github_enterprise_config")),
		GitlabConfig:           expandCloudbuildv2ConnectionGitlabConfig(d.Get("gitlab_config")),
		Project:                dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Connection %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLCloudbuildv2Client(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteConnection(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Connection: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Connection %q", d.Id())
	return nil
}

func resourceCloudbuildv2ConnectionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/connections/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/connections/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandCloudbuildv2ConnectionGithubConfig(o interface{}) *cloudbuildv2.ConnectionGithubConfig {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGithubConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGithubConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGithubConfig{
		AppInstallationId:    dcl.Int64(int64(obj["app_installation_id"].(int))),
		AuthorizerCredential: expandCloudbuildv2ConnectionGithubConfigAuthorizerCredential(obj["authorizer_credential"]),
	}
}

func flattenCloudbuildv2ConnectionGithubConfig(obj *cloudbuildv2.ConnectionGithubConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"app_installation_id":   obj.AppInstallationId,
		"authorizer_credential": flattenCloudbuildv2ConnectionGithubConfigAuthorizerCredential(obj.AuthorizerCredential),
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGithubConfigAuthorizerCredential(o interface{}) *cloudbuildv2.ConnectionGithubConfigAuthorizerCredential {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGithubConfigAuthorizerCredential
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGithubConfigAuthorizerCredential
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGithubConfigAuthorizerCredential{
		OAuthTokenSecretVersion: dcl.String(obj["oauth_token_secret_version"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGithubConfigAuthorizerCredential(obj *cloudbuildv2.ConnectionGithubConfigAuthorizerCredential) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"oauth_token_secret_version": obj.OAuthTokenSecretVersion,
		"username":                   obj.Username,
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGithubEnterpriseConfig(o interface{}) *cloudbuildv2.ConnectionGithubEnterpriseConfig {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGithubEnterpriseConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGithubEnterpriseConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGithubEnterpriseConfig{
		HostUri:                    dcl.String(obj["host_uri"].(string)),
		AppId:                      dcl.Int64(int64(obj["app_id"].(int))),
		AppInstallationId:          dcl.Int64(int64(obj["app_installation_id"].(int))),
		AppSlug:                    dcl.String(obj["app_slug"].(string)),
		PrivateKeySecretVersion:    dcl.String(obj["private_key_secret_version"].(string)),
		ServiceDirectoryConfig:     expandCloudbuildv2ConnectionGithubEnterpriseConfigServiceDirectoryConfig(obj["service_directory_config"]),
		SslCa:                      dcl.String(obj["ssl_ca"].(string)),
		WebhookSecretSecretVersion: dcl.String(obj["webhook_secret_secret_version"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGithubEnterpriseConfig(obj *cloudbuildv2.ConnectionGithubEnterpriseConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"host_uri":                      obj.HostUri,
		"app_id":                        obj.AppId,
		"app_installation_id":           obj.AppInstallationId,
		"app_slug":                      obj.AppSlug,
		"private_key_secret_version":    obj.PrivateKeySecretVersion,
		"service_directory_config":      flattenCloudbuildv2ConnectionGithubEnterpriseConfigServiceDirectoryConfig(obj.ServiceDirectoryConfig),
		"ssl_ca":                        obj.SslCa,
		"webhook_secret_secret_version": obj.WebhookSecretSecretVersion,
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGithubEnterpriseConfigServiceDirectoryConfig(o interface{}) *cloudbuildv2.ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGithubEnterpriseConfigServiceDirectoryConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGithubEnterpriseConfigServiceDirectoryConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGithubEnterpriseConfigServiceDirectoryConfig{
		Service: dcl.String(obj["service"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGithubEnterpriseConfigServiceDirectoryConfig(obj *cloudbuildv2.ConnectionGithubEnterpriseConfigServiceDirectoryConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"service": obj.Service,
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGitlabConfig(o interface{}) *cloudbuildv2.ConnectionGitlabConfig {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGitlabConfig{
		AuthorizerCredential:       expandCloudbuildv2ConnectionGitlabConfigAuthorizerCredential(obj["authorizer_credential"]),
		ReadAuthorizerCredential:   expandCloudbuildv2ConnectionGitlabConfigReadAuthorizerCredential(obj["read_authorizer_credential"]),
		WebhookSecretSecretVersion: dcl.String(obj["webhook_secret_secret_version"].(string)),
		HostUri:                    dcl.StringOrNil(obj["host_uri"].(string)),
		ServiceDirectoryConfig:     expandCloudbuildv2ConnectionGitlabConfigServiceDirectoryConfig(obj["service_directory_config"]),
		SslCa:                      dcl.String(obj["ssl_ca"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGitlabConfig(obj *cloudbuildv2.ConnectionGitlabConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"authorizer_credential":         flattenCloudbuildv2ConnectionGitlabConfigAuthorizerCredential(obj.AuthorizerCredential),
		"read_authorizer_credential":    flattenCloudbuildv2ConnectionGitlabConfigReadAuthorizerCredential(obj.ReadAuthorizerCredential),
		"webhook_secret_secret_version": obj.WebhookSecretSecretVersion,
		"host_uri":                      obj.HostUri,
		"service_directory_config":      flattenCloudbuildv2ConnectionGitlabConfigServiceDirectoryConfig(obj.ServiceDirectoryConfig),
		"ssl_ca":                        obj.SslCa,
		"server_version":                obj.ServerVersion,
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGitlabConfigAuthorizerCredential(o interface{}) *cloudbuildv2.ConnectionGitlabConfigAuthorizerCredential {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfigAuthorizerCredential
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfigAuthorizerCredential
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGitlabConfigAuthorizerCredential{
		UserTokenSecretVersion: dcl.String(obj["user_token_secret_version"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGitlabConfigAuthorizerCredential(obj *cloudbuildv2.ConnectionGitlabConfigAuthorizerCredential) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"user_token_secret_version": obj.UserTokenSecretVersion,
		"username":                  obj.Username,
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGitlabConfigReadAuthorizerCredential(o interface{}) *cloudbuildv2.ConnectionGitlabConfigReadAuthorizerCredential {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfigReadAuthorizerCredential
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfigReadAuthorizerCredential
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGitlabConfigReadAuthorizerCredential{
		UserTokenSecretVersion: dcl.String(obj["user_token_secret_version"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGitlabConfigReadAuthorizerCredential(obj *cloudbuildv2.ConnectionGitlabConfigReadAuthorizerCredential) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"user_token_secret_version": obj.UserTokenSecretVersion,
		"username":                  obj.Username,
	}

	return []interface{}{transformed}

}

func expandCloudbuildv2ConnectionGitlabConfigServiceDirectoryConfig(o interface{}) *cloudbuildv2.ConnectionGitlabConfigServiceDirectoryConfig {
	if o == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfigServiceDirectoryConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return cloudbuildv2.EmptyConnectionGitlabConfigServiceDirectoryConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &cloudbuildv2.ConnectionGitlabConfigServiceDirectoryConfig{
		Service: dcl.String(obj["service"].(string)),
	}
}

func flattenCloudbuildv2ConnectionGitlabConfigServiceDirectoryConfig(obj *cloudbuildv2.ConnectionGitlabConfigServiceDirectoryConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"service": obj.Service,
	}

	return []interface{}{transformed}

}

func flattenCloudbuildv2ConnectionInstallationState(obj *cloudbuildv2.ConnectionInstallationState) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"action_uri": obj.ActionUri,
		"message":    obj.Message,
		"stage":      obj.Stage,
	}

	return []interface{}{transformed}

}
