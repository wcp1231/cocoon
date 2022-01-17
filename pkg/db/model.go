package db

import (
	"cocoon/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	COLLECTION_APPLICATION         = "applications"
	COLLECTION_APPLICATION_SESSION = "app_sessions"
	COLLECTION_TCP_TRAFFIC         = "tcp_traffics"
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
	Direction   *model.Direction   `json:"direction,omitempty" bson:"direction,omitempty"`
	Seq         uint64             `json:"seq,omitempty" bson:"seq,omitempty"`
	Timestamp   time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Raw         []byte             `json:"raw,omitempty" bson:"raw,omitempty"`
}
