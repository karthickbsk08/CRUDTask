package handler

import (
	"fmt"
	"net/http"
	"tasks/common"
	"tasks/constants"
	"tasks/helpers"
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
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, constants.DeletionSuccessQuote)

	lDebug.Log(helpers.Statement, "DeleteTaskByID_API(-)")

}

func DeleteTaskByID(pDebug *helpers.HelperStruct, pTaskId string) error {
	pDebug.Log(helpers.Statement, "DeleteTaskByID(+)")

	pDebug.Log(helpers.Statement, "DeleteTaskByID(-)")
	return nil
}
