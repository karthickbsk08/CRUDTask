package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	connectdb "tasks/ConnectDB"
	govalidatorpkg "tasks/GovalidatorPkg"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
	"time"
)

func UpdateTaskByAPI(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "UpdateTaskByAPI(+)")

	var lUpdateTaskReq models.CreateTask
	var lUpdateTaskResp models.CreateTask
	var lTaskID string
	var lErr error

	if lTaskID = common.ReadUrlId(r, w); lTaskID == "" {
		return
	}
	// Decode JSON directly into struct
	lErr = json.NewDecoder(r.Body).Decode(&lUpdateTaskReq)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT001 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lUpdateTaskResp))
		return
	}

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, &lUpdateTaskReq)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT002 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lUpdateTaskResp))
		return
	}

	lErr = common.TypeChecker[models.CreateTask](&lUpdateTaskReq)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT003 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lUpdateTaskResp))
		return
	}
	// Convert lTaskID from string to int
	var lTaskIDInt int
	lTaskIDInt, lErr = strconv.Atoi(lTaskID)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT004 : Invalid Task ID"
		fmt.Fprint(w, common.DoMarshall(lUpdateTaskResp))
		return
	}
	lUpdateTaskReq.ID = lTaskIDInt

	lUpdateTaskResp, lErr = UpdateTaskByID(lDebug, &lUpdateTaskReq)
	if lErr != nil {
		lUpdateTaskResp.APIStatus = constants.ErrorCode
		lUpdateTaskResp.APIError = "HCT005 : " + lErr.Error()
		fmt.Fprint(w, common.DoMarshall(lUpdateTaskResp))
		return
	}

	fmt.Fprint(w, common.DoMarshall(lUpdateTaskResp))
	lDebug.Log(helpers.Statement, "UpdateTaskByAPI(-)")

}

func UpdateTaskByID(pDebug *helpers.HelperStruct, pUpdateReq *models.CreateTask) (models.CreateTask, error) {
	pDebug.Log(helpers.Statement, "UpdateTaskByID(+)")

	var pUpdateRespRec models.Tasks
	var lUpdateTaskResp models.CreateTask

	lUpdateResult := connectdb.GDB.GRMPostgres.Table(`tasks`).
		Where("id = ?", pUpdateReq.ID).
		Select(`id`, `title`, `description`, `status`, `duedate`, `createdat`, `updatedat`).
		UpdateColumns(map[string]any{
			"title":       pUpdateReq.Title,
			"description": pUpdateReq.Description,
			"status":      pUpdateReq.Status,
			"duedate":     pUpdateReq.DueDate,
			"createdat":   time.Now().Format(time.RFC3339),
			"updatedat":   time.Now().Format(time.RFC3339),
		})
	if lUpdateResult.Error != nil {
		pDebug.Log(helpers.Elog, "UpdateTaskByID(-) Error: "+lUpdateResult.Error.Error())
		return lUpdateTaskResp, helpers.ErrReturn(fmt.Errorf("UTI001 : %w", lUpdateResult.Error))
	}

	lSelectResult := connectdb.GDB.GRMPostgres.Table(`tasks`).
		Where("id = ?", pUpdateReq.ID).
		First(&pUpdateRespRec)
	if lSelectResult.Error != nil {
		pDebug.Log(helpers.Elog, "UpdateTaskByID(-) Error: "+lSelectResult.Error.Error())
		return lUpdateTaskResp, helpers.ErrReturn(fmt.Errorf("UTI002 : %w", lSelectResult.Error))
	}

	lUpdateTaskResp.ID = pUpdateRespRec.ID
	lUpdateTaskResp.Title = pUpdateRespRec.Title
	lUpdateTaskResp.Description = pUpdateRespRec.Description
	lUpdateTaskResp.Status = pUpdateRespRec.Status
	lUpdateTaskResp.DueDate = pUpdateRespRec.Duedate.Format(time.RFC3339)
	lUpdateTaskResp.CreatedAt = pUpdateRespRec.CreatedAt.Format(time.RFC3339)
	lUpdateTaskResp.UpdatedAt = pUpdateRespRec.UpdatedAt.Format(time.RFC3339)

	pDebug.Log(helpers.Statement, "UpdateTaskByID(-)")
	return lUpdateTaskResp, nil

}
