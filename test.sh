#!/bin/bash

printf "===\nCorrect headers:\n"
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -H "Authorization: Bearer 123" localhost:3000

printf "\n\n===\nUnauthorized headers:\n"
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -H "Authorization: Bearer failiure" localhost:3000
