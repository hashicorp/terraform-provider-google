package google

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/googleapi"
)

func instanceSchedulingNodeAffinitiesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operator": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IN", "NOT"}, false),
			},
			"values": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func expandAliasIpRanges(ranges []interface{}) []*computeBeta.AliasIpRange {
	ipRanges := make([]*computeBeta.AliasIpRange, 0, len(ranges))
	for _, raw := range ranges {
		data := raw.(map[string]interface{})
		ipRanges = append(ipRanges, &computeBeta.AliasIpRange{
			IpCidrRange:         data["ip_cidr_range"].(string),
			SubnetworkRangeName: data["subnetwork_range_name"].(string),
		})
	}
	return ipRanges
}

func flattenAliasIpRange(ranges []*computeBeta.AliasIpRange) []map[string]interface{} {
	rangesSchema := make([]map[string]interface{}, 0, len(ranges))
	for _, ipRange := range ranges {
		rangesSchema = append(rangesSchema, map[string]interface{}{
			"ip_cidr_range":         ipRange.IpCidrRange,
			"subnetwork_range_name": ipRange.SubnetworkRangeName,
		})
	}
	return rangesSchema
}

func expandScheduling(v interface{}) (*computeBeta.Scheduling, error) {
	if v == nil {
		// We can't set default values for lists.
		return &computeBeta.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}, nil
	}

	ls := v.([]interface{})
	if len(ls) == 0 {
		// We can't set default values for lists
		return &computeBeta.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}, nil
	}

	if len(ls) > 1 || ls[0] == nil {
		return nil, fmt.Errorf("expected exactly one scheduling block")
	}

	original := ls[0].(map[string]interface{})
	scheduling := &computeBeta.Scheduling{
		ForceSendFields: make([]string, 0, 4),
	}

	if v, ok := original["automatic_restart"]; ok {
		scheduling.AutomaticRestart = googleapi.Bool(v.(bool))
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "AutomaticRestart")
	}

	if v, ok := original["preemptible"]; ok {
		scheduling.Preemptible = v.(bool)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "Preemptible")

	}

	if v, ok := original["on_host_maintenance"]; ok {
		scheduling.OnHostMaintenance = v.(string)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "OnHostMaintenance")
	}

	if v, ok := original["node_affinities"]; ok && v != nil {
		naSet := v.(*schema.Set).List()
		scheduling.NodeAffinities = make([]*computeBeta.SchedulingNodeAffinity, len(ls))
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "NodeAffinities")
		for _, nodeAffRaw := range naSet {
			if nodeAffRaw == nil {
				continue
			}
			nodeAff := nodeAffRaw.(map[string]interface{})
			transformed := &computeBeta.SchedulingNodeAffinity{
				Key:      nodeAff["key"].(string),
				Operator: nodeAff["operator"].(string),
				Values:   convertStringArr(nodeAff["values"].(*schema.Set).List()),
			}
			scheduling.NodeAffinities = append(scheduling.NodeAffinities, transformed)
		}
	}

	return scheduling, nil
}

func flattenScheduling(resp *computeBeta.Scheduling) []map[string]interface{} {
	schedulingMap := map[string]interface{}{
		"on_host_maintenance": resp.OnHostMaintenance,
		"preemptible":         resp.Preemptible,
	}

	if resp.AutomaticRestart != nil {
		schedulingMap["automatic_restart"] = *resp.AutomaticRestart
	}

	nodeAffinities := schema.NewSet(schema.HashResource(instanceSchedulingNodeAffinitiesElemSchema()), nil)
	for _, na := range resp.NodeAffinities {
		nodeAffinities.Add(map[string]interface{}{
			"key":      na.Key,
			"operator": na.Operator,
			"values":   schema.NewSet(schema.HashString, convertStringArrToInterface(na.Values)),
		})
	}
	schedulingMap["node_affinities"] = nodeAffinities

	return []map[string]interface{}{schedulingMap}
}

