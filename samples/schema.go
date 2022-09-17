package main

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

type MysqlSample struct {
	ID            int         `json:"id"`
	TinyType      int8        `json:"tinyType"`
	IntType       int32       `json:"intType"`
	BigIntType    int64       `json:"bigIntType"`
	FloatType     float32     `json:"floatType"`
	DoubleType    float64     `json:"doubleType"`
	DecimalType   float64     `json:"decimalType"`
	DateType      string      `json:"dateType"`
	TimeType      string      `json:"timeType"`
	YearType      string      `json:"yearType"`
	DatetimeType  string      `json:"datetimeType"`
	TimestampType time.Time   `json:"timestampType"`
	CharType      byte        `json:"charType"`
	VarCharType   string      `json:"varCharType"`
	TinyBlobType  []byte      `json:"tiny_blob_type"`
	TinyTextType  string      `json:"tinyTextType"`
	BlobType      []byte      `json:"blobType"`
	TextType      string      `json:"textType"`
	NullType      interface{} `json:"nullType"`
}

func NewRandomMysqlSample() *MysqlSample {
	date := RandomDate()
	return &MysqlSample{
		TinyType:      int8(RandomInt64Range(0, 256)),
		IntType:       int32(RandomInt64Range(0, math.MaxInt32)),
		BigIntType:    RandomInt64Range(0, math.MaxInt64),
		FloatType:     float32(RandomFloat64Range(0, math.MaxFloat32)),
		DoubleType:    RandomFloat64Range(0, math.MaxFloat64),
		DecimalType:   RandomFloat64Range(0, 10000),
		DateType:      date.Format("2006-01-02"),
		TimeType:      date.Format("15:04:05"),
		YearType:      date.Format("2006"),
		DatetimeType:  date.Format("2006-01-02T15:04:05"),
		TimestampType: date,
		CharType:      byte(RandomInt64Range(32, 127)),
		VarCharType:   RandomAlphaString(int(RandomInt64Range(1, 64))),
		TinyBlobType:  []byte(RandomAlphaString(int(RandomInt64Range(1, 64)))),
		TinyTextType:  RandomAlphaString(int(RandomInt64Range(1, 64))),
		BlobType:      []byte(RandomAlphaString(int(RandomInt64Range(64, 128)))),
		TextType:      RandomAlphaString(int(RandomInt64Range(64, 128))),
	}
}

func NewMysqlSampleFromScan(row *sql.Rows) (*MysqlSample, error) {
	s := &MysqlSample{}
	var char string
	err := row.Scan(&s.ID, &s.TinyType, &s.IntType, &s.BigIntType, &s.FloatType, &s.DoubleType, &s.DecimalType,
		&s.DateType, &s.TimeType, &s.YearType, &s.DatetimeType, &s.TimestampType,
		&char, &s.VarCharType, &s.TinyBlobType, &s.TinyTextType, &s.BlobType, &s.TextType, &s.NullType)
	s.CharType = []byte(char)[0]
	return s, err
}

func (m *MysqlSample) InsertSQL() string {
	insert := `INSERT INTO samples(
tiny_type,int_type,big_int_type,float_type,double_type,decimal_type,
date_type,time_type,year_type,datetime_type,timestamp_type,
char_type,varchar_type,tinyblob_type,tinytext_type,blob_type,text_type)
VALUES (%d,%d,%d,%f,%f,%8.4f,
'%s','%s','%s','%s',NOW(),
'%c','%s','%s','%s','%s','%s')`
	return fmt.Sprintf(insert, m.TinyType, m.IntType, m.BigIntType, m.FloatType, m.DoubleType, m.DecimalType,
		m.DateType, m.TimeType, m.YearType, m.DatetimeType,
		m.CharType, m.VarCharType, m.TinyBlobType, m.TinyTextType, m.BlobType, m.TextType)
}
