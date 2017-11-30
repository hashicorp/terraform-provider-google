# https://cloud.google.com/vpc/docs/shared-vpc

provider "google" {
  region  = "${var.region}"
  project = "host-project-${var.project_base_id}"
}

resource "google_project" "host_project" {
  name            = "Host Project"
  project_id      = "host-project-${var.project_base_id}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
  skip_delete     = "true"
}

resource "google_project" "service_project_1" {
  name            = "Service Project 1"
  project_id      = "service-project-${var.project_base_id}-1"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
  skip_delete     = "true"
}

resource "google_project" "service_project_2" {
  name            = "Service Project 2"
  project_id      = "service-project-${var.project_base_id}-2"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
  skip_delete     = "true"
}

resource "google_project" "standalone_project" {
  name            = "Standalone Project"
  project_id      = "standalone-${var.project_base_id}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account_id}"
  skip_delete     = "true"
}

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

resource "google_compute_shared_vpc_host_project" "host_project" {
  project    = "${google_project.host_project.project_id}"
  depends_on = ["google_project_service.host_project"]
}

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

resource "google_compute_network" "shared_network" {
  name                    = "shared-network"
  auto_create_subnetworks = "true"
  project                 = "${google_compute_shared_vpc_host_project.host_project.project}"

  depends_on = ["google_compute_shared_vpc_service_project.service_project_1",
    "google_compute_shared_vpc_service_project.service_project_2",
  ]
}

resource "google_compute_firewall" "shared_network" {
  name    = "allow-ssh-and-icmp"
  network = "${google_compute_network.shared_network.self_link}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "80"]
  }
}

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
