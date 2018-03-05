package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/Polarishq/middleware/framework/log"
)

// SwaggerUIPath is the static URL prefix from which the Swagger UI is served
const SwaggerUIPath = "/swagger-ui"

// SwaggerUIHandler provides a middleware handler which serves the swagger UI
type SwaggerUIHandler struct {
	handler http.Handler
	dir     string
}

// NewSwaggerUIHandler creates a new middleware handler which serves the swagger UI
func NewSwaggerUIHandler(dir string, handler http.Handler) *SwaggerUIHandler {
	return &SwaggerUIHandler{dir: dir, handler: handler}
}

func (s *SwaggerUIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.URL.Path, SwaggerUIPath) == 0 {
		if _, err := os.Stat(s.dir); os.IsNotExist(err) {
			http.Error(w, s.dir+" is not found", http.StatusNotFound)
		}

		http.StripPrefix("", http.FileServer(http.Dir(s.dir))).ServeHTTP(w, r)
		log.Debug("request served by SwaggerUIHandler handler")
		return
	}

	s.handler.ServeHTTP(w, r)
}
