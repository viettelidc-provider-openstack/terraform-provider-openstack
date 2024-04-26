resource "viettelidc_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = file("${var.ssh_key_file}.pub")
}

resource "viettelidc_compute_instance_v2" "multi" {
  count           = var.instance_count
  name            = format("${var.instance_prefix}-%02d", count.index + 1)
  image_name      = var.image
  flavor_name     = var.flavor
  key_pair        = viettelidc_compute_keypair_v2.terraform.name
  security_groups = ["default"]
  network {
    name = var.network_name
  }
}

resource "viettelidc_networking_floatingip_v2" "fip" {
  count = var.instance_count
  pool  = var.pool
}

resource "viettelidc_compute_floatingip_associate_v2" "fip" {
  count       = var.instance_count
  instance_id = viettelidc_compute_instance_v2.multi.*.id[count.index]
  floating_ip = viettelidc_networking_floatingip_v2.fip.*.address[count.index]
}
