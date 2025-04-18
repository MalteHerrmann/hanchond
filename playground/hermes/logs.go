package hermes

import (
	"os"
	"strings"
)

func getErrorFromHermesLogs(logsFile string) string {
	bz, err := os.ReadFile(logsFile)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(bz), "\n")
	foundErrors := make([]string, 0, len(lines)) // TODO: check if reallocating per new found error or pre-allocating to much space is worse for performance; doesn't really matter though
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "error") {
			foundErrors = append(foundErrors, line)
		}
	}

	return strings.Join(foundErrors, "\n")
}
