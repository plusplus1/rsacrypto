#!/bin/bash

go build -mod vendor -ldflags "-s -w" -o bin/parallel_rsa main.go
