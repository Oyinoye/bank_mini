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

