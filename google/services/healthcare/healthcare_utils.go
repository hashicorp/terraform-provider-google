// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package healthcare

import (
	"fmt"
	"regexp"
	"strings"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

type HealthcareDatasetId struct {
	Project  string
	Location string
	Name     string
}

func (s *HealthcareDatasetId) DatasetId() string {
	return fmt.Sprintf("projects/%s/locations/%s/datasets/%s", s.Project, s.Location, s.Name)
}

func (s *HealthcareDatasetId) TerraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Location, s.Name)
}

func ParseHealthcareDatasetId(id string, config *transport_tpg.Config) (*HealthcareDatasetId, error) {
	parts := strings.Split(id, "/")

	datasetIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})$")
	datasetIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})$")
	datasetRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})$")

	if datasetIdRegex.MatchString(id) {
		return &HealthcareDatasetId{
			Project:  parts[0],
			Location: parts[1],
			Name:     parts[2],
		}, nil
	}

	if datasetIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{datasetName}` id format.")
		}

		return &HealthcareDatasetId{
			Project:  config.Project,
			Location: parts[0],
			Name:     parts[1],
		}, nil
	}

	if parts := datasetRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &HealthcareDatasetId{
			Project:  parts[1],
			Location: parts[2],
			Name:     parts[3],
		}, nil
	}
	return nil, fmt.Errorf("Invalid Dataset id format, expecting `{projectId}/{locationId}/{datasetName}` or `{locationId}/{datasetName}.`")
}

type healthcareFhirStoreId struct {
	DatasetId HealthcareDatasetId
	Name      string
}

func (s *healthcareFhirStoreId) FhirStoreId() string {
	return fmt.Sprintf("%s/fhirStores/%s", s.DatasetId.DatasetId(), s.Name)
}

func (s *healthcareFhirStoreId) TerraformId() string {
	return fmt.Sprintf("%s/%s", s.DatasetId.TerraformId(), s.Name)
}

func ParseHealthcareFhirStoreId(id string, config *transport_tpg.Config) (*healthcareFhirStoreId, error) {
	parts := strings.Split(id, "/")
	fhirStoreIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	fhirStoreIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	fhirStoreRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})/fhirStores/([a-zA-Z0-9_-]{1,256})$")

	if fhirStoreIdRegex.MatchString(id) {
		return &healthcareFhirStoreId{
			DatasetId: HealthcareDatasetId{
				Project:  parts[0],
				Location: parts[1],
				Name:     parts[2],
			},
			Name: parts[3],
		}, nil
	}

	if fhirStoreIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{datasetName}/{fhirStoreName}` id format.")
		}

		return &healthcareFhirStoreId{
			DatasetId: HealthcareDatasetId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := fhirStoreRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareFhirStoreId{
			DatasetId: HealthcareDatasetId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}
	return nil, fmt.Errorf("Invalid FhirStore id format, expecting `{projectId}/{locationId}/{datasetName}/{fhirStoreName}` or `{locationId}/{datasetName}/{fhirStoreName}.`")
}

type healthcareHl7V2StoreId struct {
	DatasetId HealthcareDatasetId
	Name      string
}

func (s *healthcareHl7V2StoreId) Hl7V2StoreId() string {
	return fmt.Sprintf("%s/hl7V2Stores/%s", s.DatasetId.DatasetId(), s.Name)
}

func (s *healthcareHl7V2StoreId) TerraformId() string {
	return fmt.Sprintf("%s/%s", s.DatasetId.TerraformId(), s.Name)
}

func ParseHealthcareHl7V2StoreId(id string, config *transport_tpg.Config) (*healthcareHl7V2StoreId, error) {
	parts := strings.Split(id, "/")
	hl7V2StoreIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	hl7V2StoreIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	hl7V2StoreRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})/hl7V2Stores/([a-zA-Z0-9_-]{1,256})$")

	if hl7V2StoreIdRegex.MatchString(id) {
		return &healthcareHl7V2StoreId{
			DatasetId: HealthcareDatasetId{
				Project:  parts[0],
				Location: parts[1],
				Name:     parts[2],
			},
			Name: parts[3],
		}, nil
	}

	if hl7V2StoreIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{datasetName}/{hl7V2StoreName}` id format.")
		}

		return &healthcareHl7V2StoreId{
			DatasetId: HealthcareDatasetId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := hl7V2StoreRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareHl7V2StoreId{
			DatasetId: HealthcareDatasetId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}
	return nil, fmt.Errorf("Invalid Hl7V2Store id format, expecting `{projectId}/{locationId}/{datasetName}/{hl7V2StoreName}` or `{locationId}/{datasetName}/{hl7V2StoreName}.`")
}

type healthcareDicomStoreId struct {
	DatasetId HealthcareDatasetId
	Name      string
}

func (s *healthcareDicomStoreId) DicomStoreId() string {
	return fmt.Sprintf("%s/dicomStores/%s", s.DatasetId.DatasetId(), s.Name)
}

func (s *healthcareDicomStoreId) TerraformId() string {
	return fmt.Sprintf("%s/%s", s.DatasetId.TerraformId(), s.Name)
}

func ParseHealthcareDicomStoreId(id string, config *transport_tpg.Config) (*healthcareDicomStoreId, error) {
	parts := strings.Split(id, "/")
	dicomStoreIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	dicomStoreIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	dicomStoreRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})/dicomStores/([a-zA-Z0-9_-]{1,256})$")

	if dicomStoreIdRegex.MatchString(id) {
		return &healthcareDicomStoreId{
			DatasetId: HealthcareDatasetId{
				Project:  parts[0],
				Location: parts[1],
				Name:     parts[2],
			},
			Name: parts[3],
		}, nil
	}

	if dicomStoreIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{datasetName}/{dicomStoreName}` id format.")
		}

		return &healthcareDicomStoreId{
			DatasetId: HealthcareDatasetId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := dicomStoreRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareDicomStoreId{
			DatasetId: HealthcareDatasetId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}
	return nil, fmt.Errorf("Invalid DicomStore id format, expecting `{projectId}/{locationId}/{datasetName}/{dicomStoreName}` or `{locationId}/{datasetName}/{dicomStoreName}.`")
}
