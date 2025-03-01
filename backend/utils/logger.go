package utils

import "log"

// Info logs informational messages.
func Info(message string) {
	log.Println("[INFO]", message)
}

// Error logs error messages.
func Error(message string) {
	log.Println("[ERROR]", message)
}

// Warn logs warning messages.
func Warn(message string) {
	log.Println("[WARN]", message)
}
