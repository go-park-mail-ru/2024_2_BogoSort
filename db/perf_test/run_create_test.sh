#!/bin/bash

echo "Запуск нагрузочного тестирования на создание объявлений..."
echo "POST http://localhost:8080/api/adverts" | vegeta attack -duration=60s -rate=1000 > perf_test/create_adverts.bin
vegeta report -type=text < perf_test/create_adverts.bin > perf_test/create_adverts_report.txt
vegeta report -type=json < perf_test/create_adverts.bin > perf_test/create_adverts_report.json
