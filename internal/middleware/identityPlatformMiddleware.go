package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"sit-iot-message-mng-api/config"

	"github.com/gin-gonic/gin"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserEmailKey contextKey = "userEmail"

func IdentityPlatformMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract the token from the header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			log.Println("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Validate the token using Identity Platform's REST API
		apiURL := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:lookup?key=%s", cfg.AuthApiKey)
		req, err := http.NewRequestWithContext(context.Background(), "POST", apiURL, strings.NewReader(fmt.Sprintf(`{"idToken":"%s"}`, token)))
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Token validation failed with status code: %d, response: %s", resp.StatusCode, string(body))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Parse response to get user info
		var tokenResponse struct {
			Users []struct {
				LocalID string `json:"localId"`
				Email   string `json:"email"`
			} `json:"users"`
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			c.Abort()
			return
		}

		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			log.Printf("Failed to parse response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			c.Abort()
			return
		}

		if len(tokenResponse.Users) == 0 {
			log.Println("No user found in token response")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userID := tokenResponse.Users[0].LocalID
		userEmail := tokenResponse.Users[0].Email

		if userID == "" {
			log.Println("User ID not found in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.Printf("User authenticated: ID=%s, Email=%s", userID, userEmail)

		// Add userID and userEmail to the context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserEmailKey, userEmail)

		c.Request = c.Request.WithContext(ctx)
		// Proceed to the next handler
		c.Next()
	}
}
