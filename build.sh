#! /bin/sh
# run goimports
find ./ -name "*.go" | xargs goimports -w=true

# run go build
go build
