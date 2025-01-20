#!/bin/bash
set -x

apt update && apt install zip curl -y

go mod tidy

go build -buildvcs=false -o server

chmod +x server

echo "DB_HOST=localhost" >> .env
echo "DB_PORT=5432"  >> .env
echo "DB_USER=validator"  >> .env
echo "DB_PASSWORD=val1dat0r" >> .env
echo "DB_NAME=project-sem-1"  >> .env

export PGPASSWORD=val1dat0r

psql -U validator -d project-sem-1 -h localhost -c 'CREATE TABLE prices("id" serial primary key, "name" varchar(20), "category" varchar(40), "price" float, "create_date" timestamp)'