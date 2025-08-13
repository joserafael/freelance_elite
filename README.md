# Go Authentication Example

This project is a basic authentication example built with Go, demonstrating user registration, login, and secure logout with JWT (JSON Web Tokens) and a token blacklisting mechanism.

## Features
- User Registration
- User Login (JWT-based)
- User Logout (with JWT Blacklisting)
- Protected Routes

## Technologies Used
- **Go**: The primary programming language.
- **Echo**: A high-performance, minimalist Go web framework.
- **GORM**: An ORM (Object-Relational Mapping) library for Go, used for database interactions.
- **go-jwt/jwt/v5**: Go package for JSON Web Tokens.
- **bcrypt**: For secure password hashing.
- **godotenv**: For loading environment variables from `.env` files.
- **stretchr/testify**: A Go testing framework.

## Setup

### Prerequisites
- Go (version 1.18 or higher)
- MySQL database

### Environment Variables
Create a `.env` file in the project root with the following content (for development):
```
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=freelance_elite_dev
DB_HOST=localhost
DB_PORT=3306
JWT_SECRET=your_jwt_secret_key
```
Create a `.env.test` file in the project root with the following content (for testing):
```
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=freelance_elite_test
DB_HOST=localhost
DB_PORT=3306
JWT_SECRET=your_test_jwt_secret_key
```
*Remember to replace placeholder values with your actual database credentials and strong secret keys.*

### Database Setup
1. Ensure your MySQL server is running.
2. Run the auto-migration tool to create and update your database schema:
    ```bash
    go run cmd/automigrate/main.go
    ```

## Running the Application
```bash
go run main.go
```
The application will start on `http://localhost:1323`.

## Testing
To run the unit tests:
```bash
go test -v ./handlers
```

## API Endpoints
- `POST /register`: Register a new user.
    **Request Body:**
    ```json
    {
      "username": "john_doe",
      "email": "john.doe@example.com",
      "password": "securepassword123"
    }
    ```
- `POST /login`: Authenticate a user and receive a JWT.
    **Request Body:**
    ```json
    {
      "email": "john.doe@example.com",
      "password": "securepassword123"
    }
    ```
- `POST /logout`: Invalidate the current JWT (requires Authorization header).
- `GET /profile`: Access protected user profile information (requires valid JWT in Authorization header).

## Feedback

If you have any suggestions, encounter issues, or have ideas for improvements, please feel free to open an issue on the GitHub repository or contact me directly.
