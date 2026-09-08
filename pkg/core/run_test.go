package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNeedsVarFile(t *testing.T) {
	t.Run("returns true for commands that accept var-file", func(t *testing.T) {
		commands := []string{"plan", "apply", "destroy", "refresh", "import", "console"}

		for _, cmd := range commands {
			assert.True(t, NeedsVarFile(cmd), "expected NeedsVarFile(%q) to be true", cmd)
		}
	})

	t.Run("returns false for commands that do not accept var-file", func(t *testing.T) {
		commands := []string{"fmt", "validate", "state", "output", "init", "show", "graph", "providers", "taint", "untaint"}

		for _, cmd := range commands {
			assert.False(t, NeedsVarFile(cmd), "expected NeedsVarFile(%q) to be false", cmd)
		}
	})
}

func TestBuildArgs(t *testing.T) {
	t.Run("injects var-file for plan", func(t *testing.T) {
		args := BuildArgs("prod", "plan", nil)
		assert.Equal(t, []string{"plan", "-var-file", "variables/prod.tfvars"}, args)
	})

	t.Run("injects var-file for apply with extra args", func(t *testing.T) {
		args := BuildArgs("staging", "apply", []string{"-target", "module.db", "-auto-approve"})
		assert.Equal(t, []string{"apply", "-var-file", "variables/staging.tfvars", "-target", "module.db", "-auto-approve"}, args)
	})

	t.Run("does not inject var-file for state", func(t *testing.T) {
		args := BuildArgs("prod", "state", []string{"list"})
		assert.Equal(t, []string{"state", "list"}, args)
	})

	t.Run("does not inject var-file for fmt", func(t *testing.T) {
		args := BuildArgs("prod", "fmt", nil)
		assert.Equal(t, []string{"fmt"}, args)
	})

	t.Run("injects var-file for destroy with extra args", func(t *testing.T) {
		args := BuildArgs("prod", "destroy", []string{"-auto-approve"})
		assert.Equal(t, []string{"destroy", "-var-file", "variables/prod.tfvars", "-auto-approve"}, args)
	})
}
