package apigate

import (
	"net/http"
	"tasks/helpers"
	"time"
)

type ApiLogCapture struct {
	RequestId        string                `gorm:"column:request_id"`
	RespBody         string                `gorm:"column:respbody"`
	ResponseStatus   int64                 `gorm:"column:responseStatus"`
	ReqDateTime      time.Time             `gorm:"column:reqDateTime"`
	RealIP           string                `gorm:"column:realip"`
	ForwardedIP      string                `gorm:"column:forwardedip"`
	Method           string                `gorm:"column:method"`
	Path             string                `gorm:"column:path"`
	Host             string                `gorm:"column:host"`
	RemoteAddr       string                `gorm:"column:remoteaddr"`
	Header           string                `gorm:"column:header"`
	Endpoint         string                `gorm:"column:endpoint"`
	RespDateTime     time.Time             `gorm:"column:respDateTime"`
	ReqBody          string                `gorm:"column:reqbody"`
	RequestUnixTime  int64                 `gorm:"column:requesttime"`
	ResponseUnixTime int64                 `gorm:"column:responsetime"`
	PDebug           *helpers.HelperStruct `gorm:"-"`
}

type RequestorDetails struct {
	RealIP      string
	ForwardedIP string
	Method      string
	Path        string
	Host        string
	RemoteAddr  string
	Header      string
	Body        string
	EndPoint    string
	RequestType string
}

type ResponseCaptureWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

type HeaderDetails struct {
	Key   string
	Value string
}

type ResponseStruct struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}
