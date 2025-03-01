package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
)

type LoanSubmissionHandler struct {
	SubmissionStore datastore.LoanSubmissionStore
}

func NewLoanSubmissionHandler(
	submissionStore datastore.LoanSubmissionStore) *LoanSubmissionHandler {
	return &LoanSubmissionHandler{
		SubmissionStore: submissionStore,
	}
}

func (h *LoanSubmissionHandler) HandleGetAllLoanSubmissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	loanSubmissionRows, err := h.SubmissionStore.GetAllLoanSubmissions()

	if err != nil {
		errMsg := "Failed to get all loan submissions"
		responseBodyErr := GetAllLoanSubmissionsResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseBodyErr)

		return
	}

	loanSubmissions := make([]LoanSubmission, 0, len(loanSubmissionRows))
	for _, row := range loanSubmissionRows {
		loanSubmissions = append(loanSubmissions, LoanSubmission{
			SubmissionID:            row.SubmissionID,
			VehicleType:             row.VehicleType,
			VehicleBrand:            row.VehicleBrand,
			VehicleModel:            row.VehicleModel,
			VehicleLicenseNumber:    row.VehicleLicenseNumber,
			VehicleOdometer:         row.VehicleOdometer,
			ManufacturingYear:       row.ManufacturingYear,
			ProposedLoanAmount:      row.ProposedLoanAmount,
			ProposedLoanTenureMonth: row.ProposedLoanTenure,
			IsCommercialVehicle:     row.IsCommercialVehicle,
		})
	}

	responseBody := GetAllLoanSubmissionsResponse{
		Data: &loanSubmissions,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func validateLoanSubmissionID(w http.ResponseWriter, loanSubmissionID string) bool {
	if loanSubmissionID == "" {
		errMsg := "Missing submission_id query parameter"
		responseBodyErr := LoanSubmissionTrackStatusResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return false
	}

	if !IsValidUUID(loanSubmissionID) {
		errMsg := "Invalid submission_id: " + loanSubmissionID
		responseBodyErr := LoanSubmissionTrackStatusResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return false
	}

	return true
}

func (h *LoanSubmissionHandler) HandleTrackLoanSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	loanSubmissionID := r.URL.Query().Get("loan_submission_id")
	if !validateLoanSubmissionID(w, loanSubmissionID) {
		return
	}

	loanSubmissionRow, err := h.SubmissionStore.GetLoanSubmissionByID(loanSubmissionID)

	if err != nil {
		errMsg := "Failed to get loan submission"
		responseBodyErr := LoanSubmissionTrackStatusResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	if loanSubmissionRow == nil {
		errMsg := "Loan submission not found"
		responseBodyErr := LoanSubmissionTrackStatusResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	loanSubmission := LoanSubmission{
		SubmissionID:            loanSubmissionRow.SubmissionID,
		VehicleType:             loanSubmissionRow.VehicleType,
		VehicleBrand:            loanSubmissionRow.VehicleBrand,
		VehicleModel:            loanSubmissionRow.VehicleModel,
		VehicleLicenseNumber:    loanSubmissionRow.VehicleLicenseNumber,
		VehicleOdometer:         loanSubmissionRow.VehicleOdometer,
		ManufacturingYear:       loanSubmissionRow.ManufacturingYear,
		ProposedLoanAmount:      loanSubmissionRow.ProposedLoanAmount,
		ProposedLoanTenureMonth: loanSubmissionRow.ProposedLoanTenure,
		IsCommercialVehicle:     loanSubmissionRow.IsCommercialVehicle,
	}

	responseBody := LoanSubmissionTrackStatusResponse{
		Data: &loanSubmission,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}
