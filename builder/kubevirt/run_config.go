package kubevirt

import (
	"fmt"
	"net"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type RunConfig struct {
	communicator.SSH        `mapstructure:",squash"`
	SourceImage             string `mapstructure:"source_image"`
	SourceServerWaitTime    int    `mapstructure:"source_server_wait_time"`
	CloudInitDataVolumeName string `mapstructure:"cloud_init_data_volume_name"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context, k8sConfig *K8sConfig, comm *communicator.Config) []error {
	var errs []error

	comm.SSH = c.SSH
	comm.Type = "ssh"
	comm.SSHHost = c.SSHHost

	port, err := getAvailablePort()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting available port: %s", err))
	} else {
		comm.SSHPort = port
	}

	if len(c.SourceImage) < 1 {
		errs = append(errs, fmt.Errorf("the 'source_image' property must be specified"))
	}

	if c.SourceServerWaitTime < 0 {
		errs = append(errs, fmt.Errorf("the 'source_server_wait_time' property must be a positive integer"))
	} else if c.SourceServerWaitTime == 0 {
		c.SourceServerWaitTime = 30
	}

	if len(c.CloudInitDataVolumeName) == 0 {
		c.CloudInitDataVolumeName = "cloudinit"
	}

	return errs
}

func getAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
