# Example

Example of e2e communication between QIX Engine and a PostgreSQL DB using then new GRPC protocol.

### Steps to get the example up and running

Run in a \*nix environment (or Git Bash if on Windows):

```bash
$ cd posgres-image
$ ./build
$ cd ..
$ docker-compose up -d
$ cd reload-runner
$ npm install
$ npm start
```
