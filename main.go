package main

import (
	"log"
	"net/http"
	"os"
	connectdb "tasks/ConnectDB"
	"tasks/apigate"
	"tasks/beequeue"
	"tasks/catching"
	"tasks/constants"
	"tasks/handler"
	"tasks/helpers"
	"time"

	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
)

// main initializes and starts the Tasks API Server.
//
// It performs the following steps:
//   - Initializes a custom logger and sets up file-based logging with timestamps.
//   - Establishes global database connections.
//   - Initializes a Redis cache client.
//   - Starts background workers for Asynq (task queue) and beequeue (custom task sync).
//   - Sets up HTTP routes for task management and login using Gorilla Mux.
//   - Wraps the router with a middleware for API call logging.
//   - Starts the HTTP server on port 29098.
//
// The function will terminate the application if any critical initialization fails.

/*
Package main implements the entry point for the Tasks API Server.

This server provides endpoints for task management and user authentication. It utilizes a PostgreSQL database for persistent storage, Redis for caching and background task processing, and supports both synchronous and asynchronous task handling via Asynq and a custom beequeue implementation.

Key Features:
- RESTful API endpoints for creating, retrieving, updating, and deleting tasks.
- User login endpoint.
- Middleware for logging API calls.
- Background workers for processing asynchronous tasks.
- Structured logging to file with timestamps and source information.

Dependencies:
- github.com/gorilla/mux for HTTP routing.
- github.com/hibiken/asynq for background task processing.
- Custom packages: ConnectDB, apigate, beequeue, catching, constants, handler, helpers.

To run the server:
 1. Ensure PostgreSQL and Redis are running and accessible.
 2. Build and run the application:
    go build -o tasks
    ./tasks

API Endpoints:
- POST/GET    /tasks         : Create or list tasks.
- PUT/GET/DELETE /tasks/{id} : Update, retrieve, or delete a task by ID.
- POST        /login         : User authentication.

The server listens on port 29098 by default.
*/
func main() {
	// lDebug is an instance of HelperStruct from the helpers package, used to provide helper methods or utilities within the application.
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "main(+), Starting Tasks API Server")

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
		return
	}
	catching.CreateCacheClient(lDebug, constants.RedisServerAddress)
	go startAsynqWorker()
	// go beequeue.TaskSync(lDebug)

	router := mux.NewRouter()

	router.HandleFunc("/tasks", handler.InterfaceAPForAllTasks).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

	router.HandleFunc("/tasks/{id}", handler.InterfaceAPITasksByID).Methods(http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/login", handler.Login).Methods(http.MethodOptions, http.MethodPost)

	//Initiate Queue to process API Incoming and outgoing Log
	apigate.ApiCallLogChannel = apigate.InitiateApiCallLog()
	handler := apigate.RequestMiddleWare(router)

	// ln, err := net.Listen("tcp", ":29098")
	// if err != nil {
	// 	log.Fatalf("Failed to listen: %v", err)
	// }

	srv := &http.Server{
		Addr:    ":29098",
		Handler: handler,
	}

	srv.ListenAndServe()
}

// startAsynqWorker initializes and starts an Asynq worker server with specified configuration.
// It sets up a Redis-backed task queue with defined concurrency and queue priorities,
// registers the task handler for the given task identifier, and begins processing tasks.
// If the server encounters an error during startup or execution, it logs the error and terminates the process.
func startAsynqWorker() {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	lDebug.Log(helpers.Statement, "startAsynqWorker(+), Starting Asynq Worker")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: constants.RedisServerAddress},
		asynq.Config{
			Concurrency: 10,
			Queues:      map[string]int{"default": 1},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(constants.UniqueTaskIndentifier, beequeue.HandleSendReminderTask) // make sure this function exists

	log.Println("📦 Asynq Worker started...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Asynq server error: %v", err)
	}
}
