// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/compute/v1"
)

func instanceSchedulingNodeAffinitiesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN"}, false),
			},
			"values": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func expandAliasIpRanges(ranges []interface{}) []*compute.AliasIpRange {
	ipRanges := make([]*compute.AliasIpRange, 0, len(ranges))
	for _, raw := range ranges {
		data := raw.(map[string]interface{})
		ipRanges = append(ipRanges, &compute.AliasIpRange{
			IpCidrRange:         data["ip_cidr_range"].(string),
			SubnetworkRangeName: data["subnetwork_range_name"].(string),
		})
	}
	return ipRanges
}

func flattenAliasIpRange(ranges []*compute.AliasIpRange) []map[string]interface{} {
	rangesSchema := make([]map[string]interface{}, 0, len(ranges))
	for _, ipRange := range ranges {
		rangesSchema = append(rangesSchema, map[string]interface{}{
			"ip_cidr_range":         ipRange.IpCidrRange,
			"subnetwork_range_name": ipRange.SubnetworkRangeName,
		})
	}
	return rangesSchema
}

func expandScheduling(v interface{}) (*compute.Scheduling, error) {
	if v == nil {
		// We can't set default values for lists.
		return &compute.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}, nil
	}

	ls := v.([]interface{})
	if len(ls) == 0 {
		// We can't set default values for lists
		return &compute.Scheduling{
			AutomaticRestart: googleapi.Bool(true),
		}, nil
	}

	if len(ls) > 1 || ls[0] == nil {
		return nil, fmt.Errorf("expected exactly one scheduling block")
	}

	original := ls[0].(map[string]interface{})
	scheduling := &compute.Scheduling{
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
		scheduling.NodeAffinities = make([]*compute.SchedulingNodeAffinity, len(ls))
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "NodeAffinities")
		for _, nodeAffRaw := range naSet {
			if nodeAffRaw == nil {
				continue
			}
			nodeAff := nodeAffRaw.(map[string]interface{})
			transformed := &compute.SchedulingNodeAffinity{
				Key:      nodeAff["key"].(string),
				Operator: nodeAff["operator"].(string),
				Values:   tpgresource.ConvertStringArr(nodeAff["values"].(*schema.Set).List()),
			}
			scheduling.NodeAffinities = append(scheduling.NodeAffinities, transformed)
		}
	}

	if v, ok := original["min_node_cpus"]; ok {
		scheduling.MinNodeCpus = int64(v.(int))
	}
	if v, ok := original["provisioning_model"]; ok {
		scheduling.ProvisioningModel = v.(string)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "ProvisioningModel")
	}
	if v, ok := original["instance_termination_action"]; ok {
		scheduling.InstanceTerminationAction = v.(string)
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "InstanceTerminationAction")
	}
	if v, ok := original["local_ssd_recovery_timeout"]; ok {
		transformedLocalSsdRecoveryTimeout, err := expandComputeLocalSsdRecoveryTimeout(v)
		if err != nil {
			return nil, err
		}
		scheduling.LocalSsdRecoveryTimeout = transformedLocalSsdRecoveryTimeout
		scheduling.ForceSendFields = append(scheduling.ForceSendFields, "LocalSsdRecoveryTimeout")
	}
	return scheduling, nil
}

func expandComputeLocalSsdRecoveryTimeout(v interface{}) (*compute.Duration, error) {
	l := v.([]interface{})
	duration := compute.Duration{}
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	transformedNanos, err := expandComputeLocalSsdRecoveryTimeoutNanos(original["nanos"])
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNanos); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Nanos = int64(transformedNanos.(int))
	}

	transformedSeconds, err := expandComputeLocalSsdRecoveryTimeoutSeconds(original["seconds"])
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSeconds); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		duration.Seconds = int64(transformedSeconds.(int))
	}
	return &duration, nil
}

func expandComputeLocalSsdRecoveryTimeoutNanos(v interface{}) (interface{}, error) {
	return v, nil
}

func expandComputeLocalSsdRecoveryTimeoutSeconds(v interface{}) (interface{}, error) {
	return v, nil
}

