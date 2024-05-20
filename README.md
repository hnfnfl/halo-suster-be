# Halo Suster Backend REST API

## Description

This project is a handy web app that helps manage info about users, patients, and nurses. It's built with Go and the Gin web framework, and it uses a PostgreSQL database to keep all the data safe and sound.

## Prerequisites

- Go 1.20 or later
- PostgreSQL 13 or later
- Docker (optional)
- AWS S3 bucket (optional; for image upload)

## Installation

<!-- TODO: complete the installation instructions -->

Here are the steps to install the project:

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
4. Build the binary using Makefile:
   ```bash
   make build
   ```
5. Set up the configuration file in `local_configuration/config.yaml`

6. Migrate the database schema:
   ```bash
   make migrate-up
   ```
7. Run the server:
   ```bash
    make run
   ```
8. The server should be running on `localhost:8080`

Also, you can run the server using Docker:

1. Build the Docker image:

   ```bash
   make docker-build
   ```

2. Update the environment variables in `docker-compose.yml`

3. Run the Docker container:
   ```bash
   make docker-run
   ```

## Configuration

The configuration file is located in `local_configuration/config.yaml`. This file contains the configuration for the database, logging, JWT, and AWS S3.

Here's an example of the configuration file:

```yaml
Environment: development # development, production
LogLevel: debug # debug, info, warn, error
AUTHEXPIRY: 1 # in hours

DB:
  Name: mydb # database name
  Port: 5432 # database port
  Host: localhost # database host
  Username: postgres # database username
  Password: admin # database password
  Params: sslmode=disable # database connection parameters
  MaxIdleConns: 20 # maximum idle connections
  MaxOpenConns: 20 # maximum open connections

Jwt:
  Secret: mysecret # JWT secret key
  BcryptSalt: 10 # Bcrypt salt

AWS:
  Access.Key.ID: myaccesskey # AWS access key ID
  Secret.Access.Key: mysecretaccesskey # AWS secret access key
  S3.Bucket.Name: mybucketname # AWS S3 bucket name
  Region: ap-southeast-1 # AWS region
```

## Endpoints

The following endpoints are available:

- `POST /v1/user/it/register`: Register a new user with the IT role
- `POST /v1/user/it/login`: Log in for an IT staff
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
