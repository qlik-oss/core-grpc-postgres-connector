FROM postgres:10.14
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
