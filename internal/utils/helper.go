package utils

import (
	"encoding/json"
	"strings"
)

// ExtractBearerToken extracts the Bearer token from the Authorization header
func ExtractBearerToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

// ParseJSON parses a JSON string into the provided destination variable.
// Example: ParseJSON(`[0,9]`, &arr)
func ParseJSON(input string, dest interface{}) error {
	return json.Unmarshal([]byte(input), dest)
}
