#!/bin/sh
go build -o bin/observer app/observer/observer.go 
go build -o bin/node app/node/node.go 
go build -o bin/master app/master/master.go