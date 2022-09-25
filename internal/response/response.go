package response

import (
	"encoding/json"
	"net/http"
)

type StatusType string

const (
	Ok      StatusType = "ok"
	Error              = "error"
	Warning            = "warning"
)

type Response struct {
	Status  StatusType `json:"status"`
	Message string     `json:"message"`
	Data    any        `json:"data"`
}

func WriteResponse(w http.ResponseWriter, e Response) {
	json.NewEncoder(w).Encode(e)
}
