// Package google - implement CRUD operations for Container Registry Build Triggers
// https://cloud.google.com/container-builder/docs/api/reference/rest/v1/projects.triggers#BuildTrigger
package google

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudbuild/v1"
)

func resourceCloudBuildTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudbuildBuildTriggerCreate,
		Read:   resourceCloudbuildBuildTriggerRead,
		Delete: resourceCloudbuildBuildTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudBuildTriggerImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"filename": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"build"},
			},
			"build": {
				Type:        schema.TypeList,
				Description: "Contents of the build template.",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"images": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"step": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"args": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"substitutions": &schema.Schema{
				Optional: true,
				Type:     schema.TypeMap,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"trigger_template": &schema.Schema{
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"commit_sha": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"dir": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"project": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"repo_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"tag_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudbuildBuildTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the address parameter
	buildTrigger := &cloudbuild.BuildTrigger{}

	if v, ok := d.GetOk("description"); ok {
		buildTrigger.Description = v.(string)
	}

	if v, ok := d.GetOk("filename"); ok {
		buildTrigger.Filename = v.(string)
	} else {
		buildTrigger.Build = expandCloudbuildBuildTriggerBuild(d)
	}

	buildTrigger.TriggerTemplate = expandCloudbuildBuildTriggerTemplate(d, project)
	buildTrigger.Substitutions = expandStringMap(d, "substitutions")

	tstr, err := json.Marshal(buildTrigger)
	if err != nil {
		return err
	}
	log.Printf("[INFO] build trigger request: %s", string(tstr))
	trigger, err := config.clientBuild.Projects.Triggers.Create(project, buildTrigger).Do()
	if err != nil {
		return fmt.Errorf("Error creating build trigger: %s", err)
	}

	d.SetId(trigger.Id)

	return resourceCloudbuildBuildTriggerRead(d, meta)
}

func resourceCloudbuildBuildTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	ID := d.Id()
	buildTrigger, err := config.clientBuild.Projects.Triggers.Get(project, ID).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Cloudbuild Trigger %q", ID))
	}

	d.Set("description", buildTrigger.Description)
	d.Set("substitutions", buildTrigger.Substitutions)

	if buildTrigger.TriggerTemplate != nil {
		d.Set("trigger_template", flattenCloudbuildBuildTriggerTemplate(d, config, buildTrigger.TriggerTemplate))
	}

	if buildTrigger.Filename != "" {
		d.Set("filename", buildTrigger.Filename)
	} else if buildTrigger.Build != nil {
		d.Set("build", flattenCloudbuildBuildTriggerBuild(d, config, buildTrigger.Build))
	}

	return nil
}

func expandCloudbuildBuildTriggerTemplate(d *schema.ResourceData, project string) *cloudbuild.RepoSource {
	if d.Get("trigger_template.#").(int) == 0 {
		return nil
	}
	tmpl := &cloudbuild.RepoSource{}
	if v, ok := d.GetOk("trigger_template.0.project"); ok {
		tmpl.ProjectId = v.(string)
	} else {
		tmpl.ProjectId = project
	}
	if v, ok := d.GetOk("trigger_template.0.branch_name"); ok {
		tmpl.BranchName = v.(string)
	}
	if v, ok := d.GetOk("trigger_template.0.commit_sha"); ok {
		tmpl.CommitSha = v.(string)
	}
	if v, ok := d.GetOk("trigger_template.0.dir"); ok {
		tmpl.Dir = v.(string)
	}
	if v, ok := d.GetOk("trigger_template.0.repo_name"); ok {
		tmpl.RepoName = v.(string)
	}
	if v, ok := d.GetOk("trigger_template.0.tag_name"); ok {
		tmpl.TagName = v.(string)
	}
	return tmpl
}

func flattenCloudbuildBuildTriggerTemplate(d *schema.ResourceData, config *Config, t *cloudbuild.RepoSource) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 1)

	flattened[0] = map[string]interface{}{
		"branch_name": t.BranchName,
		"commit_sha":  t.CommitSha,
		"dir":         t.Dir,
		"project":     t.ProjectId,
		"repo_name":   t.RepoName,
		"tag_name":    t.TagName,
	}

	return flattened
}

func expandCloudbuildBuildTriggerBuild(d *schema.ResourceData) *cloudbuild.Build {
	if d.Get("build.#").(int) == 0 {
		return nil
	}

	build := &cloudbuild.Build{}
	if v, ok := d.GetOk("build.0.images"); ok {
		build.Images = convertStringArr(v.([]interface{}))
	}
	if v, ok := d.GetOk("build.0.tags"); ok {
		build.Tags = convertStringArr(v.([]interface{}))
	}
	stepCount := d.Get("build.0.step.#").(int)
	build.Steps = make([]*cloudbuild.BuildStep, 0, stepCount)
	for s := 0; s < stepCount; s++ {
		step := &cloudbuild.BuildStep{
			Name: d.Get(fmt.Sprintf("build.0.step.%d.name", s)).(string),
		}
		if v, ok := d.GetOk(fmt.Sprintf("build.0.step.%d.args", s)); ok {
			step.Args = strings.Split(v.(string), " ")
		}
		build.Steps = append(build.Steps, step)
	}
	return build
}

func flattenCloudbuildBuildTriggerBuild(d *schema.ResourceData, config *Config, b *cloudbuild.Build) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 1)

	flattened[0] = map[string]interface{}{}

	if b.Images != nil {
		flattened[0]["images"] = convertStringArrToInterface(b.Images)
	}
	if b.Tags != nil {
		flattened[0]["tags"] = convertStringArrToInterface(b.Tags)
	}
	if b.Steps != nil {
		steps := make([]map[string]interface{}, len(b.Steps))
		for i, step := range b.Steps {
			steps[i] = map[string]interface{}{}
			steps[i]["name"] = step.Name
			steps[i]["args"] = strings.Join(step.Args, " ")
		}
		flattened[0]["step"] = steps
	}

	return flattened
}

func resourceCloudbuildBuildTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the build trigger
	log.Printf("[DEBUG] build trigger delete request")
	_, err = config.clientBuild.Projects.Triggers.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting build trigger: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceCloudBuildTriggerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) == 1 {
		return []*schema.ResourceData{d}, nil
	} else if len(parts) == 2 {
		d.Set("project", parts[0])
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, nil
	} else {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {trigger_name} or {project}/{trigger_name}", d.Id())
	}
}
