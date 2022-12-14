package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SampleApp struct {
	db    *sql.DB
	http  *http.Client
	redis *redis.Client
	mongo *mongo.Client
}

func responseErr(w http.ResponseWriter, err error) {
	data, _ := json.Marshal(struct{ Error string }{Error: err.Error()})
	w.WriteHeader(400)
	_, _ = w.Write(data)
}

func responseOk(w http.ResponseWriter, result interface{}) {
	data, err := json.Marshal(result)
	if err != nil {

	}
	w.WriteHeader(200)
	_, _ = w.Write(data)
}

func main() {
	const proxyAddr = "127.0.0.1:7820"
	dialer, err := newSocksDialer(proxyAddr)
	if err != nil {
		panic(err)
	}
	httpClient, err := newHttpClient(proxyAddr)
	if err != nil {
		panic(err)
	}
	mysqlUrl := os.Getenv("mysql")
	db := connectMysql(mysqlUrl, socksProxy(dialer))
	redisUrl := os.Getenv("redis")
	redis, err := newRedisClient(redisUrl, proxyAddr)
	if err != nil {
		panic(err)
	}
	mongoUrl := os.Getenv("mongo")
	mongo, err := newMongoClient(mongoUrl, proxyAddr)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	app := SampleApp{db: db, http: httpClient, redis: redis, mongo: mongo}

	http.HandleFunc("/http/get", app.httpGet)
	http.HandleFunc("/http/post", app.httpPost)
	http.HandleFunc("/http/status", app.httpStatus)
	http.HandleFunc("/mongo/find", app.mongoFind)
	http.HandleFunc("/mongo/insert", app.mongoInsert)
	http.HandleFunc("/mongo/remove", app.mongoRemove)
	http.HandleFunc("/mysql/select", app.mysqlSelect)
	http.HandleFunc("/mysql/insert", app.mysqlInsert)
	http.HandleFunc("/mysql/update", app.mysqlUpdate)
	http.HandleFunc("/mysql/delete", app.mysqlDelete)
	http.HandleFunc("/mysql/prepared/select", app.mysqlPreparedSelect)
	http.HandleFunc("/mysql/prepared/insert", app.mysqlPreparedInsert)
	http.HandleFunc("/redis/string", app.redisString)
	http.HandleFunc("/redis/zset", app.redisZSet)
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
