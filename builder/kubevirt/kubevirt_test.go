package kubevirt

import (
	"testing"
)

func TestCommandRunner_Start(t *testing.T) {

	cases := []struct {
		name    string
		arg     []string
		wantErr bool
	}{
		{
			name:    "ls",
			arg:     []string{},
			wantErr: false,
		},
		{
			name:    "invalid_command",
			arg:     []string{},
			wantErr: true,
		},
		{
			name:    "ls",
			arg:     []string{"-la"},
			wantErr: false,
		},
		{
			name:    "sleep",
			arg:     []string{"3"},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cr := &CommandRunner{}
			err := cr.Start(tc.name, tc.arg...)
			if (err != nil) != tc.wantErr {
				t.Errorf("CommandRunner.Start() error = %v, wantErr %v", err, tc.wantErr)
			}
			err = cr.Wait()
			if (err != nil) != tc.wantErr {
				t.Errorf("CommandRunner.Wait() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
