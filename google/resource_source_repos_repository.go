package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/sourcerepo/v1"
	"log"
)

func resourceSourceReposRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceSourceReposRepositoryCreate,
		Read:   resourceSourceReposRepositoryRead,
		Delete: resourceSourceReposRepositoryDelete,
		//Update: not supported,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
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

func resourceSourceReposRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if _, ok := d.GetOk("project"); ok {
		log.Printf(d.Get("project").(string))
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	repoName := d.Get("name").(string)
	name := buildRepositoryName(project, repoName)

	repo := &sourcerepo.Repo{
		Name: name,
	}

	project = "projects/" + project

	job, err := config.clientSourceRepos.Projects.Repos.Create(project, repo).Do()
	if err != nil {
		return err
	}
	d.SetId(job.Name)

	return nil
}

func resourceSourceReposRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	//project = "projects/" + project
	//name := project + "/repos/" + d.Get("name").(string)

	repoName := d.Get("name").(string)
	name := buildRepositoryName(project, repoName)

	_, err = config.clientSourceRepos.Projects.Repos.Get(name).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Source Repo %q", d.Id()))
	}

	return nil
}

func resourceSourceReposRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	repoName := d.Get("name").(string)
	name := buildRepositoryName(project, repoName)

	_, err = config.clientSourceRepos.Projects.Repos.Delete(name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting Source Repos Repository: %s", err)
	}

	return nil
}

func buildRepositoryName(project, name string) string {
	repositoryName := "projects/" + project + "/repos/" + name
	return repositoryName
}