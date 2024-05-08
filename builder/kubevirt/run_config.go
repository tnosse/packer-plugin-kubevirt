package kubevirt

import "github.com/hashicorp/packer-plugin-sdk/communicator"

type RunConfig struct {
	Comm communicator.Config `mapstructure:",squash"`
}
