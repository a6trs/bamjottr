package grass

import (
	"../soil"
	"net/http"
	"strconv"
)

// @url /sight
func SightHandler(w http.ResponseWriter, r *http.Request) {
	tgttype, err := strconv.Atoi(r.FormValue("tgttype"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tgtid, err := strconv.Atoi(r.FormValue("tgtid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	level, err := strconv.Atoi(r.FormValue("level"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: Move this map (or, array) somewhere else
	tblname := []string{"sights_projects", "sights_posts"}
	sight := &soil.Sight{Account: accountInSession(w, r), Target: tgtid, TableName: tblname[tgttype]}
	err = sight.Load(soil.KEY_Sight_AccountAndTarget)
	if err != nil {
		sight = &soil.Sight{ID: -1, Account: accountInSession(w, r), Target: tgtid, Level: level, TableName: tblname[tgttype]}
	} else {
		sight.Level = level
	}
	if err = sight.Save(soil.KEY_Sight_ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
