package google

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeSnapshot() *schema.Resource {

	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeSnapshot().Schema)

	dsSchema["filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	dsSchema["most_recent"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "name", "filter", "most_recent", "project")

	dsSchema["name"].ExactlyOneOf = []string{"name", "filter"}
	dsSchema["filter"].ExactlyOneOf = []string{"name", "filter"}

	return &schema.Resource{
		Read:   dataSourceGoogleComputeSnapshotRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		return retrieveSnapshot(d, meta, project, v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		projectGetCall := config.NewResourceManagerClient(userAgent).Projects.Get(project)

		if config.UserProjectOverride {
			billingProject := project

			// err == nil indicates that the billing_project value was found
			if bp, err := getBillingProject(d, config); err == nil {
				billingProject = bp
			}
			projectGetCall.Header().Add("X-Goog-User-Project", billingProject)
		}

		//handling the pagination locally
		allSnapshots := make([]*compute.Snapshot, 0)
		token := ""
		for paginate := true; paginate; {
			snapshots, err := config.NewComputeClient(userAgent).Snapshots.List(project).Filter(v.(string)).PageToken(token).Do()
			if err != nil {
				return fmt.Errorf("error retrieving list of snapshots: %s", err)

			}
			allSnapshots = append(allSnapshots, snapshots.Items...)

			token = snapshots.NextPageToken
			paginate = token != ""
		}

		mostRecent := d.Get("most_recent").(bool)
		if mostRecent {
			sort.Sort(ByCreationTimestampOfSnapshot(allSnapshots))
		}

		count := len(allSnapshots)
		if count == 1 || count > 1 && mostRecent {
			return retrieveSnapshot(d, meta, project, allSnapshots[0].Name)
		}

		return fmt.Errorf("your filter has returned %d snapshot(s). Please refine your filter or set most_recent to return exactly one snapshot", len(allSnapshots))

	}

	return fmt.Errorf("one of name or filter must be set")
}

func retrieveSnapshot(d *schema.ResourceData, meta interface{}, project, name string) error {
	d.SetId("projects/" + project + "/global/snapshots/" + name)
	d.Set("name", name)
	return resourceComputeSnapshotRead(d, meta)
}

// ByCreationTimestamp implements sort.Interface for []*Snapshot based on
// the CreationTimestamp field.
type ByCreationTimestampOfSnapshot []*compute.Snapshot

func (a ByCreationTimestampOfSnapshot) Len() int      { return len(a) }
func (a ByCreationTimestampOfSnapshot) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCreationTimestampOfSnapshot) Less(i, j int) bool {
	return a[i].CreationTimestamp > a[j].CreationTimestamp
}
