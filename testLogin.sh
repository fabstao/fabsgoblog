#!/bin/bash
# *************************************************
# (C) Fabs 2020
# 
# Args:
# $1 - email
# $2 - password
# $3 - url base host (i.e. http://localhost:8019 )
#
# *************************************************

if [ "$#" -le "3" ]; then                                                                                                              
       echo "Usage: $0 <email> <password> <url basehost (i.e. http://localhost:8019)>"
       exit 1                                                                                                                          
fi      

curl -v -H "Content-Type: application/json" -X POST -d '{ "email": "$1", "password": "$2" }' $3/api/login
