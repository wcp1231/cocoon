package server

import (
	"cocoon/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type CocoonHttpHandler struct {
	logger *zap.Logger
	server *CocoonServer

	srv *http.Server
}

func NewCocoonHttpHandler(logger *zap.Logger, server *CocoonServer) *CocoonHttpHandler {
	httpServer := &CocoonHttpHandler{
		logger: logger,
		server: server,
		srv:    &http.Server{},
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/records/{session}", httpServer.getRecordsBySession).Methods("GET")

	httpServer.srv.Handler = r
	return httpServer
}

func (c *CocoonHttpHandler) getRecordsBySession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	session, ok := vars["session"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"missing session\"}")
		return
	}
	cursor, err := c.server.database.ReadRecordsBySession(session)
	if err != nil {
		c.logger.Error("Get records by session failed", zap.String("session", session), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"some error\"}")
		return
	}
	ctx := context.Background()
	var records []db.Record
	err = cursor.All(ctx, &records)
	if err != nil {
		c.logger.Error("Get all records failed", zap.String("session", session), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"some error\"}")
		return
	}
	resp, err := json.Marshal(records)
	if err != nil {
		c.logger.Error("Marshal records failed", zap.String("session", session), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"some error\"}")
		return
	}
	w.Write(resp)
}
