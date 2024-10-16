package tasks

import (
	"database/sql"
	"net/http"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
