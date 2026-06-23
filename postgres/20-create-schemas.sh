#!/usr/bin/env bash

set -e

export PGPASSWORD=test
psql -U postgres -d films -c "GRANT ALL PRIVILEGES ON SCHEMA public TO program;"
psql -U postgres -d cinema -c "GRANT ALL PRIVILEGES ON SCHEMA public TO program;"
psql -U postgres -d tickets -c "GRANT ALL PRIVILEGES ON SCHEMA public TO program;"
