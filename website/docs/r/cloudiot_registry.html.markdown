---
subcategory: "Cloud IoT Core"
layout: "google"
page_title: "Google: google_cloudiot_registry"
sidebar_current: "docs-google-cloudiot-registry-x"
description: |-
  Creates a device registry in Google's Cloud IoT Core platform
---

# google\_cloudiot\_registry

 Creates a device registry in Google's Cloud IoT Core platform. For more information see
[the official documentation](https://cloud.google.com/iot/docs/) and
[API](https://cloud.google.com/iot/docs/reference/cloudiot/rest/v1/projects.locations.registries).


## Example Usage

```hcl
resource "google_pubsub_topic" "default-devicestatus" {
  name = "default-devicestatus"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "default-telemetry"
}

resource "google_cloudiot_registry" "default-registry" {
  name = "default-registry"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.default-telemetry.id
  }

  state_notification_config = {
    pubsub_topic_name = google_pubsub_topic.default-devicestatus.id
  }

  http_config = {
    http_enabled_state = "HTTP_ENABLED"
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_ENABLED"
  }

  credentials {
    public_key_certificate = {
      format      = "X509_CERTIFICATE_PEM"
      certificate = file("rsa_cert.pem")
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by device registry.
    Changing this forces a new resource to be created.

- - -

* `project` - (Optional) The project in which the resource belongs. If it is not provided, the provider project is used.

* `region` - (Optional) The Region in which the created address should reside. If it is not provided, the provider region is used.

* `event_notification_configs` - (Optional) List of configurations for event notification, such as
PubSub topics to publish device events to. Structure is documented below.

* `state_notification_config` - (Optional) A PubSub topic to publish device state updates. Structure is documented below.

* `mqtt_config` - (Optional) Activate or deactivate MQTT. Structure is documented below.
* `http_config` - (Optional) Activate or deactivate HTTP. Structure is documented below.

* `credentials` - (Optional) List of public key certificates to authenticate devices. Structure is documented below. 


The `event_notification_configs` block supports:

* `pubsub_topic_name` - (Required) PubSub topic name to publish device events.

* `subfolder_matches` - (Optional) If the subfolder name matches this string
   exactly, this configuration will be used. The string must not include the
   leading '/' character. If empty, all strings are matched. Empty value can
   only be used for the last `event_notification_configs` item.

The `state_notification_config` block supports:

* `pubsub_topic_name` - (Required) PubSub topic name to publish device state updates.

The `mqtt_config` block supports:

* `mqtt_enabled_state` - (Required) The field allows `MQTT_ENABLED` or `MQTT_DISABLED`.

The `http_config` block supports:

* `http_enabled_state` - (Required) The field allows `HTTP_ENABLED` or `HTTP_DISABLED`.

The `credentials` block supports:

* `public_key_certificate` - (Required) The certificate format and data.

The `public_key_certificate` block supports:

* `format` - (Required) The field allows only  `X509_CERTIFICATE_PEM`.
* `certificate` - (Required) The certificate data.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{region}}/registries/{{name}}`

## Import

A device registry can be imported using the `name`, e.g.

```
$ terraform import google_cloudiot_registry.default-registry projects/{project}/locations/{region}/registries/{name}
```
