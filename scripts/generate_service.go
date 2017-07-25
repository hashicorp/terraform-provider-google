package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

type ComputeApiGenVersion int

const (
	none ComputeApiGenVersion = iota
	v1
	v0beta
)

var OrderedComputeApiVersions = []ComputeApiGenVersion{
	v0beta,
	v1,
}

type VersionInfo struct {
	ClientNameSuffix string
	ServiceName      string
}

type NormalisedData struct {
	Type          string
	LowestVersion ComputeApiGenVersion
	UpdateVersion ComputeApiGenVersion
	ExtraParams   []string
}

type TemplateData struct {
	// Name for the resource in the new service API
	ResourceName string

	// Type of the resource object in the client. Generally but not always the ResourceName.
	ClientType string

	// Name of the Service of the client. Generally but not always the ResourceName pluralised.
	ServiceType string

	// Versions supported by this resource.
	Versions []VersionInfo

	// Versions supported by this resource for update.
	UpdateVersions []VersionInfo

	// Parameters to add to signatures/calls denormalised.
	ExtraParamsSignature string
	ExtraParamsCall      string
}

func main() {
	clientTypeFlag := flag.String("type", "", "resource to generate a service for.")
	scopeFlag := flag.String("scope", "", "scope if non-global. one of `region` or `zone`.")
	parentFlag := flag.String("parent", "", "name of parent resource, if any.")
	lowestVersionFlag := flag.String("lowestversion", "v1", "lowest version of the resource. one of `v1` or `v0beta`. defaults to v1.")
	updateVersionFlag := flag.String("updateversion", "", "version at which update is supported. one of `v1` or `v0beta`.")
	flag.Parse()

	if *clientTypeFlag == "" {
		flag.PrintDefaults()
		log.Fatal("usage: go run generate_service.go -type $TYPE -scope $SCOPE -parent $PARENT -lowestversion $VERSION -updateversion $VERSION")
	}

	if *scopeFlag != "" && !(*scopeFlag == "region" || *scopeFlag == "zone") {
		log.Fatal("scope must be `region` or `zone`.")
	}

	lowestVersion := v1
	if !(*lowestVersionFlag == "v1" || *lowestVersionFlag == "v0beta") {
		log.Fatal("lowestversion must be `v1` or `v0beta`.")
	} else if *lowestVersionFlag == "v0beta" {
		lowestVersion = v0beta
	}

	updateVersion := none
	if *updateVersionFlag != "" && !(*updateVersionFlag == "v1" || *updateVersionFlag == "v0beta") {
		log.Fatal("updateversion must be `v1` or `v0beta`.")
	} else if *updateVersionFlag == "v0beta" {
		updateVersion = v0beta
	} else if *updateVersionFlag == "v1" {
		updateVersion = v1
	}

	// This is the information we would expect users to provide
	// TODO: investigate whether we need more params than scope+parent
	normalisedData := &NormalisedData{
		Type:          *clientTypeFlag,
		LowestVersion: lowestVersion,
		ExtraParams:   []string{},
		UpdateVersion: updateVersion,
	}

	if *scopeFlag != "" {
		normalisedData.ExtraParams = append(normalisedData.ExtraParams, *scopeFlag)
	}

	if *parentFlag != "" {
		normalisedData.ExtraParams = append(normalisedData.ExtraParams, *parentFlag)
	}

	// Override in cases where multiple API resources are backed by the same client struct
	overrides := map[string]string{
		"GlobalAddress": "Address",
	}

	clientType := normalisedData.Type
	if v, ok := overrides[clientType]; ok {
		clientType = v
	}

	templateData := TemplateData{
		ResourceName: normalisedData.Type,
		ClientType:   clientType,
	}

	// eg. Address -> Addresses instead of Address -> Addresss
	if strings.HasSuffix(templateData.ClientType, "s") {
		templateData.ServiceType = normalisedData.Type + "es"
	} else {
		templateData.ServiceType = normalisedData.Type + "s"
	}

	// Add the strings required for each version that the resource supports.
	for _, version := range OrderedComputeApiVersions {
		clientNameSuffix := ""
		serviceName := "v1"

		if version == v0beta {
			clientNameSuffix = "Beta"
			serviceName = "v0beta"
		}

		templateData.Versions = append(templateData.Versions, VersionInfo{ClientNameSuffix: clientNameSuffix, ServiceName: serviceName})
		if version == normalisedData.LowestVersion {
			break
		}
	}

	// If a resource can be updated, add information about how to update it
	if normalisedData.UpdateVersion != none {
		for _, version := range OrderedComputeApiVersions {
			clientNameSuffix := ""
			serviceName := "v1"

			if version == v0beta {
				clientNameSuffix = "Beta"
				serviceName = "v0beta"
			}

			templateData.UpdateVersions = append(templateData.UpdateVersions, VersionInfo{ClientNameSuffix: clientNameSuffix, ServiceName: serviceName})
			if version == normalisedData.UpdateVersion {
				break
			}
		}
	}

	// If the resource has extra parameters like a scope or parent resource, add them now.
	// This works under the assumption that they are always strings.
	if len(normalisedData.ExtraParams) > 0 {
		eCall := ""

		for _, param := range normalisedData.ExtraParams {
			eCall += fmt.Sprintf(" %s,", param)
		}

		templateData.ExtraParamsCall = eCall
		templateData.ExtraParamsSignature = eCall[:len(eCall)-1] + " string,"
	}

	f, err := os.Create(fmt.Sprintf("gen-%s.go", *clientTypeFlag))
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	t := template.Must(template.New("temp").Parse(multiversionServiceTemplate))
	t.Execute(buf, templateData)

	fmtd, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("Formatting error: %s", err)
		fmtd = buf.Bytes()
	}

	if _, err := f.Write(fmtd); err != nil {
		log.Fatal(err)
	}
}

