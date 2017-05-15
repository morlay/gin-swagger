#!/usr/bin/env sh

go install
cd ./example

govendor generate +l

gin-swagger

gin-swagger client --input swagger.json --name service