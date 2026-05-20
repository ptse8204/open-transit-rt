package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestAPIContractCheckScript(t *testing.T) {
	root := releaseCandidateRepoRoot(t)
	cmd := exec.Command("sh", "scripts/api-contract-check.sh")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("api contract check failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "api contract check passed") {
		t.Fatalf("unexpected api contract output: %s", out)
	}
}
