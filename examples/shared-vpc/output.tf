output "status_page_public_ip" {
  value = "${google_compute_instance.project_2_vm.network_interface.0.access_config.0.assigned_nat_ip}"
}
