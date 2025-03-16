# Todo List REST API

A simple REST API for managing todo tasks built with Go 1.22+. Features in-memory storage, middleware, graceful shutdown, and testing.

## Features

- RESTful endpoints for task management
- In-memory storage with thread-safe operations
- Request logging middleware
- Graceful shutdown handling
- Proper HTTP status codes and error handling
- JSON request/response format
- Concurrent request handling

## Installation

### Prerequisites
- Go 1.22 or higher

## Running the Application

Choose one of the following methods to start the application:

### 1. Using Go Run
```bash
go run ./cmd/api
```

### 2. Build and Run
```bash
go build -o todo-api ./cmd/api

./todo-api
```

### 3. Air (Live Reload, good for development enviroments)
```bash
go install github.com/air-verse/air@latest

air
```

### 4. Docker compose
```bash
docker compose up -d
```

## Testing the application

```bash
go test -v ./internal/handlers/...
```