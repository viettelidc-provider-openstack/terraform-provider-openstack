resource "viettelidc_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = file("${var.ssh_key_file}.pub")
}

resource "viettelidc_compute_instance_v2" "my_instance" {
  name            = "my_instance"
  image_name      = var.image
  flavor_name     = var.flavor
  key_pair        = viettelidc_compute_keypair_v2.terraform.name
  security_groups = ["default"]
  network {
    name = var.network_name
  }
}

resource "viettelidc_networking_floatingip_v2" "fip" {
  pool = var.pool
}

resource "viettelidc_compute_floatingip_associate_v2" "fip" {
  instance_id = viettelidc_compute_instance_v2.my_instance.id
  floating_ip = viettelidc_networking_floatingip_v2.fip.address
  connection {
    host        = viettelidc_networking_floatingip_v2.fip.address
    user        = var.ssh_user_name
    private_key = file(var.ssh_key_file)
  }

  provisioner "local-exec" {
    command = "echo ${viettelidc_networking_floatingip_v2.fip.address} > instance_ip.txt"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mkfs.ext4 ${viettelidc_compute_volume_attach_v2.attached.device}",
      "sudo mkdir /mnt/volume",
      "sudo mount ${viettelidc_compute_volume_attach_v2.attached.device} /mnt/volume",
      "sudo df -h /mnt/volume",
    ]
  }
}
