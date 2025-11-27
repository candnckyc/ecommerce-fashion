# Backend API

Go backend following 3-layer architecture:
- Handlers (HTTP layer)
- Services (Business logic)
- Repository (Data access)

## Setup

\`\`\`bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Start server
go run cmd/api/main.go
\`\`\`

## Project Structure

\`\`\`
backend/
├── cmd/api/          # Entry point
├── internal/         # Application code
│   ├── handlers/     # HTTP handlers
│   ├── services/     # Business logic
│   ├── repository/   # Database access
│   ├── models/       # Data structures
│   ├── middleware/   # HTTP middleware
│   └── utils/        # Helpers
├── migrations/       # Database migrations
└── tests/           # Tests
\`\`\`
