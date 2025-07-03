package common

import (
	"encoding/json"
	"fmt"
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

// In a common package
func TypeChecker[T any](val any) error {
	var zero T
	_, ok := val.(*T)
	if !ok {
		return fmt.Errorf("error: expected type %T, got %T", zero, val)
	}
	return nil
}
