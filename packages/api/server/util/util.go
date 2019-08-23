package util

import (
	"encoding/json"
	"github.com/honerlaw/mentordoc/server/model"
	"log"
	"net/http"
	"time"
)

// gets unix time in milliseconds
func NowUnix() int64 {
	return time.Now().UnixNano() / 1e6
}

func WriteJsonToResponse(w http.ResponseWriter, status int, model interface{}) {
	data, err := json.Marshal(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Panic(err)
	}
}

func WriteHttpError(w http.ResponseWriter, err interface{}) {
	if httpError, ok := err.(*model.HttpError); ok {
		WriteJsonToResponse(w, httpError.Status, httpError)
	} else {
		WriteJsonToResponse(w, 500, &model.HttpError{
			Errors: []string{err.(error).Error()},
		})
	}
}