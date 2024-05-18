package kubevirt

import "testing"

func TestPortForward_StartAndStop(t *testing.T) {
	pf := PortForward{}
	err := pf.Start("sleep", "10")
	if err != nil {
		t.Fatalf("Failed to start port forward: %s", err)
	}
	err = pf.Stop()
	if err != nil {
		t.Fatalf("Failed to stop port forward: %s", err)
	}
}

func TestPortForward_StartAndWait(t *testing.T) {
	pf := PortForward{}
	err := pf.Start("sleep", "2")
	if err != nil {
		t.Fatalf("Failed to start port forward: %s", err)
	}
	err = pf.Wait()
	if err != nil {
		t.Fatalf("Failed to wait for stop: %s", err)
	}
}
