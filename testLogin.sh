#!/bin/bash

curl -v -H "Content-Type: application/json" -X POST -d '{ "email": "fsalaman@gmail.com", "password": "password" }' http://localhost:8019/api/login
