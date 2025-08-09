# SIT IoT Message Management API

A REST API service for managing IoT messages, built with Go, MongoDB, and Firebase Authentication.

## Features

- **Message Management**: CRUD operations for IoT messages with support for different message types (telemetry, commands, events, alerts)
- **Authentication**: Firebase Identity Platform integration for secure API access
- **Authorization**: User-based access control with context propagation
- **React Admin Compatible**: Endpoints support pagination, sorting, and filtering for React Admin frontend
- **Project & Device Filtering**: Messages can be filtered by project and device IDs
- **CORS Support**: Configured for web frontend integration

## API Endpoints

### Messages
- `POST /api/message` - Create a new message
- `GET /api/message/:id` - Get message by ID
- `PUT /api/message/:id` - Update message
- `DELETE /api/message/:id` - Delete message
- `GET /api/message` - List messages with pagination and filtering

### Project-specific Messages
- `GET /api/project/:projectId/message` - List messages for a specific project

### Device-specific Messages
- `GET /api/device/:deviceId/message` - List messages for a specific device

## Data Models

### Message
- `id` - Unique identifier (MongoDB ObjectID)
- `projectId` - Associated project ID
- `deviceId` - Source device ID
- `type` - Message type (telemetry, command, event, alert)
- `status` - Message status (pending, delivered, failed, processed)
- `payload` - Message content (flexible JSON)
- `metadata` - Additional key-value pairs
- `timestamp` - Message timestamp
- `createdAt` - Creation timestamp
- `updatedAt` - Last update timestamp
- `createdBy` - User who created the message

## Environment Variables

Create a `.env` file or set the following environment variables:

```bash
PORT=8080
DATABASE_URL=mongodb://localhost:27017/sit_iot_message_mng
FIREBASE_CREDENTIALS_PATH=/path/to/firebase-credentials.json
AUTH_API_KEY=your_firebase_auth_api_key
AUDIENCE=your_firebase_project_id.firebaseapp.com
```

## Development Setup

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Start MongoDB**
   Make sure MongoDB is running on your system.

3. **Configure Firebase**
   - Create a Firebase project
   - Download the service account credentials JSON file
   - Set the `FIREBASE_CREDENTIALS_PATH` environment variable

4. **Run the Application**
   ```bash
   go run main.go
   ```

The server will start on the configured port (default: 8080).

## Project Structure

```
sit-iot-message-mng-api/
├── main.go                              # Application entry point
├── config/
│   └── config.go                        # Configuration management
├── database/
│   └── database.go                      # DB and Firebase initialization
└── internal/
    ├── controllers/
    │   └── message_controller.go        # HTTP request handlers
    ├── services/
    │   └── message_service.go           # Business logic
    ├── repositories/
    │   └── message_repository.go        # Data access layer
    ├── models/
    │   └── message.go                   # Data models
    ├── middleware/
    │   └── identityPlatformMiddleware.go # Authentication middleware
    ├── routes/
    │   └── routes.go                    # Route definitions
    └── utils/
        └── helper.go                    # Utility functions
```

## Authentication

The API uses Firebase Authentication with Bearer tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your_firebase_id_token>
```

## React Admin Integration

The API is designed to work with React Admin. Query parameters supported:

- `filter` - JSON object for filtering (e.g., `{"type":"telemetry"}`)
- `range` - Array for pagination (e.g., `[0,9]`)
- `sort` - Array for sorting (e.g., `["timestamp","DESC"]`)

The API returns the `Content-Range` header required by React Admin for pagination.

## License

MIT License
