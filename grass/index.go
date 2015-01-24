package grass

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "index", map[string]interface{}{})
}
