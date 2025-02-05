# routes Directory

This directory contains the routing configuration for the SkillSwap API.

## Files

- **routes.go**:  
  Defines and registers API endpoints. It groups routes into public and protected endpoints and applies middleware (such as JWT authentication) where needed.

## Usage

Routes are registered in the `main.go` file during the initialization of the Gin engine.

