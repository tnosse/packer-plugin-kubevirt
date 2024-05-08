package kubevirt

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"k8s.io/client-go/kubernetes"
)

var _ multistep.Step = &StepRunSourceServer{}

type StepRunSourceServer struct {
	Client *kubernetes.Clientset
}

func (s *StepRunSourceServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *StepRunSourceServer) Cleanup(state multistep.StateBag) {
}
