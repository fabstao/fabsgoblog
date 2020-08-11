#!/bin/bash
# *************************************************
# (C) Fabs 2020
# 
# Args:
# $1 - Token
# $2 - url base host (i.e. http://localhost:8019 )
#
# *************************************************

if [ "$#" -le "3" ]; then                                                                                                              
       echo "Usage: $0 <token> <url basehost (i.e. http://localhost:8019)>"
       exit 1                                                                                                                          
fi      

curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $1" \
 -d '{"titulo": "Prueba API", "texto":"![](https://i.imgur.com/FwLPvQa.png) \n # Es una prueba API \n ## REST \n * Esta entrada fue creada usando cURL \n * JSON \n * REST"}' \
 $2/sapi/ 
