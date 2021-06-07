package main

import (
	"context"
	"dbforum/application"
	"dbforum/infrastructure/persistance"
	"dbforum/interfaces"
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

	userRepo := persistance.NewUserRepo(postgresConn)
	forumRepo := persistance.NewForumRepo(postgresConn)
	threadRepo := persistance.NewThreadRepo(postgresConn)

	userApp := application.NewUserApp(userRepo)
	forumApp := application.NewForumApp(forumRepo)
	threadApp := application.NewThreadApp(threadRepo)

	userInfo := interfaces.NewUserInfo(userApp)
	forumInfo := interfaces.NewForumInfo(forumApp)
	threadInfo := interfaces.NewThreadInfo(threadApp)

	router := router.New()

	router.POST("/user/{username}/create", userInfo.CreateUser)
	router.GET("/user/{username}/profile", userInfo.GetUser)
	router.POST("/user/{username}/profile", userInfo.EditUser)

	router.POST("/forum/create", forumInfo.CreateForum)
	router.GET("/forum/{forumname}/details", forumInfo.GetForum)
	router.GET("/forum/{forumname}/users", forumInfo.GetForumUsers)
	router.GET("/forum/{forumname}/threads", forumInfo.GetForumThreads)

	router.POST("/forum/{forumname}/create", threadInfo.CreateThread)
	router.GET("/thread/{threadnameOrID}/details", threadInfo.GetThread)
	router.POST("/thread/{threadnameOrID}/details", threadInfo.EditThread)
	router.GET("/thread/{threadnameOrID}/posts", threadInfo.GetThreadPosts)
	router.POST("/thread/{threadnameOrID}/vote", threadInfo.VoteThread)

	fmt.Printf("Starting server at localhost%s\n", addr)
	fasthttp.ListenAndServe(addr, router.Handler)
}

func main() {
	runServer(":5050")
}
