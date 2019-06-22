FROM postgres:10.9
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
