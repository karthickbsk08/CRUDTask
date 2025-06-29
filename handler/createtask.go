package handler

import (
	"encoding/json"
	"log"
	"net/http"
	connectdb "tasks/ConnectDB"
	govalidatorpkg "tasks/GovalidatorPkg"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
	"time"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "CreateTask(+)")

	var lCreateTaskReq models.CreateTask
	var lCreateTaskResp models.CreateTask
	var lVal any

	// Decode JSON directly into struct
	lErr := json.NewDecoder(r.Body).Decode(&lCreateTaskReq)
	if lErr != nil {
		lCreateTaskResp.APIStatus = constants.ErrorCode
		lCreateTaskResp.APIError = "HCT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskResp, w)
		return
	}
	log.Println("input : ", common.DoMarshall(lCreateTaskReq))

	lVal, lErr = govalidatorpkg.CleanAndValidateStruct(lCreateTaskReq)
	if lErr != nil {
		lCreateTaskResp.APIStatus = constants.ErrorCode
		lCreateTaskResp.APIError = "HCT002 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskResp, w)
		return
	}
	v, ok := lVal.(models.CreateTask)
	if !ok {
		lCreateTaskResp.APIStatus = constants.ErrorCode
		lCreateTaskResp.APIError = "HCT003: Invalid struct type"
		common.DoMarshallnWriteResponse(lCreateTaskResp, w)
		return
	}
	lCreateTaskReq = v

	lCreateTaskResp, lErr = InsertTaskInDB(lCreateTaskReq)
	if lErr != nil {
		lCreateTaskResp.APIStatus = constants.ErrorCode
		lCreateTaskResp.APIError = "HCT003 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskResp, w)
		return
	}
	common.DoMarshallnWriteResponse(lCreateTaskResp, w)
	lDebug.Log(helpers.Statement, "CreateTask(-)")

}

// create table tasks (
//    ID SERIAL PRIMARY KEY,
//    Title varchar(200) not null,
//    Description text,
//    Status status_enum NOT NULL DEFAULT 'Pending',
//    Duedate TIMESTAMP,                -- Optional datetime
//    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//    UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//    CreatedBy varchar(100) not null,
//    UpdatedBy varchar(100) not null
// );

func InsertTaskInDB(pCreateTaskReq models.CreateTask) (lCreateTaskResp models.CreateTask, lErr error) {

	var lTask models.Tasks
	lTask.Title = pCreateTaskReq.Title
	lTask.Status = pCreateTaskReq.Status
	lTask.Description = pCreateTaskReq.Description
	lTask.Duedate, _ = time.Parse(time.RFC3339, pCreateTaskReq.DueDate)
	lTask.CreatedBy = "AUTOBOT"
	lTask.UpdatedBy = "AUTOBOT"

	lResult := connectdb.GDB.GRMPostgres.Table(`tasks`).Create(&lTask)
	if lResult.Error != nil {
		lErr = lResult.Error
		log.Println("Error on while insert the task : ", lErr.Error())
		return
	}

	lCreateTaskResp.ID = lTask.ID
	lCreateTaskResp.Title = pCreateTaskReq.Title
	lCreateTaskResp.Description = pCreateTaskReq.Description
	lCreateTaskResp.Status = pCreateTaskReq.Status
	lCreateTaskResp.DueDate = pCreateTaskReq.DueDate
	lCreateTaskResp.CreatedAt = pCreateTaskReq.DueDate
	lCreateTaskResp.UpdatedAt = lTask.CreatedAt.Format(time.RFC3339)
	lCreateTaskResp.UpdatedAt = lTask.UpdatedAt.Format(time.RFC3339)
	return
}
