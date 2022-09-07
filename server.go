package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Creates a new company
// Http request: POST /companies
func addCompany(w http.ResponseWriter, req *http.Request) {
	// Decode JSON in message body into instance of type company
	messageBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "addCompany: %v", err)
	}
	var company Company
	err = json.Unmarshal(messageBody, &company)
	if err != nil {
		fmt.Fprintf(w, "addCompany: %v", err)
	}

	// add company to databases
	result, err := db.Exec("INSERT INTO companies (name, country, website, phone) VALUES (?, ?, ?, ?)", company.Name, company.Country, company.Website, company.Phone)
	if err != nil {
		fmt.Fprintf(w, "addCompany: %v", err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Fprintf(w, "addCompany: %v", err)
	}
	fmt.Fprintf(w, "Added the company with id no. %v:\n"+string(messageBody), id)
}

// Retrieve all companies
// Http request: GET /companies
// NOTE: Client should be able to filter list response by
// company properties using request query
func getCompanies(w http.ResponseWriter, req *http.Request) {
	var query string = req.URL.Query().Encode()
	var propertyFilters map[string]string = parseQuery(query)
	var filteredCompanies []Company // Company slice to hold data returned from rows

	if len(propertyFilters) == 0 {
		rows, err := db.Query("SELECT * FROM companies")
		if err != nil {
			fmt.Fprintf(w, "getCompanies: %v", err)
		}
		defer rows.Close()
		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var company Company
			if err := rows.Scan(&company.Code, &company.Name, &company.Country, &company.Website, &company.Phone); err != nil {
				fmt.Fprintf(w, "getCompanies: %v", err)
			}
			filteredCompanies = append(filteredCompanies, company)
		}
		// Format output into JSON
		b, err := json.Marshal(filteredCompanies)
		if err != nil {
			fmt.Fprintf(w, "getCompanies: %v", err)
		}
		fmt.Fprintf(w, "Companies found:\n%v", string(b))
	} else {
		var properties []string = make([]string, len(propertyFilters))
		var values []string = make([]string, len(propertyFilters))

		var index int
		index = 0
		for k := range propertyFilters {
			properties[index] = k
			index++
		}

		index = 0
		for _, v := range propertyFilters {
			values[index] = v
			index++
		}
		var condition string = ""
		for i := 0; i < len(properties); i++ {
			condition += properties[i] + "='" + values[i] + "'"
			if i != len(properties)-1 {
				condition += " AND "
			}
		}

		rows, err := db.Query("SELECT * FROM companies WHERE " + condition)
		if rows == nil { // Check if rows are empty to avoid nil pointer derefrence
			fmt.Fprintf(w, "getCompanies: Could not find any companies satisfying the query, "+condition)
			return
		}
		if err != nil {
			fmt.Fprintf(w, "getCompanies: %v", err)
			return
		}
		defer rows.Close()
		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var company Company
			if err := rows.Scan(&company.Name, &company.Code, &company.Country, &company.Website, &company.Phone); err != nil {
				fmt.Fprintf(w, "getCompanies: %v", err)
			}
			filteredCompanies = append(filteredCompanies, company)
		}
		// Format output into JSON
		b, err := json.Marshal(filteredCompanies)
		if err != nil {
			fmt.Fprintf(w, "getCompanies: %v", err)
		}
		fmt.Fprintf(w, "Companies found:\n%v", string(b))
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
	stringId := strings.ReplaceAll(url, "/companies/", "")
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "getCompany: %v", err)
	}

	var company Company // A company to hold data from returned row
	row := db.QueryRow("SELECT * FROM companies WHERE id = ?", id)
	if err := row.Scan(&company.Code, &company.Name, &company.Country, &company.Website, &company.Phone); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "getCompany %v: no such company", id)
		} else {
			fmt.Fprintf(w, "getCompany %d: %v", id, err)
		}
	}
	// Format output into JSON
	b, err := json.Marshal(company)
	if err != nil {
		fmt.Fprintf(w, "getCompany: %v", err)
	}
	fmt.Fprintf(w, "%v", string(b))
}

// Update details of company with <id> if it exists
// Http request: PUT /companies/id
func updateCompany(w http.ResponseWriter, req *http.Request) {
	// Retrieve the id from URL
	url := req.URL.String()
	stringId := strings.ReplaceAll(url, "/companies/", "")
	id, err := strconv.ParseInt(stringId, 10, 64)

	var company Company
	// Retrieve updated company details from body of request message
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "updateCompany: %v", err)
	}
	err = json.Unmarshal(requestBody, &company)
	if err != nil {
		fmt.Fprintf(w, "updateCompany: %v", err)
	}

	// Update old company with new company details
	_, err = db.Exec("UPDATE companies SET name=?, country=?, website=?, phone=? WHERE id=?", company.Name, company.Country, company.Website, company.Phone, id)
	if err != nil {
		fmt.Fprintf(w, "updateCompany: %v", err)
	} else {
		fmt.Fprintf(w, "Updated company with id %d:\n"+string(requestBody), id)
	}
}

// Remove company with <id>
// Http request: DELETE /companies/id
func deleteCompany(w http.ResponseWriter, req *http.Request) {
	// Retrieve the id from URL
	var url string = req.URL.String()
	var stringId string = strings.ReplaceAll(url, "/companies/", "")
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "deleteCompany: %v", err)
	}

	// Retrieve and store company data that will be deleted to be able to print to user
	var company Company
	row := db.QueryRow("SELECT * FROM companies WHERE id = ?", id)
	if err := row.Scan(&company.Code, &company.Name, &company.Country, &company.Website, &company.Phone); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "deleteCompany: No company with the id %d", id)
			return
		}
		fmt.Fprintf(w, "deleteCompany: %v", err)
	}

	// Delete company with id from database
	_, err = db.Exec("DELETE FROM companies WHERE id = ?", id)
	if err != nil {
		fmt.Fprintf(w, "deleteCompany: %v", err)
	} else {
		b, err := json.Marshal(company)
		if err != nil {
			fmt.Fprintf(w, "deleteCompany: %v", err)
		}
		fmt.Fprintf(w, "Deleted company with id %d:\n%v", id, string(b))
	}
}
