package sqx

var defaultLogger Logger = nil

// SetDefaultLogger sets the logger that should be used to log information.
// If you need to change the logger for a specific request, use WithLogger
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

type logger struct {
	printf func(format string, v ...any)
}

// Printf forwards the format string and args to the logger's saved printf function given in the MakeLogger call.
func (l logger) Printf(format string, v ...any) {
	l.printf(format, v...)
}

// MakeLogger creates a new Logger interface by using the provided function as the logger's printf function.
func MakeLogger(printf func(format string, v ...any)) Logger {
	return logger{printf: printf}
}

// Logger is a simple interface that can be used to log events in the sqx package. It contains a single Printf function
// that takes a format string and arguments, much like fmt.Printf.
type Logger interface {
	// Printf prints output using the provided logger
	// Arguments are passed in the style of fmt.Printf.
	Printf(format string, v ...any)
}
