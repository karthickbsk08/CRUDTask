package common

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func DoMarshallnWriteResponse(pData any, w http.ResponseWriter) {
	if lErr := json.NewEncoder(w).Encode(pData); lErr != nil {
		log.Println("Marshalling error")
	}
}

func DoMarshall(pData any) string {
	lData, _ := json.Marshal(pData)
	return string(lData)

}

func ReadUrlId(r *http.Request, w http.ResponseWriter) string {
	variables := mux.Vars(r)
	if lval, ok := variables["id"]; !ok || lval == "" {
		http.NotFound(w, r)
		return ""
	}
	return variables["id"]
}
