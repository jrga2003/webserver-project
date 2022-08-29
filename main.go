package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/companies", handleCompany)
	http.ListenAndServe(":8090", nil)
}

func handleCompany(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	if url == "/companies" {
		if req.Method == "POST" {
			addCompany(w, req)
		} else {
			getCompanies(w, req)
		}
	} else {
		if req.Method == "" { // For client request "" means GET
			getCompany(w, req)
		} else if req.Method == "PUT" {
			updateCompany(w, req)
		} else {
			deleteCompany(w, req)
		}
	}
}
