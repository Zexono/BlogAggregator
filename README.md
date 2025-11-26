# Blog Aggregator

```
Learning Goals
    Learn how to integrate a Go application with a PostgreSQL database
    Practice using your SQL skills to query and migrate a database (using sqlc and goose, two lightweight tools for typesafe SQL in Go)
    Learn how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database
```
# How to use it (not finish yet)

you need Postgres and Go installed to run the program.

Linux / WSL for Postgres 
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```
Make sure you have Go installed 
```
https://go.dev/doc/install 
```

then Run 
```
go install github.com/Zexono/blogagg/@latest.
```
To install the gator

After that you need to create a config file
Create a .gatorconfig.json
```
{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user_name":"current_user_name"}
```

# Command



