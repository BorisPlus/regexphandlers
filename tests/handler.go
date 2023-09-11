package regexhandlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type DefaultHandler struct{}

func (h DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/version", http.StatusTemporaryRedirect)
}

type VersionHandler struct{}

func (h VersionHandler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
	_, err := response.Write([]byte("1.0.0"))
	if err != nil {
		panic(fmt.Sprintf("VersionHandler err %s\n", err))
	}
}

type GetHandler struct{}

func (h GetHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	responsed := make(map[string]string)
	responsed["parent_id"] = request.Form.Get("parent_id")
	responsed["child_name"] = request.Form.Get("child_name")
	jsonResponsed, err := json.Marshal(responsed)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(jsonResponsed)
	if err != nil {
		panic(fmt.Sprintf("GetHandler err %s\n", err))
	}
}
