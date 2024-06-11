package kubevirt

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type K8sConfig struct {
	Namespace          string `mapstructure:"namespace"`
	UseServiceNodePort bool   `mapstructure:"use_service_node_port"`
	NodeHost           string `mapstructure:"node_host"`
}

func (c *K8sConfig) Prepare(ctx *interpolate.Context) []error {
	errs := []error{}
	if c.Namespace == "" {
		c.Namespace = "default"
	}

	if c.UseServiceNodePort && c.NodeHost == "" {
		errs = append(errs, fmt.Errorf("node_host cannot be empty when use_service_node_port"))
	}
	return errs
}
