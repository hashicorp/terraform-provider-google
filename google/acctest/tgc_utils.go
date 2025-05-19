// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type ResourceMetadata struct {
	CaiAssetName    string         `json:"cai_asset_name"`
	ResourceType    string         `json:"resource_type"`
	ResourceAddress string         `json:"resource_address"`
	ImportMetadata  ImportMetadata `json:"import_metadata,omitempty"`
	Service         string         `json:"service"`
}

type ImportMetadata struct {
	Id            string   `json:"id,omitempty"`
	IgnoredFields []string `json:"ignored_fields,omitempty"`
}

type TgcMetadataPayload struct {
	TestName         string                      `json:"test_name"`
	RawConfig        string                      `json:"raw_config"`
	ResourceMetadata map[string]ResourceMetadata `json:"resource_metadata"`
	PrimaryResource  string                      `json:"primary_resource"`
}

// Hardcode the Terraform resource name -> API service name mapping temporarily.
// TODO: [tgc] read the mapping from the resource metadata files.
var ApiServiceNames = map[string]string{
	"google_compute_instance": "compute.googleapis.com",
	"google_project":          "cloudresourcemanager.googleapis.com",
}

// encodeToBase64JSON converts a struct to base64-encoded JSON
func encodeToBase64JSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// CollectAllTgcMetadata collects metadata for all resources in a test step
func CollectAllTgcMetadata(tgcPayload TgcMetadataPayload) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Process each resource to get CAI asset names and resolve auto IDs
		for address, metadata := range tgcPayload.ResourceMetadata {
			// If there is import metadata update our primary resource
			if metadata.ImportMetadata.Id != "" {
				tgcPayload.PrimaryResource = address
			}

			rState := s.RootModule().Resources[address]
			if rState == nil || rState.Primary == nil {
				log.Printf("[DEBUG]TGC Terraform error: resource state unavailable for %s, skipping", address)
				continue
			}

			// Resolve the CAI asset name
			if apiServiceName, ok := ApiServiceNames[metadata.ResourceType]; ok {
				var rName string
				switch metadata.ResourceType {
				case "google_project":
					rName = fmt.Sprintf("projects/%s", rState.Primary.Attributes["number"])
				default:
					rName = rState.Primary.ID
				}
				metadata.CaiAssetName = fmt.Sprintf("//%s/%s", apiServiceName, rName)
			} else {
				metadata.CaiAssetName = "unknown"
			}

			// Resolve auto IDs in import metadata
			if metadata.ImportMetadata.Id != "" {
				metadata.ImportMetadata.Id = strings.Replace(metadata.ImportMetadata.Id, "<AUTO_ID>", rState.Primary.ID, 1)
			}

			// Update the metadata in the map
			tgcPayload.ResourceMetadata[address] = metadata
		}

		// Encode the entire payload to base64 JSON
		encodedData, err := encodeToBase64JSON(tgcPayload)
		if err != nil {
			log.Printf("[DEBUG]TGC Terraform error: %v", err)
		} else {
			log.Printf("[DEBUG]TGC Terraform metadata: %s", encodedData)
		}

		return nil
	}
}

// parseResources extracts all resources from a Terraform configuration string
func parseResources(config string) []string {
	// This regex matches resource blocks in Terraform configurations
	resourceRegex := regexp.MustCompile(`resource\s+"([^"]+)"\s+"([^"]+)"`)
	matches := resourceRegex.FindAllStringSubmatch(config, -1)

	var resources []string
	for _, match := range matches {
		if len(match) >= 3 {
			// Combine resource type and name to form the address
			resources = append(resources, fmt.Sprintf("%s.%s", match[1], match[2]))
		}
	}

	return resources
}

// getServicePackage determines the service package for a resource type
func getServicePackage(resourceType string) string {
	var ServicePackages = map[string]string{
		"google_compute_":   "compute",
		"google_storage_":   "storage",
		"google_sql_":       "sql",
		"google_container_": "container",
		"google_bigquery_":  "bigquery",
		"google_project":    "resourcemanager",
		"google_cloud_run_": "cloudrun",
	}

	// Check for exact matches first
	if service, ok := ServicePackages[resourceType]; ok {
		return service
	}

	// Check for prefix matches
	for prefix, service := range ServicePackages {
		if strings.HasPrefix(resourceType, prefix) {
			return service
		}
	}

	// Default to "unknown" if no match found
	return "unknown"
}

