package web

import (
	"net/http"
)

type Options struct {
	Handler http.Handler
}
