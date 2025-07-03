package handler

import (
	"encoding/json"
	"net/http"
	govalidatorpkg "tasks/GovalidatorPkg"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
)

func UpdateTaskByAPI(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "UpdateTaskByAPI(+)")

	var lUpdateTaskReq models.CreateTask
	var lUpdateTaskResp models.CreateTask
	var lTaskID string
	var lErr error
	var lVal any

	if lTaskID = common.ReadUrlId(r, w); lTaskID == "" {
		return
	}
	// Decode JSON directly into struct
	lErr = json.NewDecoder(r.Body).Decode(&lUpdateTaskReq)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
		return
	}

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, lUpdateTaskReq)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
		return
	}
	v, ok := lVal.(models.CreateTask)
	if !ok {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT001 : " + "Invalid struct type"
		common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
		return
	}
	lUpdateTaskReq = v

	lUpdateTaskResp, lErr = UpdateTaskByID(lDebug, lTaskID)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT003 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
		return
	}

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, lUpdateTaskResp)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
		return
	}
	val, ok := lVal.(models.CreateTask)
	if !ok {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT001 : " + "Invalid struct type"
		common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
		return
	}
	lUpdateTaskResp = val

	common.DoMarshallnWriteResponse(lUpdateTaskResp, w)
	lDebug.Log(helpers.Statement, "UpdateTaskByAPI(-)")

}

func UpdateTaskByID(pDebug *helpers.HelperStruct, pTaskId string) (models.CreateTask, error) {
	pDebug.Log(helpers.Statement, "UpdateTaskByID(+)")

	pDebug.Log(helpers.Statement, "UpdateTaskByID(-)")
	return models.CreateTask{}, nil

}
