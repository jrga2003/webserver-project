package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Capture connection properties
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "companies",
		AllowNativePasswords: true,
	}
	// Get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	// Register http handlers
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
