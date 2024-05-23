package kubevirt

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"net"
)

type RunConfig struct {
	SSHUsername string `mapstructure:"ssh_username"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context, comm *communicator.Config) []error {
	var errs []error

	comm.SSHUsername = c.SSHUsername
	comm.SSHHost = "localhost"
	comm.Type = "ssh"

	port, err := getAvailablePort()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting available port: %s", err))
	} else {
		comm.SSHPort = port
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
