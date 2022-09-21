package proto

// Flags sent from the MySQL C API
type MySqlFlag int32

const (
	MySqlFlag_EMPTY                 MySqlFlag = 0
	MySqlFlag_NOT_NULL_FLAG         MySqlFlag = 1
	MySqlFlag_PRI_KEY_FLAG          MySqlFlag = 2
	MySqlFlag_UNIQUE_KEY_FLAG       MySqlFlag = 4
	MySqlFlag_MULTIPLE_KEY_FLAG     MySqlFlag = 8
	MySqlFlag_BLOB_FLAG             MySqlFlag = 16
	MySqlFlag_UNSIGNED_FLAG         MySqlFlag = 32
	MySqlFlag_ZEROFILL_FLAG         MySqlFlag = 64
	MySqlFlag_BINARY_FLAG           MySqlFlag = 128
	MySqlFlag_ENUM_FLAG             MySqlFlag = 256
	MySqlFlag_AUTO_INCREMENT_FLAG   MySqlFlag = 512
	MySqlFlag_TIMESTAMP_FLAG        MySqlFlag = 1024
	MySqlFlag_SET_FLAG              MySqlFlag = 2048
	MySqlFlag_NO_DEFAULT_VALUE_FLAG MySqlFlag = 4096
	MySqlFlag_ON_UPDATE_NOW_FLAG    MySqlFlag = 8192
	MySqlFlag_NUM_FLAG              MySqlFlag = 32768
	MySqlFlag_PART_KEY_FLAG         MySqlFlag = 16384
	MySqlFlag_GROUP_FLAG            MySqlFlag = 32768
	MySqlFlag_UNIQUE_FLAG           MySqlFlag = 65536
	MySqlFlag_BINCMP_FLAG           MySqlFlag = 131072
)

var MySqlFlag_name = map[int32]string{
	0:     "EMPTY",
	1:     "NOT_NULL_FLAG",
	2:     "PRI_KEY_FLAG",
	4:     "UNIQUE_KEY_FLAG",
	8:     "MULTIPLE_KEY_FLAG",
	16:    "BLOB_FLAG",
	32:    "UNSIGNED_FLAG",
	64:    "ZEROFILL_FLAG",
	128:   "BINARY_FLAG",
	256:   "ENUM_FLAG",
	512:   "AUTO_INCREMENT_FLAG",
	1024:  "TIMESTAMP_FLAG",
	2048:  "SET_FLAG",
	4096:  "NO_DEFAULT_VALUE_FLAG",
	8192:  "ON_UPDATE_NOW_FLAG",
	32768: "NUM_FLAG",
	16384: "PART_KEY_FLAG",
	// Duplicate value: 32768: "GROUP_FLAG",
	65536:  "UNIQUE_FLAG",
	131072: "BINCMP_FLAG",
}
var MySqlFlag_value = map[string]int32{
	"EMPTY":                 0,
	"NOT_NULL_FLAG":         1,
	"PRI_KEY_FLAG":          2,
	"UNIQUE_KEY_FLAG":       4,
	"MULTIPLE_KEY_FLAG":     8,
	"BLOB_FLAG":             16,
	"UNSIGNED_FLAG":         32,
	"ZEROFILL_FLAG":         64,
	"BINARY_FLAG":           128,
	"ENUM_FLAG":             256,
	"AUTO_INCREMENT_FLAG":   512,
	"TIMESTAMP_FLAG":        1024,
	"SET_FLAG":              2048,
	"NO_DEFAULT_VALUE_FLAG": 4096,
	"ON_UPDATE_NOW_FLAG":    8192,
	"NUM_FLAG":              32768,
	"PART_KEY_FLAG":         16384,
	"GROUP_FLAG":            32768,
	"UNIQUE_FLAG":           65536,
	"BINCMP_FLAG":           131072,
}

type Flag int32

const (
	Flag_NONE       Flag = 0
	Flag_ISINTEGRAL Flag = 256
	Flag_ISUNSIGNED Flag = 512
	Flag_ISFLOAT    Flag = 1024
	Flag_ISQUOTED   Flag = 2048
	Flag_ISTEXT     Flag = 4096
	Flag_ISBINARY   Flag = 8192
)

var Flag_name = map[int32]string{
	0:    "NONE",
	256:  "ISINTEGRAL",
	512:  "ISUNSIGNED",
	1024: "ISFLOAT",
	2048: "ISQUOTED",
	4096: "ISTEXT",
	8192: "ISBINARY",
}
var Flag_value = map[string]int32{
	"NONE":       0,
	"ISINTEGRAL": 256,
	"ISUNSIGNED": 512,
	"ISFLOAT":    1024,
	"ISQUOTED":   2048,
	"ISTEXT":     4096,
	"ISBINARY":   8192,
}

