package yougoclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
		fmt.Fprintln(w, data)
		h.ID = data.ID

	case "GET":
		//create response payload, post to page
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.ID)
		fmt.Fprint(w)
	}
}

func ServeStaticSite(w http.ResponseWriter, r *http.Request) {

}
