package grass

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := sstore.Get(r, "account-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s := sess.Values["id"]
	if s == nil {
		s = -1
	}
	renderTemplate(w, "index", map[string]interface{}{"aid": s})
}