func flattenScheduling(resp *compute.Scheduling) []map[string]interface{} {
	schedulingMap := map[string]interface{}{
		"on_host_maintenance":         resp.OnHostMaintenance,
		"preemptible":                 resp.Preemptible,
		"min_node_cpus":               resp.MinNodeCpus,
		"provisioning_model":          resp.ProvisioningModel,
		"instance_termination_action": resp.InstanceTerminationAction,
	}

	if resp.AutomaticRestart != nil {
		schedulingMap["automatic_restart"] = *resp.AutomaticRestart
	}

	if resp.LocalSsdRecoveryTimeout != nil {
		schedulingMap["local_ssd_recovery_timeout"] = flattenComputeLocalSsdRecoveryTimeout(resp.LocalSsdRecoveryTimeout)
	}

	nodeAffinities := schema.NewSet(schema.HashResource(instanceSchedulingNodeAffinitiesElemSchema()), nil)
	for _, na := range resp.NodeAffinities {
		nodeAffinities.Add(map[string]interface{}{
			"key":      na.Key,
			"operator": na.Operator,
			"values":   schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(na.Values)),
		})
	}
	schedulingMap["node_affinities"] = nodeAffinities

	return []map[string]interface{}{schedulingMap}
}

func flattenComputeLocalSsdRecoveryTimeout(v *compute.Duration) []interface{} {
	if v == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["nanos"] = v.Nanos
	transformed["seconds"] = v.Seconds
	return []interface{}{transformed}
}

func flattenAccessConfigs(accessConfigs []*compute.AccessConfig) ([]map[string]interface{}, string) {
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

func flattenIpv6AccessConfigs(ipv6AccessConfigs []*compute.AccessConfig) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(ipv6AccessConfigs))
	for i, ac := range ipv6AccessConfigs {
		flattened[i] = map[string]interface{}{
			"network_tier": ac.NetworkTier,
		}
		flattened[i]["public_ptr_domain_name"] = ac.PublicPtrDomainName
		flattened[i]["external_ipv6"] = ac.ExternalIpv6
		flattened[i]["external_ipv6_prefix_length"] = strconv.FormatInt(ac.ExternalIpv6PrefixLength, 10)
		flattened[i]["name"] = ac.Name
	}
	return flattened
}

func flattenNetworkInterfaces(d *schema.ResourceData, config *transport_tpg.Config, networkInterfaces []*compute.NetworkInterface) ([]map[string]interface{}, string, string, string, error) {
	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var region, internalIP, externalIP string

	for i, iface := range networkInterfaces {
		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigs(iface.AccessConfigs)

		subnet, err := tpgresource.ParseSubnetworkFieldValue(iface.Subnetwork, d, config)
		if err != nil {
			return nil, "", "", "", err
		}
		region = subnet.Region

		flattened[i] = map[string]interface{}{
			"network_ip":         iface.NetworkIP,
			"network":            tpgresource.ConvertSelfLinkToV1(iface.Network),
			"subnetwork":         tpgresource.ConvertSelfLinkToV1(iface.Subnetwork),
			"subnetwork_project": subnet.Project,
			"access_config":      ac,
			"alias_ip_range":     flattenAliasIpRange(iface.AliasIpRanges),
			"nic_type":           iface.NicType,
			"stack_type":         iface.StackType,
			"ipv6_access_config": flattenIpv6AccessConfigs(iface.Ipv6AccessConfigs),
			"queue_count":        iface.QueueCount,
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

func expandAccessConfigs(configs []interface{}) []*compute.AccessConfig {
	acs := make([]*compute.AccessConfig, len(configs))
	for i, raw := range configs {
		acs[i] = &compute.AccessConfig{}
		acs[i].Type = "ONE_TO_ONE_NAT"
		if raw != nil {
			data := raw.(map[string]interface{})
			acs[i].NatIP = data["nat_ip"].(string)
			acs[i].NetworkTier = data["network_tier"].(string)
			if ptr, ok := data["public_ptr_domain_name"]; ok && ptr != "" {
				acs[i].SetPublicPtr = true
				acs[i].PublicPtrDomainName = ptr.(string)
			}
		}
	}
	return acs
}

func expandIpv6AccessConfigs(configs []interface{}) []*compute.AccessConfig {
	iacs := make([]*compute.AccessConfig, len(configs))
	for i, raw := range configs {
		iacs[i] = &compute.AccessConfig{}
		if raw != nil {
			data := raw.(map[string]interface{})
			iacs[i].NetworkTier = data["network_tier"].(string)
			if ptr, ok := data["public_ptr_domain_name"]; ok && ptr != "" {
				iacs[i].PublicPtrDomainName = ptr.(string)
			}
			if eip, ok := data["external_ipv6"]; ok && eip != "" {
				iacs[i].ExternalIpv6 = eip.(string)
			}
			if eipl, ok := data["external_ipv6_prefix_length"]; ok && eipl != "" {
				if strVal, ok := eipl.(string); ok {
					if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
						iacs[i].ExternalIpv6PrefixLength = intVal
					}
				}
			}
			if name, ok := data["name"]; ok && name != "" {
				iacs[i].Name = name.(string)
			}
			iacs[i].Type = "DIRECT_IPV6" // Currently only type supported
		}
	}
	return iacs
}

func expandNetworkInterfaces(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]*compute.NetworkInterface, error) {
	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*compute.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)
		if network == "" && subnetwork == "" {
			return nil, fmt.Errorf("exactly one of network or subnetwork must be provided")
		}

		nf, err := tpgresource.ParseNetworkFieldValue(network, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", network, err)
		}

		subnetProjectField := fmt.Sprintf("network_interface.%d.subnetwork_project", i)
		sf, err := tpgresource.ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", subnetwork, err)
		}

		ifaces[i] = &compute.NetworkInterface{
			NetworkIP:         data["network_ip"].(string),
			Network:           nf.RelativeLink(),
			Subnetwork:        sf.RelativeLink(),
			AccessConfigs:     expandAccessConfigs(data["access_config"].([]interface{})),
			AliasIpRanges:     expandAliasIpRanges(data["alias_ip_range"].([]interface{})),
			NicType:           data["nic_type"].(string),
			StackType:         data["stack_type"].(string),
			QueueCount:        int64(data["queue_count"].(int)),
			Ipv6AccessConfigs: expandIpv6AccessConfigs(data["ipv6_access_config"].([]interface{})),
		}
	}
	return ifaces, nil
}

