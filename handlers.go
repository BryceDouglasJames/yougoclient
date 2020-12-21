package yougoclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestResponseFormat struct {
	Message string `json:"message"`
}

func (h *SearchRequest) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var data SearchRequest

	switch r.Method {
	case "POST":
		//request buffer 100 KB
		r.Body = http.MaxBytesReader(w, r.Body, 100000)

		//create decoder
		decoder := json.NewDecoder(r.Body)

		//STRICT request scope
		decoder.DisallowUnknownFields()

		//init decoder
		err := decoder.Decode(&data)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			fmt.Println()
			return
		}

		//Makes sure there is only ONE json object
		err = decoder.Decode(&struct{}{})
		if err != io.EOF {
			msg := "Request body must only contain a single JSON object"
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		//response write the payload
		//w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.Header().Set("Content-Type", "application/json")
		h.SetRequest(data.ID)
		fmt.Fprintln(w, h.GetRequest())

	case "GET":
		//create response payload, post to page
		response, err := json.Marshal(h)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			fmt.Println()
			return
		}

		//w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		json.NewEncoder(w).Encode(h.GetRequest())
		fmt.Fprint(w)

		w.Write(response)
	}
}

func (h *Users) ServeArray(w http.ResponseWriter, r *http.Request) {

	response, err := json.Marshal(h)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		fmt.Println()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(h.Searches)
	//fmt.Fprint(w)
	w.Write(response)
}
