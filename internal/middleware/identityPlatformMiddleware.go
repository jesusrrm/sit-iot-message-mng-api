package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"sit-iot-message-mng-api/config"
	"sit-iot-message-mng-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserEmailKey contextKey = "userEmail"
const TokenKey contextKey = "tokenKey"

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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		tokenStr := utils.ExtractBearerToken(authHeader)
		userID := utils.GetUserIDFromToken(tokenStr)
		userEmail := utils.GetUserEmailFromToken(tokenStr)

		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Add userID to the context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserEmailKey, userEmail)
		ctx = context.WithValue(ctx, TokenKey, tokenStr)

		c.Request = c.Request.WithContext(ctx)

		// Proceed to the next handler
		c.Next()
	}
}
