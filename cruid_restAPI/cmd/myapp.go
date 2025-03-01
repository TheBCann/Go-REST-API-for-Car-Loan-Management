package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
	"github.com/alphaloan/vehicle/handler"
)

func main() {
	datastore.InitializeDatabase("db/migration", "sqlite3://alphaloan.db")

	db, err := sql.Open("sqlite3", "alphaloan.db")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")

	if err != nil {
		log.Fatal("Failed to enable foreign key support:", err)
	}

	defer db.Close()

	loanCustomerStore := datastore.NewLoanCustomerStore(db)
	loanSubmissionStore := datastore.NewLoanSubmissionStore(db)

	loanSubmitHandler := handler.NewLoanSubmitHandler(*loanCustomerStore, *loanSubmissionStore)

	http.HandleFunc("/api/loan/submit", loanSubmitHandler.HandleSubmitLoan)

	loanSubmissionHandler := handler.NewLoanSubmissionHandler(*loanSubmissionStore)

	http.HandleFunc("/api/loan/submissions", loanSubmissionHandler.HandleGetAllLoanSubmissions)

	http.HandleFunc("/api/loan/submission/track", loanSubmissionHandler.HandleTrackLoanSubmission)

	loanCustomerHandler := handler.NewLoanCustomerHandler(*loanCustomerStore, *loanSubmissionStore)
	http.HandleFunc("/api/loan/customers", loanCustomerHandler.HandleGetAllCustomers)

	http.HandleFunc("/api/loan/customer/{customer_id}/info", loanCustomerHandler.HandleGetCustomerInfo)

	http.HandleFunc("/api/loan/customer/{customer_id}/update", loanCustomerHandler.HandleUpdateCustomer)

	http.HandleFunc("/api/loan/customer/{customer_id}/delete", loanCustomerHandler.HandleDeleteCustomer)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
