# Project Name

Halo Suster Backend REST API

## Description

This project is a handy web app that helps manage info about users, patients, and nurses. It's built with Go and the Gin web framework, and it uses a PostgreSQL database to keep all the data safe and sound.

## Prerequisites

- Go 1.16 or later
- PostgreSQL 13 or later
- Docker (optional)
- AWS S3 bucket (optional)

## Installation

<!-- TODO: complete the installation instructions -->

1. Clone the repository:
   ```bash
   git clone https://github.com/hnfnfl/halo-suster-be.git
   ```
2. Navigate to the project directory:
   ```bash
   cd halo-suster-be
   ```
3. Install the required dependencies:
   ```bash
   go mod tidy
   ```
4.

## Usage

To start the server, run the following command in the project directory:

```bash
go run .
```

## Endpoints

The following endpoints are available:

- `POST /v1/user/it/register`: Register a new user with the IT role
- `POST /v1/user/it/login`: Log in for a IT staff
- `POST /v1/user/nurse/register`: Register a new nurse
- `POST /v1/user/nurse/login`: Log in for a nurse
- `PUT /v1/user/nurse/:userId`: Update a nurse's info
- `DELETE /v1/user/nurse/:userId`: Delete a nurse from the database
- `GET /v1/user/`: Get all users
- `POST /v1/medical/patient`: Register a new patient
- `GET /v1/medical/patient`: Get all patients
- `POST /v1/medical/record`: Create a new medical record for a patient
- `GET /v1/medical/record`: Get all medical records for a patient
- `POST /v1/image`: Upload an image to an AWS S3 bucket
