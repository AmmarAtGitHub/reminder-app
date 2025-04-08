package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AmmarAtGitHub/reminder-app/models"
)

// TasksHandler handles both GET and POST requests on "/tasks"
func TasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllTasks(w, r, db) // Call GetAllTasks logic
		case http.MethodPost:
			addTask(w, r, db) // Call AddTaskHandler logic
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Extracted GetAllTasks logic
func getAllTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/tasks" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	rows, err := db.Query("SELECT id, title, description, reminder_date FROM tasks")
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.ReminderDate); err != nil {
			http.Error(w, "Failed to read task", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func addTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Verify endpoint and method
	if r.URL.Path != "/tasks" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return // Crucial return to stop execution
	}

	// Decode JSON payload
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//Validate Reminder Date. Check if ReminderDate is populated before inserting into the database
	if task.ReminderDate == nil {
		log.Println("Reminder date is nil!")
		http.Error(w, "Reminder date is required", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded Task: %+v", task) // Debugging output

	// Insert into database with all required fields
	query := `INSERT INTO tasks 
        (title, description, reminder_date) 
        VALUES ($1, $2, $3) 
        RETURNING id`

	err := db.QueryRow(query,
		task.Title,
		task.Description,
		task.ReminderDate,
	).Scan(&task.ID)

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Failed to add task", http.StatusInternalServerError)
	}

	// Return created task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// DeleteTaskHandler handles deleting a task by ID from the URL
func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the task ID from the URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		// Convert ID to an integer
		taskID, err := strconv.Atoi(pathParts[2])
		if err != nil {
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		// Perform the delete operation
		query := "DELETE FROM tasks WHERE id = $1"
		result, err := db.Exec(query, taskID)
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}

		// Check if any row was actually deleted
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		// Successful deletion
		w.WriteHeader(http.StatusNoContent) // 204 No Content (Success, no response body)
	}
}
