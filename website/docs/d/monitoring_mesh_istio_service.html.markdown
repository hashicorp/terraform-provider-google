---
subcategory: "Cloud (Stackdriver) Monitoring"
description: |-
  An Monitoring Service resource created automatically by GCP to monitor an
  Mesh Istio service.
---

# google_monitoring_mesh_istio_service

A Monitoring Service is the root resource under which operational aspects of a
generic service are accessible. A service is some discrete, autonomous, and
network-accessible unit, designed to solve an individual concern

An Mesh Istio monitoring service is automatically created by GCP to monitor
Mesh Istio services.


To get more information about Service, see:

* [API documentation](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/services)
* How-to Guides
    * [Service Monitoring](https://cloud.google.com/monitoring/service-monitoring)
    * [Monitoring API Documentation](https://cloud.google.com/monitoring/api/v3/)

## Example Usage - Monitoring Mesh Istio Service


```hcl
# Monitors the default MeshIstio service
data "google_monitoring_mesh_istio_service" "default" {
        mesh_uid = "proj-573164786102"
        service_namespace = "istio-system" 
        service_name = "prometheus"
}
```

## Argument Reference

The arguments of this data source act as filters for identifying a given -created service.

The given filters must match exactly one service whose data will be exported as attributes. The following arguments are supported:

The following fields must be specified:

* `mesh_uid` - (Required) Identifier for the mesh in which this Istio service is defined.
  Corresponds to the meshUid metric label in Istio metrics.

* `service_namespace` - (Required) The namespace of the Istio service underlying this service.
  Corresponds to the destination_service_namespace metric label in Istio metrics.

* `service_name` - (Required) The name of the Istio service underlying this service.
  Corresponds to the destination_service_name metric label in Istio metrics.
  
- - -

Other optional fields include:

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` -
  The full REST resource name for this channel. The syntax is:
  `projects/[PROJECT_ID]/services/[SERVICE_ID]`.

* `display_name` -
  Name used for UI elements listing this (Monitoring) Service.

* `telemetry` -
  Configuration for how to query telemetry on the Service. Structure is documented below.

The `telemetry` block includes:

* `resource_name` -
  (Optional)
  The full name of the resource that defines this service.
  Formatted as described in
  https://cloud.google.com/apis/design/resource_names.
