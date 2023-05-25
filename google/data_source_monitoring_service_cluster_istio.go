// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No tests are added in this PR as currently there is no TF-supported method that can be used to
// enable both services (Cluster Istio and Mesh Istio) in GKE
func DataSourceMonitoringServiceClusterIstio() *schema.Resource {
	ciSchema := map[string]*schema.Schema{
		"location": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The location of the Kubernetes cluster in which this Istio service is defined. 
                        Corresponds to the location resource label in k8s_cluster resources.`,
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The name of the Kubernetes cluster in which this Istio service is defined. 
                        Corresponds to the clusterName resource label in k8s_cluster resources.`,
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
	filter := `cluster_istio.cluster_name="{{cluster_name}}" AND
		cluster_istio.service_namespace="{{service_namespace}}" AND
		cluster_istio.service_name="{{service_name}}" AND
		cluster_istio.location="{{location}}"`
	return dataSourceMonitoringServiceType(ciSchema, filter, dataSourceMonitoringServiceClusterIstioRead)
}

func dataSourceMonitoringServiceClusterIstioRead(res map[string]interface{}, d *schema.ResourceData, meta interface{}) error {
	var clusterIstio map[string]interface{}
	if v, ok := res["cluster_istio"]; ok {
		clusterIstio = v.(map[string]interface{})
	}
	if len(clusterIstio) == 0 {
		return nil
	}

	if err := d.Set("location", clusterIstio["location"]); err != nil {
		return err
	}
	if err := d.Set("service_name", clusterIstio["service_name"]); err != nil {
		return err
	}
	if err := d.Set("service_namespace", clusterIstio["service_namespace"]); err != nil {
		return err
	}
	if err := d.Set("cluster_name", clusterIstio["cluster_name"]); err != nil {
		return err
	}
	return nil
}
