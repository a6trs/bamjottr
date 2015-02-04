package grass

import (
	"../soil"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// @url /account_search/invite/{prjid:[0-9]+}/{q}
func AccountSearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Uncomment the following line and import "time" to test spinners on pages locally.
	// time.Sleep(1000 * time.Millisecond)
	// Load the project
	prjid, _ := strconv.Atoi(vars["prjid"])
	prj := &soil.Project{ID: prjid}
	if err := prj.Load(soil.KEY_Project_ID); err != nil {
		fmt.Fprintf(w, `{"error": "Invalid project ID. `+err.Error()+`"}`)
		return
	}
	// Show the list.
	allaccounts, err := soil.FindAccounts(prjid, accountInSession(w, r), vars["q"])
	if err != nil {
		fmt.Fprintf(w, `{"error": "Cannot fetch the list TAT `+err.Error()+`"}`)
		return
	}
	ret, err := json.Marshal(allaccounts)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Cannot convert the results to JSON :X `+err.Error()+`"}`)
		return
	}
	fmt.Fprintf(w, string(ret))
}