func flattenServiceAccounts(serviceAccounts []*compute.ServiceAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, serviceAccount := range serviceAccounts {
		result[i] = map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": schema.NewSet(tpgresource.StringScopeHashcode, tpgresource.ConvertStringArrToInterface(serviceAccount.Scopes)),
		}
	}
	return result
}

func expandServiceAccounts(configs []interface{}) []*compute.ServiceAccount {
	accounts := make([]*compute.ServiceAccount, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		accounts[i] = &compute.ServiceAccount{
			Email:  data["email"].(string),
			Scopes: tpgresource.CanonicalizeServiceScopes(tpgresource.ConvertStringSet(data["scopes"].(*schema.Set))),
		}

		if accounts[i].Email == "" {
			accounts[i].Email = "default"
		}
	}
	return accounts
}

func flattenGuestAccelerators(accelerators []*compute.AcceleratorConfig) []map[string]interface{} {
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, accelerator := range accelerators {
		acceleratorsSchema[i] = map[string]interface{}{
			"count": accelerator.AcceleratorCount,
			"type":  accelerator.AcceleratorType,
		}
	}
	return acceleratorsSchema
}

func resourceInstanceTags(d tpgresource.TerraformResourceData) *compute.Tags {
	// Calculate the tags
	var tags *compute.Tags
	if v := d.Get("tags"); v != nil {
		vs := v.(*schema.Set)
		tags = new(compute.Tags)
		tags.Items = make([]string, vs.Len())
		for i, v := range vs.List() {
			tags.Items[i] = v.(string)
		}

		tags.Fingerprint = d.Get("tags_fingerprint").(string)
	}

	return tags
}

func expandShieldedVmConfigs(d tpgresource.TerraformResourceData) *compute.ShieldedInstanceConfig {
	if _, ok := d.GetOk("shielded_instance_config"); !ok {
		return nil
	}

	prefix := "shielded_instance_config.0"
	return &compute.ShieldedInstanceConfig{
		EnableSecureBoot:          d.Get(prefix + ".enable_secure_boot").(bool),
		EnableVtpm:                d.Get(prefix + ".enable_vtpm").(bool),
		EnableIntegrityMonitoring: d.Get(prefix + ".enable_integrity_monitoring").(bool),
		ForceSendFields:           []string{"EnableSecureBoot", "EnableVtpm", "EnableIntegrityMonitoring"},
	}
}

