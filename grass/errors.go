package grass

import (
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := sstore.Get(r, "flash")
	flash := sess.Flashes("errmsg")
	msg := ""
	if flash != nil {
		msg = flash[0].(string)
	}
	sess.Save(r, w)
	renderTemplate(w, r, "error", map[string]interface{}{"msg": msg})
}
