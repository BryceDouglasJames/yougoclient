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

func (h *Users) SearchHandler(w http.ResponseWriter, r *http.Request) {
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

		CurrentSearch = data.ID

		finished := make(chan bool)
		go worker(finished)
		<-finished
		fmt.Println("Main: Completed")

		//w.Header().Set("Content-Type", "application/json")
		//fmt.Fprintln(w, data.ID)

	case "GET":
		//create response payload, post to page
		response, err := json.Marshal(h.Searches)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			fmt.Println()
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		json.NewEncoder(w).Encode(response)
		w.Write(response)
	}
}

func (h *Users) AddUser(w http.ResponseWriter, r *http.Request) {
	var data UserAdd
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

		for _, user := range ClientList {
			if user.UserName == data.ID {
				w.WriteHeader(401)
				w.Write([]byte("User exists"))
				fmt.Println()
				return
			}
		}

		finished := make(chan bool)
		go adduser(finished, data.ID)
		<-finished

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, data.ID)

	case "GET":
		response, err := json.Marshal(ClientList)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			fmt.Println()
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		json.NewEncoder(w).Encode(response)
		w.Write(response)
	}

}

func (h *Users) ServeArray(w http.ResponseWriter, r *http.Request) {

	response, err := json.Marshal(h.Searches)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		fmt.Println()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
