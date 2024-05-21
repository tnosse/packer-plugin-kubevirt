# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

source "kubevirt" "basic-example" {
  communicator = "ssh"
  ssh_port = 2222
  ssh_host = "localhost"
  ssh_username = "ubuntu"
  output = "${path.root}/basic-example.img"
}

build {
  sources = [
    "source.kubevirt.basic-example"
  ]

  provisioner "file" {
    source = "${path.root}/test-assets/test-file.txt"
    destination = "/tmp/test-file.txt"
  }

  provisioner "file" {
    source = "/tmp/test-file.txt"
    destination = "${path.root}/test-assets/test-file-download.txt"
    direction = "download"
  }

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
