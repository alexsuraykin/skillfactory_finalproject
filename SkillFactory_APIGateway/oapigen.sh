#!/usr/bin/env sh

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,gin,spec,client -package oapi -o internal/api/oapi/oapi.gen.go swagger.yaml