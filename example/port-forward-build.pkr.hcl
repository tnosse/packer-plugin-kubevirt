
# packer {
#   required_plugins {
#     kubevirt = {
#       version = ">=v0.0.7"
#       source  = "github.com/tnosse/kubevirt"
#     }
#   }
# }

source "kubevirt" "port-forward-example" {
  // we use localhost, ssh_port will be generated
  ssh_host           = "localhost"
  ssh_username       = "ubuntu"
  output_image_file  = "${path.root}/basic-example.img"
  skip_extract_image = true
  source_image       = "quay.io/containerdisks/ubuntu:22.04"
  storage            = "3Gi"
  memory             = "1Gi"
}

build {
  sources = [
    "source.kubevirt.port-forward-example"
  ]

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
