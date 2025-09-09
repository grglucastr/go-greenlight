# Greenlight API
An Open Source API for a Movies app. This is part of the course of Golang called “Let’s go further“ by Alex Edwards.

## 🚀 Features
- **User Registration:** Users can register in order to use the API.
- **Authentication/Authorization:** Users can login to access the reading resources. Other writing endpoints requires special permissions.
- **Movies Management:** Users can manage movies by adding, editing, listing, filtering and deleting.

## 🛠️ Prerequisites
Before you get started, ensure you have the following installed on your machine:

- **Go:** A recent version (e.g., go >= 1.24)
- **Git:** For cloning the repository.
- **PostgresSQL:** The backend and the connectors are ready for PostgreSQL usage.
- **Make:** This API comes with a Makefile containing a set of automations to help run common tasks.

## 📋 Installation
1. Clone the repository:

```
git clone [https://github.com/grglucastr/go-greenlight.git](https://github.com/grglucastr/go-greenlight.git)

cd go-greenlight
```

Install dependencies:

```
go mod tidy
```

## ⚙️ Configuration
Configure PostgreSQL Server
Create the database and enable a few extensions.

### After login to the database server, create a database like bellow
```
CREATE DATABASE greenlight;

# Enable the citext extension
CREATE EXTENSION IF NOT EXISTS citext
```

Create a .envrc file in the root of the project to set your environment variables.

### Example .envrc file content

```
export GREENLIGHT_DB_DSN=postgres://user:password@localhost/greenlight?sslmode=disable
```

## 🏃 Running the Application
### Install Make

**Linux/MacOS**
If you’re on Linux/Mac OS then Make should be available for you.

**Windows**
If you’re on Windows, then probably you’ll have to install make before try to run the API.

```
choco install make
```

### Run
To run the application, execute the following command:

```
make run/api
```

The API will be available at http://localhost:4000.

## 🌐 API Endpoints
This section should detail the API's endpoints. Use a clear format to describe each one.

```[GET] /api/v1/[resource]```
- **Description:** Get a list of all [resources].
- **URL:** http://localhost:[PORT]/api/v1/[resource]
- **Query Parameters:**
- - ```limit```: ```int```, optional - The maximum number of items to return.

- **Response:**
```
{
  "data": [
    {
      "id": 1,
      "name": "Example Resource"
    }
  ]
}
```

```[POST] /api/v1/[resource]```
- **Description:** Create a new [resource].
- **URL:** http://localhost:[PORT]/api/v1/[resource]
- **Request Body:**
```
{
  "name": "New Resource"
}
```

- **Response:** ```201 Created```

## 👋 Contributing
Contributions are always welcome! Please read the ```CONTRIBUTING.md``` file for details on our code of conduct and the process for submitting pull requests.

## 📄 License
This project is licensed under the [License Name] - see the ```LICENSE.md``` file for details.