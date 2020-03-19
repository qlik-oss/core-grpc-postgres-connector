FROM postgres:10.11
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
