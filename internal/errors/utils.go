package errors

import "fmt"

func buildErrorMsg(msg string, err string) string {
	switch err {
	case FileErr:
		return fmt.Sprintf("File error - %s", msg)
	case ServiceErr:
		return fmt.Sprintf("Service error: %s", msg)
	case SystemErr:
		return fmt.Sprintf("System error: %s", msg)
	default:
		return fmt.Sprintf("Unexpected error: %s", msg)
	}
}
