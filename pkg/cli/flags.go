package cli

const (
	// DebugFlagName is the flag name for debug log level.
	DebugFlagName = "debug"

	// VerboseFlagName is the flag name for verbose log level.
	VerboseFlagName = "verbose"
)

var (
	// DebugFlag increases log level verbosity to Debug.
	DebugFlag bool

	// VerboseFlag increases log level verbosity to Info.
	VerboseFlag bool
)
