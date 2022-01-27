package main

import (
	log "cocoon/pkg/logger"
	"cocoon/pkg/server"
	"flag"
	"go.uber.org/zap"
)

var (
	dbUri string
)

func main() {
	flag.StringVar(&dbUri, "db", "", "Mongodb Uri")
	flag.Parse()

	logger := log.NewLogger()
	s := server.NewCocoonServer(logger, dbUri)

	go func() {
		s.StartHttp(":8080")
	}()

	err := s.Start(":7070")
	if err != nil {
		logger.Error("Start cocoon server failed", zap.Error(err))
	}
}
