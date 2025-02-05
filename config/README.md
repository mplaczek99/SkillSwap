# config Directory

This directory holds the configuration and initialization code for the application.

## Files

- **config.go**:  
  - Loads environment variables using `godotenv` (if a `.env` file is present).
  - Sets default values for configuration parameters such as `DB_DRIVER`, `DB_SOURCE`, `SERVER_PORT`, and `JWT_SECRET`.
  - Initializes the database connection using GORM with connection pooling.

## Usage

Ensure you have a valid `.env` file or proper environment variables set before running the application.

