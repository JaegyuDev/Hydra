package main

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	clerkClient, err := clerk.NewClient(app.config.clerk.clientToken)
	if err != nil {
		app.logger.Error("Could not start clerk client:", err)
		os.Exit(400)
	}

	injectActiveSession := clerk.WithSessionV2(clerkClient)

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.HandlerFunc("GET", "/status", app.status)
	mux.Handler("GET", "/testAuth", injectActiveSession(app.testAuth(clerkClient)))

	return app.recoverPanic(mux)
}
