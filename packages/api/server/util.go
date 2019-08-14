package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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