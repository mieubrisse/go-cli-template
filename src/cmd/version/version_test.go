package version

import (
	"bytes"
	"testing"

	"github.com/owner-replaceme/project-replaceme/internal/buildinfo"
)

func TestVersionCommand(t *testing.T) {
	original := buildinfo.Version
	buildinfo.Version = "1.2.3"
	defer func() { buildinfo.Version = original }()

	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)

	Cmd.SetArgs([]string{})
	if err := Cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want := "1.2.3\n"
	if got != want {
		t.Errorf("version output = %q, want %q", got, want)
	}
}
