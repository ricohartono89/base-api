package handler

import (
	jsonEncoding "encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/ricohartono89/base-api/internal/json"
)

// DefaultDecoder ...
var DefaultDecoder = schema.NewDecoder()

// EndpointHandler ...
type EndpointHandler func(http.ResponseWriter, *http.Request) json.HTTPResponse

func (fn EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	res := fn(w, r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.GetStatus())

	if res.HasError() {
		handleErrorResponse(w, r, res)
	} else {
		handleOKResponse(w, res)
	}
}

func handleOKResponse(w http.ResponseWriter, res json.HTTPResponse) {
	data := res.GetData()
	if data == nil {
		data = map[string]string{
			"message": "Success",
		}
	}
	encodeResponse(w, data)
}

func handleErrorResponse(w http.ResponseWriter, r *http.Request, res json.HTTPResponse) {
	log.Println(res.GetErrorMessageVerbose())

	data := map[string]interface{}{
		"code":    res.GetErrCode(),
		"message": res.GetErrorMessage(),
	}
	encodeResponse(w, data)
}

func encodeResponse(w http.ResponseWriter, data interface{}) {
	err := jsonEncoding.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error encode:", err.Error())
		http.Error(w, "Error encode response", http.StatusInternalServerError)
	}
}
