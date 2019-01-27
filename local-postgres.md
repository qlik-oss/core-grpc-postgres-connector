# Guide for using a local postgres database with QAE and core-grpc-postgres-connector

## First step: Setting up the example

Verify that it´s possible to spin up the docker environment and run the example as described in the [README.md](./README.md).

After running all the steps a sample of the airports table should be visible in the console.

## Second step: Dump the airports table

This section can be skipped and the [airports.tar](./example/postgres-image/airports.tar) be used in the next step

1. Get the ID of the database container `docker ps`
2. Connect to the database container `docker exec -it <CONTAINER ID> sh`
3. Execute `pg_dump --host localhost --port 5432 --username postgres --format tar --verbose --file airports.tar --table public.airports postgres` in the shell of the database container
4. Exit that container
5. Copy the tar from the container with `docker cp <CONTAINER ID>:/airports.tar airports.tar`

## Third step: Setting up local postgres

1. Use [pdAdmin](https://www.pgadmin.org/) to connect to your local (on the development computer). (localhost and port 5432 by default)
2. Create a new database (rightClick **Create > Databases** under **Servers** and name it **test**)
3. RightClick the **test** database and select **Restore...**
4. Locale the airports.tar and press **Restore**

By default the local postgres is only listening on the localhost. Follow the guide [Configure PostgreSQL to allow remote connection](https://blog.bigbinary.com/2016/01/23/configure-postgresql-to-allow-remote-connection.html)

## Fourth step: Adjusting the core-grps-postgres-connector example

The next step will only work for the docker-desktop-for(mac/windows). There is an alias/route for connecting to the localhost on the development computer. `host.docker.internal` will map back to `localhost` of the development computer / host.

The only thing we need to change is the **qConnectionString** that contains the information of what host the **core-grpc-postgres-connector** shall connect to and the name of the new database that we created **test**.

Change the host part of the **qConnectionString** to point at `host.docker.internal` (if other port than the default 5432 is used for postgres on the localhost that has to be changed as well).

Also change the database name to **test** (from postgres).

./example/reload-runner/index.js
```
 qConnectionString:
        'CUSTOM CONNECT TO "provider=postgres-grpc-connector;host=host.docker.internal;port=5432;database=test"',
```

Rerun the test from the first step `npm start` and you should get the same result since we only moved the same table to another postgres.

## Debug
If you didn´t get the same result as the first step we need to do some investigation and the best way of doing this is to try to connect to the database from inside the **core-grps-postgres-connector** container

1. Get the ID of the core-grps-postgres-connector container `docker ps`
2. Connect to the database container `docker exec -it <CONTAINER ID> sh`
3. Since we want to keep the containers as slim as possible no extra tools is installed by default but it´s easy to do. The container is built using *'Alpine** that has its own packagesystem. Just run the `apk add postgresql-client` to install the tool we need.
4. Next we will try to connect to the database on the development computer by running `psql -h host.docker.internal -U postgres`. If you get an error here the container couldn't connect to the database and therefore wasn't able to load the airports table.

Other good commands to try is:
- `\l` list databases
- `\c test`  change to the test database
- `\dt+` list tables

It´s also possible to execute SQL statemets directly (you may need to change to correct database as described above)
- `SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_SCHEMA='public';`
- `SELECT * FROM airports;`

## Using corectl

Instead of using the programatic approach the tool [corectl](https://github.com/qlik-oss/corectl/releases) can be used to desribe the app you want to create.

Start by creating a configuration file for corectl (in this case against the local postgres):
```
engine: localhost:19076
script: ./script.qvs
app: usingCorectl.qvf

connections:
  postgres-grpc-connector:
      type: postgres-grpc-connector
      username: postgres
      password: postgres
      settings:
        database: test
        host: host.docker.internal
        port: 5432
```

Next you need to add the [Qlik script](https://core.qlik.com/services/qix-engine/script_reference/introduction/) to read from the database:
```
lib connect to 'postgres-grpc-connector';

Airports:
sql SELECT rowID,Airport,City,Country,IATACode,ICAOCode,Latitude,Longitude,Altitude,TimeZone,DST,TZ, clock_timestamp() FROM airports ORDER BY Airport;
```