// Type defines the various supported data types in bind vars
// and query results.
type Type int32

const (
	// NULL_TYPE specifies a NULL type.
	Type_NULL_TYPE Type = 0
	// INT8 specifies a TINYINT type.
	// Properties: 1, IsNumber.
	Type_INT8 Type = 257
	// UINT8 specifies a TINYINT UNSIGNED type.
	// Properties: 2, IsNumber, IsUnsigned.
	Type_UINT8 Type = 770
	// INT16 specifies a SMALLINT type.
	// Properties: 3, IsNumber.
	Type_INT16 Type = 259
	// UINT16 specifies a SMALLINT UNSIGNED type.
	// Properties: 4, IsNumber, IsUnsigned.
	Type_UINT16 Type = 772
	// INT24 specifies a MEDIUMINT type.
	// Properties: 5, IsNumber.
	Type_INT24 Type = 261
	// UINT24 specifies a MEDIUMINT UNSIGNED type.
	// Properties: 6, IsNumber, IsUnsigned.
	Type_UINT24 Type = 774
	// INT32 specifies a INTEGER type.
	// Properties: 7, IsNumber.
	Type_INT32 Type = 263
	// UINT32 specifies a INTEGER UNSIGNED type.
	// Properties: 8, IsNumber, IsUnsigned.
	Type_UINT32 Type = 776
	// INT64 specifies a BIGINT type.
	// Properties: 9, IsNumber.
	Type_INT64 Type = 265
	// UINT64 specifies a BIGINT UNSIGNED type.
	// Properties: 10, IsNumber, IsUnsigned.
	Type_UINT64 Type = 778
	// FLOAT32 specifies a FLOAT type.
	// Properties: 11, IsFloat.
	Type_FLOAT32 Type = 1035
	// FLOAT64 specifies a DOUBLE or REAL type.
	// Properties: 12, IsFloat.
	Type_FLOAT64 Type = 1036
	// TIMESTAMP specifies a TIMESTAMP type.
	// Properties: 13, IsQuoted.
	Type_TIMESTAMP Type = 2061
	// DATE specifies a DATE type.
	// Properties: 14, IsQuoted.
	Type_DATE Type = 2062
	// TIME specifies a TIME type.
	// Properties: 15, IsQuoted.
	Type_TIME Type = 2063
	// DATETIME specifies a DATETIME type.
	// Properties: 16, IsQuoted.
	Type_DATETIME Type = 2064
	// YEAR specifies a YEAR type.
	// Properties: 17, IsNumber, IsUnsigned.
	Type_YEAR Type = 785
	// DECIMAL specifies a DECIMAL or NUMERIC type.
	// Properties: 18, None.
	Type_DECIMAL Type = 18
	// TEXT specifies a TEXT type.
	// Properties: 19, IsQuoted, IsText.
	Type_TEXT Type = 6163
	// BLOB specifies a BLOB type.
	// Properties: 20, IsQuoted, IsBinary.
	Type_BLOB Type = 10260
	// VARCHAR specifies a VARCHAR type.
	// Properties: 21, IsQuoted, IsText.
	Type_VARCHAR Type = 6165
	// VARBINARY specifies a VARBINARY type.
	// Properties: 22, IsQuoted, IsBinary.
	Type_VARBINARY Type = 10262
	// CHAR specifies a CHAR type.
	// Properties: 23, IsQuoted, IsText.
	Type_CHAR Type = 6167
	// BINARY specifies a BINARY type.
	// Properties: 24, IsQuoted, IsBinary.
	Type_BINARY Type = 10264
	// BIT specifies a BIT type.
	// Properties: 25, IsQuoted.
	Type_BIT Type = 2073
	// ENUM specifies an ENUM type.
	// Properties: 26, IsQuoted.
	Type_ENUM Type = 2074
	// SET specifies a SET type.
	// Properties: 27, IsQuoted.
	Type_SET Type = 2075
	// TUPLE specifies a tuple. This cannot
	// be returned in a QueryResult, but it can
	// be sent as a bind var.
	// Properties: 28, None.
	Type_TUPLE Type = 28
	// GEOMETRY specifies a GEOMETRY type.
	// Properties: 29, IsQuoted.
	Type_GEOMETRY Type = 2077
	// JSON specifies a JSON type.
	// Properties: 30, IsQuoted.
	Type_JSON Type = 2078
	// EXPRESSION specifies a SQL expression.
	// This type is for internal use only.
	// Properties: 31, None.
	Type_EXPRESSION Type = 31
)

