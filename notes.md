## MINI BANK - Comprehensive Simple Bank App Using GO


### New Connection in Docker pg container with new credentials for app:
 - `docker exec -it f173c5... psql -U su_name -d db_name -c "ALTER ROLE root WITH LOGIN PASSWORD 'password' SUPERUSER;"`
 - Commands:
    * `docker stop container_name/id` - stop given container.
    * `docker ps` - show all running containers.
    * `docker ps -a` - show all containers whether running or not.
    * `docker start container_name/id` - start given container
    
    * `docker exec -it postgres_container /bin/sh` - access postgres on pod's shell
    * `createdb --username=root --owner=root simple_bank` - create db from container's bash
    * `psql simple_bank` - access to the created db. `\q` to exit it
    * `dropdb simple_bank` - to delete db.
    * `exit` - to exit container.
    
    * `docker exec -it postgres_container createdb --username=root --owner=root simple_bank` - direct from outside.
    * `docker exec -it postgres_container psql -U root simple_bank` - to exit container.
    
    * `history | grep "docker run"` - search terminal history for command used to run docker.
    
 
 ### DB Diagram, Schema Design and SQL gen.
 - https://dbdiagram.io
 - migration using `[text](https://github.com/golang-migrate/migrate)`
 - Commands:
    * `migrate create -ext sql -dir db/migration` - create up and down migration files.
    * `migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up` - run the migration in migration folder.
    * `migrate create -ext sql -dir db/migration -seq add_users` - example migration for updating db schema for added users table


### Options for Running SQL querie in Go
    * 1. DATABASE/SQL   
        - Fast and straightforward.
        - Manual mapping SQL fields to variables
        - Easy to make mitakes, not caught until runtime.
        
    * 2. GORM   
        - CRUD functions already implemented, very short production code.
        - Must lean to write queries using gorm's function
        - Runs slowly on high load.
        
    * 3. SQLX   
        - Quite fast and easy to use.
        - Fields mapping via query text and struct tags.
        - Failure won't occur until runtime.
        
    * 4. SQLC   
        - Very fast and easy to use.
        - Automatic code generation.
        - Catch SQL query errors before generating codes.
        - Full support Postgres. MySQL is experimental.


### SQLC for GO with codegen 
    * 1. Codegen: after installing SQLC:
        - `sqlc init` - Create sqlc.yaml file.
        - `docker ps` - Impute .
        - `make sqlc` - generates code from the sql queries.
        - `go mod init github.com/Oyinoye/bank_mini` - initializes a new go project
        - `go mod tidy` - install go dependencies.
        
    * 2. lib pq: Go prostgres driver for Go' database/sql package (github.com/lib/pq).
        - `go get github.com/lib/pq` - install pq.
    
### Testing with GO 
    Normal testing can be done using conditionals. But testify package (www.github.com/stretchr/testify) can be used to get check test result
        - `go get github.com/stretchr/testify` - install testify package.
        - `go test -v ./db/sqlc -run TestTransferTx` - example to run tests verbosely.
        
    Using Mockgen package for mocking
        -  `mockgen -package mockdb -destination db/mock/store.go github.com/Oyinoye/bank_mini/db/sqlc Store`

### Testing for blocking 
    Open psql in two tabs and use the following commands:
    
        - `Begin;` - Begin transaction on either instance
        - `SELECT * FROM accounts WHERE id = 1;` - Works on both tabs (non-blocking). Compare with next:
        - `SELECT * FROM accounts WHERE id = 1 FOR UPDATE;` - Now select won't run in other tab (blocked). Has to wait for first transaction to commit or rollback.
        - `ROLLBACK;` - Rollback the transasction.
        - `COMMIT;` - Commit the transaction.
        
        Now consistency can be guaranteed:
        - `UPDATE accounts SET balance = 500 WHERE id = 1;` - We can run this now and then commit to unblock second transaction.
        
        Adding 'NO KEY' will enable postgres free up ID column since it won't be updated. This prevent deadlocks as in example below:
        - `-- name: GetAccountForUpdate :one
            SELECT * FROM accounts
            WHERE id = $1 LIMIT 1
            FOR NO KEY UPDATE;`

