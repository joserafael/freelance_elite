# Freelance Elite API

A comprehensive RESTful API built with Go for managing freelance profiles, user authentication, and related data. This project provides a robust backend solution for freelance platforms with user management, profile creation, and geographic data handling.

## 🚀 Features

### Authentication & Security
- **User Registration** with validation and password confirmation
- **JWT-based Authentication** with secure token generation
- **Token Blacklisting** for secure logout functionality
- **Protected Routes** with middleware authentication
- **Password Hashing** using bcrypt for security

### Profile Management
- **Complete Profile System** with personal information
- **Advanced Search & Filtering** by name, age, gender, and country
- **Age-based Filtering** with dynamic date calculations
- **Profile Relationships** with users, genders, and countries

### Geographic Data
- **Countries Management** with regions and subregions
- **Regional Filtering** and search capabilities
- **Country-based Profile Organization**

### Gender Management
- **Gender Categories** with active/inactive status
- **CRUD Operations** for gender management

## 🛠 Technologies Used

- **Go 1.25** - Primary programming language
- **Echo v4** - High-performance web framework
- **GORM** - ORM for database interactions
- **MySQL** - Primary database
- **JWT (golang-jwt/jwt/v5)** - JSON Web Token authentication
- **bcrypt** - Password hashing
- **godotenv** - Environment variable management
- **testify** - Testing framework

## 📋 Prerequisites

- Go 1.25 or higher
- MySQL 8.0 or higher
- Git

## ⚙️ Setup & Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd freelance_elite
```

### 2. Environment Configuration
Create a `.env` file in the project root:

```env
# Common Configuration
DB_HOST=localhost
DB_PORT=3306
JWT_SECRET=your_super_secret_jwt_key_here

# Development Environment
APP_ENV=development
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=freelance_elite_dev

# Test Environment
TEST_DB_USER=your_mysql_user
TEST_DB_PASSWORD=your_mysql_password
TEST_DB_NAME=freelance_elite_test
TEST_JWT_SECRET=your_test_jwt_secret_key
```

> **Security Note**: Replace placeholder values with strong, unique credentials. Never commit real credentials to version control.

### 3. Database Setup

#### Create Databases
```sql
CREATE DATABASE freelance_elite_dev;
CREATE DATABASE freelance_elite_test;
```

#### Run Migrations
```bash
go run cmd/automigrate/main.go
```

#### Seed Initial Data
```bash
go run cmd/seed/main.go
```

### 4. Install Dependencies
```bash
go mod download
```

## 🚀 Running the Application

### Development Server
```bash
go run main.go
```
The server will start on `http://localhost:1323`

### Alternative Methods
```bash
# Build and run
go build
./freelance_elite

# Run with go run
go run .
```

## 🧪 Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Package Tests
```bash
go test -v ./handlers
go test -v ./models
```

### Test with Coverage
```bash
go test -cover ./...
```

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "password_confirmation": "securepassword123"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

#### Logout
```http
POST /logout
Authorization: Bearer <jwt_token>
```

### Profile Endpoints

#### Get User Profile (Protected)
```http
GET /profile
Authorization: Bearer <jwt_token>
```

#### Create Profile
```http
POST /profiles
Content-Type: application/json

{
  "name": "John",
  "last_name": "Doe",
  "date_birth": "1990-01-15T00:00:00Z",
  "user_id": 1,
  "gender_id": 1,
  "country_id": 1
}
```

#### Get All Profiles
```http
GET /profiles?gender_id=1&country_id=1&search=john&min_age=25&max_age=35
```

#### Get Profile by ID
```http
GET /profiles/:id
```

#### Update Profile
```http
PUT /profiles/:id
Content-Type: application/json

{
  "name": "John Updated",
  "last_name": "Doe",
  "date_birth": "1990-01-15T00:00:00Z",
  "user_id": 1,
  "gender_id": 1,
  "country_id": 1
}
```

#### Delete Profile
```http
DELETE /profiles/:id
```

### Gender Endpoints

#### Get All Genders
```http
GET /genders
```

#### Get Gender by ID
```http
GET /genders/:id
```

#### Create Gender
```http
POST /genders
Content-Type: application/json

{
  "name": "Non-binary",
  "is_active": true
}
```

#### Update Gender
```http
PUT /genders/:id
Content-Type: application/json

{
  "name": "Updated Gender",
  "is_active": true
}
```

#### Delete Gender
```http
DELETE /genders/:id
```

### Country Endpoints

#### Get All Countries
```http
GET /countries?region=Europe&is_active=true&search=spain
```

#### Get Country by ID
```http
GET /countries/:id
```

#### Get Countries by Region
```http
GET /countries/region/:region
```

#### Create Country
```http
POST /countries
Content-Type: application/json

{
  "name": "Spain",
  "code": "ES",
  "region": "Europe",
  "subregion": "Southern Europe",
  "capital": "Madrid",
  "is_active": true
}
```

#### Update Country
```http
PUT /countries/:id
Content-Type: application/json

{
  "name": "Spain",
  "code": "ES",
  "region": "Europe",
  "subregion": "Southern Europe",
  "capital": "Madrid",
  "is_active": true
}
```

#### Delete Country
```http
DELETE /countries/:id
```

## 📁 Project Structure

```
freelance_elite/
├── cmd/                    # Command-line tools
│   ├── automigrate/       # Database migration tool
│   └── seed/              # Database seeding tool
├── config/                # Configuration files
│   ├── database.go        # Database configuration
│   └── routes.go          # Route definitions
├── db/                    # Database connection
│   └── database.go        # Database initialization
├── handlers/              # HTTP request handlers
│   ├── auth.go           # Authentication handlers
│   ├── profiles.go       # Profile management
│   ├── gender.go         # Gender management
│   ├── countries.go      # Country management
│   └── *_test.go         # Test files
├── models/               # Data models
│   ├── user.go          # User model
│   ├── profile.go       # Profile model
│   ├── gender.go        # Gender model
│   ├── countries.go     # Country model
│   └── blacklisted_token.go # Token blacklist
├── main.go              # Application entry point
├── go.mod               # Go module definition
└── README.md            # Project documentation
```

## 🔧 Development Commands

### Database Operations
```bash
# Run migrations
go run cmd/automigrate/main.go

# Seed database
go run cmd/seed/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./handlers

# Run tests with coverage
go test -cover ./...
```

### Building
```bash
# Build for current platform
go build

# Build for specific platform
GOOS=linux GOARCH=amd64 go build
```

## 🔒 Security Features

- **JWT Token Authentication** with configurable expiration
- **Token Blacklisting** for secure logout
- **Password Hashing** using bcrypt
- **Input Validation** and sanitization
- **SQL Injection Protection** via GORM
- **Environment-based Configuration** for sensitive data

## 🚀 Deployment

### Environment Variables for Production
```env
APP_ENV=production
DB_HOST=your_production_db_host
DB_PORT=3306
DB_USER=your_production_db_user
DB_PASSWORD=your_production_db_password
DB_NAME=freelance_elite_prod
JWT_SECRET=your_super_secure_production_jwt_secret
```

### Docker Support (Optional)
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o freelance_elite

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/freelance_elite .
CMD ["./freelance_elite"]
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

If you encounter any issues or have questions:

1. Check the existing issues on GitHub
2. Create a new issue with detailed information
3. Contact the development team

## 🔄 Changelog

### Version 1.0.0
- Initial release with authentication system
- Profile management functionality
- Geographic data management
- Comprehensive test suite
- Database migration and seeding tools
