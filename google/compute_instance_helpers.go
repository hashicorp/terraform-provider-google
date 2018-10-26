package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

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

func flattenScheduling(scheduling *computeBeta.Scheduling) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	schedulingMap := map[string]interface{}{
		"on_host_maintenance": scheduling.OnHostMaintenance,
		"preemptible":         scheduling.Preemptible,
	}
	if scheduling.AutomaticRestart != nil {
		schedulingMap["automatic_restart"] = *scheduling.AutomaticRestart
	}
	result = append(result, schedulingMap)
	return result
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

func expandNetworkInterfaces(d *schema.ResourceData, config *Config) ([]*computeBeta.NetworkInterface, error) {
	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*computeBeta.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)
		if (network == "" && subnetwork == "") || (network != "" && subnetwork != "") {
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

func resourceInstanceTags(d *schema.ResourceData) *computeBeta.Tags {
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
