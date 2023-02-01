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

type App struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := App{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	dbConn, err := db.NewDatabase()
	if err != nil {
		app.errorLog.Fatalf("Could not initialize database connection")
	}
	defer dbConn.Close()

	err = dbConn.MigrateDB()
	if err != nil {
		app.errorLog.Fatalf("Could not migrate schema: %s", err)
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
			app.errorLog.Fatalf("Listen and Serve: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit
	app.infoLog.Printf("Listen and Serve: %s\n", err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		app.infoLog.Printf("Server Shutdown: %s\n", err)
	}
	app.infoLog.Printf("Server exiting")

}
