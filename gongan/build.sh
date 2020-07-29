#!/bin/sh
go build -o 1put_info 1put_info.go type.go
go build -o 2get_info 2get_info.go type.go
go build -o 3get_history 3get_history.go type.go



