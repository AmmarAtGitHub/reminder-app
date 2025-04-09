package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Test for GET /tasks endpoint
// This test checks if the endpoint returns a 200 OK status and the expected task data.
func TestGetAllTasks_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Prepare mock rows
	rows := sqlmock.NewRows([]string{"id", "title", "description", "reminder_date"}).
		AddRow(1, "Test Task", "Description", nil)

	mock.ExpectQuery("SELECT id, title, description, reminder_date FROM tasks").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	getAllTasks(rr, req, db)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Test Task")
}

// test POST /tasks with missing reminder date
// This test checks if the endpoint returns a 400 Bad Request status when the reminder date is missing.
func TestAddTask_MissingReminderDate(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Missing reminder_date
	body := `{"title": "Test Task", "description": "Missing date"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := TasksHandler(db)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Reminder date is required")
}

// test DELETE /tasks/{id} success
// This test checks if the endpoint returns a 200 OK status when the task is successfully deleted.
func TestDeleteTask_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rr := httptest.NewRecorder()

	handler := DeleteTask(db)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
