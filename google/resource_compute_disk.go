package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	computeDiskUserRegexString = "^(?:https://www.googleapis.com/compute/v1/projects/)?([-_a-zA-Z0-9]*)/zones/([-_a-zA-Z0-9]*)/instances/([-_a-zA-Z0-9]*)$"
)

var (
	computeDiskUserRegex = regexp.MustCompile(computeDiskUserRegexString)
)

func resourceComputeDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeDiskCreate,
		Read:   resourceComputeDiskRead,
		Update: resourceComputeDiskUpdate,
		Delete: resourceComputeDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"disk_encryption_key_raw": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"disk_encryption_key_sha256": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"image": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: diskImageDiffSuppress,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: linkDiffSuppress,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pd-standard",
				ForceNew: true,
			},

			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"label_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	z, err := getZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Loading zone: %s", z)
	zone, err := config.clientCompute.Zones.Get(
		project, z).Do()
	if err != nil {
		return fmt.Errorf(
			"Error loading zone '%s': %s", z, err)
	}

	// Build the disk parameter
	disk := &compute.Disk{
		Name:   d.Get("name").(string),
		SizeGb: int64(d.Get("size").(int)),
	}

	// If we were given a source image, load that.
	if v, ok := d.GetOk("image"); ok {
		log.Printf("[DEBUG] Resolving image name: %s", v.(string))
		imageUrl, err := resolveImage(config, project, v.(string))
		if err != nil {
			return fmt.Errorf(
				"Error resolving image name '%s': %s",
				v.(string), err)
		}

		disk.SourceImage = imageUrl
		log.Printf("[DEBUG] Image name resolved to: %s", imageUrl)
	}

	if v, ok := d.GetOk("type"); ok {
		log.Printf("[DEBUG] Loading disk type: %s", v.(string))
		diskType, err := readDiskType(config, zone, project, v.(string))
		if err != nil {
			return fmt.Errorf(
				"Error loading disk type '%s': %s",
				v.(string), err)
		}

		disk.Type = diskType.SelfLink
	}

	if v, ok := d.GetOk("snapshot"); ok {
		snapshotName := v.(string)
		match, _ := regexp.MatchString("^https://www.googleapis.com/compute", snapshotName)
		if match {
			disk.SourceSnapshot = snapshotName
		} else {
			log.Printf("[DEBUG] Loading snapshot: %s", snapshotName)
			snapshotData, err := config.clientCompute.Snapshots.Get(
				project, snapshotName).Do()

			if err != nil {
				return fmt.Errorf(
					"Error loading snapshot '%s': %s",
					snapshotName, err)
			}
			disk.SourceSnapshot = snapshotData.SelfLink
		}
	}

	if v, ok := d.GetOk("disk_encryption_key_raw"); ok {
		disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{}
		disk.DiskEncryptionKey.RawKey = v.(string)
	}

	if _, ok := d.GetOk("labels"); ok {
		disk.Labels = expandLabels(d)
	}

	op, err := config.clientCompute.Disks.Insert(
		project, z, disk).Do()
	if err != nil {
		return fmt.Errorf("Error creating disk: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(disk.Name)

	err = computeOperationWaitTime(config.clientCompute, op, project, "Creating Disk", int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if err != nil {
		return err
	}
	return resourceComputeDiskRead(d, meta)
}

func resourceComputeDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	z, err := getZone(d, config)
	if err != nil {
		return err
	}
	d.Partial(true)
	if d.HasChange("size") {
		rb := &compute.DisksResizeRequest{
			SizeGb: int64(d.Get("size").(int)),
		}
		op, err := config.clientCompute.Disks.Resize(
			project, z, d.Id(), rb).Do()
		if err != nil {
			return fmt.Errorf("Error resizing disk: %s", err)
		}
		d.SetPartial("size")

		err = computeOperationWaitTime(config.clientCompute, op, project, "Resizing Disk", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if err != nil {
			return err
		}
	}

	if d.HasChange("labels") {
		zslr := compute.ZoneSetLabelsRequest{
			Labels:           expandLabels(d),
			LabelFingerprint: d.Get("label_fingerprint").(string),
		}
		op, err := config.clientCompute.Disks.SetLabels(
			project, z, d.Id(), &zslr).Do()
		if err != nil {
			return fmt.Errorf("Error when setting labels: %s", err)
		}
		d.SetPartial("labels")

		err = computeOperationWaitTime(config.clientCompute, op, project, "Setting labels on disk", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if err != nil {
			return err
		}
	}
	d.Partial(false)

	return resourceComputeDiskRead(d, meta)
}

func resourceComputeDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	getDisk := func(zone string) (interface{}, error) {
		return config.clientCompute.Disks.Get(project, zone, d.Id()).Do()
	}

	var disk *compute.Disk
	if zone, _ := getZone(d, config); zone != "" {
		disk, err = config.clientCompute.Disks.Get(
			project, zone, d.Id()).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Disk %q", d.Get("name").(string)))
		}
	} else {
		// If the resource was imported, the only info we have is the ID. Try to find the resource
		// by searching in the region of the project.
		var resource interface{}
		resource, err = getZonalResourceFromRegion(getDisk, region, config.clientCompute, project)

		if err != nil {
			return err
		}

		disk = resource.(*compute.Disk)
	}

	d.Set("name", disk.Name)
	d.Set("self_link", disk.SelfLink)
	d.Set("type", GetResourceNameFromSelfLink(disk.Type))
	d.Set("zone", GetResourceNameFromSelfLink(disk.Zone))
	d.Set("size", disk.SizeGb)
	d.Set("users", disk.Users)
	if disk.DiskEncryptionKey != nil && disk.DiskEncryptionKey.Sha256 != "" {
		d.Set("disk_encryption_key_sha256", disk.DiskEncryptionKey.Sha256)
	}

	d.Set("image", disk.SourceImage)
	d.Set("snapshot", disk.SourceSnapshot)
	d.Set("labels", disk.Labels)
	d.Set("label_fingerprint", disk.LabelFingerprint)
	d.Set("project", project)

	return nil
}

func resourceComputeDiskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	z, err := getZone(d, config)
	if err != nil {
		return err
	}

	// if disks are attached, they must be detached before the disk can be deleted
	if instances, ok := d.Get("users").([]interface{}); ok {
		type detachArgs struct{ project, zone, instance, deviceName string }
		var detachCalls []detachArgs
		self := d.Get("self_link").(string)
		for _, instance := range instances {
			if !computeDiskUserRegex.MatchString(instance.(string)) {
				return fmt.Errorf("Unknown user %q of disk %q", instance, self)
			}
			matches := computeDiskUserRegex.FindStringSubmatch(instance.(string))
			instanceProject := matches[1]
			instanceZone := matches[2]
			instanceName := matches[3]
			i, err := config.clientCompute.Instances.Get(instanceProject, instanceZone, instanceName).Do()
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					log.Printf("[WARN] instance %q not found, not bothering to detach disks", instance.(string))
					continue
				}
				return fmt.Errorf("Error retrieving instance %s: %s", instance.(string), err.Error())
			}
			for _, disk := range i.Disks {
				if disk.Source == self {
					detachCalls = append(detachCalls, detachArgs{
						project:    project,
						zone:       GetResourceNameFromSelfLink(i.Zone),
						instance:   i.Name,
						deviceName: disk.DeviceName,
					})
				}
			}
		}
		for _, call := range detachCalls {
			op, err := config.clientCompute.Instances.DetachDisk(call.project, call.zone, call.instance, call.deviceName).Do()
			if err != nil {
				return fmt.Errorf("Error detaching disk %s from instance %s/%s/%s: %s", call.deviceName, call.project,
					call.zone, call.instance, err.Error())
			}
			err = computeOperationWait(config.clientCompute, op, call.project,
				fmt.Sprintf("Detaching disk from %s/%s/%s", call.project, call.zone, call.instance))
			if err != nil {
				return err
			}
		}
	}

	// Delete the disk
	op, err := config.clientCompute.Disks.Delete(
		project, z, d.Id()).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing Disk %q because it's gone", d.Get("name").(string))
			// The resource doesn't exist anymore
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting disk: %s", err)
	}

	err = computeOperationWaitTime(config.clientCompute, op, project, "Deleting Disk", int(d.Timeout(schema.TimeoutDelete).Minutes()))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// We cannot suppress the diff for the case when family name is not part of the image name since we can't
