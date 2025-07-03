package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	govalidatorpkg "tasks/GovalidatorPkg"
	jwttokengen "tasks/JwtTokenGen"
	"tasks/helpers"
	"tasks/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "Login(+)")

	var lLoginReq models.LoginDetails
	var lToken string
	var lval any

	// Decode JSON directly into struct
	lErr := json.NewDecoder(r.Body).Decode(&lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return

	}

	lErr = govalidatorpkg.CleanAndValidateStruct(lDebug, lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return

	}
	v, ok := lval.(models.LoginDetails)
	if !ok {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	lLoginReq = v

	lToken, lErr = jwttokengen.GenerateJWT(lLoginReq.Username, lLoginReq.Password)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	lErr = InsertLoginUserDetlInDB(lLoginReq)
	if lErr != nil {
		log.Println("Error on decoding")
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, lToken)
	lDebug.Log(helpers.Statement, "Login(-)")

}

func InsertLoginUserDetlInDB(pData models.LoginDetails) error {
	return nil

}
