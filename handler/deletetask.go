package handler

import (
	"fmt"
	"net/http"
	connectdb "tasks/ConnectDB"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
	"tasks/models"
)

func DeleteTaskByID_API(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "DeleteTaskByID_API(+)")

	var lErr error
	var lTaskId string

	if lTaskId = common.ReadUrlId(r, w); lTaskId == "" {
		return
	}

	lErr = DeleteTaskByID(lDebug, lTaskId)
	if lErr != nil {
		lDebug.Log(helpers.Elog, "Error on while delete the task : ", lErr.Error())
		if lErr.Error() == "HD404" {
			lDebug.Log(helpers.Elog, http.StatusNotFound)
			http.NotFound(w, r)
		} else {
			fmt.Fprint(w, http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprint(w, constants.DeletionSuccessQuote)
	lDebug.Log(helpers.Statement, "DeleteTaskByID_API(-)")

}

func DeleteTaskByID(pDebug *helpers.HelperStruct, pTaskId string) error {
	pDebug.Log(helpers.Statement, "DeleteTaskByID(+)")

	lResult := connectdb.GDB.GRMPostgres.Table(`tasks`).Where("id = ?", pTaskId).Delete(&models.Tasks{})
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, "Error on while delete the task : ", lResult.Error.Error())
		return fmt.Errorf("HDT001 : %s", lResult.Error.Error())
	}

	if lResult.RowsAffected == 0 {
		pDebug.Log(helpers.Elog, "Error on while delete the task : HD404")
		return fmt.Errorf("HD404")
	}

	pDebug.Log(helpers.Statement, "DeleteTaskByID(-)")
	return nil
}
