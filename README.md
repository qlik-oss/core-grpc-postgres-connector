# GRPC Go Custom connector without QVX

A POC to validate that it is possible to write a custom connector without the qvx protocol.

## Start Grpc Postgres go services

1. Install tools using install-tools.sh 
1. Run "dep ensure" to download/update dependencies
2. go run src\server\postgres_grpc_server.go
3. Move GrpcPostgres into C:\Users\seb\AppData\Local\Programs\Common Files\Qlik\Custom Data (change to your username)

## Build engine branch hsv/grpccol

1. Select build target Sense Release/x64
2. Make QlikMain default start up project
3. Command line args: -P 9076 --WSPath "C:\Users\seb\AppData\Local\Programs\Qlik\Sense\Client" --MigrationPort 9074 --DataPrepPort 9072 --NPrintingPort 9073 -S EnableGrpcCustomConnectors=1 -S GrpcConnectorPlugins="GrpcPostgres,selun-seb1.qliktech.com:50051" -S EnableConnectivityService=0 -N
(change user specifics seb to your own username)
4. Disable engine in C:\Users\seb\AppData\Local\Programs\Qlik\Sense services.conf-file in desktop
5. Run engine
6. Start desktop

