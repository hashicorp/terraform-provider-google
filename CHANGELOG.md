## 0.1.1 (June 21, 2017)

BUG FIXES: 

* Restrict the number of health_checks in Backend Service resources to 1. ([#145](https://github.com/terraform-providers/terraform-provider-google/145))

## 0.1.0 (June 20, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `compute_disk.image`: shorhand for disk images is no longer supported, and will display a diff if used ([#1](https://github.com/terraform-providers/terraform-provider-google/1))

IMPROVEMENTS:

* Add an additional delay when checking for sql operations [[#15170](https://github.com/terraform-providers/terraform-provider-google/15170)](https://github.com/hashicorp/terraform/pull/15170)
* Add support for importing `compute_backend_service` ([#40](https://github.com/terraform-providers/terraform-provider-google/40))
* Wait for disk resizes to complete ([#1](https://github.com/terraform-providers/terraform-provider-google/1))
* Support `connection_draining_timeout_sec` in `google_compute_region_backend_service` ([#101](https://github.com/terraform-providers/terraform-provider-google/101))
* Add support for labels and tags on GKE node_config ([#7](https://github.com/terraform-providers/terraform-provider-google/7))
* Made `path_rule` optional in `google_compute_url_map`'s `path_matcher` block ([#118](https://github.com/terraform-providers/terraform-provider-google/118))

BUG FIXES:

* Changed `google_compute_instance_group_manager` `target_size` default to 0 ([#65](https://github.com/terraform-providers/terraform-provider-google/65))
* Represent GCS Bucket locations as uppercase in state. ([#117](https://github.com/terraform-providers/terraform-provider-google/117))
