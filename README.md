# Project Tinc1

_This is a basic API built with Go, Gin and Gorm_

### Installation 📋

_This implementation has been written and tested using Go 1.16_

You need to install Go in order to run this Project.
https://golang.org/doc/install

Install Gin, Gorm and the JWT middleware:

    $ go get -u github.com/gin-gonic/gin
    $ go get -u gorm.io/gorm
    $ go get github.com/dgrijalva/jwt-go
    
Once within the root directory, install dependencies:

    /root$ go mod tidy
    
### Running the API 📋

Within the root directory:

    /root$ go run api.go
    
The API will run on port 8000
    
### Running tests 📋

Run:

    /root$ go test

### Endpoints 📋

- Login endpoint:

    localhost:8000/login

    Receives "email" and "password" as form data. Returns a JWT and a user object if auth is ok. 401 otherwise.

The following endpoints require an authentication header (BEARER TOKEN)

- Inboundfiles endpoint:

    localhost:8000/api/inboundfiles

    Receives a JSON in the request body:
    {
        "id": "YOURKEY"
    }

- Outboundfiles endpoint:

    localhost:8000/api/outboundfiles

    Receives a JSON in the request body:
    {
        "id": "YOURKEY"
    }