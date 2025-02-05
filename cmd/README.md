# cmd Directory

This directory contains the entry point(s) for the SkillSwap backend application.

## Files

- **main.go**:  
  - Loads the configuration and initializes the database.
  - Sets up the Gin HTTP server with middleware.
  - Registers API routes.
  - Starts the server on the configured port.

## Running

Start the application from the root directory:
```bash
go run cmd/main.go
```