func expandConfidentialInstanceConfig(d tpgresource.TerraformResourceData) *compute.ConfidentialInstanceConfig {
	if _, ok := d.GetOk("confidential_instance_config"); !ok {
		return nil
	}

	prefix := "confidential_instance_config.0"
	return &compute.ConfidentialInstanceConfig{
		EnableConfidentialCompute: d.Get(prefix + ".enable_confidential_compute").(bool),
	}
}

func flattenConfidentialInstanceConfig(ConfidentialInstanceConfig *compute.ConfidentialInstanceConfig) []map[string]bool {
	if ConfidentialInstanceConfig == nil {
		return nil
	}

	return []map[string]bool{{
		"enable_confidential_compute": ConfidentialInstanceConfig.EnableConfidentialCompute,
	}}
}

func expandAdvancedMachineFeatures(d tpgresource.TerraformResourceData) *compute.AdvancedMachineFeatures {
	if _, ok := d.GetOk("advanced_machine_features"); !ok {
		return nil
	}

	prefix := "advanced_machine_features.0"
	return &compute.AdvancedMachineFeatures{
		EnableNestedVirtualization: d.Get(prefix + ".enable_nested_virtualization").(bool),
		ThreadsPerCore:             int64(d.Get(prefix + ".threads_per_core").(int)),
		VisibleCoreCount:           int64(d.Get(prefix + ".visible_core_count").(int)),
	}
}

func flattenAdvancedMachineFeatures(AdvancedMachineFeatures *compute.AdvancedMachineFeatures) []map[string]interface{} {
	if AdvancedMachineFeatures == nil {
		return nil
	}
	return []map[string]interface{}{{
		"enable_nested_virtualization": AdvancedMachineFeatures.EnableNestedVirtualization,
		"threads_per_core":             AdvancedMachineFeatures.ThreadsPerCore,
		"visible_core_count":           AdvancedMachineFeatures.VisibleCoreCount,
	}}
}

func flattenShieldedVmConfig(shieldedVmConfig *compute.ShieldedInstanceConfig) []map[string]bool {
	if shieldedVmConfig == nil {
		return nil
	}

	return []map[string]bool{{
		"enable_secure_boot":          shieldedVmConfig.EnableSecureBoot,
		"enable_vtpm":                 shieldedVmConfig.EnableVtpm,
		"enable_integrity_monitoring": shieldedVmConfig.EnableIntegrityMonitoring,
	}}
}

func expandDisplayDevice(d tpgresource.TerraformResourceData) *compute.DisplayDevice {
	if _, ok := d.GetOk("enable_display"); !ok {
		return nil
	}
	return &compute.DisplayDevice{
		EnableDisplay:   d.Get("enable_display").(bool),
		ForceSendFields: []string{"EnableDisplay"},
	}
}

func flattenEnableDisplay(displayDevice *compute.DisplayDevice) interface{} {
	if displayDevice == nil {
		return nil
	}

	return displayDevice.EnableDisplay
}

// Node affinity updates require a reboot
func schedulingHasChangeRequiringReboot(d *schema.ResourceData) bool {
	o, n := d.GetChange("scheduling")
	oScheduling := o.([]interface{})[0].(map[string]interface{})
	newScheduling := n.([]interface{})[0].(map[string]interface{})

	return hasNodeAffinitiesChanged(oScheduling, newScheduling)
}

// Terraform doesn't correctly calculate changes on schema.Set, so we do it manually
// https://github.com/hashicorp/terraform-plugin-sdk/issues/98
func schedulingHasChangeWithoutReboot(d *schema.ResourceData) bool {
	if !d.HasChange("scheduling") {
		// This doesn't work correctly, which is why this method exists
		// But it is here for posterity
		return false
	}
	o, n := d.GetChange("scheduling")
	oScheduling := o.([]interface{})[0].(map[string]interface{})
	newScheduling := n.([]interface{})[0].(map[string]interface{})

	if schedulingHasChangeRequiringReboot(d) {
		return false
	}

	if oScheduling["automatic_restart"] != newScheduling["automatic_restart"] {
		return true
	}

	if oScheduling["preemptible"] != newScheduling["preemptible"] {
		return true
	}

	if oScheduling["on_host_maintenance"] != newScheduling["on_host_maintenance"] {
		return true
	}

	if oScheduling["min_node_cpus"] != newScheduling["min_node_cpus"] {
		return true
	}

	if oScheduling["provisioning_model"] != newScheduling["provisioning_model"] {
		return true
	}

	if oScheduling["instance_termination_action"] != newScheduling["instance_termination_action"] {
		return true
	}

	return false
}

