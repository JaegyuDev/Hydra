package main

import (
	"encoding/json"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"net/http"
)

type Event struct {
	Data   *json.RawMessage `json:"data"`
	Object string           `json:"object"`
	Type   string           `json:"type"`
}

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

		var event Event
		err = json.Unmarshal(payload, &event)
		if err != nil {
			app.logger.Error("could not parse the event from the webhook:", err)
			http.Error(w, "Error parsing webhook", http.StatusUnprocessableEntity)
			return
		}

		switch event.Type {
		case "user.created":
			err := parseEventUserCreated(event)
			if err != nil {
				app.logger.Error("could not parse the data from the webhook:", err)
				http.Error(w, "Error parsing webhook", http.StatusInternalServerError)
				return
			}

		}
	}
}

// TODO: set up some emitter/register on event handler
func parseEventUserCreated(event Event) error {
	var user clerk.User
	err := json.Unmarshal(*event.Data, &user)
	if err != nil {
		return err
	}

	return nil
}
