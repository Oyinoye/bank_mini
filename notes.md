## MINI BANK - Comprehensive Simple Bank App Using GO


### New Connection in Docker pg container with new credentials for app:
 - `docker exec -it f173c5... psql -U su_name -d db_name -c "ALTER ROLE root WITH LOGIN PASSWORD 'password' SUPERUSER;"`
 
 ### DB Diagram, Schema Design and SQL gen.
 - https://dbdiagram.io
 - migration using `[text](https://github.com/golang-migrate/migrate)`
