
# packer {
#   required_plugins {
#     kubevirt = {
#       version = ">=v0.0.8"
#       source  = "github.com/tnosse/kubevirt"
#     }
#   }
# }

source "kubevirt" "node-port-example" {
  // ssh_host must be set to a node in the cluster
  ssh_host     = "node-127-0-0-1.xip.io"
  ssh_username = "ubuntu"
  service_type = "NodePort"
  // service_port is optional with NodePort
  // service_port       = 31000
  output_image_file  = "${path.root}/basic-example.img"
  skip_extract_image = true
  source_image       = "quay.io/containerdisks/ubuntu:22.04"
  storage            = "3Gi"
  memory             = "1Gi"
}

build {
  sources = [
    "source.kubevirt.node-port-example"
  ]

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
