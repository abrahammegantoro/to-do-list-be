# Todo List API

A RESTful API for managing todos built with Go Echo framework and PostgreSQL database.

## Features

- User authentication (register/login)
- CRUD operations for todos
- Category management
- JWT token-based authentication
- PostgreSQL database integration

## Tech Stack

- [Go](https://golang.org/) - Programming language
- [Echo](https://echo.labstack.com/) - Web framework
- [PostgreSQL](https://www.postgresql.org/) - Database
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT authentication

## API Endpoints

All endpoints are prefixed with `/api/v1`

### Authentication
```
POST /login    - User login
POST /register - User registration
```

### Todos
```
GET    /todos          - Get all todos for authenticated user
GET    /todos/:id      - Get single todo
POST   /todos          - Create new todo
PUT    /todos/:id      - Update existing todo
DELETE /todos/:id      - Delete todo
```

### Categories
```
GET    /todos/categories - Get all categories
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Git

### Installation

1. Clone the repository
```bash
git clone 
cd 
```

2. Set up environment variables
```bash
cp .env.example .env
```

Update the `.env` file with your configuration:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=todo_db
JWT_SECRET=your_jwt_secret
```

3. Install dependencies
```bash
go mod download
```

4. Run database migrations
```bash
go run db/migrate.go
```

5. Start the server
```bash
go run app/main.go
```

The server will start on `http://localhost:8080` (or your configured port)