# About weak-jwt
This project is an API server intentionally design with common vulnerabilities associated with JWT.

By: Jiacheng Wu (https://github.com/jasonwubz)

# Behind the scenes
The exploits created for this project is a result of analysis done on existing CVEs related to JWT during the time of this writing. Why CVE analysis important? It is because existing CVEs give a broad coverage of vulnerabilities.

## Dependencies
This project is written in Go. You must install go before continuing.
This server uses third party packages: 
- [Echo framework](https://github.com/labstack/echo) - for easy API development 
- [Viper](https://github.com/spf13/viper) - for parsing configuration
- [JWT-GO](https://github.com/dgrijalva/jwt-go) - for the most obvious reasons
- [gosqlite-3](https://github.com/mattn/go-sqlite3) - for simulation of global shared variables related to secret rotation
- [base36](github.com/martinlindhe/base36) - for the weak secret challenge demonstration

Note: If I missed out a library do let me know. I only list direct dependencies.

## Rules of play
Although the intention of this server is to demonstrate vulnerabilities, you can gamify this demo. The rules of play is to avoid reading the source code for answers!

## How to run
After cloning this project, run the following
```sh
go run main.go
```

You can use a tool such as Postman for testing the API endpoints.

## Endpoints

To access the endpoints, use POST method on the url `http://127.0.0.1:1323/api/xxx` where `xxx` is the name of the endpoint. Below is a table of all endpoints. Some are challenges, some demonstrate a concept.

|Challenge/Demo|Endpoint(s)|Description|
|---|---|---|
|None Algorithm|/none|Returns an expired token|
|   |/none-answer|Its vulnerability is that it accepts unsign token. Research about 'none' algorithm.|
|Weak Secret|/weak|Returns an expired token|
|   |/weak-answer|The token is easy to brute-force. Secret is hard-coded in config file. Avoid reading config file and try the challenge. You must be familiar with dictionary attacks.|
|Timing Attack|/timing|Returns an expired token|
|   |/timing-answer|An exploit script is included but the script requires refinement. Sleep time is added to reduce difficulty of the challenge. Understanding of timing attacks is needed.|
|Bad PRG|/math-rand|Returns an expired token|
|   |/math-rand-answer|Secret is a math.random string. Requires a powerful hardware to brute-force all 10 characters that are in base 36 representation.|
|Rotate|/rotate|Returns a valid token|
|   |/rotate-answer|This is not a challenge. It simply demonstrates how to rotate the secret as an added layer of extra security. The solution suggested in this design is not final as the SQLite database is not encrypted. Do not copy the same solution to your production environment. This is a demo only|

# Future plans
If this project gains significant traction, I may continue to enhance it with score tracking and Capture the flag type activities. 

# Disclaimers
Use this project for educational purposes only. I do not take any responsibilities for any damages caused. Also, do not implement the exact demo solutions as they also have other security flaws.