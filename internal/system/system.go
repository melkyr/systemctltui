// package system
package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Unit represents a single systemd unit fetched from systemctl, with more details.
type Unit struct {
	Name        string
	Load        string
	Active      string
	Sub         string
	Description string
	Type        string // e.g., "service", "device", "mount"
}

// FetchUnits calls 'systemctl list-units' and parses the output into structured Unit data.
func FetchUnits() ([]Unit, error) {
	// Use --no-legend to get just the data rows
	// Use --plain to ensure consistent space separation
	out, err := exec.Command("systemctl", "list-units", "--all", "--no-legend", "--plain").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute systemctl: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	units := make([]Unit, 0, len(lines))

	for _, line := range lines {
		// Use strings.FieldsFunc to split by one or more spaces
		fields := strings.Fields(line)

		if len(fields) >= 5 { // Expect at least 5 fields: UNIT LOAD ACTIVE SUB DESCRIPTION
			unitName := fields[0]

			// Attempt to extract unit type from the name suffix (e.g., ".service")
			unitType := "unknown"
			if lastDot := strings.LastIndex(unitName, "."); lastDot != -1 && lastDot < len(unitName)-1 {
				unitType = unitName[lastDot+1:]
			}

			// The description might contain spaces, so join the remaining fields
			description := strings.Join(fields[4:], " ")


			unit := Unit{
				Name:        unitName,
				Load:        fields[1],
				Active:      fields[2],
				Sub:         fields[3],
				Description: description,
				Type:        unitType,
			}
			units = append(units, unit)
		} else if len(fields) > 0 {
             // Handle lines with fewer than expected fields (shouldn't happen with --plain?)
             units = append(units, Unit{Name: fields[0], Description: "Error parsing unit data"})
        }
	}
	return units, nil
}

// TODO: Add other systemctl commands here (e.g., Status, Start, Stop)
// func Status(unitName string) (string, error) { ... }
