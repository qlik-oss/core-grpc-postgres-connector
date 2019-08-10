FROM postgres:10.10
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
