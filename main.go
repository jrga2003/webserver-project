package main

import (
	"net/http"
)

func main() {
	http.ListenAndServe(":8090", nil)
}

func handleCompany(req http.Request) {
	url := req.URL.String()
	if url == "/companies" {
		if req.Method == "POST" {
			http.HandleFunc(url, addCompany)
		} else {
			http.HandleFunc(url, getCompanies)
		}
	} else {
		if req.Method == "" { // For client request "" means GET
			http.HandleFunc(url, getCompany)
		} else if req.Method == "PUT" {
			http.HandleFunc(url, updateCompany)
		} else {
			http.HandleFunc(url, deleteCompany)
		}
	}
}
