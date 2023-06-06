// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No tests are added in this PR as currently there is no TF-supported method that can be used to
// enable both services (Cluster Istio and Mesh Istio) in GKE
func DataSourceMonitoringServiceMeshIstio() *schema.Resource {
	miSchema := map[string]*schema.Schema{
		"mesh_uid": {
			Type:     schema.TypeString,
			Required: true,
			Description: `Identifier for the mesh in which this Istio service is defined.
                        Corresponds to the meshUid metric label in Istio metrics.`,
		},
		"service_namespace": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The namespace of the Istio service underlying this service.
                        Corresponds to the destination_service_namespace metric label in Istio metrics.`,
		},
		"service_name": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The name of the Istio service underlying this service. 
                        Corresponds to the destination_service_name metric label in Istio metrics.`,
		},
	}
	t := `mesh_istio.mesh_uid="{{mesh_uid}}" AND 
            mesh_istio.service_name="{{service_name}}" AND 
            mesh_istio.service_namespace="{{service_namespace}}"`
	return dataSourceMonitoringServiceType(miSchema, t, dataSourceMonitoringServiceMeshIstioRead)
}

func dataSourceMonitoringServiceMeshIstioRead(res map[string]interface{}, d *schema.ResourceData, meta interface{}) error {
	var meshIstio map[string]interface{}
	if v, ok := res["mesh_istio"]; ok {
		meshIstio = v.(map[string]interface{})
	}
	if len(meshIstio) == 0 {
		return nil
	}
	if err := d.Set("service_name", meshIstio["service_name"]); err != nil {
		return err
	}
	if err := d.Set("service_namespace", meshIstio["service_namespace"]); err != nil {
		return err
	}
	if err := d.Set("mesh_name", meshIstio["mesh_name"]); err != nil {
		return err
	}
	return nil
}
