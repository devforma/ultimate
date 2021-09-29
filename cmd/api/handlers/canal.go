package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devforma/ultimate/internal/database"
	"github.com/devforma/ultimate/internal/util"
	"github.com/go-playground/validator/v10"
)

type Canal struct {
	ID       int    `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	Duration int    `db:"duration" json:"duration"`
}

type NewCanal struct {
	Title    string `json:"title" validate:"required,max=6" `
	Duration int    `json:"duration" validate:"required,min=12,max=24"`
}

type CanalAPI struct {
	DB *database.DB
}

func (c CanalAPI) NewCanal(w http.ResponseWriter, r *http.Request) {
	var nc NewCanal

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(util.StringToBytes("read request body failed"))
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &nc)
	if err != nil {
		w.Write(util.StringToBytes("unmarshal request body failed"))
		return
	}

	vali := validator.New()
	err = vali.Struct(&nc)
	if err != nil {
		errStrings := processErr(err)
		w.Write(util.StringToBytes(strings.Join(errStrings, "\n")))
		return
	}
}

func processErr(err error) []string {
	var errorStrings []string
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		errorStrings = append(errorStrings, invalid.Error())
		return errorStrings
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		for idx := range validationErrors {
			ers := validationErrors[idx].Field() + validationErrors[idx].Param() + validationErrors[idx].Tag()
			errorStrings = append(errorStrings, ers)
		}
	}

	return errorStrings
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
