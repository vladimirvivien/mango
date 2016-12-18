package gobuilder

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	testProjectPath = "./testproject"
)

func TestBuilderNew(t *testing.T) {
	b := New()
	if b == nil {
		t.Error(b)
	}
	if b.ShouldInstall != false {
		t.Error("expected ShouldInstall to be false")
	}
	if b.ShouldBeVerbose != false {
		t.Error("expected ShouldBeVerbose to be false")
	}
	if b.ShouldDetectRace != false {
		t.Error("expected ShouldDetecRace to be false")
	}
	if b.ParallelBuilds != runtime.NumCPU() {
		t.Errorf("expected ParallelBuilds %d, but got %d",
			runtime.NumCPU(),
			b.ParallelBuilds,
		)
	}
}

func TestBuilder_Simple(t *testing.T) {
	binPath := filepath.Join(testProjectPath, "bin")
	defer func() {
		os.RemoveAll(binPath)
	}()
	b := New()
	b.SrcRoot = testProjectPath
	b.Packages = []string{"./hello"}
	b.Build()
	if _, err := os.Stat(binPath); err != nil {
		t.Fatalf("expected bin directory in project path: %v", err)
	}
}

// TODO add more tests for multiple package builds, simple tests and multi-package tests.
