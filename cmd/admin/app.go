package main

import (
	"context"
	"errors"
	"events/cmd/admin/handler/adminevents"
	"events/cmd/admin/middleware"
	"events/server"
	"events/server/handler"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type app struct {
	serviceName string
	server      *http.Server
}

var errServiceStopped = errors.New("service is already stopped")
var apiKey = "b4a8cfb4-71e3-47c6-9d19-afee216e09d6:app"

// newApp returns an http server
// - first sets up the app internals
// - add handlers to a web server, using internals and connections
// - returns the web server with handlers
func newApp(ctx context.Context) (*app, error) {

	// start AWS config

	// connect to AWS services

	a := &app{
		serviceName: appName + ".app",
	}

	h := a.getHandler(ctx, appEnv)

	a.server = &http.Server{
		Addr:              appAddr,
		Handler:           h,
		ReadHeaderTimeout: time.Second * 9,
	}

	return a, nil
}

func (a *app) getHandler(ctx context.Context, env string) http.Handler {
	engine := server.NewBasicEngine(a.serviceName, env)
	initiatorMiddleware, err := middleware.NewInitiatorMiddleware(apiKey) // TODO: where to get the apiKey from
	if err != nil {
		logrus.Fatalf("unable to init uuid middleware, error: %s", err)
	}

	apiAdmin := engine.Group("/admin")
	apiAdmin.GET("/ping", handler.NewPingHandler(a.serviceName)) // so it could take ping
	apiAdminSecure := apiAdmin.Group("", initiatorMiddleware)    // with initiator middleware
	{
		h := adminevents.NewHandler()
		h.BindEndpoint(apiAdminSecure)
	}

	return engine
}

func (a *app) Run() error {
	if a.server == nil {
		return errServiceStopped
	}

	logrus.Infof("Starting API server on %s", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("API server error: %s", err)
		}
	}()

	return nil
}

func (a *app) Shutdown(ctx context.Context) error {
	if a.server == nil {
		return errServiceStopped
	}

	toCtx, toCancel := context.WithTimeout(ctx, 60*time.Second)
	defer toCancel()

	if err := a.server.Shutdown(toCtx); err != nil {
		logrus.Errorf("API Server shutdown failed: %s", err)
	}

	a.server = nil

	return nil
}
