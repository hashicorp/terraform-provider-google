package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/sourcerepo/v1"
)

func resourceSourceRepoRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceSourceRepoRepositoryCreate,
		Read:   resourceSourceRepoRepositoryRead,
		Delete: resourceSourceRepoRepositoryDelete,
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

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSourceRepoRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	repoName := d.Get("name").(string)
	name := buildRepositoryName(project, repoName)

	repo := &sourcerepo.Repo{
		Name: name,
	}

	parent := "projects/" + project

	op, err := config.clientSourceRepo.Projects.Repos.Create(parent, repo).Do()
	if err != nil {
		return fmt.Errorf("Error creating the Source Repo: %s", err)
	}
	d.SetId(op.Name)

	return nil
}

func resourceSourceRepoRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	repoName := d.Get("name").(string)
	name := buildRepositoryName(project, repoName)

	repo, err := config.clientSourceRepo.Projects.Repos.Get(name).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Source Repo %q", d.Id()))
	}

	d.Set("size", repo.Size)

	return nil
}

func resourceSourceRepoRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	repoName := d.Get("name").(string)
	name := buildRepositoryName(project, repoName)

	_, err = config.clientSourceRepo.Projects.Repos.Delete(name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting the Source Repo: %s", err)
	}

	return nil
}

func buildRepositoryName(project, name string) string {
	repositoryName := "projects/" + project + "/repos/" + name
	return repositoryName
}
