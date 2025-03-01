package handler

import (
	"database/sql"
	"time"

	"github.com/alphaloan/vehicle/datastore"
	"github.com/google/uuid"
)

type LoanCustomer struct {
	CustomerID    string  `json:"customer_id"`
	IDCardNumber  string  `json:"id_card_number"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`
	PhoneNumber   string  `json:"phone_number"`
	Email         *string `json:"email"`
	MonthlyIncome float64 `json:"monthly_income"`
	AddressStreet string  `json:"address_street"`
	AddressCity   string  `json:"address_city"`
}

type LoanSubmission struct {
	SubmissionID            string `json:"submission_id"`
	VehicleType             string `json:"vehicle_type"`
	VehicleBrand            string `json:"vehicle_brand"`
	VehicleModel            string `json:"vehicle_model"`
	VehicleLicenseNumber    string `json:"vehicle_license_number"`
	VehicleOdometer         int    `json:"vehicle_odometer"`
	ManufacturingYear       int    `json:"manufacturing_year"`
	ProposedLoanAmount      int    `json:"proposed_loan_amount"`
	ProposedLoanTenureMonth int    `json:"proposed_loan_tenure_month"`
	IsCommercialVehicle     bool   `json:"is_commercial_vehicle"`
}

type LoanSubmitRequest struct {
	Customer     LoanCustomer   `json:"customer"`
	ProposedLoan LoanSubmission `json:"proposed_loan"`
}

type LoanSubmitResponse struct {
	CustomerID   *string `json:"customer_id"`
	SubmissionID *string `json:"submission_id"`
}

func convertLoanCustomer(loanCustomer *LoanCustomer) *datastore.LoanCustomerRow {
	if loanCustomer == nil {
		return nil
	}

	parsedEmail := ""

	if loanCustomer.Email != nil {
		parsedEmail = *loanCustomer.Email
	}

	return &datastore.LoanCustomerRow{
		CustomerID:   uuid.New().String(),
		IDCardNumber: loanCustomer.IDCardNumber,
		FullName:     loanCustomer.FullName,
		BirthDate:    loanCustomer.BirthDate,
		PhoneNumber:  loanCustomer.PhoneNumber,
		Email: sql.NullString{
			String: parsedEmail,
			Valid:  loanCustomer.Email != nil && *loanCustomer.Email != "",
		},
		MonthlyIncome: loanCustomer.MonthlyIncome,
		AddressStreet: loanCustomer.AddressStreet,
		AddressCity:   loanCustomer.AddressCity,
	}
}

func convertLoanProposal(loanProposal *LoanSubmission, customerID string) *datastore.LoanSubmissionRow {
	if loanProposal == nil {
		return nil
	}

	now := time.Now().Unix()

	return &datastore.LoanSubmissionRow{
		SubmissionID:         uuid.New().String(),
		VehicleType:          loanProposal.VehicleType,
		VehicleBrand:         loanProposal.VehicleBrand,
		VehicleModel:         loanProposal.VehicleModel,
		VehicleLicenseNumber: loanProposal.VehicleLicenseNumber,
		VehicleOdometer:      loanProposal.VehicleOdometer,
		ManufacturingYear:    loanProposal.ManufacturingYear,
		ProposedLoanAmount:   loanProposal.ProposedLoanAmount,
		ProposedLoanTenure:   loanProposal.ProposedLoanTenureMonth,
		LoanStatus:           "NEW",
		IsCommercialVehicle:  loanProposal.IsCommercialVehicle,
		CreatedAt:            now,
		UpdatedAt:            now,
		CustomerID:           customerID,
	}
}

type GetAllLoanSubmissionsResponse struct {
	ErrorMessage *string           `json:"error_message"`
	Data         *[]LoanSubmission `json:"data"`
}

type LoanSubmissionTrackStatusResponse struct {
	ErrorMessage *string         `json:"error_message"`
	Data         *LoanSubmission `json:"data"`
}

type GetAllLoanCustomersResponse struct {
	ErrorMessage *string         `json:"error_message"`
	Data         *[]LoanCustomer `json:"data"`
}

type LoanCustomerWithAllSubmissions struct {
	Customer    *LoanCustomer     `json:"customer"`
	Submissions *[]LoanSubmission `json:"loan_submissions"`
}

type GetCustomerInfoResponse struct {
	ErrorMessage *string                         `json:"error_message"`
	Data         *LoanCustomerWithAllSubmissions `json:"data"`
}

type UpdateCustomerResponse struct {
	ErrorMessage *string `json:"error_message"`
	CustomerID   *string `json:"customer_id"`
	Updated      bool    `json:"updated"`
}

type DeleteCustomerResponse struct {
	ErrorMessage *string `json:"error_message"`
	CustomerID   *string `json:"customer_id"`
	Deleted      bool    `json:"deleted"`
}
