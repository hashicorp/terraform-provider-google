variable "region" {
  default = "us-central1"
}

variable "region_zone" {
  default = "us-central1-f"
}

variable "project_base_id" {
  description = "A string to use in the middle of the generated project names.  If you delete a project, its name cannot be used again, so if you run this example repeatedly, you might need to modify this."
}

variable "org_id" {
  description = "The ID of the Google Cloud Organization."
}

variable "billing_account_id" {
  description = "The ID of the associated billing account (optional)."
  default     = ""
}
