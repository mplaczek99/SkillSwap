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
2. **PostgreSQL** (auto-installed if you use the `--setup-db` flag on macOS or Linux)  
3. (Optional) **Docker** for containerized deployments  

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

3. **(Optional) Configure .env**
    -   Create a .env file to set environment variables like POSTGRES_SUPERUSER, POSTGRES_SUPERPASS, APP_DB_USER, APP_DB_PASS, etc.
    - If you skip this step, defaults will be used.

4. **Database Setup**
    You can automatically install and configure PostgreSQL on macOS, Arch Linux, Debian, or Ubuntu with:
    ```bash
    go run cmd/main.go --setup-db
    ```
    
    This will:
    - Install PostgreSQL if not present  
    - Initialize the data directory (`initdb`)  
    - Start the PostgreSQL service  
    - Create your app’s database and user, then grant permissions

    Once your database is set up, start the application by simply running:
    ```bash
    go run cmd/main.go
    ```

    By default, the server listens on port 8080. You can override this using the SERVER_PORT environment variable.

