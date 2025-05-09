#!/bin/bash

mkdir -p perf_test/results

# Тестирование создания объявлений
bash perf_test/run_create_test.sh
cp perf_test/create_adverts_report.txt perf_test/results/
cp perf_test/create_adverts_report.json perf_test/results/
vegeta plot -type=line < perf_test/create_adverts.bin > perf_test/create_adverts_plot.html

# Тестирование чтения объявлений
bash perf_test/run_read_test.sh
cp perf_test/read_adverts_report.txt perf_test/results/
cp perf_test/read_adverts_report.json perf_test/results/
vegeta plot -type=line < perf_test/read_adverts.bin > perf_test/read_adverts_plot.html

echo "Все тесты выполнены и результаты сохранены в папке perf_test/results."
