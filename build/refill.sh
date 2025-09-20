#!/usr/bin/env bash

./runsql.sh drop.sql

cd ..
go run cmd/migrate/main.go
cd build

./runsql.sh fill.sql

echo "Refill done."
