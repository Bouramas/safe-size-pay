package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := NewLogger()

	cfg, err := newConfig()
	if err != nil {
		logger.Fatalln("could not setup config, ", err)
	}

	// start Server
	var wg sync.WaitGroup
	s, err := NewServer(cfg, logger)
	if err != nil {
		logger.Fatalln("Failed to start server, ", err)
	}

	// Handle graceful HTTP server shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		gracefulTerm := make(chan os.Signal, 1)
		signal.Notify(gracefulTerm, syscall.SIGINT, syscall.SIGTERM)
		sig := <-gracefulTerm
		s.Log.Infof("Server notified, %+v", sig)

		if err := s.Shutdown(context.Background()); err != nil {
			s.Log.Errorf("Failed to properly shutdown the server, %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.ListenAndServe()
		if err != nil {
			s.Log.Errorf("%v", err)
		}
	}()

	s.Log.Infof("Listening on %s...", s.srv.Addr)
	wg.Wait()
	_ = cfg.DB.Close()
	s.Log.Infof("Process gracefully terminated")
}

func NewLogger() *logrus.Logger {

	logger := logrus.New()
	logger.SetFormatter(
		&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		})
	return logger
}
