# https://cloud.google.com/vpc/docs/shared-vpc

provider "google" {
  region      = "${var.region}"
  credentials = "${file("${var.credentials_file_path}")}"
}

provider "random" {}

resource "random_id" "host_project_name" {
  byte_length = 8
}

resource "random_id" "service_project_1_name" {
  byte_length = 8
}

resource "random_id" "service_project_2_name" {
  byte_length = 8
}

resource "random_id" "standalone_project_name" {
  byte_length = 8
}

# The project which owns the VPC.
resource "google_project" "host_project" {
  name            = "Host Project"
  project_id      = "tf-vpc-${random_id.host_project_name.hex}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
}

# One project which will use the VPC.
resource "google_project" "service_project_1" {
  name            = "Service Project 1"
  project_id      = "tf-vpc-${random_id.service_project_1_name.hex}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
}

# The other project which will use the VPC.
resource "google_project" "service_project_2" {
  name            = "Service Project 2"
  project_id      = "tf-vpc-${random_id.service_project_2_name.hex}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
}

# A project which will not use the VPC, for the sake of demonstration.
resource "google_project" "standalone_project" {
  name            = "Standalone Project"
  project_id      = "tf-vpc-${random_id.standalone_project_name.hex}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
}

# Compute service needs to be enabled for all four new projects.
resource "google_project_service" "host_project" {
  project = "${google_project.host_project.project_id}"
  service = "compute.googleapis.com"
}

resource "google_project_service" "service_project_1" {
  project = "${google_project.service_project_1.project_id}"
  service = "compute.googleapis.com"
}

resource "google_project_service" "service_project_2" {
  project = "${google_project.service_project_2.project_id}"
  service = "compute.googleapis.com"
}

resource "google_project_service" "standalone_project" {
  project = "${google_project.standalone_project.project_id}"
  service = "compute.googleapis.com"
}

# Enable shared VPC hosting in the host project.
resource "google_compute_shared_vpc_host_project" "host_project" {
  project    = "${google_project.host_project.project_id}"
  depends_on = ["google_project_service.host_project"]
}

# Enable shared VPC in the two service projects - explicitly depend on the host
# project enabling it, because enabling shared VPC will fail if the host project
# is not yet hosting.
resource "google_compute_shared_vpc_service_project" "service_project_1" {
  host_project    = "${google_project.host_project.project_id}"
  service_project = "${google_project.service_project_1.project_id}"

  depends_on = ["google_compute_shared_vpc_host_project.host_project",
    "google_project_service.service_project_1",
  ]
}

resource "google_compute_shared_vpc_service_project" "service_project_2" {
  host_project    = "${google_project.host_project.project_id}"
  service_project = "${google_project.service_project_2.project_id}"

  depends_on = ["google_compute_shared_vpc_host_project.host_project",
    "google_project_service.service_project_2",
  ]
}

# Create the hosted network.
resource "google_compute_network" "shared_network" {
  name                    = "shared-network"
  auto_create_subnetworks = "true"
  project                 = "${google_compute_shared_vpc_host_project.host_project.project}"

  depends_on = ["google_compute_shared_vpc_service_project.service_project_1",
    "google_compute_shared_vpc_service_project.service_project_2",
  ]
}

# Allow the hosted network to be hit over ICMP, SSH, and HTTP.
resource "google_compute_firewall" "shared_network" {
  name    = "allow-ssh-and-icmp"
  network = "${google_compute_network.shared_network.self_link}"
  project = "${google_compute_network.shared_network.project}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "80"]
  }
}

# Create a standalone network with the same firewall rules.
resource "google_compute_network" "standalone_network" {
  name                    = "standalone-network"
  auto_create_subnetworks = "true"
  project                 = "${google_project.standalone_project.project_id}"
  depends_on              = ["google_project_service.standalone_project"]
}

resource "google_compute_firewall" "standalone_network" {
  name    = "allow-ssh-and-icmp"
  network = "${google_compute_network.standalone_network.self_link}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "80"]
  }

  project = "${google_project.standalone_project.project_id}"
}

# Create a VM which hosts a web page stating its identity ("VM1")
resource "google_compute_instance" "project_1_vm" {
  name         = "tf-project-1-vm"
  project      = "${google_project.service_project_1.project_id}"
  machine_type = "f1-micro"
  zone         = "${var.region_zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/family/debian-8"
    }
  }

  metadata_startup_script = "VM_NAME=VM1\n${file("scripts/install-vm.sh")}"

  network_interface {
    network = "${google_compute_network.shared_network.self_link}"

    access_config {
      // Ephemeral IP
    }
  }

  service_account {
    scopes = ["https://www.googleapis.com/auth/compute.readonly"]
  }

  depends_on = ["google_project_service.service_project_1"]
}

# Create a VM which hosts a web page demonstrating the example networking.
resource "google_compute_instance" "project_2_vm" {
  name         = "tf-project-2-vm"
  machine_type = "f1-micro"
  project      = "${google_project.service_project_2.project_id}"
  zone         = "${var.region_zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/family/debian-8"
    }
  }

  metadata_startup_script = <<EOF
VM1_EXT_IP=${google_compute_instance.project_1_vm.network_interface.0.access_config.0.assigned_nat_ip}
ST_VM_EXT_IP=${google_compute_instance.standalone_project_vm.network_interface.0.access_config.0.assigned_nat_ip}
VM1_INT_IP=${google_compute_instance.project_1_vm.network_interface.0.address}
ST_VM_INT_IP=${google_compute_instance.standalone_project_vm.network_interface.0.address}
${file("scripts/install-network-page.sh")}
EOF

  network_interface {
    network = "${google_compute_network.shared_network.self_link}"

    access_config {
      // Ephemeral IP
    }
  }

  service_account {
    scopes = ["https://www.googleapis.com/auth/compute.readonly"]
  }

  depends_on = ["google_project_service.service_project_2"]
}

# Create a VM which hosts a web page stating its identity ("standalone").
resource "google_compute_instance" "standalone_project_vm" {
  name         = "tf-standalone-vm"
  machine_type = "f1-micro"
  project      = "${google_project.standalone_project.project_id}"
  zone         = "${var.region_zone}"

  boot_disk {
    initialize_params {
      image = "projects/debian-cloud/global/images/family/debian-8"
    }
  }

  metadata_startup_script = "VM_NAME=standalone\n${file("scripts/install-vm.sh")}"

  network_interface {
    network = "${google_compute_network.standalone_network.self_link}"

    access_config {
      // Ephemeral IP
    }
  }

  service_account {
    scopes = ["https://www.googleapis.com/auth/compute.readonly"]
  }

  depends_on = ["google_project_service.standalone_project"]
}
