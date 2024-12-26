#!/bin/bash
# Установка Vegeta

curl -LO https://github.com/tsenart/vegeta/releases/download/v12.8.4/vegeta-12.8.4-linux-amd64.tar.gz
tar -xzf vegeta-12.8.4-linux-amd64.tar.gz
sudo mv vegeta /usr/local/bin/
rm vegeta-12.8.4-linux-amd64.tar.gz
