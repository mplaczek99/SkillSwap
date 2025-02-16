# SkillSwap Backend

## System Architecture

- **Backend:** Implemented in Golang for performance and scalability.
- **Communication:** Exposes RESTful APIs to integrate with any frontend client (e.g., web or mobile).
- **Database:** Utilizes a relational database (MySQL or PostgreSQL) to manage user data, transactions, and schedules securely.

## Features

1. **User Profiles**
   - Display name, profile picture, bio
   - Skills offered and sought
   - Ratings and feedback

2. **Skill Exchange System**
   - SkillPoints as a virtual currency
   - Users earn points by teaching, spend points by learning
   - Encourages active participation and balanced value exchange

3. **Search and Match**
   - Search for other users by skills, availability, location, etc.
   - Advanced filters and sorting for efficient matching

4. **Scheduling and Notifications**
   - Calendar system to manage sessions
   - Email and in-app notifications (if integrated on the frontend) for reminders

5. **Real-Time Communication**
   - Secure messaging system endpoints (frontends can implement UIs to connect)
   - Facilitates session details and rapport-building

6. **Administrative Tools**
   - Admin dashboard endpoints for monitoring platform usage
   - Dispute resolution and content moderation

## Technical Implementation

### Backend

- **Language/Framework:** Golang (e.g., using Gin or Echo)
- **Database:** MySQL or PostgreSQL
  - Stores user data, skills, transactions, schedules
- **Authentication:** JWT for secure sessions; optional OAuth 2.0 for social logins
- **Security:**
  - TLS encryption (HTTPS)
  - bcrypt or similar for password hashing
  - Role-based access controls for admins

### Deployment

- **Containerization:** Docker for consistent development and deployment
- **Cloud Hosting:** AWS, GCP, or Azure (or any preferred provider)
- **CI/CD:** Automated pipelines for building, testing, and deploying

## Expected Outcomes

The backend will:

- Handle user profiles, SkillPoint transactions, and scheduling efficiently
- Provide clear documentation and endpoints for easy frontend or mobile integration
- Scale to accommodate growing user engagement

---

## Getting Started

### Prerequisites

- **Go** (1.18+ recommended)
- **MySQL or PostgreSQL** for data persistence
- (Optional) **Docker** for containerized deployments

### Installation & Running

```bash
# Clone the repository (main branch)
git clone https://github.com/mplaczek99/SkillSwap.git
cd SkillSwap

# Tidy up dependencies
go mod tidy

# Run the application
go run cmd/main.go

# Or build and run with Docker
docker-compose up --build
```
---
> ⚠️ **Note:**
> - This README focuses on the backend. The frontend is maintained separately.
> - This documentation may evolve as the project matures.
```
