variable "region" {
  default = "us-west1"
}

variable "region_zone" {
  default = "us-west1-a"
}

variable "project_name" {
  description = "The ID of the Google Cloud project"
}

variable "credentials_file_path" {
  description = "Path to the JSON file used to describe your account credentials"
  default     = "~/.gcloud/Terraform.json"
}
