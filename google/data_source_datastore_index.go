package google

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatastoreIndex() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatastoreIndexRead,

		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The entity kind which the index applies to.`,
			},
			"ancestor": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateEnum([]string{"NONE", "ALL_ANCESTORS", ""}),
				Description:  `Policy for including ancestors in the index. Default value: "NONE" Possible values: ["NONE", "ALL_ANCESTORS"]`,
				Default:      "NONE",
			},
			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `An ordered list of properties to index on.`,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"direction": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateEnum([]string{"ASCENDING", "DESCENDING"}),
							Description:  `The direction the index should optimize for sorting. Possible values: ["ASCENDING", "DESCENDING"]`,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The property name to index.`,
						},
					},
				},
			},
			"index_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The index id.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func dataSourceDatastoreIndexRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{DatastoreBasePath}}projects/{{project}}/indexes/{{index_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Index: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	obj := make(map[string]interface{})
	kindProp, err := expandDatastoreIndexKind(d.Get("kind"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kind"); !isEmptyValue(reflect.ValueOf(kindProp)) && (ok || !reflect.DeepEqual(v, kindProp)) {
		obj["kind"] = kindProp
	}
	ancestorProp, err := expandDatastoreIndexAncestor(d.Get("ancestor"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ancestor"); !isEmptyValue(reflect.ValueOf(ancestorProp)) && (ok || !reflect.DeepEqual(v, ancestorProp)) {
		obj["ancestor"] = ancestorProp
	}
	propertiesProp, err := expandDatastoreIndexProperties(d.Get("properties"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("properties"); !isEmptyValue(reflect.ValueOf(propertiesProp)) && (ok || !reflect.DeepEqual(v, propertiesProp)) {
		obj["properties"] = propertiesProp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil, datastoreIndex409Contention)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("DatastoreIndex %q", d.Id()))
	}

	if !reflect.DeepEqual(obj["kind"], res["kind"]) {
		return fmt.Errorf("Expected different Kind: %s. Actual: %s", obj["kind"], res["kind"])
	}
	if !reflect.DeepEqual(obj["ancestor"], res["ancestor"]) {
		return fmt.Errorf("Expected different Ancestor: %s. Actual: %s", obj["ancestor"], res["ancestor"])
	}
	if !reflect.DeepEqual(obj["properties"], res["properties"]) {
		return fmt.Errorf("Expected different Ancestor: %s. Actual: %s", obj["properties"], res["properties"])
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Index: %s", err)
	}

	return nil
}
