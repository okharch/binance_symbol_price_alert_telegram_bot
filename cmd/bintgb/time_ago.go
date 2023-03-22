package main

import (
	"fmt"
	"time"
)

func timeAgo(t time.Time) string {
	duration := time.Since(t)
	switch {
	case duration.Seconds() < 60:
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	case duration.Minutes() < 60:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration.Hours() < 48:
		return "yesterday"
	default:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}
