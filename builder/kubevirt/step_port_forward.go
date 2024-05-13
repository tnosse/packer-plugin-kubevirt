package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"os/exec"
)

var _ multistep.Step = &StepPortForward{}

type PortForward struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func (p *PortForward) Start(name string, arg ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)
	p.cmd = cmd
	p.cancel = cancel
	return cmd.Start()
}

func (p *PortForward) Stop() error {
	p.cancel()
	return p.cmd.Wait()
}

type StepPortForward struct {
}

func (s *StepPortForward) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	pf := &PortForward{}
	err := pf.Start("kubectl", "port-forward", "-n", "default", "source-server", "2222:2222")

	ui.Say("Starting port forward...")
	if err != nil {
		err := fmt.Errorf("Port-forward finished with error: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("port-forward", pf)

	return multistep.ActionContinue
}

func (s *StepPortForward) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	pf := state.Get("port-forward").(*PortForward)
	ui.Say("Stopping port forward...")
	if err := pf.Stop(); err != nil && err.Error() != "signal: killed" {
		ui.Error(err.Error())
	}
}
