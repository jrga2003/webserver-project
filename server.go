package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Creates a new company
// Http request: POST /companies
func addCompany(w http.ResponseWriter, req *http.Request) {
	messageBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "addCompany: "+err.Error())
	}
	var company Company
	err = json.Unmarshal(messageBody, &company)
	if err != nil {
		fmt.Fprintf(w, "addCompany: "+err.Error())
	}
	companyList = append(companyList, company)
	fmt.Fprintf(w, "Succesfully added the company:\n"+string(messageBody))
}

// Retrieve all companies
// Http request: GET /companies
// NOTE: Client should be able to filter list response by
// company properties using request query
func getCompanies(w http.ResponseWriter, req *http.Request) {
	var query string = req.URL.Query().Encode()
	var propertyFilters map[string]string = parseQuery(query)
	var companiesToReturn []Company = companyList

	for filter, value := range propertyFilters {
		var updatedCompanies []Company
		for _, company := range companiesToReturn {
			switch filter {
			case "Name":
				if company.Name == value {
					updatedCompanies = append(updatedCompanies, company)
				}
			case "Code":
				if company.Code == value {
					updatedCompanies = append(updatedCompanies, company)
				}
			case "Country":
				if company.Country == value {
					updatedCompanies = append(updatedCompanies, company)
				}
			case "Website":
				if company.Website == value {
					updatedCompanies = append(updatedCompanies, company)
				}
			case "Phone":
				if company.Phone == value {
					updatedCompanies = append(updatedCompanies, company)
				}
			default:
				//fmt.Fprintf(w, filter+" : "+value+"\n")
			}
		}
		companiesToReturn = updatedCompanies
	}
	if len(companiesToReturn) == 0 {
		fmt.Fprintf(w, "Found no companies\n")
	} else {
		fmt.Fprintf(w, "Companies found: %+v\n", companiesToReturn)
	}
}

// The following is a helper function for the getCompanies function
// above. It is used to convert the query as part of a http request
// to a map which matches company properties to values.
// E.g.
// q: "Code=1&Website=HomeDepot" --> map["Code":"1" "Website":"HomeDepot"]
// q: "Name=R" --> map["Name":"R"]
// q: "Country=L&Website=randomwebsite.com" --> map["Country":"L" "Website":"randomwebsite.com"]
// q: "Name=Penguin&Website=penguin.com" --> map["Name":"Penguin" "Website":"penguin.com"]
func parseQuery(q string) map[string]string {
	var chars []rune = []rune(q)
	var m map[string]string = make(map[string]string)

	for index, letter := range chars {
		if letter == '=' {
			var property string // traverse backwards until you find ampersand character
			for i := index - 1; i >= 0; i-- {
				if i == 0 {
					property = string(chars[i:index])
					break
				}
				if chars[i] == '&' {
					property = string(chars[i+1 : index])
					break
				}
			}

			var value string // traverse forwards until you find ampersand character
			for i := index + 1; i < len(chars); i++ {
				if chars[i] == '&' {
					value = string(chars[index+1 : i])
					break
				}
				if i == len(chars)-1 {
					value = string(chars[index+1 : i+1])
					break
				}
			}
			m[property] = value
		}
	}
	return m
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
		if id == companyId {
			fmt.Fprintf(w, "Company found: %v\n", company)
			found = true
		}
	}
	if !found {
		fmt.Fprintf(w, "getCompany: Could not find the company with identifier "+id)
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
		fmt.Fprintf(w, err.Error())
	}
	err = json.Unmarshal(requestBody, &updatedCompany)
	if err != nil {
		fmt.Fprintf(w, err.Error())
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
		fmt.Fprintf(w, "updateCompany: Could not find the company with the id "+id)
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
			fmt.Fprintf(w, "deleteCompany: "+err.Error())
		} else {
			fmt.Fprintf(w, "Succesfully removed company: \n"+string(b))
		}
	} else {
		fmt.Fprintf(w, "deleteCompany: Could not find company with identifier "+id)
	}
}

// Given a list of companies, and in index to remove,
// the following function returns a copy of the original
// with company at the given index removed
func remove(s []Company, i int) []Company {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
