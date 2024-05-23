
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

  provisioner "file" {
    source      = "${path.root}/test-assets/test-file.txt"
    destination = "/tmp/test-file.txt"
  }

  provisioner "file" {
    source      = "/tmp/test-file.txt"
    destination = "${path.root}/test-assets/test-file-download.txt"
    direction   = "download"
  }

  provisioner "shell" {
    inline = [
      "uname -a",
    ]
  }
}
