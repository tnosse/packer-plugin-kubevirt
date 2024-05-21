package kubevirt

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ multistep.Step = &StepDownloadImage{}

type StepDownloadImage struct {
	client.Client
}

func (s *StepDownloadImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	//config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Closing server...")
	pod := state.Get("server").(*v1.Pod)
	_ = s.Client.Delete(context.TODO(), pod)

	converterPod := s.createConverterPod()
	err := s.Client.Create(context.TODO(), converterPod)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Start downloading image...")
	cp := &CommandRunner{}
	err = cp.Start("kubectl", "-n", "default", "cp", "--retries=-1", "--container=converter", "/image/disk.img", "disk.img")
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Waiting for download to complete...")
	err = cp.cmd.Wait()
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *StepDownloadImage) Cleanup(state multistep.StateBag) {
	//TODO implement me
	panic("implement me")
}

func (*StepDownloadImage) createConverterPod() *v1.Pod {
	return &v1.Pod{
		// Pod definition goes here
	}
}
