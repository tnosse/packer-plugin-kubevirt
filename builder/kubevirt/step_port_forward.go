package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	virtv1 "kubevirt.io/api/core/v1"
	"os/exec"
	"strings"
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
	err := p.cmd.Wait()
	if err != nil && !strings.Contains(err.Error(), "signal: killed") {
		return err
	}
	return nil
}

func (p *PortForward) Wait() error {
	return p.cmd.Wait()
}

type StepPortForward struct {
}

func (s *StepPortForward) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("server").(*virtv1.VirtualMachine)

	pf := &PortForward{}
	err := pf.Start("virtctl", "port-forward", fmt.Sprintf("vm/%s.%s", vm.Name, vm.Namespace), "2222:22")

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
