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
// NOTE: Client should be able to filter list response by 
// company properties using request query 
func getCompanies(w http.ResponseWriter, req *http.Request) {
	var query string = req.URL.Query().Encode()
	var propertyFilters map[string]string = parseQuery(query)

	if propertyFilters == nil {
		stringToPrint := fmt.Sprintf("%+v", companyList)
		fmt.Fprintf(w, stringToPrint)
	} else {
		var companiesToReturn []Company
		for filter, value := range propertyFilters {
			for _, company := range companyList {
				switch filter {
				case "Name":
					if company.Name == value {
						companiesToReturn = append(companiesToReturn, company)
					}
				case "Code":
					if company.Code == value {
						companiesToReturn = append(companiesToReturn, company)
					}
				case "Country":
					if company.Country == value {
						companiesToReturn = append(companiesToReturn, company)
					}
				case "Website":
					if company.Website == value {
						companiesToReturn = append(companiesToReturn, company)
					}
				case "Phone":
					if company.Phone == value {
						companiesToReturn = append(companiesToReturn, company)
					}
				default:
					fmt.Fprintf(w, "The property " + filter + " does not exist for the entity Company")
				}
			}
		}
	}
}

// The following is a helper function for the getCompanies function
// above. It is used to convert the query as part of a http request
// to a map which matches company properties to values.
// E.g.
// q: "Code=1&Website=HomeDepot" --> map["Code":"1" "Website":"HomeDepot"]
func parseQuery(q string) map[string]string {
	var chars []rune = []rune(q)
	var m map[string]string = make(map[string]string)

	for i, letter := range chars {
		if letter == '=' {
			var property string =  // traverse backwards until you find ampersand character
			var value string = // traverse forwards until you find ampersand character
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
