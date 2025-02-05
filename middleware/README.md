# middleware Directory

This directory contains custom middleware functions for the application.

## Files

- **auth_middleware.go**:  
  Implements JWT authentication middleware. It:
  - Extracts the JWT token from the `Authorization` header.
  - Validates the token using utility functions.
  - Aborts the request if the token is missing or invalid.

## Usage

This middleware is applied to routes that require authentication. It ensures that only authorized users can access protected endpoints.

