package fileutil

// Args holds the command-line arguments parsed from the CLI.
type Args struct {
	// ComparisonCriteria is the field used as the map key for matching files
	// across both lists. Accepted values: "filename" (default) or "path".
	ComparisonCriteria string
}

// Conf holds all runtime configuration for the comparison run.
type Conf struct {
	// Arg contains the parsed CLI arguments.
	Arg Args

	// TrimIndexPathA is the number of leading directory segments to strip
	// from paths in list_A before comparison.
	TrimIndexPathA int

	// TrimIndexPathB is the number of leading directory segments to strip
	// from paths in list_B before comparison.
	TrimIndexPathB int

	// TrimPath indicates whether path trimming is active for this run.
	TrimPath bool

	// CommonDirSep is the path separator detected from the first list
	// ("/" on Unix, "\\" on Windows). Used when reconstructing trimmed paths.
	CommonDirSep string
}

// Config is the single shared configuration instance.
// It must be fully populated before any goroutine reads it.
var Config Conf
