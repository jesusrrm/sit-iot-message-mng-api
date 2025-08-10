package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

import "golang.org/x/crypto/bcrypt"

// GetUserIDFromToken extracts the UID from a verified auth token in Gin context.
func GetUserIDFromToken(token string) string {
    if strings.TrimSpace(token) == "" {
        log.Println("Token is empty")
        return ""
    }
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        log.Println("Invalid token format")
        return ""
    }
    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        log.Println("Failed to decode token payload:", err)
        return ""
    }
    var claims map[string]interface{}
    if err := json.Unmarshal(payload, &claims); err != nil {
        log.Println("Failed to unmarshal token payload:", err)
        return ""
    }
    if uid, ok := claims["user_id"].(string); ok {
        return uid
    }
    if sub, ok := claims["sub"].(string); ok {
        return sub // fallback to "sub" claim if "user_id" is not present
    }
    log.Println("UID not found in token claims")
    return ""
}

// GetUserEmailFromToken extracts the email from a verified auth token in Gin context.
func GetUserEmailFromToken(token string) string {
    if strings.TrimSpace(token) == "" {
        log.Println("Token is empty")
        return ""
    }
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        log.Println("Invalid token format")
        return ""
    }
    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        log.Println("Failed to decode token payload:", err)
        return ""
    }
    var claims map[string]interface{}
    if err := json.Unmarshal(payload, &claims); err != nil {
        log.Println("Failed to unmarshal token payload:", err)
        return ""
    }
    if email, ok := claims["email"].(string); ok {
        return email
    }

    log.Println("Email not found in token claims")
    return ""
}

func ParseJSON(input string, dest interface{}) error {
    return json.Unmarshal([]byte(input), dest)
}

func ExtractBearerToken(authHeader string) string {
    parts := strings.Split(authHeader, " ")
    if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
        return parts[1]
    }
    return ""
}

// HashPassword generates a bcrypt hash for the given password, ensuring $2a$ prefix.
func HashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Ensure the hash uses $2a$ prefix for compatibility (e.g., with VerneMQ)
	// The Go library usually produces $2a$, but this handles cases where it might produce $2b$
	// or if a different library version behaves differently.
	hashedPassword := string(hashedPasswordBytes)
	if strings.HasPrefix(hashedPassword, "$2b$") {
		hashedPassword = strings.Replace(hashedPassword, "$2b$", "$2a$", 1)
	}
	return hashedPassword, nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
// It expects the hash to potentially have a $2a$ or $2b$ prefix.
func CheckPasswordHash(password, hash string) bool {
	// The bcrypt.CompareHashAndPassword function in Go's library
	// can typically handle both $2a$ and $2b$ prefixes correctly.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}