package main

import (
	"context"
	"dbforum/application"
	"dbforum/infrastructure/persistance"
	"dbforum/interfaces/forum"
	"dbforum/interfaces/user"
	"fmt"
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func runServer(addr string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load .env file", zap.String("error", err.Error()))
	}

	postgresConnectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	postgresConn, err := pgxpool.Connect(context.Background(), postgresConnectionString)
	if err != nil {
		log.Fatal("Could not connect to postgres database", zap.String("error", err.Error()))
		return
	}

	repoUser := persistance.NewUserRepo(postgresConn)
	repoForum := persistance.NewForumRepo(postgresConn)

	userApp := application.NewUserApp(repoUser)
	forumApp := application.NewForumApp(repoForum)

	userInfo := user.NewUserInfo(userApp)
	forumInfo := forum.NewForumInfo(forumApp)

	router := router.New()

	router.POST("/user/{username}/create", userInfo.CreateUser)
	router.GET("/user/{username}/profile", userInfo.GetUser)
	router.POST("/user/{username}/profile", userInfo.EditUser)

	router.POST("/forum/create", forumInfo.CreateForum)
	router.GET("/forum/{forumname}/details", forumInfo.GetForum)

	fmt.Printf("Starting server at localhost%s\n", addr)
	fasthttp.ListenAndServe(addr, router.Handler)
}

func main() {
	runServer(":5050")
}
