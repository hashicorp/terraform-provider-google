// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	"log"
	"time"

	"google.golang.org/api/composer/v1"
)

func ResourceComposerUserWorkloadsSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceComposerUserWorkloadsSecretCreate,
		Read:   resourceComposerUserWorkloadsSecretRead,
		Update: resourceComposerUserWorkloadsSecretUpdate,
		Delete: resourceComposerUserWorkloadsSecretDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComposerUserWorkloadsSecretImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.DefaultProviderRegion,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateGCEName,
				Description:  `Name of the environment.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The location or Compute Engine region for the environment.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
			"environment": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateGCEName,
				Description:  `Name of the environment.`,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Sensitive:   true,
				Description: `A map of the secret data.`,
			},
		},
	}
}

func resourceComposerUserWorkloadsSecretCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretName, err := resourceComposerUserWorkloadsSecretName(d, config)
	if err != nil {
		return err
	}

	secret := &composer.UserWorkloadsSecret{
		Name: secretName.ResourceName(),
		Data: tpgresource.ConvertStringMap(d.Get("data").(map[string]interface{})),
	}

	log.Printf("[DEBUG] Creating new UserWorkloadsSecret %q", secretName.ParentName())
	resp, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.UserWorkloadsSecrets.Create(secretName.ParentName(), secret).Do()
	if err != nil {
		return fmt.Errorf("Error creating UserWorkloadsSecret: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	respJson, _ := resp.MarshalJSON()
	log.Printf("[DEBUG] Finished creating UserWorkloadsSecret %q: %#v", d.Id(), string(respJson))

	return resourceComposerUserWorkloadsSecretRead(d, meta)
}

func resourceComposerUserWorkloadsSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretName, err := resourceComposerUserWorkloadsSecretName(d, config)
	if err != nil {
		return err
	}

	res, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.UserWorkloadsSecrets.Get(secretName.ResourceName()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("UserWorkloadsSecret %q", d.Id()))
	}

	if err := d.Set("project", secretName.Project); err != nil {
		return fmt.Errorf("Error setting UserWorkloadsSecret Project: %s", err)
	}
	if err := d.Set("region", secretName.Region); err != nil {
		return fmt.Errorf("Error setting UserWorkloadsSecret Region: %s", err)
	}
	if err := d.Set("environment", secretName.Environment); err != nil {
		return fmt.Errorf("Error setting UserWorkloadsSecret Environment: %s", err)
	}
	if err := d.Set("name", tpgresource.GetResourceNameFromSelfLink(res.Name)); err != nil {
		return fmt.Errorf("Error setting UserWorkloadsSecret Name: %s", err)
	}
	return nil
}

func resourceComposerUserWorkloadsSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretName, err := resourceComposerUserWorkloadsSecretName(d, config)
	if err != nil {
		return err
	}

	if d.HasChange("data") {
		secret := &composer.UserWorkloadsSecret{
			Name: secretName.ResourceName(),
			Data: tpgresource.ConvertStringMap(d.Get("data").(map[string]interface{})),
		}

		secretJson, _ := secret.MarshalJSON()
		log.Printf("[DEBUG] Updating UserWorkloadsSecret %q: %s", d.Id(), string(secretJson))

		resp, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.UserWorkloadsSecrets.Update(secretName.ResourceName(), secret).Do()
		if err != nil {
			return err
		}

		respJson, _ := resp.MarshalJSON()
		log.Printf("[DEBUG] Finished updating UserWorkloadsSecret %q: %s", d.Id(), string(respJson))
	}

	return nil
}

func resourceComposerUserWorkloadsSecretDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretName, err := resourceComposerUserWorkloadsSecretName(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting UserWorkloadsSecret %q", d.Id())
	_, err = config.NewComposerClient(userAgent).Projects.Locations.Environments.UserWorkloadsSecrets.Delete(secretName.ResourceName()).Do()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Finished deleting UserWorkloadsSecret %q", d.Id())

	return nil
}

func resourceComposerUserWorkloadsSecretImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{"projects/(?P<project>[^/]+)/locations/(?P<region>[^/]+)/environments/(?P<environment>[^/]+)/userWorkloadsSecrets/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<environment>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// retrieve "data" in advance, because Read function won't do it.
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	res, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.UserWorkloadsSecrets.Get(id).Do()
	if err != nil {
		return nil, err
	}

	if err := d.Set("data", res.Data); err != nil {
		return nil, fmt.Errorf("Error setting UserWorkloadsSecret Data: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceComposerUserWorkloadsSecretName(d *schema.ResourceData, config *transport_tpg.Config) (*UserWorkloadsSecretsName, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return nil, err
	}

	return &UserWorkloadsSecretsName{
		Project:     project,
		Region:      region,
		Environment: d.Get("environment").(string),
		Secret:      d.Get("name").(string),
	}, nil
}

type UserWorkloadsSecretsName struct {
	Project     string
	Region      string
	Environment string
	Secret      string
}

func (n *UserWorkloadsSecretsName) ResourceName() string {
	return fmt.Sprintf("projects/%s/locations/%s/environments/%s/userWorkloadsSecrets/%s", n.Project, n.Region, n.Environment, n.Secret)
}

func (n *UserWorkloadsSecretsName) ParentName() string {
	return fmt.Sprintf("projects/%s/locations/%s/environments/%s", n.Project, n.Region, n.Environment)
}
