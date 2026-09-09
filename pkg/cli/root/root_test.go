package root

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitArgs(t *testing.T) {
	subcommands := map[string]bool{"version": true, "help": true}

	tests := []struct {
		name        string
		args        []string
		wantOtfArgs []string
		wantTfArgs  []string
	}{
		{
			name:        "no args",
			args:        []string{},
			wantOtfArgs: nil,
			wantTfArgs:  nil,
		},
		{
			name:        "env and command only",
			args:        []string{"prod", "plan"},
			wantOtfArgs: []string{"prod", "plan"},
			wantTfArgs:  []string{},
		},
		{
			name:        "terraform destroy flag passes through",
			args:        []string{"prod", "apply", "-destroy"},
			wantOtfArgs: []string{"prod", "apply"},
			wantTfArgs:  []string{"-destroy"},
		},
		{
			name:        "terraform auto-approve flag passes through",
			args:        []string{"prod", "apply", "-auto-approve"},
			wantOtfArgs: []string{"prod", "apply"},
			wantTfArgs:  []string{"-auto-approve"},
		},
		{
			name:        "multiple terraform flags pass through",
			args:        []string{"prod", "apply", "-destroy", "-auto-approve"},
			wantOtfArgs: []string{"prod", "apply"},
			wantTfArgs:  []string{"-destroy", "-auto-approve"},
		},
		{
			name:        "debug flag with terraform args",
			args:        []string{"--debug", "prod", "apply", "-destroy"},
			wantOtfArgs: []string{"--debug", "prod", "apply"},
			wantTfArgs:  []string{"-destroy"},
		},
		{
			name:        "short debug flag with terraform args",
			args:        []string{"-d", "prod", "apply", "-destroy"},
			wantOtfArgs: []string{"-d", "prod", "apply"},
			wantTfArgs:  []string{"-destroy"},
		},
		{
			name:        "verbose flag",
			args:        []string{"-v", "staging", "plan"},
			wantOtfArgs: []string{"-v", "staging", "plan"},
			wantTfArgs:  []string{},
		},
		{
			name:        "multiple otf flags",
			args:        []string{"--debug", "--verbose", "prod", "apply", "-destroy"},
			wantOtfArgs: []string{"--debug", "--verbose", "prod", "apply"},
			wantTfArgs:  []string{"-destroy"},
		},
		{
			name:        "subcommand routes entirely to cobra",
			args:        []string{"version"},
			wantOtfArgs: []string{"version"},
			wantTfArgs:  nil,
		},
		{
			name:        "subcommand with flag routes entirely to cobra",
			args:        []string{"--debug", "version"},
			wantOtfArgs: []string{"--debug", "version"},
			wantTfArgs:  nil,
		},
		{
			name:        "help flag only",
			args:        []string{"--help"},
			wantOtfArgs: []string{"--help"},
			wantTfArgs:  nil,
		},
		{
			name:        "env only without command",
			args:        []string{"prod"},
			wantOtfArgs: []string{"prod"},
			wantTfArgs:  nil,
		},
		{
			name:        "terraform var flag passes through",
			args:        []string{"prod", "plan", "-var", "foo=bar"},
			wantOtfArgs: []string{"prod", "plan"},
			wantTfArgs:  []string{"-var", "foo=bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOtfArgs, gotTfArgs := splitArgs(tt.args, subcommands)
			assert.Equal(t, tt.wantOtfArgs, gotOtfArgs)
			assert.Equal(t, tt.wantTfArgs, gotTfArgs)
		})
	}
}
