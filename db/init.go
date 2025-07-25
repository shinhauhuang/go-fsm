package main

import (
	"context"
	"fmt"
	"go-fsm/ent"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// --- Database Configuration ---
	dbDriver := os.Getenv("DB_DRIVER")
	var client *ent.Client

	switch dbDriver {
	case "mariadb":
		dsn := os.Getenv("DB_DSN")
		if dsn == "" {
			log.Fatal("DB_DRIVER is 'mariadb' but DB_DSN is not set. Please set the MariaDB DSN.")
		}
		client, err = ent.Open("mysql", dsn)
	default:
		log.Fatal("This script is intended for initializing a persistent database like MariaDB. Please set DB_DRIVER=mariadb.")
	}

	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool to create the schema
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	fmt.Println("Database schema initialized successfully.")
}
