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
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceComputeImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeImageCreate,
		Read:   resourceComputeImageRead,
		Update: resourceComputeImageUpdate,
		Delete: resourceComputeImageDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeImageImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Name of the resource; provided by the client when the resource is
created. The name must be 1-63 characters long, and comply with
RFC1035. Specifically, the name must be 1-63 characters long and
match the regular expression '[a-z]([-a-z0-9]*[a-z0-9])?' which means
the first character must be a lowercase letter, and all following
characters must be a dash, lowercase letter, or digit, except the
last character, which cannot be a dash.`,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `An optional description of this resource. Provide this property when
you create the resource.`,
			},
			"disk_size_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `Size of the image when restored onto a persistent disk (in GB).`,
			},
			"family": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `The name of the image family to which this image belongs. You can
create disks by specifying an image family instead of a specific
image name. The image family always returns its latest image that is
not deprecated. The name of the image family must comply with
RFC1035.`,
			},
			"guest_os_features": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Description: `A list of features to enable on the guest operating system.
Applicable only for bootable images.`,
				Elem: computeImageGuestOsFeaturesSchema(),
				// Default schema.HashSchema is used.
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Labels to apply to this Image.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"licenses": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `Any applicable license URI.`,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},
			"raw_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The parameters of the raw disk image.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							Description: `The full Google Cloud Storage URL where disk storage is stored
You must provide either this property or the sourceDisk property
but not both.`,
						},
						"container_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"TAR", ""}, false),
							Description: `The format used to encode and transmit the block device, which
should be TAR. This is just a container and transmission format
and not a runtime format. Provided by the client when the disk
image is created. Default value: "TAR" Possible values: ["TAR"]`,
							Default: "TAR",
						},
						"sha1": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: `An optional SHA1 checksum of the disk image before unpackaging.
This is provided by the client when the disk image is created.`,
						},
					},
				},
			},
			"source_disk": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description: `The source disk to create this image based on.
You must provide either this property or the
rawDisk.source property but not both to create an image.`,
			},
			"source_image": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description: `URL of the source image used to create this image. In order to create an image, you must provide the full or partial
URL of one of the following:

The selfLink URL
This property
The rawDisk.source URL
The sourceDisk URL`,
			},
			"source_snapshot": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description: `URL of the source snapshot used to create this image.

In order to create an image, you must provide the full or partial URL of one of the following:

The selfLink URL
This property
The sourceImage URL
The rawDisk.source URL
The sourceDisk URL`,
			},
			"archive_size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `Size of the image tar.gz archive stored in Google Cloud Storage (in
bytes).`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Creation timestamp in RFC3339 text format.`,
			},
			"label_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The fingerprint used for optimistic locking of this resource. Used
internally during updates.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func computeImageGuestOsFeaturesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MULTI_IP_SUBNET", "SECURE_BOOT", "UEFI_COMPATIBLE", "VIRTIO_SCSI_MULTIQUEUE", "WINDOWS"}, false),
				Description:  `The type of supported feature. Read [Enabling guest operating system features](https://cloud.google.com/compute/docs/images/create-delete-deprecate-private-images#guest-os-features) to see a list of available options. Possible values: ["MULTI_IP_SUBNET", "SECURE_BOOT", "UEFI_COMPATIBLE", "VIRTIO_SCSI_MULTIQUEUE", "WINDOWS"]`,
			},
		},
	}
}

func resourceComputeImageCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeImageDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	diskSizeGbProp, err := expandComputeImageDiskSizeGb(d.Get("disk_size_gb"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disk_size_gb"); !isEmptyValue(reflect.ValueOf(diskSizeGbProp)) && (ok || !reflect.DeepEqual(v, diskSizeGbProp)) {
		obj["diskSizeGb"] = diskSizeGbProp
	}
	familyProp, err := expandComputeImageFamily(d.Get("family"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("family"); !isEmptyValue(reflect.ValueOf(familyProp)) && (ok || !reflect.DeepEqual(v, familyProp)) {
		obj["family"] = familyProp
	}
	guestOsFeaturesProp, err := expandComputeImageGuestOsFeatures(d.Get("guest_os_features"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("guest_os_features"); !isEmptyValue(reflect.ValueOf(guestOsFeaturesProp)) && (ok || !reflect.DeepEqual(v, guestOsFeaturesProp)) {
		obj["guestOsFeatures"] = guestOsFeaturesProp
	}
	labelsProp, err := expandComputeImageLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	labelFingerprintProp, err := expandComputeImageLabelFingerprint(d.Get("label_fingerprint"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("label_fingerprint"); !isEmptyValue(reflect.ValueOf(labelFingerprintProp)) && (ok || !reflect.DeepEqual(v, labelFingerprintProp)) {
		obj["labelFingerprint"] = labelFingerprintProp
	}
	licensesProp, err := expandComputeImageLicenses(d.Get("licenses"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("licenses"); !isEmptyValue(reflect.ValueOf(licensesProp)) && (ok || !reflect.DeepEqual(v, licensesProp)) {
		obj["licenses"] = licensesProp
	}
	nameProp, err := expandComputeImageName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	rawDiskProp, err := expandComputeImageRawDisk(d.Get("raw_disk"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("raw_disk"); !isEmptyValue(reflect.ValueOf(rawDiskProp)) && (ok || !reflect.DeepEqual(v, rawDiskProp)) {
		obj["rawDisk"] = rawDiskProp
	}
	sourceDiskProp, err := expandComputeImageSourceDisk(d.Get("source_disk"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_disk"); !isEmptyValue(reflect.ValueOf(sourceDiskProp)) && (ok || !reflect.DeepEqual(v, sourceDiskProp)) {
		obj["sourceDisk"] = sourceDiskProp
	}
	sourceImageProp, err := expandComputeImageSourceImage(d.Get("source_image"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_image"); !isEmptyValue(reflect.ValueOf(sourceImageProp)) && (ok || !reflect.DeepEqual(v, sourceImageProp)) {
		obj["sourceImage"] = sourceImageProp
	}
	sourceSnapshotProp, err := expandComputeImageSourceSnapshot(d.Get("source_snapshot"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_snapshot"); !isEmptyValue(reflect.ValueOf(sourceSnapshotProp)) && (ok || !reflect.DeepEqual(v, sourceSnapshotProp)) {
		obj["sourceSnapshot"] = sourceSnapshotProp
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/images")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Image: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Image: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/global/images/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(
		config, res, project, "Creating Image",
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Image: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Image %q: %#v", d.Id(), res)

	return resourceComputeImageRead(d, meta)
}

func resourceComputeImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/images/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeImage %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}

	if err := d.Set("archive_size_bytes", flattenComputeImageArchiveSizeBytes(res["archiveSizeBytes"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("creation_timestamp", flattenComputeImageCreationTimestamp(res["creationTimestamp"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("description", flattenComputeImageDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("disk_size_gb", flattenComputeImageDiskSizeGb(res["diskSizeGb"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("family", flattenComputeImageFamily(res["family"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("guest_os_features", flattenComputeImageGuestOsFeatures(res["guestOsFeatures"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("labels", flattenComputeImageLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("label_fingerprint", flattenComputeImageLabelFingerprint(res["labelFingerprint"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("licenses", flattenComputeImageLicenses(res["licenses"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("name", flattenComputeImageName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("source_disk", flattenComputeImageSourceDisk(res["sourceDisk"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("source_image", flattenComputeImageSourceImage(res["sourceImage"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("source_snapshot", flattenComputeImageSourceSnapshot(res["sourceSnapshot"], d, config)); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Image: %s", err)
	}

	return nil
}

func resourceComputeImageUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	d.Partial(true)

	if d.HasChange("labels") || d.HasChange("label_fingerprint") {
		obj := make(map[string]interface{})

		labelsProp, err := expandComputeImageLabels(d.Get("labels"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
			obj["labels"] = labelsProp
		}
		labelFingerprintProp, err := expandComputeImageLabelFingerprint(d.Get("label_fingerprint"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("label_fingerprint"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelFingerprintProp)) {
			obj["labelFingerprint"] = labelFingerprintProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/images/{{name}}/setLabels")
		if err != nil {
			return err
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Image %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Image %q: %#v", d.Id(), res)
		}

		err = computeOperationWaitTime(
			config, res, project, "Updating Image",
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceComputeImageRead(d, meta)
}

func resourceComputeImageDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/images/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Image %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Image")
	}

	err = computeOperationWaitTime(
		config, res, project, "Deleting Image",
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Image %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeImageImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/global/images/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/global/images/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeImageArchiveSizeBytes(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
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

func flattenComputeImageCreationTimestamp(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageDiskSizeGb(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
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

func flattenComputeImageFamily(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageGuestOsFeatures(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := schema.NewSet(schema.HashResource(computeImageGuestOsFeaturesSchema()), []interface{}{})
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed.Add(map[string]interface{}{
			"type": flattenComputeImageGuestOsFeaturesType(original["type"], d, config),
		})
	}
	return transformed
}
func flattenComputeImageGuestOsFeaturesType(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageLabelFingerprint(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageLicenses(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return convertAndMapStringArr(v.([]interface{}), ConvertSelfLinkToV1)
}

func flattenComputeImageName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeImageSourceDisk(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeImageSourceImage(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeImageSourceSnapshot(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func expandComputeImageDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageDiskSizeGb(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageFamily(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageGuestOsFeatures(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedType, err := expandComputeImageGuestOsFeaturesType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !isEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeImageGuestOsFeaturesType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandComputeImageLabelFingerprint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageLicenses(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			return nil, fmt.Errorf("Invalid value for licenses: nil")
		}
		f, err := parseGlobalFieldValue("licenses", raw.(string), "project", d, config, true)
		if err != nil {
			return nil, fmt.Errorf("Invalid value for licenses: %s", err)
		}
		req = append(req, f.RelativeLink())
	}
	return req, nil
}

func expandComputeImageName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageRawDisk(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedContainerType, err := expandComputeImageRawDiskContainerType(original["container_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedContainerType); val.IsValid() && !isEmptyValue(val) {
		transformed["containerType"] = transformedContainerType
	}

	transformedSha1, err := expandComputeImageRawDiskSha1(original["sha1"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSha1); val.IsValid() && !isEmptyValue(val) {
		transformed["sha1Checksum"] = transformedSha1
	}

	transformedSource, err := expandComputeImageRawDiskSource(original["source"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSource); val.IsValid() && !isEmptyValue(val) {
		transformed["source"] = transformedSource
	}

	return transformed, nil
}

func expandComputeImageRawDiskContainerType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageRawDiskSha1(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageRawDiskSource(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeImageSourceDisk(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseZonalFieldValue("disks", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for source_disk: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeImageSourceImage(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("images", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for source_image: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeImageSourceSnapshot(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("snapshots", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for source_snapshot: %s", err)
	}
	return f.RelativeLink(), nil
}
