<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# SIT IoT Message Management API

This is a Go-based REST API for managing IoT messages with the following architecture:

## Project Structure
- `main.go` - Entry point of the application
- `config/` - Configuration management
- `database/` - Database and Firebase initialization
- `internal/`
  - `controllers/` - HTTP request handlers
  - `services/` - Business logic layer
  - `repositories/` - Data access layer
  - `models/` - Data models and structs
  - `middleware/` - HTTP middleware (authentication, etc.)
  - `routes/` - Route definitions
  - `utils/` - Utility functions

## Key Features
- Firebase Authentication integration
- MongoDB database with proper ObjectID handling
- RESTful API with CRUD operations for IoT messages
- React Admin compatible endpoints with pagination, sorting, and filtering
- User-based authorization and context propagation
- CORS support for web frontends

## Coding Guidelines
- Follow the repository pattern for data access
- Use context.Context to propagate user information across layers
- Implement proper error handling and logging
- Use MongoDB ObjectIDs for document IDs
- Support React Admin query parameters (filter, range, sort)
- Include proper CORS headers for frontend integration
- Validate user permissions at service layer

## Dependencies
- Gin web framework for HTTP routing
- MongoDB Go driver for database operations
- Firebase Admin SDK for authentication
- CORS middleware for cross-origin requests

## Environment Variables
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - MongoDB connection string
- `FIREBASE_CREDENTIALS_PATH` - Path to Firebase service account credentials
- `AUTH_API_KEY` - Firebase Auth API key
- `AUDIENCE` - Firebase project audience
