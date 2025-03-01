# Docker Development Instructions

This guide explains how to use Docker to develop and test the SkillSwap application with the newly added Swagger documentation.

## Prerequisites

- Docker and Docker Compose installed
- Git (to clone the repository)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/mplaczek99/SkillSwap.git
   cd SkillSwap
   ```

2. Generate Swagger documentation before building the Docker images:
   ```bash
   cd backend
   
   # Install swag if you don't have it
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Generate documentation
   make swagger
   
   # Go back to root directory
   cd ..
   ```

3. Build and start the containers:
   ```bash
   docker-compose up --build
   ```

4. Access the application:
   - Frontend: http://localhost:8081
   - Backend API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html

## Updating Swagger Documentation

If you make changes to the API endpoints or models:

1. Update the Swagger annotations in your Go files
2. Regenerate the documentation:
   ```bash
   # Either connect to the running backend container
   docker exec -it skillswap_backend bash
   cd /root
   make swagger
   
   # Or stop the containers and rebuild
   docker-compose down
   cd backend
   make swagger
   cd ..
   docker-compose up --build
   ```

3. Refresh the Swagger UI page to see the changes

## Testing Authentication

To test protected endpoints in Swagger UI:

1. First, use the `/api/auth/login` endpoint to get a token
2. Click the "Authorize" button at the top of the Swagger UI
3. Enter your token in the format `Bearer YOUR_TOKEN_HERE`
4. Click "Authorize"
5. Now you can access protected endpoints

## Troubleshooting

- **Swagger page shows "Failed to load API definition"**: Make sure you generated the Swagger documentation before building the Docker image
- **Authorization not working**: Check that you're using the correct format (`Bearer TOKEN`) and that your token is valid
- **Container startup issues**: Check the Docker logs with `docker-compose logs`
