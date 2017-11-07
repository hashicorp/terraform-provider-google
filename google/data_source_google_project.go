package google

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func dataSourceGoogleProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleProjectRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	filter := fmt.Sprintf("name:%s", d.Get("name").(string))
	firstCall := config.clientResourceManager.Projects.List().Filter(filter)

	resp, err := firstCall.Do()
	if err != nil {
		return fmt.Errorf("Error reading Project: %s", err)
	}

	if len(resp.Projects) < 1 {
		// The project doesn't exist
		return fmt.Errorf("Project Not Found")
	}

	var activeProjects []*cloudresourcemanager.Project

	// filter any project that is not in "ACTIVE" state
	for _, project := range resp.Projects {
		if project.LifecycleState == "ACTIVE" {
			activeProjects = append(activeProjects, project)
		}
	}

	// fetch the entire projects collection but only append those which are active
	nextToken := resp.NextPageToken

	for nextToken != "" {
		call := firstCall.PageToken(nextToken)

		resp, err := call.Do()
		if err != nil {
			return fmt.Errorf("Error reading Project: %s", err)
		}

		for _, project := range resp.Projects {
			if project.LifecycleState == "ACTIVE" {
				activeProjects = append(activeProjects, project)
			}
		}

		nextToken = resp.NextPageToken
	}

	// sort by ascending order of creation
	sort.Slice(activeProjects, func(i, j int) bool {
		iCreateTime, _ := time.Parse(time.RFC3339, activeProjects[i].CreateTime)
		jCreateTime, _ := time.Parse(time.RFC3339, activeProjects[j].CreateTime)

		return iCreateTime.Before(jCreateTime)
	})

	// the last element is the most recent one
	lastActiveProject := activeProjects[len(activeProjects)-1]

	log.Printf("[DEBUG] Google project found: %q", lastActiveProject)

	d.Set("project_id", lastActiveProject.ProjectId)
	d.Set("number", strconv.FormatInt(int64(lastActiveProject.ProjectNumber), 10))
	d.SetId(lastActiveProject.ProjectId)

	return nil
}
