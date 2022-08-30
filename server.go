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
	messageBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	var company Company
	err = json.Unmarshal(messageBody, &company)
	if err != nil {
		log.Fatal(err)
	}
	companyList = append(companyList, company)
	fmt.Fprintf(w, "Succesfully added the company: \n"+string(messageBody))
}

// Retrieve all companies
// Http request: GET /companies
func getCompanies(w http.ResponseWriter, req *http.Request) {
	stringToPrint := fmt.Sprintf("%+v", companyList)
	fmt.Fprintf(w, stringToPrint)
}

// Retrieves the details of company with <id>
// Http request: GET /companies/id
func getCompany(w http.ResponseWriter, req *http.Request) {
	// Retrieve the id from URL
	url := req.URL.String()
	id := strings.ReplaceAll(url, "/companies/", "")

	var found bool
	// find id from companyList
	for _, company := range companyList {
		companyId := company.Code
		if strings.EqualFold(id, companyId) {
			fmt.Fprintf(w, "%v", company)
			found = true
		}
	}
	if !found {
		fmt.Fprintf(w, "Could not find the company with identifier "+id)
	}
}

// Update details of company with <id> if it exists
// Http request: PUT /companies/id
func updateCompany(w http.ResponseWriter, req *http.Request) {
	var updatedCompany Company

	// Retrieve the id from URL
	var url string = req.URL.String()
	var id string = strings.ReplaceAll(url, "/companies/", "")

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
	var found bool
	for i, company := range companyList {
		companyId := company.Code
		if companyId == id {
			companyList[i] = updatedCompany
			found = true
			fmt.Fprintf(w, "Succesfully updated company details: \n"+string(requestBody))
		}
	}
	if !found {
		fmt.Fprintf(w, "Not able to update company details")
	}
}

// Remove company with <id>
// Http request: DELETE /companies/id
func deleteCompany(w http.ResponseWriter, req *http.Request) {
	// Retrieve the id from URL
	var url string = req.URL.String()
	var id string = strings.ReplaceAll(url, "/companies/", "")

	var found bool
	var indexToDelete int
	for i, company := range companyList {
		companyId := company.Code
		if companyId == id {
			indexToDelete = i
			found = true
		}
	}
	// Delete element at indexToDelete
	if found {
		var removedCompany Company = companyList[indexToDelete]
		companyList = remove(companyList, indexToDelete)
		b, err := json.Marshal(removedCompany)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Fprintf(w, "Succesfully removed company: \n"+string(b))
		}
	} else {
		fmt.Fprintf(w, "Could not find company with identifier "+id)
	}
}

// Given a list of companies, and in index to remove,
// the following function returns a copy of the original
// with company at the given index removed
func remove(s []Company, i int) []Company {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
