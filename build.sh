#!/bin/bash

echo "get packages..."
go get github.com/emicklei/go-restful
go get github.com/gambol99/go-marathon
go get github.com/hashicorp/hcl
echo "get packages finished"

echo "build..."
go build -o monitor
echo "build finished"
