// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No tests are added in this PR as currently there is no TF-supported method that can be used to
// enable either services (Cluster Istio, Mesh Istio and Istio Canonical Service) in GKE
func DataSourceMonitoringIstioCanonicalService() *schema.Resource {
	csSchema := map[string]*schema.Schema{
		"mesh_uid": {
			Type:     schema.TypeString,
			Required: true,
			Description: `Identifier for the Istio mesh in which this canonical service is defined.
                        Corresponds to the meshUid metric label in Istio metrics.`,
		},
		"canonical_service_namespace": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The namespace of the canonical service underlying this service.
                        Corresponds to the destination_service_namespace metric label in Istio metrics.`,
		},
		"canonical_service": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The name of the canonical service underlying this service.. 
                        Corresponds to the destination_service_name metric label in Istio metrics.`,
		},
	}
	t := `istio_canonical_service.mesh_uid="{{mesh_uid}}" AND 
			istio_canonical_service.canonical_service="{{canonical_service}}" AND 
			istio_canonical_service.canonical_service_namespace="{{canonical_service_namespace}}"`
	return dataSourceMonitoringServiceType(csSchema, t, dataSourceMonitoringIstioCanonicalServiceRead)
}

func dataSourceMonitoringIstioCanonicalServiceRead(res map[string]interface{}, d *schema.ResourceData, meta interface{}) error {
	var istioCanonicalService map[string]interface{}
	if v, ok := res["istio_canonical_service"]; ok {
		istioCanonicalService = v.(map[string]interface{})
	}
	if len(istioCanonicalService) == 0 {
		return nil
	}
	if err := d.Set("canonical_service", istioCanonicalService["canonical_service"]); err != nil {
		return err
	}
	if err := d.Set("canonical_service_namespace", istioCanonicalService["canonical_service_namespace"]); err != nil {
		return err
	}
	if err := d.Set("mesh_name", istioCanonicalService["mesh_name"]); err != nil {
		return err
	}
	return nil
}
