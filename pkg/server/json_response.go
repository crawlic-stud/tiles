package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	StatusCode int
	Data       any
	Error      string
}

func (r JSONResponse) Write(w http.ResponseWriter, _ *http.Request) {
	status := r.StatusCode
	if status == 0 {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if r.Error != "" {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": r.Error,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(r.Data)
}

func JSON(status int, data any) Response {
	return JSONResponse{
		StatusCode: status,
		Data:       data,
	}
}

func JSONErrorf(status int, err string, args ...any) Response {
	return JSONResponse{
		StatusCode: status,
		Error:      fmt.Sprintf(err, args...),
	}
}

func JSONError(status int, err string) Response {
	return JSONResponse{
		StatusCode: status,
		Error:      err,
	}
}
