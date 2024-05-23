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
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1b1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

var _ multistep.Step = &StepRunSourceServer{}

type StepRunSourceServer struct {
	Client client.Client
}

func (s *StepRunSourceServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	pvcName := GenerateUniqueName("source-data-")

	ui.Say("Creating source server vm...")
	cfg := &CloudInitConfig{
		SSHAuthorizedKeys: []string{
			string(config.Comm.SSHPublicKey),
		},
	}
	vm := s.createSourceServerVm(config, pvcName, cfg)
	state.Put("server", vm)
	ui.Say(fmt.Sprintf("Launching source server from image %s ...", *vm.Spec.DataVolumeTemplates[0].Spec.Source.Registry.URL))
	err := s.Client.Create(context.Background(), vm)
	if err != nil {
		err := fmt.Errorf("Error launching source server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	key := types.NamespacedName{Namespace: vm.Namespace, Name: vm.Name}
	for {
		select {
		case <-ctx.Done():
			ui.Error("Build is canceled, stops waiting for source server VM")
			return multistep.ActionHalt
		default:
			err = s.Client.Get(context.Background(), key, vm)
			if err != nil {
				ui.Say("Waiting for source server vm to be running...")
				continue
			}
			if vm.Status.Ready {
				ui.Say("Source server is running")
				time.Sleep(30 * time.Second)
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepRunSourceServer) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	vm := state.Get("server").(*virtv1.VirtualMachine)
	err := s.Client.Delete(context.Background(), vm)
	if err != nil {
		ui.Error(fmt.Sprintf("Error deleting source server vm: %s", err))
	}
}

func (s *StepRunSourceServer) createSourceServerVm(config *Config, pvcName string, cloudInit *CloudInitConfig) *virtv1.VirtualMachine {

	image := fmt.Sprintf("docker://%s", strings.ReplaceAll(config.RunConfig.SourceImage, "docker://", ""))
	var running = true

	return &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "source-server-",
			Namespace:    config.K8sConfig.Namespace,
			Labels: map[string]string{
				"packerBuildName": config.PackerBuildName,
			},
		},
		Spec: virtv1.VirtualMachineSpec{
			Running: &running,
			Template: &virtv1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"packerBuildName": config.PackerBuildName,
					},
				},
				Spec: virtv1.VirtualMachineInstanceSpec{
					Domain: virtv1.DomainSpec{
						CPU: &virtv1.CPU{
							//DedicatedCPUPlacement: true,
						},
						Devices: virtv1.Devices{
							Disks: []virtv1.Disk{
								{
									Name: "datavolumedisk",
									DiskDevice: virtv1.DiskDevice{
										Disk: &virtv1.DiskTarget{
											Bus: "virtio",
										},
									},
								},
								{
									Name: "cloudinit",
									DiskDevice: virtv1.DiskDevice{
										Disk: &virtv1.DiskTarget{
											Bus: "virtio",
										},
									},
								},
							},
							Rng: &virtv1.Rng{},
						},
						Resources: virtv1.ResourceRequirements{
							Requests: v1.ResourceList{
								"cpu":    resource.MustParse(config.ResourceConfig.Cpu),
								"memory": resource.MustParse(config.ResourceConfig.Memory),
							},
						},
					},
					Volumes: []virtv1.Volume{
						{
							Name: "datavolumedisk",
							VolumeSource: virtv1.VolumeSource{
								DataVolume: &virtv1.DataVolumeSource{
									Name: pvcName,
								},
							},
						},
						{
							Name: "cloudinit",
							VolumeSource: virtv1.VolumeSource{
								CloudInitNoCloud: &virtv1.CloudInitNoCloudSource{
									UserData: cloudInit.String(),
								},
							},
						},
					},
				},
			},
			DataVolumeTemplates: []virtv1.DataVolumeTemplateSpec{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pvcName,
						Namespace: "default",
					},
					Spec: cdiv1b1.DataVolumeSpec{
						Source: &cdiv1b1.DataVolumeSource{
							Registry: &cdiv1b1.DataVolumeSourceRegistry{
								URL: &image,
							},
						},
						PVC: &v1.PersistentVolumeClaimSpec{
							AccessModes: []v1.PersistentVolumeAccessMode{
								v1.ReadWriteOnce,
							},
							Resources: v1.VolumeResourceRequirements{
								Requests: v1.ResourceList{
									"storage": resource.MustParse(config.ResourceConfig.Storage),
								},
							},
						},
					},
				},
			},
		},
	}
}
