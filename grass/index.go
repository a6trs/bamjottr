package grass

import (
	_ "../soil"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}
