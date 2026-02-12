# Bank Management System


A comprehensive banking system backend built with Go, Gin, GORM, and PostgreSQL. This system supports multiple banks and branches, customer account management, loan management, and transaction tracking.

> ![ER Diagram](internal/banksystem%20-%20er%20diagram.png)


## Features

- **Multi-Bank Support**: Manage multiple banks and their branches
- **Customer Management**: Create and manage customer profiles
- **Account Management**:
  - Create savings accounts
  - Deposit and withdraw funds
  - View account details and balance
- **Loan Management**:
  - Take loans at 12% interest rate
  - Repay loans
  - View loan details including pending amount and interest calculations
  - Interest is automatically accrued and added to the remaining loan amount based on elapsed time (years and months)
- **Transaction History**: Track all transactions (deposits, withdrawals, loan transactions)
- **Joint Accounts**: Support for joint account ownership and roles
- **Branch and Bank Management**: CRUD for banks and branches
- **Comprehensive API**: RESTful endpoints for all major operations

## API Endpoints (Summary)

- **Health Check**
  - `GET /check` — System health
- **Banks**
  - `POST /api/banks` — Create bank
  - `GET /api/banks` — List banks
  - `GET /api/banks/:id` — Bank details
- **Branches**
  - `POST /api/branches` — Create branch
  - `GET /api/branches` — List branches
  - `GET /api/branches/:id` — Branch details
  - `GET /api/branches/bank/:bank_id` — Branches by bank
- **Customers**
  - `POST /api/customers` — Create customer
  - `GET /api/customers` — List customers
  - `GET /api/customers/:id` — Customer details
  - `PUT /api/customers/:id` — Update customer
- **Accounts**
  - `POST /api/accounts` — Create account
  - `GET /api/accounts/:id` — Account details
  - `GET /api/accounts/customer/:customer_id` — Accounts by customer
  - `POST /api/accounts/:id/deposit` — Deposit
  - `POST /api/accounts/:id/withdraw` — Withdraw
- **Joint Accounts**
  - `POST /api/joint-accounts` — Create joint account
  - `GET /api/joint-accounts/:account_id` — Joint account details
  - `GET /api/joint-accounts/:account_id/owners` — Account owners
  - `GET /api/joint-accounts/customer/:customer_id` — Joint accounts by customer
  - `POST /api/joint-accounts/:account_id/customers` — Add customer to joint account
  - `DELETE /api/joint-accounts/:account_id/customers/:customer_id` — Remove customer
  - `PUT /api/joint-accounts/:account_id/customers/:customer_id/role` — Update customer role
- **Loans**
  - `POST /api/loans` — Create loan
  - `GET /api/loans/:id` — Loan details
  - `GET /api/loans/:id/details` — Comprehensive loan details
  - `GET /api/loans/customer/:customer_id` — Loans by customer
  - `POST /api/loans/:id/repay` — Repay loan
- **Transactions**
  - `GET /api/transactions/account/:id` — Account transactions
  - `GET /api/transactions/customer/:customer_id` — Customer transactions
  - `GET /api/transactions/loan/:loan_id` — Loan transactions

## Tech Stack

- **Go 1.21+**: Programming language
- **Gin**: Web framework
- **GORM**: ORM for database operations
- **PostgreSQL**: Database


