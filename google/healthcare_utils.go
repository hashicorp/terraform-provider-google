package google

import (
	"fmt"
	"regexp"
	"strings"
)

type healthcareDatasetId struct {
	Project  string
	Location string
	Name     string
}

func (s *healthcareDatasetId) datasetId() string {
	return fmt.Sprintf("projects/%s/locations/%s/datasets/%s", s.Project, s.Location, s.Name)
}

func (s *healthcareDatasetId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Location, s.Name)
}

func parseHealthcareDatasetId(id string, config *Config) (*healthcareDatasetId, error) {
	parts := strings.Split(id, "/")

	datasetIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})$")
	datasetIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})$")
	datasetRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})$")

	if datasetIdRegex.MatchString(id) {
		return &healthcareDatasetId{
			Project:  parts[0],
			Location: parts[1],
			Name:     parts[2],
		}, nil
	}

	if datasetIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{datasetName}` id format.")
		}

		return &healthcareDatasetId{
			Project:  config.Project,
			Location: parts[0],
			Name:     parts[1],
		}, nil
	}

	if parts := datasetRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareDatasetId{
			Project:  parts[1],
			Location: parts[2],
			Name:     parts[3],
		}, nil
	}
	return nil, fmt.Errorf("Invalid Dataset id format, expecting `{projectId}/{locationId}/{datasetName}` or `{locationId}/{datasetName}.`")
}

type healthcareFhirStoreId struct {
	DatasetId healthcareDatasetId
	Name      string
}

func (s *healthcareFhirStoreId) fhirStoreId() string {
	return fmt.Sprintf("%s/fhirStores/%s", s.DatasetId.datasetId(), s.Name)
}

func (s *healthcareFhirStoreId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.DatasetId.terraformId(), s.Name)
}

func parseHealthcareFhirStoreId(id string, config *Config) (*healthcareFhirStoreId, error) {
	parts := strings.Split(id, "/")

	fhirStoreIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	fhirStoreIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	fhirStoreRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})/fhirStores/([a-zA-Z0-9_-]{1,256})$")

	if fhirStoreIdRegex.MatchString(id) {
		return &healthcareFhirStoreId{
			DatasetId: healthcareDatasetId{
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
			DatasetId: healthcareDatasetId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := fhirStoreRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareFhirStoreId{
			DatasetId: healthcareDatasetId{
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
	DatasetId healthcareDatasetId
	Name      string
}

func (s *healthcareHl7V2StoreId) hl7V2StoreId() string {
	return fmt.Sprintf("%s/hl7V2Stores/%s", s.DatasetId.datasetId(), s.Name)
}

func (s *healthcareHl7V2StoreId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.DatasetId.terraformId(), s.Name)
}

func parseHealthcareHl7V2StoreId(id string, config *Config) (*healthcareHl7V2StoreId, error) {
	parts := strings.Split(id, "/")

	hl7V2StoreIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	hl7V2StoreIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	hl7V2StoreRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})/hl7V2Stores/([a-zA-Z0-9_-]{1,256})$")

	if hl7V2StoreIdRegex.MatchString(id) {
		return &healthcareHl7V2StoreId{
			DatasetId: healthcareDatasetId{
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
			DatasetId: healthcareDatasetId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := hl7V2StoreRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareHl7V2StoreId{
			DatasetId: healthcareDatasetId{
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
	DatasetId healthcareDatasetId
	Name      string
}

func (s *healthcareDicomStoreId) dicomStoreId() string {
	return fmt.Sprintf("%s/dicomStores/%s", s.DatasetId.datasetId(), s.Name)
}

func (s *healthcareDicomStoreId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.DatasetId.terraformId(), s.Name)
}

func parseHealthcareDicomStoreId(id string, config *Config) (*healthcareDicomStoreId, error) {
	parts := strings.Split(id, "/")

	dicomStoreIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	dicomStoreIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,256})/([a-zA-Z0-9_-]{1,256})$")
	dicomStoreRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/datasets/([a-zA-Z0-9_-]{1,256})/dicomStores/([a-zA-Z0-9_-]{1,256})$")

	if dicomStoreIdRegex.MatchString(id) {
		return &healthcareDicomStoreId{
			DatasetId: healthcareDatasetId{
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
			DatasetId: healthcareDatasetId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := dicomStoreRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &healthcareDicomStoreId{
			DatasetId: healthcareDatasetId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}
	return nil, fmt.Errorf("Invalid DicomStore id format, expecting `{projectId}/{locationId}/{datasetName}/{dicomStoreName}` or `{locationId}/{datasetName}/{dicomStoreName}.`")
}
