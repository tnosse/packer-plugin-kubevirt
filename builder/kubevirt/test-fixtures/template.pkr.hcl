# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

source "kubevirt" "basic-example" {
  mock = "mock-config"
  communicator = "ssh"
  ssh_port = 2222
  ssh_host = "localhost"
  ssh_username = "ubuntu"
}

build {
  sources = [
    "source.kubevirt.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo build generated data: ${build.GeneratedMockData}",
    ]
  }

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
