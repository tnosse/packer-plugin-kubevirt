package kubevirt

import (
	"bytes"
	"text/template"
)

const cloudInitTemplate = `#cloud-config
{{- if eq .Hostname nil }}
hostname: {{.Hostname}}
{{- end }}
{{- if eq .Hostname nil }}
fqdn: {{.FQDN}}
{{- end }}
ssh_authorized_keys:
{{- range .SSHAuthorizedKeys }}
- {{ . }}
{{- end }}
`

type CloudInitConfig struct {
	Hostname          string
	FQDN              string
	Commands          []string
	SSHAuthorizedKeys []string
}

func (cfg *CloudInitConfig) String() string {
	tmpl, err := template.New("cloudInit").Parse(cloudInitTemplate)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, cfg); err != nil {
		panic(err)
	}
	return tpl.String()
}
