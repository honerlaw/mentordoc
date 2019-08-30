package util

import (
	"encoding/json"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// gets unix time in milliseconds
func NowUnix() int64 {
	return time.Now().UnixNano() / 1e6
}

func WriteJsonToResponse(w http.ResponseWriter, status int, model interface{}) {
	var data []byte
	if model != nil {
		byteArr, err := json.Marshal(model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data = byteArr
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		log.Panic(err)
	}
}

func WriteHttpError(w http.ResponseWriter, err interface{}) {
	if httpError, ok := err.(*shared.HttpError); ok {
		WriteJsonToResponse(w, httpError.Status, httpError)
	} else {
		WriteJsonToResponse(w, 500, &shared.HttpError{
			Errors: []string{err.(error).Error()},
		})
	}
}

func BuildSqlPlaceholderArray(slice interface{}) string {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("a slice is required")
	}

	// convert the passed models to an interface array so we can work with it...
	placeholders := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		placeholders[i] = "?"
	}

	return strings.Join(placeholders, ", ")
}

func ConvertStringArrayToInterfaceArray(slice []string) []interface{} {
	newSlice := make([]interface{}, len(slice))
	for i := 0; i < len(slice); i++ {
		newSlice[i] = slice[i]
	}
	return newSlice
}