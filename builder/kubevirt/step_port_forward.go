package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	virtv1 "kubevirt.io/api/core/v1"
)

var _ multistep.Step = &StepPortForward{}

type StepPortForward struct {
}

func (s *StepPortForward) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	vm := state.Get("server").(*virtv1.VirtualMachine)

	pf := &CommandRunner{}
	err := pf.Start(
		"virtctl",
		"port-forward",
		fmt.Sprintf("vm/%s.%s", vm.Name, vm.Namespace),
		fmt.Sprintf("%d:22", config.Comm.SSHPort))

	ui.Say(fmt.Sprintf("Starting port forward using local port %d ...", config.Comm.SSHPort))
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
	pf := state.Get("port-forward").(*CommandRunner)
	ui.Say("Stopping port forward...")
	if err := pf.Stop(); err != nil && err.Error() != "signal: killed" {
		ui.Error(err.Error())
	}
}
