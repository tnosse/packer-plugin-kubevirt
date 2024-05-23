package kubevirt

import "testing"

func TestCloudInitConfig_String(t *testing.T) {
	tests := []func() (string, CloudInitConfig, string){
		testSSHAuthorizedKey,
	}
	for _, tt := range tests {
		name, cfg, want := tt()
		t.Run(name, func(t *testing.T) {
			if got := cfg.String(); got != want {
				t.Errorf("String() = %v, want %v", got, want)
			}
		})
	}
}

func testSSHAuthorizedKey() (string, CloudInitConfig, string) {

	return "SSHAuthorizedKeys", CloudInitConfig{
			SSHAuthorizedKeys: []string{
				"AKey",
			},
		}, `#cloud-config
ssh_authorized_keys:
- AKey
`
}
