#!/usr/bin/env bash
protoc --proto_path=src/qlik --go_out=plugins=grpc:src/qlik    grpc_server.proto