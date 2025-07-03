package main

import (
	"log"
	"net/http"
	"os"
	connectdb "tasks/ConnectDB"
	"tasks/handler"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	lFile, lErr := os.OpenFile(`./log/logfile`+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatalf("error opening file: %v", lErr)
	}
	defer lFile.Close()
	log.SetOutput(lFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	lErr = connectdb.Build_Global_Connections()
	if lErr != nil {
		log.Fatalf("error building global db connections: %v", lErr)
	}

	router := mux.NewRouter()

	router.HandleFunc("/tasks", handler.InterfaceAPForAllTasks).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

	router.HandleFunc("/tasks/{id}", handler.InterfaceAPITasksByID).Methods(http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/login", handler.Login).Methods(http.MethodOptions, http.MethodPost)

	srv := &http.Server{
		Addr:    ":29098",
		Handler: router,
	}
	srv.ListenAndServe()
}
