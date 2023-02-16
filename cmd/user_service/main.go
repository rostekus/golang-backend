package main

import (
	"log"
	"rostekus/golang-backend/db"
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/server"
	"rostekus/golang-backend/internal/user"
	"rostekus/golang-backend/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection")
	}
	defer dbConn.Close()

	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userHadler := user.NewHandler(userService)
	healthHandler := health.NewHandler()
	router := router.NewUserServiceRouter(userHadler, healthHandler)
	srv := server.NewServer("23450", router)
	srv.Run()

}
