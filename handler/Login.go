package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	connectdb "tasks/ConnectDB"
	govalidatorpkg "tasks/GovalidatorPkg"
	jwttokengen "tasks/JwtTokenGen"
	"tasks/common"
	"tasks/encryption"
	"tasks/helpers"
	"tasks/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "Login(+)")

	var lLoginReq models.LoginDetails
	var lToken string

	// Decode JSON directly into struct
	lErr := json.NewDecoder(r.Body).Decode(&lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	log.Println("input json : ", common.DoMarshall(lLoginReq))

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, &lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	lErr = common.TypeChecker[models.LoginDetails](&lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	lToken, lErr = jwttokengen.GenerateJWT(lDebug, lLoginReq.Username, lLoginReq.Password)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	// catching.SetToCache(lDebug, catching.GRedisClient, lLoginReq.Username, lToken, jwttokengen.JWT_EXPIRATION)

	lErr = InsertLoginUserDetlInDB(lDebug, &lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, lToken)
	lDebug.Log(helpers.Statement, "Login(-)")

}

func InsertLoginUserDetlInDB(pDebug *helpers.HelperStruct, pNewUserRec *models.LoginDetails) error {
	pDebug.Log(helpers.Statement, "InsertLoginUserDetlInDB(+)")

	var lUserRec models.Users
	lUserRec.UserName = pNewUserRec.Username
	lUserRec.CreatedBy = pNewUserRec.Username
	lUserRec.UpdatedBy = pNewUserRec.Username

	lBytePassword, lErr := encryption.HashPassword(pNewUserRec.Password)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "HGTBI001 : ", lErr.Error())
		return helpers.ErrReturn(lErr)
	}
	lUserRec.PasswordHash = string(lBytePassword)

	lResult := connectdb.GDB.GRMPostgres.Table(`users`).Create(&lUserRec)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, "HGTBI001 : ", lResult.Error.Error())
		return helpers.ErrReturn(lResult.Error)
	}
	pDebug.Log(helpers.Statement, "InsertLoginUserDetlInDB(-)")
	return nil

}
