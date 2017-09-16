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

func resourceKmsKeyRingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	location := d.Get("location").(string)

	parent := fmt.Sprintf("projects/%s/locations/%s", project, location)

	_, err = config.clientKms.Projects.Locations.KeyRings.
		Create(parent, &cloudkms.KeyRing{}).
		KeyRingId(name).
		Do()

	if err != nil {
		return fmt.Errorf("Error creating KeyRing: %s", err)
	}

	d.SetId(name)
	d.Set("location", location)

	return resourceKmsKeyRingRead(d, meta)
}

func resourceKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()

	log.Printf("[DEBUG] Executing read for KMS KeyRing %s", name)

	keyRing, err := config.clientKms.Projects.Locations.KeyRings.
		Get(name).
		Do()

	if err != nil {
		return fmt.Errorf("Error reading KeyRing: %s", err)
	}

	d.SetId(keyRing.Name)

	return nil
}

func resourceKmsKeyRingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
