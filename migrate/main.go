package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type CustomLogger struct {
	*log.Logger
}

func (c *CustomLogger) Verbose() bool {
	return true
}
func main() {
	// Database connection string
	dbURL := "postgres://dbuser:password@localhost:5432/mydb?sslmode=disable"

	// Path to the migrations directory
	migrationsPath := "file://migrations"

	// Create a new migrate instance
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Unable to create migrate instance: %v", err)
	}
	defer m.Close()

	// Get the command (up or down) from the command line
	if len(os.Args) < 2 {
		log.Fatal("Please specify 'up' or 'down' as the first argument")
	}
	command := os.Args[1]

	// Number of steps (default is 1)
	steps := 1
	if len(os.Args) > 2 {
		_, err := fmt.Sscanf(os.Args[2], "%d", &steps)
		if err != nil {
			log.Fatalf("Invalid number of steps: %v", err)
		}
	}
	// Enable logging
	m.Log = &CustomLogger{
		Logger: log.New(os.Stdout, "Migrate:", log.LstdFlags|log.Lshortfile),
	} // Use the built-in logger

	// Run migrations up
	// Run the migration command
	switch command {
	case "up":
		err = m.Steps(steps)
	case "down":
		err = m.Steps(-steps)
	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down'", command)
	}

	// Handle errors
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No changes applied.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Printf("Migration %s successful!\n", command)
	}
}
