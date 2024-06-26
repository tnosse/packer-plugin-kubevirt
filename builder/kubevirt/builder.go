// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1b1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	cfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

const BuilderId = "tnosse.kubevirt"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config
	K8sConfig           K8sConfig      `mapstructure:",squash"`
	ImageConfig         ImageConfig    `mapstructure:",squash"`
	RunConfig           RunConfig      `mapstructure:",squash"`
	ResourceConfig      ResourceConfig `mapstructure:",squash"`

	ctx interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) (generatedVars []string, warnings []string, err error) {
	err = config.Decode(&b.config, &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)

	if err != nil {
		return nil, nil, err
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.K8sConfig.Prepare(&b.config.ctx)...)
	errs = packer.MultiErrorAppend(errs, b.config.ImageConfig.Prepare(&b.config.ctx)...)
	errs = packer.MultiErrorAppend(errs, b.config.RunConfig.Prepare(&b.config.ctx, &b.config.K8sConfig, &b.config.Comm)...)
	errs = packer.MultiErrorAppend(errs, b.config.ResourceConfig.Prepare(&b.config.ctx, &b.config.Comm)...)

	if b.config.Comm.Type != "ssh" {
		return nil, nil, fmt.Errorf("Only 'ssh' is supported for now")
	}

	var buildGeneratedData []string
	return buildGeneratedData, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {

	// Create a kubernetes client
	k8sConfig, err := cfg.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting kubernetes config: %s", err)
	}

	k8sClient, err := client.New(k8sConfig, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client: %s", err)
	}
	virtv1.AddToScheme(k8sClient.Scheme())
	cdiv1b1.AddToScheme(k8sClient.Scheme())

	// Setup the state bag and initial state for the steps
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{}
	steps = append(steps,
		&communicator.StepSSHKeyGen{
			CommConf:            &b.config.Comm,
			SSHTemporaryKeyPair: b.config.Comm.SSHTemporaryKeyPair,
		})
	steps = append(steps, &communicator.StepDumpSSHKey{Path: "debug-key.crt", SSH: &b.config.Comm.SSH})
	steps = append(steps, &StepRunSourceServer{Client: k8sClient})

	if b.config.K8sConfig.ServiceType == "" {
		steps = append(steps, &StepPortForward{})
	}

	steps = append(steps, &communicator.StepConnect{
		Config: &b.config.Comm,
		Host: func(bag multistep.StateBag) (string, error) {
			if b.config.K8sConfig.ServiceType == "" {
				return "localhost", nil
			} else {
				return b.config.Comm.SSHHost, nil
			}
		},
		SSHConfig: b.config.Comm.SSHConfigFunc(),
	})

	steps = append(steps, &commonsteps.StepProvision{})
	steps = append(steps, &StepCopyImage{Client: k8sClient})

	state.Put("image", map[string]interface{}{
		"image": &b.config.ImageConfig.OutputImageFile,
	})

	// Run!
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	// If there was an error, return that
	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	artifact := &Artifact{
		// Add the builder generated data to the artifact StateData so that post-processors
		// can access them.
		StateData: map[string]interface{}{"image": state.Get("image")},
	}
	return artifact, nil
}
