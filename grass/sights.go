package grass

import (
	"../soil"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func SightHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tgttype := vars["tgttype"]
	tgtid, err := strconv.Atoi(vars["tgtid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: Check whether the target exists
	// TODO: Map `tgttype` to table names
	sight := &soil.Sight{ID: -1, Account: accountInSession(w, r), Target: tgtid, TableName: tgttype}
	if err := sight.Save(soil.KEY_Sight_ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
