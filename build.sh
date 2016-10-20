#!/bin/bash

echo "get packages..."
go get github.com/gambol99/go-marathon
go get github.com/hashicorp/hcl
go get github.com/bmizerany/assert
echo "get packages finished"

echo "build..."
go build -o monitor
echo "build finished"
