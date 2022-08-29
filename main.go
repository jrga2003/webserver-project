package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/companies", handleCompanies)
	http.HandleFunc("/companies/", handleCompany)
	http.ListenAndServe(":8090", nil)
}

func handleCompanies(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		addCompany(w, req)
	} else {
		getCompanies(w, req)
	}
}

func handleCompany(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" { // For client request "" means GET
		getCompany(w, req)
	} else if req.Method == "PUT" {
		updateCompany(w, req)
	} else {
		deleteCompany(w, req)
	}
}
