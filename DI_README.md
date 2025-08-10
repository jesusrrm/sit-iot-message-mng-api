# Database Dependency Injection

This project now supports multiple database backends through a dependency injection system. You can choose between MongoDB and Firestore as your database provider.

## Configuration

Add the following environment variable to configure your database provider:

```bash
# Choose your database provider (default: mongo)
DATABASE_PROVIDER=mongo  # or "firestore"
```

### MongoDB Configuration (default)
```bash
DATABASE_PROVIDER=mongo
DB_URI_MESSAGE_MNG=mongodb://localhost:27017/sit-iot-message-mng
DB_NAME_MESSAGE_MNG=sit-iot-messages-mng
```

### Firestore Configuration
```bash
DATABASE_PROVIDER=firestore
FIREBASE_CREDENTIALS_PATH=/path/to/your/firebase-credentials.json
# Or set GOOGLE_APPLICATION_CREDENTIALS for default auth
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_PROVIDER` | Database provider to use (`mongo` or `firestore`) | `mongo` | No |
| `DB_URI_MESSAGE_MNG` | MongoDB connection URI | `mongodb://localhost:27017/sit-iot-message-mng` | Only for MongoDB |
| `DB_NAME_MESSAGE_MNG` | MongoDB database name | `sit-iot-messages-mng` | Only for MongoDB |
| `FIREBASE_CREDENTIALS_PATH` | Path to Firebase service account JSON | - | Only for Firestore |

## Switching Between Databases

To switch from MongoDB to Firestore:

1. Set the environment variable:
   ```bash
   export DATABASE_PROVIDER=firestore
   export FIREBASE_CREDENTIALS_PATH=/path/to/firebase-credentials.json
   ```

2. Restart your application

To switch back to MongoDB:

1. Set the environment variable:
   ```bash
   export DATABASE_PROVIDER=mongo
   export DB_URI_MESSAGE_MNG=mongodb://localhost:27017/sit-iot-message-mng
   ```

2. Restart your application

## Implementation Details

### Repository Pattern
The application uses the Repository pattern with a common interface (`MessageRepository`) that both MongoDB and Firestore implementations satisfy.

### Factory Pattern
A `RepositoryFactory` is used to create the appropriate repository instance based on the configuration.

### ID Handling
- **MongoDB**: Uses `primitive.ObjectID`
- **Firestore**: Uses string IDs
- The `Message` model now has an `interface{}` ID field with helper methods to handle both types

### Key Files
- `config/config.go` - Configuration with database provider selection
- `database/database.go` - Database initialization for both providers
- `internal/repositories/repository_factory.go` - Factory for creating repositories
- `internal/repositories/message_repository.go` - MongoDB implementation
- `internal/repositories/message_repository_firestore.go` - Firestore implementation
- `internal/models/message.go` - Updated model with flexible ID handling

## Dependencies

### For MongoDB support:
```go
go.mongodb.org/mongo-driver/mongo
go.mongodb.org/mongo-driver/bson
go.mongodb.org/mongo-driver/bson/primitive
```

### For Firestore support:
```go
cloud.google.com/go/firestore
firebase.google.com/go/v4
google.golang.org/api/option
```

Make sure to add these dependencies to your `go.mod`:

```bash
# For Firestore support
go get cloud.google.com/go/firestore
go get firebase.google.com/go/v4
```

## Testing

You can test with different database providers by setting the environment variable before running your tests:

```bash
# Test with MongoDB
DATABASE_PROVIDER=mongo go test ./...

# Test with Firestore
DATABASE_PROVIDER=firestore FIREBASE_CREDENTIALS_PATH=/path/to/creds.json go test ./...
```
