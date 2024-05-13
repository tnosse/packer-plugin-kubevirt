package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ multistep.Step = &StepCreatePVC{}

type StepCreatePVC struct {
	Client client.Client
}

func (s *StepCreatePVC) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Creating pvc for source server pod...")
	_, pvc := s.createPVCManifest()
	err := s.Client.Create(context.Background(), pvc)
	if err != nil {
		err := fmt.Errorf("Error creating PVC for source server pod: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("pvc", pvc)
	return multistep.ActionContinue
}

func (s *StepCreatePVC) Cleanup(state multistep.StateBag) {
	if pvc := state.Get("pvc"); pvc != nil {
		_ = s.Client.Delete(context.Background(), pvc.(client.Object))
	}
}

func (s *StepCreatePVC) createPVCManifest() (types.NamespacedName, *v1.PersistentVolumeClaim) {
	key := types.NamespacedName{Namespace: "default", Name: "source-server"}
	return key, &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					"storage": resource.MustParse("100Mi"),
				},
			},
		},
	}
}
