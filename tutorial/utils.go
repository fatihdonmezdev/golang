package utils

import "strings"

func extractDateFromResHead(resHead string) string {
	lines := strings.Split(resHead, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Date:") {
			parts := strings.Split(line, "Date:")
			if len(parts) > 1 {
				// Trim any leading/trailing spaces from the date substring
				date := strings.TrimSpace(parts[1])
				return date
			}
		}
	}
	return "" // Return empty string if date not found
}
