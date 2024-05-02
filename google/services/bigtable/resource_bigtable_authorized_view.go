// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var familySubsetSchemaElem *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"family_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: `Name of the column family to be included in the authorized view.`,
		},
		"qualifiers": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: `Base64-encoded individual exact column qualifiers of the column family to be included in the authorized view.`,
		},
		"qualifier_prefixes": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: `Base64-encoded prefixes for qualifiers of the column family to be included in the authorized view. Every qualifier starting with one of these prefixes is included in the authorized view. To provide access to all qualifiers, include the empty string as a prefix ("").`,
		},
	},
}

func ResourceBigtableAuthorizedView() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableAuthorizedViewCreate,
		Read:   resourceBigtableAuthorizedViewRead,
		Update: resourceBigtableAuthorizedViewUpdate,
		Delete: resourceBigtableAuthorizedViewDestroy,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableAuthorizedViewImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the authorized view. Must be 1-50 characters and must only contain hyphens, underscores, periods, letters and numbers.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"instance_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `The name of the Bigtable instance in which the authorized view belongs.`,
			},

			"table_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `The name of the Bigtable table in which the authorized view belongs.`,
			},

			"deletion_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROTECTED", "UNPROTECTED"}, false),
				Description: `A field to make the authorized view protected against data loss i.e. when set to PROTECTED, deleting the authorized view, the table containing the authorized view, and the instance containing the authorized view would be prohibited.
If not provided, currently deletion protection will be set to UNPROTECTED as it is the API default value. Note this field configs the deletion protection provided by the API in the backend, and should not be confused with Terraform-side deletion protection.`,
			},

			"subset_view": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `An AuthorizedView permitting access to an explicit subset of a Table.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"row_prefixes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Base64-encoded row prefixes to be included in the authorized view. To provide access to all rows, include the empty string as a prefix ("").`,
						},
						"family_subsets": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `Subsets of column families to be included in the authorized view.`,
							Elem:        familySubsetSchemaElem,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceBigtableAuthorizedViewCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	if err := d.Set("instance_name", instanceName); err != nil {
		return fmt.Errorf("Error setting instance_name: %s", err)
	}

	defer c.Close()

	authorizedViewId := d.Get("name").(string)
	tableId := d.Get("table_name").(string)
	authorizedViewConf := bigtable.AuthorizedViewConf{
		AuthorizedViewID: authorizedViewId,
		TableID:          tableId,
	}

	// Check if deletion protection is given
	// If not given, currently tblConf.DeletionProtection will be set to false in the API
	deletionProtection := d.Get("deletion_protection")
	if deletionProtection == "PROTECTED" {
		authorizedViewConf.DeletionProtection = bigtable.Protected
	} else if deletionProtection == "UNPROTECTED" {
		authorizedViewConf.DeletionProtection = bigtable.Unprotected
	}

	subsetView, ok := d.GetOk("subset_view")
	if !ok || len(subsetView.([]interface{})) != 1 {
		return fmt.Errorf("subset_view must be specified for authorized view %s", authorizedViewId)
	}
	subsetViewConf, err := generateSubsetViewConfig(subsetView.([]interface{}))
	if err != nil {
		return err
	}
	authorizedViewConf.AuthorizedView = subsetViewConf

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel() // Always call cancel.
	err = c.CreateAuthorizedView(ctxWithTimeout, &authorizedViewConf)
	if err != nil {
		return fmt.Errorf("Error creating authorized view. %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{table_name}}/authorizedViews/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceBigtableAuthorizedViewRead(d, meta)
}

func resourceBigtableAuthorizedViewRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	authorizedViewId := d.Get("name").(string)
	tableId := d.Get("table_name").(string)
	authorizedViewInfo, err := c.AuthorizedViewInfo(ctx, tableId, authorizedViewId)
	if err != nil {
		if tpgresource.IsNotFoundGrpcError(err) {
			log.Printf("[WARN] Removing %s because it's gone", authorizedViewId)
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	deletionProtection := authorizedViewInfo.DeletionProtection
	if deletionProtection == bigtable.Protected {
		if err := d.Set("deletion_protection", "PROTECTED"); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	} else if deletionProtection == bigtable.Unprotected {
		if err := d.Set("deletion_protection", "UNPROTECTED"); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	} else {
		return fmt.Errorf("Error setting deletion_protection, it should be either PROTECTED or UNPROTECTED")
	}

	if sv, ok := authorizedViewInfo.AuthorizedView.(*bigtable.SubsetViewInfo); ok {
		subsetView := flattenSubsetViewInfo(sv)
		if err := d.Set("subset_view", subsetView); err != nil {
			return fmt.Errorf("Error setting subset_view: %s", err)
		}
	} else {
		return fmt.Errorf("Error parsing server returned subset_view since it's empty")
	}

	return nil
}

func resourceBigtableAuthorizedViewUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	defer c.Close()

	authorizedViewId := d.Get("name").(string)
	tableId := d.Get("table_name").(string)
	authorizedViewConf := bigtable.AuthorizedViewConf{
		AuthorizedViewID: authorizedViewId,
		TableID:          tableId,
	}

	if d.HasChange("subset_view") {
		subsetView := d.Get("subset_view")
		if len(subsetView.([]interface{})) != 1 {
			return fmt.Errorf("subset_view must be specified for authorized view %s", authorizedViewId)
		}
		subsetViewConf, err := generateSubsetViewConfig(subsetView.([]interface{}))
		if err != nil {
			return err
		}
		authorizedViewConf.AuthorizedView = subsetViewConf
	}

	if d.HasChange("deletion_protection") {
		deletionProtection := d.Get("deletion_protection")
		if deletionProtection == "PROTECTED" {
			authorizedViewConf.DeletionProtection = bigtable.Protected
		} else if deletionProtection == "UNPROTECTED" {
			authorizedViewConf.DeletionProtection = bigtable.Unprotected
		}
	}

	updateAuthorizedViewConf := bigtable.UpdateAuthorizedViewConf{
		AuthorizedViewConf: authorizedViewConf,
		IgnoreWarnings:     true,
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel() // Always call cancel.
	err = c.UpdateAuthorizedView(ctxWithTimeout, updateAuthorizedViewConf)
	if err != nil {
		return fmt.Errorf("Error updating authorized view. %s", err)
	}

	return resourceBigtableAuthorizedViewRead(d, meta)
}

func resourceBigtableAuthorizedViewDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	authorizedViewId := d.Get("name").(string)
	tableId := d.Get("table_name").(string)
	err = c.DeleteAuthorizedView(ctx, tableId, authorizedViewId)
	if err != nil {
		return fmt.Errorf("Error deleting authorized view. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceBigtableAuthorizedViewImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance_name>[^/]+)/tables/(?P<table_name>[^/]+)/authorizedViews/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance_name>[^/]+)/(?P<table_name>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance_name>[^/]+)/(?P<table_name>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{table_name}}/authorizedViews/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func generateSubsetViewConfig(subsetView []interface{}) (*bigtable.SubsetViewConf, error) {
	subsetViewConf := bigtable.SubsetViewConf{}

	if len(subsetView) == 0 {
		return nil, fmt.Errorf("Error constructing SubsetViewConfig; empty subset_view list")
	}
	if subsetView[0] == nil {
		return &subsetViewConf, nil
	}
	sv, ok := subsetView[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error constructing SubsetViewConfig; element in subset_view list has wrong type: %s", reflect.TypeOf(subsetView[0]))
	}
	if rowPrefixes, ok := sv["row_prefixes"]; ok {
		for _, rowPrefix := range rowPrefixes.(*schema.Set).List() {
			decodedRowPrefix, err := base64.StdEncoding.DecodeString(rowPrefix.(string))
			if err != nil {
				return nil, err
			}
			subsetViewConf.AddRowPrefix(decodedRowPrefix)
		}
	}
	if familySubsets, ok := sv["family_subsets"]; ok {
		for _, familySubset := range familySubsets.(*schema.Set).List() {
			familySubsetElement := familySubset.(map[string]interface{})
			familyName := familySubsetElement["family_name"].(string)
			if qualifiers, ok := familySubsetElement["qualifiers"]; ok {
				for _, qualifier := range qualifiers.(*schema.Set).List() {
					decodedQualifier, err := base64.StdEncoding.DecodeString(qualifier.(string))
					if err != nil {
						return nil, err
					}
					subsetViewConf.AddFamilySubsetQualifier(familyName, decodedQualifier)
				}
			}
			if qualifierPrefixes, ok := familySubsetElement["qualifier_prefixes"]; ok {
				for _, qualifierPrefix := range qualifierPrefixes.(*schema.Set).List() {
					decodedQualifierPrefix, err := base64.StdEncoding.DecodeString(qualifierPrefix.(string))
					if err != nil {
						return nil, err
					}
					subsetViewConf.AddFamilySubsetQualifierPrefix(familyName, decodedQualifierPrefix)
				}
			}
		}
	}
	return &subsetViewConf, nil
}

func flattenSubsetViewInfo(subsetViewInfo *bigtable.SubsetViewInfo) []map[string]interface{} {
	subsetView := make([]map[string]interface{}, 1)

	rowPrefixes := []string{}
	for _, prefix := range subsetViewInfo.RowPrefixes {
		encodedRowPrefix := base64.StdEncoding.EncodeToString(prefix)
		rowPrefixes = append(rowPrefixes, encodedRowPrefix)
	}
	familySubsets := []map[string]interface{}{}
	for k, v := range subsetViewInfo.FamilySubsets {
		familySubsetElement := make(map[string]interface{})
		familySubsetElement["family_name"] = k
		qualifiers := []string{}
		for _, qualifier := range v.Qualifiers {
			encodedQualifier := base64.StdEncoding.EncodeToString(qualifier)
			qualifiers = append(qualifiers, encodedQualifier)
		}
		if len(qualifiers) > 0 {
			familySubsetElement["qualifiers"] = qualifiers
		}
		qualifierPrefixes := []string{}
		for _, qualifierPrefix := range v.QualifierPrefixes {
			encodedQualifierPrefix := base64.StdEncoding.EncodeToString(qualifierPrefix)
			qualifierPrefixes = append(qualifierPrefixes, encodedQualifierPrefix)
		}
		if len(qualifierPrefixes) > 0 {
			familySubsetElement["qualifier_prefixes"] = qualifierPrefixes
		}
		familySubsets = append(familySubsets, familySubsetElement)
	}
	subsetView[0] = make(map[string]interface{})
	if len(rowPrefixes) > 0 {
		subsetView[0]["row_prefixes"] = rowPrefixes
	}
	if len(familySubsets) > 0 {
		subsetView[0]["family_subsets"] = familySubsets
	}

	return subsetView
}
