# Example Postgres GRPC Connector

The Example Postgres GRPC Connector shows how to load data into QIX Engine from Postgres using a 
dockerized connector built in Golang.

# Example

The `/example` directory defines a simple stack consisting of
* QIX Engine
* Postgres GRPC Connector
* Postgres Database

Using the reload runner in [example/reload-runner](example/reload-runner) you can trigger a reload the QIX Engine that 
loads an example table (originally defined in
[example/postgres-image/airports.csv](example/postgres-image/airports.csv)). 
 
### Steps to get the example up and running

Run in a \*nix environment (or Git Bash if on Windows):

```bash
$ cd example
$ docker-compose up -d --build
$ cd reload-runner
$ npm install
$ npm start
```
