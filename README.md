# SkillSwap

## Overview

SkillSwap is a modern web platform where users can exchange skills and knowledge with others using a virtual currency called SkillPoints. The platform enables skill sharing, scheduling of learning sessions, direct messaging between users, and video content sharing.

## Features

- **User Authentication**: Secure registration and login system
- **Skill Marketplace**: Browse, search, and filter skills by category and experience level
- **Real-time Chat**: Direct messaging between users
- **Video Tutorials**: Upload and share instructional videos
- **Session Scheduling**: Book and manage skill exchange sessions
- **SkillPoints Economy**: Virtual currency for the skill exchange marketplace
- **Job Board**: Post and find job opportunities matching your skills
- **User Feedback System**: Rate and review learning experiences

## Technology Stack

### Frontend
- **Framework**: Vue.js 3
- **State Management**: Vuex
- **Routing**: Vue Router
- **UI**: Custom design system with responsive components
- **Icons**: Font Awesome
- **Testing**: Jest

### Backend
- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Database ORM**: GORM
- **Authentication**: JWT
- **API Documentation**: Swagger
- **Testing**: Go testing package

### Database
- PostgreSQL

### DevOps
- Docker
- Docker Compose

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Node.js (for local frontend development)
- Go (for local backend development)

### Running with Docker

1. Clone the repository
   ```
   git clone https://github.com/your-username/skillswap.git
   cd skillswap
   ```

2. Start the application
   ```
   docker-compose up --build
   ```

3. Access the application:
   - **Frontend**: http://localhost:8081
   - **Backend API**: http://localhost:8080
   - **API Documentation**: http://localhost:8080/swagger/index.html

### Default Test User
The system automatically creates a test user:
- **Email**: `test@example.com`
- **Password**: `somepassword`

### Local Development

#### Frontend Development
```
cd frontend
npm install
npm run serve
```

#### Backend Development
```
cd backend
go mod download
go run cmd/main.go
```

## Project Structure

```
SkillSwap/
├── frontend/           # Vue.js frontend
│   ├── src/            # Source code
│   │   ├── components/ # Vue components
│   │   ├── store/      # Vuex store
│   │   ├── router/     # Vue Router
│   │   ├── models/     # Data models
│   │   ├── services/   # API services
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

## Key Components

### Frontend Components
- **Dashboard**: Main user interface with activity feed and skill showcases
- **Profile**: User profiles with skills and ratings
- **Search**: Advanced search with filters for finding skills and users
- **Chat**: Real-time messaging system
- **Schedule**: Calendar for booking skill exchange sessions
- **VideoUpload**: Interface for uploading tutorial videos
- **Transactions**: SkillPoints transaction history
- **JobPostings**: Job board with posting and application features

### Backend API Endpoints
- Auth: `/api/auth/register`, `/api/auth/login`
- Search: `/api/search`
- Schedule: `/api/schedule`
- Videos: `/api/videos/upload`, `/api/videos`
- Protected routes require JWT Authentication
- Admin routes require Admin role

## Configuration

The application can be configured using environment variables in a `.env` file:

```
# Database Configuration
DB_DRIVER=postgres
DB_SOURCE=host=db port=5432 user=techie password=techiestrongpassword dbname=skillswap_db sslmode=disable

# Backend Configuration
JWT_SECRET=your_secret_key_should_be_long_and_secure_in_production
SERVER_PORT=8080

# Frontend Configuration
VUE_APP_API_URL=http://backend:8080
```

## Code Quality

### Formatting

#### Frontend
The frontend uses Prettier for code formatting. To format your code:
```
cd frontend
npm run format
```

#### Backend
The backend uses `gofmt` for code formatting. To format your Go code:
```
cd backend
go fmt ./...
```

### Linting

#### Frontend
ESLint is used for linting the Vue.js code. To run the linter:
```
cd frontend
npm run lint
```

#### Backend
For the Go backend, you can use:
```
cd backend
go vet ./...
```

## Testing

### Frontend
```
cd frontend
npm test
```

### Backend
```
cd backend
go test ./...
```
