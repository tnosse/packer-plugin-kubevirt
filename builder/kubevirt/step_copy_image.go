package kubevirt

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"io"
	v1 "k8s.io/api/core/v1"
	virtv1 "kubevirt.io/api/core/v1"
	"os"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ multistep.Step = &StepCopyImage{}

type StepCopyImage struct {
	Client client.Client
}

func (s *StepCopyImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)
	vm := state.Get("server").(*virtv1.VirtualMachine)

	if config.ImageConfig.SkipExtractImage {
		ui.Say("Skipping extract of VM image")
		return multistep.ActionContinue
	}

	ui.Say("Pausing vm...")
	pauseCmd := &CommandRunner{}
	err := pauseCmd.Start("virtctl", "-n", vm.Namespace, "pause", "vm", vm.Name)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to pause vm: %s", err))
		return multistep.ActionHalt
	}
	err = pauseCmd.Wait()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to wait for pause cmd: %s", err))
		return multistep.ActionHalt
	}

	ui.Say("Getting virt-handler pod...")
	list := &v1.PodList{}
	err = s.Client.List(context.TODO(), list, client.InNamespace(vm.Namespace), client.MatchingLabels{
		"vm.kubevirt.io/name": vm.Name,
	})
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get virt-handler pod: %s", err))
		return multistep.ActionHalt
	}
	podName := list.Items[0].Name
	if len(list.Items) != 1 {
		ui.Error(fmt.Sprintf("Expected 1 virt-handler pod, but found %d", len(list.Items)))
		return multistep.ActionHalt
	}
	ui.Say(fmt.Sprintf("Using %s for copying image", podName))

	ui.Say("Copying vm image...")
	tmpDir, err := os.MkdirTemp("", "packer-kubevirt-image")
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create temporary directory: %s", err))
		return multistep.ActionHalt
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	cmd := exec.Command("kubectl",
		"--namespace", vm.Namespace,
		"cp",
		"--retries", "-1",
		fmt.Sprintf("%s:/var/run/kubevirt-private/vmi-disks/datavolumedisk", podName),
		tmpDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to copy image: %s\n%s", err, string(out)))
		ui.Say(fmt.Sprintf("cmd args: %s", cmd.Args))
		return multistep.ActionHalt
	}

	err = moveFile(tmpDir+string(os.PathSeparator)+"disk.img", config.ImageConfig.OutputImageFile)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to move temporary file: %s", err))
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Successfully copied image: %s", string(out)))

	return multistep.ActionContinue
}

func (s *StepCopyImage) Cleanup(state multistep.StateBag) {}

func moveFile(source, destination string) error {
	// First attempt a simple rename
	err := os.Rename(source, destination)
	if err == nil {
		return nil
	}

	// If rename fails (e.g., cross-device link), fall back to copy
	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the contents
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(destination) // Clean up the partial file
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	// Ensure all data is written to disk
	if err := dst.Sync(); err != nil {
		dst.Close()
		os.Remove(destination)
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	// Close files before removing source
	src.Close()
	dst.Close()

	// Remove the source file
	if err := os.Remove(source); err != nil {
		// If we can't remove the source, we should still return success
		// since the copy was successful, but log the error
		return fmt.Errorf("file copied successfully but failed to remove source file: %w", err)
	}

	return nil
}
