
# 📝 Reminder App – Task Scheduler & Email Reminder (REST API)

A simple and clean task reminder application written in Go. This app allows users to add tasks with reminder dates, and it sends email notifications when those tasks are due. The system uses a RESTful API to manage tasks and a background scheduler that triggers reminder emails at regular intervals.



## 🚀 Features

- ✅ REST API to create, list, and delete tasks via HTTP endpoints
- 📬 Schedule and send Email reminders (with Gmail SMTP) for tasks when due using Go routines
- ⏳ Background scheduler running every minute
- 🐘 PostgreSQL for persistent storage  
- 🔐 Secure credentials using environment variables



## ⚙️ Setup Instructions

#### 1. Clone the Repository
`git clone https://github.com/AmmarAtGitHub/reminder-app.git`  
`cd reminder-app`

#### 2. Set Environment Variables
Create a .env file in the project root:  
```
DB_HOST=localhost  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=yourpassword  
DB_NAME=reminder
``` 

⚠️ Do not commit .env or creds.json to Git.  

#### 3. Initialize Database
Make sure PostgreSQL is running, then:  
`createdb reminder`  
`psql -U postgres -d reminder -f migrations/001_create_task_table.sql`

#### 4. Install Go Dependencies
`go mod tidy`

#### 5. Run the Application
``go run ./main.go``





### 📬 Gmail App Password Setup
The app uses Gmail to send reminders. For security, Gmail requires an App Password instead of your actual password.

Steps to Generate an App Password:  
1- Visit: https://myaccount.google.com/security  
2- Enable 2-Step Verification  
3- Go to App Passwords  
4- Select App: "Mail", Device: "Other"  
5- Copy the generated 16-character password  
6- Use this in email/reminder.go:  

``smtp.PlainAuth("", "your-email@gmail.com", "your-app-password", "smtp.gmail.com")``  
💡 You can also configure a different SMTP provider like Outlook, SendGrid, etc.





### 📅 Reminder Scheduler
The scheduler runs every minute and performs the following:  
1- Queries tasks with reminder_date <= NOW() and notified = false  
2- Sends an email to the configured address  
3- Marks the task as notified = true  
This logic is found in:  

``go StartReminderScheduler(db)``  
Defined in scheduler/checker.go.  



🔗 API Specification  
📥 Create a Task  
POST /tasks  
~~~json
{
  "title": "Finish report",
  "description": "Annual performance report",
  "reminder_date": "2025-04-10T14:30:00Z"
}
~~~
Response: 
~~~json
 
{
  "id": 1,
  "title": "Finish report",
  "description": "Annual performance report",
  "reminder_date": "2025-04-10T14:30:00Z"
}
~~~  
📤 Get All Tasks  
GET /tasks  

Response:  
~~~json
[
  {
    "id": 1,
    "title": "Finish report",
    "description": "Annual performance report",
    "reminder_date": "2025-04-10T14:30:00Z"
  }
]
~~~  
❌ Delete a Task  
DELETE /tasks/{id}  

``Response: HTTP 204 No Content`` 



### 🔒 Security Notes  
- Never expose your .env or credentials in version control.  
- Use HTTPS and token-based authentication if deploying publicly.  
