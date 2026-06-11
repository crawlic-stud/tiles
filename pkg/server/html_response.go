package server

import (
	"fmt"
	"net/http"
	"tiles/templates"

	"github.com/a-h/templ"
)

type HTMLResponse struct {
	StatusCode int
	Component  templ.Component
}

func (r HTMLResponse) Write(w http.ResponseWriter, req *http.Request) {
	status := r.StatusCode
	if status == 0 {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	templ.Handler(
		r.Component,
		templ.WithStreaming(),
	).ServeHTTP(w, req)
}

func HTML(component templ.Component) Response {
	return HTMLResponse{
		Component: component,
	}
}

func HTMLErrorf(status int, message string, args ...any) Response {
	return HTMLResponse{
		StatusCode: status,
		Component:  templates.ErrorPage(status, fmt.Sprintf(message, args...)),
	}
}

func HTMLError(status int, message string) Response {
	return HTMLResponse{
		StatusCode: status,
		Component:  templates.ErrorPage(status, message),
	}
}
