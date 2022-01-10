package ldap

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var members = allMember{
	{ID: 11, Name: "Zahphod Beeblebrox", DN: "CN=Zahphod Beeblebrox,OU=People,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 12, Name: "Arthur Dent", DN: "CN=Arthur Dent,OU=People,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 13, Name: "Ford Prefect", DN: "CN=Ford Prefect,OU=People,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 14, Name: "Tricia McMillan", DN: "CN=Tricia McMillan,OU=People,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 15, Name: "Marvin", DN: "CN=Marvin,OU=People,DC=lübthen,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 16, Name: "Darth Vader", DN: "CN=Darth Vader,OU=People,DC=redefin,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 17, Name: "Anakyn Skywalker", DN: "CN=Anakyn Skywalker,OU=People,DC=redefin,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 18, Name: "Luke Skywalker", DN: "CN=Luke Skywalker,OU=People,DC=redefin,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 19, Name: "Meister Joda", DN: "CN=Meister Joda,OU=People,DC=redefin,DC=,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
	{ID: 20, Name: "Han Solo", DN: "CN=Han Solo,OU=People,DC=redefin,DC=hagenow,DC=parchim,DC=elkm,DC=de"},
}

type allMember []Member

func getOneMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	memberID := mux.Vars(r)["id"]

	for _, member := range members {
		smemberid, _ := strconv.Atoi(memberID)
		if member.ID == smemberid {
			json.NewEncoder(w).Encode(member)
		}
	}
}

func getAllMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(members)
}

func searchMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pattern := r.FormValue("pattern")
	var result []Member
	for _, member := range members {
		if strings.Contains(member.DN, pattern) {
			result = append(result, member)

		}
	}
	json.NewEncoder(w).Encode(result)
}
