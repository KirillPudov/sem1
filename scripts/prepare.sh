#!/bin/bash
set -x
go mod tidy
export GOBIN=/usr/local/bin
go install -buildvcs=false

chmod +x /usr/local/bin/project_sem

export PGPASSWORD=val1dat0r

psql -U validator -d project-sem-1 -h localhost -c 'CREATE TABLE prices("id" int primary key, "name" varchar(20), "category" varchar(40), "price" float, "create_date" date)'