package google

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func DataSourceSqlDatabases() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceSqlDatabasesRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project that contains the instance.`,
			},
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Cloud SQL database instance in which the database belongs.`,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(ResourceSQLDatabase().Schema),
				},
			},
		},
	}
}

func dataSourceSqlDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	var databases *sqladmin.DatabasesListResponse
	err = RetryTimeDuration(func() (rerr error) {
		databases, rerr = config.NewSqlAdminClient(userAgent).Databases.List(project, d.Get("instance").(string)).Do()
		return rerr
	}, d.Timeout(schema.TimeoutRead), IsSqlOperationInProgressError)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Databases in %q instance", d.Get("instance").(string)))
	}
	flattenedDatabases := flattenDatabases(databases.Items)

	//client-side sorting to provide consistent ordering of the databases
	sort.SliceStable(flattenedDatabases, func(i, j int) bool {
		return strings.Compare(flattenedDatabases[i]["name"].(string), flattenedDatabases[j]["name"].(string)) < 1
	})
	if err := d.Set("databases", flattenedDatabases); err != nil {
		return fmt.Errorf("Error setting databases: %s", err)
	}
	d.SetId(fmt.Sprintf("project/%s/instance/%s/databases", project, d.Get("instance").(string)))
	return nil
}

func flattenDatabases(fetchedDatabases []*sqladmin.Database) []map[string]interface{} {
	if fetchedDatabases == nil {
		return make([]map[string]interface{}, 0)
	}

	databases := make([]map[string]interface{}, 0, len(fetchedDatabases))
	for _, rawDatabase := range fetchedDatabases {
		database := make(map[string]interface{})
		database["name"] = rawDatabase.Name
		database["instance"] = rawDatabase.Instance
		database["project"] = rawDatabase.Project
		database["charset"] = rawDatabase.Charset
		database["collation"] = rawDatabase.Collation
		database["self_link"] = rawDatabase.SelfLink

		databases = append(databases, database)
	}
	return databases
}
