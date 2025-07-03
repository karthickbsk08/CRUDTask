package handler

import (
	"net/http"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
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

	var lTaskResp models.CreateTask
	var lTaskID string
	var lErr error

	if lTaskID = common.ReadUrlId(r, w); lTaskID == "" {
		return
	}

	lTaskResp, lErr = GetTaskByID(lDebug, lTaskID)
	if lErr != nil {
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "HCT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lTaskResp, w)
		return
	}

	common.DoMarshallnWriteResponse(lTaskResp, w)
	lDebug.Log(helpers.Statement, "GetTaskByID_API(-)")

}

func GetTaskByID(pDebug *helpers.HelperStruct, pTaskId string) (models.CreateTask, error) {
	pDebug.Log(helpers.Statement, "GetTaskByID(+)")

	pDebug.Log(helpers.Statement, "GetTaskByID(-)")
	return models.CreateTask{}, nil

}
