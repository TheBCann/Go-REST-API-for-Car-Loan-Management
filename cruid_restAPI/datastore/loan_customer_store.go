package datastore

import (
	"database/sql"
	"fmt"
)

type LoanCustomerRow struct {
	CustomerID    string
	IDCardNumber  string
	FullName      string
	BirthDate     string
	PhoneNumber   string
	Email         sql.NullString
	MonthlyIncome float64
	AddressStreet string
	AddressCity   string
}

type LoanCustomerStore struct {
	db *sql.DB
}

func NewLoanCustomerStore(db *sql.DB) *LoanCustomerStore {
	return &LoanCustomerStore{
		db: db,
	}
}

const sqlUpsertCustomer = `
	INSERT INTO loan_customers (
		customer_id,	
		id_card_number,	
		full_name,	
		birth_date,		
		phone_number,	
		email,		
		monthly_income,	
		address_street,	
		address_city
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	) ON CONFLICT (id_card_number) DO UPDATE SET
		full_name = EXCLUDED.full_name,    
        birth_date = EXCLUDED.birth_date,    
        phone_number = EXCLUDED.phone_number,   
	    email = EXCLUDED.email, 
        monthly_income = EXCLUDED.monthly_income,	      
        address_street = EXCLUDED.address_street,
        address_city = EXCLUDED.address_city
	RETURNING customer_id;
`

func (s *LoanCustomerStore) UpsertCustomer(customer *LoanCustomerRow) (string, error) {
	var customerID string
	err := s.db.QueryRow(sqlUpsertCustomer,
		customer.CustomerID,
		customer.IDCardNumber,
		customer.FullName,
		customer.BirthDate,
		customer.PhoneNumber,
		customer.Email,
		customer.MonthlyIncome,
		customer.AddressStreet,
		customer.AddressCity,
	).Scan(&customerID)

	if err != nil {
		return "", err
	}

	return customerID, nil
}

const sqlGetAllCustomers = `
SELECT
	customer_id, id_card_number,
	full_name, birth_date,
	phone_number, email,
	monthly_income, address_street,
	address_city
FROM loan_customers
ORDER BY full_name;
`

func (s *LoanCustomerStore) GetAllCustomers() ([]*LoanCustomerRow, error) {
	rows, err := s.db.Query(sqlGetAllCustomers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*LoanCustomerRow
	for rows.Next() {
		var customer LoanCustomerRow
		err := rows.Scan(
			&customer.CustomerID,
			&customer.IDCardNumber,
			&customer.FullName,
			&customer.BirthDate,
			&customer.PhoneNumber,
			&customer.Email,
			&customer.MonthlyIncome,
			&customer.AddressStreet,
			&customer.AddressCity,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

const sqlGetCustomerByID = `
SELECT
	c.customer_id, c.id_card_number,
	c.full_name, c.birth_date,
	c.phone_number, c.email,
	c.monthly_income, c.address_street,
	c.address_city, s.submission_id,
	s.vehicle_brand, s.vehicle_type,
	s.vehicle_model, s.vehicle_license_number,
	s.vehicle_odometer, s.manufacturing_year,
	s.proposed_loan_amount, s.proposed_loan_tenure_month,
	s.is_commercial_vehicle, s.created_at,
	s.updated_at
FROM loan_customers c
INNER JOIN loan_submissions s
ON c.customer_id = s.customer_id
WHERE c.customer_id = $1
ORDER BY s.created_at DESC;
`

type LoanCustomerWithAllSubmissionsRow struct {
	LoanCustomerRow *LoanCustomerRow
	LoanSubmissions []*LoanSubmissionRow
}

func (s *LoanCustomerStore) GetCustomerByID(customerID string) (*LoanCustomerWithAllSubmissionsRow, error) {
	rows, err := s.db.Query(sqlGetCustomerByID, customerID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customer *LoanCustomerRow
	var submissions []*LoanSubmissionRow

	for rows.Next() {
		var submission LoanSubmissionRow
		if customer == nil {
			customer = &LoanCustomerRow{}
			err := rows.Scan(
				&customer.CustomerID,
				&customer.IDCardNumber,
				&customer.FullName,
				&customer.BirthDate,
				&customer.PhoneNumber,
				&customer.Email,
				&customer.MonthlyIncome,
				&customer.AddressStreet,
				&customer.AddressCity,
				&submission.SubmissionID,
				&submission.VehicleBrand,
				&submission.VehicleType,
				&submission.VehicleModel,
				&submission.VehicleLicenseNumber,
				&submission.VehicleOdometer,
				&submission.ManufacturingYear,
				&submission.ProposedLoanAmount,
				&submission.ProposedLoanTenure,
				&submission.IsCommercialVehicle,
				&submission.CreatedAt,
				&submission.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := rows.Scan(
				new(string), new(string), new(string), new(string), new(string), new(sql.NullString), new(float64), new(string), new(string),
				&submission.SubmissionID,
				&submission.VehicleBrand,
				&submission.VehicleType,
				&submission.VehicleModel,
				&submission.VehicleLicenseNumber,
				&submission.VehicleOdometer,
				&submission.ManufacturingYear,
				&submission.ProposedLoanAmount,
				&submission.ProposedLoanTenure,
				&submission.IsCommercialVehicle,
				&submission.CreatedAt,
				&submission.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
		submissions = append(submissions, &submission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	return &LoanCustomerWithAllSubmissionsRow{
		LoanCustomerRow: customer,
		LoanSubmissions: submissions,
	}, nil
}

const sqlUpdateCustomerByID = `
UPDATE loan_customers
SET
	full_name = COALESCE($1, full_name),
	birth_date = COALESCE($2, birth_date),
	phone_number = COALESCE($3, phone_number),
	email = COALESCE($4, email),
	monthly_income = COALESCE($5, monthly_income),
	address_street = COALESCE($6, address_street),
	address_city = COALESCE($7, address_city),
	id_card_number = COALESCE($8, id_card_number)
WHERE customer_id = $9
RETURNING customer_id;
`

func (s *LoanCustomerStore) UpdateCustomerByID(customer *LoanCustomerRow, customerIDToUpdate string) (string, error) {
	var customerID string
	err := s.db.QueryRow(sqlUpdateCustomerByID,
		customer.FullName,
		customer.BirthDate,
		customer.PhoneNumber,
		customer.Email,
		customer.MonthlyIncome,
		customer.AddressStreet,
		customer.AddressCity,
		customer.IDCardNumber,
		customerIDToUpdate,
	).Scan(&customerID)

	if err != nil {
		return "", err
	}

	return customerID, nil
}

const sqlDeleteCustomerByCustomerID = `
DELETE FROM	loan_customers
WHERE customer_id = $1
RETURNING customer_id;
`

func (s *LoanCustomerStore) DeleteCustomerByID(customerIDToDelete string) (string, error) {
	var customerID string
	err := s.db.QueryRow(sqlDeleteCustomerByCustomerID, customerIDToDelete).Scan(&customerID)

	if err != nil {
		return "", err
	}

	return customerID, nil
}
