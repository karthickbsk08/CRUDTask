package beequeue

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"tasks/helpers"
)

func EmailReminder(pDebug *helpers.HelperStruct, pDueOverTaskRec TaskPayload) error {
	pDebug.Log(helpers.Statement, "EmailReminder(+), Starting Email Reminder Function")

	var tp bytes.Buffer

	template, lErr := template.ParseFiles("html/taskReminder.html")
	if lErr != nil {
		pDebug.Log(helpers.Elog, "error parsing template:", lErr)
		return lErr
	}

	lErr = template.Execute(&tp, pDueOverTaskRec)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "error mappint struct in template:", lErr)
		return lErr
	}

	lEmailbody := tp.String()

	lErr = email(lEmailbody, pDueOverTaskRec.Email)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "error sharing mail:", lErr)
		return lErr
	}

	pDebug.Log(helpers.Statement, "EmailReminder(-), Starting Email Reminder Function")
	return nil
}

func email(message string, pTo string) error {
	// Sender data
	from := "karthickbsk08@gmail.com"
	password := "ksyv roco sgrb injo" // Not your Gmail password

	// Receiver email address
	to := []string{pTo}

	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message body
	// message := []byte("Subject: Hello from Go!\r\n\r\nThis is a test email sent from Go using net/smtp.")

	// This sets the email body as HTML
	message1 := []byte("Subject: Task Reminder\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" + message)
	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message1))
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}

/*
[App Code]

	│
	├─> Compose Message
	│
	├─> Connect to SMTP Server (smtp.gmail.com:587)
	│
	├─> Authenticate (username + password)
	│
	├─> Send Mail (to address)
	│
	└─> Response:
	     ├─ Success (200 OK)
	     └─ Error (e.g., 535 Auth failed)

		 Email Sending Flow

│
├── 1. Compose Email
│   ├── From (Sender's Email)
│   ├── To (Recipient Email(s))
│   ├── Subject
│   └── Body (Plain text or HTML)
│
├── 2. Connect to SMTP Server
│   ├── SMTP Host (e.g., smtp.gmail.com)
│   ├── Port (usually 587 or 465)
│   └── TLS/SSL Encryption (for security)
│
├── 3. Authenticate
│   ├── Email Address (Username)
│   └── Password or App Password (Token)
│
├── 4. Send Email
│   └── Use SMTP Commands (handled by libraries like net/smtp)
│       ├── HELO/EHLO
│       ├── AUTH
│       ├── MAIL FROM
│       ├── RCPT TO
│       ├── DATA
│       └── QUIT
│
└── 5. Confirmation

	├── Success: Email delivered to SMTP server
	└── Error: Log error (e.g., auth failed, connection refused, etc.)
*/
