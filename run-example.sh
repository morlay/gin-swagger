#!/usr/bin/env sh

go install
cd ./example/service
gin-swagger
