FROM postgres:10.7
COPY ./airports.csv /
COPY ./init-airports-data.sql /docker-entrypoint-initdb.d/
