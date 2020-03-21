FROM postgres:10.12
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
