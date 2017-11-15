# Example Postgres GRPC Connector

The Example Postgres GRPC Connector shows how to load data into QIX Engine from Postgres using a
dockerized connector built in Golang. It streams the data asynchronously using go channels though
the following components before sending it onto QIX Engine.
* postgres_reader - reads the data from the database into reasonably sized SQL data chunks.
* async_translator - takes the SQL data chunks and translates them into GRPC data chunks.
* async_stream_writer - takes the GRPC data chunks and writes them onto the GRPC stream.

The reason for the division is to be able to utilize multiple CPU cores to process the different stages simultaneously.

## Example

The `/example` directory defines a simple stack of services using docker-compose:
* QIX Engine
* Postgres GRPC Connector
* Postgres Database
* Node Test Runner (only used for automated testing)

The script in [example/reload-runner](example/reload-runner) is used to instruct QIX Engine to load the example
data (originally defined in [example/postgres-image/airports.csv](example/postgres-image/airports.csv))
using the connector.

### Steps to run the example

Run in a \*nix environment (or Git Bash if on Windows):

```bash
$ cd example
$ docker-compose up -d --build
$ cd reload-runner
$ npm install
$ npm start
```