func flattenAccessConfigs(accessConfigs []*computeBeta.AccessConfig) ([]map[string]interface{}, string) {
	flattened := make([]map[string]interface{}, len(accessConfigs))
	natIP := ""
	for i, ac := range accessConfigs {
		flattened[i] = map[string]interface{}{
			"nat_ip":       ac.NatIP,
			"network_tier": ac.NetworkTier,
		}
		if ac.SetPublicPtr {
			flattened[i]["public_ptr_domain_name"] = ac.PublicPtrDomainName
		}
		if natIP == "" {
			natIP = ac.NatIP
		}
	}
	return flattened, natIP
}

func flattenNetworkInterfaces(d *schema.ResourceData, config *Config, networkInterfaces []*computeBeta.NetworkInterface) ([]map[string]interface{}, string, string, string, error) {
	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var region, internalIP, externalIP string

	for i, iface := range networkInterfaces {
		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigs(iface.AccessConfigs)

		subnet, err := ParseSubnetworkFieldValue(iface.Subnetwork, d, config)
		if err != nil {
			return nil, "", "", "", err
		}
		region = subnet.Region

		flattened[i] = map[string]interface{}{
			"network_ip":         iface.NetworkIP,
			"network":            ConvertSelfLinkToV1(iface.Network),
			"subnetwork":         ConvertSelfLinkToV1(iface.Subnetwork),
			"subnetwork_project": subnet.Project,
			"access_config":      ac,
			"alias_ip_range":     flattenAliasIpRange(iface.AliasIpRanges),
		}
		// Instance template interfaces never have names, so they're absent
		// in the instance template network_interface schema. We want to use the
		// same flattening code for both resource types, so we avoid trying to
		// set the name field when it's not set at the GCE end.
		if iface.Name != "" {
			flattened[i]["name"] = iface.Name
		}
		if internalIP == "" {
			internalIP = iface.NetworkIP
		}
	}
	return flattened, region, internalIP, externalIP, nil
}

func expandAccessConfigs(configs []interface{}) []*computeBeta.AccessConfig {
	acs := make([]*computeBeta.AccessConfig, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})
		acs[i] = &computeBeta.AccessConfig{
			Type:        "ONE_TO_ONE_NAT",
			NatIP:       data["nat_ip"].(string),
			NetworkTier: data["network_tier"].(string),
		}
		if ptr, ok := data["public_ptr_domain_name"]; ok && ptr != "" {
			acs[i].SetPublicPtr = true
			acs[i].PublicPtrDomainName = ptr.(string)
		}
	}
	return acs
}

func expandNetworkInterfaces(d TerraformResourceData, config *Config) ([]*computeBeta.NetworkInterface, error) {
	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*computeBeta.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)
		if network == "" && subnetwork == "" {
			return nil, fmt.Errorf("exactly one of network or subnetwork must be provided")
		}

		nf, err := ParseNetworkFieldValue(network, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", network, err)
		}

		subnetProjectField := fmt.Sprintf("network_interface.%d.subnetwork_project", i)
		sf, err := ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", subnetwork, err)
		}

		ifaces[i] = &computeBeta.NetworkInterface{
			NetworkIP:     data["network_ip"].(string),
			Network:       nf.RelativeLink(),
			Subnetwork:    sf.RelativeLink(),
			AccessConfigs: expandAccessConfigs(data["access_config"].([]interface{})),
			AliasIpRanges: expandAliasIpRanges(data["alias_ip_range"].([]interface{})),
		}

	}
	return ifaces, nil
}

func flattenServiceAccounts(serviceAccounts []*computeBeta.ServiceAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, serviceAccount := range serviceAccounts {
		result[i] = map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": schema.NewSet(stringScopeHashcode, convertStringArrToInterface(serviceAccount.Scopes)),
		}
	}
	return result
}

