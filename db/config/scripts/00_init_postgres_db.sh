#!/bin/bash
set -e

# Check if PostgreSQL is already initialized
if [ ! -f "/var/lib/postgresql/data/PG_VERSION" ]; then
    # Switch to postgres user for initialization
    su postgres -c "initdb --username=postgres --pwfile=<(echo \"$POSTGRES_PASSWORD\") /var/lib/postgresql/data"

    cd /var/lib/postgresql/data

    # Generate SSL certificates with proper CA
    # Generate CA key and certificate
    openssl genrsa -out ca.key 4096
    openssl req -new -x509 -days 365 -key ca.key -out root.crt \
        -subj "/CN=postgres-ca"

    # Generate server key and CSR
    openssl genrsa -out server.key 2048
    openssl req -new -key server.key -out server.csr \
        -subj "/CN=postgres"

    # Sign the server certificate with CA
    openssl x509 -req -days 365 -in server.csr \
        -CA root.crt -CAkey ca.key -CAcreateserial \
        -out server.crt \
        -extfile <(echo "subjectAltName = DNS:postgres,DNS:localhost,IP:127.0.0.1")

    # Set proper permissions
    chown postgres:postgres server.key server.crt root.crt ca.key
    chmod 600 server.key ca.key
    chmod 644 server.crt root.crt

    # Start PostgreSQL temporarily to create database and user
    su postgres -c "pg_ctl -D /var/lib/postgresql/data start"
    
    # Wait for PostgreSQL to start
    until su postgres -c "pg_isready"; do
        echo "Waiting for PostgreSQL to start..."
        sleep 1
    done

    # Create database and service user
    su postgres -c "createdb ${POSTGRES_DB}"
    su postgres -c "psql -c \"CREATE USER ${SERVICE_USER} WITH PASSWORD '${SERVICE_PASSWORD}';\""

    # Stop PostgreSQL
    su postgres -c "pg_ctl -D /var/lib/postgresql/data stop"
fi

# Start PostgreSQL with custom config as postgres user
exec su postgres -c "postgres -c config_file=/etc/postgresql/postgresql.conf -c hba_file=/etc/postgresql/pg_hba.conf"