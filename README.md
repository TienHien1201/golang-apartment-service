# Apartment business Service

## Introduction
This project is a Go-based service for managing booking apartment vinhome and social. The Makefile included provides a set of commands to streamline development, testing, and deployment.

## Prerequisites
Ensure you have the following installed:
- Go (latest version)
- Docker (for SonarQube and Swagger UI)
- Git
- `golangci-lint`
- `swag`

## Setup
To initialize the development environment, run:
```sh
make init
```
This installs necessary tools such as `golangci-lint` and `swag` if they are not already available.

## Running the Application
To start the application with default settings:
```sh
make run
```
To run the application in different environments:
```sh
make dev   # Runs the application in the DEV environment
make qc    # Runs the application in the QC environment
```

## Building the Application
To build the application for different environments:
```sh
make build-dev   # Build for DEV environment
make build-qc    # Build for QC environment
make build-prod  # Build for PROD environment
```
The build output is placed in `build_<ENV>/` with the necessary configuration files.

## API Documentation
To generate Swagger documentation:
```sh
make gen-swagger
```
To run a local Swagger UI:
```sh
make run-swagger
```
To push the generated Swagger JSON to the remote API documentation server:
```sh
make push-swagger
```

## Code Quality and Linting
To check the code formatting and linting:
```sh
make run
```
This runs `golangci-lint` on the codebase.

## Static Code Analysis with SonarQube
To run a SonarQube scan:
```sh
make sonar
```
Ensure SonarQube is properly configured and accessible before running this command.

## Database Migrations
To run database migrations:
```sh
make migrate ENV=dev DB=mysql COMMAND=up STEPS=1
```
Other migration commands:
- `COMMAND=down` to roll back a migration
- `COMMAND=force STEPS=10` to force a migration to step 10
- `COMMAND=version` to check the current migration version

## Help
To see all available commands:
```sh
make help
```
This will display a list of supported commands along with their descriptions.

## License
This project is licensed under the MIT License.