### API with GO 
    Gin is a robust, popular and widely used package.
        -  Check Gin docs for more info.
               
### Loading from .env
    Viper is a great package for configuration. Search and installation.
        - 


### Security (PASETO vs JWT)
    Paseto is more advantageous than JWT as it is more secure and mitigates security risks of JWT.
        - Check docs for more info.


### Use DBdocs to generate documentation.
    Installation and usage instructions- [dbdocs](https://dbdocs.io/)
        - Put in the SQL code and it's good to go.
        - To acces us `dbdocs login` command as specified in the docs.
        - `dbdocs build doc/db.dbml` - will generate the documentation
        - `dbdocs pasword --sset <password> --project <projectname>` sets the password.
        
    Using DBML cli:
        - `npm install -g @dbml/cli`
        - `dbml2sql --postgres -o doc/schema.sql doc/db.dbml`


## gRPC with Go

Remote Procedure calls have the features of being highly performant (built on HTTP/2 and supports binary framing, multiplexing, header compression, bidirectional communication), allows strong API contract and Automatic code generation. Read more on the [documentation](https://grpc.io/docs)

### Types of gRPC
    - Unary gRPC
    - Client streaming gRPC
    - Server streamig gRPC
    - Bidirectional streaming gRPC


### Quickstart
Quickstart with Go instructions [found here](https://grpc.io/docs/languages/go/) 
    - `protoc --version` - check if Protocol buffer compiler, protoc is installed...
    - `protoc-gen-go --version` - check if Protocol buffer compiler, protoc is installed...
    
Installation command
    - `go get google.golang.org/grpc` - Run at go project route.
    - `brew install protobuf`

Additional prerequisites for generating code:
If these imports are appearing in code you are automatically generating from a .proto file, you will also need the protobuf compiler and the Go code generation plugins
    - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
    - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
    
After running, to ensure everything syncs correctly, also run:
    - `go mod tidy`
    
### gRPC Client
Evans is a good rpc client.
    - `protoc --version` - check if Protocol buffer compiler, protoc is installed...
    - `protoc-gen-go --version` - check if Protocol buffer compiler, protoc is installed...
Install Evans client for grpc API testing at [Evans' github](https://github.com/ktr0731/evans)
    - `evans --host localhost --port 9090 -r repl` or `evans -r repl`
Evans commands
    - `show service` - displays rpc services.
    - `show package` - shows the available packages
    - `package pb` - package has to be selected. Since tis proto usess pb, service is pb.SimpleBank.
    - `service SimpleBank` - select SimpleBank service.
    - `call CreateUser` - call the create user function.
    - `exit` - exits the console.

### gRPC Gateway
A plugin of protobuf compiler. It allows both grpc and normal http clients to connect respectively to the GRPC server and Gateway resçectively from a single code. It generates HTTP proxy codes from protobuf definition. Conversion takes place between Gateway and GRPC (in-process translation).
This in-process translation works for only unary rpc. A separate proxy server is needed for streaming as well as a combination of unary and streaming. 
[Read more.](https://github.com/grpc-ecosystem/grpc-gateway)

    - `protoc --version` - check if Protocol buffer compiler, protoc is installed...
    - `protoc-gen-go --version` - check if Protocol buffer compiler, protoc is installed...
    - `protoc-gen-grpc-gateway --help to show other options that can be added (check protoc command in make file)çç
    
### gRPC Gateway wagger Documentation
Documentation instructions found in the link from previous section (grpc gateway). 
Serve swagger documentation using static file server. [Statik](https://github.com/rakyll/statik) library is a great option with many advantages over just reading and serving the files as is.

### Other tools used
Pgx - PostgreSQL Driver and Toolkit.
    - `$ go get github.com/jackc/pgx/v5` - check [PGX docs](https://github.com/jackc/pgx) 
zerologs - Good logs for structured logs in prod [here](https://github.com/rs/zerolog).
    - `go get -u github.com/rs/zerolog/log`

