package kubevirt

import (
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

// K8sConfigTestCases container for table driven tests
type K8sConfigTestCases struct {
	name     string
	kcInput  *K8sConfig
	expected *K8sConfig
	errors   []error
}

func TestK8sConfigPrepare(t *testing.T) {
	testCases := []K8sConfigTestCases{
		{
			name: "EmptyNamespace",
			kcInput: &K8sConfig{
				Namespace:   "",
				ServiceType: "",
			},
			expected: &K8sConfig{
				Namespace:   "default",
				ServiceType: "",
			},
			errors: []error{},
		},
		{
			name: "NonEmptyNamespace",
			kcInput: &K8sConfig{
				Namespace:   "default",
				ServiceType: v1.ServiceTypeNodePort,
			},
			expected: &K8sConfig{
				Namespace:   "default",
				ServiceType: v1.ServiceTypeNodePort,
			},
			errors: []error{},
		},
		{
			name: "DefaultServicePort",
			kcInput: &K8sConfig{
				Namespace:   "default",
				ServiceType: v1.ServiceTypeLoadBalancer,
				ServicePort: 0,
			},
			expected: &K8sConfig{
				Namespace:   "default",
				ServiceType: v1.ServiceTypeLoadBalancer,
				ServicePort: 22,
			},
			errors: []error{},
		},
		{
			name: "DefaultServicePort2",
			kcInput: &K8sConfig{
				Namespace:   "default",
				ServiceType: v1.ServiceTypeClusterIP,
				ServicePort: 0,
			},
			expected: &K8sConfig{
				Namespace:   "default",
				ServiceType: v1.ServiceTypeClusterIP,
				ServicePort: 22,
			},
			errors: []error{},
		},
	}
	ctx := &interpolate.Context{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.kcInput.Prepare(ctx)
			assert.EqualValues(t, tc.expected, tc.kcInput)
			assert.EqualValues(t, tc.errors, errs)
		})
	}
}
