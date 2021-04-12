package google

import (
	"fmt"
	"log"
)

type dclLogger struct{}

// Fatal records Fatal errors.
func (l dclLogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Fatalf records Fatal errors with added arguments.
func (l dclLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(fmt.Sprintf("[DEBUG][DCL FATAL] %s", format), args...)
}

// Info records Info errors.
func (l dclLogger) Info(args ...interface{}) {
	log.Print(args...)
}

// Infof records Info errors with added arguments.
func (l dclLogger) Infof(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[DEBUG][DCL INFO] %s", format), args...)
}

// Warningf records Warning errors with added arguments.
func (l dclLogger) Warningf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[DEBUG][DCL WARNING] %s", format), args...)
}

// Warning records Warning errors.
func (l dclLogger) Warning(args ...interface{}) {
	log.Print(args...)
}
