package testPackageLogger

import (
	"log"
)

func LogInfo(message string) {
	log.Printf("Info %v", message)
}

func LogError(message string) {
	log.Printf("Error %v", message)
}

func LogWarning(message string) {
	log.Printf("Warning %v", message)
}
