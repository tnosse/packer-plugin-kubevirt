package kubevirt

import (
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"testing"
)

func TestResourceConfig_Prepare(t *testing.T) {
	type fields struct {
		Storage string
		Memory  string
		Cpu     string
	}
	tests := []struct {
		name   string
		fields fields
		wantErrors bool
	}{
		{
			name:   "No errors with correct quantities",
			fields: fields{Storage: "500Mi", Memory: "50Mi", Cpu: "100m"},
			wantErrors: false,
		},
		{
			name:   "No errors with empty quantities",
			fields: fields{Storage: "", Memory: "", Cpu: ""},
			wantErrors: false,
		},
		{
			name:   "Error with incorrect Storage quantity",
			fields: fields{Storage: "incorrect", Memory: "50Mi", Cpu: "100m"},
			wantErrors: true,
		},
		{
			name:   "Error with incorrect Memory quantity",
			fields: fields{Storage: "500Mi", Memory: "incorrect", Cpu: "100m"},
			wantErrors: true,
		},
		{
			name:   "Error with incorrect Cpu quantity",
			fields: fields{Storage: "500Mi", Memory: "50Mi", Cpu: "incorrect"},
			wantErrors: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ResourceConfig{
				Storage: tt.fields.Storage,
				Memory:  tt.fields.Memory,
				Cpu:     tt.fields.Cpu,
			}
			i := interpolate.Context{}
			c2 := communicator.Config{}
			errs := c.Prepare(&i, &c2)
			if !tt.wantErrors && len(errs) > 0 {
				t.Errorf("Unexpected errors: %v", errs)
			} else if tt.wantErrors && len(errs) == 0 {
				t.Errorf("Expected errors, but got none")
			}
		})
	}
}