package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	type Test struct {
		Name      string
		Env       string
		Command   string
		ExtraArgs []string
		Want      []string
	}

	var tests = []Test{
		{
			Name:    "InjectVarFileForPlan",
			Env:     "prod",
			Command: "plan",
			Want:    []string{"plan", "-var-file", "variables/prod.tfvars"},
		},
		{
			Name:      "InjectVarFileForApplyWithExtraArgs",
			Env:       "staging",
			Command:   "apply",
			ExtraArgs: []string{"-target", "module.db", "-auto-approve"},
			Want:      []string{"apply", "-var-file", "variables/staging.tfvars", "-target", "module.db", "-auto-approve"},
		},
		{
			Name:      "DoNotInjectVarFileForState",
			Env:       "prod",
			Command:   "state",
			ExtraArgs: []string{"list"},
			Want:      []string{"state", "list"},
		},
		{
			Name:    "DoNotInjectVarFileForFmt",
			Env:     "prod",
			Command: "fmt",
			Want:    []string{"fmt"},
		},
		{
			Name:      "InjectVarFileForDestroyWithExtraArgs",
			Env:       "prod",
			Command:   "destroy",
			ExtraArgs: []string{"-auto-approve"},
			Want:      []string{"destroy", "-var-file", "variables/prod.tfvars", "-auto-approve"},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var got = BuildArgs(test.Env, test.Command, test.ExtraArgs)

			assert.Equal(t, test.Want, got, `want: '%v', got: '%v'`, test.Want, got)
		})
	}
}
