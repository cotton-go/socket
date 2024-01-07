package http

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
)

type H map[string]any

func JSON(w http.ResponseWriter, code int, data any) error {
	w.WriteHeader(code)

	r := render.JSON{Data: data}
	if !bodyAllowedForStatus(code) {
		r.WriteContentType(w)
		return nil
	}

	if err := r.Render(w); err != nil {
		return err
	}

	return nil
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