// TODO: copy ForceSendFields and NullFields into the new structs so they only need to be defined once.
// TODO: copy Imports
// TODO: create an interface for a $TYPEService?
const multiversionServiceTemplate = `
package google

func (s *ComputeMultiversionService) Insert{{ $.ResourceName }}(project string,{{ $.ExtraParamsSignature }} resource *computeBeta.{{ $.ClientType }}, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	{{ range .Versions }}
		case {{ .ServiceName  }}:
		{{ .ServiceName }}Resource := &compute{{ .ClientNameSuffix }}.{{ $.ClientType }}{}
		err := Convert(resource, {{ .ServiceName }}Resource)
		if err != nil {
			return nil, err
		}

		{{ .ServiceName }}Op, err := s.{{ .ServiceName }}.{{ $.ServiceType }}.Insert(project,{{ $.ExtraParamsCall }} {{ .ServiceName }}Resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert({{ .ServiceName }}Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) Get{{ $.ResourceName }}(project string,{{ $.ExtraParamsSignature }} resource string, version ComputeApiVersion) (*computeBeta.{{ $.ClientType }}, error) {
	res := &computeBeta.{{ .ClientType }}{}
	switch version {
	{{ range .Versions }}
		case {{ .ServiceName }}:
		r, err := s.{{ .ServiceName }}.{{ $.ServiceType }}.Get(project,{{ $.ExtraParamsCall }} resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(r, res)
		if err != nil {
			return nil, err
		}

		return res, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}
{{ if ne (len .UpdateVersions) 0 }}
func (s *ComputeMultiversionService) Update{{ $.ResourceName }}(project string,{{ $.ExtraParamsSignature }} resourceName string, resource *computeBeta.{{ $.ClientType }}, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	{{ range .UpdateVersions }}
		case {{ .ServiceName }}:
		{{ .ServiceName }}Resource := &compute{{ .ClientNameSuffix }}.{{ $.ClientType }}{}
		err := Convert(resource, {{ .ServiceName }}Resource)
		if err != nil {
			return nil, err
		}
		{{ .ServiceName }}Op, err := s.{{ .ServiceName }}.{{ $.ServiceType }}.Update(project,{{ $.ExtraParamsCall }} resourceName, {{ .ServiceName }}Resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert({{ .ServiceName }}Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}
{{ end }}
func (s *ComputeMultiversionService) Delete{{ $.ResourceName }}(project string,{{ $.ExtraParamsSignature }} resource string, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	{{ range .Versions }}
		case {{ .ServiceName }}:
		{{ .ServiceName }}Op, err := s.{{ .ServiceName }}.{{ $.ServiceType }}.Delete(project,{{ $.ExtraParamsCall }} resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert({{ .ServiceName }}Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	{{ end }}
	}

	return nil, fmt.Errorf("Unknown API version.")
}

`
