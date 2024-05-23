
packer {
  required_plugins {
    kubevirt = {
      version = ">=v0.0.3"
      source  = "github.com/tnosse/kubevirt"
    }
  }
}

source "kubevirt" "basic-example" {
  ssh_username       = "ubuntu"
  output_image_file  = "${path.root}/basic-example.img"
  skip_extract_image = true
  source_image       = "quay.io/containerdisks/ubuntu:22.04"
  storage            = "3Gi"
  memory             = "1Gi"
}

build {
  sources = [
    "source.kubevirt.basic-example"
  ]

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
