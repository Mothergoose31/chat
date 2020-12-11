
package main

import (
	"database/sql"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

type database struct {
	db        *sql.DB
	insertban chan *dbInsertBan
	deleteban chan *dbDeleteBan
	sync.Mutex
}

type dbInsertBan struct {
	uid       Userid
	targetuid Userid
	ipaddress *sql.NullString
	reason    string
	starttime time.Time
	endtime   *mysql.NullTime
	retries   uint8
}

type dbDeleteBan struct {
	uid Userid
