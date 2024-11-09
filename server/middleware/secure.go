package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func NewSecureHeaders() gin.HandlerFunc {
	return secure.New(secure.Config{
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		IENoOpen:             true,
		IsDevelopment:        false,
		ReferrerPolicy:       "strict-origin-when-cross-origin",
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func NewCorsHeaders(env string) gin.HandlerFunc {
	cfg := cors.DefaultConfig()

	/*
		cfg.AllowOrigins = viper.GetStringSlice("cors.origins")
		cfg.AllowOrigins = append(cfg.AllowOrigins, viper.GetString("endpoint"))
		cfg.AllowMethods = viper.GetStringSlice("cors.methods")
		cfg.AllowHeaders = viper.GetStringSlice("cors.allow_headers")
		cfg.ExposeHeaders = viper.GetStringSlice("cors.expose_headers")
		cfg.AllowCredentials = viper.GetBool("cors.allow_credentials")
		cfg.MaxAge = viper.GetDuration("cors.max_age") * time.Hour
	*/
	cfg.AllowBrowserExtensions = true
	cfg.AllowWildcard = true

	// for test env only
	cfg.AllowFiles = true

	return cors.New(cfg)
}
