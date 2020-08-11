#!/bin/bash

curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $1" \
 -d '{"titulo": "Prueba API", "texto":"![](https://i.imgur.com/FwLPvQa.png) # Es una prueba API \n ## REST \n * JSON \n * REST"}' \
 localhost:8019/sapi/ 
