# Notes REST API Service
[Read this in Russian](./README.ru.md)

## Description

The project is a REST API service developed in the Go language, which provides an interface for working with notes. The service allows you to add notes, display a list of notes and manage access to them through authentication and authorization (JWT). Validation of spelling errors in notes is carried out through integration with an external service (Yandex.Speller).

## Architecture

- **Programming Language:** Go
- **Web Server Framework:** "net/http"
- **Routing:** chi
- **Database:** PostgreSQL
- **Authentication:** JWT (JSON Web Token)
- **Spelling Validation:** Integration with Yandex.Speller
- **Containerization:** Docker

## Project Structure

- **cmd/** - Entry point of the application.
- **internal/** - Business logic and core application code.
  - **config/** - Application configurations.
  - **database/** - Database connection and migrations.
  - **domain/** - Interface definitions.
  - **handlers/** - HTTP handlers for request processing.
  - **middleware/** - Middleware for request processing.
  - **models/** - Data model definitions.
  - **repository/** - Implementation of database interactions.
  - **services/** - Implementation of services (JWT, validation).
  - **usecases/** - Business logic and data handling.
- **migrations/** - SQL migrations for creating and updating database schema.

## Running the Project

### Prerequisites

- **Go:** 1.18 or higher
- **PostgreSQL:** 13.x or higher
- **Docker:** 20.10 or higher
- **Docker Compose:** 1.27 or higher
- **Postman:** 7.0 or higher (for working with the Postman collection)

### Commands to Run

1. Clone the repository and navigate to the project directory:
    ```bash
    git clone https://github.com/ananikitina/notes-rest.git
    cd notes-rest
    ```

2. Start the containers:
    ```bash
    docker-compose up --build
    ```

3. The API will be available at `http://localhost:8080`.

## Request Formats

- **POST /register** - Create a new user.
- **POST /login** - User authentication and JWT acquisition (the token returned in the response must be saved).
- **POST /note** - Add a note for the user.
- **GET /notes** - Retrieve the user's notes.
- **GET /allnotes** - Retrieve all notes (admin only).

Input data should be in JSON format.

When testing the /note, /notes, and /allnotes routes, include an **Authorization** header with the value **Bearer token**, where token is the JWT obtained during login.

If a spelling error is detected, the note will not be saved to the database, and an error with detailed validation results will be returned.
## Postman Collection

You can import the Postman collection for convenient API testing. The collection file is available in the root of the project: `./Notes REST.postman_collection.json`.

## Preconfigured Users

- **Admin:**
  - **Email:** `admin@ex.com`
  - **Password:** `adminpassword`
  - **Role:** `admin`

- **User:**
  - **Email:** `user@ex.com`
  - **Password:** `userpassword`
  - **Role:** `user`

