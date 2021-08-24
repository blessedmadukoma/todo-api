## TODO API - GOLANG AND POSTGRESQL

# A Todo REST API built using Golang and Postgresql

# To run:

1. Clone the repo
2. Create a .env file with the following data:

- DB_USER (database username: default is postgres)
- DB_PASSWORD (database password)
- DB_NAME=todos
- DIALECT (dialect being the driver: postgres)
- DB_PORT (database port)
- DB_HOST (database host: use localhost or 127.0.0.1)
- PORT (8000 or 8080)

3. Run `go run main.go` to start the application
   The endpoints:

- /todos :GET method -> get all todos
- /todo/create :POST method -> create a new todo
- /todo/{id} :GET method -> get a single todo
- /todo/{id} :PUT method -> update an existing todo
- /todo/{id} :DELETE method -> delete an existing todo
- /todo/search/{key} :POST method -> return all existing todo items with a keyword (any text involved)
