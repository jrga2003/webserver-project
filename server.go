package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Creates a new company
// Http request: POST /companies
func addCompany(w http.ResponseWriter, req *http.Request) {
	if strings.EqualFold(req.Method, "POST") {
		messageBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err) // print error
		}
		var company Company
		err = json.Unmarshal(messageBody, &company)
		if err != nil {
			log.Fatal(err)
		}
		companyList = append(companyList, company)
		fmt.Fprintf(w, "Succesfully added the company")
	}
}

// Retrieve all companies
// Http request: GET /companies
func getCompanies(w http.ResponseWriter, req *http.Request) {
	stringToPrint := fmt.Sprintf("%+v", companyList)
	fmt.Fprintf(w, stringToPrint)
}

// Retrieves the details of company with <id>
// Http request: GET /companies/:id
func getCompany(w http.ResponseWriter, req *http.Request) {
	// Retrieve the id from URL
	url := req.URL.String()
	id := strings.ReplaceAll(url, "/companies/", "")

	// find id from companyList
	for _, company := range companyList {
		companyId := company.Code
		if strings.EqualFold(id, companyId) {
			return fmt.Fprintf(w, "%v", company)
		}
	}
}

// Update details of company with <id> if it exists
// Http request: PUT /companies/:id
func updateCompany(w http.ResponseWriter, req *http.Request) {
	var updatedCompany Company

	// Retrieve the id from URL
	url := req.URL.String()
	id := strings.ReplaceAll(url, "/companies/", "")

	// Retrieve updated company details from body of request message
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(requestBody, &updatedCompany)
	if err != nil {
		log.Fatal(err)
	}

	// find id in companyList
	for i, company := range companyList {
		companyId := company.Code
		if strings.EqualFold(id, companyId) {
			companyList[i] = updatedCompany
			return fmt.Fprintf(w, "Succesfully updated company details")
		}
	}
}

// Remove company with <id>
// Http request: DELETE /companies/:id
func deleteCompany(w http.ResponseWriter, req *http.Request) {
	// Retrieve the id from URL
	url := req.URL.String()
	id := strings.ReplaceAll(url, "/companies/", "")

	var indexToDelete int
	for i, company := range companyList {
		companyId := company.Code
		if strings.EqualFold(id, companyId) {
			indexToDelete = i
		}
	}
	// Delete element at indexToDelete
	companyList = remove(companyList, indexToDelete)
}

// Given a list of companies, and in index to remove,
// the following function removes the company at the given
// index
func remove(s []Company, i int) []Company {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
