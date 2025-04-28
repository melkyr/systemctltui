// package messages
package messages

//import "fmt" // Often useful for error messages

// CommandFinishedMsg is a custom message sent when a system command finishes.
type CommandFinishedMsg struct {
	Output string // The command's standard output and standard error
	Err    error  // Any error that occurred during execution
}

// Add other application-level messages here in the future if needed.
// type SomeOtherAppMsg struct { ... }