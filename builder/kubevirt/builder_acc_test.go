// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubevirt

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/port-forward.pkr.hcl
var testBuilderHCL2Basic string

//go:embed test-fixtures/node-port.pkr.hcl
var testBuilderHCL2NodePort string

//go:embed test-fixtures/lb-port.pkr.hcl
var testBuilderHCL2LbPort string

// Run with: PACKER_ACC=1 go test -count 1 -v ./builder/kubevirt/builder_acc_test.go  -timeout=120m
func TestAccKubevirtBuilderPortForward(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "kubevirt_builder_port_forward_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testBuilderHCL2Basic,
		Type:     "kubevirt",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			//logsBytes, err := ioutil.ReadAll(logs)
			_, err = ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			//logsString := string(logsBytes)

			//buildGeneratedDataLog := "kubevirt.basic-example: build generated data: mock-build-data"
			//if matched, _ := regexp.MatchString(buildGeneratedDataLog+".*", logsString); !matched {
			//	t.Fatalf("logs doesn't contain expected foo value %q", logsString)
			//}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccKubevirtBuilderNodePort(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "kubevirt_builder_node_port_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testBuilderHCL2NodePort,
		Type:     "kubevirt",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			//logsBytes, err := ioutil.ReadAll(logs)
			_, err = ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			//logsString := string(logsBytes)

			//buildGeneratedDataLog := "kubevirt.basic-example: build generated data: mock-build-data"
			//if matched, _ := regexp.MatchString(buildGeneratedDataLog+".*", logsString); !matched {
			//	t.Fatalf("logs doesn't contain expected foo value %q", logsString)
			//}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccKubevirtBuilderLbPort(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "kubevirt_builder_node_port_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testBuilderHCL2LbPort,
		Type:     "kubevirt",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			//logsBytes, err := ioutil.ReadAll(logs)
			_, err = ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			//logsString := string(logsBytes)

			//buildGeneratedDataLog := "kubevirt.basic-example: build generated data: mock-build-data"
			//if matched, _ := regexp.MatchString(buildGeneratedDataLog+".*", logsString); !matched {
			//	t.Fatalf("logs doesn't contain expected foo value %q", logsString)
			//}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
