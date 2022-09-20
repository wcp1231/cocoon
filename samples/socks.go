package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"runtime"
	"time"
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

func newHttpClient(addr, user, password string) (*http.Client, error) {
	dialSocksProxy, err := newSocksDialer(addr, user, password)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating SOCKS5 proxy. %v", err))
	}
	if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
		dialContext := contextDialer.DialContext
		return &http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           dialContext,
				MaxIdleConns:          10,
				IdleConnTimeout:       60 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		}, nil
	}
	return nil, errors.New("Failed type assertion to DialContext")
}
