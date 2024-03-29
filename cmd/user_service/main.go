package main

import (
	"log"
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/router"
	"rostekus/golang-backend/internal/server"
	"rostekus/golang-backend/internal/user"
	"rostekus/golang-backend/pkg/db"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	dbConn, err := db.NewPostgresDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection")
	}
	defer dbConn.Close()

	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userHadler := user.NewHandler(userService)
	healthHandler := health.NewHandler()
	router := router.NewUserServiceRouter(userHadler, healthHandler)
	srv := server.NewServer("23450", router.Router)
	srv.Run()

}
