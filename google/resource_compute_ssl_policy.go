package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	compute "google.golang.org/api/compute/v1"
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

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"custom_features": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, // currently, the beta patch API call does not allow updating the description
			},

			"min_tls_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TLS_1_0",
				// Although compute-gen.go says that TLS_1_3 is a valid value, the API currently (26 Mar 2018)
				// responds with an HTTP 200 but doesn't actually create/update the policy. Open bug for this:
				// https://issuetracker.google.com/issues/76433946
				ValidateFunc: validation.StringInSlice([]string{"TLS_1_0", "TLS_1_1", "TLS_1_2"}, false),
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

			"enabled_features": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

	sslPolicy := &compute.SslPolicy{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Profile:        d.Get("profile").(string),
		MinTlsVersion:  d.Get("min_tls_version").(string),
		CustomFeatures: convertStringSet(d.Get("custom_features").(*schema.Set)),
	}

	op, err := config.clientCompute.SslPolicies.Insert(project, sslPolicy).Do()
	if err != nil {
		return fmt.Errorf("Error creating SSL Policy: %s", err)
	}

	d.SetId(sslPolicy.Name)

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), "Creating SSL Policy")
	if err != nil {
		d.SetId("") // if insert fails, remove from state
		return err
	}

	return resourceComputeSslPolicyRead(d, meta)
}

func resourceComputeSslPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Id()

	sslPolicy, err := config.clientCompute.SslPolicies.Get(project, name).Do()
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
	d.Set("enabled_features", convertStringArrToInterface(sslPolicy.EnabledFeatures))
	d.Set("custom_features", convertStringArrToInterface(sslPolicy.CustomFeatures))

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

	sslPolicy := &compute.SslPolicy{
		Fingerprint:   d.Get("fingerprint").(string),
		Profile:       d.Get("profile").(string),
		MinTlsVersion: d.Get("min_tls_version").(string),
	}

	if v, ok := d.Get("custom_features").(*schema.Set); ok {
		if v.Len() > 0 {
			sslPolicy.CustomFeatures = convertStringSet(v)
		} else {
			sslPolicy.NullFields = append(sslPolicy.NullFields, "CustomFeatures")
		}
	}

	op, err := config.clientCompute.SslPolicies.Patch(project, name, sslPolicy).Do()
	if err != nil {
		return fmt.Errorf("Error updating SSL Policy: %s", err)
	}

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Updating SSL Policy")
	if err != nil {
		return err
	}

	return resourceComputeSslPolicyRead(d, meta)
}

func resourceComputeSslPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	op, err := config.clientCompute.SslPolicies.Delete(project, name).Do()
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
