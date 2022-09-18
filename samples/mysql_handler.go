package main

import (
	"errors"
	"net/http"
)

func (s *SampleApp) mysqlSelect(w http.ResponseWriter, req *http.Request) {
	res, err := s.db.Query("SELECT * FROM samples")
	if err != nil {
		responseErr(w, err)
		return
	}
	defer res.Close()
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

func (s *SampleApp) mysqlPreparedSelect(w http.ResponseWriter, req *http.Request) {
	res, err := s.db.Query("SELECT * FROM samples WHERE id > ?", 0)
	if err != nil {
		responseErr(w, err)
		return
	}
	defer res.Close()
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
