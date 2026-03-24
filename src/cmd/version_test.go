package cmd

import (
	"bytes"
	"testing"

	"github.com/owner-replaceme/project-replaceme/internal/version"
)

func TestVersionCommand(t *testing.T) {
	version.Version = "1.2.3"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{versionCmdStr})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want := "1.2.3\n"
	if got != want {
		t.Errorf("version output = %q, want %q", got, want)
	}
}
