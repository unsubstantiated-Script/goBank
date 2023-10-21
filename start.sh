#!/bin/sh

set -e
#This whole file can actually go bye-bye
#echo "run db migration"
#source /app/app.env
#/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up



echo "start this app foo!"
exec "$@"
