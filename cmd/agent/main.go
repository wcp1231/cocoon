package main

import (
	"cocoon/pkg/agent"
	log "cocoon/pkg/logger"
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DUMP_CMD = "dump"
	MOCK_CMD = "mock"
)

var (
	dumpCmd *flag.FlagSet
	mockCmd *flag.FlagSet

	appname string
	session string
	listen  string
	remote  string

	logger = log.NewLogger()
)

func init() {
	dumpCmd = flag.NewFlagSet(DUMP_CMD, flag.ExitOnError)
	dumpCmd.StringVar(&appname, "app", "Application", "Application name")
	dumpCmd.StringVar(&appname, "session", "", "Application session")
	dumpCmd.StringVar(&listen, "listen", "0.0.0.0:1820", "Listen address")
	dumpCmd.StringVar(&remote, "remote", "", "Remote agent address")
	mockCmd = flag.NewFlagSet(MOCK_CMD, flag.ExitOnError)
	mockCmd.StringVar(&appname, "app", "", "Application name")
	mockCmd.StringVar(&listen, "listen", "0.0.0.0:1820", "Listen address")
	mockCmd.StringVar(&remote, "remote", "", "Remote agent address")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'dump' or 'mock' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case DUMP_CMD:
		err := dumpCmd.Parse(os.Args[2:])
		if err != nil {
			dumpCmd.PrintDefaults()
			os.Exit(1)
		}
		listenAddr, err := net.ResolveTCPAddr("tcp", listen)
		if err != nil {
			logger.Fatal("error", zap.Error(err))
			os.Exit(1)
		}

		ensureSession()

		logger.Info("Start proxy mode.",
			zap.String("app", appname),
			zap.String("session", session),
			zap.String("listen", listen),
			zap.String("remote", remote))
		s := agent.NewServer(context.Background(), logger, appname, session, remote, listenAddr)
		go s.Start()

		signalChan := make(chan os.Signal, 1)
		signal.Ignore()
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sc := <-signalChan

		switch sc {
		case syscall.SIGINT:
			logger.Info("Shutting down proxy...")
			s.Shutdown()
			s.Wg.Wait()
			<-s.ClosedChan
		case syscall.SIGQUIT, syscall.SIGTERM:
			logger.Info("Graceful Shutting down proxy...")
			s.GracefulShutdown()
			s.Wg.Wait()
			<-s.ClosedChan
		default:
			logger.Info("Unexpected signal")
			os.Exit(1)
		}
	case MOCK_CMD:
		err := mockCmd.Parse(os.Args[1:])
		if err != nil {
			mockCmd.PrintDefaults()
			os.Exit(1)
		}
		fmt.Println("Start mock mode")
	default:
		fmt.Printf("Expected 'dump' or 'mock' subcommand")
		os.Exit(1)
	}
}

func ensureSession() {
	if session == "" {
		session = time.Now().Format("2006-01-02_15:04:05")
	}
}