func hasNodeAffinitiesChanged(oScheduling, newScheduling map[string]interface{}) bool {
	oldNAs := oScheduling["node_affinities"].(*schema.Set).List()
	newNAs := newScheduling["node_affinities"].(*schema.Set).List()
	if len(oldNAs) != len(newNAs) {
		return true
	}
	for i := range oldNAs {
		oldNodeAffinity := oldNAs[i].(map[string]interface{})
		newNodeAffinity := newNAs[i].(map[string]interface{})
		if oldNodeAffinity["key"] != newNodeAffinity["key"] {
			return true
		}
		if oldNodeAffinity["operator"] != newNodeAffinity["operator"] {
			return true
		}

		// ConvertStringSet will sort the set into a slice, allowing DeepEqual
		if !reflect.DeepEqual(tpgresource.ConvertStringSet(oldNodeAffinity["values"].(*schema.Set)), tpgresource.ConvertStringSet(newNodeAffinity["values"].(*schema.Set))) {
			return true
		}
	}

	return false
}

func expandReservationAffinity(d *schema.ResourceData) (*compute.ReservationAffinity, error) {
	_, ok := d.GetOk("reservation_affinity")
	if !ok {
		return nil, nil
	}

	prefix := "reservation_affinity.0"
	reservationAffinityType := d.Get(prefix + ".type").(string)

	affinity := compute.ReservationAffinity{
		ConsumeReservationType: reservationAffinityType,
		ForceSendFields:        []string{"ConsumeReservationType"},
	}

	_, hasSpecificReservation := d.GetOk(prefix + ".specific_reservation")
	if (reservationAffinityType == "SPECIFIC_RESERVATION") != hasSpecificReservation {
		return nil, fmt.Errorf("specific_reservation must be set when reservation_affinity is SPECIFIC_RESERVATION, and not set otherwise")
	}

	prefix = prefix + ".specific_reservation.0"
	if hasSpecificReservation {
		affinity.Key = d.Get(prefix + ".key").(string)
		affinity.ForceSendFields = append(affinity.ForceSendFields, "Key", "Values")

		for _, v := range d.Get(prefix + ".values").([]interface{}) {
			affinity.Values = append(affinity.Values, v.(string))
		}
	}

	return &affinity, nil
}

func flattenReservationAffinity(affinity *compute.ReservationAffinity) []map[string]interface{} {
	if affinity == nil {
		return nil
	}

	flattened := map[string]interface{}{
		"type": affinity.ConsumeReservationType,
	}

	if affinity.ConsumeReservationType == "SPECIFIC_RESERVATION" {
		flattened["specific_reservation"] = []map[string]interface{}{{
			"key":    affinity.Key,
			"values": affinity.Values,
		}}
	}

	return []map[string]interface{}{flattened}
}

func expandNetworkPerformanceConfig(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*compute.NetworkPerformanceConfig, error) {
	configs, ok := d.GetOk("network_performance_config")
	if !ok {
		return nil, nil
	}

	npcSlice := configs.([]interface{})
	if len(npcSlice) > 1 {
		return nil, fmt.Errorf("cannot specify multiple network_performance_configs")
	}

	if len(npcSlice) == 0 || npcSlice[0] == nil {
		return nil, nil
	}
	npc := npcSlice[0].(map[string]interface{})
	return &compute.NetworkPerformanceConfig{
		TotalEgressBandwidthTier: npc["total_egress_bandwidth_tier"].(string),
	}, nil
}

func flattenNetworkPerformanceConfig(c *compute.NetworkPerformanceConfig) []map[string]interface{} {
	if c == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"total_egress_bandwidth_tier": c.TotalEgressBandwidthTier,
		},
	}
}
