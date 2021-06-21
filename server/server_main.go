package main

import (
	"context"
	"dbforum/application"
	"dbforum/infrastructure/persistance"
	"dbforum/interfaces"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func loggerMid(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		if end.Sub(begin) > 90*time.Millisecond {
			log.Printf("%s - %s",
				string(ctx.Request.URI().FullURI()),
				end.Sub(begin).String())
		}
	})
}

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
	postRepo := persistance.NewPostRepo(postgresConn)
	threadRepo := persistance.NewThreadRepo(postgresConn)
	serviceRepo := persistance.NewServiceRepo(postgresConn)

	serviceApp := application.NewServiceApp(serviceRepo)
	userApp := application.NewUserApp(userRepo, serviceApp)
	forumApp := application.NewForumApp(forumRepo, serviceApp)
	postApp := application.NewPostApp(postRepo, userRepo, threadRepo, forumRepo, serviceApp)
	threadApp := application.NewThreadApp(threadRepo, postRepo, serviceApp)

	userInfo := interfaces.NewUserInfo(userApp)
	forumInfo := interfaces.NewForumInfo(forumApp)
	postInfo := interfaces.NewPostInfo(postApp, threadApp)
	threadInfo := interfaces.NewThreadInfo(threadApp)
	serviceInfo := interfaces.NewServiceInfo(serviceApp)

	router := router.New()

	prefix := "/api"
	router.POST(prefix+"/user/{username}/create", userInfo.CreateUser)
	router.GET(prefix+"/user/{username}/profile", userInfo.GetUser)
	router.POST(prefix+"/user/{username}/profile", userInfo.EditUser)

	router.POST(prefix+"/forum/create", forumInfo.CreateForum)
	router.GET(prefix+"/forum/{forumname}/details", forumInfo.GetForum)
	router.GET(prefix+"/forum/{forumname}/users", forumInfo.GetForumUsers)
	router.GET(prefix+"/forum/{forumname}/threads", forumInfo.GetForumThreads)

	router.POST(prefix+"/forum/{forumname}/create", threadInfo.CreateThread)
	router.GET(prefix+"/thread/{threadnameOrID}/details", threadInfo.GetThread)
	router.POST(prefix+"/thread/{threadnameOrID}/details", threadInfo.EditThread)
	router.GET(prefix+"/thread/{threadnameOrID}/posts", threadInfo.GetThreadPosts)
	router.POST(prefix+"/thread/{threadnameOrID}/vote", threadInfo.VoteThread)

	router.POST(prefix+"/thread/{threadnameOrID}/create", postInfo.CreatePost)
	router.GET(prefix+"/post/{postID}/details", postInfo.GetPost)
	router.POST(prefix+"/post/{postID}/details", postInfo.EditPost)

	router.GET(prefix+"/service/status", serviceInfo.GetForumStats)
	router.POST(prefix+"/service/clear", serviceInfo.ClearForum)

	fmt.Printf("Starting server at localhost%s\n", addr)
	fasthttp.ListenAndServe(addr, loggerMid(router.Handler))
}

func main() {
	runServer(":5000")
}
