package google

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

var runtimeConfigFullName *regexp.Regexp = regexp.MustCompile("^projects/([^/]+)/configs/(.+)$")

func resourceRuntimeconfigConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuntimeconfigConfigCreate,
		Read:   resourceRuntimeconfigConfigRead,
		Update: resourceRuntimeconfigConfigUpdate,
		Delete: resourceRuntimeconfigConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceRuntimeconfigConfigImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp("[0-9A-Za-z](?:[_.A-Za-z0-9-]{0,62}[_.A-Za-z0-9])?"),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceRuntimeconfigConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	fullName := resourceRuntimeconfigFullName(project, name)
	runtimeConfig := runtimeconfig.RuntimeConfig{
		Name: fullName,
	}

	if val, ok := d.GetOk("description"); ok {
		runtimeConfig.Description = val.(string)
	}

	_, err = config.clientRuntimeconfig.Projects.Configs.Create("projects/"+project, &runtimeConfig).Do()

	if err != nil {
		return err
	}
	d.SetId(fullName)

	return nil
}

func resourceRuntimeconfigConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	fullName := d.Id()
	runConfig, err := config.clientRuntimeconfig.Projects.Configs.Get(fullName).Do()
	if err != nil {
		return err
	}

	project, name, err := resourceRuntimeconfigParseFullName(runConfig.Name)
	if err != nil {
		return err
	}
	// Check to see if project matches our current defined value - if it doesn't, we'll explicitly set it
	curProject, err := getProject(d, config)
	if err != nil {
		return err
	}
	if project != curProject {
		d.Set("project", project)
	}

	d.Set("name", name)
	d.Set("description", runConfig.Description)
	d.Set("project", project)

	return nil
}

func resourceRuntimeconfigConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Update works more like an 'overwrite' method - we build a new runtimeconfig.RuntimeConfig struct and it becomes
	// the new config. This means our Update logic looks an awful lot like Create (and hence, doesn't use
	// schema.ResourceData.hasChange()).
	fullName := d.Id()
	runtimeConfig := runtimeconfig.RuntimeConfig{
		Name: fullName,
	}
	if v, ok := d.GetOk("description"); ok {
		runtimeConfig.Description = v.(string)
	}

	_, err := config.clientRuntimeconfig.Projects.Configs.Update(fullName, &runtimeConfig).Do()
	if err != nil {
		return err
	}
	return nil
}

func resourceRuntimeconfigConfigDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	fullName := d.Id()

	_, err := config.clientRuntimeconfig.Projects.Configs.Delete(fullName).Do()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceRuntimeconfigConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/configs/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/configs/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

// resourceRuntimeconfigFullName turns a given project and a 'short name' for a runtime config into a full name
// (e.g. projects/my-project/configs/my-config).
func resourceRuntimeconfigFullName(project, name string) string {
	return fmt.Sprintf("projects/%s/configs/%s", project, name)
}

// resourceRuntimeconfigParseFullName parses a full name (e.g. projects/my-project/configs/my-config) by parsing out the
// project and the short name. Returns "", "", nil upon error.
func resourceRuntimeconfigParseFullName(fullName string) (project, name string, err error) {
	matches := runtimeConfigFullName.FindStringSubmatch(fullName)
	if matches == nil {
		return "", "", fmt.Errorf("Given full name doesn't match expected regexp; fullname = '%s'", fullName)
	}
	return matches[1], matches[2], nil
}
