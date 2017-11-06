// Contains functions that don't really belong anywhere else.

package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// getRegionFromZone returns the region from a zone for Google cloud.
func getRegionFromZone(zone string) string {
	if zone != "" && len(zone) > 2 {
		region := zone[:len(zone)-2]
		return region
	}
	return ""
}

// getRegion reads the "region" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getRegion(d *schema.ResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("region")
	if !ok {
		if config.Region != "" {
			return config.Region, nil
		}
		return "", fmt.Errorf("region: required field is not set")
	}
	return res.(string), nil
}

func getRegionFromInstanceState(is *terraform.InstanceState, config *Config) (string, error) {
	res, ok := is.Attributes["region"]

	if ok && res != "" {
		return res, nil
	}

	if config.Region != "" {
		return config.Region, nil
	}

	return "", fmt.Errorf("region: required field is not set")
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d *schema.ResourceData, config *Config) (string, error) {
	return getProjectFromSchema("project", d, config)
}

func getProjectFromInstanceState(is *terraform.InstanceState, config *Config) (string, error) {
	res, ok := is.Attributes["project"]

	if ok && res != "" {
		return res, nil
	}

	if config.Project != "" {
		return config.Project, nil
	}

	return "", fmt.Errorf("project: required field is not set")
}

func getZonalResourceFromRegion(getResource func(string) (interface{}, error), region string, compute *compute.Service, project string) (interface{}, error) {
	zoneList, err := compute.Zones.List(project).Do()
	if err != nil {
		return nil, err
	}
	var resource interface{}
	for _, zone := range zoneList.Items {
		if strings.Contains(zone.Name, region) {
			resource, err = getResource(zone.Name)
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					// Resource was not found in this zone
					continue
				}
				return nil, fmt.Errorf("Error reading Resource: %s", err)
			}
			// Resource was found
			return resource, nil
		}
	}
	// Resource does not exist in this region
	return nil, nil
}

func getZonalBetaResourceFromRegion(getResource func(string) (interface{}, error), region string, compute *computeBeta.Service, project string) (interface{}, error) {
	zoneList, err := compute.Zones.List(project).Do()
	if err != nil {
		return nil, err
	}
	var resource interface{}
	for _, zone := range zoneList.Items {
		if strings.Contains(zone.Name, region) {
			resource, err = getResource(zone.Name)
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					// Resource was not found in this zone
					continue
				}
				return nil, fmt.Errorf("Error reading Resource: %s", err)
			}
			// Resource was found
			return resource, nil
		}
	}
	// Resource does not exist in this region
	return nil, nil
}

// getNetworkLink takes a "network" field and if the value:
// - is a resource URL, returns the string unchanged
// - is the network name only, then looks up the resource URL using the google client
func getNetworkLink(config *Config, project, network string) (string, error) {
	if network == "" {
		return "", nil
	}

	if strings.HasPrefix(network, "https://www.googleapis.com/compute/") {
		return network, nil
	}

	networkData, err := config.clientCompute.Networks.Get(project, network).Do()
	if err != nil {
		return "", fmt.Errorf("Error reading network: %s", err)
	}
	return networkData.SelfLink, nil
}

// getSubnetworkLink takes the "subnetwork" field and if the value is:
// - a resource URL, returns the string unchanged
// - a subnetwork name, looks up the resource URL using the google client.
//
// If `subnetworkField` is a resource url, `subnetworkProjectField` cannot be set.
// If `subnetworkField` is a subnetwork name, `subnetworkProjectField` will be used
// 	as the project if set. If not, we fallback on the default project.
func getSubnetworkLink(config *Config, defaultProject, region, subnetworkProject, subnetwork string) (string, error) {
	if subnetwork == "" {
		return "", nil
	}

	if regexp.MustCompile(SubnetworkLinkRegex).MatchString(subnetwork) {
		return subnetwork, nil
	}

	project := defaultProject
	if subnetworkProject != "" {
		project = subnetworkProject
	}
	subnetworkData, err := config.clientCompute.Subnetworks.Get(project, region, subnetwork).Do()
	if err != nil {
		return "", fmt.Errorf("Error referencing subnetwork '%s' in region '%s': %s", subnetwork, region, err)
	}
	return subnetworkData.SelfLink, nil
}

// getNetworkName reads the "network" field from the given resource data and if the value:
// - is a resource URL, extracts the network name from the URL and returns it
// - is the network name only (i.e not prefixed with http://www.googleapis.com/compute/...), is returned unchanged
func getNetworkName(d *schema.ResourceData, field string) (string, error) {
	if v, ok := d.GetOk(field); ok {
		network := v.(string)
		return getNetworkNameFromSelfLink(network)
	}
	return "", nil
}

func getNetworkNameFromSelfLink(network string) (string, error) {
	if !strings.HasPrefix(network, "https://www.googleapis.com/compute/") {
		return network, nil
	}
	// extract the network name from SelfLink URL
	networkName := network[strings.LastIndex(network, "/")+1:]
	if networkName == "" {
		return "", fmt.Errorf("network url not valid")
	}
	return networkName, nil
}

