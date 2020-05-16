FROM postgres:10.13
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
