package setup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	_ "github.com/lib/pq"
)

// Environment-driven configuration.
var (
	postgresSuperuser = getenv("POSTGRES_SUPERUSER", "postgres")
	postgresSuperPass = getenv("POSTGRES_SUPERPASS", "")
	// Default PostgreSQL data directory for Linux; will be overridden on macOS.
	postgresDataDir = getenv("POSTGRES_DATA_DIR", "/var/lib/postgres/data")
	// DSN for the PostgreSQL superuser.
	postgresSuperDSN = fmt.Sprintf("user=%s password=%s host=localhost sslmode=disable", postgresSuperuser, postgresSuperPass)

	// Application-level credentials.
	appUser     = getenv("APP_DB_USER", "techie")
	appPassword = getenv("APP_DB_PASS", "techiestrongpassword")
	appDB       = getenv("APP_DB_NAME", "skillswap_db")
	appDSN      = fmt.Sprintf("user=%s password=%s dbname=%s host=localhost sslmode=disable", appUser, appPassword, appDB)
)

// init adjusts defaults based on OS.
func init() {
	if runtime.GOOS == "darwin" && postgresDataDir == "/var/lib/postgres/data" {
		// Homebrew’s default data directory.
		postgresDataDir = "/usr/local/var/postgres"
	}
}

// getenv reads an environment variable or returns a fallback.
func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// runCmd executes a command, piping its output to Stdout/Stderr.
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CheckAndInstallPostgres ensures that PostgreSQL is installed.
// On macOS, it uses Homebrew; on Arch/Debian it uses the respective package managers.
func CheckAndInstallPostgres() error {
	if _, err := exec.LookPath("psql"); err == nil {
		log.Println("PostgreSQL is already installed.")
		return nil
	}

	fmt.Println("PostgreSQL not found. Installing...")

	switch {
	case runtime.GOOS == "darwin":
		if err := runCmd("brew", "install", "postgresql"); err != nil {
			return fmt.Errorf("failed to install postgres via brew: %w", err)
		}

	case fileExists("/etc/arch-release"):
		if err := runCmd("sudo", "pacman", "-S", "--noconfirm", "postgresql"); err != nil {
			return fmt.Errorf("failed to install postgres via pacman: %w", err)
		}

	case fileExists("/etc/debian_version"):
		if err := runCmd("sudo", "apt-get", "update"); err != nil {
			return fmt.Errorf("failed to update apt-get: %w", err)
		}
		if err := runCmd("sudo", "apt-get", "install", "-y", "postgresql"); err != nil {
			return fmt.Errorf("failed to install postgres via apt-get: %w", err)
		}

	default:
		log.Println("Unsupported OS. Please install PostgreSQL manually.")
		return nil
	}

	return nil
}

// InitializePostgres runs initdb if the data directory does not exist.
// On macOS, it runs as the current user; on Linux, it uses the 'postgres' user.
func InitializePostgres() error {
	if !fileExists(postgresDataDir) {
		fmt.Printf("Initializing PostgreSQL data directory at %s...\n", postgresDataDir)
		var cmd *exec.Cmd
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("initdb", "-D", postgresDataDir)
		} else {
			cmd = exec.Command("sudo", "-u", "postgres", "initdb", "-D", postgresDataDir)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initdb: %w", err)
		}
	} else {
		fmt.Printf("PostgreSQL data directory already exists at %s.\n", postgresDataDir)
	}
	return nil
}

// StartPostgres starts the PostgreSQL service if it is not already running.
// On macOS, it uses Homebrew services; on Linux, it uses systemd.
func StartPostgres() error {
	if runtime.GOOS == "darwin" {
		fmt.Println("Starting PostgreSQL via Homebrew services...")
		if err := runCmd("brew", "services", "start", "postgresql"); err != nil {
			return fmt.Errorf("failed to start postgres via brew services: %w", err)
		}
	} else {
		// Check if the PostgreSQL service is active.
		if err := exec.Command("systemctl", "is-active", "--quiet", "postgresql").Run(); err != nil {
			fmt.Println("Starting PostgreSQL via systemd...")
			if err := runCmd("sudo", "systemctl", "start", "postgresql"); err != nil {
				return fmt.Errorf("failed to start postgresql: %w", err)
			}

			fmt.Println("Enabling PostgreSQL service on boot...")
			if err := runCmd("sudo", "systemctl", "enable", "postgresql"); err != nil {
				// Not fatal if enabling fails.
				log.Println("Warning: failed to enable postgresql on boot:", err)
			}
		} else {
			fmt.Println("PostgreSQL service is already running.")
		}
	}
	return nil
}

// SetupDatabase creates the application database and user (if they do not exist)
// and then sets the appropriate schema permissions.
func SetupDatabase() error {
	// Connect as the superuser to the default database.
	db, err := sql.Open("postgres", postgresSuperDSN)
	if err != nil {
		return fmt.Errorf("failed to connect using postgres DSN: %w", err)
	}
	defer db.Close()

	log.Println("Ensuring database and user exist...")

	// Create the database. Ignore errors if it already exists.
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", appDB)); err != nil {
		log.Printf("Note: %v", err)
	}

	// Create the user. Ignore errors if it already exists.
	if _, err := db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", appUser, appPassword)); err != nil {
		log.Printf("Note: %v", err)
	}

	// Grant privileges on the database.
	if _, err := db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", appDB, appUser)); err != nil {
		log.Printf("Note: %v", err)
	}

	// Connect to the newly created app database as the superuser.
	postgresAppDSN := fmt.Sprintf("user=%s password=%s host=localhost dbname=%s sslmode=disable", postgresSuperuser, postgresSuperPass, appDB)
	dbApp, err := sql.Open("postgres", postgresAppDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to app database as superuser: %w", err)
	}
	defer dbApp.Close()

	log.Println("Granting schema permissions to the app user...")

	queries := []string{
		fmt.Sprintf("ALTER SCHEMA public OWNER TO %s;", appUser),
		fmt.Sprintf("GRANT USAGE, CREATE ON SCHEMA public TO %s;", appUser),
		fmt.Sprintf("GRANT ALL ON SCHEMA public TO %s;", appUser),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO %s;", appUser),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO %s;", appUser),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO %s;", appUser),
		fmt.Sprintf("GRANT ALL ON ALL TABLES IN SCHEMA public TO %s;", appUser),
		fmt.Sprintf("GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO %s;", appUser),
		fmt.Sprintf("GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO %s;", appUser),
	}

	for _, q := range queries {
		if _, err := dbApp.Exec(q); err != nil {
			log.Printf("Note: %v", err)
		}
	}

	log.Println("✅ Database setup completed successfully!")
	return nil
}

// fileExists checks if the given path exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
