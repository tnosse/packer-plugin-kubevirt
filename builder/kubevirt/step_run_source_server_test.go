package kubevirt

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	virtv1 "kubevirt.io/api/core/v1"
	"testing"
)

func TestStepRunSourceServer_createService(t *testing.T) {
	vm := &virtv1.VirtualMachine{ObjectMeta: metav1.ObjectMeta{Name: "vm", Namespace: "default"}}
	ipPolicy := v1.IPFamilyPolicySingleStack
	testCases := map[string]struct {
		config *Config
		vm     *virtv1.VirtualMachine
		want   *v1.Service
	}{
		"defaultConfig": {
			config: &Config{
				K8sConfig: K8sConfig{},
			},
			vm: vm,
			want: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeClusterIP,
					Ports: []v1.ServicePort{
						{
							Port:     22,
							Protocol: v1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								IntVal: 22,
							},
						},
					},
					IPFamilies:     []v1.IPFamily{v1.IPv4Protocol},
					IPFamilyPolicy: &ipPolicy,
				},
			},
		},
		"nodePortConfig": {
			config: &Config{
				K8sConfig: K8sConfig{
					ServiceType: v1.ServiceTypeNodePort,
				},
			},
			vm: vm,
			want: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeNodePort,
					Ports: []v1.ServicePort{
						{
							Port:     22,
							NodePort: 0,
							Protocol: v1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								IntVal: 22,
							},
						},
					},
					IPFamilies:     []v1.IPFamily{v1.IPv4Protocol},
					IPFamilyPolicy: &ipPolicy,
				},
			},
		},
		"nodePortConfigWithServicePort": {
			config: &Config{
				K8sConfig: K8sConfig{
					ServiceType: v1.ServiceTypeNodePort,
					ServicePort: 30000,
				},
			},
			vm: vm,
			want: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeNodePort,
					Ports: []v1.ServicePort{
						{
							Port:     22,
							NodePort: 30000,
							Protocol: v1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								IntVal: 22,
							},
						},
					},
					IPFamilies:     []v1.IPFamily{v1.IPv4Protocol},
					IPFamilyPolicy: &ipPolicy,
				},
			},
		},
		"loadBalancerConfig": {
			config: &Config{
				K8sConfig: K8sConfig{
					ServicePort: 0,
					ServiceType: v1.ServiceTypeLoadBalancer,
				},
			},
			vm: vm,
			want: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Port:     22,
							Protocol: v1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								IntVal: 22,
							},
						},
					},
					IPFamilies:     []v1.IPFamily{v1.IPv4Protocol},
					IPFamilyPolicy: &ipPolicy,
				},
			},
		},
		"loadBalancerConfigWithServicePort": {
			config: &Config{
				K8sConfig: K8sConfig{
					ServicePort: 2222,
					ServiceType: v1.ServiceTypeLoadBalancer,
				},
			},
			vm: vm,
			want: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Port:     2222,
							Protocol: v1.ProtocolTCP,
							TargetPort: intstr.IntOrString{
								IntVal: 22,
							},
						},
					},
					IPFamilies:     []v1.IPFamily{v1.IPv4Protocol},
					IPFamilyPolicy: &ipPolicy,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			s := &StepRunSourceServer{}
			assert.Len(t, tc.config.K8sConfig.Prepare(nil), 0, "should have been empty")
			service := s.createService(tc.config, tc.vm)
			assert.Equal(t, tc.want.Spec, service.Spec)
		})
	}
}
