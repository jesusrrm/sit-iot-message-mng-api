package repositories

import (
	"sit-iot-message-mng-api/config"
	"testing"
)

func TestRepositoryFactory(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		wantErr  bool
	}{
		{
			name:     "MongoDB provider",
			provider: "mongo",
			wantErr:  false,
		},
		{
			name:     "MongoDB provider (alternative)",
			provider: "mongodb",
			wantErr:  false,
		},
		{
			name:     "Firestore provider",
			provider: "firestore",
			wantErr:  false,
		},
		{
			name:     "Invalid provider",
			provider: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DatabaseProvider: tt.provider,
			}
			factory := NewRepositoryFactory(cfg)

			if factory.GetDatabaseProvider() != tt.provider {
				t.Errorf("GetDatabaseProvider() = %v, want %v", factory.GetDatabaseProvider(), tt.provider)
			}

			// Test repository creation (will fail without actual DB clients, but should validate provider)
			_, err := factory.CreateMessageRepository(nil, nil)

			if tt.wantErr && err == nil {
				t.Errorf("CreateMessageRepository() expected error but got none")
			}

			if !tt.wantErr && err == nil {
				t.Errorf("CreateMessageRepository() should fail without DB clients")
			}
		})
	}
}

func TestGetDatabaseProvider(t *testing.T) {
	cfg := &config.Config{
		DatabaseProvider: "mongo",
	}
	factory := NewRepositoryFactory(cfg)

	if got := factory.GetDatabaseProvider(); got != "mongo" {
		t.Errorf("GetDatabaseProvider() = %v, want %v", got, "mongo")
	}
}
