#!/bin/bash
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <up/down>"
    exit 1
fi
source ../.env
cd schema
goose postgres "$DB_URL" $1
