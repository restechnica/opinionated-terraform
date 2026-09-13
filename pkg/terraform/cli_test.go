package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/restechnica/opinionated-terraform/pkg/cli"
)

func TestNeedsVarFile(t *testing.T) {
	type Test struct {
		Name    string
		Command string
		Want    bool
	}

	var tests = []Test{
		{Name: "PlanNeedsVarFile", Command: "plan", Want: true},
		{Name: "ApplyNeedsVarFile", Command: "apply", Want: true},
		{Name: "DestroyNeedsVarFile", Command: "destroy", Want: true},
		{Name: "RefreshNeedsVarFile", Command: "refresh", Want: true},
		{Name: "ImportNeedsVarFile", Command: "import", Want: true},
		{Name: "ConsoleNeedsVarFile", Command: "console", Want: true},
		{Name: "FmtDoesNotNeedVarFile", Command: "fmt", Want: false},
		{Name: "ValidateDoesNotNeedVarFile", Command: "validate", Want: false},
		{Name: "StateDoesNotNeedVarFile", Command: "state", Want: false},
		{Name: "OutputDoesNotNeedVarFile", Command: "output", Want: false},
		{Name: "InitDoesNotNeedVarFile", Command: "init", Want: false},
		{Name: "ShowDoesNotNeedVarFile", Command: "show", Want: false},
		{Name: "GraphDoesNotNeedVarFile", Command: "graph", Want: false},
		{Name: "ProvidersDoesNotNeedVarFile", Command: "providers", Want: false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var got = NeedsVarFile(test.Command)

			assert.Equal(t, test.Want, got, `want: '%v', got: '%v'`, test.Want, got)
		})
	}
}

func TestBuildArgs(t *testing.T) {
	t.Run("InjectVarFileForPlanWhenFileExists", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "prod.tfvars"), []byte(""), 0644))

		got, err := BuildArgs("prod", "plan", nil)
		require.NoError(t, err)
		var want = []string{"plan", "-var-file", "variables/prod.tfvars"}

		assert.Equal(t, want, got, `want: '%v', got: '%v'`, want, got)
	})

	t.Run("InjectVarFileForApplyWithExtraArgs", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "staging.tfvars"), []byte(""), 0644))

		got, err := BuildArgs("staging", "apply", []string{"-target", "module.db", "-auto-approve"})
		require.NoError(t, err)
		var want = []string{"apply", "-var-file", "variables/staging.tfvars", "-target", "module.db", "-auto-approve"}

		assert.Equal(t, want, got, `want: '%v', got: '%v'`, want, got)
	})

	t.Run("OmitVarFileWhenVariablesDirMissing", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))

		got, err := BuildArgs("prod", "plan", nil)
		require.NoError(t, err)
		var want = []string{"plan"}

		assert.Equal(t, want, got, `want: '%v', got: '%v'`, want, got)
	})

	t.Run("ErrorWhenVariablesDirExistsButFileIsMissing", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))

		_, err := BuildArgs("stagin", "apply", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stagin")
	})

	t.Run("DoNotInjectVarFileForState", func(t *testing.T) {
		got, err := BuildArgs("prod", "state", []string{"list"})
		require.NoError(t, err)
		var want = []string{"state", "list"}

		assert.Equal(t, want, got, `want: '%v', got: '%v'`, want, got)
	})

	t.Run("DoNotInjectVarFileForFmt", func(t *testing.T) {
		got, err := BuildArgs("prod", "fmt", nil)
		require.NoError(t, err)
		var want = []string{"fmt"}

		assert.Equal(t, want, got, `want: '%v', got: '%v'`, want, got)
	})

	t.Run("InjectVarFileForDestroyWithExtraArgs", func(t *testing.T) {
		var dir = t.TempDir()
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.MkdirAll(cli.DefaultVariablesDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(cli.DefaultVariablesDir, "prod.tfvars"), []byte(""), 0644))

		got, err := BuildArgs("prod", "destroy", []string{"-auto-approve"})
		require.NoError(t, err)
		var want = []string{"destroy", "-var-file", "variables/prod.tfvars", "-auto-approve"}

		assert.Equal(t, want, got, `want: '%v', got: '%v'`, want, got)
	})
}
