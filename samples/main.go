package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SampleApp struct {
	db *sql.DB
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

func (s *SampleApp) mysqlSelect(w http.ResponseWriter, req *http.Request) {
	res, err := s.db.Query("SELECT * FROM samples")
	if err != nil {
		responseErr(w, err)
		return
	}
	var samples []*MysqlSample
	for res.Next() {
		sample, err := NewMysqlSampleFromScan(res)
		if err != nil {
			responseErr(w, err)
			return
		}
		samples = append(samples, sample)
	}
	responseOk(w, samples)
}
func (s *SampleApp) mysqlInsert(w http.ResponseWriter, req *http.Request) {
	sample := NewRandomMysqlSample()
	sql := sample.InsertSQL()
	res, err := s.db.Exec(sql)
	if err != nil {
		responseErr(w, err)
		return
	}
	var lastInsertId, rowAffected int64
	if lastInsertId, err = res.LastInsertId(); err != nil {
		responseErr(w, err)
		return
	}
	if rowAffected, err = res.RowsAffected(); err != nil {
		responseErr(w, err)
		return
	}
	responseOk(w, struct {
		LastInsertId int64
		RowsAffected int64
	}{
		LastInsertId: lastInsertId,
		RowsAffected: rowAffected,
	})
}
func (s *SampleApp) mysqlUpdate(w http.ResponseWriter, req *http.Request) {
	responseErr(w, errors.New("not yet implement"))
}
func (s *SampleApp) mysqlDelete(w http.ResponseWriter, req *http.Request) {
	res, err := s.db.Exec("DELETE FROM samples")
	if err != nil {
		responseErr(w, err)
		return
	}
	var lastInsertId, rowAffected int64
	if lastInsertId, err = res.LastInsertId(); err != nil {
		responseErr(w, err)
		return
	}
	if rowAffected, err = res.RowsAffected(); err != nil {
		responseErr(w, err)
		return
	}
	responseOk(w, struct {
		LastInsertId int64
		RowsAffected int64
	}{
		LastInsertId: lastInsertId,
		RowsAffected: rowAffected,
	})
}

func main() {
	dialer, err := newSocksDialer("127.0.0.1:7820", "", "")
	if err != nil {
		panic(err)
	}
	mysqlUrl := os.Getenv("mysql")
	db := connectMysql(mysqlUrl, socksProxy(dialer))

	rand.Seed(time.Now().UnixNano())
	app := SampleApp{db: db}

	http.HandleFunc("/mysql/select", app.mysqlSelect)
	http.HandleFunc("/mysql/insert", app.mysqlInsert)
	http.HandleFunc("/mysql/update", app.mysqlUpdate)
	http.HandleFunc("/mysql/delete", app.mysqlDelete)
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
