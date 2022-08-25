package main

import (
	"fmt"
)

func main() {
	// Submits entity to specified resource, changing its state
	//resp, err := http.Post("companies.com", Company, ...)
	// Requests a representation of the specified resource
	// Company's attributes should be available to filter
	//resp, err := http.Get("companies.com")
	//http.HandleFunc("/addCompany", addCompany)

	var monk Company = Company{"Monk", "1", "Poland", "www.monk.io", "+447736705394"}
	//fmt.Printf("%+v", monk)
	//fmt.Print(monk)
	companyList = append(companyList, monk)
	fmt.Println(companyList)
	a := fmt.Sprintf("%+v", companyList)
	fmt.Println(a)
}
