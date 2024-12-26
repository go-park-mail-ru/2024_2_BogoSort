#!/bin/bash

API_URL="http://localhost:8080/api/adverts"
TOTAL=100000
BATCH=1000

for ((i=1; i<=TOTAL; i+=BATCH)); do
    echo "Создание объявлений: $i до $((i + BATCH - 1))"
    for ((j=1; j<=BATCH && i+j-1<=TOTAL; j++)); do
        TITLE="Объявление $((i + j -1))"
        DESCRIPTION="Описание объявления $((i + j -1))"
        PRICE=$((RANDOM % 10000 + 100))
        LOCATION="Москва"
        STATUS="active"
        
        curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "{\"title\":\"$TITLE\",\"description\":\"$DESCRIPTION\",\"price\":$PRICE,\"location\":\"$LOCATION\",\"status\":\"$STATUS\"}" &
    done
    wait
done
