package main

import (
	"github.com/JaegyuDev/Hydra/internal/json"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"net/http"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := json.ResponseJSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) testAuth(client clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessClaims, ok := clerk.SessionFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		user, err := client.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			panic(err)
		}

		w.Write([]byte("Welcome " + user.ID))
	}
}
