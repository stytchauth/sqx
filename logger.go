package sqx

var defaultLogger Logger = nil

// SetDefaultLogger sets the logger that should be used to log information.
// If you need to change the logger for a specific request, use WithLogger
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

type Logger interface {
	// Printf prints output using the provided logger
	// Arguments are passed in the style of fmt.Printf.
	Printf(format string, v ...any)
}
