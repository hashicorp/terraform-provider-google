package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudkms/v1"
	"log"
	"regexp"
	"strings"
)

func resourceKmsKeyRing() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsKeyRingCreate,
		Read:   resourceKmsKeyRingRead,
		Delete: resourceKmsKeyRingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKmsKeyRingImportState,
		},

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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

type kmsKeyRingId struct {
	Project  string
	Location string
	Name     string
}

func (s *kmsKeyRingId) keyRingId() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", s.Project, s.Location, s.Name)
}

func (s *kmsKeyRingId) parentString() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.Project, s.Location)
}

func resourceKmsKeyRingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	keyRingId := &kmsKeyRingId{
		Project:  project,
		Location: d.Get("location").(string),
		Name:     d.Get("name").(string),
	}

	keyRing, err := config.clientKms.Projects.Locations.KeyRings.Create(keyRingId.parentString(), &cloudkms.KeyRing{}).KeyRingId(keyRingId.Name).Do()

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

	_, err := config.clientKms.Projects.Locations.KeyRings.Get(keyRingName).Do()

	if err != nil {
		return fmt.Errorf("Error reading KeyRing: %s", err)
	}

	return nil
}

/*
	Because KMS KeyRing resources cannot be deleted on GCP, we are only going to remove it from state.
	Re-creation of this resource through Terraform will produce an error.
*/

func resourceKmsKeyRingDelete(d *schema.ResourceData, meta interface{}) error {
	keyRingName := d.Id()

	log.Printf("[WARNING] KMS KeyRing resources cannot be deleted from GCP. This KeyRing %s will be removed from Terraform state, but will still be present on the server.", keyRingName)

	d.SetId("")

	return nil
}

func parseKmsKeyRingId(id string) (*kmsKeyRingId, error) {
	keyRingIdRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/keyRings/(.+)$")

	if !keyRingIdRegex.MatchString(id) {
		return nil, fmt.Errorf("Invalid KeyRing id format, expecting projects/{projectId}/locations/{locationId}/keyRings/{keyRingName}")
	}

	parts := strings.Split(id, "/")

	return &kmsKeyRingId{
		Project:  parts[1],
		Location: parts[3],
		Name:     parts[5],
	}, nil
}

func resourceKmsKeyRingImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id, err := parseKmsKeyRingId(d.Id())

	if err != nil {
		return nil, err
	}

	d.Set("name", id.Name)
	d.Set("location", id.Location)

	d.SetId(id.keyRingId())

	return []*schema.ResourceData{d}, nil
}
