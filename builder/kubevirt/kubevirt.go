package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/random"
	"os/exec"
	"strings"
)

func GenerateUniqueName(prefix string) string {
	randPart := random.AlphaNumLower(6) // generates a random number between 0 and 999
	uniqueName := fmt.Sprintf("%s%s", prefix, randPart)
	return uniqueName
}

type CommandRunner struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func (p *CommandRunner) Start(name string, arg ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)
	p.cmd = cmd
	p.cancel = cancel
	return cmd.Start()
}

func (p *CommandRunner) Stop() error {
	p.cancel()
	err := p.cmd.Wait()
	if err != nil && !strings.Contains(err.Error(), "signal: killed") {
		return err
	}
	return nil
}

func (p *CommandRunner) Wait() error {
	return p.cmd.Wait()
}
