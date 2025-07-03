package handler

import (
	"encoding/json"
	"net/http"
	connectdb "tasks/ConnectDB"
	govalidatorpkg "tasks/GovalidatorPkg"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
	"time"
)

// CreateTask handles the HTTP request for creating a new task.
// It performs the following steps:
//   1. Initializes a debug helper for logging.
//   2. Decodes the incoming JSON request body into a CreateTask struct.
//   3. Validates and cleans the input using a validator package.
//   4. Checks the type of the request struct.
//   5. Inserts the new task into the database.
//   6. Writes the response back to the client, including any errors encountered.
//
// The function expects a JSON payload matching the CreateTask model and responds
// with the created task details or an error message in case of failure.
//
// Parameters:
//   - w: http.ResponseWriter for writing the HTTP response.
//   - r: *http.Request containing the HTTP request data.
//
// Error Codes:
//   - HCT001: Error decoding JSON request.
//   - HCT002: Validation error.
//   - HCT003: Type checking error.
//   - HCT004: Database insertion error.

// InsertTaskInDB inserts a new task record into the database.
// It maps the fields from the CreateTask request to the Tasks model,
// sets audit fields, and saves the record using GORM.
//
// Parameters:
//   - pDebug: Pointer to HelperStruct for logging.
//   - pCreateTaskReq: Pointer to CreateTask struct containing the task data.
//
// Returns:
//   - lErr: Error if the database operation fails, otherwise nil.
//
// On success, updates the CreateTask request struct with the new task's ID,
// CreatedAt, and UpdatedAt timestamps.
func CreateTask(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "CreateTask(+)")

	var lCreateTaskReq models.CreateTask
	var lErr error
	// var lCreateTaskResp models.CreateTask
	lCreateTaskReq.APIStatus = constants.SuccessCode
	lCreateTaskReq.APIError = ""

	// Decode JSON directly into struct
	lErr = json.NewDecoder(r.Body).Decode(&lCreateTaskReq)
	if lErr != nil {
		lCreateTaskReq.APIStatus = constants.ErrorCode
		lCreateTaskReq.APIError = "HCT001 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskReq, w)
		return
	}

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, &lCreateTaskReq)
	if lErr != nil {
		lCreateTaskReq.APIStatus = constants.ErrorCode
		lCreateTaskReq.APIError = "HCT002 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskReq, w)
		return
	}
	lErr = common.TypeChecker[models.CreateTask](&lCreateTaskReq)
	if lErr != nil {
		lCreateTaskReq.APIStatus = constants.ErrorCode
		lCreateTaskReq.APIError = "HCT003 :" + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskReq, w)
		return
	}

	lErr = InsertTaskInDB(lDebug, &lCreateTaskReq)
	if lErr != nil {
		lCreateTaskReq.APIStatus = constants.ErrorCode
		lCreateTaskReq.APIError = "HCT004 : " + lErr.Error()
		common.DoMarshallnWriteResponse(lCreateTaskReq, w)
		return
	}
	// constants.Offset_Auto_increment = constants.Offset_Auto_increment + 1
	common.DoMarshallnWriteResponse(lCreateTaskReq, w)
	lDebug.Log(helpers.Statement, "CreateTask(-)")

}

func InsertTaskInDB(pDebug *helpers.HelperStruct, pCreateTaskReq *models.CreateTask) (lErr error) {
	pDebug.Log(helpers.Statement, "InsertTaskInDB(+)")

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
		pDebug.Log(helpers.Elog, "Error on while insert the task : ", lErr.Error())
		return
	}

	pCreateTaskReq.ID = lTask.ID
	pCreateTaskReq.CreatedAt = lTask.CreatedAt.Format(time.RFC3339)
	pCreateTaskReq.UpdatedAt = lTask.UpdatedAt.Format(time.RFC3339)
	pDebug.Log(helpers.Statement, "InsertTaskInDB(-)")
	return
}
