#!/bin/bash

echo "Запуск нагрузочного тестирования на чтение объявлений..."
echo "GET http://localhost:8080/api/adverts/\$ID" | vegeta attack -duration=60s -rate=1000 > perf_test/read_adverts.bin
vegeta report -type=text < perf_test/read_adverts.bin > perf_test/read_adverts_report.txt
vegeta report -type=json < perf_test/read_adverts.bin > perf_test/read_adverts_report.json
