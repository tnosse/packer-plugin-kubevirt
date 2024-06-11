package kubevirt

import (
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type K8sConfig struct {
	Namespace          string `mapstructure:"namespace"`
	UseServiceNodePort bool   `mapstructure:"use_service_node_port"`
}

func (c *K8sConfig) Prepare(ctx *interpolate.Context) []error {
	errs := []error{}

	if c.Namespace == "" {
		c.Namespace = "default"
	}

	return errs
}
