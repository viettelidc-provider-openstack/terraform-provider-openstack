output "floating_ip" {
  value = viettelidc_networking_floatingip_v2.fip.address
}

output "volume_devices" {
  value = viettelidc_compute_volume_attach_v2.attached.device
}
