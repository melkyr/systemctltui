// Modify internal/system/commands.go
package system

import (
	"bytes"
	"os/exec"
	//"strings"
    "fmt" // Needed for fmt.Errorf

	tea "github.com/charmbracelet/bubbletea"
	"systemctltui/internal/messages" // <--- Import the new messages package
)

// ExecuteCommandAsync runs a systemctl command asynchronously and sends a CommandFinishedMsg
func ExecuteCommandAsync(command string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(command, args...)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run() // This is a blocking call

		output := stdout.String()
		if stderrStr := stderr.String(); stderrStr != "" {
			// Append stderr unless it's just a status message on some systems
            // Simple heuristic: if stderr contains "Active:" or similar status, maybe don't append?
            // For robustness, let's append it always or make it conditional based on Err
            if err != nil || stderrStr != "" { // Append if there was an error or stderr is not empty
               if output != "" {
                  output += "\n--- STDERR ---\n" // Separator
               }
               output += stderrStr
            }
		}


		// Send the result back to the main update loop
		return messages.CommandFinishedMsg{Output: output, Err: err} // <--- Use messages.CommandFinishedMsg
	}
}

// Helper function to construct the systemctl command for common actions
func SystemctlCommand(command string, unit string) tea.Cmd {
    // Basic validation
    if command == "" {
        return func() tea.Msg {
             return messages.CommandFinishedMsg{Output: "", Err: fmt.Errorf("no command specified")} // <--- Use messages.CommandFinishedMsg
        }
    }

    args := []string{command}

    // Logic to add the unit argument based on command
    unitRequired := true // Assume most commands need a unit by default
    switch command {
        case "list-units", "list-timers", "list-sockets", "list-jobs", "list-dependencies", "list-unit-files":
            unitRequired = false // These commands don't take a unit as the *first* argument like start/stop
        case "status", "is-active", "is-enabled", "is-failed":
             // These can take a unit, and are much more useful with one,
             // but 'status' without a unit shows overall system status.
             // Let's make the unit optional for status, is-active, etc.
             unitRequired = false // Unit is optional
             if unit != "" {
                 args = append(args, unit)
             }
         default:
             // start, stop, restart, enable, disable, mask, unmask, kill, cat, edit, show etc.
             // These typically require a unit.
             unitRequired = true
             if unit != "" {
                 args = append(args, unit)
             }
    }


    if unitRequired && unit == "" {
         // If the command *requires* a unit but none was provided
          return func() tea.Msg {
             return messages.CommandFinishedMsg{Output: "", Err: fmt.Errorf("command '%s' requires a unit", command)} // <--- Use messages.CommandFinishedMsg
        }
    }


	return ExecuteCommandAsync("systemctl", args...)
}
