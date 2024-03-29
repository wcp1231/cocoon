package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
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
func newSocksDialer(addr string) (proxy.Dialer, error) {
	return proxy.SOCKS5("tcp", addr, &proxy.Auth{User: "", Password: ""}, proxy.Direct)
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

func newHttpClient(addr string) (*http.Client, error) {
	dialSocksProxy, err := newSocksDialer(addr)
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
func newRedisClient(redisAddr, proxyAddr string) (*redis.Client, error) {
	dialSocksProxy, err := newSocksDialer(proxyAddr)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialSocksProxy.Dial(network, addr)
		},
	}), nil
}

func newMongoClient(mongoUrl, proxyAddr string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(mongoUrl)
	dialSocksProxy, err := newSocksDialer(proxyAddr)
	if err != nil {
		return nil, err
	}
	dialer := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialSocksProxy.Dial(network, address)
	}
	opts.SetDialer(topology.DialerFunc(dialer))
	return mongo.Connect(context.Background(), opts)
}
