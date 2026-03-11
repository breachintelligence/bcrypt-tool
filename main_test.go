package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "bcrypt-tool-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "bcrypt-tool")
	cmd := exec.Command("go", "build", "-o", binaryPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		panic("failed to build binary: " + string(out))
	}

	os.Exit(m.Run())
}

func runTool(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error running bcrypt-tool: %v", err)
		}
	}
	return strings.TrimSpace(string(out)), exitCode
}

func TestHashAndMatch(t *testing.T) {
	password := "testpassword123"

	hash, code := runTool(t, "hash", password)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.HasPrefix(hash, "$2a$") {
		t.Fatalf("expected hash to start with $2a$, got %q", hash)
	}

	// Verify the generated hash matches the original password
	out, code := runTool(t, "match", password, hash)
	if code != 0 || out != "yes" {
		t.Fatalf("expected match to return 'yes' (exit 0), got %q (exit %d)", out, code)
	}

	// Verify a wrong password does not match
	out, code = runTool(t, "match", "wrongpassword", hash)
	if code != 1 || out != "no" {
		t.Fatalf("expected mismatch to return 'no' (exit 1), got %q (exit %d)", out, code)
	}
}

func TestHashWithCost(t *testing.T) {
	password := "costtest"
	costVal := "4"

	hash, code := runTool(t, "hash", password, costVal)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	// Verify cost is reported correctly
	out, code := runTool(t, "cost", hash)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if out != costVal {
		t.Fatalf("expected cost %q, got %q", costVal, out)
	}

	// Verify the hash still matches
	out, code = runTool(t, "match", password, hash)
	if code != 0 || out != "yes" {
		t.Fatalf("expected match to return 'yes' (exit 0), got %q (exit %d)", out, code)
	}
}

func TestNoArgs(t *testing.T) {
	_, code := runTool(t, "hash")
	if code != 2 {
		t.Fatalf("expected exit code 2 for missing args, got %d", code)
	}
}