func getRouterLockName(region string, router string) string {
	return fmt.Sprintf("router/%s/%s", region, router)
}

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return fmt.Errorf("Error reading %s: %s", resource, err)
}

func isConflictError(err error) bool {
	if e, ok := err.(*googleapi.Error); ok && e.Code == 409 {
		return true
	} else if !ok && errwrap.ContainsType(err, &googleapi.Error{}) {
		e := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
		if e.Code == 409 {
			return true
		}
	}
	return false
}

func linkDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	parts := strings.Split(old, "/")
	if parts[len(parts)-1] == new {
		return true
	}
	return false
}

func optionalPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return prefix+old == new || prefix+new == old
	}
}

func ipCidrRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// The range may be a:
	// A) single IP address (e.g. 10.2.3.4)
	// B) CIDR format string (e.g. 10.1.2.0/24)
	// C) netmask (e.g. /24)
	//
	// For A) and B), no diff to suppress, they have to match completely.
	// For C), The API picks a network IP address and this creates a diff of the form:
	// network_interface.0.alias_ip_range.0.ip_cidr_range: "10.128.1.0/24" => "/24"
	// We should only compare the mask portion for this case.
	if len(new) > 0 && new[0] == '/' {
		oldNetmaskStartPos := strings.LastIndex(old, "/")

		if oldNetmaskStartPos != -1 {
			oldNetmask := old[strings.LastIndex(old, "/"):]
			if oldNetmask == new {
				return true
			}
		}
	}

	return false
}

// Port range '80' and '80-80' is equivalent.
// `old` is read from the server and always has the full range format (e.g. '80-80', '1024-2048').
// `new` can be either a single port or a port range.
func portRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if old == new+"-"+new {
		return true
	}
	return false
}

// Single-digit hour is equivalent to hour with leading zero e.g. suppress diff 1:00 => 01:00.
// Assume either value could be in either format.
func rfc3339TimeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if (len(old) == 4 && "0"+old == new) || (len(new) == 4 && "0"+new == old) {
		return true
	}
	return false
}

// expandLabels pulls the value of "labels" out of a schema.ResourceData as a map[string]string.
func expandLabels(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "labels")
}

// expandStringMap pulls the value of key out of a schema.ResourceData as a map[string]string.
func expandStringMap(d *schema.ResourceData, key string) map[string]string {
	v, ok := d.GetOk(key)

	if !ok {
		return map[string]string{}
	}

	return convertStringMap(v.(map[string]interface{}))
}

func convertStringMap(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, val := range v {
		m[k] = val.(string)
	}
	return m
}

func convertStringArr(ifaceArr []interface{}) []string {
	return convertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

func extractLastResourceFromUri(uri string) string {
	rUris := strings.Split(uri, "/")
	return rUris[len(rUris)-1]
}

func convertStringArrToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}

func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	return s
}

func convertArrToMap(ifaceArr []interface{}) map[string]struct{} {
	sm := make(map[string]struct{})
	for _, s := range ifaceArr {
		sm[s.(string)] = struct{}{}
	}
	return sm
}

func mergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func retry(retryFunc func() error) error {
	return retryTime(retryFunc, 1)
}

func retryTime(retryFunc func() error, minutes int) error {
	return resource.Retry(time.Duration(minutes)*time.Minute, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && (gerr.Code == 429 || gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503) {
			return resource.RetryableError(gerr)
		}
		return resource.NonRetryableError(err)
	})
}

func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	if len(m) == 0 {
		return map[string]interface{}{}
	}

	return m[0].(map[string]interface{})
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

func resourceInstanceMetadata(d *schema.ResourceData) (*compute.Metadata, error) {
	m := &compute.Metadata{}
	mdMap := d.Get("metadata").(map[string]interface{})
	if v, ok := d.GetOk("metadata_startup_script"); ok && v.(string) != "" {
		mdMap["startup-script"] = v
	}
	if len(mdMap) > 0 {
		m.Items = make([]*compute.MetadataItems, 0, len(mdMap))
		for key, val := range mdMap {
			v := val.(string)
			m.Items = append(m.Items, &compute.MetadataItems{
				Key:   key,
				Value: &v,
			})
		}

		// Set the fingerprint. If the metadata has never been set before
		// then this will just be blank.
		m.Fingerprint = d.Get("metadata_fingerprint").(string)
	}

	return m, nil
}

func flattenMetadata(metadata *compute.Metadata) map[string]string {
	metadataMap := make(map[string]string)
	for _, item := range metadata.Items {
		metadataMap[item.Key] = *item.Value
	}
	return metadataMap
}

