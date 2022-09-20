package proto

import (
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/dubbo"
	"cocoon/pkg/proto/http"
	"cocoon/pkg/proto/mongo"
	"cocoon/pkg/proto/mysql"
	"cocoon/pkg/proto/redis"
)

func NewRequestDissector(protocol *common.Protocol, reqC, respC chan common.Message) api.ClientDissector {
	switch protocol {
	case common.PROTOCOL_HTTP:
		return http.NewRequestDissector(reqC, respC)
	case common.PROTOCOL_REDIS:
		return redis.NewRequestDissector(reqC, respC)
	case common.PROTOCOL_MONGO:
		return mongo.NewRequestDissector(reqC, respC)
	case common.PROTOCOL_DUBBO:
		return dubbo.NewRequestDissector(reqC, respC)
	}
	return nil
}

func NewMysqlDissector(reqC, respC chan common.Message) *mysql.Dissector {
	return mysql.NewRequestDissector(reqC, respC)
}
