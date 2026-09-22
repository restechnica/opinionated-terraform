package root

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetSubcommandMap(t *testing.T) {
	t.Run("IncludesCompletionCommands", func(t *testing.T) {
		cmd := &cobra.Command{Use: "root"}

		got := GetSubcommandMap(cmd)

		assert.True(t, got["__complete"])
		assert.True(t, got["__completeNoDesc"])
	})

	t.Run("IncludesSubcommandNames", func(t *testing.T) {
		cmd := &cobra.Command{Use: "root"}
		cmd.AddCommand(&cobra.Command{Use: "version"})

		got := GetSubcommandMap(cmd)

		assert.True(t, got["version"])
	})

	t.Run("IncludesAliases", func(t *testing.T) {
		cmd := &cobra.Command{Use: "root"}
		cmd.AddCommand(&cobra.Command{Use: "version", Aliases: []string{"v"}})

		got := GetSubcommandMap(cmd)

		assert.True(t, got["version"])
		assert.True(t, got["v"])
	})

	t.Run("DoesNotIncludeUnknownCommands", func(t *testing.T) {
		cmd := &cobra.Command{Use: "root"}

		got := GetSubcommandMap(cmd)

		assert.False(t, got["unknown"])
	})
}

func TestNewArgSplitter(t *testing.T) {
	subcommands := map[string]bool{"version": true, "help": true, "__complete": true, "__completeNoDesc": true}

	tests := []struct {
		name        string
		args        []string
		wantOTFArgs []string
		wantTFArgs  []string
	}{
		{
			name:        "no args",
			args:        []string{"otf"},
			wantOTFArgs: nil,
			wantTFArgs:  nil,
		},
		{
			name:        "env and command only",
			args:        []string{"otf", "prod", "plan"},
			wantOTFArgs: []string{"prod", "plan"},
			wantTFArgs:  []string{},
		},
		{
			name:        "terraform destroy flag passes through",
			args:        []string{"otf", "prod", "apply", "-destroy"},
			wantOTFArgs: []string{"prod", "apply"},
			wantTFArgs:  []string{"-destroy"},
		},
		{
			name:        "terraform auto-approve flag passes through",
			args:        []string{"otf", "prod", "apply", "-auto-approve"},
			wantOTFArgs: []string{"prod", "apply"},
			wantTFArgs:  []string{"-auto-approve"},
		},
		{
			name:        "multiple terraform flags pass through",
			args:        []string{"otf", "prod", "apply", "-destroy", "-auto-approve"},
			wantOTFArgs: []string{"prod", "apply"},
			wantTFArgs:  []string{"-destroy", "-auto-approve"},
		},
		{
			name:        "debug flag with terraform args",
			args:        []string{"otf", "--debug", "prod", "apply", "-destroy"},
			wantOTFArgs: []string{"--debug", "prod", "apply"},
			wantTFArgs:  []string{"-destroy"},
		},
		{
			name:        "short debug flag with terraform args",
			args:        []string{"otf", "-d", "prod", "apply", "-destroy"},
			wantOTFArgs: []string{"-d", "prod", "apply"},
			wantTFArgs:  []string{"-destroy"},
		},
		{
			name:        "verbose flag",
			args:        []string{"otf", "-v", "staging", "plan"},
			wantOTFArgs: []string{"-v", "staging", "plan"},
			wantTFArgs:  []string{},
		},
		{
			name:        "multiple otf flags",
			args:        []string{"otf", "--debug", "--verbose", "prod", "apply", "-destroy"},
			wantOTFArgs: []string{"--debug", "--verbose", "prod", "apply"},
			wantTFArgs:  []string{"-destroy"},
		},
		{
			name:        "subcommand routes entirely to cobra",
			args:        []string{"otf", "version"},
			wantOTFArgs: []string{"version"},
			wantTFArgs:  nil,
		},
		{
			name:        "subcommand with flag routes entirely to cobra",
			args:        []string{"otf", "--debug", "version"},
			wantOTFArgs: []string{"--debug", "version"},
			wantTFArgs:  nil,
		},
		{
			name:        "help flag only",
			args:        []string{"otf", "--help"},
			wantOTFArgs: []string{"--help"},
			wantTFArgs:  nil,
		},
		{
			name:        "env only without command",
			args:        []string{"otf", "prod"},
			wantOTFArgs: []string{"prod"},
			wantTFArgs:  nil,
		},
		{
			name:        "terraform var flag passes through",
			args:        []string{"otf", "prod", "plan", "-var", "foo=bar"},
			wantOTFArgs: []string{"prod", "plan"},
			wantTFArgs:  []string{"-var", "foo=bar"},
		},
		{
			name:        "__complete routes entirely to cobra",
			args:        []string{"otf", "__complete", "prod", ""},
			wantOTFArgs: []string{"__complete", "prod", ""},
			wantTFArgs:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			splitter := NewArgSplitter(subcommands)
			assert.Equal(t, tt.wantOTFArgs, splitter.OTFArgs)
			assert.Equal(t, tt.wantTFArgs, splitter.TFArgs)
		})
	}
}
