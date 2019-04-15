// Headers func and struct
package api

import (
	"encoding/json"
	"github.com/Atluss/TestTaskElma/lib"
	"net/http"
)

const (
	V1Api = "v1"
)

type HeadRequest interface {
	Request() // execute FetchTask
}

type RequestRest struct {
	HeadRequest
	w *http.ResponseWriter
	r *http.Request
}

type ReplayStatus struct {
	Status      int    `json:"Status"`
	Description string `json:"Description"`
}

func (obj *ReplayStatus) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&obj)
	if !lib.LogOnError(err, "error: can't encode ReplayStatus") {
		return err
	}
	return nil
}

func SetDefaultHeadersHttp(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
}

func SetDefaultHeadersV1API(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json")
}
