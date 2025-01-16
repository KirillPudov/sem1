#!/bin/bash
set -e
go mod tidy
#go build -buildvcs=false -o server

export PGPASSWORD=val1dat0r

psql -U validator -d project-sem-1 -h localhost -c 'CREATE TABLE prices("id" int primary key, "name" varchar(20), "category" varchar(40), "price" float, "create_date" date)'