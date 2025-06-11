package kubevirt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/martian/v3/log"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1b1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
				time.Sleep(time.Duration(config.RunConfig.SourceServerWaitTime) * time.Second)

				if config.K8sConfig.ServiceType != "" {
					ui.Sayf("Creating VM service, type %s, for SSH provisioning.", config.K8sConfig.ServiceType)
					svc := s.createService(config, vm)
					if err := s.Client.Create(context.Background(), svc); err != nil {
						ui.Say("Failed to create VM service")
						ui.Error(err.Error())
						return multistep.ActionHalt
					}

					// Wait for lb IP
					if config.K8sConfig.ServiceType == v1.ServiceTypeLoadBalancer {
						for {
							select {
							case <-ctx.Done():
								ui.Error("Build is canceled, stops waiting for source server VM service")
								return multistep.ActionHalt
							default:
								err = s.Client.Get(context.Background(), key, svc)
								if err != nil {
									ui.Errorf("Error waiting for VM service: %s", err.Error())
									return multistep.ActionHalt
								}
								if len(svc.Status.LoadBalancer.Ingress) < 1 {
									ui.Say("VM service LoadBalancer IP is not ready")
									time.Sleep(3 * time.Second)
								} else {
									config.Comm.SSHPort = config.K8sConfig.ServicePort
									config.Comm.SSHHost = svc.Status.LoadBalancer.Ingress[0].IP
									return multistep.ActionContinue
								}
							}
						}
					} else if config.K8sConfig.ServiceType == v1.ServiceTypeNodePort {
						config.Comm.SSHPort = int(svc.Spec.Ports[0].NodePort)
					} else {
						config.Comm.SSHPort = config.K8sConfig.ServicePort
						config.Comm.SSHHost = svc.Spec.ClusterIP
					}
				}
				connectStrategy := "port-forward"
				if config.K8sConfig.ServiceType != "" {
					connectStrategy = string(config.K8sConfig.ServiceType) + " service"
				}
				ui.Sayf("Using %s %s:%d for SSH.", connectStrategy, config.Comm.SSHHost, config.Comm.SSHPort)
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepRunSourceServer) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)
	vm := state.Get("server").(*virtv1.VirtualMachine)

	err := s.Client.Delete(context.Background(), vm)
	if err != nil {
		ui.Error(fmt.Sprintf("Error deleting source server vm: %s", err))
	}
	if config.K8sConfig.ServiceType != "" {
		err = s.Client.Delete(context.Background(), &v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: vm.Namespace, Name: vm.Name}})
		if err != nil {
			ui.Error(fmt.Sprintf("Error deleting source server vm service: %s", err))
		}
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
					Networks: []virtv1.Network{
						*virtv1.DefaultPodNetwork(),
					},
					Domain: virtv1.DomainSpec{
						CPU: &virtv1.CPU{
							//DedicatedCPUPlacement: true,
						},
						Devices: virtv1.Devices{
							Interfaces: []virtv1.Interface{
								*virtv1.DefaultMasqueradeNetworkInterface(),
							},
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
									Name: config.RunConfig.CloudInitDataVolumeName,
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
							Name: config.RunConfig.CloudInitDataVolumeName,
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
						Namespace: config.K8sConfig.Namespace,
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

func (s *StepRunSourceServer) createService(config *Config, vm *virtv1.VirtualMachine) *v1.Service {
	ipFamilyPolicy := v1.IPFamilyPolicySingleStack
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vm.Name,
			Namespace: vm.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeClusterIP,
			Selector: vm.Labels,
			Ports: []v1.ServicePort{
				{
					Port:     22,
					Protocol: v1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: 22,
					},
				},
			},
			IPFamilies: []v1.IPFamily{
				v1.IPv4Protocol,
			},
			IPFamilyPolicy: &ipFamilyPolicy,
		},
	}

	switch config.K8sConfig.ServiceType {
	case v1.ServiceTypeNodePort:
		svc.Spec.Type = v1.ServiceTypeNodePort
		if config.K8sConfig.ServicePort > 0 {
			svc.Spec.Ports[0].NodePort = int32(config.K8sConfig.ServicePort)
		}
	case v1.ServiceTypeLoadBalancer:
		svc.Spec.Type = v1.ServiceTypeLoadBalancer
		svc.Spec.Ports[0].Port = int32(config.K8sConfig.ServicePort)
	case v1.ServiceTypeClusterIP:
		svc.Spec.Ports[0].Port = int32(config.K8sConfig.ServicePort)
	default:
		log.Errorf("unknown service type: %s", config.K8sConfig.ServiceType)
	}

	return svc
}
