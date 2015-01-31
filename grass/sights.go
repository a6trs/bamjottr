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
	silvl, err := strconv.Atoi(vars["silvl"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: Check whether the target exists
	// TODO: Map `tgttype` to table names
	sight := &soil.Sight{Account: accountInSession(w, r), Target: tgtid, TableName: "sights_" + tgttype}
	err = sight.Load(soil.KEY_Sight_AccountAndTarget)
	if err != nil {
		sight = &soil.Sight{ID: -1, Account: accountInSession(w, r), Target: tgtid, Level: silvl, TableName: "sights_" + tgttype}
	} else {
		sight.Level = silvl
	}
	if err = sight.Save(soil.KEY_Sight_ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
