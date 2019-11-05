# This will use the `google` provider
module "ip" {
  source = "./ip"
  name   = "ipv4"
}

# The following modules will use the `google-beta` provider
# Because it has been aliased to the `google` name
module "ip-beta" {
  source = "./ip"

  name = "ipv4-beta"
  labels = {
    "hello" = "world"
    "foo"   = "bar"
  }

  providers = {
    google = google-beta
  }
}

module "ip-beta-no-labels" {
  source = "./ip"

  name = "ipv4-beta-no-labels"

  providers = {
    google = google-beta
  }
}

# Using the `google-beta` provider in a config requires
# the `google-beta` provider block
provider "google-beta" {
}

# Display outputs from each block
output "ip_address" {
  value = module.ip.address
}

output "ip_address_beta" {
  value = module.ip-beta.address
}

output "ip_address_beta_no_labels" {
  value = module.ip-beta-no-labels.address
}
