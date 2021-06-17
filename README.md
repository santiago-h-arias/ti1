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

Within the root directory:

    /root$ go run api.go
