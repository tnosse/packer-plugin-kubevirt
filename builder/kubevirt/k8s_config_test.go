package kubevirt

import (
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
			},
			expected: &K8sConfig{
				Namespace:          "default",
				UseServiceNodePort: false,
			},
			errors: []error{},
		},
		{
			name: "NonEmptyNamespace",
			kcInput: &K8sConfig{
				Namespace:          "kube-system",
				UseServiceNodePort: true,
			},
			expected: &K8sConfig{
				Namespace:          "kube-system",
				UseServiceNodePort: true,
			},
			errors: []error{},
		},
		{
			name: "EmptyNodeHost",
			kcInput: &K8sConfig{
				Namespace:          "default",
				UseServiceNodePort: true,
			},
			expected: &K8sConfig{
				Namespace:          "default",
				UseServiceNodePort: true,
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
