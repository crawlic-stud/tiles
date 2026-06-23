package server

import (
	"fmt"
	"net/http"
	"tiles/pkg/helpers"
	"tiles/templates"

	"github.com/a-h/templ"
)

type HTMLResponse struct {
	StatusCode int
	Component  templ.Component
}

func (r HTMLResponse) Write(w http.ResponseWriter, req *http.Request) {
	handler := templ.Handler(
		r.Component,
		templ.WithStreaming(),
	)
	handler.Status = helpers.Ternary(r.StatusCode == 0, int(http.StatusOK), r.StatusCode)
	handler.ServeHTTP(w, req)
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
