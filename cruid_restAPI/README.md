# Vehicle Loan Management API

A comprehensive REST API for managing vehicle loans, customer information, and loan submissions.

## Project Description

This API provides a complete solution for loan management systems, allowing users to:
- Create and manage customer profiles
- Submit and track loan applications
- Manage loan approvals and processing

Built with Go and using SQLite for data storage, this API offers a lightweight yet powerful solution for loan management needs.

## Setup Instructions

### Prerequisites

- Go 1.16 or higher
- SQLite3

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/alphaloan/vehicle.git
   cd vehicle
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build the application:
   ```
   make build
   ```

4. Run the server:
   ```
   make run
   ```

The server will start on the configured port (default: 8080).

### Database Migration

The application automatically handles database migrations on startup. The initial migration creates the necessary tables for customers, loans, and loan submissions.

To manually run migrations:
```
make migrate
```

## API Endpoints

### Customer Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/loan/customers` | Get all customers |
| GET | `/api/loan/customers/:id` | Get customer by ID |
| POST | `/api/loan/customers` | Create a new customer |
| PUT | `/api/loan/customers/:id` | Update an existing customer |
| DELETE | `/api/loan/customers/:id` | Delete a customer |

### Loan Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/loan/submissions` | Get all loan submissions |
| GET | `/api/loan/submissions/:id` | Get loan submission by ID |
| POST | `/api/loan/submissions` | Create a new loan submission |
| PUT | `/api/loan/submissions/:id` | Update a loan submission |
| DELETE | `/api/loan/submissions/:id` | Delete a loan submission |

### Loan Submission Process

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/loan/submit` | Submit a new loan application |

## Data Models

### Customer

```json
{
  "id": "UUID",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone_number": "string",
  "address": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Loan Submission

```json
{
  "id": "UUID",
  "customer_id": "UUID",
  "vehicle_type": "string",
  "vehicle_make": "string",
  "vehicle_model": "string",
  "vehicle_year": "int",
  "loan_amount": "float",
  "loan_term_months": "int",
  "interest_rate": "float",
  "status": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Usage Examples

### Creating a New Customer

**Request:**
```bash
curl -X POST http://localhost:8080/api/loan/customers \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "555-123-4567",
    "address": "123 Main St, Anytown, USA"
  }'
```

**Response:**
```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "555-123-4567",
  "address": "123 Main St, Anytown, USA",
  "created_at": "2023-04-12T15:30:45Z",
  "updated_at": "2023-04-12T15:30:45Z"
}
```

### Submitting a Loan Application

**Request:**
```bash
curl -X POST http://localhost:8080/api/loan/submit \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "vehicle_type": "car",
    "vehicle_make": "Toyota",
    "vehicle_model": "Camry",
    "vehicle_year": 2022,
    "loan_amount": 25000,
    "loan_term_months": 60,
    "interest_rate": 3.5
  }'
```

**Response:**
```json
{
  "id": "7fa85f64-5717-4562-b3fc-2c963f66afa9",
  "customer_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "vehicle_type": "car",
  "vehicle_make": "Toyota",
  "vehicle_model": "Camry",
  "vehicle_year": 2022,
  "loan_amount": 25000,
  "loan_term_months": 60,
  "interest_rate": 3.5,
  "status": "pending",
  "created_at": "2023-04-12T15:35:22Z",
  "updated_at": "2023-04-12T15:35:22Z"
}
```

## Development

### Running Tests

```
make test
```

### Code Formatting

```
make fmt
```

### Building for Production

```
make build-prod
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

# Welcome to Udemy Go Labs!

Go labs are based on Go ??? in the Ubuntu distribution. You can practice Go coding as you follow the lab tasks.

Following commands are also supported: 
* vim 
* wget 
* zsh

