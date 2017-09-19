#!/usr/bin/env bash
protoc --proto_path=qlik --go_out=plugins=grpc:qlik    grpc_server.proto