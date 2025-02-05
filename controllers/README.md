# controllers Directory

This directory contains HTTP controllers that handle incoming API requests for the SkillSwap backend.

## Files

- **auth_controller.go**:  
  Handles authentication endpoints (e.g., user login and JWT token generation).

- **user_controller.go**:  
  Manages user-related endpoints such as registration and profile retrieval.

- **skill_controller.go**:  
  Processes endpoints related to skill operations (e.g., adding or retrieving skills).

## How It Works

Controllers receive HTTP requests, validate input data, and delegate business logic to the services layer. They return JSON responses based on the outcome of each operation.

