package tpgiamresource

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var iamAuditConfigSchema = map[string]*schema.Schema{
	"service": {
		Type:        schema.TypeString,
		Required:    true,
		Description: `Service which will be enabled for audit logging. The special value allServices covers all services.`,
	},
	"audit_log_config": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: `The configuration for logging of each type of permission. This can be specified multiple times.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"log_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: `Permission type for which logging is to be configured. Must be one of DATA_READ, DATA_WRITE, or ADMIN_READ.`,
				},
				"exempted_members": {
					Type:        schema.TypeSet,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Description: `Identities that do not cause logging for this type of permission. Each entry can have one of the following values:user:{emailid}: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com. serviceAccount:{emailid}: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com. group:{emailid}: An email address that represents a Google group. For example, admins@example.com. domain:{domain}: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.`,
				},
			},
		},
	},
	"etag": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The etag of iam policy`,
	},
}

func ResourceIamAuditConfig(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc NewResourceIamUpdaterFunc, resourceIdParser ResourceIdParserFunc, options ...func(*IamSettings)) *schema.Resource {
	settings := NewIamSettings(options...)

	return &schema.Resource{
		Create: resourceIamAuditConfigCreateUpdate(newUpdaterFunc, settings.EnableBatching),
		Read:   resourceIamAuditConfigRead(newUpdaterFunc),
		Update: resourceIamAuditConfigCreateUpdate(newUpdaterFunc, settings.EnableBatching),
		Delete: resourceIamAuditConfigDelete(newUpdaterFunc, settings.EnableBatching),
		Schema: tpgresource.MergeSchemas(iamAuditConfigSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamAuditConfigImport(resourceIdParser),
		},
		UseJSONNumber: true,
	}
}

func resourceIamAuditConfigRead(newUpdaterFunc NewResourceIamUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		eAuditConfig := getResourceIamAuditConfig(d)
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("AuditConfig for %s on %q", eAuditConfig.Service, updater.DescribeResource()))
		}
		log.Printf("[DEBUG]: Retrieved policy for %s: %+v", updater.DescribeResource(), p)

		var ac *cloudresourcemanager.AuditConfig
		for _, b := range p.AuditConfigs {
			if b.Service != eAuditConfig.Service {
				continue
			}
			ac = b
			break
		}
		if ac == nil {
			log.Printf("[DEBUG]: AuditConfig for service %q not found in policy for %s, removing from state file.", eAuditConfig.Service, updater.DescribeResource())
			d.SetId("")
			return nil
		}

		if err := d.Set("etag", p.Etag); err != nil {
			return fmt.Errorf("Error setting etag: %s", err)
		}
		err = d.Set("audit_log_config", flattenAuditLogConfigs(ac.AuditLogConfigs))
		if err != nil {
			return fmt.Errorf("Error flattening audit log config: %s", err)
		}
		if err := d.Set("service", ac.Service); err != nil {
			return fmt.Errorf("Error setting service: %s", err)
		}
		return nil
	}
}

func iamAuditConfigImport(resourceIdParser ResourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*transport_tpg.Config)
		s := strings.Fields(d.Id())
		if len(s) != 2 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to AuditConfig id %s; expected 'resource_name service'.", s)
		}
		id, service := s[0], s[1]

		// Set the ID only to the first part so all IAM types can share the same ResourceIdParserFunc.
		d.SetId(id)
		if err := d.Set("service", service); err != nil {
			return nil, fmt.Errorf("Error setting service: %s", err)
		}
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the ResourceIdParserFunc.
		d.SetId(d.Id() + "/audit_config/" + service)
		return []*schema.ResourceData{d}, nil
	}
}

func resourceIamAuditConfigCreateUpdate(newUpdaterFunc NewResourceIamUpdaterFunc, enableBatching bool) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		ac := getResourceIamAuditConfig(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			cleaned := removeAllAuditConfigsWithService(ep.AuditConfigs, ac.Service)
			ep.AuditConfigs = append(cleaned, ac)
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config, fmt.Sprintf(
				"Overwrite audit config for service %s on resource %q", ac.Service, updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/audit_config/" + ac.Service)
		return resourceIamAuditConfigRead(newUpdaterFunc)(d, meta)
	}
}

func resourceIamAuditConfigDelete(newUpdaterFunc NewResourceIamUpdaterFunc, enableBatching bool) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		ac := getResourceIamAuditConfig(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			ep.AuditConfigs = removeAllAuditConfigsWithService(ep.AuditConfigs, ac.Service)
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config, fmt.Sprintf(
				"Delete audit config for service %s on resource %q", ac.Service, updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Resource %s with IAM audit config %q", updater.DescribeResource(), d.Id()))
		}

		return resourceIamAuditConfigRead(newUpdaterFunc)(d, meta)
	}
}

func getResourceIamAuditConfig(d *schema.ResourceData) *cloudresourcemanager.AuditConfig {
	auditLogConfigSet := d.Get("audit_log_config").(*schema.Set)
	auditLogConfigs := make([]*cloudresourcemanager.AuditLogConfig, auditLogConfigSet.Len())
	for x, y := range auditLogConfigSet.List() {
		logConfig := y.(map[string]interface{})
		auditLogConfigs[x] = &cloudresourcemanager.AuditLogConfig{
			LogType:         logConfig["log_type"].(string),
			ExemptedMembers: tpgresource.ConvertStringArr(logConfig["exempted_members"].(*schema.Set).List()),
		}
	}
	return &cloudresourcemanager.AuditConfig{
		AuditLogConfigs: auditLogConfigs,
		Service:         d.Get("service").(string),
	}
}

func flattenAuditLogConfigs(configs []*cloudresourcemanager.AuditLogConfig) *schema.Set {
	auditLogConfigSchema := iamAuditConfigSchema["audit_log_config"].Elem.(*schema.Resource)
	exemptedMemberSchema := auditLogConfigSchema.Schema["exempted_members"].Elem.(*schema.Schema)
	res := schema.NewSet(schema.HashResource(auditLogConfigSchema), []interface{}{})
	for _, conf := range configs {
		res.Add(map[string]interface{}{
			"log_type":         conf.LogType,
			"exempted_members": schema.NewSet(schema.HashSchema(exemptedMemberSchema), tpgresource.ConvertStringArrToInterface(conf.ExemptedMembers)),
		})
	}
	return res
}
