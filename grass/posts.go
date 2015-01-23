package grass

import (
	"net/http"
	"fmt"
)

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "post_create", map[string]interface{}{"aid": accountInSession(w, r)})
	} else {
		title := r.FormValue("title")
		body := r.FormValue("body")
		fmt.Println(title)
		fmt.Println(body)
	}
}
