FROM postgres:10.8
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
