package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/devforma/ultimate/internal/database"
)

type Canal struct {
	ID       int    `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	Duration int    `db:"duration" json:"duration"`
}

type CanalAPI struct {
	DB *database.DB
}

func (c CanalAPI) SingleCanal(w http.ResponseWriter, r *http.Request) {
	var canal Canal
	c.DB.Get(&canal, "SELECT * FROM `canal` WHERE `id`=?", r.URL.Query().Get("id"))

	time.Sleep(8 * time.Second)

	if data, err := json.Marshal(canal); err == nil {
		w.Write(data)
	}
}

func (c CanalAPI) ListCanal(w http.ResponseWriter, r *http.Request) {
	var canals []Canal
	c.DB.Select(&canals, "SELECT * FROM `canal`")
	if data, err := json.Marshal(canals); err == nil {
		w.Write(data)
	}
}
