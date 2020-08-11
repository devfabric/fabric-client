#!/bin/sh
go build -o 1register_org 1register_org.go type.go
go build -o 2get_orginfo 2get_orginfo.go type.go
go build -o 3put_info 3put_info.go type.go
go build -o 4put_baseinfo 4put_baseinfo.go type.go
go build -o 5update_baseinfo 5update_baseinfo.go type.go
go build -o 6put_favorinfo 6put_favorinfo.go type.go
go build -o 7update_favorinfo 7update_favorinfo.go type.go
go build -o 8confirm_baseinfo 8confirm_baseinfo.go type.go
go build -o 9get_info 9get_info.go type.go
go build -o 10query_bycond 10query_bycond.go type.go
go build -o 11query_simple 11query_simple.go type.go
go build -o 12get_history 12get_history.go type.go

go build -o 19set_sharedkv 19set_sharedkv.go type.go




