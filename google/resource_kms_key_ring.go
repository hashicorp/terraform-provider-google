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

func (s *kmsKeyRingId) parentId() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.Project, s.Location)
}

func (s *kmsKeyRingId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Location, s.Name)
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

	keyRing, err := config.clientKms.Projects.Locations.KeyRings.Create(keyRingId.parentId(), &cloudkms.KeyRing{}).KeyRingId(keyRingId.Name).Do()

	if err != nil {
		return fmt.Errorf("Error creating KeyRing: %s", err)
	}

	log.Printf("[DEBUG] Created KeyRing %s", keyRing.Name)

	d.SetId(keyRingId.terraformId())

	return resourceKmsKeyRingRead(d, meta)
}

func resourceKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyRingId, err := parseKmsKeyRingId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Executing read for KMS KeyRing %s", keyRingId.keyRingId())

	_, err = config.clientKms.Projects.Locations.KeyRings.Get(keyRingId.keyRingId()).Do()

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
	config := meta.(*Config)

	keyRingId, err := parseKmsKeyRingId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[WARNING] KMS KeyRing resources cannot be deleted from GCP. This KeyRing %s will be removed from Terraform state, but will still be present on the server.", keyRingId.keyRingId())

	d.SetId("")

	return nil
}

func parseKmsKeyRingId(id string, config *Config) (*kmsKeyRingId, error) {
	parts := strings.Split(id, "/")

	keyRingIdRegex := regexp.MustCompile("^([a-z0-9-]+)/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	keyRingIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")

	if keyRingIdRegex.MatchString(id) {
		return &kmsKeyRingId{
			Project:  parts[0],
			Location: parts[1],
			Name:     parts[2],
		}, nil
	}

	if keyRingIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{keyRingName}` id format.")
		}

		return &kmsKeyRingId{
			Project:  config.Project,
			Location: parts[0],
			Name:     parts[1],
		}, nil
	}

	return nil, fmt.Errorf("Invalid KeyRing id format, expecting `{projectId}/{locationId}/{keyRingName}` or `{locationId}/{keyRingName}.`")
}

func resourceKmsKeyRingImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	keyRingId, err := parseKmsKeyRingId(d.Id(), config)
	if err != nil {
		return nil, err
	}

	d.Set("name", keyRingId.Name)
	d.Set("location", keyRingId.Location)

	if config.Project != keyRingId.Project {
		d.Set("project", keyRingId.Project)
	}

	d.SetId(keyRingId.terraformId())

	return []*schema.ResourceData{d}, nil
}
