package datastore

import (
	"database/sql"
)

type LoanSubmissionRow struct {
	SubmissionID         string
	VehicleType          string
	VehicleBrand         string
	VehicleModel         string
	VehicleLicenseNumber string
	VehicleOdometer      int
	ManufacturingYear    int
	ProposedLoanAmount   int
	ProposedLoanTenure   int
	LoanStatus           string
	IsCommercialVehicle  bool
	CreatedAt            int64
	UpdatedAt            int64
	CustomerID           string
}

type LoanSubmissionStore struct {
	db *sql.DB
}

func NewLoanSubmissionStore(db *sql.DB) *LoanSubmissionStore {
	return &LoanSubmissionStore{
		db: db,
	}
}

const sqlUpsertSubmission = `
	INSERT INTO loan_submissions (
		submission_id,		
		vehicle_type,		
		vehicle_brand,		
		vehicle_model,		
		vehicle_license_number,
		vehicle_odometer,
		manufacturing_year,	
		proposed_loan_amount,
		proposed_loan_tenure_month,	
		loan_status,			
		is_commercial_vehicle,	
		created_at,			
		updated_at,			
		customer_id			
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	) ON CONFLICT (submission_id) DO UPDATE SET
		vehicle_type = EXCLUDED.vehicle_type, 
		vehicle_brand = EXCLUDED.vehicle_brand,
		vehicle_model = EXCLUDED.vehicle_model,
		vehicle_license_number = EXCLUDED.vehicle_license_number,
		vehicle_odometer = EXCLUDED.vehicle_odometer,
		manufacturing_year = EXCLUDED.manufacturing_year,
		proposed_loan_amount = EXCLUDED.proposed_loan_amount,
		proposed_loan_tenure_month = EXCLUDED.proposed_loan_tenure_month,
		loan_status = EXCLUDED.loan_status,
		is_commercial_vehicle = EXCLUDED.is_commercial_vehicle,
		created_at = EXCLUDED.created_at,	
	    updated_at = EXCLUDED.updated_at,	
        customer_id = EXCLUDED.customer_id
	RETURNING submission_id;
`

func (s *LoanSubmissionStore) UpsertSubmission(submission *LoanSubmissionRow) (string, error) {
	var submissionID string

	err := s.db.QueryRow(sqlUpsertSubmission,
		submission.SubmissionID,
		submission.VehicleType,
		submission.VehicleBrand,
		submission.VehicleModel,
		submission.VehicleLicenseNumber,
		submission.VehicleOdometer,
		submission.ManufacturingYear,
		submission.ProposedLoanAmount,
		submission.ProposedLoanTenure,
		submission.LoanStatus,
		submission.IsCommercialVehicle,
		submission.CreatedAt,
		submission.UpdatedAt,
		submission.CustomerID,
	).Scan(&submissionID)

	if err != nil {
		return "", err
	}

	return submissionID, nil
}

const sqlGetAllLoanSubmissions = `
SELECT
	submission_id, vehicle_type,
	vehicle_brand, vehicle_model,
	vehicle_license_number, vehicle_odometer,
	manufacturing_year, proposed_loan_amount,
	proposed_loan_tenure_month, loan_status,
	is_commercial_vehicle, created_at,
	updated_at, customer_id
FROM loan_submissions
ORDER BY created_at DESC;
`

func (s *LoanSubmissionStore) GetAllLoanSubmissions() ([]*LoanSubmissionRow, error) {
	rows, err := s.db.Query(sqlGetAllLoanSubmissions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*LoanSubmissionRow
	for rows.Next() {
		submission := &LoanSubmissionRow{}
		err := rows.Scan(
			&submission.SubmissionID,
			&submission.VehicleType,
			&submission.VehicleBrand,
			&submission.VehicleModel,
			&submission.VehicleLicenseNumber,
			&submission.VehicleOdometer,
			&submission.ManufacturingYear,
			&submission.ProposedLoanAmount,
			&submission.ProposedLoanTenure,
			&submission.LoanStatus,
			&submission.IsCommercialVehicle,
			&submission.CreatedAt,
			&submission.UpdatedAt,
			&submission.CustomerID,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return submissions, nil
}

const sqlGetLoanSubmissionByID = `
SELECT
	submission_id, vehicle_type,
	vehicle_brand, vehicle_model,
	vehicle_license_number, vehicle_odometer,
	manufacturing_year, proposed_loan_amount,
	proposed_loan_tenure_month, loan_status,
	is_commercial_vehicle, created_at,
	updated_at, customer_id
FROM loan_submissions
WHERE submission_id = $1;
`

func (s *LoanSubmissionStore) GetLoanSubmissionByID(submissionID string) (*LoanSubmissionRow, error) {
	row := s.db.QueryRow(sqlGetLoanSubmissionByID, submissionID)
	submission := &LoanSubmissionRow{}
	err := row.Scan(
		&submission.SubmissionID,
		&submission.VehicleType,
		&submission.VehicleBrand,
		&submission.VehicleModel,
		&submission.VehicleLicenseNumber,
		&submission.VehicleOdometer,
		&submission.ManufacturingYear,
		&submission.ProposedLoanAmount,
		&submission.ProposedLoanTenure,
		&submission.LoanStatus,
		&submission.IsCommercialVehicle,
		&submission.CreatedAt,
		&submission.UpdatedAt,
		&submission.CustomerID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return submission, nil
}
