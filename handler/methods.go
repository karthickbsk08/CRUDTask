package handler

import (
	"net/http"
)

func InterfaceAPForAllTasks(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		CreateTask(w, r)
	case http.MethodGet:
		GetAllTasks(w, r)
	}
}

func InterfaceAPITasksByID(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPut:
		UpdateTaskByAPI(w, r)
	case http.MethodGet:
		GetTaskByID_API(w, r)
	case http.MethodDelete:
		DeleteTaskByID_API(w, r)
	}
}
