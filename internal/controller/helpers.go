package controller

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/maknahar/go-web-skeleton/internal/utils/logger"
	"github.com/maknahar/go-web-skeleton/public/dtos"
)

func LogAndWriteResponse(w http.ResponseWriter, r *http.Request, h func(l *log.Entry) ([]byte, int)) {
	start := time.Now()
	l := logger.GetLogger()

	response, statusCode := h(l)

	if statusCode >= 400 {
		response = jsonifyErrMessage(string(response), statusCode, nil)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(response)
	if err != nil {
		l.WithError(err).Error("Could not encode response")
	}
	l.Printf("Status Code: %d, Response time: %v, Response: %s", statusCode, time.Since(start), string(response))

}

func jsonifyErrMessage(errMsg string, code int, err error) []byte {
	errObj := new(dtos.ErrorResponseDTO)
	errObj.Message = errMsg
	if err != nil {
		errObj.Message += " " + err.Error()
	}

	errObj.Code = code
	errObj.Status = "ERROR"
	data, _ := json.Marshal(errObj)
	return data
}
