package config

import (
	"os"
)

type Config struct {
	Port                    string
	DatabaseURL             string
	FirebaseCredentialsPath string
	AuthApiKey              string
	Audience                string
	ProjectServiceApiUrl    string
	DBName   				string

}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:                    getEnv("PORT", "8080"),
		DatabaseURL:             getEnv("DB_URI_MESSAGE_MNG", "mongodb://localhost:27017/sit-iot-message-mng"),
		FirebaseCredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		AuthApiKey:              getEnv("AUTH_API_KEY", ""),
		Audience:                getEnv("AUDIENCE", ""),
		ProjectServiceApiUrl:    getEnv("PROJECT_SERVICE_API_URL", "http://localhost"),
		DBName:   				 getEnv("DB_NAME", "sit-iot-message-mng"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
