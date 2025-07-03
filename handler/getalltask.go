package handler

import (
	"fmt"
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

	var queryParam models.TaskQueryParams
	var lAllTaskResp models.GetAllTaskResp
	var lErr error
	lAllTaskResp.APIStatus = constants.SuccessCode
	lAllTaskResp.APIError = ""

	// Parse query parameters
	lErr = govalidatorpkg.Decoder.Decode(&queryParam, r.URL.Query())
	if lErr != nil {
		lAllTaskResp.APIStatus = constants.ErrorCode
		lAllTaskResp.APIError = "GAT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lAllTaskResp, w)
		return
	}

	//Testing
	common.DoMarshall(queryParam)

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, &queryParam)
	if lErr != nil {
		lAllTaskResp.APIStatus = constants.ErrorCode
		lAllTaskResp.APIError = "GAT002 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lAllTaskResp, w)
		return
	}

	lErr = common.TypeChecker[models.TaskQueryParams](&queryParam)
	if lErr != nil {
		lAllTaskResp.APIStatus = constants.ErrorCode
		lAllTaskResp.APIError = "GAT003" + lErr.Error()
		common.DoMarshallnWriteResponse(lAllTaskResp, w)
		return
	}

	lAllTaskResp.Tasks, lErr = GetAllTask_InDB(lDebug, queryParam)
	if lErr != nil {
		lAllTaskResp.APIStatus = constants.ErrorCode
		lAllTaskResp.APIError = "GAT004 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lAllTaskResp, w)
		return
	}
	lAllTaskResp.Page = queryParam.Page
	lAllTaskResp.Limit = queryParam.Limit
	lAllTaskResp.Total = len(lAllTaskResp.Tasks)

	common.DoMarshallnWriteResponse(lAllTaskResp, w)
	lDebug.Log(helpers.Statement, "GetAllTasks(-)")

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
