FROM mysql:8.0.23

COPY ./pkg/server/database/mysqlDB/*.sql /docker-entrypoint-initdb.d/