func expandServiceAccounts(configs []interface{}) []*computeBeta.ServiceAccount {
	accounts := make([]*computeBeta.ServiceAccount, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		accounts[i] = &computeBeta.ServiceAccount{
			Email:  data["email"].(string),
			Scopes: canonicalizeServiceScopes(convertStringSet(data["scopes"].(*schema.Set))),
		}

		if accounts[i].Email == "" {
			accounts[i].Email = "default"
		}
	}
	return accounts
}

func flattenGuestAccelerators(accelerators []*computeBeta.AcceleratorConfig) []map[string]interface{} {
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, accelerator := range accelerators {
		acceleratorsSchema[i] = map[string]interface{}{
			"count": accelerator.AcceleratorCount,
			"type":  accelerator.AcceleratorType,
		}
	}
	return acceleratorsSchema
}

func resourceInstanceTags(d TerraformResourceData) *computeBeta.Tags {
	// Calculate the tags
	var tags *computeBeta.Tags
	if v := d.Get("tags"); v != nil {
		vs := v.(*schema.Set)
		tags = new(computeBeta.Tags)
		tags.Items = make([]string, vs.Len())
		for i, v := range vs.List() {
			tags.Items[i] = v.(string)
		}

		tags.Fingerprint = d.Get("tags_fingerprint").(string)
	}

	return tags
}

func expandShieldedVmConfigs(d TerraformResourceData) *computeBeta.ShieldedVmConfig {
	if _, ok := d.GetOk("shielded_instance_config"); !ok {
		return nil
	}

	prefix := "shielded_instance_config.0"
	return &computeBeta.ShieldedVmConfig{
		EnableSecureBoot:          d.Get(prefix + ".enable_secure_boot").(bool),
		EnableVtpm:                d.Get(prefix + ".enable_vtpm").(bool),
		EnableIntegrityMonitoring: d.Get(prefix + ".enable_integrity_monitoring").(bool),
		ForceSendFields:           []string{"EnableSecureBoot", "EnableVtpm", "EnableIntegrityMonitoring"},
	}
}

func flattenShieldedVmConfig(shieldedVmConfig *computeBeta.ShieldedVmConfig) []map[string]bool {
	if shieldedVmConfig == nil {
		return nil
	}

	return []map[string]bool{{
		"enable_secure_boot":          shieldedVmConfig.EnableSecureBoot,
		"enable_vtpm":                 shieldedVmConfig.EnableVtpm,
		"enable_integrity_monitoring": shieldedVmConfig.EnableIntegrityMonitoring,
	}}
}

func expandDisplayDevice(d TerraformResourceData) *computeBeta.DisplayDevice {
	if _, ok := d.GetOk("enable_display"); !ok {
		return nil
	}
	return &computeBeta.DisplayDevice{
		EnableDisplay:   d.Get("enable_display").(bool),
		ForceSendFields: []string{"EnableDisplay"},
	}
}

func flattenEnableDisplay(displayDevice *computeBeta.DisplayDevice) interface{} {
	if displayDevice == nil {
		return nil
	}

	return displayDevice.EnableDisplay
}

// Terraform doesn't correctly calculate changes on schema.Set, so we do it manually
// https://github.com/hashicorp/terraform-plugin-sdk/issues/98
func schedulingHasChange(d *schema.ResourceData) bool {
	if !d.HasChange("scheduling") {
		// This doesn't work correctly, which is why this method exists
		// But it is here for posterity
		return false
	}
	o, n := d.GetChange("scheduling")
	oScheduling := o.([]interface{})[0].(map[string]interface{})
	newScheduling := n.([]interface{})[0].(map[string]interface{})
	originalNa := oScheduling["node_affinities"].(*schema.Set)
	newNa := newScheduling["node_affinities"].(*schema.Set)
	if oScheduling["automatic_restart"] != newScheduling["automatic_restart"] {
		return true
	}

	if oScheduling["preemptible"] != newScheduling["preemptible"] {
		return true
	}

	if oScheduling["on_host_maintenance"] != newScheduling["on_host_maintenance"] {
		return true
	}

	return reflect.DeepEqual(newNa, originalNa)
}
