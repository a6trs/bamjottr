package soil

import (
	"strconv"
	"strings"
)

// The abstract method for finding what 'you may also like'.
// e.g. Find 3 most gazed projects by whom gazed project #3:
//   Recommend(1, "projects", 3)
func Recommend(fromid int, table string, limit int) []int {
	var list []int
	// TODO: Improve this algorithm whenever possible
	// Here we just retrieve all the people whose sight level for this project
	//   is not zero and find what else they stared at (or you can say starred)
	rs1, err := db.Query(`SELECT account FROM sights_`+table+` WHERE target = ? AND level <> 0`, fromid)
	if err != nil {
		return nil
	}
	defer rs1.Close()
	gazers := []string{}
	for rs1.Next() {
		var a int
		if rs1.Scan(&a) == nil {
			gazers = append(gazers, strconv.Itoa(a))
		}
	}
	// stackoverflow.com/q/1503959
	rs2, err := db.Query(`SELECT target FROM sights_`+table+` WHERE account IN (`+strings.Join(gazers, ",")+`) AND target <> ? GROUP BY target ORDER BY count(*) DESC LIMIT ?`, fromid, limit)
	if err != nil {
		return nil
	}
	defer rs2.Close()
	for rs2.Next() {
		var a int
		if rs2.Scan(&a) == nil {
			list = append(list, a)
		}
	}
	return list
}
