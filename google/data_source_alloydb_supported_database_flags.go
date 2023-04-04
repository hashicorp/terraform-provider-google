package google

import (
	alloydb "cloud.google.com/go/alloydb/apiv1"
	alloydbpb "cloud.google.com/go/alloydb/apiv1/alloydbpb"
	"fmt"
	gax "github.com/googleapis/gax-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/iterator"
)

func DataSourceAlloydbSupportedDatabaseFlags() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceAlloydbSupportedDatabaseFlagsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The canonical id for the location. For example: "us-east1".`,
			},
			"supported_database_flags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The name of the flag resource, following Google Cloud conventions, e.g.: * projects/{project}/locations/{location}/flags/{flag} This field currently has no semantic meaning.`,
						},
						"flag_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The name of the database flag, e.g. "max_allowed_packets". The is a possibly key for the Instance.database_flags map field.`,
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `ValueType describes the semantic type of the value that the flag accepts. The supported values are:- 'VALUE_TYPE_UNSPECIFIED', 'STRING', 'INTEGER', 'FLOAT', 'NONE'.`,
						},
						"accepts_multiple_values": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: `Whether the database flag accepts multiple values. If true, a comma-separated list of stringified values may be specified.`,
						},
						"supported_db_versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Major database engine versions for which this flag is supported. Supported values are:- 'DATABASE_VERSION_UNSPECIFIED', and 'POSTGRES_14'.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"requires_db_restart": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: `Whether setting or updating this flag on an Instance requires a database restart. If a flag that requires database restart is set, the backend will automatically restart the database (making sure to satisfy any availability SLO's).`,
						},
						"string_restrictions": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Restriction on STRING type value.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_values": {
										Type:        schema.TypeList,
										Computed:    true,
										Optional:    true,
										Description: `The list of allowed values, if bounded. This field will be empty if there is a unbounded number of allowed values.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"integer_restrictions": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Restriction on INTEGER type value.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										Description: `The minimum value that can be specified, if applicable.`,
									},
									"max_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										Description: `The maximum value that can be specified, if applicable.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlloydbSupportedDatabaseFlagsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}
	billingProject := project
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	location := ""
	if v, ok := d.GetOk("location"); ok {
		location = v.(string)
	}
	if location == "" {
		return fmt.Errorf("Location cannot be empty")
	}
	var supportedDatabaseFlagIterator *alloydb.SupportedDatabaseFlagIterator
	dbFlagsReq := new(alloydbpb.ListSupportedDatabaseFlagsRequest)
	alloydbClient := config.NewAlloydbClient(userAgent)
	if alloydbClient == nil {
		return fmt.Errorf("Failed to call the API to fetch the supported database flags")
	}
	err = nil
	err = retryTime(func() error {
		url := fmt.Sprintf("v1/projects/%s/locations/%s/supportedDatabaseFlags", billingProject, location)
		supportedDatabaseFlagIterator = alloydbClient.ListSupportedDatabaseFlags(config.context, dbFlagsReq, gax.WithPath(url))
		return nil
	}, 5)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Supported Database flags %q", d.Id()))
	}

	var supportedDatabaseFlags []map[string]interface{}
	for {
		supportedDatabaseFlag := make(map[string]interface{})
		flag, err := supportedDatabaseFlagIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to fetch the supported database flags for the provided location: %s", location))
		}
		if flag.Name != "" {
			supportedDatabaseFlag["name"] = flag.Name
		}
		if flag.FlagName != "" {
			supportedDatabaseFlag["flag_name"] = flag.FlagName
		}
		supportedDatabaseFlag["value_type"] = flag.ValueType.String()
		supportedDatabaseFlag["accepts_multiple_values"] = flag.AcceptsMultipleValues
		supportedDatabaseFlag["requires_db_restart"] = flag.RequiresDbRestart
		if flag.SupportedDbVersions != nil {
			dbVersions := make([]string, 0, len(flag.SupportedDbVersions))
			for _, supDbVer := range flag.SupportedDbVersions {
				dbVersions = append(dbVersions, supDbVer.String())
			}
			supportedDatabaseFlag["supported_db_versions"] = dbVersions
		}

		if flag.Restrictions != nil {
			if stringRes, ok := flag.Restrictions.(*alloydbpb.SupportedDatabaseFlag_StringRestrictions_); ok {
				restrictions := make([]map[string][]string, 0, 1)
				fetchedAllowedValues := stringRes.StringRestrictions.AllowedValues
				if fetchedAllowedValues != nil {
					allowedValues := make([]string, 0, len(fetchedAllowedValues))
					for _, val := range fetchedAllowedValues {
						allowedValues = append(allowedValues, val)
					}
					stringRestrictions := map[string][]string{
						"allowed_values": allowedValues,
					}
					restrictions = append(restrictions, stringRestrictions)
					supportedDatabaseFlag["string_restrictions"] = restrictions
				}
			}
			if integerRes, ok := flag.Restrictions.(*alloydbpb.SupportedDatabaseFlag_IntegerRestrictions_); ok {
				restrictions := make([]map[string]int64, 0, 1)
				minValue := integerRes.IntegerRestrictions.MinValue
				maxValue := integerRes.IntegerRestrictions.MaxValue
				integerRestrictions := map[string]int64{
					"min_value": minValue.GetValue(),
					"max_value": maxValue.GetValue(),
				}
				restrictions = append(restrictions, integerRestrictions)
				supportedDatabaseFlag["integer_restrictions"] = restrictions
			}
		}
		supportedDatabaseFlags = append(supportedDatabaseFlags, supportedDatabaseFlag)
	}
	if err := d.Set("supported_database_flags", supportedDatabaseFlags); err != nil {
		return fmt.Errorf("Error setting supported_database_flags: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/supportedDbFlags", project, location))
	return nil
}