// make a network call in a DiffSuppressFunc.
func diskImageDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	// 'old' is read from the API.
	// It always has the format 'https://www.googleapis.com/compute/v1/projects/(%s)/global/images/(%s)'
	matches := resolveImageLink.FindStringSubmatch(old)
	if matches == nil {
		// Image read from the API doesn't have the expected format. In practice, it should never happen
		return false
	}
	oldProject := matches[1]
	oldName := matches[2]

	// Partial or full self link family
	if resolveImageProjectFamily.MatchString(new) {
		// Value matches pattern "projects/{project}/global/images/family/{family-name}$"
		matches := resolveImageProjectFamily.FindStringSubmatch(new)
		newProject := matches[1]
		newFamilyName := matches[2]

		return diskImageProjectNameEquals(oldProject, newProject) && diskImageFamilyEquals(oldName, newFamilyName)
	}

	// Partial or full self link image
	if resolveImageProjectImage.MatchString(new) {
		// Value matches pattern "projects/{project}/global/images/{image-name}$"
		matches := resolveImageProjectImage.FindStringSubmatch(new)
		newProject := matches[1]
		newImageName := matches[2]

		return diskImageProjectNameEquals(oldProject, newProject) && diskImageEquals(oldName, newImageName)
	}

	// Partial link without project family
	if resolveImageGlobalFamily.MatchString(new) {
		// Value is "global/images/family/{family-name}"
		matches := resolveImageGlobalFamily.FindStringSubmatch(new)
		familyName := matches[1]

		return diskImageFamilyEquals(oldName, familyName)
	}

	// Partial link without project image
	if resolveImageGlobalImage.MatchString(new) {
		// Value is "global/images/{image-name}"
		matches := resolveImageGlobalImage.FindStringSubmatch(new)
		imageName := matches[1]

		return diskImageEquals(oldName, imageName)
	}

	// Family shorthand
	if resolveImageFamilyFamily.MatchString(new) {
		// Value is "family/{family-name}"
		matches := resolveImageFamilyFamily.FindStringSubmatch(new)
		familyName := matches[1]

		return diskImageFamilyEquals(oldName, familyName)
	}

	// Shorthand for image or family
	if resolveImageProjectImageShorthand.MatchString(new) {
		// Value is "{project}/{image-name}" or "{project}/{family-name}"
		matches := resolveImageProjectImageShorthand.FindStringSubmatch(new)
		newProject := matches[1]
		newName := matches[2]

		return diskImageProjectNameEquals(oldProject, newProject) &&
			(diskImageEquals(oldName, newName) || diskImageFamilyEquals(oldName, newName))
	}

	// Image or family only
	if diskImageEquals(oldName, new) || diskImageFamilyEquals(oldName, new) {
		// Value is "{image-name}" or "{family-name}"
		return true
	}

	return false
}

