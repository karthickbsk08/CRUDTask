package handler

import (
	"net/http"
	connectdb "tasks/ConnectDB"
	govalidatorpkg "tasks/GovalidatorPkg"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "GetAllTasks(+)")

	var q models.TaskQueryParams
	var lTaskResp models.GetAllTaskResp
	var lVal any

	// Parse query parameters
	lErr := govalidatorpkg.Decoder.Decode(&q, r.URL.Query())
	if lErr != nil {
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "GAT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lTaskResp, w)
		return
	}

	lVal, lErr = govalidatorpkg.CleanAndValidateStruct(lTaskResp)
	if lErr != nil {
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "GAT002 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lTaskResp, w)
		return
	}

	v, ok := lVal.(models.GetAllTaskResp)
	if !ok {
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "GAT003: Error invalid struct type"
		common.DoMarshallnWriteResponse(lTaskResp, w)
		return
	}
	lTaskResp = v

	lTaskResp.Tasks, lErr = GetAllTask_InDB(q)
	if lErr != nil {
		lTaskResp.APIStatus = constants.ErrorCode
		lTaskResp.APIError = "GAT004 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lTaskResp, w)
		return
	}
	lTaskResp.Page = q.Page
	lTaskResp.Limit = q.Limit
	lTaskResp.Total = len(lTaskResp.Tasks)

	common.DoMarshallnWriteResponse(lTaskResp, w)
	lDebug.Log(helpers.Statement, "GetAllTasks(-)")

}

func GetAllTask_InDB(q models.TaskQueryParams) ([]models.CreateTask, error) {

	connectdb.GDB.GRMPostgres.Table(`tasks`).Select(``)

	return []models.CreateTask{}, nil
}
