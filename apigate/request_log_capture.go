package apigate

var ApiCallLogChannel chan<- ApiLogCapture

// // This method is used to initiate the channel for capture the logs
// func InitiateApiCallLog() chan<- ApiLogCapture {
// 	LogRespChannel := make(chan ApiLogCapture, 100)
// 	go func() {
// 		count := 0
// 		for value := range LogRespChannel {
// 			count++
// 			// ApiCallLogCapture(count, value)
// 			//time.Sleep(time.Millisecond * 500)
// 		}
// 	}()
// 	return LogRespChannel
// }

// func ApiCallLogCapture(count int, pRequest ApiLogCapture) {
// 	pRequest.PDebug.Log(helpers.Statement, "ApiCallLogCapture (+)", count, ftdb.GDB.Maria.Stats())
// 	pRequest.PDebug.Log(helpers.Details, pRequest, "pRequest")
// 	// Insert the record into the database using GORM
// 	// Replace the <table_name>
// 	lResult := ftdb.GDB.GRMTxlSrvMaria.Table("xxapi_req_log").Create(&pRequest)
// 	if lResult.Error != nil {
// 		pRequest.PDebug.Log(helpers.Elog, "LogCapture Error", lResult.Error.Error())
// 	}

// 	pRequest.PDebug.Log(helpers.Statement, "ApiCallLogCapture (-)")

// }
