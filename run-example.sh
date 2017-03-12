#!/usr/bin/env sh

go install
cd ./example/service

gin-swagger enum
#gin-swagger client -input swagger.json -name service