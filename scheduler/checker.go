package scheduler

import (
	"database/sql"
	"log"

	"github.com/AmmarAtGitHub/reminder-app/email"
	"github.com/AmmarAtGitHub/reminder-app/models"

	"time"
)

func GetDueTasks(db *sql.DB) ([]models.Task, error) {
	rows, err := db.Query(`
	select id, title, description, reminder_date, notified, is_completed, created_at
	from tasks
	where reminder_date <= NOW() AND notified = false
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.ReminderDate, &task.Notified, &task.IsCompleted, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
func MarkTaskAsNotified(db *sql.DB, taskID int) error {
	_, err := db.Exec(`
		UPDATE tasks
		SET notified = TRUE
		WHERE id = $1
	`, taskID)
	return err
}
func StartReminderScheduler(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Checking for due tasks...")

		tasks, err := GetDueTasks(db)
		if err != nil {
			log.Printf("Error fetching due tasks: %v", err)
			continue
		}

		for _, task := range tasks {
			log.Printf("Sending reminder for task: %s", task.Title)

			err := email.SendReminder("ammar.shams.eddin@gmail.com", task.Title)
			if err != nil {
				log.Printf("Error sending email for task %d: %v", task.ID, err)
				continue
			}

			err = MarkTaskAsNotified(db, task.ID)
			if err != nil {
				log.Printf("Error updating notified status for task %d: %v", task.ID, err)
				continue
			}

			log.Printf("Reminder sent and status updated for task: %s", task.Title)
		}
	}
}
