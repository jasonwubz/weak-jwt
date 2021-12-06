# About weak-jwt
This project is an API server intentionally design with common vulnerabilities associated with JWT.

By: Jiacheng Wu (https://github.com/jasonwubz)

## Dependencies
This project is written in Go. You must install go before continuing.
This server uses third party packages: 
- [Echo framework](https://github.com/labstack/echo) - for easy API development 
- [Viper](https://github.com/spf13/viper) - for parsing configuration
- [JWT-GO](https://github.com/dgrijalva/jwt-go) - for the most obvious reasons

## Rules of play
Although the intention of this server is to demonstrate vulnerabilities, you can gamify this demo. The rules of play is to avoid reading the source code for answers!

## How to run
After cloning this project, run the following
```sh
go run main.go
```

You can use a tool such as Postman for testing the API endpoints.