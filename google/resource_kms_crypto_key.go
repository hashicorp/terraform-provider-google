// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceKmsCryptoKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsCryptoKeyCreate,
		Read:   resourceKmsCryptoKeyRead,
		Update: resourceKmsCryptoKeyUpdate,
		Delete: resourceKmsCryptoKeyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceKmsCryptoKeyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Update: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"key_ring": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: kmsCryptoKeyRingsEquivalent,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"purpose": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENCRYPT_DECRYPT", "ASYMMETRIC_SIGN", "ASYMMETRIC_DECRYPT", ""}, false),
				Default:      "ENCRYPT_DECRYPT",
			},
			"rotation_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: orEmpty(validateKmsCryptoKeyRotationPeriod),
			},
			"version_template": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protection_level": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"SOFTWARE", "HSM", ""}, false),
							Default:      "SOFTWARE",
						},
					},
				},
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKmsCryptoKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	purposeProp, err := expandKmsCryptoKeyPurpose(d.Get("purpose"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("purpose"); !isEmptyValue(reflect.ValueOf(purposeProp)) && (ok || !reflect.DeepEqual(v, purposeProp)) {
		obj["purpose"] = purposeProp
	}
	rotationPeriodProp, err := expandKmsCryptoKeyRotationPeriod(d.Get("rotation_period"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rotation_period"); !isEmptyValue(reflect.ValueOf(rotationPeriodProp)) && (ok || !reflect.DeepEqual(v, rotationPeriodProp)) {
		obj["rotationPeriod"] = rotationPeriodProp
	}
	versionTemplateProp, err := expandKmsCryptoKeyVersionTemplate(d.Get("version_template"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("version_template"); !isEmptyValue(reflect.ValueOf(versionTemplateProp)) && (ok || !reflect.DeepEqual(v, versionTemplateProp)) {
		obj["versionTemplate"] = versionTemplateProp
	}

	obj, err = resourceKmsCryptoKeyEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{KmsBasePath}}projects/{{project}}/keyRings/{{key_ring}}/cryptoKeys?cryptoKeyId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new CryptoKey: %#v", obj)
	res, err := sendRequestWithTimeout(config, "POST", url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating CryptoKey: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating CryptoKey %q: %#v", d.Id(), res)

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{KmsBasePath}}projects/{{projects}}/keyRings/{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("KmsCryptoKey %q", d.Id()))
	}

	res, err = resourceKmsCryptoKeyDecoder(d, meta, res)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}

	if err := d.Set("purpose", flattenKmsCryptoKeyPurpose(res["purpose"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	if err := d.Set("rotation_period", flattenKmsCryptoKeyRotationPeriod(res["rotationPeriod"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	if err := d.Set("version_template", flattenKmsCryptoKeyVersionTemplate(res["versionTemplate"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}

	return nil
}

func resourceKmsCryptoKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	rotationPeriodProp, err := expandKmsCryptoKeyRotationPeriod(d.Get("rotation_period"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rotation_period"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, rotationPeriodProp)) {
		obj["rotationPeriod"] = rotationPeriodProp
	}
	versionTemplateProp, err := expandKmsCryptoKeyVersionTemplate(d.Get("version_template"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("version_template"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, versionTemplateProp)) {
		obj["versionTemplate"] = versionTemplateProp
	}

	obj, err = resourceKmsCryptoKeyUpdateEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{KmsBasePath}}projects/{{projects}}/keyRings/{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating CryptoKey %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("rotation_period") {
		updateMask = append(updateMask, "rotationPeriod,nextRotationTime")
	}

	if d.HasChange("version_template") {
		updateMask = append(updateMask, "versionTemplate.algorithm")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
	_, err = sendRequestWithTimeout(config, "PATCH", url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating CryptoKey %q: %s", d.Id(), err)
	}

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf(`
[WARNING] KMS CryptoKey resources cannot be deleted from GCP. The CryptoKey %s will be removed from Terraform state,
and all its CryptoKeyVersions will be destroyed, but it will still be present on the server.`, cryptoKeyId.cryptoKeyId())

	// Delete all versions of the key
	if err := clearCryptoKeyVersions(cryptoKeyId, config); err != nil {
		return err
	}

	// Make sure automatic key rotation is disabled if set
	if d.Get("rotation_period") != "" {
		if err := disableCryptoKeyRotation(cryptoKeyId, config); err != nil {
			return fmt.Errorf(
				"While cryptoKeyVersions were cleared, Terraform was unable to disable automatic rotation of key due to an error: %s."+
					"Please retry or manually disable automatic rotation to prevent creation of a new version of this key.", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceKmsCryptoKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return nil, err
	}

	d.Set("key_ring", cryptoKeyId.KeyRingId.keyRingId())
	d.Set("name", cryptoKeyId.Name)

	return []*schema.ResourceData{d}, nil
}

func flattenKmsCryptoKeyPurpose(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyRotationPeriod(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionTemplate(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["algorithm"] =
		flattenKmsCryptoKeyVersionTemplateAlgorithm(original["algorithm"], d)
	transformed["protection_level"] =
		flattenKmsCryptoKeyVersionTemplateProtectionLevel(original["protectionLevel"], d)
	return []interface{}{transformed}
}
func flattenKmsCryptoKeyVersionTemplateAlgorithm(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionTemplateProtectionLevel(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandKmsCryptoKeyPurpose(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKmsCryptoKeyRotationPeriod(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKmsCryptoKeyVersionTemplate(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAlgorithm, err := expandKmsCryptoKeyVersionTemplateAlgorithm(original["algorithm"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAlgorithm); val.IsValid() && !isEmptyValue(val) {
		transformed["algorithm"] = transformedAlgorithm
	}

	transformedProtectionLevel, err := expandKmsCryptoKeyVersionTemplateProtectionLevel(original["protection_level"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProtectionLevel); val.IsValid() && !isEmptyValue(val) {
		transformed["protectionLevel"] = transformedProtectionLevel
	}

	return transformed, nil
}

func expandKmsCryptoKeyVersionTemplateAlgorithm(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKmsCryptoKeyVersionTemplateProtectionLevel(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceKmsCryptoKeyEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// if rotationPeriod is set, nextRotationTime must also be set.
	if d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return nil, fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		obj["nextRotationTime"] = nextRotation
	}

	return obj, nil
}

func resourceKmsCryptoKeyUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// if rotationPeriod is changed, nextRotationTime must also be set.
	if d.HasChange("rotation_period") && d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return nil, fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		obj["nextRotationTime"] = nextRotation
	}

	return obj, nil
}

func resourceKmsCryptoKeyDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	// Take the returned long form of the name and use it as `self_link`.
	// Then modify the name to be the user specified form.
	// We can't just ignore_read on `name` as the linter will
	// complain that the returned `res` is never used afterwards.
	// Some field needs to be actually set, and we chose `name`.
	d.Set("self_link", res["name"].(string))
	res["name"] = d.Get("name").(string)
	return res, nil
}