func resourceInstanceTags(d *schema.ResourceData) *compute.Tags {
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

func flattenScheduling(scheduling *compute.Scheduling) []map[string]interface{} {
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

func getProjectAndRegionFromSubnetworkLink(subnetwork string) (string, string) {
	r := regexp.MustCompile(SubnetworkLinkRegex)
	if !r.MatchString(subnetwork) {
		return "", ""
	}

	matches := r.FindStringSubmatch(subnetwork)
	return matches[1], matches[2]
}

func getProjectFromSubnetworkLink(subnetwork string) string {
	project, _ := getProjectAndRegionFromSubnetworkLink(subnetwork)
	return project
}

func flattenAccessConfigs(accessConfigs []*compute.AccessConfig) ([]map[string]interface{}, string) {
	flattened := make([]map[string]interface{}, len(accessConfigs))
	natIP := ""
	for i, ac := range accessConfigs {
		flattened[i] = map[string]interface{}{
			"nat_ip":          ac.NatIP,
			"assigned_nat_ip": ac.NatIP,
		}
		if natIP == "" {
			natIP = ac.NatIP
		}
	}
	return flattened, natIP
}

func flattenNetworkInterfaces(networkInterfaces []*compute.NetworkInterface) ([]map[string]interface{}, string, string, string) {
	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var region, internalIP, externalIP string

	for i, iface := range networkInterfaces {
		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigs(iface.AccessConfigs)

		var project string
		project, region = getProjectAndRegionFromSubnetworkLink(iface.Subnetwork)

		flattened[i] = map[string]interface{}{
			"address":            iface.NetworkIP,
			"network_ip":         iface.NetworkIP,
			"network":            iface.Network,
			"subnetwork":         iface.Subnetwork,
			"subnetwork_project": project,
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
	return flattened, region, internalIP, externalIP
}

func expandAccessConfigs(configs []interface{}) []*compute.AccessConfig {
	acs := make([]*compute.AccessConfig, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})
		acs[i] = &compute.AccessConfig{
			Type:  "ONE_TO_ONE_NAT",
			NatIP: data["nat_ip"].(string),
		}
	}
	return acs
}

func expandNetworkInterfaces(d *schema.ResourceData, config *Config) ([]*compute.NetworkInterface, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	configs := d.Get("network_interface").([]interface{})
	ifaces := make([]*compute.NetworkInterface, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		network := data["network"].(string)
		subnetwork := data["subnetwork"].(string)
		if (network == "" && subnetwork == "") || (network != "" && subnetwork != "") {
			return nil, fmt.Errorf("exactly one of network or subnetwork must be provided")
		}

		networkLink, err := getNetworkLink(config, project, network)
		if err != nil {
			return nil, fmt.Errorf("cannot determine selflink for subnetwork '%s': %s", subnetwork, err)
		}

		subnetworkProject := data["subnetwork_project"].(string)
		subnetLink, err := getSubnetworkLink(config, project, region, subnetworkProject, subnetwork)
		if err != nil {
			return nil, fmt.Errorf("cannot determine selflink for subnetwork '%s': %s", subnetwork, err)
		}

		ifaces[i] = &compute.NetworkInterface{
			NetworkIP:     data["network_ip"].(string),
			Network:       networkLink,
			Subnetwork:    subnetLink,
			AccessConfigs: expandAccessConfigs(data["access_config"].([]interface{})),
			AliasIpRanges: expandAliasIpRanges(data["alias_ip_range"].([]interface{})),
		}

		// network_ip is deprecated. We want address to win if both are set.
		if data["address"].(string) != "" {
			ifaces[i].NetworkIP = data["address"].(string)
		}

	}
	return ifaces, nil
}

func flattenServiceAccounts(serviceAccounts []*compute.ServiceAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, serviceAccount := range serviceAccounts {
		result[i] = map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": schema.NewSet(stringScopeHashcode, convertStringArrToInterface(serviceAccount.Scopes)),
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
			Scopes: canonicalizeServiceScopes(convertStringSet(data["scopes"].(*schema.Set))),
		}

		if accounts[i].Email == "" {
			accounts[i].Email = "default"
		}
	}
	return accounts
}

func expandGuestAccelerators(zone string, configs []interface{}) []*compute.AcceleratorConfig {
	guestAccelerators := make([]*compute.AcceleratorConfig, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})
		guestAccelerators[i] = &compute.AcceleratorConfig{
			AcceleratorCount: int64(data["count"].(int)),
			AcceleratorType:  acceleratorRef(zone, data["type"].(string)),
		}
	}

	return guestAccelerators
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

// Instances want a partial URL, but instance templates want the bare
// accelerator name without zone (despite the docs saying otherwise).
//
// Using a partial URL on an instance template results in:
// Invalid value for field 'resource.properties.guestAccelerators[0].acceleratorType':
// 'zones/us-east1-b/acceleratorTypes/nvidia-tesla-k80'.
// Accelerator type 'zones/us-east1-b/acceleratorTypes/nvidia-tesla-k80'
// must be a valid resource name (not an url).
func acceleratorRef(zone, accelerator string) string {
	if strings.HasPrefix(accelerator, "zones/") || zone == "" {
		return accelerator
	}
	return fmt.Sprintf("zones/%s/acceleratorTypes/%s", zone, accelerator)
}
