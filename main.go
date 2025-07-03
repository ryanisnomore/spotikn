package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/chromedp/chromedp"
)

var addr = os.Getenv("SPOTIFY_TOKENER_ADDR")

func main() {
	if addr == "" {
		addr = "0.0.0.0:8080"
	}

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	mux := http.NewServeMux()

	s := &server{
		ctx: ctx,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}

	mux.HandleFunc("/api/token", s.handleToken)

	go s.Start()
	defer s.Stop()

	slog.Info("Server started", slog.String("address", s.server.Addr))
	defer slog.Info("Server stopped", slog.String("address", s.server.Addr))
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

type server struct {
	ctx    context.Context
	server *http.Server
}

func (s *server) Start() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start server", slog.Any("err", err))
	}
}

func (s *server) Stop() {
	if err := s.server.Close(); err != nil {
		slog.Error("Failed to stop server", slog.Any("err", err))
	}
}
