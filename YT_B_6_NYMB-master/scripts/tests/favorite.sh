#!/bin/bash

echo "STARTING FAVORITE TESTS"

echo -ne "getting token ... "
response=$(curl -sb -X POST -H "Accept: application/json Content-Type: application/json" \
         -d '{"username":"admin@admin.com","password":"admin"}' localhost:8080/login)
if [[ $response =~ ({\"expire\":\"[T:0-9\-]+\",\"token\":\"([a-zA-z0-9.\-]+)\"}) ]]; then
  token=${BASH_REMATCH[2]}
  echo "ok"
else
  echo "FAIL: response=$response"
fi
