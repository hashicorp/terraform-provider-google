## 0.1.0 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `compute_disk.image`: shorhand for disk images is no longer supported, and will display a diff if used [GH-1]

IMPROVEMENTS:

* Add an additional delay when checking for sql operations [GH-15170](https://github.com/hashicorp/terraform/pull/15170)
* Add support for importing `compute_backend_service` [GH-40]
* Wait for disk resizes to complete [GH-1]
