package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	Server   *http.Server
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewServer(port string, router http.Handler) *server {
	port = fmt.Sprintf(":%s", port)
	return &server{
		Server: &http.Server{
			Addr:    port,
			Handler: router,
		},
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)}
}
func (s *server) Run() {
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.ErrorLog.Fatalf("Listen and Serve: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit
	s.InfoLog.Printf("Listen and Serve")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		s.InfoLog.Printf("Server Shutdown: %s\n", err)
	}
	s.InfoLog.Printf("Server exiting")

}
