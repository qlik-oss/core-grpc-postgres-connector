FROM postgres:10.6
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