func diskImageProjectNameEquals(project1, project2 string) bool {
	// Convert short project name to full name
	// For instance, centos => centos-cloud
	fullProjectName, ok := imageMap[project2]
	if ok {
		project2 = fullProjectName
	}

	return project1 == project2
}

func diskImageEquals(oldImageName, newImageName string) bool {
	return oldImageName == newImageName
}

func diskImageFamilyEquals(imageName, familyName string) bool {
	// Handles the case when the image name includes the family name
	// e.g. image name: debian-9-drawfork-v20180109, family name: debian-9
	if strings.Contains(imageName, familyName) {
		return true
	}

	if suppressCanonicalFamilyDiff(imageName, familyName) {
		return true
	}

	if suppressWindowsSqlFamilyDiff(imageName, familyName) {
		return true
	}

	if suppressWindowsFamilyDiff(imageName, familyName) {
		return true
	}

	return false
}

// e.g. image: ubuntu-1404-trusty-v20180122, family: ubuntu-1404-lts
func suppressCanonicalFamilyDiff(imageName, familyName string) bool {
	parts := canonicalUbuntuLtsImage.FindStringSubmatch(imageName)
	if len(parts) == 2 {
		f := fmt.Sprintf("ubuntu-%s-lts", parts[1])
		if f == familyName {
			return true
		}
	}

	return false
}

// e.g. image: sql-2017-standard-windows-2016-dc-v20180109, family: sql-std-2017-win-2016
// e.g. image: sql-2017-express-windows-2012-r2-dc-v20180109, family: sql-exp-2017-win-2012-r2
func suppressWindowsSqlFamilyDiff(imageName, familyName string) bool {
	parts := windowsSqlImage.FindStringSubmatch(imageName)
	if len(parts) == 5 {
		edition := parts[2] // enterprise, standard or web.
		sqlVersion := parts[1]
		windowsVersion := parts[3]

		// Translate edition
		switch edition {
		case "enterprise":
			edition = "ent"
		case "standard":
			edition = "std"
		case "express":
			edition = "exp"
		}

		var f string
		if revision := parts[4]; revision != "" {
			// With revision
			f = fmt.Sprintf("sql-%s-%s-win-%s-r%s", edition, sqlVersion, windowsVersion, revision)
		} else {
			// No revision
			f = fmt.Sprintf("sql-%s-%s-win-%s", edition, sqlVersion, windowsVersion)
		}

		if f == familyName {
			return true
		}
	}

	return false
}

// e.g. image: windows-server-1709-dc-core-v20180109, family: windows-1709-core
// e.g. image: windows-server-1709-dc-core-for-containers-v20180109, family: "windows-1709-core-for-containers
func suppressWindowsFamilyDiff(imageName, familyName string) bool {
	updatedFamilyString := strings.Replace(familyName, "windows-", "windows-server-", 1)
	updatedFamilyString = strings.Replace(updatedFamilyString, "-core", "-dc-core", 1)

	if strings.Contains(imageName, updatedFamilyString) {
		return true
	}

	return false
}
