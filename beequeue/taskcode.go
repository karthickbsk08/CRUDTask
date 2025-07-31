package beequeue

import (
	"fmt"
	connectdb "tasks/ConnectDB"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
	"time"

	"github.com/hibiken/asynq"
)

var (
	Gclient = asynq.NewClient(asynq.RedisClientOpt{Addr: constants.RedisServerAddress})
)

type TaskPayload struct {
	Email   string
	Title   string
	DueDate string
}

func PushTasksIntoRedis(pDebug *helpers.HelperStruct, pPayLoad TaskPayload) {

	// var q models.TaskQueryParams
	// q.DueDateBefore = time.Now().Format(time.RFC3339)
	var info *asynq.TaskInfo
	var lErr error

	taskInfo := asynq.NewTask(constants.UniqueTaskIndentifier, []byte(common.DoMarshall(pPayLoad)))
	if taskInfo != nil {
		info, lErr = Gclient.Enqueue(taskInfo, asynq.Queue("default"))
		if lErr != nil {
			fmt.Println("Enqueue error:", lErr)
		}
	}

	fmt.Println("Task enqueued:", info.ID)
}
func GetAllTask_InDB(pDebug *helpers.HelperStruct, q models.TaskQueryParams) ([]models.Tasks, error) {
	pDebug.Log(helpers.Statement, "GetAllTask_InDB(+)")

	var lAlltasks []models.Tasks
	lQuery := ``
	lOrderQuery := ``

	if q.Page > 0 && q.Limit > 0 {
		// q.Limit += constants.Offset_Auto_increment
		lQuery += fmt.Sprintf(` and id between (%d) and (%d) `, ((q.Page - 1) * q.Limit), (((q.Page - 1) * q.Limit) + q.Limit))
	}

	if q.Status != "" {
		lQuery += fmt.Sprintf(` and status = '%s' `, q.Status)
	}

	if q.DueDateAfter != "" {
		lQuery += fmt.Sprintf(` and duedate > '%s' `, q.DueDateAfter)
	}

	if q.DueDateBefore != "" {
		lQuery += fmt.Sprintf(` and duedate < '%s' `, q.DueDateBefore)
	}

	if q.SortOrder != "" && q.SortBy != "" {
		lOrderQuery += fmt.Sprintf(`%s %s `, q.SortBy, q.SortOrder)
	}

	lResult := connectdb.GDB.GRMPostgres.Table(`tasks`).Where(`1=1` + lQuery).Order(lOrderQuery).Scan(&lAlltasks)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, "HGATD001", lResult.Error)
		return lAlltasks, helpers.ErrReturn(lResult.Error)
	}

	pDebug.Log(helpers.Statement, "GetAllTask_InDB(-)")
	return lAlltasks, nil
}

func PushTasksAsSchedular(pDebug *helpers.HelperStruct) {
	// Example task
	var info *asynq.TaskInfo
	var lErr error

	var q models.TaskQueryParams
	q.DueDateBefore = time.Now().Format(time.RFC3339)

	lDueCrossedTaskArr, lErr := GetAllTask_InDB(pDebug, q)
	if lErr != nil {
		fmt.Println("Enqueue error:", lErr)
		return
	}

	for _, task := range lDueCrossedTaskArr {
		payload := TaskPayload{
			Email:   "karthick.b@flattrade.co.in",
			Title:   task.Title,
			DueDate: task.Duedate.Format(time.RFC3339),
		}

		taskInfo := asynq.NewTask(constants.UniqueTaskIndentifier, []byte(common.DoMarshall(payload)))
		if taskInfo != nil {
			info, lErr = Gclient.Enqueue(taskInfo, asynq.Queue("default"))
			if lErr != nil {
				fmt.Println("Enqueue error:", lErr)
				continue
			}
		}
	}
	fmt.Println("Task enqueued:", info.ID)

}

func TaskSync(pDebug *helpers.HelperStruct) {

	lTickerChan := time.NewTicker(1 * time.Minute)

	for range lTickerChan.C {
		PushTasksAsSchedular(pDebug)
	}
}
