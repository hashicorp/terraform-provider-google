package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudkms/v1"
	"log"
)

func resourceKmsKeyRing() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsKeyRingCreate,
		Read:   resourceKmsKeyRingRead,
		Delete: resourceKmsKeyRingDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func kmsResourceParentString(project, location string) string {
	return fmt.Sprintf("projects/%s/locations/%s", project, location)
}

func kmsResourceParentKeyRingName(project, location, name string) string {
	return fmt.Sprintf("%s/keyRings/%s", kmsResourceParentString(project, location), name)
}

func resourceKmsKeyRingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	location := d.Get("location").(string)

	parent := kmsResourceParentString(project, location)

	keyRing, err := config.clientKms.Projects.Locations.KeyRings.
		Create(parent, &cloudkms.KeyRing{}).
		KeyRingId(name).
		Do()

	if err != nil {
		return fmt.Errorf("Error creating KeyRing: %s", err)
	}

	log.Printf("[DEBUG] Created KeyRing %s", keyRing.Name)

	d.SetId(keyRing.Name)

	return resourceKmsKeyRingRead(d, meta)
}

func resourceKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyRingName := d.Id()

	log.Printf("[DEBUG] Executing read for KMS KeyRing %s", keyRingName)

	_, err := config.clientKms.Projects.Locations.KeyRings.
		Get(keyRingName).
		Do()

	if err != nil {
		return fmt.Errorf("Error reading KeyRing: %s", err)
	}

	return nil
}

func resourceKmsKeyRingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
