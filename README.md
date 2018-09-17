# Example Postgres GRPC Connector

The Example Postgres GRPC Connector shows how to load data into Qlik Associative Engine from Postgres using a
dockerized connector built in Golang. It streams the data asynchronously using go channels though
the following components before sending it onto Qlik Associative Engine.
* postgres_reader - reads the data from the database into reasonably sized SQL data chunks.
* async_translator - takes the SQL data chunks and translates them into GRPC data chunks.
* async_stream_writer - takes the GRPC data chunks and writes them onto the GRPC stream.

The reason for the division is to be able to utilize multiple CPU cores to process the different stages simultaneously.

## Example

The `/example` directory defines a simple stack of services using docker-compose:
* Qlik Associative Engine
* Postgres GRPC Connector
* Postgres Database
* Node Test Runner (only used for automated testing)

The script in [example/reload-runner](example/reload-runner) is used to instruct Qlik Associative Engine to load the example
data (originally defined in [example/postgres-image/airports.csv](example/postgres-image/airports.csv))
using the connector.

### Steps to run the example

Run in a \*nix environment (or Git Bash if on Windows), note that you must accept the
[Qlik Core EULA](https://core.qlik.com/eula/) by setting the `ACCEPT_EULA`
environment variable:

```bash
$ cd example
$ ACCEPT_EULA=yes docker-compose up -d --build
$ cd reload-runner
$ npm install
$ npm start
```
