#!/usr/bin/env sh

go install
cd ./example

gin-swagger enum
#gin-swagger error
gin-swagger
gin-swagger client --input swagger.json --name service