package handlers

import (
	"encoding/json"
	"net/http"
)

func CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode("111")
}