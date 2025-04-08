package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AmmarAtGitHub/reminder-app/handlers"
	_ "github.com/AmmarAtGitHub/reminder-app/models"
	"github.com/AmmarAtGitHub/reminder-app/scheduler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// loadEnv loads environment variables from a .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// connectDB establishes a connection to the PostgreSQL database using environment variables
func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr) //create a connetion pool
	if err != nil {
		return nil, err
	}
	// Verify the connection is working
	if err := db.Ping(); err != nil {
		fmt.Println("Successfully connected to the database!")
		return nil, err
	}

	return db, nil
}

func main() {
	// Load environment variables from .env file
	loadEnv()
	// Connect to the database
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	// Ensure the database connection is closed when the program exits
	defer db.Close()
	fmt.Println("Connected to the database!")

	// Run the scheduler in a separate goroutine
	go scheduler.StartReminderScheduler(db)
	fmt.Printf("Starting server at port 8080\n")

	// Set up HTTP handlers
	http.HandleFunc("/tasks", handlers.TasksHandler(db))
	http.HandleFunc("/tasks/", handlers.DeleteTask(db))

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
