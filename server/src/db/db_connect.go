package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB //created outside to make it global.

// GetDB returns the global database instance
func GetDB() *sql.DB {
	return Db
}

func ConnectDatabase() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	host := getEnv("DB_HOST", "HOST")
	portStr := getEnv("DB_PORT", "PORT")
	user := getEnv("DB_USER", "DB_USERNAME", "USER")
	dbname := getEnv("DB_NAME")
	pass := getEnv("DB_PASSWORD", "PASSWORD")

	if host == "" || portStr == "" || user == "" || dbname == "" {
		panic("database configuration is incomplete; check .env values")
	}

	port, errPort := strconv.Atoi(portStr)
	if errPort != nil {
		fmt.Println("Invalid database port", errPort)
		panic(errPort)
	}

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {

		fmt.Println("There is an error while connecting to the database ", errSql)
		panic(errSql)
	}

	if errPing := db.Ping(); errPing != nil {
		fmt.Println("Unable to reach the database, please confirm the service is running", errPing)
		panic(errPing)
	}

	ensureSchema(db)
	ensureMessagesTable(db)
	ensureBioDataColumns(db)
	ensureProfileColumns(db)
	ensureLastSeenColumn(db)

	Db = db
	fmt.Println("Successfully connected to database!")
}

func ensureSchema(db *sql.DB) {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)`,
		"users",
	).Scan(&exists)
	if err != nil {
		fmt.Println("Failed to check existing tables", err)
		panic(err)
	}

	if exists {
		fmt.Println("Database tables already exist, skipping init script")
		return
	}

	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine current file path for init script")
	}

	scriptPath := filepath.Join(filepath.Dir(callerFile), "init-script.sql")
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Println("Failed to read init script", err)
		panic(err)
	}

	if _, err = db.Exec(string(script)); err != nil {
		fmt.Println("Failed to execute init script", err)
		panic(err)
	}

	fmt.Println("Database schema created from init script")
}

func ensureMessagesTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_id INT NOT NULL REFERENCES users(id),
		receiver_id INT NOT NULL REFERENCES users(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Failed to create messages table", err)
		panic(err)
	}
	fmt.Println("Messages table ensured")
}

func ensureBioDataColumns(db *sql.DB) {
	// Check and add max_distance_km if missing
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'bio_data' AND column_name = 'max_distance_km'
		)
	`).Scan(&exists)

	if err == nil && !exists {
		_, err = db.Exec(`ALTER TABLE bio_data ADD COLUMN max_distance_km INT`)
		if err != nil {
			fmt.Println("Failed to add max_distance_km column", err)
		} else {
			fmt.Println("Added max_distance_km column to bio_data")
		}
	}

	// Check and add age_range_min if missing
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'bio_data' AND column_name = 'age_range_min'
		)
	`).Scan(&exists)

	if err == nil && !exists {
		_, err = db.Exec(`ALTER TABLE bio_data ADD COLUMN age_range_min INT`)
		if err != nil {
			fmt.Println("Failed to add age_range_min column", err)
		} else {
			fmt.Println("Added age_range_min column to bio_data")
		}
	}

	// Check and add age_range_max if missing
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'bio_data' AND column_name = 'age_range_max'
		)
	`).Scan(&exists)

	if err == nil && !exists {
		_, err = db.Exec(`ALTER TABLE bio_data ADD COLUMN age_range_max INT`)
		if err != nil {
			fmt.Println("Failed to add age_range_max column", err)
		} else {
			fmt.Println("Added age_range_max column to bio_data")
		}
	}
}

func ensureProfileColumns(db *sql.DB) {
	// Check and add latitude if missing
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'profiles' AND column_name = 'latitude'
		)
	`).Scan(&exists)

	if err == nil && !exists {
		_, err = db.Exec(`ALTER TABLE profiles ADD COLUMN latitude FLOAT`)
		if err != nil {
			fmt.Println("Failed to add latitude column", err)
		} else {
			fmt.Println("Added latitude column to profiles")
		}
	}

	// Check and add longitude if missing
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'profiles' AND column_name = 'longitude'
		)
	`).Scan(&exists)

	if err == nil && !exists {
		_, err = db.Exec(`ALTER TABLE profiles ADD COLUMN longitude FLOAT`)
		if err != nil {
			fmt.Println("Failed to add longitude column", err)
		} else {
			fmt.Println("Added longitude column to profiles")
		}
	}
}

func ensureLastSeenColumn(db *sql.DB) {
	// Check and add last_seen if missing
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'last_seen'
		)
	`).Scan(&exists)

	if err == nil && !exists {
		_, err = db.Exec(`ALTER TABLE users ADD COLUMN last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)
		if err != nil {
			fmt.Println("Failed to add last_seen column", err)
		} else {
			fmt.Println("Added last_seen column to users")
		}
	}
}

func getEnv(keys ...string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok && value != "" {
			return value
		}
	}
	return ""
}
