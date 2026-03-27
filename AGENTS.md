# AGENTS.md - Developer Guidelines for Freelance Elite

This document provides guidelines and instructions for agentic coding agents working on this codebase.

## Project Overview

Freelance Elite is a Go-based RESTful API for managing freelance profiles, user authentication, and geographic data. It uses:
- **Go 1.23** - Language
- **Echo v4** - Web framework
- **GORM** - ORM for database
- **MySQL** - Database
- **testify** - Testing framework

---

## Build, Lint, and Test Commands

### Running the Application

```bash
# Development server (runs on http://localhost:1323)
go run main.go

# Build binary
go build

# Run built binary
./freelance_elite
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./handlers

# Run tests with coverage
go test -cover ./...

# Run a specific test function
go test -v -run TestRegisterSuccess ./handlers

# Run tests for a specific package
go test -v ./models
go test -v ./handlers
```

### Database Commands

```bash
# Run migrations
go run cmd/automigrate/main.go

# Seed database (compile all seed files)
go run cmd/seed/*.go
```

---

## Code Style Guidelines

### Formatting

- Use `gofmt` (automatic with most editors) - 4-space indentation
- Keep lines under 100 characters when practical
- Add blank lines between functions and between import groups
- Use vertical whitespace to group related logic

### Imports

Imports should be organized in three groups with blank lines between:

1. Standard library packages
2. External/third-party packages
3. Internal/local packages

```go
import (
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "freelance_elite/db"
    "freelance_elite/models"
)
```

### Naming Conventions

- **Files**: Use snake_case (e.g., `auth_handler.go`, `user_model.go`)
- **Types/Structs**: Use PascalCase (e.g., `User`, `ProfileHandler`)
- **Variables/Functions**: Use camelCase (e.g., `userID`, `Register`)
- **Constants**: Use PascalCase or snake_case with all caps for exported (e.g., `MaxAge`, `DEFAULT_LIMIT`)
- **Packages**: Use short, lowercase names (e.g., `handlers`, `models`, `config`)
- **Database tables**: Use snake_case, plural (e.g., `users`, `profiles`)

### Types and Declarations

- Prefer explicit type declarations over type inference for clarity
- Use struct tags for JSON serialization and GORM mapping
- Group related struct fields together

```go
type User struct {
    ID        int            `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Username  string         `json:"username"`
    Email     string         `json:"email" gorm:"unique"`
    Password  string         `json:"-"`
}
```

### Error Handling

- Return errors as JSON responses with appropriate HTTP status codes
- Check errors immediately after operations
- Provide meaningful error messages
- Handle database errors with `gorm.ErrRecordNotFound` specifically
- Use string matching for SQL-specific errors (e.g., "Duplicate entry")

```go
if err := db.DB.Create(&user).Error; err != nil {
    if strings.Contains(err.Error(), "Duplicate entry") {
        return c.JSON(http.StatusConflict, map[string]string{"error": "Username or email already exists"})
    }
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user: " + err.Error()})
}
```

### Middleware and Authentication

- Use middleware for cross-cutting concerns (logging, recovery, auth)
- Check for authorization header and handle missing tokens appropriately
- Return proper HTTP status codes (401 for unauthorized, 403 for forbidden)

### Database Operations

- Use GORM methods: `Create`, `First`, `Where`, `Save`, `Delete`
- Always check for `gorm.ErrRecordNotFound` when querying single records
- Use struct tags for column definitions, indexes, and constraints

### Response Format

Return JSON responses consistently:

```go
// Success response
return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})

// Error response
return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})

// Data response
return c.JSON(http.StatusOK, user)
```

### Comments

- Add comments for exported functions explaining purpose and parameters
- Use doc comments (starting with function name) for public APIs
- Keep comments concise and descriptive
- Avoid obvious comments (e.g., `// Increment i`)

### Testing Conventions

- Use testify suite for organized tests
- Follow naming pattern: `*_test.go` for test files
- Test structure:
  - `SetupSuite()` - runs once before all tests (setup DB, routes)
  - `TearDownSuite()` - runs once after all tests (cleanup)
  - `SetupTest()` - runs before each test (clean tables)
  - Individual test functions: `TestXxx()`

```go
type AuthTestSuite struct {
    suite.Suite
    e *echo.Echo
}

func (s *AuthTestSuite) SetupSuite() {
    // Setup test database
    db.InitDB(...)
    s.e = echo.New()
    s.e.POST("/register", Register)
}

func (s *AuthTestSuite) TestRegisterSuccess() {
    // Test implementation
    assert.Equal(s.T(), http.StatusCreated, rec.Code)
}
```

---

## Project Structure

```
freelance_elite/
├── cmd/
│   ├── automigrate/      # Database migration tool
│   └── seed/            # Database seeding tool
├── config/              # Configuration files
│   ├── database.go
│   └── routes.go
├── db/                  # Database connection
│   └── database.go
├── handlers/            # HTTP request handlers
│   ├── auth.go
│   ├── profiles.go
│   ├── gender.go
│   └── countries.go
├── models/              # Data models
│   ├── user.go
│   ├── profile.go
│   ├── gender.go
│   └── countries.go
└── main.go              # Entry point
```

---

## Environment Variables

Required variables in `.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=freelance_elite_dev
JWT_SECRET=your_jwt_secret
APP_ENV=development
```

For testing:
```env
TEST_DB_USER=your_mysql_user
TEST_DB_PASSWORD=your_mysql_password
TEST_DB_NAME=freelance_elite_test
TEST_JWT_SECRET=your_test_jwt_secret
```

---

## Key Conventions Summary

1. **Always** run tests before committing changes
2. **Use** testify suite pattern for organized tests
3. **Handle** errors explicitly with proper HTTP status codes
4. **Organize** imports in three groups (stdlib, external, internal)
5. **Follow** Go naming conventions (PascalCase for exports, camelCase for locals)
6. **Use** struct tags for JSON/GORM configuration
7. **Return** consistent JSON response format
8. **Never** commit secrets or credentials to version control
