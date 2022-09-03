package main

var companyList []Company = make([]Company, 0)

type Company struct {
	Name    string
	Code    string
	Country string
	Website string
	Phone   string
}
