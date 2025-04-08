package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// email/sender.go
func SendReminder(to string, taskTitle string) error {
	e := email.NewEmail()
	e.From = "Task Reminder <sender_email@gmail.com>"
	e.To = []string{to}
	e.Subject = "Task Reminder: " + taskTitle
	e.Text = []byte(fmt.Sprintf("Reminder for task: %s", taskTitle))

	return e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", "receiver_email_here", "your_password_string", "smtp.gmail.com"),
	)
}
