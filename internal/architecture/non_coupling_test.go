package architecture

import (
	"os/exec"
	"strings"
	"testing"
)

func TestTripUpdatesDoesNotCoupleIntoExistingPhaseEntrypoints(t *testing.T) {
	for _, pkg := range []string{
		"open-transit-rt/cmd/feed-vehicle-positions",
		"open-transit-rt/cmd/telemetry-ingest",
		"open-transit-rt/cmd/gtfs-studio",
	} {
		output, err := exec.Command("go", "list", "-deps", pkg).CombinedOutput()
		if err != nil {
			t.Fatalf("go list %s: %v\n%s", pkg, err, output)
		}
		deps := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, dep := range deps {
			switch dep {
			case "open-transit-rt/internal/prediction", "open-transit-rt/internal/feed/tripupdates":
				t.Fatalf("%s unexpectedly depends on %s", pkg, dep)
			}
		}
	}
}
