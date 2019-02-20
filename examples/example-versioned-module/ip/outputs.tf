output "address" {
  description = "The generated address of the ip address"

  value = "${coalesce(
		element(concat(google_compute_address.ip_address.*.address, list("")), 0),
		element(concat(google_compute_address.ip_address_beta.*.address, list("")), 0)
		)}"
}
