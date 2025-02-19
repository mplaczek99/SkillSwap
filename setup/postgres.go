package setup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

// We read these from ENV so they're not hard-coded.
// Fallbacks can be set if not found in ENV.
var (
	postgresSuperuser = getenv("POSTGRES_SUPERUSER", "postgres")
	postgresSuperPass = getenv("POSTGRES_SUPERPASS", "")
	postgresDataDir   = getenv("POSTGRES_DATA_DIR", "/var/lib/postgres/data")
	// Connection for the postgres superuser; uses "peer" or password auth depending on system setup.
	// If local "peer" auth is configured, you can omit the password.
	postgresSuperDSN = fmt.Sprintf("user=%s password=%s host=localhost sslmode=disable", postgresSuperuser, postgresSuperPass)

	// App-level credentials
	appUser     = getenv("APP_DB_USER", "techie")
	appPassword = getenv("APP_DB_PASS", "techiestrongpassword")
	appDB       = getenv("APP_DB_NAME", "skillswap_db")
	appDSN      = fmt.Sprintf("user=%s password=%s dbname=%s host=localhost sslmode=disable", appUser, appPassword, appDB)
)

// getenv reads an environment variable or uses fallback if not set.
func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// CheckAndInstallPostgres ensures `psql` is installed. Currently supports Arch & Debian/Ubuntu.
func CheckAndInstallPostgres() error {
	_, err := exec.LookPath("psql")
	if err == nil {
		log.Println("PostgreSQL is already installed.")
		return nil
	}

	fmt.Println("PostgreSQL not found. Installing...")

	switch {
	case fileExists("/etc/arch-release"):
		// For Arch Linux
		cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "postgresql")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install postgres via pacman: %w", err)
		}

	case fileExists("/etc/debian_version"):
		// For Debian/Ubuntu
		cmd := exec.Command("sudo", "apt-get", "update")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to update apt-get: %w", err)
		}

		cmd = exec.Command("sudo", "apt-get", "install", "-y", "postgresql")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install postgres via apt-get: %w", err)
		}

	default:
		log.Println("Unsupported or unknown Linux distribution. Please install PostgreSQL manually.")
		return nil
	}

	return nil
}

// InitializePostgres runs initdb if the data directory doesn't exist (Arch default shown).
func InitializePostgres() error {
	if !fileExists(postgresDataDir) {
		fmt.Printf("Initializing PostgreSQL data directory at %s...\n", postgresDataDir)
		cmd := exec.Command("sudo", "-u", "postgres", "initdb", "-D", postgresDataDir)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initdb: %w", err)
		}
	} else {
		fmt.Printf("PostgreSQL data directory already exists at %s.\n", postgresDataDir)
	}
	return nil
}

// StartPostgres starts the PostgreSQL service if it's not running.
// Optionally enables it so it starts on boot.
func StartPostgres() error {
	// Check if active
	err := exec.Command("systemctl", "is-active", "--quiet", "postgresql").Run()
	if err != nil {
		fmt.Println("Starting PostgreSQL via systemd...")
		cmd := exec.Command("sudo", "systemctl", "start", "postgresql")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to start postgresql: %w", err)
		}

		// Enable on boot (optional)
		fmt.Println("Enabling PostgreSQL service on boot...")
		cmd = exec.Command("sudo", "systemctl", "enable", "postgresql")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			// Not fatal if enabling fails, but log it
			log.Println("Warning: failed to enable postgresql on boot:", err)
		}
	} else {
		fmt.Println("PostgreSQL service is already running.")
	}
	return nil
}

// SetupDatabase creates the app database & user if they don't exist, and then fixes schema permissions.
func SetupDatabase() error {
	// Connect as the superuser (e.g., "postgres") to the default database.
	db, err := sql.Open("postgres", postgresSuperDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres DSN: %w", err)
	}
	defer db.Close()

	log.Println("Ensuring database and user exist...")

	// Create the database if it doesn’t exist.
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", appDB)); err != nil {
		log.Printf("Note: %v", err)
	}

	// Create the user if it doesn’t exist.
	if _, err := db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", appUser, appPassword)); err != nil {
		log.Printf("Note: %v", err)
	}

	// Grant privileges on the database.
	if _, err := db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", appDB, appUser)); err != nil {
		log.Printf("Note: %v", err)
	}

	// Build a new DSN that connects to the newly created app database as the superuser.
	postgresAppDSN := fmt.Sprintf("user=%s password=%s host=localhost dbname=%s sslmode=disable", postgresSuperuser, postgresSuperPass, appDB)

	// Close the previous connection and open a new one to the app database as superuser.
	db.Close()
	db, err = sql.Open("postgres", postgresAppDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to app database as superuser: %w", err)
	}
	defer db.Close()

	log.Println("Granting schema permissions to the app user...")

	// Run schema-altering commands as superuser.
	if _, err := db.Exec("ALTER SCHEMA public OWNER TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("GRANT USAGE, CREATE ON SCHEMA public TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("GRANT ALL ON SCHEMA public TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("GRANT ALL ON ALL TABLES IN SCHEMA public TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}
	if _, err := db.Exec("GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO techie;"); err != nil {
		log.Printf("Note: %v", err)
	}

	log.Println("✅ Database setup completed successfully!")
	return nil
}

// fileExists is a helper to check if a path exists and is not a directory
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
