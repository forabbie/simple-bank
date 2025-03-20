# Simple Bank

From [Backend Master Class](https://bit.ly/backendmaster) course by [TECH SCHOOL](https://bit.ly/m/techschool). A course that teaches step-by-step how to design, develop, and deploy a backend web service from scratch.

## Overview

Simple Bank provides APIs for the frontend to perform the following operations:

- Create and manage bank accounts.
- Record all balance changes to each account.
- Perform money transfers between accounts.

## Features

1. **Bank Accounts**: Users can create and manage accounts with an owner's name, balance, and currency.
2. **Transaction Records**: Every deposit or withdrawal is recorded for accountability.
3. **Money Transfers**: Secure transactions between accounts are handled within a transaction to ensure consistency.

## Setup Local Development

### Install Required Tools

Ensure you have the following tools installed:

- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

  ```bash
  brew install golang-migrate
  ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

  ```bash
  brew install sqlc
  ```

- [Gomock](https://github.com/golang/mock)

  ```bash
  go install github.com/golang/mock/mockgen@v1.6.0
  ```

### Setup Infrastructure

1. Create the bank network:

   ```bash
   make network
   ```

2. Start the PostgreSQL container:

   ```bash
   make postgres
   ```

3. Create the `simple_bank` database:

   ```bash
   make createdb
   ```

4. Run database migrations:

   ```bash
   make migrateup
   ```

5. Run a single migration:

   ```bash
   make migrateup1
   ```

6. Rollback all migrations:

   ```bash
   make migratedown
   ```

7. Rollback a single migration:

   ```bash
   make migratedown1
   ```

### Code Generation

- Generate schema SQL file:

  ```bash
  make db_schema
  ```

- Generate SQL CRUD operations:

  ```bash
  make sqlc
  ```

- Generate database mock:

  ```bash
  make mock
  ```

- Create a new database migration:

  ```bash
  make new_migration name=<migration_name>
  ```

### Running the Project

- Start the server:

  ```bash
  make server
  ```

- Run tests:

  ```bash
  make test
  ```

## API Documentation 📖

### User Management

#### Register a New User

```
POST /api/v1/users
```

**Request Body:**

```json
{
  "username": "exampleUser",
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Login User

```
POST /api/v1/users/login
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Renew Access Token

```
POST /api/v1/tokens/renew_access
```

**Request Body:**

```json
{
  "refresh_token": "your-refresh-token"
}
```

#### Update User

```
PATCH /api/v1/users/:id
```

**Request Body:**

```json
{
  "email": "newemail@example.com"
}
```

### Account Management

#### Create Account

```
POST /api/v1/accounts
```

**Request Body:**

```json
{
  "owner": "exampleUser",
  "balance": 1000
}
```

#### Get Account by ID

```
GET /api/v1/accounts/:id
```

#### List Accounts

```
GET /api/v1/accounts?page={page}&limit={limit}
```

### Transfer Management

#### Create Transfer

```
POST /api/v1/transfers
```

**Request Body:**

```json
{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 500
}
```

#### List Transfers

```
GET /api/v1/transfers?page={page}&limit={limit}
```

---
