# SkillSwap Platform

SkillSwap is a skill exchange platform where users can share their expertise and learn from others using a virtual currency called SkillPoints.

## System Architecture

- **Frontend:** Vue.js application  
- **Backend:** Go (Golang) with Gin framework  
- **Database:** PostgreSQL

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running the Application

**Clone the repository:**
```bash
git clone https://github.com/mplaczek99/SkillSwap.git
cd SkillSwap
```

**Start the services:**
```bash
docker-compose up --build
```

**Access the application:**

- **Frontend:** http://localhost:8081  
- **Backend API:** http://localhost:8080  
- **API Documentation:** http://localhost:8080/swagger/index.html

## Default Test User

The system automatically creates a test user:

- **Email:** `test@example.com`  
- **Password:** `somepassword`

## Development

### Frontend Development (Vue.js)

The frontend is built with Vue.js 3 and uses:

- Vuex for state management
- Vue Router for navigation
- Font Awesome for icons
- Custom CSS with design system

To run the frontend in development mode:
```bash
cd frontend
npm install
npm run serve
```

### Backend Development (Go)

The backend is built with Go and uses:

- Gin for HTTP routing
- GORM for database operations
- JWT for authentication
- Swagger for API documentation

To run the backend in development mode:
```bash
cd backend
go mod download
go run cmd/main.go
```

## Features

- User authentication (register, login)
- User profiles
- Skill listing and searching
- Chat functionality
- Video upload and processing
- Scheduling for skill exchange sessions

## Project Structure

```
SkillSwap/
├── frontend/           # Vue.js frontend
│   ├── src/            # Source code
│   │   ├── components/ # Vue components
│   │   ├── store/      # Vuex store
│   │   ├── router/     # Vue Router
│   │   └── assets/     # Static assets
│   ├── tests/          # Frontend tests
│   └── public/         # Public assets
├── backend/            # Go backend
│   ├── cmd/            # Entry point
│   ├── config/         # Configuration
│   ├── controllers/    # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   ├── repositories/   # Data access
│   ├── routes/         # API routes
│   ├── services/       # Business logic
│   └── utils/          # Utilities
└── docker-compose.yml  # Docker configuration
```

## Environment Variables

The application uses environment variables for configuration. You can set them in a `.env` file at the root of the project.

**Key variables:**

- `DB_SOURCE`: PostgreSQL connection string
- `JWT_SECRET`: Secret for JWT token generation
- `VUE_APP_API_URL`: Backend API URL for the frontend

