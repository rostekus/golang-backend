package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rostekus/golang-backend/db"
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/user"
	"rostekus/golang-backend/router"
	"syscall"
	"time"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}
	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userHadler := user.NewHandler(userService)
	healthHandler := health.NewHandler()
	router := router.NewRouter(userHadler, healthHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
