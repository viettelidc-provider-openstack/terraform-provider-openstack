output "floating_ip_addresses" {
  value = viettelidc_networking_floatingip_v2.fip.*.address
}