// determineImportMetadata checks if the next step is an import step and extracts all import metadata
func determineImportMetadata(steps []resource.TestStep, currentStepIndex int, resourceName string) ImportMetadata {
	var metadata ImportMetadata

	// Check if there's a next step and if it's an import step
	if currentStepIndex+1 < len(steps) {
		nextStep := steps[currentStepIndex+1]

		// Check if it's an import step for our resource
		if nextStep.ImportState && (nextStep.ResourceName == resourceName ||
			strings.HasSuffix(nextStep.ResourceName, "."+strings.Split(resourceName, ".")[1])) {
			// Capture ignored fields if present
			if nextStep.ImportStateVerify && len(nextStep.ImportStateVerifyIgnore) > 0 {
				metadata.IgnoredFields = nextStep.ImportStateVerifyIgnore
			}

			// If ImportStateId is explicitly set, use that
			if nextStep.ImportStateId != "" {
				metadata.Id = nextStep.ImportStateId
				return metadata
			}

			// If ImportStateIdPrefix is set, note it
			if nextStep.ImportStateIdPrefix != "" {
				metadata.Id = fmt.Sprintf("%s<AUTO_ID>", nextStep.ImportStateIdPrefix)
				return metadata
			}

			// If ImportStateIdFunc is set, get function info
			if nextStep.ImportStateIdFunc != nil {
				metadata.Id = "<DYNAMIC_IMPORT_ID>"
				return metadata
			}

			// Default case - the ID will be automatically determined
			metadata.Id = "<AUTO_ID>"
			return metadata
		}
	}

	return metadata
}

// extendWithTGCData adds TGC metadata check function to the last non-plan config entry
func extendWithTGCData(t *testing.T, c resource.TestCase) resource.TestCase {
	var updatedSteps []resource.TestStep

	// Find the last non-plan config step
	lastNonPlanConfigStep := -1
	for i := len(c.Steps) - 1; i >= 0; i-- {
		step := c.Steps[i]
		if step.Config != "" && !step.PlanOnly {
			lastNonPlanConfigStep = i
			break
		}
	}

	// Process all steps
	for i, step := range c.Steps {
		// If this is the last non-plan config step, add our TGC check
		if i == lastNonPlanConfigStep {
			// Parse resources from the config
			resources := parseResources(step.Config)

			// Skip if no resources found
			if len(resources) == 0 {
				updatedSteps = append(updatedSteps, step)
				continue
			}

			// Determine the service package from the first resource
			firstResource := resources[0]
			parts := strings.Split(firstResource, ".")
			if len(parts) < 2 {
				updatedSteps = append(updatedSteps, step)
				continue
			}

			// Collect metadata for all resources
			resourceMetadata := make(map[string]ResourceMetadata)

			// Create the consolidated TGC payload
			tgcPayload := TgcMetadataPayload{
				TestName:         t.Name(),
				RawConfig:        step.Config,
				ResourceMetadata: resourceMetadata,
			}

			for _, res := range resources {
				parts := strings.Split(res, ".")
				if len(parts) >= 2 {
					resourceType := parts[0]

					// Determine import metadata if the next step is an import step
					importMeta := determineImportMetadata(c.Steps, i, res)

					// Create metadata for this resource
					resourceMetadata[res] = ResourceMetadata{
						ResourceType:    resourceType,
						ResourceAddress: res,
						ImportMetadata:  importMeta,
						Service:         getServicePackage(resourceType),
						// CaiAssetName will be populated at runtime in the check function
					}
				}
			}

			// Add a single consolidated TGC check for all resources
			tgcCheck := CollectAllTgcMetadata(tgcPayload)

			// If there's an existing check function, wrap it with our consolidated check
			if step.Check != nil {
				existingCheck := step.Check
				step.Check = resource.ComposeTestCheckFunc(
					existingCheck,
					tgcCheck,
				)
			} else {
				// Otherwise, just use our consolidated check
				step.Check = tgcCheck
			}
		}

		updatedSteps = append(updatedSteps, step)
	}

	c.Steps = updatedSteps
	return c
}
