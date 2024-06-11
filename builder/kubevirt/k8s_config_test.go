package kubevirt

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/stretchr/testify/assert"
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
				Namespace:          "",
				UseServiceNodePort: false,
				NodeHost:           "localhost",
			},
			expected: &K8sConfig{
				Namespace:          "default",
				UseServiceNodePort: false,
				NodeHost:           "localhost",
			},
			errors: []error{},
		},
		{
			name: "NonEmptyNamespace",
			kcInput: &K8sConfig{
				Namespace:          "kube-system",
				UseServiceNodePort: true,
				NodeHost:           "localhost",
			},
			expected: &K8sConfig{
				Namespace:          "kube-system",
				UseServiceNodePort: true,
				NodeHost:           "localhost",
			},
			errors: []error{},
		},
		{
			name: "EmptyNodeHost",
			kcInput: &K8sConfig{
				Namespace:          "default",
				UseServiceNodePort: true,
				NodeHost:           "",
			},
			expected: &K8sConfig{
				Namespace:          "default",
				UseServiceNodePort: true,
				NodeHost:           "",
			},
			errors: []error{
				fmt.Errorf("node_host cannot be empty when use_service_node_port"),
			},
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
