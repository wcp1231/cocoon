package db

import (
	"cocoon/pkg/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	COLLECTION_APPLICATION         = "applications"
	COLLECTION_APPLICATION_SESSION = "app_sessions"
	COLLECTION_TCP_TRAFFIC         = "tcp_traffics"
	COLLECTION_RECORDS             = "records"
)

// Application 应用信息
type Application struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

// ApplicationSession
type ApplicationSession struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AppID   string             `json:"app,omitempty" bson:"app,omitempty"`
	Session string             `json:"session,omitempty" bson:"session,omitempty"`
}

// TcpTraffic 原始的 TCP 数据
type TcpTraffic struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Session     string             `json:"session,omitempty" bson:"session,omitempty"`
	Source      string             `json:"src,omitempty" bson:"src,omitempty"`
	Destination string             `json:"dst,omitempty" bson:"dst,omitempty"`
	IsOutgoing  bool               `json:"outgoing,omitempty" bson:"outgoing,omitempty"`
	Direction   *common.Direction  `json:"direction,omitempty" bson:"direction,omitempty"`
	Seq         uint64             `json:"seq,omitempty" bson:"seq,omitempty"`
	Timestamp   time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Size        int                `json:"size,omitempty" bson:"size,omitempty"`
	Raw         []byte             `json:"raw,omitempty" bson:"raw,omitempty"`
}

// Record 记录 request 和 response
// 暂时不考虑同样的 request 发给不同 server 返回结果不一致的情况
// 比如 get key 发给不同的 redis 集群返回不同的业务数据
type Record struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Session    string             `json:"session,omitempty" bson:"session,omitempty"`
	IsOutgoing bool               `json:"outgoing,omitempty" bson:"outgoing,omitempty"`
	Proto      string             `json:"proto,omitempty" bson:"proto,omitempty"`
	ReqHeader  map[string]string  `json:"req_header,omitempty" bson:"req_header,omitempty"`
	RespHeader map[string]string  `json:"resp_header,omitempty" bson:"resp_header,omitempty"`
	ReqBody    string             `json:"req_body,omitempty" bson:"req_body,omitempty"`
	RespBody   string             `json:"resp_body,omitempty" bson:"resp_body,omitempty"`
}
