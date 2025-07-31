package beequeue

import (
	"context"
	"encoding/json"
	"fmt"
	"tasks/helpers"
	"time"

	"github.com/hibiken/asynq"
)

func HandleSendReminderTask(ctx context.Context, task *asynq.Task) error {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "HandleSendReminderTask(+), Starting Task Handler for Email Reminder")
	var payload TaskPayload
	var lErr error
	if lErr = json.Unmarshal(task.Payload(), &payload); lErr != nil {
		return lErr
	}

	dueTime, _ := time.Parse(time.RFC3339, payload.DueDate)
	if time.Now().After(dueTime) {
		fmt.Printf("📧 Sending email to %s: Task '%s' is overdue\n", payload.Email, payload.Title)
		lErr = EmailReminder(lDebug, payload)
		if lErr != nil {
			lDebug.Log(helpers.Elog, " Error sending email reminder:", lErr)
			return lErr
		}
	} else {
		fmt.Printf("✅ Task '%s' is not overdue. Skipping.\n", payload.Title)
	}
	return nil
}
