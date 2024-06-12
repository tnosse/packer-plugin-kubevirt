package kubevirt

import (
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	v1 "k8s.io/api/core/v1"
)

type K8sConfig struct {
	Namespace   string         `mapstructure:"namespace"`
	ServiceType v1.ServiceType `mapstructure:"service_type"`
	ServicePort int            `mapstructure:"service_port"`
}

func (c *K8sConfig) Prepare(ctx *interpolate.Context) []error {
	errs := []error{}

	if c.Namespace == "" {
		c.Namespace = "default"
	}

	if c.ServicePort == 0 && (c.ServiceType == v1.ServiceTypeLoadBalancer || c.ServiceType == v1.ServiceTypeClusterIP) {
		c.ServicePort = 22
	}

	return errs
}
