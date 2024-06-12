
# packer {
#   required_plugins {
#     kubevirt = {
#       version = ">=v0.0.7"
#       source  = "github.com/tnosse/kubevirt"
#     }
#   }
# }

source "kubevirt" "lb-example" {
  service_type = "LoadBalancer"
  // optional, default is 22
  service_port       = 22
  ssh_username       = "ubuntu"
  output_image_file  = "${path.root}/basic-example.img"
  skip_extract_image = true
  source_image       = "quay.io/containerdisks/ubuntu:22.04"
  storage            = "3Gi"
  memory             = "1Gi"
}

build {
  sources = [
    "source.kubevirt.lb-example"
  ]

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
