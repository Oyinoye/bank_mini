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
    
