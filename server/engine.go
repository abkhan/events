package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"events/server/handler"
	"events/server/middleware"
)

// NewEngine configure a basic gin engine.
func NewEngine(serviceName, env string) *gin.Engine {
	router := NewBasicEngine(serviceName, env)

	router.Use(
	//middleware.NewHeaderDetector(), // if we need to detect and do anything with the header
	)

	return router
}

// NewBasicEngine returns a new basic middleware configured gin engine that is fit for any microservices.
func NewBasicEngine(serviceName, env string) *gin.Engine {

	//TODO: make configurable
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(
		gin.Logger(),
		middleware.DetectClientIP,
		//middleware.NewCorsHeaders(env),
		middleware.NewSecureHeaders(),
	)

	router.GET("/ping", handler.NewPingHandler(serviceName))
	router.GET("/metrics", handler.NewMetricsHandler("put-correct-version-here"))

	return router
}
