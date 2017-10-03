// Contains functions that don't really belong anywhere else.

package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/errwrap"
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
	res, ok := d.GetOk("project")
	if !ok {
		if config.Project != "" {
			return config.Project, nil
		}
		return "", fmt.Errorf("project: required field is not set")
	}
	return res.(string), nil
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

// getNetworkLink reads the "network" field from the given resource data and if the value:
// - is a resource URL, returns the string unchanged
// - is the network name only, then looks up the resource URL using the google client
func getNetworkLink(d *schema.ResourceData, config *Config, field string) (string, error) {
	if v, ok := d.GetOk(field); ok {
		network := v.(string)

		project, err := getProject(d, config)
		if err != nil {
			return "", err
		}

		if !strings.HasPrefix(network, "https://www.googleapis.com/compute/") {
			// Network value provided is just the name, lookup the network SelfLink
			networkData, err := config.clientCompute.Networks.Get(
				project, network).Do()
			if err != nil {
				return "", fmt.Errorf("Error reading network: %s", err)
			}
			network = networkData.SelfLink
		}

		return network, nil

	} else {
		return "", nil
	}
}

// Reads the "subnetwork" fields from the given resource data and if the value is:
// - a resource URL, returns the string unchanged
// - a subnetwork name, looks up the resource URL using the google client.
//
// If `subnetworkField` is a resource url, `subnetworkProjectField` cannot be set.
// If `subnetworkField` is a subnetwork name, `subnetworkProjectField` will be used
// 	as the project if set. If not, we fallback on the default project.
func getSubnetworkLink(d *schema.ResourceData, config *Config, subnetworkField, subnetworkProjectField, zoneField string) (string, error) {
	if v, ok := d.GetOk(subnetworkField); ok {
		subnetwork := v.(string)
		r := regexp.MustCompile(SubnetworkLinkRegex)
		if r.MatchString(subnetwork) {
			return subnetwork, nil
		}

		var project string
		if subnetworkProject, ok := d.GetOk(subnetworkProjectField); ok {
			project = subnetworkProject.(string)
		} else {
			var err error
			project, err = getProject(d, config)
			if err != nil {
				return "", err
			}
		}

		region := getRegionFromZone(d.Get(zoneField).(string))

		subnet, err := config.clientCompute.Subnetworks.Get(project, region, subnetwork).Do()
		if err != nil {
			return "", fmt.Errorf(
				"Error referencing subnetwork '%s' in region '%s': %s",
				subnetwork, region, err)
		}

		return subnet.SelfLink, nil
	}
	return "", nil
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
	if strings.HasPrefix(network, "https://www.googleapis.com/compute/") {
		// extract the network name from SelfLink URL
		networkName := network[strings.LastIndex(network, "/")+1:]
		if networkName == "" {
			return "", fmt.Errorf("network url not valid")
		}
		return networkName, nil
	}

	return network, nil
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

// expandLabels pulls the value of "labels" out of a schema.ResourceData as a map[string]string.
func expandLabels(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "labels")
}

// expandStringMap pulls the value of key out of a schema.ResourceData as a map[string]string.
func expandStringMap(d *schema.ResourceData, key string) map[string]string {
	mp := map[string]string{}
	if v, ok := d.GetOk(key); ok {
		labelMap := v.(map[string]interface{})
		for k, v := range labelMap {
			mp[k] = v.(string)
		}
	}
	return mp
}

func convertStringArr(ifaceArr []interface{}) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, v.(string))
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
