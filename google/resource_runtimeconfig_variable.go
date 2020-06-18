package google

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

func resourceRuntimeconfigVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuntimeconfigVariableCreate,
		Read:   resourceRuntimeconfigVariableRead,
		Update: resourceRuntimeconfigVariableUpdate,
		Delete: resourceRuntimeconfigVariableDelete,

		Importer: &schema.ResourceImporter{
			State: resourceRuntimeconfigVariableImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the variable to manage. Note that variable names can be hierarchical using slashes (e.g. "prod-variables/hostname").`,
			},

			"parent": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the RuntimeConfig resource containing this variable.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"value": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"text"},
			},

			"text": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"value"},
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds, representing when the variable was last updated. Example: "2016-10-09T12:33:37.578138407Z".`,
			},
		},
	}
}

func resourceRuntimeconfigVariableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	variable, parent, err := newRuntimeconfigVariableFromResourceData(d, project)
	if err != nil {
		return err
	}

	createdVariable, err := config.clientRuntimeconfig.Projects.Configs.Variables.Create(resourceRuntimeconfigFullName(project, parent), variable).Do()
	if err != nil {
		return err
	}
	d.SetId(createdVariable.Name)

	return setRuntimeConfigVariableToResourceData(d, *createdVariable)
}

func resourceRuntimeconfigVariableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	fullName := d.Id()
	createdVariable, err := config.clientRuntimeconfig.Projects.Configs.Variables.Get(fullName).Do()
	if err != nil {
		return err
	}

	return setRuntimeConfigVariableToResourceData(d, *createdVariable)
}

func resourceRuntimeconfigVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Update works more like an 'overwrite' method - we build a new runtimeconfig.Variable struct and it becomes the
	// new config. This means our Update logic looks an awful lot like Create (and hence, doesn't use
	// schema.ResourceData.hasChange()).

	variable, _, err := newRuntimeconfigVariableFromResourceData(d, project)
	if err != nil {
		return err
	}

	createdVariable, err := config.clientRuntimeconfig.Projects.Configs.Variables.Update(variable.Name, variable).Do()
	if err != nil {
		return err
	}

	return setRuntimeConfigVariableToResourceData(d, *createdVariable)
}

func resourceRuntimeconfigVariableDelete(d *schema.ResourceData, meta interface{}) error {
	fullName := d.Id()
	config := meta.(*Config)

	_, err := config.clientRuntimeconfig.Projects.Configs.Variables.Delete(fullName).Do()
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceRuntimeconfigVariableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/configs/(?P<parent>[^/]+)/variables/(?P<name>[^/]+)", "(?P<parent>[^/]+)/(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/configs/{{parent}}/variables/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

// resourceRuntimeconfigVariableFullName turns a given project, runtime config name, and a 'short name' for a runtime
// config variable into a full name (e.g. projects/my-project/configs/my-config/variables/my-variable).
func resourceRuntimeconfigVariableFullName(project, config, name string) string {
	return fmt.Sprintf("projects/%s/configs/%s/variables/%s", project, config, name)
}

// resourceRuntimeconfigVariableParseFullName parses a full name
// (e.g. projects/my-project/configs/my-config/variables/my-variable) by parsing out the
// project, runtime config name, and the short name. Returns "", "", "", err upon error.
func resourceRuntimeconfigVariableParseFullName(fullName string) (project, config, name string, err error) {
	re := regexp.MustCompile("^projects/([^/]+)/configs/([^/]+)/variables/(.+)$")
	matches := re.FindStringSubmatch(fullName)
	if matches == nil {
		return "", "", "", fmt.Errorf("Given full name doesn't match expected regexp; fullname = '%s'", fullName)
	}
	return matches[1], matches[2], matches[3], nil
}

// newRuntimeconfigVariableFromResourceData builds a new runtimeconfig.Variable struct from the data stored in a
// schema.ResourceData. Also returns the full name of the parent. Returns nil, "", err upon error.
func newRuntimeconfigVariableFromResourceData(d *schema.ResourceData, project string) (variable *runtimeconfig.Variable, parent string, err error) {
	// Validate that both text and value are not set
	text, textSet := d.GetOk("text")
	value, valueSet := d.GetOk("value")

	if !textSet && !valueSet {
		return nil, "", fmt.Errorf("You must specify one of value or text.")
	}

	// TODO(selmanj) here we assume it's a simple name, not a full name. Should probably support full name as well
	parent = d.Get("parent").(string)
	name := d.Get("name").(string)

	fullName := resourceRuntimeconfigVariableFullName(project, parent, name)

	variable = &runtimeconfig.Variable{
		Name: fullName,
	}

	if textSet {
		variable.Text = text.(string)
	} else {
		variable.Value = value.(string)
	}

	return variable, parent, nil
}

// setRuntimeConfigVariableToResourceData stores a provided runtimeconfig.Variable struct inside a schema.ResourceData.
func setRuntimeConfigVariableToResourceData(d *schema.ResourceData, variable runtimeconfig.Variable) error {
	varProject, parent, name, err := resourceRuntimeconfigVariableParseFullName(variable.Name)
	if err != nil {
		return err
	}
	d.Set("name", name)
	d.Set("parent", parent)
	d.Set("project", varProject)
	d.Set("value", variable.Value)
	d.Set("text", variable.Text)
	d.Set("update_time", variable.UpdateTime)

	return nil
}
