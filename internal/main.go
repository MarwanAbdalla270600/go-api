package main

import (
	"context"
	"go-api/internal/controller"
	"go-api/internal/middleware"
	"go-api/internal/repo"
	"go-api/internal/service"
	"go-api/internal/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func main() {
	router := gin.Default()

	//dotenv
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	//database
	dsn := utils.CeateConnectionString()
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	//redis store
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	//initialize singletons for user
	userRepo := repo.NewUserRepo(db, rdb, ctx)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router.GET("/users", middleware.AuthMiddleware(userRepo), userController.GetAll)
	router.POST("/auth/register", userController.RegisterUser)
	router.POST("/auth/login", userController.Login)
	router.POST("/auth/logout", userController.Logout)
	router.Run(":8080")
}
