package kubevirt

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceConfig struct {
	Storage string `mapstructure:"storage"`
	Memory  string `mapstructure:"memory"`
	Cpu     string `mapstructure:"cpu"`
}

func (c *ResourceConfig) Prepare(i *interpolate.Context, c2 *communicator.Config) []error {
	var errs []error

	if c.Storage == "" {
		c.Storage = "500Mi"
	}

	if c.Memory == "" {
		c.Memory = "50Mi"
	}

	if c.Cpu == "" {
		c.Cpu = "100m"
	}

	if _, err := resource.ParseQuantity(c.Storage); err != nil {
		errs = append(errs, fmt.Errorf("error parsing Storage quantity %q: %v", c.Storage, err))
	}

	if _, err := resource.ParseQuantity(c.Memory); err != nil {
		errs = append(errs, fmt.Errorf("error parsing Memory quantity %q: %v", c.Memory, err))
	}

	if _, err := resource.ParseQuantity(c.Cpu); err != nil {
		errs = append(errs, fmt.Errorf("error parsing Cpu quantity %q: %v", c.Cpu, err))
	}

	return errs
}
