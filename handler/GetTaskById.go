package handler

import (
	"net/http"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
)

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
