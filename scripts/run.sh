#!/bin/bash
set -x

exec env POSTGRES_HOST=localhost POSTGRES_PORT=5432 POSTGRES_DB=project-sem-1 POSTGRES_USER=validator POSTGRES_PASSWORD=val1dat0r project_sem