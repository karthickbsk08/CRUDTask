package apigate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	ratelimiter "tasks/apigate/rateLimiter"
	"tasks/constants"
	"tasks/helpers"
	"time"
)

func RequestMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Step 1 : Debugging instance and new session id creation
		lDebug := new(helpers.HelperStruct)
		lDebug.Init()
		lDebug.Log(helpers.Statement, "RequestMiddleWare(+)", lDebug.Sid)

		//Step 2 : `Rate Limiting Logic (Token Bucket Algorithm)
		lrateLimiter := ratelimiter.AssignRateLimitValue()

		if !lrateLimiter.Allow() {
			lMessage := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&lMessage)
			return
		}

		//Step 3 : Request time capture
		RequestUnixTime := time.Now().UnixMilli()
		ReqDateTime := time.Now()

		//Step 4 : request instance reading
		requestorDetail := GetRequestorDetail(r)

		//Step 5 : Prepare instance of response writer which will be return once the api writered its response
		captureWriter := &ResponseCaptureWriter{ResponseWriter: w}

		//Step 6 : CORS Headers checking
		// Uncomment the following lines if you need to handle these common headers.
		// Set up CROS credentails
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Credentials", "true")
		(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Set up the Header API for common
		// If you need additional header you need to config here itself
		(w).Header().Set("Access-Control-Allow-Headers", "sid,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

		//Step 6 : Allow Options method directly with status 200
		if strings.EqualFold("OPTIONS", r.Method) {
			captureWriter.WriteHeader(http.StatusOK)
			return
		}

		//Step 7 : Check Session id given in the request else proceed to create a new once and use it rest of process
		lSid := r.Header.Get("sid") // session id
		// If sid present in API request HEADER then reassign the value
		if !strings.EqualFold(lSid, "") {
			lDebug.Sid = lSid
		}

		//Step 8 :  Create context instance which helps to pass the created Session id throughout the process (Middleware --> Process API --> Return to Middleware)
		ctx := context.WithValue(r.Context(), helpers.RequestIDKey, lDebug.Sid)
		r = r.WithContext(ctx)

		//Step 9 : Only Read the body if the PUT,POST methods for others skip
		if strings.Contains(constants.MethodsWithBody, r.Method) {

			//reading body and reinitialize in request itself to read inside the respective api's
			lBody, lErr := io.ReadAll(r.Body)
			if lErr != nil {
				lDebug.Log(helpers.Elog, "Error @ body reading  in middleware "+lErr.Error())
			} else {
				requestorDetail.Body = string(lBody)
				r.Body = io.NopCloser(bytes.NewBuffer(lBody))
				lDebug.Log(helpers.Details, fmt.Sprintf("request body of api %s : %s", requestorDetail.EndPoint, requestorDetail.Body))
			}

		}

		// Step 10 : serve the handler to the respective endpoint
		next.ServeHTTP(captureWriter, r)

		//Step 11 : Capture response time and prepare the logging mechanism
		ResponseUnixTime := time.Now().UnixMilli()
		RespDateTime := time.Now()
		RespStatus := captureWriter.Status()
		RespData := captureWriter.Body()

		LogEntry := ApiLogCapture{
			RequestId:        lDebug.Sid,
			ReqDateTime:      ReqDateTime,
			RequestUnixTime:  RequestUnixTime,
			RespDateTime:     RespDateTime,
			ResponseUnixTime: ResponseUnixTime,
			RealIP:           requestorDetail.RealIP,
			ForwardedIP:      requestorDetail.ForwardedIP,
			Method:           requestorDetail.Method,
			Path:             requestorDetail.Path,
			Host:             requestorDetail.Host,
			RemoteAddr:       requestorDetail.RemoteAddr,
			Header:           requestorDetail.Header,
			Endpoint:         requestorDetail.EndPoint,
			ReqBody:          requestorDetail.Body,
			RespBody:         string(RespData),
			ResponseStatus:   int64(RespStatus),
			PDebug:           lDebug,
		}
		ApiCallLogChannel <- LogEntry

		lDebug.Log(helpers.Statement, "RequestMiddleWare(-)", lDebug.Sid)

	})
}
