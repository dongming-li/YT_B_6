#!/bin/bash

echo "STARTING USER TESTS"

echo -ne "\ttesting create\t"
STATUS="$(curl -s -o /dev/null -w '%{http_code}' \
        -X POST -F 'username=test' -F 'password=test' \
        -F 'firstname=firsttest' -F 'lastname=lasttest' \
        -F 'email=test@test.com' localhost:8080/user)"
if [ ! $STATUS -eq 201 ]; then
  echo "FAIL: status=$STATUS"
else
  echo "pass"
fi

echo -ne "\tgetting token\t"
response=$(curl -sb -X POST -H "Accept: application/json Content-Type: application/json" \
         -d '{"username":"test@test.com","password":"test"}' localhost:8080/login)
if [[ $response =~ ({\"expire\":\"[T:0-9\-]+\",\"token\":\"([a-zA-z0-9.\-]+)\"}) ]]; then
  token=${BASH_REMATCH[2]}
  echo "ok"
else
  echo "FAIL: response=$response"
fi

echo -ne "\ttesting list\t"
STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X GET \
        -H "Authorization: Bearer $token" localhost:8080/user)
if [ ! $STATUS -eq 200 ]; then
  echo "FAIL: status=$STATUS"
else
  echo "pass"
fi


echo -ne "\ttesting read\t"
STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X GET \
        -H "Authorization: Bearer $token" localhost:8080/user/3)
if [ ! $STATUS -eq 200 ]; then
  echo "FAIL: status=$STATUS"
else
  echo "pass"
fi

echo -ne "\ttesting update\t"
STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X PUT -F 'firstname=testname' \
        -H "Authorization: Bearer $token" localhost:8080/user/3)
if [ ! $STATUS -eq 200 ]; then
  echo "FAIL: status=$STATUS"
else
  echo "pass"
fi

echo -ne "\ttesting delete\t"
STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X DELETE \
        -H "Authorization: Bearer $token" localhost:8080/user/3)
if [ ! $STATUS -eq 200 ]; then
  echo "FAIL: status=$STATUS"
else
  echo "pass"
fi
