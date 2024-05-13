package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ multistep.Step = &StepRunSourceServer{}

type StepRunSourceServer struct {
	Client client.Client
}

func (s *StepRunSourceServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	pvc := state.Get("pvc").(*v1.PersistentVolumeClaim)

	ui.Say("Creating source server pod...")
	key, pod := s.createSourceServerManifest(config, pvc.Name)
	err := s.Client.Create(context.Background(), pod)
	if err != nil {
		err := fmt.Errorf("Error launching source server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	for {
		pod := &v1.Pod{}
		err = s.Client.Get(context.Background(), key, pod)
		if err != nil {
			ui.Say("Waiting for source server pod to be running...")
			continue
		}
		if pod.Status.Phase == v1.PodRunning {
			ui.Say("Source server is running")
			state.Put("server", pod)
			break
		}
	}

	return multistep.ActionContinue
}

func (s *StepRunSourceServer) Cleanup(state multistep.StateBag) {
	pod := state.Get("server").(*v1.Pod)
	_ = s.Client.Delete(context.Background(), pod)
}

func (s *StepRunSourceServer) createSourceServerManifest(config *Config, pvcName string) (types.NamespacedName, *v1.Pod) {
	key := types.NamespacedName{Namespace: "default", Name: "source-server"}
	return key, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "server",
					Image: "linuxserver/openssh-server:latest",
					Env: []v1.EnvVar{
						{
							Name:  "PUBLIC_KEY",
							Value: string(config.Comm.SSHPublicKey),
						},
						{
							Name:  "SUDO_ACCESS",
							Value: "true",
						},
						{
							Name:  "PUID",
							Value: "1000",
						},
						{
							Name:  "PGID",
							Value: "1000",
						},
						{
							Name:  "TZ",
							Value: "Europe/Stockholm",
						},
						{
							Name:  "USER_NAME",
							Value: "packer",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "disk",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				},
			},
		},
	}
}
