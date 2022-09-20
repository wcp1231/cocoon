package main

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SampleApp struct {
	db   *sql.DB
	http *http.Client
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
	dialer, err := newSocksDialer("127.0.0.1:7820", "", "")
	if err != nil {
		panic(err)
	}
	httpClient, err := newHttpClient("127.0.0.1:7820", "", "")
	if err != nil {
		panic(err)
	}
	mysqlUrl := os.Getenv("mysql")
	db := connectMysql(mysqlUrl, socksProxy(dialer))

	rand.Seed(time.Now().UnixNano())
	app := SampleApp{db: db, http: httpClient}

	http.HandleFunc("/http/get", app.httpGet)
	http.HandleFunc("/http/post", app.httpPost)
	http.HandleFunc("/mysql/select", app.mysqlSelect)
	http.HandleFunc("/mysql/insert", app.mysqlInsert)
	http.HandleFunc("/mysql/update", app.mysqlUpdate)
	http.HandleFunc("/mysql/delete", app.mysqlDelete)
	http.HandleFunc("/mysql/prepared/select", app.mysqlPreparedSelect)
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
