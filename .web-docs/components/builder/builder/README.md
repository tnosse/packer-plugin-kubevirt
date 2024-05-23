  Include a short description about the builder. This is a good place
  to call out what the builder does, and any requirements for the given
  builder environment. See https://www.packer.io/docs/builder/null
-->

The KubeVirt Packer Builder is a tool designed to automate the building of raw VM images from container disk images
in a KubeVirt environment. To use it, you'll need to have `kubectl` and `virtctl` CLI utilities installed on your system.

This builder relies on container disk images as the source for creating raw VM images, providing a repeatable process
for producing VM images for KubeVirt and KVM.

One of the key features of the KubeVirt Packer Builder is its integration within the Packer ecosystem.
This allows the use of Packer's various functionalities throughout the image building process, contributing to
a consistent and reliable outcome.

In the scope of continuous delivery pipelines that involve KubeVirt VM deployments, the KubeVirt Packer Builder's
ability to simplify and automate image creation can be a significant advantage.
But remember, before you start using this plugin, make sure you have kubectl and virtctl installed on your system.

<!-- Builder Configuration Fields -->

**Required**

- `ssh_username` (string) - The ssh username to use for the base image.
- `output_image_file` - The name of the image output file.

<!--
  Optional Configuration Fields

  Configuration options that are not required or have reasonable defaults
  should be listed under the optionals section. Defaults values should be
  noted in the description of the field
-->

**Optional**

- `skip_extract_image` (boolean) - If we should skip creating the `output_image_file`, used for testing.

<!--
  A basic example on the usage of the builder. Multiple examples
  can be provided to highlight various build configurations.

-->
### Example Usage


```hcl
source "kubevirt" "basic-example" {
  ssh_username = "ubuntu"
  output_image_file = "${path.root}/basic-example.img"
  skip_extract_image = true
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
```
