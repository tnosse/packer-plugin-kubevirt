package kubevirt

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"os"
)

type ImageConfig struct {
	Storage          string `mapstructure:"storage"`
	SkipExtractImage bool   `mapstructure:"skip_extract_image"`
	OutputImageFile  string `mapstructure:"output_image_file"`
}

func (c *ImageConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error

	if c.OutputImageFile == "" && !c.SkipExtractImage {
		return append(errs, fmt.Errorf("output_image_file is required"))
	}

	if _, err := os.Stat(c.OutputImageFile); err == nil {
		return append(errs, fmt.Errorf("output_image_file %s already exists", c.OutputImageFile))
	}

	outputFile, err := os.Create(c.OutputImageFile)
	if err != nil {
		return append(errs, fmt.Errorf("invalid output_image_file path: %s", err))
	}
	err = outputFile.Close()
	if err != nil {
		return append(errs, fmt.Errorf("failed to close output_image_file: %s", err))
	}
	err = os.RemoveAll(outputFile.Name())
	if err != nil {
		return append(errs, fmt.Errorf("failed to remove tmp output_image_file: %s", err))
	}

	return errs
}
