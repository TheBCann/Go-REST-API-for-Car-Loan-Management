package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
)

type LoanCustomerHandler struct {
	CustomerStore   datastore.LoanCustomerStore
	SubmissionStore datastore.LoanSubmissionStore
}

func NewLoanCustomerHandler(
	customerStore datastore.LoanCustomerStore,
	submissionStore datastore.LoanSubmissionStore) *LoanCustomerHandler {
	return &LoanCustomerHandler{
		CustomerStore:   customerStore,
		SubmissionStore: submissionStore,
	}
}

func (h *LoanCustomerHandler) HandleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	loanCustomerRows, err := h.CustomerStore.GetAllCustomers()

	if err != nil {
		errMsg := "Failed to get all loan customers"
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	loanCustomers := make([]LoanCustomer, 0, len(loanCustomerRows))
	for _, row := range loanCustomerRows {
		loanCustomers = append(loanCustomers, LoanCustomer{
			CustomerID:    row.CustomerID,
			IDCardNumber:  row.IDCardNumber,
			FullName:      row.FullName,
			BirthDate:     row.BirthDate,
			PhoneNumber:   row.PhoneNumber,
			Email:         &row.Email.String,
			MonthlyIncome: row.MonthlyIncome,
			AddressStreet: row.AddressStreet,
			AddressCity:   row.AddressCity,
		})
	}

	responseBody := GetAllLoanCustomersResponse{
		Data: &loanCustomers,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func validateCustomerID(w http.ResponseWriter, customerID string) bool {
	if customerID == "" {
		errMsg := "Missing customer_id path variable"
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return false
	}

	if !IsValidUUID(customerID) {
		errMsg := "Invalid customer_id: " + customerID
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return false
	}

	return true
}

func (h *LoanCustomerHandler) HandleGetCustomerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	customerID := r.PathValue("customer_id")
	if !validateCustomerID(w, customerID) {
		return
	}

	loanCustomerWithAllSubmissionsRow, err := h.CustomerStore.GetCustomerByID(customerID)

	if loanCustomerWithAllSubmissionsRow == nil {
		errMsg := "Customer not found"
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	if err != nil {
		errMsg := "Failed to get loan customer"
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	loanSubmissions := make([]LoanSubmission, 0, len(loanCustomerWithAllSubmissionsRow.LoanSubmissions))
	for _, submissionRow := range loanCustomerWithAllSubmissionsRow.LoanSubmissions {
		loanSubmissions = append(loanSubmissions, LoanSubmission{
			SubmissionID:            submissionRow.SubmissionID,
			VehicleType:             submissionRow.VehicleType,
			VehicleBrand:            submissionRow.VehicleBrand,
			VehicleModel:            submissionRow.VehicleModel,
			VehicleLicenseNumber:    submissionRow.VehicleLicenseNumber,
			VehicleOdometer:         submissionRow.VehicleOdometer,
			ManufacturingYear:       submissionRow.ManufacturingYear,
			ProposedLoanAmount:      submissionRow.ProposedLoanAmount,
			ProposedLoanTenureMonth: submissionRow.ProposedLoanTenure,
			IsCommercialVehicle:     submissionRow.IsCommercialVehicle,
		})
	}

	loanCustomer := LoanCustomerWithAllSubmissions{
		Customer: &LoanCustomer{
			CustomerID:    loanCustomerWithAllSubmissionsRow.LoanCustomerRow.CustomerID,
			IDCardNumber:  loanCustomerWithAllSubmissionsRow.LoanCustomerRow.IDCardNumber,
			FullName:      loanCustomerWithAllSubmissionsRow.LoanCustomerRow.FullName,
			BirthDate:     loanCustomerWithAllSubmissionsRow.LoanCustomerRow.BirthDate,
			PhoneNumber:   loanCustomerWithAllSubmissionsRow.LoanCustomerRow.PhoneNumber,
			Email:         nil,
			MonthlyIncome: loanCustomerWithAllSubmissionsRow.LoanCustomerRow.MonthlyIncome,
			AddressStreet: loanCustomerWithAllSubmissionsRow.LoanCustomerRow.AddressStreet,
			AddressCity:   loanCustomerWithAllSubmissionsRow.LoanCustomerRow.AddressCity,
		},
		Submissions: &loanSubmissions,
	}

	if loanCustomerWithAllSubmissionsRow.LoanCustomerRow.Email.Valid {
		loanCustomer.Customer.Email = &loanCustomerWithAllSubmissionsRow.LoanCustomerRow.Email.String
	}

	responseBody := GetCustomerInfoResponse{
		Data: &loanCustomer,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func (h *LoanCustomerHandler) HandleUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Only PATCH method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	customerID := r.PathValue("customer_id")
	if !validateCustomerID(w, customerID) {
		return
	}

	var request LoanCustomer
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	LoanCustomerRow := convertLoanCustomer(&request)

	updatedCustomerID, err := h.CustomerStore.UpdateCustomerByID(LoanCustomerRow, customerID)

	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "Customer ID not found"
			responseBodyErr := UpdateCustomerResponse{
				ErrorMessage: &errMsg,
				CustomerID:   &customerID,
				Updated:      false,
			}

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(responseBodyErr)
		} else {
			errMsg := "Failed to update customer"
			responseBodyErr := UpdateCustomerResponse{
				ErrorMessage: &errMsg,
				CustomerID:   &customerID,
				Updated:      false,
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseBodyErr)
		}
		return
	}

	responseBody := UpdateCustomerResponse{
		CustomerID: &updatedCustomerID,
		Updated:    true,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func (h *LoanCustomerHandler) HandleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	customerID := r.PathValue("customer_id")
	if !validateCustomerID(w, customerID) {
		return
	}

	deleteCustomerID, err := h.CustomerStore.DeleteCustomerByID(customerID)

	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "Customer ID not found"
			responseBodyErr := DeleteCustomerResponse{
				ErrorMessage: &errMsg,
				CustomerID:   &customerID,
				Deleted:      false,
			}

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(responseBodyErr)
		} else {
			errMsg := "Failed to delete customer"
			responseBodyErr := DeleteCustomerResponse{
				ErrorMessage: &errMsg,
				CustomerID:   &customerID,
				Deleted:      false,
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseBodyErr)
		}
		return
	}

	responseBody := DeleteCustomerResponse{
		CustomerID: &deleteCustomerID,
		Deleted:    true,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}
