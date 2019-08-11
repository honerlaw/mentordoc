package util

import (
	"encoding/json"
	"net/http"
	"time"
)

func NowUnix() int64 {
	return time.Now().UnixNano() / 1e6
}

func WriteJsonToResponse(w http.ResponseWriter, model interface{}) {
	data, err := json.Marshal(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}