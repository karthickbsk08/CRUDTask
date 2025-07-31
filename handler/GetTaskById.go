package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	connectdb "tasks/ConnectDB"
	"tasks/catching"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
	"time"
)

// GetTaskByID_API handles HTTP requests to retrieve a task by its ID.
// It reads the task ID from the URL, fetches the corresponding task using GetTaskByID,
// and writes the response as JSON. If an error occurs, it returns an error response.
// The function also logs entry and exit points for debugging purposes.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request containing the HTTP request.
//
// Response:
//   - On success: JSON representation of the requested task.
//   - On failure: JSON with error code and message.
func GetTaskByID_API(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "GetTaskByID_API(+)")

	// ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// defer cancel()

	ctx := context.Background()

	var lTaskResp models.CreateTask
	lTaskResp.APIStatus = constants.SuccessCode
	lTaskResp.APIError = ""
	var lTaskID string
	var lCatchedJson string
	var lErr error

	if lTaskID = common.ReadUrlId(r, w); lTaskID == "" {
		return
	}
	lDebug.Log(helpers.Details, "GetTaskByID_API(+) : lTaskID = "+lTaskID)

	lCatchedJson, lErr = catching.GetFromCache(lDebug, catching.GRedisClient, ctx, fmt.Sprintf("task:%s", lTaskID))
	if lErr != nil {
		lDebug.Log(helpers.Elog, "HGTBI001 : ", lErr.Error())
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "HGTBI001 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lTaskResp))
		return
	}

	if lCatchedJson != "" {
		lDebug.Log(helpers.Details, "GetTaskByID_API(+) : Cache hit for task ID "+lTaskID)
		lErr = json.Unmarshal([]byte(lCatchedJson), &lTaskResp)
		if lErr != nil {
			lDebug.Log(helpers.Elog, "HGTBI002 : ", lErr.Error())
			lTaskResp.APIStatus = constants.ErrorCode
			lTaskResp.APIError = "HGTBI002 : " + lErr.Error()
			fmt.Fprint(w, common.DoMarshall(lTaskResp))
			return
		}
		fmt.Fprint(w, common.DoMarshall(lTaskResp))
		fmt.Printf("%s", "get from redis cache")
		lDebug.Log(helpers.Statement, "GetTaskByID_API(-)")
		return
	}
	lErr = GetTaskByID(lDebug, lTaskID, &lTaskResp)
	if lErr != nil {
		lDebug.Log(helpers.Elog, "HGTBI003 : ", lErr.Error())
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "HGTBI003 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lTaskResp))
		return
	}

	lErr = catching.SetToCache(lDebug, catching.GRedisClient, ctx, fmt.Sprintf("task:%d", lTaskResp.ID), common.DoMarshall(lTaskResp), time.Duration(10*time.Second))
	if lErr != nil {
		lDebug.Log(helpers.Elog, "HGTBI004 : ", lErr.Error())
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "HGTBI004 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lTaskResp))
		return
	}
	fmt.Fprint(w, common.DoMarshall(lTaskResp))
	lDebug.Log(helpers.Statement, "GetTaskByID_API(-)")

}

func GetTaskByID(pDebug *helpers.HelperStruct, pTaskId string, pTaskRec *models.CreateTask) error {
	pDebug.Log(helpers.Statement, "GetTaskByID(+)")

	var lTaskRec models.Tasks

	lResult := connectdb.GDB.GRMPostgres.Table(`tasks`).Where("id = ?", pTaskId).Scan(&lTaskRec)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, "HGTBI001 : ", lResult.Error.Error())
		return helpers.ErrReturn(lResult.Error)
	}

	pTaskRec.ID = lTaskRec.ID
	pTaskRec.Title = lTaskRec.Title
	pTaskRec.Description = lTaskRec.Description
	pTaskRec.Status = lTaskRec.Status
	pTaskRec.DueDate = lTaskRec.Duedate.Format("2006-01-02")
	pTaskRec.CreatedAt = lTaskRec.CreatedAt.Format("2006-01-02 15:04:05")
	pTaskRec.UpdatedAt = lTaskRec.UpdatedAt.Format("2006-01-02 15:04:05")

	pDebug.Log(helpers.Statement, "GetTaskByID(-)")
	return nil

}
