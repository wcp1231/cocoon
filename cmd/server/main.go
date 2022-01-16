package main

import (
	log "cocoon/pkg/logger"
	"cocoon/pkg/server"
	"go.uber.org/zap"
)

func main() {
	logger := log.NewLogger()
	s := server.NewRpcServer(logger)
	err := s.Start(":7070")
	if err != nil {
		logger.Error("Start cocoon server failed", zap.Error(err))
	}
}
