package starter

import "context"

// Ping returns a simple greeting. Replace with your business logic.
func Ping(_ context.Context) (string, error) {
	return "pong", nil
}
