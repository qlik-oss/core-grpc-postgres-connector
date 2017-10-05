# GRPC Go Custom connector without QVX

A POC to validate that it is possible to write a custom connector without the qvx protocol using GRPC.

## Example

[An example of a reload with the QIX Engine docker using GRPC to a PostgreSQL connector](/example/README.md)

### Get it working with Qlik Sense (For internal use and nothing public that supports this)

1. Install tools using install-tools.sh 
2. Run "dep ensure" to download/update dependencies
3. go run src\server\postgres_grpc_server.go
4. Move GrpcPostgres into C:\Users\username\AppData\Local\Programs\Common Files\Qlik\Custom Data (change to your username)
5. Start QIX Engine with command line args: -P 9076 --WSPath "C:\Users\username\AppData\Local\Programs\Qlik\Sense\Client" --MigrationPort 9074 --DataPrepPort 9072 --NPrintingPort 9073 -S EnableGrpcCustomConnectors=1 -S GrpcConnectorPlugins="GrpcPostgres,host:50051" -S EnableConnectivityService=0 -N
(change to your own username)
