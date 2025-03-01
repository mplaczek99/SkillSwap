# SkillSwap Backend

## System Architecture

- **Backend**  
  Written in Golang for performance, simplicity, and scalability.
- **Communication**  
  Exposes RESTful APIs that any frontend (web or mobile) can integrate with.
- **Database**  
  Utilizes **PostgreSQL** to securely manage user data, transactions, and schedules.

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
   - Find users by skills, availability, location, etc.  
   - Advanced filters and sorting for efficient matching  

4. **Scheduling and Notifications**  
   - Calendar system to manage sessions  
   - Email and in-app notifications (if implemented by the frontend)  

5. **Real-Time Communication**  
   - Secure messaging endpoint (UI integration up to the frontend)  
   - Facilitates session details and rapport-building  

6. **Administrative Tools**  
   - Admin dashboard endpoints for platform monitoring  
   - Dispute resolution and content moderation  

## Technical Implementation

### Backend

- **Language/Framework**  
  Golang (Gin used in this project)
- **Database**  
  PostgreSQL  
- **Authentication**  
  - JWT for secure sessions
  - (Optional) OAuth 2.0 for social logins
- **Security**  
  - TLS encryption (HTTPS)  
  - bcrypt (or similar) for password hashing  
  - Role-based access controls for admins  

### Deployment

- **Containerization**  
  Docker for consistent development and deployment
- **Cloud Hosting**  
  AWS, GCP, Azure, or any preferred provider
- **CI/CD**  
  Automated pipelines for building, testing, and deploying

## Expected Outcomes

The backend will:

- Handle user profiles, SkillPoint transactions, and scheduling efficiently
- Provide clear API documentation for easy frontend/mobile integration
- Scale effectively as user engagement grows

---

## Getting Started

### Prerequisites

1. **Go** (v1.18+ recommended)
2. **Docker** & **Docker Compose** (Recommended approach)

### Installation & Running

1. **Clone the Repository**  
   ```bash
   git clone https://github.com/mplaczek99/SkillSwap.git
   cd SkillSwap
    ```

2. **Install Dependencies**
    ```bash
    go mod tidy
    ```

3. **Running with Docker** (Recommended)
    - **Ensure local Postgres isn't conflicting**: If you have a local PostgreSQL service running on port 5432, stop it or change its port.

    - **Start the containers**
    ```bash
    docker-compose up --build
    ```
    This will sping up two containers:
    - **db**: A PostgreSQL container
    - **backend**: Your Go application container (listening on port 8080)

    **Access**
        - The API is available at "http://localhost:8080"
        - The database is running inside Docket at "db:5432"

### 📝 Notes

✅ The `--setup-db` flag and automated PostgreSQL installation have been **removed**.

✅ **Port Binding:**  
   - **Database:** Docker binds PostgreSQL to **host port 5432**.  
   - **Backend:** The Go application runs on **host port 8080**.  

⚠️ **Troubleshooting:**  
   - If you encounter **"port already in use"** errors, check that no other service is running on ports **5432** (PostgreSQL) or **8080** (backend).  
   - Use `docker ps` to verify running containers.  
