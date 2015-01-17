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
	s := sess.Values["cookie-id"]
	if s == nil {
		s = ""
	}
	renderTemplate(w, "index", map[string]interface{}{"authcookie": s})
}
