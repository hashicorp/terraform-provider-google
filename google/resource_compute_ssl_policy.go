package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

func resourceComputeSslPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSslPolicyCreate,
		Read:   resourceComputeSslPolicyRead,
		Update: resourceComputeSslPolicyUpdate,
		Delete: resourceComputeSslPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"custom_features": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"min_tls_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TLS_1_0",
				ValidateFunc: validation.StringInSlice([]string{"TLS_1_0", "TLS_1_1", "TLS_1_2", "TLS_1_3"}, false),
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"profile": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "COMPATIBLE",
				ValidateFunc: validation.StringInSlice([]string{"COMPATIBLE", "MODERN", "RESTRICTED", "CUSTOM"}, false),
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			profile := diff.Get("profile")
			customFeaturesCount := diff.Get("custom_features.#")

			// Validate that policy configs aren't incompatible during all phases
			// CUSTOM profile demands non-zero custom_features, and other profiles (i.e., not CUSTOM) demand zero custom_features
			if diff.HasChange("profile") || diff.HasChange("custom_features") {
				if profile.(string) == "CUSTOM" {
					if customFeaturesCount.(int) == 0 {
						return fmt.Errorf("Error in SSL Policy %s: the profile is set to %s but no custom_features are set.", diff.Get("name"), profile.(string))
					}
				} else {
					if customFeaturesCount != 0 {
						return fmt.Errorf("Error in SSL Policy %s: the profile is set to %s but using custom_features requires the profile to be CUSTOM.", diff.Get("name"), profile.(string))
					}
				}
				return nil
			}
			return nil
		},
	}
}

func resourceComputeSslPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sslPolicy := &computeBeta.SslPolicy{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Profile:        d.Get("profile").(string),
		MinTlsVersion:  d.Get("min_tls_version").(string),
		CustomFeatures: convCustomFeaturesListToSlice(d.Get("custom_features").([]interface{})),
	}

	op, err := config.clientComputeBeta.SslPolicies.Insert(project, sslPolicy).Do()
	if err != nil {
		return fmt.Errorf("Error creating SSL Policy: %s", err)
	}

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutDelete).Minutes()), "Creating SSL Policy")
	if err != nil {
		return err
	}

	d.SetId(sslPolicy.Name)

	return resourceComputeSslPolicyRead(d, meta)
}

func resourceComputeSslPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	// handle resource inputs when only the resource id (name) is given
	if name == "" {
		name = d.Id()
	}

	sslPolicy, err := config.clientComputeBeta.SslPolicies.Get(project, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SSL Policy %q", name))
	}

	d.Set("name", sslPolicy.Name)
	d.Set("description", sslPolicy.Description)
	d.Set("min_tls_version", sslPolicy.MinTlsVersion)
	d.Set("profile", sslPolicy.Profile)
	d.Set("fingerprint", sslPolicy.Fingerprint)
	d.Set("project", project)
	d.Set("self_link", ConvertSelfLinkToV1(sslPolicy.SelfLink))

	if sslPolicy.CustomFeatures != nil {
		d.Set("custom_features", convertStringArrToInterface(sslPolicy.CustomFeatures))
	}

	d.SetId(sslPolicy.Name)

	return nil
}

func resourceComputeSslPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	d.Partial(true)

	if d.HasChange("min_tls_version") {
		sslPolicy := &computeBeta.SslPolicy{
			MinTlsVersion: d.Get("min_tls_version").(string),
			Fingerprint:   d.Get("fingerprint").(string),
		}

		op, err := config.clientComputeBeta.SslPolicies.Patch(project, name, sslPolicy).Do()
		if err != nil {
			return fmt.Errorf("Error updating SSL Policy: %s", err)
		}

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Updating SSL Policy")
		if err != nil {
			return err
		}

		d.SetPartial("min_tls_version")
	}

	if d.HasChange("profile") {
		sslPolicy := &computeBeta.SslPolicy{
			Profile:     d.Get("profile").(string),
			Fingerprint: d.Get("fingerprint").(string),
		}

		op, err := config.clientComputeBeta.SslPolicies.Patch(project, name, sslPolicy).Do()
		if err != nil {
			return fmt.Errorf("Error updating SSL Policy: %s", err)
		}

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Updating SSL Policy")
		if err != nil {
			return err
		}

		d.SetPartial("profile")
	}

	if d.HasChange("custom_features") {
		sslPolicy := &computeBeta.SslPolicy{
			CustomFeatures: convCustomFeaturesListToSlice(d.Get("custom_features").([]interface{})),
			Fingerprint:    d.Get("fingerprint").(string),
		}

		op, err := config.clientComputeBeta.SslPolicies.Patch(project, name, sslPolicy).Do()
		if err != nil {
			return fmt.Errorf("Error updating SSL Policy: %s", err)
		}

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Updating SSL Policy")
		if err != nil {
			return err
		}

		d.SetPartial("custom_features")
	}

	d.Partial(false)

	return resourceComputeSslPolicyRead(d, meta)
}

func resourceComputeSslPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	op, err := config.clientComputeBeta.SslPolicies.Delete(project, name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting SSL Policy: %s", err)
	}

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutDelete).Minutes()), "Deleting Subnetwork")
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func convCustomFeaturesListToSlice(customFeaturesList []interface{}) []string {
	customFeaturesSlice := make([]string, 0, len(customFeaturesList))

	for _, v := range customFeaturesList {
		customFeaturesSlice = append(customFeaturesSlice, v.(string))
	}

	return customFeaturesSlice
}
