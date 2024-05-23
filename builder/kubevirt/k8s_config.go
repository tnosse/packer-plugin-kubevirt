package kubevirt

import (
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type K8sConfig struct {
	Namespace string `mapstructure:"namespace"`
}

func (c *K8sConfig) Prepare(ctx *interpolate.Context) []error {
	if c.Namespace == "" {
		c.Namespace = "default"
	}
	return []error{}
}
