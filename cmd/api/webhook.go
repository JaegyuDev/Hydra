package main

import (
	"encoding/json"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"net/http"
)

type ClerkEvent struct {
	Data   *json.RawMessage `json:"data"`
	Object string           `json:"object"`
	Type   ClerkEventType   `json:"type"`
}

// clerkWebhook should run on /webhooks/clerk. This is set in Clerk's dashboard and
// if this is changed it WILL break things
func (app *application) clerkWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := make([]byte, 2048)
		_, err := r.Body.Read(payload)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusUnprocessableEntity)
			return
		}

		err = app.clerk.svixWebHook.Verify(payload, r.Header)
		if err != nil {
			http.Error(w, "Invalid webhook response", http.StatusBadRequest)
			return
		}

		var event ClerkEvent
		err = json.Unmarshal(payload, &event)
		if err != nil {
			app.logger.Error("could not parse the event from the webhook:", err)
			http.Error(w, "Error parsing webhook", http.StatusUnprocessableEntity)
			return
		}

		// These two types will have different event.Data structures
		switch event.Type {
		case UserCreated:
			var user clerk.User
			err := json.Unmarshal(*event.Data, &user)
			if err != nil {
				app.logger.Error("could not parse the data from the webhook:", err)
				http.Error(w, "Error parsing webhook", http.StatusInternalServerError)
				return
			}

			app.clerk.eventEmitter.Trigger(event.Type, user)
		}
	}
}
