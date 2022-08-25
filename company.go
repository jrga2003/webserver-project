package main

var companyList []Company = make([]Company, 0)

type Company struct {
	Name    string
	Code    string
	Country string
	Website string
	Phone   string
}

func (c Company) String() string {
	return "{" + "Name: " + c.Name + " Code: " + c.Code + " Country: " + c.Country + " Website: " + c.Website + " Phone: " + c.Phone + "}"
}
