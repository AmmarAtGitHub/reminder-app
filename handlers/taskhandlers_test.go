package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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
