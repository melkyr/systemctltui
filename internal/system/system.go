// package system
package system

import (
	"fmt" // <--- Add this import
	"os/exec"
	"strings"
)

// Unit represents a single systemd unit fetched from systemctl.
type Unit struct {
	Name string
	// Add other fields if you parse them later (e.g., Load, Active, Sub, Description)
}

// FetchUnits calls 'systemctl list-units' and parses the output.
func FetchUnits() ([]Unit, error) {
	// Using --no-legend removes the header and footer
	out, err := exec.Command("systemctl", "list-units", "--all", "--no-legend").Output()
	if err != nil {
		// It's better to return the error and let the caller handle it
		return nil, fmt.Errorf("failed to execute systemctl: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	units := make([]Unit, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			// Assuming the first field is always the unit name
			units = append(units, Unit{Name: fields[0]})
			// TODO: Parse other fields if needed
		}
	}
	return units, nil
}

// TODO: Add other systemctl commands here (e.g., Status, Start, Stop)
// func Status(unitName string) (string, error) { ... }
