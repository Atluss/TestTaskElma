// Headers func and struct
package api

import (
	"encoding/json"
	"github.com/Atluss/TestTaskElma/pkg/v1"
	"net/http"
)

const (
	V1Api = "v1"
)

// HeadRequest for api request
type HeadRequest interface {
	Request() // execute FetchTask
}

// RequestRest http request
type RequestRest struct {
	HeadRequest
	w *http.ResponseWriter
	r *http.Request
}

// ReplayStatus struct for answer request
type ReplayStatus struct {
	Status      int    `json:"Status"`
	Description string `json:"Description"`
}

// Encode to json
func (obj *ReplayStatus) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&obj)
	if !v1.LogOnError(err, "error: can't encode ReplayStatus") {
		return err
	}
	return nil
}

// SetDefaultHeadersHttp for http connections
func SetDefaultHeadersHttp(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
}

// SetDefaultHeadersV1API for json answers
func SetDefaultHeadersV1API(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json")
}