var Type_name = map[int32]string{
	0:     "NULL_TYPE",
	257:   "INT8",
	770:   "UINT8",
	259:   "INT16",
	772:   "UINT16",
	261:   "INT24",
	774:   "UINT24",
	263:   "INT32",
	776:   "UINT32",
	265:   "INT64",
	778:   "UINT64",
	1035:  "FLOAT32",
	1036:  "FLOAT64",
	2061:  "TIMESTAMP",
	2062:  "DATE",
	2063:  "TIME",
	2064:  "DATETIME",
	785:   "YEAR",
	18:    "DECIMAL",
	6163:  "TEXT",
	10260: "BLOB",
	6165:  "VARCHAR",
	10262: "VARBINARY",
	6167:  "CHAR",
	10264: "BINARY",
	2073:  "BIT",
	2074:  "ENUM",
	2075:  "SET",
	28:    "TUPLE",
	2077:  "GEOMETRY",
	2078:  "JSON",
	31:    "EXPRESSION",
}
var Type_value = map[string]int32{
	"NULL_TYPE":  0,
	"INT8":       257,
	"UINT8":      770,
	"INT16":      259,
	"UINT16":     772,
	"INT24":      261,
	"UINT24":     774,
	"INT32":      263,
	"UINT32":     776,
	"INT64":      265,
	"UINT64":     778,
	"FLOAT32":    1035,
	"FLOAT64":    1036,
	"TIMESTAMP":  2061,
	"DATE":       2062,
	"TIME":       2063,
	"DATETIME":   2064,
	"YEAR":       785,
	"DECIMAL":    18,
	"TEXT":       6163,
	"BLOB":       10260,
	"VARCHAR":    6165,
	"VARBINARY":  10262,
	"CHAR":       6167,
	"BINARY":     10264,
	"BIT":        2073,
	"ENUM":       2074,
	"SET":        2075,
	"TUPLE":      28,
	"GEOMETRY":   2077,
	"JSON":       2078,
	"EXPRESSION": 31,
}

func (x Type) String() string {
	return Type_name[int32(x)]
}

// EventToken is a structure that describes a point in time in a
// replication stream on one shard. The most recent known replication
// position can be retrieved from vttablet when executing a query. It
// is also sent with the replication streams from the binlog service.
type EventToken struct {
	// timestamp is the MySQL timestamp of the statements. Seconds since Epoch.
	Timestamp int64 `json:"timestamp,omitempty"`
	// The shard name that applied the statements. Note this is not set when
	// streaming from a vttablet. It is only used on the client -> vtgate link.
	Shard string `json:"shard,omitempty"`
	// The position on the replication stream after this statement was applied.
	// It is not the transaction ID / GTID, but the position / GTIDSet.
	Position string `json:"position,omitempty"`
}

// Value represents a typed value.
type Value struct {
	Type  Type   `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

// BindVariable represents a single bind variable in a Query.
type BindVariable struct {
	Type  Type   `json:"type,omitempty"`
	Value []byte `json:"value,omitempty"`
	// values are set if type is TUPLE.
	Values []*Value `json:"values,omitempty"`
}

// Field describes a single column returned by a query
type Field struct {
	// name of the field as returned by mysql C API
	Name string `json:"name,omitempty"`
	// vitess-defined type. Conversion function is in sqltypes package.
	Type Type `json:"type,omitempty"`
	// Remaining Fields from mysql C API.
	// These Fields are only populated when ExecuteOptions.included_fields
	// is set to IncludedFields.ALL.
	Table    string `json:"table,omitempty"`
	OrgTable string `json:"org_table,omitempty"`
	Database string `json:"database,omitempty"`
	OrgName  string `json:"org_name,omitempty"`
	// column_length is really a uint32. All 32 bits can be used.
	ColumnLength uint32 `json:"column_length,omitempty"`
	// charset is actually a uint16. Only the lower 16 bits are used.
	Charset uint32 `json:"charset,omitempty"`
	// decimals is actualy a uint8. Only the lower 8 bits are used.
	Decimals uint32 `json:"decimals,omitempty"`
	// flags is actually a uint16. Only the lower 16 bits are used.
	Flags uint32 `json:"flags,omitempty"`
}

// Row is a database row.
type Row struct {
	// lengths contains the length of each value in values.
	// A length of -1 means that the field is NULL. While
	// reading values, you have to accummulate the length
	// to know the offset where the next value begins in values.
	Lengths []int64 `json:"lengths,omitempty"`
	// values contains a concatenation of all values in the row.
	Values []byte `json:"values,omitempty"`
}

// ResultExtras contains optional out-of-band information. Usually the
// extras are requested by adding ExecuteOptions flags.
type ResultExtras struct {
	// event_token is populated if the include_event_token flag is set
	// in ExecuteOptions.
	EventToken *EventToken `json:"event_token,omitempty"`
	// If set, it means the data returned with this result is fresher
	// than the compare_token passed in the ExecuteOptions.
	Fresher bool `json:"fresher,omitempty"`
}
