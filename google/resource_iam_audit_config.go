package google

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var iamAuditConfigSchema = map[string]*schema.Schema{
	"service": {
		Type:     schema.TypeString,
		Required: true,
	},
	"audit_log_config": {
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"log_type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"exempted_members": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
				},
			},
		},
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ResourceIamAuditConfig(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc) *schema.Resource {
	return &schema.Resource{
		Create: resourceIamAuditConfigCreate(newUpdaterFunc),
		Read:   resourceIamAuditConfigRead(newUpdaterFunc),
		Update: resourceIamAuditConfigUpdate(newUpdaterFunc),
		Delete: resourceIamAuditConfigDelete(newUpdaterFunc),
		Schema: mergeSchemas(iamAuditConfigSchema, parentSpecificSchema),
	}
}

func ResourceIamAuditConfigWithImport(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) *schema.Resource {
	r := ResourceIamAuditConfig(parentSpecificSchema, newUpdaterFunc)
	r.Importer = &schema.ResourceImporter{
		State: iamAuditConfigImport(resourceIdParser),
	}
	return r
}

func resourceIamAuditConfigCreate(newUpdaterFunc newResourceIamUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		p := getResourceIamAuditConfig(d)
		err = iamPolicyReadModifyWrite(updater, func(ep *cloudresourcemanager.Policy) error {
			ep.AuditConfigs = mergeAuditConfigs(append(ep.AuditConfigs, p))
			return nil
		})
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/audit_config/" + p.Service)
		return resourceIamAuditConfigRead(newUpdaterFunc)(d, meta)
	}
}

func resourceIamAuditConfigRead(newUpdaterFunc newResourceIamUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		eAuditConfig := getResourceIamAuditConfig(d)
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG]: AuditConfig for service %q not found for non-existent resource %s, removing from state file.", eAuditConfig.Service, updater.DescribeResource())
				d.SetId("")
				return nil
			}

			return err
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
		d.Set("etag", p.Etag)
		err = d.Set("audit_log_config", flattenAuditLogConfigs(ac.AuditLogConfigs))
		if err != nil {
			return fmt.Errorf("Error flattening audit log config: %s", err)
		}
		d.Set("service", ac.Service)
		return nil
	}
}

func iamAuditConfigImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Fields(d.Id())
		if len(s) != 2 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to AuditConfig id %s; expected 'resource_name service'.", s)
		}
		id, service := s[0], s[1]

		// Set the ID only to the first part so all IAM types can share the same resourceIdParserFunc.
		d.SetId(id)
		d.Set("service", service)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the resourceIdParserFunc.
		d.SetId(d.Id() + "/audit_config/" + service)
		return []*schema.ResourceData{d}, nil
	}
}

func resourceIamAuditConfigUpdate(newUpdaterFunc newResourceIamUpdaterFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		ac := getResourceIamAuditConfig(d)
		err = iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error {
			var found bool
			for pos, b := range p.AuditConfigs {
				if b.Service != ac.Service {
					continue
				}
				found = true
				p.AuditConfigs[pos] = ac
				break
			}
			if !found {
				p.AuditConfigs = append(p.AuditConfigs, ac)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceIamAuditConfigRead(newUpdaterFunc)(d, meta)
	}
}

func resourceIamAuditConfigDelete(newUpdaterFunc newResourceIamUpdaterFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		ac := getResourceIamAuditConfig(d)
		err = iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error {
			toRemove := -1
			for pos, b := range p.AuditConfigs {
				if b.Service != ac.Service {
					continue
				}
				toRemove = pos
				break
			}
			if toRemove < 0 {
				log.Printf("[DEBUG]: Policy audit configs for %s did not include an audit config for service %q", updater.DescribeResource(), ac.Service)
				return nil
			}

			p.AuditConfigs = append(p.AuditConfigs[:toRemove], p.AuditConfigs[toRemove+1:]...)
			return nil
		})
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG]: Resource %s is missing or deleted, marking policy audit config as deleted", updater.DescribeResource())
				return nil
			}
			return err
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
			ExemptedMembers: convertStringArr(logConfig["exempted_members"].(*schema.Set).List()),
		}
	}
	return &cloudresourcemanager.AuditConfig{
		AuditLogConfigs: auditLogConfigs,
		Service:         d.Get("service").(string),
	}
}

func flattenAuditLogConfigs(configs []*cloudresourcemanager.AuditLogConfig) *schema.Set {
	res := schema.NewSet(schema.HashResource(iamAuditConfigSchema["audit_log_config"].Elem.(*schema.Resource)), []interface{}{})
	for _, conf := range configs {
		res.Add(map[string]interface{}{
			"log_type":         conf.LogType,
			"exempted_members": schema.NewSet(schema.HashSchema(iamAuditConfigSchema["audit_log_config"].Elem.(*schema.Resource).Schema["exempted_members"].Elem.(*schema.Schema)), convertStringArrToInterface(conf.ExemptedMembers)),
		})
	}
	return res
}
