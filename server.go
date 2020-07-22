package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/db"
	dbmodel "github.com/ryakosh/wishlist/lib/db/model"
	"github.com/ryakosh/wishlist/lib/graph"
	"github.com/ryakosh/wishlist/lib/graph/generated"
)

const (
	defaultPort              = "8080"
	logsDir                  = "./logs/"
	defaultRequestComplexity = 10
)

var accessLog *log.Logger

func graphqlHandler() gin.HandlerFunc {
	config := generated.Config{Resolvers: &graph.Resolver{DB: db.DB}}
	config.Directives.AuthRequired = dbmodel.AuthRequired
	config.Directives.EmailVerificationRequired = dbmodel.EmailVerificationRequired
	calcComplexity(&config.Complexity)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	h.Use(extension.FixedComplexityLimit(100))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func calcComplexity(complexityRoot *generated.ComplexityRoot) {
	calcUsersComplexity := func(childComplexity int, _ int, limit int) int {
		return (childComplexity * limit) + defaultRequestComplexity
	}

	complexityRoot.Users.Query = calcUsersComplexity
	complexityRoot.Wishes.Query = calcUsersComplexity
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func accessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteAddr := c.Request.RemoteAddr
		method := c.Request.Method
		path := c.Request.URL.Path

		accessLog.Printf("%s - %s - %s\n", remoteAddr, method, path)
	}
}

func initLogs() {
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(logsDir, 0700); err != nil {
			lib.LogError(lib.LFatal, "Could not create logs directory", err)
		}
	} else if err != nil {
		lib.LogError(lib.LFatal, "Could not create logs directory", err)
	}

	serverLog, err := os.OpenFile(logsDir+"server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)

	}
	log.SetOutput(serverLog)

	gin.DisableConsoleColor()
	ginLog, err := os.OpenFile(logsDir+"gin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)
	}
	gin.DefaultWriter = ginLog

	accessLogFile, err := os.OpenFile(logsDir+"access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)

	}

	accessLog = log.New(accessLogFile, "", log.LstdFlags)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()
	r.Use(accessLogger(), lib.GinCtxToCtx())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run()
}

func init() {
	initLogs()
}
