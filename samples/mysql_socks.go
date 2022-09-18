package main

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/net/proxy"
	"net"
)

type Option func(*sql.DB)

func socksProxy(dialer proxy.Dialer) Option {
	return func(d *sql.DB) {
		mysql.RegisterDialContext("tcp", func(ctx context.Context, addr string) (net.Conn, error) {
			return dialer.Dial("tcp", addr)
		})
	}
}
func newSocksDialer(addr, user, password string) (proxy.Dialer, error) {
	return proxy.SOCKS5("tcp", addr, &proxy.Auth{User: user, Password: password}, proxy.Direct)
}
func connectMysql(dsn string, opts ...Option) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		opt(db)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
