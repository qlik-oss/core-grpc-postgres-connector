FROM postgres:9.6.5
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
