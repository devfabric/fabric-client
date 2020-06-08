#!/bin/sh
go build -o 1check 1check.go
go build -o 2put_baseinfo 2put_baseinfo.go type.go
go build -o 3update_baseinfo 3update_baseinfo.go type.go
go build -o 4put_info 4put_info.go type.go
go build -o 5update_info 5update_info.go type.go
go build -o 6get_info 6get_info.go type.go
go build -o 7query_bycond 7query_bycond.go type.go
go build -o 8query_simple 8query_simple.go type.go
go build -o 9get_history 9get_history.go type.go







