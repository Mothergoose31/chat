
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



	var db = &database{
		insertban: make(chan *dbInsertBan, 10),
		deleteban: make(chan *dbDeleteBan, 10),
	}
	
	func initDatabase(dbtype string, dbdsn string) {
		var err error
		conn, err := sql.Open(dbtype, dbdsn)
		if err != nil {
			B("Could not open database: ", err)
			time.Sleep(time.Second)
			initDatabase(dbtype, dbdsn)
			return
		}
		err = conn.Ping()
		if err != nil {
			B("Could not connect to database: ", err)
			time.Sleep(time.Second)
			initDatabase(dbtype, dbdsn)
			return
		}
	
		db.db = conn
		go db.runInsertBan()
		go db.runDeleteBan()
	}
	
	func (db *database) getStatement(name string, sql string) *sql.Stmt {
		db.Lock()
		stmt, err := db.db.Prepare(sql)
		db.Unlock()
		if err != nil {
			D("Unable to create", name, "statement:", err)
			time.Sleep(100 * time.Millisecond)
			return db.getStatement(name, sql)
		}
		return stmt
	}
	
	func (db *database) getInsertBanStatement() *sql.Stmt {
		return db.getStatement("insertBan", `
			INSERT INTO bans
			SET
				userid         = ?,
				targetuserid   = ?,
				ipaddress      = ?,
				reason         = ?,
				starttimestamp = ?,
				endtimestamp   = ?
		`)
	}