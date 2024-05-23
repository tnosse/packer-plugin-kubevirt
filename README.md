# Packer Plugin for KubeVirt

This repository houses the Packer KubeVirt plugin. This plugin allows Packer to create KubeVirt images.

Packer is a tool from HashiCorp that allows you to create identical machine images for multiple platforms from a single source configuration file.

KubeVirt is an add-on to Kubernetes, which enables running virtual machines on top of Kubernetes.

## Quick Start
1. Download and install [Packer](https://www.packer.io/downloads.html)
2. [Install](#Installation) the plugin
3. Refer to the [Documentation](#Documentation) and [Examples](#Examples)

## Requirements
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and [`virtctl`](https://kubevirt.io/user-guide/operations/virtctl_client_tool/) CLIs
- Access to a Kubernetes cluster with KubeVirt installed
- [Packer](https://www.packer.io/downloads.html)

## Installation
Detailed instructions coming soon!

## Usage
Detailed usage instructions coming soon!

## Documentation
You can see detailed plugin documentation [here](docs).

## Examples
The [examples](example) directory contains example JSON and HCL configurations.

## Developing
If you wish to work on the Packer KubeVirt Plugin, you'll first need [Go](http://www.golang.org) installed on your machine.

Detailed instructions coming soon!

## Contributing
Your contributions are always welcome! Please take a look at the [contribution guidelines](CONTRIBUTING.md) first.

## License
This project is licensed under the [MIT License](LICENSE.md).
